package models

import "time"

type Player struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Gamertag      string     `json:"gamertag" gorm:"not null;size:100;index"`
	FirstName     string     `json:"first_name" gorm:"size:100"`
	LastName      string     `json:"last_name" gorm:"size:100"`
	Country       string     `json:"country" gorm:"size:3"`
	Birthdate     *time.Time `json:"birthdate"`
	Role          string     `json:"role" gorm:"size:50"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	LiquipediaURL string     `json:"source_profile_url,omitempty" gorm:"column:liquipedia_url"`
	TwitterHandle string     `json:"twitter_handle" gorm:"size:100"`
	AvatarURL     string     `json:"avatar_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Player) TableName() string { return "players" }

type PlayerTransfer struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PlayerID     uint      `json:"player_id" gorm:"index"`
	FromTeamID   *uint     `json:"from_team_id" gorm:"index"`
	ToTeamID     *uint     `json:"to_team_id" gorm:"index"`
	TransferDate time.Time `json:"transfer_date"`
	TransferType string    `json:"transfer_type" gorm:"size:50"`
	Role         string    `json:"role" gorm:"size:50"`
	GameCode     string    `json:"game_code" gorm:"size:10"`
	Season       string    `json:"season" gorm:"size:50"`
	Description  string    `json:"description" gorm:"size:500"`

	RawFromTeamName string `json:"raw_from_team_name" gorm:"size:200;column:raw_from_team_name"`
	RawToTeamName   string `json:"raw_to_team_name" gorm:"size:200;column:raw_to_team_name"`

	CreatedAt time.Time `json:"created_at"`

	Player   Player `json:"player" gorm:"foreignKey:PlayerID"`
	FromTeam *Team  `json:"from_team" gorm:"foreignKey:FromTeamID"`
	ToTeam   *Team  `json:"to_team" gorm:"foreignKey:ToTeamID"`
}

func (PlayerTransfer) TableName() string { return "player_transfers" }
