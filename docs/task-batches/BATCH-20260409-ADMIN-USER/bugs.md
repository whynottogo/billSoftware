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

### BUG-BATCH-20260409-ADMIN-USER-001 管理端用户列表页未命中预期 DOM，导致启停回归中断

```yaml
bug_id: BUG-BATCH-20260409-ADMIN-USER-001
primary_req_id: ADMIN-USER-002
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-01
related_test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-01-01
status: 已关闭
severity: 高
summary: Playwright 在 /admin/users 回归时无法定位 .status-pill，启用/禁用用例中断
environment: frontend dev server + Playwright wrapper + http://localhost:9000/admin/users
steps:
  - 打开 /admin/users
  - 执行状态切换回归脚本
expected:
  - 页面包含用户状态标签并可执行启停交互
actual:
  - 初次回归时 Playwright 定位 .status-pill 超时
  - 修复后 /admin/users 渲染用户列表正常，`tc01-admin-users-toggle.txt` 已验证状态切换成功
evidence:
  - assets/playwright/tc01-admin-users-open.txt
  - assets/playwright/tc01-admin-users-toggle.txt
suspected_root_cause: 并行开发导致 AdminUsers 页面被占位内容覆盖，未加载用户列表实现
owner_ai: Codex
last_updated_at: 2026-04-09 22:04
```

## 当前状态

- 当前批次缺陷已全部关闭。
- 若出现控制台 error、路由异常、4xx/5xx、空白页，先登记到本文件再排查。
