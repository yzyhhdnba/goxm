# PILIPILI Go 后端技术报告（含八股与面试问答）

## 0. 报告定位

本报告面向以下目标：

1. 系统化梳理 PILIPILI Go 后端已实现能力与工程设计。
2. 整理该项目涉及的 Go 后端核心技术点与常见八股问题。
3. 提供可直接用于实习面试的项目深挖问答。
4. 给出当前短板与下一步补强路线。

本报告基于仓库当前代码与文档事实编写，不夸大未完成能力。

## 1. 项目总体概览

PILIPILI Go 是从前端业务原型演进而来的 Go 后端重构项目，目标是形成“可运行、可联调、可扩展”的在线视频互动平台后端。

当前后端已经形成完整主链路：

1. 认证：注册、登录、刷新、登出、用户信息。
2. 内容：分区、推荐流、热门流、关注流、分区流、视频详情。
3. 互动：点赞、收藏、评论、回复、评论点赞、关注关系。
4. 创作者：投稿、编辑、源视频上传、封面上传、稿件列表。
5. 管理：待审/已审列表、通过/驳回、今日统计、分区统计。
6. 站内：搜索、历史记录、通知读取。
7. 工程：MySQL + Redis 探活、CORS、优雅停机、统一响应、路由级测试与 E2E 脚本。

## 2. 后端技术栈全景

语言与框架：

1. Go 1.25
2. Gin
3. GORM
4. JWT（github.com/golang-jwt/jwt/v5）
5. Redis 客户端（github.com/redis/go-redis/v9）
6. YAML 配置加载（gopkg.in/yaml.v3）

数据库与中间件：

1. MySQL（生产主库）
2. SQLite（路由级测试）
3. Redis（健康探活与后续缓存扩展基础）

工程能力：

1. Docker Compose（MySQL + Redis）
2. 自动迁移 AutoMigrate
3. 默认分区 Seed
4. 优雅停机（signal + server.Shutdown）
5. 端到端脚本（bash + curl + jq）

## 3. 后端目录架构与分层设计

目录形态采用模块化单体 + 业务域拆分。

