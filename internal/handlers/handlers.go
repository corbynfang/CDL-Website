package handlers

// handlers.go — Handler struct, constructor, and shared HTTP utilities.
// Business logic lives in internal/services. Handlers parse params, call services, return JSON.
//
// Handler file structure:
//   handlers.go   — this file: Handler struct, New constructor, HTTP helpers
//   seasons.go    — GetSeasons, GetSeason, GetActiveSeason
//   teams.go      — GetTeams, GetTeam, GetTeamPlayers, GetTeamStats
//   franchises.go — GetFranchises, GetFranchise
//   players.go    — GetPlayers, GetPlayer, GetPlayerStats, GetPlayerKDStats,
//                   GetPlayerMatches, GetPlayerFranchiseCareer
//   matches.go    — GetMatch
//   tournaments.go— GetTournaments, GetTournamentBySlug, GetTournament, GetTournamentBracket,
//                   GetTournamentMatches, GetTournamentTeams, GetTournamentStats
//   transfers.go  — GetTransfers
//   stats.go      — GetTopKDPlayers, GetAllPlayersKDStats

import (
	"math"
	"strconv"

	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/corbynfang/CDL-Website/internal/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	players     *services.PlayerService
	teams       *services.TeamService
	seasons     *services.SeasonService
	franchises  *services.FranchiseService
	matches     *services.MatchService
	tournaments *services.TournamentService
	transfers   *services.TransferService
	stats       *services.StatsService
	users       *services.UserService
	threads     *services.ThreadService
}

func New(db *gorm.DB) *Handler {
	playerStore := store.NewGormPlayerStore(db)
	seasonStore := store.NewGormSeasonStore(db)
	teamStore := store.NewGormTeamStore(db)
	franchiseStore := store.NewGormFranchiseStore(db)
	matchStore := store.NewGormMatchStore(db)
	tournamentStore := store.NewGormTournamentStore(db)
	transferStore := store.NewGormTransferStore(db)
	statsStore := store.NewGormStatsStore(db)
	userStore := store.NewGormUserStore(db)
	threadStore := store.NewGormThreadStore(db)

	return &Handler{
		db:          db,
		players:     services.NewPlayerService(playerStore),
		teams:       services.NewTeamService(teamStore, seasonStore),
		seasons:     services.NewSeasonService(seasonStore),
		franchises:  services.NewFranchiseService(franchiseStore),
		matches:     services.NewMatchService(matchStore),
		tournaments: services.NewTournamentService(tournamentStore),
		transfers:   services.NewTransferService(transferStore),
		stats:       services.NewStatsService(statsStore),
		users:       services.NewUserService(userStore),
		threads:     services.NewThreadService(threadStore),
	}
}

func validateID(id string) (int, error) {
	return strconv.Atoi(id)
}

func noCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

func shortCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=60, s-maxage=300")
}

func longCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=300, s-maxage=3600")
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

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

func buildMeta(page, limit, total int) PaginationMeta {
	pages := max(int(math.Ceil(float64(total)/float64(limit))), 1)
	return PaginationMeta{Page: page, Limit: limit, Total: total, TotalPages: pages}
}
