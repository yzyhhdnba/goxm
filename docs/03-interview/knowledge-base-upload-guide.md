# PILIPILI Go 面试知识库上传指南

## 1. 这份文档的用途

这份文档专门服务于“把仓库文档上传为面试项目知识库”这个场景。

目标不是把所有 Markdown 一股脑上传，而是：

- 优先上传最能支撑项目介绍、追问回答、技术细节表达的文档
- 避免把低信息密度、重复说明、历史遗留或本地脚手架说明混进知识库
- 让知识库更像“项目答辩资料”，而不是“仓库杂项说明集合”

---

## 2. 推荐上传策略

建议按三层来组织：

### 第一层：必传

这些文件最适合作为面试知识库的核心内容。

1. `README.md`
2. `docs/01-contracts/blueprint.md`
3. `docs/01-contracts/schema.md`
4. `docs/01-contracts/api.md`
5. `docs/02-development/technical-architecture.md`
6. `docs/02-development/key-code-implementation.md`
7. `docs/03-interview/interview-talk-track.md`
8. `docs/03-interview/resume-project-writeup-and-qa.md`

原因：

- 这 8 份文档基本覆盖了“项目是什么、怎么设计、怎么实现、怎么讲出来”四个维度
- 既能回答项目总览问题，也能支撑深入追问
- 内容之间有明显分工，重复度相对可控

### 第二层：选传

这些文件适合在你希望知识库更完整时补充上传。

1. `docs/02-development/backend-tech-report.md`
2. `docs/02-development/task-breakdown.md`
3. `AGENTS.md`

原因：

- `backend-tech-report.md` 适合补充后端八股、设计权衡和延伸问答
- `task-breakdown.md` 适合回答“项目怎么推进”“阶段目标是什么”“你如何规划迭代”
- `AGENTS.md` 适合补充“项目边界、协作约束、哪些能力属于事实、哪些仍属规划”

注意：

- 这三份不是面试主入口
- 如果你的知识库容量有限，优先级明显低于“第一层必传”

### 第三层：不建议上传

以下文件通常不适合作为面试知识库主内容：

1. `docs/04-frontend/frontend-vue-cli-readme.md`
2. `docs/04-frontend/backend-frontend-legacy-readme.md`
3. `docs/04-frontend/legacy-patch-scripts.md`

原因：

- 信息密度低，偏本地脚手架或历史遗留说明
- 对面试表达帮助很小
- 容易稀释知识库质量，让检索结果出现无关内容

---

## 3. 最推荐的上传清单

如果你只想上传一套“够用、干净、面试友好”的知识库，推荐直接用下面这份清单：

1. `README.md`
2. `docs/01-contracts/blueprint.md`
3. `docs/01-contracts/schema.md`
4. `docs/01-contracts/api.md`
5. `docs/02-development/technical-architecture.md`
6. `docs/02-development/key-code-implementation.md`
7. `docs/03-interview/interview-talk-track.md`
8. `docs/03-interview/resume-project-writeup-and-qa.md`

这是当前最平衡的一组：

- `README.md` 负责入口与当前事实
- `01-contracts` 负责契约与系统边界
- `technical-architecture.md` 负责系统设计视角
- `key-code-implementation.md` 负责实现细节与源码入口
- `03-interview` 两份文档负责“如何讲”和“如何回答”

---

## 4. 如果你想压缩到最小集合

如果知识库只想放最核心内容，建议保留这 5 份：

1. `README.md`
2. `docs/01-contracts/blueprint.md`
3. `docs/02-development/technical-architecture.md`
4. `docs/02-development/key-code-implementation.md`
5. `docs/03-interview/resume-project-writeup-and-qa.md`

适用情况：

- 上传额度有限
- 你希望知识库更偏“讲项目”和“回答追问”
- 你不想上传太多契约细节和任务拆解

---

## 5. 上传顺序建议

如果上传工具支持分批导入，建议按下面顺序上传：

1. `README.md`
2. `docs/01-contracts/blueprint.md`
3. `docs/01-contracts/schema.md`
4. `docs/01-contracts/api.md`
5. `docs/02-development/technical-architecture.md`
6. `docs/02-development/key-code-implementation.md`
7. `docs/03-interview/interview-talk-track.md`
8. `docs/03-interview/resume-project-writeup-and-qa.md`

这个顺序的好处是：

- 先让知识库理解项目定位和边界
- 再补结构、接口、实现细节
- 最后补面试表达层内容

---

## 6. 上传时的实际建议

### 建议保留的文档风格

- 项目事实
- 当前已实现能力
- 技术取舍
- 主链路实现
- 面试问答模板

### 建议避免的内容

- 纯本地脚手架命令
- 历史遗留补丁说明
- 与面试无关的低层运行噪音
- 重复表述太多但信息增量很低的文档

### 如果知识库支持标签

建议按下面方式给文档打标签：

- `project-overview`
- `architecture`
- `schema`
- `api`
- `implementation`
- `interview`
- `qa`

---

## 7. 一句话结论

如果你的目标是“把这套文档当成面试时可检索、可追问、可展开的项目知识库”，最推荐上传的是：

- `README.md`
- `docs/01-contracts/blueprint.md`
- `docs/01-contracts/schema.md`
- `docs/01-contracts/api.md`
- `docs/02-development/technical-architecture.md`
- `docs/02-development/key-code-implementation.md`
- `docs/03-interview/interview-talk-track.md`
- `docs/03-interview/resume-project-writeup-and-qa.md`
