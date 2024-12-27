package common

import (
	"log"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(msg string) {
	log.Printf("[INFO] [%s] %s", l.prefix, msg)
}

func (l *Logger) Warn(msg string) {
	log.Printf("[WARN] [%s] %s", l.prefix, msg)
}

func (l *Logger) Error(err error) {
	log.Printf("[ERROR] [%s] %v", l.prefix, err)
}

func (l *Logger) Fatal(err error) {
	log.Fatalf("[FATAL] [%s] %v", l.prefix, err)
}
