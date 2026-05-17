package handlers

// handlers.go — shared utilities used across all handler files.
// This is the only file in the handlers package that doesn't own specific routes.
//
// Handler file structure:
//   handlers.go   — this file: shared helpers (validateID, getContext, etc.)
//   seasons.go    — GetSeasons, GetSeason, GetActiveSeason
//   teams.go      — GetTeams, GetTeam, GetTeamPlayers, GetTeamStats
//   franchises.go — GetFranchises, GetFranchise
//   players.go    — GetPlayers, GetPlayer, GetPlayerStats, GetPlayerKDStats,
//                   GetPlayerMatches, GetPlayerFranchiseCareer
//   matches.go    — GetMatch, GetTournaments, GetTournament, GetTournamentBracket
//   transfers.go  — GetTransfers
//   stats.go      — GetTopKDPlayers, GetTopKDPlayersNew, GetAllPlayersKDStats, GetDatabaseValidation

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func validateID(id string) (int, error) {
	return strconv.Atoi(id)
}

func getContext(seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// noCacheHeaders prevents browsers from caching API responses.
// Used on leaderboard endpoints where staleness is most visible.
func noCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

func calculateKD(kills, deaths int) float64 {
	if deaths == 0 {
		return 0
	}
	return float64(kills) / float64(deaths)
}

// PaginationMeta is included in every paginated response so the frontend
// knows how many total pages exist and which page it's on.
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// parsePagination reads ?page and ?limit from the request, applies safe
// defaults, caps the limit, and returns the offset + a partial meta struct.
// The caller must set meta.Total before sending the response.
func parsePagination(c *gin.Context) (page, limit, offset int) {
	page = 1
	limit = 25

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset = (page - 1) * limit
	return
}

// buildMeta constructs the final PaginationMeta given a total row count.
func buildMeta(page, limit, total int) PaginationMeta {
	pages := int(math.Ceil(float64(total) / float64(limit)))
	if pages < 1 {
		pages = 1
	}
	return PaginationMeta{Page: page, Limit: limit, Total: total, TotalPages: pages}
}

// applyPagination applies LIMIT/OFFSET to a GORM query chain.
func applyPagination(q *gorm.DB, limit, offset int) *gorm.DB {
	return q.Limit(limit).Offset(offset)
}
