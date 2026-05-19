package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

// deleteG2MinnBadInserts removes the placeholder matches inserted during the first
// patch run when "G2 Minnesota" (team_id=85) was incorrectly looked up instead of
// "Minnesota RØKKR" (team_id=17). These rows have a bracket_patch: dedup key and
// carry no real match data. The historical matches with RØKKR (id=17) remain intact.
func deleteG2MinnBadInserts() {
	database.ConnectDatabase()
	db := database.DB
	result := db.Exec(`
		DELETE FROM matches
		WHERE (team1_id = 85 OR team2_id = 85)
		  AND liquipedia_url LIKE 'bracket_patch:%'
	`)
	if result.Error != nil {
		log.Fatalf("delete failed: %v", result.Error)
	}
	log.Printf("[cleanup] deleted %d G2 Minnesota bracket_patch placeholder rows", result.RowsAffected)
}

// ─── EWC 2025 fix ────────────────────────────────────────────────────────────
//
// EWC 2025 (tournament_id=52) was double-seeded:
//   - phase2 era_finals: stored round names like "group_play_a_winners_round_1" (wrong canonical)
//   - phase3 enriched:   stored "opening_match" without group prefix (correct base name, wrong key)
//
// Fix: delete era_finals BP duplicates, then update enriched records to:
//   - bracket_round  = "group_play_a_opening_match" (group prefix + canonical name)
//   - bracket_position = 1…N within each (round, group) block sorted by series_match_id

const ewcTournamentID = 52

type ewcRow struct {
	SeriesMatchID string
	Group         string // "A"…"D", "" for playoffs
	RoundName     string
	Datetime      string
	Team1         string
	Team2         string
	Score1        int
	Score2        int
}

func ewcCanonicalRound(group, roundName string) string {
	var base string
	switch roundName {
	case "Opening Match":
		base = "opening_match"
	case "Winners Match":
		base = "winners_match"
	case "Elimination Match":
		base = "elimination_match"
	case "Decider Match":
		base = "decider_match"
	case "Quarterfinal":
		return "quarterfinal"
	case "Semifinal":
		return "semifinal"
	case "Third Place Match":
		return "third_place_match"
	case "Grand Final":
		return "grand_finals"
	default:
		base = strings.ToLower(strings.ReplaceAll(roundName, " ", "_"))
	}
	if group != "" {
		return fmt.Sprintf("group_play_%s_%s", strings.ToLower(group), base)
	}
	return base
}

func readEWC2025CSV() []ewcRow {
	f, err := os.Open("database/enriched_series_matches.csv")
	if err != nil {
		log.Fatalf("open enriched_series_matches.csv: %v", err)
	}
	defer f.Close()
	rdr := csv.NewReader(f)
	records, _ := rdr.ReadAll()
	header := records[0]
	col := func(rec []string, name string) string {
		for i, h := range header {
			if h == name && i < len(rec) {
				return rec[i]
			}
		}
		return ""
	}
	atoi := func(s string) int {
		n := 0
		for _, c := range s {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			}
		}
		return n
	}
	var out []ewcRow
	for _, rec := range records[1:] {
		if col(rec, "event_slug") != "esports-world-cup-2025" {
			continue
		}
		out = append(out, ewcRow{
			SeriesMatchID: col(rec, "series_match_id"),
			Group:         col(rec, "group_name"),
			RoundName:     col(rec, "round_name"),
			Datetime:      col(rec, "match_datetime"),
			Team1:         col(rec, "team_1_canonical"),
			Team2:         col(rec, "team_2_canonical"),
			Score1:        atoi(col(rec, "team_1_map_wins")),
			Score2:        atoi(col(rec, "team_2_map_wins")),
		})
	}
	return out
}

type matchDetail struct {
	ID              uint
	Team1ID         uint
	Team2ID         uint
	Team1Score      int
	Team2Score      int
	BracketRound    string
	BracketPosition int
	LiquipediaURL   string
}

func childCounts(db *gorm.DB, matchID uint) (pms, pMapS, mm int64) {
	db.Table("player_match_stats").Where("match_id = ?", matchID).Count(&pms)
	db.Table("player_map_stats").Where("match_id = ?", matchID).Count(&pMapS)
	db.Table("match_maps").Where("match_id = ?", matchID).Count(&mm)
	return
}

