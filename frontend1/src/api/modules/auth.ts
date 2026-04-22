import { apiClient } from '@/api/client';
import type { LoginInput, LoginResponse, User } from '@/shared/types/domain';

export const AuthAPI = {
  login(payload: LoginInput) {
    return apiClient.post<never, LoginResponse>('/auth/login', payload);
  },
  getCurrentUser() {
    return apiClient.get<never, User>('/users/me');
  },
  logout() {
    return apiClient.post('/auth/logout');
  },
};
