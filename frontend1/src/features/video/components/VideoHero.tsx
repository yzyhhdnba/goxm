import { formatCount, formatDate, resolveMediaUrl } from '@/shared/utils/format';
import type { VideoDetail } from '@/shared/types/domain';

type VideoHeroProps = {
  detail: VideoDetail;
  isAuthenticated: boolean;
  onLike: () => void;
  onFavorite: () => void;
  onFollow: () => void;
};

export function VideoHero({ detail, isAuthenticated, onLike, onFavorite, onFollow }: VideoHeroProps) {
  return (
    <section className="overflow-hidden rounded-md bg-white shadow-sm">
      <div className="aspect-video bg-ink">
        {detail.play_url ? (
          <video className="h-full w-full object-cover" controls poster={resolveMediaUrl(detail.cover_url)} src={resolveMediaUrl(detail.play_url)} />
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-white/70">当前视频还没有可播放资源</div>
        )}
      </div>

      <div className="space-y-5 p-6">
        <div className="space-y-3">
          <h1 className="text-3xl font-semibold text-ink">{detail.title}</h1>
          <div className="flex flex-wrap gap-3 text-sm text-slate-500">
            <span>{formatCount(detail.view_count)} 播放</span>
            <span>{formatCount(detail.comment_count)} 评论</span>
            <span>{formatDate(detail.published_at)}</span>
            <span>作者：{detail.author.username}</span>
          </div>
        </div>

        <p className="rounded-md bg-paper px-4 py-4 text-sm leading-7 text-slate-700">
          {detail.description || '作者还没有填写简介。'}
        </p>

        <div className="flex flex-wrap gap-3">
          <button
            className={`rounded-full px-5 py-3 text-sm font-medium transition ${detail.viewer_state.liked ? 'bg-accent text-white' : 'bg-mist text-sea'}`}
            onClick={onLike}
            type="button"
          >
            {detail.viewer_state.liked ? '已点赞' : '点赞'} · {formatCount(detail.like_count)}
          </button>
          <button
            className={`rounded-full px-5 py-3 text-sm font-medium transition ${detail.viewer_state.favorited ? 'bg-ink text-white' : 'bg-mist text-sea'}`}
            onClick={onFavorite}
            type="button"
          >
            {detail.viewer_state.favorited ? '已收藏' : '收藏'} · {formatCount(detail.favorite_count)}
          </button>
          <button
            className="rounded-full border border-slate-200 px-5 py-3 text-sm font-medium text-slate-700 transition hover:border-accent hover:text-accent"
            onClick={onFollow}
            type="button"
          >
            {detail.viewer_state.followed ? '已关注作者' : '关注作者'}
          </button>
          {!isAuthenticated ? <span className="self-center text-sm text-slate-400">登录后可执行互动操作</span> : null}
        </div>
      </div>
    </section>
  );
}
