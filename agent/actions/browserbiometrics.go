package actions

import (
	"encoding/base64"
	"fmt"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func handleGetBiometricsKey(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx sockets.CallingContext) (response interface{}, err error) {
	if approved, err := systemauth.GetApproval("Approve Credential Access", fmt.Sprintf("%s on %s>%s>%s is trying to access your vault encryption key for browser biometric unlock.", ctx.UserName, ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName)); err != nil || !approved {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	masterKey, err := cfg.GetMasterKey()
	masterKeyB64 := base64.StdEncoding.EncodeToString(masterKey)
	response, err = ipc.IPCMessageFromPayload(ipc.GetBiometricsKeyResponse{
		Key: masterKeyB64,
	})
	return response, err
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageTypeGetBiometricsKeyRequest, ensureEverything(systemauth.BrowserBiometrics, handleGetBiometricsKey))
}
