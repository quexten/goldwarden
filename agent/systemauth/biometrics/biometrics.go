package biometrics

import (
	"os"

	"github.com/quexten/goldwarden/logging"
)

var log = logging.GetLogger("Goldwarden", "Biometrics")

var biometricsDisabled = false

func init() {
	if os.Getenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED") == "true" {
		biometricsDisabled = true
	}
}

type Approval string

const (
	AccessCredential  Approval = "com.quexten.goldwarden.accesscredential"
	ChangePin         Approval = "com.quexten.goldwarden.changepin"
	SSHKey            Approval = "com.quexten.goldwarden.usesshkey"
	ModifyVault       Approval = "com.quexten.goldwarden.modifyvault"
	BrowserBiometrics Approval = "com.quexten.goldwarden.browserbiometrics"
)

func (a Approval) String() string {
	return string(a)
}
