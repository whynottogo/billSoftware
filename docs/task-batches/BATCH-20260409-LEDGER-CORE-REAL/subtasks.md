# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-LEDGER-CORE-REAL-01 | 用户级数据隔离与禁用会话失效 | SYS-DATA-001 | ["USER-LEDGER-001","USER-LEDGER-002","USER-LEDGER-003","USER-LEDGER-004"] | 已完成 | Codex | ["SYS-AUTH-001","SYS-AUTH-003"] | ["TC-ST-BATCH-20260409-LEDGER-CORE-REAL-01-01"] | ["BUG-BATCH-20260409-LEDGER-CORE-REAL-001"] | 2026-04-10 00:18 |
| ST-BATCH-20260409-LEDGER-CORE-REAL-02 | 当月收支真实读取与新增 | USER-LEDGER-001 | ["USER-LEDGER-002","USER-LEDGER-003"] | 已完成 | Codex | ["SYS-DATA-001","USER-AUTH-002"] | ["TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-01","TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-02"] | ["BUG-BATCH-20260409-LEDGER-CORE-REAL-001"] | 2026-04-10 00:18 |
| ST-BATCH-20260409-LEDGER-CORE-REAL-03 | 收支分类真实管理 | USER-LEDGER-004 | ["USER-LEDGER-001","USER-LEDGER-002","USER-LEDGER-003"] | 已完成 | Codex | ["SYS-DATA-001","USER-AUTH-002"] | ["TC-ST-BATCH-20260409-LEDGER-CORE-REAL-03-01"] | ["BUG-BATCH-20260409-LEDGER-CORE-REAL-001"] | 2026-04-10 00:18 |

## 子任务说明

### ST-BATCH-20260409-LEDGER-CORE-REAL-01 用户级数据隔离与禁用会话失效

```yaml
subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-01
title: 用户级数据隔离与禁用会话失效
primary_req_id: SYS-DATA-001
related_req_ids: ["USER-LEDGER-001","USER-LEDGER-002","USER-LEDGER-003","USER-LEDGER-004"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-AUTH-001","SYS-AUTH-003"]
implementation_notes:
  - 所有 ledger/category 查询与写入均按中间件注入的 user_id 过滤
  - 管理员禁用用户时同步失效该用户活跃 user_sessions
  - 已验证用户 1 无法看到用户 2 的流水，且被禁用后的旧会话访问用户接口返回 401
linked_test_case_ids: ["TC-ST-BATCH-20260409-LEDGER-CORE-REAL-01-01"]
linked_bug_ids: ["BUG-BATCH-20260409-LEDGER-CORE-REAL-001"]
last_updated_at: 2026-04-10 00:18
```

### ST-BATCH-20260409-LEDGER-CORE-REAL-02 当月收支真实读取与新增

```yaml
subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-02
title: 当月收支真实读取与新增
primary_req_id: USER-LEDGER-001
related_req_ids: ["USER-LEDGER-002","USER-LEDGER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001","USER-AUTH-002"]
implementation_notes:
  - 已接通 GET /api/user/ledger?month=YYYY-MM 与 POST /api/user/ledger
  - 响应结构按 UserHome 当前 `summary/groups/categories/overview` 形状返回
  - 已验证新增支出与新增收入后，月汇总、日分组和常用分类同步更新
linked_test_case_ids: ["TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-01","TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-02"]
linked_bug_ids: ["BUG-BATCH-20260409-LEDGER-CORE-REAL-001"]
last_updated_at: 2026-04-10 00:18
```

### ST-BATCH-20260409-LEDGER-CORE-REAL-03 收支分类真实管理

```yaml
subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-03
title: 收支分类真实管理
primary_req_id: USER-LEDGER-004
related_req_ids: ["USER-LEDGER-001","USER-LEDGER-002","USER-LEDGER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001","USER-AUTH-002"]
implementation_notes:
  - 已接通 GET /api/user/categories、POST /api/user/categories、DELETE /api/user/categories/:id
  - 已验证默认分类按钮禁用、自定义分类新增成功、自定义分类删除成功
  - 前端分类管理对话框已使用真实接口
linked_test_case_ids: ["TC-ST-BATCH-20260409-LEDGER-CORE-REAL-03-01"]
linked_bug_ids: ["BUG-BATCH-20260409-LEDGER-CORE-REAL-001"]
last_updated_at: 2026-04-10 00:18
```
