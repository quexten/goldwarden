package actions

import (
	"fmt"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func handleGetCliCredentials(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx sockets.CallingContext) (response interface{}, err error) {
	req := request.ParsedPayload().(ipc.GetCLICredentialsRequest)

	if approved, err := systemauth.GetApproval("Approve Credential Access", fmt.Sprintf("%s on %s>%s>%s is trying to access credentials for %s", ctx.UserName, ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName, req.ApplicationName)); err != nil || !approved {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	env, found := vault.GetEnvCredentialForExecutable(req.ApplicationName)
	if !found {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "no credentials found for " + req.ApplicationName,
		})
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	response, err = ipc.IPCMessageFromPayload(ipc.GetCLICredentialsResponse{
		Env: env,
	})

	return
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageTypeGetCLICredentialsRequest, ensureEverything(systemauth.AccessCredential, handleGetCliCredentials))
}
