package search

import (
	"context"
	"errors"
	"strings"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 50
)

var ErrInvalidKeyword = errors.New("invalid keyword")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchVideos(ctx context.Context, keyword string, page int, pageSize int) (VideoSearchResponse, error) {
	normalizedKeyword := strings.TrimSpace(keyword)
	if normalizedKeyword == "" {
		return VideoSearchResponse{}, ErrInvalidKeyword
	}

	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.SearchVideos(ctx, normalizedKeyword, page, pageSize)
	if err != nil {
		return VideoSearchResponse{}, err
	}

	list := make([]VideoItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, VideoItem{
			ID:              row.ID,
			AreaID:          row.AreaID,
			Title:           row.Title,
			Description:     row.Description,
			CoverURL:        row.CoverURL,
			PlayURL:         row.PlayURL,
			DurationSeconds: row.DurationSeconds,
			ViewCount:       row.ViewCount,
			CommentCount:    row.CommentCount,
			LikeCount:       row.LikeCount,
			FavoriteCount:   row.FavoriteCount,
			PublishedAt:     row.PublishedAt,
			Author: AuthorPreview{
				ID:        row.AuthorID,
				Username:  row.AuthorUsername,
				AvatarURL: row.AuthorAvatarURL,
			},
		})
	}

	return VideoSearchResponse{
		List: list,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *Service) SearchUsers(ctx context.Context, keyword string, page int, pageSize int) (UserSearchResponse, error) {
	normalizedKeyword := strings.TrimSpace(keyword)
	if normalizedKeyword == "" {
		return UserSearchResponse{}, ErrInvalidKeyword
	}

	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.SearchUsers(ctx, normalizedKeyword, page, pageSize)
	if err != nil {
		return UserSearchResponse{}, err
	}

	list := make([]UserItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, UserItem{
			ID:            row.ID,
			Username:      row.Username,
			AvatarURL:     row.AvatarURL,
			Bio:           row.Bio,
			FollowerCount: row.FollowerCount,
			VideoCount:    row.VideoCount,
		})
	}

	return UserSearchResponse{
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
