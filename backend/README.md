# Backend Notes

`backend/` 是当前 Go 后端主工作区。

## 目录用途

- `cmd/api/`
  API 进程入口，负责配置加载、依赖初始化、迁移、路由装配与优雅停机。

- `cmd/worker/`
  Worker 进程入口。
  当前阶段主要完成运行骨架、配置接入与 RocketMQ 基础连通能力，尚未接入完整业务消费链路。

- `configs/`
  配置模板与本地运行配置入口。
  `config.example.yaml` 适合提交到仓库，`config.yaml` 仅用于本地联调。

- `internal/`
  业务域与基础设施代码。

- `pkg/`
  可跨业务域复用的公共包。

- `scripts/`
  本地联调与 E2E 验证脚本。

- `storage/`
  本地媒体落盘目录，默认不作为仓库契约的一部分。

- `bin/`
  本地构建产物目录，默认不提交。

## 当前约束

- 当前默认以后端主链路可跑通为优先，不把 Worker / MQ 规划写成已完成业务能力。
- 新增后端逻辑时，优先保持 `handler -> service -> repository` 分层。
- 当前评论、回复仍采用统一建模，不要拆成两套表和两套服务。

## 启动提示

API：

```bash
cd backend
GOCACHE=/tmp/pilipili-go-build-cache GOMODCACHE=/tmp/pilipili-go-mod-cache go run ./cmd/api -config configs/config.yaml
```

Worker：

```bash
cd backend
GOCACHE=/tmp/pilipili-go-build-cache GOMODCACHE=/tmp/pilipili-go-mod-cache go run ./cmd/worker -config configs/config.yaml
```
