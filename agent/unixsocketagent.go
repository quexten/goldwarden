package agent

import (
	"context"
	"crypto/subtle"
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
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/quexten/goldwarden/logging"
)

const (
	FullSyncInterval     = 60 * time.Minute
	TokenRefreshInterval = 10 * time.Minute
)

var log = logging.GetLogger("Goldwarden", "Agent")

func writeError(c net.Conn, errMsg error) {
	payload := messages.ActionResponse{
		Success: false,
		Message: errMsg.Error(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Warn("Could not json marshall: %s", err.Error())
		return
	}
	_, err = c.Write(payloadBytes)
	if err != nil {
		log.Warn("Could not write payload: %s", err.Error())
	}
}

func serveAgentSession(c net.Conn, vault *vault.Vault, cfg *config.Config) {
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

		// todo refactor to other file
		if msg.Type == messages.MessageTypeForEmptyPayload(messages.SessionAuthRequest{}) {
			if cfg.ConfigFile.RuntimeConfig.DaemonAuthToken == "" {
				return
			}

			req := messages.ParsePayload(msg).(messages.SessionAuthRequest)
			verified := subtle.ConstantTimeCompare([]byte(cfg.ConfigFile.RuntimeConfig.DaemonAuthToken), []byte(req.Token)) == 1

			payload := messages.SessionAuthResponse{
				Verified: verified,
			}
			log.Info("Verified: %t", verified)
			callingContext := sockets.GetCallingContext(c)
			if verified {
				systemauth.CreatePinSession(callingContext, 365*24*time.Hour) // permanent session
			}

			responsePayload, err := messages.IPCMessageFromPayload(payload)
			if err != nil {
				writeError(c, err)
				continue
			}
			payloadBytes, err := json.Marshal(responsePayload)
			if err != nil {
				writeError(c, err)
				continue
			}

			_, err = c.Write(payloadBytes)
			if err != nil {
				log.Error("Failed writing to socket " + err.Error())
			}
			continue
		}

		// todo refactor to other file
		if msg.Type == messages.MessageTypeForEmptyPayload(messages.PinentryRegistrationRequest{}) {
			// todo lockdown this method better
			if cfg.ConfigFile.RuntimeConfig.DaemonAuthToken == "" {
				return
			}

			log.Info("Received pinentry registration request")

			getPasswordChan := make(chan struct {
				title       string
				description string
			})
			getPasswordReturnChan := make(chan struct {
				password string
				err      error
			})
			getApprovalChan := make(chan struct {
				title       string
				description string
			})
			getApprovalReturnChan := make(chan struct {
				approved bool
				err      error
			})

			pe := pinentry.Pinentry{
				GetPassword: func(title string, description string) (string, error) {
					getPasswordChan <- struct {
						title       string
						description string
					}{title, description}
					returnValue := <-getPasswordReturnChan
					return returnValue.password, returnValue.err
				},
				GetApproval: func(title string, description string) (bool, error) {
					getApprovalChan <- struct {
						title       string
						description string
					}{title, description}
					returnValue := <-getApprovalReturnChan
					return returnValue.approved, returnValue.err
				},
			}

			pinentrySetError := pinentry.SetExternalPinentry(pe)
			payload := messages.PinentryRegistrationResponse{
				Success: pinentrySetError == nil,
			}
			log.Info("Pinentry registration success: %t", payload.Success)

			responsePayload, err := messages.IPCMessageFromPayload(payload)
			if err != nil {
				writeError(c, err)
				continue
			}
			payloadBytes, err := json.Marshal(responsePayload)
			if err != nil {
				writeError(c, err)
				continue
			}

			_, err = c.Write(payloadBytes)
			if err != nil {
				log.Error("Failed writing to socket " + err.Error())
			}
			_, err = c.Write([]byte("\n"))
			if err != nil {
				log.Error("Failed writing to socket " + err.Error())
			}
			time.Sleep(50 * time.Millisecond) //todo fix properly

			if pinentrySetError != nil {
				return
			}

			for {
				fmt.Println("Waiting for pinentry request")
				select {
				case getPasswordRequest := <-getPasswordChan:
					log.Info("Received getPassword request")
					payload := messages.PinentryPinRequest{
						Message: getPasswordRequest.description,
					}
					payloadPayload, err := messages.IPCMessageFromPayload(payload)
					if err != nil {
						writeError(c, err)
						continue
					}

					payloadBytes, err := json.Marshal(payloadPayload)
					if err != nil {
						writeError(c, err)
						continue
					}

					_, err = c.Write(payloadBytes)
					if err != nil {
						log.Error("Failed writing to socket " + err.Error())
					}

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

					if msg.Type == messages.MessageTypeForEmptyPayload(messages.PinentryPinResponse{}) {
						getPasswordResponse := messages.ParsePayload(msg).(messages.PinentryPinResponse)
						getPasswordReturnChan <- struct {
							password string
							err      error
						}{getPasswordResponse.Pin, nil}
					}
				case getApprovalRequest := <-getApprovalChan:
					log.Info("Received getApproval request")
					payload := messages.PinentryApprovalRequest{
						Message: getApprovalRequest.description,
					}
					payloadPayload, err := messages.IPCMessageFromPayload(payload)
					if err != nil {
						writeError(c, err)
						continue
					}

					payloadBytes, err := json.Marshal(payloadPayload)
					if err != nil {
						writeError(c, err)
						continue
					}

					_, err = c.Write(payloadBytes)
					if err != nil {
						log.Error("Failed writing to socket " + err.Error())
					}

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

					if msg.Type == messages.MessageTypeForEmptyPayload(messages.PinentryApprovalResponse{}) {
						getApprovalResponse := messages.ParsePayload(msg).(messages.PinentryApprovalResponse)
						getApprovalReturnChan <- struct {
							approved bool
							err      error
						}{getApprovalResponse.Approved, nil}
					}
				}
			}
		}

		var responseBytes []byte
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
		err = cfg.WriteConfig()
		if err != nil {
			log.Warn("Could not write config: %s", err.Error())
		}
	}
	cfg.ConfigFile.RuntimeConfig = runtimeConfig
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
					var protectedUserSymmetricKey crypto.SymmetricEncryptionKey
					if vault.Keyring.IsMemguard {
						protectedUserSymmetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(userSymmetricKey)
					} else {
						protectedUserSymmetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(userSymmetricKey)
					}
					if err != nil {
						log.Error("could not get encryption key from bytes: %s", err.Error())
					}

					err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, &protectedUserSymmetricKey, true)
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

	err = processsecurity.DisableDumpable()
	if err != nil {
		log.Warn("Could not disable dumpable: %s", err.Error())
	}

	go func() {
		err = processsecurity.MonitorLocks(func() {
			cfg.Lock()
			vault.Clear()
			vault.Keyring.Lock()
			systemauth.WipeSessions()
		})
		if err != nil {
			log.Warn("Could not monitor screensaver: %s", err.Error())
		}
	}()
	go func() {
		err = processsecurity.MonitorIdle(func() {
			cfg.Lock()
			vault.Clear()
			vault.Keyring.Lock()
			systemauth.WipeSessions()
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
						var protectedUserSymmetricKey crypto.SymmetricEncryptionKey
						if vault.Keyring.IsMemguard {
							protectedUserSymmetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(userSymmetricKey)
						} else {
							protectedUserSymmetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(userSymmetricKey)
						}
						if err != nil {
							log.Error("could not get encryption key from bytes: %s", err.Error())
						}

						err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, &protectedUserSymmetricKey, true)
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

				err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, &cfg, nil, false)
				if err != nil {
					log.Warn("Could not do full sync: %s", err.Error())
					continue
				}
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
		fmt.Println("listen error", err.Error())
		return err
	}
	defer l.Close()
	log.Info("Agent listening on %s...", path)

	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				fmt.Println("accept error", err.Error())
			}

			go serveAgentSession(fd, vault, &cfg)
		}
	}()

	mainloop()
	return nil
}
