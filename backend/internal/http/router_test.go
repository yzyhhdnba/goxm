package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pilipili-go/backend/internal/account"
	"pilipili-go/backend/internal/admin"
	"pilipili-go/backend/internal/area"
	"pilipili-go/backend/internal/comment"
	appconfig "pilipili-go/backend/internal/config"
	"pilipili-go/backend/internal/history"
	apphttp "pilipili-go/backend/internal/http"
	"pilipili-go/backend/internal/notice"
	"pilipili-go/backend/internal/search"
	"pilipili-go/backend/internal/social"
	"pilipili-go/backend/internal/video"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type envelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type loginData struct {
	AccessToken           string               `json:"access_token"`
	RefreshToken          string               `json:"refresh_token"`
	AccessTokenExpiresIn  int64                `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64                `json:"refresh_token_expires_in"`
	User                  account.UserResponse `json:"user"`
}

type tokenData struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiresIn  int64  `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

type availabilityData struct {
	Available bool `json:"available"`
}

func TestAuthFlowRefreshRotationAndLogout(t *testing.T) {
	router, _ := newTestApp(t)

	emptyUsernameResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/auth/check-username?username=", nil, "")
	if emptyUsernameResp.Code != stdhttp.StatusBadRequest {
		t.Fatalf("empty username status = %d, want %d, body=%s", emptyUsernameResp.Code, stdhttp.StatusBadRequest, emptyUsernameResp.Body.String())
	}
	emptyUsernameResult := decodeBody[envelope[map[string]any]](t, emptyUsernameResp)
	if emptyUsernameResult.Code != 2012 {
		t.Fatalf("empty username error code = %d, want %d", emptyUsernameResult.Code, 2012)
	}

	emptyEmailResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/auth/check-email?email=", nil, "")
	if emptyEmailResp.Code != stdhttp.StatusBadRequest {
		t.Fatalf("empty email status = %d, want %d, body=%s", emptyEmailResp.Code, stdhttp.StatusBadRequest, emptyEmailResp.Body.String())
	}
	emptyEmailResult := decodeBody[envelope[map[string]any]](t, emptyEmailResp)
	if emptyEmailResult.Code != 2013 {
		t.Fatalf("empty email error code = %d, want %d", emptyEmailResult.Code, 2013)
	}

	registerBody := map[string]any{
		"username": "alice_batch_b_test",
		"email":    "alice_batch_b_test@example.com",
		"password": "12345678",
	}
	registerResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/register", registerBody, "")
	if registerResp.Code != stdhttp.StatusCreated {
		t.Fatalf("register status = %d, want %d, body=%s", registerResp.Code, stdhttp.StatusCreated, registerResp.Body.String())
	}

	usernameCheckResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/auth/check-username?username=alice_batch_b_test", nil, "")
	if usernameCheckResp.Code != stdhttp.StatusOK {
		t.Fatalf("check username status = %d, want %d, body=%s", usernameCheckResp.Code, stdhttp.StatusOK, usernameCheckResp.Body.String())
	}
	usernameCheck := decodeBody[envelope[availabilityData]](t, usernameCheckResp)
	if usernameCheck.Data.Available {
		t.Fatalf("expected username to be unavailable after registration")
	}

	emailCheckResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/auth/check-email?email=alice_batch_b_test@example.com", nil, "")
	if emailCheckResp.Code != stdhttp.StatusOK {
		t.Fatalf("check email status = %d, want %d, body=%s", emailCheckResp.Code, stdhttp.StatusOK, emailCheckResp.Body.String())
	}
	emailCheck := decodeBody[envelope[availabilityData]](t, emailCheckResp)
	if emailCheck.Data.Available {
		t.Fatalf("expected email to be unavailable after registration")
	}

	loginBody := map[string]any{
		"username": "alice_batch_b_test",
		"password": "12345678",
	}
	loginResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/login", loginBody, "")
	if loginResp.Code != stdhttp.StatusOK {
		t.Fatalf("login status = %d, want %d, body=%s", loginResp.Code, stdhttp.StatusOK, loginResp.Body.String())
	}
	loginResult := decodeBody[envelope[loginData]](t, loginResp)

	meResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/me", nil, loginResult.Data.AccessToken)
	if meResp.Code != stdhttp.StatusOK {
		t.Fatalf("me status = %d, want %d, body=%s", meResp.Code, stdhttp.StatusOK, meResp.Body.String())
	}

	refreshBody := map[string]any{
		"refresh_token": loginResult.Data.RefreshToken,
	}
	refreshResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/refresh", refreshBody, "")
	if refreshResp.Code != stdhttp.StatusOK {
		t.Fatalf("refresh status = %d, want %d, body=%s", refreshResp.Code, stdhttp.StatusOK, refreshResp.Body.String())
	}
	refreshResult := decodeBody[envelope[tokenData]](t, refreshResp)
	if refreshResult.Data.RefreshToken == loginResult.Data.RefreshToken {
		t.Fatalf("expected refresh token to rotate")
	}

	oldAccessResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/me", nil, loginResult.Data.AccessToken)
	if oldAccessResp.Code != stdhttp.StatusOK {
		t.Fatalf("old access token should remain valid until expiry, got status=%d body=%s", oldAccessResp.Code, oldAccessResp.Body.String())
	}

	oldRefreshResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/refresh", refreshBody, "")
	if oldRefreshResp.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("old refresh token status = %d, want %d, body=%s", oldRefreshResp.Code, stdhttp.StatusUnauthorized, oldRefreshResp.Body.String())
	}
	oldRefreshResult := decodeBody[envelope[map[string]any]](t, oldRefreshResp)
	if oldRefreshResult.Code != 2009 {
		t.Fatalf("old refresh error code = %d, want %d", oldRefreshResult.Code, 2009)
	}

	logoutResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/logout", nil, refreshResult.Data.AccessToken)
	if logoutResp.Code != stdhttp.StatusOK {
		t.Fatalf("logout status = %d, want %d, body=%s", logoutResp.Code, stdhttp.StatusOK, logoutResp.Body.String())
	}

	revokedAccessResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/me", nil, refreshResult.Data.AccessToken)
	if revokedAccessResp.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("revoked access token status = %d, want %d, body=%s", revokedAccessResp.Code, stdhttp.StatusUnauthorized, revokedAccessResp.Body.String())
	}
	revokedAccessResult := decodeBody[envelope[map[string]any]](t, revokedAccessResp)
	if revokedAccessResult.Code != 4003 {
		t.Fatalf("revoked access error code = %d, want %d", revokedAccessResult.Code, 4003)
	}
}

