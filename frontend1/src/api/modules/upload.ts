import { apiClient } from '@/api/client';
import type {
  CreateVideoInput,
  CreateVideoResponse,
  CreatorVideoListResponse,
  UploadCoverResponse,
  UploadSourceResponse,
} from '@/shared/types/domain';

export const UploadAPI = {
  createVideo(payload: CreateVideoInput) {
    return apiClient.post<never, CreateVideoResponse>('/videos', payload);
  },
  uploadSource(videoId: string | number, file: File) {
    const formData = new FormData();
    formData.append('file', file);
    return apiClient.post<never, UploadSourceResponse>(`/videos/${videoId}/source`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  uploadCover(videoId: string | number, file: File) {
    const formData = new FormData();
    formData.append('file', file);
    return apiClient.post<never, UploadCoverResponse>(`/videos/${videoId}/cover`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  listCreatorVideos(params?: { review_status?: string; page?: number; page_size?: number }) {
    return apiClient.get<never, CreatorVideoListResponse>('/creator/videos', { params });
  },
};
