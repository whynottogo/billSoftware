package model

import "time"

type UserSession struct {
	ID           uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	UserID       uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	SessionToken string    `json:"session_token" xorm:"not null unique VARCHAR(128) 'session_token'"`
	ClientType   string    `json:"client_type" xorm:"not null default 'pc' VARCHAR(32) 'client_type'"`
	IsActive     int       `json:"is_active" xorm:"not null default 1 TINYINT 'is_active'"`
	LoginAt      time.Time `json:"login_at" xorm:"DATETIME 'login_at'"`
	CreatedAt    time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

