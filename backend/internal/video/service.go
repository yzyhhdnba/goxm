package video

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"pilipili-go/backend/internal/media"
)

const (
	defaultRecommendLimit = 8
	maxRecommendLimit     = 30
	defaultFeedLimit      = 20
	maxFeedLimit          = 50
	defaultPage           = 1
	defaultPageSize       = 20
	maxPageSize           = 100
)

var (
	ErrInvalidCursor       = errors.New("invalid cursor")
	ErrInvalidSort         = errors.New("invalid sort")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidReviewStatus = errors.New("invalid review status")
	ErrVideoNotFound       = errors.New("video not found")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrAreaNotFound        = errors.New("area not found")
)

const areaSortLatest = "latest"

type Service struct {
	repo          *Repository
	followChecker FollowChecker
	storage       MediaStorage
}

type FollowChecker interface {
	IsFollowing(ctx context.Context, followerID uint64, followeeID uint64) (bool, error)
}

type MediaStorage interface {
	SaveVideoSource(videoID uint64, file *multipart.FileHeader) (media.StoredObject, error)
	SaveVideoCover(videoID uint64, file *multipart.FileHeader) (media.StoredObject, error)
}

func NewService(repo *Repository, followChecker FollowChecker, storage MediaStorage) *Service {
	return &Service{
		repo:          repo,
		followChecker: followChecker,
		storage:       storage,
	}
}

// ListRecommend 对应文档中的“推荐流 / cursor 分页”章节。
// 关键点是多取一条数据来判断 has_more，并把最后一条记录编码成 next_cursor。
func (s *Service) ListRecommend(ctx context.Context, rawCursor string, limit int) (FeedResponse, error) {
	cursor, err := parseCursor(rawCursor)
	if err != nil {
		return FeedResponse{}, err
	}

	normalizedLimit := normalizeRecommendLimit(limit)
	rows, err := s.repo.ListRecommend(ctx, cursor, normalizedLimit+1)
	if err != nil {
		return FeedResponse{}, err
	}

	return buildFeedResponse(rows, normalizedLimit, func(last videoRow) string {
		return encodeCursor(last.PublishedAt, last.ID)
	}), nil
}

func (s *Service) ListHot(ctx context.Context, rawCursor string, limit int) (FeedResponse, error) {
	cursor, err := parseHotCursor(rawCursor)
	if err != nil {
		return FeedResponse{}, err
	}

	normalizedLimit := normalizeFeedLimit(limit)
	rows, err := s.repo.ListHot(ctx, cursor, normalizedLimit+1)
	if err != nil {
		return FeedResponse{}, err
	}

	return buildFeedResponse(rows, normalizedLimit, func(last videoRow) string {
		return encodeHotCursor(last.HotScore, last.PublishedAt, last.ID)
	}), nil
}

func (s *Service) ListFollowing(ctx context.Context, userID uint64, rawCursor string, limit int) (FeedResponse, error) {
	cursor, err := parseCursor(rawCursor)
	if err != nil {
		return FeedResponse{}, err
	}

	normalizedLimit := normalizeFeedLimit(limit)
	rows, err := s.repo.ListFollowing(ctx, userID, cursor, normalizedLimit+1)
	if err != nil {
		return FeedResponse{}, err
	}

	return buildFeedResponse(rows, normalizedLimit, func(last videoRow) string {
		return encodeCursor(last.PublishedAt, last.ID)
	}), nil
}

func (s *Service) ListByArea(ctx context.Context, areaID uint64, sort string, rawCursor string, limit int) (FeedResponse, error) {
	normalizedSort, err := normalizeAreaSort(sort)
	if err != nil {
		return FeedResponse{}, err
	}

	cursor, err := parseCursor(rawCursor)
	if err != nil {
		return FeedResponse{}, err
	}

	normalizedLimit := normalizeFeedLimit(limit)

	var rows []videoRow
	switch normalizedSort {
	case areaSortLatest:
		rows, err = s.repo.ListByArea(ctx, areaID, cursor, normalizedLimit+1)
	default:
		return FeedResponse{}, ErrInvalidSort
	}
	if err != nil {
		if IsNotFound(err) {
			return FeedResponse{}, ErrAreaNotFound
		}
		return FeedResponse{}, err
	}

	return buildFeedResponse(rows, normalizedLimit, func(last videoRow) string {
		return encodeCursor(last.PublishedAt, last.ID)
	}), nil
}

