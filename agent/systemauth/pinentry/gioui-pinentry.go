//go:build windows || darwin

package pinentry

import (
	"github.com/quexten/goldwarden/agent/systemauth/pinentry/giouipinentry"
)

func GetPassword(title string, description string) (string, error) {
	var resultChan = make(chan string)
	giouipinentry.GetPin(title, description, func(pin string) {
		resultChan <- pin
	})
	return <-resultChan, nil
}

func GetApproval(title string, description string) (bool, error) {
	var resultChan = make(chan bool)
	giouipinentry.GetApproval(title, description, func(approved bool) {
		resultChan <- approved
	})
	return <-resultChan, nil
}
