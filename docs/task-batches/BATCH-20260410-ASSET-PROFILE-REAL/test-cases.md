# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 8 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-01 资产接口结构与用户隔离验证

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-01
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-01
primary_req_id: USER-ASSET-001
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 请求资产总览接口
  - 新增账户并再次读取总览
  - 请求账户详情接口并写入余额操作
expected:
  - 总览接口返回 summary + categories
  - 详情接口返回 account + category + records
  - 所有数据仅属于当前用户
actual:
  - `GET /api/user/assets` 初次返回空的真实资产总览与分类结构，创建账户后再次读取返回 1 个流动资产账户，字段结构满足 `summary + categories`。
  - `GET /api/user/assets/1` 返回 `account + category + records`，随后余额操作后再次读取详情，记录列表与账户余额都只反映当前登录用户的数据。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-assets-empty.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-assets-create.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-assets-after-create.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-assets-detail-before-op.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-assets-operation.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-assets-detail-after-op.json
executor: Codex
last_executed_at: 2026-04-10 14:21
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-02 账户新增编辑与余额操作接口验证

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-02
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-01
primary_req_id: USER-ASSET-002
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 创建账户
  - 编辑账户
  - 提交调整/增加/减少余额操作
expected:
  - 三类写操作全部成功并即时回显
actual:
  - 真实页面编辑现有账户时，`PUT /api/user/assets/1` 返回 200，账户备注更新为 `资产联调测试账户-UI编辑` 并立即回显到总览与详情页。
  - 详情页继续执行 `POST /api/user/assets/1/operations` 增加余额 `15.00` 后，账户余额更新为 `635.50`，记录列表新增 `UI回归增加15`。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-page.txt
executor: Codex
last_executed_at: 2026-04-10 15:52
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-01 个人信息读取与更新接口验证

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-01
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-02
primary_req_id: USER-PROFILE-001
status: 已通过
preconditions:
  - 后端已启动
  - 已获取用户 token
steps:
  - 请求个人信息接口
  - 更新昵称、邮箱、头像字段
expected:
  - 账号与手机号只读返回
  - 昵称、邮箱、头像可更新并持久化
actual:
  - `GET /api/user/profile` 返回账号、手机号、昵称、邮箱和头像相关字段；账号与手机号只读字段与页面展示一致。
  - `PUT /api/user/profile` 可持久化昵称、邮箱与头像字段，后续接口回读与重新登录拿到的 profile 数据一致。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-get.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-update.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-verify-after-precise-input.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-avatar-verify.json
executor: Codex
last_executed_at: 2026-04-10 16:10
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-02 修改密码与会话失效接口验证

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-02
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-02
primary_req_id: USER-PROFILE-002
status: 已通过
preconditions:
  - 用户已登录
steps:
  - 输入错误原密码尝试修改
  - 输入正确原密码和新密码修改
  - 使用旧会话访问受保护接口
expected:
  - 错误原密码被拒绝
  - 修改成功后旧会话失效并要求重新登录
actual:
  - 使用错误 `current_password` 调用改密接口时返回 400，错误原密码不会修改真实密码。
  - 使用正确原密码改密成功后，旧 token 访问受保护接口返回 401，新密码登录成功；随后又将密码恢复到 `123456`，恢复后的旧密码可再次登录。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-password-wrong-current.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-password-change-to-new.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-after-password-change-old-token.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-new-password.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-password-restore-old.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-after-password-restore-old-token.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-restored-old-password.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-restored-after-ui-password.json
executor: Codex
last_executed_at: 2026-04-10 16:03
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-01 资产总览与新增编辑真实回归

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-01
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-03
primary_req_id: USER-ASSET-001
status: 已通过
preconditions:
  - 用户账号可登录
  - 资产接口已联通
steps:
  - 打开资产总览页
  - 新增账户并编辑账户
expected:
  - 总览卡片与分类列表展示真实数据
  - 新增/编辑账户成功且列表即时刷新
actual:
  - 资产总览页成功展示真实净资产 `635.50`、总资产/总负债和分类列表。
  - 账户编辑弹窗保存后，账户备注更新为 `资产联调测试账户-UI编辑`，列表在同页即时刷新，无需重新登录。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-screenshot.txt
executor: Codex
last_executed_at: 2026-04-10 15:54
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-02 账户详情与余额记录真实回归

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-02
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-03
primary_req_id: USER-ASSET-003
status: 已通过
preconditions:
  - 用户账号可登录
  - 账户详情接口已联通
