# 测试计划

```yaml
batch_id: BATCH-20260409-REQ-STATUS-REALIGN
status: 已完成
environment:
  workspace: /Users/shynin/software/billSoftware
  frontend: source audit only
  backend: source audit only
  database: model / handler review only
test_scope:
  - docs/requirements/bill-software-requirements.md
  - frontend/src/router/index.js
  - frontend/src/utils/request.js
  - frontend/src/api/userAuth.js
  - frontend/src/pages/user
  - frontend/src/pages/admin
  - frontend/src/utils/*Mock.js
  - backend/internal/router/router.go
  - backend/internal/handler
  - backend/internal/middleware
entry_points:
  - docs/requirements/bill-software-requirements.md#8 协同任务清单
  - docs/task-batches/BATCH-20260409-*/overview.md
  - frontend/src/router/index.js
  - frontend/src/utils/request.js
  - backend/internal/router/router.go
verification_strategy:
  - 逐条对照 req 的 Detail / Acceptance Criteria 与真实源码实现
  - 若 Notes 明确允许 mock 或原型级页面，则 mock 页面可作为 V1 完成依据
  - 若仅完成前端原型或本地 mock，且需求未明确允许 mock，则状态回写为开发中
  - 若前后端均有实现但前端尚未接线，则状态回写为联调中
risks:
  - 少数需求存在“页面表现已完成但真实数据未接入”的灰区，需要严格依赖 Notes 是否允许 mock 判定
  - 本批次不启动服务做实时联调，结论依赖源码与批次记录的一致性
last_updated_at: 2026-04-09 23:18
```

## 测试范围

- 覆盖页面：用户端登录/注册、个人信息、账单/预算/资产/图表/家庭、管理端登录/门户/用户管理
- 覆盖接口：用户注册、用户登录、管理员登录、管理端用户列表、用户启停、用户汇总
- 覆盖核心流程：状态口径审计、需求回退判断、模块总览状态收口

## 入口与前置条件

- 页面入口：不启动浏览器，直接审查 Vue 路由、页面实现和 mock util
- 登录账号：不适用
- 初始化数据：不适用
- 环境依赖：可读取完整工作区源码与批次文档

## 执行顺序

1. 读取需求任务池与既有批次记录，确认当前状态来源
2. 审查前端 API 接线、mock util 使用范围和路由守卫
3. 审查后端 router、handler、middleware 覆盖范围
4. 按验收项判断 req 的真实完成度并形成修正结论
5. 回写需求文档、模块状态和状态口径说明

## 验证策略

- 文档纠偏批次以源码审计为主，不做 Playwright 页面回归
- 认证类 req 重点验证真实接口和真实会话逻辑是否存在
- 业务类 req 重点验证是否仍依赖 `*Mock.js` 或 `localStorage` 本地演示
- 管理端 req 额外区分“Notes 允许 mock 的完成态”和“需真实用户数据的联调态”

## 风险记录

- 风险 1：如果后续有人继续沿用旧口径回写状态，需求总表会再次失真
- 风险 2：现有批次文档大多聚焦前端原型落地，后续补做接口联调时需要单独创建新批次承接
