package actions

import (
	"fmt"

	"github.com/quexten/goldwarden/cli/agent/config"
	"github.com/quexten/goldwarden/cli/agent/sockets"
	"github.com/quexten/goldwarden/cli/agent/systemauth"
	"github.com/quexten/goldwarden/cli/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/cli/agent/vault"
	"github.com/quexten/goldwarden/cli/ipc/messages"
)

func handleGetCliCredentials(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	req := messages.ParsePayload(request).(messages.GetCLICredentialsRequest)

	if approved, err := pinentry.GetApproval("Approve Credential Access", fmt.Sprintf("%s on %s>%s>%s is trying to access credentials for %s", ctx.UserName, ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName, req.ApplicationName)); err != nil || !approved {
		response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		if err != nil {
			return messages.IPCMessage{}, err
		}
		return response, nil
	}

	env, found := vault.GetEnvCredentialForExecutable(req.ApplicationName)
	if !found {
		response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: "no credentials found for " + req.ApplicationName,
		})
		if err != nil {
			return messages.IPCMessage{}, err
		}
		return response, nil
	}

	response, err = messages.IPCMessageFromPayload(messages.GetCLICredentialsResponse{
		Env: env,
	})

	return
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.GetCLICredentialsRequest{}), ensureEverything(systemauth.AccessVault, handleGetCliCredentials))
}
