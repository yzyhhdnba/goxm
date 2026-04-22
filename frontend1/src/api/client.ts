import axios from 'axios';
import { notifyUnauthorized } from '@/api/interceptors';
import { getAccessToken } from '@/shared/utils/auth-storage';
import { ApiError, type ApiEnvelope } from '@/shared/types/api';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use((config) => {
  const token = getAccessToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

apiClient.interceptors.response.use(
  (response) => {
    const envelope = response.data as ApiEnvelope<unknown>;
    if (typeof envelope?.code === 'number' && envelope.code !== 0 && envelope.code !== 200) {
      throw new ApiError(envelope.message || '请求失败', envelope.code, response.status);
    }

    return envelope?.data ?? response.data;
  },
  (error) => {
    const status = error?.response?.status as number | undefined;
    const code = error?.response?.data?.code as number | undefined;
    const message = error?.response?.data?.message || error?.message || '网络异常';

    if (status === 401) {
      notifyUnauthorized();
    }

    return Promise.reject(new ApiError(message, code, status));
  },
);

export { apiClient };
