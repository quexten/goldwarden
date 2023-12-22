package bitwarden

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/awnumar/memguard"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/bitwarden/twofactor"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/logging"
	"golang.org/x/crypto/pbkdf2"
)

var authLog = logging.GetLogger("Goldwarden", "Auth")

type preLoginRequest struct {
	Email string `json:"email"`
}

type preLoginResponse struct {
	KDF            int
	KDFIterations  int
	KDFMemory      int
	KDFParallelism int
}

type LoginResponseToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Key          string `json:"key"`
}

const (
	deviceName       = "goldwarden"
	loginScope       = "api offline_access"
	loginApiKeyScope = "api"
)

func deviceType() string {
	switch runtime.GOOS {
	case "linux":
		return "8"
	case "darwin":
		return "7"
	case "windows":
		return "6"
	default:
		return "14"
	}
}

func LoginWithMasterpassword(ctx context.Context, email string, cfg *config.Config, vault *vault.Vault) (LoginResponseToken, crypto.MasterKey, string, error) {
	var preLogin preLoginResponse
	if err := authenticatedHTTPPost(ctx, cfg.ConfigFile.ApiUrl+"/accounts/prelogin", &preLogin, preLoginRequest{
		Email: email,
	}); err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("could not pre-login: %v", err)
	}

	var values url.Values
	var masterKey crypto.MasterKey
	var hashedPassword string

	password, err := pinentry.GetPassword("Bitwarden Password", "Enter your Bitwarden password")
	if err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", err
	}

	masterKey, err = crypto.DeriveMasterKey(*memguard.NewBufferFromBytes([]byte(strings.Clone(password))), email, crypto.KDFConfig{Type: crypto.KDFType(preLogin.KDF), Iterations: uint32(preLogin.KDFIterations), Memory: uint32(preLogin.KDFMemory), Parallelism: uint32(preLogin.KDFParallelism)})
	if err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", err
	}

	hashedPassword = b64enc.EncodeToString(pbkdf2.Key(masterKey.GetBytes(), []byte(password), 1, 32, sha256.New))

	values = urlValues(
		"grant_type", "password",
		"username", email,
		"password", string(hashedPassword),
		"scope", loginScope,
		"client_id", "connector",
		"deviceType", deviceType(),
		"deviceName", deviceName,
		"deviceIdentifier", cfg.ConfigFile.DeviceUUID,
	)

	var loginResponseToken LoginResponseToken
	err = authenticatedHTTPPost(ctx, cfg.ConfigFile.IdentityUrl+"/connect/token", &loginResponseToken, values)
	errsc, ok := err.(*errStatusCode)
	if ok && bytes.Contains(errsc.body, []byte("TwoFactor")) {
		loginResponseToken, err = Perform2FA(values, errsc, cfg, ctx)
		if err != nil {
			return LoginResponseToken{}, crypto.MasterKey{}, "", err
		}
	} else if err != nil && strings.Contains(err.Error(), "Captcha required.") {
		return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("captcha required, please login via the web interface")

	} else if err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("could not login via password: %v", err)
	}

	authLog.Info("Logged in")
	return loginResponseToken, masterKey, hashedPassword, nil
}

