package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

// Security logging function
func logSecurityEvent(event, details, clientIP string) {
	log.Printf("[SECURITY] %s | %s | %s | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		event,
		details,
		clientIP)
}

// Enhanced input validation
func validateID(id string) (int, error) {
	// Check if ID is numeric and within reasonable bounds
	if id == "" {
		return 0, fmt.Errorf("empty ID")
	}

	// Use regex to ensure only digits
	matched, _ := regexp.MatchString(`^\d+$`, id)
	if !matched {
		return 0, fmt.Errorf("non-numeric ID")
	}

	// Convert to int
	num, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}

	// Check reasonable bounds (1 to 1 million)
	if num < 1 || num > 1000000 {
		return 0, fmt.Errorf("ID out of bounds")
	}

	return num, nil
}

// Sanitize query parameters
func sanitizeQueryParam(param string) string {
	// Remove potentially dangerous characters
	param = strings.TrimSpace(param)
	param = strings.ReplaceAll(param, ";", "")
	param = strings.ReplaceAll(param, "--", "")
	param = strings.ReplaceAll(param, "/*", "")
	param = strings.ReplaceAll(param, "*/", "")
	param = strings.ReplaceAll(param, "xp_", "")
	param = strings.ReplaceAll(param, "sp_", "")
	return param
}

func GetTeams(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetTeams", c.ClientIP())

	var teams []database.Team

	// Use context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch teams with their current players
	if err := database.DB.WithContext(ctx).
		Preload("Players", "team_rosters.end_date IS NULL").
		Where("is_active = ?", true).
		Find(&teams).Error; err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetTeams failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	// For each team, get their current players manually
	for i := range teams {
		var players []database.Player
		if err := database.DB.WithContext(ctx).
			Joins("JOIN team_rosters ON players.id = team_rosters.player_id").
			Where("team_rosters.team_id = ? AND team_rosters.end_date IS NULL", teams[i].ID).
			Find(&players).Error; err == nil {
			teams[i].Players = players
		}
	}

	c.JSON(http.StatusOK, teams)
}

