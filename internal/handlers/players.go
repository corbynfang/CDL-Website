package handlers

// players.go — handlers for the /players endpoints.
// Includes the new GetPlayerFranchiseCareer which aggregates a player's match stats
// across all teams, grouped by franchise slot so career continuity is preserved even
// when a franchise renames itself (e.g. Minnesota RØKKR → G2 Minnesota).

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

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

// GetPlayerStats returns per-match aggregate stats for a player.
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

// GetPlayerKDStats returns a player's overall and per-mode K/D, plus tournament breakdown.
func GetPlayerKDStats(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	var player database.Player
	if err := database.DB.WithContext(ctx).First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	var tournamentStats []database.PlayerTournamentStats
	if err := database.DB.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("Tournament").
		Order("tournament_id DESC").
		Find(&tournamentStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	var totalKills, totalDeaths, totalAssists int
	var totalHpK, totalHpD, totalSndK, totalSndD, totalCtlK, totalCtlD int

	tournamentList := make([]gin.H, 0, len(tournamentStats))
	for _, stat := range tournamentStats {
		totalKills += stat.TotalKills
		totalDeaths += stat.TotalDeaths
		totalAssists += stat.TotalAssists
		totalHpK += stat.HpKills
		totalHpD += stat.HpDeaths
		totalSndK += stat.SndKills
		totalSndD += stat.SndDeaths
		totalCtlK += stat.ControlKills
		totalCtlD += stat.ControlDeaths
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
		"hp_kd_ratio":      calculateKD(totalHpK, totalHpD),
		"snd_kd_ratio":     calculateKD(totalSndK, totalSndD),
		"control_kd_ratio": calculateKD(totalCtlK, totalCtlD),
		"tournament_stats": tournamentList,
	})
}

// GetPlayerMatches returns match history for a player, grouped by tournament.
// Each match includes match_id so the frontend can link to the match detail page.
func GetPlayerMatches(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

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

	eventsMap := make(map[uint]gin.H)
	for _, stat := range matchStats {
		match := stat.Match
		tournamentID := match.TournamentID
		if eventsMap[tournamentID] == nil {
			eventsMap[tournamentID] = gin.H{
				"event":         match.Tournament.Name,
				"year":          match.Tournament.StartDate.Year(),
				"tournament_id": tournamentID,
				"matches":       []gin.H{},
			}
		}

		var opponent, opponentAbbr, result string
		if stat.TeamID == match.Team1ID {
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
			"match_id":      match.ID,
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

// GetPlayerFranchiseCareer aggregates a player's match stats across all teams they've
// played for, grouped by franchise slot. This preserves career continuity even when
// franchises rename: OpTic Chicago + OpTic Texas appear as one franchise career.
// Non-CDL appearances (challenger teams, EWC with no CDL franchise) are grouped separately.
func GetPlayerFranchiseCareer(c *gin.Context) {
	playerID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	var player database.Player
	if err := database.DB.WithContext(ctx).First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	// One row per (franchise, team, game_code) combination.
	// Franchise can be null for non-CDL teams.
	type eraRow struct {
		FranchiseID   *uint
		FranchiseKey  string
		FranchiseName string
		TeamID        uint
		TeamName      string
		GameCode      string
		SeasonName    string
		Matches       int
		Maps          int
		Kills         int
		Deaths        int
	}

	var rows []eraRow
	database.DB.WithContext(ctx).Raw(`
		SELECT
			t.franchise_id,
			COALESCE(f.franchise_key, '')      AS franchise_key,
			COALESCE(f.name, t.name)           AS franchise_name,
			pms.team_id                        AS team_id,
			t.name                             AS team_name,
			COALESCE(t.game_code, '')          AS game_code,
			COALESCE(s.name, '')               AS season_name,
			COUNT(DISTINCT pms.match_id)       AS matches,
			SUM(pms.maps_played)               AS maps,
			SUM(pms.total_kills)               AS kills,
			SUM(pms.total_deaths)              AS deaths
		FROM player_match_stats pms
		JOIN teams t ON t.id = pms.team_id
		LEFT JOIN franchises f ON f.id = t.franchise_id
		LEFT JOIN matches m ON m.id = pms.match_id
		LEFT JOIN tournaments tour ON tour.id = m.tournament_id
		LEFT JOIN seasons s ON s.id = tour.season_id
		WHERE pms.player_id = ?
		GROUP BY t.franchise_id, f.franchise_key, f.name, pms.team_id, t.name, t.game_code, s.name
		ORDER BY t.franchise_id NULLS LAST, t.game_code
	`, playerID).Scan(&rows)

	type eraOut struct {
		TeamID     uint    `json:"team_id"`
		TeamName   string  `json:"team_name"`
		GameCode   string  `json:"game_code"`
		SeasonName string  `json:"season_name"`
		Matches    int     `json:"matches"`
		Maps       int     `json:"maps"`
		Kills      int     `json:"kills"`
		Deaths     int     `json:"deaths"`
		KD         float64 `json:"kd"`
	}
	type franchiseOut struct {
		FranchiseKey  string    `json:"franchise_key"`
		FranchiseName string    `json:"franchise_name"`
		Eras          []eraOut  `json:"eras"`
		TotalMatches  int       `json:"total_matches"`
		TotalMaps     int       `json:"total_maps"`
		TotalKills    int       `json:"total_kills"`
		TotalDeaths   int       `json:"total_deaths"`
		CareerKD      float64   `json:"career_kd"`
	}

	franchiseMap := map[string]*franchiseOut{}
	var franchiseOrder []string

	for _, r := range rows {
		key := r.FranchiseKey
		if key == "" {
			key = "misc"
		}
		if _, ok := franchiseMap[key]; !ok {
			name := r.FranchiseName
			if key == "misc" {
				name = "Non-CDL / Other"
			}
			franchiseMap[key] = &franchiseOut{
				FranchiseKey:  key,
				FranchiseName: name,
				Eras:          []eraOut{},
			}
			franchiseOrder = append(franchiseOrder, key)
		}

		kd := calculateKD(r.Kills, r.Deaths)
		franchiseMap[key].Eras = append(franchiseMap[key].Eras, eraOut{
			TeamID:     r.TeamID,
			TeamName:   r.TeamName,
			GameCode:   r.GameCode,
			SeasonName: r.SeasonName,
			Matches:    r.Matches,
			Maps:       r.Maps,
			Kills:      r.Kills,
			Deaths:     r.Deaths,
			KD:         kd,
		})
		franchiseMap[key].TotalMatches += r.Matches
		franchiseMap[key].TotalMaps += r.Maps
		franchiseMap[key].TotalKills += r.Kills
		franchiseMap[key].TotalDeaths += r.Deaths
	}

	result := make([]franchiseOut, 0, len(franchiseOrder))
	for _, key := range franchiseOrder {
		f := franchiseMap[key]
		f.CareerKD = calculateKD(f.TotalKills, f.TotalDeaths)
		result = append(result, *f)
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":  playerID,
		"gamertag":   player.Gamertag,
		"franchises": result,
	})
}
