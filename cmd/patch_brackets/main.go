package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

var bracketSlugAliases = map[string]string{
	"cdl-2021-stage-1-major": "cdl-major-1-tournament-2021",
	"cdl-2021-stage-2-major": "cdl-major-2-tournament-2021",
	"cdl-2021-stage-3-major": "cdl-major-3-tournament-2021",
	"cdl-2021-stage-4-major": "cdl-major-4-tournament-2021",
	"cdl-2021-stage-5-major": "cdl-major-5-tournament-2021",
}

var bracketCSVs = []string{
	"database/cdl_cw_stage_brackets.csv",
	"database/cdl_champs_brackets.csv",
	"database/cdl_major_brackets.csv",
}

type bracketRow struct {
	TournamentSlug string
	SourceRound    string
	CanonicalRound string
	Position       int
	Team1Name      string
	Team2Name      string
	Team1Score     int
	Team2Score     int
	WinnerName     string
	MatchDate      string
}

func readBracketCSV(path string) []bracketRow {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}
	if len(records) < 2 {
		return nil
	}
	header := records[0]
	col := func(record []string, name string) string {
		for i, h := range header {
			if h == name {
				if i < len(record) {
					return record[i]
				}
			}
		}
		return ""
	}
	atoi := func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}
	var rows []bracketRow
	for _, rec := range records[1:] {
		rows = append(rows, bracketRow{
			TournamentSlug: col(rec, "tournament_slug"),
			SourceRound:    col(rec, "source_round_name"),
			CanonicalRound: col(rec, "canonical_round_key"),
			Position:       atoi(col(rec, "bracket_position")),
			Team1Name:      col(rec, "team1_name"),
			Team2Name:      col(rec, "team2_name"),
			Team1Score:     atoi(col(rec, "team1_score")),
			Team2Score:     atoi(col(rec, "team2_score")),
			WinnerName:     col(rec, "winner_name"),
			MatchDate:      col(rec, "match_date"),
		})
	}
	return rows
}

func buildTeamLookup(db *gorm.DB) map[string]uint {
	var teams []struct {
		ID   uint
		Name string
	}
	db.Table("teams").Select("id, name").Scan(&teams)
	m := map[string]uint{}
	for _, t := range teams {
		m[t.Name] = t.ID
	}
	log.Printf("[lookup] %d teams loaded", len(m))
	return m
}

func buildTournamentLookup(db *gorm.DB) map[string]uint {
	var tournaments []struct {
		ID   uint
		Slug string
	}
	db.Table("tournaments").Select("id, slug").Scan(&tournaments)
	m := map[string]uint{}
	for _, t := range tournaments {
		m[t.Slug] = t.ID
	}
	log.Printf("[lookup] %d tournaments loaded", len(m))
	return m
}

func applyBracketCSV(db *gorm.DB, teamLookup map[string]uint, tournamentBySlug map[string]uint, path string) (updated, inserted, skipped int) {
	rows := readBracketCSV(path)
	log.Printf("[%s] %d rows to process", path, len(rows))

	for _, r := range rows {
		dbSlug := r.TournamentSlug
		if alias, ok := bracketSlugAliases[r.TournamentSlug]; ok {
			dbSlug = alias
		}
		tournamentID := tournamentBySlug[dbSlug]
		if tournamentID == 0 {
			log.Printf("WARN: tournament not found for slug %q — skipping", dbSlug)
			skipped++
			continue
		}
		team1ID := teamLookup[r.Team1Name]
		team2ID := teamLookup[r.Team2Name]
		if team1ID == 0 || team2ID == 0 {
			log.Printf("WARN: team not found (%q or %q) — skipping", r.Team1Name, r.Team2Name)
			skipped++
			continue
		}

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
			if existing.BracketRound == r.CanonicalRound && existing.BracketPosition == r.Position {
				// Already correct — count as updated (idempotent re-run)
				updated++
				continue
			}
			db.Model(&existing).Updates(map[string]interface{}{
				"bracket_round":    r.CanonicalRound,
				"bracket_position": r.Position,
			})
			updated++
			continue
		}

		log.Printf("WARN: no existing match for %s [%s pos=%d] %s vs %s (%d-%d) — inserting",
			dbSlug, r.CanonicalRound, r.Position, r.Team1Name, r.Team2Name, r.Team1Score, r.Team2Score)
		var winnerID *uint
		if wid := teamLookup[r.WinnerName]; wid != 0 {
			winnerID = &wid
		}
		dedupKey := fmt.Sprintf("bracket_patch:%s:%s:%d", dbSlug, r.CanonicalRound, r.Position)
		m := database.Match{
			TournamentID:    tournamentID,
			Team1ID:         team1ID,
			Team2ID:         team2ID,
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

func main() {
	if len(os.Args) > 1 && os.Args[1] == "audit" {
		runAudit()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "cleanup" {
		deleteG2MinnBadInserts()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "ewc2025" {
		runEWC2025Fix()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "ewc2024" {
		runEWC2024PositionPatch()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "toronto-koi" {
		runTorontoKoiRebrand()
		return
	}

	database.ConnectDatabase()
	db := database.DB

	teamLookup := buildTeamLookup(db)
	tournamentBySlug := buildTournamentLookup(db)

	var totalUpdated, totalInserted, totalSkipped int
	for _, path := range bracketCSVs {
		u, i, s := applyBracketCSV(db, teamLookup, tournamentBySlug, path)
		totalUpdated += u
		totalInserted += i
		totalSkipped += s
	}

	log.Printf("[patch_brackets] DONE — updated=%d  inserted=%d  skipped=%d",
		totalUpdated, totalInserted, totalSkipped)
}
