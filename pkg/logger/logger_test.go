package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Helper function to create a test logger that logs to an in-memory buffer.
func newTestLogger() (*zapLogger, *bytes.Buffer) {
	// Create an in-memory buffer to capture log output.
	var buf bytes.Buffer
	writer := zapcore.AddSync(&buf)

	// Create a custom encoder and core for the test.
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	})
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)

	// Create a Sugared logger from the core.
	logger := zap.New(core).Sugar()

	return &zapLogger{log: logger}, &buf
}

func TestZapLogger_Info(t *testing.T) {
	logger, buf := newTestLogger()

	logger.Info("info message", "key", "value")

	assert.Contains(t, buf.String(), `"msg":"info message"`)
	assert.Contains(t, buf.String(), `"level":"info"`)
	assert.Contains(t, buf.String(), `"key":"value"`)
}

func TestZapLogger_Error(t *testing.T) {
	logger, buf := newTestLogger()

	logger.Error("error message", "key", "value")

	assert.Contains(t, buf.String(), `"msg":"error message"`)
	assert.Contains(t, buf.String(), `"level":"error"`)
	assert.Contains(t, buf.String(), `"key":"value"`)
}
