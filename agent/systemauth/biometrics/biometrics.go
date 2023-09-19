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
	AccessVault       Approval = "com.quexten.goldwarden.accessvault"
	SSHKey            Approval = "com.quexten.goldwarden.usesshkey"
	BrowserBiometrics Approval = "com.quexten.goldwarden.browserbiometrics"
)

func (a Approval) String() string {
	return string(a)
}
