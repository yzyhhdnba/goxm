# PILIPILI Go 化蓝图

## 1. 目标定位

把当前的 `PILIPILI` 从“已有完整交互雏形的 Vue 视频站前端项目”升级为“可用于简历和面试展示的 Go 后端完整项目”。

目标不是简单把接口语言改成 Go，而是完成下面这次升级：

- 从前端主导项目，升级成前后端协作的完整系统
- 从零散接口调用，升级成清晰的业务域和分层架构
- 从单机演示风格，升级成具备缓存、异步、容器化能力的工程项目
- 从“视频网站页面实现”，升级成“在线视频与互动平台后端系统”

最终推荐的项目定位：

> 基于 Go 的在线视频与弹幕互动平台，支持用户鉴权、视频发布与播放、评论、回复、点赞、收藏、关注、历史记录、通知、审核管理与 Feed/热门推荐。

---

## 2. 当前仓库审计结论

仓库：`https://gitee.com/yzynba/pilipili`

当前项目是一个典型的 Vue 3 + TypeScript 单前端工程：

- `Vue 3`
- `Vue Router`
- `Vuex`
- `Element Plus`
- `Axios`
- `DPlayer`
- `HLS.js`
- `ECharts`

### 2.1 当前目录结构说明

现有仓库的业务页面已经不算少，说明它非常适合做 Go 项目的前端基底：

- 首页：`Home.vue`
- 视频详情页：`views/video/index.vue`
- 搜索页：`views/search/*`
- 个人中心：`views/personal/*`
- 上传中心：`views/upload/*`
- 历史记录：`views/history/index.vue`
- 通知页：`views/notice/index.vue`
- 后台审核管理：`views/management/*`

说明这个项目不是“只有首页和登录弹窗”的壳子，而是已经形成了较完整的业务外观。

### 2.2 当前已经具备的业务能力

从前端调用的接口和页面结构看，当前业务能力已经覆盖：

- 用户注册、登录
- 首页视频列表
- 视频详情
- 视频播放
- 视频点赞、收藏
- 评论
- 评论回复
- 评论点赞
- 关注/取关
- 搜索视频、搜索用户
- 历史记录
- 通知读取
- 个人中心数据
- 视频上传
- 上传后审核状态查询
- 后台视频审核
- 后台数据统计

### 2.3 当前前端依赖的后端接口清单

前端代码里已出现的接口包括：

- `login`
- `auth/register`
- `auth/check/name`
- `auth/check/email`
- `videoList`
- `video`
- `video/like`
- `video/dislike`
- `video/liked`
- `video/collect`
- `video/cancel`
- `video/collected`
- `comment`
- `comment/make`
- `comment/like`
- `comment/dislike`
- `comment/liked`
- `reply`
- `reply/make`
- `reply/like`
- `reply/dislike`
- `reply/liked`
- `follow`
- `unfollow`
- `followed`
- `histories`
- `notice`
- `user/read`
- `search/video`
- `search/user`
- `getMyIndexInfo`
- `area`
- `upload/Video`
- `upload/VideoInfo`
- `upload/VideoCover`
- `video/processed/video`
- `creator/Approved`
- `creator/ToBeApproved`
- `creator/NoApproved`
- `admin/video`
- `admin/approve`
- `admin/deny`
- `admin/TodayData`
- `admin/AreaData`
- `cover/cover`

这份清单非常关键，因为它说明你不是从零设计业务，而是已经有一套前端视角的“事实接口规范”。

### 2.4 当前项目的问题

虽然业务面不错，但当前后端视角上还有明显问题：

- 前端直接写死了后端地址 `http://172.20.10.6:8081`
- 接口命名不统一，REST 风格和动作式命名混杂
- 登录态主要靠 `localStorage`，缺少正式 token 体系
- 没有明显的统一请求层、统一错误码、统一鉴权
- 评论、回复、点赞、收藏、关注等领域边界尚未抽象成后端模块
- 上传、封面、HLS 播放、审核、通知之间缺少统一的数据链路
- 缺少工程化目录、容器化和后端代码结构

结论：

