import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { UserAPI } from '@/api/modules/user';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { UserSearchCard } from '@/shared/components/search/UserSearchCard';
import { VideoCard } from '@/shared/components/video/VideoCard';
import type { DashboardResponse } from '@/shared/types/domain';

export function DashboardPage() {
  const [dashboard, setDashboard] = useState<DashboardResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError(null);
      try {
        const response = await UserAPI.getDashboard();
        setDashboard(response);
      } catch (requestError) {
        setError('个人仪表盘加载失败');
      } finally {
        setLoading(false);
      }
    }

    void load();
  }, []);

  return (
    <AppShell>
      <section className="mb-6 rounded-md bg-white p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-sea">Dashboard</p>
            <h1 className="text-3xl font-semibold text-ink">个人中心</h1>
          </div>
          <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to="/history">
            查看历史记录
          </Link>
        </div>
      </section>

      {loading ? (
        <div className="space-y-6">
          <LoadingBlock lines={6} />
          <LoadingBlock lines={6} />
        </div>
      ) : null}

      {!loading && error ? <EmptyState title={error} /> : null}

      {!loading && dashboard ? (
        <div className="space-y-6">
          <section className="grid gap-5 lg:grid-cols-[0.8fr_1.2fr]">
            <div className="rounded-md bg-white p-6 shadow-sm">
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Profile</p>
              <h2 className="mt-3 text-2xl font-semibold text-ink">{dashboard.user.username}</h2>
              <p className="mt-3 text-sm leading-7 text-slate-600">{dashboard.user.bio || '这个账号还没有填写简介。'}</p>
              <div className="mt-5 grid grid-cols-2 gap-3 text-sm text-slate-500">
                <div className="rounded-md bg-paper p-4">{dashboard.user.follower_count} 粉丝</div>
                <div className="rounded-md bg-paper p-4">{dashboard.user.following_count} 关注</div>
                <div className="rounded-md bg-paper p-4">{dashboard.user.video_count} 视频</div>
                <div className="rounded-md bg-paper p-4">{dashboard.stats.total_view_count} 总播放</div>
              </div>
            </div>

            <div className="rounded-md bg-white p-6 shadow-sm">
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Quick Access</p>
              <div className="mt-4 grid gap-3 md:grid-cols-3">
                <Link className="rounded-md bg-paper px-4 py-4 text-sm text-slate-700" to={`/users/${dashboard.user.id}?tab=videos&page=1`}>
                  查看我的主页
                </Link>
                <Link className="rounded-md bg-paper px-4 py-4 text-sm text-slate-700" to="/history">
                  查看历史记录
                </Link>
                <Link className="rounded-md bg-paper px-4 py-4 text-sm text-slate-700" to={`/users/${dashboard.user.id}?tab=following&page=1`}>
                  查看关注列表
                </Link>
              </div>
            </div>
          </section>

          <section className="space-y-4">
            <div>
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Recent Videos</p>
              <h2 className="text-2xl font-semibold text-ink">最近公开视频</h2>
            </div>
            {dashboard.recent_videos.length === 0 ? (
              <EmptyState title="暂时还没有公开视频" />
            ) : (
              <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
                {dashboard.recent_videos.map((item) => (
                  <VideoCard item={item} key={item.id} to={`/video/${item.id}?from=${encodeURIComponent('/me')}`} />
                ))}
              </div>
            )}
          </section>

          <section className="space-y-4">
            <div>
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Favorites</p>
              <h2 className="text-2xl font-semibold text-ink">最近收藏</h2>
            </div>
            {dashboard.favorite_videos.length === 0 ? (
              <EmptyState title="暂时还没有收藏视频" />
            ) : (
              <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
                {dashboard.favorite_videos.map((item) => (
                  <VideoCard item={item} key={item.id} to={`/video/${item.id}?from=${encodeURIComponent('/me')}`} />
                ))}
              </div>
            )}
          </section>

          <section className="space-y-4">
            <div>
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Following</p>
              <h2 className="text-2xl font-semibold text-ink">关注的人</h2>
            </div>
            {dashboard.following_users.length === 0 ? (
              <EmptyState title="暂时还没有关注任何人" />
            ) : (
              <div className="space-y-4">
                {dashboard.following_users.map((item) => (
                  <UserSearchCard item={item} key={item.id} to={`/users/${item.id}?from=${encodeURIComponent('/me')}`} />
                ))}
              </div>
            )}
          </section>
        </div>
      ) : null}
    </AppShell>
  );
}
