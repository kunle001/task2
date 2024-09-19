package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(s string, keysAndValues ...interface{})
}

type zapLogger struct {
	log *zap.SugaredLogger
}

func NewLogger() (Logger, error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "service.log"),
		MaxSize:    10, // this is 10 megabytes.
		MaxAge:     30, // keeping for 30 days.
		MaxBackups: 3,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{})
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileLogger), zapcore.DebugLevel),
	)

	logger := zap.New(core)

	return &zapLogger{log: logger.Sugar()}, nil
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log.Fatalw(msg, keysAndValues...)
}
