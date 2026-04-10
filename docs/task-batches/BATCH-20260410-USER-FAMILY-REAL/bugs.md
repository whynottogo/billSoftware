# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 5 |
| 已挂起 | 0 |

### BUG-BATCH-20260410-USER-FAMILY-REAL-001 联调账号 `e2e0410000401` 使用预期密码登录失败，阻塞双用户家庭闭环

```yaml
bug_id: BUG-BATCH-20260410-USER-FAMILY-REAL-001
primary_req_id: USER-FAMILY-001
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-04
related_test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01
status: 已关闭
severity: medium
summary: family 双用户闭环预检时，联调账号 `e2e0410000401` 使用当前掌握的密码 `123456` 登录返回 401，阻塞“创建 -> 加入 -> 多成员统计 -> 退出”主链路
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
steps:
  - 启动前后端服务
  - 使用用户 `e2e0410000402` 和 `e2e0410000401` 分别调用 `POST /api/user/auth/login`
  - 用户 2 使用密码 `123456` 登录成功
  - 用户 1 使用同一预期密码登录
expected:
  - 双用户联调账号都能用批次约定密码登录，以便继续家庭创建/加入闭环
actual:
  - `e2e0410000402` 登录成功并返回 token
  - `e2e0410000401` 返回 `{\"code\":401,\"message\":\"password is incorrect\"}`
  - 之后通过真实注册新联调用户 `e2efamily0410`、管理员启用并顺序化登录，双用户闭环恢复
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-login-user2.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-login-user1.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-register-user3.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-enable-user3.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-login-user3-retry.json
suspected_root_cause: 历史批次中的测试账号密码与当前批次掌握的口令不一致，属于联调环境/测试数据失配
owner_ai: Codex
last_updated_at: 2026-04-10 20:11
```

### BUG-BATCH-20260410-USER-FAMILY-REAL-002 workaround 期间并发启用后立即登录，用户 3 命中旧禁用状态

```yaml
bug_id: BUG-BATCH-20260410-USER-FAMILY-REAL-002
primary_req_id: USER-FAMILY-001
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-04
related_test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01
status: 已关闭
severity: low
summary: 为绕开用户1密码失配而注册新联调用户时，脚本并发执行“启用用户3”和“用户3登录”，登录请求先命中旧禁用状态返回 403
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
steps:
  - 管理员登录成功
  - 注册新联调用户 `e2efamily0410`
  - 并发执行 `PUT /api/admin/users/3/status` 和 `POST /api/user/auth/login`
expected:
  - workaround 期间应在启用完成后再登录新用户，避免命中旧状态
actual:
  - 启用接口返回成功
  - 并发登录请求返回 `{\"code\":403,\"message\":\"user is disabled\"}`
  - 顺序化执行“启用 -> 登录”后，`20260410-family-api-login-user3-retry.json` 返回成功
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-enable-user3.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-login-user3.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-login-user3-retry.json
suspected_root_cause: 联调脚本并发顺序不当，属于测试时序问题，不是 family 功能本身实现错误
owner_ai: Codex
last_updated_at: 2026-04-10 20:11
```

### BUG-BATCH-20260410-USER-FAMILY-REAL-003 `UserFamilies.vue` 使用 `async/await` 触发 `regeneratorRuntime is not defined`，阻塞真实页面联调

```yaml
bug_id: BUG-BATCH-20260410-USER-FAMILY-REAL-003
primary_req_id: USER-FAMILY-001
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-02
related_test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-01
status: 已关闭
severity: high
summary: Playwright 打开 `http://localhost:9000/user/login` 时，前端 bundle 在 `frontend/src/pages/user/UserFamilies.vue` 抛出 `ReferenceError: regeneratorRuntime is not defined`，家庭页面真实联调被阻塞
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
steps:
  - 启动前后端服务
  - 用 Playwright 打开 `http://localhost:9000/user/login`
  - 观察初始控制台输出
expected:
  - 用户登录页可正常加载，并继续执行家庭页面回归
actual:
  - 初始回归中控制台报错 `ReferenceError: regeneratorRuntime is not defined`
  - 将 `UserFamilies.vue` 与 `UserFamilyDetail.vue` 的 `async/await` 改为 Promise 写法后，登录页与家庭列表页恢复正常加载
  - 修复后控制台无 error，页面可继续联调
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-ui-runtime-error.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-ui-runtime-error-console.log
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-list-ui-success-console.txt
suspected_root_cause: 当前前端运行环境未注入 `regeneratorRuntime`，页面脚本直接使用 `async/await` 导致浏览器运行时崩溃
owner_ai: Codex
last_updated_at: 2026-04-10 20:11
```

### BUG-BATCH-20260410-USER-FAMILY-REAL-004 家庭详情接口已返回成功，但页面仍落入失败态

```yaml
bug_id: BUG-BATCH-20260410-USER-FAMILY-REAL-004
primary_req_id: USER-FAMILY-002
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-03
related_test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-03-01
status: 已关闭
severity: high
summary: 详情接口 `GET /api/user/families/:familyId` 已返回 200，但 `UserFamilyDetail.vue` 页面仍显示“家庭详情加载失败，请稍后重试”，阻塞家庭详情统计真实联调
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
steps:
  - 使用用户 `e2e0410000402` 登录
  - 打开 `http://localhost:9000/user/families`
  - 点击家庭卡片“查看详情”
expected:
  - 页面展示家庭月度与年度汇总，以及成员占比面板
actual:
  - 初次回归时页面进入详情路由，但因 Promise 改写遗留的 `this` 上下文引用错误而落入失败态
  - 修复后，详情页成功展示家庭汇总，四种成员占比切换和网络请求全部恢复正常
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-error.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-error-network.txt
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success.png
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-detail-ui-success-network.txt
suspected_root_cause: `UserFamilyDetail.vue` Promise 改写后仍残留错误的 `this` 上下文引用，接口成功结果被前端自身异常吞没
owner_ai: Codex
last_updated_at: 2026-04-10 20:11
```

### BUG-BATCH-20260410-USER-FAMILY-REAL-005 清理测试数据时，创建人旧 token 调用退出接口返回 401

```yaml
bug_id: BUG-BATCH-20260410-USER-FAMILY-REAL-005
primary_req_id: USER-FAMILY-001
related_subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-04
related_test_case_id: TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01
status: 已关闭
severity: low
summary: 在清理本批次测试家庭时，用户 `e2e0410000402` 使用旧 token 调用 `POST /api/user/families/:familyId/leave` 返回 `401 user session is invalid`
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
steps:
  - 用户 3 正常退出家庭
  - 使用此前缓存的用户 2 token 调用 `POST /api/user/families/FAM-18A4FBF2C4ADE858-2/leave`
expected:
  - 创建人可退出并删除仅剩自己的测试家庭
actual:
  - 旧 token 首次调用时返回 `401 user session is invalid`
  - 改用浏览器当前有效 `bill_user_token` 后，创建人退出成功并删除测试家庭，列表接口返回 `familyCount: 0`
evidence:
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user2-leave-cleanup.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user2-leave-cleanup-browser-token.json
  - docs/task-batches/BATCH-20260410-USER-FAMILY-REAL/assets/playwright/20260410-family-api-user2-list-after-cleanup-browser-token.json
suspected_root_cause: 缓存的用户 2 token 已失效，属于测试会话时效问题，不是家庭退出接口本身的业务错误
owner_ai: Codex
last_updated_at: 2026-04-10 20:11
```

## 当前状态

- 当前批次已登记的 5 条缺陷均已关闭：
  - `BUG-BATCH-20260410-USER-FAMILY-REAL-001`：联调账号密码失配，已通过新注册联调用户绕开
  - `BUG-BATCH-20260410-USER-FAMILY-REAL-002`：测试脚本时序问题，顺序化执行后已恢复
  - `BUG-BATCH-20260410-USER-FAMILY-REAL-003`：family 页面脚本运行时依赖 `regeneratorRuntime`，已改为 Promise 写法并回归通过
  - `BUG-BATCH-20260410-USER-FAMILY-REAL-004`：family 详情页 Promise 上下文异常，已修复并回归通过
  - `BUG-BATCH-20260410-USER-FAMILY-REAL-005`：测试清理阶段使用了失效 token，已改用浏览器有效登录态完成收尾
