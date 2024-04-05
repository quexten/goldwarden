package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"strconv"

	"github.com/awnumar/memguard"
	"github.com/quexten/goldwarden/logging"
)

var cryptoLog = logging.GetLogger("Goldwarden", "Crypto")

type EncString struct {
	Type        EncStringType
	IV, CT, MAC []byte
}

type EncStringType int

const (
	AesCbc256_B64                     EncStringType = 0
	AesCbc128_HmacSha256_B64          EncStringType = 1
	AesCbc256_HmacSha256_B64          EncStringType = 2
	Rsa2048_OaepSha256_B64            EncStringType = 3
	Rsa2048_OaepSha1_B64              EncStringType = 4
	Rsa2048_OaepSha256_HmacSha256_B64 EncStringType = 5
	Rsa2048_OaepSha1_HmacSha256_B64   EncStringType = 6
)

func (t EncStringType) HasMAC() bool {
	return t != AesCbc256_B64
}

func (s *EncString) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	i := bytes.IndexByte(data, '.')
	if i < 0 {
		return errors.New("invalid cipher string format, missign type. total length: " + strconv.Itoa(len(data)))
	}

	typStr := string(data[:i])
	var err error
	if t, err := strconv.Atoi(typStr); err != nil {
		return errors.New("invalid cipher string type, could not parse, length: " + strconv.Itoa(len(data)))
	} else {
		s.Type = EncStringType(t)
	}

	switch s.Type {
	case AesCbc128_HmacSha256_B64, AesCbc256_HmacSha256_B64, AesCbc256_B64:
	default:
		return errors.New("invalid cipher string type, unknown type: " + strconv.Itoa(int(s.Type)))
	}

	data = data[i+1:]
	parts := bytes.Split(data, []byte("|"))
	if len(parts) != 3 {
		return errors.New("invalid cipher string format, missing parts, length: " + strconv.Itoa(len(data)) + "type: " + strconv.Itoa(int(s.Type)))
	}

	if s.IV, err = b64decode(parts[0]); err != nil {
		return err
	}
	if s.CT, err = b64decode(parts[1]); err != nil {
		return err
	}
	if s.Type.HasMAC() {
		if s.MAC, err = b64decode(parts[2]); err != nil {
			return err
		}
	}
	return nil
}

func (s EncString) MarshalText() ([]byte, error) {
	if s.Type == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(int(s.Type)))
	buf.WriteByte('.')
	buf.Write(b64encode(s.IV))
	buf.WriteByte('|')
	buf.Write(b64encode(s.CT))
	if s.Type.HasMAC() {
		buf.WriteByte('|')
		buf.Write(b64encode(s.MAC))
	}
	return buf.Bytes(), nil
}

func (s EncString) IsNull() bool {
	return len(s.IV) == 0 && len(s.CT) == 0 && len(s.MAC) == 0
}

func b64decode(src []byte) ([]byte, error) {
	dst := make([]byte, b64enc.DecodedLen(len(src)))
	n, err := b64enc.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	dst = dst[:n]
	return dst, nil
}

func b64encode(src []byte) []byte {
	dst := make([]byte, b64enc.EncodedLen(len(src)))
	b64enc.Encode(dst, src)
	return dst
}

func DecryptWith(s EncString, key SymmetricEncryptionKey) ([]byte, error) {
	encKeyData, err := key.EncryptionKeyBytes()
	if err != nil {
		return nil, err
	}
	macKeyData, err := key.MacKeyBytes()
	if err != nil {
		return nil, err
	}

	switch s.Type {
	case AesCbc256_B64, AesCbc256_HmacSha256_B64:
		break
	default:
		return nil, fmt.Errorf("decrypt: unsupported cipher type %q", s.Type)
	}

	if s.Type == AesCbc256_HmacSha256_B64 {
		if len(s.MAC) == 0 || len(macKeyData) == 0 {
			return nil, fmt.Errorf("decrypt: cipher string type expects a MAC")
		}
		var msg []byte
		msg = append(msg, s.IV...)
		msg = append(msg, s.CT...)
		if !isMacValid(msg, s.MAC, macKeyData) {
			return nil, fmt.Errorf("decrypt: MAC mismatch")
		}
	} else if s.Type == AesCbc256_B64 {
		return nil, fmt.Errorf("decrypt: cipher of unsupported type %q", s.Type)
	}

	dst, err := decryptAESCBC256(s.IV, s.CT, encKeyData)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func EncryptWith(data []byte, encType EncStringType, key SymmetricEncryptionKey) (EncString, error) {
	encKeyData, err := key.EncryptionKeyBytes()
	if err != nil {
		return EncString{}, err
	}
	macKeyData, err := key.MacKeyBytes()
	if err != nil {
		return EncString{}, err
	}

	s := EncString{}
	switch encType {
	case AesCbc256_B64, AesCbc256_HmacSha256_B64:
	default:
		return s, fmt.Errorf("encrypt: unsupported cipher type %q", s.Type)
	}
	s.Type = encType

	iv, ciphertext, err := encryptAESCBC256(data, encKeyData)
	if err != nil {
		return s, err
	}
	s.CT = ciphertext
	s.IV = iv

	if encType == AesCbc256_HmacSha256_B64 {
		if len(macKeyData) == 0 {
			return s, fmt.Errorf("encrypt: cipher string type expects a MAC")
		}
		var macMessage []byte
		macMessage = append(macMessage, s.IV...)
		macMessage = append(macMessage, s.CT...)
		mac := hmac.New(sha256.New, macKeyData)
		mac.Write(macMessage)
		s.MAC = mac.Sum(nil)
	}

	return s, nil
}

func EncryptWithToString(data []byte, encType EncStringType, key SymmetricEncryptionKey) (string, error) {
	s, err := EncryptWith(data, encType, key)
	if err != nil {
		return "", err
	}

	marshalled, err := s.MarshalText()
	if err != nil {
		return "", err
	}

	return string(marshalled), nil
}

func GenerateAsymmetric(useMemguard bool) (AsymmetricEncryptionKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return MemoryAsymmetricEncryptionKey{}, err
	}

	encKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return MemoryAsymmetricEncryptionKey{}, err
	}

	if useMemguard {
		return MemguardAsymmetricEncryptionKey{memguard.NewEnclave(encKey)}, nil
	} else {
		return MemoryAsymmetricEncryptionKey{encKey}, nil
	}
}

func DecryptWithAsymmetric(s []byte, asymmetrickey AsymmetricEncryptionKey) ([]byte, error) {
	key, err := asymmetrickey.PrivateBytes()
	if err != nil {
		return nil, err
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	rawKey, err := b64decode(s[2:])
	if err != nil {
		return nil, err
	}

	res, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, parsedKey.(*rsa.PrivateKey), rawKey, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func EncryptWithAsymmetric(s []byte, asymmbetrickey AsymmetricEncryptionKey) ([]byte, error) {
	key, err := asymmbetrickey.PrivateBytes()
	if err != nil {
		return nil, err
	}

	parsedKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	res, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, parsedKey.(*rsa.PublicKey), s, nil)
	if err != nil {
		return nil, err
	}

	resB64 := b64encode(res)
	res = append([]byte("4."), resB64...)

	return res, nil
}
