package logger

import (
    "log"
)

type stdLogger struct{}

func newStdLogger() Logger {
    return &stdLogger{}
}

func (l *stdLogger) Info(msg string) {
    log.Printf("[INFO] %s", msg)
}

func (l *stdLogger) Error(msg string) {
    log.Printf("[ERROR] %s", msg)
}

func (l *stdLogger) Infof(format string, args ...interface{}) {
    log.Printf("[INFO] "+format, args...)
}

func (l *stdLogger) Errorf(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}
