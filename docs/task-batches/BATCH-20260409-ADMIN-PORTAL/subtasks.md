# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-ADMIN-PORTAL-01 | 后台数据概览 Make 化 | ADMIN-DASHBOARD-001 | ["ADMIN-APPROVAL-001","ADMIN-FAMILY-001"] | 已完成 | Codex | ["ADMIN-AUTH-001"] | ["TC-ST-BATCH-20260409-ADMIN-PORTAL-01-01","TC-ST-BATCH-20260409-ADMIN-PORTAL-01-02"] | ["BUG-BATCH-20260409-ADMIN-PORTAL-001"] | 2026-04-09 21:46 |
| ST-BATCH-20260409-ADMIN-PORTAL-02 | 后台待审批用户页 Make 化 | ADMIN-APPROVAL-001 | ["ADMIN-DASHBOARD-001","ADMIN-USER-001"] | 已完成 | Codex | ["ADMIN-AUTH-001","ADMIN-USER-001"] | ["TC-ST-BATCH-20260409-ADMIN-PORTAL-02-01"] | ["BUG-BATCH-20260409-ADMIN-PORTAL-002"] | 2026-04-09 21:46 |
| ST-BATCH-20260409-ADMIN-PORTAL-03 | 后台家庭管理页 Make 化 | ADMIN-FAMILY-001 | ["ADMIN-DASHBOARD-001"] | 已完成 | Codex | ["ADMIN-AUTH-001","USER-FAMILY-001"] | ["TC-ST-BATCH-20260409-ADMIN-PORTAL-03-01"] | [] | 2026-04-09 21:46 |

## 子任务说明模板

### ST-BATCH-20260409-ADMIN-PORTAL-01 后台数据概览 Make 化

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-01
title: 后台数据概览 Make 化
primary_req_id: ADMIN-DASHBOARD-001
related_req_ids: ["ADMIN-APPROVAL-001","ADMIN-FAMILY-001"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-AUTH-001"]
implementation_notes:
  - 以 Figma Make 的 AdminDashboard 为参考，重建后台概览信息层级、统计卡和趋势区块
  - 使用稳定的 Vue/SVG 图表能力替代易引发路由切换异常的方案
  - 页面保留进入用户管理入口，便于后续共享壳集成
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-PORTAL-01-01","TC-ST-BATCH-20260409-ADMIN-PORTAL-01-02"]
linked_bug_ids: ["BUG-BATCH-20260409-ADMIN-PORTAL-001"]
last_updated_at: 2026-04-09 21:46
```

Implementation Notes:

- 该子任务承担后台门户主视觉基调，决定其他两页的复用样式。

### ST-BATCH-20260409-ADMIN-PORTAL-02 后台待审批用户页 Make 化

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-02
title: 后台待审批用户页 Make 化
primary_req_id: ADMIN-APPROVAL-001
related_req_ids: ["ADMIN-DASHBOARD-001","ADMIN-USER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-AUTH-001","ADMIN-USER-001"]
implementation_notes:
  - 以 Figma Make 的 PendingApproval 为参考，构建审批列表、状态标识和批量操作区域
  - 审批动作先提供原型级反馈，不接入真实审批接口
  - 页面结构保持列表密度和后台工作台节奏
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-PORTAL-02-01"]
linked_bug_ids: ["BUG-BATCH-20260409-ADMIN-PORTAL-002"]
last_updated_at: 2026-04-09 21:46
```

Implementation Notes:

- 该子任务优先保证信息可读性和操作反馈，不扩展到审批业务联调。

### ST-BATCH-20260409-ADMIN-PORTAL-03 后台家庭管理页 Make 化

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-03
title: 后台家庭管理页 Make 化
primary_req_id: ADMIN-FAMILY-001
related_req_ids: ["ADMIN-DASHBOARD-001"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-AUTH-001","USER-FAMILY-001"]
implementation_notes:
  - 以 Figma Make 的 FamilyManagement 为参考，重建家庭统计与家庭卡片列表
  - 详情、更多操作和分页先做交互壳层与状态提示
  - 总资产与成员指标采用 mock 数据，保持后续接口接入扩展点
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-PORTAL-03-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:46
```

Implementation Notes:

- 该子任务用于补齐后台菜单中家庭管理入口的可用页面承载。
