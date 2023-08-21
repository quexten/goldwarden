package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/awnumar/memguard"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

type KDFType int

const (
	PBKDF2   KDFType = 0
	Argon2ID KDFType = 1
)

type KDFConfig struct {
	Type        KDFType
	Iterations  uint32
	Memory      uint32
	Parallelism uint32
}

type MasterKey struct {
	encKey *memguard.Enclave
}

func (masterKey MasterKey) GetBytes() []byte {
	defer debug.FreeOSMemory()

	buffer, err := masterKey.encKey.Open()
	if err != nil {
		panic(err)
	}
	defer buffer.Destroy()

	return bytes.Clone(buffer.Bytes())
}

func DeriveMasterKey(password memguard.LockedBuffer, email string, kdfConfig KDFConfig) (MasterKey, error) {
	defer debug.FreeOSMemory()

	var key []byte
	switch kdfConfig.Type {
	case PBKDF2:
		key = pbkdf2.Key(password.Bytes(), []byte(strings.ToLower(email)), int(kdfConfig.Iterations), 32, sha256.New)
	case Argon2ID:
		var salt [32]byte = sha256.Sum256([]byte(strings.ToLower(email)))
		key = argon2.IDKey(password.Bytes(), salt[:], kdfConfig.Iterations, kdfConfig.Memory*1024, uint8(kdfConfig.Parallelism), 32)
	default:
		password.Destroy()
		return MasterKey{}, fmt.Errorf("unsupported KDF type %d", kdfConfig.Type)
	}
	password.Destroy()

	return MasterKey{memguard.NewEnclave(key)}, nil
}

func MasterKeyFromBytes(key []byte) MasterKey {
	return MasterKey{memguard.NewEnclave(key)}
}
