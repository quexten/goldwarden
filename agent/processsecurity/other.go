//go:build windows || darwin

package processsecurity

func DisableDumpale() error {
	// no additional dumping protection
	return nil
}
