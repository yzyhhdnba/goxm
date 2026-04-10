# PILIPILI Go API 清单

## 1. 统一约定

Base URL：

```text
/api/v1
```

统一响应建议：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

弱动态列表分页建议：

- `page`
- `page_size`

动态 Feed 分页建议：

- 首页推荐流、关注流、热门榜、分区流统一优先使用 Cursor-based 分页
- 推荐参数统一为 `cursor + limit`
- 搜索、通知、历史、后台列表等弱动态列表可继续使用 `page + page_size`

认证方式：

- `Authorization: Bearer <token>`

---

## 2. 认证与账户

## 2.1 注册

`POST /api/v1/auth/register`

请求体：

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "123456"
}
```

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user": {
      "id": 1,
      "username": "alice"
    }
  }
}
```

## 2.2 登录

`POST /api/v1/auth/login`

请求体：

```json
{
  "username": "alice",
  "password": "123456"
}
```

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "access_token": "jwt-access-token",
    "refresh_token": "jwt-refresh-token",
    "access_token_expires_in": 7200,
    "refresh_token_expires_in": 1209600,
    "user": {
      "id": 1,
      "username": "alice",
      "role": "user"
    }
  }
}
```

说明：

- `access_token` 建议短期有效，例如 2 小时
- `refresh_token` 建议长期有效，例如 14 天
- 第一阶段默认采用单设备可轮换的 refresh token 机制
- 当前实现中 `username` 字段第一阶段兼容“用户名或邮箱”登录

## 2.3 刷新 Token

`POST /api/v1/auth/refresh`

请求体：

```json
{
  "refresh_token": "jwt-refresh-token"
}
```

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "access_token": "new-jwt-access-token",
    "refresh_token": "new-jwt-refresh-token",
    "access_token_expires_in": 7200,
    "refresh_token_expires_in": 1209600
  }
}
```

说明：

- refresh 成功后旧 `refresh_token` 应立即失效
- refresh 成功后旧 `access_token` 不立即失效，而是保留到自然过期或显式登出

## 2.4 退出登录

`POST /api/v1/auth/logout`

说明：

- 退出登录时应使当前 access token 和 refresh token 一并失效
- 当前实现通过 access token 鉴权后执行退出

## 2.5 检查用户名

`GET /api/v1/auth/check-username?username=alice`

## 2.6 检查邮箱

`GET /api/v1/auth/check-email?email=alice@example.com`

## 2.7 获取当前用户

`GET /api/v1/users/me`

## 2.8 更新当前用户资料

`PATCH /api/v1/users/me`

请求体：

```json
{
  "avatar_url": "https://...",
  "bio": "hello"
}
```

## 2.9 修改密码

`PATCH /api/v1/users/password`

请求体：

