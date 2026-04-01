package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	appconfig "pilipili-go/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

var ErrInvalidTokenType = errors.New("invalid token type")

type TokenManager struct {
	cfg appconfig.JWTConfig
}

type Claims struct {
	UserID       uint64 `json:"user_id"`
	TokenVersion uint   `json:"token_version"`
	TokenType    string `json:"token_type"`
	RefreshID    string `json:"refresh_id,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresIn  int64
	RefreshTokenExpiresIn int64
}

func NewTokenManager(cfg appconfig.JWTConfig) *TokenManager {
	return &TokenManager{cfg: cfg}
}

// IssueTokenPair 是双 Token 会话模型的核心入口。
// 对应文档《KEY_CODE_IMPLEMENTATION.md》中的“JWT 登录态、双 Token 和无感刷新”。
func (m *TokenManager) IssueTokenPair(userID uint64, tokenVersion uint) (TokenPair, string, error) {
	now := time.Now()
	refreshID, err := randomID()
	if err != nil {
		return TokenPair{}, "", fmt.Errorf("generate refresh id: %w", err)
	}

	accessClaims := Claims{
		UserID:       userID,
		TokenVersion: tokenVersion,
		TokenType:    TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.cfg.AccessTTLMinute) * time.Minute)),
		},
	}

	refreshClaims := Claims{
		UserID:       userID,
		TokenVersion: tokenVersion,
		TokenType:    TokenTypeRefresh,
		RefreshID:    refreshID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.cfg.RefreshTTLHour) * time.Hour)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(m.cfg.AccessSecret))
	if err != nil {
		return TokenPair{}, "", fmt.Errorf("sign access token: %w", err)
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(m.cfg.RefreshSecret))
	if err != nil {
		return TokenPair{}, "", fmt.Errorf("sign refresh token: %w", err)
	}

	return TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  int64(time.Duration(m.cfg.AccessTTLMinute) * time.Minute / time.Second),
		RefreshTokenExpiresIn: int64(time.Duration(m.cfg.RefreshTTLHour) * time.Hour / time.Second),
	}, HashRefreshID(refreshID), nil
}

func (m *TokenManager) ParseAccessToken(tokenString string) (*Claims, error) {
	return m.parse(tokenString, m.cfg.AccessSecret, TokenTypeAccess)
}

func (m *TokenManager) ParseRefreshToken(tokenString string) (*Claims, error) {
	return m.parse(tokenString, m.cfg.RefreshSecret, TokenTypeRefresh)
}

func (m *TokenManager) parse(tokenString string, secret string, expectedType string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.TokenType != expectedType {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func HashRefreshID(refreshID string) string {
	sum := sha256.Sum256([]byte(refreshID))
	return hex.EncodeToString(sum[:])
}

func randomID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
