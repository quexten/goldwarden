package bitwarden

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
	"time"

	"github.com/quexten/goldwarden/cli/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/cli/agent/config"
	"github.com/quexten/goldwarden/cli/agent/vault"
	"golang.org/x/crypto/hkdf"
)

type SendFileMetadata struct {
	FileName string `json:"fileName"`
	Id       string `json:"id"`
	Size     int    `json:"size"`
	SizeName string `json:"sizeName"`
}

type SendTextMetadata struct {
	Hidden   bool    `json:"hidden"`
	Response *string `json:"response"`
	Text     string  `json:"text"`
}

type SendMetadata struct {
	CreatorIdentifier string
	ExpirationDate    string
	File              SendFileMetadata
	Id                string
	Name              string
	Object            string
	Text              SendTextMetadata
	Type              int
}

type SendCreateRequest struct {
	AccessCount    *int             `json:"accessCount"`
	AccessId       *string          `json:"accessId"`
	DeletionDate   string           `json:"deletionDate"`
	Disabled       bool             `json:"disabled"`
	ExpirationDate *string          `json:"expirationDate"`
	HideEmail      bool             `json:"hideEmail"`
	Key            string           `json:"key"`
	MaxAccessCount *int             `json:"maxAccessCount"`
	Name           string           `json:"name"`
	Notes          *string          `json:"notes"`
	Text           SendTextMetadata `json:"text"`
	Type           int              `json:"type"`
}

func CreateSend(ctx context.Context, cfg *config.Config, vault *vault.Vault, name string, text string) (string, error) {
	timestampIn14Days := time.Now().AddDate(0, 0, 14)
	timestampIn14DaysStr := timestampIn14Days.Format("2006-01-02T15:04:05Z")

	// generate 32 byte key
	sendSourceKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, sendSourceKey)
	if err != nil {
		return "", err
	}

	encryptedSendSourceKey, err := crypto.EncryptWithToString(sendSourceKey, crypto.AesCbc256_HmacSha256_B64, vault.Keyring.GetAccountKey())
	if err != nil {
		return "", err
	}

	sendUseKeyPairBytes := make([]byte, 64)
	_, err = hkdf.New(sha256.New, sendSourceKey, []byte("bitwarden-send"), []byte("send")).Read(sendUseKeyPairBytes)
	if err != nil {
		return "", err
	}

	sendUseKeyPair, err := crypto.MemorySymmetricEncryptionKeyFromBytes(sendUseKeyPairBytes)
	if err != nil {
		return "", err
	}

	encryptedName, err := crypto.EncryptWithToString([]byte(name), crypto.AesCbc256_HmacSha256_B64, sendUseKeyPair)
	if err != nil {
		return "", err
	}

	encryptedText, err := crypto.EncryptWithToString([]byte(text), crypto.AesCbc256_HmacSha256_B64, sendUseKeyPair)
	if err != nil {
		return "", err
	}

	sendRequest := SendCreateRequest{
		DeletionDate: timestampIn14DaysStr,
		Disabled:     false,
		HideEmail:    false,
		Key:          encryptedSendSourceKey,
		Name:         encryptedName,
		Text: SendTextMetadata{
			Hidden: false,
			Text:   encryptedText,
		},
		Type: 0,
	}

	var result SendCreateRequest
	err = authenticatedHTTPPost(ctx, cfg.ConfigFile.ApiUrl+"/sends", &result, sendRequest)
	if err != nil {
		return "", err
	}

	return cfg.ConfigFile.VaultUrl + "/#/send/" + *result.AccessId + "/" + strings.ReplaceAll(base64.RawURLEncoding.EncodeToString(sendSourceKey), "+", "-"), nil
}
