package browserbiometrics

const templateMozilla = `{
    "name": "com.8bit.bitwarden",
    "description": "Bitwarden desktop <-> browser bridge",
    "path": "PATH",
    "type": "stdio",
    "allowed_extensions": [
      "{446900e4-71c2-419f-a6a7-df9c091e268b}"
    ]
}`

const templateChrome = `{
	"name": "com.8bit.bitwarden",
	"description": "Bitwarden desktop <-> browser bridge",
	"path": "PATH",
	"type": "stdio",
	"allowed_origins": [
	  "chrome-extension://nngceckbapebfimnlniiiahkandclblb/",
	  "chrome-extension://jbkfoedolllekgbhcbcoahefnbanhhlh/",
	  "chrome-extension://ccnckbpmaceehanjmeomladnmlffdjgn/"
	]
  }`

const proxyScript = `#!/usr/bin/env bash

# Check if the "com.quexten.Goldwarden" Flatpak is installed
if flatpak list | grep -q "com.quexten.Goldwarden"; then
  flatpak run --command=goldwarden com.quexten.Goldwarden "$@"
else
  # If not installed, attempt to run the local version
  goldwarden "$@"
fi
`
