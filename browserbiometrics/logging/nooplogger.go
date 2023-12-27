//go:build !debuglogging

package logging

func Debugf(format string, args ...interface{}) {
}

func Errorf(format string, args ...interface{}) {
}

func Panicf(format string, args ...interface{}) {
}
