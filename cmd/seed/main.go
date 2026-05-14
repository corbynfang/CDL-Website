package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

// seasonConfig links each pair of CSV files to a season.
type seasonConfig struct {
	MatchFile   string // per-match player stats
	PlayerFile  string // season-aggregate player stats (empty string = none)
	Name        string
	GameTitle   string
	StartYear   int
}

var seasons = []seasonConfig{
	{
		MatchFile:  "database/bo6_season_stats.csv",
		PlayerFile: "",
		Name:       "Black Ops 6 2024-25",
		GameTitle:  "Black Ops 6",
		StartYear:  2024,
	},
	{
		MatchFile:  "database/cdl_mw3_season_stats.csv",
		PlayerFile: "database/cdl_mw3_player_stats.csv",
		Name:       "Modern Warfare III 2023-24",
		GameTitle:  "Modern Warfare III",
		StartYear:  2023,
	},
	{
		MatchFile:  "database/cdl_mw2_season_stats.csv",
		PlayerFile: "database/cdl_mw2_players_stats.csv",
		Name:       "Modern Warfare II 2022-23",
		GameTitle:  "Modern Warfare II",
		StartYear:  2022,
	},
	{
		MatchFile:  "database/cdl_vanguard_season_stats.csv",
		PlayerFile: "database/cdl_vanguard_players_stats.csv",
		Name:       "Vanguard 2021-22",
		GameTitle:  "Vanguard",
		StartYear:  2021,
	},
	{
		MatchFile:  "database/cdl_coldwar_season_stats.csv",
		PlayerFile: "database/cdl_coldwar_players_stats.csv",
		Name:       "Black Ops Cold War 2020-21",
		GameTitle:  "Black Ops Cold War",
		StartYear:  2020,
	},
}

// badGamertags are CSV artifacts that are not real CDL players.
var badGamertags = []string{"5aLDx"}

func main() {
	database.ConnectDatabase()
	database.AutoMigrate()

	db := database.DB
	cleanupBadPlayers(db)
	ensureUnaffiliatedTeam(db)

	for _, cfg := range seasons {
		season := createSeason(db, cfg)
		if cfg.MatchFile != "" {
			seedMatchStats(db, cfg, season)
		}
		if cfg.PlayerFile != "" {
			seedPlayerStats(db, cfg, season)
		}
	}

	log.Println("==> Seeding complete.")
}

// cleanupBadPlayers removes known CSV artifacts from the database.
func cleanupBadPlayers(db *gorm.DB) {
	for _, tag := range badGamertags {
		var player database.Player
		if err := db.Where("gamertag = ?", tag).First(&player).Error; err != nil {
			continue // already gone
		}
		// Delete related stats rows first to avoid FK violations
		db.Where("player_id = ?", player.ID).Delete(&database.PlayerMatchStats{})
		db.Where("player_id = ?", player.ID).Delete(&database.PlayerTournamentStats{})
		db.Delete(&player)
		log.Printf("Removed bad player: %s (id=%d)", tag, player.ID)
	}
}

// ── Season ────────────────────────────────────────────────────────────────────

func createSeason(db *gorm.DB, cfg seasonConfig) database.Season {
	season := database.Season{
		Name:      cfg.Name,
		GameTitle: cfg.GameTitle,
		StartDate: time.Date(cfg.StartYear, 9, 1, 0, 0, 0, 0, time.UTC),
	}
	db.Where("name = ?", cfg.Name).FirstOrCreate(&season)
	log.Printf("[%s] season id=%d", cfg.Name, season.ID)
	return season
}

// ── Match-level stats ─────────────────────────────────────────────────────────

type matchRow struct {
	MatchID   string
	Event     string
	Date      string
	TeamA     string
	TeamB     string
	Winner    string
	Score     string
	Team      string
	Player    string
	Kills     int
	Deaths    int
	KD        float64
	PlusMinus int
	Damage    int
	BPRating  float64
}

