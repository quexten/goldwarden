package logging

import (
	"os"

	"github.com/LlamaNite/llamalog"
)

type SilentWriter struct {
}

func (w SilentWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func GetLogger(packageDir ...string) *llamalog.Logger {
	if os.Getenv("GOLDWARDEN_SILENT_LOGGING") == "true" {
		return llamalog.NewLoggerFromWriter(SilentWriter{}, packageDir...)
	} else {
		return llamalog.NewLogger(packageDir...)
	}
}
