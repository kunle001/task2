package http

import (
	"net/http"

	"file-mod-tracker/internal/ports"
)

type httpHandler struct {
	fileMonitorService ports.FileMonitorService
}

func NewHandler(fileMonitorService ports.FileMonitorService) *httpHandler {
	return &httpHandler{
		fileMonitorService: fileMonitorService,
	}
}

func (h *httpHandler) GetFileStats(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

func (h *httpHandler) EnqueueCommands(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

func (h *httpHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

func (h *httpHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	// Implementation
}
