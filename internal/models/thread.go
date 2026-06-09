package models

import "time"

type MatchThread struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MatchID   uint      `json:"match_id" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Match     Match     `json:"match,omitempty" gorm:"foreignKey:MatchID"`
}

func (MatchThread) TableName() string { return "match_threads" }

type ThreadPost struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	ThreadID  uint       `json:"thread_id" gorm:"index;not null"`
	UserID    uint       `json:"user_id" gorm:"index;not null"`
	Body      string     `json:"body" gorm:"type:text;not null"`
	Edited    bool       `json:"edited" gorm:"default:false"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
}

func (ThreadPost) TableName() string { return "thread_posts" }
