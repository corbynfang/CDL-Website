package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMatch(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	detail, err := h.matches.GetMatchDetail(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *Handler) GetTournaments(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	tournaments, err := h.matches.ListTournaments(ctx, c.Query("season_id"))
	if err != nil {
		log.Printf("GetTournaments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

func (h *Handler) GetTournamentBySlug(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	detail, err := h.matches.GetTournamentBySlug(ctx, c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *Handler) GetTournament(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := getContext(10)
	defer cancel()

	tournament, err := h.matches.GetTournamentByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

func (h *Handler) GetTournamentBracket(c *gin.Context) {
	tournamentID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	result, err := h.matches.AssembleBracket(ctx, tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetTournamentMatches(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	matches, err := h.matches.ListTournamentMatches(ctx, id)
	if err != nil {
		log.Printf("GetTournamentMatches error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch matches"})
		return
	}
	c.JSON(http.StatusOK, matches)
}

func (h *Handler) GetTournamentTeams(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	teams, err := h.matches.GetTournamentTeams(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournament teams"})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *Handler) GetTournamentStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	stats, err := h.matches.GetTournamentStats(ctx, id)
	if err != nil {
		log.Printf("GetTournamentStats error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
