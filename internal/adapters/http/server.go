package http

import (
	"encoding/json"
	"net/http"

	"file-mod-tracker/internal/ports"
	"file-mod-tracker/pkg/logger"
)

type Server struct {
	fileMonitorService ports.FileMonitorService
	logger             logger.Logger
	workerAdapter      ports.WorkerAdapter
}

func NewServer(fileMonitorService ports.FileMonitorService, logger logger.Logger, workerAdapter ports.WorkerAdapter) *Server {
	return &Server{
		fileMonitorService: fileMonitorService,
		logger:             logger,
		workerAdapter:      workerAdapter,
	}
}

func (s *Server) Start(port string) error {
	http.HandleFunc("/file-stats", s.handleFileStats)
	http.HandleFunc("/enqueue-commands", s.handleEnqueueCommands)
	http.HandleFunc("/health", s.handleHealthCheck)
	http.HandleFunc("/logs", s.handleGetLogs)

	s.logger.Info("Starting HTTP server", "port", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *Server) handleFileStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	directory := r.URL.Query().Get("directory")
	if directory == "" {
		http.Error(w, "Directory parameter is required", http.StatusBadRequest)
		return
	}

	stats, err := s.fileMonitorService.GetFileStats(directory)
	if err != nil {
		s.logger.Error("Failed to get file stats", "error", err)
		http.Error(w, "Failed to get file stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleEnqueueCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var commands []string
	if err := json.NewDecoder(r.Body).Decode(&commands); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.fileMonitorService.EnqueueCommands(commands); err != nil {
		s.logger.Error("Failed to enqueue commands", "error", err)
		http.Error(w, "Failed to enqueue commands", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Implement health check logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileChanges := s.workerAdapter.GetFileChanges()
	json.NewEncoder(w).Encode(fileChanges)
}
