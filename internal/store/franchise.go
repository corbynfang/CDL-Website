package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type FranchiseStore interface {
	List(ctx context.Context) ([]models.Franchise, error)
	GetByKey(ctx context.Context, key string) (*models.Franchise, error)
	GetTeamsByFranchiseID(ctx context.Context, franchiseID uint) ([]models.Team, error)
}

type gormFranchiseStore struct{ db *gorm.DB }

func NewGormFranchiseStore(db *gorm.DB) FranchiseStore { return &gormFranchiseStore{db: db} }

func (s *gormFranchiseStore) List(ctx context.Context) ([]models.Franchise, error) {
	var franchises []models.Franchise
	err := s.db.WithContext(ctx).
		Where("franchise_key != ''").
		Order("name ASC").
		Find(&franchises).Error
	return franchises, err
}

func (s *gormFranchiseStore) GetByKey(ctx context.Context, key string) (*models.Franchise, error) {
	var franchise models.Franchise
	if err := s.db.WithContext(ctx).Where("franchise_key = ?", key).First(&franchise).Error; err != nil {
		return nil, err
	}
	return &franchise, nil
}

func (s *gormFranchiseStore) GetTeamsByFranchiseID(ctx context.Context, franchiseID uint) ([]models.Team, error) {
	var teams []models.Team
	err := s.db.WithContext(ctx).
		Where("franchise_id = ?", franchiseID).
		Order("valid_from ASC").
		Find(&teams).Error
	return teams, err
}
