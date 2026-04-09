package model

import "time"

type User struct {
	ID               uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	Username         string    `json:"username" xorm:"not null unique VARCHAR(64) 'username'"`
	Nickname         string    `json:"nickname" xorm:"not null VARCHAR(64) 'nickname'"`
	Phone            string    `json:"phone" xorm:"not null unique VARCHAR(20) 'phone'"`
	Email            string    `json:"email" xorm:"VARCHAR(128) 'email'"`
	PasswordHash     string    `json:"-" xorm:"not null VARCHAR(255) 'password_hash'"`
	Status           int       `json:"status" xorm:"not null default 0 TINYINT 'status'"`
	AvatarOriginal   string    `json:"avatar_original,omitempty" xorm:"LONGTEXT 'avatar_original'"`
	AvatarCompressed string    `json:"avatar_compressed,omitempty" xorm:"LONGTEXT 'avatar_compressed'"`
	CreatedAt        time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt        time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (User) TableName() string {
	return "users"
}

