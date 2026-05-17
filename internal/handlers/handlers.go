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
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
