# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 2 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-USER-CHART-01-01 支出图表页验证

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-CHART-01-01
related_subtask_id: ST-BATCH-20260409-USER-CHART-01
primary_req_id: USER-CHART-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:21
evidence:
  - assets/playwright/tc01-expense-assert.txt
  - assets/playwright/tc01-expense.png
  - assets/playwright/tc01-expense-after-switch-snapshot.yml
  - assets/playwright/tc01-expense-url.txt
  - assets/playwright/tc01-expense-console.log
  - assets/playwright/tc01-expense-network.log
```

Preconditions:

- 前端已构建完成并通过独立静态服务 `http://localhost:9101` 访问
- 已注入用户端登录态

Steps:

1. 访问 `/user/charts/expense`
2. 验证支出页包含月统计折线、年统计折线与前十排行榜
3. 点击年份切换到下一年并确认数据刷新

Expected:

- 支出页结构完整，包含双折线与前十排行
- 年份切换后标题和数据同步变化
- 控制台无 error

Actual:

- `tc01-expense-assert.txt` 显示 `yearChanged=true`、`rankingCount=10`
- 页面 URL 正确且控制台无 error

Evidence:

- `assets/playwright/tc01-expense-assert.txt`
- `assets/playwright/tc01-expense.png`
- `assets/playwright/tc01-expense-after-switch-snapshot.yml`
- `assets/playwright/tc01-expense-url.txt`
- `assets/playwright/tc01-expense-console.log`
- `assets/playwright/tc01-expense-network.log`

### TC-ST-BATCH-20260409-USER-CHART-02-01 收入图表页验证

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-CHART-02-01
related_subtask_id: ST-BATCH-20260409-USER-CHART-02
primary_req_id: USER-CHART-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:21
evidence:
  - assets/playwright/tc02-income-assert.txt
  - assets/playwright/tc02-income.png
  - assets/playwright/tc02-income-after-switch-snapshot.yml
  - assets/playwright/tc02-income-url.txt
  - assets/playwright/tc02-income-console.log
  - assets/playwright/tc02-income-network.log
```

Preconditions:

- 前端已构建完成并通过独立静态服务 `http://localhost:9101` 访问
- 已注入用户端登录态

Steps:

1. 从支出页切换到 `/user/charts/income`
2. 验证收入页包含年收入折线和前十排行
3. 点击年份切换并确认数据刷新

Expected:

- 收入页路由可达且结构完整
- 年份切换后标题和数据同步变化
- 控制台无 error

Actual:

- `tc02-income-assert.txt` 显示 `incomeRoute=/user/charts/income`、`yearChanged=true`、`rankingCount=10`
- 页面 URL 正确且控制台无 error

Evidence:

- `assets/playwright/tc02-income-assert.txt`
- `assets/playwright/tc02-income.png`
- `assets/playwright/tc02-income-after-switch-snapshot.yml`
- `assets/playwright/tc02-income-url.txt`
- `assets/playwright/tc02-income-console.log`
- `assets/playwright/tc02-income-network.log`
