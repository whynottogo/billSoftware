# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 6 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-01 家庭接口结构验证

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-01
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-01
primary_req_id: USER-FAMILY-001
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 请求家庭列表接口
  - 请求创建家庭接口并记录返回的家庭 ID
  - 请求按家庭 ID 加入、按邀请链接加入和退出接口
expected:
  - 家庭列表包含家庭名称、创建人、家庭 ID 与成员信息
  - 创建接口返回独立家庭 ID 与邀请链接
  - 加入和退出接口成功且只作用于当前登录用户
actual:
  - 用户 2 成功创建家庭 `FAM-18A4FBF2C4ADE858-2`
  - 用户 3 按家庭 ID 加入成功，首次退出成功，再按邀请链接加入成功
  - 用户 3 最终退出后再次访问详情返回 `404 family not found`
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-phase1-summary.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-create-family.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-join-by-id.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user3-leave-first.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-join-by-link.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user3-leave-final.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user3-detail-after-leave.json
executor: Codex
last_executed_at: 2026-04-10 20:11
```

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-02 家庭详情与成员占比接口结构验证

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-02
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-01
primary_req_id: USER-FAMILY-002
status: 已通过
preconditions:
  - 后端已启动
  - 已获取家庭详情可访问的 familyId
steps:
  - 请求家庭详情接口
  - 按月收入、月支出、年收入、年支出四种口径请求成员占比接口
expected:
  - 详情接口返回 monthOptions、yearOptions、members、inviteLink 等详情页所需字段
  - 成员占比接口返回 title、total、rows，且口径切换正确
actual:
  - 家庭详情接口返回成员列表、`monthOptions`、`yearOptions`、`inviteLink` 和 UI 所需汇总字段
  - 四种成员占比接口均返回 `200`，并正确反映用户 2 的真实流水与用户 3 的 0 值占比
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-family-detail.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-member-share-month-income.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-member-share-month-expense.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-member-share-year-income.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-member-share-year-expense.json
executor: Codex
last_executed_at: 2026-04-10 20:11
```

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-01 家庭列表真实展示与详情跳转

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-01
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-02
primary_req_id: USER-FAMILY-001
status: 已通过
preconditions:
  - 用户账号可登录
  - 家庭列表接口已联通
steps:
  - 登录用户端并进入 /user/families
  - 检查家庭名称、创建人、成员数和家庭 ID
  - 点击任一家庭查看详情
expected:
  - 列表展示真实家庭数据，不再依赖 mock
  - 点击后能进入 /user/families/:familyId
actual:
  - 用户 2 登录后，`/user/families` 正确展示真实家庭名称、创建人、成员数和家庭 ID
  - 点击“查看详情”后成功进入 `http://localhost:9000/user/families/FAM-18A4FBF2C4ADE858-2`
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success-console.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success-network.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success.png
executor: Codex
last_executed_at: 2026-04-10 20:11
```

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-02 创建加入退出家庭真实闭环

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-02
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-02
primary_req_id: USER-FAMILY-001
status: 已通过
preconditions:
  - 用户账号可登录
  - 家庭创建、加入、退出接口已联通
steps:
  - 在 /user/families 创建一个新家庭
  - 使用家庭 ID 加入家庭
  - 使用邀请链接加入家庭
  - 退出一个已加入的家庭
expected:
  - 创建后页面提示并展示独立家庭 ID
  - 按 ID 与按链接加入都成功
  - 退出后列表、统计和详情入口同步更新
actual:
  - 真实链路已完成“创建家庭 -> 按 ID 加入 -> 退出 -> 按邀请链接加入 -> 最终退出 -> 创建人删除测试家庭”
  - 清理完成后，列表接口和页面均回到 0 条家庭记录的空态
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-create-family.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-join-by-id.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-join-by-link.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user3-leave-final.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user2-leave-cleanup-browser-token.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user2-list-after-cleanup-browser-token.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-empty-after-cleanup.png
executor: Codex
last_executed_at: 2026-04-10 20:11
```

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-03-01 家庭详情真实汇总与成员占比切换

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-03-01
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-03
primary_req_id: USER-FAMILY-002
status: 已通过
preconditions:
  - 用户账号可登录
  - 家庭详情和成员占比接口已联通
steps:
  - 打开 /user/families/:familyId
  - 检查月度与年度收入、支出、结余
  - 点击四个占比入口观察图表和列表变化
expected:
  - 页面展示真实月度和年度汇总
  - 点击不同入口后成员占比标题、总额和 rows 正确变化
actual:
  - 详情页成功展示 2026 年 4 月和 2026 年的真实收入、支出、结余汇总
  - 页面依次触发月收入、月支出、年收入、年支出四种成员占比请求，均返回 `200` 且标题、总额和成员 rows 正确切换
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success-console.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success-network.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-family-detail.json
executor: Codex
last_executed_at: 2026-04-10 20:11
```

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01 家庭真实页面回归

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-04
primary_req_id: USER-FAMILY-001
status: 已通过
preconditions:
  - 前后端服务已启动
  - 家庭前后端改造已完成
steps:
  - 回归登录 -> 家庭列表 -> 创建 -> 详情 -> 占比切换 -> 返回列表
  - 抓取截图、URL、控制台和关键网络请求
expected:
  - 全链路可用，页面无未预期 error
  - 关键网络请求均返回 2xx
actual:
  - 已完成“登录 -> 家庭列表真实展示 -> 详情页汇总与四种占比切换 -> 用户退出 -> 详情访问回收 -> 创建人清理 -> 列表空态恢复”全链路回归
  - Playwright 控制台无 error，关键接口请求均返回 2xx；访问回收验证返回预期 `404 family not found`
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success-console.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success-network.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user3-detail-after-leave.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-empty-after-cleanup.png
executor: Codex
last_executed_at: 2026-04-10 20:11
```

### TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-02 构建与服务健康检查

```yaml
test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-02
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-04
primary_req_id: USER-FAMILY-002
status: 已通过
preconditions:
  - 改造代码已落盘
steps:
  - 在 frontend 执行 npm run build
  - 在 backend 执行 go test ./...
  - 请求 GET /api/health
expected:
  - 前端构建通过
  - 后端测试通过
  - 健康检查返回成功
actual:
  - 前端构建通过，仅保留 webpack 体积告警
  - 后端 `go test ./...` 全部通过
  - `GET /api/health` 返回服务正常
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-frontend-build.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-backend-go-test.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-health-final.json
executor: Codex
last_executed_at: 2026-04-10 20:11
```
