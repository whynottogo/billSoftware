# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 3 |
| 已挂起 | 0 |

### BUG-BATCH-20260410-ASSET-PROFILE-REAL-003 资料保存后昵称未更新且邮箱串值

```yaml
bug_id: BUG-BATCH-20260410-ASSET-PROFILE-REAL-003
primary_req_id: USER-PROFILE-001
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-04
related_test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-01
status: 已关闭
severity: medium
summary: 个人信息页初次自动化回归出现昵称/邮箱串值，后续确认为 Playwright wrapper `fill` 定位偏差导致的误报
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
steps:
  - 登录测试账号 `e2e0410000402`
  - 打开 `http://localhost:9000/user/profile`
  - 将昵称改为 `E2E验证0402-资料页保存`
  - 将邮箱改为 `e2e0410000402+ui@example.com`
  - 点击 `保存资料`
expected:
  - 昵称与邮箱分别按输入值独立保存，并在刷新后保持一致
actual:
  - 初次使用基于角色名的 `fill` 回归时，页面出现昵称/邮箱错位，并把异常值写入测试数据。
  - 随后改用 placeholder 与文件输入精确选择器重跑资料保存、头像上传和接口回读，昵称 `E2E验证0402-资料页二次验证`、邮箱 `e2e0410000402+verify@example.com` 与头像都保存成功，问题未再复现。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-save-anomaly-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-save-anomaly-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-save-anomaly-screenshot.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-save-anomaly-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-after-profile-save-anomaly.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-avatar-verify.json
suspected_root_cause: Playwright wrapper 通过角色名执行 `fill` 时发生定位偏差，误把两次输入写到了错误字段；产品真实保存逻辑本身正常
owner_ai: Codex
last_updated_at: 2026-04-10 16:13
```

### BUG-BATCH-20260410-ASSET-PROFILE-REAL-002 Playwright wrapper 会话落到 about:blank

```yaml
bug_id: BUG-BATCH-20260410-ASSET-PROFILE-REAL-002
primary_req_id: SYS-DATA-001
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-05
related_test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-02
status: 已关闭
severity: medium
summary: 用户登录页一度出现白屏，根因是新增资产/资料页面使用 `async/await` 后缺少 `regeneratorRuntime` 运行时，导致 wrapper 快照落空
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
steps:
  - 使用 wrapper 打开 `http://localhost:9000/user/login`
  - 立即抓取 snapshot、console、network 和 screenshot
expected:
  - 页面正常停留在用户登录页并可继续执行真实页面回归
actual:
  - 破坏性版本下登录页快照为空，控制台记录 `ReferenceError: regeneratorRuntime is not defined`，页面无法继续真实回归。
  - 将新增用户资产/资料页面中的 `async/await` 改回 Promise 链后，`npm run build` 通过，登录页、资产页、资料页与资产详情页均恢复正常渲染。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-snapshot.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-console-raw.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-screenshot.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-snapshot-fixed.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-console-fixed.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-page.txt
suspected_root_cause: 新增资产/资料页面中的 `async/await` 在当前 webpack/babel 运行时缺少 `regeneratorRuntime` 支持，导致登录后首轮真实页面回归白屏
owner_ai: Codex
last_updated_at: 2026-04-10 16:13
```

### BUG-BATCH-20260410-ASSET-PROFILE-REAL-001 改密回归前置登录失败

```yaml
bug_id: BUG-BATCH-20260410-ASSET-PROFILE-REAL-001
primary_req_id: USER-PROFILE-002
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-02
related_test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-02
status: 已关闭
severity: medium
summary: 改密闭环回归前，测试账号使用预期原密码 `123456` 登录返回 401，导致无法继续验证“修改密码 -> 旧会话失效 -> 新密码登录”
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
steps:
  - 使用测试账号 `e2e0410000402`
  - 调用 `POST /api/user/auth/login`，请求体为原密码 `123456`
expected:
  - 登录成功并返回可用于改密回归的有效 token
actual:
  - 接口返回 `{\"code\":401,\"message\":\"password is incorrect\"}`
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-old-password.json
suspected_root_cause: 测试账号密码可能已在之前回归中被修改但未恢复，或当前测试基线对原密码的假设已失效
owner_ai: Codex
last_updated_at: 2026-04-10 14:21
```

## 说明

- 本批次发现资产页空白、资产写操作失败、资料保存失败、改密后旧会话未失效、控制台报错或关键网络请求失败时，必须先归档证据到 `assets/playwright/`，再在此登记。
- 当前批次 3 个 Bug 均已关闭：`BUG-001` 是测试数据基线问题，`BUG-002` 是前端运行时缺少 `regeneratorRuntime`，`BUG-003` 是 Playwright wrapper 定位偏差造成的误报。
