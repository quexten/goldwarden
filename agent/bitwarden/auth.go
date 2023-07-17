package bitwarden

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/LlamaNite/llamalog"
	"github.com/awnumar/memguard"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"golang.org/x/crypto/pbkdf2"
)

var authLog = llamalog.NewLogger("Goldwarden", "Auth")

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

	password, err := systemauth.GetPassword("Bitwarden Password", "Enter your Bitwarden password")
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
		var twoFactor TwoFactorResponse
		if err := json.Unmarshal(errsc.body, &twoFactor); err != nil {
			return LoginResponseToken{}, crypto.MasterKey{}, "", err
		}
		provider, token, err := performSecondFactor(&twoFactor, cfg)
		if err != nil {
			return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("could not obtain two-factor auth token: %v", err)
		}
		values.Set("twoFactorProvider", strconv.Itoa(int(provider)))
		values.Set("twoFactorToken", string(token))
		values.Set("twoFactorRemember", "1")
		loginResponseToken = LoginResponseToken{}
		if err := authenticatedHTTPPost(ctx, cfg.ConfigFile.IdentityUrl+"/connect/token", &loginResponseToken, values); err != nil {
			return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("could not login via two-factor: %v", err)
		}
		authLog.Info("2FA login successful")
	} else if err != nil && strings.Contains(err.Error(), "Captcha required.") {
		return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("captcha required, please login via the web interface")

	} else if err != nil {
		return LoginResponseToken{}, crypto.MasterKey{}, "", fmt.Errorf("could not login via password: %v", err)
	}

	authLog.Info("Logged in")
	return loginResponseToken, masterKey, hashedPassword, nil
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
