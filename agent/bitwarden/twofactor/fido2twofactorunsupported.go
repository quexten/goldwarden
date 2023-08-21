//go:build nofido2

package twofactor

import (
	"errors"
	"github.com/quexten/goldwarden/agent/config"
)

const isFido2Enabled = false

func Fido2TwoFactor(challengeB64 string, credentials []string, config *config.Config) (string, error) {
	return "", errors.New("Fido2 is not enabled")
}
