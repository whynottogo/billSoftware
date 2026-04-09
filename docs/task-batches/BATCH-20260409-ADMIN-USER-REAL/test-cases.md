# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 2 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-ADMIN-USER-REAL-01-01 用户列表真实渲染

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-REAL-01-01
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-REAL-01
primary_req_id: ADMIN-USER-001
status: 已通过
preconditions:
  - 管理员账号可登录
  - 本地前后端均已启动
steps:
  - 管理员登录并进入 /admin/users
  - 校验列表、统计卡和状态筛选
expected:
  - 列表来自真实接口且不暴露密码字段
  - 新注册用户默认显示为禁用
actual: 管理员真实登录后进入 `/admin/users`，列表显示 `e2e0410000402`、`e2e0410000401` 两个真实用户，状态均为禁用，页面未出现任何密码字段。
evidence:
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/20260410-0007-admin-users-real.yml
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/20260410-0010-admin-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:07:40
```

### TC-ST-BATCH-20260409-ADMIN-USER-REAL-02-01 用户启停影响登录

```yaml
test_case_id: TC-ST-BATCH-20260409-ADMIN-USER-REAL-02-01
related_subtask_id: ST-BATCH-20260409-ADMIN-USER-REAL-02
primary_req_id: ADMIN-USER-002
status: 已通过
preconditions:
  - 已注册测试用户 `e2e0410000402`
  - 管理员可在用户列表执行启停
steps:
  - 启用目标用户
  - 在用户端验证可登录
  - 再次禁用并验证旧会话失效与重新登录失败
expected:
  - 启用后用户端可登录并进入 `/user/ledger`
  - 禁用后旧会话触发失效提示，重新登录失败
actual: 启用后，用户 `e2e0410000402` 成功登录到 `/user/ledger`；禁用后，旧会话访问用户接口返回 401 并触发“账号已在其他地方登录，请重新登录。”提示，随后重新登录返回 `user is disabled`。
evidence:
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/20260410-0007-admin-enable-success.yml
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/20260410-0008-admin-disable-success.yml
  - docs/task-batches/BATCH-20260409-ADMIN-USER-REAL/assets/playwright/20260410-0010-admin-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:10:13
```