func TestBatchCAreasRecommendAndDetail(t *testing.T) {
	router, db := newTestApp(t)

	seedBatchCData(t, db)

	areasResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/areas", nil, "")
	if areasResp.Code != stdhttp.StatusOK {
		t.Fatalf("areas status = %d, want %d, body=%s", areasResp.Code, stdhttp.StatusOK, areasResp.Body.String())
	}
	areasResult := decodeBody[envelope[[]area.AreaResponse]](t, areasResp)
	if len(areasResult.Data) < 6 {
		t.Fatalf("areas length = %d, want at least 6", len(areasResult.Data))
	}

	recommendResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/feed/recommend?limit=2", nil, "")
	if recommendResp.Code != stdhttp.StatusOK {
		t.Fatalf("recommend status = %d, want %d, body=%s", recommendResp.Code, stdhttp.StatusOK, recommendResp.Body.String())
	}
	recommendResult := decodeBody[envelope[video.FeedResponse]](t, recommendResp)
	if len(recommendResult.Data.Items) != 2 {
		t.Fatalf("recommend item length = %d, want %d", len(recommendResult.Data.Items), 2)
	}
	if !recommendResult.Data.HasMore {
		t.Fatalf("recommend has_more = false, want true")
	}
	if recommendResult.Data.Items[0].ID != 1003 || recommendResult.Data.Items[1].ID != 1002 {
		t.Fatalf("unexpected recommend order: %+v", recommendResult.Data.Items)
	}
	if recommendResult.Data.NextCursor == "" {
		t.Fatalf("recommend next_cursor should not be empty")
	}

	recommendNextResp := performJSONRequest(
		t,
		router,
		stdhttp.MethodGet,
		"/api/v1/feed/recommend?limit=2&cursor="+recommendResult.Data.NextCursor,
		nil,
		"",
	)
	if recommendNextResp.Code != stdhttp.StatusOK {
		t.Fatalf("recommend next status = %d, want %d, body=%s", recommendNextResp.Code, stdhttp.StatusOK, recommendNextResp.Body.String())
	}
	recommendNextResult := decodeBody[envelope[video.FeedResponse]](t, recommendNextResp)
	if len(recommendNextResult.Data.Items) != 1 {
		t.Fatalf("recommend next item length = %d, want %d", len(recommendNextResult.Data.Items), 1)
	}
	if recommendNextResult.Data.Items[0].ID != 1001 {
		t.Fatalf("recommend next first id = %d, want %d", recommendNextResult.Data.Items[0].ID, 1001)
	}
	if recommendNextResult.Data.HasMore {
		t.Fatalf("recommend next has_more = true, want false")
	}

	detailResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/videos/1002", nil, "")
	if detailResp.Code != stdhttp.StatusOK {
		t.Fatalf("detail status = %d, want %d, body=%s", detailResp.Code, stdhttp.StatusOK, detailResp.Body.String())
	}
	detailResult := decodeBody[envelope[video.DetailResponse]](t, detailResp)
	if detailResult.Data.ID != 1002 {
		t.Fatalf("detail id = %d, want %d", detailResult.Data.ID, 1002)
	}
	if detailResult.Data.Author.ID != 101 {
		t.Fatalf("detail author id = %d, want %d", detailResult.Data.Author.ID, 101)
	}
	if detailResult.Data.ViewerState.Liked || detailResult.Data.ViewerState.Favorited || detailResult.Data.ViewerState.Followed {
		t.Fatalf("viewer_state should default to false before social tables are added")
	}
}

