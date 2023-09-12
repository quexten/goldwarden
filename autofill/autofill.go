//go:build !noautofill

package autofill

import (
	"errors"

	"github.com/atotto/clipboard"
	"github.com/quexten/goldwarden/autofill/autotype"
	"github.com/quexten/goldwarden/client"
	"github.com/quexten/goldwarden/ipc"
)

func GetLoginByUUID(uuid string, client client.Client) (ipc.DecryptedLoginCipher, error) {
	resp, err := client.SendToAgent(ipc.GetLoginRequest{
		UUID: uuid,
	})
	if err != nil {
		return ipc.DecryptedLoginCipher{}, err
	}

	switch resp.(type) {
	case ipc.GetLoginResponse:
		castedResponse := (resp.(ipc.GetLoginResponse))
		return castedResponse.Result, nil
	case ipc.ActionResponse:
		castedResponse := (resp.(ipc.ActionResponse))
		return ipc.DecryptedLoginCipher{}, errors.New("Error: " + castedResponse.Message)
	default:
		return ipc.DecryptedLoginCipher{}, errors.New("Wrong response type")
	}
}

func ListLogins(client client.Client) ([]ipc.DecryptedLoginCipher, error) {
	resp, err := client.SendToAgent(ipc.ListLoginsRequest{})
	if err != nil {
		return []ipc.DecryptedLoginCipher{}, err
	}

	switch resp.(type) {
	case ipc.GetLoginsResponse:
		castedResponse := (resp.(ipc.GetLoginsResponse))
		return castedResponse.Result, nil
	case ipc.ActionResponse:
		castedResponse := (resp.(ipc.ActionResponse))
		return []ipc.DecryptedLoginCipher{}, errors.New("Error: " + castedResponse.Message)
	default:
		return []ipc.DecryptedLoginCipher{}, errors.New("Wrong response type")
	}
}

func Run(layout string, useCopyPaste bool, client client.Client) {
	logins, err := ListLogins(client)
	if err != nil {
		panic(err)
	}

	autofillEntries := []AutofillEntry{}
	for _, login := range logins {
		autofillEntries = append(autofillEntries, AutofillEntry{
			Name:     login.Name,
			Username: login.Username,
			UUID:     login.UUID,
		})
	}

	RunAutofill(autofillEntries, func(uuid string, c chan bool) {
		login, err := GetLoginByUUID(uuid, client)
		if err != nil {
			panic(err)
		}

		if useCopyPaste {
			clipboard.WriteAll(string(login.Username))
			autotype.Paste(layout)
			autotype.TypeString("\t", layout)
			clipboard.WriteAll(login.Password)
			autotype.Paste(layout)
		} else {
			autotype.TypeString(string(login.Username), layout)
			autotype.TypeString("\t", layout)
			autotype.TypeString(string(login.Password), layout)
		}

		clipboard.WriteAll(login.TwoFactorCode)
		c <- true
	})
}
