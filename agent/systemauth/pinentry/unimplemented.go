//go:build !linux && !windows && !darwin && !freebsd

package pinentry

import "errors"

func getPassword(title string, description string) (string, error) {
	log.Info("Asking for password is not implemented on this platform")
	return "", errors.New("Not implemented")
}

func getApproval(title string, description string) (bool, error) {
	log.Info("Asking for approval is not implemented on this platform")
	return true, errors.New("Not implemented")
}
