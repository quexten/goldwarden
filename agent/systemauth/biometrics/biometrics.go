package biometrics

import (
	"github.com/quexten/goldwarden/logging"
)

var log = logging.GetLogger("Goldwarden", "Biometrics")

type Approval string

const (
	AccessVault       Approval = "com.quexten.goldwarden.accessvault"
	SSHKey            Approval = "com.quexten.goldwarden.usesshkey"
	BrowserBiometrics Approval = "com.quexten.goldwarden.browserbiometrics"
)

func (a Approval) String() string {
	return string(a)
}
