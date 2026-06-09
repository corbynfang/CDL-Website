package store

import (
	"context"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type UserStore interface {
	GetBySupabaseUID(ctx context.Context, uid string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

type gormUserStore struct{ db *gorm.DB }

func NewGormUserStore(db *gorm.DB) UserStore { return &gormUserStore{db: db} }

func (s *gormUserStore) GetBySupabaseUID(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).Where("supabase_uid = ? AND deleted_at IS NULL", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *gormUserStore) Create(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *gormUserStore) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ThreadPost{}).
			Where("user_id = ? AND deleted_at IS NULL", id).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		return tx.Model(&models.User{}).Where("id = ?", id).Update("deleted_at", now).Error
	})
}
