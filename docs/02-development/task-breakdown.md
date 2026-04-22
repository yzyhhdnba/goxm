# PILIPILI Go 任务拆解

## 1. 文档目的

本文件用于把 `PILIPILI Go` 的总体规划拆成可执行任务，方便：

- 人类开发者按阶段推进
- AI 代理按边界执行
- 后续把任务迁移到 issue、看板或里程碑

本文件默认服务于当前阶段的真实状态：

- 项目仍处于 Go 化规划与重构阶段
- 当前仓库已经完成 Batch A 与 Batch B 的后端落地
- 任务拆解优先面向 MVP 与主链路落地

---

## 2. 使用方式

建议按下面顺序使用本文件：

1. 先读 `README.md`
2. 再读 `docs/01-contracts/blueprint.md`
3. 再读 `docs/01-contracts/schema.md` 和 `docs/01-contracts/api.md`
4. 最后按本文件逐项推进

任务状态建议：

- `[ ]` 未开始
- `[~]` 进行中
- `[x]` 已完成

如果后续要接入项目管理工具，可直接按章节拆成 Epic 和 Story。

---

## 3. 总体执行原则

- 先做文档和骨架，再做业务，再做增强
- 先跑通主链路，再做热榜、缓存和异步
- 先保留原前端业务语义，再逐步优化接口风格
- 任何重大方向变更先问用户

当前推荐执行顺序：

1. 仓库整理与基础设施
2. 后端初始化与鉴权
3. 首页/详情/评论/点赞/关注主链路
4. 上传/审核/历史/通知
5. 缓存、Feed、Worker、MQ

---

## 4. 里程碑总览

### M0 文档与仓库基线

目标：

- 文档、命名和目录目标稳定

完成标准：

- `README.md`
- `AGENTS.md`
- `docs/01-contracts/blueprint.md`
- `docs/01-contracts/schema.md`
- `docs/01-contracts/api.md`
- `docs/02-development/task-breakdown.md`

全部存在并互相不冲突

### M1 后端骨架可启动

目标：

- `backend/` 初始化完成
- API 服务能本地启动

### M2 用户与主链路可用

目标：

- 注册/登录
- 首页视频流
- 视频详情
- 评论与回复
- 点赞与收藏
- 关注关系

### M3 创作者与后台闭环

目标：

- 上传
- 审核
- 个人中心
- 历史记录
- 通知

### M4 工程增强

目标：

- Redis
- Feed / 热榜
- RocketMQ / Worker
- 媒体转码
- Docker Compose 全链路

---

## 5. M0 文档与仓库基线

### Epic 0.1 文档收口

- [x] 编写 `README.md`
- [x] 编写 `AGENTS.md`
- [x] 编写 Go 化蓝图
- [x] 编写表结构草案
- [x] 编写 API 清单
- [x] 编写任务拆解文档

验收标准：

- 所有文档都明确区分“当前状态”和“未来目标”
- 不存在把未实现能力写成已完成事实的表述

### Epic 0.2 仓库规划收口

- [ ] 明确根目录目标结构
- [ ] 决定 `frontend/` 的迁移策略
- [ ] 决定 `backend/` 的初始化方式
- [ ] 决定 `deploy/` 和 `docs/` 的长期用途

验收标准：

- 目标目录结构能指导后续初始化
- `AGENTS.md` 与 `README.md` 中的目录口径一致

---

## 6. M1 后端骨架可启动

## Epic 1.1 `backend/` 初始化

- [x] 创建 `backend/go.mod`
- [x] 创建 `backend/cmd/api/main.go`
- [x] 创建 `backend/internal/config`
- [x] 创建 `backend/internal/db`
- [x] 创建 `backend/internal/http`
- [x] 创建 `backend/pkg/response`
- [x] 创建 `backend/configs/config.yaml`
- [x] 创建 `backend/configs/config.example.yaml`

前置依赖：

- M0 完成

验收标准：

- `go run ./cmd/api` 可以启动
- 启动后提供健康检查接口

## Epic 1.2 基础配置与环境

- [x] 定义配置结构体
- [x] 接入配置加载逻辑
- [ ] 支持区分本地开发和容器环境
- [x] 设计配置样例
- [x] 增加启动配置强校验

验收标准：

- 启动时能成功读取数据库、Redis、服务端口等配置
- 配置缺失时有清晰错误提示

## Epic 1.3 数据库接入

- [x] 初始化 MySQL 连接
- [x] 定义基础 GORM 配置
- [ ] 决定迁移方式
- [ ] 建立第一批模型注册逻辑

