package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTeams(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	teams, err := h.teams.List(ctx, c.Query("season_id"), c.Query("scope"))
	if errors.Is(err, services.ErrInvalidSeason) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season"})
		return
	}
	if err != nil {
		log.Printf("GetTeams error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.Header("Cache-Control", "public, max-age=60, s-maxage=300")
	c.JSON(http.StatusOK, teams)
}

func (h *Handler) GetTeam(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	team, err := h.teams.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *Handler) GetTeamPlayers(c *gin.Context) {
	teamID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	seasonID := c.Query("season_id")
	if seasonID != "" {
		if _, err := strconv.Atoi(seasonID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season_id"})
			return
		}
	}

	scope := c.DefaultQuery("scope", "current")
	fetch := h.teams.GetCurrentRoster
	switch scope {
	case "current":
	case "all":
		fetch = h.teams.GetPlayers
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scope"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var players []models.Player
	players, err = fetch(ctx, teamID, seasonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team players"})
		return
	}
	c.JSON(http.StatusOK, players)
}

func (h *Handler) GetTeamStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	stats, err := h.teams.GetStats(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
