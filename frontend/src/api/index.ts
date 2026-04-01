import request from '../utils/request';

/**
 * 认证相关接口 (Batch B)
 */
export const AuthAPI = {
  login: (data: any) => request.post('/auth/login', data),
  register: (data: any) => request.post('/auth/register', data),
  logout: () => request.post('/auth/logout'),
  getCurrentUser: () => request.get('/users/me'),
  checkUsername: (username: string) => request.get('/auth/check-username', { params: { username } }),
  checkEmail: (email: string) => request.get('/auth/check-email', { params: { email } }),
};

export const UserAPI = {
  getProfile: (id: string | number) => request.get(`/users/${id}/profile`),
  getDashboard: () => request.get('/users/me/dashboard'),
};

/**
 * 视频流与详情接口 (Batch C)
 */
export const VideoAPI = {
  getRecommend: (params?: any) => request.get('/feed/recommend', { params }),
  getHot: (params?: any) => request.get('/feed/hot', { params }),
  getVideoDetail: (id: string | number) => request.get(`/videos/${id}`),
  getVideoListByArea: (areaId: string | number, params?: any) => request.get(`/areas/${areaId}/videos`, { params }),
  getAreas: () => request.get('/areas'),
};

/**
 * 互动相关接口 (Batch D)
 */
export const InteractAPI = {
  // 视频点赞与收藏
  likeVideo: (id: string | number) => request.post(`/videos/${id}/likes`),
  unlikeVideo: (id: string | number) => request.delete(`/videos/${id}/likes`),
  collectVideo: (id: string | number) => request.post(`/videos/${id}/favorites`),
  uncollectVideo: (id: string | number) => request.delete(`/videos/${id}/favorites`),
  
  // 评论与回复
  getComments: (videoId: string | number, params?: any) => request.get(`/videos/${videoId}/comments`, { params }),
  postComment: (videoId: string | number, data: any) => request.post(`/videos/${videoId}/comments`, data),
  likeComment: (commentId: string | number) => request.post(`/comments/${commentId}/likes`),
  unlikeComment: (commentId: string | number) => request.delete(`/comments/${commentId}/likes`),
  
  getReplies: (commentId: string | number, params?: any) => request.get(`/comments/${commentId}/replies`, { params }),
  postReply: (commentId: string | number, data: any) => request.post(`/comments/${commentId}/replies`, data),

  // 用户与关注
  followUser: (userId: string | number) => request.post(`/users/${userId}/follow`),
  unfollowUser: (userId: string | number) => request.delete(`/users/${userId}/follow`),
  getFollowStatus: (userId: string | number) => request.get(`/users/${userId}/follow-status`),
};

export const SearchAPI = {
  searchVideos: (params?: any) => request.get('/search/videos', { params }),
  searchUsers: (params?: any) => request.get('/search/users', { params }),
};

export const HistoryAPI = {
  list: (params?: any) => request.get('/histories', { params }),
  report: (data: any) => request.post('/histories', data),
};

export const NoticeAPI = {
  list: (params?: any) => request.get('/notices', { params }),
  markRead: (noticeId: string | number) => request.patch(`/notices/${noticeId}/read`),
};

export const UploadAPI = {
  createVideoMetadata: (data: any) => request.post('/videos', data),
  uploadSource: (videoId: string | number, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return request.post(`/videos/${videoId}/source`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
  uploadCover: (videoId: string | number, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return request.post(`/videos/${videoId}/cover`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};

export const CreatorAPI = {
  listVideos: (params?: any) => request.get('/creator/videos', { params }),
  updateVideo: (videoId: string | number, data: any) => request.patch(`/videos/${videoId}`, data),
};

export const AdminAPI = {
  listVideos: (params?: any) => request.get('/admin/videos', { params }),
  listPendingVideos: (params?: any) => request.get('/admin/videos/pending', { params }),
  approveVideo: (videoId: string | number) => request.post(`/admin/videos/${videoId}/approve`),
  rejectVideo: (videoId: string | number, data: any) => request.post(`/admin/videos/${videoId}/reject`, data),
  getTodayStats: () => request.get('/admin/stats/today'),
  getAreaStats: () => request.get('/admin/stats/area'),
};
