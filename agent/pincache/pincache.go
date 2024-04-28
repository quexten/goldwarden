package pincache

import (
	"errors"

	"github.com/awnumar/memguard"
	"github.com/quexten/goldwarden/agent/systemauth/biometrics"
)

var cachedPin *memguard.Enclave

func SetPin(useMemguard bool, pin []byte) {
	cachedPin = memguard.NewEnclave(pin)
}

func GetPin() ([]byte, error) {
	approved := biometrics.CheckBiometrics(biometrics.SSHKey)
	if approved {
		buffer, err := cachedPin.Open()
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	} else {
		return nil, errors.New("biometrics not approved")
	}
}

func HasPin() bool {
	return cachedPin != nil
}

func ClearPin() {
	pin, err := cachedPin.Open()
	if err != nil {
		cachedPin = nil
		return
	}
	pin.Destroy()
	cachedPin = nil
}
