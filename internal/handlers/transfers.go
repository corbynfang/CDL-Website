package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTransfers(c *gin.Context) {
	noCacheHeaders(c)

	ctx, cancel := getContext(15)
	defer cancel()

	transfers, err := h.transfers.List(ctx, services.TransferFilters{
		Season:   c.Query("season"),
		GameCode: c.Query("game_code"),
		TeamID:   c.Query("team_id"),
		PlayerID: c.Query("player_id"),
	})
	if err != nil {
		log.Printf("GetTransfers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transfers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"transfers": transfers,
		"count":     len(transfers),
	})
}
