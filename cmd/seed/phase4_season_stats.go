package main

// phase4_season_stats.go — seeds season-level aggregate player stats.
// The *_player_stats.csv files contain pre-aggregated K/D, HP/SND/CTL breakdowns
// per player per season. These power the Stats leaderboard page.
// This is separate from per-map stats because the aggregation was done externally
// and includes stats not derivable from our map-level data alone (e.g. hp_k/10m).

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedSeasonStats(db *gorm.DB, cfg seasonStatCfg, seasonID uint, teamLookup map[string]uint, playerLookup map[string]uint) {
	f, err := os.Open(cfg.PlayerFile)
	if err != nil {
		log.Printf("[%s] skipping season stats: %v", cfg.GameCode, err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		return
	}

	headers := normalizeHeaders(records[0])
	get := func(rec []string, name string) string {
		i, ok := headers[name]
		if !ok || i >= len(rec) {
			return ""
		}
		return strings.TrimSpace(rec[i])
	}

	// A virtual "Season Stats" tournament holds all aggregate rows.
	// It's not a real event — just a DB container so PlayerTournamentStats has a tournament_id.
	summaryTournament := database.Tournament{
		SeasonID:       seasonID,
		Name:           cfg.Name + " — Season Stats",
		Slug:           cfg.GameCode + "-season-stats",
		TournamentType: "season_summary",
		StartDate:      time.Date(cfg.StartYear, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	db.Where("slug = ? AND season_id = ?", summaryTournament.Slug, seasonID).FirstOrCreate(&summaryTournament)

	var statsBatch []database.PlayerTournamentStats

	for _, rec := range records[1:] {
		gamertag := get(rec, "player")
		if gamertag == "" {
			continue
		}
		isBad := false
		for _, bad := range badGamertags {
			if strings.EqualFold(gamertag, bad) {
				isBad = true
				break
			}
		}
		if isBad {
			continue
		}

		playerID := resolvePlayer(gamertag, playerLookup, db)
		if playerID == 0 {
			continue
		}
		teamID := dominantTeam(db, playerID, seasonID)
		if teamID == 0 {
			teamID = ensureUnaffiliatedTeam(db, teamLookup)
		}

		rank := atoi(get(rec, "rank"))
		statsBatch = append(statsBatch, database.PlayerTournamentStats{
			PlayerID:        playerID,
			TeamID:          teamID,
			TournamentID:    summaryTournament.ID,
			Rank:            &rank,
			TotalKills:      atoi(get(rec, "kills")),
			TotalDeaths:     atoi(get(rec, "deaths")),
			KDRatio:         atof(get(rec, "k/d")),
			OverallMaps:     atoi(get(rec, "series_played")),
			HpKills:         atoi(get(rec, "hp_kills")),
			HpDeaths:        atoi(get(rec, "hp_deaths")),
			HpKDRatio:       atof(get(rec, "hp_k/d")),
			HpKPerMap:       atof(get(rec, "hp_k/10m")),
			HpMaps:          atoi(get(rec, "hp_maps_played")),
			SndKills:        atoi(get(rec, "snd_kills")),
			SndDeaths:       atoi(get(rec, "snd_deaths")),
			SndKDRatio:      atof(get(rec, "snd_k/d")),
			SndKPerMap:      atof(get(rec, "snd_kpr")),
			SndMaps:         atoi(get(rec, "snd_maps_played")),
			ControlKDRatio:  atof(get(rec, "ctl_kd")),
			ControlKPerMap:  atof(get(rec, "ctl_k/10m")),
			ControlCaptures: atoi(get(rec, "ctl_ticks")),
			ControlMaps:     atoi(get(rec, "ctl_maps_played")),
		})
	}

	if len(statsBatch) > 0 {
		db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(statsBatch, 500)
	}
	log.Printf("[%s] season stats seeded: %d rows", cfg.GameCode, len(statsBatch))
}
