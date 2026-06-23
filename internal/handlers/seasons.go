package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSeasons(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	seasons, err := h.seasons.List(ctx)
	if err != nil {
		log.Printf("GetSeasons error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seasons"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, seasons)
}

func (h *Handler) GetSeason(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	season, err := h.seasons.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, season)
}

func (h *Handler) GetActiveSeason(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	season, err := h.seasons.GetActive(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active season found"})
		return
	}
	shortCacheHeaders(c)
	c.JSON(http.StatusOK, season)
}
