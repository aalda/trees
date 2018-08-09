package log

import (
	"fmt"
	"io"
	"log"
)

type errorLogger struct {
	log.Logger
}

func newError(out io.Writer, prefix string, flag int) *errorLogger {
	var l errorLogger

	l.SetOutput(out)
	l.SetPrefix(prefix)
	l.SetFlags(flag)
	return &l
}

// A impl 'l errorLogger' qed/log.Logger
func (l errorLogger) Error(v ...interface{}) {
	l.Output(caller, fmt.Sprint(v...))
	osExit(1)
}

func (l errorLogger) Errorf(format string, v ...interface{}) {
	l.Output(caller, fmt.Sprintf(format, v...))
	osExit(1)
}

func (l errorLogger) Info(v ...interface{})                  { return }
func (l errorLogger) Debug(v ...interface{})                 { return }
func (l errorLogger) Infof(format string, v ...interface{})  { return }
func (l errorLogger) Debugf(format string, v ...interface{}) { return }
