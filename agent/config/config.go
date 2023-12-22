package config

import (
	"bytes"
	cryptoSubtle "crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"runtime/debug"
	"sync"

	"github.com/google/uuid"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/tink-crypto/tink-go/v2/aead/subtle"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

const (
	KDFIterations     = 2
	KDFMemory         = 2 * 1024 * 1024
	KDFThreads        = 8
	DefaultConfigPath = "~/.config/goldwarden.json"
)

type RuntimeConfig struct {
	DisableAuth           bool
	DisablePinRequirement bool
	AuthMethod            string
	DoNotPersistConfig    bool
	ConfigDirectory       string
	DisableSSHAgent       bool
	WebsocketDisabled     bool
	ApiURI                string
	IdentityURI           string
	NotificationsURI      string
	SingleProcess         bool
	DeviceUUID            string
	User                  string
	Password              string
	Pin                   string
	UseMemguard           bool
}

type ConfigFile struct {
	IdentityUrl                 string
	ApiUrl                      string
	NotificationsUrl            string
	DeviceUUID                  string
	ConfigKeyHash               string
	EncryptedToken              string
	EncryptedUserSymmetricKey   string
	EncryptedMasterPasswordHash string
	EncryptedMasterKey          string
	RuntimeConfig               RuntimeConfig `json:"-"`
}

type LoginToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Key          string `json:"key"`
}

type Config struct {
	useMemguard bool
	key         *LockedBuffer
	ConfigFile  ConfigFile
	mu          sync.Mutex
}

func DefaultConfig(useMemguard bool) Config {
	deviceUUID, _ := uuid.NewUUID()
	keyBuffer := NewBuffer(32, useMemguard)
	return Config{
		useMemguard,
		&keyBuffer,
		ConfigFile{
			IdentityUrl:                 "https://vault.bitwarden.com/identity",
			ApiUrl:                      "https://vault.bitwarden.com/api",
			NotificationsUrl:            "https://notifications.bitwarden.com",
			DeviceUUID:                  deviceUUID.String(),
			ConfigKeyHash:               "",
			EncryptedToken:              "",
			EncryptedUserSymmetricKey:   "",
			EncryptedMasterPasswordHash: "",
			EncryptedMasterKey:          "",
			RuntimeConfig:               RuntimeConfig{},
		},
		sync.Mutex{},
	}
}

func (c *Config) IsLocked() bool {
	key := (*c.key).Bytes()
	return bytes.Equal(key, make([]byte, 32)) && c.HasPin()
}

func (c *Config) IsLoggedIn() bool {
	return c.ConfigFile.EncryptedMasterPasswordHash != ""
}

func (c *Config) Unlock(password string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsLocked() {
		return true
	}

	key := argon2.Key([]byte(password), []byte(c.ConfigFile.DeviceUUID), KDFIterations, KDFMemory, KDFThreads, 32)
	debug.FreeOSMemory()
	keyHash := sha3.Sum256(key)
	configKeyHash := hex.EncodeToString(keyHash[:])
	if cryptoSubtle.ConstantTimeCompare([]byte(configKeyHash), []byte(c.ConfigFile.ConfigKeyHash)) != 1 {
		return false
	}

	keyBuffer := NewBufferFromBytes(key, c.useMemguard)
	c.key = &keyBuffer
	return true
}

func (c *Config) VerifyPin(password string) bool {
	key := argon2.Key([]byte(password), []byte(c.ConfigFile.DeviceUUID), KDFIterations, KDFMemory, KDFThreads, 32)
	debug.FreeOSMemory()
	keyHash := sha3.Sum256(key)
	configKeyHash := hex.EncodeToString(keyHash[:])
	if cryptoSubtle.ConstantTimeCompare([]byte(configKeyHash), []byte(c.ConfigFile.ConfigKeyHash)) != 1 {
		return false
	} else {
		return true
	}
}

func (c *Config) Lock() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.IsLocked() {
		return
	}
	(*c.key).Wipe()
}

