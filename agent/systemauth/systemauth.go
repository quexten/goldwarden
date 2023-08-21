package systemauth

import "os"

var systemAuthDisabled = false

func init() {
	if os.Getenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED") == "true" {
		systemAuthDisabled = true
	}
}
