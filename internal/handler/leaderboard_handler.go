// Package handler is a controller for handling the requests
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"su-server/internal/service"

	"github.com/go-chi/chi/v5"
)

type LeaderboardHandler struct {
	service *service.LeaderboardService
}

func NewLeaderboardHandler(service *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

// GetLeaderboard handles GET /leaderboard
func (h *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.GetLeaderboard(r.Context())
	if err != nil {
		http.Error(w, "Failed to get leaderboard: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// GetUserRank handles GET /leaderboard/{userID}
func (h *LeaderboardHandler) GetUserRank(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	entry, err := h.service.GetUserRank(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user rank: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// UpdateEntry handles POST /leaderboard/update
func (h *LeaderboardHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID    int `json:"user_id"`
		StepCount int `json:"step_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		slog.Error("Error Information", " = ", err)
		return
	}

	result, err := h.service.UpdateEntry(r.Context(), body.UserID, body.StepCount)
	if err != nil {
		http.Error(w, "Failed to update leaderboard: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Reset handles POST /leaderboard/reset
func (h *LeaderboardHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Reset(r.Context()); err != nil {
		http.Error(w, "Failed to reset leaderboard: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Error Information", " = ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Leaderboard reset",
		Status:  true,
	})
}
