package model

import "time"

const (
	FileBizTypeLedgerImage = "ledger_image"
	FileStatusUploaded     = "uploaded"
	FileStatusBound        = "bound"
	FileStatusDeleted      = "deleted"
)

type FileUpload struct {
	ID                 uint64    `json:"id" xorm:"pk autoincr BIGINT UNSIGNED 'id'"`
	UserID             uint64    `json:"user_id" xorm:"not null index BIGINT UNSIGNED 'user_id'"`
	BizType            string    `json:"biz_type" xorm:"not null VARCHAR(32) 'biz_type'"`
	OriginalName       string    `json:"original_name" xorm:"not null VARCHAR(255) 'original_name'"`
	MimeType           string    `json:"mime_type" xorm:"not null VARCHAR(100) 'mime_type'"`
	SizeBytes          uint64    `json:"size_bytes" xorm:"not null BIGINT UNSIGNED 'size_bytes'"`
	OriginalObjectKey  string    `json:"original_object_key" xorm:"not null VARCHAR(500) 'original_object_key'"`
	ThumbnailObjectKey string    `json:"thumbnail_object_key" xorm:"VARCHAR(500) 'thumbnail_object_key'"`
	Width              int       `json:"width" xorm:"INT 'width'"`
	Height             int       `json:"height" xorm:"INT 'height'"`
	Status             string    `json:"status" xorm:"not null VARCHAR(32) 'status'"`
	CreatedAt          time.Time `json:"created_at" xorm:"created DATETIME 'created_at'"`
	UpdatedAt          time.Time `json:"updated_at" xorm:"updated DATETIME 'updated_at'"`
}

func (FileUpload) TableName() string {
	return "file_uploads"
}
