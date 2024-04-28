package actions

import (
	"encoding/json"
	"io"
	"net/http"

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

func handleSetVaultURL(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	vaultURL := messages.ParsePayload(request).(messages.SetVaultURLRequest).Value
	cfg.ConfigFile.VaultUrl = vaultURL
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

type ConfigResponse struct {
	Version     string `json:"version"`
	GitHash     string `json:"gitHash"`
	Environment struct {
		Vault         string `json:"vault"`
		Api           string `json:"api"`
		Identity      string `json:"identity"`
		Notifications string `json:"notifications"`
	}
}

func handleSetURLsAutomatically(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	autoconfigBaseURL := messages.ParsePayload(request).(messages.SetURLsAutomaticallyRequest).Value

	// make http request
	resp, err := http.Get(autoconfigBaseURL + "/api/config")
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// parse response
	var configResponse ConfigResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	err = json.Unmarshal(body, &configResponse)
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	cfg.ConfigFile.ApiUrl = configResponse.Environment.Api
	cfg.ConfigFile.IdentityUrl = configResponse.Environment.Identity
	cfg.ConfigFile.NotificationsUrl = configResponse.Environment.Notifications
	cfg.ConfigFile.VaultUrl = configResponse.Environment.Vault

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

func handleGetConfigEnvironment(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	return messages.IPCMessageFromPayload(messages.GetConfigEnvironmentResponse{
		ApiURL:           cfg.ConfigFile.ApiUrl,
		IdentityURL:      cfg.ConfigFile.IdentityUrl,
		NotificationsURL: cfg.ConfigFile.NotificationsUrl,
		VaultURL:         cfg.ConfigFile.VaultUrl,
	})
}

func handleSetClientID(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	clientID := messages.ParsePayload(request).(messages.SetClientIDRequest).Value
	err = cfg.SetClientID(clientID)
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

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

func handleSetClientSecret(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	clientSecret := messages.ParsePayload(request).(messages.SetClientSecretRequest).Value
	err = cfg.SetClientSecret(clientSecret)
	if err != nil {
		return messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: err.Error(),
		})
	}

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

func handleGetRuntimeConfig(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	return messages.IPCMessageFromPayload(messages.GetRuntimeConfigResponse{
		UseMemguard:          cfg.ConfigFile.RuntimeConfig.UseMemguard,
		SSHAgentSocketPath:   cfg.ConfigFile.RuntimeConfig.SSHAgentSocketPath,
		GoldwardenSocketPath: cfg.ConfigFile.RuntimeConfig.GoldwardenSocketPath,
	})
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetIdentityURLRequest{}), handleSetIdentity)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetApiURLRequest{}), handleSetApiURL)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetNotificationsURLRequest{}), handleSetNotifications)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetVaultURLRequest{}), handleSetVaultURL)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetURLsAutomaticallyRequest{}), handleSetURLsAutomatically)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.GetConfigEnvironmentRequest{}), handleGetConfigEnvironment)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.GetRuntimeConfigRequest{}), handleGetRuntimeConfig)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetClientIDRequest{}), handleSetClientID)
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.SetClientSecretRequest{}), handleSetClientSecret)
}
