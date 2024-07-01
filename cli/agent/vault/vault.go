package vault

import (
	"errors"
	"strings"
	"sync"

	"github.com/quexten/goldwarden/cli/agent/bitwarden/crypto"
	"github.com/quexten/goldwarden/cli/agent/bitwarden/models"
	"github.com/quexten/goldwarden/cli/logging"
	"golang.org/x/exp/slices"
)

var vaultLog = logging.GetLogger("Goldwarden", "Vault")

type Vault struct {
	Keyring            *crypto.Keyring
	logins             map[string]models.Cipher
	secureNotes        map[string]models.Cipher
	sshKeyNoteIDs      []string
	envCredentials     map[string]string
	lastSynced         int64
	websocketConnected bool
	mu                 sync.Mutex
}

func NewVault(keyring *crypto.Keyring) *Vault {
	return &Vault{
		Keyring:            keyring,
		logins:             make(map[string]models.Cipher),
		secureNotes:        make(map[string]models.Cipher),
		sshKeyNoteIDs:      make([]string, 0),
		envCredentials:     make(map[string]string),
		lastSynced:         0,
		websocketConnected: false,
	}
}

func (vault *Vault) lockMutex() {
	vault.mu.Lock()
}

func (vault *Vault) unlockMutex() {
	vault.mu.Unlock()
}

func (vault *Vault) Clear() {
	vault.lockMutex()
	vault.logins = make(map[string]models.Cipher)
	vault.secureNotes = make(map[string]models.Cipher)
	vault.sshKeyNoteIDs = make([]string, 0)
	vault.envCredentials = make(map[string]string)
	vault.lastSynced = 0
	vault.unlockMutex()
}

func (vault *Vault) AddOrUpdateLogin(cipher models.Cipher) {
	vault.lockMutex()
	vault.logins[cipher.ID.String()] = cipher
	vault.unlockMutex()
}

func (vault *Vault) DeleteCipher(uuid string) {
	vault.lockMutex()
	delete(vault.logins, uuid)
	delete(vault.envCredentials, uuid)

	newSecureNotes := make(map[string]models.Cipher)
	for _, noteID := range vault.sshKeyNoteIDs {
		if noteID != uuid {
			newSecureNotes[noteID] = vault.secureNotes[noteID]
		}
	}
	vault.secureNotes = newSecureNotes
	vault.unlockMutex()
}

func (vault *Vault) AddOrUpdateSecureNote(cipher models.Cipher) {
	vault.lockMutex()
	vault.secureNotes[cipher.ID.String()] = cipher

	if vault.isSSHKey(cipher) {
		if !slices.Contains(vault.sshKeyNoteIDs, cipher.ID.String()) {
			vault.sshKeyNoteIDs = append(vault.sshKeyNoteIDs, cipher.ID.String())
		}
	} else if executableName, isEnv := vault.isEnv(cipher); isEnv {
		vault.envCredentials[executableName] = cipher.ID.String()
	}

	vault.unlockMutex()
}

func (vault *Vault) isEnv(cipher models.Cipher) (string, bool) {
	if cipher.Type != models.CipherNote {
		return "", false
	}

	if !cipher.DeletedDate.IsZero() {
		return "", false
	}

	key, err := cipher.GetKeyForCipher(*vault.Keyring)
	if err != nil {
		vaultLog.Error("Failed to get key for cipher "+cipher.ID.String(), err.Error())
		return "", false
	}

	isEnv := false
	executableName := ""

	for _, field := range cipher.Fields {
		fieldName, err := crypto.DecryptWith(field.Name, key)
		if err != nil {
			continue
		}
		fieldValue, err := crypto.DecryptWith(field.Value, key)
		if err != nil {
			continue
		}

		if string(fieldName) == "custom-type" && string(fieldValue) == "env" {
			isEnv = true
		} else if string(fieldName) == "executable" {
			executableName = string(fieldValue)
		}
	}

	return executableName, isEnv
}

