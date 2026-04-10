package model

import "time"

type FamilyMember struct {
	ID       uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	FamilyID uint64    `json:"family_id" xorm:"not null BIGINT UNSIGNED 'family_id'"`
	UserID   uint64    `json:"user_id" xorm:"not null BIGINT UNSIGNED 'user_id'"`
	JoinedAt time.Time `json:"joined_at" xorm:"created DATETIME 'joined_at'"`
}

func (FamilyMember) TableName() string {
	return "family_members"
}
