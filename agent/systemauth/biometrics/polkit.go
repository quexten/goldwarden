//go:build linux || freebsd

package biometrics

import (
	"github.com/amenzhinsky/go-polkit"
)

const POLICY = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE policyconfig PUBLIC
 "-//freedesktop//DTD PolicyKit Policy Configuration 1.0//EN"
 "http://www.freedesktop.org/software/polkit/policyconfig-1.dtd">
<policyconfig>

  <action id="com.quexten.goldwarden.accessvault">
    <description>Allow access to the vault</description>
    <message>Allows access to the vault</message>
    <defaults>
      <allow_any>auth_self</allow_any>
      <allow_inactive>auth_self</allow_inactive>
      <allow_active>auth_self</allow_active>
    </defaults>
  </action>
  <action id="com.quexten.goldwarden.usesshkey">
    <description>Use SSH Key</description>
    <message>Authenticate to use an SSH Key from your vault</message>
    <defaults>
      <allow_any>auth_self</allow_any>
      <allow_inactive>auth_self</allow_inactive>
      <allow_active>auth_self</allow_active>
    </defaults>
  </action>
  <action id="com.quexten.goldwarden.browserbiometrics">
    <description>Browser Biometrics</description>
    <message>Authenticate to allow Goldwarden to unlock your browser</message>
    <defaults>
      <allow_any>auth_self</allow_any>
      <allow_inactive>auth_self</allow_inactive>
      <allow_active>auth_self</allow_active>
    </defaults>
  </action>

</policyconfig>`

func CheckBiometrics(approvalType Approval) bool {
	if biometricsDisabled {
		return true
	}

	log.Info("Checking biometrics for %s", approvalType.String())

	authority, err := polkit.NewAuthority()
	if err != nil {
		log.Error("Failed to create polkit authority: %s", err.Error())
		return false
	}

	result, err := authority.CheckAuthorization(
		approvalType.String(),
		nil,
		uint32(polkit.AuthenticationRequiredRetained), "",
	)

	if err != nil {
		log.Error("Failed to create polkit authority: %s", err.Error())
		return false
	}

	log.Info("Biometrics result: %t", result.IsAuthorized)

	return result.IsAuthorized
}

func BiometricsWorking() bool {
	if biometricsDisabled {
		return false
	}

	authority, err := polkit.NewAuthority()
	if err != nil {
		log.Warn("Failed to create polkit authority: %s", err.Error())
		return false
	}

	result, err := authority.EnumerateActions("en")
	if err != nil {
		log.Warn("Failed to enumerate polkit actions: %s", err.Error())
		return false
	}

	if len(result) == 0 {
		log.Warn("No polkit actions found")
		return false
	}

	testFor := AccessVault
	for _, action := range result {
		if Approval(action.ActionID) == testFor {
			return true
		}
	}

	return false
}