func (vault *Vault) isSSHKey(cipher models.Cipher) bool {
	if cipher.Type != models.CipherNote {
		return false
	}

	if !cipher.DeletedDate.IsZero() {
		return false
	}

	key, err := cipher.GetKeyForCipher(*vault.Keyring)
	if err != nil {
		vaultLog.Error("Failed to get key for cipher "+cipher.ID.String(), err.Error())
		return false
	}

	for _, field := range cipher.Fields {
		fieldName, err := crypto.DecryptWith(field.Name, key)
		if err != nil {
			cipherID := cipher.ID.String()
			if cipher.OrganizationID != nil {
				orgID := cipher.OrganizationID.String()
				vaultLog.Error("Failed to decrypt field name with on cipher %s in organization %s: %s", cipherID, orgID, err.Error())
			} else {
				vaultLog.Error("Failed to decrypt field name with on cipher %s: %s", cipherID, err.Error())
			}
			continue
		}
		fieldValue, err := crypto.DecryptWith(field.Value, key)
		if err != nil {
			continue
		}

		if string(fieldName) == "custom-type" && string(fieldValue) == "ssh-key" {
			return true
		}
	}

	return false
}

type SSHKey struct {
	Name      string
	Key       string
	PublicKey string
}

func (vault *Vault) GetSSHKeys() []SSHKey {
	vault.lockMutex()
	defer vault.unlockMutex()

	var sshKeys []SSHKey
	for _, id := range vault.sshKeyNoteIDs {
		privateKey := ""
		publicKey := ""

		key, err := vault.secureNotes[id].GetKeyForCipher(*vault.Keyring)
		if err != nil {
			continue
		}

		for _, field := range vault.secureNotes[id].Fields {
			fieldName, err := crypto.DecryptWith(field.Name, key)
			if err != nil {
				continue
			}
			if string(fieldName) == "private-key" {
				pk, err := crypto.DecryptWith(field.Value, key)
				if err != nil {
					continue
				} else {
					privateKey = string(pk)
				}
			}
			if string(fieldName) == "public-key" {
				pk, err := crypto.DecryptWith(field.Value, key)
				if err != nil {
					continue
				} else {
					publicKey = string(pk)
				}
			}
		}

		privateKey = strings.Replace(privateKey, "-----BEGIN OPENSSH PRIVATE KEY-----", "", 1)
		privateKey = strings.Replace(privateKey, "-----END OPENSSH PRIVATE KEY-----", "", 1)

		pkParts := strings.Join(strings.Split(privateKey, " "), "\n")
		privateKeyString := "-----BEGIN OPENSSH PRIVATE KEY-----" + pkParts + "-----END OPENSSH PRIVATE KEY-----"

		decryptedTitle, err := crypto.DecryptWith(vault.secureNotes[id].Name, key)
		if err != nil {
			continue
		}

		sshKeys = append(sshKeys, SSHKey{
			Name:      string(decryptedTitle),
			Key:       string(privateKeyString),
			PublicKey: string(publicKey),
		})
	}
	return sshKeys
}

func (vault *Vault) GetEnvCredentialForExecutable(executableName string) (map[string]string, bool) {
	vault.lockMutex()
	defer vault.unlockMutex()

	env := make(map[string]string)

	if id, ok := vault.envCredentials[executableName]; ok {
		key, err := vault.secureNotes[id].GetKeyForCipher(*vault.Keyring)
		if err != nil {
			vaultLog.Error("Failed to get key for cipher " + id)
			return make(map[string]string), false
		}

		for _, field := range vault.secureNotes[id].Fields {
			fieldName, err := crypto.DecryptWith(field.Name, key)
			if err != nil {
				continue
			}
			fieldValue, err := crypto.DecryptWith(field.Value, key)
			if err != nil {
				continue
			}

			if string(fieldName) == "custom-type" || string(fieldName) == "executable" {
				continue
			}

			env[string(fieldName)] = string(fieldValue)
		}
		return env, true
	}
	return make(map[string]string), false
}

func (vault *Vault) GetLogins() []models.Cipher {
	vault.lockMutex()
	defer vault.unlockMutex()

	var logins []models.Cipher
	for _, cipher := range vault.logins {
		if cipher.Type != models.CipherLogin {
			continue
		}
		if !cipher.DeletedDate.IsZero() {
			continue
		}
		logins = append(logins, cipher)
	}
	return logins
}

