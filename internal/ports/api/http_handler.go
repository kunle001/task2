package api

import "net/http"

type HTTPHandler interface {
	GetFileStats(w http.ResponseWriter, r *http.Request)
	EnqueueCommands(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetLogs(w http.ResponseWriter, r *http.Request)
}

type HTTPServer interface {
	Start()
}
