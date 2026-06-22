// Package service represents all the services of the application
package service

import (
	"context"
	"su-server/internal/model"
	"su-server/internal/repository"
)

type LeaderboardService struct {
	repo *repository.LeaderboardRepository
}

func NewLeaderboardService(repo *repository.LeaderboardRepository) *LeaderboardService {
	return &LeaderboardService{repo: repo}
}

// GetLeaderboard returns all entries ranked by step count
func (s *LeaderboardService) GetLeaderboard(ctx context.Context) ([]model.LeaderboardEntry, error) {
	return s.repo.GetRanked(ctx)
}

// GetUserRank returns a single user's rank
func (s *LeaderboardService) GetUserRank(ctx context.Context, userID int) (*model.LeaderboardEntry, error) {
	return s.repo.GetUserRank(ctx, userID)
}

// UpdateEntry upserts a user's step count in the leaderboard
func (s *LeaderboardService) UpdateEntry(ctx context.Context, userID int, stepCount int) (*model.Leaderboard, error) {
	return s.repo.Upsert(ctx, userID, stepCount)
}

// Reset clears the leaderboard for weekly reset
func (s *LeaderboardService) Reset(ctx context.Context) error {
	return s.repo.Reset(ctx)
}