func (vault *Vault) GetNotes() []models.Cipher {
	vault.lockMutex()
	defer vault.unlockMutex()

	var notes []models.Cipher
	for _, cipher := range vault.secureNotes {
		if cipher.Type != models.CipherNote {
			continue
		}
		if !cipher.DeletedDate.IsZero() {
			continue
		}
		notes = append(notes, cipher)
	}
	return notes
}

func (vault *Vault) GetLoginByFilter(uuid string, orgId string, name string, username string) (models.Cipher, error) {
	vault.lockMutex()
	defer vault.unlockMutex()

	for _, cipher := range vault.logins {
		if uuid != "" && cipher.ID.String() != uuid {
			continue
		}
		if orgId != "" && cipher.OrganizationID.String() != orgId {
			continue
		}

		key, err := cipher.GetKeyForCipher(*vault.Keyring)
		if err != nil {
			vaultLog.Error("Failed to get key for cipher " + cipher.ID.String())
			continue
		}
		if name != "" {
			if cipher.Name.IsNull() {
				continue
			}

			decryptedName, err := crypto.DecryptWith(cipher.Name, key)
			if err != nil {
				vaultLog.Error("Failed to decrypt name for cipher " + cipher.ID.String())
				continue
			}
			if name != "" && string(decryptedName) != name {
				continue
			}
		}

		if username != "" {
			if cipher.Login.Username.IsNull() {
				continue
			}

			decryptedUsername, err := crypto.DecryptWith(cipher.Login.Username, key)
			if err != nil {
				vaultLog.Error("Failed to decrypt username for cipher " + cipher.ID.String())
				continue
			}
			if username != "" && string(decryptedUsername) != username {
				continue
			}
		}

		return cipher, nil
	}

	return models.Cipher{}, errors.New("cipher not found")
}

func (vault *Vault) GetNoteByFilter(uuid string, orgId string, name string) (models.Cipher, error) {
	vault.lockMutex()
	defer vault.unlockMutex()

	for _, cipher := range vault.secureNotes {
		if uuid != "" && cipher.ID.String() != uuid {
			continue
		}
		if orgId != "" && cipher.OrganizationID.String() != orgId {
			continue
		}

		key, err := cipher.GetKeyForCipher(*vault.Keyring)
		if err != nil {
			vaultLog.Error("Failed to get key for cipher "+cipher.ID.String(), err.Error())
			continue
		}
		decryptedName, err := crypto.DecryptWith(cipher.Name, key)
		if err != nil {
			vaultLog.Error("Failed to decrypt name for cipher "+cipher.ID.String(), err.Error())
			continue
		}
		if name != "" && string(decryptedName) != name {
			continue
		}
		return cipher, nil
	}

	return models.Cipher{}, errors.New("cipher not found")
}

func (vault *Vault) GetLogin(uuid string) (models.Cipher, error) {
	vault.lockMutex()
	defer vault.unlockMutex()

	for _, cipher := range vault.logins {
		if cipher.ID.String() == uuid {
			return cipher, nil
		}
	}

	return models.Cipher{}, errors.New("cipher not found")
}

func (vault *Vault) GetSecureNote(uuid string) (models.Cipher, error) {
	vault.lockMutex()
	defer vault.unlockMutex()

	for _, cipher := range vault.secureNotes {
		if cipher.ID.String() == uuid {
			return cipher, nil
		}
	}

	return models.Cipher{}, errors.New("cipher not found")
}

func (vault *Vault) SetLastSynced(lastSynced int64) {
	vault.lockMutex()
	vault.lastSynced = lastSynced
	vault.unlockMutex()
}

func (vault *Vault) GetLastSynced() int64 {
	vault.lockMutex()
	defer vault.unlockMutex()

	return vault.lastSynced
}

func (vault *Vault) SetWebsocketConnected(connected bool) {
	vault.lockMutex()
	vault.websocketConnected = connected
	vault.unlockMutex()
}

func (vault *Vault) IsWebsocketConnected() bool {
	vault.lockMutex()
	defer vault.unlockMutex()

	return vault.websocketConnected
}
