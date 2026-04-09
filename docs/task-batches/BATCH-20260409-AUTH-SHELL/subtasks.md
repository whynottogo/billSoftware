# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-AUTH-SHELL-01 | 用户端登录注册页 Make 化 | USER-AUTH-001 | ["USER-AUTH-002","SYS-AUTH-001"] | 已完成 | Codex | ["SYS-AUTH-001"] | ["TC-ST-BATCH-20260409-AUTH-SHELL-01-01","TC-ST-BATCH-20260409-AUTH-SHELL-01-02"] | [] | 2026-04-09 18:17 |
| ST-BATCH-20260409-AUTH-SHELL-02 | 用户端首页与用户壳替换 | USER-LEDGER-001 | ["SYS-ARCH-001","SYS-AUTH-001"] | 已完成 | Codex | ["USER-AUTH-002"] | ["TC-ST-BATCH-20260409-AUTH-SHELL-02-01"] | [] | 2026-04-09 18:17 |
| ST-BATCH-20260409-AUTH-SHELL-03 | 管理员登录页 Make 化 | ADMIN-AUTH-001 | ["SYS-AUTH-003","SYS-ARCH-001"] | 已完成 | Codex | ["SYS-AUTH-003"] | ["TC-ST-BATCH-20260409-AUTH-SHELL-03-01"] | [] | 2026-04-09 18:17 |
| ST-BATCH-20260409-AUTH-SHELL-04 | 管理端用户列表与后台壳替换 | ADMIN-USER-001 | ["SYS-ARCH-001","SYS-AUTH-003"] | 已完成 | Codex | ["ADMIN-AUTH-001"] | ["TC-ST-BATCH-20260409-AUTH-SHELL-04-01"] | [] | 2026-04-09 18:17 |

## 子任务说明模板

### ST-BATCH-20260409-AUTH-SHELL-01 用户端登录注册页 Make 化

```yaml
subtask_id: ST-BATCH-20260409-AUTH-SHELL-01
title: 用户端登录注册页 Make 化
primary_req_id: USER-AUTH-001
related_req_ids: ["USER-AUTH-002","SYS-AUTH-001"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-AUTH-001"]
implementation_notes:
  - 以 Figma Make 中 UserLogin 与 UserRegister 为主要视觉和信息层级参考
  - 保留当前接口调用与 Vue 路由方式，不直接迁移 React 代码
  - 确保“用户名或手机号登录”与“注册后待管理员启用”文案清晰可见
linked_test_case_ids: ["TC-ST-BATCH-20260409-AUTH-SHELL-01-01","TC-ST-BATCH-20260409-AUTH-SHELL-01-02"]
linked_bug_ids: []
last_updated_at: 2026-04-09 18:17
```

Implementation Notes:

- 需要同时处理用户登录入口、用户注册入口、认证成功后的路由跳转和双端入口切换提示。

### ST-BATCH-20260409-AUTH-SHELL-02 用户端首页与用户壳替换

```yaml
subtask_id: ST-BATCH-20260409-AUTH-SHELL-02
title: 用户端首页与用户壳替换
primary_req_id: USER-LEDGER-001
related_req_ids: ["SYS-ARCH-001","SYS-AUTH-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-AUTH-002"]
implementation_notes:
  - 以 Figma Make 中 MonthlyLedger 与 UserNavigation 为参考重做用户端布局
  - 保留现有 /user 路由结构与 token 存储 key
  - 当前批次先替换首页骨架与导航观感，不覆盖记账弹窗与分类管理的完整业务
linked_test_case_ids: ["TC-ST-BATCH-20260409-AUTH-SHELL-02-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 18:17
```

Implementation Notes:

- 该子任务侧重 Make 原型的桌面端外壳和信息层级，不在本批次引入完整记账录入流程。

### ST-BATCH-20260409-AUTH-SHELL-03 管理员登录页 Make 化

```yaml
subtask_id: ST-BATCH-20260409-AUTH-SHELL-03
title: 管理员登录页 Make 化
primary_req_id: ADMIN-AUTH-001
related_req_ids: ["SYS-AUTH-003","SYS-ARCH-001"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-AUTH-003"]
implementation_notes:
  - 以 Figma Make 中 AdminLogin 为视觉基线重做管理员登录页
  - 管理员登录态继续使用 bill_admin_token
  - 明确管理员与普通用户入口隔离，不混入用户注册能力
linked_test_case_ids: ["TC-ST-BATCH-20260409-AUTH-SHELL-03-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 18:17
```

Implementation Notes:

- 管理端页面需要沿用品牌一致性，同时比用户端更克制、信息更后台化。

### ST-BATCH-20260409-AUTH-SHELL-04 管理端用户列表与后台壳替换

```yaml
subtask_id: ST-BATCH-20260409-AUTH-SHELL-04
title: 管理端用户列表与后台壳替换
primary_req_id: ADMIN-USER-001
related_req_ids: ["SYS-ARCH-001","SYS-AUTH-003"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-AUTH-001"]
implementation_notes:
  - 以 Figma Make 中 UserManagement 为主要参考重做后台用户列表骨架
  - 保留现有 /admin/users 路由
  - 当前批次先完成信息架构、摘要卡、筛选位和列表外观，不强行补齐后端联调细节
linked_test_case_ids: ["TC-ST-BATCH-20260409-AUTH-SHELL-04-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 18:17
```

Implementation Notes:

- 该页面是后续用户启用禁用、详情查看的后台承接页，先完成 Make 风格落地和关键字段分组。
