package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetThread(c *gin.Context) {
	matchID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	page, limit, _ := parsePagination(c)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	posts, total, threadID, err := h.threads.GetThread(ctx, uint(matchID), page, limit)
	if err != nil {
		log.Printf("GetThread error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch thread"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"thread_id":  threadID,
		"data":       posts,
		"pagination": buildMeta(page, limit, int(total)),
	})
}

func (h *Handler) CreatePost(c *gin.Context) {
	matchID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	uid := c.GetString("supabase_uid")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	user, err := h.users.GetBySupabaseUID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "complete profile setup first"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		}
		return
	}

	_, _, threadID, err := h.threads.GetThread(ctx, uint(matchID), 1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load thread"})
		return
	}

	var body struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body is required"})
		return
	}

	post, err := h.threads.CreatePost(ctx, threadID, user.ID, body.Body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPostEmpty):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPostTooLong):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		}
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *Handler) EditPost(c *gin.Context) {
	postID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	uid := c.GetString("supabase_uid")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	user, err := h.users.GetBySupabaseUID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "complete profile setup first"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		}
		return
	}

	var body struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body is required"})
		return
	}

	err = h.threads.EditPost(ctx, uint(postID), user.ID, body.Body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPostEmpty), errors.Is(err, services.ErrPostTooLong):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to edit post"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}

func (h *Handler) DeletePost(c *gin.Context) {
	postID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	uid := c.GetString("supabase_uid")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	user, err := h.users.GetBySupabaseUID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "complete profile setup first"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		}
		return
	}

	err = h.threads.DeletePost(ctx, uint(postID), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
