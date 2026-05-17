package main

import "time"

// types.go — all struct types used across the seeder.
// Keeping types in one place means any phase file can use them without hunting through 2,000 lines.

// ─── CSV row types ────────────────────────────────────────────────────────────
// Each struct maps 1:1 to a CSV file we read.

type brandingRow struct {
	GameCode          string
	SeasonYear        string
	RawTeamName       string
	CanonicalTeamName string
	TeamSlug          string
	FranchiseKey      string
	ValidFrom         string
	ValidTo           string
	Notes             string
}

type nonCDLRow struct {
	RawTeamName        string
	CanonicalTeamName  string
	TeamSlug           string
	IsCDLFranchise     bool
	TeamClassification string
	Region             string
	LinkedCDLTeam      string
	RelationshipType   string
	DoNotMerge         bool
	NeedsManualReview  bool
	Notes              string
}

type playerAliasRow struct {
	PlayerName          string
	CanonicalPlayerName string
	PlayerSlug          string
	TeamContext         string
	GameCode            string
	SeasonYear          string
	NeedsManualReview   bool
	Notes               string
}

type eventAliasRow struct {
	RawEventName         string
	CanonicalEventName   string
	EventSlug            string
	GameCode             string
	SeasonYear           string
	EventType            string
	StartDate            string
	EndDate              string
	BreakingPointEventID string
	SourceURL            string
	FandomURL            string
	StatsOnly            bool
	HasBracket           bool
	HasQualifierStage    bool
	HasStats             bool
}

// seriesRow is one match from an era_finals *_series_final.csv.
type seriesRow struct {
	MatchID       int
	SourceURL     string
	MatchDatetime string
	BestOf        int
	Status        string
	TeamAID       int // BP internal team ID
	TeamAName     string
	TeamBID       int // BP internal team ID
	TeamBName     string
	TeamAScore    int
	TeamBScore    int
	WinnerID      int // BP internal
	WinnerName    string
	BPRoundName   string
	SeriesFormat  string
	SourceType    string
}

// mapRow is one map entry from an era_finals *_match_maps_final.csv.
type mapRow struct {
	MatchID     int
	MapNumber   int
	MapName     string
	ModeName    string
	TeamAID     int
	TeamBID     int
	ScoreA      int
	ScoreB      int
	WinnerID    int
	WinnerName  string
	Played      bool
	DurationMin int
	DurationSec int
	SourceType  string
}

// playerStatRow is one player-per-map entry from an era_finals *_player_map_stats_final.csv.
type playerStatRow struct {
	MatchID              int
	MapNumber            int
	PlayerID             int // BP internal
	PlayerTag            string
	TeamID               int // BP internal
	Kills                int
	Deaths               int
	KD                   float64
	Damage               int
	Assists              int
	BPRating             float64
	HillTime             int
	SndRounds            int
	PlantCount           int
	DefuseCount          int
	SnipeCount           int
	FirstBloodCount      int
	FirstDeathCount      int
	ZoneTierCaptureCount int
	CtlAttackRounds      int
	CtlDefenseRounds     int
	NonTradedKills       int
	HighestStreak        int
	DataQualityNote      string
	SourceType           string
}

type transferRow struct {
	Date         string
	Player       string
	FromTeam     string
	ToTeam       string
	Role         string
	TransferType string
}

// ─── Enriched CSV row types ───────────────────────────────────────────────────
// The enriched_*.csv files (EWC events, Major 1 2023 wiki) have a slightly different
// schema than era_finals, so they get their own types.

type enrichedSeriesRow struct {
	SeriesMatchID     string
	EventName         string
	EventSlug         string
	GameCode          string
	SeasonYear        string
	StageName         string
	StageType         string
	GroupName         string
	RoundName         string
	MatchDatetime     string
	Team1             string
	Team2             string
	Team1Canonical    string
	Team2Canonical    string
	Team1FranchiseKey string
	Team2FranchiseKey string
	Team1MapWins      int
	Team2MapWins      int
	MapsPlayed        int
	Winner            string
	WinnerCanonical   string
	SeriesFormat      string
	Source            string
	SourceURL         string
}

type enrichedMapRow struct {
	SeriesMatchID string
	EventSlug     string
	GameCode      string
	SeasonYear    string
	MapNumber     int
	Team1         string
	Team2         string
	MapName       string
	Mode          string
	Score1        int
	Score2        int
	MapWinner     string
	Played        bool
	Duration      string
	Source        string
}

type enrichedStatRow struct {
	SeriesMatchID   string
	MapNumber       int
	EventSlug       string
	GameCode        string
	SeasonYear      string
	MapName         string
	Mode            string
	Team            string
	Player          string
	Kills           int
	Deaths          int
	KD              float64
	HillTime        int
	Captures        int
	Plants          int
	Defuses         int
	FirstKills      int
	FirstDeaths     int
	DataQualityNote string
	Source          string
}

// ─── Internal helper types ────────────────────────────────────────────────────

// eventRange is used by findTournamentForMatch to locate a tournament by date overlap.
type eventRange struct {
	Slug      string
	GameCode  string
	StartDate time.Time
	EndDate   time.Time
}

// bpMatchTeams holds BP team IDs → team names for resolving which team a player belongs to.
// The era_finals player stats file only has BP internal team IDs, not names.
type bpMatchTeams struct {
	TeamABPID int
	TeamAName string
	TeamBBPID int
	TeamBName string
}

// unresolvedTeamEntry records the outcome of every transfer team name resolution.
// Written to database/unresolved_transfer_teams.csv after Phase 5.
type unresolvedTeamEntry struct {
	RawName            string
	Count              int
	FirstSeen          time.Time
	LastSeen           time.Time
	SuggestedCanonical string
	ResolutionStatus   string // resolved_existing | auto_created
	NeedsManualReview  bool
}
