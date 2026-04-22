import { FormEvent, useMemo, useState } from 'react';
import { CommentAPI } from '@/api/modules/comment';
import { formatDate } from '@/shared/utils/format';
import type { CommentItem } from '@/shared/types/domain';

const replyPageSize = 20;

type ReplyThreadProps = {
  parentCommentId: number;
  replies: CommentItem[];
  totalReplies: number;
  disabled: boolean;
  loading: boolean;
  error?: string | null;
  page: number;
  totalPages: number;
  onLoadPage: (page: number) => Promise<void>;
  onReplyCreated: (created: CommentItem) => void;
};

export function ReplyThread({
  parentCommentId,
  replies,
  totalReplies,
  disabled,
  loading,
  error,
  page,
  totalPages,
  onLoadPage,
  onReplyCreated,
}: ReplyThreadProps) {
  const [activeReplyTarget, setActiveReplyTarget] = useState<number | null>(null);
  const [content, setContent] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [items, setItems] = useState<CommentItem[]>(replies);

  const replyItems = useMemo(() => items, [items]);

  const targetLabel = useMemo(() => {
    if (activeReplyTarget === null || activeReplyTarget === parentCommentId) {
      return '回复楼主';
    }
    const target = replyItems.find((item) => item.id === activeReplyTarget);
    return target ? `回复 ${target.user.username}` : '回复评论';
  }, [activeReplyTarget, parentCommentId, replyItems]);

  const openReplyBox = (targetId: number) => {
    setActiveReplyTarget(targetId);
    setContent('');
  };

  const closeReplyBox = () => {
    setActiveReplyTarget(null);
    setContent('');
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const targetId = activeReplyTarget ?? parentCommentId;
    const nextContent = content.trim();
    if (!nextContent || disabled) {
      return;
    }

    setSubmitting(true);
    try {
      const created = await CommentAPI.createReply(targetId, { content: nextContent });
      setItems((current) => [...current, created]);
      onReplyCreated(created);
      closeReplyBox();
    } finally {
      setSubmitting(false);
    }
  };

  const toggleLike = async (replyId: number) => {
    if (disabled) {
      return;
    }

    const target = replyItems.find((item) => item.id === replyId);
    if (!target) {
      return;
    }

    const nextLiked = !target.viewer_state.liked;
    setItems((current) =>
      current.map((item) =>
        item.id === replyId
          ? {
              ...item,
              like_count: Math.max(0, item.like_count + (nextLiked ? 1 : -1)),
              viewer_state: {
                ...item.viewer_state,
                liked: nextLiked,
              },
            }
          : item,
      ),
    );

    try {
      if (nextLiked) {
        await CommentAPI.like(replyId);
      } else {
        await CommentAPI.unlike(replyId);
      }
    } catch (error) {
      setItems((current) =>
        current.map((item) =>
          item.id === replyId
            ? {
                ...item,
                like_count: Math.max(0, item.like_count + (nextLiked ? -1 : 1)),
                viewer_state: {
                  ...item.viewer_state,
                  liked: !nextLiked,
                },
              }
            : item,
        ),
      );
    }
  };

  return (
    <div className="space-y-3 rounded-md bg-paper p-4">
      {loading ? <p className="text-sm text-slate-400">回复加载中...</p> : null}
      {!loading && error ? <p className="text-sm text-red-500">{error}</p> : null}

      {replyItems.map((reply) => (
        <div key={reply.id} className="rounded-md bg-white p-4">
          <div className="flex items-start justify-between gap-3">
            <div className="space-y-2">
              <div className="flex flex-wrap items-center gap-3">
                <p className="text-sm font-semibold text-ink">{reply.user.username}</p>
                <span className="text-xs text-slate-400">{formatDate(reply.created_at)}</span>
                {reply.viewer_state.liked ? (
                  <span className="rounded-full bg-red-50 px-2 py-1 text-xs text-red-500">已点赞</span>
                ) : null}
              </div>
              <p className="whitespace-pre-line text-sm leading-6 text-slate-700">{reply.content}</p>
              <div className="flex flex-wrap items-center gap-3 text-xs text-slate-500">
                <button
                  className={`rounded-full px-3 py-2 transition ${
                    reply.viewer_state.liked ? 'bg-red-50 text-red-500' : 'border border-slate-200 text-slate-600'
                  }`}
                  disabled={disabled}
                  onClick={() => void toggleLike(reply.id)}
                  type="button"
                >
                  {reply.viewer_state.liked ? '取消点赞' : '点赞'} · {reply.like_count}
                </button>
              </div>
            </div>
            <button
              className="rounded-full border border-slate-200 px-3 py-2 text-xs text-slate-600 transition hover:border-accent hover:text-accent"
              onClick={() => openReplyBox(reply.id)}
              type="button"
            >
              回复
            </button>
          </div>
        </div>
      ))}

      <div className="flex gap-3">
        <button
          className="rounded-full border border-slate-200 px-4 py-2 text-xs text-slate-600 transition hover:border-accent hover:text-accent"
          onClick={() => openReplyBox(parentCommentId)}
          type="button"
        >
          回复楼主
        </button>
      </div>

      {!loading && !error && totalReplies > replyPageSize ? (
        <div className="flex items-center justify-between rounded-md bg-white px-4 py-3">
          {page > 1 ? (
            <button
              className="rounded-full border border-slate-200 px-3 py-2 text-xs text-slate-600"
              onClick={() => void onLoadPage(page - 1)}
              type="button"
            >
              上一页回复
            </button>
          ) : (
            <span className="rounded-full border border-slate-100 px-3 py-2 text-xs text-slate-300">上一页回复</span>
          )}
          <span className="text-xs text-slate-500">
            第 {page} / {totalPages} 页 · 共 {totalReplies} 条
          </span>
          {page < totalPages ? (
            <button
              className="rounded-full border border-slate-200 px-3 py-2 text-xs text-slate-600"
              onClick={() => void onLoadPage(page + 1)}
              type="button"
            >
              下一页回复
            </button>
          ) : (
            <span className="rounded-full border border-slate-100 px-3 py-2 text-xs text-slate-300">下一页回复</span>
          )}
        </div>
      ) : null}

      {activeReplyTarget !== null ? (
        <form className="space-y-3 rounded-md border border-slate-200 bg-white px-4 py-4" onSubmit={handleSubmit}>
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-ink">{targetLabel}</p>
            <button className="text-xs text-slate-400" onClick={closeReplyBox} type="button">
              取消
            </button>
          </div>
          <textarea
            className="min-h-24 w-full rounded-md border border-slate-200 bg-slate-50 px-3 py-3 outline-none transition focus:border-accent disabled:cursor-not-allowed"
            disabled={disabled || submitting}
            onChange={(event) => setContent(event.target.value)}
            placeholder="请输入回复内容"
            value={content}
          />
          <button
            className="rounded-full bg-ink px-4 py-2 text-sm font-medium text-white transition hover:bg-sea disabled:cursor-not-allowed disabled:opacity-60"
            disabled={disabled || submitting || !content.trim()}
            type="submit"
          >
            {submitting ? '发布中...' : '发布回复'}
          </button>
        </form>
      ) : null}
    </div>
  );
}
