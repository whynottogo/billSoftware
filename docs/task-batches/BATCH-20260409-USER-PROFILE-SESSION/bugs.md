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

### BUG-BATCH-20260409-USER-PROFILE-SESSION-001 首轮构建失败导致测试阻塞

```yaml
bug_id: BUG-BATCH-20260409-USER-PROFILE-SESSION-001
primary_req_id: USER-PROFILE-001
related_subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-01
related_test_case_id: TC-ST-BATCH-20260409-USER-PROFILE-SESSION-01-01
status: 已关闭
severity: 高
summary: Worker C 首轮构建时出现路由并行冲突与 mock 语法兼容错误，导致页面验证无法继续
environment: frontend build + webpack + babel + http://localhost:9000
steps:
  - 在 frontend 目录执行 npm run build
  - 观察 webpack 与 babel 输出
expected:
  - 工程构建通过并可继续执行 Playwright 页面验证
actual:
  - 首轮构建报错无法解析 @/pages/user/UserHome.vue 且提示 mock 语法兼容问题
  - 修复后构建通过，Playwright 回归通过
evidence:
  - assets/playwright/profile-session-build-error.log
  - assets/playwright/profile-session-build-success.log
  - assets/playwright/tc01-user-profile-console.log
suspected_root_cause: 并行开发时共享路由引用瞬时不一致，且 profile mock 在首版中触发了当前编译链兼容问题
owner_ai: Codex
last_updated_at: 2026-04-09 21:50
```

Steps:

1. 在 `frontend` 目录执行 `npm run build`
2. 查看终端输出中的 webpack 与 babel 报错

Expected:

- 构建通过并可继续页面验证

Actual:

- 首轮构建失败并阻塞测试
- 修复后构建通过，且 4 个 Playwright 用例均通过

Evidence:

- `assets/playwright/profile-session-build-error.log`
- `assets/playwright/profile-session-build-success.log`
- `assets/playwright/tc01-user-profile-console.log`

Retest Result:

- 状态流转：`新建 -> 已确认 -> 修复中 -> 待回归 -> 已关闭`
- 回归结论：`npm run build` 成功，`TC-01/02/03/04` 全部通过，缺陷关闭。
