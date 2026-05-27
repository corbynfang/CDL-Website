package services

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

type TransferFilters struct {
	Season   string
	GameCode string
	TeamID   string
	PlayerID string
}

type TransferService struct {
	store store.TransferStore
}

func NewTransferService(s store.TransferStore) *TransferService {
	return &TransferService{store: s}
}

func (ts *TransferService) List(ctx context.Context, f TransferFilters) ([]models.PlayerTransfer, error) {
	return ts.store.List(ctx, f.Season, f.GameCode, f.TeamID, f.PlayerID)
}
