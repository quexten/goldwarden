//go:build linux || freebsd

package processsecurity

import "golang.org/x/sys/unix"

func DisableDumpable() error {
	return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
}
