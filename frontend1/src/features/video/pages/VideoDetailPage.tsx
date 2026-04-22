import { useEffect, useMemo, useState } from 'react';
import { Link, useParams, useSearchParams } from 'react-router-dom';
import { CommentAPI } from '@/api/modules/comment';
import { VideoAPI } from '@/api/modules/video';
import { CommentComposer } from '@/features/video/components/CommentComposer';
import { CommentList } from '@/features/video/components/CommentList';
import { VideoHero } from '@/features/video/components/VideoHero';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { VideoCard } from '@/shared/components/video/VideoCard';
import { useAppSelector } from '@/store/hooks';
import type { CommentItem, FeedItem, VideoDetail } from '@/shared/types/domain';

export function VideoDetailPage() {
  const { id } = useParams();
  const [searchParams] = useSearchParams();
  const isAuthenticated = useAppSelector((state) => state.auth.isAuthenticated);
  const [detail, setDetail] = useState<VideoDetail | null>(null);
  const [recommend, setRecommend] = useState<FeedItem[]>([]);
  const [comments, setComments] = useState<CommentItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submittingComment, setSubmittingComment] = useState(false);

  useEffect(() => {
    if (!id) {
      setError('缺少视频 ID');
      setLoading(false);
      return;
    }
    const videoId = id;

    async function load() {
      setLoading(true);
      setError(null);

      try {
        const [video, commentList, recommendList] = await Promise.all([
          VideoAPI.getDetail(videoId),
          CommentAPI.getList(videoId, { page: 1, page_size: 20 }),
          VideoAPI.getRecommend({ limit: 6 }),
        ]);

        setDetail(video);
        setComments(commentList.list);
        setRecommend(recommendList.items.filter((item) => item.id !== Number(videoId)).slice(0, 4));
      } catch (requestError) {
        setError('视频详情加载失败');
      } finally {
        setLoading(false);
      }
    }

    void load();
  }, [id]);

  const interactionBlocked = useMemo(() => !detail || !isAuthenticated, [detail, isAuthenticated]);

  const ensureDetail = () => {
    if (!detail) {
      return null;
    }

    if (!isAuthenticated) {
      setError('当前交互需要先登录');
      return null;
    }

    return detail;
  };

  const updateViewerState = (patch: Partial<VideoDetail['viewer_state']>, patchCounts?: Partial<VideoDetail>) => {
    setDetail((current) => {
      if (!current) {
        return current;
      }

      return {
        ...current,
        ...patchCounts,
        viewer_state: {
          ...current.viewer_state,
          ...patch,
        },
      };
    });
  };

  const handleLike = async () => {
    const current = ensureDetail();
    if (!current) {
      return;
    }

    try {
      if (current.viewer_state.liked) {
        await VideoAPI.unlike(current.id);
        updateViewerState({ liked: false }, { like_count: Math.max(0, current.like_count - 1) });
        return;
      }

      await VideoAPI.like(current.id);
      updateViewerState({ liked: true }, { like_count: current.like_count + 1 });
    } catch (requestError) {
      setError('点赞操作失败');
    }
  };

  const handleFavorite = async () => {
    const current = ensureDetail();
    if (!current) {
      return;
    }

    try {
      if (current.viewer_state.favorited) {
        await VideoAPI.unfavorite(current.id);
        updateViewerState({ favorited: false }, { favorite_count: Math.max(0, current.favorite_count - 1) });
        return;
      }

      await VideoAPI.favorite(current.id);
      updateViewerState({ favorited: true }, { favorite_count: current.favorite_count + 1 });
    } catch (requestError) {
      setError('收藏操作失败');
    }
  };

  const handleFollow = async () => {
    const current = ensureDetail();
    if (!current) {
      return;
    }

    try {
      if (current.viewer_state.followed) {
        await VideoAPI.unfollow(current.author.id);
        updateViewerState({ followed: false });
        return;
      }

      await VideoAPI.follow(current.author.id);
      updateViewerState({ followed: true });
    } catch (requestError) {
      setError('关注操作失败');
    }
  };

  const handleCommentSubmit = async (content: string) => {
    if (!id || interactionBlocked) {
      setError('发表评论前请先登录');
      return;
    }
    const videoId = id;

    setSubmittingComment(true);
    try {
      const created = await CommentAPI.create(videoId, { content });
      setComments((current) => [created, ...current]);
      setDetail((current) => (current ? { ...current, comment_count: current.comment_count + 1 } : current));
      setError(null);
    } catch (requestError) {
      setError('发表评论失败');
    } finally {
      setSubmittingComment(false);
    }
  };

  const backTo = searchParams.get('from') || '/';
  const currentPath = useMemo(() => `/video/${id}${searchParams.toString() ? `?${searchParams.toString()}` : ''}`, [id, searchParams]);
  const authorProfileLink = useMemo(() => {
    if (!detail) {
      return '/';
    }

    const from = searchParams.get('from');
    if (from) {
      const match = from.match(/^\/users\/(\d+)(\?.*)?$/);
      if (match && Number(match[1]) === detail.author.id) {
        return from;
      }
    }

    return `/users/${detail.author.id}?from=${encodeURIComponent(currentPath)}`;
  }, [currentPath, detail, searchParams]);

  return (
    <AppShell>
      <div className="mb-6">
        <Link className="text-sm text-slate-500 transition hover:text-ink" to={backTo}>
          ← 返回上一页
        </Link>
      </div>

      {loading ? (
        <div className="grid gap-6 lg:grid-cols-[1.5fr_0.7fr]">
          <LoadingBlock lines={8} />
          <LoadingBlock lines={6} />
        </div>
      ) : null}

      {!loading && error && !detail ? <EmptyState description="请确认视频 ID 有效，且后端接口可以访问。" title={error} /> : null}

      {!loading && detail ? (
        <div className="grid gap-6 lg:grid-cols-[1.5fr_0.7fr]">
          <div className="space-y-6">
            {error ? <div className="rounded-md bg-amber-50 px-4 py-3 text-sm text-amber-700">{error}</div> : null}
            <VideoHero detail={detail} isAuthenticated={isAuthenticated} onFavorite={() => void handleFavorite()} onFollow={() => void handleFollow()} onLike={() => void handleLike()} />
            <CommentComposer disabled={!isAuthenticated} onSubmit={handleCommentSubmit} submitting={submittingComment} />
            <CommentList comments={comments} disabled={!isAuthenticated} />
          </div>

          <aside className="space-y-4">
            <div className="rounded-md bg-white p-5 shadow-sm">
              <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Author</p>
              <h2 className="mt-3 text-xl font-semibold text-ink">{detail.author.username}</h2>
              <p className="mt-2 text-sm leading-6 text-slate-500">
                如果你是从用户主页进入详情，这里会回到原来的主页状态；否则会进入作者主页。
              </p>
              <Link
                className="mt-4 inline-flex rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-700 transition hover:border-accent hover:text-accent"
                to={authorProfileLink}
              >
                查看作者主页
              </Link>
            </div>

            <div className="space-y-4">
              <div>
                <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Recommend</p>
                <h2 className="text-xl font-semibold text-ink">推荐视频</h2>
              </div>

              {recommend.map((item) => (
                <VideoCard item={item} key={item.id} to={`/video/${item.id}?from=${encodeURIComponent(currentPath)}`} />
              ))}
            </div>
          </aside>
        </div>
      ) : null}
    </AppShell>
  );
}
