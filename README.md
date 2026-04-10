# PILIPILI Go

## 项目简介

`PILIPILI Go` 是一个处于 **Go 化规划与重构阶段** 的在线视频与互动平台项目。

它的起点不是从零开始搭建，而是基于已有的 `PILIPILI` 前端业务原型，逐步补齐 Go 后端、数据模型、接口规范、媒体链路与工程化能力，最终演进为一个适合学习、展示和求职表达的完整后端项目。

当前阶段的目标不是“宣称项目已经完成”，而是先把：

- 项目定位
- 目标架构
- 数据库模型
- API 契约
- 开发顺序

这些基础设计定清楚，再进入真正的后端开发与前后端联调。

---

## 当前状态

当前仓库已经从纯文档筹备期进入 **基础后端骨架已落地、主链路继续扩展阶段**。

现阶段仓库包含的内容主要是：

- Go 化总体蓝图
- 数据库表结构草案
- API 设计清单
- `backend/` 基础骨架
- 用户鉴权模块的第一版实现
- `areas` 与 `videos` 的第一版实现
- 根目录 `docker-compose.yml`
- MySQL 初始化脚本

当前已经可以完成下面这条最小运行链路：

- 使用 `docker-compose` 启动 `mysql` 与 `redis`
- 本地启动 Go API
- 访问 `/healthz` 与 `/api/v1/ping`
- 调用 `register -> login -> refresh -> users/me -> logout`
- 访问 `GET /api/v1/areas`
- 访问 `GET /api/v1/feed/recommend`
- 访问 `GET /api/v1/feed/hot`
- 登录后访问 `GET /api/v1/feed/following`
- 访问 `GET /api/v1/areas/:id/videos`
- 访问 `GET /api/v1/videos/:id`
- 访问 `GET /api/v1/users/:id/videos`
- 登录后访问 `GET /api/v1/search/videos`
- 登录后访问 `GET /api/v1/histories`
- 登录后访问 `POST /api/v1/videos`
- 登录后访问 `GET /api/v1/creator/videos`
- 管理员访问 `GET /api/v1/admin/videos/pending`

也就是说，这个仓库现在已经不只是设计中台，而是进入了“**设计完成 + 基础骨架可运行**”的阶段。

