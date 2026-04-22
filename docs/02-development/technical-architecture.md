# PILIPILI Go 技术设计文档

## 1. 文档定位

本文件用于单独整理 `goxm` 当前阶段的技术设计，服务于以下目标：

- 在契约文档之外，补充“系统如何组织与落地”的技术视角
- 吸收参考项目 `feedsystem_video_go` 的文档结构优点，但只保留适合当前仓库阶段的部分
- 明确区分“当前已经实现的技术事实”和“后续规划中的工程增强”
- 为后续开发、联调、面试表达与代理接手提供统一入口

使用建议：

1. 先读 `README.md`
2. 再读 `docs/01-contracts/blueprint.md`
3. 再读 `docs/01-contracts/schema.md`
4. 再读 `docs/01-contracts/api.md`
5. 最后读本文件

本文件不替代契约文档。若本文件与契约文档冲突，应以契约文档为准，并回头核对实现是否偏离。

---

## 2. 项目技术定位

`PILIPILI Go` 的技术定位不是“从零造一个全新视频站”，而是：

- 以现有 Vue 前端业务原型为产品基底
- 以 Go 模块化单体为后端主实现
- 先打通视频平台主链路，再逐步接入缓存、异步任务与媒体增强

当前技术路线仍然遵守仓库既定边界：

- 第一阶段主栈：Go、Gin、GORM、MySQL、Redis、Docker Compose
- 第二阶段增强：RocketMQ、Worker、ffmpeg、MinIO 或 OSS
- 前端继续保留 Vue，不改造成 React，不推翻现有页面交互结构

这意味着当前技术设计优先解决三个问题：

1. 如何让现有前端页面稳定命中真实 Go API
2. 如何把零散业务收敛到清晰的后端领域边界
3. 如何为后续缓存、异步和媒体链路预留演进空间

---

## 3. 当前技术现状摘要

截至当前仓库状态，已经落地的技术事实包括：

- 根目录采用 `frontend/ + backend/ + docs/ + docker-compose.yml` 的单仓布局
- `backend/` 已具备可运行的 API 入口与基础模块拆分
- MySQL 与 Redis 已通过根目录 `docker-compose.yml` 支持本地联调
- Redis 已从“客户端初始化 + 探活”扩展到视频详情与热门榜读缓存
- RocketMQ 已补齐配置结构、`cmd/worker` 骨架与可选 Compose profile
- API 已接入配置加载、依赖初始化、自动迁移、默认分区种子与优雅停机
- 路由层已接入统一响应、日志、Recovery、CORS、中间件鉴权与静态媒体暴露
- Batch B 到 Phase 4 的主链路已经具备实现与路由级集成测试覆盖
- 前端主目录已经收敛到根级 `frontend/`，并通过 `frontend/src/api/index.ts` 统一 API 出口

当前尚未落地、但已经作为技术规划存在的部分包括：

- `cmd/worker` 业务消费逻辑
- RocketMQ 事件生产与消费的完整接线
- ffmpeg 转码与 HLS 播放链路
- 对象存储替代本地落盘
- 推模式 Feed、实时推送与实时弹幕链路
- 搜索与推荐的生产级基础设施

因此，本项目当前更准确的描述是：

> 一个已经具备真实后端主链路、但仍处于渐进增强阶段的 Go 视频平台项目。

---

## 4. 技术栈分层

### 4.1 当前已落地技术栈

后端：

- Go
- Gin
- GORM
- MySQL
- Redis
- JWT
- YAML 配置加载

前端：

- Vue
- Vue Router
- Vuex
- Axios 封装请求层

联调与运行：

- Docker Compose
- 本地文件存储
- 路由级集成测试
- 脚本化 E2E 验证

### 4.2 已规划但未落地的增强栈

- RocketMQ：承接点赞、评论、关注、热度更新等异步任务
- Worker：承接回刷、转码、通知扩展等后台消费能力
- ffmpeg：承接视频处理和 HLS 切片
- MinIO / OSS：承接媒体文件对象存储
- FULLTEXT / Elasticsearch：承接搜索升级

这里的关键约束是：

> 规划中的技术栈只能写成“后续演进方向”，不能写成“当前已实现能力”。

---

## 5. 总体架构设计

当前推荐并实际收敛中的仓库结构如下：

