package social

import (
	"context"
	"errors"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrFollowTargetNotFound = errors.New("follow target not found")
	ErrCannotFollowSelf     = errors.New("cannot follow self")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Follow 对应文档“关注关系与 viewer_state.followed”的服务层入口。
// 这里负责拦截自关注这类业务非法场景，再把一致性要求交给 Repository 事务处理。
func (s *Service) Follow(ctx context.Context, followerID uint64, followeeID uint64) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}
	if err := s.repo.Follow(ctx, followerID, followeeID); err != nil {
		if IsNotFound(err) {
			return ErrFollowTargetNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Unfollow(ctx context.Context, followerID uint64, followeeID uint64) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}
	if err := s.repo.Unfollow(ctx, followerID, followeeID); err != nil {
		if IsNotFound(err) {
			return ErrFollowTargetNotFound
		}
		return err
	}
	return nil
}

// Status 是详情页和用户主页查询“是否已关注”的服务层桥接点。
// 阅读时建议和前端 videoDetail.vue 里的关注按钮初始化一起对照。
func (s *Service) Status(ctx context.Context, followerID uint64, followeeID uint64) (FollowStatusResponse, error) {
	if err := s.repo.EnsureActiveUser(ctx, followeeID); err != nil {
		if IsNotFound(err) {
			return FollowStatusResponse{}, ErrFollowTargetNotFound
		}
		return FollowStatusResponse{}, err
	}
	followed, err := s.repo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return FollowStatusResponse{}, err
	}
	return FollowStatusResponse{Followed: followed}, nil
}

// ListFollowers / ListFollowing 提供社交页所需的标准分页输出。
func (s *Service) ListFollowers(ctx context.Context, userID uint64, page int, pageSize int) (UserListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListFollowers(ctx, userID, page, pageSize)
	if err != nil {
		if IsNotFound(err) {
			return UserListResponse{}, ErrFollowTargetNotFound
		}
		return UserListResponse{}, err
	}

	return UserListResponse{
		List: mapUserRows(rows),
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *Service) ListFollowing(ctx context.Context, userID uint64, page int, pageSize int) (UserListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListFollowing(ctx, userID, page, pageSize)
	if err != nil {
		if IsNotFound(err) {
			return UserListResponse{}, ErrFollowTargetNotFound
		}
		return UserListResponse{}, err
	}

	return UserListResponse{
		List: mapUserRows(rows),
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func normalizePagination(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > maxPageSize:
		pageSize = maxPageSize
	}
	return page, pageSize
}

func mapUserRows(rows []userRow) []UserCard {
	items := make([]UserCard, 0, len(rows))
	for _, row := range rows {
		items = append(items, UserCard{
			ID:        row.ID,
			Username:  row.Username,
			AvatarURL: row.AvatarURL,
			Bio:       row.Bio,
		})
	}
	return items
}
