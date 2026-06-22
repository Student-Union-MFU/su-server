// Package handler is a controller for handling the requests
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"su-server/internal/model"
	"su-server/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
)

type StepsHandler struct {
	service *service.StepsService
}

func NewStepsHandler(service *service.StepsService) *StepsHandler {
	return &StepsHandler{service: service}
}

// GetStepsByUserID handles GET /steps/{userID}
func (h *StepsHandler) GetStepsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	steps, err := h.service.GetStepsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get steps: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(steps)
}

// GetStepsByDateRange handles GET /steps/{userID}/range?from=2026-06-01&to=2026-06-22
func (h *StepsHandler) GetStepsByDateRange(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		http.Error(w, "invalid from date, use YYYY-MM-DD", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		http.Error(w, "invalid to date, use YYYY-MM-DD", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	steps, err := h.service.GetStepsByDateRange(r.Context(), userID, from, to)
	if err != nil {
		http.Error(w, "Failed to get steps: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(steps)
}

// SyncSteps handles POST /steps/sync — single day sync
func (h *StepsHandler) SyncSteps(w http.ResponseWriter, r *http.Request) {
	var steps model.Steps
	if err := json.NewDecoder(r.Body).Decode(&steps); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	result, err := h.service.SyncSteps(r.Context(), steps)
	if err != nil {
		http.Error(w, "Failed to sync steps: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// SyncManySteps handles POST /steps/sync/bulk — sync multiple days at once
func (h *StepsHandler) SyncManySteps(w http.ResponseWriter, r *http.Request) {
	var stepsList []model.Steps
	if err := json.NewDecoder(r.Body).Decode(&stepsList); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	results, err := h.service.SyncManySteps(r.Context(), stepsList)
	if err != nil {
		http.Error(w, "Failed to sync steps: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(results)
}
