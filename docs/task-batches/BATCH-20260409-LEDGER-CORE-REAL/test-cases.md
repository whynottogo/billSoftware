# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 4 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-LEDGER-CORE-REAL-01-01 用户隔离与禁用失效

```yaml
test_case_id: TC-ST-BATCH-20260409-LEDGER-CORE-REAL-01-01
related_subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-01
primary_req_id: SYS-DATA-001
status: 已通过
preconditions:
  - 至少存在两个已启用用户
  - 其中一个用户已有活跃会话
steps:
  - 校验用户 1 和用户 2 读取到的 ledger/category 数据不同
  - 管理员禁用用户 2
  - 校验旧会话访问 `/api/user/ledger`、`/api/user/categories` 返回 401
expected:
  - 任何查询或写入都不跨用户
  - 禁用后旧会话不可继续访问受保护接口
actual: 用户 2 新增了 1 条支出和 1 条收入；用户 1 读取 `2026-04` 月账本时 `groups` 为空、`summary` 为 0，确认未串数据。管理员禁用用户 2 后，旧会话访问 `/api/user/ledger` 与 `/api/user/categories` 返回 401。
evidence:
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-user1-isolation-ledger.json
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0010-session-invalidation-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:14:41
```

### TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-01 当月收支真实读取

```yaml
test_case_id: TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-01
related_subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-02
primary_req_id: USER-LEDGER-001
status: 已通过
preconditions:
  - 已启用用户可登录
  - 目标月份存在至少一条记录
steps:
  - 登录后进入 /user/ledger
  - 校验 summary/groups/categories/overview 展示
expected:
  - 页面使用真实接口返回的数据
  - 月份切换后仍只看到本人记录
actual: 用户 `e2e0410000402` 登录后进入 `/user/ledger`，页面展示真实默认分类；新增支出和新增收入后，月汇总、日分组和常用分类都按真实接口返回结果更新。
evidence:
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0013-expense-saved.yml
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-income-saved.yml
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-ledger-backend-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:14:04
```

### TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-02 新增收支即时回显

```yaml
test_case_id: TC-ST-BATCH-20260409-LEDGER-CORE-REAL-02-02
related_subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-02
primary_req_id: USER-LEDGER-002
status: 已通过
preconditions:
  - 已启用用户可登录
  - 至少存在一个收入分类和一个支出分类
steps:
  - 在 /user/ledger 新增支出
  - 在 /user/ledger 新增收入
expected:
  - 保存后当前月份列表即时出现对应记录
  - 月汇总和日汇总同步变化
actual: 新增 `午餐测试` 支出 50 元后，本月支出变为 `¥50`；新增 `工资测试` 收入 100 元后，本月收入变为 `¥100`、结余变为 `¥50`，并在 2026-04-10 分组中即时出现两条记录。
evidence:
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0013-expense-saved.yml
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-income-saved.yml
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-ledger-backend-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:14:04
```

### TC-ST-BATCH-20260409-LEDGER-CORE-REAL-03-01 分类管理真实化

```yaml
test_case_id: TC-ST-BATCH-20260409-LEDGER-CORE-REAL-03-01
related_subtask_id: ST-BATCH-20260409-LEDGER-CORE-REAL-03
primary_req_id: USER-LEDGER-004
status: 已通过
preconditions:
  - 已启用用户可登录
steps:
  - 打开分类管理对话框
  - 新增自定义分类
  - 删除自定义分类并尝试删除默认分类
expected:
  - 自定义分类可新增删除
  - 默认分类不可删除
  - 新增分类可立即出现在记账表单中
actual: 支出分类新增 `测试` 成功并显示为“自定义分类”，随后删除成功；默认分类的删除按钮保持禁用态。后端日志显示 `POST /api/user/categories` 与 `DELETE /api/user/categories/:id` 均返回 200。
evidence:
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0012-category-added.yml
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0012-category-deleted.yml
  - docs/task-batches/BATCH-20260409-LEDGER-CORE-REAL/assets/playwright/20260410-0014-ledger-backend-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:12:56
```
