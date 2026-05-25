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
	var primaryMatches, enrichedMatches []matchDetail

	db.Table("matches").
		Where("tournament_id = ? AND liquipedia_url LIKE 'https://www.breakingpoint%'", ewcTournamentID).
		Select("id, team1_id, team2_id, team1_score, team2_score, bracket_round, bracket_position, liquipedia_url").
		Scan(&primaryMatches)

	db.Table("matches").
		Where("tournament_id = ? AND liquipedia_url LIKE 'enriched:EWC2025%'", ewcTournamentID).
		Select("id, team1_id, team2_id, team1_score, team2_score, bracket_round, bracket_position, liquipedia_url").
		Scan(&enrichedMatches)

	log.Printf("[ewc2025] %d era_finals (source URL) matches, %d enriched matches", len(primaryMatches), len(enrichedMatches))

	// Build a normalised key for matching: use LEAST/GREATEST on team IDs + scores.
	type matchKey struct{ t1, t2, s1, s2 uint }
	norm := func(m matchDetail) matchKey {
		if m.Team1ID < m.Team2ID {
			return matchKey{m.Team1ID, m.Team2ID, uint(m.Team1Score), uint(m.Team2Score)}
		}
		return matchKey{m.Team2ID, m.Team1ID, uint(m.Team2Score), uint(m.Team1Score)}
	}

	primaryByKey := map[matchKey]matchDetail{}
	for _, m := range primaryMatches {
		primaryByKey[norm(m)] = m
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
		if primary, ok := primaryByKey[k]; ok {
			pairs = append(pairs, pair{keep: primary, del: e, hasDup: true})
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

	teamLookup := buildTeamLookup(db)
	rows := readEWC2025CSV()
	log.Printf("[ewc2025] %d enriched EWC 2025 rows from CSV", len(rows))

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

const torontoUltraName = "Toronto Ultra"
const torontoKoiName   = "Toronto Koi"

func runTorontoKoiRebrand() {
	database.ConnectDatabase()
	db := database.DB

	var ultra database.Team
	if err := db.Where("name = ?", torontoUltraName).First(&ultra).Error; err != nil {
		log.Fatalf("Toronto Ultra not found: %v", err)
	}
	log.Printf("[toronto] Ultra: id=%d franchise_id=%v abbr=%s", ultra.ID, ultra.FranchiseID, ultra.Abbreviation)

	var koi database.Team
	if err := db.Where("name = ?", torontoKoiName).First(&koi).Error; err != nil {
		koi = database.Team{
			Name:               torontoKoiName,
			Abbreviation:       "TK",
			FranchiseID:        ultra.FranchiseID,
			IsCDLFranchise:     ultra.IsCDLFranchise,
			TeamClassification: ultra.TeamClassification,
			IsActive:           false,
			Source:             "ewc_rebrand",
		}
		if err := db.Create(&koi).Error; err != nil {
			log.Fatalf("Failed to create Toronto Koi: %v", err)
		}
		log.Printf("[toronto] Created Toronto Koi: id=%d", koi.ID)
	} else {
		log.Printf("[toronto] Found existing Toronto Koi: id=%d", koi.ID)
	}

	var tournamentIDs []uint
	db.Table("tournaments").Where("tournament_type = ?", "international_major").Pluck("id", &tournamentIDs)
	log.Printf("[toronto] EWC tournaments: %v", tournamentIDs)

	if len(tournamentIDs) == 0 {
		log.Println("[toronto] No EWC tournaments found — nothing to do.")
		return
	}

	var matchIDs []uint
	db.Table("matches").
		Where("tournament_id IN ? AND (team1_id = ? OR team2_id = ?)", tournamentIDs, ultra.ID, ultra.ID).
		Pluck("id", &matchIDs)
	log.Printf("[toronto] Affected matches: %v", matchIDs)

	type matchRow struct {
		ID           uint
		TournamentID uint
		Team1ID      uint
		Team2ID      uint
		WinnerID     *uint
		BracketRound string
	}
	var rows []matchRow
	db.Table("matches").Where("id IN ?", matchIDs).
		Select("id, tournament_id, team1_id, team2_id, winner_id, bracket_round").
		Scan(&rows)
	fmt.Printf("\n%-6s %-6s %-12s %-6s %-6s %-6s\n", "id", "t_id", "bracket_round", "tm1", "tm2", "win")
	for _, r := range rows {
		win := "nil"
		if r.WinnerID != nil {
			win = fmt.Sprintf("%d", *r.WinnerID)
		}
		fmt.Printf("%-6d %-6d %-12s %-6d %-6d %-6s\n",
			r.ID, r.TournamentID, r.BracketRound, r.Team1ID, r.Team2ID, win)
	}
	fmt.Println()

	if len(matchIDs) == 0 {
		log.Println("[toronto] No matches to update — nothing to do.")
		return
	}

	res := db.Table("matches").Where("id IN ? AND team1_id = ?", matchIDs, ultra.ID).Update("team1_id", koi.ID)
	log.Printf("[toronto] matches.team1_id: %d rows", res.RowsAffected)
	res = db.Table("matches").Where("id IN ? AND team2_id = ?", matchIDs, ultra.ID).Update("team2_id", koi.ID)
	log.Printf("[toronto] matches.team2_id: %d rows", res.RowsAffected)
	res = db.Table("matches").Where("id IN ? AND winner_id = ?", matchIDs, ultra.ID).Update("winner_id", koi.ID)
	log.Printf("[toronto] matches.winner_id: %d rows", res.RowsAffected)
	res = db.Table("match_maps").Where("match_id IN ? AND winner_id = ?", matchIDs, ultra.ID).Update("winner_id", koi.ID)
	log.Printf("[toronto] match_maps.winner_id: %d rows", res.RowsAffected)
	res = db.Table("player_map_stats").Where("match_id IN ? AND team_id = ?", matchIDs, ultra.ID).Update("team_id", koi.ID)
	log.Printf("[toronto] player_map_stats: %d rows", res.RowsAffected)
	res = db.Table("player_match_stats").Where("match_id IN ? AND team_id = ?", matchIDs, ultra.ID).Update("team_id", koi.ID)
	log.Printf("[toronto] player_match_stats: %d rows", res.RowsAffected)
	res = db.Table("player_tournament_stats").Where("tournament_id IN ? AND team_id = ?", tournamentIDs, ultra.ID).Update("team_id", koi.ID)
	log.Printf("[toronto] player_tournament_stats: %d rows", res.RowsAffected)
	res = db.Table("team_tournament_stats").Where("tournament_id IN ? AND team_id = ?", tournamentIDs, ultra.ID).Update("team_id", koi.ID)
	log.Printf("[toronto] team_tournament_stats: %d rows", res.RowsAffected)

	log.Printf("[toronto] Done. %d EWC matches rebranded to Toronto Koi (id=%d).", len(matchIDs), koi.ID)
}

var ewc2024Positions = map[uint]int{
	// Group A (OG / CRR / OB / LVL)
	1179: 1, // opening:  OG  vs OB
	1180: 1, // opening:  LVL vs CRR
	1181: 1, // elim:     OB  vs LVL
	1182: 1, // winners:  OG  vs CRR
	1183: 1, // decider:  CRR vs LVL

	// Group B (VS / C / GM / GE)
	1184: 2, // opening:  C   vs GM
	1185: 2, // opening:  VS  vs GE
	1186: 2, // elim:     GM  vs GE
	1187: 2, // winners:  C   vs VS
	1188: 2, // decider:  C   vs GM

	// Group C (TU / BB / S / TH)
	1189: 3, // opening:  TU  vs S
	1190: 3, // opening:  TH  vs BB
	1191: 3, // elim:     S   vs TH
	1192: 3, // winners:  TU  vs BB
	1193: 3, // decider:  BB  vs S

	// Group D (AF / T / LG / TF)
	1194: 4, // opening:  T   vs LG
	1195: 4, // opening:  AF  vs TF
	1196: 4, // elim:     LG  vs TF
	1197: 4, // winners:  T   vs AF
	1198: 4, // decider:  T   vs LG
}

func runEWC2024PositionPatch() {
	database.ConnectDatabase()
	db := database.DB

	fmt.Println("=== EWC 2024 group-stage position patch (dry-run first) ===")

	type row struct {
		ID              uint
		Team1ID, Team2ID uint
		BracketRound    string
		BracketPosition int
	}
	var current []row
	ids := make([]uint, 0, len(ewc2024Positions))
	for id := range ewc2024Positions {
		ids = append(ids, id)
	}
	db.Table("matches").
		Where("id IN ?", ids).
		Select("id, team1_id, team2_id, bracket_round, bracket_position").
		Order("bracket_round, id").
		Scan(&current)

	fmt.Printf("Found %d matches to patch (expected 20).\n\n", len(current))
	for _, r := range current {
		want := ewc2024Positions[r.ID]
		flag := ""
		if r.BracketPosition == want {
			flag = " (already correct)"
		}
		fmt.Printf("  id=%-5d  %-20s  pos %d → %d%s\n",
			r.ID, r.BracketRound, r.BracketPosition, want, flag)
	}

	if len(current) != 20 {
		fmt.Printf("\nERROR: expected 20 matches, found %d — aborting.\n", len(current))
		return
	}

	fmt.Println("\nApplying bracket_position updates...")
	updated := 0
	for _, r := range current {
		want := ewc2024Positions[r.ID]
		if r.BracketPosition == want {
			continue
		}
		res := db.Table("matches").
			Where("id = ? AND tournament_id = 53", r.ID).
			Update("bracket_position", want)
		if res.Error != nil {
			log.Printf("ERROR updating id=%d: %v", r.ID, res.Error)
		} else if res.RowsAffected == 1 {
			updated++
		}
	}

	fmt.Printf("\nDone. %d rows updated.\n", updated)
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
