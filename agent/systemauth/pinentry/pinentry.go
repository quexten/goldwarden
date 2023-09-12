package pinentry

import (
	"os"

	"github.com/quexten/goldwarden/logging"
)

var log = logging.GetLogger("Goldwarden", "Pinentry")
var systemAuthDisabled = false

func init() {
	if os.Getenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED") == "true" {
		systemAuthDisabled = true
	}
}
