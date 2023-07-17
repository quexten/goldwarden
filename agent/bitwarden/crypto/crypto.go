package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math"

	"github.com/awnumar/memguard"
)

var b64enc = base64.StdEncoding.Strict()

type SymmetricEncryptionKey struct {
	encKey *memguard.Enclave
	macKey *memguard.Enclave
}

type AsymmetricEncryptionKey struct {
	encKey *memguard.Enclave
}

func SymmetricEncryptionKeyFromBytes(key []byte) (SymmetricEncryptionKey, error) {
	if len(key) != 64 {
		memguard.WipeBytes(key)
		return SymmetricEncryptionKey{}, fmt.Errorf("invalid key length: %d", len(key))
	}
	return SymmetricEncryptionKey{memguard.NewEnclave(key[0:32]), memguard.NewEnclave(key[32:64])}, nil
}

func (key SymmetricEncryptionKey) Bytes() []byte {
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

func AssymmetricEncryptionKeyFromBytes(key []byte) (AsymmetricEncryptionKey, error) {
	k := memguard.NewEnclave(key)
	return AsymmetricEncryptionKey{k}, nil
}

func isMacValid(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func encryptAESCBC256(data, key []byte) (iv, ciphertext []byte, _ error) {
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
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)
	return iv, ciphertext, nil
}

func decryptAESCBC256(iv, ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("iv length does not match AES block size")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of AES block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext) // decrypt in-place
	data, err := unpadPKCS7(ciphertext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func unpadPKCS7(src []byte, size int) ([]byte, error) {
	n := src[len(src)-1]
	if len(src)%size != 0 {
		return nil, fmt.Errorf("expected PKCS7 padding for block size %d, but have %d bytes", size, len(src))
	}
	if len(src) <= int(n) {
		return nil, fmt.Errorf("cannot unpad %d bytes out of a total of %d", n, len(src))
	}
	src = src[:len(src)-int(n)]
	return src, nil
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
