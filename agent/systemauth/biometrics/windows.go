//go:build windows

package biometrics

func CheckBiometrics(approvalType Approval) bool {
	log.Info("Biometrics undefined on windows... skipping")
	return true
}

func BiometricsWorking() bool {
	return false
}
