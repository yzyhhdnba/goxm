import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { AuthAPI } from '@/api/modules/auth';
import { clearSession, getAccessToken, getRefreshToken, saveSession } from '@/shared/utils/auth-storage';
import type { LoginInput, LoginResponse, User } from '@/shared/types/domain';

type AuthStatus = 'idle' | 'loading' | 'authenticated' | 'guest';

type AuthState = {
  accessToken: string | null;
  refreshToken: string | null;
  user: User | null;
  isAuthenticated: boolean;
  status: AuthStatus;
  bootstrapped: boolean;
  error: string | null;
};

const initialState: AuthState = {
  accessToken: getAccessToken(),
  refreshToken: getRefreshToken(),
  user: null,
  isAuthenticated: Boolean(getAccessToken()),
  status: getAccessToken() ? 'loading' : 'guest',
  bootstrapped: false,
  error: null,
};

export const bootstrapAuth = createAsyncThunk('auth/bootstrap', async (_, { rejectWithValue }) => {
  const token = getAccessToken();
  if (!token) {
    return null;
  }

  try {
    const user = await AuthAPI.getCurrentUser();
    return user;
  } catch (error) {
    clearSession();
    return rejectWithValue('登录状态恢复失败');
  }
});

export const login = createAsyncThunk<LoginResponse, LoginInput, { rejectValue: string }>(
  'auth/login',
  async (payload, { rejectWithValue }) => {
    try {
      const response = await AuthAPI.login(payload);
      saveSession({
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
      });
      return response;
    } catch (error) {
      return rejectWithValue('账号或密码错误');
    }
  },
);

export const logout = createAsyncThunk('auth/logout', async () => {
  try {
    await AuthAPI.logout();
  } finally {
    clearSession();
  }
});

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    forceLogout(state, action: PayloadAction<string | undefined>) {
      clearSession();
      state.accessToken = null;
      state.refreshToken = null;
      state.user = null;
      state.isAuthenticated = false;
      state.status = 'guest';
      state.bootstrapped = true;
      state.error = action.payload ?? null;
    },
    clearAuthError(state) {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(bootstrapAuth.pending, (state) => {
        state.status = state.accessToken ? 'loading' : 'guest';
      })
      .addCase(bootstrapAuth.fulfilled, (state, action) => {
        state.user = action.payload;
        state.isAuthenticated = Boolean(action.payload);
        state.status = action.payload ? 'authenticated' : 'guest';
        state.bootstrapped = true;
      })
      .addCase(bootstrapAuth.rejected, (state, action) => {
        state.accessToken = null;
        state.refreshToken = null;
        state.user = null;
        state.isAuthenticated = false;
        state.status = 'guest';
        state.bootstrapped = true;
        state.error = typeof action.payload === 'string' ? action.payload : '登录状态恢复失败';
      })
      .addCase(login.pending, (state) => {
        state.status = 'loading';
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.accessToken = action.payload.access_token;
        state.refreshToken = action.payload.refresh_token;
        state.user = action.payload.user;
        state.isAuthenticated = true;
        state.status = 'authenticated';
        state.bootstrapped = true;
      })
      .addCase(login.rejected, (state, action) => {
        state.isAuthenticated = false;
        state.status = 'guest';
        state.error = action.payload ?? '登录失败';
      })
      .addCase(logout.fulfilled, (state) => {
        state.accessToken = null;
        state.refreshToken = null;
        state.user = null;
        state.isAuthenticated = false;
        state.status = 'guest';
        state.bootstrapped = true;
        state.error = null;
      });
  },
});

export const { forceLogout, clearAuthError } = authSlice.actions;
export default authSlice.reducer;
