package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// validateID validates and converts string ID to int
func validateID(id string) (int, error) {
	return strconv.Atoi(id)
}

// getContext creates a context with timeout for database operations
func getContext(seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// noCacheHeaders adds cache-busting headers to response
func noCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// calculateKD safely calculates K/D ratio
func calculateKD(kills, deaths int) float64 {
	if deaths == 0 {
		return 0
	}
	return float64(kills) / float64(deaths)
}

// =============================================================================
// TEAM HANDLERS
// =============================================================================

// GetTeams returns all active teams
func GetTeams(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var teams []database.Team
	if err := database.DB.WithContext(ctx).Where("is_active = ?", true).Find(&teams).Error; err != nil {
		log.Printf("GetTeams error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeam returns a single team by ID
func GetTeam(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var team database.Team
	if err := database.DB.WithContext(ctx).First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetTeamPlayers returns all players for a team
func GetTeamPlayers(c *gin.Context) {
	teamID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var players []database.Player
	if err := database.DB.WithContext(ctx).
		Joins("JOIN team_rosters ON players.id = team_rosters.player_id").
		Where("team_rosters.team_id = ? AND team_rosters.end_date IS NULL", teamID).
		Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetTeamStats returns tournament stats for a team
func GetTeamStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var stats []database.TeamTournamentStats
	if err := database.DB.WithContext(ctx).
		Where("team_id = ?", id).
		Preload("Tournament").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// =============================================================================
// PLAYER HANDLERS
// =============================================================================

// GetPlayers returns all players
func GetPlayers(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var players []database.Player
	if err := database.DB.WithContext(ctx).Find(&players).Error; err != nil {
		log.Printf("GetPlayers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetPlayer returns a single player by ID
func GetPlayer(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var player database.Player
	if err := database.DB.WithContext(ctx).First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

// GetPlayerStats returns match stats for a player
func GetPlayerStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var stats []database.PlayerMatchStats
	if err := database.DB.WithContext(ctx).
		Where("player_id = ?", id).
		Preload("Match").
		Preload("Team").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPlayerKDStats returns K/D statistics for a player across tournaments
func GetPlayerKDStats(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	// Get player info
	var player database.Player
	if err := database.DB.WithContext(ctx).First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	// Get tournament stats for this player
	var tournamentStats []database.PlayerTournamentStats
	if err := database.DB.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("Tournament").
		Order("tournament_id DESC").
		Find(&tournamentStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	// Aggregate stats
	var totalKills, totalDeaths, totalAssists int
	var totalHpKills, totalHpDeaths, totalSndKills, totalSndDeaths, totalCtlKills, totalCtlDeaths int

	tournamentList := make([]gin.H, 0, len(tournamentStats))

	for _, stat := range tournamentStats {
		totalKills += stat.TotalKills
		totalDeaths += stat.TotalDeaths
		totalAssists += stat.TotalAssists
		totalHpKills += stat.HpKills
		totalHpDeaths += stat.HpDeaths
		totalSndKills += stat.SndKills
		totalSndDeaths += stat.SndDeaths
		totalCtlKills += stat.ControlKills
		totalCtlDeaths += stat.ControlDeaths

		tournamentList = append(tournamentList, gin.H{
			"tournament_id":   stat.TournamentID,
			"tournament_name": stat.Tournament.Name,
			"kills":           stat.TotalKills,
			"deaths":          stat.TotalDeaths,
			"assists":         stat.TotalAssists,
			"kd_ratio":        calculateKD(stat.TotalKills, stat.TotalDeaths),
			"maps_played":     stat.OverallMaps,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":        playerID,
		"gamertag":         player.Gamertag,
		"avatar_url":       player.AvatarURL,
		"total_kills":      totalKills,
		"total_deaths":     totalDeaths,
		"total_assists":    totalAssists,
		"avg_kd":           calculateKD(totalKills, totalDeaths),
		"hp_kd_ratio":      calculateKD(totalHpKills, totalHpDeaths),
		"snd_kd_ratio":     calculateKD(totalSndKills, totalSndDeaths),
		"control_kd_ratio": calculateKD(totalCtlKills, totalCtlDeaths),
		"tournament_stats": tournamentList,
	})
}

// GetPlayerMatches returns match history for a player grouped by event
func GetPlayerMatches(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	// Get player match stats with all related data
	var matchStats []database.PlayerMatchStats
	if err := database.DB.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("Match").
		Preload("Match.Tournament").
		Preload("Match.Team1").
		Preload("Match.Team2").
		Preload("Team").
		Order("match_id DESC").
		Limit(100).
		Find(&matchStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player matches"})
		return
	}

	// Group matches by tournament (event)
	eventsMap := make(map[uint]gin.H)

	for _, stat := range matchStats {
		match := stat.Match
		tournamentID := match.TournamentID

		// Initialize event if not exists
		if eventsMap[tournamentID] == nil {
			eventsMap[tournamentID] = gin.H{
				"event":         match.Tournament.Name,
				"year":          match.Tournament.StartDate.Year(),
				"tournament_id": tournamentID,
				"matches":       []gin.H{},
			}
		}

		// Determine opponent and result
		var opponent, opponentAbbr, result string
		playerTeamID := stat.TeamID

		if playerTeamID == match.Team1ID {
			opponent = match.Team2.Name
			opponentAbbr = match.Team2.Abbreviation
			if match.Team1Score > match.Team2Score {
				result = "W"
			} else {
				result = "L"
			}
		} else {
			opponent = match.Team1.Name
			opponentAbbr = match.Team1.Abbreviation
			if match.Team2Score > match.Team1Score {
				result = "W"
			} else {
				result = "L"
			}
		}

		resultScore := result + " " + strconv.Itoa(match.Team1Score) + ":" + strconv.Itoa(match.Team2Score)

		matchData := gin.H{
			"date":          match.MatchDate.Format(time.RFC3339),
			"opponent":      opponent,
			"opponent_abbr": opponentAbbr,
			"result":        resultScore,
			"kd":            stat.KDRatio,
			"kills":         stat.TotalKills,
			"deaths":        stat.TotalDeaths,
		}

		event := eventsMap[tournamentID]
		matchesList := event["matches"].([]gin.H)
		event["matches"] = append(matchesList, matchData)
	}

	// Convert map to sorted slice
	events := make([]gin.H, 0, len(eventsMap))
	for _, event := range eventsMap {
		events = append(events, event)
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"events":    events,
		"total":     len(matchStats),
	})
}

// =============================================================================
// STATS HANDLERS
// =============================================================================

// GetTopKDPlayers returns top players by K/D ratio
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

	// Get aggregated stats from player_tournament_stats
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
		Select(`
			pts.player_id,
			MAX(p.gamertag) as gamertag,
			COALESCE(MAX(p.avatar_url), '') as avatar_url,
			COALESCE(MAX(t.abbreviation), '') as team_abbr,
			SUM(pts.total_kills) as season_kills,
			SUM(pts.total_deaths) as season_deaths,
			SUM(pts.total_assists) as season_assists
		`).
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

// GetTopKDPlayersNew is an alias for GetTopKDPlayers
func GetTopKDPlayersNew(c *gin.Context) {
	GetTopKDPlayers(c)
}

// GetAllPlayersKDStats returns K/D stats for all players
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
		Select(`
			pts.player_id,
			MAX(p.gamertag) as gamertag,
			COALESCE(MAX(p.avatar_url), '') as avatar_url,
			COALESCE(MAX(t.abbreviation), '') as team_abbr,
			SUM(pts.total_kills) as season_kills,
			SUM(pts.total_deaths) as season_deaths,
			SUM(pts.total_assists) as season_assists
		`).
		Joins("JOIN players p ON pts.player_id = p.id").
		Joins("LEFT JOIN teams t ON pts.team_id = t.id").
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
		seasonKD := calculateKD(row.SeasonKills, row.SeasonDeaths)
		players = append(players, gin.H{
			"player_id":            row.PlayerID,
			"gamertag":             row.Gamertag,
			"avatar_url":           row.AvatarURL,
			"team_abbr":            row.TeamAbbr,
			"season_kills":         row.SeasonKills,
			"season_deaths":        row.SeasonDeaths,
			"season_assists":       row.SeasonAssists,
			"season_kd":            seasonKD,
			"season_kd_plus_minus": seasonKD - 1.0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"players":   players,
		"count":     len(players),
	})
}

// =============================================================================
// TOURNAMENT HANDLERS
// =============================================================================

// GetTournaments returns all tournaments
func GetTournaments(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var tournaments []database.Tournament
	if err := database.DB.WithContext(ctx).
		Preload("Season").
		Order("start_date DESC").
		Find(&tournaments).Error; err != nil {
		log.Printf("GetTournaments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// GetTournament returns a single tournament by ID
func GetTournament(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var tournament database.Tournament
	if err := database.DB.WithContext(ctx).
		Preload("Season").
		First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// GetTournamentBracket returns bracket data for a tournament
func GetTournamentBracket(c *gin.Context) {
	tournamentID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	// Get tournament info
	var tournament database.Tournament
	if err := database.DB.WithContext(ctx).First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get all bracket matches
	var matches []database.Match
	if err := database.DB.WithContext(ctx).
		Where("tournament_id = ? AND bracket_round != ''", tournamentID).
		Preload("Team1").
		Preload("Team2").
		Order("bracket_round, bracket_position").
		Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bracket data"})
		return
	}

	// Organize by bracket round
	bracket := map[string][]gin.H{
		"winners_r1":     {},
		"winners_r2":     {},
		"winners_finals": {},
		"elim_r1":        {},
		"elim_r2":        {},
		"elim_r3":        {},
		"elim_finals":    {},
		"grand_finals":   {},
	}

	for _, match := range matches {
		matchData := gin.H{
			"id":               match.ID,
			"team1_id":         match.Team1ID,
			"team2_id":         match.Team2ID,
			"team1_name":       match.Team1.Name,
			"team1_abbr":       match.Team1.Abbreviation,
			"team1_logo":       match.Team1.LogoURL,
			"team2_name":       match.Team2.Name,
			"team2_abbr":       match.Team2.Abbreviation,
			"team2_logo":       match.Team2.LogoURL,
			"team1_score":      match.Team1Score,
			"team2_score":      match.Team2Score,
			"winner_id":        match.WinnerID,
			"bracket_position": match.BracketPosition,
			"match_date":       match.MatchDate,
		}

		if _, exists := bracket[match.BracketRound]; exists {
			bracket[match.BracketRound] = append(bracket[match.BracketRound], matchData)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tournament_id":   tournamentID,
		"tournament_name": tournament.Name,
		"bracket":         bracket,
		"total_matches":   len(matches),
	})
}

// =============================================================================
// TRANSFER HANDLERS
// =============================================================================

// GetTransfers returns player transfers
func GetTransfers(c *gin.Context) {
	noCacheHeaders(c)

	ctx, cancel := getContext(15)
	defer cancel()

	query := database.DB.WithContext(ctx).
		Preload("Player").
		Preload("FromTeam").
		Preload("ToTeam")

	// Optional filters
	if season := c.Query("season"); season != "" {
		query = query.Where("season = ?", season)
	}
	if teamID := c.Query("team_id"); teamID != "" {
		query = query.Where("from_team_id = ? OR to_team_id = ?", teamID, teamID)
	}
	if playerID := c.Query("player_id"); playerID != "" {
		query = query.Where("player_id = ?", playerID)
	}

	var transfers []database.PlayerTransfer
	if err := query.Order("transfer_date DESC").Find(&transfers).Error; err != nil {
		log.Printf("GetTransfers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transfers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"transfers": transfers,
		"count":     len(transfers),
	})
}

// =============================================================================
// DEBUG HANDLERS
// =============================================================================

// GetDatabaseValidation returns database health check
func GetDatabaseValidation(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	var playerCount, teamCount, matchCount, tournamentCount int64

	database.DB.WithContext(ctx).Model(&database.Player{}).Count(&playerCount)
	database.DB.WithContext(ctx).Model(&database.Team{}).Count(&teamCount)
	database.DB.WithContext(ctx).Model(&database.Match{}).Count(&matchCount)
	database.DB.WithContext(ctx).Model(&database.Tournament{}).Count(&tournamentCount)

	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"timestamp":   time.Now().Unix(),
		"players":     playerCount,
		"teams":       teamCount,
		"matches":     matchCount,
		"tournaments": tournamentCount,
	})
}
