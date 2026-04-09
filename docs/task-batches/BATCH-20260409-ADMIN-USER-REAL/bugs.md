# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 1 |
| 已挂起 | 0 |

## 缺陷记录

### BUG-BATCH-20260409-ADMIN-USER-REAL-001 后端服务不可达导致真实联调阻塞

```yaml
bug_id: BUG-BATCH-20260409-ADMIN-USER-REAL-001
primary_req_id: ADMIN-USER-002
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-REAL-02
related_test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-REAL-02-01
status: 已关闭
severity: S2
summary: 本地环境后端 http://localhost:8080 不可达，导致管理端用户真实接口联调无法执行
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
steps:
  - 启动前端并访问管理端登录页
  - 执行后端健康检查请求
  - 观察接口不可达返回
expected:
  - /api/health 返回 200，管理端可继续联调 /api/admin/users 和状态修改接口
actual:
  - 2026-04-09 23:46 首次检查时 `/api/health` 不可达
  - 2026-04-10 00:00 已补齐本地 Go 环境、生成 `go.sum` 并初始化数据库 schema，随后 `/api/health` 恢复 200
  - 回归已通过，Bug 关闭
evidence:
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/2026-04-09-admin-user-real-validation-blocker.md
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/20260410-0010-admin-auth-evidence.txt
suspected_root_cause: 当前执行机最初缺少 Go 运行环境，且数据库 schema 未初始化
owner_ai: Codex
last_updated_at: 2026-04-10 00:18
```

## 说明

- 本批次在真实联调或 Playwright 回归中如发现异常，先归档到 `assets/playwright/`，再在此处登记。