验收标准：

- 服务启动后能连上 MySQL
- 能完成第一批模型迁移或初始化

## Epic 1.4 基础中间件

- [x] 统一响应结构
- [ ] 统一错误码
- [x] 请求日志中间件
- [x] panic recover
- [x] CORS
- [x] 优雅停机

验收标准：

- 成功与失败响应结构统一
- 500 错误不会导致服务直接退出

## Epic 1.5 健康检查与基础路由

- [x] `GET /healthz`
- [x] `GET /api/v1/ping`
- [x] API 路由分组初始化
- [x] `healthz` 覆盖 MySQL 与 Redis 探活

验收标准：

- 本地启动后能通过浏览器或 curl 验证服务在线

---

## 7. M2 用户与主链路可用

## Epic 2.1 用户系统

### Story 2.1.1 表与模型

- [x] 创建 `users` 表模型
- [x] 加入用户名唯一约束
- [x] 加入邮箱唯一约束
- [x] 定义用户 DTO

### Story 2.1.2 注册

- [x] 实现 `POST /api/v1/auth/register`
- [x] 用户名合法性校验
- [x] 邮箱合法性校验
- [x] 密码加密存储

### Story 2.1.3 登录

- [x] 实现 `POST /api/v1/auth/login`
- [x] 校验密码
- [x] 签发 access token
- [x] 签发 refresh token
- [x] 返回用户基础信息

### Story 2.1.4 Refresh Token

- [x] 实现 `POST /api/v1/auth/refresh`
- [x] 校验 refresh token
- [x] 旋转 refresh token
- [x] 返回新的 token 对

### Story 2.1.5 当前用户与登出

- [x] 实现 `GET /api/v1/users/me`
- [x] 实现 `POST /api/v1/auth/logout`
- [x] 设计 token 失效策略

### Story 2.1.6 用户名与邮箱可用性检查

- [x] 实现 `GET /api/v1/auth/check-username`
- [x] 实现 `GET /api/v1/auth/check-email`

验收标准：

- 前端登录/注册弹窗可以切到真实 Go 接口
- 成功登录后可访问受保护资源
- access token 过期后可通过 refresh token 刷新
- refresh 成功后旧 refresh token 立即失效
- refresh 不会立刻吊销已签发 access token
- logout 后当前 access token 和 refresh token 均失效

## Epic 2.2 分区系统

- [x] 创建 `areas` 表
- [x] 写入初始分区数据
- [x] 实现 `GET /api/v1/areas`

验收标准：

- 上传页和首页分区逻辑可以使用真实分区数据

## Epic 2.3 视频系统基础

### Story 2.3.1 表与模型

- [x] 创建 `videos` 表
- [x] 定义视频状态和审核状态枚举
- [x] 定义视频详情 DTO

### Story 2.3.2 首页视频流

- [x] 实现 `GET /api/v1/feed/recommend`
- [x] 使用 `cursor + limit` 分页
- [x] 支持按 `published_at desc, id desc` 倒序

### Story 2.3.3 视频详情

- [x] 实现 `GET /api/v1/videos/:id`
- [x] 返回作者信息
- [x] 返回点赞/收藏/关注 viewer_state

### Story 2.3.4 作者视频列表

- [x] 实现 `GET /api/v1/users/:id/videos`

验收标准：

- 首页和详情页可脱离旧后端运行
- 未登录和已登录状态都能正常返回
- 当前阶段 `viewer_state` 已切到真实关系查询

## Epic 2.4 视频点赞与收藏

### Story 2.4.1 表与模型

- [x] 创建 `video_likes`
- [x] 创建 `video_favorites`

### Story 2.4.2 点赞接口

- [x] 实现 `POST /api/v1/videos/:id/likes`
- [x] 实现 `DELETE /api/v1/videos/:id/likes`
- [x] 实现 `GET /api/v1/videos/:id/likes/me`

### Story 2.4.3 收藏接口

- [x] 实现 `POST /api/v1/videos/:id/favorites`
- [x] 实现 `DELETE /api/v1/videos/:id/favorites`
- [x] 实现 `GET /api/v1/videos/:id/favorites/me`

验收标准：

- 点赞、取消点赞、收藏、取消收藏都能落库
- 计数字段同步更新

## Epic 2.5 评论与回复

### Story 2.5.1 表与模型

- [x] 创建 `comments`
- [x] 创建 `comment_likes`
- [x] 确认统一评论/回复建模

