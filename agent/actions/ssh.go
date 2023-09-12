package actions

import (
	"context"
	"strings"

	"github.com/quexten/goldwarden/agent/bitwarden"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/ssh"
	"github.com/quexten/goldwarden/agent/systemauth/biometrics"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
	"github.com/quexten/goldwarden/logging"
)

var actionsLog = logging.GetLogger("Goldwarden", "Actions")

func handleAddSSH(msg ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response ipc.IPCMessage, err error) {
	req := msg.ParsedPayload().(ipc.CreateSSHKeyRequest)

	cipher, publicKey := ssh.NewSSHKeyCipher(req.Name, vault.Keyring)
	response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
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

	response, err = ipc.IPCMessageFromPayload(ipc.CreateSSHKeyResponse{
		Digest: strings.ReplaceAll(publicKey, "\n", "") + " " + req.Name,
	})

	return
}

func handleListSSH(msg ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, callingContext *sockets.CallingContext) (response ipc.IPCMessage, err error) {
	keys := vault.GetSSHKeys()
	keyStrings := make([]string, 0)
	for _, key := range keys {
		keyStrings = append(keyStrings, strings.ReplaceAll(key.PublicKey+" "+key.Name, "\n", ""))
	}

	response, err = ipc.IPCMessageFromPayload(ipc.GetSSHKeysResponse{
		Keys: keyStrings,
	})
	return
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageTypeCreateSSHKeyRequest, ensureEverything(biometrics.SSHKey, handleAddSSH))
	AgentActionsRegistry.Register(ipc.IPCMessageTypeGetSSHKeysRequest, ensureIsNotLocked(ensureIsLoggedIn(handleListSSH)))
}
