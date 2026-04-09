# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 5 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-BILL-BUDGET-01-01 月账单列表页展示与跳转

```yaml
test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-01-01
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-01
primary_req_id: USER-BILL-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:41
evidence:
  - assets/playwright/tc01-user-bills-month-open.txt
  - assets/playwright/tc01-user-bills-month-snapshot.txt
  - assets/playwright/tc01-user-bills-month.png
  - assets/playwright/tc01-user-bills-month-url.txt
  - assets/playwright/tc01-user-bills-month-console-error.txt
  - assets/playwright/tc01-user-bills-month-console.txt
  - assets/playwright/tc01-user-bills-month-network.txt
  - assets/playwright/tc01-user-bills-month-click.txt
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/bills/month`

Steps:

1. 打开 `/user/bills/month`
2. 检查年度汇总卡与月份列表
3. 点击月份进入月账单详情页

Expected:

- 页面展示本年结余、年收入、年支出
- 月份列表展示月收入、月支出、月结余
- 点击月份可进入对应月账单详情页

Actual:

- 页面成功展示年度汇总与月份列表
- 点击首个账单卡片后成功跳转至 `http://localhost:9000/user/bills/month/2026-04`
- 控制台 error 为 0，网络请求与路由跳转均正常

Evidence:

- `assets/playwright/tc01-user-bills-month-open.txt`
- `assets/playwright/tc01-user-bills-month-snapshot.txt`
- `assets/playwright/tc01-user-bills-month.png`
- `assets/playwright/tc01-user-bills-month-url.txt`
- `assets/playwright/tc01-user-bills-month-console-error.txt`
- `assets/playwright/tc01-user-bills-month-console.txt`
- `assets/playwright/tc01-user-bills-month-network.txt`
- `assets/playwright/tc01-user-bills-month-click.txt`

### TC-ST-BATCH-20260409-BILL-BUDGET-02-01 月账单详情页展示

```yaml
test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-02-01
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-02
primary_req_id: USER-BILL-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:42
evidence:
  - assets/playwright/tc02-user-bill-detail-open.txt
  - assets/playwright/tc02-user-bill-detail-snapshot.txt
  - assets/playwright/tc02-user-bill-detail.png
  - assets/playwright/tc02-user-bill-detail-url.txt
  - assets/playwright/tc02-user-bill-detail-console-error.txt
  - assets/playwright/tc02-user-bill-detail-console.txt
  - assets/playwright/tc02-user-bill-detail-network.txt
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/bills/month/2026-04`

Steps:

1. 打开 `/user/bills/month/2026-04`
2. 检查月汇总、分类构成、排行、趋势和成就区块
3. 验证页面切换月份时内容能同步变化

Expected:

- 页面显示本月结余、上月结余、本月收入、本月支出
- 页面展示支出分类、排行、趋势和对比内容
- 页面整体结构与账单列表页保持同一设计体系

Actual:

- 页面完整展示月汇总、分类构成、排行、趋势和成就区块
- 控制台 error 为 0，未复现历史 ECharts TypeError
- 页面结构与账单列表视觉体系保持一致，满足回归通过标准

Evidence:

- `assets/playwright/tc02-user-bill-detail-open.txt`
- `assets/playwright/tc02-user-bill-detail-snapshot.txt`
- `assets/playwright/tc02-user-bill-detail.png`
- `assets/playwright/tc02-user-bill-detail-url.txt`
- `assets/playwright/tc02-user-bill-detail-console-error.txt`
- `assets/playwright/tc02-user-bill-detail-console.txt`
- `assets/playwright/tc02-user-bill-detail-network.txt`

### TC-ST-BATCH-20260409-BILL-BUDGET-03-01 年账单页展示

```yaml
test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-03-01
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-03
primary_req_id: USER-BILL-003
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:44
evidence:
  - assets/playwright/tc03-user-bills-year-open.txt
  - assets/playwright/tc03-user-bills-year-switch.txt
  - assets/playwright/tc03-user-bills-year-snapshot.txt
  - assets/playwright/tc03-user-bills-year.png
  - assets/playwright/tc03-user-bills-year-url.txt
  - assets/playwright/tc03-user-bills-year-console-error.txt
  - assets/playwright/tc03-user-bills-year-console.txt
  - assets/playwright/tc03-user-bills-year-network.txt
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/bills/year`

Steps:

1. 打开 `/user/bills/year`
2. 检查年度总览、年份切换和历年列表
3. 验证年维度摘要和趋势内容完整

