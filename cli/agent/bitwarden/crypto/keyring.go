package crypto

import (
	"errors"

	"github.com/quexten/goldwarden/cli/logging"
)

var keyringLog = logging.GetLogger("Goldwarden", "Keyring")

type Keyring struct {
	isLocked                 bool
	accountKey               SymmetricEncryptionKey
	AsymmetricEncyryptionKey AsymmetricEncryptionKey
	IsMemguard               bool
	OrganizationKeys         map[string]string
}

func NewMemoryKeyring(accountKey *MemorySymmetricEncryptionKey) Keyring {
	keyringLog.Info("Creating new memory keyring")
	return Keyring{
		isLocked:   accountKey == nil,
		accountKey: accountKey,
	}
}

func NewMemguardKeyring(accountKey *MemguardSymmetricEncryptionKey) Keyring {
	keyringLog.Info("Creating new memguard keyring")
	return Keyring{
		isLocked:   accountKey == nil,
		accountKey: accountKey,
	}
}

func (keyring Keyring) IsLocked() bool {
	return keyring.isLocked
}

func (keyring *Keyring) Lock() {
	keyringLog.Info("Locking keyring")
	keyring.isLocked = true
	keyring.accountKey = nil
	keyring.AsymmetricEncyryptionKey = MemoryAsymmetricEncryptionKey{}
	keyring.OrganizationKeys = nil
}

func (keyring *Keyring) UnlockWithAccountKey(accountKey SymmetricEncryptionKey) {
	keyringLog.Info("Unlocking keyring with account key")
	keyring.isLocked = false
	keyring.accountKey = accountKey
}

func (keyring *Keyring) GetAccountKey() SymmetricEncryptionKey {
	return keyring.accountKey
}

func (keyring *Keyring) GetSymmetricKeyForOrganization(uuid string) (SymmetricEncryptionKey, error) {
	if key, ok := keyring.OrganizationKeys[uuid]; ok {
		decryptedOrgKey, err := DecryptWithAsymmetric([]byte(key), keyring.AsymmetricEncyryptionKey)
		if err != nil {
			return MemorySymmetricEncryptionKey{}, err
		}

		return MemorySymmetricEncryptionKeyFromBytes(decryptedOrgKey)
	}
	return MemorySymmetricEncryptionKey{}, errors.New("no key found for organization")
}
