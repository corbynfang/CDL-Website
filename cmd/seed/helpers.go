package main

// helpers.go — utility functions shared across all seeder phases.
// Anything that isn't tied to a specific phase lives here: date parsers,
// string helpers, DB resolution helpers, and the tournament-date lookup logic.

import (
	"strconv"
	"strings"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

// ─── Date/time parsers ────────────────────────────────────────────────────────

// parseISOTime parses era_finals series timestamps: "2024-12-06T20:00:00+00:00"
func parseISOTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		return t
	}
	return time.Time{}
}

// parseFlexDate parses dates from event_aliases and enriched files.
// Handles: "2021-08-19 2:00 pm", "2024-08-15 14:00", "2022-12-15"
func parseFlexDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, f := range []string{
		"2006-01-02 3:04 pm",
		"2006-01-02 15:04",
		"2006-01-02",
	} {
		if t, err := time.Parse(f, strings.ToLower(s)); err == nil {
			return t
		}
	}
	return time.Time{}
}

// parseTransferDate parses transfer CSV dates: "Oct 8 2025"
func parseTransferDate(s string) time.Time {
	s = strings.TrimSpace(s)
	for _, f := range []string{"Jan 2 2006", "Jan 02 2006"} {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// parseDurationString converts "11:14" (mm:ss) from enriched match maps to total seconds.
func parseDurationString(s string) int {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0
	}
	m, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	sec, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return m*60 + sec
}

// ─── String/number helpers ────────────────────────────────────────────────────

func atoi(s string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(s))
	return v
}

func atof(s string) float64 {
	s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "%"))
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// makeAbbr generates an up-to-10-char abbreviation from capital letters in a team name.
func makeAbbr(name string) string {
	var b strings.Builder
	for _, w := range strings.Fields(name) {
		if len(w) > 0 && w[0] >= 'A' && w[0] <= 'Z' {
			b.WriteByte(w[0])
		}
	}
	r := strings.ToUpper(b.String())
	if r == "" {
		r = strings.ToUpper(name[:min(3, len(name))])
	}
	if len(r) > 10 {
		r = r[:10]
	}
	return r
}

// normalizeHeaders converts CSV header rows to a lowercase snake_case lookup map.
// Handles MW3's mixed-case headers ("Rank", "K/D") alongside other seasons' snake_case.
func normalizeHeaders(row []string) map[string]int {
	m := map[string]int{}
	for i, h := range row {
		key := strings.ToLower(strings.TrimSpace(h))
		key = strings.ReplaceAll(key, " ", "_")
		m[key] = i
		m[strings.ReplaceAll(key, "_", "")] = i
	}
	return m
}

// mergeInto adds all entries from src into dst without overwriting existing keys.
func mergeInto(dst, src map[string]uint) {
	for k, v := range src {
		if _, exists := dst[k]; !exists {
			dst[k] = v
		}
	}
}

// ─── Bracket round mapping ────────────────────────────────────────────────────

// bpRoundToDBRound converts BreakingPoint bp_round_name values to short snake_case identifiers
// stored in matches.bracket_round. Unknown values are snake_cased automatically.
func bpRoundToDBRound(raw string) string {
	switch raw {
	case "Major Qualifier":
		return "major_qualifier"
	case "Winners Round 1":
		return "winners_r1"
	case "Winners Round 2":
		return "winners_r2"
	case "Winners Round 3":
		return "winners_r3"
	case "Winners Finals":
		return "winners_finals"
	case "Elimination Round 1":
		return "elim_r1"
	case "Elimination Round 2":
		return "elim_r2"
	case "Elimination Round 3":
		return "elim_r3"
	case "Elimination Round 4":
		return "elim_r4"
	case "Elimination Finals":
		return "elim_finals"
	case "Finals":
		return "finals"
	case "Grand Finals":
		return "grand_finals"
	case "3rd Place Decider":
		return "3rd_place"
	default:
		return strings.ToLower(strings.ReplaceAll(raw, " ", "_"))
	}
}

// ─── DB resolution helpers ────────────────────────────────────────────────────

// resolvePlayer returns the DB player ID for a gamertag.
// If the tag isn't in the lookup yet, it creates a minimal Player record and caches it.
func resolvePlayer(tag string, lookup map[string]uint, db *gorm.DB) uint {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return 0
	}
	if id, ok := lookup[tag]; ok {
		return id
	}
	lower := strings.ToLower(tag)
	for name, id := range lookup {
		if strings.ToLower(name) == lower {
			lookup[tag] = id
			return id
		}
	}
	p := database.Player{Gamertag: tag}
	db.Where("gamertag = ?", tag).FirstOrCreate(&p)
	lookup[tag] = p.ID
	return p.ID
}

