package log

import (
	"fmt"
	"strings"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

func (l Level) string() (lvl string) {
	switch l {
	case DebugLevel:
		lvl = "D"
	case InfoLevel:
		lvl = "I"
	case WarningLevel:
		lvl = "W"
	case ErrorLevel:
		lvl = "E"
	case FatalLevel:
		lvl = "F"
	default:
		lvl = "U"
	}

	return
}

func parseLevel(lvl string) (level Level, err error) {
	switch strings.ToLower(lvl) {
	case "d":
		fallthrough
	case "debug":
		level = DebugLevel
	case "i":
		fallthrough
	case "info":
		level = InfoLevel
	case "w":
		fallthrough
	case "warning":
		level = WarningLevel
	case "e":
		fallthrough
	case "error":
		level = ErrorLevel
	case "f":
		fallthrough
	case "fatal":
		level = FatalLevel
	default:
		err = fmt.Errorf("Invalid log level: %s", lvl)
	}

	return
}