func GetTeam(c *gin.Context) {
	// Enhanced input validation
	id, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid team ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var team database.Team

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

func GetTeamPlayers(c *gin.Context) {
	// Enhanced input validation
	teamID, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid team ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var players []database.Player

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Joins("JOIN team_rosters ON players.id = team_rosters.player_id").
		Where("team_rosters.team_id = ? AND team_rosters.end_date IS NULL", teamID).
		Find(&players).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetTeamPlayers failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

func GetPlayers(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetPlayers", c.ClientIP())

	var players []database.Player

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Find(&players).Error; err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetPlayers failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	c.JSON(http.StatusOK, players)
}

func GetPlayer(c *gin.Context) {
	// Enhanced input validation
	id, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid player ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var player database.Player

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func GetPlayerStats(c *gin.Context) {
	// Enhanced input validation
	id, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid player ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var stats []database.PlayerMatchStats

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Where("player_id = ?", id).
		Preload("Match").
		Preload("Team").
		Find(&stats).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetPlayerStats failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetTeamStats(c *gin.Context) {
	// Enhanced input validation
	id, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid team ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var stats []database.TeamTournamentStats

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Where("team_id = ?", id).
		Preload("Tournament").
		Find(&stats).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetTeamStats failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Tournament handlers
func GetTournaments(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetTournaments", c.ClientIP())

	var tournaments []database.Tournament

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Preload("Season").Find(&tournaments).Error; err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetTournaments failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

func GetTournament(c *gin.Context) {
	// Enhanced input validation
	id, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid tournament ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	var tournament database.Tournament

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// GetPlayerKDStats returns K/D statistics for a specific player
func GetPlayerKDStats(c *gin.Context) {
	// Enhanced input validation
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid player ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var stats []database.PlayerTournamentStats

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Where("player_id = ?", playerID).
		Preload("Tournament").
		Preload("Team").
		Order("tournament_id DESC").
		Find(&stats).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetPlayerKDStats failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player K/D stats"})
		return
	}

	// Log raw stats for debugging
	log.Printf("Raw stats for player %d: %+v", playerID, stats)

	// Group stats by tournament
	tournamentStats := make(map[int]gin.H)
	var totalKills, totalDeaths, totalAssists int
	var totalMaps int

	for _, stat := range stats {
		tournamentID := int(stat.TournamentID)

		if tournamentStats[tournamentID] == nil {
			// Get tournament name
			var tournament database.Tournament
			if err := database.DB.WithContext(ctx).First(&tournament, tournamentID).Error; err != nil {
				log.Printf("Error fetching tournament %d: %v", tournamentID, err)
				tournamentStats[tournamentID] = gin.H{
					"tournament_id":   tournamentID,
					"tournament_name": "Unknown Tournament",
					"kills":           0,
					"deaths":          0,
					"assists":         0,
					"maps_played":     0,
					"matches":         0,
				}
			} else {
				tournamentStats[tournamentID] = gin.H{
					"tournament_id":   tournamentID,
					"tournament_name": tournament.Name,
					"kills":           0,
					"deaths":          0,
					"assists":         0,
					"maps_played":     0,
					"matches":         0,
				}
			}
		}

		tournament := tournamentStats[tournamentID]
		tournament["kills"] = tournament["kills"].(int) + stat.TotalKills
		tournament["deaths"] = tournament["deaths"].(int) + stat.TotalDeaths
		tournament["assists"] = tournament["assists"].(int) + stat.TotalAssists
		tournament["maps_played"] = tournament["maps_played"].(int) + 1 // Each tournament stat represents one tournament
		tournament["matches"] = tournament["matches"].(int) + 1

		totalKills += stat.TotalKills
		totalDeaths += stat.TotalDeaths
		totalAssists += stat.TotalAssists
		totalMaps += 1
	}

	// Calculate KD ratios for each tournament
	var tournamentStatsList []gin.H
	for tournamentID, stats := range tournamentStats {
		tournament := stats
		kills := tournament["kills"].(int)
		deaths := tournament["deaths"].(int)
		assists := tournament["assists"].(int)

		var kdRatio, kdaRatio float64
		if deaths > 0 {
			kdRatio = float64(kills) / float64(deaths)
			kdaRatio = float64(kills+assists) / float64(deaths)
		}

		tournament["kd_ratio"] = kdRatio
		tournament["kda_ratio"] = kdaRatio
		tournament["tournament_id"] = tournamentID
		tournamentStatsList = append(tournamentStatsList, tournament)
	}

	// Calculate overall statistics
	var avgKD, avgKDA, avgADR float64
	if totalDeaths > 0 {
		avgKD = float64(totalKills) / float64(totalDeaths)
		avgKDA = float64(totalKills+totalAssists) / float64(totalDeaths)
	}

	if totalMaps > 0 {
		avgADR = float64(totalKills*100) / float64(totalMaps) // Simplified ADR calculation
	}

	// Log calculated stats for debugging
	log.Printf("Calculated stats for player %d: Kills=%d, Deaths=%d, KD=%.2f, KDA=%.2f",
		playerID, totalKills, totalDeaths, avgKD, avgKDA)

	// Create mock match stats from tournament stats for compatibility
	var matchStats []gin.H
	for _, stat := range stats {
		// Recalculate KD for each stat to ensure accuracy
		var calculatedKD, calculatedKDA float64
		if stat.TotalDeaths > 0 {
			calculatedKD = float64(stat.TotalKills) / float64(stat.TotalDeaths)
			calculatedKDA = float64(stat.TotalKills+stat.TotalAssists) / float64(stat.TotalDeaths)
		}

		matchStats = append(matchStats, gin.H{
			"id":            stat.ID,
			"match_id":      stat.TournamentID,
			"player_id":     stat.PlayerID,
			"team_id":       stat.TeamID,
			"maps_played":   1,
			"total_kills":   stat.TotalKills,
			"total_deaths":  stat.TotalDeaths,
			"total_assists": stat.TotalAssists,
			"kd_ratio":      calculatedKD,
			"kda_ratio":     calculatedKDA,
			"db_kd_ratio":   stat.KDRatio, // Keep original for comparison
			"db_kda_ratio":  stat.KDARatio,
		})
	}

	// Get player info for response
	var player database.Player
	if err := database.DB.WithContext(ctx).First(&player, playerID).Error; err != nil {
		log.Printf("Error fetching player %d: %v", playerID, err)
	}

	// Get EWC2025 detailed stats if available
	var ewcStats database.PlayerTournamentStats
	var ewcDetailed gin.H
	if err := database.DB.WithContext(ctx).Where("player_id = ? AND tournament_id = 7", playerID).First(&ewcStats).Error; err == nil {
		ewcDetailed = gin.H{
			"ewc_snd_kills":            ewcStats.SndKills,
			"ewc_snd_deaths":           ewcStats.SndDeaths,
			"ewc_snd_kd_ratio":         ewcStats.SndKDRatio,
			"ewc_snd_plus_minus":       ewcStats.SndPlusMinus,
			"ewc_snd_k_per_map":        ewcStats.SndKPerMap,
			"ewc_snd_first_kills":      ewcStats.SndFirstKills,
			"ewc_snd_maps":             ewcStats.SndMaps,
			"ewc_hp_kills":             ewcStats.HpKills,
			"ewc_hp_deaths":            ewcStats.HpDeaths,
			"ewc_hp_kd_ratio":          ewcStats.HpKDRatio,
			"ewc_hp_plus_minus":        ewcStats.HpPlusMinus,
			"ewc_hp_k_per_map":         ewcStats.HpKPerMap,
			"ewc_hp_time_milliseconds": ewcStats.HpTimeMilliseconds,
			"ewc_hp_maps":              ewcStats.HpMaps,
			"ewc_control_kills":        ewcStats.ControlKills,
			"ewc_control_deaths":       ewcStats.ControlDeaths,
			"ewc_control_kd_ratio":     ewcStats.ControlKDRatio,
			"ewc_control_plus_minus":   ewcStats.ControlPlusMinus,
			"ewc_control_k_per_map":    ewcStats.ControlKPerMap,
			"ewc_control_captures":     ewcStats.ControlCaptures,
			"ewc_control_maps":         ewcStats.ControlMaps,
		}
	}

	response := gin.H{
		"player_id":        playerID,
		"gamertag":         player.Gamertag,
		"avatar_url":       player.AvatarURL,
		"total_matches":    len(stats),
		"total_maps":       totalMaps,
		"total_kills":      totalKills,
		"total_deaths":     totalDeaths,
		"total_assists":    totalAssists,
		"avg_kd":           avgKD,
		"avg_kda":          avgKDA,
		"avg_adr":          avgADR,
		"tournament_stats": tournamentStatsList,
		"match_stats":      matchStats,
	}

	// Add EWC2025 detailed stats if available
	for key, value := range ewcDetailed {
		response[key] = value
	}

	c.JSON(http.StatusOK, response)
}

// GetTopKDPlayers returns the top K/D players for 2025
func GetTopKDPlayers(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetTopKDPlayers", c.ClientIP())

	var results []gin.H

	log.Printf("Starting GetTopKDPlayers query")

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Use tournament stats for consistent KD calculation
	rows, err := database.DB.WithContext(ctx).Raw(`
		SELECT 
			pts.player_id, 
			pts.team_id, 
			pts.total_kills,
			pts.total_deaths,
			pts.total_assists,
			pts.kd_ratio,
			pts.kda_ratio,
			p.gamertag,
			t.name as team_name,
			t.abbreviation as team_abbreviation
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5,7) 
		AND pts.total_deaths > 0
		ORDER BY pts.kd_ratio DESC 
		LIMIT 20
	`).Rows()

	if err != nil {
		log.Printf("Error executing query: %v", err)
		logSecurityEvent("DB_ERROR", "GetTopKDPlayers failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}
	defer rows.Close()

	log.Printf("Query executed successfully, scanning rows")

	for rows.Next() {
		var playerID, teamID uint
		var totalKills, totalDeaths, totalAssists int
		var kdRatio, kdaRatio float64
		var gamertag, teamName, teamAbbreviation string

		err := rows.Scan(&playerID, &teamID, &totalKills, &totalDeaths, &totalAssists, &kdRatio, &kdaRatio, &gamertag, &teamName, &teamAbbreviation)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Recalculate KD to ensure accuracy
		var calculatedKD float64
		if totalDeaths > 0 {
			calculatedKD = float64(totalKills) / float64(totalDeaths)
		}

		log.Printf("Found player: %s, Kills: %d, Deaths: %d, Calculated KD: %f, DB KD: %f",
			gamertag, totalKills, totalDeaths, calculatedKD, kdRatio)

		results = append(results, gin.H{
			"player_id":         playerID,
			"gamertag":          gamertag,
			"team_name":         teamName,
			"team_abbreviation": teamAbbreviation,
			"total_kills":       totalKills,
			"total_deaths":      totalDeaths,
			"total_assists":     totalAssists,
			"kd_ratio":          calculatedKD,
			"kda_ratio":         kdaRatio,
		})
	}

	log.Printf("Returning %d results", len(results))
	c.JSON(http.StatusOK, results)
}

// GetTopKDPlayersNew is a new version of the top KD players handler with aggregated stats
func GetTopKDPlayersNew(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetTopKDPlayersNew", c.ClientIP())

	var results []gin.H

	log.Printf("Starting GetTopKDPlayersNew query")

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Aggregate stats across all tournaments for each player
	rows, err := database.DB.WithContext(ctx).Raw(`
		SELECT 
			pts.player_id, 
			SUM(pts.total_kills) as total_kills,
			SUM(pts.total_deaths) as total_deaths,
			SUM(pts.total_assists) as total_assists,
			p.gamertag,
			t.name as team_name,
			t.abbreviation as team_abbreviation,
			COUNT(DISTINCT pts.tournament_id) as tournaments_played
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5,7)
		GROUP BY pts.player_id, p.gamertag, t.name, t.abbreviation
		HAVING SUM(pts.total_deaths) > 0
		ORDER BY (SUM(pts.total_kills) * 1.0 / SUM(pts.total_deaths)) DESC 
		LIMIT 10
	`).Rows()

	if err != nil {
		log.Printf("Error executing query: %v", err)
		logSecurityEvent("DB_ERROR", "GetTopKDPlayersNew failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}
	defer rows.Close()

	log.Printf("Query executed successfully, scanning rows")

	for rows.Next() {
		var playerID uint
		var totalKills, totalDeaths, totalAssists, tournamentsPlayed int
		var gamertag, teamName, teamAbbreviation string

		err := rows.Scan(&playerID, &totalKills, &totalDeaths, &totalAssists, &gamertag, &teamName, &teamAbbreviation, &tournamentsPlayed)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Calculate KD and KDA
		var kdRatio, kdaRatio float64
		if totalDeaths > 0 {
			kdRatio = float64(totalKills) / float64(totalDeaths)
			kdaRatio = float64(totalKills+totalAssists) / float64(totalDeaths)
		}

		log.Printf("Found player: %s, Kills: %d, Deaths: %d, KD: %f, Tournaments: %d",
			gamertag, totalKills, totalDeaths, kdRatio, tournamentsPlayed)

		results = append(results, gin.H{
			"player_id":          playerID,
			"gamertag":           gamertag,
			"team_name":          teamName,
			"team_abbreviation":  teamAbbreviation,
			"total_kills":        totalKills,
			"total_deaths":       totalDeaths,
			"total_assists":      totalAssists,
			"kd_ratio":           kdRatio,
			"kda_ratio":          kdaRatio,
			"tournaments_played": tournamentsPlayed,
		})
	}

	log.Printf("Returning %d results", len(results))
	c.JSON(http.StatusOK, results)
}

// GetAllPlayersKDStats returns KD and KD+/- for all players for the season, and KD for each major tournament
func GetAllPlayersKDStats(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetAllPlayersKDStats", c.ClientIP())

	// Add cache-busting headers for Railway
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("X-Railway-Cache", "disabled")

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all players with tournament stats in one optimized query
	type KDRow struct {
		PlayerID     uint
		TeamID       uint
		TournamentID uint
		TotalKills   int
		TotalDeaths  int
		TotalAssists int
		KDRatio      float64
		KDARatio     float64
		Gamertag     string
		TeamAbbr     string
		AvatarURL    string
	}

	var kdRows []KDRow
	if err := database.DB.WithContext(ctx).Raw(`
		SELECT 
			pts.player_id, 
			pts.team_id, 
			pts.tournament_id, 
			pts.total_kills, 
			pts.total_deaths, 
			pts.total_assists,
			pts.kd_ratio,
			pts.kda_ratio,
			p.gamertag, 
			COALESCE(t.abbreviation, 'N/A') as team_abbr, 
			p.avatar_url
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		LEFT JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5,7)
		ORDER BY pts.player_id, pts.tournament_id
	`).Scan(&kdRows).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetAllPlayersKDStats failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	// Build a map: playerID -> {majors: {tournamentID: KD}, ...}
	playerMap := make(map[uint]gin.H)

	// Coaches to exclude from the KD stats
	excludedCoaches := map[string]bool{
		"Accuracy": true,
		"Crimsix":  true,
	}

	// Initialize player map from kdRows
	for _, row := range kdRows {
		if excludedCoaches[row.Gamertag] {
			continue
		}

		if playerMap[row.PlayerID] == nil {
			playerMap[row.PlayerID] = gin.H{
				"player_id":      row.PlayerID,
				"gamertag":       row.Gamertag,
				"avatar_url":     row.AvatarURL,
				"team_abbr":      row.TeamAbbr,
				"majors":         map[uint]gin.H{},
				"season_kills":   0,
				"season_deaths":  0,
				"season_assists": 0,
			}
		}
	}

	// Process tournament stats
	for _, row := range kdRows {
		if playerMap[row.PlayerID] != nil {
			// Calculate KD for this tournament
			var tournamentKD float64
			if row.TotalDeaths > 0 {
				tournamentKD = float64(row.TotalKills) / float64(row.TotalDeaths)
			}

			// Store tournament data
			playerMap[row.PlayerID]["majors"].(map[uint]gin.H)[row.TournamentID] = gin.H{
				"kd_ratio": tournamentKD,
				"kills":    row.TotalKills,
				"deaths":   row.TotalDeaths,
				"assists":  row.TotalAssists,
			}

			// Accumulate season totals
			playerMap[row.PlayerID]["season_kills"] = playerMap[row.PlayerID]["season_kills"].(int) + row.TotalKills
			playerMap[row.PlayerID]["season_deaths"] = playerMap[row.PlayerID]["season_deaths"].(int) + row.TotalDeaths
			playerMap[row.PlayerID]["season_assists"] = playerMap[row.PlayerID]["season_assists"].(int) + row.TotalAssists
		}
	}

	// Only include players who have stats for at least one tournament
	// and ensure every player has a KD for every major (1-5)
	for _, p := range playerMap {
		majors := p["majors"].(map[uint]gin.H)
		// Only include players who have at least one tournament stat
		if len(majors) == 0 {
			continue
		}
		for i := 1; i <= 5; i++ {
			if _, ok := majors[uint(i)]; !ok {
				majors[uint(i)] = gin.H{
					"kd_ratio": 0.0,
					"kills":    0,
					"deaths":   0,
					"assists":  0,
				}
			}
		}
	}

	// Build response - use a map to ensure no duplicates by player_id
	uniquePlayers := make(map[uint]gin.H)
	for _, p := range playerMap {
		seasonKills := p["season_kills"].(int)
		seasonDeaths := p["season_deaths"].(int)
		seasonAssists := p["season_assists"].(int)

		var seasonKD, seasonKDA float64
		if seasonDeaths > 0 {
			seasonKD = float64(seasonKills) / float64(seasonDeaths)
			seasonKDA = float64(seasonKills+seasonAssists) / float64(seasonDeaths)
		}
		seasonKDPlusMinus := seasonKD - 1.0

		playerID := p["player_id"].(uint)
		uniquePlayers[playerID] = gin.H{
			"player_id":            playerID,
			"gamertag":             p["gamertag"],
			"avatar_url":           p["avatar_url"],
			"team_abbr":            p["team_abbr"],
			"season_kills":         seasonKills,
			"season_deaths":        seasonDeaths,
			"season_assists":       seasonAssists,
			"season_kd":            seasonKD,
			"season_kda":           seasonKDA,
			"season_kd_plus_minus": seasonKDPlusMinus,
			"majors":               p["majors"],
		}
	}

	// Convert map to slice
	var result []gin.H
	for _, player := range uniquePlayers {
		result = append(result, player)
	}

	// Add timestamp to response for cache busting
	response := gin.H{
		"timestamp": time.Now().Unix(),
		"players":   result,
		"count":     len(result),
	}

	c.JSON(http.StatusOK, response)
}

func GetTransfers(c *gin.Context) {
	// Log request for security monitoring
	logSecurityEvent("API_ACCESS", "GetTransfers", c.ClientIP())

	// Add cache-busting headers for Railway
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("X-Railway-Cache", "disabled")

	var transfers []database.PlayerTransfer

	// Sanitize query parameters
	season := sanitizeQueryParam(c.Query("season"))
	teamID := sanitizeQueryParam(c.Query("team_id"))
	transferType := sanitizeQueryParam(c.Query("type"))
	playerID := sanitizeQueryParam(c.Query("player_id"))

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := database.DB.WithContext(ctx).Preload("Player").Preload("FromTeam").Preload("ToTeam")

	// Add filters if provided
	if season != "" {
		query = query.Where("season = ?", season)
	}

	if teamID != "" {
		// Validate teamID is numeric
		if _, err := validateID(teamID); err != nil {
			logSecurityEvent("INVALID_INPUT", "Invalid team_id in query: "+teamID, c.ClientIP())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team_id parameter"})
			return
		}
		query = query.Where("from_team_id = ? OR to_team_id = ?", teamID, teamID)
	}

	if transferType != "" {
		query = query.Where("transfer_type = ?", transferType)
	}

	if playerID != "" {
		// Validate playerID is numeric
		if _, err := validateID(playerID); err != nil {
			logSecurityEvent("INVALID_INPUT", "Invalid player_id in query: "+playerID, c.ClientIP())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player_id parameter"})
			return
		}
		query = query.Where("player_id = ?", playerID)
	}

	// Order by transfer date (most recent first)
	if err := query.Order("transfer_date DESC").Find(&transfers).Error; err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetTransfers failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transfers"})
		return
	}

	// Add timestamp to response for cache busting
	response := gin.H{
		"timestamp": time.Now().Unix(),
		"transfers": transfers,
		"count":     len(transfers),
	}

	c.JSON(http.StatusOK, response)
}

// Database validation utility
func validateDatabaseStats() gin.H {
	var issues []string

	// Check for players with zero deaths but non-zero kills (shouldn't happen)
	var zeroDeathPlayers []gin.H
	database.DB.Raw(`
		SELECT player_id, total_kills, total_deaths, total_assists
		FROM player_tournament_stats 
		WHERE total_deaths = 0 AND total_kills > 0
		LIMIT 10
	`).Scan(&zeroDeathPlayers)

	if len(zeroDeathPlayers) > 0 {
		issues = append(issues, fmt.Sprintf("Found %d players with zero deaths but non-zero kills", len(zeroDeathPlayers)))
	}

	// Check for players with negative stats
	var negativeStats []gin.H
	database.DB.Raw(`
		SELECT player_id, total_kills, total_deaths, total_assists
		FROM player_tournament_stats 
		WHERE total_kills < 0 OR total_deaths < 0 OR total_assists < 0
		LIMIT 10
	`).Scan(&negativeStats)

	if len(negativeStats) > 0 {
		issues = append(issues, fmt.Sprintf("Found %d records with negative stats", len(negativeStats)))
	}

	// Check for KD ratio inconsistencies
	var kdInconsistencies []gin.H
	database.DB.Raw(`
		SELECT player_id, total_kills, total_deaths, kd_ratio,
		       CASE 
		           WHEN total_deaths > 0 THEN ROUND((total_kills * 1.0 / total_deaths)::numeric, 2)
		           ELSE 0 
		       END as calculated_kd
		FROM player_tournament_stats 
		WHERE total_deaths > 0 
		AND ABS(kd_ratio - (total_kills * 1.0 / total_deaths)) > 0.01
		LIMIT 10
	`).Scan(&kdInconsistencies)

	if len(kdInconsistencies) > 0 {
		issues = append(issues, fmt.Sprintf("Found %d records with KD ratio inconsistencies", len(kdInconsistencies)))
	}

	return gin.H{
		"issues":             issues,
		"zero_death_players": zeroDeathPlayers,
		"negative_stats":     negativeStats,
		"kd_inconsistencies": kdInconsistencies,
		"timestamp":          time.Now().Unix(),
	}
}

// GetDatabaseValidation returns database validation results
func GetDatabaseValidation(c *gin.Context) {
	logSecurityEvent("API_ACCESS", "GetDatabaseValidation", c.ClientIP())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use context for database operations
	database.DB.WithContext(ctx)

	validation := validateDatabaseStats()
	c.JSON(http.StatusOK, validation)
}

// GetPlayerMatches returns all matches for a specific player
func GetPlayerMatches(c *gin.Context) {
	// Enhanced input validation
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		logSecurityEvent("INVALID_INPUT", "Invalid player ID: "+c.Param("id"), c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	logSecurityEvent("API_ACCESS", "GetPlayerMatches for player "+c.Param("id"), c.ClientIP())

	// Get query parameters
	tournamentID := sanitizeQueryParam(c.Query("tournament_id"))
	limit := sanitizeQueryParam(c.Query("limit"))

	// Set default limit if not provided
	if limit == "" {
		limit = "50"
	}

	// Validate limit
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		limitInt = 50
	}

	// Build the query
	query := database.DB.Table("player_match_stats").
		Select(`
			player_match_stats.*,
			matches.match_date,
			matches.team1_score,
			matches.team2_score,
			matches.match_type,
			matches.format,
			tournaments.name as tournament_name,
			t1.name as team1_name,
			t1.abbreviation as team1_abbr,
			t2.name as team2_name,
			t2.abbreviation as team2_abbr,
			players.gamertag,
			teams.name as player_team_name,
			teams.abbreviation as player_team_abbr
		`).
		Joins("JOIN matches ON player_match_stats.match_id = matches.id").
		Joins("JOIN tournaments ON matches.tournament_id = tournaments.id").
		Joins("JOIN teams t1 ON matches.team1_id = t1.id").
		Joins("JOIN teams t2 ON matches.team2_id = t2.id").
		Joins("JOIN players ON player_match_stats.player_id = players.id").
		Joins("JOIN teams ON player_match_stats.team_id = teams.id").
		Where("player_match_stats.player_id = ?", playerID)

	// Add tournament filter if provided
	if tournamentID != "" {
		if _, err := validateID(tournamentID); err != nil {
			logSecurityEvent("INVALID_INPUT", "Invalid tournament_id in query: "+tournamentID, c.ClientIP())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id parameter"})
			return
		}
		query = query.Where("matches.tournament_id = ?", tournamentID)
	}

	// Execute query with limit and ordering
	var matches []gin.H
	if err := query.Order("matches.match_date DESC").Limit(limitInt).Find(&matches).Error; err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetPlayerMatches failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player matches"})
		return
	}

	// Transform the data to include calculated fields
	var result []gin.H
	for _, match := range matches {
		// Determine if player's team won
		playerTeamID := match["team_id"].(uint)
		team1ID := match["team1_id"].(uint)
		team1Score := match["team1_score"].(int)
		team2Score := match["team2_score"].(int)

		var matchResult string
		var opponent string
		var opponentAbbr string

		if playerTeamID == team1ID {
			if team1Score > team2Score {
				matchResult = "W"
			} else {
				matchResult = "L"
			}
			opponent = match["team2_name"].(string)
			opponentAbbr = match["team2_abbr"].(string)
		} else {
			if team2Score > team1Score {
				matchResult = "W"
			} else {
				matchResult = "L"
			}
			opponent = match["team1_name"].(string)
			opponentAbbr = match["team1_abbr"].(string)
		}

		// Format score
		score := fmt.Sprintf("%d:%d", team1Score, team2Score)

		// Calculate game mode KDs (placeholder for now)
		hpKD := 0.0
		sndKD := 0.0
		ctlKD := 0.0

		// Add to result
		result = append(result, gin.H{
			"match_id":      match["match_id"],
			"date":          match["match_date"],
			"tournament":    match["tournament_name"],
			"opponent":      opponent,
			"opponent_abbr": opponentAbbr,
			"result":        matchResult,
			"score":         score,
			"kd":            match["kd_ratio"],
			"kills":         match["total_kills"],
			"deaths":        match["total_deaths"],
			"hp_kd":         hpKD,
			"snd_kd":        sndKD,
			"ctl_kd":        ctlKD,
			"slayer_rating": 0.0, // Placeholder
			"maps_played":   match["maps_played"],
			"match_type":    match["match_type"],
			"format":        match["format"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"matches":   result,
		"total":     len(result),
	})
}