当前已实现的 Batch B 能力：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/check-username`
- `GET /api/v1/auth/check-email`
- `GET /api/v1/users/me`

当前 Batch B 的鉴权语义已经明确为：

- 采用 `access token + refresh token`
- `refresh` 成功后旧 `refresh token` 立即失效
- 已签发的 `access token` 不会因 `refresh` 立刻失效，而是在自然过期或 `logout` 后失效
- `logout` 通过 `token_version + refresh_token_hash` 统一吊销当前会话
- `refresh` 落库时携带旧 `refresh_token_hash` 条件，避免并发续签穿透

当前已实现的 Batch C 能力：

- `GET /api/v1/areas`
- `GET /api/v1/feed/recommend`
- `GET /api/v1/videos/:id`
- `areas` 默认分区种子写入
- `viewer_state` 已接入软鉴权返回结构

当前已实现的 Batch D 能力：

- `POST /api/v1/videos/:id/likes`
- `DELETE /api/v1/videos/:id/likes`
- `GET /api/v1/videos/:id/likes/me`
- `POST /api/v1/videos/:id/favorites`
- `DELETE /api/v1/videos/:id/favorites`
- `GET /api/v1/videos/:id/favorites/me`
- `GET /api/v1/videos/:id/comments`
- `POST /api/v1/videos/:id/comments`
- `GET /api/v1/comments/:id/replies`
- `POST /api/v1/comments/:id/replies`
- `POST /api/v1/comments/:id/likes`
- `DELETE /api/v1/comments/:id/likes`
- `GET /api/v1/comments/:id/likes/me`
- `POST /api/v1/users/:id/follow`
- `DELETE /api/v1/users/:id/follow`
- `GET /api/v1/users/:id/follow-status`
- `GET /api/v1/users/:id/followers`
- `GET /api/v1/users/:id/following`
- `videos/:id` 的 `viewer_state` 已接入真实点赞、收藏、关注状态查询

当前基础设施新增能力：

- 启动时同时连接 MySQL 与 Redis
- `/healthz` 同时探测 MySQL 与 Redis
- 主程序已接入优雅停机
- 配置加载阶段会强校验数据库、Redis 和 JWT 密钥
- 已接入开发期 CORS 中间件，支持前后端分离联调

当前已实现的 Feed / 作者页扩展能力：

- `GET /api/v1/feed/hot`
- `GET /api/v1/feed/following`
- `GET /api/v1/areas/:id/videos`
- `GET /api/v1/users/:id/videos`
- `feed/hot` 当前按 `hot_score desc, published_at desc, id desc` 排序
- `feed/following` 当前按 `published_at desc, id desc` 返回已关注作者的公开视频
- 分区流当前已实现，默认只支持 `sort=latest`
- 分区流 `next_cursor` 当前复用 `unix_timestamp:video_id`
- 作者视频列表当前保留 `page + page_size`
- 热门榜、关注流、分区流都已经补上测试覆盖

当前已实现的 Phase 4 / 创作者后台能力：

- `GET /api/v1/search/videos`
- `GET /api/v1/search/users`
- `POST /api/v1/histories`
- `GET /api/v1/histories`
- `GET /api/v1/notices`
- `PATCH /api/v1/notices/:id/read`
- `GET /api/v1/users/:id/profile`
- `GET /api/v1/users/me/dashboard`
- `POST /api/v1/videos`
- `PATCH /api/v1/videos/:id`
- `POST /api/v1/videos/:id/source`
- `POST /api/v1/videos/:id/cover`
- `GET /api/v1/creator/videos`
- `GET /api/v1/admin/videos/pending`
- `GET /api/v1/admin/videos?review_status=reviewed`
- `POST /api/v1/admin/videos/:id/approve`
- `POST /api/v1/admin/videos/:id/reject`
- `GET /api/v1/admin/stats/today`
- `GET /api/v1/admin/stats/area`
- 本地媒体文件当前默认落盘到 `backend/storage/`
- API 当前通过 Gin 静态路由暴露 `/uploads/...`
- 管理员审核通过 / 驳回后，当前会自动向创作者写入站内通知
- 编辑稿件当前已经打通，作者更新标题 / 分区 / 简介后会自动回到 `pending` 待审状态
- 个人中心、通知页、上传编辑弹窗、创作者列表、后台审核页、搜索页、历史页都已切到根目录 `frontend/src/api/index.ts` 统一 API 层
- 视频详情页的评论 / 回复前端当前已切到统一 request 返回结构，不再依赖旧版 `res.data`
- 评论、回复、评论点赞等前端登录态判断当前统一以 `access_token` 为准
- 头部菜单、管理页背景、管理导航头像等残留老资源当前已改为本地常量和本地静态资源，不再依赖 `101.35.142.191`
- 前端当前已支持 `VUE_APP_API_BASE_URL` 与 `VUE_APP_API_PROXY_TARGET` 环境变量
- 路由级集成测试当前已经覆盖 Search / History / Notice / Personal Dashboard / Upload / Admin / Edit 主链路
- 搜索第一阶段仍基于 MySQL `LIKE`，后续数据规模上来后需要升级为 FULLTEXT 或 Elasticsearch

---

## 项目来源

`PILIPILI Go` 的业务基底来自原始前端项目：

- 原项目仓库：[PILIPILI](https://gitee.com/yzynba/pilipili)

原项目已经具备比较完整的视频站前端业务雏形，包括：

- 首页视频展示
- 视频详情页
- 评论与回复
- 点赞与收藏
- 用户登录注册
- 搜索
- 历史记录
- 通知
- 上传中心
- 后台审核管理

这意味着 `PILIPILI Go` 的改造重点不是重新设计一个全新产品，而是：

- 保留原有业务语义和页面流程
- 用 Go 重建后端边界
- 把零散接口升级为稳定的系统设计

---

## 参考项目

本项目在后端结构设计上会参考以下 Go 项目：

- [feedsystem_video_go](https://github.com/LeoninCS/feedsystem_video_go)

主要借鉴方向包括：

- `backend + frontend + docker-compose` 的单仓组织方式
- `cmd/api + cmd/worker` 的双进程思路
- 按业务域拆分 `account / video / social / feed`
- `JWT + Redis` 的登录态设计
- `Redis + MQ + Worker` 的渐进式工程增强

但本项目不会直接照搬参考仓库的全部实现，而会优先围绕 `PILIPILI` 自身的业务结构推进。

---

## 项目目标

`PILIPILI Go` 的目标形态是一个基于 Go 的在线视频与互动平台，支持：

- 用户注册、登录、退出与个人资料管理
- 视频发布、播放、详情、推荐与分区浏览
- 评论、回复、点赞、收藏
- 用户关注、粉丝与个人主页
- 历史记录与站内通知
- 视频上传、封面上传与审核流转
- 热门推荐、关注流与基础缓存优化
- 容器化本地联调与后续 Worker 异步化扩展

---

## 推荐技术栈

第一阶段推荐技术栈：

- Go
- Gin
- GORM
- MySQL
- Redis
- Docker Compose

第二阶段规划增强：

- RabbitMQ
- Worker
- ffmpeg
- MinIO 或对象存储

当前这些能力仍属于规划范围，不应被表述为“已经全部完成”。

---

## 文档索引

当前仓库的核心文档如下：

- 文档总览：[docs/README.md](/Users/hhd/Desktop/test/goxm/docs/README.md)
- 总体蓝图：[blueprint.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/blueprint.md)
- 表结构草案：[schema.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/schema.md)
- API 清单：[api.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/api.md)
- 任务拆解：[task-breakdown.md](/Users/hhd/Desktop/test/goxm/docs/02-development/task-breakdown.md)
- 协作约束：[AGENTS.md](/Users/hhd/Desktop/test/goxm/AGENTS.md)

建议阅读顺序：

1. `README.md`
2. `docs/01-contracts/blueprint.md`
3. `docs/01-contracts/schema.md`
4. `docs/01-contracts/api.md`
5. `docs/02-development/task-breakdown.md`
6. `AGENTS.md`

---

## 仓库整理说明

- 根目录历史遗留的一次性前端 JS 补丁脚本，已统一迁移到 `frontend/scripts/legacy-patches/`
- 这些脚本仅用于阶段性批量改写，不属于日常构建与运行链路
- 后续如需修改前端逻辑，优先直接改动源码，避免重复执行旧补丁脚本

---

## 当前可运行内容

当前阶段已经验证通过的最小启动方式：

1. 启动依赖

```bash
docker-compose up -d mysql redis
```

2. 本地启动 API

```bash
cd backend
GOCACHE=/tmp/pilipili-go-build-cache GOMODCACHE=/tmp/pilipili-go-mod-cache go run ./cmd/api -config configs/config.yaml
```

3. 前端构建校验

在当前本机 Node 版本下，Vue CLI 构建需要加上 OpenSSL 兼容参数：

```bash
cd frontend
NODE_OPTIONS=--openssl-legacy-provider npm run build
```

4. 验证接口

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/api/v1/ping
curl http://127.0.0.1:8080/api/v1/areas
curl "http://127.0.0.1:8080/api/v1/feed/recommend?limit=8"
```

