# PILIPILI Go 面试讲解顺序

这份文档是 [KEY_CODE_IMPLEMENTATION.md](./KEY_CODE_IMPLEMENTATION.md) 的配套讲解稿，目标不是“把所有代码讲一遍”，而是帮助你在面试里按时间长短稳定输出。

适用场景：

- 面试官让你用 3 到 5 分钟快速介绍一个项目
- 面试官追问你“挑一个最熟的项目详细讲讲”
- 二面 / 交叉面想听你从架构、业务链路、技术细节、改进点四层展开

建议使用方式：

- 先背熟 `5 分钟版本`，保证任何情况下都能先把主线讲完整
- 再准备 `15 分钟版本`，覆盖登录、视频详情、评论、投稿审核这些高频链路
- 最后把 `30 分钟版本` 当成“被持续追问时的展开地图”

---

## 一、5 分钟版本：先把项目价值和主链路讲清楚

### 1. 开场模板

你可以直接按下面这个顺序讲：

1. 项目是什么
2. 解决了什么问题
3. 你负责了什么
4. 技术上最有代表性的点是什么
5. 结果怎么样

参考话术：

> 我做的是一个仿 B 站的视频平台后端重构项目，叫 **PILIPILI Go**。项目目标是把原来偏课程作业式的功能，重构成一个有清晰分层、可持续扩展的前后端系统。  
> 后端主要用 **Go + Gin + GORM + MySQL + Redis**，前端是 **Vue + Element Plus**。  
> 我重点梳理和实现了几个核心链路：登录鉴权、视频详情页、评论树、关注关系、投稿上传、后台审核。  
> 这个项目里我最想强调的是三件事：第一是后端分层比较清楚，请求从 Router 到 Service 到 Repository 的边界明确；第二是登录态做了双 Token 和 refresh 轮换；第三是视频详情页把点赞、收藏、关注这些 `viewer_state` 一次性聚合返回，前端消费会更顺。  
> 如果面试官愿意，我可以继续展开讲登录链路，或者讲视频详情和评论这条最完整的业务主链路。

### 2. 5 分钟必须覆盖的四个点

#### 点 1：技术栈和分层

- 后端：Go、Gin、GORM、MySQL、Redis、JWT
- 前端：Vue、Axios、Element Plus
- 分层：`router -> handler -> service -> repository -> db`

#### 点 2：最核心业务链路

优先讲这一条：

- 用户打开视频详情页
- 前端并发拉视频详情、评论列表、推荐视频
- 后端详情接口返回基础信息和 `viewer_state`
- 用户可以继续评论、点赞、收藏、关注作者

#### 点 3：一个代表性技术设计

优先讲这个：

- JWT 双 Token
- access token 短期有效
- refresh token 带 `refresh_id`
- 数据库保存 `refresh_token_hash`
- 刷新时做 `token_version + refresh hash` 校验

#### 点 4：你做了什么优化

- 把接口响应统一成 `Envelope`
- 把 viewer state 聚合到详情接口，减少前端额外请求
- 评论和回复列表用批量查询/JOIN，避免明显 N+1
- 审核流程放到事务里，保证状态和计数一致

### 3. 5 分钟版最容易被追问的问题

- 为什么要双 Token，而不是只用一个 JWT
- 为什么 refresh token 要存 hash，不直接存明文
- 为什么视频详情接口要返回 viewer_state
- 评论树怎么建模
- 审核为什么要做事务

---

## 二、15 分钟版本：按“架构 + 两条主链路 + 一个设计亮点”展开

15 分钟的目标，不是把所有模块都铺开，而是让面试官觉得你对系统真的有整体掌控。

### 1. 第一段：项目目标和整体架构（2 分钟）

建议顺序：

1. 先讲业务目标
2. 再讲技术栈
3. 再讲分层结构
4. 再讲你重点负责的部分

可以围绕下面这些文件展开：

