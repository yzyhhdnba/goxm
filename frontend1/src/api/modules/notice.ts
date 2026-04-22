import { apiClient } from '@/api/client';
import type { NoticeItem, NoticeListResponse } from '@/shared/types/domain';

export const NoticeAPI = {
  list(params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, NoticeListResponse>('/notices', { params });
  },
  markRead(noticeId: string | number) {
    return apiClient.patch<never, NoticeItem>(`/notices/${noticeId}/read`);
  },
};
