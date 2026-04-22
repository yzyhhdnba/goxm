# AGENTS.md

## 1. 文档定位

本文件主要写给 **AI 代理 / 编程助手**，同时也兼顾人类开发者阅读。

它的作用不是介绍项目，而是约束后续协作行为，确保所有人围绕同一条路线推进 `PILIPILI Go`，避免：

- 擅自改方向
- 把规划写成已实现
- 在未确认的地方过度设计
- 偏离当前阶段目标

如果你是新接手这个仓库的代理，请把这份文档看成“进入项目后的第一份操作约束”。

---

## 2. 项目当前事实

当前仓库处于 **Go 化规划与重构阶段**。

这意味着：

- 已经存在 `backend/` 基础骨架
- 已经存在根目录 `docker-compose.yml`
- 已经存在 MySQL 初始化脚本
- 已经验证最小 API 启动链路可用
- 已经实现用户鉴权的第一版主链路
- 已经完成 Batch B 的注册、登录、刷新、登出、可用性检查与 `users/me`
- 已经完成 Batch C 的 `areas`、`feed/recommend` 与 `videos/:id`
- 已经完成 Batch D 的点赞、收藏、评论、回复、评论点赞与关注关系
- 已经完成 Batch E 的 `feed/hot`、`feed/following`、`areas/:id/videos` 与 `users/:id/videos`
- 已经完成 Phase 4 的 `search`、`history`、`creator upload`、`admin review` 主链路
- 已经修复视频详情页评论 / 回复前端触发链路，并清理关键页面残留的老旧外部资源硬编码
- 已经接入 Redis 客户端初始化与 `healthz` 探活
- 已经接入视频详情与热门榜的 Redis 读缓存，并在关键写路径上做主动失效
- 已经接入主程序优雅停机
- 已经接入开发期 CORS 中间件
- 已经把原有前端迁入当前工作区，并在 `frontend/src/api/index.ts` 收口核心 API
- 已经新增根目录 `frontend1/` 作为并行 React 前端工程，用于在不替换现有 Vue 前端的前提下验证 React 技术栈重构
- 当前工作区已经初始化为根级 Git 仓库
- 仓库里仍保留 `backend/frontend/` 历史副本，当前主前端目录以根目录 `frontend/` 为准
- 当前前端构建在本机 Node 环境下需要 `NODE_OPTIONS=--openssl-legacy-provider`

任何代理或协作者都 **不得** 在未实现前把功能描述为“已完成”，但已经验证通过的基础骨架能力可以如实说明。

---

## 3. 项目目标

`PILIPILI Go` 的目标是把原始前端业务原型，升级成一个具备真实后端能力的在线视频与互动平台项目。

目标能力包括：

- 用户系统
- 视频系统
- 评论与回复
- 点赞与收藏
- 关注关系
- 历史与通知
- 上传与审核
- Feed / 热门推荐
- Redis / Worker / MQ 的渐进式工程增强

注意：

- 当前目标是 **渐进重构**
- 不是一步到位上所有高级特性
- 优先级永远高于“炫技式架构”

---

## 4. 代理必须先阅读的文档

开始任何实质性工作前，请按以下顺序阅读：

1. [README.md](/Users/hhd/Desktop/test/goxm/README.md)
2. [blueprint.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/blueprint.md)
3. [schema.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/schema.md)
4. [api.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/api.md)
5. [technical-architecture.md](/Users/hhd/Desktop/test/goxm/docs/02-development/technical-architecture.md)
6. [AGENTS.md](/Users/hhd/Desktop/test/goxm/AGENTS.md)

如果代码和文档未来出现冲突，优先处理方式不是擅自重写文档，而是：

1. 先确认当前实现是否偏离规划
2. 判断偏离是否合理
3. 若涉及重大方向变化，先询问用户

---

## 5. 决策优先级

在本项目中，决策优先级如下：

1. 用户明确指令
2. 当前仓库文档契约
3. 当前阶段的开发目标
4. 参考项目的设计思路
5. 代理自己的默认偏好

这意味着：

- `feedsystem_video_go` 只能作为参考，不能压过当前项目自己的文档
- 代理不能因为“更现代”或“更通用”就擅自改掉既定方案

---

