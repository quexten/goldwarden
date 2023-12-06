//go:build !noautofill

package autofill

import (
	"errors"

	"github.com/atotto/clipboard"
	"github.com/quexten/goldwarden/autofill/autotype"
	"github.com/quexten/goldwarden/client"
	"github.com/quexten/goldwarden/ipc/messages"
)

func GetLoginByUUID(uuid string, client client.Client) (messages.DecryptedLoginCipher, error) {
	resp, err := client.SendToAgent(messages.GetLoginRequest{
		UUID: uuid,
	})
	if err != nil {
		return messages.DecryptedLoginCipher{}, err
	}

	switch resp.(type) {
	case messages.GetLoginResponse:
		castedResponse := (resp.(messages.GetLoginResponse))
		return castedResponse.Result, nil
	case messages.ActionResponse:
		castedResponse := (resp.(messages.ActionResponse))
		return messages.DecryptedLoginCipher{}, errors.New("Error: " + castedResponse.Message)
	default:
		return messages.DecryptedLoginCipher{}, errors.New("Wrong response type")
	}
}

func ListLogins(client client.Client) ([]messages.DecryptedLoginCipher, error) {
	resp, err := client.SendToAgent(messages.ListLoginsRequest{})
	if err != nil {
		return []messages.DecryptedLoginCipher{}, err
	}

	switch resp.(type) {
	case messages.GetLoginsResponse:
		castedResponse := (resp.(messages.GetLoginsResponse))
		return castedResponse.Result, nil
	case messages.ActionResponse:
		castedResponse := (resp.(messages.ActionResponse))
		return []messages.DecryptedLoginCipher{}, errors.New("Error: " + castedResponse.Message)
	default:
		return []messages.DecryptedLoginCipher{}, errors.New("Wrong response type")
	}
}

func Run(layout string, client client.Client) {
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

		autotype.TypeString(string(login.Username)+"\t"+string(login.Password), layout)

		clipboard.WriteAll(login.TwoFactorCode)
		c <- true
	})
}
