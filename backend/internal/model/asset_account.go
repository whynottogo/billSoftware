package model

import "time"

type AssetAccount struct {
	ID            uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	UserID        uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	AccountType   string    `json:"account_type" xorm:"not null VARCHAR(32) 'account_type'"`
	SubType       string    `json:"sub_type" xorm:"VARCHAR(32) 'sub_type'"`
	Name          string    `json:"name" xorm:"not null VARCHAR(128) 'name'"`
	Remark        string    `json:"remark" xorm:"VARCHAR(255) 'remark'"`
	CardNo        string    `json:"card_no" xorm:"VARCHAR(128) 'card_no'"`
	BalanceNature string    `json:"balance_nature" xorm:"not null default 'asset' VARCHAR(32) 'balance_nature'"`
	Balance       float64   `json:"balance" xorm:"not null default 0 DECIMAL(18,2) 'balance'"`
	CreatedAt     time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt     time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (AssetAccount) TableName() string {
	return "asset_accounts"
}

type AssetAccountLog struct {
	ID           uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	AccountID    uint64    `json:"account_id" xorm:"not null index BIGINT UNSIGNED 'account_id'"`
	UserID       uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	ChangeType   string    `json:"change_type" xorm:"not null VARCHAR(32) 'change_type'"`
	AmountBefore float64   `json:"amount_before" xorm:"not null default 0 DECIMAL(18,2) 'amount_before'"`
	AmountChange float64   `json:"amount_change" xorm:"not null default 0 DECIMAL(18,2) 'amount_change'"`
	AmountAfter  float64   `json:"amount_after" xorm:"not null default 0 DECIMAL(18,2) 'amount_after'"`
	Remark       string    `json:"remark" xorm:"VARCHAR(255) 'remark'"`
	LogDate      time.Time `json:"log_date" xorm:"DATETIME 'log_date'"`
	CreatedAt    time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
}

func (AssetAccountLog) TableName() string {
	return "asset_account_logs"
}
