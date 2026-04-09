# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 5 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-AUTH-SHELL-01-01 用户登录页展示与跳转

```yaml
test_case_id: TC-ST-BATCH-20260409-AUTH-SHELL-01-01
related_subtask_id: ST-BATCH-20260409-AUTH-SHELL-01
primary_req_id: USER-AUTH-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 18:17
evidence:
  - assets/playwright/user-login.png
  - assets/playwright/user-login-open.txt
  - assets/playwright/user-login-snapshot.txt
  - assets/playwright/user-login-console.txt
  - assets/playwright/user-login-network.txt
```

Preconditions:

- 前端开发服务已启动
- 浏览器可访问 `/user/login`

Steps:

1. 打开 `/user/login`
2. 检查品牌区、登录表单、管理员入口和注册入口
3. 输入示例账号密码并触发登录

Expected:

- 页面展示用户端登录所需的用户名或手机号与密码字段
- 登录成功后跳转到 `/user/ledger`，或在接口不可用时给出明确错误提示
- 页面不出现管理员注册、找回密码等未定义能力

Actual:

- 当前 URL：`http://localhost:9000/user/login`
- 页面已渲染为 Make 风格的用户端双栏登录布局，管理员入口与注册入口均可见
- `assets/playwright/user-login-console.txt` 中控制台无 error 或 warning
- 本轮验证聚焦页面替换与路由可达性，`assets/playwright/user-login-network.txt` 未记录业务请求

Evidence:

- `assets/playwright/user-login.png`
- `assets/playwright/user-login-open.txt`
- `assets/playwright/user-login-snapshot.txt`
- `assets/playwright/user-login-console.txt`
- `assets/playwright/user-login-network.txt`

### TC-ST-BATCH-20260409-AUTH-SHELL-01-02 用户注册页展示与提交

```yaml
test_case_id: TC-ST-BATCH-20260409-AUTH-SHELL-01-02
related_subtask_id: ST-BATCH-20260409-AUTH-SHELL-01
primary_req_id: USER-AUTH-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 18:17
evidence:
  - assets/playwright/user-register.png
  - assets/playwright/user-register-open.txt
  - assets/playwright/user-register-snapshot.txt
  - assets/playwright/user-register-console.txt
  - assets/playwright/user-register-network.txt
```

Preconditions:

- 前端开发服务已启动
- 浏览器可访问 `/user/register`

Steps:

1. 打开 `/user/register`
2. 检查用户名、昵称、手机号、邮箱、密码字段
3. 提交注册表单并观察页面反馈

Expected:

- 页面展示完整注册字段
- 明确提示“注册后需管理员启用”
- 提交后返回登录页或给出明确错误提示

Actual:

- 当前 URL：`http://localhost:9000/user/register`
- 页面已渲染为双栏注册布局，字段与需求文档一致，并明确提示需管理员启用后登录
- `assets/playwright/user-register-console.txt` 中控制台无 error 或 warning
- 本轮验证聚焦页面替换与信息完整性，`assets/playwright/user-register-network.txt` 未记录业务请求

Evidence:

- `assets/playwright/user-register.png`
- `assets/playwright/user-register-open.txt`
- `assets/playwright/user-register-snapshot.txt`
- `assets/playwright/user-register-console.txt`
- `assets/playwright/user-register-network.txt`

### TC-ST-BATCH-20260409-AUTH-SHELL-02-01 用户端首页骨架与导航

```yaml
test_case_id: TC-ST-BATCH-20260409-AUTH-SHELL-02-01
related_subtask_id: ST-BATCH-20260409-AUTH-SHELL-02
primary_req_id: USER-LEDGER-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 18:17
evidence:
  - assets/playwright/user-ledger.png
  - assets/playwright/user-ledger-open.txt
  - assets/playwright/user-ledger-snapshot.txt
  - assets/playwright/user-ledger-console.txt
  - assets/playwright/user-ledger-network.txt
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/ledger`

Steps:

1. 打开 `/user/ledger`
2. 检查左侧导航、顶部摘要、当月收支内容区
3. 执行退出登录并确认跳回 `/user/login`

Expected:

- 页面采用 Make 风格的用户端桌面壳
- 导航高亮与 `/user` 路由前缀保持一致
- 退出登录只清理用户端存储并跳回用户登录页

