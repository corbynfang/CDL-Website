package handlers

// transfers.go — handler for the /transfers endpoint.
// Returns player transfer history with optional filters for season, team, and player.
// The raw_from_team_name / raw_to_team_name fields are always populated so the frontend
// can display "Free Agent" correctly even when the FK is null.

import (
	"log"
	"net/http"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

func GetTransfers(c *gin.Context) {
	noCacheHeaders(c)

	ctx, cancel := getContext(15)
	defer cancel()

	query := database.DB.WithContext(ctx).
		Preload("Player").
		Preload("FromTeam").
		Preload("ToTeam")

	if season := c.Query("season"); season != "" {
		query = query.Where("season = ?", season)
	}
	if gameCode := c.Query("game_code"); gameCode != "" {
		query = query.Where("game_code = ?", gameCode)
	}
	if teamID := c.Query("team_id"); teamID != "" {
		query = query.Where("from_team_id = ? OR to_team_id = ?", teamID, teamID)
	}
	if playerID := c.Query("player_id"); playerID != "" {
		query = query.Where("player_id = ?", playerID)
	}

	var transfers []database.PlayerTransfer
	if err := query.Order("transfer_date DESC").Find(&transfers).Error; err != nil {
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
