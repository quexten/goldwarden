package actions

import (
	"context"
	"fmt"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func handleLogin(msg ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext sockets.CallingContext) (response interface{}, err error) {
	if !cfg.HasPin() {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "No pin set. Set a pin first!",
		})
		if err != nil {
			return nil, err
		}

		return
	}

	req := msg.ParsedPayload().(ipc.DoLoginRequest)

	ctx := context.Background()
	token, masterKey, masterpasswordHash, err := bitwarden.LoginWithMasterpassword(ctx, req.Email, cfg, vault)
	if err != nil {
		var payload = ipc.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not login: %s", err.Error()),
		}
		response, err = ipc.IPCMessageFromPayload(payload)
		if err != nil {
			return nil, err
		}
		return
	}

	cfg.SetToken(config.LoginToken{
		AccessToken:  token.AccessToken,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Key:          token.Key,
	})

	profile, err := bitwarden.Sync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), cfg)
	if err != nil {
		var payload = ipc.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not sync vault: %s", err.Error()),
		}
		response, err = ipc.IPCMessageFromPayload(payload)
		if err != nil {
			return nil, err
		}
		return
	}

	var orgKeys map[string]string = make(map[string]string)
	for _, org := range profile.Profile.Organizations {
		orgId := org.Id.String()
		orgKeys[orgId] = org.Key
	}

	err = crypto.InitKeyringFromMasterKey(vault.Keyring, profile.Profile.Key, profile.Profile.PrivateKey, orgKeys, masterKey)
	if err != nil {
		var payload = ipc.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not sync vault: %s", err.Error()),
		}
		response, err = ipc.IPCMessageFromPayload(payload)
		if err != nil {
			return nil, err
		}
		return
	}

	cfg.SetUserSymmetricKey(vault.Keyring.AccountKey.Bytes())
	cfg.SetMasterPasswordHash([]byte(masterpasswordHash))
	cfg.SetMasterKey([]byte(masterKey.GetBytes()))
	protectedUserSymetricKey, err := crypto.SymmetricEncryptionKeyFromBytes(vault.Keyring.AccountKey.Bytes())
	if err != nil {
		var payload = ipc.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not sync vault: %s", err.Error()),
		}
		response, err = ipc.IPCMessageFromPayload(payload)
		if err != nil {
			return nil, err
		}
		return
	}
	err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, cfg, &protectedUserSymetricKey, false)

	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
	if err != nil {
		panic(err)
	}

	return
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageTypeDoLoginRequest, ensureIsNotLocked(handleLogin))
}
