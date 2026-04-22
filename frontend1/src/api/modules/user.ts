import { apiClient } from '@/api/client';
import type { DashboardResponse, PagedFeedResponse, UserProfile } from '@/shared/types/domain';

export const UserAPI = {
  getProfile(userId: string | number) {
    return apiClient.get<never, UserProfile>(`/users/${userId}/profile`);
  },
  getVideos(userId: string | number, params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, PagedFeedResponse>(`/users/${userId}/videos`, { params });
  },
  getDashboard() {
    return apiClient.get<never, DashboardResponse>('/users/me/dashboard');
  },
  follow(userId: string | number) {
    return apiClient.post(`/users/${userId}/follow`);
  },
  unfollow(userId: string | number) {
    return apiClient.delete(`/users/${userId}/follow`);
  },
};
