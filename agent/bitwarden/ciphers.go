package bitwarden

import (
	"context"

	"github.com/quexten/goldwarden/agent/bitwarden/models"
	"github.com/quexten/goldwarden/agent/config"
)

func PostCipher(ctx context.Context, cipher models.Cipher, cfg *config.Config) (models.Cipher, error) {
	var resultingCipher models.Cipher
	err := authenticatedHTTPPost(ctx, cfg.ConfigFile.ApiUrl+"/ciphers", &resultingCipher, cipher)
	return resultingCipher, err
}

func GetCipher(ctx context.Context, uuid string, cfg *config.Config) (models.Cipher, error) {
	var cipher models.Cipher
	err := authenticatedHTTPGet(ctx, cfg.ConfigFile.ApiUrl+"/ciphers/"+uuid, &cipher)
	return cipher, err
}

func DeleteCipher(ctx context.Context, uuid string, cfg *config.Config) error {
	var result interface{}
	err := authenticatedHTTPDelete(ctx, cfg.ConfigFile.ApiUrl+"/ciphers/"+uuid, &result)
	return err
}

func PutCipher(ctx context.Context, uuid string, cipher models.Cipher, cfg *config.Config) (models.Cipher, error) {
	var resultingCipher models.Cipher
	err := authenticatedHTTPPut(ctx, cfg.ConfigFile.ApiUrl+"/ciphers/"+uuid, &resultingCipher, cipher)
	return resultingCipher, err
}
