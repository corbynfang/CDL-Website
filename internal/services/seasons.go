package services

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

type SeasonService struct {
	store store.SeasonStore
}

func NewSeasonService(s store.SeasonStore) *SeasonService {
	return &SeasonService{store: s}
}

func (ss *SeasonService) List(ctx context.Context) ([]models.Season, error) {
	return ss.store.List(ctx)
}

func (ss *SeasonService) GetByID(ctx context.Context, id int) (*models.Season, error) {
	return ss.store.GetByID(ctx, id)
}

func (ss *SeasonService) GetActive(ctx context.Context) (*models.Season, error) {
	return ss.store.GetActive(ctx)
}
