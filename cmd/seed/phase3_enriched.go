package main

// phase3_enriched.go — seeds EWC 2024, EWC 2025, and CDL Major 1 2023 (wiki source).
// These events aren't in era_finals because they either had different data sources or
// different formats (group stages, international teams). Rows sourced from
// "bo6_season_stats_breakingpoint" are skipped — era_finals already covers those.

import (
	"log"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedEnrichedMatches(
	db *gorm.DB,
	teamLookup map[string]uint,
	playerLookup map[string]uint,
	tournamentBySlug map[string]uint,
) {
	seriesRows := readEnrichedSeriesCSV("database/enriched_series_matches.csv")
	mapRows := readEnrichedMapCSV("database/enriched_match_maps.csv")
	statRows := readEnrichedStatCSV("database/enriched_player_map_stats.csv")

	mapsByID := map[string][]enrichedMapRow{}
	for _, mr := range mapRows {
		mapsByID[mr.SeriesMatchID] = append(mapsByID[mr.SeriesMatchID], mr)
	}
	statsByID := map[string][]enrichedStatRow{}
	for _, sr := range statRows {
		statsByID[sr.SeriesMatchID] = append(statsByID[sr.SeriesMatchID], sr)
	}

	seriesSeeded := 0
	var matchMapsBatch []models.MatchMap
	var playerMapStatsBatch []models.PlayerMapStats
	var playerMatchStatsBatch []models.PlayerMatchStats

	for _, s := range seriesRows {
		if s.Source == "bo6_season_stats_breakingpoint" {
			continue
		}

		tournamentID := tournamentBySlug[s.EventSlug]
		if tournamentID == 0 {
			log.Printf("[enriched] WARN: no tournament for slug %q (series %s)", s.EventSlug, s.SeriesMatchID)
			continue
		}

		team1ID := resolveTeamID(teamLookup, s.Team1Canonical, s.GameCode)
		team2ID := resolveTeamID(teamLookup, s.Team2Canonical, s.GameCode)
		if team1ID == 0 {
			team1ID = ensureUnknownTeam(db, s.Team1Canonical, teamLookup)
		}
		if team2ID == 0 {
			team2ID = ensureUnknownTeam(db, s.Team2Canonical, teamLookup)
		}

		var winnerID *uint
		if wid := resolveTeamID(teamLookup, s.WinnerCanonical, s.GameCode); wid != 0 {
			winnerID = &wid
		}

		matchDate := parseFlexDateCtx(s.MatchDatetime, s.SeriesMatchID)
		dedupKey := "enriched:" + s.SeriesMatchID
		m := models.Match{
			TournamentID:  tournamentID,
			Team1ID:       team1ID,
			Team2ID:       team2ID,
			MatchDate:     matchDate,
			Format:        s.SeriesFormat,
			Team1Score:    s.Team1MapWins,
			Team2Score:    s.Team2MapWins,
			WinnerID:      winnerID,
			LiquipediaURL: dedupKey,
			BracketRound:  rawRoundToDBRound(s.RoundName),
		}
		// Match is kept as FirstOrCreate — we need m.ID immediately for child rows.
		db.Where("liquipedia_url = ?", dedupKey).FirstOrCreate(&m)
		// If previously seeded before PST timezone parsing was supported, match_date
		// may be zero in the DB — correct it now that parseFlexDateCtx handles it.
		if !matchDate.IsZero() && m.MatchDate.IsZero() {
			log.Printf("[enriched] correcting match_date for %s: 0001-01-01 → %s", s.SeriesMatchID, matchDate.UTC().Format("2006-01-02"))
			db.Model(&m).Update("match_date", matchDate)
			m.MatchDate = matchDate
		}
		seriesSeeded++

		for _, mr := range mapsByID[s.SeriesMatchID] {
			var mapWinnerID *uint
			if wid := resolveTeamID(teamLookup, mr.MapWinner, mr.GameCode); wid != 0 {
				mapWinnerID = &wid
			}
			matchMapsBatch = append(matchMapsBatch, models.MatchMap{
				MatchID:     m.ID,
				MapNumber:   mr.MapNumber,
				MapName:     mr.MapName,
				Mode:        mr.Mode,
				Score1:      mr.Score1,
				Score2:      mr.Score2,
				WinnerID:    mapWinnerID,
				Played:      mr.Played,
				DurationSec: parseDurationString(mr.Duration),
				Source:      mr.Source,
			})
		}

		type enrichedAgg struct {
			PlayerID uint
			TeamID   uint
			Kills    int
			Deaths   int
			Maps     int
		}
		enrichedAggs := map[uint]*enrichedAgg{}

		for _, st := range statsByID[s.SeriesMatchID] {
			playerID := resolvePlayer(st.Player, playerLookup, db)
			if playerID == 0 {
				continue
			}
			teamID := resolveTeamID(teamLookup, st.Team, st.GameCode)
			if teamID == 0 {
				teamID = ensureUnknownTeam(db, st.Team, teamLookup)
			}

			playerMapStatsBatch = append(playerMapStatsBatch, models.PlayerMapStats{
				MatchID:         m.ID,
				MapNumber:       st.MapNumber,
				PlayerID:        playerID,
				TeamID:          teamID,
				Kills:           st.Kills,
				Deaths:          st.Deaths,
				KDRatio:         st.KD,
				HillTime:        st.HillTime,
				PlantCount:      st.Plants,
				DefuseCount:     st.Defuses,
				FirstBloodCount: st.FirstKills,
				FirstDeathCount: st.FirstDeaths,
				DataQualityNote: st.DataQualityNote,
				Source:          st.Source,
			})

			if _, ok := enrichedAggs[playerID]; !ok {
				enrichedAggs[playerID] = &enrichedAgg{PlayerID: playerID, TeamID: teamID}
			}
			agg := enrichedAggs[playerID]
			agg.Kills += st.Kills
			agg.Deaths += st.Deaths
			agg.Maps++
		}

		for _, agg := range enrichedAggs {
			kd := 0.0
			if agg.Deaths > 0 {
				kd = float64(agg.Kills) / float64(agg.Deaths)
			}
			playerMatchStatsBatch = append(playerMatchStatsBatch, models.PlayerMatchStats{
				MatchID:     m.ID,
				PlayerID:    agg.PlayerID,
				TeamID:      agg.TeamID,
				MapsPlayed:  agg.Maps,
				TotalKills:  agg.Kills,
				TotalDeaths: agg.Deaths,
				KDRatio:     kd,
			})
		}
	}

	if len(matchMapsBatch) > 0 {
		db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(matchMapsBatch, 500)
	}
	if len(playerMapStatsBatch) > 0 {
		db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(playerMapStatsBatch, 500)
	}
	if len(playerMatchStatsBatch) > 0 {
		db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(playerMatchStatsBatch, 500)
	}

	log.Printf("[enriched] series=%d  maps=%d  playerStats=%d", seriesSeeded, len(matchMapsBatch), len(playerMapStatsBatch))
}
