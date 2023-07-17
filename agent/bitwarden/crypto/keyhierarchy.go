package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/awnumar/memguard"
	"golang.org/x/crypto/hkdf"
)

func InitKeyringFromMasterPassword(keyring *Keyring, accountKey EncString, accountPrivateKey EncString, orgKeys map[string]string, password memguard.LockedBuffer, email string, kdfConfig KDFConfig) error {
	masterKey, err := DeriveMasterKey(password, email, kdfConfig)
	if err != nil {
		return err
	}

	return InitKeyringFromMasterKey(keyring, accountKey, accountPrivateKey, orgKeys, masterKey)
}

func InitKeyringFromMasterKey(keyring *Keyring, accountKey EncString, accountPrivateKey EncString, orgKeys map[string]string, masterKey MasterKey) error {
	var accountSymmetricKeyByteArray []byte

	switch accountKey.Type {
	case AesCbc256_HmacSha256_B64:
		stretchedMasterKey, err := stretchKey(masterKey)
		if err != nil {
			return err
		}

		accountSymmetricKeyByteArray, err = DecryptWith(accountKey, stretchedMasterKey)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported account key type: %d", accountKey.Type)
	}

	accountSymmetricKey, err := SymmetricEncryptionKeyFromBytes(accountSymmetricKeyByteArray)
	if err != nil {
		return err
	}

	keyring.AccountKey = &accountSymmetricKey

	pkcs8PrivateKey, err := DecryptWith(accountPrivateKey, accountSymmetricKey)
	if err != nil {
		return err
	}
	keyring.AsymmetricEncyryptionKey = AsymmetricEncryptionKey{memguard.NewEnclave(pkcs8PrivateKey)}
	keyring.OrganizationKeys = orgKeys

	return nil
}

func InitKeyringFromUserSymmetricKey(keyring *Keyring, accountSymmetricKey SymmetricEncryptionKey, accountPrivateKey EncString, orgKeys map[string]string) error {
	keyring.AccountKey = &accountSymmetricKey

	pkcs8PrivateKey, err := DecryptWith(accountPrivateKey, accountSymmetricKey)
	if err != nil {
		return err
	}
	keyring.AsymmetricEncyryptionKey = AsymmetricEncryptionKey{memguard.NewEnclave(pkcs8PrivateKey)}
	keyring.OrganizationKeys = orgKeys

	return nil
}

func stretchKey(masterKey MasterKey) (SymmetricEncryptionKey, error) {
	key := make([]byte, 32)
	macKey := make([]byte, 32)

	buffer, err := masterKey.encKey.Open()
	if err != nil {
		return SymmetricEncryptionKey{}, err
	}

	var r io.Reader
	r = hkdf.Expand(sha256.New, buffer.Data(), []byte("enc"))
	r.Read(key)
	r = hkdf.Expand(sha256.New, buffer.Data(), []byte("mac"))
	r.Read(macKey)
	return SymmetricEncryptionKey{memguard.NewEnclave(key), memguard.NewEnclave(macKey)}, nil
}
