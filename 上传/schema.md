# PILIPILI Go Schema 草案

## 1. 设计原则

- 先保证业务闭环，再考虑复杂拆分
- 第一阶段以 MySQL 单库为主
- 统计字段保留在主表用于查询，但写路径优先考虑 Redis 原子累加再回刷 MySQL
- 点赞、收藏、关注等关系表全部做唯一索引
- 评论和回复合并建模，避免两套表
- Redis 只做缓存和登录态辅助，不承担主存储职责
- `status` 负责业务状态，`deleted_at` 负责 GORM 软删除

---

## 2. 表结构总览

第一期建议建设以下表：

1. `users`
2. `areas`
3. `videos`
4. `video_likes`
5. `video_favorites`
6. `comments`
7. `comment_likes`
8. `follows`
9. `view_histories`
10. `notices`
11. `video_reviews`

---

## 3. 详细表设计

## 3.1 `users`

用途：用户账号、身份和展示资料

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 用户 ID |
| username | varchar(64) | unique, not null | 用户名 |
| email | varchar(128) | unique, not null | 邮箱 |
| password_hash | varchar(255) | not null | bcrypt 哈希 |
| token_version | int unsigned | not null, default 1 | access/refresh token 版本 |
| refresh_token_hash | varchar(128) | null | 当前有效 refresh token 哈希 |
| avatar_url | varchar(255) | null | 头像 |
| bio | varchar(255) | null | 简介 |
| role | varchar(16) | not null, default `user` | `user/admin` |
| status | varchar(16) | not null, default `active` | `active/disabled` |
| follower_count | int unsigned | not null, default 0 | 粉丝数 |
| following_count | int unsigned | not null, default 0 | 关注数 |
| video_count | int unsigned | not null, default 0 | 已发布视频数 |
| created_at | datetime | not null | 创建时间 |
| updated_at | datetime | not null | 更新时间 |
| deleted_at | datetime | null | GORM 软删除时间 |

索引建议：

- `uniq_users_username`
- `uniq_users_email`
- `idx_users_role_status`

---

## 3.2 `areas`

用途：视频分区

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 分区 ID |
| name | varchar(64) | unique, not null | 分区名称 |
| slug | varchar(64) | unique, not null | 英文标识 |
| sort_order | int | not null, default 0 | 排序 |
| status | varchar(16) | not null, default `active` | 状态 |
| created_at | datetime | not null | 创建时间 |
| updated_at | datetime | not null | 更新时间 |
| deleted_at | datetime | null | GORM 软删除时间 |

预置数据建议：

- 动画
- 番剧
- 音乐
- 游戏
- 科技
- 生活

---

## 3.3 `videos`

用途：视频核心实体

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 视频 ID |
| author_id | bigint unsigned | not null | 作者 ID |
| area_id | bigint unsigned | not null | 分区 ID |
| title | varchar(128) | not null | 标题 |
| description | text | null | 简介 |
| cover_url | varchar(255) | null | 封面地址 |
| source_path | varchar(255) | null | 原始视频路径 |
| play_url | varchar(255) | null | 播放地址 |
| duration_seconds | int unsigned | not null, default 0 | 时长 |
| status | varchar(16) | not null, default `visible` | `visible/hidden/deleted` |
| review_status | varchar(16) | not null, default `pending` | `pending/approved/rejected` |
| review_reason | varchar(255) | null | 驳回原因 |
| published_at | datetime | null | 发布时间 |
| like_count | int unsigned | not null, default 0 | 点赞数 |
| favorite_count | int unsigned | not null, default 0 | 收藏数 |
| comment_count | int unsigned | not null, default 0 | 评论数 |
| view_count | int unsigned | not null, default 0 | 播放数 |
| hot_score | bigint | not null, default 0 | 热度分数 |
| created_at | datetime | not null | 创建时间 |
| updated_at | datetime | not null | 更新时间 |
| deleted_at | datetime | null | GORM 软删除时间 |

索引建议：

- `idx_videos_author_id`
- `idx_videos_area_id_review_status`
- `idx_videos_published_at`
- `idx_videos_hot_score`
- `idx_videos_area_id_published_at`

说明：

- 第一阶段保留 `videos` 统计字段，方便详情和列表查询
- 建议写路径优先进入 Redis，再由后台任务、ticker 或异步链路批量回刷 MySQL
- 如果某一阶段还未接入 Redis 计数，则必须明确这只是临时简化实现

---

## 3.4 `video_likes`

用途：视频点赞关系

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| video_id | bigint unsigned | not null | 视频 ID |
| user_id | bigint unsigned | not null | 用户 ID |
| created_at | datetime | not null | 点赞时间 |

索引建议：

- 唯一索引 `(video_id, user_id)`
- 索引 `(user_id, created_at)`

---

## 3.5 `video_favorites`

用途：视频收藏关系

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| video_id | bigint unsigned | not null | 视频 ID |
| user_id | bigint unsigned | not null | 用户 ID |
| created_at | datetime | not null | 收藏时间 |

索引建议：

- 唯一索引 `(video_id, user_id)`
- 索引 `(user_id, created_at)`

---

## 3.6 `comments`

用途：统一承载一级评论和回复

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 评论 ID |
| video_id | bigint unsigned | not null | 视频 ID |
| user_id | bigint unsigned | not null | 评论用户 |
| root_id | bigint unsigned | not null, default 0 | 一级评论 ID |
| parent_id | bigint unsigned | not null, default 0 | 直接父评论 ID |
| content | text | not null | 评论内容 |
| reply_count | int unsigned | not null, default 0 | 回复数 |
| like_count | int unsigned | not null, default 0 | 点赞数 |
| status | varchar(16) | not null, default `visible` | `visible/deleted/hidden` |
| created_at | datetime | not null | 创建时间 |
| updated_at | datetime | not null | 更新时间 |
| deleted_at | datetime | null | GORM 软删除时间 |