- 程序入口：[../backend/cmd/api/main.go](../backend/cmd/api/main.go)
- 路由装配：[../backend/internal/http/router.go](../backend/internal/http/router.go)
- 统一请求层：[../frontend/src/utils/request.ts](../frontend/src/utils/request.ts)

你要强调：

- 启动时会加载配置、初始化 MySQL/Redis、自动迁移、注册路由
- 路由分公开接口、软鉴权接口、强鉴权接口
- 前端请求层统一处理 token 注入、响应解包、401 提示

### 2. 第二段：登录鉴权链路（4 分钟）

推荐按“前端 -> service -> token -> repo -> 返回前端”顺序讲。

关键文件：

- [../frontend/src/components/header/login.vue](../frontend/src/components/header/login.vue)
- [../backend/internal/account/service.go](../backend/internal/account/service.go)
- [../backend/internal/auth/token.go](../backend/internal/auth/token.go)
- [../backend/internal/middleware/auth/auth.go](../backend/internal/middleware/auth/auth.go)

你要讲清楚：

- 前端登录成功后会把 access token、refresh token、userInfo 写入本地
- 后端登录时先校验账号密码，再签发双 Token
- refresh token 带唯一 `refresh_id`
- 数据库存的是 `refresh_id` 的 hash
- 刷新 token 时做三重校验：token 合法、token_version 匹配、refresh hash 匹配
- 退出登录时会清空 refresh hash 并提升 token_version

高频加分点：

- 解释 `RotateRefreshTokenHash` 为什么像 CAS
- 解释为什么 access token 不必每次落库校验

### 3. 第三段：视频详情 + 评论 + 关注这条完整链路（5 分钟）

这是最适合讲“前后端协作”的部分。

关键文件：

- [../frontend/src/components/detail/videoDetail.vue](../frontend/src/components/detail/videoDetail.vue)
- [../frontend/src/components/video/videoMore.vue](../frontend/src/components/video/videoMore.vue)
- [../frontend/src/components/detail/floor.vue](../frontend/src/components/detail/floor.vue)
- [../backend/internal/video/service.go](../backend/internal/video/service.go)
- [../backend/internal/comment/repo.go](../backend/internal/comment/repo.go)
- [../backend/internal/social/repo.go](../backend/internal/social/repo.go)
- [../backend/internal/history/repo.go](../backend/internal/history/repo.go)

建议讲法：

1. 详情页初始化时并发拉详情、评论、推荐
2. 详情接口会返回 `viewer_state`，包括 liked/favorited/followed
3. 播放区组件读取 viewer_state 初始化按钮状态
4. 评论楼层组件会懒加载回复列表
5. 历史上报是 `(user_id, video_id)` 维度 upsert
6. 关注状态优先使用详情接口返回值，缺失时再回退单查

这一段里最值得面试官继续追问的点：

- 为什么要把 viewer_state 聚合返回
- 评论树为什么用 `root_id + parent_id`
- 历史记录为什么用 upsert 而不是 insert
- 关注计数为什么要和关系写入放在同一个事务

### 4. 第四段：投稿上传和审核流（3 分钟）

关键文件：

- [../frontend/src/views/upload/form.vue](../frontend/src/views/upload/form.vue)
- [../frontend/src/components/upload/videoUpload.vue](../frontend/src/components/upload/videoUpload.vue)
- [../frontend/src/views/management/videomanage1.vue](../frontend/src/views/management/videomanage1.vue)
- [../backend/internal/video/service.go](../backend/internal/video/service.go)
- [../backend/internal/admin/repo.go](../backend/internal/admin/repo.go)

核心讲法：

- 投稿拆成两步：先创建元数据，再上传媒体文件
- 这样可以避免大文件上传失败时元数据和媒体状态耦合太死
- 审核通过/驳回放在一个事务里完成：更新稿件状态、作者计数、通知、审核记录

### 5. 第五段：最后用“问题与改进”收尾（1 分钟）

