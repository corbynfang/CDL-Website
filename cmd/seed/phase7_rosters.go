package main

import (
	"database/sql"
	"log"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type rosterStint struct {
	PlayerID  uint
	TeamID    uint
	SeasonID  uint
	StartDate sql.NullTime
	EndDate   sql.NullTime
	MapCount  int
}

func inferRosterStints(db *gorm.DB) ([]rosterStint, error) {
	var stints []rosterStint
	err := db.Raw(`
		SELECT
			pms.player_id  AS player_id,
			pms.team_id    AS team_id,
			tour.season_id AS season_id,
			MIN(m.match_date) FILTER (WHERE m.match_date > '0001-01-02'::timestamptz) AS start_date,
			MAX(m.match_date) FILTER (WHERE m.match_date > '0001-01-02'::timestamptz) AS end_date,
			COUNT(*)       AS map_count
		FROM player_map_stats pms
		JOIN matches m       ON m.id = pms.match_id
		JOIN tournaments tour ON tour.id = m.tournament_id
		LEFT JOIN match_maps mm
			ON mm.match_id = pms.match_id
			AND mm.map_number = pms.map_number
		WHERE pms.team_id <> 0
		  AND tour.season_id <> 0
		  AND (mm.id IS NULL OR mm.played = true)
		GROUP BY pms.player_id, pms.team_id, tour.season_id
		ORDER BY tour.season_id, pms.team_id, pms.player_id
	`).Scan(&stints).Error
	return stints, err
}

func seedRosters(db *gorm.DB) {
	stints, err := inferRosterStints(db)
	if err != nil {
		log.Printf("roster inference failed: %v", err)
		return
	}

	if err := db.Exec("DELETE FROM team_rosters").Error; err != nil {
		log.Printf("failed to clear team_rosters: %v", err)
		return
	}

	rosters := make([]models.TeamRoster, 0, len(stints))
	for _, st := range stints {
		r := models.TeamRoster{
			TeamID:    st.TeamID,
			PlayerID:  st.PlayerID,
			SeasonID:  st.SeasonID,
			IsStarter: true,
		}
		if st.StartDate.Valid {
			r.StartDate = st.StartDate.Time
		}
		if st.EndDate.Valid {
			end := st.EndDate.Time
			r.EndDate = &end
		}
		rosters = append(rosters, r)
	}

	if len(rosters) > 0 {
		if err := db.CreateInBatches(rosters, 500).Error; err != nil {
			log.Printf("failed to insert team_rosters: %v", err)
			return
		}
	}
	log.Printf("rosters inferred: %d stints", len(rosters))
}
