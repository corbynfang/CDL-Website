package database

import "time"

// Franchise is a CDL franchise slot stable across team rebrands.
// Minnesota RØKKR and G2 Minnesota are separate Team rows both pointing to the same Franchise.
type Franchise struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FranchiseKey string    `json:"franchise_key" gorm:"uniqueIndex;not null;size:100"`
	Name         string    `json:"name" gorm:"not null;size:200"` // most recent/active team name
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Franchise) TableName() string { return "franchises" }

// Team represents one historical team name/branding at a specific point in time.
// CDL franchises have FranchiseID set; non-CDL teams (challengers, parent orgs, academy) have FranchiseID nil.
type Team struct {
	ID             uint       `json:"id" gorm:"primaryKey;column:id"`
	Name           string     `json:"name" gorm:"not null;size:200;column:name"`
	Abbreviation   string     `json:"abbreviation" gorm:"not null;size:10;column:abbreviation"`
	City           string     `json:"city" gorm:"size:100;column:city"`
	LogoURL        string     `json:"logo_url" gorm:"column:logo_url"`
	PrimaryColor   string     `json:"primary_color" gorm:"size:7;column:primary_color"`
	SecondaryColor string     `json:"secondary_color" gorm:"size:7;column:secondary_color"`
	FoundedDate    *time.Time `json:"founded_date" gorm:"column:founded_date"`
	IsActive       bool       `json:"is_active" gorm:"default:true;column:is_active"`

	// Franchise continuity — nil for non-CDL teams
	FranchiseID *uint `json:"franchise_id" gorm:"column:franchise_id;index"`

	// Which CDL game era this branding belongs to (BO6, CW, MW2, MW3, VG).
	// Populated from cdl_team_branding_by_season.csv. Empty for non-CDL teams.
	GameCode string `json:"game_code" gorm:"size:10;column:game_code"`

	// Classification
	IsCDLFranchise     bool   `json:"is_cdl_franchise" gorm:"default:false;column:is_cdl_franchise;index"`
	TeamClassification string `json:"team_classification" gorm:"size:60;column:team_classification"`
	// cdl_franchise | challenger | academy | non_cdl_org | orgless | unknown | unknown_challenger_or_regional
	DoNotMerge bool `json:"do_not_merge" gorm:"default:false;column:do_not_merge"`

	// Date range this name/branding was valid
	ValidFrom *time.Time `json:"valid_from" gorm:"column:valid_from"`
	ValidTo   *time.Time `json:"valid_to" gorm:"column:valid_to"`

	// Data quality
	NeedsManualReview bool   `json:"needs_manual_review" gorm:"default:false;column:needs_manual_review"`
	Source            string `json:"source" gorm:"size:100;column:source"` // branding_csv | non_cdl_alias | transfer_csv

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	Franchise *Franchise `json:"franchise" gorm:"foreignKey:FranchiseID"`
	Players   []Player   `json:"players" gorm:"many2many:team_rosters;"`
}

func (Team) TableName() string { return "teams" }

type Player struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Gamertag      string     `json:"gamertag" gorm:"not null;size:100;index"`
	FirstName     string     `json:"first_name" gorm:"size:100"`
	LastName      string     `json:"last_name" gorm:"size:100"`
	Country       string     `json:"country" gorm:"size:3"`
	Birthdate     *time.Time `json:"birthdate"`
	Role          string     `json:"role" gorm:"size:50"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	LiquipediaURL string     `json:"liquipedia_url"`
	TwitterHandle string     `json:"twitter_handle" gorm:"size:100"`
	AvatarURL     string     `json:"avatar_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Player) TableName() string { return "players" }

type TeamRoster struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	TeamID    uint       `json:"team_id" gorm:"index"`
	PlayerID  uint       `json:"player_id" gorm:"index"`
	SeasonID  uint       `json:"season_id" gorm:"index"`
	Role      string     `json:"role" gorm:"size:50"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsStarter bool       `json:"is_starter" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Season Season `json:"season" gorm:"foreignKey:SeasonID"`
}