## 6. 当前阶段的工作边界

当前阶段允许做的事情：

- 完善项目文档
- 细化表结构与 API
- 初始化 Go 后端骨架
- 规划前后端目录
- 搭建基础配置和 Docker Compose
- 从第一阶段主链路开始实现
- 运行依赖并验证最小 API 链路

当前阶段不应该优先做的事情：

- 一上来拆微服务
- 一上来做 Kafka / ES / Kubernetes
- 一上来实现复杂推荐算法
- 在 MVP 未完成前铺开全量 MQ 事件总线
- 在没有必要时推翻原有前端业务结构

核心原则：

> 先做“完整可跑通的视频平台主链路”，再做“漂亮的工程增强”。

---

## 7. 目标技术路线

默认路线如下，除非用户明确要求变更：

### 第一阶段

- Go
- Gin
- GORM
- MySQL
- Redis
- Docker Compose

### 第二阶段

- RocketMQ
- Worker
- ffmpeg
- MinIO 或 OSS

### 默认前端策略

- 保留原有 Vue 前端
- 不重写为 React
- 不在当前阶段放弃现有交互结构

### 当前已采纳的关键设计

- 鉴权采用 `access token + refresh token`
- refresh token 第一阶段采用单设备可轮换策略
- token 失效基于 `token_version + refresh_token_hash`
- refresh 只轮换 refresh token，不立即吊销当前 access token
- logout 负责提升 `token_version` 并清空 `refresh_token_hash`
- refresh 落库必须携带旧 `refresh_token_hash` 条件，避免并发续签竞态
- 动态 Feed 文档统一向 cursor 分页收敛
- 业务状态与技术软删除分离：`status + deleted_at`
- 视频计数查询保留主表字段，但写路径目标是 Redis 优先再回刷 MySQL
- 当前 Redis 已落地的部分是视频详情与热门榜读缓存；“写 Redis 再异步回刷 MySQL”仍属于后续增强
- `feed/recommend` 当前只返回 `visible + approved + published_at 非空` 的公开视频
- `feed/hot` 当前已实现，排序规则为 `hot_score desc, published_at desc, id desc`
- `feed/following` 当前已实现，只返回当前用户已关注作者的公开视频
- `areas/:id/videos` 当前已实现，默认只支持 `sort=latest`
- `users/:id/videos` 当前已实现，并继续保留 `page + page_size`
- `search/videos` 与 `search/users` 当前已实现，并继续保留 `page + page_size`
- 搜索当前仍基于 MySQL `LIKE` 模糊匹配，大数据量阶段应升级到 FULLTEXT 或 Elasticsearch
- `histories` 当前已实现，接口需要登录
- `notices` 当前已实现，管理员审核通过 / 驳回时会写入创作者通知
- `PATCH /notices/:id/read` 当前已实现，用于站内通知已读
- `users/:id/profile` 当前已实现，支持可选鉴权下的 `viewer_state.followed`
- `users/me/dashboard` 当前已实现，承接原前端 `getMyIndexInfo` 聚合需求
- `videos` 创建元数据、源文件上传、封面上传当前已实现，默认落盘到 `backend/storage/`
- API 当前通过 Gin 静态路由暴露 `/uploads/...`
- `PATCH /videos/:id` 当前已实现，作者编辑稿件后会回到 `pending`
- `creator/videos` 当前已实现，支持 `approved/pending/rejected/all`
- `admin/videos/pending` 与 `admin/videos?review_status=reviewed` 当前已实现，且只允许 `role=admin`
- `admin/videos/:id/approve` 与 `admin/videos/:id/reject` 当前已实现，审核结果会驱动通知写入
- admin 请求当前会复用 middleware 写入的 `authctx.CurrentUser`，避免每次额外查一遍 `users`
- `admin/stats/today` 当前返回 `active_user_count / submitted_video_count / approved_video_count / play_count / comment_count`
- `admin/stats/area` 当前按分区返回 `approved_count / pending_count / rejected_count / total_count`
- 当前播放器阶段优先支持直接播放落盘的 mp4 文件，后续再升级到 ffmpeg + HLS
- 当前 RocketMQ / Worker 仅完成配置、连通性探测与运行骨架，不要把它表述成“业务异步链路已完成”
- `videos/:id` 当前通过软鉴权返回真实 `viewer_state`
- 评论列表与回复列表当前使用 `page + page_size`
- 热门榜 `next_cursor` 当前使用 `hot_score:unix_timestamp:video_id`
- 推荐流、关注流、分区流 `next_cursor` 当前使用 `unix_timestamp:video_id`
- handler 中当前用户上下文统一走 `pkg/authctx`
- handler 中页码分页参数统一走 `pkg/request.ParsePagination`
- 公共公开视频查询 builder 当前只保留单一 `Model/Select` 来源，避免重复 `Table` 链式覆盖
- 评论列表与回复列表当前已收敛为单次连表查询，不再重复查询 `comments` 主表
- 配置加载阶段必须校验数据库、Redis 与 JWT 密钥，不允许空值或占位值启动
- 主程序默认使用 `http.Server + signal.NotifyContext` 做优雅停机
- 前后端分离开发默认允许本地常见端口的 CORS 访问
- Phase 4 及后续前端页面改造优先走 `frontend/src/api/index.ts`，不要在 Personal / Notice / Search / History / Upload / Management 目录继续新增硬编码 Axios
- 前端 API 地址统一通过 `frontend/.env.*` 中的 `VUE_APP_API_BASE_URL` / `VUE_APP_API_PROXY_TARGET` 管理
- 前端通过 `frontend/src/utils/request.ts` 发请求时，当前约定是“直接返回真实 data 载荷”，不要在新代码里继续按旧 Axios 风格读取 `res.data`
- 评论、回复、点赞、关注等前端交互当前统一以 `access_token` 作为登录态判断，不要把 `userInfo` 当成唯一鉴权依据
- 头部菜单、管理页背景、管理导航头像等静态资源当前统一走本地常量 / 本地静态资源，不要重新引入 `101.35.142.191` 之类外部硬编码资源
- 路由级集成测试当前已经覆盖 Search / History / Notice / Personal Dashboard / Upload / Admin / Edit 主链路
- 若联调中发现旧视频标题乱码，优先判断为历史脏数据问题；在证据不足时，不要先回退当前 `utf8mb4` 配置

