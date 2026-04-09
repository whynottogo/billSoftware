# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-USER-CHART-01 | 支出图表页与排行榜落地 | USER-CHART-001 | ["USER-LEDGER-004"] | 已完成 | Codex | ["USER-LEDGER-001","USER-LEDGER-004"] | ["TC-ST-BATCH-20260409-USER-CHART-01-01"] | [] | 2026-04-09 22:21 |
| ST-BATCH-20260409-USER-CHART-02 | 收入图表页与排行榜落地 | USER-CHART-002 | ["USER-LEDGER-004"] | 已完成 | Codex | ["USER-LEDGER-001","USER-LEDGER-004"] | ["TC-ST-BATCH-20260409-USER-CHART-02-01"] | [] | 2026-04-09 22:21 |

## 子任务说明模板

### ST-BATCH-20260409-USER-CHART-01 支出图表页与排行榜落地

```yaml
subtask_id: ST-BATCH-20260409-USER-CHART-01
title: 支出图表页与排行榜落地
primary_req_id: USER-CHART-001
related_req_ids: ["USER-LEDGER-004"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001","USER-LEDGER-004"]
implementation_notes:
  - 在 /user/charts/expense 页面提供月统计折线图、年统计折线图和支出前十排行
  - 优先复用稳定 SVG 图表组件，不引入 ECharts 运行风险
  - 图表与排行榜数据来自本地 chart mock
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-CHART-01-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 22:21
```

Implementation Notes:

- 该子任务重点是结构与交互完整度，非接口真实性。

### ST-BATCH-20260409-USER-CHART-02 收入图表页与排行榜落地

```yaml
subtask_id: ST-BATCH-20260409-USER-CHART-02
title: 收入图表页与排行榜落地
primary_req_id: USER-CHART-002
related_req_ids: ["USER-LEDGER-004"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001","USER-LEDGER-004"]
implementation_notes:
  - 新增 /user/charts/income 路由与页面表现
  - 提供按年收入趋势折线和收入前十排行
  - 收入页支持年份切换并保持与支出页一致的视觉语言
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-CHART-02-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 22:21
```

Implementation Notes:

- 收入页是独立路由，不与支出页共用单页切换壳。
