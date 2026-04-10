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

### BUG-BATCH-20260410-BUDGET-CHART-REAL-001 前端构建无法解析 UserBudgetMonth 页面

```yaml
bug_id: BUG-BATCH-20260410-BUDGET-CHART-REAL-001
primary_req_id: USER-BUDGET-001
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-02
related_test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-01
status: 已关闭
severity: high
summary: 前端基线构建时无法解析 `@/pages/user/UserBudgetMonth.vue`，预算与图表联调被阻塞
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
steps:
  - 在 `frontend` 目录执行 `npm run build`
  - 等待 webpack 编译路由与页面模块
expected:
  - 前端构建通过，预算页面模块可正常解析并进入后续联调
actual:
  - webpack 报错 `Module not found: Can't resolve '@/pages/user/UserBudgetMonth.vue'`，构建终止
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-precheck.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-after-budget-fix.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-final.txt
suspected_root_cause: 预算月页面文件缺失、路径大小写不匹配，或并发实现期间页面文件被暂时移动
owner_ai: Codex
last_updated_at: 2026-04-10 12:54
```

Steps:

1. 在 `frontend` 目录执行 `npm run build`
2. 等待 webpack 编译路由与页面模块

Expected:

- 前端构建通过，预算页面模块可正常解析并进入后续联调

Actual:

- webpack 报错 `Module not found: Can't resolve '@/pages/user/UserBudgetMonth.vue'`，构建终止

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-precheck.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-after-budget-fix.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-final.txt`
- 当前 URL：`N/A`
- 控制台摘要：`npm run build` 输出 `ERROR in ./src/router/index.js 11:0-63`
- 关键请求：`N/A`

Suspected Root Cause:

- 预算月页面文件缺失、路径大小写不匹配，或并发实现期间页面文件被暂时移动

Retest Result:

- 2026-04-10 12:26 已补回预算页文件并重跑 `npm run build`，构建通过，阻塞解除。
- 2026-04-10 12:54 再次执行最终构建验证，前端构建继续通过。

### BUG-BATCH-20260410-BUDGET-CHART-REAL-002 Playwright 浏览器工具初始化目录失败

```yaml
bug_id: BUG-BATCH-20260410-BUDGET-CHART-REAL-002
primary_req_id: USER-BUDGET-001
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-05
related_test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-02
status: 已关闭
severity: medium
summary: Playwright 浏览器工具初始化时尝试创建 `/.playwright-mcp` 目录失败，已切换 wrapper 继续执行真实页面回归
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
steps:
  - 调用浏览器工具打开 `http://localhost:9000/user/login`
  - 观察浏览器工具初始化输出
expected:
  - Playwright 浏览器工具正常打开页面并允许继续交互
actual:
  - 工具返回 `Error: ENOENT: no such file or directory, mkdir '/.playwright-mcp'`
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-playwright-browser-init-error.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income-network.txt
suspected_root_cause: 工具默认尝试在只读根目录创建 `/.playwright-mcp`，当前桌面环境不允许写入该路径
owner_ai: Codex
last_updated_at: 2026-04-10 12:54
```

Steps:

1. 调用浏览器工具打开 `http://localhost:9000/user/login`
2. 观察浏览器工具初始化输出

Expected:

- Playwright 浏览器工具正常打开页面并允许继续交互

Actual:

- 工具返回 `Error: ENOENT: no such file or directory, mkdir '/.playwright-mcp'`

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-playwright-browser-init-error.txt`
- 当前 URL：`http://localhost:9000/user/login`
- 控制台摘要：`Playwright browser tool init failed before page load`
- 关键请求：`N/A`

Suspected Root Cause:

- 工具默认尝试在只读根目录创建 `/.playwright-mcp`，当前桌面环境不允许写入该路径

Retest Result:

- 已切换本机 Playwright wrapper 继续执行真实浏览器回归，不再阻塞本批次测试。
- 2026-04-10 12:54 月预算、年预算、支出图表、收入图表真实页面回归全部通过，验收完成。

## 说明

- 本批次发现预算保存失败、图表空白、控制台报错、关键网络请求失败或接口结构不符时，必须先归档证据到 `assets/playwright/`，再在此登记。
- 当前结论：`BUG-BATCH-20260410-BUDGET-CHART-REAL-001`、`BUG-BATCH-20260410-BUDGET-CHART-REAL-002` 均已关闭，本批次无未关闭 Bug。
