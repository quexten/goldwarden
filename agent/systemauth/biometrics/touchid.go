//go:build darwin

package biometrics

import (
	touchid "github.com/lox/go-touchid"
)

func CheckBiometrics(approvalType Approval) bool {
	ok, err := touchid.Authenticate(approvalType.String())
	if err != nil {
		log.Error(err)
	}

	if ok {
		log.Info("Authenticated")
		return true
	} else {
		log.Error("Failed to authenticate")
		return false
	}
}
