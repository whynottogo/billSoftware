# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 4 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-USER-PROFILE-SESSION-01-01 个人信息页查看编辑与头像占位

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-PROFILE-SESSION-01-01
related_subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-01
primary_req_id: USER-PROFILE-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:46
evidence:
  - assets/playwright/tc01-user-profile.png
  - assets/playwright/tc01-user-profile-url.txt
  - assets/playwright/tc01-user-profile-snapshot.txt
  - assets/playwright/tc01-user-profile-console.log
  - assets/playwright/tc01-user-profile-network.txt
```

Preconditions:

- 已写入用户端登录态
- 浏览器可访问 `/user/profile`

Steps:

1. 打开 `/user/profile`
2. 检查账户和手机号只读展示，昵称和邮箱可编辑
3. 修改昵称和邮箱并点击“保存资料”

Expected:

- 个人信息页加载正常
- 只读字段不可编辑，可编辑字段可保存
- 页面结构包含头像上传占位入口

Actual:

- 页面成功打开到 `http://localhost:9000/user/profile`
- 通过自动化修改昵称与邮箱并点击“保存资料”，页面保持稳定
- 控制台无 error/warning，网络日志未出现 4xx/5xx 业务请求

Evidence:

- `assets/playwright/tc01-user-profile.png`
- `assets/playwright/tc01-user-profile-url.txt`
- `assets/playwright/tc01-user-profile-snapshot.txt`
- `assets/playwright/tc01-user-profile-console.log`
- `assets/playwright/tc01-user-profile-network.txt`

### TC-ST-BATCH-20260409-USER-PROFILE-SESSION-02-01 修改密码页校验与重新登录

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-PROFILE-SESSION-02-01
related_subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-02
primary_req_id: USER-PROFILE-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:47
evidence:
  - assets/playwright/tc02-password-reset-login.png
  - assets/playwright/tc02-password-reset-url.txt
  - assets/playwright/tc02-password-reset-snapshot.txt
  - assets/playwright/tc02-password-reset-console.log
  - assets/playwright/tc02-password-reset-network.txt
```

Preconditions:

- 已写入用户端登录态
- 浏览器可访问 `/user/profile/password`

Steps:

1. 打开 `/user/profile/password`
2. 提交错误密码组合，验证校验分支
3. 提交正确密码组合并观察跳转

Expected:

- 校验失败时有明确提示
- 提交成功后清空本地登录态并跳回登录页
- 登录页显示“密码已修改，请重新登录”提示

Actual:

- 错误组合提交后页面未跳转
- 正确组合提交后跳转到 `http://localhost:9000/user/login?reason=password_reset`
- 控制台无 error/warning，网络日志无 4xx/5xx

Evidence:

- `assets/playwright/tc02-password-reset-login.png`
- `assets/playwright/tc02-password-reset-url.txt`
- `assets/playwright/tc02-password-reset-snapshot.txt`
- `assets/playwright/tc02-password-reset-console.log`
- `assets/playwright/tc02-password-reset-network.txt`

### TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-01 主动退出提示承接

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-01
related_subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-03
primary_req_id: USER-AUTH-003
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:48
evidence:
  - assets/playwright/tc03-logout-login.png
  - assets/playwright/tc03-logout-url.txt
  - assets/playwright/tc03-logout-snapshot.txt
  - assets/playwright/tc03-logout-console.log
  - assets/playwright/tc03-logout-network.txt
```

Preconditions:

- 已写入用户端登录态
- 浏览器可访问 `/user/profile`

Steps:

1. 打开 `/user/profile`
2. 点击“主动退出登录”
3. 检查登录页提示承接

Expected:

- 登录态被清空
- 页面跳转到 `/user/login`
- 登录页显示统一会话提示文案

Actual:

- 点击退出后跳转到 `http://localhost:9000/user/login?reason=logout`
- 登录页显示退出提示承接
- 控制台无 error/warning，网络日志无 4xx/5xx

Evidence:

- `assets/playwright/tc03-logout-login.png`
- `assets/playwright/tc03-logout-url.txt`
- `assets/playwright/tc03-logout-snapshot.txt`
- `assets/playwright/tc03-logout-console.log`
- `assets/playwright/tc03-logout-network.txt`

### TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-02 被挤下线提示承接

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-02
related_subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-03
primary_req_id: USER-AUTH-003
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:49
evidence:
  - assets/playwright/tc04-kickout-page.png
  - assets/playwright/tc04-kickout-login.png
  - assets/playwright/tc04-kickout-url.txt
  - assets/playwright/tc04-kickout-snapshot.txt
  - assets/playwright/tc04-kickout-console.log
  - assets/playwright/tc04-kickout-network.txt
```

Preconditions:

- 浏览器可访问 `/user/session-kickout?reason=kicked`

Steps:

1. 打开 `/user/session-kickout?reason=kicked`
2. 检查统一会话失效文案
3. 验证最终登录页承接

Expected:

- 页面展示会话失效提示与承接说明
- 最终跳转到登录页并显示统一提示文案
- 流程不依赖真实后端接口

Actual:

- 会话失效页可正常打开并展示统一文案
- 页面承接到 `http://localhost:9000/user/login?reason=kicked`
- 控制台无 error/warning，网络日志无 4xx/5xx

Evidence:

- `assets/playwright/tc04-kickout-page.png`
- `assets/playwright/tc04-kickout-login.png`
- `assets/playwright/tc04-kickout-url.txt`
- `assets/playwright/tc04-kickout-snapshot.txt`
- `assets/playwright/tc04-kickout-console.log`
- `assets/playwright/tc04-kickout-network.txt`