func TestBatchDInteractionsCommentsSocialAndCORS(t *testing.T) {
	router, db := newTestApp(t)
	seedBatchCData(t, db)

	preflightReq := httptest.NewRequest(stdhttp.MethodOptions, "/api/v1/ping", nil)
	preflightReq.Header.Set("Origin", "http://localhost:5173")
	preflightReq.Header.Set("Access-Control-Request-Method", stdhttp.MethodGet)
	preflightReq.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type")
	preflightResp := httptest.NewRecorder()
	router.ServeHTTP(preflightResp, preflightReq)
	if preflightResp.Code != stdhttp.StatusNoContent {
		t.Fatalf("preflight status = %d, want %d, body=%s", preflightResp.Code, stdhttp.StatusNoContent, preflightResp.Body.String())
	}
	if got := preflightResp.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("allow origin = %q, want %q", got, "http://localhost:5173")
	}

	viewerToken, viewerUser := registerAndLoginTestUser(t, router, "viewer_batch_d", "viewer_batch_d@example.com")

	followResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/users/101/follow", nil, viewerToken)
	if followResp.Code != stdhttp.StatusOK {
		t.Fatalf("follow status = %d, want %d, body=%s", followResp.Code, stdhttp.StatusOK, followResp.Body.String())
	}

	followStatusResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/follow-status", nil, viewerToken)
	if followStatusResp.Code != stdhttp.StatusOK {
		t.Fatalf("follow status api = %d, want %d, body=%s", followStatusResp.Code, stdhttp.StatusOK, followStatusResp.Body.String())
	}
	followStatus := decodeBody[envelope[social.FollowStatusResponse]](t, followStatusResp)
	if !followStatus.Data.Followed {
		t.Fatalf("follow status should be true after follow")
	}

	followersResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/followers?page=1&page_size=20", nil, "")
	if followersResp.Code != stdhttp.StatusOK {
		t.Fatalf("followers status = %d, want %d, body=%s", followersResp.Code, stdhttp.StatusOK, followersResp.Body.String())
	}
	followersResult := decodeBody[envelope[social.UserListResponse]](t, followersResp)
	if len(followersResult.Data.List) != 1 || followersResult.Data.List[0].ID != viewerUser.ID {
		t.Fatalf("unexpected followers response: %+v", followersResult.Data.List)
	}

	followingResp := performJSONRequest(t, router, stdhttp.MethodGet, fmt.Sprintf("/api/v1/users/%d/following?page=1&page_size=20", viewerUser.ID), nil, "")
	if followingResp.Code != stdhttp.StatusOK {
		t.Fatalf("following status = %d, want %d, body=%s", followingResp.Code, stdhttp.StatusOK, followingResp.Body.String())
	}
	followingResult := decodeBody[envelope[social.UserListResponse]](t, followingResp)
	if len(followingResult.Data.List) != 1 || followingResult.Data.List[0].ID != 101 {
		t.Fatalf("unexpected following response: %+v", followingResult.Data.List)
	}

	likeResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/videos/1002/likes", nil, viewerToken)
	if likeResp.Code != stdhttp.StatusOK {
		t.Fatalf("like video status = %d, want %d, body=%s", likeResp.Code, stdhttp.StatusOK, likeResp.Body.String())
	}

	likeStatusResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/videos/1002/likes/me", nil, viewerToken)
	if likeStatusResp.Code != stdhttp.StatusOK {
		t.Fatalf("video like status api = %d, want %d, body=%s", likeStatusResp.Code, stdhttp.StatusOK, likeStatusResp.Body.String())
	}
	likeStatus := decodeBody[envelope[video.LikeStatusResponse]](t, likeStatusResp)
	if !likeStatus.Data.Liked {
		t.Fatalf("video like status should be true after like")
	}

	favoriteResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/videos/1002/favorites", nil, viewerToken)
	if favoriteResp.Code != stdhttp.StatusOK {
		t.Fatalf("favorite video status = %d, want %d, body=%s", favoriteResp.Code, stdhttp.StatusOK, favoriteResp.Body.String())
	}

	favoriteStatusResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/videos/1002/favorites/me", nil, viewerToken)
	if favoriteStatusResp.Code != stdhttp.StatusOK {
		t.Fatalf("video favorite status api = %d, want %d, body=%s", favoriteStatusResp.Code, stdhttp.StatusOK, favoriteStatusResp.Body.String())
	}
	favoriteStatus := decodeBody[envelope[video.FavoriteStatusResponse]](t, favoriteStatusResp)
	if !favoriteStatus.Data.Favorited {
		t.Fatalf("video favorite status should be true after favorite")
	}

	detailResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/videos/1002", nil, viewerToken)
	if detailResp.Code != stdhttp.StatusOK {
		t.Fatalf("detail with viewer status = %d, want %d, body=%s", detailResp.Code, stdhttp.StatusOK, detailResp.Body.String())
	}
	detailResult := decodeBody[envelope[video.DetailResponse]](t, detailResp)
	if !detailResult.Data.ViewerState.Liked || !detailResult.Data.ViewerState.Favorited || !detailResult.Data.ViewerState.Followed {
		t.Fatalf("viewer_state should reflect like/favorite/follow, got %+v", detailResult.Data.ViewerState)
	}
	if detailResult.Data.LikeCount != 7 || detailResult.Data.FavoriteCount != 3 {
		t.Fatalf("unexpected detail counters after like/favorite: like=%d favorite=%d", detailResult.Data.LikeCount, detailResult.Data.FavoriteCount)
	}

	createCommentResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/videos/1002/comments", map[string]any{
		"content": "batch d root comment",
	}, viewerToken)
	if createCommentResp.Code != stdhttp.StatusOK {
		t.Fatalf("create comment status = %d, want %d, body=%s", createCommentResp.Code, stdhttp.StatusOK, createCommentResp.Body.String())
	}
	createdComment := decodeBody[envelope[comment.Item]](t, createCommentResp)

	createReplyResp := performJSONRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/comments/%d/replies", createdComment.Data.ID), map[string]any{
		"content": "batch d reply",
	}, viewerToken)
	if createReplyResp.Code != stdhttp.StatusOK {
		t.Fatalf("create reply status = %d, want %d, body=%s", createReplyResp.Code, stdhttp.StatusOK, createReplyResp.Body.String())
	}
	createdReply := decodeBody[envelope[comment.Item]](t, createReplyResp)
	if createdReply.Data.RootID != createdComment.Data.ID {
		t.Fatalf("reply root id = %d, want %d", createdReply.Data.RootID, createdComment.Data.ID)
	}

	commentLikeResp := performJSONRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/comments/%d/likes", createdComment.Data.ID), nil, viewerToken)
	if commentLikeResp.Code != stdhttp.StatusOK {
		t.Fatalf("comment like status = %d, want %d, body=%s", commentLikeResp.Code, stdhttp.StatusOK, commentLikeResp.Body.String())
	}

	commentLikeStatusResp := performJSONRequest(t, router, stdhttp.MethodGet, fmt.Sprintf("/api/v1/comments/%d/likes/me", createdComment.Data.ID), nil, viewerToken)
	if commentLikeStatusResp.Code != stdhttp.StatusOK {
		t.Fatalf("comment like status api = %d, want %d, body=%s", commentLikeStatusResp.Code, stdhttp.StatusOK, commentLikeStatusResp.Body.String())
	}
	commentLikeStatus := decodeBody[envelope[comment.LikeStatusResponse]](t, commentLikeStatusResp)
	if !commentLikeStatus.Data.Liked {
		t.Fatalf("comment like status should be true after like")
	}

	commentListResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/videos/1002/comments?page=1&page_size=20", nil, viewerToken)
	if commentListResp.Code != stdhttp.StatusOK {
		t.Fatalf("comment list status = %d, want %d, body=%s", commentListResp.Code, stdhttp.StatusOK, commentListResp.Body.String())
	}
	commentList := decodeBody[envelope[comment.ListResponse]](t, commentListResp)
	if len(commentList.Data.List) != 1 {
		t.Fatalf("comment list length = %d, want %d", len(commentList.Data.List), 1)
	}
	if !commentList.Data.List[0].ViewerState.Liked || commentList.Data.List[0].ReplyCount != 1 {
		t.Fatalf("unexpected comment list item: %+v", commentList.Data.List[0])
	}

	replyListResp := performJSONRequest(t, router, stdhttp.MethodGet, fmt.Sprintf("/api/v1/comments/%d/replies?page=1&page_size=20", createdComment.Data.ID), nil, viewerToken)
	if replyListResp.Code != stdhttp.StatusOK {
		t.Fatalf("reply list status = %d, want %d, body=%s", replyListResp.Code, stdhttp.StatusOK, replyListResp.Body.String())
	}
	replyList := decodeBody[envelope[comment.ListResponse]](t, replyListResp)
	if len(replyList.Data.List) != 1 || replyList.Data.List[0].ID != createdReply.Data.ID {
		t.Fatalf("unexpected reply list response: %+v", replyList.Data.List)
	}

	commentUnlikeResp := performJSONRequest(t, router, stdhttp.MethodDelete, fmt.Sprintf("/api/v1/comments/%d/likes", createdComment.Data.ID), nil, viewerToken)
	if commentUnlikeResp.Code != stdhttp.StatusOK {
		t.Fatalf("comment unlike status = %d, want %d, body=%s", commentUnlikeResp.Code, stdhttp.StatusOK, commentUnlikeResp.Body.String())
	}
	commentLikeStatusResp = performJSONRequest(t, router, stdhttp.MethodGet, fmt.Sprintf("/api/v1/comments/%d/likes/me", createdComment.Data.ID), nil, viewerToken)
	commentLikeStatus = decodeBody[envelope[comment.LikeStatusResponse]](t, commentLikeStatusResp)
	if commentLikeStatus.Data.Liked {
		t.Fatalf("comment like status should be false after unlike")
	}

	unfavoriteResp := performJSONRequest(t, router, stdhttp.MethodDelete, "/api/v1/videos/1002/favorites", nil, viewerToken)
	if unfavoriteResp.Code != stdhttp.StatusOK {
		t.Fatalf("unfavorite status = %d, want %d, body=%s", unfavoriteResp.Code, stdhttp.StatusOK, unfavoriteResp.Body.String())
	}
	unlikeResp := performJSONRequest(t, router, stdhttp.MethodDelete, "/api/v1/videos/1002/likes", nil, viewerToken)
	if unlikeResp.Code != stdhttp.StatusOK {
		t.Fatalf("unlike video status = %d, want %d, body=%s", unlikeResp.Code, stdhttp.StatusOK, unlikeResp.Body.String())
	}
	unfollowResp := performJSONRequest(t, router, stdhttp.MethodDelete, "/api/v1/users/101/follow", nil, viewerToken)
	if unfollowResp.Code != stdhttp.StatusOK {
		t.Fatalf("unfollow status = %d, want %d, body=%s", unfollowResp.Code, stdhttp.StatusOK, unfollowResp.Body.String())
	}

	followStatusResp = performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/follow-status", nil, viewerToken)
	followStatus = decodeBody[envelope[social.FollowStatusResponse]](t, followStatusResp)
	if followStatus.Data.Followed {
		t.Fatalf("follow status should be false after unfollow")
	}

	detailResp = performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/videos/1002", nil, viewerToken)
	detailResult = decodeBody[envelope[video.DetailResponse]](t, detailResp)
	if detailResult.Data.ViewerState.Liked || detailResult.Data.ViewerState.Favorited || detailResult.Data.ViewerState.Followed {
		t.Fatalf("viewer_state should reset after cleanup, got %+v", detailResult.Data.ViewerState)
	}
	if detailResult.Data.CommentCount != 6 {
		t.Fatalf("detail comment count = %d, want %d", detailResult.Data.CommentCount, 6)
	}
}

