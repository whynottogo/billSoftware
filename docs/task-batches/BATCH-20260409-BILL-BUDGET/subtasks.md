# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-BILL-BUDGET-01 | 用户月账单列表 Make 化 | USER-BILL-001 | ["USER-BILL-002","USER-BILL-003"] | 已完成 | Codex | ["USER-LEDGER-001"] | ["TC-ST-BATCH-20260409-BILL-BUDGET-01-01"] | [] | 2026-04-09 21:47 |
| ST-BATCH-20260409-BILL-BUDGET-02 | 用户月账单详情 Make 化 | USER-BILL-002 | ["USER-BILL-001"] | 已完成 | Codex | ["USER-BILL-001"] | ["TC-ST-BATCH-20260409-BILL-BUDGET-02-01"] | ["BUG-BATCH-20260409-BILL-BUDGET-001"] | 2026-04-09 21:47 |
| ST-BATCH-20260409-BILL-BUDGET-03 | 用户年账单 Make 化 | USER-BILL-003 | ["USER-BILL-001","USER-BILL-002"] | 已完成 | Codex | ["USER-BILL-001"] | ["TC-ST-BATCH-20260409-BILL-BUDGET-03-01"] | ["BUG-BATCH-20260409-BILL-BUDGET-002"] | 2026-04-09 21:47 |
| ST-BATCH-20260409-BILL-BUDGET-04 | 用户月预算 Make 化 | USER-BUDGET-001 | ["USER-BUDGET-002"] | 已完成 | Codex | ["USER-LEDGER-001"] | ["TC-ST-BATCH-20260409-BILL-BUDGET-04-01"] | [] | 2026-04-09 21:47 |
| ST-BATCH-20260409-BILL-BUDGET-05 | 用户年预算 Make 化 | USER-BUDGET-002 | ["USER-BUDGET-001"] | 已完成 | Codex | ["USER-BUDGET-001"] | ["TC-ST-BATCH-20260409-BILL-BUDGET-05-01"] | [] | 2026-04-09 21:47 |

## 子任务说明模板

### ST-BATCH-20260409-BILL-BUDGET-01 用户月账单列表 Make 化

```yaml
subtask_id: ST-BATCH-20260409-BILL-BUDGET-01
title: 用户月账单列表 Make 化
primary_req_id: USER-BILL-001
related_req_ids: ["USER-BILL-002","USER-BILL-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001"]
implementation_notes:
  - 以 Figma Make 中 MonthlyBillList 为主要参考，重做账单首页的信息节奏与月份卡片
  - 页面顶部展示年度总览，并提供进入月账单详情页与年账单页的入口
  - 保留 Vue Router 结构，不在本批次强接真实账单接口
linked_test_case_ids: ["TC-ST-BATCH-20260409-BILL-BUDGET-01-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:47
```

Implementation Notes:

- 该子任务重点是把账单首页从占位页替换为可浏览的年度总览和月份列表。

### ST-BATCH-20260409-BILL-BUDGET-02 用户月账单详情 Make 化

```yaml
subtask_id: ST-BATCH-20260409-BILL-BUDGET-02
title: 用户月账单详情 Make 化
primary_req_id: USER-BILL-002
related_req_ids: ["USER-BILL-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BILL-001"]
implementation_notes:
  - 以 Figma Make 中 MonthlyBillDetail 为主要参考，补齐月汇总、分类构成、排行、趋势和成就区块
  - 使用静态数据驱动图表和卡片，保留后续替换为真实接口的空间
  - 通过动态路由承接月份详情浏览
linked_test_case_ids: ["TC-ST-BATCH-20260409-BILL-BUDGET-02-01"]
linked_bug_ids: ["BUG-BATCH-20260409-BILL-BUDGET-001"]
last_updated_at: 2026-04-09 21:47
```

Implementation Notes:

- 该页面需要兼顾信息密度和可读性，视觉上延续第一批的品牌语言。
- 已完成图表组件替换并通过回归，`BUG-BATCH-20260409-BILL-BUDGET-001` 已关闭。

### ST-BATCH-20260409-BILL-BUDGET-03 用户年账单 Make 化

```yaml
subtask_id: ST-BATCH-20260409-BILL-BUDGET-03
title: 用户年账单 Make 化
primary_req_id: USER-BILL-003
related_req_ids: ["USER-BILL-001","USER-BILL-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BILL-001"]
implementation_notes:
  - 以 Figma Make 中 YearlyBill 为主要参考，重做年度切换、总览英雄区和历年账单列表
  - 当前批次使用前端静态数据模拟多年份切换
  - 页面必须与月账单详情保持差异化，不做简单复用复制
linked_test_case_ids: ["TC-ST-BATCH-20260409-BILL-BUDGET-03-01"]
linked_bug_ids: ["BUG-BATCH-20260409-BILL-BUDGET-002"]
last_updated_at: 2026-04-09 21:47
```

Implementation Notes:

- 年账单页面已完成年份切换实测，受共享编译阻塞影响的 `BUG-BATCH-20260409-BILL-BUDGET-002` 已关闭。

### ST-BATCH-20260409-BILL-BUDGET-04 用户月预算 Make 化

```yaml
subtask_id: ST-BATCH-20260409-BILL-BUDGET-04
title: 用户月预算 Make 化
primary_req_id: USER-BUDGET-001
related_req_ids: ["USER-BUDGET-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001"]
implementation_notes:
  - 以 Figma Make 中 MonthlyBudget 为参考重做本月预算总览和分类预算卡片
  - 页面强调本月预算、本月支出和剩余额度，并保留细分预算设置入口
  - 当前批次不引入真实表单提交，只提供清晰的交互占位和信息结构
linked_test_case_ids: ["TC-ST-BATCH-20260409-BILL-BUDGET-04-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:47
```

Implementation Notes:

- 预算页需要体现“控制”和“剩余额度”的节奏，和账单页形成用途区分。

### ST-BATCH-20260409-BILL-BUDGET-05 用户年预算 Make 化

```yaml
subtask_id: ST-BATCH-20260409-BILL-BUDGET-05
title: 用户年预算 Make 化
primary_req_id: USER-BUDGET-002
related_req_ids: ["USER-BUDGET-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BUDGET-001"]
implementation_notes:
  - 以 Figma Make 中 YearlyBudget 为参考重做年度总览、月度执行和年度分类表现
  - 页面提供年份切换和年预算维度的汇总，不复制月预算布局
  - 保持与用户端整体壳层和导航样式一致
linked_test_case_ids: ["TC-ST-BATCH-20260409-BILL-BUDGET-05-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:47
```

Implementation Notes:

- 年预算页强调年度分配节奏和执行趋势，避免变成月预算的平铺复制。