当前已验证返回：

- `/healthz` 返回 `code=0`
- `/api/v1/ping` 返回 `pong`
- 注册、登录、刷新、登出和 `users/me` 链路已实测通过
- `areas`、`feed/recommend`、`videos/:id` 已由路由级集成测试覆盖
- 点赞、收藏、评论、回复、评论点赞、关注与 `viewer_state` 已由路由级集成测试覆盖
- Search、History、Upload、Creator、Admin 已由路由级集成测试覆盖
- CORS 预检已由路由级集成测试覆盖
- 浏览器自动化已验证管理后台、搜索、历史、上传、视频详情主链路可命中当前 Go API 并返回 `200`
- 视频详情页“发布评论”按钮本轮已修复前端触发链路，评论请求会直接命中 `POST /api/v1/videos/:id/comments`
- `frontend/` 当前已能完成生产构建，但在当前 Node 环境下需附带 `NODE_OPTIONS=--openssl-legacy-provider`
- 当前前端生产构建已通过，但仍保留既有 CSS 顺序与包体积 warning，属于后续体验优化项
- 如果数据库里还没有已审核且已发布的视频，`feed/recommend` 会返回空列表

当前联调备注：

- 历史测试数据里仍可能存在个别标题乱码；当前更接近旧数据污染问题，而不是现有 Go API / MySQL `utf8mb4` 配置问题
- 本轮已清理关键页面残留的外部静态资源硬编码，因此浏览器中的 `502 / ERR_FAILED` 旧资源报错已不再是主链路阻塞项

