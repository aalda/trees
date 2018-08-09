package log

import (
	"log"
)

type silentLogger struct {
	log.Logger
}

func newSilent() *silentLogger {
	return &silentLogger{}
}

// A impl 'l Nologger' qed/log.Logger
func (l silentLogger) Error(v ...interface{})                 { osExit(1) }
func (l silentLogger) Info(v ...interface{})                  { return }
func (l silentLogger) Debug(v ...interface{})                 { return }
func (l silentLogger) Errorf(format string, v ...interface{}) { osExit(1) }
func (l silentLogger) Infof(format string, v ...interface{})  { return }
func (l silentLogger) Debugf(format string, v ...interface{}) { return }
