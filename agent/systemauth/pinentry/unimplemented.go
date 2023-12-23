//go:build !linux

package pinentry

import "errors"

func GetPassword(title string, description string) (string, error) {
	log.Info("Asking for password is not implemented on this platform")
	return "", errors.New("Not implemented")
}

func GetApproval(title string, description string) (bool, error) {
	log.Info("Asking for approval is not implemented on this platform")
	return true, errors.New("Not implemented")
}
