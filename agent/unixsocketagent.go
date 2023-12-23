package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/quexten/goldwarden/agent/actions"
	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/processsecurity"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/ssh"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/quexten/goldwarden/logging"
)

const (
	FullSyncInterval     = 60 * time.Minute
	TokenRefreshInterval = 30 * time.Minute
)

var log = logging.GetLogger("Goldwarden", "Agent")

func writeError(c net.Conn, errMsg error) error {
	payload := messages.ActionResponse{
		Success: false,
		Message: errMsg.Error(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.Write(payloadBytes)
	if err != nil {
		return err
	}
	return nil
}

func serveAgentSession(c net.Conn, ctx context.Context, vault *vault.Vault, cfg *config.Config) {
	for {
		buf := make([]byte, 1024*1024)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]

		var msg messages.IPCMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			writeError(c, err)
			continue
		}

		responseBytes := []byte{}
		if action, actionFound := actions.AgentActionsRegistry.Get(msg.Type); actionFound {
			callingContext := sockets.GetCallingContext(c)
			payload, err := action(msg, cfg, vault, &callingContext)
			if err != nil {
				writeError(c, err)
				continue
			}
			responseBytes, err = json.Marshal(payload)
			if err != nil {
				writeError(c, err)
				continue
			}
		} else {
			payload := messages.ActionResponse{
				Success: false,
				Message: "Action not found",
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				writeError(c, err)
				continue
			}
			responseBytes = payloadBytes
		}

		_, err = c.Write(responseBytes)
		if err != nil {
			log.Error("Failed writing to socket " + err.Error())
		}
	}
}

type AgentState struct {
}

func StartUnixAgent(path string, runtimeConfig config.RuntimeConfig) error {
	ctx := context.Background()

	// check if exists
	keyring := crypto.NewKeyring(nil)
	var vault = vault.NewVault(&keyring)
	cfg, err := config.ReadConfig(runtimeConfig)
	if err != nil {
		log.Warn("Could not read config: %s", err.Error())
		cfg = config.DefaultConfig()
		cfg.ConfigFile.RuntimeConfig = runtimeConfig
		cfg.WriteConfig()
	}
	cfg.ConfigFile.RuntimeConfig = runtimeConfig
	if cfg.ConfigFile.RuntimeConfig.ApiURI != "" {
		cfg.ConfigFile.ApiUrl = cfg.ConfigFile.RuntimeConfig.ApiURI
	}
	if cfg.ConfigFile.RuntimeConfig.IdentityURI != "" {
		cfg.ConfigFile.IdentityUrl = cfg.ConfigFile.RuntimeConfig.IdentityURI
	}
	if cfg.ConfigFile.RuntimeConfig.NotificationsURI != "" {
		cfg.ConfigFile.NotificationsUrl = cfg.ConfigFile.RuntimeConfig.NotificationsURI
	}
	if cfg.ConfigFile.RuntimeConfig.DeviceUUID != "" {
		cfg.ConfigFile.DeviceUUID = cfg.ConfigFile.RuntimeConfig.DeviceUUID
	}

	if !cfg.IsLocked() {
		log.Warn("Config is not locked. SET A PIN!!")
		token, err := cfg.GetToken()
		if err == nil {
			if token.AccessToken != "" {
				// attempt to sync every minute until successful
				for {
					bitwarden.RefreshToken(ctx, &cfg)
					userSymmetricKey, err := cfg.GetUserSymmetricKey()
					if err != nil {
						fmt.Println(err)
						time.Sleep(60 * time.Second)
						continue
					}
					protectedUserSymetricKey, err := crypto.SymmetricEncryptionKeyFromBytes(userSymmetricKey)

					err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, &protectedUserSymetricKey, true)
					if err != nil {
						fmt.Println(err)
						time.Sleep(60 * time.Second)
						continue
					} else {
						break
					}
				}
			}
		}
	}

	processsecurity.DisableDumpable()
	err = processsecurity.MonitorLocks(func() {
		cfg.Lock()
		vault.Clear()
		vault.Keyring.Lock()
	})
	if err != nil {
		log.Warn("Could not monitor screensaver: %s", err.Error())
	}

	if !runtimeConfig.WebsocketDisabled {
		go bitwarden.RunWebsocketDaemon(ctx, vault, &cfg)
	}

	if !runtimeConfig.DisableSSHAgent {
		vaultAgent := ssh.NewVaultAgent(vault, &cfg)
		vaultAgent.SetUnlockRequestAction(func() bool {
			err := cfg.TryUnlock(vault)
			if err == nil {
				token, err := cfg.GetToken()
				if err == nil {
					if token.AccessToken != "" {
						bitwarden.RefreshToken(ctx, &cfg)
						userSymmetricKey, err := cfg.GetUserSymmetricKey()
						if err != nil {
							fmt.Println(err)
						}
						protectedUserSymetricKey, err := crypto.SymmetricEncryptionKeyFromBytes(userSymmetricKey)

						err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, &protectedUserSymetricKey, true)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
				return true
			}
			return false
		})
		go vaultAgent.Serve()
	}

	go func() {
		for {
			time.Sleep(TokenRefreshInterval)
			if !cfg.IsLocked() {
				bitwarden.RefreshToken(ctx, &cfg)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(FullSyncInterval)
			if !cfg.IsLocked() {
				token, err := cfg.GetToken()
				if err != nil {
					log.Warn("Could not get token: %s", err.Error())
					continue
				}

				bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token), vault, &cfg, nil, false)
			}
		}
	}()

	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	l, err := net.Listen("unix", path)
	if err != nil {
		println("listen error", err.Error())
		return err
	}
	log.Info("Agent listening on %s...", path)

	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				println("accept error", err.Error())
			}

			go serveAgentSession(fd, ctx, vault, &cfg)
		}
	}()

	mainloop()
	return nil
}
