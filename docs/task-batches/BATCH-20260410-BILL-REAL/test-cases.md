# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 8 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260410-BILL-REAL-01-01 用户账单接口结构验证

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-01-01
related_subtask_id: ST-BATCH-20260410-BILL-REAL-01
primary_req_id: USER-BILL-001
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 请求用户端年度账单列表接口
  - 请求用户端月账单详情接口
expected:
  - 接口返回结构满足前端页面消费字段
  - 统计结果仅包含当前用户数据
actual: 使用用户 `e2e0410000402` 的真实 token 校验 `/api/user/bills/years`、`/api/user/bills/year/2026` 与 `/api/user/bills/month/2025-12`，返回包含 `years/history`、年度 `summary/months`、月度 `categorySplit/ranking/dailyTrend/comparison/achievements` 等页面消费字段，且仅返回该用户的 2025/2026 账单数据。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-user-bills-years.json
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-user-bills-year-2026.json
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-user-bills-month-2025-12.json
executor: Codex
last_executed_at: 2026-04-10 02:18
```

### TC-ST-BATCH-20260410-BILL-REAL-01-02 管理端账单接口结构验证

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-01-02
related_subtask_id: ST-BATCH-20260410-BILL-REAL-01
primary_req_id: ADMIN-USER-003
status: 已通过
preconditions:
  - 后端已启动
  - 已获取管理员 token
steps:
  - 请求管理端单用户月账单接口
  - 请求管理端单用户年账单接口
expected:
  - 接口返回结构满足管理端详情页消费字段
  - 与用户端同月份/同年份统计口径一致
actual: 使用管理员 token 请求 `/api/admin/users/2/bills/overview`，返回 `profile/monthOptions/yearOptions/monthly/yearly` 结构；其中 `2026-04` 月数据为 `3500/350/3150/5`，`2025` 年数据为 `6000/1130/4870/2个月`，与用户端口径一致。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-admin-user-2-overview.json
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-user-bills-year-2026.json
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-user-bills-month-2025-12.json
executor: Codex
last_executed_at: 2026-04-10 02:18
```

### TC-ST-BATCH-20260410-BILL-REAL-02-01 用户月账单列表真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-02-01
related_subtask_id: ST-BATCH-20260410-BILL-REAL-02
primary_req_id: USER-BILL-001
status: 已通过
preconditions:
  - 已启用用户可登录
  - 目标用户已有跨月 ledger 数据
steps:
  - 登录用户端并进入 `/user/bills/month`
  - 切换年份并点击最近月份
expected:
  - 顶部显示真实年度汇总
  - 月份列表显示真实月收入、月支出、月结余与记录情况
actual: 通过 Playwright 打开 `/user/bills/month`，页面真实展示 2026 年汇总 `收入 ¥11,200 / 支出 ¥1,450 / 结余 ¥9,750` 与 4 个真实月份；点击左箭头后正确切换到 2025 年，并显示 12 月与 11 月两条月账单，证明年份切换与月份列表均已接到真实接口。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-month-v2.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-month-v2-network.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-40-04-690Z.png
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-40-53-139Z.yml
executor: Codex
last_executed_at: 2026-04-10 02:41
```

### TC-ST-BATCH-20260410-BILL-REAL-02-02 用户年账单真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-02-02
related_subtask_id: ST-BATCH-20260410-BILL-REAL-02
primary_req_id: USER-BILL-003
status: 已通过
preconditions:
  - 已启用用户可登录
steps:
  - 进入 `/user/bills/year`
  - 切换不同年份查看总汇总和历年列表
expected:
  - 页面展示真实年度收入、支出、结余与月度趋势
actual: 通过 Playwright 打开 `/user/bills/year`，页面默认展示 2026 年总汇总与 1-4 月趋势；点击 2025 年历年卡片后，页面切换到 2025 年并展示 `收入 ¥6,000 / 支出 ¥1,130 / 结余 ¥4,870`，完成真实年度切换验证。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-year.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-year-network.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-42-35-571Z.png
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-43-17-403Z.yml
executor: Codex
last_executed_at: 2026-04-10 02:43
```

