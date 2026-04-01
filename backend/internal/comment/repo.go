package comment

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const reviewStatusApproved = "approved"
const activeUserStatus = "active"

type Repository struct {
	db *gorm.DB
}

type row struct {
	ID            uint64    `gorm:"column:id"`
	RootID        uint64    `gorm:"column:root_id"`
	ParentID      uint64    `gorm:"column:parent_id"`
	Content       string    `gorm:"column:content"`
	LikeCount     uint      `gorm:"column:like_count"`
	ReplyCount    uint      `gorm:"column:reply_count"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UserID        uint64    `gorm:"column:user_id"`
	Username      string    `gorm:"column:username"`
	UserAvatarURL string    `gorm:"column:avatar_url"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Comment{}, &CommentLike{})
}

func (r *Repository) ListComments(ctx context.Context, videoID uint64, page int, pageSize int) ([]row, int64, error) {
	if err := ensurePublicVideo(r.db.WithContext(ctx), videoID); err != nil {
		return nil, 0, err
	}

	countQuery := r.db.WithContext(ctx).
		Model(&Comment{}).
		Where("video_id = ? AND root_id = 0 AND status = ?", videoID, StatusVisible)

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows, err := r.listRows(ctx, buildRowsQuery(r.db.WithContext(ctx), videoID, 0), page, pageSize, "comments.created_at DESC", "comments.id DESC")
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) ListReplies(ctx context.Context, commentID uint64, page int, pageSize int) ([]row, int64, error) {
	baseComment, err := r.findVisibleComment(ctx, commentID)
	if err != nil {
		return nil, 0, err
	}

	rootID := baseComment.ID
	if baseComment.RootID != 0 {
		rootID = baseComment.RootID
	}

	countQuery := r.db.WithContext(ctx).
		Model(&Comment{}).
		Where("root_id = ? AND status = ?", rootID, StatusVisible)

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows, err := r.listRows(ctx, buildRowsQuery(r.db.WithContext(ctx), 0, rootID), page, pageSize, "comments.created_at ASC", "comments.id ASC")
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) CreateRootComment(ctx context.Context, videoID uint64, userID uint64, content string) (*row, error) {
	var comment Comment

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePublicVideo(tx, videoID); err != nil {
			return err
		}
		if err := ensureActiveUser(tx, userID); err != nil {
			return err
		}

		comment = Comment{
			VideoID:  videoID,
			UserID:   userID,
			Content:  content,
			Status:   StatusVisible,
			RootID:   0,
			ParentID: 0,
		}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}

		return incrementVideoCommentCount(tx, videoID, 1)
	})
	if err != nil {
		return nil, err
	}

	return r.loadRowByID(ctx, comment.ID)
}

// CreateReply 是评论树写路径里最关键的函数。
// 对应文档“评论树与回复链路”，这里会同时计算 root_id、写入回复并维护 reply_count / comment_count。
func (r *Repository) CreateReply(ctx context.Context, parentCommentID uint64, userID uint64, content string) (*row, error) {
	var reply Comment

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		parent, err := findVisibleCommentTx(tx, parentCommentID)
		if err != nil {
			return err
		}
		if err := ensureActiveUser(tx, userID); err != nil {
			return err
		}

		rootID := parent.ID
		if parent.RootID != 0 {
			rootID = parent.RootID
		}

		reply = Comment{
			VideoID:  parent.VideoID,
			UserID:   userID,
			RootID:   rootID,
			ParentID: parent.ID,
			Content:  content,
			Status:   StatusVisible,
		}
		if err := tx.Create(&reply).Error; err != nil {
			return err
		}

		if err := tx.Model(&Comment{}).
			Where("id = ?", rootID).
			Update("reply_count", gorm.Expr("reply_count + 1")).
			Error; err != nil {
			return err
		}

		return incrementVideoCommentCount(tx, parent.VideoID, 1)
	})
	if err != nil {
		return nil, err
	}

	return r.loadRowByID(ctx, reply.ID)
}

