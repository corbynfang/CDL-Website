package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type SeasonStore interface {
	List(ctx context.Context) ([]models.Season, error)
	GetByID(ctx context.Context, id int) (*models.Season, error)
	GetActive(ctx context.Context) (*models.Season, error)
}

type gormSeasonStore struct{ db *gorm.DB }

func NewGormSeasonStore(db *gorm.DB) SeasonStore { return &gormSeasonStore{db: db} }

func (s *gormSeasonStore) List(ctx context.Context) ([]models.Season, error) {
	var seasons []models.Season
	err := s.db.WithContext(ctx).Order("start_date DESC").Find(&seasons).Error
	return seasons, err
}

func (s *gormSeasonStore) GetByID(ctx context.Context, id int) (*models.Season, error) {
	var season models.Season
	if err := s.db.WithContext(ctx).First(&season, id).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

func (s *gormSeasonStore) GetActive(ctx context.Context) (*models.Season, error) {
	var season models.Season
	if err := s.db.WithContext(ctx).Where("is_active = ?", true).First(&season).Error; err != nil {
		return nil, err
	}
	return &season, nil
}