func TestBatchEAuthorVideosHotAndFollowingFeeds(t *testing.T) {
	router, db := newTestApp(t)
	seedBatchCData(t, db)
	seedBatchEFeedData(t, db)

	viewerToken, viewerUser := registerAndLoginTestUser(t, router, "viewer_batch_e", "viewer_batch_e@example.com")

	for _, followeeID := range []uint64{101, 102} {
		followResp := performJSONRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/users/%d/follow", followeeID), nil, viewerToken)
		if followResp.Code != stdhttp.StatusOK {
			t.Fatalf("follow %d status = %d, want %d, body=%s", followeeID, followResp.Code, stdhttp.StatusOK, followResp.Body.String())
		}
	}

	authorVideosResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/videos?page=1&page_size=2", nil, "")
	if authorVideosResp.Code != stdhttp.StatusOK {
		t.Fatalf("author videos status = %d, want %d, body=%s", authorVideosResp.Code, stdhttp.StatusOK, authorVideosResp.Body.String())
	}
	authorVideosResult := decodeBody[envelope[video.VideoListResponse]](t, authorVideosResp)
	if authorVideosResult.Data.Pagination.Total != 3 {
		t.Fatalf("author videos total = %d, want %d", authorVideosResult.Data.Pagination.Total, 3)
	}
	if len(authorVideosResult.Data.List) != 2 {
		t.Fatalf("author videos length = %d, want %d", len(authorVideosResult.Data.List), 2)
	}
	if authorVideosResult.Data.List[0].ID != 1003 || authorVideosResult.Data.List[1].ID != 1002 {
		t.Fatalf("unexpected author videos order: %+v", authorVideosResult.Data.List)
	}

	authorVideosNextResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/videos?page=2&page_size=2", nil, "")
	if authorVideosNextResp.Code != stdhttp.StatusOK {
		t.Fatalf("author videos next status = %d, want %d, body=%s", authorVideosNextResp.Code, stdhttp.StatusOK, authorVideosNextResp.Body.String())
	}
	authorVideosNext := decodeBody[envelope[video.VideoListResponse]](t, authorVideosNextResp)
	if len(authorVideosNext.Data.List) != 1 || authorVideosNext.Data.List[0].ID != 1001 {
		t.Fatalf("unexpected author videos next page: %+v", authorVideosNext.Data.List)
	}

	hotResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/feed/hot?limit=2", nil, "")
	if hotResp.Code != stdhttp.StatusOK {
		t.Fatalf("hot feed status = %d, want %d, body=%s", hotResp.Code, stdhttp.StatusOK, hotResp.Body.String())
	}
	hotResult := decodeBody[envelope[video.FeedResponse]](t, hotResp)
	if len(hotResult.Data.Items) != 2 {
		t.Fatalf("hot feed length = %d, want %d", len(hotResult.Data.Items), 2)
	}
	if hotResult.Data.Items[0].ID != 2002 || hotResult.Data.Items[1].ID != 1002 {
		t.Fatalf("unexpected hot feed order: %+v", hotResult.Data.Items)
	}
	if !hotResult.Data.HasMore {
		t.Fatalf("hot feed has_more = false, want true")
	}
	if strings.Count(hotResult.Data.NextCursor, ":") != 2 {
		t.Fatalf("hot feed next_cursor format = %q, want score:unix:id", hotResult.Data.NextCursor)
	}

	hotNextResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/feed/hot?limit=2&cursor="+hotResult.Data.NextCursor, nil, "")
	if hotNextResp.Code != stdhttp.StatusOK {
		t.Fatalf("hot next status = %d, want %d, body=%s", hotNextResp.Code, stdhttp.StatusOK, hotNextResp.Body.String())
	}
	hotNextResult := decodeBody[envelope[video.FeedResponse]](t, hotNextResp)
	if len(hotNextResult.Data.Items) != 2 {
		t.Fatalf("hot next length = %d, want %d", len(hotNextResult.Data.Items), 2)
	}
	if hotNextResult.Data.Items[0].ID != 1003 || hotNextResult.Data.Items[1].ID != 2001 {
		t.Fatalf("unexpected hot next order: %+v", hotNextResult.Data.Items)
	}

	followingResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/feed/following?limit=2", nil, viewerToken)
	if followingResp.Code != stdhttp.StatusOK {
		t.Fatalf("following feed status = %d, want %d, body=%s", followingResp.Code, stdhttp.StatusOK, followingResp.Body.String())
	}
	followingResult := decodeBody[envelope[video.FeedResponse]](t, followingResp)
	if len(followingResult.Data.Items) != 2 {
		t.Fatalf("following feed length = %d, want %d", len(followingResult.Data.Items), 2)
	}
	if followingResult.Data.Items[0].ID != 2002 || followingResult.Data.Items[1].ID != 1003 {
		t.Fatalf("unexpected following feed order: %+v", followingResult.Data.Items)
	}
	if !followingResult.Data.HasMore {
		t.Fatalf("following feed has_more = false, want true")
	}
	if strings.Count(followingResult.Data.NextCursor, ":") != 1 {
		t.Fatalf("following feed next_cursor format = %q, want unix:id", followingResult.Data.NextCursor)
	}

	followingNextResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/feed/following?limit=2&cursor="+followingResult.Data.NextCursor, nil, viewerToken)
	if followingNextResp.Code != stdhttp.StatusOK {
		t.Fatalf("following next status = %d, want %d, body=%s", followingNextResp.Code, stdhttp.StatusOK, followingNextResp.Body.String())
	}
	followingNext := decodeBody[envelope[video.FeedResponse]](t, followingNextResp)
	if len(followingNext.Data.Items) != 2 {
		t.Fatalf("following next length = %d, want %d", len(followingNext.Data.Items), 2)
	}
	if followingNext.Data.Items[0].ID != 1002 || followingNext.Data.Items[1].ID != 2001 {
		t.Fatalf("unexpected following next order: %+v", followingNext.Data.Items)
	}

	areaVideosResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/areas/1/videos?limit=2&sort=latest", nil, "")
	if areaVideosResp.Code != stdhttp.StatusOK {
		t.Fatalf("area videos status = %d, want %d, body=%s", areaVideosResp.Code, stdhttp.StatusOK, areaVideosResp.Body.String())
	}
	areaVideosResult := decodeBody[envelope[video.FeedResponse]](t, areaVideosResp)
	if len(areaVideosResult.Data.Items) != 2 {
		t.Fatalf("area videos length = %d, want %d", len(areaVideosResult.Data.Items), 2)
	}
	if areaVideosResult.Data.Items[0].ID != 1003 || areaVideosResult.Data.Items[1].ID != 1002 {
		t.Fatalf("unexpected area videos order: %+v", areaVideosResult.Data.Items)
	}
	if !areaVideosResult.Data.HasMore {
		t.Fatalf("area videos has_more = false, want true")
	}
	if strings.Count(areaVideosResult.Data.NextCursor, ":") != 1 {
		t.Fatalf("area videos next_cursor format = %q, want unix:id", areaVideosResult.Data.NextCursor)
	}

	areaVideosNextResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/areas/1/videos?limit=2&sort=latest&cursor="+areaVideosResult.Data.NextCursor, nil, "")
	if areaVideosNextResp.Code != stdhttp.StatusOK {
		t.Fatalf("area videos next status = %d, want %d, body=%s", areaVideosNextResp.Code, stdhttp.StatusOK, areaVideosNextResp.Body.String())
	}
	areaVideosNext := decodeBody[envelope[video.FeedResponse]](t, areaVideosNextResp)
	if len(areaVideosNext.Data.Items) != 1 || areaVideosNext.Data.Items[0].ID != 1001 {
		t.Fatalf("unexpected area videos next page: %+v", areaVideosNext.Data.Items)
	}

	invalidAreaSortResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/areas/1/videos?sort=hot", nil, "")
	if invalidAreaSortResp.Code != stdhttp.StatusBadRequest {
		t.Fatalf("invalid area sort status = %d, want %d, body=%s", invalidAreaSortResp.Code, stdhttp.StatusBadRequest, invalidAreaSortResp.Body.String())
	}
	invalidAreaSort := decodeBody[envelope[map[string]any]](t, invalidAreaSortResp)
	if invalidAreaSort.Code != 3115 {
		t.Fatalf("invalid area sort code = %d, want %d", invalidAreaSort.Code, 3115)
	}

	missingAreaResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/areas/999/videos?sort=latest", nil, "")
	if missingAreaResp.Code != stdhttp.StatusNotFound {
		t.Fatalf("missing area status = %d, want %d, body=%s", missingAreaResp.Code, stdhttp.StatusNotFound, missingAreaResp.Body.String())
	}
	missingArea := decodeBody[envelope[map[string]any]](t, missingAreaResp)
	if missingArea.Code != 3116 {
		t.Fatalf("missing area code = %d, want %d", missingArea.Code, 3116)
	}

	if viewerUser.ID == 0 {
		t.Fatalf("viewer user id should not be zero")
	}
}

