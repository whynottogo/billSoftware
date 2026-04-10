package model

import "time"

type Budget struct {
	ID          uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	UserID      uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	BudgetType  string    `json:"budget_type" xorm:"not null VARCHAR(16) 'budget_type'"`
	PeriodKey   string    `json:"period_key" xorm:"not null VARCHAR(16) 'period_key'"`
	TotalAmount float64   `json:"total_amount" xorm:"not null default 0 DECIMAL(18,2) 'total_amount'"`
	UsedAmount  float64   `json:"used_amount" xorm:"not null default 0 DECIMAL(18,2) 'used_amount'"`
	CreatedAt   time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt   time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (Budget) TableName() string {
	return "budgets"
}

type BudgetItem struct {
	ID         uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	BudgetID   uint64    `json:"budget_id" xorm:"not null index BIGINT UNSIGNED 'budget_id'"`
	CategoryID uint64    `json:"category_id" xorm:"not null index BIGINT UNSIGNED 'category_id'"`
	Amount     float64   `json:"amount" xorm:"not null default 0 DECIMAL(18,2) 'amount'"`
	UsedAmount float64   `json:"used_amount" xorm:"not null default 0 DECIMAL(18,2) 'used_amount'"`
	CreatedAt  time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt  time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (BudgetItem) TableName() string {
	return "budget_items"
}
