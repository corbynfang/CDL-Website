package services

import (
	"context"
	"errors"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
	"gorm.io/gorm"
)

type UserService struct {
	store store.UserStore
}

func NewUserService(s store.UserStore) *UserService {
	return &UserService{store: s}
}

func (us *UserService) SyncProfile(ctx context.Context, supabaseUID, username string) (*models.User, error) {
	user, err := us.store.GetBySupabaseUID(ctx, supabaseUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	newUser := &models.User{SupabaseUID: supabaseUID, Username: username}
	if err := us.store.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (us *UserService) GetBySupabaseUID(ctx context.Context, uid string) (*models.User, error) {
	return us.store.GetBySupabaseUID(ctx, uid)
}

func (us *UserService) Delete(ctx context.Context, id uint) error {
	return us.store.Delete(ctx, id)
}