func TestPhase4SearchHistoryUploadAndAdmin(t *testing.T) {
	router, db := newTestApp(t)
	seedBatchCData(t, db)

	searchVideosResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/search/videos?keyword=Batch%20C&page=1&page_size=2", nil, "")
	if searchVideosResp.Code != stdhttp.StatusOK {
		t.Fatalf("search videos status = %d, want %d, body=%s", searchVideosResp.Code, stdhttp.StatusOK, searchVideosResp.Body.String())
	}
	searchVideosResult := decodeBody[envelope[search.VideoSearchResponse]](t, searchVideosResp)
	if searchVideosResult.Data.Pagination.Total != 3 {
		t.Fatalf("search videos total = %d, want %d", searchVideosResult.Data.Pagination.Total, 3)
	}
	if len(searchVideosResult.Data.List) != 2 || searchVideosResult.Data.List[0].ID != 1003 {
		t.Fatalf("unexpected search videos response: %+v", searchVideosResult.Data.List)
	}

	searchUsersResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/search/users?keyword=tom&page=1&page_size=20", nil, "")
	if searchUsersResp.Code != stdhttp.StatusOK {
		t.Fatalf("search users status = %d, want %d, body=%s", searchUsersResp.Code, stdhttp.StatusOK, searchUsersResp.Body.String())
	}
	searchUsersResult := decodeBody[envelope[search.UserSearchResponse]](t, searchUsersResp)
	if len(searchUsersResult.Data.List) != 1 || searchUsersResult.Data.List[0].ID != 101 {
		t.Fatalf("unexpected search users response: %+v", searchUsersResult.Data.List)
	}

	creatorToken, creatorUser := registerAndLoginTestUser(t, router, "creator_phase4", "creator_phase4@example.com")

	profileResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/profile", nil, "")
	if profileResp.Code != stdhttp.StatusOK {
		t.Fatalf("public profile status = %d, want %d, body=%s", profileResp.Code, stdhttp.StatusOK, profileResp.Body.String())
	}
	profileResult := decodeBody[envelope[account.ProfileResponse]](t, profileResp)
	if profileResult.Data.ID != 101 || profileResult.Data.ViewerState.Followed {
		t.Fatalf("unexpected public profile response: %+v", profileResult.Data)
	}

	reportHistoryResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/histories", map[string]any{
		"video_id":         1002,
		"progress_seconds": 120,
	}, creatorToken)
	if reportHistoryResp.Code != stdhttp.StatusOK {
		t.Fatalf("report history status = %d, want %d, body=%s", reportHistoryResp.Code, stdhttp.StatusOK, reportHistoryResp.Body.String())
	}

	historyResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/histories?page=1&page_size=20", nil, creatorToken)
	if historyResp.Code != stdhttp.StatusOK {
		t.Fatalf("history list status = %d, want %d, body=%s", historyResp.Code, stdhttp.StatusOK, historyResp.Body.String())
	}
	historyResult := decodeBody[envelope[history.ListResponse]](t, historyResp)
	if len(historyResult.Data.List) != 1 || historyResult.Data.List[0].VideoID != 1002 {
		t.Fatalf("unexpected histories response: %+v", historyResult.Data.List)
	}
	if historyResult.Data.List[0].AreaName != "动画" {
		t.Fatalf("history area name = %q, want %q", historyResult.Data.List[0].AreaName, "动画")
	}

	followResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/users/101/follow", nil, creatorToken)
	if followResp.Code != stdhttp.StatusOK {
		t.Fatalf("follow for dashboard status = %d, want %d, body=%s", followResp.Code, stdhttp.StatusOK, followResp.Body.String())
	}

	profileResp = performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/101/profile", nil, creatorToken)
	if profileResp.Code != stdhttp.StatusOK {
		t.Fatalf("authed profile status = %d, want %d, body=%s", profileResp.Code, stdhttp.StatusOK, profileResp.Body.String())
	}
	profileResult = decodeBody[envelope[account.ProfileResponse]](t, profileResp)
	if !profileResult.Data.ViewerState.Followed {
		t.Fatalf("profile viewer follow state should be true after follow")
	}

	favoriteResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/videos/1002/favorites", nil, creatorToken)
	if favoriteResp.Code != stdhttp.StatusOK {
		t.Fatalf("favorite for dashboard status = %d, want %d, body=%s", favoriteResp.Code, stdhttp.StatusOK, favoriteResp.Body.String())
	}

	createVideoResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/videos", map[string]any{
		"area_id":     3,
		"title":       "Phase 4 Creator Video",
		"description": "phase 4 upload pipeline",
	}, creatorToken)
	if createVideoResp.Code != stdhttp.StatusOK {
		t.Fatalf("create video status = %d, want %d, body=%s", createVideoResp.Code, stdhttp.StatusOK, createVideoResp.Body.String())
	}
	createdVideo := decodeBody[envelope[video.CreateVideoResponse]](t, createVideoResp)
	if createdVideo.Data.ReviewStatus != video.ReviewStatusPending {
		t.Fatalf("created video review status = %q, want %q", createdVideo.Data.ReviewStatus, video.ReviewStatusPending)
	}

	uploadSourceResp := performMultipartRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/videos/%d/source", createdVideo.Data.ID), "file", "demo.mp4", []byte("fake mp4 data"), creatorToken)
	if uploadSourceResp.Code != stdhttp.StatusOK {
		t.Fatalf("upload source status = %d, want %d, body=%s", uploadSourceResp.Code, stdhttp.StatusOK, uploadSourceResp.Body.String())
	}
	uploadSourceResult := decodeBody[envelope[video.SourceUploadResponse]](t, uploadSourceResp)
	expectedPlayURL := fmt.Sprintf("/uploads/videos/%d/source.mp4", createdVideo.Data.ID)
	if uploadSourceResult.Data.PlayURL != expectedPlayURL {
		t.Fatalf("play url = %q, want %q", uploadSourceResult.Data.PlayURL, expectedPlayURL)
	}

	uploadCoverResp := performMultipartRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/videos/%d/cover", createdVideo.Data.ID), "file", "cover.jpg", []byte("fake jpg data"), creatorToken)
	if uploadCoverResp.Code != stdhttp.StatusOK {
		t.Fatalf("upload cover status = %d, want %d, body=%s", uploadCoverResp.Code, stdhttp.StatusOK, uploadCoverResp.Body.String())
	}
	uploadCoverResult := decodeBody[envelope[video.CoverUploadResponse]](t, uploadCoverResp)
	expectedCoverURL := fmt.Sprintf("/uploads/videos/%d/cover.jpg", createdVideo.Data.ID)
	if uploadCoverResult.Data.CoverURL != expectedCoverURL {
		t.Fatalf("cover url = %q, want %q", uploadCoverResult.Data.CoverURL, expectedCoverURL)
	}

	creatorPendingResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/creator/videos?review_status=pending&page=1&page_size=20", nil, creatorToken)
	if creatorPendingResp.Code != stdhttp.StatusOK {
		t.Fatalf("creator pending status = %d, want %d, body=%s", creatorPendingResp.Code, stdhttp.StatusOK, creatorPendingResp.Body.String())
	}
	creatorPendingResult := decodeBody[envelope[video.CreatorVideoListResponse]](t, creatorPendingResp)
	if !containsCreatorVideo(creatorPendingResult.Data.List, createdVideo.Data.ID) {
		t.Fatalf("creator pending list should contain created video, got %+v", creatorPendingResult.Data.List)
	}

	rejectVideoResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/videos", map[string]any{
		"area_id":     1,
		"title":       "Phase 4 Rejected Video",
		"description": "to be rejected",
	}, creatorToken)
	if rejectVideoResp.Code != stdhttp.StatusOK {
		t.Fatalf("create reject video status = %d, want %d, body=%s", rejectVideoResp.Code, stdhttp.StatusOK, rejectVideoResp.Body.String())
	}
	rejectVideo := decodeBody[envelope[video.CreateVideoResponse]](t, rejectVideoResp)

	adminToken, adminUser := registerAndLoginTestUser(t, router, "admin_phase4", "admin_phase4@example.com")
	if err := db.Model(&account.User{}).Where("id = ?", adminUser.ID).Update("role", account.RoleAdmin).Error; err != nil {
		t.Fatalf("promote admin role: %v", err)
	}

	nonAdminPendingResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/admin/videos/pending?page=1&page_size=20", nil, creatorToken)
	if nonAdminPendingResp.Code != stdhttp.StatusForbidden {
		t.Fatalf("non-admin pending status = %d, want %d, body=%s", nonAdminPendingResp.Code, stdhttp.StatusForbidden, nonAdminPendingResp.Body.String())
	}

	adminPendingResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/admin/videos/pending?page=1&page_size=20", nil, adminToken)
	if adminPendingResp.Code != stdhttp.StatusOK {
		t.Fatalf("admin pending status = %d, want %d, body=%s", adminPendingResp.Code, stdhttp.StatusOK, adminPendingResp.Body.String())
	}
	adminPendingResult := decodeBody[envelope[admin.VideoListResponse]](t, adminPendingResp)
	if !containsAdminVideo(adminPendingResult.Data.List, createdVideo.Data.ID) || !containsAdminVideo(adminPendingResult.Data.List, rejectVideo.Data.ID) {
		t.Fatalf("admin pending list should contain both new videos, got %+v", adminPendingResult.Data.List)
	}

	approveResp := performJSONRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/admin/videos/%d/approve", createdVideo.Data.ID), nil, adminToken)
	if approveResp.Code != stdhttp.StatusOK {
		t.Fatalf("approve status = %d, want %d, body=%s", approveResp.Code, stdhttp.StatusOK, approveResp.Body.String())
	}
	approveResult := decodeBody[envelope[admin.VideoItem]](t, approveResp)
	if approveResult.Data.ReviewStatus != admin.ReviewStatusApproved || approveResult.Data.PublishedAt == nil {
		t.Fatalf("unexpected approve result: %+v", approveResult.Data)
	}

	rejectResp := performJSONRequest(t, router, stdhttp.MethodPost, fmt.Sprintf("/api/v1/admin/videos/%d/reject", rejectVideo.Data.ID), map[string]any{
		"reason": "cover not compliant",
	}, adminToken)
	if rejectResp.Code != stdhttp.StatusOK {
		t.Fatalf("reject status = %d, want %d, body=%s", rejectResp.Code, stdhttp.StatusOK, rejectResp.Body.String())
	}
	rejectResult := decodeBody[envelope[admin.VideoItem]](t, rejectResp)
	if rejectResult.Data.ReviewStatus != admin.ReviewStatusRejected || rejectResult.Data.ReviewReason != "cover not compliant" {
		t.Fatalf("unexpected reject result: %+v", rejectResult.Data)
	}

	creatorApprovedResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/creator/videos?review_status=approved&page=1&page_size=20", nil, creatorToken)
	if creatorApprovedResp.Code != stdhttp.StatusOK {
		t.Fatalf("creator approved status = %d, want %d, body=%s", creatorApprovedResp.Code, stdhttp.StatusOK, creatorApprovedResp.Body.String())
	}
	creatorApprovedResult := decodeBody[envelope[video.CreatorVideoListResponse]](t, creatorApprovedResp)
	if !containsCreatorVideo(creatorApprovedResult.Data.List, createdVideo.Data.ID) {
		t.Fatalf("creator approved list should contain approved video, got %+v", creatorApprovedResult.Data.List)
	}

	creatorRejectedResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/creator/videos?review_status=rejected&page=1&page_size=20", nil, creatorToken)
	if creatorRejectedResp.Code != stdhttp.StatusOK {
		t.Fatalf("creator rejected status = %d, want %d, body=%s", creatorRejectedResp.Code, stdhttp.StatusOK, creatorRejectedResp.Body.String())
	}
	creatorRejectedResult := decodeBody[envelope[video.CreatorVideoListResponse]](t, creatorRejectedResp)
	if !containsCreatorVideo(creatorRejectedResult.Data.List, rejectVideo.Data.ID) {
		t.Fatalf("creator rejected list should contain rejected video, got %+v", creatorRejectedResult.Data.List)
	}

	adminReviewedResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/admin/videos?review_status=reviewed&page=1&page_size=20", nil, adminToken)
	if adminReviewedResp.Code != stdhttp.StatusOK {
		t.Fatalf("admin reviewed status = %d, want %d, body=%s", adminReviewedResp.Code, stdhttp.StatusOK, adminReviewedResp.Body.String())
	}
	adminReviewedResult := decodeBody[envelope[admin.VideoListResponse]](t, adminReviewedResp)
	if !containsAdminVideo(adminReviewedResult.Data.List, createdVideo.Data.ID) || !containsAdminVideo(adminReviewedResult.Data.List, rejectVideo.Data.ID) {
		t.Fatalf("admin reviewed list should contain approved and rejected videos, got %+v", adminReviewedResult.Data.List)
	}

	todayStatsResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/admin/stats/today", nil, adminToken)
	if todayStatsResp.Code != stdhttp.StatusOK {
		t.Fatalf("today stats status = %d, want %d, body=%s", todayStatsResp.Code, stdhttp.StatusOK, todayStatsResp.Body.String())
	}
	todayStats := decodeBody[envelope[admin.TodayStats]](t, todayStatsResp)
	if todayStats.Data.SubmittedVideoCount < 2 || todayStats.Data.ApprovedVideoCount < 1 || todayStats.Data.PlayCount < 1 {
		t.Fatalf("unexpected today stats: %+v", todayStats.Data)
	}

	areaStatsResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/admin/stats/area", nil, adminToken)
	if areaStatsResp.Code != stdhttp.StatusOK {
		t.Fatalf("area stats status = %d, want %d, body=%s", areaStatsResp.Code, stdhttp.StatusOK, areaStatsResp.Body.String())
	}
	areaStats := decodeBody[envelope[[]admin.AreaStatsItem]](t, areaStatsResp)
	if !containsAreaStats(areaStats.Data, 3, 1, 0) {
		t.Fatalf("expected music area stats to include approved upload, got %+v", areaStats.Data)
	}

	noticesResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/notices?page=1&page_size=20", nil, creatorToken)
	if noticesResp.Code != stdhttp.StatusOK {
		t.Fatalf("creator notices status = %d, want %d, body=%s", noticesResp.Code, stdhttp.StatusOK, noticesResp.Body.String())
	}
	noticesResult := decodeBody[envelope[notice.ListResponse]](t, noticesResp)
	if !containsNoticeTitle(noticesResult.Data.List, "稿件审核通过") || !containsNoticeTitle(noticesResult.Data.List, "稿件审核未通过") {
		t.Fatalf("creator notices should contain approval and rejection notices, got %+v", noticesResult.Data.List)
	}

	noticeItem, ok := findNoticeByTitle(noticesResult.Data.List, "稿件审核通过")
	if !ok {
		t.Fatalf("approve notice not found in %+v", noticesResult.Data.List)
	}

	readNoticeResp := performJSONRequest(t, router, stdhttp.MethodPatch, fmt.Sprintf("/api/v1/notices/%d/read", noticeItem.ID), nil, creatorToken)
	if readNoticeResp.Code != stdhttp.StatusOK {
		t.Fatalf("mark read status = %d, want %d, body=%s", readNoticeResp.Code, stdhttp.StatusOK, readNoticeResp.Body.String())
	}
	readNoticeResult := decodeBody[envelope[notice.Item]](t, readNoticeResp)
	if !readNoticeResult.Data.Read || readNoticeResult.Data.ReadAt == nil {
		t.Fatalf("notice should be marked as read, got %+v", readNoticeResult.Data)
	}

	dashboardResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/me/dashboard", nil, creatorToken)
	if dashboardResp.Code != stdhttp.StatusOK {
		t.Fatalf("dashboard status = %d, want %d, body=%s", dashboardResp.Code, stdhttp.StatusOK, dashboardResp.Body.String())
	}
	dashboardResult := decodeBody[envelope[account.DashboardResponse]](t, dashboardResp)
	if dashboardResult.Data.User.ID != creatorUser.ID || dashboardResult.Data.User.VideoCount != 1 {
		t.Fatalf("unexpected dashboard user response: %+v", dashboardResult.Data.User)
	}
	if !containsDashboardVideo(dashboardResult.Data.RecentVideos, createdVideo.Data.ID) {
		t.Fatalf("dashboard recent videos should contain approved video, got %+v", dashboardResult.Data.RecentVideos)
	}
	if !containsDashboardVideo(dashboardResult.Data.FavoriteVideos, 1002) {
		t.Fatalf("dashboard favorite videos should contain favorited video, got %+v", dashboardResult.Data.FavoriteVideos)
	}
	if !containsDashboardUser(dashboardResult.Data.FollowingUsers, 101) {
		t.Fatalf("dashboard following users should contain followed author, got %+v", dashboardResult.Data.FollowingUsers)
	}

	updateVideoResp := performJSONRequest(t, router, stdhttp.MethodPatch, fmt.Sprintf("/api/v1/videos/%d", createdVideo.Data.ID), map[string]any{
		"area_id":     2,
		"title":       "Phase 4 Creator Video Updated",
		"description": "edited and resubmitted",
	}, creatorToken)
	if updateVideoResp.Code != stdhttp.StatusOK {
		t.Fatalf("update video status = %d, want %d, body=%s", updateVideoResp.Code, stdhttp.StatusOK, updateVideoResp.Body.String())
	}
	updateVideoResult := decodeBody[envelope[video.CreatorVideoItem]](t, updateVideoResp)
	if updateVideoResult.Data.ReviewStatus != video.ReviewStatusPending || updateVideoResult.Data.PublishedAt != nil {
		t.Fatalf("updated video should return to pending, got %+v", updateVideoResult.Data)
	}

	creatorPendingAfterEditResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/creator/videos?review_status=pending&page=1&page_size=20", nil, creatorToken)
	if creatorPendingAfterEditResp.Code != stdhttp.StatusOK {
		t.Fatalf("creator pending after edit status = %d, want %d, body=%s", creatorPendingAfterEditResp.Code, stdhttp.StatusOK, creatorPendingAfterEditResp.Body.String())
	}
	creatorPendingAfterEdit := decodeBody[envelope[video.CreatorVideoListResponse]](t, creatorPendingAfterEditResp)
	if !containsCreatorVideo(creatorPendingAfterEdit.Data.List, createdVideo.Data.ID) {
		t.Fatalf("pending list should contain edited video, got %+v", creatorPendingAfterEdit.Data.List)
	}

	dashboardAfterEditResp := performJSONRequest(t, router, stdhttp.MethodGet, "/api/v1/users/me/dashboard", nil, creatorToken)
	if dashboardAfterEditResp.Code != stdhttp.StatusOK {
		t.Fatalf("dashboard after edit status = %d, want %d, body=%s", dashboardAfterEditResp.Code, stdhttp.StatusOK, dashboardAfterEditResp.Body.String())
	}
	dashboardAfterEdit := decodeBody[envelope[account.DashboardResponse]](t, dashboardAfterEditResp)
	if dashboardAfterEdit.Data.User.VideoCount != 0 {
		t.Fatalf("dashboard user video count after edit = %d, want %d", dashboardAfterEdit.Data.User.VideoCount, 0)
	}
	if containsDashboardVideo(dashboardAfterEdit.Data.RecentVideos, createdVideo.Data.ID) {
		t.Fatalf("dashboard recent videos should not contain edited pending video, got %+v", dashboardAfterEdit.Data.RecentVideos)
	}

	if creatorUser.ID == 0 {
		t.Fatalf("creator user id should not be zero")
	}
}