steps:
  - 进入账户详情页
  - 执行调整、增加、减少至少一种余额操作
  - 验证记录列表按月展示
expected:
  - 余额操作写入成功
  - 详情统计和记录列表即时更新
actual:
  - 账户详情页成功展示真实余额 `635.50`、本月净变动和 3 条历史记录。
  - 余额操作选择“增加余额”后写入 `UI回归增加15`，详情页即时出现新记录并把余额更新到 `635.50`。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-screenshot.txt
executor: Codex
last_executed_at: 2026-04-10 15:52
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-01 个人资料真实回归

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-01
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-04
primary_req_id: USER-PROFILE-001
status: 已通过
preconditions:
  - 用户账号可登录
  - 个人信息接口已联通
steps:
  - 打开个人信息页
  - 修改昵称、邮箱和头像
  - 保存并刷新页面
expected:
  - 账号与手机号只读展示
  - 昵称、邮箱和头像更新后回显一致
actual:
  - 资料页使用精确 placeholder/selector 输入重跑后，昵称 `E2E验证0402-资料页二次验证`、邮箱 `e2e0410000402+verify@example.com` 与头像上传都可通过真实接口保存，并在页面刷新后回显一致。
  - 初次 `fill` 导致的昵称/邮箱串值未再次复现，后续已确认 `BUG-BATCH-20260410-ASSET-PROFILE-REAL-003` 属于 Playwright wrapper 定位偏差而非产品逻辑缺陷。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-screenshot.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-avatar-verify.json
executor: Codex
last_executed_at: 2026-04-10 16:10
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-02 修改密码真实回归

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-02
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-04
primary_req_id: USER-PROFILE-002
status: 已通过
preconditions:
  - 用户账号可登录
  - 改密接口已联通
steps:
  - 打开修改密码页
  - 完成原密码校验和密码修改
  - 验证返回登录页并重新登录
expected:
  - 表单校验符合需求
  - 修改成功后跳回登录页并能使用新密码登录
actual:
  - UI 先验证错误原密码会被拒绝，随后使用正确原密码把密码改为 `654321`，页面立即跳回 `/user/login?reason=password_reset`。
  - 改密后旧 token 访问 `/api/user/profile` 返回 401，随后可用新密码登录；最后已通过 UI 把密码恢复到 `123456` 并再次登录成功。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-after-password-reset-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-old-token-after-ui-password-change.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-restored-after-ui-password.json
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-password-wrong-current.json
executor: Codex
last_executed_at: 2026-04-10 15:42
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-01 资产与个人信息闭环回归

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-01
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-05
primary_req_id: SYS-DATA-001
status: 已通过
preconditions:
  - 本批次功能已联调完成
steps:
  - 登录用户端
  - 跑通资产总览、账户详情、个人信息、修改密码闭环
expected:
  - 全链路使用真实数据
  - 旧会话失效和重新登录承接正常
actual:
  - 已跑通 `登录 -> 资产总览编辑 -> 账户详情增加余额 -> 个人资料保存与头像上传 -> 修改密码 -> 旧会话 401 -> 新密码登录 -> 恢复旧密码 -> 旧密码重新登录` 全链路。
  - 闭环结束时测试账号仍可使用基线密码 `123456` 登录，适合继续复用到后续批次。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-after-password-reset-page.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-profile-old-token-after-ui-password-change.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-login-restored-after-ui-password.json
executor: Codex
last_executed_at: 2026-04-10 16:13
```

### TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-02 控制台与关键请求回归

```yaml
test_case_id: TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-02
related_subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-05
primary_req_id: SYS-DATA-001
status: 已通过
preconditions:
  - 本批次功能已联调完成
steps:
  - 在资产和个人信息相关页面抓取控制台与网络请求
  - 检查是否存在未预期报错、4xx 或 5xx
expected:
  - 关键请求全部成功
  - 无未预期控制台 error
actual:
  - 登录页修复后不再出现 `regeneratorRuntime` 白屏；资产总览、资产详情、个人资料页的关键请求均返回 200，相关页面未出现未预期 console error。
  - 唯一出现的 400 来自故意验证“错误原密码”的改密用例，属于预期结果，已单独在密码接口用例中记录。
evidence:
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-snapshot-fixed.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-login-console-fixed.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-assets-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-asset-detail-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-console.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-user-profile-avatar-verify-network.txt
  - docs/task-batches/BATCH-20260410-ASSET-PROFILE-REAL/assets/playwright/20260410-api-user-password-wrong-current.json
executor: Codex
last_executed_at: 2026-04-10 16:13
```
