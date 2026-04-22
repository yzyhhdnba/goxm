import { apiClient } from '@/api/client';
import type { FeedResponse, VideoDetail } from '@/shared/types/domain';

export const VideoAPI = {
  getDetail(id: string | number) {
    return apiClient.get<never, VideoDetail>(`/videos/${id}`);
  },
  like(id: string | number) {
    return apiClient.post(`/videos/${id}/likes`);
  },
  unlike(id: string | number) {
    return apiClient.delete(`/videos/${id}/likes`);
  },
  favorite(id: string | number) {
    return apiClient.post(`/videos/${id}/favorites`);
  },
  unfavorite(id: string | number) {
    return apiClient.delete(`/videos/${id}/favorites`);
  },
  getRecommend(params?: { cursor?: string; limit?: number }) {
    return apiClient.get<never, FeedResponse>('/feed/recommend', { params });
  },
  follow(userId: string | number) {
    return apiClient.post(`/users/${userId}/follow`);
  },
  unfollow(userId: string | number) {
    return apiClient.delete(`/users/${userId}/follow`);
  },
};
