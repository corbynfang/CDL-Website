package handlers

import (
	"context"
	"database/sql"
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

func parsePaginationParams(c *gin.Context, defaultLimit, maxLimit int) (int, int) {
	limit := defaultLimit
	offset := 0

	if rawLimit := sanitizeQueryParam(c.Query("limit")); rawLimit != "" {
		if parsedLimit, err := strconv.Atoi(rawLimit); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	if rawOffset := sanitizeQueryParam(c.Query("offset")); rawOffset != "" {
		if parsedOffset, err := strconv.Atoi(rawOffset); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}

func parseBoolQuery(c *gin.Context, key string) bool {
	value := strings.ToLower(sanitizeQueryParam(c.Query(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

type majorStat struct {
	KDRatio float64 `json:"kd_ratio"`
	Kills   int     `json:"kills"`
	Deaths  int     `json:"deaths"`
	Assists int     `json:"assists"`
}

type playerKDResponse struct {
	PlayerID          uint               `json:"player_id"`
	Gamertag          string             `json:"gamertag"`
	AvatarURL         string             `json:"avatar_url,omitempty"`
	TeamAbbr          string             `json:"team_abbr"`
	SeasonKills       int                `json:"season_kills"`
	SeasonDeaths      int                `json:"season_deaths"`
	SeasonAssists     int                `json:"season_assists"`
	SeasonKD          float64            `json:"season_kd"`
	SeasonKDA         float64            `json:"season_kda"`
	SeasonKDPlusMinus float64            `json:"season_kd_plus_minus"`
	Majors            map[uint]majorStat `json:"majors,omitempty"`
}

func GetTeams(c *gin.Context) {
	logSecurityEvent("API_ACCESS", "GetTeams", c.ClientIP())

	var teams []database.Team

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.DB.WithContext(ctx).Where("is_active = ?", true).Find(&teams).Error; err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetTeams failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
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

	// Calculate KD ratios for each tournament ** VERY IMPORTANT **
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

	var totalHpKills, totalHpDeaths, totalSndKills, totalSndDeaths, totalCtlKills, totalCtlDeaths int
	for _, stat := range stats {
		totalHpKills += stat.HpKills
		totalHpDeaths += stat.HpDeaths
		totalSndKills += stat.SndKills
		totalSndDeaths += stat.SndDeaths
		totalCtlKills += stat.ControlKills
		totalCtlDeaths += stat.ControlDeaths
	}

	var hpKDRatio, sndKDRatio, ctlKDRatio float64
	if totalHpDeaths > 0 {
		hpKDRatio = float64(totalHpKills) / float64(totalHpDeaths)
	}
	if totalSndDeaths > 0 {
		sndKDRatio = float64(totalSndKills) / float64(totalSndDeaths)
	}
	if totalCtlDeaths > 0 {
		ctlKDRatio = float64(totalCtlKills) / float64(totalCtlDeaths)
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
		"player_id":            playerID,
		"gamertag":             player.Gamertag,
		"avatar_url":           player.AvatarURL,
		"total_matches":        len(stats),
		"total_maps":           totalMaps,
		"total_kills":          totalKills,
		"total_deaths":         totalDeaths,
		"total_assists":        totalAssists,
		"avg_kd":               avgKD,
		"avg_kda":              avgKDA,
		"avg_adr":              avgADR,
		"hp_kd_ratio":          hpKDRatio,
		"snd_kd_ratio":         sndKDRatio,
		"control_kd_ratio":     ctlKDRatio,
		"ewc_hp_kd_ratio":      hpKDRatio,  // Use aggregated for display
		"ewc_snd_kd_ratio":     sndKDRatio, // Use aggregated for display
		"ewc_control_kd_ratio": ctlKDRatio, // Use aggregated for display
		"tournament_stats":     tournamentStatsList,
		"match_stats":          matchStats,
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

	limit, offset := parsePaginationParams(c, 25, 100)

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	type topKDRow struct {
		PlayerID      uint
		Gamertag      string
		TeamAbbr      string
		SeasonKills   int
		SeasonDeaths  int
		SeasonAssists int
	}

	query := `
		SELECT
			pts.player_id,
			MAX(p.gamertag) AS gamertag,
			COALESCE(MAX(t.abbreviation), '') AS team_abbr,
			SUM(pts.total_kills) AS season_kills,
			SUM(pts.total_deaths) AS season_deaths,
			SUM(pts.total_assists) AS season_assists
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		LEFT JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5,7)
		GROUP BY pts.player_id
		HAVING SUM(pts.total_deaths) > 0
		ORDER BY (SUM(pts.total_kills)::decimal / NULLIF(SUM(pts.total_deaths), 0)) DESC,
			MAX(p.gamertag)
		LIMIT ? OFFSET ?
	`

	var rows []topKDRow
	if err := database.DB.WithContext(ctx).Raw(query, limit, offset).Scan(&rows).Error; err != nil {
		log.Printf("Error executing query: %v", err)
		logSecurityEvent("DB_ERROR", "GetTopKDPlayers failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}

	players := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		seasonKD := 0.0
		seasonKDA := 0.0
		if row.SeasonDeaths > 0 {
			seasonKD = float64(row.SeasonKills) / float64(row.SeasonDeaths)
			seasonKDA = float64(row.SeasonKills+row.SeasonAssists) / float64(row.SeasonDeaths)
		}

		players = append(players, gin.H{
			"player_id":            row.PlayerID,
			"gamertag":             row.Gamertag,
			"team_abbr":            row.TeamAbbr,
			"season_kills":         row.SeasonKills,
			"season_deaths":        row.SeasonDeaths,
			"season_assists":       row.SeasonAssists,
			"season_kd":            seasonKD,
			"season_kda":           seasonKDA,
			"season_kd_plus_minus": seasonKD - 1.0,
		})
	}

	response := gin.H{
		"timestamp": time.Now().Unix(),
		"players":   players,
		"count":     len(players),
		"limit":     limit,
		"offset":    offset,
	}

	c.JSON(http.StatusOK, response)
}

// GetTopKDPlayersNew is a new version of the top KD players handler with aggregated stats
func GetTopKDPlayersNew(c *gin.Context) {
	logSecurityEvent("API_ACCESS", "GetTopKDPlayersNew", c.ClientIP())

	limit, offset := parsePaginationParams(c, 50, 200)

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	type aggregatedKDRow struct {
		PlayerID          uint
		Gamertag          string
		TeamName          string
		TeamAbbr          string
		SeasonKills       int
		SeasonDeaths      int
		SeasonAssists     int
		TournamentsPlayed int
	}

	query := `
		SELECT
			pts.player_id,
			MAX(p.gamertag) AS gamertag,
			COALESCE(MAX(t.name), '') AS team_name,
			COALESCE(MAX(t.abbreviation), '') AS team_abbr,
			SUM(pts.total_kills) AS season_kills,
			SUM(pts.total_deaths) AS season_deaths,
			SUM(pts.total_assists) AS season_assists,
			COUNT(DISTINCT pts.tournament_id) AS tournaments_played
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		LEFT JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5,7)
		GROUP BY pts.player_id
		HAVING SUM(pts.total_deaths) > 0
		ORDER BY (SUM(pts.total_kills)::decimal / NULLIF(SUM(pts.total_deaths), 0)) DESC,
			MAX(p.gamertag)
		LIMIT ? OFFSET ?
	`

	var rows []aggregatedKDRow
	if err := database.DB.WithContext(ctx).Raw(query, limit, offset).Scan(&rows).Error; err != nil {
		log.Printf("Error executing query: %v", err)
		logSecurityEvent("DB_ERROR", "GetTopKDPlayersNew failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}

	players := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		seasonKD := 0.0
		seasonKDA := 0.0
		if row.SeasonDeaths > 0 {
			seasonKD = float64(row.SeasonKills) / float64(row.SeasonDeaths)
			seasonKDA = float64(row.SeasonKills+row.SeasonAssists) / float64(row.SeasonDeaths)
		}

		players = append(players, gin.H{
			"player_id":            row.PlayerID,
			"gamertag":             row.Gamertag,
			"team_name":            row.TeamName,
			"team_abbr":            row.TeamAbbr,
			"season_kills":         row.SeasonKills,
			"season_deaths":        row.SeasonDeaths,
			"season_assists":       row.SeasonAssists,
			"season_kd":            seasonKD,
			"season_kda":           seasonKDA,
			"season_kd_plus_minus": seasonKD - 1.0,
			"tournaments_played":   row.TournamentsPlayed,
		})
	}

	response := gin.H{
		"timestamp": time.Now().Unix(),
		"players":   players,
		"count":     len(players),
		"limit":     limit,
		"offset":    offset,
	}

	c.JSON(http.StatusOK, response)
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

	limit, offset := parsePaginationParams(c, 100, 500)
	includeMajors := parseBoolQuery(c, "include_majors")

	excludedGametags := []string{"Accuracy", "Crimsix"}
	tournamentFilter := []int{1, 2, 3, 4, 5, 7}

	type playerSeasonRow struct {
		PlayerID      uint
		Gamertag      string
		AvatarURL     string
		TeamAbbr      string
		SeasonKills   int
		SeasonDeaths  int
		SeasonAssists int
	}

	var seasonRows []playerSeasonRow
	seasonQuery := `
		SELECT
			pts.player_id,
			MAX(p.gamertag) AS gamertag,
			COALESCE(MAX(p.avatar_url), '') AS avatar_url,
		COALESCE(MAX(t.abbreviation), '') AS team_abbr,
			SUM(pts.total_kills) AS season_kills,
			SUM(pts.total_deaths) AS season_deaths,
			SUM(pts.total_assists) AS season_assists
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		LEFT JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5,7)
		  AND COALESCE(p.gamertag, '') NOT IN (?, ?)
		GROUP BY pts.player_id
		HAVING SUM(pts.total_kills) > 0 OR SUM(pts.total_deaths) > 0
		ORDER BY
			CASE WHEN SUM(pts.total_deaths) > 0
				THEN SUM(pts.total_kills)::decimal / SUM(pts.total_deaths)
				ELSE 0
			END DESC,
			MAX(p.gamertag)
		LIMIT ? OFFSET ?
	`

	if err := database.DB.WithContext(ctx).Raw(seasonQuery, excludedGametags[0], excludedGametags[1], limit, offset).Scan(&seasonRows).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetAllPlayersKDStats failed (season query)", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	var totalPlayers int64
	if err := database.DB.WithContext(ctx).
		Table("player_tournament_stats AS pts").
		Joins("JOIN players p ON pts.player_id = p.id").
		Where("pts.tournament_id IN ?", tournamentFilter).
		Where("COALESCE(p.gamertag, '') NOT IN ?", excludedGametags).
		Distinct("pts.player_id").
		Count(&totalPlayers).Error; err != nil {
		logSecurityEvent("DB_ERROR", "GetAllPlayersKDStats failed (count)", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats count"})
		return
	}

	playerIDs := make([]uint, 0, len(seasonRows))
	for _, row := range seasonRows {
		playerIDs = append(playerIDs, row.PlayerID)
	}

	majorsByPlayer := make(map[uint]map[uint]majorStat)
	if includeMajors && len(playerIDs) > 0 {
		type majorRow struct {
			PlayerID     uint
			TournamentID uint
			TotalKills   int
			TotalDeaths  int
			TotalAssists int
		}

		var majorRows []majorRow
		if err := database.DB.WithContext(ctx).
			Table("player_tournament_stats").
			Select("player_id, tournament_id, total_kills, total_deaths, total_assists").
			Where("player_id IN ?", playerIDs).
			Where("tournament_id IN ?", tournamentFilter).
			Find(&majorRows).Error; err != nil {
			logSecurityEvent("DB_ERROR", "GetAllPlayersKDStats failed (major breakdown)", c.ClientIP())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player tournament stats"})
			return
		}

		for _, row := range majorRows {
			if majorsByPlayer[row.PlayerID] == nil {
				majorsByPlayer[row.PlayerID] = make(map[uint]majorStat)
			}

			kd := 0.0
			if row.TotalDeaths > 0 {
				kd = float64(row.TotalKills) / float64(row.TotalDeaths)
			}

			majorsByPlayer[row.PlayerID][row.TournamentID] = majorStat{
				KDRatio: kd,
				Kills:   row.TotalKills,
				Deaths:  row.TotalDeaths,
				Assists: row.TotalAssists,
			}
		}
	}

	players := make([]playerKDResponse, 0, len(seasonRows))
	for _, row := range seasonRows {
		seasonKD := 0.0
		seasonKDA := 0.0
		if row.SeasonDeaths > 0 {
			seasonKD = float64(row.SeasonKills) / float64(row.SeasonDeaths)
			seasonKDA = float64(row.SeasonKills+row.SeasonAssists) / float64(row.SeasonDeaths)
		}

		player := playerKDResponse{
			PlayerID:          row.PlayerID,
			Gamertag:          row.Gamertag,
			AvatarURL:         row.AvatarURL,
			TeamAbbr:          row.TeamAbbr,
			SeasonKills:       row.SeasonKills,
			SeasonDeaths:      row.SeasonDeaths,
			SeasonAssists:     row.SeasonAssists,
			SeasonKD:          seasonKD,
			SeasonKDA:         seasonKDA,
			SeasonKDPlusMinus: seasonKD - 1.0,
		}

		if includeMajors {
			player.Majors = majorsByPlayer[row.PlayerID]
		}

		players = append(players, player)
	}

	response := gin.H{
		"timestamp": time.Now().Unix(),
		"players":   players,
		"count":     len(players),
		"total":     totalPlayers,
		"limit":     limit,
		"offset":    offset,
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
	limitInt, _ := parsePaginationParams(c, 50, 100)

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Build raw SQL query
	sqlQuery := `
		SELECT
			pms.match_id,
			pms.team_id,
			pms.total_kills,
			pms.total_deaths,
			pms.total_assists,
			pms.kd_ratio,
			pms.maps_played,
			m.match_date,
			m.team1_score,
			m.team2_score,
			m.team1_id,
			m.team2_id,
			m.tournament_id,
			m.match_type,
			m.format,
			t.name as tournament_name,
			t.start_date as tournament_start_date,
			t1.name as team1_name,
			t1.abbreviation as team1_abbr,
			t2.name as team2_name,
			t2.abbreviation as team2_abbr,
			p.gamertag,
			t3.name as player_team_name,
			t3.abbreviation as player_team_abbr
		FROM player_match_stats pms
		JOIN matches m ON pms.match_id = m.id
		JOIN tournaments t ON m.tournament_id = t.id
		JOIN teams t1 ON m.team1_id = t1.id
		JOIN teams t2 ON m.team2_id = t2.id
		JOIN players p ON pms.player_id = p.id
		JOIN teams t3 ON pms.team_id = t3.id
		WHERE pms.player_id = ?
	`

	args := []interface{}{playerID}

	// Add tournament filter if provided
	if tournamentID != "" {
		if _, err := validateID(tournamentID); err != nil {
			logSecurityEvent("INVALID_INPUT", "Invalid tournament_id in query: "+tournamentID, c.ClientIP())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id parameter"})
			return
		}
		sqlQuery += " AND m.tournament_id = ?"
		args = append(args, tournamentID)
	}

	sqlQuery += " ORDER BY m.match_date DESC LIMIT ?"
	args = append(args, limitInt)

	log.Printf("Executing Raw SQL query for player %d with limit %d", playerID, limitInt)
	log.Printf("SQL: %s", sqlQuery)

	// Execute query
	rows, err := database.DB.Raw(sqlQuery, args...).Rows()
	if err != nil {
		log.Printf("Database Error: %v", err)
		logSecurityEvent("DB_ERROR", "GetPlayerMatches failed", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player matches"})
		return
	}
	defer rows.Close()

	var matches []map[string]interface{}

	// Scan rows into matches slice
	for rows.Next() {
		var matchID, teamID, totalKills, totalDeaths, totalAssists, mapsPlayed, tournamentID int
		var kdRatio float64
		var tournamentName, team1Name, team1Abbr, team2Name, team2Abbr, gamertag, playerTeamName, playerTeamAbbr string
		var matchType, format sql.NullString
		var team1Score, team2Score, team1ID, team2ID int
		var matchDate time.Time
		var tournamentStartDate time.Time

		err := rows.Scan(
			&matchID, &teamID, &totalKills, &totalDeaths, &totalAssists, &kdRatio, &mapsPlayed,
			&matchDate, &team1Score, &team2Score, &team1ID, &team2ID, &tournamentID, &matchType, &format,
			&tournamentName, &tournamentStartDate, &team1Name, &team1Abbr, &team2Name, &team2Abbr, &gamertag,
			&playerTeamName, &playerTeamAbbr,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		matchTypeStr := ""
		if matchType.Valid {
			matchTypeStr = matchType.String
		}
		formatStr := ""
		if format.Valid {
			formatStr = format.String
		}

		matches = append(matches, map[string]interface{}{
			"match_id":         matchID,
			"team_id":          teamID,
			"total_kills":      totalKills,
			"total_deaths":     totalDeaths,
			"total_assists":    totalAssists,
			"kd_ratio":         kdRatio,
			"maps_played":      mapsPlayed,
			"match_date":       matchDate,
			"team1_score":      team1Score,
			"team2_score":      team2Score,
			"team1_id":         team1ID,
			"team2_id":         team2ID,
			"tournament_id":    tournamentID,
			"tournament_name":  tournamentName,
			"tournament_year":  tournamentStartDate.Year(),
			"match_type":       matchTypeStr,
			"format":           formatStr,
			"team1_name":       team1Name,
			"team1_abbr":       team1Abbr,
			"team2_name":       team2Name,
			"team2_abbr":       team2Abbr,
			"gamertag":         gamertag,
			"player_team_name": playerTeamName,
			"player_team_abbr": playerTeamAbbr,
		})
	}

	tournamentStatsMap := make(map[uint]database.PlayerTournamentStats)
	var tournamentStats []database.PlayerTournamentStats
	database.DB.WithContext(ctx).Where("player_id = ?", playerID).Find(&tournamentStats)
	for _, ts := range tournamentStats {
		tournamentStatsMap[ts.TournamentID] = ts
	}

	eventsMap := make(map[int]gin.H)
	tournamentYears := make(map[int]int)

	for _, match := range matches {
		tournamentID := match["tournament_id"].(int)
		tournamentYear := match["tournament_year"].(int)

		if eventsMap[tournamentID] == nil {
			tournamentYears[tournamentID] = tournamentYear
			eventsMap[tournamentID] = gin.H{
				"event":   match["tournament_name"].(string),
				"year":    tournamentYear,
				"matches": []gin.H{},
			}
		}

		// Determine if the player's team won
		playerTeamID := match["team_id"].(int)
		team1ID := match["team1_id"].(int)
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

		// Format score and result
		score := fmt.Sprintf("%d:%d", team1Score, team2Score)
		resultScore := fmt.Sprintf("%s %s", matchResult, score)

		matchKills := match["total_kills"].(int)
		matchDeaths := match["total_deaths"].(int)
		matchKD := match["kd_ratio"].(float64)

		// Calculate slayer rating (simplified: kills per map)
		slayerRating := 0.0
		if mapsPlayed, ok := match["maps_played"].(int); ok && mapsPlayed > 0 {
			slayerRating = float64(matchKills) / float64(mapsPlayed)
		}

		// Calculate rating (simplified: KD * 10)
		rating := matchKD * 10.0

		// Parse match date
		matchDateTime := match["match_date"].(time.Time)
		matchDate := matchDateTime.In(time.UTC)
		matchDateFormatted := ""
		if !matchDate.IsZero() {
			matchDateFormatted = matchDate.Format(time.RFC3339)
		}

		// Add match to event
		event := eventsMap[tournamentID]
		matchesList := event["matches"].([]gin.H)

		var hpKD *float64
		var sndKD *float64
		var ctlKD *float64
		if ts, ok := tournamentStatsMap[uint(tournamentID)]; ok {
			if ts.HpKDRatio > 0 {
				value := ts.HpKDRatio
				hpKD = &value
			}
			if ts.SndKDRatio > 0 {
				value := ts.SndKDRatio
				sndKD = &value
			}
			if ts.ControlKDRatio > 0 {
				value := ts.ControlKDRatio
				ctlKD = &value
			}
		}
		matchesList = append(matchesList, gin.H{
			"date":          matchDateFormatted,
			"opponent":      opponent,
			"opponent_abbr": opponentAbbr,
			"result":        resultScore,
			"kd":            matchKD,
			"kills":         matchKills,
			"deaths":        matchDeaths,
			"hp_kd":         hpKD,
			"snd_kd":        sndKD,
			"ctl_kd":        ctlKD,
			"slayer_rating": slayerRating,
			"rating":        rating,
		})
		event["matches"] = matchesList
	}

	// Convert map to slice and sort by year/tournament ID (most recent first)
	var events []gin.H
	for tournamentID, event := range eventsMap {
		event["tournament_id"] = tournamentID
		events = append(events, event)
	}

	// Sort events by year (descending)
	for i := 0; i < len(events); i++ {
		for j := i + 1; j < len(events); j++ {
			yearI := events[i]["year"].(int)
			yearJ := events[j]["year"].(int)
			if yearJ > yearI || (yearJ == yearI && events[j]["tournament_id"].(int) > events[i]["tournament_id"].(int)) {
				events[i], events[j] = events[j], events[i]
			}
		}
	}

	// Sort matches within each event by date (descending)
	for _, event := range events {
		matchesList := event["matches"].([]gin.H)
		for i := 0; i < len(matchesList); i++ {
			for j := i + 1; j < len(matchesList); j++ {
				if matchesList[j]["date"].(string) > matchesList[i]["date"].(string) {
					matchesList[i], matchesList[j] = matchesList[j], matchesList[i]
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"events":    events,
		"total":     len(matches),
	})
}
