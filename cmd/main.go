package main

import (
	"file-mod-tracker/internal/adapters/config"
	"file-mod-tracker/internal/adapters/http"
	"file-mod-tracker/internal/adapters/osquery"
	"file-mod-tracker/internal/adapters/ui"
	"file-mod-tracker/internal/adapters/worker"
	"file-mod-tracker/internal/core/service"
	"file-mod-tracker/pkg/logger"
	"fmt"
	"os"
)

func main() {
	// Initialize logger
	log, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		log.Fatal("Failed to initialize logger")
		os.Exit(1)
	}

	log.Info("logger initialized")
	log.Info("Current working directory: ", "dir", getCurrentDirectory())

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", "error", err)
	}

	// Initialize adapters
	osqueryAdapter := osquery.NewAdapter(log)
	workerAdapter := worker.NewAdapter(log, osqueryAdapter, cfg.MonitoredDir, cfg.CheckFrequency)

	// Initialize core service
	fileMonitorService := service.NewFileMonitorService(osqueryAdapter, workerAdapter, log)

	// Initialize HTTP server
	server := http.NewServer(fileMonitorService, log, workerAdapter)

	// Initialize UI
	ui := ui.NewMacOSUI(fileMonitorService, workerAdapter)

	// Start worker threads
	workerAdapter.Start()

	// Start HTTP server in a goroutine
	go func() {
		if err := server.Start(cfg.ServerPort); err != nil {
			log.Fatal("Failed to start HTTP server", "error", err)
		}
	}()

	// Show UI (this will block until the UI is closed)
	ui.Show()

	// Stop worker threads when the application exits
	workerAdapter.Stop()
}

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("Failed to get current working directory: %e", err)
	}

	return dir
}
