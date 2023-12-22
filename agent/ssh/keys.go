package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/mikesmitty/edkey"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/bitwarden/models"
	"golang.org/x/crypto/ssh"
)

func NewSSHKeyCipher(name string, keyring *crypto.Keyring) (models.Cipher, string) {

	var reader io.Reader = rand.Reader
	pub, priv, err := ed25519.GenerateKey(reader)

	if err != nil {
		panic(err)
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(priv)
	privBlock := pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privateKey),
	}

	privatePEM := pem.EncodeToMemory(&privBlock)
	publicKey, err := ssh.NewPublicKey(pub)

	encryptedName, _ := crypto.EncryptWith([]byte(name), crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)
	encryptedPublicKeyKey, _ := crypto.EncryptWith([]byte("public-key"), crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)
	encryptedPublicKeyValue, _ := crypto.EncryptWith([]byte(string(ssh.MarshalAuthorizedKey(publicKey))), crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)
	encryptedCustomTypeKey, _ := crypto.EncryptWith([]byte("custom-type"), crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)
	encryptedCustomTypeValue, _ := crypto.EncryptWith([]byte("ssh-key"), crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)
	encryptedPrivateKeyKey, _ := crypto.EncryptWith([]byte("private-key"), crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)
	encryptedPrivateKeyValue, _ := crypto.EncryptWith(privatePEM, crypto.AesCbc256_HmacSha256_B64, keyring.AccountKey)

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

	return cipher, string(ssh.MarshalAuthorizedKey(publicKey))
}
