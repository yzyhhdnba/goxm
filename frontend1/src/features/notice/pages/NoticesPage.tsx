import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { NoticeAPI } from '@/api/modules/notice';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { formatDate } from '@/shared/utils/format';
import type { NoticeItem } from '@/shared/types/domain';

export function NoticesPage() {
  const [items, setItems] = useState<NoticeItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        const response = await NoticeAPI.list({ page: 1, page_size: 20 });
        setItems(response.list);
      } catch (requestError) {
        setError('通知加载失败');
      } finally {
        setLoading(false);
      }
    }

    void load();
  }, []);

  const markRead = async (noticeId: number) => {
    const updated = await NoticeAPI.markRead(noticeId);
    setItems((current) => current.map((item) => (item.id === noticeId ? updated : item)));
  };

  return (
    <AppShell>
      <section className="mb-6 rounded-md bg-white p-6 shadow-sm">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-sea">Notices</p>
            <h1 className="mt-2 text-3xl font-semibold text-ink">系统通知</h1>
          </div>
          <Link className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" to="/me">
            返回个人中心
          </Link>
        </div>
      </section>
      {loading ? (
        <div className="space-y-4">
          {Array.from({ length: 4 }).map((_, index) => (
            <LoadingBlock key={index} lines={4} />
          ))}
        </div>
      ) : error ? (
        <EmptyState title={error} />
      ) : items.length === 0 ? (
        <EmptyState title="当前还没有系统通知" />
      ) : (
        <div className="space-y-4">
          {items.map((item) => (
            <article className="rounded-md bg-white p-5 shadow-sm" key={item.id}>
              <div className="flex items-start justify-between gap-4">
                <div className="space-y-2">
                  <div className="flex flex-wrap items-center gap-3">
                    <h2 className="text-lg font-semibold text-ink">{item.title}</h2>
                    {!item.read ? <span className="rounded-full bg-rose-100 px-2 py-1 text-xs text-rose-500">未读</span> : null}
                  </div>
                  <p className="text-sm leading-7 text-slate-600">{item.content}</p>
                  <p className="text-xs text-slate-400">{formatDate(item.created_at)}</p>
                </div>
                {!item.read ? (
                  <button className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" onClick={() => void markRead(item.id)} type="button">
                    标记已读
                  </button>
                ) : null}
              </div>
            </article>
          ))}
        </div>
      )}
    </AppShell>
  );
}
