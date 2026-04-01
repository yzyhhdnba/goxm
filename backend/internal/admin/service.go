package admin

import (
	"context"
	"errors"
	"strings"
	"time"

	"pilipili-go/backend/internal/account"
	"pilipili-go/backend/pkg/authctx"

	"gorm.io/gorm"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidStatus  = errors.New("invalid review status")
	ErrVideoNotFound  = errors.New("video not found")
	ErrInvalidPayload = errors.New("invalid payload")
)

type Service struct {
	repo        *Repository
	accountRepo *account.Repository
}

func NewService(repo *Repository, accountRepo *account.Repository) *Service {
	return &Service{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (s *Service) ListVideos(ctx context.Context, operatorID uint64, reviewStatus string, page int, pageSize int) (VideoListResponse, error) {
	if err := s.ensureAdmin(ctx, operatorID); err != nil {
		return VideoListResponse{}, err
	}

	normalizedStatus, err := normalizeReviewStatus(reviewStatus)
	if err != nil {
		return VideoListResponse{}, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListVideos(ctx, normalizedStatus, page, pageSize)
	if err != nil {
		return VideoListResponse{}, err
	}

	return VideoListResponse{
		List: mapVideoItems(rows),
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *Service) Approve(ctx context.Context, operatorID uint64, videoID uint64) (VideoItem, error) {
	if videoID == 0 {
		return VideoItem{}, ErrInvalidPayload
	}
	if err := s.ensureAdmin(ctx, operatorID); err != nil {
		return VideoItem{}, err
	}

	row, err := s.repo.Review(ctx, videoID, operatorID, ReviewStatusApproved, "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return VideoItem{}, ErrVideoNotFound
		}
		return VideoItem{}, err
	}
	return mapVideoItem(*row), nil
}

func (s *Service) Reject(ctx context.Context, operatorID uint64, videoID uint64, input ReviewInput) (VideoItem, error) {
	if videoID == 0 || strings.TrimSpace(input.Reason) == "" {
		return VideoItem{}, ErrInvalidPayload
	}
	if err := s.ensureAdmin(ctx, operatorID); err != nil {
		return VideoItem{}, err
	}

	row, err := s.repo.Review(ctx, videoID, operatorID, ReviewStatusRejected, strings.TrimSpace(input.Reason))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return VideoItem{}, ErrVideoNotFound
		}
		return VideoItem{}, err
	}
	return mapVideoItem(*row), nil
}

func (s *Service) GetTodayStats(ctx context.Context, operatorID uint64) (TodayStats, error) {
	if err := s.ensureAdmin(ctx, operatorID); err != nil {
		return TodayStats{}, err
	}
	return s.repo.GetTodayStats(ctx, time.Now())
}

func (s *Service) GetAreaStats(ctx context.Context, operatorID uint64) ([]AreaStatsItem, error) {
	if err := s.ensureAdmin(ctx, operatorID); err != nil {
		return nil, err
	}
	return s.repo.GetAreaStats(ctx)
}

func (s *Service) ensureAdmin(ctx context.Context, userID uint64) error {
	if currentUser, ok := authctx.GetCurrentUserFromContext(ctx); ok {
		if currentUser.ID != 0 && currentUser.ID != userID {
			return ErrForbidden
		}
		if currentUser.Role != account.RoleAdmin {
			return ErrForbidden
		}
		return nil
	}

	user, err := s.accountRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Role != account.RoleAdmin {
		return ErrForbidden
	}
	return nil
}

func normalizeReviewStatus(raw string) (string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return ReviewStatusPending, nil
	}
	switch normalized {
	case ReviewStatusPending, ReviewStatusApproved, ReviewStatusRejected, ReviewStatusReviewed, ReviewStatusAll:
		return normalized, nil
	default:
		return "", ErrInvalidStatus
	}
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

func mapVideoItems(rows []videoRow) []VideoItem {
	items := make([]VideoItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapVideoItem(row))
	}
	return items
}

func mapVideoItem(row videoRow) VideoItem {
	return VideoItem{
		ID:              row.ID,
		AreaID:          row.AreaID,
		AreaName:        row.AreaName,
		Title:           row.Title,
		Description:     row.Description,
		CoverURL:        row.CoverURL,
		PlayURL:         row.PlayURL,
		SourcePath:      row.SourcePath,
		DurationSeconds: row.DurationSeconds,
		ReviewStatus:    row.ReviewStatus,
		ReviewReason:    row.ReviewReason,
		ViewCount:       row.ViewCount,
		CommentCount:    row.CommentCount,
		LikeCount:       row.LikeCount,
		FavoriteCount:   row.FavoriteCount,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		PublishedAt:     row.PublishedAt,
		AuthorID:        row.AuthorID,
		AuthorUsername:  row.AuthorUsername,
	}
}
