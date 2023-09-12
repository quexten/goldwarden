package twofactor

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/logging"
)

var twofactorLog = logging.GetLogger("Goldwarden", "TwoFactor")

func PerformSecondFactor(resp *TwoFactorResponse, cfg *config.Config) (TwoFactorProvider, []byte, error) {
	if provider, isInMap := resp.TwoFactorProviders2[WebAuthn]; isInMap {
		if isFido2Enabled {
			chall := provider["challenge"].(string)

			var creds []string
			for _, credential := range provider["allowCredentials"].([]interface{}) {
				publicKey := credential.(map[string]interface{})["id"].(string)
				creds = append(creds, publicKey)
			}

			result, err := Fido2TwoFactor(chall, creds, cfg)
			if err != nil {
				return WebAuthn, nil, err
			}
			return WebAuthn, []byte(result), err
		} else {
			twofactorLog.Warn("WebAuthn is enabled for the account but goldwarden is not compiled with FIDO2 support")
		}
	}
	if _, isInMap := resp.TwoFactorProviders2[Authenticator]; isInMap {
		token, err := pinentry.GetPassword("Authenticator Second Factor", "Enter your two-factor auth code")
		return Authenticator, []byte(token), err
	}
	if _, isInMap := resp.TwoFactorProviders2[Email]; isInMap {
		token, err := pinentry.GetPassword("Email Second Factor", "Enter your two-factor auth code")
		return Email, []byte(token), err
	}

	return Authenticator, []byte{}, errors.New("no second factor available")
}

type TwoFactorProvider int

const (
	Authenticator         TwoFactorProvider = 0
	Email                 TwoFactorProvider = 1
	Duo                   TwoFactorProvider = 2 //Not supported
	YubiKey               TwoFactorProvider = 3 //Not supported
	U2f                   TwoFactorProvider = 4 //Not supported
	Remember              TwoFactorProvider = 5 //Not supported
	OrganizationDuo       TwoFactorProvider = 6 //Not supported
	WebAuthn              TwoFactorProvider = 7
	_TwoFactorProviderMax                   = 8 //Not supported
)

func (t *TwoFactorProvider) UnmarshalText(text []byte) error {
	i, err := strconv.Atoi(string(text))
	if err != nil || i < 0 || i >= _TwoFactorProviderMax {
		return fmt.Errorf("invalid two-factor auth provider: %q", text)
	}
	*t = TwoFactorProvider(i)
	return nil
}

type TwoFactorResponse struct {
	TwoFactorProviders2 map[TwoFactorProvider]map[string]interface{}
}
