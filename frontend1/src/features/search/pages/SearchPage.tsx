import { FormEvent, useEffect, useMemo, useState } from 'react';
import { Link, useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import { SearchAPI } from '@/api/modules/search';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { UserSearchCard } from '@/shared/components/search/UserSearchCard';
import { VideoCard } from '@/shared/components/video/VideoCard';
import type { SearchUserItem, SearchVideoItem } from '@/shared/types/domain';

type SearchTab = 'videos' | 'users';

const pageSize = 12;

export function SearchPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const routeTab = useMemo<SearchTab>(() => {
    if (location.pathname === '/seuser') {
      return 'users';
    }
    if (location.pathname === '/sevideo') {
      return 'videos';
    }
    const tab = searchParams.get('tab');
    return tab === 'users' ? 'users' : 'videos';
  }, [location.pathname, searchParams]);
  const routeKeyword = useMemo(() => searchParams.get('keyword') || searchParams.get('text') || '', [searchParams]);
  const routePage = useMemo(() => {
    const raw = Number(searchParams.get('page') || '1');
    return Number.isFinite(raw) && raw > 0 ? raw : 1;
  }, [searchParams]);

  const [keywordInput, setKeywordInput] = useState(routeKeyword);
  const [videos, setVideos] = useState<SearchVideoItem[]>([]);
  const [users, setUsers] = useState<SearchUserItem[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setKeywordInput(routeKeyword);
  }, [routeKeyword]);

  useEffect(() => {
    if (!routeKeyword.trim()) {
      setVideos([]);
      setUsers([]);
      setTotal(0);
      setLoading(false);
      setError(null);
      return;
    }

    async function load() {
      setLoading(true);
      setError(null);

      try {
        if (routeTab === 'videos') {
          const response = await SearchAPI.searchVideos({
            keyword: routeKeyword,
            page: routePage,
            page_size: pageSize,
          });
          setVideos(response.list);
          setUsers([]);
          setTotal(response.pagination.total);
          return;
        }

        const response = await SearchAPI.searchUsers({
          keyword: routeKeyword,
          page: routePage,
          page_size: pageSize,
        });
        setUsers(response.list);
        setVideos([]);
        setTotal(response.pagination.total);
      } catch (requestError) {
        setError('搜索失败，请稍后再试');
      } finally {
        setLoading(false);
      }
    }

    void load();
  }, [routeKeyword, routePage, routeTab]);

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const trimmed = keywordInput.trim();
    if (!trimmed) {
      return;
    }

    if (routeTab === 'users') {
      navigate(`/seuser?keyword=${encodeURIComponent(trimmed)}&page=1`);
      return;
    }

    navigate(`/sevideo?keyword=${encodeURIComponent(trimmed)}&page=1`);
  };

  const switchTab = (nextTab: SearchTab) => {
    const trimmed = keywordInput.trim() || routeKeyword.trim();
    if (nextTab === 'users') {
      navigate(trimmed ? `/seuser?keyword=${encodeURIComponent(trimmed)}&page=1` : '/search?tab=users');
      return;
    }

    navigate(trimmed ? `/sevideo?keyword=${encodeURIComponent(trimmed)}&page=1` : '/search?tab=videos');
  };

  const createPageLink = (nextPage: number) => {
    const trimmed = routeKeyword.trim();
    if (!trimmed) {
      return routeTab === 'users' ? '/search?tab=users' : '/search?tab=videos';
    }
    return routeTab === 'users'
      ? `/seuser?keyword=${encodeURIComponent(trimmed)}&page=${nextPage}`
      : `/sevideo?keyword=${encodeURIComponent(trimmed)}&page=${nextPage}`;
  };

  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const hasPrev = routePage > 1;
  const hasNext = routePage < totalPages;
  const fromSearch = `${location.pathname}${location.search}`;

  return (
    <AppShell>
      <section className="mb-8 rounded-md bg-white p-6 shadow-sm">
        <div className="mb-6 space-y-2">
          <p className="text-sm uppercase tracking-[0.3em] text-sea">Search</p>
          <h1 className="text-3xl font-semibold text-ink">搜索视频或用户</h1>
          <p className="text-sm text-slate-500">当前按后端真实实现接入：匿名可搜索，分页使用 `page + page_size`。</p>
        </div>

        <form className="flex flex-col gap-4 md:flex-row" onSubmit={handleSubmit}>
          <div className="flex rounded-full bg-paper p-1">
            <button
              className={`rounded-full px-4 py-2 text-sm font-medium transition ${routeTab === 'videos' ? 'bg-ink text-white' : 'text-slate-600'}`}
              onClick={() => switchTab('videos')}
              type="button"
            >
              视频
            </button>
            <button
              className={`rounded-full px-4 py-2 text-sm font-medium transition ${routeTab === 'users' ? 'bg-ink text-white' : 'text-slate-600'}`}
              onClick={() => switchTab('users')}
              type="button"
            >
              用户
            </button>
          </div>

          <input
            className="flex-1 rounded-full border border-slate-200 bg-slate-50 px-5 py-3 outline-none transition focus:border-accent"
            onChange={(event) => setKeywordInput(event.target.value)}
            placeholder={routeTab === 'videos' ? '输入视频关键词' : '输入用户名关键词'}
            value={keywordInput}
          />

          <button className="rounded-full bg-accent px-5 py-3 text-sm font-medium text-white" type="submit">
            开始搜索
          </button>
        </form>
      </section>

      {!routeKeyword.trim() ? <EmptyState description="先输入关键词，再决定搜视频还是搜用户。" title="还没有开始搜索" /> : null}

      {routeKeyword.trim() ? (
        <section className="space-y-6">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">
                {routeTab === 'videos' ? 'Video Results' : 'User Results'}
              </p>
              <h2 className="text-2xl font-semibold text-ink">
                “{routeKeyword}” 共找到 {total} 条结果
              </h2>
            </div>
            <div className="text-sm text-slate-500">第 {routePage} / {totalPages} 页</div>
          </div>

          {loading ? (
            <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
              {Array.from({ length: 6 }).map((_, index) => (
                <LoadingBlock key={index} lines={5} />
              ))}
            </div>
          ) : null}

          {!loading && error ? <EmptyState title={error} /> : null}

          {!loading && !error && routeTab === 'videos' && videos.length === 0 ? (
            <EmptyState title="没有找到相关视频" />
          ) : null}

          {!loading && !error && routeTab === 'users' && users.length === 0 ? (
            <EmptyState title="没有找到相关用户" />
          ) : null}

          {!loading && !error && routeTab === 'videos' && videos.length > 0 ? (
            <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
              {videos.map((item) => (
                <VideoCard
                  item={item}
                  key={item.id}
                  to={`/video/${item.id}?from=${encodeURIComponent(fromSearch)}`}
                />
              ))}
            </div>
          ) : null}

          {!loading && !error && routeTab === 'users' && users.length > 0 ? (
            <div className="space-y-4">
              {users.map((item) => (
                <UserSearchCard
                  item={item}
                  key={item.id}
                  to={`/users/${item.id}?from=${encodeURIComponent(fromSearch)}`}
                />
              ))}
            </div>
          ) : null}

          {!loading && !error && total > 0 ? (
            <div className="flex items-center justify-between rounded-md bg-white px-5 py-4 shadow-sm">
              {hasPrev ? (
                <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to={createPageLink(routePage - 1)}>
                  上一页
                </Link>
              ) : (
                <span className="rounded-full border border-slate-100 px-4 py-2 text-sm text-slate-300">上一页</span>
              )}
              {hasNext ? (
                <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to={createPageLink(routePage + 1)}>
                  下一页
                </Link>
              ) : (
                <span className="rounded-full border border-slate-100 px-4 py-2 text-sm text-slate-300">下一页</span>
              )}
            </div>
          ) : null}
        </section>
      ) : null}
    </AppShell>
  );
}
