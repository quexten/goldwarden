//go:build linux || freebsd

package processsecurity

func DisableDumpable() error {
	// return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
	return nil
}
