package actions

import (
	"context"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/quexten/goldwarden/logging"
)

var actionsLog = logging.GetLogger("Goldwarden", "Actions")
var AgentActionsRegistry = newActionsRegistry()

type Action func(messages.IPCMessage, *config.Config, *vault.Vault, *sockets.CallingContext) (messages.IPCMessage, error)
type ActionsRegistry struct {
	actions map[messages.IPCMessageType]Action
}

func newActionsRegistry() *ActionsRegistry {
	return &ActionsRegistry{
		actions: make(map[messages.IPCMessageType]Action),
	}
}

func (registry *ActionsRegistry) Register(messageType messages.IPCMessageType, action Action) {
	registry.actions[messageType] = action
}

func (registry *ActionsRegistry) Get(messageType messages.IPCMessageType) (Action, bool) {
	action, ok := registry.actions[messageType]
	return action, ok
}

func ensureIsLoggedIn(action Action) Action {
	return func(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (messages.IPCMessage, error) {
		if hash, err := cfg.GetMasterPasswordHash(); err != nil || len(hash) == 0 {
			actionsLog.Error("EnsureIsLoggedIn - %s", err.Error())
			return messages.IPCMessageFromPayload(messages.ActionResponse{
				Success: false,
				Message: "Not logged in",
			})
		}

		return action(request, cfg, vault, ctx)
	}
}

func sync(ctx context.Context, vault *vault.Vault, cfg *config.Config) bool {
	token, err := cfg.GetToken()
	if err == nil {
		if token.AccessToken != "" {
			refreshed := bitwarden.RefreshToken(ctx, cfg)
			if !refreshed {
				return false
			}

			userSymmetricKey, err := cfg.GetUserSymmetricKey()
			if err != nil {
				return false
			}

			var protectedUserSymmetricKey crypto.SymmetricEncryptionKey
			if vault.Keyring.IsMemguard {
				protectedUserSymmetricKey, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(userSymmetricKey)
			} else {
				protectedUserSymmetricKey, err = crypto.MemorySymmetricEncryptionKeyFromBytes(userSymmetricKey)
			}
			if err != nil {
				return false
			}

			err = bitwarden.DoFullSync(context.WithValue(ctx, bitwarden.AuthToken{}, token.AccessToken), vault, cfg, &protectedUserSymmetricKey, true)
			if err != nil {
				return false
			}
		}
	}
	return true
}

func ensureIsNotLocked(action Action) Action {
	return func(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (messages.IPCMessage, error) {
		if cfg.IsLocked() {
			err := cfg.TryUnlock(vault)
			ctx1 := context.Background()
			success := sync(ctx1, vault, cfg)
			if err != nil || !success {
				if err != nil {
					return messages.IPCMessageFromPayload(messages.ActionResponse{
						Success: false,
						Message: err.Error(),
					})
				} else {
					return messages.IPCMessageFromPayload(messages.ActionResponse{
						Success: false,
						Message: "Could not sync vault",
					})
				}
			}

			systemauth.CreatePinSession(*ctx, systemauth.SSHTTL)
		}

		return action(request, cfg, vault, ctx)
	}
}

func ensureBiometricsAuthorized(approvalType systemauth.SessionType, action Action) Action {
	return func(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (messages.IPCMessage, error) {
		if permission, err := systemauth.GetPermission(approvalType, *ctx, cfg); err != nil || !permission {
			return messages.IPCMessageFromPayload(messages.ActionResponse{
				Success: false,
				Message: "Polkit authorization failed required",
			})
		}

		return action(request, cfg, vault, ctx)
	}
}

func ensureEverything(approvalType systemauth.SessionType, action Action) Action {
	return ensureIsNotLocked(ensureIsLoggedIn(ensureBiometricsAuthorized(approvalType, action)))
}
