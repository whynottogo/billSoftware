# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 8 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-01 预算接口结构验证

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-01
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-01
primary_req_id: USER-BUDGET-001
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 请求当前月预算读取接口
  - 请求当前月预算保存接口
  - 请求年度预算年份列表与年份详情接口
  - 请求当前年预算保存接口
expected:
  - 接口返回结构满足预算页面消费字段
  - 月预算仅允许当前月份保存，年预算仅允许当前年份保存
  - 返回的分类预算仅包含当前用户的支出分类
actual:
  - GET/PUT 当前月预算接口返回和保存均成功，保存后回读为总预算 1800
  - 年预算年份列表、详情和当前年保存接口均成功，保存后回读为总预算 20000
  - 分类预算仅返回当前用户支出分类，字段满足页面消费结构
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-month-current.json
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-month-after-save.json
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-year-options.json
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-year-2026.json
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-year-2026-after-save.json
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 后端已启动
- 已获取用户 token

Steps:

1. 请求当前月预算读取接口
2. 请求当前月预算保存接口
3. 请求年度预算年份列表与年份详情接口
4. 请求当前年预算保存接口

Expected:

- 接口返回结构满足预算页面消费字段
- 月预算仅允许当前月份保存，年预算仅允许当前年份保存
- 返回的分类预算仅包含当前用户的支出分类

Actual:

- `GET /api/user/budgets/month/current` 返回当前月预算总览、分类预算、支出与剩余金额。
- `PUT /api/user/budgets/month/current` 保存后再次读取，月预算更新为总预算 `1800`，分类预算包含购物 `360`、交通 `240` 等当前用户支出分类。
- `GET /api/user/budgets/year/options` 返回可选年份，`GET /api/user/budgets/year/2026` 返回年度预算、月度执行图与分类预算。
- `PUT /api/user/budgets/year/current` 保存后再次读取，年度预算更新为总预算 `20000`，接口契约满足页面消费结构。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-month-current.json`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-month-after-save.json`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-year-options.json`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-year-2026.json`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-budget-year-2026-after-save.json`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-02 图表接口结构验证

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-01-02
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-01
primary_req_id: USER-CHART-001
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 请求图表年份列表接口
  - 请求支出图表接口
  - 请求收入图表接口
expected:
  - 支出接口返回 summary / monthTrend / yearTrend / ranking
  - 收入接口返回 summary / yearTrend / ranking
  - 排行榜最多返回 10 条，且仅包含当前用户数据
actual:
  - 图表年份接口返回 2026
  - 支出图表返回年总支出 1450，收入图表返回年总收入 11200
  - 两个排行榜都在 10 条以内且按当前用户隔离
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-charts-years.json
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-charts-expense-2026.json
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-charts-income-2026.json
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 后端已启动
- 已获取用户 token

Steps:

1. 请求图表年份列表接口
2. 请求支出图表接口
3. 请求收入图表接口

Expected:

- 支出接口返回 `summary / monthTrend / yearTrend / ranking`
- 收入接口返回 `summary / yearTrend / ranking`
- 排行榜最多返回 10 条，且仅包含当前用户数据

Actual:

- `GET /api/user/charts/years` 正常返回年份列表，当前测试数据覆盖 `2026`。
- `GET /api/user/charts/expense/2026` 返回 `summary / monthTrend / yearTrend / ranking`，年总支出为 `1450`，排行榜数量在 10 条以内。
- `GET /api/user/charts/income/2026` 返回 `summary / yearTrend / ranking`，年总收入为 `11200`，排行榜数量在 10 条以内。
- 返回数据均按当前登录用户隔离，满足页面消费结构。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-charts-years.json`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-charts-expense-2026.json`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-api-charts-income-2026.json`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-01 用户月预算真实设置与回显

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-01
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-02
primary_req_id: USER-BUDGET-001
status: 已通过
preconditions:
  - 用户账号可登录
  - 预算接口已联通
steps:
  - 登录用户端并进入 /user/budget/month
  - 打开预算设置弹窗并填写总预算及分类预算
  - 保存后等待页面回显
expected:
  - 页面显示真实本月预算、本月支出和剩余预算
  - 分类预算卡显示真实预算占用和状态
  - 保存后网络请求成功且页面立即刷新为新值
