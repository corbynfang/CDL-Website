package main

// phase6_bracket_patches.go — patches bracket_round and bracket_position for
// all CDL eras from CSV files in database/.
//
// Each CSV must use the standard bracket patch format:
//   tournament_slug, source_round_name, canonical_round_key, bracket_position,
//   team1_name, team2_name, team1_score, team2_score, winner_name, match_date
//
// tournament_slug must match the slug stored in the tournaments table, OR be a
// known alias listed in bracketSlugAliases below.
//
// For each row this phase either:
//   UPDATE — finds an existing match by (tournament, team pair, scores) and
//            sets bracket_round + bracket_position to the CSV values.
//   INSERT — if no match is found, inserts it with a stable dedup key so
//            re-runs are safe.
//
// bracket_edges (nextMatchId relationships for SVG connector lines) are out of
// scope here and will be handled in a later phase once all bracket data is complete.

import (
	"fmt"
	"log"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

// bracketSlugAliases maps source CSV slug values to the tournament slug stored
// in the DB. Only needed when the two differ — add new aliases here as more
// historical data is sourced.
var bracketSlugAliases = map[string]string{
	// CW 2021 used a "Stage N Major" naming convention; DB normalises these to "Major N Tournament".
	"cdl-2021-stage-1-major": "cdl-major-1-tournament-2021",
	"cdl-2021-stage-2-major": "cdl-major-2-tournament-2021",
	"cdl-2021-stage-3-major": "cdl-major-3-tournament-2021",
	"cdl-2021-stage-4-major": "cdl-major-4-tournament-2021",
	"cdl-2021-stage-5-major": "cdl-major-5-tournament-2021",
}

// bracketPatchCSVs is the ordered list of bracket patch files this phase reads.
// Add a new entry here whenever bracket data for a new event is sourced.
var bracketPatchCSVs = []string{
	"database/cdl_cw_stage_brackets.csv",
	"database/cdl_champs_brackets.csv",
	"database/cdl_major_brackets.csv",
}

func seedBracketPatches(
	db *gorm.DB,
	teamLookup map[string]uint,
	tournamentBySlug map[string]uint,
) {
	// Remove any stub matches that a previous Phase 6 run inserted because it
	// couldn't find the correct Phase 2 match (due to the tournament-date bug).
	// These stubs have no match_maps or player stats, so deleting them is safe.
	// Phase 2 now corrects tournament_id in-place, so on re-seed Phase 6 will
	// find and UPDATE the real match instead of inserting a new stub.
	result := db.Where("liquipedia_url LIKE ?", "bracket_patch:%").Delete(&database.Match{})
	log.Printf("[bracket_patches] purged %d stale bracket-patch stub matches", result.RowsAffected)

	var totalUpdated, totalInserted, totalSkipped int

	for _, path := range bracketPatchCSVs {
		u, i, s := applyBracketCSV(db, teamLookup, tournamentBySlug, path)
		totalUpdated += u
		totalInserted += i
		totalSkipped += s
	}

	log.Printf("[bracket_patches] total: updated=%d  inserted=%d  skipped=%d",
		totalUpdated, totalInserted, totalSkipped)
}

func applyBracketCSV(
	db *gorm.DB,
	teamLookup map[string]uint,
	tournamentBySlug map[string]uint,
	path string,
) (updated, inserted, skipped int) {
	rows := readBracketCSV(path)

	for _, r := range rows {
		dbSlug := r.TournamentSlug
		if alias, ok := bracketSlugAliases[r.TournamentSlug]; ok {
			dbSlug = alias
		}

		tournamentID := tournamentBySlug[dbSlug]
		if tournamentID == 0 {
			log.Printf("[bracket_patches] WARN: tournament not found for slug %q — skipping", dbSlug)
			skipped++
			continue
		}

		team1ID := teamLookup[r.Team1Name]
		team2ID := teamLookup[r.Team2Name]
		if team1ID == 0 || team2ID == 0 {
			log.Printf("[bracket_patches] WARN: team not found (%q or %q) — skipping", r.Team1Name, r.Team2Name)
			skipped++
			continue
		}

		// Match by team pair + scores in either orientation.
		// Scores are more reliable than dates because source data sometimes has
		// off-by-one-day differences due to timezone handling in the era_finals seeder.
		var existing database.Match
		err := db.Where(`
			tournament_id = ? AND (
				(team1_id = ? AND team2_id = ? AND team1_score = ? AND team2_score = ?) OR
				(team1_id = ? AND team2_id = ? AND team1_score = ? AND team2_score = ?)
			)`,
			tournamentID,
			team1ID, team2ID, r.Team1Score, r.Team2Score,
			team2ID, team1ID, r.Team2Score, r.Team1Score,
		).First(&existing).Error

		if err == nil {
			db.Model(&existing).Updates(map[string]interface{}{
				"bracket_round":    r.CanonicalRound,
				"bracket_position": r.Position,
			})
			updated++
			continue
		}

		// Not found — insert with a stable dedup key so re-runs are idempotent.
		var winnerID *uint
		if wid := teamLookup[r.WinnerName]; wid != 0 {
			winnerID = &wid
		}
		dedupKey := fmt.Sprintf("bracket_patch:%s:%s:%d", dbSlug, r.CanonicalRound, r.Position)
		m := database.Match{
			TournamentID:    tournamentID,
			Team1ID:         team1ID,
			Team2ID:         team2ID,
			MatchDate:       parseFlexDate(r.MatchDate),
			Team1Score:      r.Team1Score,
			Team2Score:      r.Team2Score,
			WinnerID:        winnerID,
			BracketRound:    r.CanonicalRound,
			BracketPosition: r.Position,
			LiquipediaURL:   dedupKey,
		}
		db.Where("liquipedia_url = ?", dedupKey).FirstOrCreate(&m)
		inserted++
	}

	return
}
