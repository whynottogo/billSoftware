# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-ADMIN-USER-01 | 管理端用户列表启用/禁用交互实现 | ADMIN-USER-002 | ["ADMIN-USER-003"] | 已完成 | Codex | ["ADMIN-USER-001"] | ["TC-ST-BATCH-20260409-ADMIN-USER-01-01"] | ["BUG-BATCH-20260409-ADMIN-USER-001"] | 2026-04-09 22:04 |
| ST-BATCH-20260409-ADMIN-USER-02 | 管理端用户账单详情页（月/年 tab）实现 | ADMIN-USER-003 | ["ADMIN-USER-002","USER-BILL-001","USER-BILL-002","USER-BILL-003"] | 已完成 | Codex | ["ADMIN-USER-001"] | ["TC-ST-BATCH-20260409-ADMIN-USER-02-01"] | [] | 2026-04-09 22:04 |
| ST-BATCH-20260409-ADMIN-USER-03 | 管理端用户模块回归测试与状态收口 | ADMIN-USER-003 | ["ADMIN-USER-002"] | 已完成 | Codex | ["ST-BATCH-20260409-ADMIN-USER-01","ST-BATCH-20260409-ADMIN-USER-02"] | ["TC-ST-BATCH-20260409-ADMIN-USER-03-01"] | ["BUG-BATCH-20260409-ADMIN-USER-001"] | 2026-04-09 22:04 |

## 子任务说明模板

### ST-BATCH-20260409-ADMIN-USER-01 管理端用户列表启用/禁用交互实现

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-USER-01
title: 管理端用户列表启用/禁用交互实现
primary_req_id: ADMIN-USER-002
related_req_ids: ["ADMIN-USER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-USER-001"]
implementation_notes:
  - 补齐用户列表启用/禁用按钮交互，切换状态后即时更新列表行和顶部统计
  - 交互提示使用前端消息提示，数据以 mock 维护
  - 用户详情按钮跳转到 /admin/users/:userId
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-USER-01-01"]
linked_bug_ids: ["BUG-BATCH-20260409-ADMIN-USER-001"]
last_updated_at: 2026-04-09 22:04
```

### ST-BATCH-20260409-ADMIN-USER-02 管理端用户账单详情页（月/年 tab）实现

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-USER-02
title: 管理端用户账单详情页（月/年 tab）实现
primary_req_id: ADMIN-USER-003
related_req_ids: ["ADMIN-USER-002","USER-BILL-001","USER-BILL-002","USER-BILL-003"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-USER-001"]
implementation_notes:
  - 新增 AdminUserDetail 页面，顶部展示用户基础信息和账单汇总
  - 页面包含“月账单”“年账单”两个 tab，并提供月份/年份切换
  - 数据以 mock 实现并与用户端账单汇总口径保持一致（收入/支出/结余）
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-USER-02-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 22:04
```

### ST-BATCH-20260409-ADMIN-USER-03 管理端用户模块回归测试与状态收口

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-USER-03
title: 管理端用户模块回归测试与状态收口
primary_req_id: ADMIN-USER-003
related_req_ids: ["ADMIN-USER-002"]
status: 已完成
owner_ai: Codex
dependencies: ["ST-BATCH-20260409-ADMIN-USER-01","ST-BATCH-20260409-ADMIN-USER-02"]
implementation_notes:
  - 使用 Playwright wrapper 完成列表交互、详情页 tab 切换和路由实测
  - 归档截图、URL、console、network 证据
  - 按链路回写 test cases / bugs / subtasks / overview / requirements 状态
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-USER-03-01"]
linked_bug_ids: ["BUG-BATCH-20260409-ADMIN-USER-001"]
last_updated_at: 2026-04-09 22:04
```
