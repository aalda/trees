package log

import (
	"fmt"
	"log"
	"os"
)

// Log levels constants
const (
	SILENT = "silent"
	ERROR  = "error"
	INFO   = "info"
	DEBUG  = "debug"

	caller = 4
)

// Private interface for the std variable.
type logger interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

// The default logger is an log.ERROR level.
var std logger = newError(os.Stdout, "Qed: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

// To allow mocking we require a switchable variable.
var osExit = os.Exit

// Below is the public interface for the logger, a proxy for the switchable
// implementation defined in std.

// Error is the public log function to write to stdOut and stop execution.
func Error(v ...interface{}) {
	std.Error(v...)
}

var (

	// Fatal is the public log function to write to stdOut and stop execution
	// Same as Error.
	Fatal func(...interface{}) = Error

	// Panic is the public log function to write to stdOut and stop execution
	// Same as Error.
	Panic func(...interface{}) = Error
)

// Errorf is the public log function with params to write to stdOut and stop
// execution.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

var (

	// Fatalf is the public log function with params to write to stdOut and
	// stop execution. Same as Errorf
	Fatalf func(string, ...interface{}) = Errorf

	// Panicf is the public log function with params to write to stdOut and
	// stop execution. Same as Errorf
	Panicf func(string, ...interface{}) = Errorf
)

// Info is the public log function to write information relative to the usage
// of the qed package.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Info is the public log function to write information with params relative
// to the usage of the qed package.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Debug is the public log function to write information relative to internal
// debug information.
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugf is the public log function to write information with params relative
// to internal debug information.
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// SetLogger is a function that switches between verbosity loggers. Default
// is error level. Available levels are "silent", "debug", "info" and "error".
func SetLogger(namespace, level string) {

	prefix := fmt.Sprintf("%s: ", namespace)

	switch level {
	case SILENT:
		std = newSilent()
	case ERROR:
		std = newError(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	case INFO:
		std = newInfo(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	case DEBUG:
		std = newDebug(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	default:
		l := newInfo(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		l.Infof("Incorrect level of verbosity (%v) fallback to log.INFO", level)
		std = l
	}

}