func runEWC2025Fix() {
	database.ConnectDatabase()
	db := database.DB

	// Load all EWC 2025 matches split by source.
	var bpMatches, enrichedMatches []matchDetail

	db.Table("matches").
		Where("tournament_id = ? AND liquipedia_url LIKE 'https://www.breakingpoint%'", ewcTournamentID).
		Select("id, team1_id, team2_id, team1_score, team2_score, bracket_round, bracket_position, liquipedia_url").
		Scan(&bpMatches)

	db.Table("matches").
		Where("tournament_id = ? AND liquipedia_url LIKE 'enriched:EWC2025%'", ewcTournamentID).
		Select("id, team1_id, team2_id, team1_score, team2_score, bracket_round, bracket_position, liquipedia_url").
		Scan(&enrichedMatches)

	log.Printf("[ewc2025] %d era_finals (BP URL) matches, %d enriched matches", len(bpMatches), len(enrichedMatches))

	// Build a normalised key for matching: use LEAST/GREATEST on team IDs + scores.
	type matchKey struct{ t1, t2, s1, s2 uint }
	norm := func(m matchDetail) matchKey {
		if m.Team1ID < m.Team2ID {
			return matchKey{m.Team1ID, m.Team2ID, uint(m.Team1Score), uint(m.Team2Score)}
		}
		return matchKey{m.Team2ID, m.Team1ID, uint(m.Team2Score), uint(m.Team1Score)}
	}

	bpByKey := map[matchKey]matchDetail{}
	for _, m := range bpMatches {
		bpByKey[norm(m)] = m
	}

	type pair struct {
		keep   matchDetail // era_finals (has stats) — or enriched if no BP duplicate
		del    matchDetail // enriched duplicate to remove
		hasDup bool
	}
	var pairs []pair
	var enrichedOnly []matchDetail

	for _, e := range enrichedMatches {
		k := norm(e)
		if bp, ok := bpByKey[k]; ok {
			pairs = append(pairs, pair{keep: bp, del: e, hasDup: true})
		} else {
			enrichedOnly = append(enrichedOnly, e)
		}
	}

	// Print plan before mutating anything.
	fmt.Printf("\n=== Duplicate pair plan (%d pairs) ===\n", len(pairs))
	for _, p := range pairs {
		keepPMS, keepPMapS, keepMM := childCounts(db, p.keep.ID)
		delPMS, delPMapS, delMM := childCounts(db, p.del.ID)
		fmt.Printf("  KEEP  id=%-5d (BP URL) pms=%-3d pmap=%-3d mm=%-3d rnd=%q\n",
			p.keep.ID, keepPMS, keepPMapS, keepMM, p.keep.BracketRound)
		fmt.Printf("  DEL   id=%-5d (enrich) pms=%-3d pmap=%-3d mm=%-3d rnd=%q\n",
			p.del.ID, delPMS, delPMapS, delMM, p.del.BracketRound)
		if delPMS > 0 {
			fmt.Printf("  *** enriched match has %d player_match_stats — must migrate before delete\n", delPMS)
		}
		fmt.Println()
	}
	fmt.Printf("=== Enriched-only (no BP duplicate): %d matches ===\n", len(enrichedOnly))
	for _, e := range enrichedOnly {
		fmt.Printf("  id=%-5d rnd=%q\n", e.ID, e.BracketRound)
	}

	// Now load the CSV to compute correct bracket_round + bracket_position.
	teamLookup := buildTeamLookup(db)
	rows := readEWC2025CSV()
	log.Printf("[ewc2025] %d enriched EWC 2025 rows from CSV", len(rows))

	// 2. Compute canonical bracket_round and position for each row.
	//    Position = rank within (canonical_round) sorted by series_match_id (lexicographic = chronological).
	byRound := map[string][]ewcRow{}
	for _, row := range rows {
		cr := ewcCanonicalRound(row.Group, row.RoundName)
		byRound[cr] = append(byRound[cr], row)
	}
	for k := range byRound {
		sort.Slice(byRound[k], func(i, j int) bool {
			return byRound[k][i].SeriesMatchID < byRound[k][j].SeriesMatchID
		})
	}
	type patchInfo struct{ round string; position int }
	patches := map[string]patchInfo{}
	for round, ms := range byRound {
		for pos, m := range ms {
			patches[m.SeriesMatchID] = patchInfo{round, pos + 1}
		}
	}

	// Build lookup: (team1_id, team2_id, score1, score2) → patchInfo for BP kept matches.
	// Also build enriched dedup-key → patchInfo for enriched-only matches.
	type teamKey struct{ t1, t2 uint; s1, s2 int }
	teamKeyToPatch := map[teamKey]patchInfo{}
	enrichedKeyToPatch := map[string]patchInfo{}
	for _, row := range rows {
		t1 := teamLookup[row.Team1]
		t2 := teamLookup[row.Team2]
		if t1 == 0 || t2 == 0 {
			log.Printf("WARN: team not found (%q or %q) for %s", row.Team1, row.Team2, row.SeriesMatchID)
			continue
		}
		p := patches[row.SeriesMatchID]
		teamKeyToPatch[teamKey{t1, t2, row.Score1, row.Score2}] = p
		teamKeyToPatch[teamKey{t2, t1, row.Score2, row.Score1}] = p
		enrichedKeyToPatch["enriched:"+row.SeriesMatchID] = p
	}

	// 3. For each duplicate pair: confirm enriched has no player_match_stats, delete
	//    its child rows and itself, then update the kept BP match with correct round/pos.
	updated, deleted, skipped := 0, 0, 0

	for _, pr := range pairs {
		pms, _, _ := childCounts(db, pr.del.ID)
		if pms > 0 {
			log.Printf("SKIP pair keep=%d del=%d: enriched match has %d player_match_stats — not deleting",
				pr.keep.ID, pr.del.ID, pms)
			skipped++
			continue
		}
		for _, step := range []struct{ table string }{{"player_map_stats"}, {"match_maps"}} {
			if r := db.Exec(`DELETE FROM `+step.table+` WHERE match_id = ?`, pr.del.ID); r.Error != nil {
				log.Printf("WARN: delete %s for match %d: %v", step.table, pr.del.ID, r.Error)
			}
		}
		if r := db.Exec(`DELETE FROM matches WHERE id = ?`, pr.del.ID); r.Error != nil {
			log.Printf("WARN: delete enriched match %d: %v", pr.del.ID, r.Error)
			skipped++
			continue
		}
		deleted++
		p, ok := teamKeyToPatch[teamKey{pr.keep.Team1ID, pr.keep.Team2ID, pr.keep.Team1Score, pr.keep.Team2Score}]
		if !ok {
			log.Printf("WARN: no patch info for kept BP match %d", pr.keep.ID)
			skipped++
			continue
		}
		if r := db.Exec(`UPDATE matches SET bracket_round = ?, bracket_position = ? WHERE id = ?`,
			p.round, p.position, pr.keep.ID); r.Error != nil {
			log.Printf("WARN: update BP match %d: %v", pr.keep.ID, r.Error)
			skipped++
		} else {
			updated++
		}
	}

	// 4. Enriched-only matches (no BP duplicate): update bracket_round / bracket_position.
	for _, e := range enrichedOnly {
		p, ok := enrichedKeyToPatch[e.LiquipediaURL]
		if !ok {
			log.Printf("WARN: no patch info for enriched-only match %d (%s)", e.ID, e.LiquipediaURL)
			skipped++
			continue
		}
		if r := db.Exec(`UPDATE matches SET bracket_round = ?, bracket_position = ? WHERE id = ?`,
			p.round, p.position, e.ID); r.Error != nil {
			log.Printf("WARN: update enriched-only match %d: %v", e.ID, r.Error)
			skipped++
		} else {
			updated++
		}
	}

	log.Printf("[ewc2025] DONE — updated=%d  deleted_enriched=%d  skipped=%d", updated, deleted, skipped)
	printEWC2025Summary(db)
}

func printEWC2025Summary(db *gorm.DB) {
	type sumRow struct {
		BracketRound    string
		BracketPosition int
		Count           int
	}
	var rows []sumRow
	db.Raw(`SELECT bracket_round, bracket_position, COUNT(*) as count
		FROM matches WHERE tournament_id = ?
		GROUP BY bracket_round, bracket_position
		ORDER BY bracket_round, bracket_position`, ewcTournamentID).Scan(&rows)
	fmt.Println("\n=== EWC 2025 bracket_round / position breakdown ===")
	for _, r := range rows {
		flag := ""
		if r.BracketPosition == 0 {
			flag = " ← pos=0"
		}
		fmt.Printf("  %-48s pos=%-2d count=%d%s\n", r.BracketRound, r.BracketPosition, r.Count, flag)
	}
}
