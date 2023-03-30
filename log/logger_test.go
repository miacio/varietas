package log_test

import (
	"testing"

	"github.com/miacio/varietas/log"
)

func TestLogger(t *testing.T) {
	logParam := log.LoggerParam{
		Path:       "./log",
		MaxSize:    256,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   false,
	}
	logLevels := log.Logs{
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"error": log.ErrorLevel,
	}
	log := logParam.New(logLevels)

	log.Infoln("init success")
}