### Story 2.5.2 评论接口

- [x] 实现 `GET /api/v1/videos/:id/comments`
- [x] 实现 `POST /api/v1/videos/:id/comments`

### Story 2.5.3 回复接口

- [x] 实现 `GET /api/v1/comments/:id/replies`
- [x] 实现 `POST /api/v1/comments/:id/replies`

### Story 2.5.4 评论点赞

- [x] 实现 `POST /api/v1/comments/:id/likes`
- [x] 实现 `DELETE /api/v1/comments/:id/likes`
- [x] 实现 `GET /api/v1/comments/:id/likes/me`

验收标准：

- 楼层评论与楼中楼回复都能跑通
- 与原前端 `comment` / `reply` 交互场景能对应

## Epic 2.6 关注关系

- [x] 创建 `follows`
- [x] 实现 `POST /api/v1/users/:id/follow`
- [x] 实现 `DELETE /api/v1/users/:id/follow`
- [x] 实现 `GET /api/v1/users/:id/follow-status`
- [x] 实现 `GET /api/v1/users/:id/followers`
- [x] 实现 `GET /api/v1/users/:id/following`

验收标准：

- 详情页和个人主页能展示关注状态
- 关注与取消关注都能正确生效

---

## 8. M3 创作者与后台闭环

## Epic 3.1 个人中心

- [x] 实现 `GET /api/v1/users/:id/profile`
- [x] 实现 `GET /api/v1/users/me/dashboard`
- [x] 对齐原前端 `getMyIndexInfo` 需求

验收标准：

- 个人主页与“我的”页面具备真实数据来源

## Epic 3.2 历史记录

- [x] 创建 `view_histories`
- [x] 实现 `POST /api/v1/histories`
- [x] 实现 `GET /api/v1/histories`

验收标准：

- 播放视频后能产生浏览历史
- 历史页面可分页展示

## Epic 3.3 通知系统

- [x] 创建 `notices`
- [x] 实现 `GET /api/v1/notices`
- [x] 实现 `PATCH /api/v1/notices/:id/read`
- [x] 审核通过 / 驳回时自动生成创作者通知

验收标准：

- 通知页可读取通知
- 已读状态可持久化

## Epic 3.4 上传系统

### Story 3.4.1 视频元数据

- [x] 实现 `POST /api/v1/videos`
- [x] 支持标题、简介、分区

### Story 3.4.2 视频文件上传

- [x] 实现 `POST /api/v1/videos/:id/source`
- [x] 限制视频类型和大小
- [x] 约定文件存储目录

### Story 3.4.3 封面上传

- [x] 实现 `POST /api/v1/videos/:id/cover`

### Story 3.4.4 编辑稿件

- [x] 实现 `PATCH /api/v1/videos/:id`
- [x] 编辑后重置为 `pending`

验收标准：

- 上传页可创建真实视频记录
- 上传的视频与封面能保存到磁盘或约定目录
- 作者可编辑已有稿件并重新提审

## Epic 3.5 审核与创作者后台

- [x] 创建 `video_reviews`
- [x] 实现 `GET /api/v1/creator/videos`
- [x] 实现 `GET /api/v1/admin/videos/pending`
- [x] 实现 `POST /api/v1/admin/videos/:id/approve`
- [x] 实现 `POST /api/v1/admin/videos/:id/reject`
- [x] 实现 `GET /api/v1/admin/stats/today`
- [x] 实现 `GET /api/v1/admin/stats/area`

验收标准：

- 创作者可查看投稿状态
- 管理员可审核视频
- 后台统计页可展示基础数据

---

## 9. M4 工程增强

## Epic 4.1 前端工程对接

- [x] 迁移原 Vue 工程到 `frontend/`
- [x] 提取统一请求客户端
- [x] 接入环境变量 `API_BASE_URL`
- [x] 接入 `VUE_APP_API_PROXY_TARGET`
- [~] 清理硬编码地址
- [x] 逐步对接新接口

验收标准：

- 不再依赖硬编码 IP
- 新旧接口切换有明确策略

## Epic 4.2 Redis

- [x] 接入 Redis 基础客户端
- [ ] 设计 token 相关 key
- [x] 为视频详情增加缓存
- [x] 为热门榜增加缓存

验收标准：

- 热点详情不完全依赖数据库
- Redis 异常时主链路仍可回退

## Epic 4.3 Feed 与热门榜

- [x] 实现 `GET /api/v1/feed/hot`
- [x] 实现 `GET /api/v1/feed/following`
- [x] 实现 `GET /api/v1/areas/:id/videos`
- [x] 设计基础排序策略

