package models

import "time"

type Match struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TournamentID uint      `json:"tournament_id" gorm:"index"`
	Team1ID      uint      `json:"team1_id" gorm:"index"`
	Team2ID      uint      `json:"team2_id" gorm:"index"`
	MatchDate    time.Time `json:"match_date"`
	MatchType    string    `json:"match_type" gorm:"size:50"`
	Format       string    `json:"format" gorm:"size:20"`
	Team1Score   int       `json:"team1_score" gorm:"default:0"`
	Team2Score   int       `json:"team2_score" gorm:"default:0"`
	WinnerID     *uint     `json:"winner_id"`
	DurationMins *int      `json:"duration_minutes"`
	VodURL       string    `json:"vod_url"`
	LiquipediaURL string   `json:"-" gorm:"column:liquipedia_url"`

	BreakingPointMatchID *int `json:"-" gorm:"column:breaking_point_match_id"`

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

type MatchMap struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	MatchID   uint `json:"match_id" gorm:"not null;uniqueIndex:idx_match_map_unique"`
	MapNumber int  `json:"map_number" gorm:"not null;uniqueIndex:idx_match_map_unique"`

	MapName     string `json:"map_name" gorm:"size:100"`
	Mode        string `json:"mode" gorm:"size:50"`
	Score1      int    `json:"score_1" gorm:"default:0"`
	Score2      int    `json:"score_2" gorm:"default:0"`
	WinnerID    *uint  `json:"winner_id"`
	Played      bool   `json:"played" gorm:"default:true"`
	DurationSec int    `json:"duration_sec" gorm:"default:0"`
	Source      string `json:"source" gorm:"size:50"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Match  Match `json:"match" gorm:"foreignKey:MatchID"`
	Winner *Team `json:"winner" gorm:"foreignKey:WinnerID"`
}

func (MatchMap) TableName() string { return "match_maps" }
