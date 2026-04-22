import { apiClient } from '@/api/client';
import type { AdminVideoItem, AdminVideoListResponse, AreaStatsItem, TodayStats } from '@/shared/types/domain';

export const AdminAPI = {
  listPending(params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, AdminVideoListResponse>('/admin/videos/pending', { params });
  },
  listReviewed(params?: { review_status?: string; page?: number; page_size?: number }) {
    return apiClient.get<never, AdminVideoListResponse>('/admin/videos', { params });
  },
  approve(videoId: string | number) {
    return apiClient.post<never, AdminVideoItem>(`/admin/videos/${videoId}/approve`);
  },
  reject(videoId: string | number, payload: { reason: string }) {
    return apiClient.post<never, AdminVideoItem>(`/admin/videos/${videoId}/reject`, payload);
  },
  getTodayStats() {
    return apiClient.get<never, TodayStats>('/admin/stats/today');
  },
  getAreaStats() {
    return apiClient.get<never, AreaStatsItem[]>('/admin/stats/area');
  },
};
