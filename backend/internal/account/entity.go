package account

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleAdmin    = "admin"
	RoleUser     = "user"
	StatusActive = "active"
)

type User struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	Username         string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Email            string         `gorm:"size:128;not null;uniqueIndex" json:"email"`
	PasswordHash     string         `gorm:"size:255;not null" json:"-"`
	TokenVersion     uint           `gorm:"not null;default:1" json:"-"`
	RefreshTokenHash string         `gorm:"size:128" json:"-"`
	AvatarURL        string         `gorm:"size:255" json:"avatar_url"`
	Bio              string         `gorm:"size:255" json:"bio"`
	Role             string         `gorm:"size:16;not null;default:user;index" json:"role"`
	Status           string         `gorm:"size:16;not null;default:active;index" json:"status"`
	FollowerCount    uint           `gorm:"not null;default:0" json:"follower_count"`
	FollowingCount   uint           `gorm:"not null;default:0" json:"following_count"`
	VideoCount       uint           `gorm:"not null;default:0" json:"video_count"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type UserResponse struct {
	ID             uint64    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	AvatarURL      string    `json:"avatar_url"`
	Bio            string    `json:"bio"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	FollowerCount  uint      `json:"follower_count"`
	FollowingCount uint      `json:"following_count"`
	VideoCount     uint      `json:"video_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ViewerState struct {
	Followed bool `json:"followed"`
}

type DashboardAuthorPreview struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type DashboardVideoItem struct {
	ID              uint64                 `json:"id"`
	AreaID          uint64                 `json:"area_id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	CoverURL        string                 `json:"cover_url"`
	PlayURL         string                 `json:"play_url"`
	DurationSeconds uint                   `json:"duration_seconds"`
	ViewCount       uint                   `json:"view_count"`
	CommentCount    uint                   `json:"comment_count"`
	LikeCount       uint                   `json:"like_count"`
	FavoriteCount   uint                   `json:"favorite_count"`
	PublishedAt     time.Time              `json:"published_at"`
	Author          DashboardAuthorPreview `json:"author"`
}

type DashboardUserCard struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type ProfileResponse struct {
	ID             uint64      `json:"id"`
	Username       string      `json:"username"`
	AvatarURL      string      `json:"avatar_url"`
	Bio            string      `json:"bio"`
	FollowerCount  uint        `json:"follower_count"`
	FollowingCount uint        `json:"following_count"`
	VideoCount     uint        `json:"video_count"`
	ViewerState    ViewerState `json:"viewer_state"`
}

type DashboardStats struct {
	TotalViewCount int64 `json:"total_view_count"`
}

type DashboardResponse struct {
	User           UserResponse         `json:"user"`
	Stats          DashboardStats       `json:"stats"`
	RecentVideos   []DashboardVideoItem `json:"recent_videos"`
	FavoriteVideos []DashboardVideoItem `json:"favorite_videos"`
	FollowingUsers []DashboardUserCard  `json:"following_users"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		AvatarURL:      u.AvatarURL,
		Bio:            u.Bio,
		Role:           u.Role,
		Status:         u.Status,
		FollowerCount:  u.FollowerCount,
		FollowingCount: u.FollowingCount,
		VideoCount:     u.VideoCount,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
