package database

import (
	"time"
)

type Season struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null;size:100"`
	GameTitle string     `json:"game_title" gorm:"not null;size:100"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsActive  bool       `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Season) TableName() string {
	return "seasons"
}

type Team struct {
	ID             uint       `json:"id" gorm:"primaryKey;column:id"`
	Name           string     `json:"name" gorm:"not null;size:100;column:name"`
	Abbreviation   string     `json:"abbreviation" gorm:"not null;size:10;column:abbreviation"`
	City           string     `json:"city" gorm:"size:100;column:city"`
	LogoURL        string     `json:"logo_url" gorm:"column:logo_url"`
	PrimaryColor   string     `json:"primary_color" gorm:"size:7;column:primary_color"`
	SecondaryColor string     `json:"secondary_color" gorm:"size:7;column:secondary_color"`
	FoundedDate    *time.Time `json:"founded_date" gorm:"column:founded_date"`
	IsActive       bool       `json:"is_active" gorm:"default:true;column:is_active"`
	CreatedAt      time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Team) TableName() string {
	return "teams"
}

type Player struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Gamertag      string     `json:"gamertag" gorm:"not null;size:100"`
	FirstName     string     `json:"first_name" gorm:"size:100"`
	LastName      string     `json:"last_name" gorm:"size:100"`
	Country       string     `json:"country" gorm:"size:3"`
	Birthdate     *time.Time `json:"birthdate"`
	Role          string     `json:"role" gorm:"size:50"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	LiquipediaURL string     `json:"liquipedia_url"`
	TwitterHandle string     `json:"twitter_handle" gorm:"size:100"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Player) TableName() string {
	return "players"
}

type TeamRoster struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	TeamID    uint       `json:"team_id"`
	PlayerID  uint       `json:"player_id"`
	SeasonID  uint       `json:"season_id"`
	Role      string     `json:"role" gorm:"size:50"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsStarter bool       `json:"is_starter" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Season Season `json:"season" gorm:"foreignKey:SeasonID"`
}

type Tournament struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	SeasonID         uint       `json:"season_id"`
	Name             string     `json:"name" gorm:"not null;size:200"`
	TournamentType   string     `json:"tournament_type" gorm:"size:50"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	PrizePool        *float64   `json:"prize_pool" gorm:"type:decimal(12,2)"`
	Location         string     `json:"location" gorm:"size:100"`
	TournamentFormat string     `json:"tournament_format" gorm:"size:50"`
	LiquipediaURL    string     `json:"liquipedia_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relationships
	Season Season `json:"season" gorm:"foreignKey:SeasonID"`
}

func (Tournament) TableName() string {
	return "tournaments"
}

type Match struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TournamentID  uint      `json:"tournament_id"`
	Team1ID       uint      `json:"team1_id"`
	Team2ID       uint      `json:"team2_id"`
	MatchDate     time.Time `json:"match_date"`
	MatchType     string    `json:"match_type" gorm:"size:50"`
	Format        string    `json:"format" gorm:"size:20"`
	Team1Score    int       `json:"team1_score" gorm:"default:0"`
	Team2Score    int       `json:"team2_score" gorm:"default:0"`
	WinnerID      *uint     `json:"winner_id"`
	DurationMins  *int      `json:"duration_minutes"`
	VodURL        string    `json:"vod_url"`
	LiquipediaURL string    `json:"liquipedia_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Team1      Team       `json:"team1" gorm:"foreignKey:Team1ID"`
	Team2      Team       `json:"team2" gorm:"foreignKey:Team2ID"`
	Winner     *Team      `json:"winner" gorm:"foreignKey:WinnerID"`
}

func (Match) TableName() string {
	return "matches"
}

type PlayerMatchStats struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	MatchID      uint      `json:"match_id"`
	PlayerID     uint      `json:"player_id"`
	TeamID       uint      `json:"team_id"`
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

func (PlayerMatchStats) TableName() string {
	return "player_match_stats"
}

type TeamTournamentStats struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TournamentID  uint      `json:"tournament_id"`
	TeamID        uint      `json:"team_id"`
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

	// Relationships
	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Team       Team       `json:"team" gorm:"foreignKey:TeamID"`
}

func (TeamTournamentStats) TableName() string {
	return "team_tournament_stats"
}

// Coaches table
// Represents a coach for a team in a given season
// Matches the schema in schema.sql
// Table name: coaches

type Coach struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;size:100"`
	TeamID   uint   `json:"team_id"`
	SeasonID uint   `json:"season_id"`
}

func (Coach) TableName() string {
	return "coaches"
}
