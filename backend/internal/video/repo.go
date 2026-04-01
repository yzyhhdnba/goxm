package video

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

var errAreaNotFound = errors.New("area not found")

type FeedCursor struct {
	PublishedAt time.Time
	ID          uint64
	HasID       bool
}

type HotFeedCursor struct {
	HotScore    int64
	PublishedAt time.Time
	ID          uint64
}

const activeAreaStatus = "active"

type videoRow struct {
	ID              uint64    `gorm:"column:id"`
	AreaID          uint64    `gorm:"column:area_id"`
	Title           string    `gorm:"column:title"`
	Description     string    `gorm:"column:description"`
	CoverURL        string    `gorm:"column:cover_url"`
	PlayURL         string    `gorm:"column:play_url"`
	DurationSeconds uint      `gorm:"column:duration_seconds"`
	HotScore        int64     `gorm:"column:hot_score"`
	ViewCount       uint      `gorm:"column:view_count"`
	CommentCount    uint      `gorm:"column:comment_count"`
	LikeCount       uint      `gorm:"column:like_count"`
	FavoriteCount   uint      `gorm:"column:favorite_count"`
	PublishedAt     time.Time `gorm:"column:published_at"`
	AuthorID        uint64    `gorm:"column:author_id"`
	AuthorUsername  string    `gorm:"column:author_username"`
	AuthorAvatarURL string    `gorm:"column:author_avatar_url"`
}

type creatorRow struct {
	ID              uint64     `gorm:"column:id"`
	AreaID          uint64     `gorm:"column:area_id"`
	AreaName        string     `gorm:"column:area_name"`
	Title           string     `gorm:"column:title"`
	Description     string     `gorm:"column:description"`
	CoverURL        string     `gorm:"column:cover_url"`
	PlayURL         string     `gorm:"column:play_url"`
	SourcePath      string     `gorm:"column:source_path"`
	DurationSeconds uint       `gorm:"column:duration_seconds"`
	ReviewStatus    string     `gorm:"column:review_status"`
	ReviewReason    string     `gorm:"column:review_reason"`
	ViewCount       uint       `gorm:"column:view_count"`
	CommentCount    uint       `gorm:"column:comment_count"`
	LikeCount       uint       `gorm:"column:like_count"`
	FavoriteCount   uint       `gorm:"column:favorite_count"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	PublishedAt     *time.Time `gorm:"column:published_at"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Video{}, &VideoLike{}, &VideoFavorite{})
}

