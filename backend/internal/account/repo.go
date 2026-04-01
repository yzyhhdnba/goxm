package account

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	publicVideoStatus   = "visible"
	approvedReviewState = "approved"
)

type Repository struct {
	db *gorm.DB
}

type dashboardVideoRow struct {
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

type dashboardUserRow struct {
	ID        uint64 `gorm:"column:id"`
	Username  string `gorm:"column:username"`
	AvatarURL string `gorm:"column:avatar_url"`
	Bio       string `gorm:"column:bio"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&User{})
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) FindByID(ctx context.Context, id uint64) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindActiveByID(ctx context.Context, id uint64) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("status = ?", StatusActive).
		First(&user).
		Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", strings.ToLower(email)).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) SaveRefreshTokenHash(ctx context.Context, userID uint64, refreshTokenHash string) error {
	return r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"refresh_token_hash": refreshTokenHash,
		}).Error
}

func (r *Repository) InvalidateTokens(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"refresh_token_hash": "",
			"token_version":      gorm.Expr("token_version + 1"),
		}).Error
}

func (r *Repository) RotateRefreshTokenHash(ctx context.Context, userID uint64, tokenVersion uint, previousRefreshTokenHash string, nextRefreshTokenHash string) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ? AND token_version = ? AND refresh_token_hash = ?", userID, tokenVersion, previousRefreshTokenHash).
		Updates(map[string]any{
			"refresh_token_hash": nextRefreshTokenHash,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) FindByUsernameOrEmail(ctx context.Context, identifier string) (*User, error) {
	var user User
	query := r.db.WithContext(ctx).Model(&User{})
	if strings.Contains(identifier, "@") {
		query = query.Where("email = ?", strings.ToLower(identifier))
	} else {
		query = query.Where("username = ?", identifier)
	}
	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) IsFollowing(ctx context.Context, followerID uint64, followeeID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("follows").
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListRecentPublicVideosByAuthor(ctx context.Context, authorID uint64, limit int) ([]dashboardVideoRow, error) {
	var rows []dashboardVideoRow
	if err := r.dashboardVideoQuery(ctx).
		Where("videos.author_id = ?", authorID).
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Limit(limit).
		Scan(&rows).
		Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) ListRecentFavoriteVideosByUser(ctx context.Context, userID uint64, limit int) ([]dashboardVideoRow, error) {
	var rows []dashboardVideoRow
	if err := r.dashboardVideoQuery(ctx).
		Joins("JOIN video_favorites ON video_favorites.video_id = videos.id").
		Where("video_favorites.user_id = ?", userID).
		Order("video_favorites.created_at DESC").
		Order("video_favorites.id DESC").
		Limit(limit).
		Scan(&rows).
		Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) ListFollowingPreview(ctx context.Context, userID uint64, limit int) ([]dashboardUserRow, error) {
	var rows []dashboardUserRow
	if err := r.db.WithContext(ctx).
		Table("follows").
		Joins("JOIN users ON users.id = follows.followee_id").
		Where("follows.follower_id = ?", userID).
		Where("users.deleted_at IS NULL").
		Where("users.status = ?", StatusActive).
		Select("users.id, users.username, users.avatar_url, users.bio").
		Order("follows.created_at DESC").
		Order("follows.id DESC").
		Limit(limit).
		Scan(&rows).
		Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) SumPublicVideoViewsByAuthor(ctx context.Context, authorID uint64) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Table("videos").
		Where("videos.deleted_at IS NULL").
		Where("videos.author_id = ?", authorID).
		Where("videos.status = ?", publicVideoStatus).
		Where("videos.review_status = ?", approvedReviewState).
		Where("videos.published_at IS NOT NULL").
		Select("COALESCE(SUM(videos.view_count), 0)").
		Scan(&total).
		Error; err != nil {
		return 0, err
	}
	return total, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *Repository) dashboardVideoBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Table("videos").
		Joins("JOIN users ON users.id = videos.author_id").
		Where("videos.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Where("videos.status = ?", publicVideoStatus).
		Where("videos.review_status = ?", approvedReviewState).
		Where("videos.published_at IS NOT NULL").
		Where("users.status = ?", StatusActive)
}

func (r *Repository) dashboardVideoQuery(ctx context.Context) *gorm.DB {
	return r.dashboardVideoBase(ctx).Select(`
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
	`)
}
