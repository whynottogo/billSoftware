# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-USER-LEDGER-ACTIONS-01 | 新增支出动作壳与即时回显 | USER-LEDGER-002 | ["USER-LEDGER-001"] | 已完成 | Codex | ["USER-LEDGER-001","SYS-DATA-001"] | ["TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-01-01"] | ["BUG-BATCH-20260409-USER-LEDGER-ACTIONS-001"] | 2026-04-09 22:08 |
| ST-BATCH-20260409-USER-LEDGER-ACTIONS-02 | 新增收入动作壳与即时回显 | USER-LEDGER-003 | ["USER-LEDGER-001"] | 已完成 | Codex | ["USER-LEDGER-001","SYS-DATA-001"] | ["TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-02-01"] | [] | 2026-04-09 22:08 |
| ST-BATCH-20260409-USER-LEDGER-ACTIONS-03 | 收支分类管理壳层与默认分类 | USER-LEDGER-004 | ["USER-LEDGER-002","USER-LEDGER-003"] | 已完成 | Codex | ["USER-LEDGER-002","USER-LEDGER-003"] | ["TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-03-01"] | ["BUG-BATCH-20260409-USER-LEDGER-ACTIONS-002"] | 2026-04-09 22:08 |

## 子任务说明模板

### ST-BATCH-20260409-USER-LEDGER-ACTIONS-01 新增支出动作壳与即时回显

```yaml
subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-01
title: 新增支出动作壳与即时回显
primary_req_id: USER-LEDGER-002
related_req_ids: ["USER-LEDGER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001","SYS-DATA-001"]
implementation_notes:
  - 在 /user/ledger 页面提供新增支出入口与表单弹层
  - 日期默认当天，图片仅支持 1 张占位
  - 保存后通过本地 mock 直接回写账单列表，不调用后端接口
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-01-01"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-LEDGER-ACTIONS-001"]
last_updated_at: 2026-04-09 22:08
```

Implementation Notes:

- 该子任务聚焦“可操作、可保存、可见结果”的原型闭环，不扩展真实上传和资产联动。

### ST-BATCH-20260409-USER-LEDGER-ACTIONS-02 新增收入动作壳与即时回显

```yaml
subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-02
title: 新增收入动作壳与即时回显
primary_req_id: USER-LEDGER-003
related_req_ids: ["USER-LEDGER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001","SYS-DATA-001"]
implementation_notes:
  - 在 /user/ledger 页面提供新增收入入口与表单弹层
  - 复用新增记录表单骨架，保持收入与支出表单字段一致
  - 保存后通过本地 mock 直接回写账单列表，不调用后端接口
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-02-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 22:08
```

Implementation Notes:

- 收入和支出表单统一交互结构，降低后续真实联调时的行为分歧。

### ST-BATCH-20260409-USER-LEDGER-ACTIONS-03 收支分类管理壳层与默认分类

```yaml
subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-03
title: 收支分类管理壳层与默认分类
primary_req_id: USER-LEDGER-004
related_req_ids: ["USER-LEDGER-002","USER-LEDGER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-002","USER-LEDGER-003"]
implementation_notes:
  - 在 /user/ledger 页面提供分类管理抽屉，区分收入分类和支出分类
  - 初始化加载默认分类，支持新增分类和删除自定义分类
  - 分类数据在本地 mock 中按用户会话隔离
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-03-01"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-LEDGER-ACTIONS-002"]
last_updated_at: 2026-04-09 22:08
```

Implementation Notes:

- 默认分类口径按需求文档设置，不复用历史“收支写反”的旧口径。
