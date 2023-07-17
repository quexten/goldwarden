package bitwarden

import (
	"context"
	"fmt"

	"github.com/LlamaNite/llamalog"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/bitwarden/models"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/vault"
)

var log = llamalog.NewLogger("Goldwarden", "Bitwarden API")

func Sync(ctx context.Context, config *config.Config) (models.SyncData, error) {
	var sync models.SyncData
	if err := authenticatedHTTPGet(ctx, config.ConfigFile.ApiUrl+"/sync", &sync); err != nil {
		return models.SyncData{}, fmt.Errorf("could not sync: %v", err)
	}
	return sync, nil
}

func SyncToVault(ctx context.Context, vault *vault.Vault, config *config.Config, userSymmetricKey *crypto.SymmetricEncryptionKey) error {
	log.Info("Performing full sync...")

	sync, err := Sync(ctx, config)
	if err != nil {
		return err
	}

	if userSymmetricKey != nil {
		var orgKeys map[string]string = make(map[string]string)
		for _, org := range sync.Profile.Organizations {
			orgId := org.Id.String()
			orgKeys[orgId] = org.Key
		}
		crypto.InitKeyringFromUserSymmetricKey(vault.Keyring, *userSymmetricKey, sync.Profile.PrivateKey, orgKeys)
	}

	vault.Clear()
	for _, cipher := range sync.Ciphers {
		switch cipher.Type {
		case models.CipherLogin:
			vault.AddOrUpdateLogin(cipher)
			break
		case models.CipherNote:
			vault.AddOrUpdateSecureNote(cipher)
			break
		}
	}

	return nil
}