func (r *Repository) Like(ctx context.Context, commentID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := findVisibleCommentTx(tx, commentID); err != nil {
			return err
		}

		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&CommentLike{
			CommentID: commentID,
			UserID:    userID,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return tx.Model(&Comment{}).
			Where("id = ?", commentID).
			Update("like_count", gorm.Expr("like_count + 1")).
			Error
	})
}

func (r *Repository) Unlike(ctx context.Context, commentID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := findVisibleCommentTx(tx, commentID); err != nil {
			return err
		}

		result := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&CommentLike{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return tx.Model(&Comment{}).
			Where("id = ?", commentID).
			Update("like_count", gorm.Expr("CASE WHEN like_count > 0 THEN like_count - 1 ELSE 0 END")).
			Error
	})
}

func (r *Repository) HasLike(ctx context.Context, commentID uint64, userID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) LikeMap(ctx context.Context, userID uint64, commentIDs []uint64) (map[uint64]bool, error) {
	result := make(map[uint64]bool, len(commentIDs))
	if userID == 0 || len(commentIDs) == 0 {
		return result, nil
	}

	var likes []CommentLike
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND comment_id IN ?", userID, commentIDs).
		Find(&likes).
		Error; err != nil {
		return nil, err
	}

	for _, item := range likes {
		result[item.CommentID] = true
	}
	return result, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *Repository) loadRowByID(ctx context.Context, commentID uint64) (*row, error) {
	rows, err := r.listRows(
		ctx,
		r.db.WithContext(ctx).
			Table("comments").
			Where("comments.id = ?", commentID),
		1,
		1,
		"comments.created_at ASC",
		"comments.id ASC",
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &rows[0], nil
}

func (r *Repository) listRows(ctx context.Context, query *gorm.DB, page int, pageSize int, orders ...string) ([]row, error) {
	var rows []row
	base := query.
		Select(`
			comments.id,
			comments.root_id,
			comments.parent_id,
			comments.content,
			comments.like_count,
			comments.reply_count,
			comments.created_at,
			users.id as user_id,
			users.username,
			users.avatar_url
		`).
		Joins("JOIN users ON users.id = comments.user_id")
	for _, order := range orders {
		base = base.Order(order)
	}
	if err := base.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) findVisibleComment(ctx context.Context, commentID uint64) (*Comment, error) {
	return findVisibleCommentTx(r.db.WithContext(ctx), commentID)
}

func findVisibleCommentTx(tx *gorm.DB, commentID uint64) (*Comment, error) {
	var item Comment
	if err := tx.Where("id = ? AND status = ?", commentID, StatusVisible).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensurePublicVideo(query *gorm.DB, videoID uint64) error {
	var id uint64
	return query.Table("videos").
		Select("id").
		Where("id = ?", videoID).
		Where("status = ?", "visible").
		Where("review_status = ?", reviewStatusApproved).
		Where("published_at IS NOT NULL").
		Take(&id).
		Error
}

func ensureActiveUser(query *gorm.DB, userID uint64) error {
	var id uint64
	return query.Table("users").
		Select("id").
		Where("id = ? AND status = ?", userID, activeUserStatus).
		Take(&id).
		Error
}

func incrementVideoCommentCount(tx *gorm.DB, videoID uint64, delta int) error {
	return tx.Table("videos").
		Where("id = ?", videoID).
		Update("comment_count", gorm.Expr("comment_count + ?", delta)).
		Error
}

func rowIDs(items []row) []uint64 {
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func buildRowsQuery(db *gorm.DB, videoID uint64, rootID uint64) *gorm.DB {
	query := db.Table("comments").Where("comments.status = ?", StatusVisible)
	if videoID != 0 {
		query = query.Where("comments.video_id = ? AND comments.root_id = 0", videoID)
	}
	if rootID != 0 {
		query = query.Where("comments.root_id = ?", rootID)
	}
	return query
}
