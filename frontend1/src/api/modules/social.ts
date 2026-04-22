import { apiClient } from '@/api/client';
import type { SocialUserListResponse } from '@/shared/types/domain';

export const SocialAPI = {
  getFollowers(userId: string | number, params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, SocialUserListResponse>(`/users/${userId}/followers`, { params });
  },
  getFollowing(userId: string | number, params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, SocialUserListResponse>(`/users/${userId}/following`, { params });
  },
};
