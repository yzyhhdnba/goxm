package social

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const activeUserStatus = "active"

type Repository struct {
	db *gorm.DB
}

type userRow struct {
	ID        uint64 `gorm:"column:id"`
	Username  string `gorm:"column:username"`
	AvatarURL string `gorm:"column:avatar_url"`
	Bio       string `gorm:"column:bio"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Follow{})
}

// Follow 对应文档“关注关系与 viewer_state.followed”。
// 这里把关系写入和用户计数更新放进一个事务，避免关注关系与粉丝数出现不一致。
func (r *Repository) Follow(ctx context.Context, followerID uint64, followeeID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureActiveUserTx(tx, followerID); err != nil {
			return err
		}
		if err := ensureActiveUserTx(tx, followeeID); err != nil {
			return err
		}

		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&Follow{
			FollowerID: followerID,
			FolloweeID: followeeID,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		if err := incrementUserCounter(tx, followerID, "following_count", 1); err != nil {
			return err
		}
		return incrementUserCounter(tx, followeeID, "follower_count", 1)
	})
}

func (r *Repository) Unfollow(ctx context.Context, followerID uint64, followeeID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureActiveUserTx(tx, followerID); err != nil {
			return err
		}
		if err := ensureActiveUserTx(tx, followeeID); err != nil {
			return err
		}

		result := tx.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&Follow{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		if err := decrementUserCounter(tx, followerID, "following_count"); err != nil {
			return err
		}
		return decrementUserCounter(tx, followeeID, "follower_count")
	})
}

// IsFollowing 是视频详情页和用户主页补 viewer_state.followed 的底层查询。
// 阅读时建议和 video.Service.GetDetail、account.Service.GetProfile 一起对照。
func (r *Repository) IsFollowing(ctx context.Context, followerID uint64, followeeID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) EnsureActiveUser(ctx context.Context, userID uint64) error {
	return ensureActiveUser(r.db.WithContext(ctx), userID)
}

func (r *Repository) ListFollowers(ctx context.Context, userID uint64, page int, pageSize int) ([]userRow, int64, error) {
	if err := ensureActiveUser(r.db.WithContext(ctx), userID); err != nil {
		return nil, 0, err
	}

	base := r.db.WithContext(ctx).
		Table("follows").
		Joins("JOIN users ON users.id = follows.follower_id").
		Where("follows.followee_id = ?", userID).
		Where("users.status = ?", activeUserStatus)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []userRow
	if err := base.
		Select("users.id, users.username, users.avatar_url, users.bio").
		Order("follows.created_at DESC").
		Order("follows.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) ListFollowing(ctx context.Context, userID uint64, page int, pageSize int) ([]userRow, int64, error) {
	if err := ensureActiveUser(r.db.WithContext(ctx), userID); err != nil {
		return nil, 0, err
	}

	base := r.db.WithContext(ctx).
		Table("follows").
		Joins("JOIN users ON users.id = follows.followee_id").
		Where("follows.follower_id = ?", userID).
		Where("users.status = ?", activeUserStatus)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []userRow
	if err := base.
		Select("users.id, users.username, users.avatar_url, users.bio").
		Order("follows.created_at DESC").
		Order("follows.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).
		Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func ensureActiveUserTx(tx *gorm.DB, userID uint64) error {
	return ensureActiveUser(tx, userID)
}

func ensureActiveUser(query *gorm.DB, userID uint64) error {
	var id uint64
	return query.
		Table("users").
		Select("id").
		Where("id = ? AND status = ?", userID, activeUserStatus).
		Take(&id).
		Error
}

func incrementUserCounter(tx *gorm.DB, userID uint64, column string, delta int) error {
	return tx.Table("users").
		Where("id = ?", userID).
		Update(column, gorm.Expr(column+" + ?", delta)).
		Error
}

func decrementUserCounter(tx *gorm.DB, userID uint64, column string) error {
	return tx.Table("users").
		Where("id = ?", userID).
		Update(column, gorm.Expr("CASE WHEN "+column+" > 0 THEN "+column+" - 1 ELSE 0 END")).
		Error
}
