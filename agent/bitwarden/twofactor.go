package bitwarden

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/keys-pub/go-libfido2"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/systemauth"
)

type Fido2Response struct {
	Id         string `json:"id"`
	RawId      string `json:"rawId"`
	Type_      string `json:"type"`
	Extensions struct {
		Appid bool `json:"appid"`
	} `json:"extensions"`
	Response struct {
		AuthenticatorData string `json:"authenticatorData"`
		ClientDataJSON    string `json:"clientDataJson"`
		Signature         string `json:"signature"`
	} `json:"response"`
}

func Fido2TwoFactor(challengeB64 string, credentials []string, config *config.Config) (string, error) {
	url, err := url.Parse(config.ConfigFile.ApiUrl)
	rpid := url.Host

	locs, err := libfido2.DeviceLocations()
	if err != nil {
		return "", err
	}
	if len(locs) == 0 {
		return "", errors.New("no devices found")
	}

	path := locs[0].Path
	device, err := libfido2.NewDevice(path)
	if err != nil {
		return "", err
	}

	creds := make([][]byte, len(credentials))
	for i, cred := range credentials {
		decodedPublicKey, err := base64.RawURLEncoding.DecodeString(cred)
		if err != nil {
			websocketLog.Fatal(err.Error())
		}
		creds[i] = decodedPublicKey
	}

	clientDataJson := "{\"type\":\"webauthn.get\",\"challenge\":\"" + challengeB64 + "\",\"origin\":\"https://" + rpid + "\",\"crossOrigin\":false}"
	clientDataHash := sha256.Sum256([]byte(clientDataJson))
	clientDataJson = base64.URLEncoding.EncodeToString([]byte(clientDataJson))

	pin, err := systemauth.GetPassword("Fido2 PIN", "Enter your token's PIN")
	if err != nil {
		websocketLog.Fatal(err.Error())
	}

	assertion, err := device.Assertion(
		rpid,
		clientDataHash[:],
		creds,
		pin,
		&libfido2.AssertionOpts{
			Extensions: []libfido2.Extension{},
			UV:         libfido2.False,
		},
	)

	authDataRaw := assertion.AuthDataCBOR[2:] // first 2 bytes seem to be from cbor, don't have a proper decoder ATM but this works
	authData := base64.URLEncoding.EncodeToString(authDataRaw)
	sig := base64.URLEncoding.EncodeToString(assertion.Sig)
	credential := base64.URLEncoding.EncodeToString(assertion.CredentialID)

	resp := Fido2Response{
		Id:    credential,
		RawId: credential,
		Type_: "public-key",
		Extensions: struct {
			Appid bool `json:"appid"`
		}{Appid: false},
		Response: struct {
			AuthenticatorData string `json:"authenticatorData"`
			ClientDataJSON    string `json:"clientDataJson"`
			Signature         string `json:"signature"`
		}{
			AuthenticatorData: authData,
			ClientDataJSON:    clientDataJson,
			Signature:         sig,
		},
	}

	respjson, err := json.Marshal(resp)
	return string(respjson), nil
}

func performSecondFactor(resp *TwoFactorResponse, cfg *config.Config) (TwoFactorProvider, []byte, error) {
	if resp.TwoFactorProviders2[WebAuthn] != nil {
		chall := resp.TwoFactorProviders2[WebAuthn]["challenge"].(string)

		var creds []string
		for _, credential := range resp.TwoFactorProviders2[WebAuthn]["allowCredentials"].([]interface{}) {
			publicKey := credential.(map[string]interface{})["id"].(string)
			creds = append(creds, publicKey)
		}

		result, err := Fido2TwoFactor(chall, creds, cfg)
		if err != nil {
			return WebAuthn, nil, err
		}
		return WebAuthn, []byte(result), err
	}
	if resp.TwoFactorProviders2[Authenticator] != nil {
		token, err := systemauth.GetPassword("Authenticator Second Factor", "Enter your two-factor auth code")
		return Authenticator, []byte(token), err
	}
	if resp.TwoFactorProviders2[Email] != nil {
		token, err := systemauth.GetPassword("Email Second Factor", "Enter your two-factor auth code")
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

func urlValues(pairs ...string) url.Values {
	if len(pairs)%2 != 0 {
		panic("pairs must be of even length")
	}
	vals := make(url.Values)
	for i := 0; i < len(pairs); i += 2 {
		vals.Set(pairs[i], pairs[i+1])
	}
	return vals
}

var b64enc = base64.StdEncoding.Strict()
