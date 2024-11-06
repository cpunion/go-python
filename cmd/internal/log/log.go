package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Init initializes the logger with the given verbosity level
func Init(verbose bool) {
	config := zap.NewDevelopmentConfig()
	if !verbose {
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	// Customize output format
	config.EncoderConfig.TimeKey = ""       // Remove timestamp
	config.EncoderConfig.LevelKey = ""      // Remove log level
	config.EncoderConfig.CallerKey = ""     // Remove caller
	config.EncoderConfig.NameKey = ""       // Remove logger name
	config.EncoderConfig.StacktraceKey = "" // Remove stacktrace
	config.DisableCaller = true
	config.DisableStacktrace = true

	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	sugar = logger.Sugar()
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// Infof logs a formatted info message
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return logger.Sync()
}
