package bitwarden

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/config"
)

type AuthRequestData struct {
	CreationDate       time.Time `json:"creationDate"`
	ID                 string    `json:"id"`
	Key                string    `json:"key"`
	MasterPasswordHash string    `json:"masterPasswordHash"`
	Object             string    `json:"object"`
	Origin             string    `json:"origin"`
	PublicKey          string    `json:"publicKey"`
	RequestApproved    bool      `json:"requestApproved"`
	RequestDeviceType  string    `json:"requestDeviceType"`
	RequestIpAddress   string    `json:"requestIpAddress"`
	ResponseDate       time.Time `json:"responseDate"`
}

type AuthRequestResponseData struct {
	DeviceIdentifier   string `json:"deviceIdentifier"`
	Key                string `json:"key"`
	MasterPasswordHash string `json:"masterPasswordHash"`
	Requestapproved    bool   `json:"requestApproved"`
}

func GetAuthRequest(ctx context.Context, requestUUID string, config *config.Config) (AuthRequestData, error) {
	var authRequest AuthRequestData
	err := authenticatedHTTPGet(ctx, config.ConfigFile.ApiUrl+"/auth-requests/"+requestUUID, &authRequest)
	return authRequest, err
}

func GetAuthRequests(ctx context.Context, config *config.Config) ([]AuthRequestData, error) {
	var authRequests []AuthRequestData
	err := authenticatedHTTPGet(ctx, config.ConfigFile.ApiUrl+"/auth-requests", &authRequests)
	return authRequests, err
}

func PutAuthRequest(ctx context.Context, requestUUID string, authRequest AuthRequestData, config *config.Config) error {
	var response interface{}
	err := authenticatedHTTPPut(ctx, config.ConfigFile.ApiUrl+"/auth-requests/"+requestUUID, &response, authRequest)
	return err
}

func CreateAuthResponse(ctx context.Context, authRequest AuthRequestData, keyring *crypto.Keyring, config *config.Config) (AuthRequestResponseData, error) {
	var authRequestResponse AuthRequestResponseData

	userSymmetricKey, err := config.GetUserSymmetricKey()
	if err != nil {
		return authRequestResponse, err
	}
	masterPasswordHash, err := config.GetMasterPasswordHash()
	if err != nil {
		return authRequestResponse, err
	}

	publicKey, err := base64.StdEncoding.DecodeString(authRequest.PublicKey)
	requesterKey, err := crypto.AssymmetricEncryptionKeyFromBytes(publicKey)

	encryptedUserSymmetricKey, err := crypto.EncryptWithAsymmetric(userSymmetricKey, requesterKey)
	if err != nil {
		panic(err)
	}
	encryptedMasterPasswordHash, err := crypto.EncryptWithAsymmetric(masterPasswordHash, requesterKey)
	if err != nil {
		panic(err)
	}

	err = authenticatedHTTPPut(ctx, config.ConfigFile.ApiUrl+"/auth-requests/"+authRequest.ID, &authRequestResponse, AuthRequestResponseData{
		DeviceIdentifier:   config.ConfigFile.DeviceUUID,
		Key:                string(encryptedUserSymmetricKey),
		MasterPasswordHash: string(encryptedMasterPasswordHash),
		Requestapproved:    true,
	})
	return authRequestResponse, err
}
