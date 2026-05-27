package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPlayers(c *gin.Context) {
	page, limit, offset := parsePagination(c)

	search := c.Query("search")
	if len(search) > 50 {
		search = search[:50]
	}

	ctx, cancel := getContext(10)
	defer cancel()

	players, total, err := h.players.List(ctx, search, limit, offset)
	if err != nil {
		log.Printf("GetPlayers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       players,
		"pagination": buildMeta(page, limit, int(total)),
	})
}

func (h *Handler) GetPlayer(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}
	ctx, cancel := getContext(10)
	defer cancel()

	player, err := h.players.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}
	c.JSON(http.StatusOK, player)
}

func (h *Handler) GetPlayerStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}
	ctx, cancel := getContext(10)
	defer cancel()

	stats, err := h.players.GetMatchStats(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetPlayerKDStats(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	result, err := h.players.GetKDStats(ctx, playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetPlayerMatches(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	result, err := h.players.GetMatchHistory(ctx, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player matches"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetPlayerFranchiseCareer(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}
	ctx, cancel := getContext(15)
	defer cancel()

	result, err := h.players.GetFranchiseCareer(ctx, playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}
