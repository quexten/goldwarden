package actions

import (
	"context"
	"fmt"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
)

func handleLogin(msg messages.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response messages.IPCMessage, err error) {
	if !cfg.HasPin() && !cfg.ConfigFile.RuntimeConfig.DisablePinRequirement {
		response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: "No pin set. Set a pin first!",
		})
		if err != nil {
			return messages.IPCMessage{}, err
		}

		return
	}

	req := messages.ParsePayload(msg).(messages.DoLoginRequest)

	ctx := context.Background()
	var token bitwarden.LoginResponseToken
	var masterKey crypto.MasterKey
	var masterpasswordHash string

	if req.Passwordless {
		token, masterKey, masterpasswordHash, err = bitwarden.LoginWithDevice(ctx, req.Email, cfg, vault)
	} else {
		token, masterKey, masterpasswordHash, err = bitwarden.LoginWithMasterpassword(ctx, req.Email, cfg, vault)
	}
	if err != nil {
		var payload = messages.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not login: %s", err.Error()),
		}
		response, err = messages.IPCMessageFromPayload(payload)
		if err != nil {
			return messages.IPCMessage{}, err
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
		var payload = messages.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not sync vault: %s", err.Error()),
		}
		response, err = messages.IPCMessageFromPayload(payload)
		if err != nil {
			return messages.IPCMessage{}, err
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
		var payload = messages.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not sync vault: %s", err.Error()),
		}
		response, err = messages.IPCMessageFromPayload(payload)
		if err != nil {
			return messages.IPCMessage{}, err
		}
		return
	}

	cfg.SetUserSymmetricKey(vault.Keyring.AccountKey.Bytes())
	cfg.SetMasterPasswordHash([]byte(masterpasswordHash))
	cfg.SetMasterKey([]byte(masterKey.GetBytes()))
	var protectedUserSymetricKey crypto.SymmetricEncryptionKey
	if vault.Keyring.IsMemguard {
		protectedUserSymetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(vault.Keyring.AccountKey.Bytes())
	} else {
		protectedUserSymetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(vault.Keyring.AccountKey.Bytes())
	}
	if err != nil {
		var payload = messages.ActionResponse{
			Success: false,
			Message: fmt.Sprintf("Could not sync vault: %s", err.Error()),
		}
		response, err = messages.IPCMessageFromPayload(payload)
		if err != nil {
			return messages.IPCMessage{}, err
		}
		return
	}
	err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, cfg, &protectedUserSymetricKey, false)

	response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
		Success: true,
	})
	if err != nil {
		panic(err)
	}

	return
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.DoLoginRequest{}), ensureIsNotLocked(handleLogin))
}