func (s *Service) ListByAuthor(ctx context.Context, authorID uint64, page int, pageSize int) (VideoListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListByAuthor(ctx, authorID, page, pageSize)
	if err != nil {
		if IsNotFound(err) {
			return VideoListResponse{}, ErrAuthorNotFound
		}
		return VideoListResponse{}, err
	}

	return VideoListResponse{
		List: mapFeedItems(rows),
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

// GetDetail 是视频详情页的核心聚合入口。
// 对应文档“视频主链路：推荐流、热门流、详情页”以及“viewer_state 链路”。
func (s *Service) GetDetail(ctx context.Context, videoID uint64, viewerID uint64) (DetailResponse, error) {
	row, err := s.repo.FindPublicByID(ctx, videoID)
	if err != nil {
		if IsNotFound(err) {
			return DetailResponse{}, ErrVideoNotFound
		}
		return DetailResponse{}, err
	}

	viewerState := ViewerState{}
	if viewerID != 0 {
		liked, err := s.repo.HasLike(ctx, videoID, viewerID)
		if err != nil {
			return DetailResponse{}, err
		}
		favorited, err := s.repo.HasFavorite(ctx, videoID, viewerID)
		if err != nil {
			return DetailResponse{}, err
		}
		followed := false
		if s.followChecker != nil && viewerID != row.AuthorID {
			followed, err = s.followChecker.IsFollowing(ctx, viewerID, row.AuthorID)
			if err != nil {
				return DetailResponse{}, err
			}
		}

		viewerState = ViewerState{
			Liked:     liked,
			Favorited: favorited,
			Followed:  followed,
		}
	}

	return DetailResponse{
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
		ViewerState: viewerState,
	}, nil
}

func (s *Service) Like(ctx context.Context, videoID uint64, userID uint64) error {
	if err := s.repo.Like(ctx, videoID, userID); err != nil {
		if IsNotFound(err) {
			return ErrVideoNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Unlike(ctx context.Context, videoID uint64, userID uint64) error {
	if err := s.repo.Unlike(ctx, videoID, userID); err != nil {
		if IsNotFound(err) {
			return ErrVideoNotFound
		}
		return err
	}
	return nil
}

func (s *Service) LikeStatus(ctx context.Context, videoID uint64, userID uint64) (LikeStatusResponse, error) {
	if err := s.repo.EnsurePublic(ctx, videoID); err != nil {
		if IsNotFound(err) {
			return LikeStatusResponse{}, ErrVideoNotFound
		}
		return LikeStatusResponse{}, err
	}
	liked, err := s.repo.HasLike(ctx, videoID, userID)
	if err != nil {
		return LikeStatusResponse{}, err
	}
	return LikeStatusResponse{Liked: liked}, nil
}

func (s *Service) Favorite(ctx context.Context, videoID uint64, userID uint64) error {
	if err := s.repo.Favorite(ctx, videoID, userID); err != nil {
		if IsNotFound(err) {
			return ErrVideoNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Unfavorite(ctx context.Context, videoID uint64, userID uint64) error {
	if err := s.repo.Unfavorite(ctx, videoID, userID); err != nil {
		if IsNotFound(err) {
			return ErrVideoNotFound
		}
		return err
	}
	return nil
}

func (s *Service) FavoriteStatus(ctx context.Context, videoID uint64, userID uint64) (FavoriteStatusResponse, error) {
	if err := s.repo.EnsurePublic(ctx, videoID); err != nil {
		if IsNotFound(err) {
			return FavoriteStatusResponse{}, ErrVideoNotFound
		}
		return FavoriteStatusResponse{}, err
	}
	favorited, err := s.repo.HasFavorite(ctx, videoID, userID)
	if err != nil {
		return FavoriteStatusResponse{}, err
	}
	return FavoriteStatusResponse{Favorited: favorited}, nil
}

// CreateVideo 是投稿两段式流程的第一步：先创建元数据，再由前端继续上传源文件和封面。
// 对应文档“投稿、文件上传与媒体落盘”。
func (s *Service) CreateVideo(ctx context.Context, authorID uint64, input CreateVideoInput) (CreateVideoResponse, error) {
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)
	if input.AreaID == 0 || title == "" || len(title) > 128 {
		return CreateVideoResponse{}, ErrInvalidInput
	}

	item := &Video{
		AuthorID:     authorID,
		AreaID:       input.AreaID,
		Title:        title,
		Description:  description,
		Status:       StatusVisible,
		ReviewStatus: ReviewStatusPending,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		if IsNotFound(err) {
			return CreateVideoResponse{}, ErrAreaNotFound
		}
		return CreateVideoResponse{}, err
	}

	return CreateVideoResponse{
		ID:           item.ID,
		AreaID:       item.AreaID,
		Title:        item.Title,
		Description:  item.Description,
		ReviewStatus: item.ReviewStatus,
		CreatedAt:    item.CreatedAt,
	}, nil
}

func (s *Service) UploadSource(ctx context.Context, videoID uint64, authorID uint64, file *multipart.FileHeader) (SourceUploadResponse, error) {
	if videoID == 0 || file == nil || s.storage == nil {
		return SourceUploadResponse{}, ErrInvalidInput
	}
	if _, err := s.repo.FindOwnedByID(ctx, videoID, authorID); err != nil {
		if IsNotFound(err) {
			return SourceUploadResponse{}, ErrVideoNotFound
		}
		return SourceUploadResponse{}, err
	}

	stored, err := s.storage.SaveVideoSource(videoID, file)
	if err != nil {
		return SourceUploadResponse{}, ErrInvalidInput
	}
	if err := s.repo.UpdateSourceByOwner(ctx, videoID, authorID, stored.RelativePath, stored.PublicURL); err != nil {
		if IsNotFound(err) {
			return SourceUploadResponse{}, ErrVideoNotFound
		}
		return SourceUploadResponse{}, err
	}

	return SourceUploadResponse{
		VideoID:    videoID,
		SourcePath: stored.RelativePath,
		PlayURL:    stored.PublicURL,
	}, nil
}

func (s *Service) UploadCover(ctx context.Context, videoID uint64, authorID uint64, file *multipart.FileHeader) (CoverUploadResponse, error) {
	if videoID == 0 || file == nil || s.storage == nil {
		return CoverUploadResponse{}, ErrInvalidInput
	}
	if _, err := s.repo.FindOwnedByID(ctx, videoID, authorID); err != nil {
		if IsNotFound(err) {
			return CoverUploadResponse{}, ErrVideoNotFound
		}
		return CoverUploadResponse{}, err
	}

	stored, err := s.storage.SaveVideoCover(videoID, file)
	if err != nil {
		return CoverUploadResponse{}, ErrInvalidInput
	}
	if err := s.repo.UpdateCoverByOwner(ctx, videoID, authorID, stored.PublicURL); err != nil {
		if IsNotFound(err) {
			return CoverUploadResponse{}, ErrVideoNotFound
		}
		return CoverUploadResponse{}, err
	}

	return CoverUploadResponse{
		VideoID:  videoID,
		CoverURL: stored.PublicURL,
	}, nil
}

func (s *Service) UpdateVideo(ctx context.Context, videoID uint64, authorID uint64, input UpdateVideoInput) (CreatorVideoItem, error) {
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)
	if videoID == 0 || input.AreaID == 0 || title == "" || len(title) > 128 {
		return CreatorVideoItem{}, ErrInvalidInput
	}

	row, err := s.repo.UpdateMetadataByOwner(ctx, videoID, authorID, UpdateVideoInput{
		AreaID:      input.AreaID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		switch {
		case errors.Is(err, errAreaNotFound):
			return CreatorVideoItem{}, ErrAreaNotFound
		case IsNotFound(err):
			return CreatorVideoItem{}, ErrVideoNotFound
		default:
			return CreatorVideoItem{}, err
		}
	}

	return mapCreatorVideoItem(*row), nil
}

func (s *Service) ListCreatorVideos(ctx context.Context, authorID uint64, reviewStatus string, page int, pageSize int) (CreatorVideoListResponse, error) {
	normalizedStatus, err := normalizeCreatorReviewStatus(reviewStatus)
	if err != nil {
		return CreatorVideoListResponse{}, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	rows, total, err := s.repo.ListByCreator(ctx, authorID, normalizedStatus, page, pageSize)
	if err != nil {
		if IsNotFound(err) {
			return CreatorVideoListResponse{}, ErrAuthorNotFound
		}
		return CreatorVideoListResponse{}, err
	}

	list := make([]CreatorVideoItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, mapCreatorVideoItem(row))
	}

	return CreatorVideoListResponse{
		List: list,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func parseCursor(raw string) (*FeedCursor, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, ":")
	if len(parts) > 2 {
		return nil, ErrInvalidCursor
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	cursor := &FeedCursor{
		PublishedAt: time.Unix(timestamp, 0).UTC(),
	}

	if len(parts) == 2 {
		id, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		cursor.ID = id
		cursor.HasID = true
	}

	return cursor, nil
}

func encodeCursor(publishedAt time.Time, videoID uint64) string {
	return fmt.Sprintf("%d:%d", publishedAt.UTC().Unix(), videoID)
}

func parseHotCursor(raw string) (*HotFeedCursor, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, ":")
	if len(parts) != 3 {
		return nil, ErrInvalidCursor
	}

	hotScore, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	videoID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	return &HotFeedCursor{
		HotScore:    hotScore,
		PublishedAt: time.Unix(timestamp, 0).UTC(),
		ID:          videoID,
	}, nil
}

func encodeHotCursor(hotScore int64, publishedAt time.Time, videoID uint64) string {
	return fmt.Sprintf("%d:%d:%d", hotScore, publishedAt.UTC().Unix(), videoID)
}

func normalizeRecommendLimit(limit int) int {
	switch {
	case limit <= 0:
		return defaultRecommendLimit
	case limit > maxRecommendLimit:
		return maxRecommendLimit
	default:
		return limit
	}
}

func mapCreatorVideoItem(row creatorRow) CreatorVideoItem {
	return CreatorVideoItem{
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
	}
}

func normalizeFeedLimit(limit int) int {
	switch {
	case limit <= 0:
		return defaultFeedLimit
	case limit > maxFeedLimit:
		return maxFeedLimit
	default:
		return limit
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

func normalizeAreaSort(sort string) (string, error) {
	normalized := strings.TrimSpace(sort)
	if normalized == "" {
		return areaSortLatest, nil
	}
	if normalized != areaSortLatest {
		return "", ErrInvalidSort
	}
	return normalized, nil
}

func normalizeCreatorReviewStatus(raw string) (string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "all", nil
	}
	switch normalized {
	case "all", ReviewStatusPending, ReviewStatusApproved, ReviewStatusRejected:
		return normalized, nil
	default:
		return "", ErrInvalidReviewStatus
	}
}

func buildFeedResponse(rows []videoRow, limit int, nextCursor func(videoRow) string) FeedResponse {
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	next := ""
	if hasMore && len(rows) > 0 {
		next = nextCursor(rows[len(rows)-1])
	}

	return FeedResponse{
		Items:      mapFeedItems(rows),
		NextCursor: next,
		HasMore:    hasMore,
	}
}

func mapFeedItems(rows []videoRow) []FeedItem {
	items := make([]FeedItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapFeedItem(row))
	}
	return items
}

func mapFeedItem(row videoRow) FeedItem {
	return FeedItem{
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
	}
}
