# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260410-BUDGET-CHART-REAL-01 | 后端预算与图表共享统计接口 | USER-BUDGET-001 | ["USER-BUDGET-002","USER-CHART-001","USER-CHART-002"] | 已完成 | Codex | ["USER-LEDGER-001","USER-LEDGER-004"] | ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-01","TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-02"] | [] | 2026-04-10 12:54 |
| ST-BATCH-20260410-BUDGET-CHART-REAL-02 | 用户月预算与年预算真实接线 | USER-BUDGET-002 | ["USER-BUDGET-001"] | 已完成 | Codex | ["USER-BUDGET-001","USER-LEDGER-004"] | ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-01","TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-02"] | ["BUG-BATCH-20260410-BUDGET-CHART-REAL-001"] | 2026-04-10 12:54 |
| ST-BATCH-20260410-BUDGET-CHART-REAL-03 | 用户支出图表真实接线 | USER-CHART-001 | ["USER-BUDGET-001","USER-LEDGER-001"] | 已完成 | Codex | ["USER-BUDGET-001","USER-LEDGER-001"] | ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-03-01"] | [] | 2026-04-10 12:54 |
| ST-BATCH-20260410-BUDGET-CHART-REAL-04 | 用户收入图表真实接线 | USER-CHART-002 | ["USER-LEDGER-001"] | 已完成 | Codex | ["USER-BUDGET-001","USER-LEDGER-001"] | ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-04-01"] | [] | 2026-04-10 12:54 |
| ST-BATCH-20260410-BUDGET-CHART-REAL-05 | 预算图表闭环联调与回归 | USER-BUDGET-001 | ["USER-BUDGET-002","USER-CHART-001","USER-CHART-002"] | 已完成 | Codex | ["USER-BUDGET-001","USER-BUDGET-002","USER-CHART-001","USER-CHART-002"] | ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-01","TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-02"] | ["BUG-BATCH-20260410-BUDGET-CHART-REAL-002"] | 2026-04-10 12:54 |

## 子任务说明

### ST-BATCH-20260410-BUDGET-CHART-REAL-01 后端预算与图表共享统计接口

```yaml
subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-01
title: 后端预算与图表共享统计接口
primary_req_id: USER-BUDGET-001
related_req_ids: ["USER-BUDGET-002","USER-CHART-001","USER-CHART-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001","USER-LEDGER-004"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-01","TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 12:54
```

Implementation Notes:

- 这是唯一允许修改后端共享统计口径、预算持久化、路由接线和聚合 SQL 的子任务
- 固定接口契约如下：
  - `GET /api/user/budgets/month/current`
  - `PUT /api/user/budgets/month/current`
  - `GET /api/user/budgets/year/options`
  - `GET /api/user/budgets/year/:year`
  - `PUT /api/user/budgets/year/current`
  - `GET /api/user/charts/years`
  - `GET /api/user/charts/expense/:year`
  - `GET /api/user/charts/income/:year`
- 月预算与年预算都基于 `budgets` + `budget_items` 持久化，分类维度限定为当前用户的支出分类
- 月预算只允许设置当前月份，年预算只允许设置当前年份；图表聚合严格按当前登录 `user_id` 隔离

### ST-BATCH-20260410-BUDGET-CHART-REAL-02 用户月预算与年预算真实接线

```yaml
subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-02
title: 用户月预算与年预算真实接线
primary_req_id: USER-BUDGET-002
related_req_ids: ["USER-BUDGET-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BUDGET-001","USER-LEDGER-004"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-01","TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-02"]
linked_bug_ids: ["BUG-BATCH-20260410-BUDGET-CHART-REAL-001"]
last_updated_at: 2026-04-10 12:54
```

Implementation Notes:

- 仅负责 `UserBudgetMonth.vue`、`UserBudgetYear.vue` 及独立预算 API 模块
- 已把“设置预算 / 调整预算”从占位消息改为真实交互，并完成保存后即时刷新
- 页面消费结构保持：
  - 月预算：`{ month, label, notice, overview, highlights, categories }`
  - 年预算：`{ year, overview, monthlyExecution, categories }`
- 2026-04-10 11:52 前端构建异常已登记为 `BUG-BATCH-20260410-BUDGET-CHART-REAL-001`，修复后已完成回归关闭

### ST-BATCH-20260410-BUDGET-CHART-REAL-03 用户支出图表真实接线

```yaml
subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-03
title: 用户支出图表真实接线
primary_req_id: USER-CHART-001
related_req_ids: ["USER-BUDGET-001","USER-LEDGER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BUDGET-001","USER-LEDGER-001"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-03-01"]
linked_bug_ids: []
last_updated_at: 2026-04-10 12:54
```

Implementation Notes:

- 仅负责 `UserCharts.vue` 的支出模式和图表 API 模块中的支出部分
- 支出图表接口结构固定为 `{ year, summary, monthTrend, yearTrend, ranking }`
- 前端继续复用现有 SVG 图表组件，未引入新的图表依赖

### ST-BATCH-20260410-BUDGET-CHART-REAL-04 用户收入图表真实接线

```yaml
subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-04
title: 用户收入图表真实接线
primary_req_id: USER-CHART-002
related_req_ids: ["USER-LEDGER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BUDGET-001","USER-LEDGER-001"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-04-01"]
linked_bug_ids: []
last_updated_at: 2026-04-10 12:54
```

Implementation Notes:

- 仅负责 `UserCharts.vue` 的收入模式和图表 API 模块中的收入部分
- 收入图表接口结构固定为 `{ year, summary, yearTrend, ranking }`
- 已复用现有年份切换和模式切换，不新增页面入口

### ST-BATCH-20260410-BUDGET-CHART-REAL-05 预算图表闭环联调与回归

```yaml
subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-05
title: 预算图表闭环联调与回归
primary_req_id: USER-BUDGET-001
related_req_ids: ["USER-BUDGET-002","USER-CHART-001","USER-CHART-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BUDGET-001","USER-BUDGET-002","USER-CHART-001","USER-CHART-002"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-01","TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-02"]
linked_bug_ids: ["BUG-BATCH-20260410-BUDGET-CHART-REAL-002"]
last_updated_at: 2026-04-10 12:54
```

Implementation Notes:

- 主线程负责 Playwright 联调、测试证据归档、Bug 台账、需求状态和批次状态收口
- 核心闭环“登录 -> 设置月预算 -> 查看年预算 -> 查看支出图表 -> 切换收入图表”已通过
- 发现的浏览器 MCP 初始化问题已先登记 `BUG-BATCH-20260410-BUDGET-CHART-REAL-002`，随后切换 wrapper 完成真实回归
