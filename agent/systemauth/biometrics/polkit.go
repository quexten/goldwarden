//go:build linux || freebsd

package biometrics

import (
	"github.com/amenzhinsky/go-polkit"
	"github.com/quexten/goldwarden/logging"
)

var log = logging.GetLogger("Goldwarden", "Biometrics")

const POLICY = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE policyconfig PUBLIC
 "-//freedesktop//DTD PolicyKit Policy Configuration 1.0//EN"
 "http://www.freedesktop.org/standards/PolicyKit/1.0/policyconfig.dtd">

<policyconfig>
    <action id="com.quexten.goldwarden.accesscredential">
      <description>Allow Credential Access</description>
      <message>Authenticate to allow access to a single credential</message>
      <defaults>
        <allow_any>auth_self</allow_any>
        <allow_inactive>auth_self</allow_inactive>
        <allow_active>auth_self</allow_active>
      </defaults>
    </action>
    <action id="com.quexten.goldwarden.changepin">
      <description>Approve Pin Change</description>
      <message>Authenticate to change your Goldwarden PIN.</message>
      <defaults>
        <allow_any>auth_self</allow_any>
        <allow_inactive>auth_self</allow_inactive>
        <allow_active>auth_self</allow_active>
      </defaults>
    </action>
    <action id="com.quexten.goldwarden.usesshkey">
      <description>Use Bitwarden SSH Key</description>
      <message>Authenticate to use an SSH Key from your vault</message>
      <defaults>
        <allow_any>auth_self</allow_any>
        <allow_inactive>auth_self</allow_inactive>
        <allow_active>auth_self</allow_active>
      </defaults>
    </action>
    <action id="com.quexten.goldwarden.modifyvault">
      <description>Modify Bitwarden Vault</description>
      <message>Authenticate to allow modification of your Bitvarden vault in Goldwarden</message>
      <defaults>
        <allow_any>auth_self</allow_any>
        <allow_inactive>auth_self</allow_inactive>
        <allow_active>auth_self</allow_active>
      </defaults>
    </action>
    <action id="com.quexten.goldwarden.browserbiometrics">
      <description>Browser Biometrics</description>
      <message>Authenticate to allow Goldwarden to unlock your browser.</message>
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
