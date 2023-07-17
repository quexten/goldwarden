package main

import (
	"os"
	"strings"

	"github.com/quexten/goldwarden/browserbiometrics"
	"github.com/quexten/goldwarden/cmd"
)

func main() {
	if len(os.Args) > 1 && strings.Contains(os.Args[1], "com.8bit.bitwarden.json") {
		browserbiometrics.Main()
		return
	}

	cmd.Execute()
}
