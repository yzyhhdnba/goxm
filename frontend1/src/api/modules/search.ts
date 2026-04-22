import { apiClient } from '@/api/client';
import type { SearchUserResponse, SearchVideoResponse } from '@/shared/types/domain';

type SearchParams = {
  keyword: string;
  page?: number;
  page_size?: number;
};

export const SearchAPI = {
  searchVideos(params: SearchParams) {
    return apiClient.get<never, SearchVideoResponse>('/search/videos', { params });
  },
  searchUsers(params: SearchParams) {
    return apiClient.get<never, SearchUserResponse>('/search/users', { params });
  },
};
