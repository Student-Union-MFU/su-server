// Package model defines the data structures used throughout the application.
package model

import "time"
import "github.com/jackc/pgx/v5/pgtype"

// Steps represents a user's step count for a specific day.
type Steps struct {
	ID           int         `json:"id"`
	UserID       int         `json:"user_id"`
	StepCount    int         `json:"step_count"`
	RecordedDate pgtype.Date `json:"recorded_date"`
	SyncedAt     time.Time   `json:"synced_at"`
}
