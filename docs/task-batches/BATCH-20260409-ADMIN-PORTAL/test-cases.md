# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 4 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-ADMIN-PORTAL-01-01 数据概览页展示

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-PORTAL-01-01
related_subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-01
primary_req_id: ADMIN-DASHBOARD-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:46
evidence:
  - assets/playwright/admin-dashboard.png
  - assets/playwright/admin-dashboard-snapshot.txt
  - assets/playwright/admin-dashboard-url.txt
  - assets/playwright/admin-dashboard-console.txt
  - assets/playwright/admin-dashboard-network.txt
```

Preconditions:

- 已具备管理端登录态，或手动写入 `bill_admin_token`
- 浏览器可访问 `/admin/dashboard`

Steps:

1. 打开 `/admin/dashboard`
2. 检查统计卡、增长趋势、账单趋势、最近注册用户区块
3. 验证页面无控制台异常

Expected:

- 页面展示后台概览核心指标与图表
- 页面结构与 Make 原型一致，信息层级清晰
- 控制台不出现报错

Actual:

- 页面可正常访问并展示统计卡、双趋势图和最近注册用户区块
- 路由停留在 `http://localhost:9000/admin/dashboard`
- 控制台 `error` 数量为 0

Evidence:

- `assets/playwright/admin-dashboard.png`
- `assets/playwright/admin-dashboard-snapshot.txt`
- `assets/playwright/admin-dashboard-url.txt`
- `assets/playwright/admin-dashboard-console.txt`
- `assets/playwright/admin-dashboard-network.txt`

### TC-ST-BATCH-20260409-ADMIN-PORTAL-01-02 后台门户路由可访问性

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-PORTAL-01-02
related_subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-01
primary_req_id: ADMIN-DASHBOARD-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:46
evidence:
  - assets/playwright/admin-dashboard-url.txt
  - assets/playwright/admin-approvals-url.txt
  - assets/playwright/admin-families-url.txt
```

Preconditions:

- 已具备管理端登录态，或手动写入 `bill_admin_token`
- 临时路由入口已在本地调试环境接入

Steps:

1. 访问 `/admin/dashboard`
2. 访问 `/admin/approvals`
3. 访问 `/admin/families`

Expected:

- 三个路由都能进入对应页面
- 无重定向循环或空白页

Actual:

- 三个路由均可访问：`/admin/dashboard`、`/admin/approvals`、`/admin/families`
- 页面均可完成渲染，未出现重定向循环和空白页

Evidence:

- `assets/playwright/admin-dashboard-url.txt`
- `assets/playwright/admin-approvals-url.txt`
- `assets/playwright/admin-families-url.txt`

### TC-ST-BATCH-20260409-ADMIN-PORTAL-02-01 待审批用户页展示与操作占位

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-PORTAL-02-01
related_subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-02
primary_req_id: ADMIN-APPROVAL-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:46
evidence:
  - assets/playwright/admin-approvals.png
  - assets/playwright/admin-approvals-snapshot.txt
  - assets/playwright/admin-approvals-url.txt
  - assets/playwright/admin-approvals-console.txt
  - assets/playwright/admin-approvals-network.txt
  - assets/playwright/admin-portal-playwright-session-error.log
```

Preconditions:

- 已具备管理端登录态，或手动写入 `bill_admin_token`
- 浏览器可访问 `/admin/approvals`

Steps:

1. 打开 `/admin/approvals`
2. 检查待审批用户列表信息字段
3. 检查批准/拒绝/批量按钮入口可见

Expected:

- 页面展示待审批用户卡片列表
- 审批操作入口可见且布局完整
- 控制台无报错

Actual:

- 页面可正常展示审批列表和批量操作区域
- 批准/拒绝/批量按钮入口均可见
- 直接文本点击命令触发 Playwright 选择器误用异常，已单独记录并关闭为流程类 Bug，不影响页面可见性验收
- 控制台 `error` 数量为 0

Evidence:

- `assets/playwright/admin-approvals.png`
- `assets/playwright/admin-approvals-snapshot.txt`
- `assets/playwright/admin-approvals-url.txt`
- `assets/playwright/admin-approvals-console.txt`
- `assets/playwright/admin-approvals-network.txt`
- `assets/playwright/admin-portal-playwright-session-error.log`

### TC-ST-BATCH-20260409-ADMIN-PORTAL-03-01 家庭管理页展示

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-PORTAL-03-01
related_subtask_id: ST-BATCH-20260409-ADMIN-PORTAL-03
primary_req_id: ADMIN-FAMILY-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:46
evidence:
  - assets/playwright/admin-families.png
  - assets/playwright/admin-families-snapshot.txt
  - assets/playwright/admin-families-url.txt
  - assets/playwright/admin-families-console.txt
  - assets/playwright/admin-families-network.txt
```

Preconditions:

- 已具备管理端登录态，或手动写入 `bill_admin_token`
- 浏览器可访问 `/admin/families`

Steps:

1. 打开 `/admin/families`
2. 检查家庭统计卡与家庭列表
3. 检查查看/更多/分页按钮入口可见

Expected:

- 页面展示家庭核心统计与卡片列表
- 操作按钮入口可见且结构完整
- 控制台无报错

Actual:

- 页面可正常展示家庭统计、家庭卡片列表和分页壳
- 查看/更多/分页按钮入口可见
- 控制台 `error` 数量为 0

Evidence:

- `assets/playwright/admin-families.png`
- `assets/playwright/admin-families-snapshot.txt`
- `assets/playwright/admin-families-url.txt`
- `assets/playwright/admin-families-console.txt`
- `assets/playwright/admin-families-network.txt`