func (c *Config) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ConfigFile.EncryptedMasterPasswordHash = ""
	c.ConfigFile.EncryptedToken = ""
	c.ConfigFile.EncryptedUserSymmetricKey = ""
	c.ConfigFile.ConfigKeyHash = ""
	c.ConfigFile.EncryptedMasterKey = ""
	key := NewBuffer(32, c.useMemguard)
	c.key = &key
}

func (c *Config) HasPin() bool {
	return c.ConfigFile.ConfigKeyHash != ""
}

func (c *Config) UpdatePin(password string, write bool) {
	c.mu.Lock()

	newKey := argon2.Key([]byte(password), []byte(c.ConfigFile.DeviceUUID), KDFIterations, KDFMemory, KDFThreads, 32)
	keyHash := sha3.Sum256(newKey)
	configKeyHash := hex.EncodeToString(keyHash[:])
	debug.FreeOSMemory()

	c.ConfigFile.ConfigKeyHash = configKeyHash

	plaintextToken, err1 := c.decryptString(c.ConfigFile.EncryptedToken)
	plaintextUserSymmetricKey, err3 := c.decryptString(c.ConfigFile.EncryptedUserSymmetricKey)
	plaintextEncryptedMasterPasswordHash, err4 := c.decryptString(c.ConfigFile.EncryptedMasterPasswordHash)
	plaintextMasterKey, err5 := c.decryptString(c.ConfigFile.EncryptedMasterKey)

	key := NewBufferFromBytes(newKey, c.useMemguard)
	c.key = &key

	if err1 == nil {
		c.ConfigFile.EncryptedToken, err1 = c.encryptString(plaintextToken)
	}
	if err3 == nil {
		c.ConfigFile.EncryptedUserSymmetricKey, err3 = c.encryptString(plaintextUserSymmetricKey)
	}
	if err4 == nil {
		c.ConfigFile.EncryptedMasterPasswordHash, err4 = c.encryptString(plaintextEncryptedMasterPasswordHash)
	}
	if err5 == nil {
		c.ConfigFile.EncryptedMasterKey, err5 = c.encryptString(plaintextMasterKey)
	}
	c.mu.Unlock()

	if write {
		c.WriteConfig()
	}
}

func (c *Config) GetToken() (LoginToken, error) {
	if c.IsLocked() {
		return LoginToken{}, errors.New("config is locked")
	}
	tokenJson, err := c.decryptString(c.ConfigFile.EncryptedToken)
	if err != nil {
		return LoginToken{}, err
	}

	var token LoginToken
	err = json.Unmarshal([]byte(tokenJson), &token)
	if err != nil {
		return LoginToken{}, err
	}
	return token, nil
}

func (c *Config) SetToken(token LoginToken) error {
	if c.IsLocked() {
		return errors.New("config is locked")
	}

	tokenJson, err := json.Marshal(token)
	encryptedToken, err := c.encryptString(string(tokenJson))
	if err != nil {
		return err
	}
	// c.mu.Lock()
	c.ConfigFile.EncryptedToken = encryptedToken
	// c.mu.Unlock()
	c.WriteConfig()
	return nil
}

func (c *Config) GetUserSymmetricKey() ([]byte, error) {
	if c.IsLocked() {
		return []byte{}, errors.New("config is locked")
	}
	decrypted, err := c.decryptString(c.ConfigFile.EncryptedUserSymmetricKey)
	if err != nil {
		return []byte{}, err
	}
	return []byte(decrypted), nil
}

func (c *Config) SetUserSymmetricKey(key []byte) error {
	if c.IsLocked() {
		return errors.New("config is locked")
	}
	encryptedKey, err := c.encryptString(string(key))
	if err != nil {
		return err
	}
	// c.mu.Lock()
	c.ConfigFile.EncryptedUserSymmetricKey = encryptedKey
	// c.mu.Unlock()
	c.WriteConfig()
	return nil
}

func (c *Config) GetMasterPasswordHash() ([]byte, error) {
	if c.IsLocked() {
		return []byte{}, errors.New("config is locked")
	}
	decrypted, err := c.decryptString(c.ConfigFile.EncryptedMasterPasswordHash)
	if err != nil {
		return []byte{}, err
	}
	return []byte(decrypted), nil
}

