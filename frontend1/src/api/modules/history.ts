import { apiClient } from '@/api/client';
import type { HistoryListResponse } from '@/shared/types/domain';

export const HistoryAPI = {
  list(params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, HistoryListResponse>('/histories', { params });
  },
};
