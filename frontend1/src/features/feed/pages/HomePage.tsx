import { useEffect, useState } from 'react';
import { FeedAPI } from '@/api/modules/feed';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { VideoCard } from '@/shared/components/video/VideoCard';
import type { Area, FeedItem } from '@/shared/types/domain';

const defaultHeroCards = [
  '基于 React 重建首页主链路',
  '统一 Axios 封装与鉴权恢复',
  '首批仅保留可联调核心页面',
];

export function HomePage() {
  const [items, setItems] = useState<FeedItem[]>([]);
  const [areas, setAreas] = useState<Area[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError(null);

      try {
        const [recommend, areaList] = await Promise.all([
          FeedAPI.getRecommend({ limit: 8 }),
          FeedAPI.getAreas(),
        ]);
        setItems(recommend.items);
        setAreas(areaList);
      } catch (requestError) {
        setError('首页数据加载失败');
      } finally {
        setLoading(false);
      }
    }

    void load();
  }, []);

  return (
    <AppShell
      hero={
        <section className="overflow-hidden border-b border-white/40">
          <div className="mx-auto grid max-w-6xl gap-8 px-4 py-10 md:grid-cols-[1.3fr_0.7fr] md:px-6 md:py-16">
            <div className="space-y-5">
              <p className="text-sm uppercase tracking-[0.35em] text-sea">React Parallel Frontend</p>
              <h1 className="max-w-2xl text-4xl font-semibold leading-tight text-ink md:text-5xl">
                先把首页、登录、视频详情跑通，再逐步补完整站。
              </h1>
              <p className="max-w-xl text-base leading-7 text-slate-600">
                这一版优先稳定交付：新建 `frontend1`，不替换现有 Vue 前端；先验证真实后端接口，再扩到更复杂页面。
              </p>
            </div>

            <div className="grid gap-4">
              {defaultHeroCards.map((label) => (
                <div key={label} className="rounded-md bg-white p-5 shadow-sm">
                  <p className="text-sm uppercase tracking-[0.2em] text-slate-400">当前进度</p>
                  <p className="mt-3 text-lg font-semibold text-ink">{label}</p>
                </div>
              ))}
            </div>
          </div>
        </section>
      }
    >
      <section className="mb-8 rounded-md bg-white p-5 shadow-sm">
        <div className="flex flex-wrap items-center gap-3">
          <span className="text-sm font-semibold text-ink">分区</span>
          {areas.map((area) => (
            <span key={area.id} className="rounded-full bg-mist px-3 py-2 text-sm text-sea">
              {area.name}
            </span>
          ))}
        </div>
      </section>

      <section className="space-y-6">
        <div className="flex items-end justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Recommend</p>
            <h2 className="text-2xl font-semibold text-ink">首页推荐流</h2>
          </div>
        </div>

        {loading ? (
          <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
            {Array.from({ length: 8 }).map((_, index) => (
              <LoadingBlock key={index} lines={5} />
            ))}
          </div>
        ) : null}

        {!loading && error ? <EmptyState description="请先确认后端 API 可访问。" title={error} /> : null}

        {!loading && !error && items.length === 0 ? <EmptyState title="当前还没有推荐内容" /> : null}

        {!loading && !error && items.length > 0 ? (
          <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
            {items.map((item) => (
              <VideoCard item={item} key={item.id} />
            ))}
          </div>
        ) : null}
      </section>
    </AppShell>
  );
}
