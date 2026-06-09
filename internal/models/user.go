package models

import "time"

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SupabaseUID string     `json:"supabase_uid" gorm:"uniqueIndex;not null;size:36"`
	Username    string     `json:"username" gorm:"uniqueIndex;not null;size:30"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`
}

func (User) TableName() string { return "users" }
