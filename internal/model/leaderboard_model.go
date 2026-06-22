// Package model defines the data structures used throughout the application.
package model

import "time"

// Leaderboard represents a user's total step count for the leaderboard.
type Leaderboard struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	StepCount int       `json:"step_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LeaderboardEntry represents a leaderboard entry with user info and rank.
type LeaderboardEntry struct {
	Rank      int    `json:"rank"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	StepCount int    `json:"step_count"`
}
