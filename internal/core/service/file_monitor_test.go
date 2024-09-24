package service

import (
	"errors"
	"file-mod-tracker/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for OsqueryAdapter
type mockOsqueryAdapter struct {
	mock.Mock
}

func (m *mockOsqueryAdapter) GetFileStats(directory string) ([]domain.FileInfo, error) {
	args := m.Called(directory)
	return args.Get(0).([]domain.FileInfo), args.Error(1)
}

// Mock for WorkerAdapter
type mockWorkerAdapter struct {
	mock.Mock
}

func (m *mockWorkerAdapter) EnqueueCommands(commands []string) error {
	args := m.Called(commands)
	return args.Error(0)
}

func (m *mockWorkerAdapter) Start() {}
func (m *mockWorkerAdapter) Stop()  {}
func (m *mockWorkerAdapter) GetFileChanges() []domain.FileInfo {
	return nil
}

// Mock for Logger
type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *mockLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *mockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func TestFileMonitorService_GetFileStats(t *testing.T) {
	// Arrange
	mockOsquery := new(mockOsqueryAdapter)
	mockWorker := new(mockWorkerAdapter)
	mockLogger := new(mockLogger)
	fileMonitor := NewFileMonitorService(mockOsquery, mockWorker, mockLogger)

	expectedFileStats := []domain.FileInfo{
		{Path: "/test/file1.txt", LastModified: "2024-09-23T12:00:00Z", Size: 1234},
	}

	mockOsquery.On("GetFileStats", "/test").Return(expectedFileStats, nil)

	// Act
	fileStats, err := fileMonitor.GetFileStats("/test")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFileStats, fileStats)
	mockOsquery.AssertExpectations(t)
}

func TestFileMonitorService_GetFileStats_Error(t *testing.T) {
	// Arrange
	mockOsquery := new(mockOsqueryAdapter)
	mockWorker := new(mockWorkerAdapter)
	mockLogger := new(mockLogger)
	fileMonitor := NewFileMonitorService(mockOsquery, mockWorker, mockLogger)

	// Return an empty slice and an error
	mockOsquery.On("GetFileStats", "/test").Return([]domain.FileInfo{}, errors.New("directory not found"))

	// Act
	fileStats, err := fileMonitor.GetFileStats("/test")

	// Assert
	assert.Empty(t, fileStats)                       // Check if the slice is empty
	assert.EqualError(t, err, "directory not found") // Assert the error message
	mockOsquery.AssertExpectations(t)
}

func TestFileMonitorService_EnqueueCommands(t *testing.T) {
	// Arrange
	mockOsquery := new(mockOsqueryAdapter)
	mockWorker := new(mockWorkerAdapter)
	mockLogger := new(mockLogger)
	fileMonitor := NewFileMonitorService(mockOsquery, mockWorker, mockLogger)

	commands := []string{"command1", "command2"}

	mockWorker.On("EnqueueCommands", commands).Return(nil)

	// Act
	err := fileMonitor.EnqueueCommands(commands)

	// Assert
	assert.NoError(t, err)
	mockWorker.AssertExpectations(t)
}

func TestFileMonitorService_EnqueueCommands_Error(t *testing.T) {
	// Arrange
	mockOsquery := new(mockOsqueryAdapter)
	mockWorker := new(mockWorkerAdapter)
	mockLogger := new(mockLogger)
	fileMonitor := NewFileMonitorService(mockOsquery, mockWorker, mockLogger)

	commands := []string{"command1", "command2"}

	mockWorker.On("EnqueueCommands", commands).Return(errors.New("queue error"))

	// Act
	err := fileMonitor.EnqueueCommands(commands)

	// Assert
	assert.EqualError(t, err, "queue error")
	mockWorker.AssertExpectations(t)
}
