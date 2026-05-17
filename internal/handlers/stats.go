package handlers

// stats.go — leaderboard and aggregate stats handlers.
// These all query player_tournament_stats (season-level aggregates) rather than
// player_match_stats (per-match), so they reflect full-season performance.

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

// GetTopKDPlayers returns the top N players by career K/D across all seasons.
func GetTopKDPlayers(c *gin.Context) {
	noCacheHeaders(c)

	ctx, cancel := getContext(15)
	defer cancel()

	limit := 25
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	type PlayerAggregated struct {
		PlayerID      uint
		Gamertag      string
		AvatarURL     string
		TeamAbbr      string
		SeasonKills   int
		SeasonDeaths  int
		SeasonAssists int
	}

	var rows []PlayerAggregated
	if err := database.DB.WithContext(ctx).
		Table("player_tournament_stats pts").
		Select(`pts.player_id, MAX(p.gamertag) as gamertag,
			COALESCE(MAX(p.avatar_url), '') as avatar_url,
			COALESCE(MAX(t.abbreviation), '') as team_abbr,
			SUM(pts.total_kills) as season_kills,
			SUM(pts.total_deaths) as season_deaths,
			SUM(pts.total_assists) as season_assists`).
		Joins("JOIN players p ON pts.player_id = p.id").
		Joins("LEFT JOIN teams t ON pts.team_id = t.id").
		Group("pts.player_id").
		Having("SUM(pts.total_deaths) > 0").
		Order("(SUM(pts.total_kills)::decimal / NULLIF(SUM(pts.total_deaths), 0)) DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		log.Printf("GetTopKDPlayers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}

	players := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		players = append(players, gin.H{
			"player_id":      row.PlayerID,
			"gamertag":       row.Gamertag,
			"avatar_url":     row.AvatarURL,
			"team_abbr":      row.TeamAbbr,
			"season_kills":   row.SeasonKills,
			"season_deaths":  row.SeasonDeaths,
			"season_assists": row.SeasonAssists,
			"season_kd":      calculateKD(row.SeasonKills, row.SeasonDeaths),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"players":   players,
		"count":     len(players),
	})
}

// GetTopKDPlayersNew is an alias kept for backward compatibility.
func GetTopKDPlayersNew(c *gin.Context) {
	GetTopKDPlayers(c)
}

// GetAllPlayersKDStats returns K/D stats for all players, optionally filtered by season.
func GetAllPlayersKDStats(c *gin.Context) {
	noCacheHeaders(c)

	ctx, cancel := getContext(30)
	defer cancel()

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 500 {
			limit = parsed
		}
	}

	seasonID := c.Query("season_id")

	type PlayerAggregated struct {
		PlayerID      uint
		Gamertag      string
		AvatarURL     string
		TeamAbbr      string
		SeasonKills   int
		SeasonDeaths  int
		SeasonAssists int
	}

	query := database.DB.WithContext(ctx).
		Table("player_tournament_stats pts").
		Select(`pts.player_id, MAX(p.gamertag) as gamertag,
			COALESCE(MAX(p.avatar_url), '') as avatar_url,
			COALESCE(MAX(t.abbreviation), '') as team_abbr,
			SUM(pts.total_kills) as season_kills,
			SUM(pts.total_deaths) as season_deaths,
			SUM(pts.total_assists) as season_assists`).
		Joins("JOIN players p ON pts.player_id = p.id").
		Joins("LEFT JOIN teams t ON pts.team_id = t.id").
		Joins("JOIN tournaments tour ON pts.tournament_id = tour.id")

	if seasonID != "" {
		query = query.Where("tour.season_id = ?", seasonID)
	}

	var rows []PlayerAggregated
	if err := query.
		Group("pts.player_id").
		Having("SUM(pts.total_kills) > 0 OR SUM(pts.total_deaths) > 0").
		Order("(CASE WHEN SUM(pts.total_deaths) > 0 THEN SUM(pts.total_kills)::decimal / SUM(pts.total_deaths) ELSE 0 END) DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		log.Printf("GetAllPlayersKDStats error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	players := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		kd := calculateKD(row.SeasonKills, row.SeasonDeaths)
		players = append(players, gin.H{
			"player_id":            row.PlayerID,
			"gamertag":             row.Gamertag,
			"avatar_url":           row.AvatarURL,
			"team_abbr":            row.TeamAbbr,
			"season_kills":         row.SeasonKills,
			"season_deaths":        row.SeasonDeaths,
			"season_assists":       row.SeasonAssists,
			"season_kd":            kd,
			"season_kd_plus_minus": kd - 1.0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"players":   players,
		"count":     len(players),
	})
}

// GetDatabaseValidation returns a health check with row counts for key tables.
func GetDatabaseValidation(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var playerCount, teamCount, matchCount, tournamentCount, mapStatsCount int64
	database.DB.WithContext(ctx).Model(&database.Player{}).Count(&playerCount)
	database.DB.WithContext(ctx).Model(&database.Team{}).Count(&teamCount)
	database.DB.WithContext(ctx).Model(&database.Match{}).Count(&matchCount)
	database.DB.WithContext(ctx).Model(&database.Tournament{}).Count(&tournamentCount)
	database.DB.WithContext(ctx).Model(&database.PlayerMapStats{}).Count(&mapStatsCount)

	c.JSON(http.StatusOK, gin.H{
		"status":           "healthy",
		"timestamp":        time.Now().Unix(),
		"players":          playerCount,
		"teams":            teamCount,
		"matches":          matchCount,
		"tournaments":      tournamentCount,
		"player_map_stats": mapStatsCount,
	})
}
