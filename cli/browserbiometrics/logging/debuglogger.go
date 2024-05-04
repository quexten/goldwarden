//go:build debuglogging

package logging

import (
	"fmt"
	"os"
	"time"
)

func writeLog(level string, format string, args ...interface{}) {
	file, _ := os.OpenFile("/tmp/DELETE_ME.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// log with date and time
	file.WriteString("[" + level + "] ")
	file.WriteString(time.Now().Format("2006-01-02 15:04:05") + " ")
	file.WriteString(fmt.Sprintf(format, args...))
	file.WriteString("\n")
	file.Close()
}

func Debugf(format string, args ...interface{}) {
	writeLog("DEBUG", format, args...)
}

func Errorf(format string, args ...interface{}) {
	writeLog("ERROR", format, args...)
}

func Panicf(format string, args ...interface{}) {
	writeLog("PANIC", format, args...)
	panic("")
}