Expected:

- 页面显示总收入、总支出、总结余
- 页面可切换查看不同年份数据
- 历年列表展示每年收入、支出和结余

Actual:

- 年账单页可正常展示总览、趋势和历年列表
- 年份切换执行成功（见 switch 证据）
- 控制台 error 为 0，页面无脚本异常

Evidence:

- `assets/playwright/tc03-user-bills-year-open.txt`
- `assets/playwright/tc03-user-bills-year-switch.txt`
- `assets/playwright/tc03-user-bills-year-snapshot.txt`
- `assets/playwright/tc03-user-bills-year.png`
- `assets/playwright/tc03-user-bills-year-url.txt`
- `assets/playwright/tc03-user-bills-year-console-error.txt`
- `assets/playwright/tc03-user-bills-year-console.txt`
- `assets/playwright/tc03-user-bills-year-network.txt`

### TC-ST-BATCH-20260409-BILL-BUDGET-04-01 月预算页展示

```yaml
test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-04-01
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-04
primary_req_id: USER-BUDGET-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:45
evidence:
  - assets/playwright/tc04-user-budget-month-open.txt
  - assets/playwright/tc04-user-budget-month-snapshot.txt
  - assets/playwright/tc04-user-budget-month.png
  - assets/playwright/tc04-user-budget-month-url.txt
  - assets/playwright/tc04-user-budget-month-console-error.txt
  - assets/playwright/tc04-user-budget-month-console.txt
  - assets/playwright/tc04-user-budget-month-network.txt
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/budget/month`

Steps:

1. 打开 `/user/budget/month`
2. 检查本月预算、本月支出、本月剩余预算
3. 检查分类预算卡片和设置入口

Expected:

- 页面展示当月预算总览和剩余额度
- 页面支持按支出类型展示细分预算状态
- 页面保留清晰的预算设置入口

Actual:

- 页面成功展示本月预算、本月支出、本月剩余预算与分类预算卡片
- 设置入口与分类调整入口可交互（前端占位提示）
- 控制台 error 为 0

Evidence:

- `assets/playwright/tc04-user-budget-month-open.txt`
- `assets/playwright/tc04-user-budget-month-snapshot.txt`
- `assets/playwright/tc04-user-budget-month.png`
- `assets/playwright/tc04-user-budget-month-url.txt`
- `assets/playwright/tc04-user-budget-month-console-error.txt`
- `assets/playwright/tc04-user-budget-month-console.txt`
- `assets/playwright/tc04-user-budget-month-network.txt`

### TC-ST-BATCH-20260409-BILL-BUDGET-05-01 年预算页展示

```yaml
test_case_id: TC-ST-BATCH-20260409-BILL-BUDGET-05-01
related_subtask_id: ST-BATCH-20260409-BILL-BUDGET-05
primary_req_id: USER-BUDGET-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:45
evidence:
  - assets/playwright/tc05-user-budget-year-open.txt
  - assets/playwright/tc05-user-budget-year-switch.txt
  - assets/playwright/tc05-user-budget-year-snapshot.txt
  - assets/playwright/tc05-user-budget-year.png
  - assets/playwright/tc05-user-budget-year-url.txt
  - assets/playwright/tc05-user-budget-year-console-error.txt
  - assets/playwright/tc05-user-budget-year-console.txt
  - assets/playwright/tc05-user-budget-year-network.txt
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/budget/year`

Steps:

1. 打开 `/user/budget/year`
2. 检查年份切换、年度预算总览、月度执行和分类年度表现
3. 验证页面与月预算页的信息结构有明显区分

Expected:

- 页面展示本年预算、本年支出、本年剩余预算
- 页面支持切换查看不同年份预算数据
- 页面不是月预算页的简单复制

Actual:

- 页面成功展示本年预算、本年支出、本年剩余预算及月度执行图
- 年份切换执行成功（见 switch 证据）
- 控制台 error 为 0，页面结构与月预算页保持差异化

Evidence:

- `assets/playwright/tc05-user-budget-year-open.txt`
- `assets/playwright/tc05-user-budget-year-switch.txt`
- `assets/playwright/tc05-user-budget-year-snapshot.txt`
- `assets/playwright/tc05-user-budget-year.png`
- `assets/playwright/tc05-user-budget-year-url.txt`
- `assets/playwright/tc05-user-budget-year-console-error.txt`
- `assets/playwright/tc05-user-budget-year-console.txt`
- `assets/playwright/tc05-user-budget-year-network.txt`