actual:
  - 页面回显总预算 1800、本月支出 350、剩余预算 1450
  - 分类预算卡展示购物 360、交通 240 等真实预算值
  - 网络请求成功且控制台无未预期 error
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-dialog.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-after-save.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month.png
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 用户账号可登录
- 预算接口已联通

Steps:

1. 登录用户端并进入 `/user/budget/month`
2. 打开预算设置弹窗并填写总预算及分类预算
3. 保存后刷新页面或等待页面自动回显

Expected:

- 页面显示真实本月预算、本月支出和剩余预算
- 分类预算卡显示真实预算占用和状态
- 保存后网络请求成功且页面立即刷新为新值

Actual:

- 登录后进入月预算页，真实接口成功回填本月预算数据。
- 通过弹窗保存后，页面即时回显总预算 `1800`、本月支出 `350`、剩余预算 `1450`。
- 分类预算卡同步显示购物 `360`、交通 `240` 等预算值，关键请求成功且无未预期控制台 error。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-dialog.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-after-save.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month.png`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-02 用户年预算真实设置与回显

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-02-02
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-02
primary_req_id: USER-BUDGET-002
status: 已通过
preconditions:
  - 用户账号可登录
  - 年预算接口已联通
steps:
  - 进入 /user/budget/year
  - 调整当前年预算与分类预算
  - 验证年度总览和月度执行图更新
expected:
  - 页面展示真实年度预算、支出、剩余预算和执行进度
  - 月度预算执行图和分类预算卡与保存结果一致
  - 只能设置当前年份预算
actual:
  - 页面回显总预算 20000、本年支出 1450、剩余预算 18550
  - 住房分类预算更新为 1200，月度执行图同步更新
  - 当前年份预算保存成功，控制台无未预期 error
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-dialog.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-after-save.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year.png
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 用户账号可登录
- 年预算接口已联通

Steps:

1. 进入 `/user/budget/year`
2. 调整当前年预算与分类预算
3. 验证年度总览和月度执行图更新

Expected:

- 页面展示真实年度预算、支出、剩余预算和执行进度
- 月度预算执行图和分类预算卡与保存结果一致
- 只能设置当前年份预算

Actual:

- 年预算页真实加载 `2026` 年数据并展示年度总览、月度执行图和分类预算卡。
- 保存后页面即时回显总预算 `20000`、本年支出 `1450`、剩余预算 `18550`，住房分类预算更新为 `1200`。
- 当前年预算可正常保存，页面与接口口径一致，无路由或控制台异常。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-dialog.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-after-save.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year.png`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-03-01 用户支出图表真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-03-01
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-03
primary_req_id: USER-CHART-001
status: 已通过
preconditions:
  - 用户账号可登录
  - 支出图表接口已联通
steps:
  - 打开 /user/charts/expense
  - 验证摘要卡
  - 校验趋势图和排行榜
expected:
  - 页面全部使用真实接口数据
  - 月趋势为按天折线，年趋势为按月折线
  - 排行榜展示支出前 10 类型及占比
actual:
  - 页面展示年总支出 1450、月均支出 120.83、最高分类住房 35.86%
  - 月趋势、年趋势和排行榜均正常渲染
  - 网络请求成功且控制台无未预期 error
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-expense.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-expense.png
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 用户账号可登录
- 支出图表接口已联通

Steps:

1. 打开 `/user/charts/expense`
2. 验证年总支出、月均支出和高频分类摘要卡
3. 校验按天月趋势、按月年趋势与支出排行榜
4. 切换年份再次验证

Expected:

- 页面全部使用真实接口数据
- 月趋势为按天折线，年趋势为按月折线
- 排行榜展示支出前 10 类型及占比

Actual:

