package worker

import (
	"file-mod-tracker/internal/core/domain"
	"file-mod-tracker/internal/ports"
	"file-mod-tracker/pkg/logger"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type WorkerAdapter struct {
	commandQueue     chan string
	logger           logger.Logger
	wg               sync.WaitGroup
	stopChan         chan struct{}
	osqueryAdapter   ports.OsqueryAdapter
	monitoredDir     string
	fileChanges      []domain.FileInfo
	fileChangesMutex sync.Mutex
	checkFrequency   int
}

func NewAdapter(logger logger.Logger, osqueryAdapter ports.OsqueryAdapter, monitoredDir string, specifiedFrequency int) *WorkerAdapter {
	return &WorkerAdapter{
		commandQueue:   make(chan string, 100),
		logger:         logger,
		stopChan:       make(chan struct{}),
		osqueryAdapter: osqueryAdapter,
		monitoredDir:   monitoredDir,
		checkFrequency: specifiedFrequency,
	}
}

func (a *WorkerAdapter) EnqueueCommands(commands []string) error {
	for _, cmd := range commands {
		select {
		case a.commandQueue <- cmd:
			a.logger.Info("Command enqueued", "command", cmd)
		default:
			a.logger.Error("Command queue is full", "command", cmd)
			return fmt.Errorf("command queue is full")
		}
	}
	return nil
}

func (a *WorkerAdapter) Start() {
	a.wg.Add(2)
	go a.workerThread()
	go a.timerThread()
}

func (a *WorkerAdapter) Stop() {
	close(a.stopChan)
	a.wg.Wait()
}

func (a *WorkerAdapter) workerThread() {
	defer a.wg.Done()
	for {
		select {
		case cmd := <-a.commandQueue:
			a.logger.Info("Executing command", "command", cmd)
			parts := strings.Fields(cmd)
			if len(parts) == 0 {
				a.logger.Error("Empty command received")
				continue
			}
			command := exec.Command(parts[0], parts[1:]...)
			output, err := command.CombinedOutput()
			if err != nil {
				a.logger.Error("Command execution failed", "command", cmd, "error", err, "output", string(output))
			} else {
				a.logger.Info("Command executed successfully", "command", cmd, "output", string(output))
			}
		case <-a.stopChan:
			return
		}
	}
}

func (a *WorkerAdapter) timerThread() {
	defer a.wg.Done()
	ticker := time.NewTicker(time.Duration(a.checkFrequency) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.logger.Info("Timer thread woke up")
			stats, err := a.osqueryAdapter.GetFileStats(a.monitoredDir)
			if err != nil {
				a.logger.Error("Failed to get file stats", "error", err)
				continue
			}
			a.updateFileChanges(stats)
		case <-a.stopChan:
			return
		}
	}
}

func (a *WorkerAdapter) updateFileChanges(newStats []domain.FileInfo) {
	a.fileChangesMutex.Lock()
	defer a.fileChangesMutex.Unlock()

	// Simple implementation: just store the latest stats
	a.fileChanges = newStats
	for _, stat := range newStats {
		a.logger.Info("File info", "path", stat.Path, "lastModified", stat.LastModified, "size", stat.Size)
	}
}

func (a *WorkerAdapter) GetFileChanges() []domain.FileInfo {
	a.fileChangesMutex.Lock()
	defer a.fileChangesMutex.Unlock()
	return a.fileChanges
}
