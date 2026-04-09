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

### BUG-BATCH-20260409-BILL-BUDGET-001 月账单详情页图表切换时控制台报错

```yaml
bug_id: BUG-BATCH-20260409-BILL-BUDGET-001
primary_req_id: USER-BILL-002
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-02
related_test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-02-01
status: 已关闭
severity: 高
summary: 从月账单列表进入月账单详情页时，ECharts 在 SPA 路由切换过程中抛出 TypeError
environment: frontend dev server + Playwright wrapper + http://localhost:9000/user/bills/month/2026-04
steps:
  - 打开 /user/bills/month
  - 点击 2026年4月 月账单卡片进入详情页
  - 观察浏览器控制台输出
expected:
  - 月账单详情页图表无报错渲染
  - 控制台不出现 ECharts 相关异常
actual:
  - 旧问题中控制台曾记录 4 条 TypeError
  - 完成共享图表替换后回归通过，`tc02-user-bill-detail-console-error.txt` 显示 error 为 0
evidence:
  - assets/playwright/user-bill-detail-open.txt
  - assets/playwright/user-bill-detail-snapshot.txt
  - assets/playwright/user-bill-detail.png
  - assets/playwright/user-bill-detail-console-error.log
  - assets/playwright/user-bill-detail-network.txt
suspected_root_cause: FinanceChart 组件对 ECharts option 使用深度 watch，路由切换时被内部对象变更再次触发 setOption，导致图表状态错乱
owner_ai: Codex
last_updated_at: 2026-04-09 21:47
```

Steps:

1. 打开 `http://localhost:9000/user/bills/month`
2. 点击 2026 年 4 月账单卡片
3. 进入详情页后查看控制台

Expected:

- 图表在路由切换后正常渲染
- 控制台无 ECharts 报错

Actual:

- `assets/playwright/user-bill-detail-console-error.log` 中记录 4 条 ECharts `TypeError`
- 详情页页面可见，但图表渲染稳定性不满足通过标准

Evidence:

- `assets/playwright/user-bill-detail-open.txt`
- `assets/playwright/user-bill-detail-snapshot.txt`
- `assets/playwright/user-bill-detail.png`
- `assets/playwright/user-bill-detail-console-error.log`
- `assets/playwright/user-bill-detail-network.txt`

## 当前状态

- 当前批次缺陷已全部关闭。
- 回归完成后，月账单详情、年账单、预算月/年页面均未出现控制台 error。

### BUG-BATCH-20260409-BILL-BUDGET-002 回归过程中编译失败导致年账单页面不可用

```yaml
bug_id: BUG-BATCH-20260409-BILL-BUDGET-002
primary_req_id: USER-BILL-003
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-03
related_test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-03-01
status: 已关闭
severity: 高
summary: Playwright 回归 /user/bills/year 时出现共享模块编译错误，导致年账单页面不可用
environment: frontend dev server + Playwright wrapper + http://localhost:9000/user/bills/year
steps:
  - 打开 /user/bills/year
  - 查看控制台 error 输出
expected:
  - 项目编译成功，年账单页可正常渲染
  - 控制台不出现模块缺失错误
actual:
  - 阻塞期间出现模块缺失与语法不兼容，触发前端编译失败
  - 修复后 `frontend-build-regression.txt` 编译通过，`tc03-user-bills-year-console-error.txt` 显示 error 为 0
evidence:
  - assets/playwright/user-bills-year-console-error-regression.txt
  - assets/playwright/user-bills-year-open-regression.txt
  - assets/playwright/user-bills-year-regression.png
  - assets/playwright/frontend-build-regression.txt
  - assets/playwright/tc03-user-bills-year-console-error.txt
suspected_root_cause: 并行开发修改共享路由依赖时，出现页面文件缺失与 ES 语法兼容问题
owner_ai: Codex
last_updated_at: 2026-04-09 21:47
```
