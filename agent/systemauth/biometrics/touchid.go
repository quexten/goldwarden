//go:build darwin

package biometrics

func CheckBiometrics(approvalType Approval) bool {
	ok, err := touchid.Authenticate(approvalType.String())
	if err != nil {
		log.Fatal(err)
	}

	if ok {
		log.Printf("Authenticated")
		return true
	} else {
		log.Fatal("Failed to authenticate")
		return false
	}
}
