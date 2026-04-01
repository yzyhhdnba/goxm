package admin

import (
	"time"

	"gorm.io/gorm"
)

const (
	ReviewStatusPending  = "pending"
	ReviewStatusApproved = "approved"
	ReviewStatusRejected = "rejected"
	ReviewStatusReviewed = "reviewed"
	ReviewStatusAll      = "all"
)

type VideoReview struct {
	ID         uint64    `gorm:"primaryKey"`
	VideoID    uint64    `gorm:"not null;index"`
	ReviewerID uint64    `gorm:"not null;index"`
	Status     string    `gorm:"size:16;not null"`
	Reason     string    `gorm:"size:255"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (VideoReview) TableName() string {
	return "video_reviews"
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type VideoItem struct {
	ID              uint64     `json:"id"`
	AreaID          uint64     `json:"area_id"`
	AreaName        string     `json:"area_name"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	CoverURL        string     `json:"cover_url"`
	PlayURL         string     `json:"play_url"`
	SourcePath      string     `json:"source_path"`
	DurationSeconds uint       `json:"duration_seconds"`
	ReviewStatus    string     `json:"review_status"`
	ReviewReason    string     `json:"review_reason"`
	ViewCount       uint       `json:"view_count"`
	CommentCount    uint       `json:"comment_count"`
	LikeCount       uint       `json:"like_count"`
	FavoriteCount   uint       `json:"favorite_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	PublishedAt     *time.Time `json:"published_at"`
	AuthorID        uint64     `json:"author_id"`
	AuthorUsername  string     `json:"author_username"`
}

type VideoListResponse struct {
	List       []VideoItem `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

type ReviewInput struct {
	Reason string `json:"reason"`
}

type TodayStats struct {
	ActiveUserCount     int64 `json:"active_user_count"`
	SubmittedVideoCount int64 `json:"submitted_video_count"`
	ApprovedVideoCount  int64 `json:"approved_video_count"`
	PlayCount           int64 `json:"play_count"`
	CommentCount        int64 `json:"comment_count"`
}

type AreaStatsItem struct {
	AreaID        uint64 `json:"area_id"`
	AreaName      string `json:"area_name"`
	ApprovedCount int64  `json:"approved_count"`
	PendingCount  int64  `json:"pending_count"`
	RejectedCount int64  `json:"rejected_count"`
	TotalCount    int64  `json:"total_count"`
}

type areaStatsRow struct {
	AreaID        uint64 `gorm:"column:area_id"`
	AreaName      string `gorm:"column:area_name"`
	ApprovedCount int64  `gorm:"column:approved_count"`
	PendingCount  int64  `gorm:"column:pending_count"`
	RejectedCount int64  `gorm:"column:rejected_count"`
	TotalCount    int64  `gorm:"column:total_count"`
}

type reviewState struct {
	ID           uint64         `gorm:"column:id"`
	AuthorID     uint64         `gorm:"column:author_id"`
	Title        string         `gorm:"column:title"`
	ReviewStatus string         `gorm:"column:review_status"`
	PublishedAt  *time.Time     `gorm:"column:published_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}
