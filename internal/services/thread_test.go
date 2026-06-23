package services

import (
	"context"
	"strings"
	"testing"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockThreadStore struct {
	thread    *models.MatchThread
	posts     []models.ThreadPost
	total     int64
	createErr error
	updateErr error
	deleteErr error
}

func (m *mockThreadStore) FindThread(_ context.Context, matchID uint) (*models.MatchThread, error) {
	if m.thread != nil {
		return m.thread, nil
	}
	return &models.MatchThread{ID: 1, MatchID: matchID}, nil
}

func (m *mockThreadStore) GetOrCreateThread(_ context.Context, matchID uint) (*models.MatchThread, error) {
	if m.thread != nil {
		return m.thread, nil
	}
	return &models.MatchThread{ID: 1, MatchID: matchID}, nil
}

func (m *mockThreadStore) GetPostsByThreadID(_ context.Context, _ uint, _, _ int) ([]models.ThreadPost, int64, error) {
	return m.posts, m.total, nil
}

func (m *mockThreadStore) CreatePost(_ context.Context, post *models.ThreadPost) error {
	if m.createErr != nil {
		return m.createErr
	}
	post.ID = 99
	return nil
}

func (m *mockThreadStore) GetPost(_ context.Context, id uint) (*models.ThreadPost, error) {
	for i := range m.posts {
		if m.posts[i].ID == id {
			return &m.posts[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockThreadStore) UpdatePost(_ context.Context, _ uint, _ string) error { return m.updateErr }
func (m *mockThreadStore) SoftDeletePost(_ context.Context, _ uint) error       { return m.deleteErr }

func TestStripHTML(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain text", "plain text"},
		{"<b>bold</b>", "bold"},
		{"<script>alert(1)</script>XSS", "alert(1)XSS"},
		{"no tags here!", "no tags here!"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, stripHTML(c.in), "input: %q", c.in)
	}
}

func TestThreadService_CreatePost(t *testing.T) {
	ctx := context.Background()

	t.Run("valid body creates post", func(t *testing.T) {
		svc := NewThreadService(&mockThreadStore{})
		post, err := svc.CreatePost(ctx, 1, 7, "Great match!")
		require.NoError(t, err)
		assert.Equal(t, "Great match!", post.Body)
		assert.Equal(t, uint(7), post.UserID)
	})

	t.Run("whitespace-only body returns ErrPostEmpty", func(t *testing.T) {
		svc := NewThreadService(&mockThreadStore{})
		_, err := svc.CreatePost(ctx, 1, 7, "   ")
		assert.ErrorIs(t, err, ErrPostEmpty)
	})

	t.Run("body over 2000 chars returns ErrPostTooLong", func(t *testing.T) {
		svc := NewThreadService(&mockThreadStore{})
		_, err := svc.CreatePost(ctx, 1, 7, strings.Repeat("a", 2001))
		assert.ErrorIs(t, err, ErrPostTooLong)
	})

	t.Run("exactly 2000 chars is accepted", func(t *testing.T) {
		svc := NewThreadService(&mockThreadStore{})
		_, err := svc.CreatePost(ctx, 1, 7, strings.Repeat("a", 2000))
		assert.NoError(t, err)
	})

	t.Run("HTML is stripped from body", func(t *testing.T) {
		svc := NewThreadService(&mockThreadStore{})
		post, err := svc.CreatePost(ctx, 1, 7, "<b>sick play</b>")
		require.NoError(t, err)
		assert.Equal(t, "sick play", post.Body)
	})
}

func TestThreadService_EditPost(t *testing.T) {
	ctx := context.Background()

	t.Run("owner can edit their post", func(t *testing.T) {
		ms := &mockThreadStore{posts: []models.ThreadPost{{ID: 10, UserID: 3, Body: "original"}}}
		svc := NewThreadService(ms)
		assert.NoError(t, svc.EditPost(ctx, 10, 3, "updated"))
	})

	t.Run("non-owner gets ErrNotOwner", func(t *testing.T) {
		ms := &mockThreadStore{posts: []models.ThreadPost{{ID: 10, UserID: 3, Body: "original"}}}
		svc := NewThreadService(ms)
		assert.ErrorIs(t, svc.EditPost(ctx, 10, 99, "hacked"), ErrNotOwner)
	})

	t.Run("empty edit body returns ErrPostEmpty", func(t *testing.T) {
		ms := &mockThreadStore{posts: []models.ThreadPost{{ID: 10, UserID: 3, Body: "original"}}}
		svc := NewThreadService(ms)
		assert.ErrorIs(t, svc.EditPost(ctx, 10, 3, "  "), ErrPostEmpty)
	})

	t.Run("post not found returns an error", func(t *testing.T) {
		svc := NewThreadService(&mockThreadStore{})
		assert.Error(t, svc.EditPost(ctx, 999, 1, "body"))
	})
}

func TestThreadService_DeletePost(t *testing.T) {
	ctx := context.Background()

	t.Run("owner can delete their post", func(t *testing.T) {
		ms := &mockThreadStore{posts: []models.ThreadPost{{ID: 5, UserID: 2}}}
		svc := NewThreadService(ms)
		assert.NoError(t, svc.DeletePost(ctx, 5, 2))
	})

	t.Run("non-owner gets ErrNotOwner", func(t *testing.T) {
		ms := &mockThreadStore{posts: []models.ThreadPost{{ID: 5, UserID: 2}}}
		svc := NewThreadService(ms)
		assert.ErrorIs(t, svc.DeletePost(ctx, 5, 99), ErrNotOwner)
	})
}

func TestThreadService_GetThread(t *testing.T) {
	ctx := context.Background()
	ms := &mockThreadStore{
		posts: []models.ThreadPost{{ID: 1, Body: "first"}, {ID: 2, Body: "second"}},
		total: 2,
	}
	svc := NewThreadService(ms)
	posts, total, threadID, err := svc.GetThread(ctx, 42, 1, 25)
	require.NoError(t, err)
	assert.EqualValues(t, 2, total)
	assert.Len(t, posts, 2)
	assert.EqualValues(t, 1, threadID)
}
