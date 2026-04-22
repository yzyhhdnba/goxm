import { Link } from 'react-router-dom';
import { NoticeBell } from '@/shared/components/layout/NoticeBell';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { logout } from '@/store/slices/authSlice';

export function Header() {
  const dispatch = useAppDispatch();
  const user = useAppSelector((state) => state.auth.user);

  return (
    <header className="sticky top-0 z-50 border-b border-mist bg-white shadow-sm">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 h-16 md:px-6">
        <div className="flex items-center gap-8">
          <Link to="/" className="flex items-center gap-2 group">
            <div className="rounded-md bg-accent px-2 py-1 text-sm font-bold text-white transition group-hover:bg-hoverPink">
              BILIBILI
            </div>
            <span className="text-base font-bold text-ink transition group-hover:text-accent">React Demo</span>
          </Link>

          <nav className="hidden items-center gap-6 md:flex">
            <Link className="text-sm font-medium text-ink transition hover:text-accent" to="/">
              首页
            </Link>
            <Link className="text-sm font-medium text-ink transition hover:text-accent" to="/search?tab=videos">
              搜索
            </Link>
          </nav>
        </div>

        <div className="flex items-center gap-4">
          {user ? (
            <>
              <Link className="flex flex-col items-center justify-center text-xs text-muted transition hover:text-accent" to="/history">
                <span className="text-lg">📺</span>
                <span>历史</span>
              </Link>
              <Link className="flex flex-col items-center justify-center text-xs text-muted transition hover:text-accent" to="/creator">
                <span className="text-lg">⬆️</span>
                <span>投稿</span>
              </Link>
              {user.role === 'admin' ? (
                <Link className="flex flex-col items-center justify-center text-xs text-muted transition hover:text-accent" to="/admin">
                  <span className="text-lg">🛡️</span>
                  <span>审核中心</span>
                </Link>
              ) : null}
              <Link className="flex flex-col items-center justify-center text-xs text-muted transition hover:text-accent" to="/me">
                <span className="text-lg">👤</span>
                <span>动态</span>
              </Link>
              
              <div className="ml-2 flex items-center gap-3">
                <NoticeBell />
                <div className="flex bg-mist/50 rounded-full items-center p-1 px-3 border border-mist hover:bg-mist transition cursor-default">
                  <div className="h-6 w-6 rounded-full bg-sea text-white flex items-center justify-center font-bold text-xs uppercase shadow-sm">
                    {user.username.charAt(0)}
                  </div>
                  <span className="ml-2 text-sm text-ink max-w-[100px] truncate">{user.username}</span>
                </div>
                <button
                  className="rounded border border-mist px-3 py-1.5 text-xs text-muted transition hover:border-gray-400 hover:text-ink"
                  onClick={() => void dispatch(logout())}
                  type="button"
                >
                  退出
                </button>
              </div>
            </>
          ) : (
            <Link className="rounded-md bg-sea px-6 py-1.5 text-sm font-medium text-white transition hover:bg-hoverBlue shadow-sm" to="/login">
              登录注册
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