> `PILIPILI` 非常适合做 Go 化基底，但需要“保留业务和前端页面，重建后端架构”。

---

## 3. 参考项目 feedsystem_video_go 的可借鉴点

参考仓库：`https://github.com/LeoninCS/feedsystem_video_go`

它的优势不是“功能比你多很多”，而是“后端结构更像完整系统”。最值得借鉴的是：

- 根目录按 `backend + frontend + docker-compose` 组织
- Go 后端分成 `cmd/main.go` 和 `cmd/worker/main.go`
- `internal/account`、`internal/video`、`internal/social`、`internal/feed` 按业务域拆分
- Redis 负责 token、缓存、热榜
- RocketMQ 负责点赞、评论、关注、热度更新等异步任务
- Compose 同时拉起 MySQL、Redis、RocketMQ、Backend、Worker、Frontend

### 3.1 适合直接借鉴的部分

- 项目总体目录结构
- `cmd/api + cmd/worker` 的双进程思路
- 按业务域拆 `account/video/social/feed`
- `JWTAuth / SoftJWTAuth` 的双鉴权思路
- Redis 做 token 与热点缓存
- RocketMQ 作为第二阶段异步化升级方案
- Docker Compose 一键联调

### 3.2 不建议一开始照抄的部分

- 热榜滑动窗口快照分页
- 全链路 MQ 降级直写
- 复杂 Feed 多游标分页
- 一开始就上所有缓存策略
- 一开始就做太细的 Worker 拆分

原因很简单：`PILIPILI` 先要把“视频站的基础闭环”做扎实，再逐步升级成“带 Feed、缓存、异步”的系统。

---

## 4. 推荐的目标架构

建议把项目重构为单仓双端：

```text
pilipili/
  frontend/
  backend/
  deploy/
  docs/
  docker-compose.yml
  README.md
```

### 4.1 前端保留策略

当前 Vue 项目建议整体迁入 `frontend/`，保留：

- 页面路由
- 组件结构
- DPlayer/HLS 播放逻辑
- Element Plus UI
- 现有业务交互流程

前端重点改造内容：

- 替换硬编码接口地址为环境变量
- 提取统一 `api/client`
- 增加登录态拦截器
- 用真实接口替换现有零散请求
- 梳理数据类型和 DTO
- 逐步把动态流接口从 `page + page_size` 收敛为 cursor 分页

### 4.2 后端架构

后端建议采用模块化单体，不建议一开始做微服务：

```text
backend/
  cmd/
    api/
      main.go
    worker/
      main.go
  configs/
    config.yaml
    config.docker.yaml
  internal/
    account/
    video/
    comment/
    social/
    feed/
    search/
    history/
    notice/
    admin/
    media/
    middleware/
    config/
    db/
    http/
    worker/
  pkg/
    response/
    jwt/
    errors/
  migrations/
  Dockerfile
  go.mod
```

### 4.3 技术栈建议

第一阶段：

- Go
- Gin
- GORM
- MySQL
- Redis
- Docker Compose

第二阶段：

- RocketMQ
- Worker
- MinIO 或 OSS
- ffmpeg

第三阶段可选增强：

- Swagger / OpenAPI
- Prometheus 指标
- pprof
- 限流
- 灰度配置

---

## 5. 业务域设计

### 5.1 账户域 account

职责：

- 注册
- 登录
- 退出
- 个人资料
- 密码修改
- token 失效控制

