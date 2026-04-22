import { useEffect, useState } from 'react';
import { AdminAPI } from '@/api/modules/admin';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { formatDate, resolveMediaUrl } from '@/shared/utils/format';
import type { AdminVideoItem, AreaStatsItem, TodayStats } from '@/shared/types/domain';

type AdminTab = 'pending' | 'reviewed';

export function AdminReviewPage() {
  const [activeTab, setActiveTab] = useState<AdminTab>('pending');
  const [videos, setVideos] = useState<AdminVideoItem[]>([]);
  const [todayStats, setTodayStats] = useState<TodayStats | null>(null);
  const [areaStats, setAreaStats] = useState<AreaStatsItem[]>([]);
  const [rejectReason, setRejectReason] = useState<Record<number, string>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadVideos = async (tab: AdminTab) => {
    setLoading(true);
    setError(null);
    try {
      const response = tab === 'pending' ? await AdminAPI.listPending({ page: 1, page_size: 20 }) : await AdminAPI.listReviewed({ review_status: 'reviewed', page: 1, page_size: 20 });
      setVideos(response.list);
    } catch (requestError) {
      setError('审核列表加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadVideos(activeTab);
  }, [activeTab]);

  useEffect(() => {
    async function loadStats() {
      try {
        const [today, areas] = await Promise.all([AdminAPI.getTodayStats(), AdminAPI.getAreaStats()]);
        setTodayStats(today);
        setAreaStats(areas);
      } catch (requestError) {
        setError('统计数据加载失败');
      }
    }

    void loadStats();
  }, []);

  const approve = async (videoId: number) => {
    await AdminAPI.approve(videoId);
    await loadVideos(activeTab);
  };

  const reject = async (videoId: number) => {
    await AdminAPI.reject(videoId, { reason: rejectReason[videoId] || '内容不符合要求' });
    await loadVideos(activeTab);
  };

  return (
    <AppShell>
      <section className="mb-6 rounded-md bg-white p-6 shadow-sm">
        <p className="text-sm uppercase tracking-[0.3em] text-sea">Admin Review</p>
        <h1 className="mt-2 text-3xl font-semibold text-ink">审核管理后台</h1>
      </section>

      {todayStats ? (
        <section className="mb-6 grid gap-4 md:grid-cols-5">
          <div className="rounded-md bg-white p-4 shadow-sm text-sm text-slate-600">活跃用户 {todayStats.active_user_count}</div>
          <div className="rounded-md bg-white p-4 shadow-sm text-sm text-slate-600">投稿 {todayStats.submitted_video_count}</div>
          <div className="rounded-md bg-white p-4 shadow-sm text-sm text-slate-600">过审 {todayStats.approved_video_count}</div>
          <div className="rounded-md bg-white p-4 shadow-sm text-sm text-slate-600">播放 {todayStats.play_count}</div>
          <div className="rounded-md bg-white p-4 shadow-sm text-sm text-slate-600">评论 {todayStats.comment_count}</div>
        </section>
      ) : null}

      {areaStats.length > 0 ? (
        <section className="mb-6 rounded-md bg-white p-6 shadow-sm">
          <h2 className="text-xl font-semibold text-ink">分区统计</h2>
          <div className="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {areaStats.map((item) => (
              <div className="rounded-md bg-paper p-4 text-sm text-slate-600" key={item.area_id}>
                <p className="font-medium text-ink">{item.area_name}</p>
                <p>通过 {item.approved_count} / 待审 {item.pending_count} / 驳回 {item.rejected_count}</p>
              </div>
            ))}
          </div>
        </section>
      ) : null}

      <div className="mb-4 flex gap-3">
        {(['pending', 'reviewed'] as AdminTab[]).map((tab) => (
          <button
            className={`rounded-full px-4 py-2 text-sm transition ${activeTab === tab ? 'bg-ink text-white' : 'bg-white text-slate-600'}`}
            key={tab}
            onClick={() => setActiveTab(tab)}
            type="button"
          >
            {tab === 'pending' ? '未审核' : '已审核'}
          </button>
        ))}
      </div>

      {error ? <div className="mb-4 rounded-md bg-amber-50 px-4 py-3 text-sm text-amber-700">{error}</div> : null}
      {loading ? (
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, index) => (
            <LoadingBlock key={index} lines={5} />
          ))}
        </div>
      ) : videos.length === 0 ? (
        <EmptyState title="当前没有对应的审核稿件" />
      ) : (
        <div className="space-y-4">
          {videos.map((item) => (
            <article className="rounded-md bg-white p-5 shadow-sm" key={item.id}>
              <div className="flex flex-col gap-4 md:flex-row">
                <img alt={item.title} className="aspect-video w-full rounded-md object-cover md:w-72" src={resolveMediaUrl(item.cover_url)} />
                <div className="min-w-0 flex-1 space-y-3">
                  <div className="flex flex-wrap items-center gap-3">
                    <h3 className="text-xl font-semibold text-ink">{item.title}</h3>
                    <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">{item.author_username}</span>
                    <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">{item.area_name}</span>
                    <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">{item.review_status}</span>
                  </div>
                  <p className="text-sm text-slate-600">{item.description}</p>
                  {item.review_reason ? <p className="text-sm text-rose-500">驳回原因：{item.review_reason}</p> : null}
                  <p className="text-xs text-slate-400">创建于 {formatDate(item.created_at)}</p>
                  {activeTab === 'pending' ? (
                    <div className="space-y-3">
                      <textarea
                        className="min-h-24 w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none focus:border-accent"
                        onChange={(e) => setRejectReason((current) => ({ ...current, [item.id]: e.target.value }))}
                        placeholder="填写驳回原因"
                        value={rejectReason[item.id] || ''}
                      />
                      <div className="flex gap-3">
                        <button className="rounded-full bg-ink px-4 py-2 text-sm text-white" onClick={() => void approve(item.id)} type="button">
                          通过
                        </button>
                        <button className="rounded-full bg-rose-500 px-4 py-2 text-sm text-white" onClick={() => void reject(item.id)} type="button">
                          驳回
                        </button>
                      </div>
                    </div>
                  ) : null}
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </AppShell>
  );
}
