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

### BUG-BATCH-20260409-ADMIN-PORTAL-001 并行改动窗口导致首次构建缺失 UserHome 依赖

```yaml
bug_id: BUG-BATCH-20260409-ADMIN-PORTAL-001
primary_req_id: ADMIN-DASHBOARD-001
related_subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-01
related_test_case_id: TC-ST-BATCH-20260409-ADMIN-PORTAL-01-01
status: 已关闭
severity: 中
summary: 首次执行 npm run build 时，router 依赖的 UserHome 页面在并行窗口内未就绪导致模块解析失败
environment: frontend build + /Users/shynin/software/billSoftware/frontend
steps:
  - 在 frontend 执行 npm run build
  - 观察 webpack 输出
expected:
  - 前端工程应可稳定构建通过
actual:
  - 首次构建报错 Cannot resolve @/pages/user/UserHome.vue
  - 后续复跑构建通过
evidence:
  - assets/playwright/admin-portal-build-error.log
suspected_root_cause: 并行 worker 改动共享路由与页面文件时序重叠，首次构建命中短暂不一致状态
owner_ai: Codex
last_updated_at: 2026-04-09 21:46
```

Steps:

1. 在 `frontend` 目录执行 `npm run build`
2. 观察 webpack 输出中的模块解析结果

Expected:

- 构建输出成功，无模块缺失

Actual:

- 首次构建出现 `UserHome.vue` 解析失败
- 复跑后已恢复并通过

Evidence:

- `assets/playwright/admin-portal-build-error.log`

### BUG-BATCH-20260409-ADMIN-PORTAL-002 Playwright 会话并发调用导致工具报错

```yaml
bug_id: BUG-BATCH-20260409-ADMIN-PORTAL-002
primary_req_id: ADMIN-APPROVAL-001
related_subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-02
related_test_case_id: TC-ST-BATCH-20260409-ADMIN-PORTAL-02-01
status: 已关闭
severity: 低
summary: 对同一 Playwright session 并发执行命令以及无 ref 直接文本点击导致工具报错
environment: Playwright wrapper + session workerf
steps:
  - 对同一 session 并行执行 snapshot/screenshot/eval/console/network
  - 或直接 click 文本目标 批准
expected:
  - 测试命令应被串行执行并准确命中元素
actual:
  - 返回 browser is not open 与 does not match any elements
  - 调整为串行采集后页面验证通过
evidence:
  - assets/playwright/admin-portal-playwright-session-error.log
suspected_root_cause: 测试执行方式误用，不属于页面运行时缺陷
owner_ai: Codex
last_updated_at: 2026-04-09 21:46
```

Steps:

1. 在同一会话并发执行多条 wrapper 命令
2. 无 snapshot ref 直接执行文本点击

Expected:

- 会话稳定可用，元素命中准确

Actual:

- 工具返回会话未打开和元素未命中
- 改为串行后已恢复

Evidence:

- `assets/playwright/admin-portal-playwright-session-error.log`

## 当前状态

- 两条异常均已归档并关闭，不影响页面最终验收。
