package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) SyncProfile(c *gin.Context) {
	uid := c.GetString("supabase_uid")

	var body struct {
		Username string `json:"username" binding:"required,min=3,max=30"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required (3-30 characters)"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	user, err := h.users.SyncProfile(ctx, uid, body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync profile"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetMe(c *gin.Context) {
	uid := c.GetString("supabase_uid")

	ctx, cancel := getContext(10)
	defer cancel()

	user, err := h.users.GetBySupabaseUID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteMe(c *gin.Context) {
	uid := c.GetString("supabase_uid")

	ctx, cancel := getContext(10)
	defer cancel()

	user, err := h.users.GetBySupabaseUID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		}
		return
	}

	if err := h.users.Delete(ctx, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}
