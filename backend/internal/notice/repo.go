package notice

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Notice{})
}

func (r *Repository) ListByUser(ctx context.Context, userID uint64, page int, pageSize int) ([]Notice, int64, error) {
	base := r.db.WithContext(ctx).
		Model(&Notice{}).
		Where("user_id = ?", userID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count notices: %w", err)
	}

	var notices []Notice
	if err := base.
		Order("created_at DESC").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notices).
		Error; err != nil {
		return nil, 0, fmt.Errorf("list notices: %w", err)
	}

	return notices, total, nil
}

func (r *Repository) MarkRead(ctx context.Context, userID uint64, noticeID uint64) (*Notice, error) {
	var item Notice
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", noticeID, userID).
		First(&item).
		Error; err != nil {
		return nil, err
	}

	if item.IsRead {
		return &item, nil
	}

	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).
		Model(&Notice{}).
		Where("id = ? AND user_id = ?", noticeID, userID).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		}).
		Error; err != nil {
		return nil, err
	}

	item.IsRead = true
	item.ReadAt = &now
	return &item, nil
}