建议诚实讲 3 点：

- 现在还是单体应用，后续可以拆搜索、互动、审核模块
- 媒体存储目前偏本地化，生产应切对象存储 + CDN
- 前端部分组件还偏页面式写法，可以继续做组合式封装和状态抽离

---

## 三、30 分钟版本：把项目讲成“架构设计 + 关键权衡 + 后续演进”

30 分钟版本适合二面、主管面、交叉面。思路是从“功能实现”升级到“设计取舍”。

### 1. 第一部分：为什么要做这个项目（3 分钟）

你要回答的不是“我写了个视频网站”，而是：

- 为什么选这个题材
- 为什么值得重构
- 为什么能体现工程能力

建议表述：

- 这个题材业务不复杂，但链路很完整：登录、内容、互动、上传、审核、通知都能覆盖
- 很适合展示一个后端同学有没有从“接口能跑”走到“系统能讲清楚”
- 我做这个项目的重点不是追求炫技，而是把典型互联网业务链路做完整、做清楚

### 2. 第二部分：系统结构和边界（5 分钟）

重点讲这几个边界：

- HTTP 层负责协议和参数绑定
- Service 层负责业务编排和状态流转
- Repository 层负责事务和 SQL 访问
- 中间件负责身份解析和上下文注入
- 前端 request 层负责请求公共逻辑

你可以强调：

- 这样做的价值是代码走读成本低、职责稳定、后续替换实现成本低
- 比如以后把 GORM 换掉，影响主要在 Repository
- 把 JWT 用户信息注入上下文后，Service 层不需要反复关心 Header 细节

### 3. 第三部分：登录态设计（5 分钟）

这一部分可以按“问题 -> 方案 -> 取舍”来讲。

问题：

- 单 JWT 登录存在续签和失效控制不灵活的问题

方案：

- access token + refresh token 双 Token
- refresh token 带 refresh_id
- refresh_id 落 hash
- token_version 支持全局失效
- rotate refresh hash 防并发重复续签

取舍：

- 方案比单 JWT 复杂，但安全性和可控性更好
- 仍然不是绝对安全，生产还要加设备信息、IP 风险识别、黑名单策略

### 4. 第四部分：详情页为什么是主链路（5 分钟）

你可以告诉面试官：

- 视频详情页是内容平台访问最重的页面之一
- 它天然会把详情、评论、推荐、点赞、收藏、关注、历史这些模块串起来
- 所以我把它当成项目主链路来设计

这里重点讲两件事：

#### 设计 1：viewer_state 聚合返回

- 后端详情接口直接返回 liked/favorited/followed
- 前端首屏渲染不需要多次补请求
- 接口契约更稳定

#### 设计 2：评论树建模

- 一级评论：`root_id = 0, parent_id = 0`
- 二级回复：`root_id = 一级评论 ID, parent_id = 直接回复对象`
- 这种方式兼顾查询和结构表达，适合两层评论场景

### 5. 第五部分：数据一致性和事务边界（4 分钟）

优先讲两个事务：

- 关注/取消关注事务
- 审核事务

你要说明：

- 关注关系和粉丝数必须一起更新，否则展示会错
- 审核不仅仅是改状态，还会影响作者公开视频数、通知消息、审核审计记录
- 因此这些动作必须放在一个事务边界里处理

### 6. 第六部分：前后端协作方式（3 分钟）

这里讲你对“接口契约”的理解。

重点说：

- 后端统一返回 Envelope，前端统一解包
- 后端尽量返回对页面友好的聚合结构
- 前端对新旧字段做兼容处理，例如 `viewer_state` 缺失时回退单独查询
- 这样可以降低前后端联调摩擦

### 7. 第七部分：项目不足和下一步演进（5 分钟）

这一段特别重要，能体现你不是“做完就算了”。

建议讲这几类不足：

#### 工程层

