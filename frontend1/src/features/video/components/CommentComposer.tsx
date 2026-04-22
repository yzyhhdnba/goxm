import { FormEvent, useState } from 'react';

type CommentComposerProps = {
  disabled: boolean;
  submitting: boolean;
  onSubmit: (content: string) => Promise<void>;
};

export function CommentComposer({ disabled, submitting, onSubmit }: CommentComposerProps) {
  const [content, setContent] = useState('');

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const nextValue = content.trim();
    if (!nextValue) {
      return;
    }

    await onSubmit(nextValue);
    setContent('');
  };

  return (
    <form className="space-y-3 rounded-md bg-white p-5 shadow-sm" onSubmit={handleSubmit}>
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-ink">发表评论</h2>
        {!disabled ? null : <span className="text-sm text-slate-400">登录后可评论</span>}
      </div>
      <textarea
        className="min-h-28 w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none transition focus:border-accent disabled:cursor-not-allowed"
        disabled={disabled || submitting}
        onChange={(event) => setContent(event.target.value)}
        placeholder="请输入一条友善的评论"
        value={content}
      />
      <button
        className="rounded-full bg-ink px-5 py-3 text-sm font-medium text-white transition hover:bg-sea disabled:cursor-not-allowed disabled:opacity-60"
        disabled={disabled || submitting || !content.trim()}
        type="submit"
      >
        {submitting ? '发布中...' : '发布评论'}
      </button>
    </form>
  );
}