func LoginWithDevice(ctx context.Context, email string, cfg *config.Config, vault *vault.Vault) (LoginResponseToken, crypto.MasterKey, string, error) {
	timeout := 120 * time.Second

	// 25 random letters & numbers
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	accessCode := ""
	for i := 0; i < 25; i++ {
		accessCode += string(alphabet[rand.Intn(len(alphabet))])
	}
	publicKey, err := crypto.GenerateAsymmetric(vault.Keyring.IsMemguard)
	if err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", err
	}
	data, err := CreateAuthRequest(ctx, accessCode, cfg.ConfigFile.DeviceUUID, email, base64.StdEncoding.EncodeToString(publicKey.PublicBytes()), cfg)
	if err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", err
	}

	timeoutChan := make(chan bool)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- true
	}()

	for {
		select {
		case <-timeoutChan:
			return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("timed out waiting for device to be authorized")
		default:
			authRequestData, err := GetAuthResponse(ctx, accessCode, data.ID, cfg)
			if err != nil {
				log.Error("Could not get auth request: %s", err.Error())
			}
			if authRequestData.RequestApproved {
				masterKey, err := crypto.DecryptWithAsymmetric([]byte(authRequestData.Key), publicKey)
				masterPasswordHash, err := crypto.DecryptWithAsymmetric([]byte(authRequestData.MasterPasswordHash), publicKey)
				values := urlValues(
					"grant_type", "password",
					"username", email,
					"password", string(accessCode),
					"authRequest", authRequestData.ID,
					"scope", loginScope,
					"client_id", "connector",
					"deviceType", deviceType(),
					"deviceName", deviceName,
					"deviceIdentifier", cfg.ConfigFile.DeviceUUID,
				)

				var loginResponseToken LoginResponseToken
				err = authenticatedHTTPPost(ctx, cfg.ConfigFile.IdentityUrl+"/connect/token", &loginResponseToken, values)
				errsc, ok := err.(*errStatusCode)
				if ok && bytes.Contains(errsc.body, []byte("TwoFactor")) {
					loginResponseToken, err = Perform2FA(values, errsc, cfg, ctx)
					if err != nil {
						return LoginResponseToken{}, crypto.MasterKey{}, "", err
					}
				} else if err != nil && strings.Contains(err.Error(), "Captcha required.") {
					return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("captcha required, please login via the web interface")

				} else if err != nil {
					return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("could not login via password: %v", err)
				}
				return loginResponseToken, crypto.MasterKeyFromBytes(masterKey), string(masterPasswordHash), nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func RefreshToken(ctx context.Context, cfg *config.Config) bool {
	authLog.Info("Refreshing token")

	token, err := cfg.GetToken()
	if err != nil {
		fmt.Println("Could not get refresh token: ", err)
		return false
	}

	var loginResponseToken LoginResponseToken
	err = authenticatedHTTPPost(ctx, cfg.ConfigFile.IdentityUrl+"/connect/token", &loginResponseToken, urlValues(
		"grant_type", "refresh_token",
		"refresh_token", token.RefreshToken,
		"client_id", "connector",
	))
	if err != nil {
		fmt.Println("Could not refresh token: ", err)
		return false
	}
	cfg.SetToken(config.LoginToken{
		AccessToken:  loginResponseToken.AccessToken,
		RefreshToken: loginResponseToken.RefreshToken,
		Key:          loginResponseToken.Key,
		TokenType:    loginResponseToken.TokenType,
		ExpiresIn:    loginResponseToken.ExpiresIn,
	})

	authLog.Info("Token refreshed")

	return true
}

func Perform2FA(values url.Values, errsc *errStatusCode, cfg *config.Config, ctx context.Context) (LoginResponseToken, error) {
	var twoFactor twofactor.TwoFactorResponse
	if err := json.Unmarshal(errsc.body, &twoFactor); err != nil {
		return LoginResponseToken{}, err
	}
	provider, token, err := twofactor.PerformSecondFactor(&twoFactor, cfg)
	if err != nil {
		return LoginResponseToken{}, fmt.Errorf("could not obtain two-factor auth token: %v", err)
	}
	values.Set("twoFactorProvider", strconv.Itoa(int(provider)))
	values.Set("twoFactorToken", string(token))
	values.Set("twoFactorRemember", "1")
	loginResponseToken := LoginResponseToken{}
	if err := authenticatedHTTPPost(ctx, cfg.ConfigFile.IdentityUrl+"/connect/token", &loginResponseToken, values); err != nil {
		return LoginResponseToken{}, fmt.Errorf("could not login via two-factor: %v", err)
	}
	authLog.Info("2FA login successful")
	return loginResponseToken, nil
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
