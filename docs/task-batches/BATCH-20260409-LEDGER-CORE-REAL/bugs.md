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

## 说明

- 本批次在接口联调或 Playwright 回归中发现异常时，必须先归档 `assets/playwright/` 证据，再在此登记。

### BUG-BATCH-20260409-LEDGER-CORE-REAL-001 后端联调环境阻塞

```yaml
bug_id: BUG-BATCH-20260409-LEDGER-CORE-REAL-001
primary_req_id: SYS-DATA-001
related_subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-01
related_test_case_id: TC-ST-BATCH-20260409-LEDGER-CORE-REAL-01-01
status: 已关闭
severity: S2
summary: 执行机最初缺少 Go 运行环境且数据库 schema 未初始化，导致用户侧 ledger/category 真实联调无法执行
environment: local (macOS, 2026-04-10 00:18 CST)
steps:
  - 启动本地后端并访问 `http://localhost:8080/api/health`
  - 登录用户并访问 `/user/ledger`
expected:
  - 后端可启动，ledger/category 真实接口可联调
actual:
  - 初始阶段：`go` 命令不可用、`/api/health` 连接拒绝
  - 修复后：已安装 Go 1.26.2、生成 `backend/go.sum`、执行 schema 与 seed SQL，随后后端恢复可用且全部测试通过
evidence:
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260409-2342-backend-health-curl.txt
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260409-2354-docker-pull-golang.txt
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-ledger-backend-evidence.txt
suspected_root_cause: 当前执行机最初缺少 Go 运行环境，且远端数据库未初始化 schema
owner_ai: Codex
last_updated_at: 2026-04-10 00:18
```
