package log

import (
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level func(zapcore.Level) bool

type Logs map[string]Level

// LoggerParam
type LoggerParam struct {
	Path       string `json:"path" toml:"path" yaml:"path" xml:"path"`                         // log file path
	MaxSize    int    `json:"maxSize" toml:"maxSize" yaml:"maxSize" xml:"maxSize"`             // log file max size
	MaxBackups int    `json:"maxBackups" toml:"maxBackups" yaml:"maxBackups" xml:"maxBackups"` // log file max backups
	MaxAge     int    `json:"maxAge" toml:"maxAge" yaml:"maxAge" xml:"maxAge"`                 // log file max save day
	Compress   bool   `json:"compress" toml:"compress" yaml:"compress" xml:"compress"`         // log file whether to compress
}

// DebugLevel method is Logs map value
func DebugLevel(level zapcore.Level) bool {
	return level == zap.DebugLevel
}

// InfoLevel method is Logs map value
func InfoLevel(level zapcore.Level) bool {
	return level == zap.InfoLevel
}

// WarnLevel method is Logs map value
func WarnLevel(level zapcore.Level) bool {
	return level == zap.WarnLevel
}

// ErrorLevel method is Logs map value
func ErrorLevel(level zapcore.Level) bool {
	return level == zap.ErrorLevel
}

// DPanicLevel method is Logs map value
func DPanicLevel(level zapcore.Level) bool {
	return level == zap.DPanicLevel
}

// PanicLevel method is Logs map value
func PanicLevel(level zapcore.Level) bool {
	return level == zap.PanicLevel
}

// FatalLevel method is Logs map value
func FatalLevel(level zapcore.Level) bool {
	return level == zap.FatalLevel
}

// Default
func (log *LoggerParam) Default() *zap.SugaredLogger {
	logLevels := Logs{
		"debug":  DebugLevel,
		"info":   InfoLevel,
		"warn":   WarnLevel,
		"error":  ErrorLevel,
		"dpanic": DPanicLevel,
		"panic":  PanicLevel,
		"fatal":  FatalLevel,
	}
	return log.New(logLevels)
}

// New
func (log *LoggerParam) New(logs Logs) *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "log",
		CallerKey:     "lineNum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
			pae.AppendString(t.Format("[2006-01-02 15:04:05]"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	cores := make([]zapcore.Core, 0)

	for fileName := range logs {
		if fileName == "" {
			continue
		}
		// logger write
		write := &lumberjack.Logger{
			Filename:   getLogFilePath(log.Path, fileName),
			MaxSize:    log.MaxSize,
			MaxBackups: log.MaxBackups,
			MaxAge:     log.MaxAge,
			Compress:   log.Compress,
		}
		// the log level
		level := zap.LevelEnablerFunc(logs[fileName])
		// log core
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(write)), level)
		cores = append(cores, core)
	}
	// append default info log level
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.InfoLevel))

	core := zapcore.NewTee(cores...)
	caller := zap.AddCaller()
	development := zap.Development()

	return zap.New(core, caller, development, zap.Fields()).Sugar()
}

// logFile logger file out path
func getLogFilePath(filePath, fileName string) string {
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	fp := strings.Split(filePath, "/")
	realPath := make([]string, 0)
	for i := range fp {
		if fp[i] != "" {
			realPath = append(realPath, fp[i])
		}
	}

	filePath = strings.Join(realPath, "/")
	if filePath == "" {
		filePath = "."
	}

	if !strings.HasSuffix(fileName, ".log") {
		fileName = strings.ReplaceAll(fileName, ".", "_")
		fileName = fileName + ".log"
	}
	return filePath + "/" + fileName
}
