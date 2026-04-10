package model

import "time"

type Family struct {
	ID            uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	FamilyUID     string    `json:"family_uid" xorm:"not null unique VARCHAR(64) 'family_uid'"`
	FamilyName    string    `json:"family_name" xorm:"not null VARCHAR(64) 'family_name'"`
	CreatorUserID uint64    `json:"creator_user_id" xorm:"not null BIGINT UNSIGNED 'creator_user_id'"`
	CreatedAt     time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt     time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (Family) TableName() string {
	return "families"
}
