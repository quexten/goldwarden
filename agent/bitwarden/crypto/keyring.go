package crypto

import (
	"errors"
)

type Keyring struct {
	AccountKey               *SymmetricEncryptionKey
	AsymmetricEncyryptionKey AsymmetricEncryptionKey
	OrganizationKeys         map[string]string
}

func NewKeyring(accountKey *SymmetricEncryptionKey) Keyring {
	return Keyring{
		AccountKey: accountKey,
	}
}

func (keyring Keyring) IsLocked() bool {
	return keyring.AccountKey == nil
}

func (keyring *Keyring) Lock() {
	keyring.AccountKey = nil
	keyring.AsymmetricEncyryptionKey = AsymmetricEncryptionKey{}
	keyring.OrganizationKeys = nil
}

func (keyring *Keyring) GetSymmetricKeyForOrganization(uuid string) (SymmetricEncryptionKey, error) {
	if key, ok := keyring.OrganizationKeys[uuid]; ok {
		decryptedOrgKey, err := DecryptWithAsymmetric([]byte(key), keyring.AsymmetricEncyryptionKey)
		if err != nil {
			return SymmetricEncryptionKey{}, err
		}

		return SymmetricEncryptionKeyFromBytes(decryptedOrgKey)
	}
	return SymmetricEncryptionKey{}, errors.New("no key found for organization")
}
