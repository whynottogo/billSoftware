# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-CLOSED-LOOP-QA-01 | 注册启停登录闭环 Playwright 回归 | USER-AUTH-002 | ["USER-AUTH-001","USER-AUTH-003","SYS-AUTH-002","ADMIN-USER-001","ADMIN-USER-002"] | 已完成 | Codex | ["ADMIN-USER-001","ADMIN-USER-002"] | ["TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-01","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-02","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-03","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-04","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-05"] | ["BUG-BATCH-20260409-CLOSED-LOOP-QA-001","BUG-BATCH-20260409-CLOSED-LOOP-QA-002"] | 2026-04-10 00:18 |

## 子任务说明

### ST-BATCH-20260409-CLOSED-LOOP-QA-01 注册启停登录闭环 Playwright 回归

```yaml
subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
title: 注册启停登录闭环 Playwright 回归
primary_req_id: USER-AUTH-002
related_req_ids: ["USER-AUTH-001","USER-AUTH-003","SYS-AUTH-002","ADMIN-USER-001","ADMIN-USER-002"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-USER-001","ADMIN-USER-002"]
implementation_notes:
  - 已完成注册 -> 管理员启用/禁用 -> 用户登录/失效提示的真实闭环回归
  - 先后处理了 504（后端不可达）与 500（users 表不存在）两个阻塞
  - 最终验证结果：5 条固定用例全部通过
linked_test_case_ids: ["TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-01","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-02","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-03","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-04","TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-05"]
linked_bug_ids: ["BUG-BATCH-20260409-CLOSED-LOOP-QA-001","BUG-BATCH-20260409-CLOSED-LOOP-QA-002"]
last_updated_at: 2026-04-10 00:18
```
