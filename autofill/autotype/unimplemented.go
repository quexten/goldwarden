//go:build !linux

package autotype

func TypeString(text string, layout string) error {
	return errors.New("Not implemented")
}