type Season struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null;size:100"`
	GameTitle string     `json:"game_title" gorm:"not null;size:100"`
	GameCode  string     `json:"game_code" gorm:"size:10"` // BO6, CW, MW2, MW3, VG
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsActive  bool       `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Season) TableName() string { return "seasons" }

type Tournament struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	SeasonID         uint       `json:"season_id" gorm:"index"`
	Name             string     `json:"name" gorm:"not null;size:200"`
	Slug             string     `json:"slug" gorm:"size:200"` // event_slug from event_aliases_clean.csv
	TournamentType   string     `json:"tournament_type" gorm:"size:50"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	PrizePool        *float64   `json:"prize_pool" gorm:"type:decimal(12,2)"`
	Location         string     `json:"location" gorm:"size:100"`
	TournamentFormat string     `json:"tournament_format" gorm:"size:50"`
	LiquipediaURL    string     `json:"liquipedia_url"`
	BreakingPointURL string     `json:"breaking_point_url" gorm:"column:breaking_point_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Season Season `json:"season" gorm:"foreignKey:SeasonID"`
}

func (Tournament) TableName() string { return "tournaments" }

type Match struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TournamentID uint      `json:"tournament_id" gorm:"index"`
	Team1ID      uint      `json:"team1_id" gorm:"index"`
	Team2ID      uint      `json:"team2_id" gorm:"index"`
	MatchDate    time.Time `json:"match_date"`
	MatchType    string    `json:"match_type" gorm:"size:50"`
	Format       string    `json:"format" gorm:"size:20"` // BO5, BO7, BO9
	Team1Score   int       `json:"team1_score" gorm:"default:0"`
	Team2Score   int       `json:"team2_score" gorm:"default:0"`
	WinnerID     *uint     `json:"winner_id"`
	DurationMins *int      `json:"duration_minutes"`
	VodURL       string    `json:"vod_url"`
	LiquipediaURL string   `json:"liquipedia_url"` // general dedup key / external URL

	// Set for all era_finals sourced matches (real BP match IDs like 93815)
	BreakingPointMatchID *int `json:"breaking_point_match_id" gorm:"column:breaking_point_match_id"`

	// Bracket context
	BracketRound    string `json:"bracket_round" gorm:"size:50"`
	BracketPosition int    `json:"bracket_position" gorm:"default:0"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Team1      Team       `json:"team1" gorm:"foreignKey:Team1ID"`
	Team2      Team       `json:"team2" gorm:"foreignKey:Team2ID"`
	Winner     *Team      `json:"winner" gorm:"foreignKey:WinnerID"`
}

func (Match) TableName() string { return "matches" }

// MatchMap is a single map played within a series match.
// One Match has many MatchMaps (up to 5/7/9 depending on format).
type MatchMap struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	MatchID   uint `json:"match_id" gorm:"not null;uniqueIndex:idx_match_map_unique"`
	MapNumber int  `json:"map_number" gorm:"not null;uniqueIndex:idx_match_map_unique"`

	MapName  string `json:"map_name" gorm:"size:100"`
	Mode     string `json:"mode" gorm:"size:50"` // Hardpoint, Search and Destroy, Control
	Score1   int    `json:"score_1" gorm:"default:0"` // Team1 (team_a) score
	Score2   int    `json:"score_2" gorm:"default:0"` // Team2 (team_b) score
	WinnerID *uint  `json:"winner_id"`
	Played   bool   `json:"played" gorm:"default:true"`

	// duration_min*60 + duration_sec from BP source
	DurationSec int    `json:"duration_sec" gorm:"default:0"`
	Source      string `json:"source" gorm:"size:50"` // breakingpoint | ewc_2024 | ewc_2025 | major1_2023_wiki

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Match  Match `json:"match" gorm:"foreignKey:MatchID"`
	Winner *Team `json:"winner" gorm:"foreignKey:WinnerID"`
}

func (MatchMap) TableName() string { return "match_maps" }

