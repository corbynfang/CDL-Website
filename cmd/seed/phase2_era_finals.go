package main

// phase2_era_finals.go — seeds match data from the 15 era_finals CSV files.
// For each of the 5 eras it reads three files: series (one row per match),
// match_maps (one row per map), and player_map_stats (one row per player per map).
// It writes Match, MatchMap, PlayerMapStats, and PlayerMatchStats records.
// PlayerMatchStats is computed by summing the per-map stats — it's the aggregate
// the existing player profile "Matches" tab queries.

import (
	"log"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seedEraFinals iterates all 5 eras and seeds the full match hierarchy.
// Returns bp_match_id → DB match.ID for downstream cross-referencing.
func seedEraFinals(
	db *gorm.DB,
	teamLookup map[string]uint,
	playerLookup map[string]uint,
	seasonByCode map[string]uint,
	tournamentBySlug map[string]uint,
	eventRanges []eventRange,
) map[int]uint {

	matchByBPID := map[int]uint{}

	for _, era := range eraFinalsConfigs {
		log.Printf("[%s] Seeding era_finals…", era.GameCode)

		seriesRows := readSeriesCSV(era.SeriesFile)
		mapRows := readMapCSV(era.MapsFile)
		statRows := readPlayerStatCSV(era.StatsFile)

		// Index maps and stats by BP match_id for O(1) access inside the series loop.
		mapsByMatchID := map[int][]mapRow{}
		for _, mr := range mapRows {
			mapsByMatchID[mr.MatchID] = append(mapsByMatchID[mr.MatchID], mr)
		}
		statsByMatchID := map[int][]playerStatRow{}
		for _, sr := range statRows {
			statsByMatchID[sr.MatchID] = append(statsByMatchID[sr.MatchID], sr)
		}

		// Build source-provider team ID → team name per match.
		matchTeamByID := map[int]matchTeamContext{}
		for _, sr := range seriesRows {
			matchTeamByID[sr.MatchID] = matchTeamContext{
				TeamASourceID: sr.TeamAID, TeamAName: sr.TeamAName,
				TeamBSourceID: sr.TeamBID, TeamBName: sr.TeamBName,
			}
		}

		seasonID := seasonByCode[era.GameCode]
		seriesSeeded := 0

		// Collect child records across all matches in this era, then batch insert once.
		var matchMapsBatch []models.MatchMap
		var playerMapStatsBatch []models.PlayerMapStats
		var playerMatchStatsBatch []models.PlayerMatchStats

		for _, s := range seriesRows {
			matchTime := parseISOTime(s.MatchDatetime)
			team1ID := teamLookup[s.TeamAName]
			team2ID := teamLookup[s.TeamBName]
			if team1ID == 0 || team2ID == 0 {
				log.Printf("[%s] WARN: unresolved team in match %d (%q vs %q)", era.GameCode, s.MatchID, s.TeamAName, s.TeamBName)
			}
			var winnerID *uint
			if wid := teamLookup[s.WinnerName]; wid != 0 {
				winnerID = &wid
			}
			tournamentID := findTournamentForMatch(eventRanges, tournamentBySlug, era.GameCode, matchTime)
			if tournamentID == 0 {
				tournamentID = ensureFallbackTournament(db, seasonID, era.GameCode)
			}

			bpID := s.MatchID
			m := models.Match{
				TournamentID:         tournamentID,
				Team1ID:              team1ID,
				Team2ID:              team2ID,
				MatchDate:            matchTime,
				Format:               s.SeriesFormat,
				Team1Score:           s.TeamAScore,
				Team2Score:           s.TeamBScore,
				WinnerID:             winnerID,
				BreakingPointMatchID: &bpID,
				LiquipediaURL:        s.SourceURL,
				BracketRound:         rawRoundToDBRound(s.RoundName),
			}
			// Match is kept as FirstOrCreate — we need m.ID immediately for child rows.
			db.Where("breaking_point_match_id = ?", bpID).FirstOrCreate(&m)
			// If a previous seeder run placed this match in the fallback tournament
			// (because findTournamentForMatch used strict timestamp comparison and
			// the match ran after the CSV end-time), correct the tournament now.
			if tournamentID != 0 && m.TournamentID != tournamentID {
				db.Model(&m).Update("tournament_id", tournamentID)
				m.TournamentID = tournamentID
			}
			if team1ID != 0 && m.Team1ID != team1ID {
				log.Printf("[%s] correcting team1_id match %d: %d→%d", era.GameCode, bpID, m.Team1ID, team1ID)
				db.Model(&m).Update("team1_id", team1ID)
				m.Team1ID = team1ID
			}
			if team2ID != 0 && m.Team2ID != team2ID {
				log.Printf("[%s] correcting team2_id match %d: %d→%d", era.GameCode, bpID, m.Team2ID, team2ID)
				db.Model(&m).Update("team2_id", team2ID)
				m.Team2ID = team2ID
			}
			if winnerID != nil && m.WinnerID != nil && *m.WinnerID != *winnerID {
				log.Printf("[%s] correcting winner_id match %d: %d→%d", era.GameCode, bpID, *m.WinnerID, *winnerID)
				db.Model(&m).Update("winner_id", winnerID)
				m.WinnerID = winnerID
			}
			matchByBPID[s.MatchID] = m.ID
			seriesSeeded++

			for _, mr := range mapsByMatchID[s.MatchID] {
				var mapWinnerID *uint
				if wid := teamLookup[mr.WinnerName]; wid != 0 {
					mapWinnerID = &wid
				}
				matchMapsBatch = append(matchMapsBatch, models.MatchMap{
					MatchID:     m.ID,
					MapNumber:   mr.MapNumber,
					MapName:     mr.MapName,
					Mode:        mr.ModeName,
					Score1:      mr.ScoreA,
					Score2:      mr.ScoreB,
					WinnerID:    mapWinnerID,
					Played:      mr.Played,
					DurationSec: mr.DurationMin*60 + mr.DurationSec,
					Source:      mr.SourceType,
				})
			}

			matchTeams := matchTeamByID[s.MatchID]
			type matchAgg struct {
				PlayerID uint
				TeamID   uint
				Kills    int
				Deaths   int
				Assists  int
				Damage   int
				Maps     int
			}
			matchAggs := map[uint]*matchAgg{}

			for _, st := range statsByMatchID[s.MatchID] {
				playerID := resolvePlayer(st.PlayerTag, playerLookup, db)
				if playerID == 0 {
					continue
				}
				teamName := ""
				if st.TeamID == matchTeams.TeamASourceID {
					teamName = matchTeams.TeamAName
				} else if st.TeamID == matchTeams.TeamBSourceID {
					teamName = matchTeams.TeamBName
				}
				teamID := teamLookup[teamName]

				playerMapStatsBatch = append(playerMapStatsBatch, models.PlayerMapStats{
					MatchID:              m.ID,
					MapNumber:            st.MapNumber,
					PlayerID:             playerID,
					TeamID:               teamID,
					Kills:                st.Kills,
					Deaths:               st.Deaths,
					KDRatio:              st.KD,
					Damage:               st.Damage,
					Assists:              st.Assists,
					BPRating:             st.BPRating,
					HillTime:             st.HillTime,
					SndRounds:            st.SndRounds,
					PlantCount:           st.PlantCount,
					DefuseCount:          st.DefuseCount,
					SnipeCount:           st.SnipeCount,
					FirstBloodCount:      st.FirstBloodCount,
					FirstDeathCount:      st.FirstDeathCount,
					ZoneTierCaptureCount: st.ZoneTierCaptureCount,
					CtlAttackRounds:      st.CtlAttackRounds,
					CtlDefenseRounds:     st.CtlDefenseRounds,
					NonTradedKills:       st.NonTradedKills,
					HighestStreak:        st.HighestStreak,
					DataQualityNote:      st.DataQualityNote,
					Source:               st.SourceType,
				})

				if _, ok := matchAggs[playerID]; !ok {
					matchAggs[playerID] = &matchAgg{PlayerID: playerID, TeamID: teamID}
				}
				agg := matchAggs[playerID]
				agg.Kills += st.Kills
				agg.Deaths += st.Deaths
				agg.Assists += st.Assists
				agg.Damage += st.Damage
				agg.Maps++
			}

			for _, agg := range matchAggs {
				kd := 0.0
				if agg.Deaths > 0 {
					kd = float64(agg.Kills) / float64(agg.Deaths)
				}
				playerMatchStatsBatch = append(playerMatchStatsBatch, models.PlayerMatchStats{
					MatchID:      m.ID,
					PlayerID:     agg.PlayerID,
					TeamID:       agg.TeamID,
					MapsPlayed:   agg.Maps,
					TotalKills:   agg.Kills,
					TotalDeaths:  agg.Deaths,
					TotalAssists: agg.Assists,
					TotalDamage:  agg.Damage,
					KDRatio:      kd,
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

		log.Printf("[%s] series=%d  maps=%d  playerStats=%d", era.GameCode, seriesSeeded, len(matchMapsBatch), len(playerMapStatsBatch))
	}
	return matchByBPID
}
