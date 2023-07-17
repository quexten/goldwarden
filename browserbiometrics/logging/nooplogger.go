//go:build !logging

package logging

func Debugf(format string, args ...interface{}) {
}

func Errorf(format string, args ...interface{}) {
}

func Panicf(format string, args ...interface{}) {
}
