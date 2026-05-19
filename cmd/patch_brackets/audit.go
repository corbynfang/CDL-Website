package main

import (
	"fmt"
	"strings"

	"github.com/corbynfang/CDL-Website/internal/database"
)

func runAudit() {
	database.ConnectDatabase()
	db := database.DB

	var teams []struct {
		ID   uint
		Name string
	}
	db.Table("teams").Where("name ILIKE '%minnesota%' OR name ILIKE '%g2 minn%'").Scan(&teams)
	fmt.Println("=== Minnesota team records ===")
	for _, t := range teams {
		fmt.Printf("  id=%d  name=%q\n", t.ID, t.Name)
	}

	var skipped []struct{ ID uint; Name string }
	db.Table("teams").Where("name ILIKE '%vegas%' OR name ILIKE '%carolina royal%' OR name ILIKE '%guerrillas m8%' OR name ILIKE '%cloud9%'").Scan(&skipped)
	fmt.Println("\n=== Skipped team names in DB ===")
	for _, t := range skipped {
		fmt.Printf("  id=%d  name=%q\n", t.ID, t.Name)
	}

	var patchCount int64
	db.Table("matches").Where("liquipedia_url LIKE 'bracket_patch:%'").Count(&patchCount)
	fmt.Printf("\n=== bracket_patch: rows in matches: %d ===\n", patchCount)

	// Per-tournament bracket position audit
	fmt.Println("\n=== bracket_position=0 audit per non-CW major ===")
	type auditRow struct {
		Slug    string
		Total   int
		PosZero int
	}
	var auditRows []auditRow
	db.Raw(`
		SELECT t.slug, COUNT(*) AS total,
		       SUM(CASE WHEN m.bracket_position = 0 THEN 1 ELSE 0 END) AS pos_zero
		FROM matches m
		JOIN tournaments t ON t.id = m.tournament_id
		WHERE t.slug LIKE 'cdl-major-%'
		  AND m.bracket_round NOT IN ('', 'major_qualifier')
		GROUP BY t.slug
		ORDER BY t.slug
	`).Scan(&auditRows)
	for _, r := range auditRows {
		flag := ""
		if r.PosZero > 0 {
			flag = " ← STILL HAS ZEROS"
		}
		fmt.Printf("  %-40s total=%3d  pos_zero=%d%s\n", r.Slug, r.Total, r.PosZero, flag)
	}

	// Duplicate match detection — normalized key:
	//   tournament_id + LEAST/GREATEST(team_ids) + scores normalized to follow lower team_id
	//   + winner_id + bracket_round
	//
	// Scores are normalized so that score_lower / score_higher correspond to the lower/higher
	// team ID respectively. This prevents false positives where the same two teams produce the
	// same raw (team1_score, team2_score) values in different orientations.
	//
	// bracket_round is included so that teams legitimately meeting twice in the same tournament
	// (e.g. EWC group stage opening_match then decider_match) are never flagged.
	fmt.Println("\n=== Duplicate match check (normalized) ===")
	type dupRow struct {
		TournamentID   uint
		Slug           string
		LowerTeamID    uint
		HigherTeamID   uint
		LowerScore     int
		HigherScore    int
		WinnerID       *uint
		BracketRound   string
		Count          int
	}
	var dupRows []dupRow
	db.Raw(`
		SELECT
		    m.tournament_id,
		    t.slug,
		    LEAST(m.team1_id, m.team2_id)    AS lower_team_id,
		    GREATEST(m.team1_id, m.team2_id)  AS higher_team_id,
		    CASE WHEN m.team1_id <= m.team2_id THEN m.team1_score ELSE m.team2_score END AS lower_score,
		    CASE WHEN m.team1_id <= m.team2_id THEN m.team2_score ELSE m.team1_score END AS higher_score,
		    m.winner_id,
		    m.bracket_round,
		    COUNT(*) AS count
		FROM matches m
		JOIN tournaments t ON t.id = m.tournament_id
		WHERE m.bracket_round NOT IN ('', 'major_qualifier')
		GROUP BY
		    m.tournament_id, t.slug,
		    LEAST(m.team1_id, m.team2_id),
		    GREATEST(m.team1_id, m.team2_id),
		    CASE WHEN m.team1_id <= m.team2_id THEN m.team1_score ELSE m.team2_score END,
		    CASE WHEN m.team1_id <= m.team2_id THEN m.team2_score ELSE m.team1_score END,
		    m.winner_id,
		    m.bracket_round
		HAVING COUNT(*) > 1
		ORDER BY m.tournament_id, lower_score
	`).Scan(&dupRows)
	if len(dupRows) == 0 {
		fmt.Println("  No true duplicates found.")
	}
	for _, d := range dupRows {
		winnerStr := "nil"
		if d.WinnerID != nil {
			winnerStr = fmt.Sprintf("%d", *d.WinnerID)
		}
		isFallback := strings.Contains(d.Slug, "unmatched")
		label := "DUPLICATE"
		if isFallback {
			label = "FALLBACK-COLLISION"
		}
		fmt.Printf("  [%s] tournament %d (%s): teams %d vs %d  score %d-%d  winner=%s  rnd=%s  count=%d\n",
			label, d.TournamentID, d.Slug,
			d.LowerTeamID, d.HigherTeamID, d.LowerScore, d.HigherScore,
			winnerStr, d.BracketRound, d.Count)
		if isFallback {
			fmt.Printf("    ↳ matches from different CDL events fell into the fallback tournament; fix by correcting event date ranges in event_aliases_clean.csv\n")
		}
	}

	// What is tournament 52?
	var t52 struct{ ID uint; Slug, Name string }
	db.Table("tournaments").Where("id = 52").Scan(&t52)
	fmt.Printf("\n=== tournament 52: slug=%q ===\n", t52.Slug)

	// Breakdown of tournament 52 matches by dedup key type
	var withBPID, withEnriched, withBPURL, withOther int64
	db.Table("matches").Where("tournament_id = 52 AND breaking_point_match_id IS NOT NULL").Count(&withBPID)
	db.Table("matches").Where("tournament_id = 52 AND liquipedia_url LIKE 'enriched:EWC2025%'").Count(&withEnriched)
	db.Table("matches").Where("tournament_id = 52 AND liquipedia_url LIKE 'https://www.breakingpoint%'").Count(&withBPURL)
	db.Table("matches").Where("tournament_id = 52 AND liquipedia_url NOT LIKE 'enriched:EWC2025%' AND liquipedia_url NOT LIKE 'https://www.breakingpoint%'").Count(&withOther)
	fmt.Printf("  bp_id IS NOT NULL: %d\n  enriched:EWC2025: %d\n  breakingpoint URL: %d\n  other: %d\n", withBPID, withEnriched, withBPURL, withOther)
	type t52Row struct{ ID uint; LiquipediaURL string; BPID *int; BracketRound string }
	var t52Rows []t52Row
	db.Table("matches").Where("tournament_id = 52").Select("id, liquipedia_url, breaking_point_match_id, bracket_round").Limit(6).Scan(&t52Rows)
	fmt.Println("  sample rows:")
	for _, r := range t52Rows {
		bpid := "nil"
		if r.BPID != nil { bpid = fmt.Sprintf("%d", *r.BPID) }
		fmt.Printf("    id=%d bp=%s rnd=%q url=%q\n", r.ID, bpid, r.BracketRound, r.LiquipediaURL)
	}

	// G2 Minnesota bracket_patch inserts from this run
	type matchRow struct {
		ID, Team1ID, Team2ID uint
		Team1Score, Team2Score, BracketPosition int
		BracketRound, LiquipediaURL string
	}
	var g2Inserts []matchRow
	db.Table("matches").
		Where("(team1_id = 85 OR team2_id = 85) AND liquipedia_url LIKE 'bracket_patch:%'").
		Select("id, team1_id, team2_id, team1_score, team2_score, bracket_round, bracket_position, liquipedia_url").
		Order("liquipedia_url").Scan(&g2Inserts)
	fmt.Printf("\n=== G2 Minnesota bracket_patch inserts (%d) ===\n", len(g2Inserts))
	for _, m := range g2Inserts {
		fmt.Printf("  id=%-5d %d vs %d s=%d-%d rnd=%s pos=%d\n",
			m.ID, m.Team1ID, m.Team2ID, m.Team1Score, m.Team2Score, m.BracketRound, m.BracketPosition)
	}

	// Do originals exist with RØKKR (id=17)?
	var rokkrBracket int64
	db.Table("matches").
		Joins("JOIN tournaments t ON t.id = matches.tournament_id").
		Where("(matches.team1_id = 17 OR matches.team2_id = 17) AND t.slug LIKE 'cdl-major-%-tournament-202%' AND matches.bracket_round NOT IN ('', 'major_qualifier')").
		Count(&rokkrBracket)
	fmt.Printf("\n=== Minnesota RØKKR (id=17) bracket matches in majors: %d ===\n", rokkrBracket)

	// Spot-check the 4 requested tournament IDs
	fmt.Println("\n=== Spot-check requested tournament IDs ===")
	for _, tid := range []uint{14, 36, 46, 52} {
		var slug string
		db.Table("tournaments").Where("id = ?", tid).Select("slug").Scan(&slug)
		var total, posZero int64
		db.Table("matches").Where("tournament_id = ? AND bracket_round NOT IN ('', 'major_qualifier')", tid).Count(&total)
		db.Table("matches").Where("tournament_id = ? AND bracket_round NOT IN ('', 'major_qualifier') AND bracket_position = 0", tid).Count(&posZero)
		fmt.Printf("  tournament %-2d (%s): bracket_matches=%d  pos_zero=%d\n", tid, slug, total, posZero)
	}
}
