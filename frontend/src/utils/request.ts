import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { ElMessage } from 'element-plus';

// 创建具有统一默认配置的 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: process.env.VUE_APP_API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：统一附带 access_token。
// 对应文档“请求层封装：统一 token、统一错误、统一解包”。
request.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // 自动携带 Token
    const token = localStorage.getItem('access_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: any) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：统一解开后端 Envelope，并集中处理 401/403/500。
// 这样业务组件拿到的就是可直接消费的真实 data。
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data;
    if (code !== 200 && code !== 0 && code !== undefined) {
      ElMessage.error(message || '系统错误');
      return Promise.reject(new Error(message || 'Error'));
    }
    // 返回真实的 data 载荷
    return data !== undefined ? data : response.data;
  },
  (error: any) => {
    // 处理特定的 HTTP 状态码
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.warning('登录已过期，请重新登录');
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          // router.push('/login');
          break;
        case 403:
          ElMessage.error('权限不足');
          break;
        case 404:
          ElMessage.error('请求的资源不存在');
          break;
        case 500:
          ElMessage.error('服务器内部错误');
          break;
        default:
          ElMessage.error(`网络错误: ${error.message}`);
      }
    } else {
      ElMessage.error('网络连接异常');
    }
    return Promise.reject(error);
  }
);

export default request;
