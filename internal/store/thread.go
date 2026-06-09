package store

import (
	"context"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ThreadStore interface {
	GetOrCreateThread(ctx context.Context, matchID uint) (*models.MatchThread, error)
	GetPostsByThreadID(ctx context.Context, threadID uint, limit, offset int) ([]models.ThreadPost, int64, error)
	CreatePost(ctx context.Context, post *models.ThreadPost) error
	GetPost(ctx context.Context, id uint) (*models.ThreadPost, error)
	UpdatePost(ctx context.Context, id uint, body string) error
	SoftDeletePost(ctx context.Context, id uint) error
}

type gormThreadStore struct{ db *gorm.DB }

func NewGormThreadStore(db *gorm.DB) ThreadStore { return &gormThreadStore{db: db} }

func (s *gormThreadStore) GetOrCreateThread(ctx context.Context, matchID uint) (*models.MatchThread, error) {
	thread := models.MatchThread{MatchID: matchID}
	err := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&thread).Error
	if err != nil {
		return nil, err
	}
	if thread.ID == 0 {
		err = s.db.WithContext(ctx).Where("match_id = ?", matchID).First(&thread).Error
		if err != nil {
			return nil, err
		}
	}
	return &thread, nil
}

func (s *gormThreadStore) GetPostsByThreadID(ctx context.Context, threadID uint, limit, offset int) ([]models.ThreadPost, int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.ThreadPost{}).
		Where("thread_id = ? AND deleted_at IS NULL", threadID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var posts []models.ThreadPost
	err := s.db.WithContext(ctx).
		Where("thread_id = ? AND deleted_at IS NULL", threadID).
		Preload("User").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	return posts, total, err
}

func (s *gormThreadStore) CreatePost(ctx context.Context, post *models.ThreadPost) error {
	return s.db.WithContext(ctx).Create(post).Error
}

func (s *gormThreadStore) GetPost(ctx context.Context, id uint) (*models.ThreadPost, error) {
	var post models.ThreadPost
	err := s.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *gormThreadStore) UpdatePost(ctx context.Context, id uint, body string) error {
	return s.db.WithContext(ctx).Model(&models.ThreadPost{}).
		Where("id = ?", id).
		Updates(map[string]any{"body": body, "edited": true}).Error
}

func (s *gormThreadStore) SoftDeletePost(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&models.ThreadPost{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}
