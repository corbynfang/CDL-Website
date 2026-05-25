package main

import "time"

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
	SourceEventID string
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
	TeamAID       int // source-provider team ID
	TeamAName     string
	TeamBID       int // source-provider team ID
	TeamBName     string
	TeamAScore    int
	TeamBScore    int
	WinnerID      int // source-provider internal
	WinnerName    string
	RoundName     string
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

// matchTeamContext maps source-provider team IDs → team names for resolving which team a player belongs to.
// The era_finals player stats file only has source-provider internal team IDs, not names.
type matchTeamContext struct {
	TeamASourceID int
	TeamAName     string
	TeamBSourceID int
	TeamBName     string
}

// cwBracketRow is one match from database/cdl_cw_stage_brackets.csv.
// canonical_round_key is already the target bracket_round value; no further
// normalisation is needed. bracket_position is the 1-based slot within the round.
type cwBracketRow struct {
	TournamentSlug  string // stage slug, e.g. "cdl-2021-stage-1-major"
	SourceRoundName string // human-readable name, kept for audit trail
	CanonicalRound  string // target bracket_round value, e.g. "elim_r4"
	Position        int    // 1-based position within the round
	Team1Name       string
	Team2Name       string
	Team1Score      int
	Team2Score      int
	WinnerName      string
	MatchDate       string
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
