package video

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusVisible        = "visible"
	ReviewStatusPending  = "pending"
	ReviewStatusApproved = "approved"
	ReviewStatusRejected = "rejected"
)

type Video struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	AuthorID        uint64         `gorm:"not null;index" json:"author_id"`
	AreaID          uint64         `gorm:"not null;index:idx_videos_area_id_review_status,priority:1;index:idx_videos_area_id_published_at,priority:1" json:"area_id"`
	Title           string         `gorm:"size:128;not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	CoverURL        string         `gorm:"size:255" json:"cover_url"`
	SourcePath      string         `gorm:"size:255" json:"source_path"`
	PlayURL         string         `gorm:"size:255" json:"play_url"`
	DurationSeconds uint           `gorm:"not null;default:0" json:"duration_seconds"`
	Status          string         `gorm:"size:16;not null;default:visible;index" json:"status"`
	ReviewStatus    string         `gorm:"size:16;not null;default:pending;index:idx_videos_area_id_review_status,priority:2" json:"review_status"`
	ReviewReason    string         `gorm:"size:255" json:"review_reason"`
	PublishedAt     *time.Time     `gorm:"index;index:idx_videos_area_id_published_at,priority:2" json:"published_at"`
	LikeCount       uint           `gorm:"not null;default:0" json:"like_count"`
	FavoriteCount   uint           `gorm:"not null;default:0" json:"favorite_count"`
	CommentCount    uint           `gorm:"not null;default:0" json:"comment_count"`
	ViewCount       uint           `gorm:"not null;default:0" json:"view_count"`
	HotScore        int64          `gorm:"not null;default:0;index" json:"hot_score"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoLike struct {
	ID        uint64    `gorm:"primaryKey"`
	VideoID   uint64    `gorm:"not null;uniqueIndex:idx_video_likes_video_user,priority:1;index"`
	UserID    uint64    `gorm:"not null;uniqueIndex:idx_video_likes_video_user,priority:2;index"`
	CreatedAt time.Time `gorm:"not null"`
}

func (VideoLike) TableName() string {
	return "video_likes"
}

type VideoFavorite struct {
	ID        uint64    `gorm:"primaryKey"`
	VideoID   uint64    `gorm:"not null;uniqueIndex:idx_video_favorites_video_user,priority:1;index"`
	UserID    uint64    `gorm:"not null;uniqueIndex:idx_video_favorites_video_user,priority:2;index"`
	CreatedAt time.Time `gorm:"not null"`
}

func (VideoFavorite) TableName() string {
	return "video_favorites"
}

type AuthorPreview struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type ViewerState struct {
	Liked     bool `json:"liked"`
	Favorited bool `json:"favorited"`
	Followed  bool `json:"followed"`
}

type FeedItem struct {
	ID              uint64        `json:"id"`
	AreaID          uint64        `json:"area_id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	CoverURL        string        `json:"cover_url"`
	PlayURL         string        `json:"play_url"`
	DurationSeconds uint          `json:"duration_seconds"`
	ViewCount       uint          `json:"view_count"`
	CommentCount    uint          `json:"comment_count"`
	LikeCount       uint          `json:"like_count"`
	FavoriteCount   uint          `json:"favorite_count"`
	PublishedAt     time.Time     `json:"published_at"`
	Author          AuthorPreview `json:"author"`
}

type FeedResponse struct {
	Items      []FeedItem `json:"items"`
	NextCursor string     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type VideoListResponse struct {
	List       []FeedItem `json:"list"`
	Pagination Pagination `json:"pagination"`
}

type DetailResponse struct {
	ID              uint64        `json:"id"`
	AreaID          uint64        `json:"area_id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	CoverURL        string        `json:"cover_url"`
	PlayURL         string        `json:"play_url"`
	DurationSeconds uint          `json:"duration_seconds"`
	ViewCount       uint          `json:"view_count"`
	CommentCount    uint          `json:"comment_count"`
	LikeCount       uint          `json:"like_count"`
	FavoriteCount   uint          `json:"favorite_count"`
	PublishedAt     time.Time     `json:"published_at"`
	Author          AuthorPreview `json:"author"`
	ViewerState     ViewerState   `json:"viewer_state"`
}

type LikeStatusResponse struct {
	Liked bool `json:"liked"`
}

type FavoriteStatusResponse struct {
	Favorited bool `json:"favorited"`
}

type CreateVideoInput struct {
	AreaID      uint64 `json:"area_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateVideoInput struct {
	AreaID      uint64 `json:"area_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateVideoResponse struct {
	ID           uint64    `json:"id"`
	AreaID       uint64    `json:"area_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ReviewStatus string    `json:"review_status"`
	CreatedAt    time.Time `json:"created_at"`
}

type SourceUploadResponse struct {
	VideoID    uint64 `json:"video_id"`
	SourcePath string `json:"source_path"`
	PlayURL    string `json:"play_url"`
}

type CoverUploadResponse struct {
	VideoID  uint64 `json:"video_id"`
	CoverURL string `json:"cover_url"`
}

type CreatorVideoItem struct {
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
}

type CreatorVideoListResponse struct {
	List       []CreatorVideoItem `json:"list"`
	Pagination Pagination         `json:"pagination"`
}
