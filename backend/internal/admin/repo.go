package admin

import (
	"context"
	"fmt"
	"time"

	"pilipili-go/backend/internal/area"
	"pilipili-go/backend/internal/notice"
	"pilipili-go/backend/internal/video"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type videoRow struct {
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
	AuthorID        uint64     `gorm:"column:author_id"`
	AuthorUsername  string     `gorm:"column:author_username"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&VideoReview{})
}

func (r *Repository) ListVideos(ctx context.Context, reviewStatus string, page int, pageSize int) ([]videoRow, int64, error) {
	base := r.videoBase(ctx)
	base = applyReviewStatusFilter(base, reviewStatus)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count admin videos: %w", err)
	}

	var rows []videoRow
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
			videos.published_at,
			users.id AS author_id,
			users.username AS author_username
		`).
		Order("videos.created_at DESC").
		Order("videos.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, fmt.Errorf("list admin videos: %w", err)
	}

	return rows, total, nil
}

// Review 是审核流事务边界所在的位置。
// 对应文档“编辑稿件与审核流”，一次事务里会同时更新稿件状态、作者计数、通知和审核记录。
func (r *Repository) Review(ctx context.Context, videoID uint64, reviewerID uint64, nextStatus string, reason string) (*videoRow, error) {
	now := time.Now().UTC()

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var state reviewState
		if err := tx.Model(&video.Video{}).
			Select("id", "author_id", "title", "review_status", "published_at", "deleted_at").
			Where("id = ?", videoID).
			Take(&state).
			Error; err != nil {
			return err
		}

		wasApproved := state.ReviewStatus == video.ReviewStatusApproved && state.PublishedAt != nil
		updates := map[string]any{
			"review_status": nextStatus,
			"review_reason": reason,
			"updated_at":    now,
		}
		if nextStatus == ReviewStatusApproved {
			if state.PublishedAt == nil {
				updates["published_at"] = now
			}
		} else {
			updates["published_at"] = nil
		}

		if err := tx.Model(&video.Video{}).
			Where("id = ?", videoID).
			Updates(updates).
			Error; err != nil {
			return err
		}

		if !wasApproved && nextStatus == ReviewStatusApproved {
			if err := tx.Model(&accountUserStub{}).
				Where("id = ?", state.AuthorID).
				Update("video_count", gorm.Expr("video_count + 1")).
				Error; err != nil {
				return err
			}
		}
		if wasApproved && nextStatus != ReviewStatusApproved {
			if err := tx.Model(&accountUserStub{}).
				Where("id = ?", state.AuthorID).
				Update("video_count", gorm.Expr("CASE WHEN video_count > 0 THEN video_count - 1 ELSE 0 END")).
				Error; err != nil {
				return err
			}
		}

		noticeItem := buildReviewNotice(state, nextStatus, reason, now, videoID)
		if err := tx.Create(&noticeItem).Error; err != nil {
			return err
		}

		return tx.Create(&VideoReview{
			VideoID:    videoID,
			ReviewerID: reviewerID,
			Status:     nextStatus,
			Reason:     reason,
			CreatedAt:  now,
		}).Error
	}); err != nil {
		return nil, err
	}

	var row videoRow
	if err := r.videoBase(ctx).
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
			videos.published_at,
			users.id AS author_id,
			users.username AS author_username
		`).
		Where("videos.id = ?", videoID).
		Take(&row).
		Error; err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *Repository) GetTodayStats(ctx context.Context, now time.Time) (TodayStats, error) {
	dateKey := now.UTC().Format("2006-01-02")
	stats := TodayStats{}

	if err := r.db.WithContext(ctx).
		Table("view_histories").
		Where("DATE(watched_at) = ?", dateKey).
		Distinct("user_id").
		Count(&stats.ActiveUserCount).
		Error; err != nil {
		return TodayStats{}, fmt.Errorf("count active users: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Model(&video.Video{}).
		Where("DATE(created_at) = ?", dateKey).
		Count(&stats.SubmittedVideoCount).
		Error; err != nil {
		return TodayStats{}, fmt.Errorf("count submitted videos: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Model(&VideoReview{}).
		Where("status = ?", ReviewStatusApproved).
		Where("DATE(created_at) = ?", dateKey).
		Count(&stats.ApprovedVideoCount).
		Error; err != nil {
		return TodayStats{}, fmt.Errorf("count approved videos: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Table("view_histories").
		Where("DATE(watched_at) = ?", dateKey).
		Count(&stats.PlayCount).
		Error; err != nil {
		return TodayStats{}, fmt.Errorf("count plays: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Table("comments").
		Where("status = ?", "visible").
		Where("DATE(created_at) = ?", dateKey).
		Count(&stats.CommentCount).
		Error; err != nil {
		return TodayStats{}, fmt.Errorf("count comments: %w", err)
	}

	return stats, nil
}

func (r *Repository) GetAreaStats(ctx context.Context) ([]AreaStatsItem, error) {
	var rows []areaStatsRow
	if err := r.db.WithContext(ctx).
		Table("areas").
		Joins("LEFT JOIN videos ON videos.area_id = areas.id AND videos.deleted_at IS NULL").
		Where("areas.deleted_at IS NULL").
		Where("areas.status = ?", area.StatusActive).
		Select(`
			areas.id AS area_id,
			areas.name AS area_name,
			SUM(CASE WHEN videos.review_status = 'approved' THEN 1 ELSE 0 END) AS approved_count,
			SUM(CASE WHEN videos.review_status = 'pending' THEN 1 ELSE 0 END) AS pending_count,
			SUM(CASE WHEN videos.review_status = 'rejected' THEN 1 ELSE 0 END) AS rejected_count,
			COUNT(videos.id) AS total_count
		`).
		Group("areas.id, areas.name, areas.sort_order").
		Order("areas.sort_order ASC").
		Order("areas.id ASC").
		Scan(&rows).
		Error; err != nil {
		return nil, fmt.Errorf("list area stats: %w", err)
	}

	result := make([]AreaStatsItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, AreaStatsItem(row))
	}
	return result, nil
}

func (r *Repository) videoBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Table("videos").
		Joins("JOIN users ON users.id = videos.author_id").
		Joins("JOIN areas ON areas.id = videos.area_id").
		Where("videos.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Where("areas.deleted_at IS NULL")
}

func applyReviewStatusFilter(query *gorm.DB, reviewStatus string) *gorm.DB {
	switch reviewStatus {
	case "", ReviewStatusPending:
		return query.Where("videos.review_status = ?", ReviewStatusPending)
	case ReviewStatusApproved:
		return query.Where("videos.review_status = ?", ReviewStatusApproved)
	case ReviewStatusRejected:
		return query.Where("videos.review_status = ?", ReviewStatusRejected)
	case ReviewStatusReviewed:
		return query.Where("videos.review_status IN ?", []string{ReviewStatusApproved, ReviewStatusRejected})
	case ReviewStatusAll:
		return query
	default:
		return query.Where("1 = 0")
	}
}

type accountUserStub struct{}

func (accountUserStub) TableName() string {
	return "users"
}

func buildReviewNotice(state reviewState, nextStatus string, reason string, now time.Time, videoID uint64) notice.Notice {
	return notice.Notice{
		UserID:         state.AuthorID,
		Type:           notice.TypeVideoReview,
		Title:          reviewNoticeTitle(nextStatus),
		Content:        reviewNoticeContent(state.Title, nextStatus, reason),
		RelatedVideoID: &videoID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func reviewNoticeTitle(nextStatus string) string {
	switch nextStatus {
	case ReviewStatusApproved:
		return "稿件审核通过"
	case ReviewStatusRejected:
		return "稿件审核未通过"
	default:
		return "稿件审核状态更新"
	}
}

func reviewNoticeContent(title string, nextStatus string, reason string) string {
	switch nextStatus {
	case ReviewStatusApproved:
		return fmt.Sprintf("你投稿的视频《%s》已通过审核。", title)
	case ReviewStatusRejected:
		if reason != "" {
			return fmt.Sprintf("你投稿的视频《%s》未通过审核，原因：%s", title, reason)
		}
		return fmt.Sprintf("你投稿的视频《%s》未通过审核。", title)
	default:
		return fmt.Sprintf("你投稿的视频《%s》审核状态已更新。", title)
	}
}
