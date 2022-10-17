package log

import "time"

// Record holds information about log
type Record struct {
	Time    time.Time // time when the log produced
	Message string    // content of the log
	Line    string    // in which file and line that the log produced
	Level   Level     // level of the log
}

func NewRecord(time time.Time, message, line string, level Level) *Record {
	return &Record{
		Time:    time,
		Message: message,
		Line:    line,
		Level:   level,
	}
}
