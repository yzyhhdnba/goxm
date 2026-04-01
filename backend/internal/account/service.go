package account

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	appauth "pilipili-go/backend/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameTaken     = errors.New("username already exists")
	ErrEmailTaken        = errors.New("email already exists")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrInactiveUser      = errors.New("user is not active")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserNotFound      = errors.New("user not found")
)

const (
	dashboardVideoLimit  = 8
	dashboardFollowLimit = 12
)

type Service struct {
	repo   *Repository
	tokens *appauth.TokenManager
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiresIn  int64  `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

type LoginResponse struct {
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresIn  int64        `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64        `json:"refresh_token_expires_in"`
	User                  UserResponse `json:"user"`
}

type AvailabilityResponse struct {
	Available bool `json:"available"`
}

func NewService(repo *Repository, tokens *appauth.TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (UserResponse, error) {
	username := strings.TrimSpace(input.Username)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := input.Password

	if len(username) < 3 || len(username) > 32 || len(password) < 6 {
		return UserResponse{}, ErrInvalidInput
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return UserResponse{}, ErrInvalidInput
	}

	exists, err := s.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return UserResponse{}, fmt.Errorf("check username: %w", err)
	}
	if exists {
		return UserResponse{}, ErrUsernameTaken
	}

	exists, err = s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return UserResponse{}, fmt.Errorf("check email: %w", err)
	}
	if exists {
		return UserResponse{}, ErrEmailTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("generate password hash: %w", err)
	}

	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         RoleUser,
		Status:       StatusActive,
		TokenVersion: 1,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return user.ToResponse(), nil
}

// Login 对应“登录链路”的服务层入口：校验账号、比对密码、签发双 Token，并落库 refresh_token_hash。
// 阅读时可和前端 header/login.vue 的 login 方法一起对照。
func (s *Service) Login(ctx context.Context, input LoginInput) (LoginResponse, error) {
	identifier := strings.TrimSpace(input.Username)
	user, err := s.repo.FindByUsernameOrEmail(ctx, identifier)
	if err != nil {
		if IsNotFound(err) {
			return LoginResponse{}, ErrInvalidCredential
		}
		return LoginResponse{}, fmt.Errorf("find user: %w", err)
	}
	if user.Status != StatusActive {
		return LoginResponse{}, ErrInactiveUser
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredential
	}

	tokenPair, refreshHash, err := s.tokens.IssueTokenPair(user.ID, user.TokenVersion)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("issue token pair: %w", err)
	}
	if err := s.repo.SaveRefreshTokenHash(ctx, user.ID, refreshHash); err != nil {
		return LoginResponse{}, fmt.Errorf("save refresh token hash: %w", err)
	}

	return LoginResponse{
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresIn:  tokenPair.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: tokenPair.RefreshTokenExpiresIn,
		User:                  user.ToResponse(),
	}, nil
}

// Refresh 对应“会话刷新”主链路：重新解析 refresh token，并校验 token_version 与 refresh_token_hash。
// 这里的 RotateRefreshTokenHash 是避免旧 refresh token 并发续签的关键。
func (s *Service) Refresh(ctx context.Context, input RefreshInput) (TokenResponse, error) {
	claims, err := s.tokens.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		return TokenResponse{}, ErrInvalidToken
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		if IsNotFound(err) {
			return TokenResponse{}, ErrInvalidToken
		}
		return TokenResponse{}, fmt.Errorf("find user: %w", err)
	}
	if user.Status != StatusActive {
		return TokenResponse{}, ErrInactiveUser
	}
	if user.TokenVersion != claims.TokenVersion {
		return TokenResponse{}, ErrInvalidToken
	}
	currentRefreshHash := appauth.HashRefreshID(claims.RefreshID)
	if currentRefreshHash != user.RefreshTokenHash {
		return TokenResponse{}, ErrInvalidToken
	}

	tokenPair, refreshHash, err := s.tokens.IssueTokenPair(user.ID, user.TokenVersion)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("rotate token pair: %w", err)
	}
	if err := s.repo.RotateRefreshTokenHash(ctx, user.ID, user.TokenVersion, currentRefreshHash, refreshHash); err != nil {
		if IsNotFound(err) {
			return TokenResponse{}, ErrInvalidToken
		}
		return TokenResponse{}, fmt.Errorf("save rotated refresh token hash: %w", err)
	}

	return TokenResponse{
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresIn:  tokenPair.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: tokenPair.RefreshTokenExpiresIn,
	}, nil
}

