package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func readCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("WARN: cannot open %s: %v", path, err)
		return nil
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.LazyQuotes = true
	records, err := r.ReadAll()
	if err != nil {
		log.Printf("WARN: cannot read %s: %v", path, err)
		return nil
	}
	return records
}

func headerIndex(row []string) map[string]int {
	m := map[string]int{}
	for i, h := range row {
		m[strings.ToLower(strings.TrimSpace(h))] = i
	}
	return m
}

func col(rec []string, h map[string]int, name string) string {
	i, ok := h[name]
	if !ok || i >= len(rec) {
		return ""
	}
	return strings.TrimSpace(rec[i])
}

func colBool(rec []string, h map[string]int, name string) bool {
	v := strings.ToLower(col(rec, h, name))
	return v == "true" || v == "1" || v == "yes"
}

func readBrandingCSV(path string) []brandingRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []brandingRow
	for _, r := range records[1:] {
		rows = append(rows, brandingRow{
			GameCode:          col(r, h, "game_code"),
			SeasonYear:        col(r, h, "season_year"),
			RawTeamName:       col(r, h, "raw_team_name"),
			CanonicalTeamName: col(r, h, "canonical_team_name"),
			TeamSlug:          col(r, h, "team_slug"),
			FranchiseKey:      col(r, h, "franchise_key"),
			ValidFrom:         col(r, h, "valid_from"),
			ValidTo:           col(r, h, "valid_to"),
			Notes:             col(r, h, "notes"),
		})
	}
	return rows
}

func readNonCDLCSV(path string) []nonCDLRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []nonCDLRow
	for _, r := range records[1:] {
		rows = append(rows, nonCDLRow{
			RawTeamName:        col(r, h, "raw_team_name"),
			CanonicalTeamName:  col(r, h, "canonical_team_name"),
			TeamSlug:           col(r, h, "team_slug"),
			IsCDLFranchise:     colBool(r, h, "is_cdl_franchise"),
			TeamClassification: col(r, h, "team_classification"),
			Region:             col(r, h, "region"),
			LinkedCDLTeam:      col(r, h, "linked_cdl_team"),
			RelationshipType:   col(r, h, "relationship_type"),
			DoNotMerge:         colBool(r, h, "do_not_merge"),
			NeedsManualReview:  colBool(r, h, "needs_manual_review"),
			Notes:              col(r, h, "notes"),
		})
	}
	return rows
}

func readPlayerAliasCSV(path string) []playerAliasRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []playerAliasRow
	for _, r := range records[1:] {
		rows = append(rows, playerAliasRow{
			PlayerName:          col(r, h, "player_name"),
			CanonicalPlayerName: col(r, h, "canonical_player_name"),
			PlayerSlug:          col(r, h, "player_slug"),
			TeamContext:         col(r, h, "team_context"),
			GameCode:            col(r, h, "game_code"),
			SeasonYear:          col(r, h, "season_year"),
			NeedsManualReview:   colBool(r, h, "needs_manual_review"),
			Notes:               col(r, h, "notes"),
		})
	}
	return rows
}

func readEventAliasCSV(path string) []eventAliasRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []eventAliasRow
	for _, r := range records[1:] {
		rows = append(rows, eventAliasRow{
			RawEventName:         col(r, h, "raw_event_name"),
			CanonicalEventName:   col(r, h, "canonical_event_name"),
			EventSlug:            col(r, h, "event_slug"),
			GameCode:             col(r, h, "game_code"),
			SeasonYear:           col(r, h, "season_year"),
			EventType:            col(r, h, "event_type"),
			StartDate:            col(r, h, "start_date"),
			EndDate:              col(r, h, "end_date"),
			SourceURL:            col(r, h, "source_url"),
			FandomURL:            col(r, h, "fandom_url"),
			StatsOnly:            colBool(r, h, "stats_only"),
			HasBracket:           colBool(r, h, "has_bracket"),
			HasQualifierStage:    colBool(r, h, "has_qualifier_stage"),
			HasStats:             colBool(r, h, "has_stats"),
		})
	}
	return rows
}

