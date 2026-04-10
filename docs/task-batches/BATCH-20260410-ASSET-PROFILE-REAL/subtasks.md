# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260410-ASSET-PROFILE-REAL-01 | 后端资产域接口与数据隔离 | USER-ASSET-001 | ["USER-ASSET-002","USER-ASSET-003","SYS-DATA-001"] | 已完成 | Codex | ["SYS-DATA-001","USER-AUTH-002"] | ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-02"] | [] | 2026-04-10 16:13 |
| ST-BATCH-20260410-ASSET-PROFILE-REAL-02 | 后端个人信息与改密接口 | USER-PROFILE-001 | ["USER-PROFILE-002","SYS-DATA-001"] | 已完成 | Codex | ["SYS-DATA-001","USER-AUTH-002"] | ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-02"] | ["BUG-BATCH-20260410-ASSET-PROFILE-REAL-001"] | 2026-04-10 16:13 |
| ST-BATCH-20260410-ASSET-PROFILE-REAL-03 | 资产前端真实接线 | USER-ASSET-001 | ["USER-ASSET-002","USER-ASSET-003"] | 已完成 | Codex | ["USER-ASSET-001","USER-ASSET-002","USER-ASSET-003"] | ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-02"] | [] | 2026-04-10 16:13 |
| ST-BATCH-20260410-ASSET-PROFILE-REAL-04 | 个人信息前端真实接线 | USER-PROFILE-001 | ["USER-PROFILE-002"] | 已完成 | Codex | ["USER-PROFILE-001","USER-PROFILE-002"] | ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-02"] | ["BUG-BATCH-20260410-ASSET-PROFILE-REAL-003"] | 2026-04-10 16:13 |
| ST-BATCH-20260410-ASSET-PROFILE-REAL-05 | 资产与个人信息闭环联调回归 | SYS-DATA-001 | ["USER-ASSET-001","USER-ASSET-002","USER-ASSET-003","USER-PROFILE-001","USER-PROFILE-002"] | 已完成 | Codex | ["SYS-DATA-001","USER-ASSET-001","USER-ASSET-002","USER-ASSET-003","USER-PROFILE-001","USER-PROFILE-002"] | ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-02"] | ["BUG-BATCH-20260410-ASSET-PROFILE-REAL-002","BUG-BATCH-20260410-ASSET-PROFILE-REAL-003"] | 2026-04-10 16:13 |

## 子任务说明

### ST-BATCH-20260410-ASSET-PROFILE-REAL-01 后端资产域接口与数据隔离

```yaml
subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-01
title: 后端资产域接口与数据隔离
primary_req_id: USER-ASSET-001
related_req_ids: ["USER-ASSET-002","USER-ASSET-003","SYS-DATA-001"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001","USER-AUTH-002"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-01-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 16:13
```

Implementation Notes:

- 这是唯一允许修改资产相关后端路由、模型和数据口径的子任务
- 首批接口契约固定：
  - `GET /api/user/assets`
  - `POST /api/user/assets`
  - `PUT /api/user/assets/:id`
  - `GET /api/user/assets/:id`
  - `POST /api/user/assets/:id/operations`
- 资产总览接口需返回 `summary + categories[]`，账户详情接口需返回 `account + category + records[]`
- 所有账户、余额变动记录必须按当前登录 `user_id` 隔离

### ST-BATCH-20260410-ASSET-PROFILE-REAL-02 后端个人信息与改密接口

```yaml
subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-02
title: 后端个人信息与改密接口
primary_req_id: USER-PROFILE-001
related_req_ids: ["USER-PROFILE-002","SYS-DATA-001"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001","USER-AUTH-002"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-02-02"]
linked_bug_ids: ["BUG-BATCH-20260410-ASSET-PROFILE-REAL-001"]
last_updated_at: 2026-04-10 16:13
```

Implementation Notes:

- 这是唯一允许修改个人信息后端路由、会话失效和密码校验逻辑的子任务
- 首批接口契约固定：
  - `GET /api/user/profile`
  - `PUT /api/user/profile`
  - `PUT /api/user/profile/password`
- 资料接口返回 `account / username / nickname / phone / email / avatar_original / avatar_compressed / updated_at`
- 修改密码成功后必须立即使该用户所有活跃 `user_sessions` 失效，并要求重新登录

### ST-BATCH-20260410-ASSET-PROFILE-REAL-03 资产前端真实接线

```yaml
subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-03
title: 资产前端真实接线
primary_req_id: USER-ASSET-001
related_req_ids: ["USER-ASSET-002","USER-ASSET-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-ASSET-001","USER-ASSET-002","USER-ASSET-003"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-03-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 16:13
```

Implementation Notes:

- 仅允许修改 `UserAssets.vue`、`UserAssetDetail.vue` 和新增资产 API 模块
- 不改路由结构，不继续依赖 `userAssetMock`
- 页面消费结构保持现有页面字段：资产总览 `summary + categories`，详情页 `account + categoryName + records`

### ST-BATCH-20260410-ASSET-PROFILE-REAL-04 个人信息前端真实接线

```yaml
subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-04
title: 个人信息前端真实接线
primary_req_id: USER-PROFILE-001
related_req_ids: ["USER-PROFILE-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-PROFILE-001","USER-PROFILE-002"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-04-02"]
linked_bug_ids: ["BUG-BATCH-20260410-ASSET-PROFILE-REAL-003"]
last_updated_at: 2026-04-10 16:13
```

Implementation Notes:

- 仅允许修改 `UserProfile.vue`、`UserProfilePassword.vue` 和新增个人信息 API 模块
- 账号和手机号保持只读；昵称、邮箱、头像改走真实接口
- 改密成功后必须回到登录页，并承接现有登录态清理与提示逻辑

### ST-BATCH-20260410-ASSET-PROFILE-REAL-05 资产与个人信息闭环联调回归

```yaml
subtask_id: ST-BATCH-20260410-ASSET-PROFILE-REAL-05
title: 资产与个人信息闭环联调回归
primary_req_id: SYS-DATA-001
related_req_ids: ["USER-ASSET-001","USER-ASSET-002","USER-ASSET-003","USER-PROFILE-001","USER-PROFILE-002"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001","USER-ASSET-001","USER-ASSET-002","USER-ASSET-003","USER-PROFILE-001","USER-PROFILE-002"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-01","TC-ST-BATCH-20260410-ASSET-PROFILE-REAL-05-02"]
linked_bug_ids: ["BUG-BATCH-20260410-ASSET-PROFILE-REAL-002","BUG-BATCH-20260410-ASSET-PROFILE-REAL-003"]
last_updated_at: 2026-04-10 16:13
```

Implementation Notes:

- 核心闭环是“登录 -> 查看资产总览 -> 新增账户 -> 查看账户详情并写入余额操作 -> 修改个人资料 -> 修改密码 -> 重新登录”
- Playwright 回归必须优先覆盖真实页面、关键请求和控制台
- 发现异常时必须先落证据到 `assets/playwright/`，再登记 `bugs.md`
