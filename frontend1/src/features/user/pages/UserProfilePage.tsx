import { useEffect, useMemo, useState } from 'react';
import { Link, useParams, useSearchParams } from 'react-router-dom';
import { SocialAPI } from '@/api/modules/social';
import { UserAPI } from '@/api/modules/user';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { UserSearchCard } from '@/shared/components/search/UserSearchCard';
import { VideoCard } from '@/shared/components/video/VideoCard';
import { useAppSelector } from '@/store/hooks';
import type { FeedItem, SocialUserCard, UserProfile } from '@/shared/types/domain';

const videoPageSize = 8;
const peoplePageSize = 12;

type ProfileTab = 'videos' | 'followers' | 'following';

export function UserProfilePage() {
  const { id } = useParams();
  const [searchParams] = useSearchParams();
  const authUser = useAppSelector((state) => state.auth.user);
  const isAuthenticated = useAppSelector((state) => state.auth.isAuthenticated);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [videos, setVideos] = useState<FeedItem[]>([]);
  const [people, setPeople] = useState<SocialUserCard[]>([]);
  const [total, setTotal] = useState(0);
  const [profileLoading, setProfileLoading] = useState(true);
  const [panelLoading, setPanelLoading] = useState(true);
  const [followPending, setFollowPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const activeTab = useMemo<ProfileTab>(() => {
    const tab = searchParams.get('tab');
    if (tab === 'followers' || tab === 'following') {
      return tab;
    }
    return 'videos';
  }, [searchParams]);

  const page = useMemo(() => {
    const raw = Number(searchParams.get('page') || '1');
    return Number.isFinite(raw) && raw > 0 ? raw : 1;
  }, [searchParams]);

  const from = searchParams.get('from');
  const backTo = from || '/search?tab=users';
  const isSelf = profile && authUser ? profile.id === authUser.id : false;
  const panelPageSize = activeTab === 'videos' ? videoPageSize : peoplePageSize;
  const totalPages = Math.max(1, Math.ceil(total / panelPageSize));
  const currentPath = useMemo(() => `/users/${id}?${searchParams.toString()}`, [id, searchParams]);

  useEffect(() => {
    if (!id) {
      setError('缺少用户 ID');
      setProfileLoading(false);
      setPanelLoading(false);
      return;
    }
    const userId = id;

    async function loadProfile() {
      setProfileLoading(true);
      setError(null);

      try {
        const profileData = await UserAPI.getProfile(userId);
        setProfile(profileData);
      } catch (requestError) {
        setError('用户主页加载失败');
      } finally {
        setProfileLoading(false);
      }
    }

    void loadProfile();
  }, [id]);

  useEffect(() => {
    if (!id) {
      return;
    }
    const userId = id;

    async function loadPanel() {
      setPanelLoading(true);
      setError(null);

      try {
        if (activeTab === 'videos') {
          const videoData = await UserAPI.getVideos(userId, { page, page_size: videoPageSize });
          setVideos(videoData.list);
          setPeople([]);
          setTotal(videoData.pagination.total);
          return;
        }

        const relationData =
          activeTab === 'followers'
            ? await SocialAPI.getFollowers(userId, { page, page_size: peoplePageSize })
            : await SocialAPI.getFollowing(userId, { page, page_size: peoplePageSize });

        setPeople(relationData.list);
        setVideos([]);
        setTotal(relationData.pagination.total);
      } catch (requestError) {
        setError(activeTab === 'videos' ? '视频列表加载失败' : '关注列表加载失败');
      } finally {
        setPanelLoading(false);
      }
    }

    void loadPanel();
  }, [activeTab, id, page]);

  const buildTabLink = (tab: ProfileTab, nextPage = 1) => {
    const params = new URLSearchParams();
    params.set('tab', tab);
    params.set('page', String(nextPage));
    if (from) {
      params.set('from', from);
    }
    return `/users/${id}?${params.toString()}`;
  };

  const handleFollow = async () => {
    if (!profile || !isAuthenticated || isSelf) {
      return;
    }

    setFollowPending(true);
    try {
      if (profile.viewer_state.followed) {
        await UserAPI.unfollow(profile.id);
        setProfile({
          ...profile,
          follower_count: Math.max(0, profile.follower_count - 1),
          viewer_state: { followed: false },
        });
        return;
      }

      await UserAPI.follow(profile.id);
      setProfile({
        ...profile,
        follower_count: profile.follower_count + 1,
        viewer_state: { followed: true },
      });
    } catch (requestError) {
      setError('关注操作失败');
    } finally {
      setFollowPending(false);
    }
  };

  const panelTitle =
    activeTab === 'videos' ? '作者公开视频' : activeTab === 'followers' ? '粉丝列表' : '关注列表';

  return (
    <AppShell>
      <div className="mb-6">
        <Link className="text-sm text-slate-500 transition hover:text-ink" to={backTo}>
          ← 返回搜索结果
        </Link>
      </div>

      {profileLoading ? (
        <div className="space-y-6">
          <LoadingBlock lines={6} />
          <LoadingBlock lines={5} />
        </div>
      ) : null}

      {!profileLoading && error && !profile ? <EmptyState title={error} /> : null}

      {!profileLoading && profile ? (
        <div className="space-y-6">
          {error ? <div className="rounded-md bg-amber-50 px-4 py-3 text-sm text-amber-700">{error}</div> : null}

          <section className="rounded-md bg-white p-6 shadow-sm">
            <div className="flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
              <div className="flex items-start gap-5">
                <div className="flex h-24 w-24 items-center justify-center rounded-full bg-mist text-3xl font-semibold text-sea">
                  {profile.username.slice(0, 1).toUpperCase()}
                </div>
                <div className="space-y-3">
                  <div>
                    <p className="text-sm uppercase tracking-[0.3em] text-sea">User Profile</p>
                    <h1 className="text-3xl font-semibold text-ink">{profile.username}</h1>
                  </div>
                  <p className="max-w-2xl text-sm leading-7 text-slate-600">
                    {profile.bio || '这个用户还没有填写简介。'}
                  </p>
                  <div className="flex flex-wrap gap-3">
                    <Link
                      className={`rounded-full px-4 py-2 text-sm transition ${
                        activeTab === 'followers' ? 'bg-ink text-white' : 'bg-paper text-slate-600'
                      }`}
                      to={buildTabLink('followers')}
                    >
                      {profile.follower_count} 粉丝
                    </Link>
                    <Link
                      className={`rounded-full px-4 py-2 text-sm transition ${
                        activeTab === 'following' ? 'bg-ink text-white' : 'bg-paper text-slate-600'
                      }`}
                      to={buildTabLink('following')}
                    >
                      {profile.following_count} 关注
                    </Link>
                    <Link
                      className={`rounded-full px-4 py-2 text-sm transition ${
                        activeTab === 'videos' ? 'bg-ink text-white' : 'bg-paper text-slate-600'
                      }`}
                      to={buildTabLink('videos')}
                    >
                      {profile.video_count} 视频
                    </Link>
                  </div>
                </div>
              </div>

              {!isSelf ? (
                <button
                  className={`rounded-full px-5 py-3 text-sm font-medium transition ${
                    profile.viewer_state.followed ? 'bg-ink text-white' : 'bg-accent text-white'
                  } disabled:cursor-not-allowed disabled:opacity-60`}
                  disabled={!isAuthenticated || followPending}
                  onClick={() => void handleFollow()}
                  type="button"
                >
                  {!isAuthenticated ? '登录后可关注' : followPending ? '处理中...' : profile.viewer_state.followed ? '已关注' : '关注'}
                </button>
              ) : (
                <span className="rounded-full bg-paper px-4 py-3 text-sm text-slate-500">这是你的主页</span>
              )}
            </div>
          </section>

          <section className="space-y-5">
            <div className="flex flex-wrap items-end justify-between gap-3">
              <div>
                <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Panel</p>
                <h2 className="text-2xl font-semibold text-ink">{panelTitle}</h2>
              </div>
              <div className="text-sm text-slate-500">
                第 {page} / {totalPages} 页 · 共 {total} 条
              </div>
            </div>

            {panelLoading ? (
              <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
                {Array.from({ length: activeTab === 'videos' ? 4 : 3 }).map((_, index) => (
                  <LoadingBlock key={index} lines={5} />
                ))}
              </div>
            ) : null}

            {!panelLoading && activeTab === 'videos' && videos.length === 0 ? (
              <EmptyState title="这个用户还没有公开视频" />
            ) : null}

            {!panelLoading && activeTab !== 'videos' && people.length === 0 ? (
              <EmptyState title={activeTab === 'followers' ? '还没有粉丝' : '暂时没有关注任何人'} />
            ) : null}

            {!panelLoading && activeTab === 'videos' && videos.length > 0 ? (
              <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
                {videos.map((item) => (
                  <VideoCard item={item} key={item.id} to={`/video/${item.id}?from=${encodeURIComponent(currentPath)}`} />
                ))}
              </div>
            ) : null}

            {!panelLoading && activeTab !== 'videos' && people.length > 0 ? (
              <div className="space-y-4">
                {people.map((item) => (
                  <UserSearchCard
                    item={item}
                    key={item.id}
                    meta={
                      <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">
                        {activeTab === 'followers' ? '粉丝' : '关注'}
                      </span>
                    }
                    to={`/users/${item.id}?from=${encodeURIComponent(currentPath)}`}
                  />
                ))}
              </div>
            ) : null}

            {total > 0 ? (
              <div className="flex items-center justify-between rounded-md bg-white px-5 py-4 shadow-sm">
                {page > 1 ? (
                  <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to={buildTabLink(activeTab, page - 1)}>
                    上一页
                  </Link>
                ) : (
                  <span className="rounded-full border border-slate-100 px-4 py-2 text-sm text-slate-300">上一页</span>
                )}
                {page < totalPages ? (
                  <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to={buildTabLink(activeTab, page + 1)}>
                    下一页
                  </Link>
                ) : (
                  <span className="rounded-full border border-slate-100 px-4 py-2 text-sm text-slate-300">下一页</span>
                )}
              </div>
            ) : null}
          </section>
        </div>
      ) : null}
    </AppShell>
  );
}