---

## 目标架构

目标仓库结构规划如下：

```text
pilipili-go/
  frontend/
  backend/
  deploy/
  docs/
  docker-compose.yml
  README.md
  AGENTS.md
```

其中：

- `frontend/`：承接原有 Vue 前端页面和交互
- `backend/`：新增 Go 后端
- `deploy/`：部署脚本、Compose、环境样例
- `docs/`：后续可以沉淀架构图、流程图和接口说明
- `backend/frontend/`：当前仍保留历史副本，后续联调与新增改动默认以根目录 `frontend/` 为准

当前工作区已经具备根目录 `frontend/`、`backend/`、根目录 `docker-compose.yml` 和 MySQL 初始化脚本，后续重点是继续沿既定文档收敛业务模块与工程细节，而不是重复迁移前端目录。

---

## 开发路线

### 阶段 0：仓库与文档基线

- 明确项目定位
- 固化表结构和 API 设计
- 建立后续协作规范

### 阶段 1：后端骨架与鉴权

- 初始化 Go 后端
- 完成注册、登录、用户信息
- 统一响应和中间件
- 接入 Redis 客户端、健康检查与优雅停机

### 阶段 2：视频站最小闭环

- 首页视频流
- 视频详情
- 评论与回复
- 点赞与收藏
- 关注关系

### 阶段 3：创作者与后台

- 上传中心
- 视频审核
- 历史记录
- 通知系统

### 阶段 4：工程增强

- Redis 缓存
- Feed / 热门榜
- RabbitMQ + Worker
- 媒体转码链路
- Docker Compose 全链路联调

---

## 当前仓库内容

截至当前阶段，仓库内主要包含：

- 规划文档
- 接口设计文档
- 表结构设计文档
- `frontend/` Vue 前端工程
- `backend/` 基础骨架
- `account / video / comment / social / feed / search / history / notice / admin / media` 第一版实现
- 根目录 `docker-compose.yml`
- MySQL 初始化脚本

尚未包含：

- `cmd/worker` 的实际异步消费实现
- MQ / ffmpeg / 对象存储 等第二阶段工程增强
- 搜索与推荐的生产级基础设施
- 生产级部署配置

因此，如果你是第一次进入这个仓库，建议先读文档统一设计口径，再按 README 中的方式启动当前已经可运行的 API。

---

## 适合谁阅读这个仓库

这份仓库当前适合：

- 想把前端视频站项目升级为 Go 后端项目的人
- 想系统化梳理项目设计的人
- 想为简历准备一个更完整后端项目的人
- 想和 AI / 协作者一起推进项目重构的人

---

## 接下来会做什么

后续会按如下顺序继续推进：

1. 持续维护并校准 `README.md`、`AGENTS.md` 与设计文档
2. 继续清理前端遗留的旧数据形态与页面交互细节
3. 推进搜索能力升级、媒体链路增强与 Worker 化拆分
4. 再继续完善缓存、转码与对象存储方案
5. 视需要扩展分区流和创作者页的更多排序策略

---

## 说明

本仓库当前强调的是：

- **真实状态**
- **清晰规划**
- **渐进重构**

而不是“把未来设想写成既成事实”。

如果你希望直接开始实现，请优先阅读：

- [docs/README.md](/Users/hhd/Desktop/test/goxm/docs/README.md)
- [blueprint.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/blueprint.md)
- [schema.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/schema.md)
- [api.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/api.md)
- [task-breakdown.md](/Users/hhd/Desktop/test/goxm/docs/02-development/task-breakdown.md)
