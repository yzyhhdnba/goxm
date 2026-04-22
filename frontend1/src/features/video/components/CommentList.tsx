import { CommentThread } from '@/features/video/components/CommentThread';
import { EmptyState } from '@/shared/components/common/EmptyState';
import type { CommentItem } from '@/shared/types/domain';

type CommentListProps = {
  comments: CommentItem[];
  disabled: boolean;
};

export function CommentList({ comments, disabled }: CommentListProps) {
  if (comments.length === 0) {
    return <EmptyState description="首批实现先展示一级评论列表，回复链路后续再补。" title="还没有评论" />;
  }

  return (
    <div className="space-y-4">
      {comments.map((item) => (
        <CommentThread comment={item} disabled={disabled} key={item.id} />
      ))}
    </div>
  );
}
