package history

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrVideoNotFound = errors.New("video not found")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Report 对应文档“搜索、历史与后台统计”里的观看历史写路径。
// Service 层只做输入兜底和错误语义转换，真正的 upsert 放在 Repository。
func (s *Service) Report(ctx context.Context, userID uint64, input ReportInput) error {
	if input.VideoID == 0 {
		return ErrInvalidInput
	}
	if err := s.repo.Upsert(ctx, userID, input); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVideoNotFound
		}
		return err
	}
	return nil
}

// List 是个人中心历史列表的服务层入口。
// 它把 Repository 返回的行数据收口成前端直接消费的分页结构。
func (s *Service) List(ctx context.Context, userID uint64, page int, pageSize int) (ListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return ListResponse{}, err
	}

	list := make([]Item, 0, len(rows))
	for _, row := range rows {
		list = append(list, Item{
			VideoID:         row.VideoID,
			VideoTitle:      row.VideoTitle,
			CoverURL:        row.CoverURL,
			PlayURL:         row.PlayURL,
			AuthorID:        row.AuthorID,
			AuthorName:      row.AuthorName,
			AreaName:        row.AreaName,
			WatchedAt:       row.WatchedAt,
			ProgressSeconds: row.ProgressSeconds,
			DurationSeconds: row.DurationSeconds,
		})
	}

	return ListResponse{
		List: list,
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