代理在推进 Batch C 及以后阶段时，必须继续遵守：

- `GET /api/v1/feed/recommend`
- `GET /api/v1/feed/hot`
- `GET /api/v1/feed/following`
- `GET /api/v1/areas/:id/videos`

以上动态列表接口默认使用 `cursor + limit`，不要回退成 `page + page_size`。

---

## 8. 仓库结构目标

后续目标结构应收敛为：

```text
frontend/
backend/
deploy/
docs/
README.md
AGENTS.md
docker-compose.yml
```

后端建议结构：

```text
backend/
  cmd/
    api/
    worker/
  configs/
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
  migrations/
```

代理在新增代码或文档时，应尽量向这个结构收敛，而不是随意分散。

补充约束：

- 当前前端源码、环境变量与联调改动，默认优先落在根目录 `frontend/`
- `frontend1/` 当前视为并行 React 前端工程，可在用户明确要求 React 重构或并行实现时继续推进
- `backend/frontend/` 当前视为历史遗留副本；若无用户明确要求，不要继续把新增改动写进该目录
- `backend/storage/`、`backend/bin/`、`backend/configs/config.yaml` 属于本地运行产物或本地配置，不要把这些内容表述成仓库契约

---

## 9. 领域边界约束

### 9.1 账户域

负责：

- 注册
- 登录
- 退出
- 当前用户
- 密码修改

### 9.2 视频域

负责：

- 视频元数据
- 视频详情
- 作者视频列表
- 播放地址
- 封面地址

### 9.3 评论域

负责：

- 一级评论
- 回复
- 评论点赞

约束：

- 评论和回复采用统一建模
- 不要轻易拆成两套完全独立的表和服务

### 9.4 社交域

负责：

- 关注
- 取关
- 粉丝
- 关注列表

### 9.5 Feed 域

负责：

- 首页推荐流
- 热门榜
- 关注流
- 分区流

### 9.6 管理后台域

负责：

- 视频审核
- 审核状态查询
- 统计数据

代理在实现时，不要把不同域的职责随意混在一个大文件或大 service 里。

