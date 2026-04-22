import { apiClient } from '@/api/client';
import type { Area, FeedResponse } from '@/shared/types/domain';

export const FeedAPI = {
  getRecommend(params?: { cursor?: string; limit?: number }) {
    return apiClient.get<never, FeedResponse>('/feed/recommend', { params });
  },
  getAreas() {
    return apiClient.get<never, Area[]>('/areas');
  },
};
