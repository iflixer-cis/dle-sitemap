package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

const healthTimeout = 3 * time.Second

// health проверяет связь с БД. Не требует заголовка X-Domain-Host,
// поэтому подходит для healthcheck в swarm/docker.
func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthTimeout)
	defer cancel()

	if err := s.dbService.Ping(ctx); err != nil {
		log.Println("healthcheck failed:", err)
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