func (s *Service) GetCurrentUser(ctx context.Context, userID uint64) (UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return UserResponse{}, err
	}
	return user.ToResponse(), nil
}

func (s *Service) GetProfile(ctx context.Context, profileUserID uint64, viewerID uint64) (ProfileResponse, error) {
	user, err := s.repo.FindActiveByID(ctx, profileUserID)
	if err != nil {
		if IsNotFound(err) {
			return ProfileResponse{}, ErrUserNotFound
		}
		return ProfileResponse{}, err
	}

	followed := false
	if viewerID != 0 && viewerID != user.ID {
		followed, err = s.repo.IsFollowing(ctx, viewerID, user.ID)
		if err != nil {
			return ProfileResponse{}, err
		}
	}

	return ProfileResponse{
		ID:             user.ID,
		Username:       user.Username,
		AvatarURL:      user.AvatarURL,
		Bio:            user.Bio,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		VideoCount:     user.VideoCount,
		ViewerState: ViewerState{
			Followed: followed,
		},
	}, nil
}

func (s *Service) GetDashboard(ctx context.Context, userID uint64) (DashboardResponse, error) {
	user, err := s.repo.FindActiveByID(ctx, userID)
	if err != nil {
		if IsNotFound(err) {
			return DashboardResponse{}, ErrUserNotFound
		}
		return DashboardResponse{}, err
	}

	recentVideos, err := s.repo.ListRecentPublicVideosByAuthor(ctx, user.ID, dashboardVideoLimit)
	if err != nil {
		return DashboardResponse{}, err
	}

	favoriteVideos, err := s.repo.ListRecentFavoriteVideosByUser(ctx, user.ID, dashboardVideoLimit)
	if err != nil {
		return DashboardResponse{}, err
	}

	followingUsers, err := s.repo.ListFollowingPreview(ctx, user.ID, dashboardFollowLimit)
	if err != nil {
		return DashboardResponse{}, err
	}

	totalViewCount, err := s.repo.SumPublicVideoViewsByAuthor(ctx, user.ID)
	if err != nil {
		return DashboardResponse{}, err
	}

	return DashboardResponse{
		User:           user.ToResponse(),
		Stats:          DashboardStats{TotalViewCount: totalViewCount},
		RecentVideos:   mapDashboardVideoItems(recentVideos),
		FavoriteVideos: mapDashboardVideoItems(favoriteVideos),
		FollowingUsers: mapDashboardUserCards(followingUsers),
	}, nil
}

func (s *Service) Logout(ctx context.Context, userID uint64) error {
	return s.repo.InvalidateTokens(ctx, userID)
}

func (s *Service) CheckUsername(ctx context.Context, username string) (AvailabilityResponse, error) {
	exists, err := s.repo.ExistsByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return AvailabilityResponse{}, err
	}
	return AvailabilityResponse{Available: !exists}, nil
}

func (s *Service) CheckEmail(ctx context.Context, email string) (AvailabilityResponse, error) {
	exists, err := s.repo.ExistsByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return AvailabilityResponse{}, err
	}
	return AvailabilityResponse{Available: !exists}, nil
}

func mapDashboardVideoItems(rows []dashboardVideoRow) []DashboardVideoItem {
	items := make([]DashboardVideoItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, DashboardVideoItem{
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
			Author: DashboardAuthorPreview{
				ID:        row.AuthorID,
				Username:  row.AuthorUsername,
				AvatarURL: row.AuthorAvatarURL,
			},
		})
	}
	return items
}

func mapDashboardUserCards(rows []dashboardUserRow) []DashboardUserCard {
	items := make([]DashboardUserCard, 0, len(rows))
	for _, row := range rows {
		items = append(items, DashboardUserCard{
			ID:        row.ID,
			Username:  row.Username,
			AvatarURL: row.AvatarURL,
			Bio:       row.Bio,
		})
	}
	return items
}
