package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os/user"
	"time"

	"github.com/quexten/goldwarden/agent/actions"
	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func writeErrorToLog(err error) {
	log.Error(err.Error())
}

func serveVirtualAgent(recv chan []byte, send chan []byte, ctx context.Context, vault *vault.Vault, cfg *config.Config) {
	for {
		data := <-recv

		var msg ipc.IPCMessage
		err := json.Unmarshal(data, &msg)
		if err != nil {
			writeErrorToLog(err)
			continue
		}

		responseBytes := []byte{}
		if action, actionFound := actions.AgentActionsRegistry.Get(msg.Type); actionFound {
			user, _ := user.Current()
			process := "goldwarden"
			parent := "SINGLE_PROC_MODE"
			grandparent := "SINGLE_PROC_MODE"
			callingContext := sockets.CallingContext{
				UserName:               user.Name,
				ProcessName:            process,
				ParentProcessName:      parent,
				GrandParentProcessName: grandparent,
			}
			payload, err := action(msg, cfg, vault, callingContext)
			if err != nil {
				writeErrorToLog(err)
				continue
			}
			responseBytes, err = json.Marshal(payload)
			if err != nil {
				writeErrorToLog(err)
				continue
			}
		} else {
			payload := ipc.ActionResponse{
				Success: false,
				Message: "Action not found",
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				writeErrorToLog(err)
				continue
			}
			responseBytes = payloadBytes
		}

		send <- responseBytes
	}
}

func StartVirtualAgent(runtimeConfig config.RuntimeConfig) (chan []byte, chan []byte) {
	ctx := context.Background()

	// check if exists
	keyring := crypto.NewKeyring(nil)
	var vault = vault.NewVault(&keyring)
	cfg, err := config.ReadConfig(runtimeConfig)
	if err != nil {
		var cfg = config.DefaultConfig()
		cfg.WriteConfig()
	}
	cfg.ConfigFile.RuntimeConfig = runtimeConfig
	if cfg.ConfigFile.RuntimeConfig.ApiURI != "" {
		cfg.ConfigFile.ApiUrl = cfg.ConfigFile.RuntimeConfig.ApiURI
	}
	if cfg.ConfigFile.RuntimeConfig.IdentityURI != "" {
		cfg.ConfigFile.IdentityUrl = cfg.ConfigFile.RuntimeConfig.IdentityURI
	}
	if cfg.ConfigFile.RuntimeConfig.DeviceUUID != "" {
		cfg.ConfigFile.DeviceUUID = cfg.ConfigFile.RuntimeConfig.DeviceUUID
	}

	if !cfg.IsLocked() && !cfg.ConfigFile.RuntimeConfig.DoNotPersistConfig {
		log.Warn("Config is not locked. SET A PIN!!")
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
	}
	disableDumpable()
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

	recv := make(chan []byte)
	send := make(chan []byte)

	go func() {
		go serveVirtualAgent(recv, send, ctx, vault, &cfg)
	}()
	return recv, send
}
