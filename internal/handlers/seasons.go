package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSeasons(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	seasons, err := h.seasons.List(ctx)
	if err != nil {
		log.Printf("GetSeasons error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seasons"})
		return
	}
	c.JSON(http.StatusOK, seasons)
}

func (h *Handler) GetSeason(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}
	ctx, cancel := getContext(10)
	defer cancel()

	season, err := h.seasons.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
		return
	}
	c.JSON(http.StatusOK, season)
}

func (h *Handler) GetActiveSeason(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	season, err := h.seasons.GetActive(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active season found"})
		return
	}
	c.JSON(http.StatusOK, season)
}
