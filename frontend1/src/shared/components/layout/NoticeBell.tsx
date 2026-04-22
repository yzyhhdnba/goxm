import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { NoticeAPI } from '@/api/modules/notice';
import { formatDate } from '@/shared/utils/format';
import type { NoticeItem } from '@/shared/types/domain';

export function NoticeBell() {
  const [items, setItems] = useState<NoticeItem[]>([]);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    async function load() {
      try {
        const response = await NoticeAPI.list({ page: 1, page_size: 10 });
        setItems(response.list);
      } catch (error) {
        setItems([]);
      }
    }

    void load();
  }, []);

  const unreadCount = useMemo(() => items.filter((item) => !item.read).length, [items]);

  const markRead = async (noticeId: number) => {
    const updated = await NoticeAPI.markRead(noticeId);
    setItems((current) => current.map((item) => (item.id === noticeId ? updated : item)));
  };

  return (
    <div className="relative">
      <button className="relative rounded-full border border-slate-200 px-3 py-2 text-sm text-slate-600" onClick={() => setOpen((v) => !v)} type="button">
        通知
        {unreadCount > 0 ? <span className="ml-2 rounded-full bg-rose-500 px-2 py-0.5 text-xs text-white">{unreadCount}</span> : null}
      </button>
      {open ? (
        <div className="absolute right-0 top-12 z-30 w-96 rounded-md border border-white/60 bg-white p-4 shadow-sm">
          <div className="mb-3 flex items-center justify-between">
            <p className="text-sm font-semibold text-ink">系统通知</p>
            <Link className="text-xs text-sea" onClick={() => setOpen(false)} to="/notices">
              查看全部
            </Link>
          </div>
          <div className="space-y-3">
            {items.length === 0 ? <p className="text-sm text-slate-400">当前没有通知</p> : null}
            {items.map((item) => (
              <div className="rounded-md bg-paper p-3" key={item.id}>
                <div className="flex items-start justify-between gap-3">
                  <div className="space-y-1">
                    <p className="text-sm font-medium text-ink">{item.title}</p>
                    <p className="line-clamp-2 text-xs text-slate-500">{item.content}</p>
                    <p className="text-xs text-slate-400">{formatDate(item.created_at)}</p>
                  </div>
                  {!item.read ? (
                    <button className="rounded-full border border-slate-200 px-2 py-1 text-xs text-slate-600" onClick={() => void markRead(item.id)} type="button">
                      已读
                    </button>
                  ) : null}
                </div>
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}
