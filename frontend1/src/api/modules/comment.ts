import { apiClient } from '@/api/client';
import type { CommentItem, CommentListResponse } from '@/shared/types/domain';

export const CommentAPI = {
  getList(videoId: string | number, params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, CommentListResponse>(`/videos/${videoId}/comments`, { params });
  },
  create(videoId: string | number, payload: { content: string }) {
    return apiClient.post<never, CommentItem>(`/videos/${videoId}/comments`, payload);
  },
  getReplies(commentId: string | number, params?: { page?: number; page_size?: number }) {
    return apiClient.get<never, CommentListResponse>(`/comments/${commentId}/replies`, { params });
  },
  createReply(commentId: string | number, payload: { content: string }) {
    return apiClient.post<never, CommentItem>(`/comments/${commentId}/replies`, payload);
  },
  like(commentId: string | number) {
    return apiClient.post(`/comments/${commentId}/likes`);
  },
  unlike(commentId: string | number) {
    return apiClient.delete(`/comments/${commentId}/likes`);
  },
};
