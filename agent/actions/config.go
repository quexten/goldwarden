package actions

import (
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func handleSetApiURL(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx sockets.CallingContext) (response interface{}, err error) {
	apiURL := request.ParsedPayload().(ipc.SetApiURLRequest).Value
	cfg.ConfigFile.ApiUrl = apiURL
	err = cfg.WriteConfig()
	if err != nil {
		return ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
}

func handleSetIdentity(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx sockets.CallingContext) (response interface{}, err error) {
	identity := request.ParsedPayload().(ipc.SetIdentityURLRequest).Value
	cfg.ConfigFile.IdentityUrl = identity
	err = cfg.WriteConfig()
	if err != nil {
		return ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
}

func handleSetNotifications(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx sockets.CallingContext) (response interface{}, err error) {
	notifications := request.ParsedPayload().(ipc.SetNotificationsURLRequest).Value
	cfg.ConfigFile.NotificationsUrl = notifications
	err = cfg.WriteConfig()
	if err != nil {
		return ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ipc.IPCMessageFromPayload(ipc.ActionResponse{
		Success: true,
	})
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageTypeSetIdentityURLRequest, handleSetIdentity)
	AgentActionsRegistry.Register(ipc.IPCMessageTypeSetAPIUrlRequest, handleSetApiURL)
	AgentActionsRegistry.Register(ipc.IPCMessageTypeSetNotificationsURLRequest, handleSetNotifications)
}
