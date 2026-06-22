// Package service represents all the services of the application
package service

import (
	"context"
	"su-server/internal/model"
	"su-server/internal/repository"
	"time"
)

type StepsService struct {
	repo *repository.StepsRepository
}

func NewStepsService(repo *repository.StepsRepository) *StepsService {
	return &StepsService{repo: repo}
}

// SyncSteps is for POST Request — upserts a single day's steps
func (s *StepsService) SyncSteps(ctx context.Context, steps model.Steps) (*model.Steps, error) {
	return s.repo.Upsert(ctx, steps)
}

// SyncManySteps is for POST Request — bulk sync multiple days at once
func (s *StepsService) SyncManySteps(ctx context.Context, stepsList []model.Steps) ([]model.Steps, error) {
	return s.repo.UpsertMany(ctx, stepsList)
}

// GetStepsByUserID is for GET Request — get all steps for a user
func (s *StepsService) GetStepsByUserID(ctx context.Context, userID int) ([]model.Steps, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetStepsByDateRange is for GET Request — get steps between two dates
func (s *StepsService) GetStepsByDateRange(ctx context.Context, userID int, from, to time.Time) ([]model.Steps, error) {
	return s.repo.GetByDateRange(ctx, userID, from, to)
}