```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

---

## 3. 分区

## 3.1 获取分区列表

`GET /api/v1/areas`

返回字段建议：

- `id`
- `name`
- `slug`

说明：

- 当前实现会在启动阶段自动补齐默认分区种子

---

## 4. 视频

## 4.1 首页推荐流

`GET /api/v1/feed/recommend?cursor=1711800000:10&limit=8`

说明：

- 对应你当前首页 `videoList`
- 当前实现的 `next_cursor` 使用 `unix_timestamp:video_id` 复合游标，避免同一秒内数据乱序
- 为兼容早期调用，服务端仍可接受只有时间戳的 cursor，但推荐以后端返回的 `next_cursor` 为准
- 第一阶段可按 `published_at desc, id desc` 排序
- 第二阶段再演进为推荐流/热门混合流
- 返回建议包含 `items`、`next_cursor`、`has_more`
- 当前实现仅返回 `visible + approved + published_at 非空` 的公开视频

## 4.2 热门榜

`GET /api/v1/feed/hot?cursor=1200:1711800000:10&limit=20`

说明：

- 当前已实现
- 当前按 `hot_score desc, published_at desc, id desc` 排序
- 当前 `next_cursor` 使用 `hot_score:unix_timestamp:video_id`
- 推荐直接使用服务端返回的 `next_cursor` 继续翻页，而不是手写 cursor

## 4.3 关注流

`GET /api/v1/feed/following?cursor=1711800000:10&limit=20`

需要登录

说明：

- 当前已实现
- 当前只返回“我关注的作者”发布的公开视频
- 当前按 `published_at desc, id desc` 排序
- 当前 `next_cursor` 使用 `unix_timestamp:video_id`

## 4.4 分区视频列表

`GET /api/v1/areas/:id/videos?cursor=1711800000:10&limit=20&sort=latest`

说明：

- 当前已实现
- 分区列表也视为动态流，当前阶段默认保持 `cursor + limit`
- 当前仅支持 `sort=latest`，默认值也是 `latest`
- 当前按 `published_at desc, id desc` 排序
- 当前 `next_cursor` 使用 `unix_timestamp:video_id`

## 4.5 视频详情

`GET /api/v1/videos/:id`

返回字段建议：

```json
{
  "id": 10,
  "title": "视频标题",
  "description": "视频简介",
  "cover_url": "https://...",
  "play_url": "https://.../index.m3u8",
  "view_count": 100,
  "comment_count": 12,
  "like_count": 66,
  "favorite_count": 8,
  "published_at": "2026-03-30T10:00:00Z",
  "author": {
    "id": 2,
    "username": "tom",
    "avatar_url": "https://..."
  },
  "viewer_state": {
    "liked": true,
    "favorited": false,
    "followed": true
  }
}
```

说明：

- 这里借鉴 `feedsystem_video_go` 的软鉴权思路
- 未登录时 `viewer_state` 可返回默认 false
- 当前实现中，登录后会返回真实的点赞、收藏、关注状态

## 4.6 获取作者视频列表

`GET /api/v1/users/:id/videos?page=1&page_size=20`

说明：

- 当前已实现
- 当前返回结构使用 `list + pagination`
- 因作者主页视频列表相对稳定，当前阶段仍保留 `page + page_size`

## 4.7 创建视频元数据

`POST /api/v1/videos`

请求体：

```json
{
  "area_id": 1,
  "title": "测试视频",
  "description": "简介"
}
```

说明：

- 当前已实现
- 需要登录
- 当前采取“先创建记录，再上传文件”的三段式链路
- 创建后默认 `review_status = pending`

## 4.8 上传视频文件

`POST /api/v1/videos/:id/source`

`multipart/form-data`

字段：

- `file`

说明：

- 当前已实现
- 需要登录
- 当前只允许作者本人上传
- 当前默认落盘到 `backend/storage/videos/{video_id}/source.ext`
- 当前 `play_url` 直接返回 `/uploads/...` 静态路径，后续再升级 ffmpeg / HLS

## 4.9 编辑已有稿件

`PATCH /api/v1/videos/:id`

请求体：

```json
{
  "area_id": 2,
  "title": "更新后的标题",
  "description": "更新后的简介"
}
```

说明：

- 当前已实现
- 需要登录
- 当前只允许作者本人编辑
- 编辑后会重置为 `review_status = pending`
- 当前会清空 `review_reason`，并将 `published_at` 置空
- 如果稿件此前已经通过审核，当前会同步回退作者 `video_count`，等待重新审核通过后再恢复

## 4.10 上传封面

`POST /api/v1/videos/:id/cover`

`multipart/form-data`

字段：

- `file`

说明：

- 当前已实现
- 需要登录
- 当前只允许作者本人上传
- 当前默认落盘到 `backend/storage/videos/{video_id}/cover.ext`

## 4.11 获取播放地址

`GET /api/v1/videos/:id/play`

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "play_url": "https://.../index.m3u8"
  }
}
```

## 4.12 记录播放历史

`POST /api/v1/histories`

请求体：

```json
{
  "video_id": 10,
  "progress_seconds": 120
}
```

说明：

- 当前已实现
- 需要登录
- 当前按 `(user_id, video_id)` 做 upsert，重复观看会更新同一条历史

---

## 5. 视频互动

## 5.1 点赞视频

`POST /api/v1/videos/:id/likes`

## 5.2 取消点赞

`DELETE /api/v1/videos/:id/likes`

## 5.3 查询点赞状态

`GET /api/v1/videos/:id/likes/me`

需要登录

## 5.4 收藏视频

`POST /api/v1/videos/:id/favorites`

## 5.5 取消收藏

`DELETE /api/v1/videos/:id/favorites`

## 5.6 查询收藏状态

