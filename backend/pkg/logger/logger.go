package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.SugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new logger instance
func New() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, _ := config.Build()
	return &Logger{zapLogger.Sugar()}
}

// NewDevelopment creates a development logger
func NewDevelopment() *Logger {
	zapLogger, _ := zap.NewDevelopment()
	return &Logger{zapLogger.Sugar()}
}




