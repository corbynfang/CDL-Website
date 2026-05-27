package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type TransferStore interface {
	List(ctx context.Context, season, gameCode, teamID, playerID string) ([]models.PlayerTransfer, error)
}

type gormTransferStore struct{ db *gorm.DB }

func NewGormTransferStore(db *gorm.DB) TransferStore { return &gormTransferStore{db: db} }

func (s *gormTransferStore) List(ctx context.Context, season, gameCode, teamID, playerID string) ([]models.PlayerTransfer, error) {
	query := s.db.WithContext(ctx).
		Preload("Player").
		Preload("FromTeam").
		Preload("ToTeam")

	if season != "" {
		query = query.Where("season = ?", season)
	}
	if gameCode != "" {
		query = query.Where("game_code = ?", gameCode)
	}
	if teamID != "" {
		query = query.Where("from_team_id = ? OR to_team_id = ?", teamID, teamID)
	}
	if playerID != "" {
		query = query.Where("player_id = ?", playerID)
	}

	var transfers []models.PlayerTransfer
	err := query.Order("transfer_date DESC").Find(&transfers).Error
	return transfers, err
}
