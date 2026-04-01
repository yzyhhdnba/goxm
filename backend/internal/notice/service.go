package notice

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

var ErrNoticeNotFound = errors.New("notice not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID uint64, page int, pageSize int) (ListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)

	rows, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return ListResponse{}, err
	}

	list := make([]Item, 0, len(rows))
	for _, row := range rows {
		list = append(list, mapItem(row))
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

func (s *Service) MarkRead(ctx context.Context, userID uint64, noticeID uint64) (Item, error) {
	item, err := s.repo.MarkRead(ctx, userID, noticeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Item{}, ErrNoticeNotFound
		}
		return Item{}, err
	}
	return mapItem(*item), nil
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

func mapItem(item Notice) Item {
	return Item{
		ID:             item.ID,
		Type:           item.Type,
		Title:          item.Title,
		Content:        item.Content,
		RelatedVideoID: item.RelatedVideoID,
		Read:           item.IsRead,
		ReadAt:         item.ReadAt,
		CreatedAt:      item.CreatedAt,
	}
}
