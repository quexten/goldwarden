package systemauth

import (
	"github.com/LlamaNite/llamalog"
	"github.com/amenzhinsky/go-polkit"
)

var log = llamalog.NewLogger("Goldwarden", "Systemauth")

type Approval string

const (
	AccessCredential Approval = "com.quexten.goldwarden.accesscredential"
	ChangePin        Approval = "com.quexten.goldwarden.changepin"
	SSHKey           Approval = "com.quexten.goldwarden.usesshkey"
)

func (a Approval) String() string {
	return string(a)
}

func CheckBiometrics(approvalType Approval) bool {
	log.Info("Checking biometrics for %s", approvalType.String())

	authority, err := polkit.NewAuthority()
	if err != nil {
		return false
	}

	result, err := authority.CheckAuthorization(
		approvalType.String(),
		nil,
		polkit.CheckAuthorizationAllowUserInteraction, "",
	)

	if err != nil {
		return false
	}

	log.Info("Biometrics result: %t", result.IsAuthorized)

	return result.IsAuthorized
}
