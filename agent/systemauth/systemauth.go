package systemauth

import (
	"os"

	"github.com/quexten/goldwarden/logging"
)

var log = logging.GetLogger("Goldwarden", "Systemauth")

var systemAuthDisabled = false

func init() {
	if os.Getenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED") == "true" {
		systemAuthDisabled = true
	}
}