1. cmd/api：程序入口，负责初始化配置、DB、Redis、迁移、路由、启动与关停。
2. internal/config：配置结构、默认值、启动前强校验。
3. internal/db：GORM 连接与连接池设置。
4. internal/http：路由拼装、全局中间件、健康检查。
5. internal/middleware/auth：强鉴权与软鉴权。
6. internal/* 业务域：account、area、video、comment、social、search、history、notice、admin。
7. pkg：通用能力（authctx、request 分页解析、response 统一返回）。
8. scripts：E2E 联调脚本。

每个业务域均采用 Handler-Service-Repository 三层：

1. Handler 负责参数绑定、状态码、错误码映射。
2. Service 负责输入校验、业务编排、领域规则。
3. Repository 负责 SQL/GORM、事务、行级更新。

这个分层对面试官最重要的价值是：

1. 代码职责边界清晰。
2. 容易做单测和替换实现。
3. 便于未来拆 Worker 或拆服务。

## 4. 启动生命周期与运行时流程

启动链路（main）：

1. 读取配置并校验。
2. 连接 MySQL，设置连接池，启动前 Ping。
3. 连接 Redis，启动前 Ping。
4. 执行各业务域 AutoMigrate。
5. Seed 默认分区（幂等）。
6. 创建 Gin Router。
7. 启动 HTTP Server。
8. 监听 SIGINT/SIGTERM 并执行 10s 优雅停机。

这体现了两个后端工程点：

1. Fail Fast：配置错误和依赖不可用直接启动失败。
2. Graceful Shutdown：减少在飞请求损失。

## 5. 公共基础层技术点

### 5.1 配置管理

关键点：

1. 支持默认值回填（server、db pool、jwt ttl、cors、media）。
2. 启动前强校验：database.dsn、redis.addr、jwt secret、port、ttl。
3. 防占位值：jwt secret 含 change-me 会拒绝启动。
4. media.public_base_url 必须以斜杠开头。

面试可讲价值：

1. 避免“服务启动成功但运行异常”的软故障。
2. 降低线上错误配置导致的安全风险。

### 5.2 数据库层

关键点：

1. GORM + MySQL。
2. 连接池参数：MaxIdleConns、MaxOpenConns、ConnMaxLifetime。
3. 启动前 context 超时 Ping。

面试可讲取舍：

1. 先用 GORM 提速开发。
2. SQL 复杂查询阶段通过 Table + Joins + Select 控制结果字段，避免全量模型反序列化。

### 5.3 Redis 层

关键点：

1. 启动即连通性校验。
2. 健康检查接口中联动探活。
3. 当前先承接基础设施，后续可扩展计数缓存、榜单缓存、会话黑名单。

### 5.4 统一响应封装

统一 Envelope：

1. code
2. message
3. data

成功 code 固定为 0，失败 code 由业务域定义。

价值：

1. 前端解析逻辑统一。
2. 出错定位可直接按业务 code 检索。

### 5.5 authctx 与分页解析

authctx：

1. 既写 gin context，也写 request context。
2. admin 服务可直接从 context 读取 current user，减少重复查库。

request.ParsePagination：

1. 只做解析，不做默认值。
2. 默认值统一在业务 service 层做，保证不同域可定制。

## 6. 鉴权与会话体系（核心）

### 6.1 令牌模型

采用双令牌：

1. access token：短期、用于访问 API。
2. refresh token：长期、用于换新 token 对。

JWT Claims 字段：

1. user_id
2. token_version
3. token_type
4. refresh_id（仅 refresh token）

### 6.2 Refresh 安全机制

关键设计：

1. refresh token 中包含 refresh_id 随机值。
2. 数据库存储 refresh_id 的 sha256 hash，而不是明文 refresh token。
3. refresh 时校验：用户状态、token_version、一致的 refresh hash。
4. 续签更新采用条件更新：where user_id + token_version + old_hash。

这个设计可以回答两类高频追问：

1. 如何避免并发 refresh 导致多 token 并存。
2. token 泄露后的可控失效策略。

### 6.3 Logout 吊销策略

登出时：

1. 清空 refresh_token_hash。
2. token_version + 1。

结果：

1. 老 refresh token 立即失效。
2. 已签发 access token 因 token_version 不匹配而失效。

### 6.4 RequireAuth 与 OptionalAuth

RequireAuth：

1. 严格要求 Bearer Token。
2. 校验 token、用户状态、token_version。
3. 不通过即 401。

OptionalAuth：

1. 有 token 就尝试解析。
2. 解析失败也放行匿名访问。
3. 成功则补充 current user。

典型场景：

1. 视频详情可匿名访问。
2. 若已登录则返回真实 viewer_state。

## 7. 业务域技术细节

### 7.1 Area 域

能力：

1. 分区表迁移。
2. 启动时默认分区幂等写入（slug 冲突 DoNothing）。
3. 查询 active 分区并按 sort_order 排序。

技术点：

1. 启动期 Seed 设计。
2. 幂等初始化。

### 7.2 Video 域

能力清单：

1. 推荐流 feed/recommend（cursor + limit）。
2. 热门流 feed/hot（hot_score + 时间 + id 复合游标）。
3. 关注流 feed/following。
4. 分区流 areas/:id/videos（当前 sort=latest）。
5. 作者视频 users/:id/videos（page + page_size）。
6. 视频详情 videos/:id（支持 viewer_state）。
7. 点赞/取消点赞/状态。
8. 收藏/取消收藏/状态。
9. 投稿创建。
10. 稿件编辑（编辑后回 pending）。
11. 源视频上传、封面上传。
12. 创作者稿件列表 creator/videos。

关键技术实现：

1. 游标分页防漂移：按 published_at + id 或 hot_score + published_at + id。
2. 公共视频查询统一过滤：visible + approved + published_at 非空。
3. 点赞收藏使用事务 + OnConflict DoNothing + 计数增减。
4. 计数降到 0 时保护，使用 CASE WHEN 防负数。
5. 编辑已通过稿件后回待审，并回滚作者 video_count。
6. 上传支持扩展名白名单与大小限制。

### 7.3 Comment 域

能力清单：

1. 一级评论列表。
2. 发表评论。
3. 回复列表。
4. 发表回复。
5. 评论点赞/取消点赞/状态。

建模策略：

1. 单表承载评论与回复。
2. 一级评论 root_id=0 parent_id=0。
3. 回复 root_id=一级评论ID，parent_id=直接父评论ID。

关键实现点：

1. 评论与回复写入使用事务。
2. 回复时更新根评论 reply_count。
3. 评论写入同步更新视频 comment_count。
4. 列表查询用 JOIN users 一次取齐展示字段。
5. viewer_state.liked 通过 LikeMap 批量查询回填。

### 7.4 Social 域

能力：

1. 关注。
2. 取关。
3. 关注状态。
4. 粉丝列表。
5. 关注列表。

关键实现：

1. Follow/Unfollow 使用事务。
2. follows 关系表唯一索引保证不重复关注。
3. 关注成功后递增 following_count / follower_count。
4. 取消关注后做下界保护，避免计数负数。
5. 禁止关注自己。

### 7.5 Search 域

能力：

1. 视频搜索。
2. 用户搜索。

关键实现：

1. 第一阶段基于 MySQL LIKE。
2. 视频搜索复用公共可见视频过滤条件。
3. 用户搜索按 follower_count、video_count 排序。
4. keyword 为空直接判定无效输入。

### 7.6 History 域

能力：

1. 上报观看进度。
2. 查询历史记录。

关键实现：

1. upsert based on (user_id, video_id) 唯一键。
2. watched_at 每次更新。
3. 查询时联表 videos/users/areas 补齐展示信息。
4. 历史仅展示仍然公开可见的视频。

### 7.7 Notice 域

能力：

1. 通知列表。
2. 通知已读。

关键实现：

1. is_read + read_at 模型。
2. mark read 幂等：已读直接返回。
3. admin 审核动作自动写入 video_review 类型通知。

### 7.8 Admin 域

能力：

1. 查询待审/已审/全量稿件。
2. 审核通过。
3. 审核驳回（含原因）。
4. 今日统计。
5. 分区统计。

关键实现：

1. admin 权限双路径校验：优先 context，兜底查库。
2. 审核动作事务化：更新视频状态、写审核日志、写通知、维护作者 video_count。
3. approved 时设置 published_at。
4. rejected/回退时清空 published_at。
5. 今日统计聚合 active users、投稿量、通过量、播放量、评论量。

## 8. 媒体存储实现细节

当前阶段为本地落盘策略：

1. 视频白名单：.mp4 .m4v .mov .webm。
2. 封面白名单：.jpg .jpeg .png .webp。
3. 大小限制：视频 200MB，图片 10MB。
4. 路径规则：storage/videos/{videoID}/source.ext 与 cover.ext。
5. 对外访问：/uploads 前缀静态路由。

这个设计的价值是：

1. 先跑通投稿链路。
2. 后续可替换为对象存储，不影响上层接口。

## 9. 分页策略与查询风格

采用“双分页并存”的现实策略：

1. 动态流（推荐/热门/关注/分区）用 cursor + limit。
2. 弱动态列表（作者视频、搜索、历史、通知、后台）用 page + page_size。

这样做的核心取舍：

1. 动态流避免页码漂移。
2. 后台管理与列表查询保持简单稳定。

## 10. 事务、一致性与幂等策略

项目中可以明确讲出的工程一致性策略：

1. 点赞收藏：OnConflict DoNothing + 计数原子更新。
2. 取关取消点赞：删除行成功才递减计数。
3. 计数字段统一防负数。
4. refresh rotation：条件更新避免并发穿透。
5. history upsert：同一用户同一视频只保留一条记录。
6. admin review：审核动作事务化，保证状态、通知、日志一致。

## 11. 安全设计与风险控制

已实现安全点：

1. 密码 bcrypt 哈希。
2. JWT access/refresh 分离。
3. refresh_id 哈希落库。
4. token_version 统一吊销。
5. 配置阶段禁止占位 secret。
6. CORS 允许源可配置。
7. 需要登录接口统一 requiredAuth。

仍需补强的安全点：

1. 登录限流与防爆破。
2. 上传内容 MIME 与内容嗅探校验。
3. 关键操作审计日志（目前有审核日志，范围可扩展）。
4. 更严格的输入输出审计与安全扫描。

## 12. 测试与联调体系

当前测试形态：

1. 路由级集成测试（router_test.go）。
2. SQLite 内存库作为测试数据库。
3. 覆盖 Batch B/C/D/E 与 Phase4 主链路。
4. E2E bash 脚本覆盖真实接口联调。

测试覆盖重点：

1. refresh 轮换与 logout 吊销。
2. viewer_state 在匿名/登录态差异。
3. 热门流与关注流游标格式。
4. 评论回复与点赞取消。
5. 上传、审核、通知、dashboard 闭环。

当前不足：

1. 缺少高并发压测脚本。
2. 缺少细粒度单元测试与 mock。
3. 缺少故障注入测试（Redis 挂、DB 慢查询等）。

## 13. 后端 Go 技术点总清单（面试可背）

语言与语义：

1. 指针接收者与值接收者。
2. error 返回与 errors.Is 判型。
3. context 贯穿请求链路。
4. interface 抽象（如 storage、follow checker）。
5. map/slice 预分配优化。

并发与生命周期：

1. goroutine 启服务。
2. signal.NotifyContext 优雅停机。
3. context.WithTimeout 依赖探活超时控制。

Web 工程：

1. Gin 路由分组与中间件链。
2. 可选鉴权与强制鉴权分离。
3. 统一响应结构。
4. 统一错误码分域管理。

数据层：

1. GORM AutoMigrate。
2. Transaction 保证多写一致性。
3. OnConflict DoNothing 做幂等插入。
4. Count + List 分离查询。
5. 多表联查 Select 精准字段。

安全层：

1. bcrypt 密码哈希。
2. JWT claims 设计。
3. refresh token hash。
4. token_version 吊销机制。

存储层：

1. multipart 上传处理。
2. 扩展名白名单。
3. 大小限制。
4. 本地落盘目录组织。

## 14. 当前技术短板与真实改进路线

短板：

1. 推荐逻辑仍是规则排序，个性化不足。
2. 搜索仍基于 LIKE，规模上来后会变慢。
3. 统计字段目前以数据库主写为主，缓存回刷链路未全面落地。
4. 缺少 MQ/Worker 异步化处理。
5. 缺少完整可观测性（metrics/trace/告警）。
6. 测试偏集成，单测和压测不足。

建议改进路线：

1. P1：补 Prometheus 指标、慢 SQL 日志、基础压测。
2. P2：引入 Redis 计数缓冲 + 定时回刷。
3. P3：引入 MQ 处理审核通知、计数异步任务。
4. P4：搜索升级 FULLTEXT 或 ES。
5. P5：媒体处理升级 ffmpeg + HLS + 对象存储。

## 15. Go 后端八股文整理（高频 60 题）

### 15.1 Go 语言基础（15 题）

1. 问：Go 的值传递和引用传递如何理解。
答：Go 只有值传递；传指针只是把指针值拷贝过去，底层对象可共享修改。

2. 问：slice 扩容机制有什么特点。
答：容量不足会分配新数组并拷贝，扩容倍率在不同容量区间不一致，不能依赖固定倍数。

3. 问：map 并发安全吗。
答：不安全；并发读写会 panic，需要加锁或用 sync.Map。

4. 问：defer 的执行时机和常见坑。
答：函数返回前逆序执行；defer 参数在声明时求值，闭包变量是引用语义。

5. 问：panic 和 recover 关系。
答：panic 触发栈展开；recover 仅在 defer 中生效，用于兜底防崩。

6. 问：error 处理推荐方式。
答：显式返回 error；使用 errors.Is/As 判型，不建议字符串比较。

7. 问：interface 底层包含什么。
答：类型信息 + 数据指针。nil 接口和值为 nil 的具体类型需要区分。

8. 问：context 应该放什么。
答：放请求级控制信息（超时、取消、trace）；不要放业务大对象。

9. 问：goroutine 泄漏常见原因。
答：阻塞 channel、未取消 context、后台任务无退出条件。

10. 问：channel 关闭规则。
答：通常由发送方关闭；关闭后可继续接收零值和 ok=false。

11. 问：Go GC 对延迟影响怎么看。
答：GC 会产生停顿但较短；需关注对象分配频率、逃逸和大对象生命周期。

12. 问：逃逸分析有何意义。
答：决定对象在栈或堆分配，堆分配增加 GC 压力。

13. 问：sync.Mutex 和 RWMutex 何时选。
答：读多写少可用 RWMutex；写频繁或临界区短时 Mutex 反而更稳。

14. 问：atomic 与锁怎么取舍。
答：简单计数和标志位可 atomic；复杂复合状态仍用锁。

15. 问：Go 模块管理核心命令。
答：go mod tidy、go list、go test、go build；保持依赖最小集和可重复构建。

### 15.2 Gin 与 HTTP（10 题）

16. 问：Gin 中间件执行顺序。
答：按注册顺序进入，按栈回退顺序退出；Abort 会中断后续 handler。

17. 问：为什么要做路由分组。
答：统一前缀和中间件，降低重复定义，方便版本化管理。

18. 问：什么时候用 optional auth。
答：匿名可访问但登录后需返回个性化状态的接口，如详情 viewer_state。

19. 问：如何统一错误响应。
答：封装 response.Error，统一 code/message/data 结构，前后端协议稳定。

20. 问：CORS 配置重点。
答：AllowOrigins、AllowMethods、AllowHeaders、AllowCredentials 要与前端策略一致。

21. 问：大文件上传应关注什么。
答：大小限制、扩展名/类型校验、路径隔离、覆盖策略与权限校验。

22. 问：为什么健康检查不能只返回进程活着。
答：需要包含依赖探活（DB/Redis），否则“活着但不可用”。

23. 问：如何处理请求超时。
答：上游网关 + 应用 context timeout + 下游查询超时联动。

24. 问：HTTP 状态码与业务码关系。
答：状态码表达协议层语义，业务码表达领域语义，两者都要稳定。

25. 问：为什么需要优雅停机。
答：避免重启时中断在飞请求，减少用户错误和数据不一致。

### 15.3 GORM/MySQL（12 题）

26. 问：为什么要显式事务。
答：跨表或多步骤写入必须保证原子性，如审核动作和计数更新。

27. 问：OnConflict DoNothing 价值。
答：天然幂等，适合点赞/关注这类“重复提交不重复生效”场景。

28. 问：如何避免计数变负数。
答：SQL 使用 CASE WHEN count>0 THEN count-1 ELSE 0 END。

29. 问：为什么查询要 Select 指定字段。
答：减少网络与序列化开销，防止无关字段泄漏。

30. 问：软删除的优势。
答：可回溯、可恢复，且配合 deleted_at 自动过滤。

31. 问：联合索引设计原则。
答：按最常用过滤和排序字段排列，如 area_id + review_status。

32. 问：count 和 list 分开查的原因。
答：分页总数与列表通常执行计划不同，分开更可控。

33. 问：如何做唯一约束幂等。
答：依赖数据库唯一索引，应用层只做业务提示。

34. 问：为何先检查存在再写入还不够。
答：并发下会竞争，最终一致性仍要靠唯一索引/事务约束。

35. 问：慢查询排查先看什么。
答：执行计划、索引命中、回表、排序与临时表。

36. 问：LIKE 查询的瓶颈。
答：前置通配会破坏索引，数据量大时全表扫描明显。

37. 问：何时升级到 ES/FULLTEXT。
答：搜索量和数据量增长导致响应不可接受，且需要相关性排序时。

### 15.4 Redis/JWT/安全（12 题）

38. 问：为什么 access 和 refresh 分离。
答：兼顾安全与体验，短 access 降低泄露风险，refresh 保持登录态连续。

39. 问：refresh token 为何存 hash。
答：数据库泄露时不能直接重放 refresh token。

40. 问：token_version 的作用。
答：全局会话版本控制，可一键失效历史 access/refresh。

41. 问：并发 refresh 如何处理。
答：条件更新 old_hash + token_version，确保只有一次旋转成功。

42. 问：JWT 校验最重要的字段。
答：签名算法、过期时间、token_type、issuer、业务版本字段。

43. 问：为什么要禁用占位 secret。
答：防止默认弱密钥被误上生产。

44. 问：登录接口应加哪些安全措施。
答：限流、验证码、IP 风险评估、失败锁定策略。

45. 问：上传接口安全点。
答：大小、后缀、MIME、存储路径隔离、鉴权和审计。

46. 问：XSS 风险在本项目哪里出现。
答：评论内容和通知内容展示链路，需要前端渲染时做转义策略。

47. 问：CSRF 在 JWT 架构是否还存在。
答：若 token 存于 Authorization Header 且无 Cookie 自动带出，风险较低，但仍需结合前端策略。

48. 问：Redis 在当前项目价值。
答：基础设施连通性与后续缓存扩展前置，不让架构演进断层。

49. 问：如何做接口防重放。
答：时间戳 + nonce + 签名，或关键写操作幂等键。

### 15.5 架构与系统设计（11 题）

50. 问：为什么先模块化单体而不是微服务。
答：业务尚在快速迭代，单体更快闭环，后续可按域逐步拆分。

51. 问：为什么评论和回复用统一表。
答：减少跨表复杂度，统一索引与查询模型，便于楼中楼扩展。

52. 问：动态流为什么用 cursor。
答：防止并发写入导致页码漂移，滚动加载体验更稳定。

53. 问：为什么列表排序通常带 id 兜底。
答：当时间戳并列时用 id 保证全序和游标稳定。

54. 问：审核状态机如何设计。
答：pending -> approved/rejected，编辑后可回 pending，再次进入审核闭环。

55. 问：计数字段放主表还是关系表聚合。
答：主表读快，关系表准；常见做法是主表缓存计数 + 异步校准。

56. 问：如何理解“事实字段”和“派生字段”。
答：关系表是事实，计数字段是派生值，允许短暂延迟一致。

57. 问：为何需要统一错误码分段。
答：便于快速定位业务域和报警聚类分析。

58. 问：如何控制接口演进风险。
答：约束统一 API 层、保留兼容语义、先加后删。

59. 问：如何平衡“先跑通”和“高性能”。
答：先保证主链路完整与正确，再按瓶颈做局部优化。

60. 问：该项目下一步拆分优先级。
答：先抽异步任务（审核通知、计数回刷），再考虑 media worker。

## 16. 项目深挖面试问答（40 题）

1. 问：你为什么设计 OptionalAuth。
答：详情页属于公开资源，但登录态需要 viewer_state。OptionalAuth 既保留匿名访问，又避免前端多打一轮状态接口。

2. 问：refresh 成功后为什么旧 access 还可能有效。
答：这是故意设计。access 是短期令牌，refresh 只轮换 refresh，用户体验更好；登出时再用 token_version 全部吊销。

3. 问：如何证明 refresh 轮换安全。
答：服务端只存 refresh_id 的 hash，refresh 时比对 hash 且条件更新 old_hash，旧 token 会立即失效。

4. 问：如果用户并发点两次刷新会怎样。
答：只会有一次更新成功，另一次因 old_hash 不匹配失败，返回无效 refresh token。

5. 问：为什么点赞要先 ensurePublicVideo。
答：避免对不可见或未审核视频产生互动数据，保持业务一致性。

6. 问：为什么点赞用 OnConflict DoNothing。
答：天然幂等，解决重复点击或网络重试导致的重复写。

7. 问：为什么 unlike 时要判断 RowsAffected。
答：只有真正删除了点赞记录才应该减计数，避免计数错误。

8. 问：评论回复为什么不拆表。
答：统一 comments 表降低维护复杂度，用 root_id/parent_id 即可表达一级和楼中楼关系。

9. 问：回复列表为什么按创建时间正序。
答：楼中楼阅读更符合对话顺序，根评论列表则按时间倒序更符合信息流。

10. 问：为什么评论列表要回填 likedMap。
答：批量查询后回填 viewer_state，避免 N+1 查询。

11. 问：为什么 feed/hot 的游标是三段。
答：排序键是 hot_score desc + published_at desc + id desc，游标必须覆盖全部排序维度。

12. 问：为什么 area 流只支持 latest。
答：当前阶段先收敛接口复杂度，避免前后端契约频繁变动。

13. 问：编辑稿件后为什么回 pending。
答：内容发生变更后应重新进入审核流程，防止绕过审核直接发布。

14. 问：为什么审核通过时才设置 published_at。
答：published_at 代表对外发布时间，不能在待审态提前写入。

15. 问：为什么驳回要清空 published_at。
答：保证公开态判定一致，驳回内容不应参与对外流量分发。

16. 问：审核为什么要写 notice。
答：让创作者能感知流程状态，形成投稿-审核-反馈闭环。

17. 问：审核日志 video_reviews 的价值。
答：可追溯谁在什么时间做了什么决策，支持审计与争议排查。

18. 问：dashboard 的 video_count 为什么会变化。
答：它统计的是已通过并公开的视频数量，编辑回待审会减少。

19. 问：为什么 history 用 upsert。
答：同一用户同一视频只保留最新进度和观看时间，避免重复数据膨胀。

20. 问：history 为什么要联表过滤视频状态。
答：避免展示已经下线或未审核内容，保证前台一致性。

21. 问：search 为什么 keyword 为空直接报错。
答：无关键词搜索会放大数据库压力且语义不明确。

22. 问：后台权限如何校验。
答：先读 authctx 的 current user role，缺失时再查库兜底。

23. 问：为什么要统一 response.Envelope。
答：前端只需一套解析逻辑，且便于日志平台按 code 聚类。

24. 问：为什么 error code 分段。
答：快速定位来源域，示例：31xx 视频，41xx 评论，52xx 社交。

25. 问：启动时为什么要 AutoMigrate。
答：当前阶段强调快速迭代和本地联调效率，后续可迁移到版本化 migration。

26. 问：为什么还要 Seed 默认分区。
答：保证空库可运行，避免前端上传页/首页无分区导致链路中断。

27. 问：健康检查为什么包含 Redis。
答：虽然部分接口不依赖 Redis，但作为基础能力应可观测其可用性。

28. 问：为什么前端还能构建但需要 openssl 兼容参数。
答：是当前本机 Node 与旧前端依赖组合导致的兼容性问题，不影响 Go 后端设计本身。

29. 问：如果未来接入对象存储，改动大吗。
答：不大。Service 依赖的是 MediaStorage 接口，替换 LocalStorage 实现即可。

30. 问：如果未来接入 MQ，先改哪。
答：先把审核通知、计数回刷等异步任务从主请求中抽离。

31. 问：你如何看待当前搜索方案。
答：MVP 阶段可接受，数据量上升后必须升级 FULLTEXT/ES。

32. 问：你如何证明项目是“工程化”而不只是“接口拼接”。
答：有分层架构、统一错误码、优雅停机、探活、测试、E2E 脚本、状态机与审计链路。

33. 问：如果 Redis 挂了会怎样。
答：启动时会失败，运行时 healthz 报依赖不可用；这是当前的 fail-fast 策略。

34. 问：如果数据库短暂抖动呢。
答：请求层会报错，下一步建议加重试和断路器并配合连接池监控。

35. 问：如何控制评论恶意刷屏。
答：当前有长度限制；下一步应加频控、敏感词和风控策略。

36. 问：如何做接口压测。
答：建议先对 feed/detail/comment/create/upload 四类路径做基线压测，建立 p95/p99 指标。

37. 问：你的测试为什么选 sqlite。
答：路由级测试追求轻量和隔离；SQL 行为差异需再补 MySQL 集成测试兜底。

38. 问：你如何定义该项目 MVP 完成度。
答：用户、内容、互动、创作、审核、通知、搜索、历史都可用，具备端到端闭环。

39. 问：这个项目最大的工程风险是什么。
答：计数一致性和搜索扩展性，随着数据增长会先暴露。

40. 问：如果让你再做一次，先改什么。
答：先补可观测性与压测，再做计数异步化和搜索升级。

## 17. 面试表达模板（可直接背）

30 秒版本：

我做的是一个 Go 后端重构项目，基于现有视频站前端原型，完成了从鉴权到推荐流、评论互动、上传审核、通知历史的完整主链路。技术上采用 Gin + GORM + MySQL + Redis，重点解决了 refresh token 安全轮换、动态流 cursor 分页、防重复写幂等和审核状态闭环，并通过路由级测试和 E2E 脚本保证回归质量。

90 秒版本：

项目采用模块化单体架构，按 account/video/comment/social/search/history/notice/admin 拆分。鉴权用 access + refresh 双令牌，refresh_id 哈希落库并结合 token_version，避免并发续签竞态并支持登出全量吊销。内容侧实现推荐、热门、关注、分区四类流，动态列表统一 cursor 分页，防止页码漂移。互动侧支持点赞收藏评论回复，全部做了事务和幂等保护。创作者链路支持投稿、编辑回待审、上传源文件和封面；管理员可以审核并自动写通知。测试上用 sqlite 跑路由集成测试，覆盖 B 到 Phase4 主流程，并有 bash E2E 脚本做联调验证。

## 18. 对 Agent 开发岗的补强建议

如果目标岗位是 Agent 开发，本项目作为后端底座已合格，但建议新增两类能力：

1. Tool-Augmented Agent：做一个“审核助手 Agent”或“运营问答 Agent”，具备工具调用、错误重试和结构化输出。
2. Agent Eval：建立小型评测集，输出任务成功率、调用成本、失败样例与修复策略。

这样简历会形成“后端工程能力 + Agent 实践能力”的闭环。

---

以上内容可直接作为你的项目技术报告主文档使用。