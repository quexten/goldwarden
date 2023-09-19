package actions

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/ipc"
)

func handleGetLoginCipher(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response ipc.IPCMessage, err error) {
	req := request.ParsedPayload().(ipc.GetLoginRequest)
	login, err := vault.GetLoginByFilter(req.UUID, req.OrgId, req.Name, req.Username)
	if err != nil {
		return ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "login not found",
		})
	}

	cipherKey, err := login.GetKeyForCipher(*vault.Keyring)
	if err != nil {
		return ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "could not get cipher key",
		})
	}

	decryptedLogin := ipc.DecryptedLoginCipher{
		Name: "NO NAME FOUND",
	}
	decryptedLogin.UUID = login.ID.String()
	if login.OrganizationID != nil {
		decryptedLogin.OrgaizationID = login.OrganizationID.String()
	}

	if !login.Name.IsNull() {
		decryptedName, err := crypto.DecryptWith(login.Name, cipherKey)
		if err == nil {
			decryptedLogin.Name = string(decryptedName)
		}
	}

	if !login.Login.Username.IsNull() {
		decryptedUsername, err := crypto.DecryptWith(login.Login.Username, cipherKey)
		if err == nil {
			decryptedLogin.Username = string(decryptedUsername)
		}
	}

	if !login.Login.Password.IsNull() {
		decryptedPassword, err := crypto.DecryptWith(login.Login.Password, cipherKey)
		if err == nil {
			decryptedLogin.Password = string(decryptedPassword)
		}
	}

	if !(login.Notes == nil) && !login.Notes.IsNull() {
		decryptedNotes, err := crypto.DecryptWith(*login.Notes, cipherKey)
		if err == nil {
			decryptedLogin.Notes = string(decryptedNotes)
		}
	}
	if !login.Login.Totp.IsNull() {
		decryptedTotp, err := crypto.DecryptWith(login.Login.Totp, cipherKey)
		if err == nil {
			code, err := totp.GenerateCode(string(strings.ReplaceAll(string(decryptedTotp), " ", "")), time.Now())
			if err == nil {
				decryptedLogin.TwoFactorCode = code
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(string(decryptedTotp))
		}
	}

	if approved, err := pinentry.GetApproval("Approve Credential Access", fmt.Sprintf("%s on %s>%s>%s is trying to access credentials for user %s on entry %s", ctx.UserName, ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName, decryptedLogin.Username, decryptedLogin.Name)); err != nil || !approved {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		if err != nil {
			return ipc.IPCMessage{}, err
		}
		return response, nil
	}

	return ipc.IPCMessageFromPayload(ipc.GetLoginResponse{
		Found:  true,
		Result: decryptedLogin,
	})
}

func handleListLoginsRequest(request ipc.IPCMessage, cfg *config.Config, vault *vault.Vault, ctx *sockets.CallingContext) (response ipc.IPCMessage, err error) {
	if approved, err := pinentry.GetApproval("Approve List Credentials", fmt.Sprintf("%s on %s>%s>%s is trying to list credentials (name & username)", ctx.UserName, ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName)); err != nil || !approved {
		response, err = ipc.IPCMessageFromPayload(ipc.ActionResponse{
			Success: false,
			Message: "not approved",
		})
		if err != nil {
			return ipc.IPCMessage{}, err
		}
		return response, nil
	}

	logins := vault.GetLogins()
	decryptedLoginCiphers := make([]ipc.DecryptedLoginCipher, 0)
	for _, login := range logins {
		key, err := login.GetKeyForCipher(*vault.Keyring)
		if err != nil {
			actionsLog.Warn("Could not decrypt login:" + err.Error())
			continue
		}

		var decryptedName []byte = []byte{}
		var decryptedUsername []byte = []byte{}

		if !login.Name.IsNull() {
			decryptedName, err = crypto.DecryptWith(login.Name, key)
			if err != nil {
				actionsLog.Warn("Could not decrypt login:" + err.Error())
				continue
			}
		}

		if !login.Login.Username.IsNull() {
			decryptedUsername, err = crypto.DecryptWith(login.Login.Username, key)
			if err != nil {
				actionsLog.Warn("Could not decrypt login:" + err.Error())
				continue
			}
		}

		decryptedLoginCiphers = append(decryptedLoginCiphers, ipc.DecryptedLoginCipher{
			Name:     string(decryptedName),
			Username: string(decryptedUsername),
			UUID:     login.ID.String(),
		})

		// prevent deadlock from enclaves
		debug.FreeOSMemory()
	}

	return ipc.IPCMessageFromPayload(ipc.GetLoginsResponse{
		Found:  len(decryptedLoginCiphers) > 0,
		Result: decryptedLoginCiphers,
	})
}

func init() {
	AgentActionsRegistry.Register(ipc.IPCMessageGetLoginRequest, ensureEverything(systemauth.AccessVault, handleGetLoginCipher))
	AgentActionsRegistry.Register(ipc.IPCMessageListLoginsRequest, ensureEverything(systemauth.AccessVault, handleListLoginsRequest))
}
