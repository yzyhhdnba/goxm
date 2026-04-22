import type { ReactNode } from 'react';
import { Link } from 'react-router-dom';
import { formatCount } from '@/shared/utils/format';
import type { SearchUserItem, SocialUserCard } from '@/shared/types/domain';

type UserSearchCardProps = {
  item: SearchUserItem | SocialUserCard;
  to?: string;
  meta?: ReactNode;
};

export function UserSearchCard({ item, to, meta }: UserSearchCardProps) {
  const initials = item.username?.slice(0, 1).toUpperCase() || 'U';

  return (
    <Link
      className="flex items-start gap-4 rounded-md bg-white p-5 shadow-sm transition hover:-translate-y-1"
      to={to ?? `/users/${item.id}`}
    >
      <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-mist text-lg font-semibold text-sea">
        {initials}
      </div>
      <div className="min-w-0 flex-1 space-y-2">
        <div className="flex flex-wrap items-center gap-3">
          <h3 className="text-lg font-semibold text-ink">{item.username}</h3>
          {'follower_count' in item ? (
            <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">
              {formatCount(item.follower_count)} 粉丝
            </span>
          ) : null}
          {'video_count' in item ? (
            <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">
              {formatCount(item.video_count)} 投稿
            </span>
          ) : null}
          {meta}
        </div>
        <p className="line-clamp-2 text-sm leading-6 text-slate-600">{item.bio || '这个用户还没有填写简介。'}</p>
      </div>
    </Link>
  );
}