func (r *Repository) Create(ctx context.Context, item *Video) error {
	if err := ensureActiveArea(r.db.WithContext(ctx), item.AreaID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) ListRecommend(ctx context.Context, cursor *FeedCursor, limit int) ([]videoRow, error) {
	var rows []videoRow

	query := r.publicVideoQuery(ctx).
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Limit(limit)

	if cursor != nil {
		if cursor.HasID {
			query = query.Where(
				"(videos.published_at < ?) OR (videos.published_at = ? AND videos.id < ?)",
				cursor.PublishedAt,
				cursor.PublishedAt,
				cursor.ID,
			)
		} else {
			query = query.Where("videos.published_at < ?", cursor.PublishedAt)
		}
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) ListHot(ctx context.Context, cursor *HotFeedCursor, limit int) ([]videoRow, error) {
	var rows []videoRow

	query := r.publicVideoQuery(ctx).
		Order("videos.hot_score DESC").
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Limit(limit)

	if cursor != nil {
		query = query.Where(
			`(videos.hot_score < ?)
			 OR (videos.hot_score = ? AND videos.published_at < ?)
			 OR (videos.hot_score = ? AND videos.published_at = ? AND videos.id < ?)`,
			cursor.HotScore,
			cursor.HotScore,
			cursor.PublishedAt,
			cursor.HotScore,
			cursor.PublishedAt,
			cursor.ID,
		)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) ListFollowing(ctx context.Context, userID uint64, cursor *FeedCursor, limit int) ([]videoRow, error) {
	var rows []videoRow

	query := r.publicVideoQuery(ctx).
		Joins("JOIN follows ON follows.followee_id = videos.author_id").
		Where("follows.follower_id = ?", userID).
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Limit(limit)

	if cursor != nil {
		if cursor.HasID {
			query = query.Where(
				"(videos.published_at < ?) OR (videos.published_at = ? AND videos.id < ?)",
				cursor.PublishedAt,
				cursor.PublishedAt,
				cursor.ID,
			)
		} else {
			query = query.Where("videos.published_at < ?", cursor.PublishedAt)
		}
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) ListByArea(ctx context.Context, areaID uint64, cursor *FeedCursor, limit int) ([]videoRow, error) {
	if err := ensureActiveArea(r.db.WithContext(ctx), areaID); err != nil {
		return nil, err
	}

	var rows []videoRow
	query := r.publicVideoQuery(ctx).
		Where("videos.area_id = ?", areaID).
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Limit(limit)

	if cursor != nil {
		if cursor.HasID {
			query = query.Where(
				"(videos.published_at < ?) OR (videos.published_at = ? AND videos.id < ?)",
				cursor.PublishedAt,
				cursor.PublishedAt,
				cursor.ID,
			)
		} else {
			query = query.Where("videos.published_at < ?", cursor.PublishedAt)
		}
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) ListByAuthor(ctx context.Context, authorID uint64, page int, pageSize int) ([]videoRow, int64, error) {
	if err := ensureActiveVideoAuthor(r.db.WithContext(ctx), authorID); err != nil {
		return nil, 0, err
	}

	base := r.publicVideoBase(ctx).Where("videos.author_id = ?", authorID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []videoRow
	if err := r.publicVideoQuery(ctx).
		Where("videos.author_id = ?", authorID).
		Order("videos.published_at DESC").
		Order("videos.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) FindPublicByID(ctx context.Context, videoID uint64) (*videoRow, error) {
	var row videoRow

	if err := r.publicVideoQuery(ctx).
		Where("videos.id = ?", videoID).
		Take(&row).
		Error; err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *Repository) FindOwnedByID(ctx context.Context, videoID uint64, authorID uint64) (*Video, error) {
	var item Video
	if err := r.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", videoID, authorID).
		Take(&item).
		Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) UpdateSourceByOwner(ctx context.Context, videoID uint64, authorID uint64, sourcePath string, playURL string) error {
	result := r.db.WithContext(ctx).
		Model(&Video{}).
		Where("id = ? AND author_id = ?", videoID, authorID).
		Updates(map[string]any{
			"source_path": sourcePath,
			"play_url":    playURL,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) UpdateCoverByOwner(ctx context.Context, videoID uint64, authorID uint64, coverURL string) error {
	result := r.db.WithContext(ctx).
		Model(&Video{}).
		Where("id = ? AND author_id = ?", videoID, authorID).
		Update("cover_url", coverURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) UpdateMetadataByOwner(ctx context.Context, videoID uint64, authorID uint64, input UpdateVideoInput) (*creatorRow, error) {
	now := time.Now().UTC()

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureActiveArea(tx, input.AreaID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errAreaNotFound
			}
			return err
		}

		var current Video
		if err := tx.
			Select("id", "author_id", "review_status", "published_at").
			Where("id = ? AND author_id = ?", videoID, authorID).
			Take(&current).
			Error; err != nil {
			return err
		}

		wasApproved := current.ReviewStatus == ReviewStatusApproved && current.PublishedAt != nil

		if err := tx.Model(&Video{}).
			Where("id = ? AND author_id = ?", videoID, authorID).
			Updates(map[string]any{
				"area_id":       input.AreaID,
				"title":         input.Title,
				"description":   input.Description,
				"review_status": ReviewStatusPending,
				"review_reason": "",
				"published_at":  nil,
				"updated_at":    now,
			}).
			Error; err != nil {
			return err
		}

		if wasApproved {
			if err := tx.Table("users").
				Where("id = ?", authorID).
				Update("video_count", gorm.Expr("CASE WHEN video_count > 0 THEN video_count - 1 ELSE 0 END")).
				Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var row creatorRow
	if err := r.creatorBase(ctx).
		Select(`
			videos.id,
			videos.area_id,
			areas.name AS area_name,
			videos.title,
			videos.description,
			videos.cover_url,
			videos.play_url,
			videos.source_path,
			videos.duration_seconds,
			videos.review_status,
			videos.review_reason,
			videos.view_count,
			videos.comment_count,
			videos.like_count,
			videos.favorite_count,
			videos.created_at,
			videos.updated_at,
			videos.published_at
		`).
		Where("videos.id = ? AND videos.author_id = ?", videoID, authorID).
		Take(&row).
		Error; err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *Repository) ListByCreator(ctx context.Context, authorID uint64, reviewStatus string, page int, pageSize int) ([]creatorRow, int64, error) {
	if err := ensureActiveVideoAuthor(r.db.WithContext(ctx), authorID); err != nil {
		return nil, 0, err
	}

	base := r.creatorBase(ctx).Where("videos.author_id = ?", authorID)
	base = applyCreatorReviewStatus(base, reviewStatus)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []creatorRow
	if err := base.
		Select(`
			videos.id,
			videos.area_id,
			areas.name AS area_name,
			videos.title,
			videos.description,
			videos.cover_url,
			videos.play_url,
			videos.source_path,
			videos.duration_seconds,
			videos.review_status,
			videos.review_reason,
			videos.view_count,
			videos.comment_count,
			videos.like_count,
			videos.favorite_count,
			videos.created_at,
			videos.updated_at,
			videos.published_at
		`).
		Order("videos.created_at DESC").
		Order("videos.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) EnsurePublic(ctx context.Context, videoID uint64) error {
	return ensurePublicVideoTx(r.db.WithContext(ctx), videoID)
}

func (r *Repository) Like(ctx context.Context, videoID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePublicVideoTx(tx, videoID); err != nil {
			return err
		}

		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&VideoLike{
			VideoID: videoID,
			UserID:  userID,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return incrementVideoCounter(tx, videoID, "like_count", 1)
	})
}

func (r *Repository) Unlike(ctx context.Context, videoID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePublicVideoTx(tx, videoID); err != nil {
			return err
		}

		result := tx.Where("video_id = ? AND user_id = ?", videoID, userID).Delete(&VideoLike{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return decrementVideoCounter(tx, videoID, "like_count")
	})
}

func (r *Repository) HasLike(ctx context.Context, videoID uint64, userID uint64) (bool, error) {
	return exists(ctx, r.db.Model(&VideoLike{}).Where("video_id = ? AND user_id = ?", videoID, userID))
}

func (r *Repository) Favorite(ctx context.Context, videoID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePublicVideoTx(tx, videoID); err != nil {
			return err
		}

		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&VideoFavorite{
			VideoID: videoID,
			UserID:  userID,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return incrementVideoCounter(tx, videoID, "favorite_count", 1)
	})
}

func (r *Repository) Unfavorite(ctx context.Context, videoID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePublicVideoTx(tx, videoID); err != nil {
			return err
		}

		result := tx.Where("video_id = ? AND user_id = ?", videoID, userID).Delete(&VideoFavorite{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return decrementVideoCounter(tx, videoID, "favorite_count")
	})
}

func (r *Repository) HasFavorite(ctx context.Context, videoID uint64, userID uint64) (bool, error) {
	return exists(ctx, r.db.Model(&VideoFavorite{}).Where("video_id = ? AND user_id = ?", videoID, userID))
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// publicVideoBase 是公开视频查询的统一收口点。
// 对应文档“视频主链路：推荐流、热门流、详情页”，所有公开流量接口都依赖这里的筛选条件。
func (r *Repository) publicVideoBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Model(&Video{}).
		Joins("JOIN users ON users.id = videos.author_id").
		Where("videos.status = ?", StatusVisible).
		Where("videos.review_status = ?", ReviewStatusApproved).
		Where("videos.published_at IS NOT NULL").
		Where("users.status = ?", activeVideoAuthorStatus)
}

func (r *Repository) creatorBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Model(&Video{}).
		Joins("JOIN areas ON areas.id = videos.area_id").
		Where("areas.deleted_at IS NULL")
}

func (r *Repository) publicVideoQuery(ctx context.Context) *gorm.DB {
	return r.publicVideoBase(ctx).
		Select(`
			videos.id,
			videos.area_id,
			videos.title,
			videos.description,
			videos.cover_url,
			videos.play_url,
			videos.duration_seconds,
			videos.hot_score,
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

func ensurePublicVideoTx(tx *gorm.DB, videoID uint64) error {
	var id uint64
	return tx.Model(&Video{}).
		Select("id").
		Where("id = ?", videoID).
		Where("status = ?", StatusVisible).
		Where("review_status = ?", ReviewStatusApproved).
		Where("published_at IS NOT NULL").
		Take(&id).
		Error
}

const activeVideoAuthorStatus = "active"

func ensureActiveVideoAuthor(query *gorm.DB, authorID uint64) error {
	var id uint64
	return query.Table("users").
		Select("id").
		Where("id = ?", authorID).
		Where("status = ?", activeVideoAuthorStatus).
		Take(&id).
		Error
}

func ensureActiveArea(query *gorm.DB, areaID uint64) error {
	var id uint64
	return query.Table("areas").
		Select("id").
		Where("id = ?", areaID).
		Where("status = ?", activeAreaStatus).
		Take(&id).
		Error
}

func applyCreatorReviewStatus(query *gorm.DB, reviewStatus string) *gorm.DB {
	switch reviewStatus {
	case "", "all":
		return query
	case ReviewStatusPending, ReviewStatusApproved, ReviewStatusRejected:
		return query.Where("videos.review_status = ?", reviewStatus)
	default:
		return query.Where("1 = 0")
	}
}

func incrementVideoCounter(tx *gorm.DB, videoID uint64, column string, delta int) error {
	return tx.Model(&Video{}).
		Where("id = ?", videoID).
		Update(column, gorm.Expr(column+" + ?", delta)).
		Error
}

func decrementVideoCounter(tx *gorm.DB, videoID uint64, column string) error {
	return tx.Model(&Video{}).
		Where("id = ?", videoID).
		Update(column, gorm.Expr("CASE WHEN "+column+" > 0 THEN "+column+" - 1 ELSE 0 END")).
		Error
}

func exists(ctx context.Context, query *gorm.DB) (bool, error) {
	var count int64
	if err := query.WithContext(ctx).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
