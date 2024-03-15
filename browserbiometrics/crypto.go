package browserbiometrics

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidBlockSize    = errors.New("invalid blocksize")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data (empty or not padded)")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

func decryptStringSymmetric(key []byte, ivb64 string, data string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	iv, _ := base64.StdEncoding.DecodeString(ivb64)
	ciphertext, _ := base64.StdEncoding.DecodeString(data)
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	ciphertext, _ = pkcs7Unpad(ciphertext, aes.BlockSize)

	return string(ciphertext), nil
}

func encryptStringSymmetric(key []byte, data []byte) (EncryptedString, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptedString{}, err
	}
	data, _ = pkcs7Pad(data, block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return EncryptedString{}, err
	}
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return EncryptedString{
		IV:      base64.StdEncoding.EncodeToString(ciphertext[:aes.BlockSize]),
		Data:    base64.StdEncoding.EncodeToString(ciphertext[aes.BlockSize:]),
		EncType: 0,
	}, nil
}

func generateTransportKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func rsaEncrypt(keyB64 string, message []byte) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}

	test, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	rsatest := test.(*rsa.PublicKey)
	oaepDigest := sha1.New()
	ciphertext, _ := rsa.EncryptOAEP(oaepDigest, rand.Reader, rsatest, message, []byte{})
	b64cipherText := base64.StdEncoding.EncodeToString(ciphertext)

	return b64cipherText, nil
}
