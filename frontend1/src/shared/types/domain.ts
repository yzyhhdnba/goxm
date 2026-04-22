export type User = {
  id: number;
  username: string;
  email?: string;
  avatar_url: string;
  bio?: string;
  role?: 'user' | 'admin';
  follower_count?: number;
  following_count?: number;
  video_count?: number;
};

export type ViewerState = {
  liked: boolean;
  favorited: boolean;
  followed: boolean;
};

export type FeedItem = {
  id: number;
  area_id: number;
  title: string;
  description: string;
  cover_url: string;
  play_url: string;
  duration_seconds: number;
  view_count: number;
  comment_count: number;
  like_count: number;
  favorite_count: number;
  published_at: string;
  author: Pick<User, 'id' | 'username' | 'avatar_url'>;
};

export type FeedResponse = {
  items: FeedItem[];
  next_cursor: string;
  has_more: boolean;
};

export type VideoDetail = {
  id: number;
  area_id: number;
  title: string;
  description: string;
  cover_url: string;
  play_url: string;
  duration_seconds: number;
  view_count: number;
  comment_count: number;
  like_count: number;
  favorite_count: number;
  published_at: string;
  author: Pick<User, 'id' | 'username' | 'avatar_url'>;
  viewer_state: ViewerState;
};

export type CommentItem = {
  id: number;
  root_id: number;
  parent_id: number;
  content: string;
  like_count: number;
  reply_count: number;
  created_at: string;
  user: Pick<User, 'id' | 'username' | 'avatar_url'>;
  viewer_state: {
    liked: boolean;
  };
};

export type CommentListResponse = {
  list: CommentItem[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
  };
};

export type DashboardStats = {
  total_view_count: number;
};

export type DashboardResponse = {
  user: User;
  stats: DashboardStats;
  recent_videos: FeedItem[];
  favorite_videos: FeedItem[];
  following_users: SearchUserItem[];
};

export type HistoryItem = {
  video_id: number;
  video_title: string;
  cover_url: string;
  play_url: string;
  author_id: number;
  author_name: string;
  area_name: string;
  watched_at: string;
  progress_seconds: number;
  duration_seconds: number;
};

export type HistoryListResponse = {
  list: HistoryItem[];
  pagination: SearchPagination;
};

export type CreatorVideoItem = {
  id: number;
  area_id: number;
  area_name: string;
  title: string;
  description: string;
  cover_url: string;
  play_url: string;
  source_path: string;
  duration_seconds: number;
  review_status: 'pending' | 'approved' | 'rejected';
  review_reason: string;
  view_count: number;
  comment_count: number;
  like_count: number;
  favorite_count: number;
  created_at: string;
  updated_at: string;
  published_at?: string | null;
};

export type CreatorVideoListResponse = {
  list: CreatorVideoItem[];
  pagination: SearchPagination;
};

export type CreateVideoInput = {
  area_id: number;
  title: string;
  description: string;
};

export type CreateVideoResponse = {
  id: number;
  area_id: number;
  title: string;
  description: string;
  review_status: string;
  created_at: string;
};

export type UploadSourceResponse = {
  video_id: number;
  source_path: string;
  play_url: string;
};

export type UploadCoverResponse = {
  video_id: number;
  cover_url: string;
};

export type AdminVideoItem = {
  id: number;
  area_id: number;
  area_name: string;
  title: string;
  description: string;
  cover_url: string;
  play_url: string;
  source_path: string;
  duration_seconds: number;
  review_status: 'pending' | 'approved' | 'rejected';
  review_reason: string;
  view_count: number;
  comment_count: number;
  like_count: number;
  favorite_count: number;
  created_at: string;
  updated_at: string;
  published_at?: string | null;
  author_id: number;
  author_username: string;
};

export type AdminVideoListResponse = {
  list: AdminVideoItem[];
  pagination: SearchPagination;
};

export type TodayStats = {
  active_user_count: number;
  submitted_video_count: number;
  approved_video_count: number;
  play_count: number;
  comment_count: number;
};

export type AreaStatsItem = {
  area_id: number;
  area_name: string;
  approved_count: number;
  pending_count: number;
  rejected_count: number;
  total_count: number;
};

export type NoticeItem = {
  id: number;
  type: string;
  title: string;
  content: string;
  related_video_id?: number | null;
  read: boolean;
  read_at?: string | null;
  created_at: string;
};

export type NoticeListResponse = {
  list: NoticeItem[];
  pagination: SearchPagination;
};

export type LoginInput = {
  username: string;
  password: string;
};

export type LoginResponse = {
  access_token: string;
  refresh_token: string;
  access_token_expires_in: number;
  refresh_token_expires_in: number;
  user: User;
};

export type Area = {
  id: number;
  name: string;
  slug: string;
};

export type SearchVideoItem = FeedItem;

export type SearchUserItem = {
  id: number;
  username: string;
  avatar_url: string;
  bio: string;
  follower_count: number;
  video_count: number;
};

export type SearchPagination = {
  page: number;
  page_size: number;
  total: number;
};

export type SearchVideoResponse = {
  list: SearchVideoItem[];
  pagination: SearchPagination;
};

export type SearchUserResponse = {
  list: SearchUserItem[];
  pagination: SearchPagination;
};

export type UserProfile = {
  id: number;
  username: string;
  avatar_url: string;
  bio: string;
  follower_count: number;
  following_count: number;
  video_count: number;
  viewer_state: {
    followed: boolean;
  };
};

export type PagedFeedResponse = {
  list: FeedItem[];
  pagination: SearchPagination;
};

export type SocialUserCard = {
  id: number;
  username: string;
  avatar_url: string;
  bio: string;
};

export type SocialUserListResponse = {
  list: SocialUserCard[];
  pagination: SearchPagination;
};
