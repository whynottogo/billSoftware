# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 2 |
| 已挂起 | 0 |

## 说明

- 本批次专门维护闭环回归异常，发现问题后必须先归档到 `assets/playwright/` 再更新此文件。

### BUG-BATCH-20260409-CLOSED-LOOP-QA-001 注册接口 504 导致闭环阻塞

```yaml
bug_id: BUG-BATCH-20260409-CLOSED-LOOP-QA-001
primary_req_id: USER-AUTH-001
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
related_test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-01
status: 已关闭
severity: S2
summary: 用户注册提交时 /api/user/auth/register 返回 504，导致闭环在第 1 条用例阻塞
environment: local (macOS, 2026-04-10 00:18 CST)
steps:
  - 打开 http://localhost:9000/user/register
  - 填写注册信息并提交
expected:
  - 注册成功并返回等待管理员启用提示
actual:
  - 初始阶段：前端代理 504，后端健康检查连接拒绝
  - 修复后：本地后端恢复，注册用例通过
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260409-234831-register-submit-network.txt
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260409-234831-backend-health.txt
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
suspected_root_cause: 当时本地后端未启动
owner_ai: Codex
last_updated_at: 2026-04-10 00:18
```

### BUG-BATCH-20260409-CLOSED-LOOP-QA-002 注册接口 500：check username failed

```yaml
bug_id: BUG-BATCH-20260409-CLOSED-LOOP-QA-002
primary_req_id: USER-AUTH-001
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
related_test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-01
status: 已关闭
severity: S2
summary: 本地后端恢复后，用户注册提交返回 500，页面提示 `check username failed`
environment: local (macOS, 2026-04-10 00:18 CST)
steps:
  - 打开 http://localhost:9000/user/register
  - 填写注册信息并提交
expected:
  - 注册成功并返回等待管理员启用提示
actual:
  - 初始阶段：注册返回 500，页面提示 `check username failed`
  - 定位结果：远端数据库未初始化 `users` 等核心表
  - 修复后：执行 schema 与 seed SQL，注册与后续闭环全部通过
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0003-register-500-page.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0003-register-500-console.log
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
suspected_root_cause: 远端数据库缺少 schema 初始化
owner_ai: Codex
last_updated_at: 2026-04-10 00:18
```
