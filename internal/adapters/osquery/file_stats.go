package osquery

import (
	"file-mod-tracker/internal/core/domain"
	"file-mod-tracker/pkg/logger"
	"os"
	"path/filepath"
	"time"
)

type OsqueryAdapter struct {
	logger logger.Logger
}

func NewAdapter(logger logger.Logger) *OsqueryAdapter {
	return &OsqueryAdapter{logger: logger}
}

func (a *OsqueryAdapter) GetFileStats(directory string) ([]domain.FileInfo, error) {
	var fileInfos []domain.FileInfo

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileInfos = append(fileInfos, domain.FileInfo{
				Path:         path,
				LastModified: info.ModTime().Format(time.RFC3339),
				Size:         info.Size(),
			})
		}
		return nil
	})

	if err != nil {
		a.logger.Error("Failed to walk directory", "error", err)
		return nil, err
	}

	return fileInfos, nil
}