---

## 10. 文档是当前阶段的契约

当前阶段的设计约束主要来自这三份文档：

- [blueprint.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/blueprint.md)
- [schema.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/schema.md)
- [api.md](/Users/hhd/Desktop/test/goxm/docs/01-contracts/api.md)

代理必须遵守以下规则：

- 不要在未同步文档的情况下私自改变数据模型
- 不要在未同步文档的情况下私自改路由风格
- 不要把旧接口和新接口混成无规则状态
- 若确实需要改契约，先说明原因，再更新文档

---

## 11. 重大问题必须先问用户

以下事项属于重大决策，代理 **必须先询问用户**，不能直接做：

- 改项目名称
- 改 README 的真实性口径
- 放弃 `PILIPILI Go` 的定位
- 把 Vue 前端改成其他框架
- 把模块化单体改成微服务
- 在核心栈上从 `Gin/GORM/MySQL/Redis` 改成别的主方案
- 提前引入 Kafka、ES、K8s 等超出当前阶段的基础设施
- 放弃统一评论建模
- 大改 API 风格导致前端整体重适配
- 对外宣称某些能力“已实现”，但仓库里实际上还没有

如果你不确定某个决策是否算重大问题，默认按“先问用户”处理。

---

## 12. 可以默认推进的事情

以下事情通常可以直接做，不需要额外确认：

- 完善文档说明
- 统一命名风格
- 补充缺失的错误码说明
- 初始化基础目录
- 搭建配置加载
- 增加开发脚手架
- 根据既有文档补充 DTO / model / router 骨架
- 修正明显的不一致表述

前提是：

- 不改变既定方向
- 不夸大项目完成度
- 不破坏现有规划

---

## 13. 写文档时的约束

所有代理在编写本仓库文档时，应遵守：

- 优先使用中文
- 先写事实，再写目标
- 已完成和计划中必须明确区分
- 不要用含糊表述混淆“当前状态”和“未来目标”
- 对外文档要克制，不要营销化夸大
- 对内文档要清晰，能指导下一步开发
- 每次做完实质性变更后，同步更新 `README.md` 和 `AGENTS.md`

尤其是 `README.md`：

- 不能写成“项目已经完成 Redis/MQ/Worker/转码全链路”
- 如果功能还没做，只能写为规划目标

如果某次改动还影响了接口、表结构或任务状态，也应同步更新相应设计文档。

---

## 14. 写代码时的约束

在后续开始编码后，代理应遵守：

- 先搭主链路，再做增强能力
- 优先完成鉴权、首页、详情、评论、点赞、关注
- 不要在没有 MVP 的情况下先做复杂热榜和事件系统
- 接口命名尽量向 `docs/01-contracts/api.md` 收敛
- 模型命名尽量向 `docs/01-contracts/schema.md` 收敛
- 若实现细节与文档不完全一致，优先补充说明，不要偷偷偏移

---

## 15. 与参考项目的关系

参考仓库：

- [PILIPILI 原始前端项目](https://gitee.com/yzynba/pilipili)
- [feedsystem_video_go](https://github.com/LeoninCS/feedsystem_video_go)

使用原则：

- 原始前端项目提供业务基底
- `feedsystem_video_go` 提供 Go 后端架构参考
- 参考不等于复制
- 任何“借鉴”都应服务于 `PILIPILI Go` 自己的业务落地

---

## 16. 推荐工作流

代理在开始一个新任务前，建议按以下顺序行动：

1. 先确认当前任务属于哪个阶段
2. 阅读相关设计文档
3. 判断是否涉及重大决策
4. 若涉及重大决策，先询问用户
5. 若不涉及，直接推进并保持文档一致
6. 完成后总结“做了什么、未做什么、还有什么风险”

---

## 17. 当前最重要的事情

当前阶段最重要的目标只有两个：

1. 建立清晰、可信、可执行的文档基线
2. 为后续 `backend/` 初始化和第一阶段开发提供稳定约束

在这两个目标完成前，不要急于追求复杂实现。

---

## 18. 一句话原则

> 对 `PILIPILI Go` 来说，最重要的不是“马上写很多代码”，而是“沿着同一套真实、清晰、渐进的路线持续推进”。
