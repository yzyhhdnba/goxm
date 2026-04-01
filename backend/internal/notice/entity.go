package notice

import (
	"time"

	"gorm.io/gorm"
)

const TypeVideoReview = "video_review"

type Notice struct {
	ID             uint64         `gorm:"primaryKey"`
	UserID         uint64         `gorm:"not null;index"`
	Type           string         `gorm:"size:32;not null;default:video_review;index"`
	Title          string         `gorm:"size:128;not null"`
	Content        string         `gorm:"type:text;not null"`
	RelatedVideoID *uint64        `gorm:"index"`
	IsRead         bool           `gorm:"not null;default:false;index"`
	ReadAt         *time.Time     `gorm:"index"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notice) TableName() string {
	return "notices"
}

type Item struct {
	ID             uint64     `json:"id"`
	Type           string     `json:"type"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	RelatedVideoID *uint64    `json:"related_video_id"`
	Read           bool       `json:"read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type ListResponse struct {
	List       []Item     `json:"list"`
	Pagination Pagination `json:"pagination"`
}