```text
frontend/
backend/
deploy/
docs/
README.md
AGENTS.md
docker-compose.yml
```

其中：

- `frontend/` 是当前主前端目录，承接所有新联调与页面改造
- `backend/` 是 Go 后端主工作区
- `deploy/` 承接数据库初始化脚本与后续部署材料
- `docs/` 承接契约、技术设计、开发资料与面试资料
- `backend/frontend/` 仍然存在，但只作为历史遗留副本，不再作为默认新增改动位置

整体架构采用“前后端分离 + 模块化单体”的方式：

1. Vue 前端负责页面展示与交互编排
2. Go API 负责鉴权、业务规则、数据读写与媒体路径暴露
3. MySQL 负责主数据存储
4. Redis 当前负责基础设施接入、探活，以及视频详情 / 热门榜读缓存
5. RocketMQ 当前已补齐配置入口、健康检查接入点与 Worker 代码骨架
6. Docker Compose 负责本地依赖编排

这套结构借鉴了 `feedsystem_video_go` 的单仓双端思路，但没有照搬其更重的异步与缓存层实现。

---

## 6. 后端系统架构

### 6.1 启动与生命周期

当前 API 启动链路收敛在 `backend/cmd/api/main.go`，主流程如下：

1. 读取配置文件并校验必填项
2. 初始化 MySQL 连接与连接池
3. 初始化 Redis 客户端并执行 Ping
4. 对各领域模型执行 AutoMigrate
5. 执行默认分区 Seed
6. 组装 Gin Router、业务 Handler 与中间件
7. 启动 `http.Server`
8. 通过 `signal.NotifyContext` 接入优雅停机

这条链路体现了当前后端的两个核心取向：

- `Fail Fast`：数据库、Redis、JWT 密钥等关键配置异常时直接拒绝启动
- `Graceful Shutdown`：退出时优先等待在飞请求完成，再强制关闭

### 6.2 运行时分层

当前后端主要按以下层次组织：

```text
backend/
  cmd/api
  internal/account
  internal/admin
  internal/area
  internal/auth
  internal/comment
  internal/config
  internal/db
  internal/history
  internal/http
  internal/media
  internal/middleware
  internal/notice
  internal/redis
  internal/search
  internal/social
  internal/video
  pkg/authctx
  pkg/request
  pkg/response
```

职责划分如下：

- `cmd/api`：应用入口与生命周期管理
- `internal/config`：配置结构体、默认值与启动期校验
- `internal/db`：数据库连接与连接池设置
- `internal/redis`：Redis 客户端初始化、探活与底层访问封装
- `internal/http`：路由装配、全局中间件与健康检查
- `internal/auth`：JWT token 管理
- `internal/middleware/auth`：强鉴权与软鉴权中间件
- `internal/media`：本地媒体存储抽象
- 各业务域目录：承接 Handler、Service、Repository 与实体定义
- `pkg/authctx`：统一当前用户上下文
- `pkg/request`：统一页码分页参数解析
- `pkg/response`：统一响应 envelope

### 6.3 模块化单体的取舍

当前项目没有采用微服务，原因不是“做不到”，而是：

- 当前阶段优先级是打通主链路，不是拆分部署拓扑
- 业务规模还不足以支撑拆服务带来的复杂度
- Vue 前端仍处在快速联调阶段，单体 API 更利于迭代

因此当前的技术设计是：

> 先把模块边界拆清楚，再在必要时把热点链路抽成 Worker 或异步任务，而不是反过来先拆服务。

---

## 7. 核心领域模块设计

### 7.1 Account 域

负责：

- 注册
- 登录
- 刷新 token
- 退出登录
- 当前用户信息
- 用户名 / 邮箱可用性检查

当前设计要点：

- 使用 `access token + refresh token`
- `refresh token` 第一阶段采用单设备可轮换策略
- 失效控制基于 `token_version + refresh_token_hash`
- `logout` 同时提升 `token_version` 并清空 `refresh_token_hash`
- `refresh` 更新时必须带旧 hash 条件，避免并发续签产生竞态

### 7.2 Video / Feed 域

负责：

- 首页推荐流
- 热门榜
- 关注流
- 分区流
- 视频详情
- 作者视频列表
- 点赞与收藏
- 投稿创建、编辑与稿件列表
- 源视频上传与封面上传

当前设计要点：