func readSeriesCSV(path string) []seriesRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []seriesRow
	for _, r := range records[1:] {
		rows = append(rows, seriesRow{
			MatchID:       atoi(col(r, h, "match_id")),
			SourceURL:     col(r, h, "source_url"),
			MatchDatetime: col(r, h, "match_datetime"),
			BestOf:        atoi(col(r, h, "best_of")),
			Status:        col(r, h, "status"),
			TeamAID:       atoi(col(r, h, "team_a_id")),
			TeamAName:     col(r, h, "team_a_name"),
			TeamBID:       atoi(col(r, h, "team_b_id")),
			TeamBName:     col(r, h, "team_b_name"),
			TeamAScore:    atoi(col(r, h, "team_a_score")),
			TeamBScore:    atoi(col(r, h, "team_b_score")),
			WinnerID:      atoi(col(r, h, "winner_id")),
			WinnerName:    col(r, h, "winner_name"),
			RoundName:     col(r, h, "bp_round_name"),
			SeriesFormat:  col(r, h, "series_format"),
			SourceType:    col(r, h, "source_type"),
		})
	}
	return rows
}

func readMapCSV(path string) []mapRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []mapRow
	for _, r := range records[1:] {
		rows = append(rows, mapRow{
			MatchID:     atoi(col(r, h, "match_id")),
			MapNumber:   atoi(col(r, h, "map_number")),
			MapName:     col(r, h, "map_name"),
			ModeName:    col(r, h, "mode_name"),
			TeamAID:     atoi(col(r, h, "team_a_id")),
			TeamBID:     atoi(col(r, h, "team_b_id")),
			ScoreA:      atoi(col(r, h, "score_a")),
			ScoreB:      atoi(col(r, h, "score_b")),
			WinnerID:    atoi(col(r, h, "winner_id")),
			WinnerName:  col(r, h, "winner_name"),
			Played:      strings.ToLower(col(r, h, "played")) == "true",
			DurationMin: atoi(col(r, h, "duration_min")),
			DurationSec: atoi(col(r, h, "duration_sec")),
			SourceType:  col(r, h, "source_type"),
		})
	}
	return rows
}

func readPlayerStatCSV(path string) []playerStatRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []playerStatRow
	for _, r := range records[1:] {
		rows = append(rows, playerStatRow{
			MatchID:              atoi(col(r, h, "match_id")),
			MapNumber:            atoi(col(r, h, "map_number")),
			PlayerID:             atoi(col(r, h, "player_id")),
			PlayerTag:            col(r, h, "player_tag"),
			TeamID:               atoi(col(r, h, "team_id")),
			Kills:                atoi(col(r, h, "kills")),
			Deaths:               atoi(col(r, h, "deaths")),
			KD:                   atof(col(r, h, "kd")),
			Damage:               atoi(col(r, h, "damage")),
			Assists:              atoi(col(r, h, "assists")),
			BPRating:             atof(col(r, h, "bp_rating")),
			HillTime:             atoi(col(r, h, "hill_time")),
			SndRounds:            atoi(col(r, h, "snd_rounds")),
			PlantCount:           atoi(col(r, h, "plant_count")),
			DefuseCount:          atoi(col(r, h, "defuse_count")),
			SnipeCount:           atoi(col(r, h, "snipe_count")),
			FirstBloodCount:      atoi(col(r, h, "first_blood_count")),
			FirstDeathCount:      atoi(col(r, h, "first_death_count")),
			ZoneTierCaptureCount: atoi(col(r, h, "zone_tier_capture_count")),
			CtlAttackRounds:      atoi(col(r, h, "ctl_attack_rounds")),
			CtlDefenseRounds:     atoi(col(r, h, "ctl_defense_rounds")),
			NonTradedKills:       atoi(col(r, h, "non_traded_kills")),
			HighestStreak:        atoi(col(r, h, "highest_streak")),
			DataQualityNote:      col(r, h, "data_quality_note"),
			SourceType:           col(r, h, "source_type"),
		})
	}
	return rows
}

func readEnrichedSeriesCSV(path string) []enrichedSeriesRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []enrichedSeriesRow
	for _, r := range records[1:] {
		rows = append(rows, enrichedSeriesRow{
			SeriesMatchID:     col(r, h, "series_match_id"),
			EventName:         col(r, h, "event_name"),
			EventSlug:         col(r, h, "event_slug"),
			GameCode:          col(r, h, "game_code"),
			SeasonYear:        col(r, h, "season_year"),
			StageName:         col(r, h, "stage_name"),
			StageType:         col(r, h, "stage_type"),
			GroupName:         col(r, h, "group_name"),
			RoundName:         col(r, h, "round_name"),
			MatchDatetime:     col(r, h, "match_datetime"),
			Team1:             col(r, h, "team_1"),
			Team2:             col(r, h, "team_2"),
			Team1Canonical:    col(r, h, "team_1_canonical"),
			Team2Canonical:    col(r, h, "team_2_canonical"),
			Team1FranchiseKey: col(r, h, "team_1_franchise_key"),
			Team2FranchiseKey: col(r, h, "team_2_franchise_key"),
			Team1MapWins:      atoi(col(r, h, "team_1_map_wins")),
			Team2MapWins:      atoi(col(r, h, "team_2_map_wins")),
			MapsPlayed:        atoi(col(r, h, "maps_played")),
			Winner:            col(r, h, "winner"),
			WinnerCanonical:   col(r, h, "winner_canonical"),
			SeriesFormat:      col(r, h, "series_format"),
			Source:            col(r, h, "source"),
			SourceURL:         col(r, h, "source_url"),
		})
	}
	return rows
}

