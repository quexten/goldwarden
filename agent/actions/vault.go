package actions

import (
	"context"
	"fmt"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func handleUnlockVault(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext sockets.CallingContext) (response interface{}, err error) {
	if !cfg.HasPin() {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "No pin set",
		})
		if err != nil {
			panic(err)
		}

		return
	}

	if !cfg.IsLocked() {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: true,
			Message: "Unlocked",
		})
		if err != nil {
			panic(err)
		}

		return
	}

	err = cfg.TryUnlock(vault)
	if err != nil {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "wrong pin: " + err.Error(),
		})
		if err != nil {
			panic(err)
		}

		return
	}

	if cfg.IsLoggedIn() {
		token, err := cfg.GetToken()
		if err == nil {
			if token.AccessToken != "" {
				ctx := context.Background()
				bitwarden.RefreshToken(ctx, cfg)
				token, err := cfg.GetToken()
				err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, cfg, nil, true)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
	if err != nil {
		panic(err)
	}

	return
}

func handleLockVault(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext sockets.CallingContext) (response interface{}, err error) {
	if !cfg.HasPin() {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "No pin set",
		})
		if err != nil {
			panic(err)
		}

		return
	}

	if cfg.IsLocked() {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: true,
			Message: "Locked",
		})
		if err != nil {
			panic(err)
		}

		return
	}

	cfg.Lock()
	vault.Clear()
	vault.Keyring.Lock()

	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
	if err != nil {
		panic(err)
	}

	return
}

func handleWipeVault(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext sockets.CallingContext) (response interface{}, err error) {
	cfg.Purge()
	cfg.WriteConfig()
	vault.Clear()

	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
	if err != nil {
		panic(err)
	}

	return
}

func handleUpdateVaultPin(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext sockets.CallingContext) (response interface{}, err error) {
	pin, err := systemauth.GetPassword("Pin Change", "Enter your desired pin")
	if err != nil {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
		if err != nil {
			return nil, err
		} else {
			return response, nil
		}
	}
	cfg.UpdatePin(pin, true)

	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})

	return
}

func handlePinStatus(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext sockets.CallingContext) (response interface{}, err error) {
	var pinStatus string
	if cfg.HasPin() {
		pinStatus = "enabled"
	} else {
		pinStatus = "disabled"
	}

	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
		Message: pinStatus,
	})

	return
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageTypeUnlockVaultRequest, handleUnlockVault)
	AgentActionsRegistry.Register(ipc.IPCMessageTypeLockVaultRequest, handleLockVault)
	AgentActionsRegistry.Register(ipc.IPCMessageTypeWipeVaultRequest, handleWipeVault)
	AgentActionsRegistry.Register(ipc.IPCMessageTypeUpdateVaultPINRequest, ensureBiometricsAuthorized(systemauth.ChangePin, handleUpdateVaultPin))
	AgentActionsRegistry.Register(ipc.IPCMessageTypeGetVaultPINStatusRequest, handlePinStatus)
}
