package actions

import (
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
)

func handleSetApiURL(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	apiURL := messages.ParsePayload(request).(messages.SetApiURLRequest).Value
	cfg.ConfigFile.ApiUrl = apiURL
	err = cfg.WriteConfig()
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return messages.IPCMessageFromPayload(messages.ActionResponse{
		Success: true,
	})
}

func handleSetIdentity(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	identity := messages.ParsePayload(request).(messages.SetIdentityURLRequest).Value
	cfg.ConfigFile.IdentityUrl = identity
	err = cfg.WriteConfig()
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return messages.IPCMessageFromPayload(messages.ActionResponse{
		Success: true,
	})
}

func handleSetNotifications(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	notifications := messages.ParsePayload(request).(messages.SetNotificationsURLRequest).Value
	cfg.ConfigFile.NotificationsUrl = notifications
	err = cfg.WriteConfig()
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return messages.IPCMessageFromPayload(messages.ActionResponse{
		Success: true,
	})
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetIdentityURLRequest{}), handleSetIdentity)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetApiURLRequest{}), handleSetApiURL)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetNotificationsURLRequest{}), handleSetNotifications)
}