func readEnrichedMapCSV(path string) []enrichedMapRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []enrichedMapRow
	for _, r := range records[1:] {
		rows = append(rows, enrichedMapRow{
			SeriesMatchID: col(r, h, "series_match_id"),
			EventSlug:     col(r, h, "event_slug"),
			GameCode:      col(r, h, "game_code"),
			SeasonYear:    col(r, h, "season_year"),
			MapNumber:     atoi(col(r, h, "map_number")),
			Team1:         col(r, h, "team_1"),
			Team2:         col(r, h, "team_2"),
			MapName:       col(r, h, "map_name"),
			Mode:          col(r, h, "mode"),
			Score1:        atoi(col(r, h, "score_1")),
			Score2:        atoi(col(r, h, "score_2")),
			MapWinner:     col(r, h, "map_winner"),
			Played:        strings.ToLower(col(r, h, "played")) == "true",
			Duration:      col(r, h, "duration"),
			Source:        col(r, h, "source"),
		})
	}
	return rows
}

func readEnrichedStatCSV(path string) []enrichedStatRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []enrichedStatRow
	for _, r := range records[1:] {
		rows = append(rows, enrichedStatRow{
			SeriesMatchID:   col(r, h, "series_match_id"),
			MapNumber:       atoi(col(r, h, "map_number")),
			EventSlug:       col(r, h, "event_slug"),
			GameCode:        col(r, h, "game_code"),
			SeasonYear:      col(r, h, "season_year"),
			MapName:         col(r, h, "map_name"),
			Mode:            col(r, h, "mode"),
			Team:            col(r, h, "team"),
			Player:          col(r, h, "player"),
			Kills:           atoi(col(r, h, "kills")),
			Deaths:          atoi(col(r, h, "deaths")),
			KD:              atof(col(r, h, "kd")),
			HillTime:        atoi(col(r, h, "hill_time")),
			Captures:        atoi(col(r, h, "captures")),
			Plants:          atoi(col(r, h, "plants")),
			Defuses:         atoi(col(r, h, "defuses")),
			FirstKills:      atoi(col(r, h, "first_kills")),
			FirstDeaths:     atoi(col(r, h, "first_deaths")),
			DataQualityNote: col(r, h, "data_quality_note"),
			Source:          col(r, h, "source"),
		})
	}
	return rows
}

func readBracketCSV(path string) []cwBracketRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []cwBracketRow
	for _, r := range records[1:] {
		rows = append(rows, cwBracketRow{
			TournamentSlug:  col(r, h, "tournament_slug"),
			SourceRoundName: col(r, h, "source_round_name"),
			CanonicalRound:  col(r, h, "canonical_round_key"),
			Position:        atoi(col(r, h, "bracket_position")),
			Team1Name:       col(r, h, "team1_name"),
			Team2Name:       col(r, h, "team2_name"),
			Team1Score:      atoi(col(r, h, "team1_score")),
			Team2Score:      atoi(col(r, h, "team2_score")),
			WinnerName:      col(r, h, "winner_name"),
			MatchDate:       col(r, h, "match_date"),
		})
	}
	return rows
}

func readTransferCSV(path string) []transferRow {
	records := readCSV(path)
	if len(records) < 2 {
		return nil
	}
	h := headerIndex(records[0])
	var rows []transferRow
	for _, r := range records[1:] {
		rows = append(rows, transferRow{
			Date:         col(r, h, "date"),
			Player:       col(r, h, "player"),
			FromTeam:     col(r, h, "from_team"),
			ToTeam:       col(r, h, "to_team"),
			Role:         col(r, h, "role"),
			TransferType: col(r, h, "transfer_type"),
		})
	}
	return rows
}
