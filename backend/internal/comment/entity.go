package comment

import (
	"time"

	"gorm.io/gorm"
)

const StatusVisible = "visible"

type Comment struct {
	ID         uint64    `gorm:"primaryKey"`
	VideoID    uint64    `gorm:"not null;index:idx_comments_video_root_created,priority:1;index"`
	UserID     uint64    `gorm:"not null;index"`
	RootID     uint64    `gorm:"not null;default:0;index:idx_comments_video_root_created,priority:2"`
	ParentID   uint64    `gorm:"not null;default:0;index"`
	Content    string    `gorm:"type:text;not null"`
	ReplyCount uint      `gorm:"not null;default:0"`
	LikeCount  uint      `gorm:"not null;default:0"`
	Status     string    `gorm:"size:16;not null;default:visible;index"`
	CreatedAt  time.Time `gorm:"index:idx_comments_video_root_created,priority:3"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentLike struct {
	ID        uint64    `gorm:"primaryKey"`
	CommentID uint64    `gorm:"not null;uniqueIndex:idx_comment_likes_comment_user,priority:1;index"`
	UserID    uint64    `gorm:"not null;uniqueIndex:idx_comment_likes_comment_user,priority:2;index"`
	CreatedAt time.Time `gorm:"not null"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

type UserPreview struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type ViewerState struct {
	Liked bool `json:"liked"`
}

type Item struct {
	ID          uint64      `json:"id"`
	RootID      uint64      `json:"root_id"`
	ParentID    uint64      `json:"parent_id"`
	Content     string      `json:"content"`
	LikeCount   uint        `json:"like_count"`
	ReplyCount  uint        `json:"reply_count"`
	CreatedAt   time.Time   `json:"created_at"`
	User        UserPreview `json:"user"`
	ViewerState ViewerState `json:"viewer_state"`
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

type CreateInput struct {
	Content string `json:"content"`
}

type LikeStatusResponse struct {
	Liked bool `json:"liked"`
}