func newTestApp(t *testing.T) (stdhttp.Handler, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := account.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate users: %v", err)
	}
	areaRepo := area.NewRepository(db)
	if err := areaRepo.AutoMigrate(); err != nil {
		t.Fatalf("auto migrate areas: %v", err)
	}
	if err := areaRepo.SeedDefaults(t.Context()); err != nil {
		t.Fatalf("seed areas: %v", err)
	}
	if err := video.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate videos: %v", err)
	}
	if err := comment.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate comments: %v", err)
	}
	if err := social.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate follows: %v", err)
	}
	if err := history.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate view histories: %v", err)
	}
	if err := notice.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate notices: %v", err)
	}
	if err := admin.NewRepository(db).AutoMigrate(); err != nil {
		t.Fatalf("auto migrate video reviews: %v", err)
	}

	cfg := &appconfig.Config{
		Server: appconfig.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Mode: "test",
		},
		JWT: appconfig.JWTConfig{
			AccessSecret:    "test-access-secret",
			RefreshSecret:   "test-refresh-secret",
			Issuer:          "pilipili-go-test",
			AccessTTLMinute: 120,
			RefreshTTLHour:  24 * 14,
		},
		CORS: appconfig.CORSConfig{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
		Media: appconfig.MediaConfig{
			RootDir:       t.TempDir(),
			PublicBaseURL: "/uploads",
		},
	}

	return apphttp.NewRouter(cfg, db, nil, nil), db
}