- `feed/recommend`、`feed/hot`、`feed/following`、`areas/:id/videos` 统一采用 `cursor + limit`
- 热门榜游标使用 `hot_score:unix_timestamp:video_id`
- 其他动态流游标使用 `unix_timestamp:video_id`
- 公共可见视频统一约束为 `visible + approved + published_at 非空`
- `users/:id/videos` 仍保留 `page + page_size`
- 当前播放器优先支持直接播放落盘的 mp4 文件，HLS 仍在后续规划中

### 7.3 Comment 域

负责：

- 一级评论
- 回复
- 评论点赞

当前设计要点：

- 评论与回复统一建模在同一张 `comments` 表
- 一级评论使用 `root_id = 0`、`parent_id = 0`
- 回复使用 `root_id = 一级评论 ID`
- 评论列表与回复列表维持 `page + page_size`
- 当前查询已收敛为单次连表查询，不再反复回查主表

### 7.4 Social 域

负责：

- 关注
- 取关
- 关注状态
- 粉丝列表
- 关注列表

当前设计要点：

- 关系表通过唯一索引约束重复关注
- 关注与取关链路同步维护聚合计数字段
- `videos/:id` 与 `users/:id/profile` 通过可选鉴权返回真实关注状态

### 7.5 Search / History / Notice / Admin 域

Search：

- 提供视频搜索与用户搜索
- 当前基于 MySQL `LIKE`，属于阶段性实现

History：

- 提供播放历史写入与分页查询
- 接口需要登录

Notice：

- 提供站内通知列表与已读能力
- 当前审核通过 / 驳回时会写入创作者通知

Admin：

- 提供待审视频列表、已审列表、审核动作与统计面板
- 当前只允许 `role=admin` 调用
- 统计接口已覆盖按日与按分区两个视角

---

## 8. 数据与存储设计

当前数据设计仍以 `docs/01-contracts/schema.md` 为准，技术层面可概括为以下几点：

### 8.1 MySQL 为第一阶段主存储

当前主表包括：

- `users`
- `areas`
- `videos`
- `video_likes`
- `video_favorites`
- `comments`
- `comment_likes`
- `follows`
- `view_histories`
- `notices`
- `video_reviews`

### 8.2 业务状态与软删除分离

统一采用：

- `status` 承接业务可见性与审核语义
- `deleted_at` 承接 GORM 软删除

这样可以避免把“前台隐藏”“审核驳回”“技术删除”混为一类。

### 8.3 计数查询先保留主表字段

当前视频、用户、评论等聚合计数字段仍保留在主表，用于保证：

- 详情页查询简单
- 列表页展示稳定
- 前端联调不依赖额外聚合服务

长期目标仍然是：

- 写路径优先进入 Redis
- 再由 RocketMQ + Worker 异步链路批量回刷 MySQL

但这部分目前仍属于后续增强，不应在文档中写成既成事实。

### 8.4 Redis 当前已落地的缓存策略

当前已经落地的缓存能力包括：

- `video:detail:<id>`：视频详情读缓存
- `feed:hot`：热门榜短 TTL 读缓存

当前实现策略是：

- 读路径优先查 Redis，未命中再回退 MySQL
- Redis 异常时直接走数据库主链路，不阻断请求
- 点赞、收藏、评论、视频编辑、源文件/封面上传、管理员审核后，会主动失效视频详情与热门榜缓存

当前仍未落地的部分包括：

- 推荐流、关注流、分区流缓存
- 关注关系状态缓存
- “写 Redis，再异步回刷 MySQL”的计数主写模型

### 8.5 媒体文件当前采用本地落盘

当前视频源文件与封面文件默认写入：

- `backend/storage/`

对外访问通过 Gin 静态路由暴露：

- `/uploads/...`

这是当前阶段最务实的实现方式，便于本地联调；对象存储与转码链路属于后续演进。

---

## 9. 鉴权与安全设计

### 9.1 认证模型

当前所有需要登录的接口统一使用：

- `Authorization: Bearer <access_token>`

刷新使用：

- `POST /api/v1/auth/refresh`

### 9.2 中间件策略

当前后端中间件分两类：

- 强鉴权：必须登录才能访问
- 软鉴权：匿名可访问，但已登录时会补全 `viewer_state`

这一设计主要服务于：

