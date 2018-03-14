package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shiena/ansicolor"
)

const (
	red     string = "\x1b[91m"
	green   string = "\x1b[92m"
	yellow  string = "\x1b[93m"
	cyan    string = "\x1b[96m"
	white   string = "\x1b[39m"
	blue    string = "\x1b[94m"
	magenta string = "\x1b[95m"
	gray    string = "\x1b[90m"
)

// Logger interface for logger
type Logger interface {
	Info(v ...interface{})
	Err(v ...interface{})
	Infof(format string, v ...interface{})
	Errf(format string, v ...interface{})
	Red() *logger
	Green() *logger
	Yellow() *logger
	Cyan() *logger
	White() *logger
	Blue() *logger
	Magenta() *logger
	Gray() *logger
}

type logger struct {
	Color      string
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

func (l *logger) Info(v ...interface{}) {
	l.InfoLogger.Output(2, colorWrap(l, strings.TrimRight(fmt.Sprintln(v...), "\n")))
}

func (l *logger) Err(v ...interface{}) {
	l.ErrLogger.Output(2, colorWrap(l, strings.TrimRight(fmt.Sprintln(v...), "\n")))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.InfoLogger.Output(2, colorWrap(l, fmt.Sprintf(format, v...)))
}

func (l *logger) Errf(format string, v ...interface{}) {
	l.ErrLogger.Output(2, colorWrap(l, fmt.Sprintf(format, v...)))
}

func (l *logger) Red() *logger {
	l.Color = red
	return l
}

func (l *logger) Green() *logger {
	l.Color = green
	return l
}

func (l *logger) Yellow() *logger {
	l.Color = yellow
	return l
}

func (l *logger) Cyan() *logger {
	l.Color = cyan
	return l
}

func (l *logger) White() *logger {
	l.Color = white
	return l
}

func (l *logger) Blue() *logger {
	l.Color = blue
	return l
}
func (l *logger) Magenta() *logger {
	l.Color = magenta
	return l
}
func (l *logger) Gray() *logger {
	l.Color = gray
	return l
}

func colorWrap(l *logger, str string) string {
	result := fmt.Sprintf("%s%s%s", l.Color, str, white)
	l.Color = green
	return result
}

// NewLogger create new logger
func NewLogger() Logger {
	return &logger{
		Color:      green,
		InfoLogger: log.New(ansicolor.NewAnsiColorWriter(os.Stdout), cyan+"INFO: ", log.Ldate|log.Lshortfile),
		ErrLogger:  log.New(ansicolor.NewAnsiColorWriter(os.Stderr), red+"ERR: ", log.Ldate|log.Lshortfile),
	}
}