// PlayerMapStats holds per-player stats for a single map within a series.
// This is the most granular stat in the database — every kill, death, damage, hill second per map.
type PlayerMapStats struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	MatchID   uint `json:"match_id" gorm:"not null;uniqueIndex:idx_player_map_stat_unique"`
	MapNumber int  `json:"map_number" gorm:"not null;uniqueIndex:idx_player_map_stat_unique"`
	PlayerID  uint `json:"player_id" gorm:"not null;uniqueIndex:idx_player_map_stat_unique;index"`
	TeamID    uint `json:"team_id" gorm:"not null;index"`

	// Core stats present across all eras
	Kills   int     `json:"kills" gorm:"default:0"`
	Deaths  int     `json:"deaths" gorm:"default:0"`
	KDRatio float64 `json:"kd_ratio" gorm:"type:decimal(6,3);default:0"`
	Damage  int     `json:"damage" gorm:"default:0"`
	Assists int     `json:"assists" gorm:"default:0"`

	// BreakingPoint composite rating
	BPRating float64 `json:"bp_rating" gorm:"type:decimal(10,6);default:0"`

	// Mode-specific — 0 when not applicable to the map's mode
	HillTime             int `json:"hill_time" gorm:"default:0"`               // HP: seconds on hill
	SndRounds            int `json:"snd_rounds" gorm:"default:0"`              // SND: rounds played
	PlantCount           int `json:"plant_count" gorm:"default:0"`             // SND
	DefuseCount          int `json:"defuse_count" gorm:"default:0"`            // SND
	SnipeCount           int `json:"snipe_count" gorm:"default:0"`             // SND
	FirstBloodCount      int `json:"first_blood_count" gorm:"default:0"`       // SND
	FirstDeathCount      int `json:"first_death_count" gorm:"default:0"`       // SND
	ZoneTierCaptureCount int `json:"zone_tier_capture_count" gorm:"default:0"` // Control
	CtlAttackRounds      int `json:"ctl_attack_rounds" gorm:"default:0"`       // Control
	CtlDefenseRounds     int `json:"ctl_defense_rounds" gorm:"default:0"`      // Control

	NonTradedKills  int    `json:"non_traded_kills" gorm:"default:0"`
	HighestStreak   int    `json:"highest_streak" gorm:"default:0"`
	DataQualityNote string `json:"data_quality_note" gorm:"size:200"`
	Source          string `json:"source" gorm:"size:50"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Match  Match  `json:"match" gorm:"foreignKey:MatchID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
}

func (PlayerMapStats) TableName() string { return "player_map_stats" }

type PlayerMatchStats struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	MatchID      uint      `json:"match_id" gorm:"uniqueIndex:idx_player_match_stat_unique"`
	PlayerID     uint      `json:"player_id" gorm:"uniqueIndex:idx_player_match_stat_unique;index"`
	TeamID       uint      `json:"team_id" gorm:"index"`
	MapsPlayed   int       `json:"maps_played" gorm:"default:0"`
	TotalKills   int       `json:"total_kills" gorm:"default:0"`
	TotalDeaths  int       `json:"total_deaths" gorm:"default:0"`
	TotalAssists int       `json:"total_assists" gorm:"default:0"`
	TotalDamage  int       `json:"total_damage" gorm:"default:0"`
	KDRatio      float64   `json:"kd_ratio" gorm:"type:decimal(4,2);default:0"`
	KDARatio     float64   `json:"kda_ratio" gorm:"type:decimal(4,2);default:0"`
	ADR          float64   `json:"adr" gorm:"type:decimal(6,2);default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Match  Match  `json:"match" gorm:"foreignKey:MatchID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
}

func (PlayerMatchStats) TableName() string { return "player_match_stats" }

