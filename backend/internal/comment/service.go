package comment

import (
	"context"
	"errors"
	"strings"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
	maxContentRunes = 500
)

var (
	ErrVideoNotFound   = errors.New("video not found")
	ErrCommentNotFound = errors.New("comment not found")
	ErrInvalidInput    = errors.New("invalid input")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListComments(ctx context.Context, videoID uint64, page int, pageSize int, viewerID uint64) (ListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListComments(ctx, videoID, page, pageSize)
	if err != nil {
		if IsNotFound(err) {
			return ListResponse{}, ErrVideoNotFound
		}
		return ListResponse{}, err
	}

	return s.buildListResponse(ctx, rows, total, page, pageSize, viewerID)
}

func (s *Service) CreateComment(ctx context.Context, videoID uint64, userID uint64, input CreateInput) (Item, error) {
	content := normalizeContent(input.Content)
	if content == "" || len([]rune(content)) > maxContentRunes {
		return Item{}, ErrInvalidInput
	}

	row, err := s.repo.CreateRootComment(ctx, videoID, userID, content)
	if err != nil {
		if IsNotFound(err) {
			return Item{}, ErrVideoNotFound
		}
		return Item{}, err
	}

	return mapItem(*row, false), nil
}

func (s *Service) ListReplies(ctx context.Context, commentID uint64, page int, pageSize int, viewerID uint64) (ListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListReplies(ctx, commentID, page, pageSize)
	if err != nil {
		if IsNotFound(err) {
			return ListResponse{}, ErrCommentNotFound
		}
		return ListResponse{}, err
	}

	return s.buildListResponse(ctx, rows, total, page, pageSize, viewerID)
}

func (s *Service) CreateReply(ctx context.Context, commentID uint64, userID uint64, input CreateInput) (Item, error) {
	content := normalizeContent(input.Content)
	if content == "" || len([]rune(content)) > maxContentRunes {
		return Item{}, ErrInvalidInput
	}

	row, err := s.repo.CreateReply(ctx, commentID, userID, content)
	if err != nil {
		if IsNotFound(err) {
			return Item{}, ErrCommentNotFound
		}
		return Item{}, err
	}

	return mapItem(*row, false), nil
}

func (s *Service) Like(ctx context.Context, commentID uint64, userID uint64) error {
	if err := s.repo.Like(ctx, commentID, userID); err != nil {
		if IsNotFound(err) {
			return ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Unlike(ctx context.Context, commentID uint64, userID uint64) error {
	if err := s.repo.Unlike(ctx, commentID, userID); err != nil {
		if IsNotFound(err) {
			return ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *Service) LikeStatus(ctx context.Context, commentID uint64, userID uint64) (LikeStatusResponse, error) {
	liked, err := s.repo.HasLike(ctx, commentID, userID)
	if err != nil {
		return LikeStatusResponse{}, err
	}
	return LikeStatusResponse{Liked: liked}, nil
}

func (s *Service) buildListResponse(ctx context.Context, rows []row, total int64, page int, pageSize int, viewerID uint64) (ListResponse, error) {
	likedMap, err := s.repo.LikeMap(ctx, viewerID, rowIDs(rows))
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]Item, 0, len(rows))
	for _, item := range rows {
		items = append(items, mapItem(item, likedMap[item.ID]))
	}

	return ListResponse{
		List: items,
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

func normalizeContent(content string) string {
	return strings.TrimSpace(content)
}

func mapItem(source row, liked bool) Item {
	return Item{
		ID:         source.ID,
		RootID:     source.RootID,
		ParentID:   source.ParentID,
		Content:    source.Content,
		LikeCount:  source.LikeCount,
		ReplyCount: source.ReplyCount,
		CreatedAt:  source.CreatedAt.UTC(),
		User: UserPreview{
			ID:        source.UserID,
			Username:  source.Username,
			AvatarURL: source.UserAvatarURL,
		},
		ViewerState: ViewerState{
			Liked: liked,
		},
	}
}
