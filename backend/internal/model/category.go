package model

import "time"

type CategoryTemplate struct {
	ID           uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	CategoryType string    `json:"category_type" xorm:"not null VARCHAR(16) 'category_type'"`
	Name         string    `json:"name" xorm:"not null VARCHAR(64) 'name'"`
	SortOrder    int       `json:"sort_order" xorm:"not null default 0 INT 'sort_order'"`
	CreatedAt    time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
}

func (CategoryTemplate) TableName() string {
	return "category_templates"
}

type UserCategory struct {
	ID           uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	UserID       uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	CategoryType string    `json:"category_type" xorm:"not null VARCHAR(16) 'category_type'"`
	Name         string    `json:"name" xorm:"not null VARCHAR(64) 'name'"`
	SortOrder    int       `json:"sort_order" xorm:"not null default 0 INT 'sort_order'"`
	IsSystem     int       `json:"is_system" xorm:"not null default 1 TINYINT 'is_system'"`
	CreatedAt    time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (UserCategory) TableName() string {
	return "user_categories"
}