- 还缺更完整的 observability，比如 tracing、慢查询监控、统一审计日志
- 测试覆盖可以继续补 service/repo 边界场景

#### 架构层

- 单体结构适合当前阶段，但未来可以把搜索、审核、媒体处理拆出去
- Redis 的使用还可以更深入，比如缓存热点详情、关注状态等

#### 业务层

- 评论排序、推荐策略、历史进度更新还比较基础
- 上传链路还没有转码、封面抽帧、异步处理等能力

#### 安全层

- JWT 方案已经比简单登录更完整，但还可以加设备管理、风控、限流
- 媒体上传还可以加文件类型校验、大小限制、内容扫描

---

## 四、面试时建议优先讲的源码顺序

如果面试官让你“打开代码讲”，推荐就按这个顺序：

1. [../backend/cmd/api/main.go](../backend/cmd/api/main.go)
2. [../backend/internal/http/router.go](../backend/internal/http/router.go)
3. [../backend/internal/account/service.go](../backend/internal/account/service.go)
4. [../backend/internal/auth/token.go](../backend/internal/auth/token.go)
5. [../backend/internal/video/service.go](../backend/internal/video/service.go)
6. [../frontend/src/components/detail/videoDetail.vue](../frontend/src/components/detail/videoDetail.vue)
7. [../frontend/src/components/detail/floor.vue](../frontend/src/components/detail/floor.vue)
8. [../backend/internal/comment/repo.go](../backend/internal/comment/repo.go)
9. [../backend/internal/social/repo.go](../backend/internal/social/repo.go)
10. [../backend/internal/admin/repo.go](../backend/internal/admin/repo.go)

这个顺序的好处是：

- 先讲系统怎么启动
- 再讲请求怎么进来
- 再讲登录怎么保证身份
- 再讲核心业务怎么跑通
- 最后讲事务、一致性和后台能力

---

## 五、面试官可能追问的八个问题

### 1. 为什么 refresh token 要存 hash

回答关键词：

- 减少数据库泄露后的直接复用风险
- 类似密码不存明文的思路
- 服务端比对 hash 即可

### 2. 为什么视频详情要返回 viewer_state

回答关键词：

- 降低前端额外请求数
- 首屏状态更完整
- 接口契约更贴近页面需求

### 3. 评论为什么不做无限层级

回答关键词：

- 业务上两层最常见
- 查询和渲染复杂度更可控
- 先满足主流场景，再看是否演进

### 4. 审核为什么要用事务

回答关键词：

- 不只是更新一个字段
- 还涉及作者计数、通知、审计记录
- 任一步失败都不应该留下半成功状态

### 5. 为什么前端还要兼容旧字段

回答关键词：

- 便于渐进式重构
- 降低接口升级成本
- 有利于前后端并行开发

### 6. 这个项目最大的难点是什么

建议答：

- 不是某个算法难，而是怎么把多个看似小功能整理成一套能讲清楚的主链路和边界

### 7. 这个项目如果继续做，你最想优化什么

建议答：

- 上传转码异步化
- 缓存与监控体系补齐
- 评论/推荐/搜索做更真实的生产化设计

### 8. 你觉得这个项目最能体现你什么能力

建议答：

- 不是只会写接口，而是能把业务链路、数据结构、鉴权、安全、事务和前后端协作一起考虑

---

## 六、最后的使用建议

如果你准备把这个项目写进简历，建议至少做到三件事：

1. 先把 `5 分钟版本` 讲熟，保证不慌
2. 选定一条“王牌链路”，建议就是“登录态”或“视频详情页”
3. 准备 3 个不足点和 3 个演进方向，面试观感会明显更成熟

配套阅读：

- 简历项目描述与标准回答：[RESUME_PROJECT_WRITEUP_AND_QA.md](./RESUME_PROJECT_WRITEUP_AND_QA.md)
- 关键代码走读：[KEY_CODE_IMPLEMENTATION.md](./KEY_CODE_IMPLEMENTATION.md)
