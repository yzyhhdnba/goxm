package search

import "time"

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type AuthorPreview struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type VideoItem struct {
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

type VideoSearchResponse struct {
	List       []VideoItem `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

type UserItem struct {
	ID            uint64 `json:"id"`
	Username      string `json:"username"`
	AvatarURL     string `json:"avatar_url"`
	Bio           string `json:"bio"`
	FollowerCount uint   `json:"follower_count"`
	VideoCount    uint   `json:"video_count"`
}

type UserSearchResponse struct {
	List       []UserItem `json:"list"`
	Pagination Pagination `json:"pagination"`
}
