package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"io"

	"github.com/mikesmitty/edkey"
	"github.com/quexten/goldwarden/cli/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/cli/agent/bitwarden/models"
	"golang.org/x/crypto/ssh"
)

// todo refactor to share code
func SSHKeyCipherFromKey(name string, privateKey string, keyring *crypto.Keyring) (models.Cipher, string, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return models.Cipher{}, "", err
	}

	pubKey := signer.PublicKey()
	encryptedName, _ := crypto.EncryptWith([]byte(name), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPublicKeyKey, _ := crypto.EncryptWith([]byte("public-key"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPublicKeyValue, _ := crypto.EncryptWith([]byte(string(ssh.MarshalAuthorizedKey(pubKey))), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedCustomTypeKey, _ := crypto.EncryptWith([]byte("custom-type"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedCustomTypeValue, _ := crypto.EncryptWith([]byte("ssh-key"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPrivateKeyKey, _ := crypto.EncryptWith([]byte("private-key"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPrivateKeyValue, _ := crypto.EncryptWith([]byte(privateKey), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())

	cipher := models.Cipher{
		Type:           models.CipherNote,
		Name:           encryptedName,
		Notes:          &encryptedPublicKeyValue,
		ID:             nil,
		Favorite:       false,
		OrganizationID: nil,
		SecureNote: &models.SecureNoteCipher{
			Type: 0,
		},
		Fields: []models.Field{
			{
				Type:  0,
				Name:  encryptedCustomTypeKey,
				Value: encryptedCustomTypeValue,
			},
			{
				Type:  0,
				Name:  encryptedPublicKeyKey,
				Value: encryptedPublicKeyValue,
			},
			{
				Type:  1,
				Name:  encryptedPrivateKeyKey,
				Value: encryptedPrivateKeyValue,
			},
		},
	}

	return cipher, string(ssh.MarshalAuthorizedKey(pubKey)), nil
}

func NewSSHKeyCipher(name string, keyring *crypto.Keyring) (models.Cipher, string, error) {
	var reader io.Reader = rand.Reader
	pub, priv, err := ed25519.GenerateKey(reader)

	if err != nil {
		panic(err)
	}
	privBlock := pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(priv),
	}

	privatePEM := pem.EncodeToMemory(&privBlock)
	publicKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		log.Error("Generation of public key failed: %s", err)
	}
	_, err = ssh.ParsePrivateKey([]byte(string(privatePEM)))
	if err != nil {
		log.Error("Verification of generated private key failed: %s", err)
	}

	encryptedName, _ := crypto.EncryptWith([]byte(name), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPublicKeyKey, _ := crypto.EncryptWith([]byte("public-key"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPublicKeyValue, _ := crypto.EncryptWith([]byte(string(ssh.MarshalAuthorizedKey(publicKey))), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedCustomTypeKey, _ := crypto.EncryptWith([]byte("custom-type"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedCustomTypeValue, _ := crypto.EncryptWith([]byte("ssh-key"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPrivateKeyKey, _ := crypto.EncryptWith([]byte("private-key"), crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())
	encryptedPrivateKeyValue, _ := crypto.EncryptWith(privatePEM, crypto.AesCbc256_HmacSha256_B64, keyring.GetAccountKey())

	cipher := models.Cipher{
		Type:           models.CipherNote,
		Name:           encryptedName,
		Notes:          &encryptedPublicKeyValue,
		ID:             nil,
		Favorite:       false,
		OrganizationID: nil,
		SecureNote: &models.SecureNoteCipher{
			Type: 0,
		},
		Fields: []models.Field{
			{
				Type:  0,
				Name:  encryptedCustomTypeKey,
				Value: encryptedCustomTypeValue,
			},
			{
				Type:  0,
				Name:  encryptedPublicKeyKey,
				Value: encryptedPublicKeyValue,
			},
			{
				Type:  1,
				Name:  encryptedPrivateKeyKey,
				Value: encryptedPrivateKeyValue,
			},
		},
	}

	return cipher, string(ssh.MarshalAuthorizedKey(publicKey)), nil
}
