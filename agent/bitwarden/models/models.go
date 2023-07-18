package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/quexten/goldwarden/agent/bitwarden/crypto"
)

type SyncData struct {
	Profile Profile
	Folders []Folder
	Ciphers []Cipher
}

type Organization struct {
	Object          string
	Id              uuid.UUID
	Name            string
	UseGroups       bool
	UseDirectory    bool
	UseEvents       bool
	UseTotp         bool
	Use2fa          bool
	UseApi          bool
	UsersGetPremium bool
	SelfHost        bool
	Seats           int
	MaxCollections  int
	MaxStorageGb    int
	Key             string
	Status          int
	Type            int
	Enabled         bool
}

type Profile struct {
	ID                 uuid.UUID
	Name               string
	Email              string
	EmailVerified      bool
	Premium            bool
	MasterPasswordHint string
	Culture            string
	TwoFactorEnabled   bool
	Key                crypto.EncString
	PrivateKey         crypto.EncString
	SecurityStamp      string
	Organizations      []Organization
}

type Folder struct {
	ID           uuid.UUID
	Name         string
	RevisionDate time.Time
}

type Cipher struct {
	Type         CipherType
	ID           *uuid.UUID `json:",omitempty"`
	Name         crypto.EncString
	Edit         bool
	RevisionDate time.Time
	DeletedDate  time.Time

	FolderID            *uuid.UUID  `json:",omitempty"`
	OrganizationID      *uuid.UUID  `json:",omitempty"`
	Favorite            bool        `json:",omitempty"`
	Attachments         interface{} `json:",omitempty"`
	OrganizationUseTotp bool        `json:",omitempty"`
	CollectionIDs       []string    `json:",omitempty"`
	Fields              []Field     `json:",omitempty"`

	Card       *Card             `json:",omitempty"`
	Identity   *Identity         `json:",omitempty"`
	Login      *LoginCipher      `json:",omitempty"`
	Notes      *crypto.EncString `json:",omitempty"`
	SecureNote *SecureNoteCipher `json:",omitempty"`
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
	CardholderName crypto.EncString
	Brand          crypto.EncString
	Number         crypto.EncString
	ExpMonth       crypto.EncString
	ExpYear        crypto.EncString
	Code           crypto.EncString
}

type Identity struct {
	Title      crypto.EncString
	FirstName  crypto.EncString
	MiddleName crypto.EncString
	LastName   crypto.EncString

	Username       crypto.EncString
	Company        crypto.EncString
	SSN            crypto.EncString
	PassportNumber crypto.EncString
	LicenseNumber  crypto.EncString

	Email      crypto.EncString
	Phone      crypto.EncString
	Address1   crypto.EncString
	Address2   crypto.EncString
	Address3   crypto.EncString
	City       crypto.EncString
	State      crypto.EncString
	PostalCode crypto.EncString
	Country    crypto.EncString
}

type FieldType int
type Field struct {
	Type  FieldType
	Name  crypto.EncString
	Value crypto.EncString
}

type LoginCipher struct {
	Password crypto.EncString
	URI      crypto.EncString
	URIs     []URI
	Username crypto.EncString `json:",omitempty"`
	Totp     crypto.EncString `json:",omitempty"`
}

type URIMatch int
type URI struct {
	URI   string
	Match URIMatch
}

type SecureNoteType int
type SecureNoteCipher struct {
	Type SecureNoteType
}

func (cipher Cipher) GetKeyForCipher(keyring crypto.Keyring) (crypto.SymmetricEncryptionKey, error) {
	if cipher.OrganizationID != nil {
		return keyring.GetSymmetricKeyForOrganization(cipher.OrganizationID.String())
	}
	return *keyring.AccountKey, nil
}
