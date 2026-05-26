package handlers

import (
	"log"
	"net/http"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

func GetSeasons(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var seasons []database.Season
	if err := database.DB.WithContext(ctx).Order("start_date DESC").Find(&seasons).Error; err != nil {
		log.Printf("GetSeasons error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seasons"})
		return
	}
	c.JSON(http.StatusOK, seasons)
}

func GetSeason(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var season database.Season
	if err := database.DB.WithContext(ctx).First(&season, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
		return
	}
	c.JSON(http.StatusOK, season)
}

func GetActiveSeason(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var season database.Season
	if err := database.DB.WithContext(ctx).Where("is_active = ?", true).First(&season).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active season found"})
		return
	}
	c.JSON(http.StatusOK, season)
}
