package handlers

// teams.go — handlers for the /teams endpoints.
// The DB has one Team row per era branding (Minnesota RØKKR and G2 Minnesota are two rows,
// both linked to the same Franchise via franchise_id). GetTeams uses DISTINCT ON franchise_id
// so exactly one team per CDL slot is returned regardless of how many brandings exist.

import (
	"log"
	"net/http"
	"sort"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
)

// GetTeams returns one CDL franchise team per franchise slot.
//
// No season_id: returns the most recent branding per franchise (current identity).
// Only franchises active in the current CDL season are included.
//
// ?season_id=N: returns the era-correct branding per franchise for that season.
// Within the same game_code, valid_from ASC picks the CDL-season branding over
// post-season rebrands (e.g. Minnesota RØKKR before G2 Minnesota for BO6).
func GetTeams(c *gin.Context) {
	ctx, cancel := getContext(10)
	defer cancel()

	seasonID := c.Query("season_id")
	var teams []database.Team
	var err error

	if seasonID != "" {
		var season database.Season
		if err = database.DB.WithContext(ctx).First(&season, seasonID).Error; err != nil || season.GameCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season"})
			return
		}

		// scope=cdl (default): CDL events only.
		// scope=era: all events for the game era including EWC.
		//
		// is_cdl_franchise always filters out EWC-only org teams regardless of
		// scope — FaZe Clan, Cloud9 parent org etc. belong on event pages only.
		scope := c.Query("scope")
		var sql string
		var args []interface{}
		if scope == "era" {
			sql = `
				SELECT DISTINCT t.*
				FROM teams t
				WHERE t.is_cdl_franchise = true
				  AND t.id IN (
				    SELECT DISTINCT pms.team_id
				    FROM player_match_stats pms
				    JOIN matches m ON m.id = pms.match_id
				    JOIN tournaments trn ON trn.id = m.tournament_id
				    WHERE trn.season_id = ?
				      AND trn.tournament_type NOT IN ('season_summary','unknown')
				  )
				ORDER BY t.name ASC`
			args = []interface{}{seasonID}
		} else {
			sql = `
				SELECT DISTINCT t.*
				FROM teams t
				WHERE t.is_cdl_franchise = true
				  AND t.id IN (
				    SELECT DISTINCT pms.team_id
				    FROM player_match_stats pms
				    JOIN matches m ON m.id = pms.match_id
				    JOIN tournaments trn ON trn.id = m.tournament_id
				    WHERE trn.season_id = ?
				      AND trn.tournament_type IN (
				        'major_tournament','qualifier','championship',
				        'kickoff','minor_tournament'
				      )
				  )
				ORDER BY t.name ASC`
			args = []interface{}{seasonID}
		}
		err = database.DB.WithContext(ctx).Raw(sql, args...).Scan(&teams).Error
	} else {
		// No season filter: one team per franchise, most recent branding.
		// scope=all returns every team row (CDL + non-CDL) for admin/debug use.
		scope := c.Query("scope")
		if scope == "all" {
			err = database.DB.WithContext(ctx).
				Order("name ASC").
				Find(&teams).Error
		} else {
			err = database.DB.WithContext(ctx).Raw(`
				SELECT DISTINCT ON (franchise_id) *
				FROM teams
				WHERE is_cdl_franchise = true
				  AND franchise_id IS NOT NULL
				ORDER BY franchise_id, valid_from DESC
			`).Scan(&teams).Error
		}
	}

	if err != nil {
		log.Printf("GetTeams error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	sort.Slice(teams, func(i, j int) bool { return teams[i].Name < teams[j].Name })
	c.JSON(http.StatusOK, teams)
}

// GetTeam returns one team by ID, including its Franchise so the frontend
// can load the franchise_key and fetch the full franchise history.
func GetTeam(c *gin.Context) {
	id, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	var team database.Team
	if err := database.DB.WithContext(ctx).Preload("Franchise").First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

// GetTeamPlayers returns all players on a team, optionally filtered by season.
func GetTeamPlayers(c *gin.Context) {
	teamID, err := validateID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()

	seasonID := c.Query("season_id")
	query := database.DB.WithContext(ctx).
		Joins("JOIN team_rosters ON players.id = team_rosters.player_id").
		Where("team_rosters.team_id = ?", teamID)

	if seasonID != "" {
		query = query.Where("team_rosters.season_id = ?", seasonID)
	} else {
		query = query.Where("team_rosters.end_date IS NULL")
	}

	var players []database.Player
	if err := query.Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team players"})
		return
	}
	c.JSON(http.StatusOK, players)
}

// GetTeamStats returns tournament placement/record stats for a team.
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