索引建议：

- `idx_comments_video_root_created` on `(video_id, root_id, created_at desc)`
- `idx_comments_parent_id`
- `idx_comments_user_id`

建模规则：

- 一级评论：`root_id = 0`，`parent_id = 0`
- 回复评论：`root_id = 一级评论 ID`
- 回复对象：`parent_id = 直接回复的评论 ID`

这样可以兼容你当前前端的：

- 评论列表
- 楼中楼回复
- 回复点赞

补充说明：

- `status` 用于业务层展示和隐藏控制
- `deleted_at` 用于技术层软删除和 GORM 默认过滤

---

## 3.7 `comment_likes`

用途：评论点赞关系

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| comment_id | bigint unsigned | not null | 评论 ID |
| user_id | bigint unsigned | not null | 用户 ID |
| created_at | datetime | not null | 点赞时间 |

索引建议：

- 唯一索引 `(comment_id, user_id)`
- 索引 `(user_id, created_at)`

---

## 3.8 `follows`

用途：关注关系

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| follower_id | bigint unsigned | not null | 关注者 |
| followee_id | bigint unsigned | not null | 被关注者 |
| created_at | datetime | not null | 创建时间 |

索引建议：

- 唯一索引 `(follower_id, followee_id)`
- 索引 `(followee_id, created_at)`

---

## 3.9 `view_histories`

用途：观看历史

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| user_id | bigint unsigned | not null | 用户 ID |
| video_id | bigint unsigned | not null | 视频 ID |
| progress_seconds | int unsigned | not null, default 0 | 播放进度 |
| watched_at | datetime | not null | 最后观看时间 |
| created_at | datetime | not null | 创建时间 |
| updated_at | datetime | not null | 更新时间 |

索引建议：

- 唯一索引 `(user_id, video_id)`
- 索引 `(user_id, watched_at desc)`

说明：

- 如果用户重复观看同一个视频，更新同一条记录即可

---

## 3.10 `notices`

用途：站内通知

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| user_id | bigint unsigned | not null | 接收通知用户 |
| type | varchar(32) | not null | `comment/reply/follow/system/review` |
| title | varchar(128) | not null | 标题 |
| content | varchar(255) | not null | 内容 |
| related_id | bigint unsigned | null | 关联业务 ID |
| is_read | tinyint(1) | not null, default 0 | 是否已读 |
| created_at | datetime | not null | 创建时间 |
| updated_at | datetime | not null | 更新时间 |

索引建议：

- `idx_notices_user_read_created` on `(user_id, is_read, created_at desc)`

---

## 3.11 `video_reviews`

用途：保留审核记录

| 字段 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint unsigned | PK | 主键 |
| video_id | bigint unsigned | not null | 视频 ID |
| reviewer_id | bigint unsigned | not null | 审核人 |
| status | varchar(16) | not null | `approved/rejected` |
| reason | varchar(255) | null | 驳回原因 |
| created_at | datetime | not null | 创建时间 |

索引建议：

- `idx_video_reviews_video_id`
- `idx_video_reviews_reviewer_id`

---

## 4. Redis Key 设计

第一阶段只建议做下面这些：

### 4.1 登录态

- `auth:user:{userID}` -> token 或 token version

### 4.2 视频详情缓存

- `video:detail:{videoID}`

当前状态：

- 当前已落地第一版读缓存
- 缓存未命中时回退 MySQL
- 视频元数据、源文件、封面、评论、点赞、收藏、审核流变更后会主动失效

### 4.3 热门榜缓存

- `feed:hot`
- `feed:hot:area:{areaID}`

当前状态：

- 当前已落地 `feed:hot` 的短 TTL 读缓存
- 目前尚未实现分区热门榜缓存
- 当前尚未实现“先写 Redis 再异步回刷 MySQL”的计数主写路径

### 4.4 用户关系状态缓存

- `follow:status:{followerID}:{followeeID}`

---

## 5. 迁移顺序建议

1. `users`
2. `areas`
3. `videos`
4. `video_likes`
5. `video_favorites`
6. `comments`
7. `comment_likes`
8. `follows`
9. `view_histories`
10. `notices`
11. `video_reviews`

原因：

- 先建基础主实体
- 再建关系表
- 最后建运营和消息表

---

## 6. GORM Model 组织建议

建议按业务域组织，而不是把所有 model 扔到一个目录：

```text
internal/account/entity.go
internal/video/entity.go
internal/comment/entity.go
internal/social/entity.go
internal/history/entity.go
internal/notice/entity.go
internal/admin/entity.go
```

这样更接近 `feedsystem_video_go` 的拆法，也更适合后续扩展 service/repo。

---

## 7. 第一阶段最小可落地集合

如果要尽快做出第一个可运行版本，优先只建这 7 张表：

1. `users`
2. `areas`
3. `videos`
4. `video_likes`
5. `video_favorites`
6. `comments`
7. `follows`

原因：

- 已经足够支撑登录、首页、详情、评论、点赞、收藏、关注
- 历史、通知、审核记录可以第二批再补

---

## 8. 后续扩展方向

第二阶段扩展：

- `view_histories`
- `notices`
- `video_reviews`

第三阶段扩展：

- `outbox_events`
- `media_jobs`
- `audit_logs`

这样可以让 `PILIPILI` 从单体业务项目平滑升级到带 Worker 的工程项目。
