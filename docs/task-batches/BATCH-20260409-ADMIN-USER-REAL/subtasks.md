# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-ADMIN-USER-REAL-01 | 管理端用户列表真实接口接线 | ADMIN-USER-001 | ["ADMIN-USER-002"] | 已完成 | Codex | ["ADMIN-AUTH-001"] | ["TC-ST-BATCH-20260409-ADMIN-USER-REAL-01-01"] | ["BUG-BATCH-20260409-ADMIN-USER-REAL-001"] | 2026-04-10 00:18 |
| ST-BATCH-20260409-ADMIN-USER-REAL-02 | 管理端用户启停真实交互联调 | ADMIN-USER-002 | ["ADMIN-USER-001"] | 已完成 | Codex | ["ADMIN-USER-001"] | ["TC-ST-BATCH-20260409-ADMIN-USER-REAL-02-01"] | ["BUG-BATCH-20260409-ADMIN-USER-REAL-001"] | 2026-04-10 00:18 |

## 子任务说明

### ST-BATCH-20260409-ADMIN-USER-REAL-01 管理端用户列表真实接口接线

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-USER-REAL-01
title: 管理端用户列表真实接口接线
primary_req_id: ADMIN-USER-001
related_req_ids: ["ADMIN-USER-002"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-AUTH-001"]
implementation_notes:
  - 已复核 `GET /api/admin/users` 调用路径为 `/admin/users`（经 request baseURL 组合后为 `/api/admin/users`）
  - 已完成 id/status/created_at 的前端映射，详情按钮保持“详情联调中”
  - 2026-04-10 00:07 管理员真实登录后已验证新注册用户显示在真实列表中，状态默认为禁用
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-USER-REAL-01-01"]
linked_bug_ids: ["BUG-BATCH-20260409-ADMIN-USER-REAL-001"]
last_updated_at: 2026-04-10 00:18
```

### ST-BATCH-20260409-ADMIN-USER-REAL-02 管理端用户启停真实交互联调

```yaml
subtask_id: ST-BATCH-20260409-ADMIN-USER-REAL-02
title: 管理端用户启停真实交互联调
primary_req_id: ADMIN-USER-002
related_req_ids: ["ADMIN-USER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-USER-001"]
implementation_notes:
  - 已复核 `PUT /api/admin/users/:id/status` 请求体为 `{ status: 0|1 }`
  - 已保留启停按钮 pending 态和顶部统计联动逻辑
  - 2026-04-10 00:08 已验证启用后用户可登录，禁用后旧会话返回 401，重新登录返回 `user is disabled`
linked_test_case_ids: ["TC-ST-BATCH-20260409-ADMIN-USER-REAL-02-01"]
linked_bug_ids: ["BUG-BATCH-20260409-ADMIN-USER-REAL-001"]
last_updated_at: 2026-04-10 00:18
```