type PlayerTournamentStats struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	PlayerID     uint    `json:"player_id" gorm:"uniqueIndex:idx_player_tournament_stat_unique"`
	TeamID       uint    `json:"team_id" gorm:"index"`
	TournamentID uint    `json:"tournament_id" gorm:"uniqueIndex:idx_player_tournament_stat_unique;index"`
	TotalKills   int     `json:"total_kills"`
	TotalDeaths  int     `json:"total_deaths"`
	TotalAssists int     `json:"total_assists"`
	TotalDamage  int     `json:"total_damage"`
	KDRatio      float64 `json:"kd_ratio"`
	KDARatio     float64 `json:"kda_ratio"`

	Rank               *int    `json:"rank"`
	OverallPlusMinus   int     `json:"overall_plus_minus" gorm:"default:0"`
	OverallMaps        int     `json:"overall_maps" gorm:"default:0"`
	SndKills           int     `json:"snd_kills" gorm:"default:0"`
	SndDeaths          int     `json:"snd_deaths" gorm:"default:0"`
	SndKDRatio         float64 `json:"snd_kd_ratio" gorm:"default:0"`
	SndPlusMinus       int     `json:"snd_plus_minus" gorm:"default:0"`
	SndKPerMap         float64 `json:"snd_k_per_map" gorm:"default:0"`
	SndFirstKills      int     `json:"snd_first_kills" gorm:"default:0"`
	SndMaps            int     `json:"snd_maps" gorm:"default:0"`
	HpKills            int     `json:"hp_kills" gorm:"default:0"`
	HpDeaths           int     `json:"hp_deaths" gorm:"default:0"`
	HpKDRatio          float64 `json:"hp_kd_ratio" gorm:"default:0"`
	HpPlusMinus        int     `json:"hp_plus_minus" gorm:"default:0"`
	HpKPerMap          float64 `json:"hp_k_per_map" gorm:"default:0"`
	HpTimeMilliseconds int     `json:"hp_time_milliseconds" gorm:"default:0"`
	HpMaps             int     `json:"hp_maps" gorm:"default:0"`
	ControlKills       int     `json:"control_kills" gorm:"default:0"`
	ControlDeaths      int     `json:"control_deaths" gorm:"default:0"`
	ControlKDRatio     float64 `json:"control_kd_ratio" gorm:"default:0"`
	ControlPlusMinus   int     `json:"control_plus_minus" gorm:"default:0"`
	ControlKPerMap     float64 `json:"control_k_per_map" gorm:"default:0"`
	ControlCaptures    int     `json:"control_captures" gorm:"default:0"`
	ControlMaps        int     `json:"control_maps" gorm:"default:0"`

	Player     Player     `json:"player" gorm:"foreignKey:PlayerID"`
	Team       Team       `json:"team" gorm:"foreignKey:TeamID"`
	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
}

func (PlayerTournamentStats) TableName() string { return "player_tournament_stats" }

type TeamTournamentStats struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TournamentID  uint      `json:"tournament_id" gorm:"index"`
	TeamID        uint      `json:"team_id" gorm:"index"`
	Placement     *int      `json:"placement"`
	MatchesPlayed int       `json:"matches_played" gorm:"default:0"`
	MatchesWon    int       `json:"matches_won" gorm:"default:0"`
	MatchesLost   int       `json:"matches_lost" gorm:"default:0"`
	MapsPlayed    int       `json:"maps_played" gorm:"default:0"`
	MapsWon       int       `json:"maps_won" gorm:"default:0"`
	MapsLost      int       `json:"maps_lost" gorm:"default:0"`
	PrizeMoney    float64   `json:"prize_money" gorm:"type:decimal(10,2);default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Team       Team       `json:"team" gorm:"foreignKey:TeamID"`
}

func (TeamTournamentStats) TableName() string { return "team_tournament_stats" }

type Coach struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;size:100"`
	TeamID   uint   `json:"team_id"`
	SeasonID uint   `json:"season_id"`
}

func (Coach) TableName() string { return "coaches" }

type PlayerTransfer struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PlayerID     uint      `json:"player_id" gorm:"index"`
	FromTeamID   *uint     `json:"from_team_id" gorm:"index"`
	ToTeamID     *uint     `json:"to_team_id" gorm:"index"`
	TransferDate time.Time `json:"transfer_date"`
	TransferType string    `json:"transfer_type" gorm:"size:50"` // Signing | Transfer | Release | Retirement | Role Change
	Role         string    `json:"role" gorm:"size:50"`
	GameCode     string    `json:"game_code" gorm:"size:10"` // BO6, CW, MW2, MW3, VG
	Season       string    `json:"season" gorm:"size:50"`
	Description  string    `json:"description" gorm:"size:500"`

	// Raw source text always preserved — "Free Agent" stored as-is, from_team_id set nil.
	RawFromTeamName string `json:"raw_from_team_name" gorm:"size:200;column:raw_from_team_name"`
	RawToTeamName   string `json:"raw_to_team_name" gorm:"size:200;column:raw_to_team_name"`

	CreatedAt time.Time `json:"created_at"`

	Player   Player `json:"player" gorm:"foreignKey:PlayerID"`
	FromTeam *Team  `json:"from_team" gorm:"foreignKey:FromTeamID"`
	ToTeam   *Team  `json:"to_team" gorm:"foreignKey:ToTeamID"`
}

func (PlayerTransfer) TableName() string { return "player_transfers" }
