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
	"github.com/quexten/goldwarden/agent/notify"
	"github.com/quexten/goldwarden/agent/processsecurity"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/ssh"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/quexten/goldwarden/logging"
)

const (
	FullSyncInterval     = 60 * time.Minute
	TokenRefreshInterval = 10 * time.Minute
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

	var keyring crypto.Keyring
	if runtimeConfig.UseMemguard {
		keyring = crypto.NewMemguardKeyring(nil)
	} else {
		keyring = crypto.NewMemoryKeyring(nil)
	}

	var vault = vault.NewVault(&keyring)
	cfg, err := config.ReadConfig(runtimeConfig)
	if err != nil {
		log.Warn("Could not read config: %s", err.Error())
		cfg = config.DefaultConfig(runtimeConfig.UseMemguard)
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
					token, err := cfg.GetToken()
					if err != nil {
						log.Error("Could not get token: %s", err.Error())
					}

					userSymmetricKey, err := cfg.GetUserSymmetricKey()
					if err != nil {
						fmt.Println(err)
						time.Sleep(60 * time.Second)
						continue
					}
					var protectedUserSymetricKey crypto.SymmetricEncryptionKey
					if vault.Keyring.IsMemguard {
						protectedUserSymetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(userSymmetricKey)
					} else {
						protectedUserSymetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(userSymmetricKey)
					}

					err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, &protectedUserSymetricKey, true)
					if err != nil {
						log.Error("Could not sync: %s", err.Error())
						notify.Notify("Goldwarden", "Could not perform initial sync", "", 0, func() {})
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
	go func() {
		err = processsecurity.MonitorLocks(func() {
			cfg.Lock()
			vault.Clear()
			vault.Keyring.Lock()
		})
		if err != nil {
			log.Warn("Could not monitor screensaver: %s", err.Error())
		}
	}()
	go func() {
		err = processsecurity.MonitorIdle(func() {
			log.Warn("Idling detected but no action is implemented")
		})
		if err != nil {
			log.Warn("Could not monitor idle: %s", err.Error())
		}
	}()
	go func() {
		err = notify.ListenForNotifications()
		if err != nil {
			log.Warn("Could not listen for notifications: %s", err.Error())
		}
	}()

	go func() {
		if !runtimeConfig.WebsocketDisabled {
			for {
				// polling, switch this to signal based later
				if !cfg.IsLocked() && cfg.IsLoggedIn() {
					bitwarden.RunWebsocketDaemon(ctx, vault, &cfg)
					time.Sleep(60 * time.Second)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	if !runtimeConfig.DisableSSHAgent {
		vaultAgent := ssh.NewVaultAgent(vault, &cfg, &runtimeConfig)
		vaultAgent.SetUnlockRequestAction(func() bool {
			err := cfg.TryUnlock(vault)
			if err == nil {
				token, err := cfg.GetToken()
				if err == nil {
					if token.AccessToken != "" {
						gotToken := bitwarden.RefreshToken(ctx, &cfg)
						if !gotToken {
							log.Warn("Could not get token")
							return false
						} else {
							token, err = cfg.GetToken()
							if err != nil {
								log.Warn("Could not get token: %s", err.Error())
								return false
							}
						}
						userSymmetricKey, err := cfg.GetUserSymmetricKey()
						if err != nil {
							log.Error("Could not get user symmetric key: %s", err.Error())
						}
						var protectedUserSymetricKey crypto.SymmetricEncryptionKey
						if vault.Keyring.IsMemguard {
							protectedUserSymetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(userSymmetricKey)
						} else {
							protectedUserSymetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(userSymmetricKey)
						}

						err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, &protectedUserSymetricKey, true)
						if err != nil {
							log.Error("Could not sync: %s", err.Error())
							notify.Notify("Goldwarden", "Could not perform initial sync on ssh unlock", "", 0, func() {})
						}
					} else {
						log.Warn("Access token is empty")
					}
				} else {
					log.Error("Could not get token: %s", err.Error())
				}
				return true
			} else {
				log.Warn("Could not unlock: %s", err.Error())
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
				bitwarden.RefreshToken(ctx, &cfg)
				token, err := cfg.GetToken()
				if err != nil {
					log.Warn("Could not get token: %s", err.Error())
					continue
				}

				bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, nil, false)
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
