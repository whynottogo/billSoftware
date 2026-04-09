# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 5 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-01 注册默认禁用

```yaml
test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-01
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
primary_req_id: USER-AUTH-001
status: 已通过
preconditions:
  - 用户注册页可访问
steps:
  - 提交新用户注册
expected:
  - 注册成功并提示等待管理员启用
  - 该用户默认不可登录
actual: 用户 `e2e0410000402` 注册成功，页面回到 `/user/login?registered=1` 并提示“账号创建成功，需等待管理员启用后才能登录”；随后在管理员真实列表中看到该用户状态为禁用。
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0007-register-success.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0007-admin-users-real.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:07:06
```

### TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-02 管理端真实列表可见新用户

```yaml
test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-02
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
primary_req_id: ADMIN-USER-001
status: 已通过
preconditions:
  - 已注册测试用户
steps:
  - 管理员登录并进入 /admin/users
expected:
  - 新用户出现在真实用户列表中，状态为禁用
actual: 管理员真实登录后进入 `/admin/users`，可见 `e2e0410000402`，状态为禁用。
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0007-admin-users-real.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:07:40
```

### TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-03 启用后登录成功

```yaml
test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-03
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
primary_req_id: ADMIN-USER-002
status: 已通过
preconditions:
  - 目标用户已存在于管理员列表
steps:
  - 启用该用户
  - 在用户端登录
expected:
  - 登录成功并进入 /user/ledger
actual: 管理员启用 `e2e0410000402` 后，用户端使用该账号成功登录并进入 `/user/ledger`。
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0007-admin-enable-success.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0008-user-ledger-login-success.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:08:14
```

### TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-04 禁用后登录失败

```yaml
test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-04
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
primary_req_id: ADMIN-USER-002
status: 已通过
preconditions:
  - 目标用户已被启用并可登录
steps:
  - 禁用该用户
  - 再次在用户端登录
expected:
  - 登录失败并展示禁用或失效提示
actual: 管理员禁用用户后，重新登录页面仍停留在 `/user/login`，页面提示 `user is disabled`。
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0008-admin-disable-success.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0009-disabled-login-page.yml
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:10:13
```

### TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-05 旧会话失效回跳

```yaml
test_case_id: TC-ST-BATCH-20260409-CLOSED-LOOP-QA-01-05
related_subtask_id: ST-BATCH-20260409-CLOSED-LOOP-QA-01
primary_req_id: SYS-AUTH-002
status: 已通过
preconditions:
  - 目标用户已有活跃登录态
steps:
  - 在管理员端禁用该用户
  - 用旧会话访问受保护用户接口
expected:
  - 前端收到 401 后跳回登录页并提示会话失效
actual: 旧会话访问 `/user/ledger` 时立即弹出“账号已在其他地方登录，请重新登录。”提示；后端日志同时记录 `/api/user/categories` 与 `/api/user/ledger` 返回 401。
evidence:
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0008-session-kick-console.txt
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0008-session-kick-network.txt
  - docs/task-batches/BATCH-20260409-CLOSED-LOOP-QA/assets/playwright/20260410-0010-backend-auth-evidence.txt
executor: Codex
last_executed_at: 2026-04-10 00:08:54
```
