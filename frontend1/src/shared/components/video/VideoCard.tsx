import { Link } from 'react-router-dom';
import { LazyImage } from '@/shared/components/media/LazyImage';
import { formatCount, formatDate, resolveMediaUrl } from '@/shared/utils/format';
import type { FeedItem } from '@/shared/types/domain';

type VideoCardProps = {
  item: FeedItem;
  to?: string;
};

export function VideoCard({ item, to }: VideoCardProps) {
  return (
    <Link
      className="group overflow-hidden rounded-md bg-white shadow-sm transition hover:-translate-y-1"
      to={to ?? `/video/${item.id}`}
    >
      <div className="relative aspect-video overflow-hidden">
        <LazyImage alt={item.title} className="h-full w-full" src={resolveMediaUrl(item.cover_url)} />
        <div className="absolute inset-x-0 bottom-0 flex items-center justify-between bg-gradient-to-t from-black/70 to-transparent px-4 py-3 text-xs text-white">
          <span>{formatCount(item.view_count)} 播放</span>
          <span>{formatDate(item.published_at)}</span>
        </div>
      </div>
      <div className="space-y-3 p-4">
        <h3 className="line-clamp-2 text-sm font-semibold text-ink transition group-hover:text-accent">{item.title}</h3>
        <div className="flex items-center justify-between text-xs text-slate-500">
          <span>{item.author.username}</span>
          <span>{formatCount(item.comment_count)} 评论</span>
        </div>
      </div>
    </Link>
  );
}
