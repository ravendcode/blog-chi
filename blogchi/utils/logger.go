package utils

import (
	"fmt"
	"log"
	"os"
)

// Logger interface for logger
type Logger interface {
	Info(v ...interface{})
	Err(v ...interface{})
	Infof(format string, v ...interface{})
	Errf(format string, v ...interface{})
}

type logger struct {
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

func (l *logger) Info(v ...interface{}) {
	l.InfoLogger.Output(2, fmt.Sprintln(v...))
}

func (l *logger) Err(v ...interface{}) {
	l.ErrLogger.Output(2, fmt.Sprintln(v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.InfoLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) Errf(format string, v ...interface{}) {
	l.ErrLogger.Output(2, fmt.Sprintf(format, v...))
}

// NewLogger create new logger
func NewLogger() Logger {
	return &logger{
		InfoLogger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Lshortfile),
		ErrLogger:  log.New(os.Stderr, "ERR: ", log.Ldate|log.Lshortfile),
	}
}
