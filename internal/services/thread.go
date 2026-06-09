package services

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

var ErrPostTooLong = errors.New("post body exceeds 2000 characters")
var ErrPostEmpty = errors.New("post body cannot be empty")
var ErrNotOwner = errors.New("you can only edit your own posts")

type ThreadService struct {
	store store.ThreadStore
}

func NewThreadService(s store.ThreadStore) *ThreadService {
	return &ThreadService{store: s}
}

func (ts *ThreadService) GetThread(ctx context.Context, matchID uint, page, limit int) ([]models.ThreadPost, int64, uint, error) {
	thread, err := ts.store.GetOrCreateThread(ctx, matchID)
	if err != nil {
		return nil, 0, 0, err
	}
	offset := (page - 1) * limit
	posts, total, err := ts.store.GetPostsByThreadID(ctx, thread.ID, limit, offset)
	return posts, total, thread.ID, err
}

func (ts *ThreadService) CreatePost(ctx context.Context, threadID, userID uint, body string) (*models.ThreadPost, error) {
	body = strings.TrimSpace(stripHTML(body))
	if body == "" {
		return nil, ErrPostEmpty
	}
	if utf8.RuneCountInString(body) > 2000 {
		return nil, ErrPostTooLong
	}
	post := &models.ThreadPost{ThreadID: threadID, UserID: userID, Body: body}
	if err := ts.store.CreatePost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (ts *ThreadService) EditPost(ctx context.Context, postID, userID uint, body string) error {
	post, err := ts.store.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return ErrNotOwner
	}
	body = strings.TrimSpace(stripHTML(body))
	if body == "" {
		return ErrPostEmpty
	}
	if utf8.RuneCountInString(body) > 2000 {
		return ErrPostTooLong
	}
	return ts.store.UpdatePost(ctx, postID, body)
}

func (ts *ThreadService) DeletePost(ctx context.Context, postID, userID uint) error {
	post, err := ts.store.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return ErrNotOwner
	}
	return ts.store.SoftDeletePost(ctx, postID)
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}
