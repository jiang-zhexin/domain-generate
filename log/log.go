package log

import (
	"log"
	"os"
)

type Log struct {
	Disabled bool   `json:"disabled,omitempty"`
	Level    string `json:"level,omitempty"`
}

var (
	logDebug *log.Logger = log.New(os.Stderr, "[debug]\t", 0)
	logInfo  *log.Logger = log.New(os.Stderr, "[info]\t", 0)
	logWarn  *log.Logger = log.New(os.Stderr, "[warn]\t", 0)
	logError *log.Logger = log.New(os.Stderr, "[error]\t", 0)
)

func Init(l *Log) {
	if l.Disabled {
		logError = log.New(nullWriter{}, "[error]\t", 0)
		logWarn = log.New(nullWriter{}, "[warn]\t", 0)
		logInfo = log.New(nullWriter{}, "[info]\t", 0)
		logDebug = log.New(nullWriter{}, "[debug]\t", 0)
		return
	}
	switch l.Level {
	case "error":
		logWarn = log.New(nullWriter{}, "[warn]\t", 0)
		fallthrough
	case "warn":
		logInfo = log.New(nullWriter{}, "[info]\t", 0)
		fallthrough
	case "info":
		logDebug = log.New(nullWriter{}, "[debug]\t", 0)
	case "debug":
	default:
		Warn("未知日志等级！")
		Info("默认日志等级为：info")
		logDebug = log.New(nullWriter{}, "[debug]\t", 0)
	}
}

func Debug(format string, v ...any) {
	logDebug.Printf(format, v...)
}
func Info(format string, v ...any) {
	logInfo.Printf(format, v...)
}
func Warn(format string, v ...any) {
	logWarn.Printf(format, v...)
}
func Error(format string, v ...any) {
	logError.Printf(format, v...)
}

type nullWriter struct{}

func (w nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
