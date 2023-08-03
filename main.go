package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/quexten/goldwarden/browserbiometrics"
	"github.com/quexten/goldwarden/cmd"
)

func main() {
	if len(os.Args) > 1 && strings.Contains(os.Args[1], "com.8bit.bitwarden.json") {
		browserbiometrics.Main()
		return
	}

	if !cmd.IsPolkitSetup() {
		fmt.Println("Polkit is not setup. Run 'goldwarden setup polkit' to set it up.")
		time.Sleep(3 * time.Second)
	}

	cmd.Execute()
}
