package worker_test

import (
	"file-mod-tracker/internal/adapters/worker"
	"file-mod-tracker/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOsqueryAdapter mocks the OsqueryAdapter interface
type mockOsqueryAdapter struct {
	mock.Mock
}

func (m *mockOsqueryAdapter) GetFileStats(directory string) ([]domain.FileInfo, error) {
	args := m.Called(directory)
	return args.Get(0).([]domain.FileInfo), args.Error(1)
}

// MockLogger mocks the Logger interface
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *mockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *mockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func TestEnqueueCommands_Success(t *testing.T) {
	mockOsquery := new(mockOsqueryAdapter)
	mockLogger := new(mockLogger)

	workerAdapter := worker.NewAdapter(mockLogger, mockOsquery, "/test", 5)

	mockLogger.On("Info", "Command enqueued", []interface{}{"command", "echo test"}).Once()

	err := workerAdapter.EnqueueCommands([]string{"echo test"})
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestEnqueueCommands_QueueFull(t *testing.T) {
	mockOsquery := new(mockOsqueryAdapter)
	mockLogger := new(mockLogger)

	workerAdapter := worker.NewAdapter(mockLogger, mockOsquery, "/test", 5)

	// Fill up the queue
	for i := 0; i < 100; i++ {
		mockLogger.On("Info", "Command enqueued", []interface{}{"command", "echo test"}).Once()
		workerAdapter.EnqueueCommands([]string{"echo test"})
	}

	// Expect the error log when the queue is full
	mockLogger.On("Error", "Command queue is full", []interface{}{"command", "echo test"}).Once()

	err := workerAdapter.EnqueueCommands([]string{"echo test"})
	assert.EqualError(t, err, "command queue is full")
	mockLogger.AssertExpectations(t)
}

func TestWorkerThread_ExecutesCommandSuccessfully(t *testing.T) {
	mockOsquery := new(mockOsqueryAdapter)
	mockLogger := new(mockLogger)

	workerAdapter := worker.NewAdapter(mockLogger, mockOsquery, "/test", 5)

	mockLogger.On("Info", "Command enqueued", []interface{}{"command", "echo test"}).Once()
	mockLogger.On("Info", "Executing command", []interface{}{"command", "echo test"}).Once()
	mockLogger.On("Info", "Command executed successfully", []interface{}{"command", "echo test", "output", "test\n"}).Once()

	workerAdapter.EnqueueCommands([]string{"echo test"})
	go workerAdapter.Start()
	time.Sleep(500 * time.Millisecond) // Allow some time for the goroutine to execute
	workerAdapter.Stop()

	mockLogger.AssertExpectations(t)
}
