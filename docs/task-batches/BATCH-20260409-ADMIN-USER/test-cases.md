# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 3 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-ADMIN-USER-01-01 用户列表启用/禁用交互

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-01-01
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-01
primary_req_id: ADMIN-USER-002
status: 已通过
preconditions:
  - 已具备管理端登录态，或手动写入 bill_admin_token
  - 浏览器可访问 /admin/users
steps:
  - 打开 /admin/users
  - 对禁用用户执行启用
  - 对启用用户执行禁用
expected:
  - 行状态标签即时变化
  - 顶部启用/禁用统计同步变化
actual: 已验证启用/禁用按钮可切换状态，统计数据同步更新，并可进入用户详情页
evidence:
  - assets/playwright/tc01-admin-users-open.txt
  - assets/playwright/tc01-admin-users-toggle.txt
  - assets/playwright/tc01-admin-users-snapshot.txt
  - assets/playwright/tc01-admin-users.png
  - assets/playwright/tc01-admin-users-url.txt
  - assets/playwright/tc01-admin-users-console-error.txt
  - assets/playwright/tc01-admin-users-console.txt
  - assets/playwright/tc01-admin-users-network.txt
executor: Codex
last_executed_at: 2026-04-09 22:03
```

### TC-ST-BATCH-20260409-ADMIN-USER-02-01 用户详情页（月/年 tab）展示

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-02-01
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-02
primary_req_id: ADMIN-USER-003
status: 已通过
preconditions:
  - 已具备管理端登录态，或手动写入 bill_admin_token
  - 浏览器可访问 /admin/users/:userId
steps:
  - 从用户列表进入详情页
  - 在详情页切换“月账单”和“年账单”tab
  - 切换月份/年份检查汇总数据
expected:
  - 详情页包含月/年两个 tab
  - 可查看单用户月/年收入、支出、结余
  - tab 和筛选切换后信息更新正确
actual: 已验证用户详情页可展示月/年 tab，并完成 tab 与年份切换
evidence:
  - assets/playwright/tc02-admin-user-detail-snapshot.txt
  - assets/playwright/tc02-admin-user-detail-switch.txt
  - assets/playwright/tc02-admin-user-detail.png
  - assets/playwright/tc02-admin-user-detail-url.txt
  - assets/playwright/tc02-admin-user-detail-console-error.txt
  - assets/playwright/tc02-admin-user-detail-console.txt
  - assets/playwright/tc02-admin-user-detail-network.txt
executor: Codex
last_executed_at: 2026-04-09 22:03
```

### TC-ST-BATCH-20260409-ADMIN-USER-03-01 管理端用户模块回归（编译+控制台）

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-03-01
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-03
primary_req_id: ADMIN-USER-003
status: 已通过
preconditions:
  - 前端服务可访问
  - 管理端相关页面已完成实现
steps:
  - 执行 npm run build
  - 回归访问 /admin/users 与 /admin/users/:userId
  - 检查 console error 与关键网络请求
expected:
  - 构建通过
  - 页面可访问且无 console error
actual: 已完成构建回归且管理端用户列表/详情页控制台 error 均为 0
evidence:
  - assets/playwright/frontend-build-regression.txt
  - assets/playwright/tc01-admin-users-console-error.txt
  - assets/playwright/tc02-admin-user-detail-console-error.txt
executor: Codex
last_executed_at: 2026-04-09 22:04
```
