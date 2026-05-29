package handlers

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg gin.IRoutes, h *Handler) {
	rg.GET("/seasons", h.GetSeasons)
	rg.GET("/seasons/:id", h.GetSeason)
	rg.GET("/seasons/active", h.GetActiveSeason)

	rg.GET("/teams", h.GetTeams)
	rg.GET("/teams/:id", h.GetTeam)
	rg.GET("/teams/:id/players", h.GetTeamPlayers)
	rg.GET("/teams/:id/stats", h.GetTeamStats)

	rg.GET("/players", h.GetPlayers)
	rg.GET("/players/:id", h.GetPlayer)
	rg.GET("/players/:id/stats", h.GetPlayerStats)
	rg.GET("/players/:id/kd", h.GetPlayerKDStats)
	rg.GET("/players/:id/matches", h.GetPlayerMatches)
	rg.GET("/players/:id/franchise-career", h.GetPlayerFranchiseCareer)
	rg.GET("/players/top-kd", h.GetTopKDPlayers)

	rg.GET("/stats/all-kd-by-tournament", h.GetAllPlayersKDStats)

	rg.GET("/matches/:id", h.GetMatch)

	rg.GET("/franchises", h.GetFranchises)
	rg.GET("/franchises/:key", h.GetFranchise)

	rg.GET("/tournaments", h.GetTournaments)
	rg.GET("/tournaments/slug/:slug", h.GetTournamentBySlug)
	rg.GET("/tournaments/:id", h.GetTournament)
	rg.GET("/tournaments/:id/bracket", h.GetTournamentBracket)
	rg.GET("/tournaments/:id/matches", h.GetTournamentMatches)
	rg.GET("/tournaments/:id/teams", h.GetTournamentTeams)
	rg.GET("/tournaments/:id/stats", h.GetTournamentStats)

	rg.GET("/transfers", h.GetTransfers)
}
