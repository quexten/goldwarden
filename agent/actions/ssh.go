package actions

import (
	"context"
	"strings"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/ssh"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/quexten/goldwarden/logging"
)

var actionsLog = logging.GetLogger("Goldwarden", "Actions")

func handleAddSSH(msg messages.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response messages.IPCMessage, err error) {
	req := messages.ParsePayload(msg).(messages.CreateSSHKeyRequest)

	cipher, publicKey := ssh.NewSSHKeyCipher(req.Name, vault.Keyring)
	response, err = messages.IPCMessageFromPayload(messages.ActionResponse{
		Success: true,
	})
	if err != nil {
		panic(err)
	}

	token, err := cfg.GetToken()
	ctx := context.WithValue(context.TODO(), bitwarden.AuthToken{}, token.AccessToken)
	ciph, err := bitwarden.PostCipher(ctx, cipher, cfg)
	if err == nil {
		vault.AddOrUpdateSecureNote(ciph)
	} else {
		actionsLog.Warn("Error posting ssh key cipher: " + err.Error())
	}

	response, err = messages.IPCMessageFromPayload(messages.CreateSSHKeyResponse{
		Digest: strings.ReplaceAll(publicKey, "\n", "") + " " + req.Name,
	})

	return
}

func handleListSSH(msg messages.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response messages.IPCMessage, err error) {
	keys := vault.GetSSHKeys()
	keyStrings := make([]string, 0)
	for _, key := range keys {
		keyStrings = append(keyStrings, strings.ReplaceAll(key.PublicKey+" "+key.Name, "\n", ""))
	}

	response, err = messages.IPCMessageFromPayload(messages.GetSSHKeysResponse{
		Keys: keyStrings,
	})
	return
}

func init() {
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.CreateSSHKeyRequest{}), ensureEverything(systemauth.SSHKey, handleAddSSH))
	AgentActionsRegistry.Register(messages.MessageTypeForEmptyPayload(messages.GetSSHKeysRequest{}), ensureIsNotLocked(ensureIsLoggedIn(handleListSSH)))
}