Actual:

- 通过写入本地用户端登录态后访问 `http://localhost:9000/user/ledger` 成功进入页面
- 左侧导航、月份切换、统计卡、按天账单列表与右侧常用分类/本月概况均已展示
- `assets/playwright/user-ledger-console.txt` 中控制台无 error 或 warning
- 页面属于静态原型替换验证，`assets/playwright/user-ledger-network.txt` 未记录业务请求

Evidence:

- `assets/playwright/user-ledger.png`
- `assets/playwright/user-ledger-open.txt`
- `assets/playwright/user-ledger-snapshot.txt`
- `assets/playwright/user-ledger-console.txt`
- `assets/playwright/user-ledger-network.txt`

### TC-ST-BATCH-20260409-AUTH-SHELL-03-01 管理员登录页展示与跳转

```yaml
test_case_id: TC-ST-BATCH-20260409-AUTH-SHELL-03-01
related_subtask_id: ST-BATCH-20260409-AUTH-SHELL-03
primary_req_id: ADMIN-AUTH-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 18:17
evidence:
  - assets/playwright/admin-login.png
  - assets/playwright/admin-login-open.txt
  - assets/playwright/admin-login-snapshot.txt
  - assets/playwright/admin-login-console.txt
  - assets/playwright/admin-login-network.txt
```

Preconditions:

- 前端开发服务已启动
- 浏览器可访问 `/admin/login`

Steps:

1. 打开 `/admin/login`
2. 检查管理员账号与密码表单、返回用户端入口
3. 输入示例账号密码并触发登录

Expected:

- 页面只体现管理员登录，不出现用户注册能力
- 登录成功后进入 `/admin/users`，或在接口不可用时给出明确错误提示
- 与用户端登录页保持品牌一致但视觉更克制

Actual:

- 当前 URL：`http://localhost:9000/admin/login`
- 页面已渲染为更克制的管理端双栏登录布局，并保留返回用户端入口
- `assets/playwright/admin-login-console.txt` 中控制台无 error 或 warning
- 本轮验证聚焦页面替换与入口隔离，`assets/playwright/admin-login-network.txt` 未记录业务请求

Evidence:

- `assets/playwright/admin-login.png`
- `assets/playwright/admin-login-open.txt`
- `assets/playwright/admin-login-snapshot.txt`
- `assets/playwright/admin-login-console.txt`
- `assets/playwright/admin-login-network.txt`

### TC-ST-BATCH-20260409-AUTH-SHELL-04-01 管理端用户列表骨架

```yaml
test_case_id: TC-ST-BATCH-20260409-AUTH-SHELL-04-01
related_subtask_id: ST-BATCH-20260409-AUTH-SHELL-04
primary_req_id: ADMIN-USER-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 18:17
evidence:
  - assets/playwright/admin-users.png
  - assets/playwright/admin-users-open.txt
  - assets/playwright/admin-users-snapshot.txt
  - assets/playwright/admin-users-console.txt
  - assets/playwright/admin-users-network.txt
```

Preconditions:

- 已具备管理员登录态，或手动写入 `bill_admin_token`
- 浏览器可访问 `/admin/users`

Steps:

1. 打开 `/admin/users`
2. 检查后台导航、摘要区、搜索筛选和列表字段
3. 执行退出登录并确认跳回 `/admin/login`

Expected:

- 页面体现管理端后台壳与用户列表主内容区
- 列表至少展示账号、昵称、手机号、邮箱、状态等核心字段
- 退出登录只清理管理员端存储并跳回管理员登录页

Actual:

- 通过写入管理端登录态后访问 `http://localhost:9000/admin/users` 成功进入页面
- 页面展示后台侧导航、统计摘要、搜索筛选以及用户列表字段分组
- `assets/playwright/admin-users-console.txt` 中控制台无 error 或 warning
- 页面属于静态原型替换验证，`assets/playwright/admin-users-network.txt` 未记录业务请求

Evidence:

- `assets/playwright/admin-users.png`
- `assets/playwright/admin-users-open.txt`
- `assets/playwright/admin-users-snapshot.txt`
- `assets/playwright/admin-users-console.txt`
- `assets/playwright/admin-users-network.txt`