func (c *Config) SetMasterPasswordHash(hash []byte) error {

	if c.IsLocked() {
		return errors.New("config is locked")
	}
	encryptedHash, err := c.encryptString(string(hash))
	if err != nil {
		c.mu.Unlock()
		return err
	}

	// c.mu.Lock()
	c.ConfigFile.EncryptedMasterPasswordHash = encryptedHash
	// c.mu.Unlock()

	c.WriteConfig()
	return nil
}

func (c *Config) GetMasterKey() ([]byte, error) {
	if c.IsLocked() {
		return []byte{}, errors.New("config is locked")
	}
	decrypted, err := c.decryptString(c.ConfigFile.EncryptedMasterKey)
	if err != nil {
		return []byte{}, err
	}
	return []byte(decrypted), nil
}

func (c *Config) SetMasterKey(key []byte) error {
	if c.IsLocked() {
		return errors.New("config is locked")
	}
	encryptedKey, err := c.encryptString(string(key))
	if err != nil {
		return err
	}
	// c.mu.Lock()
	c.ConfigFile.EncryptedMasterKey = encryptedKey
	// c.mu.Unlock()
	c.WriteConfig()
	return nil
}

func (c *Config) encryptString(data string) (string, error) {
	if c.IsLocked() {
		return "", errors.New("config is locked")
	}
	ca, err := subtle.NewChaCha20Poly1305((*c.key).Bytes())
	if err != nil {
		return "", err
	}
	result, err := ca.Encrypt([]byte(data), []byte{})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

func (c *Config) decryptString(data string) (string, error) {
	if c.IsLocked() {
		return "", errors.New("config is locked")
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	ca, err := subtle.NewChaCha20Poly1305((*c.key).Bytes())
	if err != nil {
		return "", err
	}
	result, err := ca.Decrypt(decoded, []byte{})
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (config *Config) WriteConfig() error {
	if config.ConfigFile.RuntimeConfig.DoNotPersistConfig {
		return nil
	}

	config.mu.Lock()
	defer config.mu.Unlock()

	jsonBytes, err := json.Marshal(config.ConfigFile)
	if err != nil {
		return err
	}

	// write to disk
	os.Remove(config.ConfigFile.RuntimeConfig.ConfigDirectory)
	file, err := os.OpenFile(config.ConfigFile.RuntimeConfig.ConfigDirectory, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfig(rtCfg RuntimeConfig) (Config, error) {
	file, err := os.Open(rtCfg.ConfigDirectory)
	if err != nil {
		key := NewBuffer(32, rtCfg.UseMemguard)
		return Config{
			key:        &key,
			ConfigFile: ConfigFile{},
		}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := ConfigFile{}
	err = decoder.Decode(&config)
	if err != nil {
		key := NewBuffer(32, rtCfg.UseMemguard)
		return Config{
			key:        &key,
			ConfigFile: ConfigFile{},
		}, err
	}
	if config.ConfigKeyHash == "" {
		key := NewBuffer(32, rtCfg.UseMemguard)
		return Config{
			key:        &key,
			ConfigFile: config,
		}, nil
	}
	key := NewBuffer(32, rtCfg.UseMemguard)
	return Config{
		key:        &key,
		ConfigFile: config,
	}, nil
}

func (cfg *Config) TryUnlock(vault *vault.Vault) error {
	pin, err := pinentry.GetPassword("Unlock Goldwarden", "Enter the vault PIN")
	if err != nil {
		return err
	}
	success := cfg.Unlock(pin)
	if !success {
		return errors.New("invalid PIN")
	}

	if cfg.IsLoggedIn() {
		userKey, err := cfg.GetUserSymmetricKey()
		if err == nil {
			var key crypto.SymmetricEncryptionKey
			var err error
			if vault.Keyring.IsMemguard {
				key, err = crypto.MemguardSymmetricEncryptionKeyFromBytes(userKey)
			} else {
				key, err = crypto.MemorySymmetricEncryptionKeyFromBytes(userKey)
			}
			if err != nil {
				return err
			}
			vault.Keyring.AccountKey = key
		} else {
			cfg.Lock()
			return err
		}
	}

	return nil
}