建议接口：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/users/me`
- `PATCH /api/v1/users/me`
- `PATCH /api/v1/users/password`

### 5.2 视频域 video

职责：

- 视频发布
- 视频详情
- 作者视频列表
- 视频状态管理
- 播放地址
- 封面地址

建议接口：

- `POST /api/v1/videos`
- `GET /api/v1/videos/:id`
- `GET /api/v1/videos`
- `GET /api/v1/users/:id/videos`
- `GET /api/v1/videos/:id/play`
- `GET /api/v1/videos/:id/recommendations`

### 5.3 互动域 interaction

职责：

- 视频点赞
- 视频收藏
- 评论
- 评论回复
- 评论点赞

建议接口：

- `POST /api/v1/videos/:id/likes`
- `DELETE /api/v1/videos/:id/likes`
- `POST /api/v1/videos/:id/favorites`
- `DELETE /api/v1/videos/:id/favorites`
- `GET /api/v1/videos/:id/comments`
- `POST /api/v1/videos/:id/comments`
- `POST /api/v1/comments/:id/likes`
- `DELETE /api/v1/comments/:id/likes`
- `GET /api/v1/comments/:id/replies`
- `POST /api/v1/comments/:id/replies`

### 5.4 社交域 social

职责：

- 关注
- 取关
- 粉丝列表
- 关注列表

建议接口：

- `POST /api/v1/users/:id/follow`
- `DELETE /api/v1/users/:id/follow`
- `GET /api/v1/users/:id/followers`
- `GET /api/v1/users/:id/following`

### 5.5 Feed 域

职责：

- 首页推荐流
- 关注流
- 热门榜
- 分区流

建议接口：

- `GET /api/v1/feed/recommend`
- `GET /api/v1/feed/following`
- `GET /api/v1/feed/hot`
- `GET /api/v1/areas/:id/videos`

### 5.6 搜索域

职责：

- 搜视频
- 搜用户

建议接口：

- `GET /api/v1/search/videos`
- `GET /api/v1/search/users`

### 5.7 历史与通知域

职责：

- 浏览历史
- 通知拉取
- 通知已读

建议接口：

- `GET /api/v1/histories`
- `POST /api/v1/histories`
- `GET /api/v1/notices`
- `PATCH /api/v1/notices/:id/read`

### 5.8 管理后台域 admin

职责：

- 待审核视频
- 审核通过
- 审核拒绝
- 统计数据

建议接口：

- `GET /api/v1/admin/videos/pending`
- `POST /api/v1/admin/videos/:id/approve`
- `POST /api/v1/admin/videos/:id/reject`
- `GET /api/v1/admin/stats/today`
- `GET /api/v1/admin/stats/area`

---

## 6. 数据库模型设计

建议先以 MySQL 为主，Redis 做缓存和状态辅助。第一阶段不要把所有统计拆成事件表，先保证主流程稳定。

### 6.1 核心表

#### users

- `id`
- `username`
- `email`
- `password_hash`
- `avatar_url`
- `bio`
- `role`
- `status`
- `created_at`
- `updated_at`

说明：

- `role` 支持 `user/admin`
- `status` 支持 `active/disabled`

#### areas

- `id`
- `name`
- `slug`
- `sort_order`

说明：

- 对应你当前的 `area` 接口和上传分区选择

#### videos

- `id`
- `author_id`
- `area_id`
- `title`
- `description`
- `cover_url`
- `source_path`
- `play_url`
- `duration`
- `status`
- `review_status`
- `review_reason`
- `published_at`
- `created_at`
- `updated_at`
- `like_count`
- `favorite_count`
- `comment_count`
- `view_count`

说明：

- `status` 表示视频是否逻辑可见
- `review_status` 表示 `pending/approved/rejected`
- 第一阶段可以把计数字段放在 `videos` 表里，后续再考虑更复杂的统计模型

#### video_likes

- `id`
- `video_id`
- `user_id`
- `created_at`

唯一索引：

- `(video_id, user_id)`

#### video_favorites

- `id`
- `video_id`
- `user_id`
- `created_at`

唯一索引：

- `(video_id, user_id)`

#### comments

- `id`
- `video_id`
- `user_id`
- `root_id`
- `parent_id`
- `content`
- `reply_count`
- `like_count`
- `status`
- `created_at`
- `updated_at`

说明：

- 一级评论：`root_id = 0, parent_id = 0`
- 回复评论：`root_id = 一级评论 id`
- 这样可以统一评论和回复模型，替代前端现在分开的 `comment` 和 `reply`

#### comment_likes

- `id`
- `comment_id`
- `user_id`
- `created_at`

唯一索引：

- `(comment_id, user_id)`

#### follows

- `id`
- `follower_id`
- `followee_id`
- `created_at`

唯一索引：

- `(follower_id, followee_id)`

#### view_histories

- `id`
- `user_id`
- `video_id`
- `progress_seconds`
- `watched_at`

唯一索引可选：

- `(user_id, video_id)`

#### notices

- `id`
- `user_id`
- `type`
- `title`
- `content`
- `related_id`
- `is_read`
- `created_at`

### 6.2 可选表

#### video_reviews

- `id`
- `video_id`
- `reviewer_id`
- `status`
- `reason`
- `created_at`

说明：

- 如果想保留审核记录，不要只把拒绝理由写在 `videos` 表

#### outbox_events

- `id`
- `event_type`
- `biz_id`
- `payload`
- `status`
- `created_at`

说明：

- 第二阶段可用于本地事务消息/outbox

---

## 7. 媒体链路设计

你当前前端已经默认视频播放是 HLS：

- `video/processed/video-{id}/ts/index.m3u8`

这说明未来 Go 后端最好保留这条思路，但分阶段做。

### 7.1 第一阶段

- 视频文件先落本地磁盘
- 封面也先落本地
- 播放地址可以先返回原始 mp4 或者直接保留现有 HLS 目录结构

### 7.2 第二阶段

- 接入 `ffmpeg`
- 上传后异步转码成 HLS
- 抽取封面
- 播放地址统一走 `/media/...`

### 7.3 第三阶段

- 切换到 `MinIO` 或对象存储
- 后台不再直接依赖本地磁盘

结论：

> 你的前端已经对“视频平台”特征做了铺垫，所以 Go 化时一定要保留媒体链路，而不是只做普通 CRUD。

---

## 8. 鉴权与状态设计

当前前端主要靠：

- `localStorage.isLogin`
- `localStorage.userInfo`

这在 demo 阶段可以接受，但 Go 化后需要正规化。

### 8.1 建议方案

- 登录成功后签发 `access_token + refresh_token`
- 前端持久化短期 `access_token` 和长期 `refresh_token`
- access token 建议短期有效，refresh token 建议轮换刷新
- 第一阶段默认采用单设备 refresh 机制，后续如有需要再扩展多设备 session
- 中间件校验 JWT
- 退出登录/改密时使 token 失效

### 8.2 和 feedsystem_video_go 对齐的点

可以借鉴它的思路：

- 强鉴权接口用 `JWTAuth`
- 首页、视频详情、热榜等接口用软鉴权

这样你就能支持：

- 未登录也能看首页和详情
- 登录后能返回点赞/收藏/关注状态

这是比“所有接口都必须登录”更贴近真实视频站的设计。

---

## 9. 缓存与异步化策略

### 9.1 Redis 第一批要做的事情

第一阶段建议只做这三件事；其中前两项当前已经落地为第一版实现：

- token/session 辅助校验
- 视频详情缓存
- 热门视频排行榜

推荐 key：

- `auth:user:<id>`
- `video:detail:<id>`
- `feed:hot`

当前实现补充：

- `videos/:id` 已接入 Redis 读缓存，命中失败时自动回退数据库
- `feed/hot` 已接入 Redis 短 TTL 缓存
- 点赞、收藏、评论、稿件编辑、源文件/封面上传、管理员审核后会主动失效相关缓存
- 当前仍未实现“写 Redis 再异步回刷 MySQL”，这部分继续保留到 RocketMQ + Worker 阶段

### 9.2 第二阶段再做的事情

- 首页推荐流短 TTL 缓存
- 关注流缓存
- 评论页热点缓存
- 分布式锁防击穿

### 9.3 RocketMQ 的正确上车时机

不要一开始就把所有写操作异步化。推荐顺序：

1. 先同步完成点赞、评论、关注
2. 再把“计数更新、通知创建、热度更新、转码任务”改成异步

第一批异步任务推荐：

- 视频上传后的转码任务
- 点赞后的热度更新
- 评论后的通知创建
- 关注后的通知创建

第二批再做：

- 热榜聚合
- 统计分析
- 站内消息 fanout

### 9.4 Worker 的推荐职责

`cmd/worker` 第一版只要覆盖：

- `MediaWorker`
- `NoticeWorker`
- `PopularityWorker`

不要一开始就拆成太多 worker，先让项目能讲清楚。

---

## 10. 从现有接口到新接口的映射

当前项目前端已经写了很多动作式接口。Go 化时不一定要完全兼容旧接口，但建议做一张映射表，便于前端逐步替换。

### 10.1 账户类

- `login` -> `POST /api/v1/auth/login`
- `auth/register` -> `POST /api/v1/auth/register`
- `auth/check/name` -> `GET /api/v1/auth/check-username`
- `auth/check/email` -> `GET /api/v1/auth/check-email`

### 10.2 视频类

- `videoList` -> `GET /api/v1/videos` 或 `GET /api/v1/feed/recommend`
- `video` -> `GET /api/v1/videos/:id`
- `video/processed/video` -> `GET /api/v1/videos/:id/play`
- `cover/cover` -> 直接由 `cover_url` 返回，不再单独拼 URL

### 10.3 互动类

- `video/like` -> `POST /api/v1/videos/:id/likes`
- `video/dislike` -> `DELETE /api/v1/videos/:id/likes`
- `video/collect` -> `POST /api/v1/videos/:id/favorites`
- `video/cancel` -> `DELETE /api/v1/videos/:id/favorites`
- `comment` -> `GET /api/v1/videos/:id/comments`
- `comment/make` -> `POST /api/v1/videos/:id/comments`
- `reply` -> `GET /api/v1/comments/:id/replies`
- `reply/make` -> `POST /api/v1/comments/:id/replies`

### 10.4 社交类

- `follow` -> `POST /api/v1/users/:id/follow`
- `unfollow` -> `DELETE /api/v1/users/:id/follow`
- `followed` -> `GET /api/v1/users/:id/follow-status`

### 10.5 其他

- `histories` -> `GET /api/v1/histories`
- `notice` -> `GET /api/v1/notices`
- `user/read` -> `PATCH /api/v1/notices/:id/read`
- `search/video` -> `GET /api/v1/search/videos`
- `search/user` -> `GET /api/v1/search/users`

---

## 11. 完整开发流程与里程碑

下面按真正可执行的开发流程拆分。

### 阶段 0：仓库重组与基础设施

目标：

- 给项目建立可长期维护的工程结构

任务：

- 把现有 Vue 工程迁入 `frontend/`
- 新建 `backend/`
- 新建根目录 `docker-compose.yml`
- 统一 `.env`、`README`、`docs/`
- 前端提取 `API_BASE_URL`

交付标准：

- 前端仍能独立启动
- Go 后端有基础目录和空路由
- Compose 能启动 MySQL/Redis

### 阶段 1：后端骨架与鉴权闭环

目标：

- 跑通注册、登录、获取当前用户

任务：

- 初始化 `Gin + GORM + MySQL`
- 建 `users` 表
- 实现注册、登录、用户信息
- 实现 JWT 中间件
- 实现统一响应结构和错误码

交付标准：

- 前端登录注册改用真实 Go 接口
- 登录成功后能访问受保护接口
- 本地可通过 Postman 自测

### 阶段 2：视频站最小闭环

目标：

- 跑通首页、视频详情、评论、点赞、收藏

任务：

- 建 `areas/videos/video_likes/video_favorites/comments/comment_likes`
- 实现视频列表、视频详情、推荐视频
- 实现评论列表、发表评论
- 实现点赞/收藏状态查询和操作
- 完成前端首页、视频详情页对接

交付标准：

- 首页和详情页能完全脱离旧后端运行
- 登录与未登录的展示逻辑都正常
- 评论和点赞能落库

### 阶段 3：个人中心与社交关系

目标：

- 跑通用户页、关注、粉丝、历史

任务：

- 建 `follows/view_histories`
- 实现个人中心接口
- 实现关注/取关
- 实现历史记录读写
- 完成用户主页、历史页、关注列表对接

交付标准：

- 用户中心不再依赖任何旧接口
- 播放视频后能看到历史记录
- 关注关系可正确展示

### 阶段 4：上传、封面、审核

目标：

- 形成“创作者发布 -> 审核 -> 通过后可播放”的闭环

任务：

- 实现视频上传
- 实现封面上传
- 存储视频元数据
- 实现审核状态流转
- 实现后台审核与统计接口

交付标准：

- 上传页可完成真实视频发布
- 后台可审核视频
- 审核通过后首页可见

### 阶段 5：媒体处理与播放器链路

目标：

- 让视频平台像真正的视频平台，而不是普通文件上传系统

任务：

- 接入 `ffmpeg`
- 上传后转码 HLS
- 抽帧生成封面
- 标准化播放地址

交付标准：

- 详情页通过 HLS 地址播放
- 新上传视频能自动生成可播放资源

### 阶段 6：Feed、热榜、缓存

目标：

- 借鉴 `feedsystem_video_go` 的优势，做出系统亮点

任务：

- 首页推荐流接口重构
- 分区流接口
- 热门视频榜
- Redis 缓存视频详情和热点列表
- 未登录与登录用户的软鉴权支持

交付标准：

- 首页和热门页接口稳定
- 热门视频不全靠数据库计算
- 简历里可以写缓存和推荐流

当前状态补充：

- 上述阶段里的“视频详情缓存”和“热门榜缓存”已经完成第一版
- 目前仍是读缓存 + 主动失效，不是 Redis 计数主写模型

### 阶段 7：异步 Worker 与工程增强

目标：

- 从“可用”升级到“像真正后端项目”

任务：

- 接入 RocketMQ
- 视频转码改为异步任务
- 评论/关注通知异步化
- 热度更新异步化
- Docker Compose 拉起前后端和依赖

交付标准：

- `api + worker + mysql + redis + rocketmq + frontend` 可一键启动
- API 进程与 Worker 进程可分开运行
- 简历里可以写事件驱动架构

---

## 12. 具体实现优先级建议

如果时间有限，优先级应当是：

1. 登录注册
2. 首页视频列表
3. 视频详情
4. 评论
5. 点赞/收藏
6. 关注
7. 历史
8. 上传
9. 审核
10. 热门榜
11. Redis
12. RocketMQ
13. ffmpeg 转码

原因：

- 前 1 到 8 项决定项目能否称为“完整视频平台”
- 后 9 到 13 项决定项目是否有“工程亮点”

---

## 13. 对简历最有价值的版本

如果你的目标是 Go 岗，最适合写进简历的版本是：

> 基于 Go + Gin + GORM + MySQL + Redis 构建在线视频与弹幕互动平台，支持用户鉴权、视频发布与播放、评论回复、点赞收藏、关注关系、历史记录与后台审核；通过 Docker Compose 实现一键联调，并使用缓存与异步任务提升热点接口性能和系统可扩展性。

最有价值的关键词组合：

- Go
- Gin
- GORM
- MySQL
- Redis
- JWT
- Docker Compose
- RocketMQ
- Worker
- HLS / ffmpeg
- Feed / 热门榜 / 关注流

---

## 14. 最终建议

### 14.1 这个项目应该如何定位

不要把它定位成：

- Vue 视频网站练手项目
- 单纯的前端项目

应该定位成：

- Go 视频平台后端项目
- 有真实业务域、有媒体链路、有互动系统、有审核后台的完整系统

### 14.2 执行策略

最稳的策略是：

- 保留前端页面和交互
- 新建 Go 后端
- 先完成业务闭环
- 再逐步吸收 `feedsystem_video_go` 的缓存和异步亮点

### 14.3 核心原则

- 先把业务跑通，再做高并发亮点
- 先做模块化单体，再考虑更复杂拆分
- 先做真实后端，再包装简历

---

## 15. 下一步执行建议

在这份蓝图基础上，下一步应该直接进入设计产物输出：

1. 数据库表结构 SQL / GORM Model 设计
2. API 清单与请求响应 DTO
3. 后端目录初始化
4. 第一期开发任务拆解

推荐下一步优先产出：

- `docs/01-contracts/schema.md`
- `docs/01-contracts/api.md`
- `backend/` 初始化脚手架

这样我们就能从“规划”进入“真正开工”。
