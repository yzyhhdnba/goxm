import { FormEvent, useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { AppShell } from '@/shared/components/layout/AppShell';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { clearAuthError, login } from '@/store/slices/authSlice';

export function LoginPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const auth = useAppSelector((state) => state.auth);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    return () => {
      dispatch(clearAuthError());
    };
  }, [dispatch]);

  useEffect(() => {
    if (auth.isAuthenticated) {
      const target = (location.state as { from?: string } | null)?.from || '/';
      navigate(target, { replace: true });
    }
  }, [auth.isAuthenticated, location.state, navigate]);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await dispatch(login({ username, password }));
  };

  return (
    <AppShell>
      <div className="mx-auto max-w-md rounded-md bg-white p-8 shadow-sm">
        <div className="mb-8 space-y-2">
          <p className="text-sm uppercase tracking-[0.3em] text-sea">Auth</p>
          <h1 className="text-3xl font-semibold text-ink">登录到 React 版前端</h1>
          <p className="text-sm text-slate-500">当前只先打通主链路，不提前铺注册、刷新、后台等扩展页面。</p>
        </div>

        <form className="space-y-5" onSubmit={handleSubmit}>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-slate-600">用户名或邮箱</span>
            <input
              className="w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none transition focus:border-accent"
              onChange={(event) => setUsername(event.target.value)}
              placeholder="alice"
              value={username}
            />
          </label>

          <label className="block space-y-2">
            <span className="text-sm font-medium text-slate-600">密码</span>
            <input
              className="w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none transition focus:border-accent"
              onChange={(event) => setPassword(event.target.value)}
              placeholder="请输入密码"
              type="password"
              value={password}
            />
          </label>

          {auth.error ? <p className="rounded-md bg-red-50 px-4 py-3 text-sm text-red-600">{auth.error}</p> : null}

          <button
            className="w-full rounded-md bg-ink px-4 py-3 font-medium text-white transition hover:bg-sea disabled:cursor-not-allowed disabled:opacity-60"
            disabled={auth.status === 'loading'}
            type="submit"
          >
            {auth.status === 'loading' ? '登录中...' : '登录'}
          </button>
        </form>
      </div>
    </AppShell>
  );
}
