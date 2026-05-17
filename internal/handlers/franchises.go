package handlers

// franchises.go — handlers for the /franchises endpoints.
// A Franchise is a CDL slot that persists across team rebrands.
// These endpoints power the "Franchise History" sidebar on team detail pages
// and the franchise career stats on player profiles.

import (
	"log"
	"net/http"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

// GetFranchises returns all CDL franchise slots with their current name.
func GetFranchises(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var franchises []database.Franchise
	if err := database.DB.WithContext(ctx).
		Where("franchise_key != ''").
		Order("name ASC").
		Find(&franchises).Error; err != nil {
		log.Printf("GetFranchises error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch franchises"})
		return
	}
	c.JSON(http.StatusOK, franchises)
}

// GetFranchise returns a single franchise with its full era history.
// The "eras" array contains every Team row linked to this franchise, ordered
// chronologically so the frontend can render a timeline (oldest → newest).
func GetFranchise(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing franchise key"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var franchise database.Franchise
	if err := database.DB.WithContext(ctx).
		Where("franchise_key = ?", key).
		First(&franchise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Franchise not found"})
		return
	}

	var teams []database.Team
	database.DB.WithContext(ctx).
		Where("franchise_id = ?", franchise.ID).
		Order("valid_from ASC").
		Find(&teams)

	c.JSON(http.StatusOK, gin.H{
		"franchise": franchise,
		"eras":      teams,
	})
}
