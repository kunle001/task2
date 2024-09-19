package ports

import "file-mod-tracker/internal/core/domain"

type FileMonitorService interface {
	GetFileStats(directory string) ([]domain.FileInfo, error)
	EnqueueCommands(commands []string) error
}

type OsqueryAdapter interface {
	GetFileStats(directory string) ([]domain.FileInfo, error)
}

type WorkerAdapter interface {
	EnqueueCommands(commands []string) error
	Start()
	Stop()
	GetFileChanges() []domain.FileInfo
}