// ensureUnknownTeam creates a minimal Team record for a name not in the lookup.
// Used when enriched data references a team that isn't in any alias CSV.
func ensureUnknownTeam(db *gorm.DB, name string, teamLookup map[string]uint) uint {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0
	}
	if id, ok := teamLookup[name]; ok {
		return id
	}
	t := database.Team{
		Name:               name,
		Abbreviation:       makeAbbr(name),
		IsCDLFranchise:     false,
		TeamClassification: "unknown",
		DoNotMerge:         true,
		NeedsManualReview:  true,
		Source:             "enriched_csv",
	}
	db.Where("name = ? AND source = ?", name, "enriched_csv").FirstOrCreate(&t)
	teamLookup[name] = t.ID
	return t.ID
}

// findTournamentForMatch returns the tournament ID whose date range contains matchTime.
// Iterates the eventRanges slice built from event_aliases_clean.csv.
func findTournamentForMatch(ranges []eventRange, bySlug map[string]uint, gameCode string, matchTime time.Time) uint {
	if matchTime.IsZero() {
		return 0
	}
	for _, r := range ranges {
		if r.GameCode != gameCode {
			continue
		}
		if !matchTime.Before(r.StartDate) && !matchTime.After(r.EndDate) {
			if id, ok := bySlug[r.Slug]; ok {
				return id
			}
		}
	}
	return 0
}

// fallbackTournamentIDs caches catch-all tournament IDs by season so we don't create duplicates.
var fallbackTournamentIDs = map[uint]uint{}

// ensureFallbackTournament creates a catch-all tournament for matches whose date
// doesn't fall within any known event range. Every era gets at most one.
func ensureFallbackTournament(db *gorm.DB, seasonID uint, gameCode string) uint {
	if id, ok := fallbackTournamentIDs[seasonID]; ok {
		return id
	}
	slug := gameCode + "-unmatched"
	t := database.Tournament{
		SeasonID:       seasonID,
		Name:           gameCode + " Unmatched Matches",
		Slug:           slug,
		TournamentType: "unknown",
	}
	db.Where("slug = ? AND season_id = ?", slug, seasonID).FirstOrCreate(&t)
	fallbackTournamentIDs[seasonID] = t.ID
	return t.ID
}

var unaffiliatedTeamID uint

// ensureUnaffiliatedTeam returns (or creates) a placeholder team used when a player appears
// in a stats CSV but their team can't be determined from match data.
func ensureUnaffiliatedTeam(db *gorm.DB, teamLookup map[string]uint) uint {
	if unaffiliatedTeamID != 0 {
		return unaffiliatedTeamID
	}
	if id, ok := teamLookup["Unaffiliated"]; ok {
		unaffiliatedTeamID = id
		return id
	}
	t := database.Team{Name: "Unaffiliated", Abbreviation: "UNK", Source: "system"}
	db.Where("name = ?", "Unaffiliated").FirstOrCreate(&t)
	unaffiliatedTeamID = t.ID
	teamLookup["Unaffiliated"] = t.ID
	return t.ID
}

// dominantTeam returns the team ID a player appeared with most in a given season.
// Falls back to any season if the player had no matches in the target season.
func dominantTeam(db *gorm.DB, playerID uint, seasonID uint) uint {
	type result struct {
		TeamID uint
		Cnt    int
	}
	var r result
	db.Raw(`
		SELECT pms.team_id, COUNT(*) AS cnt
		FROM player_match_stats pms
		JOIN matches m ON m.id = pms.match_id
		JOIN tournaments t ON t.id = m.tournament_id
		WHERE pms.player_id = ? AND t.season_id = ?
		GROUP BY pms.team_id ORDER BY cnt DESC LIMIT 1
	`, playerID, seasonID).Scan(&r)
	if r.TeamID != 0 {
		return r.TeamID
	}
	db.Raw(`
		SELECT team_id, COUNT(*) AS cnt FROM player_match_stats
		WHERE player_id = ?
		GROUP BY team_id ORDER BY cnt DESC LIMIT 1
	`, playerID).Scan(&r)
	return r.TeamID
}
