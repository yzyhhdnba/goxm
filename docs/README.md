# Docs Index

当前仓库的文档按“契约优先、资料分组”的方式整理如下：

## 1. 根目录基线

- `README.md`：项目入口说明、当前状态与运行方式
- `AGENTS.md`：AI 代理与协作者约束
- `backend/README.md`：后端工作区说明
- `deploy/README.md`：部署与联调材料说明

## 2. 设计契约

目录：`docs/01-contracts/`

- `blueprint.md`：项目蓝图与整体技术路线
- `schema.md`：数据库模型与字段约束
- `api.md`：接口契约与分页规则

## 3. 开发资料

目录：`docs/02-development/`

- `technical-architecture.md`：独立技术设计文档，整理系统架构、模块边界与技术演进路线
- `task-breakdown.md`：阶段任务拆解
- `key-code-implementation.md`：关键实现走读
- `backend-tech-report.md`：后端技术报告与补充说明

## 4. 面试资料

目录：`docs/03-interview/`

- `knowledge-base-upload-guide.md`：面试知识库上传建议与筛选清单
- `interview-talk-track.md`：项目讲解顺序与口述稿
- `resume-project-writeup-and-qa.md`：简历写法与高频问答
- `frontend1-react-resume/`：并行 React 前端工程的简历、讲稿、问答与深度技术资料

## 5. 前端补充说明

目录：`docs/04-frontend/`

- `frontend-vue-cli-readme.md`：根目录 `frontend/` 的 Vue CLI 基础说明
- `backend-frontend-legacy-readme.md`：`backend/frontend/` 历史副本说明
- `legacy-patch-scripts.md`：前端遗留补丁脚本说明

## 6. 上传副本

目录：`上传/`

- `README.md`：上传目录用途说明
- `INDEX.md`：上传文件清单与用途索引

`上传/` 中的文件是为面试知识库准备的同步副本，不是仓库中的原始事实来源。

## 7. 使用建议

建议阅读顺序：

1. `README.md`
2. `docs/01-contracts/blueprint.md`
3. `docs/01-contracts/schema.md`
4. `docs/01-contracts/api.md`
5. `docs/02-development/technical-architecture.md`
6. `AGENTS.md`

如果后续新增 Markdown 文档，优先按以上分组收敛，不再把设计契约散落在仓库根目录。

如果你准备把这套文档上传成面试知识库，优先阅读：

1. `docs/03-interview/knowledge-base-upload-guide.md`
2. `docs/03-interview/interview-talk-track.md`
3. `docs/03-interview/resume-project-writeup-and-qa.md`
