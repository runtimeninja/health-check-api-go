package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct {
	DB *sql.DB
}

type healthResp struct {
	Status string `json:"status"`
	Time   string `json:"time"`
	DB     string `json:"db,omitempty"`
}

// /live : only checks process is up (no DB check)
func (h HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResp{
		Status: "ok",
		Time:   time.Now().Format(time.RFC3339),
	})
}

// /ready : checks DB readiness (if DB down then 503)
func (h HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// tight timeout so health never hangs
	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
	defer cancel()

	dbStatus := "ok"
	code := http.StatusOK

	if err := h.DB.PingContext(ctx); err != nil {
		dbStatus = "down"
		code = http.StatusServiceUnavailable
	}

	status := "ok"
	if code != http.StatusOK {
		status = "down"
	}

	writeJSON(w, code, healthResp{
		Status: status,
		Time:   time.Now().Format(time.RFC3339),
		DB:     dbStatus,
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
