package actions

import (
	"context"
	"fmt"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
)

func handleCreateSend(msg messages.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response messages.IPCMessage, err error) {
	token, err := cfg.GetToken()
	if err != nil {
		return messages.IPCMessage{}, fmt.Errorf("error getting token: %w", err)
	}
	parsedMsg := messages.ParsePayload(msg).(messages.CreateSendRequest)

	ctx := context.WithValue(context.TODO(), bitwarden.AuthToken{}, token.AccessToken)
	url, err := bitwarden.CreateSend(ctx, cfg, vault, parsedMsg.Name, parsedMsg.Text)
	if err != nil {
		actionsLog.Warn(err.Error())
	}

	response, err = messages.IPCMessageFromPayload(messages.CreateSendResponse{
		URL: url,
	})
	return
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.CreateSendRequest{}), ensureIsNotLocked(ensureIsLoggedIn(handleCreateSend)))
}
