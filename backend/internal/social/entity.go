package social

import "time"

type Follow struct {
	ID         uint64    `gorm:"primaryKey"`
	FollowerID uint64    `gorm:"not null;uniqueIndex:idx_follows_follower_followee,priority:1;index"`
	FolloweeID uint64    `gorm:"not null;uniqueIndex:idx_follows_follower_followee,priority:2;index"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (Follow) TableName() string {
	return "follows"
}

type FollowStatusResponse struct {
	Followed bool `json:"followed"`
}

type UserCard struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type UserListResponse struct {
	List       []UserCard `json:"list"`
	Pagination Pagination `json:"pagination"`
}
