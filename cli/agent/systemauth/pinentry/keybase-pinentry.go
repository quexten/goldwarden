//go:build windows

package pinentry

import (
	"errors"

	"github.com/keybase/client/go/logger"
	"github.com/keybase/client/go/protocol/keybase1"
	pinentry "github.com/quexten/goldwarden/cli/agent/systemauth/pinentry/keybase-pinentry"
)

func getPassword(title string, description string) (string, error) {
	pinentryInstance := pinentry.New("", logger.New(""), "")
	result, err := pinentryInstance.Get(keybase1.SecretEntryArg{
		Prompt: title,
		Desc:   description,
	})

	if err != nil {
		return "", err
	}

	if result.Canceled {
		return "", errors.New("Cancelled")
	}

	return result.Text, nil
}

func getApproval(title string, description string) (bool, error) {
	pinentryInstance := pinentry.New("", logger.New(""), "")
	result, err := pinentryInstance.Get(keybase1.SecretEntryArg{
		Prompt:     title,
		Desc:       description,
		Cancel:     "Decline",
		Ok:         "Approve",
		ShowTyping: true,
	})

	if err != nil {
		return false, err
	}

	if result.Canceled {
		return false, errors.New("Cancelled")
	}

	return true, nil
}
