package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTopKDPlayers(c *gin.Context) {
	noCacheHeaders(c)

	limit := 25
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	rows, err := h.stats.GetTopKD(ctx, limit)
	if err != nil {
		log.Printf("GetTopKDPlayers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"players":   rows,
		"count":     len(rows),
	})
}

func (h *Handler) GetAllPlayersKDStats(c *gin.Context) {
	noCacheHeaders(c)

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	rows, err := h.stats.GetAllKD(ctx, limit, c.Query("season_id"))
	if err != nil {
		log.Printf("GetAllPlayersKDStats error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	type enrichedRow struct {
		services.PlayerKDRow
		SeasonKDPlusMinus float64 `json:"season_kd_plus_minus"`
	}
	enriched := make([]enrichedRow, len(rows))
	for i, row := range rows {
		enriched[i] = enrichedRow{PlayerKDRow: row, SeasonKDPlusMinus: row.SeasonKD - 1.0}
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"players":   enriched,
		"count":     len(enriched),
	})
}
