package systemauth

import (
	"errors"
	"os"

	"github.com/twpayne/go-pinentry"
)

var authDisabled = false

func GetPassword(title string, description string) (string, error) {
	client, err := pinentry.NewClient(
		pinentry.WithBinaryNameFromGnuPGAgentConf(),
		pinentry.WithGPGTTY(),
		pinentry.WithTitle(title),
		pinentry.WithDesc(description),
		pinentry.WithPrompt(title),
	)
	log.Info("Asking for pin |%s|%s|", title, description)

	if err != nil {
		return "", err
	}
	defer client.Close()

	switch pin, fromCache, err := client.GetPIN(); {
	case pinentry.IsCancelled(err):
		log.Info("Cancelled")
		return "", errors.New("Cancelled")
	case err != nil:
		return "", err
	case fromCache:
		log.Info("Got pin from cache")
		return pin, nil
	default:
		log.Info("Got pin from user")
		return pin, nil
	}
}

func GetApproval(title string, description string) (bool, error) {
	if authDisabled {
		log.Info("Skipping approval because system auth is disabled")
		return true, nil
	}

	client, err := pinentry.NewClient(
		pinentry.WithBinaryNameFromGnuPGAgentConf(),
		pinentry.WithGPGTTY(),
		pinentry.WithTitle(title),
		pinentry.WithDesc(description),
		pinentry.WithPrompt(title),
	)
	log.Info("Asking for approval |%s|%s|", title, description)

	if err != nil {
		return false, err
	}
	defer client.Close()

	switch _, err := client.Confirm("Confirm"); {
	case pinentry.IsCancelled(err):
		log.Info("Cancelled")
		return false, errors.New("Cancelled")
	case err != nil:
		return false, err
	default:
		log.Info("Got approval from user")
		return true, nil
	}
}

func init() {
	envAuthDisabled := os.Getenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED")
	if envAuthDisabled == "true" {
		authDisabled = true
	}
}
