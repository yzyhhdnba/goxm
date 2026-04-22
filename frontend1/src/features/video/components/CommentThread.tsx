import { useState } from 'react';
import { CommentAPI } from '@/api/modules/comment';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { formatDate } from '@/shared/utils/format';
import type { CommentItem } from '@/shared/types/domain';
import { ReplyThread } from '@/features/video/components/ReplyThread';

const replyPageSize = 20;

type CommentThreadProps = {
  comment: CommentItem;
  disabled: boolean;
};

export function CommentThread({ comment, disabled }: CommentThreadProps) {
  const [item, setItem] = useState(comment);
  const [expanded, setExpanded] = useState(false);
  const [loadingReplies, setLoadingReplies] = useState(false);
  const [replies, setReplies] = useState<CommentItem[]>([]);
  const [replyCount, setReplyCount] = useState(comment.reply_count);
  const [replyPage, setReplyPage] = useState(1);
  const [replyTotal, setReplyTotal] = useState(comment.reply_count);
  const [replyError, setReplyError] = useState<string | null>(null);
  const [likePending, setLikePending] = useState(false);

  const replyTotalPages = Math.max(1, Math.ceil(replyTotal / replyPageSize));

  const loadReplies = async (page: number) => {
    setLoadingReplies(true);
    setReplyError(null);
    try {
      const response = await CommentAPI.getReplies(item.id, { page, page_size: replyPageSize });
      setReplies(response.list);
      setReplyPage(page);
      setReplyTotal(response.pagination.total);
    } catch (error) {
      setReplyError('回复加载失败');
    } finally {
      setLoadingReplies(false);
    }
  };

  const toggleReplies = async () => {
    if (expanded) {
      setExpanded(false);
      return;
    }

    setExpanded(true);
    if (replies.length > 0) {
      return;
    }
    await loadReplies(1);
  };

  const handleReplyCreated = (created: CommentItem) => {
    setReplies((current) => [...current, created]);
    setReplyCount((current) => current + 1);
    setReplyTotal((current) => current + 1);
    setItem((current) => ({
      ...current,
      reply_count: current.reply_count + 1,
    }));
    setExpanded(true);
  };

  const handleLikeToggle = async () => {
    if (disabled || likePending) {
      return;
    }

    const nextLiked = !item.viewer_state.liked;
    setLikePending(true);
    setItem((current) => ({
      ...current,
      like_count: Math.max(0, current.like_count + (nextLiked ? 1 : -1)),
      viewer_state: {
        ...current.viewer_state,
        liked: nextLiked,
      },
    }));

    try {
      if (nextLiked) {
        await CommentAPI.like(item.id);
      } else {
        await CommentAPI.unlike(item.id);
      }
    } catch (error) {
      setItem((current) => ({
        ...current,
        like_count: Math.max(0, current.like_count + (nextLiked ? -1 : 1)),
        viewer_state: {
          ...current.viewer_state,
          liked: !nextLiked,
        },
      }));
    } finally {
      setLikePending(false);
    }
  };

  return (
    <article className="rounded-md bg-white p-5 shadow-sm">
      <div className="mb-3 flex items-center justify-between">
        <div>
          <p className="text-sm font-semibold text-ink">{item.user.username}</p>
          <p className="text-xs text-slate-400">{formatDate(item.created_at)}</p>
        </div>
        <div className="flex items-center gap-3 text-xs text-slate-400">
          <span>{item.like_count} 赞</span>
          <span>{replyCount} 回复</span>
          {item.viewer_state.liked ? <span className="rounded-full bg-red-50 px-2 py-1 text-red-500">已点赞</span> : null}
        </div>
      </div>

      <p className="whitespace-pre-line text-sm leading-7 text-slate-700">{item.content}</p>

      <div className="mt-4 flex flex-wrap items-center gap-3">
        <button
          className={`rounded-full px-4 py-2 text-xs transition ${
            item.viewer_state.liked ? 'bg-red-50 text-red-500' : 'border border-slate-200 text-slate-600'
          }`}
          disabled={disabled || likePending}
          onClick={() => void handleLikeToggle()}
          type="button"
        >
          {item.viewer_state.liked ? '取消点赞' : '点赞'} · {item.like_count}
        </button>
        <button
          className="rounded-full border border-slate-200 px-4 py-2 text-xs text-slate-600 transition hover:border-accent hover:text-accent"
          onClick={() => void toggleReplies()}
          type="button"
        >
          {expanded ? '收起回复' : `查看回复${replyCount > 0 ? `（${replyCount}）` : ''}`}
        </button>
      </div>

      {expanded ? (
        <div className="mt-4">
          {!loadingReplies && !replyError && replies.length === 0 ? (
            <EmptyState description="现在可以直接在这里发布第一条回复。" title="还没有回复" />
          ) : null}
          {!replyError ? (
            <ReplyThread
              disabled={disabled}
              error={replyError}
              loading={loadingReplies}
              onReplyCreated={handleReplyCreated}
              onLoadPage={loadReplies}
              page={replyPage}
              parentCommentId={item.id}
              replies={replies}
              totalPages={replyTotalPages}
              totalReplies={replyTotal}
            />
          ) : null}
        </div>
      ) : null}
    </article>
  );
}
