package history

import (
	"time"

	"gorm.io/gorm"
)

type ViewHistory struct {
	ID              uint64    `gorm:"primaryKey"`
	UserID          uint64    `gorm:"not null;uniqueIndex:idx_view_histories_user_video,priority:1;index"`
	VideoID         uint64    `gorm:"not null;uniqueIndex:idx_view_histories_user_video,priority:2;index"`
	ProgressSeconds uint      `gorm:"not null;default:0"`
	WatchedAt       time.Time `gorm:"not null;index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (ViewHistory) TableName() string {
	return "view_histories"
}

type ReportInput struct {
	VideoID         uint64 `json:"video_id"`
	ProgressSeconds uint   `json:"progress_seconds"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type Item struct {
	VideoID         uint64    `json:"video_id"`
	VideoTitle      string    `json:"video_title"`
	CoverURL        string    `json:"cover_url"`
	PlayURL         string    `json:"play_url"`
	AuthorID        uint64    `json:"author_id"`
	AuthorName      string    `json:"author_name"`
	AreaName        string    `json:"area_name"`
	WatchedAt       time.Time `json:"watched_at"`
	ProgressSeconds uint      `json:"progress_seconds"`
	DurationSeconds uint      `json:"duration_seconds"`
}

type ListResponse struct {
	List       []Item     `json:"list"`
	Pagination Pagination `json:"pagination"`
}
