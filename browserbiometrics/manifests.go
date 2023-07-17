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
