package model

import "time"

type LedgerRecord struct {
	ID         uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	UserID     uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	RecordType string    `json:"record_type" xorm:"not null VARCHAR(16) 'record_type'"`
	CategoryID uint64    `json:"category_id" xorm:"BIGINT UNSIGNED 'category_id'"`
	Amount     float64   `json:"amount" xorm:"not null DECIMAL(18,2) 'amount'"`
	Remark     string    `json:"remark" xorm:"VARCHAR(255) 'remark'"`
	RecordDate time.Time `json:"record_date" xorm:"not null DATE 'record_date'"`
	ImageURL   string    `json:"image_url" xorm:"VARCHAR(500) 'image_url'"`
	CreatedAt  time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt  time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (LedgerRecord) TableName() string {
	return "ledger_records"
}
