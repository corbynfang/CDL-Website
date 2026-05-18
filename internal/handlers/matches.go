package handlers

// matches.go — handlers for /matches and /tournaments endpoints.
// GetMatch is the richest endpoint: it returns a full match with per-map scoreboards.
// Tournament handlers power the bracket visualization.

import (
	"log"
	"net/http"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

// GetMatch returns a single match with full map-by-map player scoreboards.
// Response: { match: {...}, maps: [ { map_number, map_name, mode, score_1, score_2,
//   team1_stats: [...], team2_stats: [...] } ] }
func GetMatch(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	var match database.Match
	if err := database.DB.WithContext(ctx).
		Preload("Team1").
		Preload("Team2").
		Preload("Winner").
		Preload("Tournament").
		Preload("Tournament.Season").
		First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	var matchMaps []database.MatchMap
	database.DB.WithContext(ctx).
		Where("match_id = ?", id).
		Order("map_number ASC").
		Find(&matchMaps)

	type playerStatRow struct {
		MapNumber       int
		PlayerID        uint
		Gamertag        string
		TeamID          uint
		Kills           int
		Deaths          int
		KDRatio         float64
		Damage          int
		Assists         int
		BPRating        float64
		HillTime        int
		SndRounds       int
		PlantCount      int
		DefuseCount     int
		FirstBloodCount int
		FirstDeathCount int
		NonTradedKills  int
		HighestStreak   int
		DataQualityNote string
	}
	var statRows []playerStatRow
	database.DB.WithContext(ctx).
		Table("player_map_stats pms").
		Select(`pms.map_number, pms.player_id, p.gamertag, pms.team_id,
			pms.kills, pms.deaths, pms.kd_ratio, pms.damage, pms.assists,
			pms.bp_rating, pms.hill_time, pms.snd_rounds, pms.plant_count,
			pms.defuse_count, pms.first_blood_count, pms.first_death_count,
			pms.non_traded_kills, pms.highest_streak, pms.data_quality_note`).
		Joins("JOIN players p ON p.id = pms.player_id").
		Where("pms.match_id = ?", id).
		Order("pms.map_number ASC, pms.kills DESC").
		Scan(&statRows)

	statsByMap := map[int][]playerStatRow{}
	for _, s := range statRows {
		statsByMap[s.MapNumber] = append(statsByMap[s.MapNumber], s)
	}

	type playerOut struct {
		PlayerID        uint    `json:"player_id"`
		Gamertag        string  `json:"gamertag"`
		Kills           int     `json:"kills"`
		Deaths          int     `json:"deaths"`
		KDRatio         float64 `json:"kd_ratio"`
		Damage          int     `json:"damage"`
		Assists         int     `json:"assists"`
		BPRating        float64 `json:"bp_rating"`
		HillTime        int     `json:"hill_time"`
		SndRounds       int     `json:"snd_rounds"`
		PlantCount      int     `json:"plant_count"`
		DefuseCount     int     `json:"defuse_count"`
		FirstBloodCount int     `json:"first_blood_count"`
		FirstDeathCount int     `json:"first_death_count"`
		NonTradedKills  int     `json:"non_traded_kills"`
		HighestStreak   int     `json:"highest_streak"`
		DataQualityNote string  `json:"data_quality_note,omitempty"`
	}
	type mapOut struct {
		MapNumber   int         `json:"map_number"`
		MapName     string      `json:"map_name"`
		Mode        string      `json:"mode"`
		Score1      int         `json:"score_1"`
		Score2      int         `json:"score_2"`
		WinnerID    *uint       `json:"winner_id"`
		DurationSec int         `json:"duration_sec"`
		Played      bool        `json:"played"`
		Team1Stats  []playerOut `json:"team1_stats"`
		Team2Stats  []playerOut `json:"team2_stats"`
	}

	maps := make([]mapOut, 0, len(matchMaps))
	for _, mm := range matchMaps {
		out := mapOut{
			MapNumber:   mm.MapNumber,
			MapName:     mm.MapName,
			Mode:        mm.Mode,
			Score1:      mm.Score1,
			Score2:      mm.Score2,
			WinnerID:    mm.WinnerID,
			DurationSec: mm.DurationSec,
			Played:      mm.Played,
			Team1Stats:  []playerOut{},
			Team2Stats:  []playerOut{},
		}
		for _, s := range statsByMap[mm.MapNumber] {
			p := playerOut{
				PlayerID:        s.PlayerID,
				Gamertag:        s.Gamertag,
				Kills:           s.Kills,
				Deaths:          s.Deaths,
				KDRatio:         s.KDRatio,
				Damage:          s.Damage,
				Assists:         s.Assists,
				BPRating:        s.BPRating,
				HillTime:        s.HillTime,
				SndRounds:       s.SndRounds,
				PlantCount:      s.PlantCount,
				DefuseCount:     s.DefuseCount,
				FirstBloodCount: s.FirstBloodCount,
				FirstDeathCount: s.FirstDeathCount,
				NonTradedKills:  s.NonTradedKills,
				HighestStreak:   s.HighestStreak,
				DataQualityNote: s.DataQualityNote,
			}
			if s.TeamID == match.Team1ID {
				out.Team1Stats = append(out.Team1Stats, p)
			} else {
				out.Team2Stats = append(out.Team2Stats, p)
			}
		}
		maps = append(maps, out)
	}

	c.JSON(http.StatusOK, gin.H{
		"match": gin.H{
			"id":                      match.ID,
			"tournament_id":           match.TournamentID,
			"tournament_name":         match.Tournament.Name,
			"tournament_slug":         match.Tournament.Slug,
			"season_name":             match.Tournament.Season.Name,
			"game_code":               match.Tournament.Season.GameCode,
			"team1_id":                match.Team1ID,
			"team1_name":              match.Team1.Name,
			"team1_abbr":              match.Team1.Abbreviation,
			"team1_logo":              match.Team1.LogoURL,
			"team2_id":                match.Team2ID,
			"team2_name":              match.Team2.Name,
			"team2_abbr":              match.Team2.Abbreviation,
			"team2_logo":              match.Team2.LogoURL,
			"team1_score":             match.Team1Score,
			"team2_score":             match.Team2Score,
			"winner_id":               match.WinnerID,
			"match_date":              match.MatchDate,
			"format":                  match.Format,
			"bracket_round":           match.BracketRound,
			"breaking_point_match_id": match.BreakingPointMatchID,
		},
		"maps": maps,
	})
}

func GetTournaments(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	seasonID := c.Query("season_id")
	// Hide data-artifact tournaments from the events page
	query := database.DB.WithContext(ctx).
		Preload("Season").
		Where("tournament_type NOT IN ('season_summary','unknown')").
		Order("start_date DESC")
	if seasonID != "" {
		query = query.Where("season_id = ?", seasonID)
	}

	var tournaments []database.Tournament
	if err := query.Find(&tournaments).Error; err != nil {
		log.Printf("GetTournaments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

// GetTournamentBySlug resolves an event by its URL-safe slug.
// Used by the frontend to load /events/:slug without knowing the numeric ID.
func GetTournamentBySlug(c *gin.Context) {
	slug := c.Param("slug")

	ctx, cancel := getContext(10)
	defer cancel()

	var tournament database.Tournament
	if err := database.DB.WithContext(ctx).
		Preload("Season").
		Where("slug = ?", slug).
		First(&tournament).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Derive team_count from distinct team IDs in matches for this tournament
	var teamCount int64
	database.DB.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM (
			SELECT team1_id AS team_id FROM matches WHERE tournament_id = ?
			UNION
			SELECT team2_id FROM matches WHERE tournament_id = ?
		) AS t
	`, tournament.ID, tournament.ID).Scan(&teamCount)

	c.JSON(http.StatusOK, gin.H{
		"tournament": tournament,
		"team_count": teamCount,
	})
}

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

// GetTournamentBracket returns all bracket matches for a tournament organized by round.
func GetTournamentBracket(c *gin.Context) {
	tournamentID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	var tournament database.Tournament
	if err := database.DB.WithContext(ctx).First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	var matches []database.Match
	database.DB.WithContext(ctx).
		Where("tournament_id = ? AND bracket_round != ''", tournamentID).
		Preload("Team1").
		Preload("Team2").
		Order("bracket_round, bracket_position").
		Find(&matches)

	format := detectBracketFormat(tournament.TournamentFormat, tournament.TournamentType)
	normalize := roundNormalizerFor(format)

	// Build bracket map: pre-initialize each known key with an empty slice so
	// rounds with 0 matches still appear in the response.
	bracketKeySet := bracketKeysFor(format)
	bracket := make(map[string][]gin.H, len(bracketKeySet))
	for k := range bracketKeySet {
		bracket[k] = []gin.H{}
	}

	// Group stage: nil for pure bracket formats; dynamic accumulator for formats
	// that separate group/play-in matches (CDL group-stage, EWC).
	var groupStage map[string][]gin.H
	if hasGroupStage(format) {
		groupStage = map[string][]gin.H{}
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
		key := normalize(match.BracketRound)
		if _, inBracket := bracket[key]; inBracket {
			bracket[key] = append(bracket[key], matchData)
		} else if groupStage != nil {
			groupStage[key] = append(groupStage[key], matchData)
		}
	}

	resp := gin.H{
		"tournament_id":   tournamentID,
		"tournament_name": tournament.Name,
		"event_format":    formatName(format),
		"total_matches":   len(matches),
		"bracket":         bracket,
	}
	if groupStage != nil {
		resp["group_stage"] = groupStage
	}
	c.JSON(http.StatusOK, resp)
}

// GetTournamentMatches returns every match for a tournament with team info,
// scores, and bracket context. Used by the Matches tab on the event detail page.
func GetTournamentMatches(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	var matches []database.Match
	if err := database.DB.WithContext(ctx).
		Where("tournament_id = ?", id).
		Preload("Team1").
		Preload("Team2").
		Preload("Winner").
		Order("match_date ASC, bracket_position ASC").
		Find(&matches).Error; err != nil {
		log.Printf("GetTournamentMatches error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch matches"})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GetTournamentTeams returns every team that played in a tournament,
// enriched with placement and record from team_tournament_stats.
func GetTournamentTeams(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	// Collect every distinct team_id that appears in this tournament's matches
	var teamIDs []uint
	database.DB.WithContext(ctx).Raw(`
		SELECT DISTINCT team_id FROM (
			SELECT team1_id AS team_id FROM matches WHERE tournament_id = ?
			UNION ALL
			SELECT team2_id FROM matches WHERE tournament_id = ?
		) AS t
	`, id, id).Scan(&teamIDs)

	if len(teamIDs) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var teams []database.Team
	database.DB.WithContext(ctx).Where("id IN ?", teamIDs).Find(&teams)

	var stats []database.TeamTournamentStats
	database.DB.WithContext(ctx).Where("tournament_id = ?", id).Find(&stats)

	statsMap := make(map[uint]database.TeamTournamentStats, len(stats))
	for _, s := range stats {
		statsMap[s.TeamID] = s
	}

	type teamOut struct {
		database.Team
		Placement   *int `json:"placement"`
		MatchesWon  int  `json:"matches_won"`
		MatchesLost int  `json:"matches_lost"`
	}

	result := make([]teamOut, 0, len(teams))
	for _, t := range teams {
		out := teamOut{Team: t}
		if s, ok := statsMap[t.ID]; ok {
			out.Placement = s.Placement
			out.MatchesWon = s.MatchesWon
			out.MatchesLost = s.MatchesLost
		}
		result = append(result, out)
	}

	c.JSON(http.StatusOK, result)
}

// GetTournamentStats returns per-player K/D stats for a single tournament.
// Used by the Stats tab on the event detail page.
func GetTournamentStats(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	ctx, cancel := getContext(15)
	defer cancel()

	var stats []database.PlayerTournamentStats
	if err := database.DB.WithContext(ctx).
		Where("tournament_id = ? AND (total_kills > 0 OR total_deaths > 0)", id).
		Preload("Player").
		Preload("Team").
		Order("(CASE WHEN total_deaths > 0 THEN CAST(total_kills AS decimal) / total_deaths ELSE 0 END) DESC").
		Find(&stats).Error; err != nil {
		log.Printf("GetTournamentStats error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