- 支出图表页成功渲染真实支出数据，摘要卡显示年总支出 `1450`、月均支出 `120.83`、高频分类 `住房` 占比 `35.86%`。
- 月趋势与年趋势图均正常渲染，排行榜按真实支出分类聚合且数量不超过 10 条。
- 关键请求成功，未发现未预期控制台 error。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-expense.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-expense.png`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-04-01 用户收入图表真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-04-01
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-04
primary_req_id: USER-CHART-002
status: 已通过
preconditions:
  - 用户账号可登录
  - 收入图表接口已联通
steps:
  - 切换到 /user/charts/income
  - 验证摘要卡
  - 校验趋势图和排行榜
expected:
  - 页面展示真实年收入趋势和收入前 10 类型排行榜
  - 模式切换与年份切换均正常
  - 摘要卡与折线图、排行榜口径一致
actual:
  - 页面展示年总收入 11200、月均收入 933.33、主要收入来源工资 86.61%
  - 收入趋势图和排行榜均正常渲染
  - 网络请求成功且控制台无未预期 error
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income.png
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 用户账号可登录
- 收入图表接口已联通

Steps:

1. 切换到 `/user/charts/income`
2. 验证年总收入、月均收入和主要收入来源摘要卡
3. 校验按月收入趋势与排行榜
4. 切换年份再次验证

Expected:

- 页面展示真实年收入趋势和收入前 10 类型排行榜
- 模式切换与年份切换均正常
- 摘要卡与折线图、排行榜口径一致

Actual:

- 收入图表页成功渲染真实收入数据，摘要卡显示年总收入 `11200`、月均收入 `933.33`、主要收入来源 `工资` 占比 `86.61%`。
- 收入趋势图与排行榜切换正常，年份切换不报错，页面消费真实接口返回数据。
- 关键请求成功，未发现未预期控制台 error。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income.png`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-01 预算与图表闭环回归

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-01
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-05
primary_req_id: USER-BUDGET-001
status: 已通过
preconditions:
  - 本批次功能已联调完成
steps:
  - 从用户端预算月页设置预算
  - 进入年预算页检查执行情况
  - 进入支出图表页和收入图表页完成切换
expected:
  - 预算设置后在预算页和图表页使用同一份真实 ledger 口径
  - 页面跳转、切换和数据回显连续一致
actual:
  - 登录到月预算设置、年预算查看、支出图表、收入图表的整条链路已跑通
  - 页面跳转、模式切换和数据回显连续稳定
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-login-snapshot.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-after-save.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-after-save.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-expense.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income.txt
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 本批次功能已联调完成

Steps:

1. 从用户端预算月页设置预算
2. 进入年预算页检查执行情况
3. 进入支出图表页和收入图表页完成切换

Expected:

- 预算设置后在预算页和图表页使用同一份真实 ledger 口径
- 页面跳转、切换和数据回显连续一致

Actual:

- 从用户登录到月预算设置、年预算查看、支出图表、收入图表的整条链路已连续跑通。
- 月预算页、年预算页和图表页对同一批真实 ledger 数据保持一致口径，预算与支出统计前后一致。
- 页面跳转、模式切换和数据回显连续稳定，闭环通过。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-login-snapshot.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-after-save.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-after-save.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-expense.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income.txt`

### TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-02 控制台与关键请求回归

```yaml
test_case_id: TC-ST-BATCH-20260410-BUDGET-CHART-REAL-05-02
related_subtask_id: ST-BATCH-20260410-BUDGET-CHART-REAL-05
primary_req_id: USER-CHART-001
status: 已通过
preconditions:
  - 本批次功能已联调完成
steps:
  - 在预算月页、预算年页、支出图表页、收入图表页抓取控制台与网络请求
  - 检查是否存在未预期报错、4xx 或 5xx
expected:
  - 关键请求全部成功
  - 无未预期控制台 error
actual:
  - 四个页面的关键接口请求均成功，未发现 4xx/5xx
  - 页面控制台未出现未预期 error
  - Playwright MCP 初始化异常已登记为 BUG-BATCH-20260410-BUDGET-CHART-REAL-002，wrapper 回归通过
evidence:
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income-network.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-console.txt
  - docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-playwright-browser-init-error.txt
executor: Codex
last_executed_at: 2026-04-10 12:54
```

Preconditions:

- 本批次功能已联调完成

Steps:

1. 在预算月页、预算年页、支出图表页、收入图表页抓取控制台与网络请求
2. 检查是否存在未预期报错、4xx 或 5xx

Expected:

- 关键请求全部成功
- 无未预期控制台 error

Actual:

- 月预算页、年预算页、支出图表页和收入图表页的关键接口请求均返回成功状态，未发现 4xx/5xx。
- 页面控制台未出现未预期 error。
- 直连 Playwright MCP 初始化目录失败已按规则登记为 `BUG-BATCH-20260410-BUDGET-CHART-REAL-002`，随后切换本机 wrapper 完成同批次真实回归，最终不阻塞验收。

Evidence:

- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-income-network.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-month-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-budget-year-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-user-charts-console.txt`
- `docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-playwright-browser-init-error.txt`
