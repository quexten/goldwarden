//go:build linux && uinput

package autotype

import "github.com/quexten/goldwarden/autofill/autotype/uinput"

func TypeString(text string, layout string) error {
	return uinput.TypeString(text, layout)
}