func performJSONRequest(t *testing.T, router stdhttp.Handler, method string, path string, body any, bearerToken string) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func performMultipartRequest(t *testing.T, router stdhttp.Handler, method string, path string, fieldName string, filename string, content []byte, bearerToken string) *httptest.ResponseRecorder {
	t.Helper()

	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)
	fileWriter, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := fileWriter.Write(content); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(method, path, &payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func decodeBody[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()

	var value T
	if err := json.Unmarshal(recorder.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode response: %v, body=%s", err, recorder.Body.String())
	}
	return value
}

func seedBatchCData(t *testing.T, db *gorm.DB) {
	t.Helper()

	author := account.User{
		ID:           101,
		Username:     "tom_author",
		Email:        "tom_author@example.com",
		PasswordHash: "hashed",
		Role:         account.RoleUser,
		Status:       account.StatusActive,
		TokenVersion: 1,
	}
	if err := db.Create(&author).Error; err != nil {
		t.Fatalf("create author: %v", err)
	}

	base := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)
	videos := []video.Video{
		{
			ID:              1001,
			AuthorID:        author.ID,
			AreaID:          1,
			Title:           "Batch C Video 1",
			Description:     "video one",
			CoverURL:        "https://example.com/1.jpg",
			PlayURL:         "https://example.com/1.m3u8",
			DurationSeconds: 120,
			Status:          video.StatusVisible,
			ReviewStatus:    video.ReviewStatusApproved,
			PublishedAt:     timePtr(base.Add(-2 * time.Hour)),
			HotScore:        20,
			ViewCount:       10,
			CommentCount:    2,
			LikeCount:       3,
			FavoriteCount:   1,
		},
		{
			ID:              1002,
			AuthorID:        author.ID,
			AreaID:          1,
			Title:           "Batch C Video 2",
			Description:     "video two",
			CoverURL:        "https://example.com/2.jpg",
			PlayURL:         "https://example.com/2.m3u8",
			DurationSeconds: 180,
			Status:          video.StatusVisible,
			ReviewStatus:    video.ReviewStatusApproved,
			PublishedAt:     timePtr(base.Add(-1 * time.Hour)),
			HotScore:        60,
			ViewCount:       20,
			CommentCount:    4,
			LikeCount:       6,
			FavoriteCount:   2,
		},
		{
			ID:              1003,
			AuthorID:        author.ID,
			AreaID:          1,
			Title:           "Batch C Video 3",
			Description:     "video three",
			CoverURL:        "https://example.com/3.jpg",
			PlayURL:         "https://example.com/3.m3u8",
			DurationSeconds: 240,
			Status:          video.StatusVisible,
			ReviewStatus:    video.ReviewStatusApproved,
			PublishedAt:     timePtr(base),
			HotScore:        40,
			ViewCount:       30,
			CommentCount:    6,
			LikeCount:       9,
			FavoriteCount:   3,
		},
	}

	if err := db.Create(&videos).Error; err != nil {
		t.Fatalf("create videos: %v", err)
	}
}

