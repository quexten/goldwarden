//go:build windows || darwin

package processsecurity

func DisableDumpable() error {
	// no additional dumping protection
	return nil
}