func seedMatchStats(db *gorm.DB, cfg seasonConfig, season database.Season) {
	rows := readMatchCSV(cfg.MatchFile)
	if len(rows) == 0 {
		return
	}

	// Collect unique names
	teamNames := uniqueStrings(func() []string {
		var s []string
		for _, r := range rows {
			s = append(s, r.TeamA, r.TeamB, r.Team)
		}
		return s
	}())
	playerNames := uniqueStrings(func() []string {
		var s []string
		for _, r := range rows {
			s = append(s, r.Player)
		}
		return s
	}())

	// Derive tournament start dates from earliest match date per event
	// (CSV is reverse-chronological: Champs first → Major 1 last)
	tournamentDates := map[string]time.Time{}
	for _, r := range rows {
		d := parseDate(r.Date)
		if d.IsZero() {
			continue
		}
		if existing, ok := tournamentDates[r.Event]; !ok || d.Before(existing) {
			tournamentDates[r.Event] = d
		}
	}

	// Create teams
	teamMap := map[string]uint{}
	for _, name := range teamNames {
		t := database.Team{Name: name, Abbreviation: makeAbbr(name)}
		db.Where("name = ?", name).FirstOrCreate(&t)
		teamMap[name] = t.ID
	}

	// Create players
	playerMap := map[string]uint{}
	for _, tag := range playerNames {
		p := database.Player{Gamertag: tag}
		db.Where("gamertag = ?", tag).FirstOrCreate(&p)
		playerMap[tag] = p.ID
	}

	// Create tournaments — start date comes from actual match dates, not hardcoded
	tournamentMap := map[string]uint{}
	for _, r := range rows {
		if _, done := tournamentMap[r.Event]; done {
			continue
		}
		startDate := tournamentDates[r.Event]
		if startDate.IsZero() {
			startDate = time.Date(cfg.StartYear, 1, 1, 0, 0, 0, 0, time.UTC)
		}
		t := database.Tournament{
			SeasonID:       season.ID,
			Name:           r.Event,
			TournamentType: detectType(r.Event),
			StartDate:      startDate,
		}
		db.Where("name = ? AND season_id = ?", r.Event, season.ID).FirstOrCreate(&t)
		tournamentMap[r.Event] = t.ID
	}

	// Create matches and per-match player stats
	matchMap := map[string]uint{} // CSV match_id → DB id

	for _, row := range rows {
		tournamentID := tournamentMap[row.Event]
		team1ID := teamMap[row.TeamA]
		team2ID := teamMap[row.TeamB]
		playerID := playerMap[row.Player]
		teamID := teamMap[row.Team]

		if tournamentID == 0 || team1ID == 0 || team2ID == 0 || playerID == 0 {
			continue
		}

		// One match per CSV match_id
		if _, exists := matchMap[row.MatchID]; !exists {
			t1Score, t2Score := parseScore(row.Score)
			var winnerID *uint
			if wid := teamMap[row.Winner]; wid != 0 {
				winnerID = &wid
			}
			m := database.Match{
				TournamentID:  tournamentID,
				Team1ID:       team1ID,
				Team2ID:       team2ID,
				MatchDate:     parseDate(row.Date),
				Team1Score:    t1Score,
				Team2Score:    t2Score,
				WinnerID:      winnerID,
				LiquipediaURL: row.MatchID, // store CSV match_id here for dedup
			}
			db.Where("liquipedia_url = ?", row.MatchID).FirstOrCreate(&m)
			matchMap[row.MatchID] = m.ID
		}

		matchID := matchMap[row.MatchID]
		if matchID == 0 || teamID == 0 {
			continue
		}

		stats := database.PlayerMatchStats{
			MatchID:     matchID,
			PlayerID:    playerID,
			TeamID:      teamID,
			TotalKills:  row.Kills,
			TotalDeaths: row.Deaths,
			TotalDamage: row.Damage,
			KDRatio:     row.KD,
		}
		db.Where("match_id = ? AND player_id = ?", matchID, playerID).FirstOrCreate(&stats)
	}

	log.Printf("[%s] match stats seeded (%d rows)", cfg.Name, len(rows))
}

func readMatchCSV(path string) []matchRow {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("skipping %s: %v", path, err)
		return nil
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		return nil
	}

	var rows []matchRow
	for _, rec := range records[1:] {
		if len(rec) < 13 {
			continue
		}
		row := matchRow{
			MatchID:   strings.TrimSpace(rec[0]),
			Event:     strings.TrimSpace(rec[1]),
			Date:      strings.TrimSpace(rec[2]),
			TeamA:     strings.TrimSpace(rec[3]),
			TeamB:     strings.TrimSpace(rec[4]),
			Winner:    strings.TrimSpace(rec[5]),
			Score:     strings.TrimSpace(rec[6]),
			Team:      strings.TrimSpace(rec[7]),
			Player:    strings.TrimSpace(rec[8]),
			Kills:     atoi(rec[9]),
			Deaths:    atoi(rec[10]),
			KD:        atof(rec[11]),
			PlusMinus: atoi(rec[12]),
		}
		if len(rec) > 13 {
			row.Damage = atoi(rec[13])
		}
		if len(rec) > 14 {
			row.BPRating = atof(rec[14])
		}
		rows = append(rows, row)
	}
	return rows
}

// ── Season-aggregate player stats ─────────────────────────────────────────────

