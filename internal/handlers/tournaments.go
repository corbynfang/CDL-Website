package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTournaments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	tournaments, err := h.tournaments.ListTournaments(ctx, c.Query("season_id"))
	if err != nil {
		log.Printf("GetTournaments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, tournaments)
}

func (h *Handler) GetTournamentBySlug(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	detail, err := h.tournaments.GetTournamentBySlug(ctx, c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, detail)
}

func (h *Handler) GetTournament(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	tournament, err := h.tournaments.GetTournamentByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, tournament)
}

func (h *Handler) GetTournamentBracket(c *gin.Context) {
	tournamentID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	result, err := h.tournaments.AssembleBracket(ctx, tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetTournamentMatches(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	matches, err := h.tournaments.ListTournamentMatches(ctx, id)
	if err != nil {
		log.Printf("GetTournamentMatches error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch matches"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, matches)
}

func (h *Handler) GetTournamentTeams(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	teams, err := h.tournaments.GetTournamentTeams(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournament teams"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, teams)
}

func (h *Handler) GetTournamentStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	stats, err := h.tournaments.GetTournamentStats(ctx, id)
	if err != nil {
		log.Printf("GetTournamentStats error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	longCacheHeaders(c)
	c.JSON(http.StatusOK, stats)
}