- 视频详情可匿名浏览
- 已登录用户可看到是否点赞、收藏、关注
- 用户主页可在匿名和登录态下返回不同视图信息

### 9.3 启动期安全校验

配置加载阶段当前会强制校验：

- 数据库 DSN
- Redis 地址
- JWT 密钥
- 基础服务端口与 TTL

目的不是“形式上的严格”，而是避免服务带着空配置或占位配置启动成功。

---

## 10. 前后端协作设计

### 10.1 前端协作边界

当前主前端目录是根级 `frontend/`，后续新增联调与 API 改造默认都应落在这里。

`backend/frontend/` 仍然存在，但只视作历史副本，不应继续写入新功能。

### 10.2 API 统一出口

当前前端核心接口统一收口到：

- `frontend/src/api/index.ts`

后续页面改造应继续遵守这一约束，不要回到页面内散落硬编码 Axios 的方式。

### 10.3 请求与环境约定

当前前端请求约定包括：

- API 地址通过 `frontend/.env.*` 管理
- 使用 `VUE_APP_API_BASE_URL` / `VUE_APP_API_PROXY_TARGET`
- `frontend/src/utils/request.ts` 直接返回真实 `data`
- 新代码不要继续按旧风格读取 `res.data`

### 10.4 登录态判断约定

评论、回复、点赞、关注等前端交互当前统一以：

- `access_token`

作为登录态判断依据，而不是依赖 `userInfo` 是否存在。

---

## 11. 联调、测试与运行方式

### 11.1 本地依赖

根目录 `docker-compose.yml` 当前默认编排两类核心依赖：

- MySQL
- Redis

同时提供可选的 `mq` profile：

- RocketMQ NameServer
- RocketMQ Broker

当前 `deploy/rocketmq/broker.conf` 默认把 `brokerIP1` 固定为 `127.0.0.1`，目的是优先支持“后端进程跑在宿主机、本地直接连 MQ”的联调方式；若后续 API / Worker 容器化，需要再按容器网络调整该配置。

这是一个刻意保持克制的选择，目的是保证：

- 主链路联调成本低
- 本地环境更容易拉起
- 文档与当前阶段事实保持一致

### 11.2 测试覆盖现状

当前已具备的验证方式包括：

- 路由级集成测试
- `backend/scripts/` 下的 E2E 脚本
- 前后端联调验证

重点已覆盖：

- 鉴权主链路
- Feed / 视频详情
- 评论、回复、点赞、收藏、关注
- Search / History / Notice
- Creator / Admin

### 11.3 当前运行注意事项

- 前端构建在本机 Node 环境下需要 `NODE_OPTIONS=--openssl-legacy-provider`
- 旧视频标题若出现乱码，优先判断为历史脏数据问题
- 不应因为历史脏数据现象直接回退当前 `utf8mb4` 配置

---

## 12. 与参考项目的关系

参考文档：

- [feedsystem_video_go 项目设计](https://github.com/LeoninCS/feedsystem_video_go/blob/main/feedsystem_video_go%E9%A1%B9%E7%9B%AE%E8%AE%BE%E8%AE%A1.md)

本仓库从中主要吸收了这些思路：

- 以“技术设计说明”补足契约文档之外的架构视角
- 采用 `frontend + backend + docker-compose` 的单仓双端组织
- 优先按领域拆分后端模块，而不是堆在一个大目录中
- 把 Redis、Worker、RocketMQ 写成阶段化演进，而不是第一天全部上齐

同时也刻意保持以下差异：

- 不直接照抄参考项目的异步复杂度
- 不把尚未落地的 Worker、MQ、转码写成当前事实
- 继续围绕 `PILIPILI` 现有 Vue 前端业务结构推进

---

## 13. 后续演进建议

在不偏离当前契约的前提下，后续技术推进建议按以下顺序进行：

1. 继续补齐统一错误码、迁移策略与配置区分
2. 把 Redis 正式接入计数与热点缓存写路径
3. 基于 `RocketMQ + cmd/worker` 承接异步回刷与媒体处理
4. 为上传链路接入 ffmpeg 转码与 HLS 播放
5. 将本地存储逐步演进到对象存储
6. 将搜索从 `LIKE` 升级到 FULLTEXT 或 Elasticsearch

所有这些动作都应继续遵守一句话原则：

> 先保证视频平台主链路完整可跑，再引入更重的工程能力。