func seedBatchEFeedData(t *testing.T, db *gorm.DB) {
	t.Helper()

	author := account.User{
		ID:           102,
		Username:     "lisa_author",
		Email:        "lisa_author@example.com",
		PasswordHash: "hashed",
		Role:         account.RoleUser,
		Status:       account.StatusActive,
		TokenVersion: 1,
	}
	if err := db.Create(&author).Error; err != nil {
		t.Fatalf("create batch e author: %v", err)
	}

	base := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)
	videos := []video.Video{
		{
			ID:              2001,
			AuthorID:        author.ID,
			AreaID:          2,
			Title:           "Batch E Video 1",
			Description:     "batch e one",
			CoverURL:        "https://example.com/e1.jpg",
			PlayURL:         "https://example.com/e1.m3u8",
			DurationSeconds: 90,
			Status:          video.StatusVisible,
			ReviewStatus:    video.ReviewStatusApproved,
			PublishedAt:     timePtr(base.Add(-90 * time.Minute)),
			HotScore:        35,
			ViewCount:       12,
			CommentCount:    1,
			LikeCount:       2,
			FavoriteCount:   1,
		},
		{
			ID:              2002,
			AuthorID:        author.ID,
			AreaID:          2,
			Title:           "Batch E Video 2",
			Description:     "batch e two",
			CoverURL:        "https://example.com/e2.jpg",
			PlayURL:         "https://example.com/e2.m3u8",
			DurationSeconds: 150,
			Status:          video.StatusVisible,
			ReviewStatus:    video.ReviewStatusApproved,
			PublishedAt:     timePtr(base.Add(30 * time.Minute)),
			HotScore:        90,
			ViewCount:       42,
			CommentCount:    8,
			LikeCount:       10,
			FavoriteCount:   4,
		},
	}

	if err := db.Create(&videos).Error; err != nil {
		t.Fatalf("create batch e videos: %v", err)
	}
}

func timePtr(value time.Time) *time.Time {
	return &value
}

func registerAndLoginTestUser(t *testing.T, router stdhttp.Handler, username string, email string) (string, account.UserResponse) {
	t.Helper()

	registerResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": username,
		"email":    email,
		"password": "12345678",
	}, "")
	if registerResp.Code != stdhttp.StatusCreated {
		t.Fatalf("register test user status = %d, want %d, body=%s", registerResp.Code, stdhttp.StatusCreated, registerResp.Body.String())
	}

	loginResp := performJSONRequest(t, router, stdhttp.MethodPost, "/api/v1/auth/login", map[string]any{
		"username": username,
		"password": "12345678",
	}, "")
	if loginResp.Code != stdhttp.StatusOK {
		t.Fatalf("login test user status = %d, want %d, body=%s", loginResp.Code, stdhttp.StatusOK, loginResp.Body.String())
	}

	loginResult := decodeBody[envelope[loginData]](t, loginResp)
	return loginResult.Data.AccessToken, loginResult.Data.User
}

func containsCreatorVideo(list []video.CreatorVideoItem, videoID uint64) bool {
	for _, item := range list {
		if item.ID == videoID {
			return true
		}
	}
	return false
}

func containsAdminVideo(list []admin.VideoItem, videoID uint64) bool {
	for _, item := range list {
		if item.ID == videoID {
			return true
		}
	}
	return false
}

func containsAreaStats(list []admin.AreaStatsItem, areaID uint64, minApproved int64, minRejected int64) bool {
	for _, item := range list {
		if item.AreaID == areaID && item.ApprovedCount >= minApproved && item.RejectedCount >= minRejected {
			return true
		}
	}
	return false
}

func containsDashboardVideo(list []account.DashboardVideoItem, videoID uint64) bool {
	for _, item := range list {
		if item.ID == videoID {
			return true
		}
	}
	return false
}

func containsDashboardUser(list []account.DashboardUserCard, userID uint64) bool {
	for _, item := range list {
		if item.ID == userID {
			return true
		}
	}
	return false
}

func containsNoticeTitle(list []notice.Item, title string) bool {
	_, ok := findNoticeByTitle(list, title)
	return ok
}

func findNoticeByTitle(list []notice.Item, title string) (notice.Item, bool) {
	for _, item := range list {
		if item.Title == title {
			return item, true
		}
	}
	return notice.Item{}, false
}