### TC-ST-BATCH-20260410-BILL-REAL-03-01 用户月账单详情真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-03-01
related_subtask_id: ST-BATCH-20260410-BILL-REAL-03
primary_req_id: USER-BILL-002
status: 已通过
preconditions:
  - 已启用用户可登录
  - 目标月份有收入和支出记录
steps:
  - 进入 `/user/bills/month/:month`
  - 校验月汇总、分类分布、排行、趋势、近 6 月对比和成就模块
expected:
  - 页面全部使用真实接口返回数据
actual: 通过月账单列表进入 `/user/bills/month/2025-12`，页面真实展示 `本月结余 ¥2,350 / 上月结余 ¥2,520 / 本月收入 ¥2,800 / 本月支出 ¥450`，并渲染分类分布、排行、支出趋势、近 6 个月对比和 4 张记账成就卡，无控制台 warning。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bill-detail-2025-12.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bill-detail-2025-12-network.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bill-detail-2025-12-console.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-41-43-640Z.png
executor: Codex
last_executed_at: 2026-04-10 02:42
```

### TC-ST-BATCH-20260410-BILL-REAL-04-01 管理端用户账单详情真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-04-01
related_subtask_id: ST-BATCH-20260410-BILL-REAL-04
primary_req_id: ADMIN-USER-003
status: 已通过
preconditions:
  - 管理员账号可登录
  - 目标用户已有跨月账单数据
steps:
  - 管理员进入 `/admin/users/:userId`
  - 切换月账单和年账单 tab
expected:
  - 页面显示真实月/年收入、支出、结余与记录数
  - 与用户端同月份或同年份统计口径一致
actual: 通过 Playwright 打开 `/admin/users/2`，月账单 tab 默认展示 `2026-04` 的 `收入 ¥3,500 / 支出 ¥350 / 结余 ¥3,150 / 5 笔`；切换到年账单 tab 后展示 2026 年汇总，再把下拉切到 2025 年后展示 `收入 ¥6,000 / 支出 ¥1,130 / 结余 ¥4,870 / 2 个月`，与用户端同口径。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-admin-user-detail.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-admin-user-detail-network.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-44-24-568Z.png
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-45-03-450Z.yml
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/.playwright-cli/page-2026-04-10T02-45-39-806Z.yml
executor: Codex
last_executed_at: 2026-04-10 02:45
```

### TC-ST-BATCH-20260410-BILL-REAL-05-01 用户端账单回归

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-05-01
related_subtask_id: ST-BATCH-20260410-BILL-REAL-05
primary_req_id: USER-BILL-001
status: 已通过
preconditions:
  - 本批次功能已联调完成
steps:
  - 从用户端月账单列表进入月账单详情，再进入年账单
expected:
  - 页面跳转正常，核心统计数据连续一致
actual: 用户端从 `/user/bills/month` 成功切换到 2025 年、进入 `2025-12` 月账单详情，再跳转到 `/user/bills/year` 查看年度总览，三个页面均保持真实数据与正常跳转。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-month-v2.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bill-detail-2025-12.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-year.txt
executor: Codex
last_executed_at: 2026-04-10 02:43
```

### TC-ST-BATCH-20260410-BILL-REAL-05-02 管理端与用户端口径对账回归

```yaml
test_case_id: TC-ST-BATCH-20260410-BILL-REAL-05-02
related_subtask_id: ST-BATCH-20260410-BILL-REAL-05
primary_req_id: ADMIN-USER-003
status: 已通过
preconditions:
  - 本批次功能已联调完成
steps:
  - 记录用户端指定月份和年份的统计值
  - 管理员打开同一用户同一月份和年份详情
expected:
  - 管理端与用户端同口径统计一致
actual: 已对同一用户完成管理端与用户端对账。用户端 2026-04 月账单为 `收入 ¥3,500 / 支出 ¥350 / 结余 ¥3,150`，管理端月账单默认值一致；用户端 2025 年账单为 `收入 ¥6,000 / 支出 ¥1,130 / 结余 ¥4,870`，管理端年账单切到 2025 后值一致。
evidence:
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-month-v2.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-user-bills-year.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-admin-user-detail.txt
  - docs/task-batches/BATCH-20260410-BILL-REAL/assets/playwright/20260410-api-admin-user-2-overview.json
executor: Codex
last_executed_at: 2026-04-10 02:45
```
