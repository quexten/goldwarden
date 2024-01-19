package actions

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/notify"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/systemauth/biometrics"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/agent/vault"

	"github.com/quexten/goldwarden/ipc/messages"
)

func handleGetBiometricsKey(request messages.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response messages.IPCMessage, err error) {
	actionsLog.Info("Browser Biometrics: Key requested, verifying biometrics...")
	if !(systemauth.VerifyPinSession(*ctx) || biometrics.CheckBiometrics(biometrics.BrowserBiometrics)) {
		response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		actionsLog.Info("Browser Biometrics: Biometrics not approved %v", err)
		if err != nil {
			return messages.IPCMessage{}, err
		}
		return response, nil
	}

	actionsLog.Info("Browser Biometrics: Biometrics verified, asking for approval...")
	if approved, err := pinentry.GetApproval("Approve Credential Access", fmt.Sprintf("%s on %s>%s>%s is trying to access your vault encryption key for browser biometric unlock.", ctx.UserName, ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName)); err != nil || !approved {
		response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		actionsLog.Info("Browser Biometrics: Biometrics not approved %v", err)
		if err != nil {
			return messages.IPCMessage{}, err
		}
		return response, nil
	}

	actionsLog.Info("Browser Biometrics: Approved, getting key...")
	masterKey, err := cfg.GetMasterKey()
	if err != nil {
		return messages.IPCMessage{}, err
	}
	masterKeyB64 := base64.StdEncoding.EncodeToString(masterKey)
	response, err = messages.IPCMessageFromPayload(messages.GetBiometricsKeyResponse{
		Key: masterKeyB64,
	})
	actionsLog.Info("Browser Biometrics: Sending key...")
	notify.Notify("Goldwarden", "Unlocked Browser Extension", "", 10*time.Second, func() {})
	return response, err
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.GetBiometricsKeyRequest{}), ensureIsNotLocked(ensureIsLoggedIn(handleGetBiometricsKey)))
}
