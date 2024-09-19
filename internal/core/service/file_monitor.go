package service

import (
	"file-mod-tracker/internal/core/domain"
	"file-mod-tracker/internal/ports"
	"file-mod-tracker/pkg/logger"
)

type fileMonitorService struct {
	osqueryAdapter ports.OsqueryAdapter
	workerAdapter  ports.WorkerAdapter
	logger         logger.Logger
}

func NewFileMonitorService(osqueryAdapter ports.OsqueryAdapter, workerAdapter ports.WorkerAdapter, logger logger.Logger) *fileMonitorService {
	return &fileMonitorService{
		osqueryAdapter: osqueryAdapter,
		workerAdapter:  workerAdapter,
		logger:         logger,
	}
}

func (s *fileMonitorService) GetFileStats(directory string) ([]domain.FileInfo, error) {
	return s.osqueryAdapter.GetFileStats(directory)
}

func (s *fileMonitorService) EnqueueCommands(commands []string) error {
	return s.workerAdapter.EnqueueCommands(commands)
}
