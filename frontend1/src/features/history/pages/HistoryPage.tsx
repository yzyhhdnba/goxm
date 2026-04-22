import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { HistoryAPI } from '@/api/modules/history';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { formatDate, formatCount, resolveMediaUrl } from '@/shared/utils/format';
import type { HistoryItem } from '@/shared/types/domain';

const pageSize = 20;

export function HistoryPage() {
  const [searchParams] = useSearchParams();
  const [items, setItems] = useState<HistoryItem[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const page = Number(searchParams.get('page') || '1') > 0 ? Number(searchParams.get('page') || '1') : 1;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError(null);
      try {
        const response = await HistoryAPI.list({ page, page_size: pageSize });
        setItems(response.list);
        setTotal(response.pagination.total);
      } catch (requestError) {
        setError('历史记录加载失败');
      } finally {
        setLoading(false);
      }
    }

    void load();
  }, [page]);

  const buildPageLink = (nextPage: number) => `/history?page=${nextPage}`;

  return (
    <AppShell>
      <section className="mb-6 rounded-md bg-white p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-sea">History</p>
            <h1 className="text-3xl font-semibold text-ink">播放历史</h1>
          </div>
          <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to="/me">
            返回个人中心
          </Link>
        </div>
      </section>

      {loading ? (
        <div className="space-y-4">
          {Array.from({ length: 4 }).map((_, index) => (
            <LoadingBlock key={index} lines={5} />
          ))}
        </div>
      ) : null}

      {!loading && error ? <EmptyState title={error} /> : null}

      {!loading && !error && items.length === 0 ? <EmptyState title="还没有播放历史" /> : null}

      {!loading && !error && items.length > 0 ? (
        <div className="space-y-4">
          {items.map((item) => (
            <article className="flex flex-col gap-4 rounded-md bg-white p-5 shadow-sm md:flex-row" key={`${item.video_id}-${item.watched_at}`}>
              <Link
                className="block overflow-hidden rounded-md md:w-72"
                to={`/video/${item.video_id}?from=${encodeURIComponent(`/history?page=${page}`)}`}
              >
                <img alt={item.video_title} className="aspect-video h-full w-full object-cover" src={resolveMediaUrl(item.cover_url)} />
              </Link>
              <div className="min-w-0 flex-1 space-y-3">
                <div className="flex flex-wrap items-center gap-3">
                  <Link
                    className="text-xl font-semibold text-ink"
                    to={`/video/${item.video_id}?from=${encodeURIComponent(`/history?page=${page}`)}`}
                  >
                    {item.video_title}
                  </Link>
                  <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">{item.area_name}</span>
                </div>
                <p className="text-sm text-slate-500">
                  作者：
                  <Link className="text-sea" to={`/users/${item.author_id}?from=${encodeURIComponent(`/history?page=${page}`)}`}>
                    {item.author_name}
                  </Link>
                </p>
                <div className="flex flex-wrap gap-3 text-sm text-slate-500">
                  <span>观看于 {formatDate(item.watched_at)}</span>
                  <span>已看 {formatCount(item.progress_seconds)} 秒</span>
                  <span>总时长 {formatCount(item.duration_seconds)} 秒</span>
                </div>
              </div>
            </article>
          ))}

          <div className="flex items-center justify-between rounded-md bg-white px-5 py-4 shadow-sm">
            {page > 1 ? (
              <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to={buildPageLink(page - 1)}>
                上一页
              </Link>
            ) : (
              <span className="rounded-full border border-slate-100 px-4 py-2 text-sm text-slate-300">上一页</span>
            )}
            <span className="text-sm text-slate-500">
              第 {page} / {totalPages} 页 · 共 {total} 条
            </span>
            {page < totalPages ? (
              <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to={buildPageLink(page + 1)}>
                下一页
              </Link>
            ) : (
              <span className="rounded-full border border-slate-100 px-4 py-2 text-sm text-slate-300">下一页</span>
            )}
          </div>
        </div>
      ) : null}
    </AppShell>
  );
}
