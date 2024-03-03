package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/notify"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
)

func handleLogin(msg messages.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response messages.IPCMessage, err error) {
	if !cfg.HasPin() {
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

	if secret, err := cfg.GetClientSecret(); err == nil && secret != "" {
		actionsLog.Info("Logging in with client secret")
		token, masterKey, masterpasswordHash, err = bitwarden.LoginWithApiKey(ctx, req.Email, cfg, vault)
	} else if req.Passwordless {
		actionsLog.Info("Logging in with passwordless")
		token, masterKey, masterpasswordHash, err = bitwarden.LoginWithDevice(ctx, req.Email, cfg, vault)
	} else {
		actionsLog.Info("Logging in with master password")
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

	_ = cfg.SetToken(config.LoginToken{
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
		defer func() {
			notify.Notify("Goldwarden", "Could not decrypt. Wrong password?", "", 10*time.Second, func() {})
			_ = cfg.SetToken(config.LoginToken{})
			vault.Clear()
		}()

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

	err = cfg.SetUserSymmetricKey(vault.Keyring.GetAccountKey().Bytes())
	err = cfg.SetMasterPasswordHash([]byte(masterpasswordHash))
	err = cfg.SetMasterKey([]byte(masterKey.GetBytes()))
	var protectedUserSymetricKey crypto.SymmetricEncryptionKey
	if vault.Keyring.IsMemguard {
		protectedUserSymetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(vault.Keyring.GetAccountKey().Bytes())
	} else {
		protectedUserSymetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(vault.Keyring.GetAccountKey().Bytes())
	}
	if err != nil {
		defer func() {
			notify.Notify("Goldwarden", "Could not decrypt. Wrong password?", "", 10*time.Second, func() {})
			_ = cfg.SetToken(config.LoginToken{})
			vault.Clear()
		}()

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
