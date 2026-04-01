package search

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	publicVideoStatus   = "visible"
	approvedReviewState = "approved"
	activeUserStatus    = "active"
)

type Repository struct {
	db *gorm.DB
}

type videoRow struct {
	ID              uint64    `gorm:"column:id"`
	AreaID          uint64    `gorm:"column:area_id"`
	Title           string    `gorm:"column:title"`
	Description     string    `gorm:"column:description"`
	CoverURL        string    `gorm:"column:cover_url"`
	PlayURL         string    `gorm:"column:play_url"`
	DurationSeconds uint      `gorm:"column:duration_seconds"`
	ViewCount       uint      `gorm:"column:view_count"`
	CommentCount    uint      `gorm:"column:comment_count"`
	LikeCount       uint      `gorm:"column:like_count"`
	FavoriteCount   uint      `gorm:"column:favorite_count"`
	PublishedAt     time.Time `gorm:"column:published_at"`
	AuthorID        uint64    `gorm:"column:author_id"`
	AuthorUsername  string    `gorm:"column:author_username"`
	AuthorAvatarURL string    `gorm:"column:author_avatar_url"`
}

type userRow struct {
	ID            uint64 `gorm:"column:id"`
	Username      string `gorm:"column:username"`
	AvatarURL     string `gorm:"column:avatar_url"`
	Bio           string `gorm:"column:bio"`
	FollowerCount uint   `gorm:"column:follower_count"`
	VideoCount    uint   `gorm:"column:video_count"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SearchVideos(ctx context.Context, keyword string, page int, pageSize int) ([]videoRow, int64, error) {
	var total int64
	countQuery := r.publicVideoBase(ctx, keyword)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count search videos: %w", err)
	}

	var rows []videoRow
	if err := r.publicVideoBase(ctx, keyword).
		Select(`
			videos.id,
			videos.area_id,
			videos.title,
			videos.description,
			videos.cover_url,
			videos.play_url,
			videos.duration_seconds,
			videos.view_count,
			videos.comment_count,
			videos.like_count,
			videos.favorite_count,
			videos.published_at,
			users.id AS author_id,
			users.username AS author_username,
			users.avatar_url AS author_avatar_url
		`).
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, fmt.Errorf("list search videos: %w", err)
	}

	return rows, total, nil
}

func (r *Repository) SearchUsers(ctx context.Context, keyword string, page int, pageSize int) ([]userRow, int64, error) {
	base := r.db.WithContext(ctx).
		Table("users").
		Where("users.deleted_at IS NULL").
		Where("users.status = ?", activeUserStatus).
		Where("users.username LIKE ? OR users.bio LIKE ?", likeKeyword(keyword), likeKeyword(keyword))

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count search users: %w", err)
	}

	var rows []userRow
	if err := base.
		Select(`
			users.id,
			users.username,
			users.avatar_url,
			users.bio,
			users.follower_count,
			users.video_count
		`).
		Order("users.follower_count DESC").
		Order("users.video_count DESC").
		Order("users.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, fmt.Errorf("list search users: %w", err)
	}

	return rows, total, nil
}

func (r *Repository) publicVideoBase(ctx context.Context, keyword string) *gorm.DB {
	return r.db.WithContext(ctx).
		Table("videos").
		Joins("JOIN users ON users.id = videos.author_id").
		Where("videos.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Where("videos.status = ?", publicVideoStatus).
		Where("videos.review_status = ?", approvedReviewState).
		Where("videos.published_at IS NOT NULL").
		Where("users.status = ?", activeUserStatus).
		Where("videos.title LIKE ? OR videos.description LIKE ?", likeKeyword(keyword), likeKeyword(keyword))
}

func likeKeyword(keyword string) string {
	return "%" + keyword + "%"
}