`GET /api/v1/videos/:id/favorites/me`

需要登录

---

## 6. 评论与回复

## 6.1 获取评论列表

`GET /api/v1/videos/:id/comments?page=1&page_size=20`

说明：

- 当前实现未登录也可访问
- 已登录时会补充评论 `viewer_state.liked`

返回结构建议：

```json
{
  "list": [
    {
      "id": 1,
      "content": "一级评论",
      "like_count": 3,
      "reply_count": 2,
      "created_at": "2026-03-30T10:00:00Z",
      "user": {
        "id": 2,
        "username": "tom"
      },
      "viewer_state": {
        "liked": false
      }
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

## 6.2 发布一级评论

`POST /api/v1/videos/:id/comments`

请求体：

```json
{
  "content": "这是一条评论"
}
```

## 6.3 获取某条评论的回复

`GET /api/v1/comments/:id/replies?page=1&page_size=20`

## 6.4 回复评论

`POST /api/v1/comments/:id/replies`

请求体：

```json
{
  "content": "这是一条回复"
}
```

## 6.5 点赞评论

`POST /api/v1/comments/:id/likes`

## 6.6 取消点赞评论

`DELETE /api/v1/comments/:id/likes`

## 6.7 查询评论点赞状态

`GET /api/v1/comments/:id/likes/me`

需要登录

---

## 7. 关注关系

## 7.1 关注用户

`POST /api/v1/users/:id/follow`

## 7.2 取消关注

`DELETE /api/v1/users/:id/follow`

## 7.3 查询关注状态

`GET /api/v1/users/:id/follow-status`

需要登录

## 7.4 粉丝列表

`GET /api/v1/users/:id/followers?page=1&page_size=20`

## 7.5 关注列表

`GET /api/v1/users/:id/following?page=1&page_size=20`

---

## 8. 搜索

## 8.1 搜视频

`GET /api/v1/search/videos?keyword=xxx&page=1&page_size=20`

说明：

- 当前已实现
- 当前返回结构使用 `list + pagination`
- 当前仍使用 MySQL `LIKE '%keyword%'` 作为第一阶段实现
- 数据量上来后应升级到 MySQL FULLTEXT 或 Elasticsearch

## 8.2 搜用户

`GET /api/v1/search/users?keyword=xxx&page=1&page_size=20`

说明：

- 当前已实现
- 当前返回结构使用 `list + pagination`
- 当前仍使用 MySQL `LIKE '%keyword%'` 作为第一阶段实现

---

## 9. 历史与通知

## 9.1 获取历史记录

`GET /api/v1/histories?page=1&page_size=20`

说明：

- 当前已实现
- 需要登录
- 当前返回视频标题、封面、作者、分区、观看时间和进度

## 9.2 获取通知列表

`GET /api/v1/notices?page=1&page_size=20`

说明：

- 当前已实现
- 需要登录
- 当前主要承接管理员审核通过 / 驳回后的创作者通知

## 9.3 标记通知已读

`PATCH /api/v1/notices/:id/read`

说明：

- 当前已实现
- 需要登录

---

## 10. 个人中心

## 10.1 获取指定用户主页

`GET /api/v1/users/:id/profile`

说明：

- 当前已实现
- 未登录也可以访问
- 已登录时会补充 `viewer_state.followed`

返回建议：

```json
{
  "id": 2,
  "username": "tom",
  "avatar_url": "https://...",
  "bio": "hello",
  "follower_count": 10,
  "following_count": 8,
  "video_count": 6,
  "viewer_state": {
    "followed": true
  }
}
```

## 10.2 获取我的主页聚合信息

`GET /api/v1/users/me/dashboard`

可包含：

- 用户信息
- 聚合统计
- 最近公开视频
- 最近收藏
- 最近关注用户

当前返回字段：

- `user`
- `stats.total_view_count`
- `recent_videos`
- `favorite_videos`
- `following_users`

说明：

- 当前已实现
- 需要登录
- 这已经承接原前端里的 `getMyIndexInfo`

---

## 11. 创作者与审核

## 11.1 获取我的投稿列表

`GET /api/v1/creator/videos?review_status=approved&page=1&page_size=20`

支持：

- `approved`
- `pending`
- `rejected`
- `all`

说明：

- 当前已实现
- 需要登录
- 当前返回结构使用 `list + pagination`

这可以承接你现在前端里的：

- `creator/Approved`
- `creator/ToBeApproved`
- `creator/NoApproved`

## 11.2 获取待审核视频

`GET /api/v1/admin/videos/pending?page=1&page_size=20`

补充：

- 当前还实现了 `GET /api/v1/admin/videos?review_status=reviewed&page=1&page_size=20`
- `review_status` 当前支持：`pending / approved / rejected / reviewed / all`
- 以上接口都需要管理员身份

## 11.3 审核通过

`POST /api/v1/admin/videos/:id/approve`

说明：

- 当前已实现
- 需要管理员身份

## 11.4 审核拒绝

`POST /api/v1/admin/videos/:id/reject`

请求体：

```json
{
  "reason": "封面不合规"
}
```

说明：

- 当前已实现
- 需要管理员身份

## 11.5 今日统计

`GET /api/v1/admin/stats/today`

当前返回字段：

- `active_user_count`
- `submitted_video_count`
- `approved_video_count`
- `play_count`
- `comment_count`

## 11.6 分区统计

`GET /api/v1/admin/stats/area`

当前返回字段：

- `area_id`
- `area_name`
- `approved_count`
- `pending_count`
- `rejected_count`
- `total_count`

---

## 12. 与旧接口的迁移映射

| 旧接口 | 新接口 |
| --- | --- |
| `login` | `POST /api/v1/auth/login` |
| - | `POST /api/v1/auth/refresh` |
| `auth/register` | `POST /api/v1/auth/register` |
| `auth/check/name` | `GET /api/v1/auth/check-username` |
| `auth/check/email` | `GET /api/v1/auth/check-email` |
| `videoList` | `GET /api/v1/feed/recommend` / `GET /api/v1/areas/:id/videos` |
| `video` | `GET /api/v1/videos/:id` |
| `video/like` | `POST /api/v1/videos/:id/likes` |
| `video/dislike` | `DELETE /api/v1/videos/:id/likes` |
| `video/collect` | `POST /api/v1/videos/:id/favorites` |
| `video/cancel` | `DELETE /api/v1/videos/:id/favorites` |
| `comment` | `GET /api/v1/videos/:id/comments` |
| `comment/make` | `POST /api/v1/videos/:id/comments` |
| `reply` | `GET /api/v1/comments/:id/replies` |
| `reply/make` | `POST /api/v1/comments/:id/replies` |
| `follow` | `POST /api/v1/users/:id/follow` |
| `unfollow` | `DELETE /api/v1/users/:id/follow` |
| `histories` | `GET /api/v1/histories` |
| `notice` | `GET /api/v1/notices` |
| `user/read` | `PATCH /api/v1/notices/:id/read` |
| `search/video` | `GET /api/v1/search/videos` |
| `search/user` | `GET /api/v1/search/users` |
| `area` | `GET /api/v1/areas` |
| `admin/video` | `GET /api/v1/admin/videos/pending` |
| `admin/approve` | `POST /api/v1/admin/videos/:id/approve` |
| `admin/deny` | `POST /api/v1/admin/videos/:id/reject` |

---

## 13. 第一阶段真正要先做的接口

如果下一步要开始搭 Go 后端，最小集合是：

1. `POST /api/v1/auth/register`
2. `POST /api/v1/auth/login`
3. `POST /api/v1/auth/refresh`
4. `GET /api/v1/users/me`
5. `POST /api/v1/auth/logout`
6. `GET /api/v1/areas`
7. `GET /api/v1/feed/recommend`
8. `GET /api/v1/videos/:id`
9. `POST /api/v1/videos/:id/likes`
10. `DELETE /api/v1/videos/:id/likes`
11. `POST /api/v1/videos/:id/favorites`
12. `DELETE /api/v1/videos/:id/favorites`
13. `GET /api/v1/videos/:id/comments`
14. `POST /api/v1/videos/:id/comments`
15. `POST /api/v1/comments/:id/replies`
16. `POST /api/v1/users/:id/follow`
17. `DELETE /api/v1/users/:id/follow`

这批接口完成后，鉴权主链路和首页、详情、评论、点赞、收藏、关注基本就都能跑起来。
