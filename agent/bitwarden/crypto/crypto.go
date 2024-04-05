package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/awnumar/memguard"
)

var b64enc = base64.StdEncoding.Strict()

type SymmetricEncryptionKey interface {
	Bytes() []byte
	EncryptionKeyBytes() ([]byte, error)
	MacKeyBytes() ([]byte, error)
}

type MemorySymmetricEncryptionKey struct {
	encKey []byte
	macKey []byte
}

type MemguardSymmetricEncryptionKey struct {
	encKey *memguard.Enclave
	macKey *memguard.Enclave
}

type AsymmetricEncryptionKey interface {
	PublicBytes() []byte
	PrivateBytes() ([]byte, error)
}

type MemoryAsymmetricEncryptionKey struct {
	encKey []byte
}

type MemguardAsymmetricEncryptionKey struct {
	encKey *memguard.Enclave
}

func MemguardSymmetricEncryptionKeyFromBytes(key []byte) (MemguardSymmetricEncryptionKey, error) {
	if len(key) != 64 {
		memguard.WipeBytes(key)
		return MemguardSymmetricEncryptionKey{}, fmt.Errorf("invalid key length: %d", len(key))
	}
	return MemguardSymmetricEncryptionKey{memguard.NewEnclave(key[0:32]), memguard.NewEnclave(key[32:64])}, nil
}

func MemorySymmetricEncryptionKeyFromBytes(key []byte) (MemorySymmetricEncryptionKey, error) {
	if len(key) != 64 {
		return MemorySymmetricEncryptionKey{}, fmt.Errorf("invalid key length: %d", len(key))
	}
	return MemorySymmetricEncryptionKey{encKey: key[0:32], macKey: key[32:64]}, nil
}

func (key MemguardSymmetricEncryptionKey) Bytes() []byte {
	k1, err := key.encKey.Open()
	if err != nil {
		panic(err)
	}
	k2, err := key.macKey.Open()
	if err != nil {
		panic(err)
	}
	keyBytes := make([]byte, 64)
	copy(keyBytes[0:32], k1.Bytes())
	copy(keyBytes[32:64], k2.Bytes())
	return keyBytes
}

func (key MemorySymmetricEncryptionKey) Bytes() []byte {
	keyBytes := make([]byte, 64)
	copy(keyBytes[0:32], key.encKey)
	copy(keyBytes[32:64], key.macKey)
	return keyBytes
}

func (key MemorySymmetricEncryptionKey) EncryptionKeyBytes() ([]byte, error) {
	return key.encKey, nil
}

func (key MemguardSymmetricEncryptionKey) EncryptionKeyBytes() ([]byte, error) {
	k, err := key.encKey.Open()
	if err != nil {
		return nil, err
	}
	keyBytes := make([]byte, 32)
	copy(keyBytes, k.Bytes())
	return keyBytes, nil
}

func (key MemorySymmetricEncryptionKey) MacKeyBytes() ([]byte, error) {
	return key.macKey, nil
}

func (key MemguardSymmetricEncryptionKey) MacKeyBytes() ([]byte, error) {
	k, err := key.macKey.Open()
	if err != nil {
		return nil, err
	}
	keyBytes := make([]byte, 32)
	copy(keyBytes, k.Bytes())
	return keyBytes, nil
}

func MemoryAssymmetricEncryptionKeyFromBytes(key []byte) (MemoryAsymmetricEncryptionKey, error) {
	return MemoryAsymmetricEncryptionKey{key}, nil
}

func MemguardAssymmetricEncryptionKeyFromBytes(key []byte) (MemguardAsymmetricEncryptionKey, error) {
	k := memguard.NewEnclave(key)
	return MemguardAsymmetricEncryptionKey{k}, nil
}

func (key MemoryAsymmetricEncryptionKey) PublicBytes() []byte {
	privateKey, err := x509.ParsePKCS8PrivateKey(key.encKey)
	if err != nil {
		panic(err)
	}
	pub := (privateKey.(*rsa.PrivateKey)).Public()
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		panic(err)
	}
	return publicKeyBytes
}

func (key MemguardAsymmetricEncryptionKey) PublicBytes() []byte {
	buffer, err := key.encKey.Open()
	if err != nil {
		panic(err)
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(buffer.Bytes())
	if err != nil {
		panic(err)
	}
	pub := (privateKey.(*rsa.PrivateKey)).Public()
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		panic(err)
	}
	return publicKeyBytes
}

func (key MemoryAsymmetricEncryptionKey) PrivateBytes() ([]byte, error) {
	return key.encKey, nil
}

func (key MemguardAsymmetricEncryptionKey) PrivateBytes() ([]byte, error) {
	buffer, err := key.encKey.Open()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func isMacValid(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func encryptAESCBC256(data, key []byte) (iv, ciphertext []byte, err error) {
	data = padPKCS7(data, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	ivSize := aes.BlockSize
	iv = make([]byte, ivSize)
	ciphertext = make([]byte, len(data))
	if _, err := io.ReadFull(cryptorand.Reader, iv); err != nil {
		return nil, nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				iv = nil
				ciphertext = nil
				err = errors.New("error encrypting AES CBC 256 data")
			}
		}
	}()

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)
	return iv, ciphertext, nil
}

func decryptAESCBC256(iv, ciphertext, key []byte) (decryptedData []byte, err error) {
	ciphertextCopy := make([]byte, len(ciphertext))
	copy(ciphertextCopy, ciphertext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("iv length does not match AES block size")
	}
	if len(ciphertextCopy)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of AES block size")
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				decryptedData = nil
				err = errors.New("error decrypting AES CBC 256 data")
			}
		}
	}()

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextCopy, ciphertextCopy) // decrypt in-place
	data, err := unpadPKCS7(ciphertextCopy, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	resultBuffer := make([]byte, len(data))
	copy(resultBuffer, data)
	return resultBuffer, nil
}

func unpadPKCS7(src []byte, size int) ([]byte, error) {
	srcCopy := make([]byte, len(src))
	copy(srcCopy, src)

	n := srcCopy[len(srcCopy)-1]
	if len(srcCopy)%size != 0 {
		return nil, fmt.Errorf("expected PKCS7 padding for block size %d, but have %d bytes", size, len(srcCopy))
	}
	if len(srcCopy) <= int(n) {
		return nil, fmt.Errorf("cannot unpad %d bytes out of a total of %d", n, len(srcCopy))
	}
	srcCopy = srcCopy[:len(srcCopy)-int(n)]

	resultCopy := make([]byte, len(srcCopy))
	copy(resultCopy, srcCopy)

	return resultCopy, nil
}

func padPKCS7(src []byte, size int) []byte {
	rem := len(src) % size
	n := size - rem
	if n > math.MaxUint8 {
		panic(fmt.Sprintf("cannot pad over %d bytes, but got %d", math.MaxUint8, n))
	}
	padded := make([]byte, len(src)+n)
	copy(padded, src)
	for i := len(src); i < len(padded); i++ {
		padded[i] = byte(n)
	}
	return padded
}
