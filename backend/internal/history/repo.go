package history

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	publicVideoStatus   = "visible"
	approvedReviewState = "approved"
	activeAreaStatus    = "active"
	activeUserStatus    = "active"
)

type Repository struct {
	db *gorm.DB
}

type row struct {
	VideoID         uint64    `gorm:"column:video_id"`
	VideoTitle      string    `gorm:"column:video_title"`
	CoverURL        string    `gorm:"column:cover_url"`
	PlayURL         string    `gorm:"column:play_url"`
	AuthorID        uint64    `gorm:"column:author_id"`
	AuthorName      string    `gorm:"column:author_name"`
	AreaName        string    `gorm:"column:area_name"`
	WatchedAt       time.Time `gorm:"column:watched_at"`
	ProgressSeconds uint      `gorm:"column:progress_seconds"`
	DurationSeconds uint      `gorm:"column:duration_seconds"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&ViewHistory{})
}

// Upsert 对应文档“搜索、历史与后台统计”中的历史上报链路。
// 它使用 (user_id, video_id) 维度的 upsert，保证重复观看只更新进度和时间，不会插入脏重复数据。
func (r *Repository) Upsert(ctx context.Context, userID uint64, input ReportInput) error {
	now := time.Now().UTC()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePublicVideo(tx, input.VideoID); err != nil {
			return err
		}

		history := ViewHistory{
			UserID:          userID,
			VideoID:         input.VideoID,
			ProgressSeconds: input.ProgressSeconds,
			WatchedAt:       now,
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "video_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"progress_seconds": input.ProgressSeconds,
				"watched_at":       now,
				"updated_at":       now,
			}),
		}).Create(&history).Error; err != nil {
			return fmt.Errorf("upsert view history: %w", err)
		}

		return nil
	})
}

// ListByUser 是个人中心“观看历史”页的数据入口。
// 这里会把 history、video、author、area 一次 JOIN 出来，避免前端列表页出现多次补查。
func (r *Repository) ListByUser(ctx context.Context, userID uint64, page int, pageSize int) ([]row, int64, error) {
	base := r.db.WithContext(ctx).
		Model(&ViewHistory{}).
		Where("view_histories.user_id = ?", userID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count histories: %w", err)
	}

	var rows []row
	if err := r.db.WithContext(ctx).
		Table("view_histories").
		Joins("JOIN videos ON videos.id = view_histories.video_id").
		Joins("JOIN users ON users.id = videos.author_id").
		Joins("JOIN areas ON areas.id = videos.area_id").
		Where("view_histories.user_id = ?", userID).
		Where("videos.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Where("areas.deleted_at IS NULL").
		Where("videos.status = ?", publicVideoStatus).
		Where("videos.review_status = ?", approvedReviewState).
		Where("videos.published_at IS NOT NULL").
		Where("users.status = ?", activeUserStatus).
		Where("areas.status = ?", activeAreaStatus).
		Select(`
			view_histories.video_id,
			videos.title AS video_title,
			videos.cover_url,
			videos.play_url,
			users.id AS author_id,
			users.username AS author_name,
			areas.name AS area_name,
			view_histories.watched_at,
			view_histories.progress_seconds,
			videos.duration_seconds
		`).
		Order("view_histories.watched_at DESC").
		Order("view_histories.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, fmt.Errorf("list histories: %w", err)
	}

	return rows, total, nil
}

func ensurePublicVideo(query *gorm.DB, videoID uint64) error {
	var id uint64
	return query.Table("videos").
		Select("id").
		Where("id = ?", videoID).
		Where("deleted_at IS NULL").
		Where("status = ?", publicVideoStatus).
		Where("review_status = ?", approvedReviewState).
		Where("published_at IS NOT NULL").
		Take(&id).
		Error
}
