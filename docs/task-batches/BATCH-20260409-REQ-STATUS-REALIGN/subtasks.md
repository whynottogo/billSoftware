# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-REQ-STATUS-REALIGN-01 | 审计系统与认证链路真实完成度 | SYS-AUTH-002 | ["SYS-ARCH-001","SYS-AUTH-001","SYS-AUTH-003","SYS-DATA-001","USER-AUTH-001","USER-AUTH-002","USER-AUTH-003","ADMIN-AUTH-001"] | 已完成 | Codex | [] | ["TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-01-01"] | [] | 2026-04-09 23:18 |
| ST-BATCH-20260409-REQ-STATUS-REALIGN-02 | 审计用户端业务模块真实完成度 | USER-LEDGER-001 | ["USER-LEDGER-002","USER-LEDGER-003","USER-LEDGER-004","USER-BILL-001","USER-BILL-002","USER-BILL-003","USER-BUDGET-001","USER-BUDGET-002","USER-ASSET-001","USER-ASSET-002","USER-ASSET-003","USER-CHART-001","USER-CHART-002","USER-PROFILE-001","USER-PROFILE-002","USER-FAMILY-001","USER-FAMILY-002"] | 已完成 | Codex | ["SYS-DATA-001"] | ["TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-02-01"] | [] | 2026-04-09 23:18 |
| ST-BATCH-20260409-REQ-STATUS-REALIGN-03 | 审计管理端模块并回写需求与模块状态 | ADMIN-USER-001 | ["ADMIN-DASHBOARD-001","ADMIN-APPROVAL-001","ADMIN-FAMILY-001","ADMIN-USER-002","ADMIN-USER-003"] | 已完成 | Codex | ["ADMIN-AUTH-001"] | ["TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-03-01"] | [] | 2026-04-09 23:18 |

## 子任务说明

### ST-BATCH-20260409-REQ-STATUS-REALIGN-01 审计系统与认证链路真实完成度

```yaml
subtask_id: ST-BATCH-20260409-REQ-STATUS-REALIGN-01
title: 审计系统与认证链路真实完成度
primary_req_id: SYS-AUTH-002
related_req_ids: ["SYS-ARCH-001","SYS-AUTH-001","SYS-AUTH-003","SYS-DATA-001","USER-AUTH-001","USER-AUTH-002","USER-AUTH-003","ADMIN-AUTH-001"]
status: 已完成
owner_ai: Codex
dependencies: []
implementation_notes:
  - 对照 router、request interceptor、auth handler 和 middleware 审查双端路由隔离、token 隔离、用户注册登录、管理员登录和单会话互踢链路
  - 将满足验收项的认证类 req 保持或修正为已完成
  - 对未覆盖完整数据权限范围的 SYS-DATA-001 维持开发中
linked_test_case_ids: ["TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-01-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 23:18
```

Implementation Notes:

- 本子任务重点纠正“互踢已做但状态仍停留开发中”的偏差，同时避免把未落地的数据权限控制误判为完成。

### ST-BATCH-20260409-REQ-STATUS-REALIGN-02 审计用户端业务模块真实完成度

```yaml
subtask_id: ST-BATCH-20260409-REQ-STATUS-REALIGN-02
title: 审计用户端业务模块真实完成度
primary_req_id: USER-LEDGER-001
related_req_ids: ["USER-LEDGER-002","USER-LEDGER-003","USER-LEDGER-004","USER-BILL-001","USER-BILL-002","USER-BILL-003","USER-BUDGET-001","USER-BUDGET-002","USER-ASSET-001","USER-ASSET-002","USER-ASSET-003","USER-CHART-001","USER-CHART-002","USER-PROFILE-001","USER-PROFILE-002","USER-FAMILY-001","USER-FAMILY-002"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001"]
implementation_notes:
  - 核对用户端页面是否接入真实 API，是否仍由 localStorage/mock util 驱动
  - 对仅完成原型级页面、未满足真实业务验收项的 req 统一回退为开发中
  - 保留页面资产不动，只纠正文档状态和口径说明
linked_test_case_ids: ["TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-02-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 23:18
```

Implementation Notes:

- 本子任务不否定已完成的前端原型工作，只是把状态从“需求已完成”校正为“业务仍在开发中”。

### ST-BATCH-20260409-REQ-STATUS-REALIGN-03 审计管理端模块并回写需求与模块状态

```yaml
subtask_id: ST-BATCH-20260409-REQ-STATUS-REALIGN-03
title: 审计管理端模块并回写需求与模块状态
primary_req_id: ADMIN-USER-001
related_req_ids: ["ADMIN-DASHBOARD-001","ADMIN-APPROVAL-001","ADMIN-FAMILY-001","ADMIN-USER-002","ADMIN-USER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-AUTH-001"]
implementation_notes:
  - 保留 Notes 明确允许 mock 的后台门户 req 为已完成
  - 将已有真实后端接口但前端尚未接线的管理端用户任务调整为联调中
  - 将真实后端尚未具备、页面仍为 mock 详情的管理端用户账单详情回退为开发中，并同步模块总览状态
linked_test_case_ids: ["TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-03-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 23:18
```

Implementation Notes:

- 本子任务区分“Notes 允许 mock 的完成态”和“需要真实用户数据但仍未接线的联调态”，避免后台任务再次混标。
