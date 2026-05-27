package services

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

type FranchiseDetail struct {
	Franchise models.Franchise `json:"franchise"`
	Eras      []models.Team    `json:"eras"`
}

type FranchiseService struct {
	store store.FranchiseStore
}

func NewFranchiseService(s store.FranchiseStore) *FranchiseService {
	return &FranchiseService{store: s}
}

func (fs *FranchiseService) List(ctx context.Context) ([]models.Franchise, error) {
	return fs.store.List(ctx)
}

func (fs *FranchiseService) GetByKey(ctx context.Context, key string) (*FranchiseDetail, error) {
	franchise, err := fs.store.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	teams, err := fs.store.GetTeamsByFranchiseID(ctx, franchise.ID)
	if err != nil {
		return nil, err
	}
	return &FranchiseDetail{Franchise: *franchise, Eras: teams}, nil
}