func seedPlayerStats(db *gorm.DB, cfg seasonConfig, season database.Season) {
	f, err := os.Open(cfg.PlayerFile)
	if err != nil {
		log.Printf("skipping %s: %v", cfg.PlayerFile, err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		return
	}

	// Build header index map, normalising to lowercase so MW3's mixed-case headers
	// ("Rank", "K/D") match the same as other seasons' snake_case ones.
	headers := normalizeHeaders(records[0])
	col := func(name string) int {
		if i, ok := headers[name]; ok {
			return i
		}
		return -1
	}

	// A virtual "Season Stats" tournament holds the aggregate rows.
	summaryTournament := database.Tournament{
		SeasonID:       season.ID,
		Name:           cfg.Name + " — Season Stats",
		TournamentType: "championship",
		StartDate:      time.Date(cfg.StartYear, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	db.Where("name = ? AND season_id = ?", summaryTournament.Name, season.ID).
		FirstOrCreate(&summaryTournament)

	count := 0
	for _, rec := range records[1:] {
		get := func(name string) string {
			i := col(name)
			if i < 0 || i >= len(rec) {
				return ""
			}
			return strings.TrimSpace(rec[i])
		}

		gamertag := get("player")
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

		// Look up the player created during match-stats seeding.
		var player database.Player
		if err := db.Where("gamertag = ?", gamertag).First(&player).Error; err != nil {
			// Player not in match data — create a minimal record.
			player = database.Player{Gamertag: gamertag}
			db.Where("gamertag = ?", gamertag).FirstOrCreate(&player)
		}

		// Find the team this player appeared with most in this season.
		// Falls back to any season, then to the "Unaffiliated" placeholder so
		// we never insert team_id = 0 which violates the foreign key constraint.
		teamID := dominantTeam(db, player.ID, season.ID)
		if teamID == 0 {
			teamID = ensureUnaffiliatedTeam(db)
		}

		rank := atoi(get("rank"))
		kills := atoi(get("kills"))
		deaths := atoi(get("deaths"))

		stats := database.PlayerTournamentStats{
			PlayerID:     player.ID,
			TeamID:       teamID,
			TournamentID: summaryTournament.ID,
			Rank:         &rank,
			TotalKills:   kills,
			TotalDeaths:  deaths,
			KDRatio:      atof(get("k/d")),
			OverallMaps:  atoi(get("series_played")),
			// HP
			HpKills:   atoi(get("hp_kills")),
			HpDeaths:  atoi(get("hp_deaths")),
			HpKDRatio: atof(get("hp_k/d")),
			HpKPerMap: atof(get("hp_k/10m")),
			HpMaps:    atoi(get("hp_maps_played")),
			// SND
			SndKills:   atoi(get("snd_kills")),
			SndDeaths:  atoi(get("snd_deaths")),
			SndKDRatio: atof(get("snd_k/d")),
			SndKPerMap: atof(get("snd_kpr")),
			SndMaps:    atoi(get("snd_maps_played")),
			// Control
			ControlKDRatio:  atof(get("ctl_kd")),
			ControlKPerMap:  atof(get("ctl_k/10m")),
			ControlCaptures: atoi(get("ctl_ticks")),
			ControlMaps:     atoi(get("ctl_maps_played")),
		}

		db.Where("player_id = ? AND tournament_id = ?", player.ID, summaryTournament.ID).
			FirstOrCreate(&stats)
		count++
	}

	log.Printf("[%s] player stats seeded (%d rows)", cfg.Name, count)
}

// dominantTeam returns the team ID this player appeared with most.
// Tries the current season first, then falls back to any season.
func dominantTeam(db *gorm.DB, playerID uint, seasonID uint) uint {
	type result struct {
		TeamID uint
		Cnt    int
	}
	var r result

	// Current season
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

	// Any season (player gamertag in stats CSV but not in this season's match CSV)
	db.Raw(`
		SELECT team_id, COUNT(*) AS cnt
		FROM player_match_stats
		WHERE player_id = ?
		GROUP BY team_id ORDER BY cnt DESC LIMIT 1
	`, playerID).Scan(&r)
	return r.TeamID
}

// unaffiliatedTeamID is cached after the first call to ensureUnaffiliatedTeam.
var unaffiliatedTeamID uint

// ensureUnaffiliatedTeam creates (or fetches) a placeholder team used when a
// player appears in the stats CSV but not in any match data, so team_id is unknown.
func ensureUnaffiliatedTeam(db *gorm.DB) uint {
	if unaffiliatedTeamID != 0 {
		return unaffiliatedTeamID
	}
	t := database.Team{Name: "Unaffiliated", Abbreviation: "UNK"}
	db.Where("name = ?", "Unaffiliated").FirstOrCreate(&t)
	unaffiliatedTeamID = t.ID
	return unaffiliatedTeamID
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// normalizeHeaders converts CSV header names to a consistent lowercase key format
// so that "K/D", "kd", and "k/d" all resolve to the same index.
func normalizeHeaders(row []string) map[string]int {
	m := map[string]int{}
	for i, h := range row {
		key := strings.ToLower(strings.TrimSpace(h))
		// Normalise common MW3 title-case variants
		key = strings.ReplaceAll(key, " ", "_")
		key = strings.ReplaceAll(key, "%", "_pct")
		m[key] = i
		// Also store the original normalised key without underscores for fuzzy lookup
		m[strings.ReplaceAll(key, "_", "")] = i
	}
	return m
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(strings.ToLower(s))
	for _, f := range []string{
		"2006-01-02 3:04 pm",
		"2006-01-02 15:04",
		"2006-01-02",
	} {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func parseScore(s string) (int, int) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, 0
	}
	a, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	b, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return a, b
}

func detectType(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "champ"):
		return "championship"
	case strings.Contains(lower, "qualifier"):
		return "qualifier"
	case strings.Contains(lower, "kickoff"):
		return "qualifier"
	case strings.Contains(lower, "major"):
		return "major"
	default:
		return "tournament"
	}
}

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

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

func atoi(s string) int {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	v, _ := strconv.Atoi(s)
	return v
}

func atof(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
