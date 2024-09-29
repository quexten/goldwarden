package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/quexten/goldwarden/cli/agent/bitwarden/crypto"
)

type SyncData struct {
	Profile Profile  `json:"profile"`
	Folders []Folder `json:"folders"`
	Ciphers []Cipher `json:"ciphers"`
}

type Organization struct {
	Object          string    `json:"object"`
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	UseGroups       bool      `json:"useGroups"`
	UseDirectory    bool      `json:"useDirectory"`
	UseEvents       bool      `json:"useEvents"`
	UseTotp         bool      `json:"useTotp"`
	Use2fa          bool      `json:"use2fa"`
	UseApi          bool      `json:"useApi"`
	UsersGetPremium bool      `json:"usersGetPremium"`
	SelfHost        bool      `json:"selfHost"`
	Seats           int       `json:"seats"`
	MaxCollections  int       `json:"maxCollections"`
	MaxStorageGb    int       `json:"maxStorageGb"`
	Key             string    `json:"key"`
	Status          int       `json:"status"`
	Type            int       `json:"type"`
	Enabled         bool      `json:"enabled"`
}

type Profile struct {
	ID                 uuid.UUID        `json:"id"`
	Name               string           `json:"name"`
	Email              string           `json:"email"`
	EmailVerified      bool             `json:"emailVerified"`
	Premium            bool             `json:"premium"`
	MasterPasswordHint string           `json:"masterPasswordHint"`
	Culture            string           `json:"culture"`
	TwoFactorEnabled   bool             `json:"twoFactorEnabled"`
	Key                crypto.EncString `json:"key"`
	PrivateKey         crypto.EncString `json:"privateKey"`
	SecurityStamp      string           `json:"securityStamp"`
	Organizations      []Organization   `json:"organizations"`
}

type Folder struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RevisionDate time.Time `json:"revisionDate"`
}

type Cipher struct {
	Type         CipherType       `json:"type,omitempty"`
	ID           *uuid.UUID       `json:"id,omitempty"`
	Name         crypto.EncString `json:"name,omitempty"`
	Edit         bool             `json:"edit,omitempty"`
	RevisionDate time.Time        `json:"revisionDate,omitempty"`
	DeletedDate  time.Time        `json:"deletedDate,omitempty"`

	FolderID            *uuid.UUID  `json:"folderId,omitempty"`
	OrganizationID      *uuid.UUID  `json:"organizationId,omitempty"`
	Favorite            bool        `json:"favorite,omitempty"`
	Attachments         interface{} `json:"attachments,omitempty"`
	OrganizationUseTotp bool        `json:"organizationUseTotp,omitempty"`
	CollectionIDs       []string    `json:"collectionIds,omitempty"`
	Fields              []Field     `json:"fields,omitempty"`

	Card       *Card             `json:"card,omitempty"`
	Identity   *Identity         `json:"identity,omitempty"`
	Login      *LoginCipher      `json:"login,omitempty"`
	Notes      *crypto.EncString `json:"notes,omitempty"`
	SecureNote *SecureNoteCipher `json:"secureNote,omitempty"`
  
  Key *crypto.EncString `json:"key,omitempty"`
}

type CipherType int

const (
	_              CipherType = iota
	CipherLogin               = 1
	CipherCard                = 3
	CipherIdentity            = 4
	CipherNote                = 2
)

type Card struct {
	CardholderName crypto.EncString `json:"cardholderName"`
	Brand          crypto.EncString `json:"brand"`
	Number         crypto.EncString `json:"number"`
	ExpMonth       crypto.EncString `json:"expMonth"`
	ExpYear        crypto.EncString `json:"expYear"`
	Code           crypto.EncString `json:"code"`
}

type Identity struct {
	Title      crypto.EncString `json:"title"`
	FirstName  crypto.EncString `json:"firstName"`
	MiddleName crypto.EncString `json:"middleName"`
	LastName   crypto.EncString `json:"lastName"`

	Username       crypto.EncString `json:"username"`
	Company        crypto.EncString `json:"company"`
	SSN            crypto.EncString `json:"ssn"`
	PassportNumber crypto.EncString `json:"passportNumber"`
	LicenseNumber  crypto.EncString `json:"licenseNumber"`

	Email      crypto.EncString `json:"email"`
	Phone      crypto.EncString `json:"phone"`
	Address1   crypto.EncString `json:"address1"`
	Address2   crypto.EncString `json:"address2"`
	Address3   crypto.EncString `json:"address3"`
	City       crypto.EncString `json:"city"`
	State      crypto.EncString `json:"state"`
	PostalCode crypto.EncString `json:"postalCode"`
	Country    crypto.EncString `json:"country"`
}

type FieldType int
type Field struct {
	Type  FieldType        `json:"type,omitempty"`
	Name  crypto.EncString `json:"name,omitempty"`
	Value crypto.EncString `json:"value,omitempty"`
}

type LoginCipher struct {
	Password crypto.EncString `json:"password,omitempty"`
	URI      crypto.EncString `json:"uri,omitempty"`
	URIs     []URI            `json:"uris,omitempty"`
	Username crypto.EncString `json:"username,omitempty"`
	Totp     crypto.EncString `json:"totp,omitempty"`
}

type URIMatch int
type URI struct {
	URI   string   `json:"uri,omitempty"`
	Match URIMatch `json:"match,omitempty"`
}

type SecureNoteType int
type SecureNoteCipher struct {
	Type SecureNoteType
}

func (cipher Cipher) GetKeyForCipher(keyring crypto.Keyring) (crypto.SymmetricEncryptionKey, error) {
	var key1 crypto.SymmetricEncryptionKey = nil
	var err error
	if cipher.OrganizationID != nil {
		key1, err = keyring.GetSymmetricKeyForOrganization(cipher.OrganizationID.String())
	} else {
		key1, err = keyring.GetAccountKey(), nil
	}

	if err != nil {
		return nil, err
	}

	if cipher.Key == nil {
		return key1, nil
	} else {
		key, err := crypto.DecryptWith(*cipher.Key, key1)
		if err != nil {
			return nil, err
		} else {
			return crypto.MemorySymmetricEncryptionKeyFromBytes(key)
		}
	}
}
