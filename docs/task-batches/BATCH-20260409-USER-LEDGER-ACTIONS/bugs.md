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

## 当前状态

### BUG-BATCH-20260409-USER-LEDGER-ACTIONS-001 新增支出保存后弹层未关闭导致后续交互被拦截

```yaml
bug_id: BUG-BATCH-20260409-USER-LEDGER-ACTIONS-001
primary_req_id: USER-LEDGER-002
related_subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-01
related_test_case_id: TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-01-01
status: 已关闭
severity: medium
summary: 新增支出保存后弹层未关闭，遮罩层拦截后续点击，导致新增收入流程无法继续
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
owner_ai: Codex
last_updated_at: 2026-04-09 22:08
```

Steps:

1. 打开 `/user/ledger`，点击“记一笔支出”
2. 填写分类、金额、备注后点击“保存”
3. 再尝试点击“记一笔收入”

Expected:

- 点击“保存”后新增支出弹层关闭
- 页面可继续执行新增收入等后续交互

Actual:

- 初始回归时新增支出弹层保持打开，遮罩层拦截后续点击
- 修复后在独立静态回归环境复测通过

Evidence:

- `assets/playwright/tc01-expense-after-save.png`
- `assets/playwright/tc01-expense-after-save-snapshot.yml`
- `assets/playwright/tc02-income-click-blocked.log`
- `assets/playwright/tc01-final-after-save-snapshot.yml`
- `assets/playwright/tc02-final-after-save-snapshot.yml`

Suspected Root Cause:

- 金额输入控件在自动化输入场景下存在兼容性差异，导致保存流程未进入成功分支

Retest Result:

- 已按 `新建 -> 已确认 -> 修复中 -> 待回归 -> 已关闭` 收口
- 修复项：`UserHome.vue` 金额输入从 `el-input-number` 改为原生 `number` 输入
- 回归结论：新增支出与新增收入流程均可保存并即时回显

### BUG-BATCH-20260409-USER-LEDGER-ACTIONS-002 共享编译异常触发 overlay 拦截分类管理点击

```yaml
bug_id: BUG-BATCH-20260409-USER-LEDGER-ACTIONS-002
primary_req_id: USER-LEDGER-004
related_subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-03
related_test_case_id: TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-03-01
status: 已关闭
severity: medium
summary: dev-server 编译错误触发 overlay iframe，拦截“分类管理”按钮点击
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
owner_ai: Codex
last_updated_at: 2026-04-09 22:08
```

Steps:

1. 打开 `/user/ledger`
2. 执行“分类管理”按钮点击
3. 观察页面存在 `webpack-dev-server-client-overlay`，点击被拦截

Expected:

- 页面可正常响应“分类管理”点击并打开分类管理弹层

Actual:

- 初始阶段被 overlay iframe 拦截
- 修复共享编译缺口后，在独立静态回归环境流程可用

Evidence:

- `assets/playwright/tc03-overlay-blocked.log`
- `assets/playwright/tc03-overlay-blocked.png`
- `assets/playwright/tc03-overlay-blocked-snapshot.yml`
- `assets/playwright/tc03-final-after-add-snapshot.yml`
- `assets/playwright/tc03-final-console.log`

Suspected Root Cause:

- 共享路由 `router/index.js` 仍引用 `@/pages/admin/AdminUsers.vue`，但文件在并行修改中缺失

Retest Result:

- 已按 `新建 -> 已确认 -> 修复中 -> 待回归 -> 已关闭` 收口
- 修复项：新增 `frontend/src/pages/admin/AdminUsers.vue` 占位文件，恢复共享编译可用性
- 回归结论：分类管理流程可执行，收入/支出分类新增可见