验收标准：

- 首页、关注流、热门榜有明确区分

## Epic 4.4 媒体链路增强

- [ ] 接入 ffmpeg
- [ ] 上传后生成 HLS 播放资源
- [ ] 约定播放目录结构
- [ ] 补充封面抽帧策略

验收标准：

- 播放链路从“静态文件上传”升级到“可标准播放的媒体资源”

## Epic 4.5 RocketMQ 与 Worker

- [x] 初始化 `cmd/worker`
- [ ] 接入 RocketMQ
- [ ] 设计第一批事件
- [ ] 把转码任务异步化
- [ ] 把通知生成异步化
- [ ] 把热视频热度更新异步化

验收标准：

- API 进程和 Worker 可独立运行
- 至少一条异步链路稳定可用

## Epic 4.6 Docker Compose

- [ ] 增加 `mysql`
- [ ] 增加 `redis`
- [ ] 增加 `backend`
- [ ] 增加 `frontend`
- [ ] 第二阶段加入 `rocketmq`
- [ ] 编写启动说明

验收标准：

- 本地一条命令可拉起核心依赖

---

## 10. 当前最优先的下一批任务

如果下一步继续推进编码，建议严格按下面顺序执行：

### Batch A

- [x] 创建 `backend/` 目录结构
- [x] 初始化 `go.mod`
- [x] 添加 `cmd/api/main.go`
- [x] 添加配置加载
- [x] 添加数据库连接
- [x] 添加健康检查

### Batch B

- [x] 创建 `users` 模型
- [x] 实现注册
- [x] 实现登录
- [x] 实现 refresh
- [x] 实现 logout
- [x] 实现用户名与邮箱可用性检查
- [x] 实现 `users/me`
- [x] 实现 JWT 中间件

### Batch C

- [x] 创建 `areas` 与 `videos` 模型
- [x] 实现 `GET /api/v1/areas`
- [x] 实现首页推荐流
- [x] 采用 `cursor + limit` 设计动态 Feed 分页
- [x] 实现视频详情
- [x] 返回 viewer_state

### Batch D

- [x] 创建评论模型
- [x] 实现评论/回复接口
- [x] 创建点赞/收藏表
- [x] 实现点赞/收藏接口
- [x] 创建关注表
- [x] 实现关注接口

到 Batch D 结束时，就已经具备第一个可演示版本。

### Phase 4

- [x] 创建 `search` 业务域
- [x] 创建 `history` 业务域
- [x] 打通上传元数据、源文件、封面三段式链路
- [x] 创建 `admin` 业务域并补上审核与统计
- [x] 在 `frontend/src/api/index.ts` 收口 Search / History / Upload / Management API
- [x] 将 Search / History / Upload / Management 旧页面切换到新 API

到 Phase 4 结束时，前后端已经具备搜索、观看历史、投稿、审核、后台统计的第一版闭环。

---

## 11. 风险与阻塞项

以下是后续最容易卡住的点：

- 前端旧接口和新接口命名不一致
- 上传链路和播放链路需要尽早约定目录结构
- 评论和回复如果不统一建模，后续会产生重复实现
- 认证方案如果过晚确定，会影响前端接入方式
- 若过早引入 MQ，会延迟 MVP 成型

处理原则：

- 先解主链路阻塞
- 复杂优化延后
- 领域建模优先于局部性能优化

---

## 12. 完成定义

### MVP 完成定义

达到以下条件时，可认为 MVP 完成：

- 用户可注册登录
- 首页可看视频列表
- 详情页可看视频信息
- 可评论、回复、点赞、收藏
- 可关注作者
- 个人页能展示基础信息

### V1 完成定义

达到以下条件时，可认为第一版完整：

- 支持上传和审核
- 支持历史记录和通知
- 前后端可联调
- Compose 可启动主要依赖

### V2 完成定义

达到以下条件时，可认为具备工程亮点：

- Redis 接入完成
- 热门榜和 Feed 具备缓存
- Worker 和 MQ 至少承担一条稳定异步链路
- 媒体处理链路具备 HLS 或等价能力

---

## 13. 后续维护建议

随着代码逐步落地，本文件应同步维护：

- 已完成任务打勾
- 新增任务放到对应里程碑
- 被放弃的任务注明原因
- 涉及重大方向变化的任务，先由用户拍板

本文件不应变成“纯展示文档”，它应该始终保持可执行性。
