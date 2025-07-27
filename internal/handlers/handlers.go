package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

func GetTeams(c *gin.Context) {
	var teams []database.Team

	if err := database.DB.Find(&teams).Error; err != nil {
		log.Printf("Database Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func GetTeam(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var team database.Team
	if err := database.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

func GetTeamPlayers(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var players []database.Player
	if err := database.DB.Joins("JOIN team_rosters ON players.id = team_rosters.player_id").
		Where("team_rosters.team_id = ? AND team_rosters.end_date IS NULL", teamID).
		Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

func GetPlayers(c *gin.Context) {
	var players []database.Player
	if err := database.DB.Find(&players).Error; err != nil {
		log.Printf("Database Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	c.JSON(http.StatusOK, players)
}

func GetPlayer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var player database.Player
	if err := database.DB.First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func GetPlayerStats(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var stats []database.PlayerMatchStats
	if err := database.DB.Where("player_id = ?", id).
		Preload("Match").
		Preload("Team").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetTeamStats(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var stats []database.TeamTournamentStats
	if err := database.DB.Where("team_id = ?", id).
		Preload("Tournament").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Tournament handlers
func GetTournaments(c *gin.Context) {
	var tournaments []database.Tournament
	if err := database.DB.Preload("Season").Find(&tournaments).Error; err != nil {
		log.Printf("Database Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

func GetTournament(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	var tournament database.Tournament
	if err := database.DB.Preload("Season").First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// GetPlayerKDStats returns K/D statistics for a specific player
func GetPlayerKDStats(c *gin.Context) {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var stats []database.PlayerTournamentStats
	if err := database.DB.Where("player_id = ?", playerID).
		Preload("Tournament").
		Preload("Team").
		Order("tournament_id DESC").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player K/D stats"})
		return
	}

	// Group stats by tournament
	tournamentStats := make(map[int]gin.H)
	var totalKills, totalDeaths, totalAssists int
	var totalMaps int

	for _, stat := range stats {
		tournamentID := int(stat.TournamentID)

		if tournamentStats[tournamentID] == nil {
			// Get tournament name
			var tournament database.Tournament
			if err := database.DB.First(&tournament, tournamentID).Error; err != nil {
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

	// Create mock match stats from tournament stats for compatibility
	var matchStats []gin.H
	for _, stat := range stats {
		matchStats = append(matchStats, gin.H{
			"id":            stat.ID,
			"match_id":      stat.TournamentID,
			"player_id":     stat.PlayerID,
			"team_id":       stat.TeamID,
			"maps_played":   1,
			"total_kills":   stat.TotalKills,
			"total_deaths":  stat.TotalDeaths,
			"total_assists": stat.TotalAssists,
			"kd_ratio":      stat.KDRatio,
			"kda_ratio":     stat.KDARatio,
		})
	}

	response := gin.H{
		"player_id":        playerID,
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

	c.JSON(http.StatusOK, response)
}

// GetTopKDPlayers returns the top K/D players for 2025
func GetTopKDPlayers(c *gin.Context) {
	var results []gin.H

	log.Printf("Starting GetTopKDPlayers query")

	// Use a simpler query first
	rows, err := database.DB.Raw(`
		SELECT 
			pms.player_id, 
			pms.team_id, 
			pms.kd_ratio,
			pms.kda_ratio,
			p.gamertag,
			t.name as team_name,
			t.abbreviation as team_abbreviation
		FROM player_match_stats pms
		JOIN players p ON pms.player_id = p.id
		JOIN teams t ON pms.team_id = t.id
		WHERE pms.kd_ratio > 0
		ORDER BY pms.kd_ratio DESC 
		LIMIT 20
	`).Rows()

	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}
	defer rows.Close()

	log.Printf("Query executed successfully, scanning rows")

	for rows.Next() {
		var playerID, teamID uint
		var kdRatio, kdaRatio float64
		var gamertag, teamName, teamAbbreviation string

		err := rows.Scan(&playerID, &teamID, &kdRatio, &kdaRatio, &gamertag, &teamName, &teamAbbreviation)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("Found player: %s, KD: %f, KDA: %f", gamertag, kdRatio, kdaRatio)

		results = append(results, gin.H{
			"player_id":         playerID,
			"gamertag":          gamertag,
			"team_name":         teamName,
			"team_abbreviation": teamAbbreviation,
			"avg_kd":            kdRatio,
			"avg_kda":           kdaRatio,
			"matches_played":    1,
		})
	}

	log.Printf("Returning %d results", len(results))
	c.JSON(http.StatusOK, results)
}

// GetTopKDPlayersNew is a new version of the top KD players handler
func GetTopKDPlayersNew(c *gin.Context) {
	var results []gin.H

	log.Printf("Starting GetTopKDPlayersNew query")

	// Use a very simple query
	rows, err := database.DB.Raw(`
		SELECT 
			pms.player_id, 
			pms.kd_ratio,
			pms.kda_ratio,
			p.gamertag,
			t.name as team_name,
			t.abbreviation as team_abbreviation
		FROM player_match_stats pms
		JOIN players p ON pms.player_id = p.id
		JOIN teams t ON pms.team_id = t.id
		WHERE pms.kd_ratio > 0
		ORDER BY pms.kd_ratio DESC 
		LIMIT 10
	`).Rows()

	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top K/D players"})
		return
	}
	defer rows.Close()

	log.Printf("Query executed successfully, scanning rows")

	for rows.Next() {
		var playerID uint
		var kdRatio, kdaRatio float64
		var gamertag, teamName, teamAbbreviation string

		err := rows.Scan(&playerID, &kdRatio, &kdaRatio, &gamertag, &teamName, &teamAbbreviation)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("Found player: %s, KD: %f, KDA: %f", gamertag, kdRatio, kdaRatio)

		results = append(results, gin.H{
			"player_id":         playerID,
			"gamertag":          gamertag,
			"team_name":         teamName,
			"team_abbreviation": teamAbbreviation,
			"avg_kd":            kdRatio,
			"avg_kda":           kdaRatio,
			"matches_played":    1,
		})
	}

	log.Printf("Returning %d results", len(results))
	c.JSON(http.StatusOK, results)
}

// GetAllPlayersKDStats returns KD and KD+/- for all players for the season, and KD for each major tournament
func GetAllPlayersKDStats(c *gin.Context) {
	// Get all players
	var players []database.Player
	if err := database.DB.Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	// Get all player_tournament_stats for Majors 1-4 and Champs (tournament_id 1-5)
	type KDRow struct {
		PlayerID     uint
		TeamID       uint
		TournamentID uint
		TotalKills   int
		TotalDeaths  int
		KDRatio      float64
		Gamertag     string
		TeamAbbr     string
	}

	var kdRows []KDRow
	db := database.DB
	db.Raw(`
		SELECT pts.player_id, pts.team_id, pts.tournament_id, pts.total_kills, pts.total_deaths, pts.kd_ratio, p.gamertag, t.abbreviation as team_abbr
		FROM player_tournament_stats pts
		JOIN players p ON pts.player_id = p.id
		JOIN teams t ON pts.team_id = t.id
		WHERE pts.tournament_id IN (1,2,3,4,5)
	`).Scan(&kdRows)

	// Build a map: playerID -> {majors: {tournamentID: KD}, ...}
	playerMap := make(map[uint]gin.H)
	for _, p := range players {
		playerMap[p.ID] = gin.H{
			"player_id":     p.ID,
			"gamertag":      p.Gamertag,
			"avatar_url":    p.AvatarURL,
			"team_abbr":     "",
			"majors":        map[uint]float64{},
			"season_kills":  0,
			"season_deaths": 0,
		}
	}
	for _, row := range kdRows {
		if playerMap[row.PlayerID] != nil {
			playerMap[row.PlayerID]["team_abbr"] = row.TeamAbbr
			playerMap[row.PlayerID]["majors"].(map[uint]float64)[row.TournamentID] = row.KDRatio
			playerMap[row.PlayerID]["season_kills"] = playerMap[row.PlayerID]["season_kills"].(int) + row.TotalKills
			playerMap[row.PlayerID]["season_deaths"] = playerMap[row.PlayerID]["season_deaths"].(int) + row.TotalDeaths
		}
	}

	// Ensure every player has a KD for every major (1-5)
	for _, p := range playerMap {
		majors := p["majors"].(map[uint]float64)
		for i := 1; i <= 5; i++ {
			if _, ok := majors[uint(i)]; !ok {
				majors[uint(i)] = 0.0 // or use null if you prefer
			}
		}
	}

	// Build response
	var result []gin.H
	for _, p := range playerMap {
		seasonKills := p["season_kills"].(int)
		seasonDeaths := p["season_deaths"].(int)
		var seasonKD float64
		if seasonDeaths > 0 {
			seasonKD = float64(seasonKills) / float64(seasonDeaths)
		}
		seasonKDPlusMinus := seasonKD - 1.0
		result = append(result, gin.H{
			"player_id":            p["player_id"],
			"gamertag":             p["gamertag"],
			"avatar_url":           p["avatar_url"],
			"team_abbr":            p["team_abbr"],
			"season_kd":            seasonKD,
			"season_kd_plus_minus": seasonKDPlusMinus,
			"majors":               p["majors"],
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetTransfers(c *gin.Context) {
	var transfers []database.PlayerTransfer

	query := database.DB.Preload("Player").Preload("FromTeam").Preload("ToTeam")

	// Add filters if provided
	if season := c.Query("season"); season != "" {
		query = query.Where("season = ?", season)
	}

	if teamID := c.Query("team_id"); teamID != "" {
		query = query.Where("from_team_id = ? OR to_team_id = ?", teamID, teamID)
	}

	if transferType := c.Query("type"); transferType != "" {
		query = query.Where("transfer_type = ?", transferType)
	}

	// Order by transfer date (most recent first)
	if err := query.Order("transfer_date DESC").Find(&transfers).Error; err != nil {
		log.Printf("Database Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transfers"})
		return
	}

	c.JSON(http.StatusOK, transfers)
}
