package services

import (
	"context"
	"errors"
	"testing"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockUserStore struct {
	user      *models.User
	getErr    error
	createErr error
	deleteErr error
	created   []*models.User
	deleted   []uint
}

func (m *mockUserStore) GetBySupabaseUID(_ context.Context, uid string) (*models.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.user != nil && m.user.SupabaseUID == uid {
		return m.user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserStore) Create(_ context.Context, user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = 42
	m.created = append(m.created, user)
	return nil
}

func (m *mockUserStore) Delete(_ context.Context, id uint) error {
	m.deleted = append(m.deleted, id)
	return m.deleteErr
}

func TestUserService_SyncProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("creates new user when none exists", func(t *testing.T) {
		ms := &mockUserStore{}
		svc := NewUserService(ms)
		user, err := svc.SyncProfile(ctx, "uid-abc", "Corbyn")
		require.NoError(t, err)
		assert.Equal(t, "uid-abc", user.SupabaseUID)
		assert.Equal(t, "Corbyn", user.Username)
		assert.Len(t, ms.created, 1)
	})

	t.Run("returns existing user without creating a duplicate", func(t *testing.T) {
		existing := &models.User{ID: 7, SupabaseUID: "uid-abc", Username: "Corbyn"}
		ms := &mockUserStore{user: existing}
		svc := NewUserService(ms)
		user, err := svc.SyncProfile(ctx, "uid-abc", "Corbyn")
		require.NoError(t, err)
		assert.Equal(t, uint(7), user.ID)
		assert.Empty(t, ms.created)
	})

	t.Run("propagates unexpected store errors", func(t *testing.T) {
		ms := &mockUserStore{getErr: errors.New("db down")}
		svc := NewUserService(ms)
		_, err := svc.SyncProfile(ctx, "uid-abc", "Corbyn")
		assert.Error(t, err)
	})
}

func TestUserService_Delete(t *testing.T) {
	ctx := context.Background()
	ms := &mockUserStore{}
	svc := NewUserService(ms)
	require.NoError(t, svc.Delete(ctx, 7))
	assert.Equal(t, []uint{7}, ms.deleted)
}
