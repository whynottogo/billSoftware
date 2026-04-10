# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260410-USER-FAMILY-REAL-01 | 后端家庭核心接口与统计口径 | USER-FAMILY-001 | ["USER-FAMILY-002","USER-LEDGER-001"] | 已完成 | Codex | ["SYS-DATA-001","USER-AUTH-002","USER-LEDGER-001"] | ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-01","TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-02"] | [] | 2026-04-10 20:11 |
| ST-BATCH-20260410-USER-FAMILY-REAL-02 | 用户家庭列表真实接线 | USER-FAMILY-001 | ["USER-FAMILY-002","USER-AUTH-002"] | 已完成 | Codex | ["USER-FAMILY-001"] | ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-01","TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-02"] | ["BUG-BATCH-20260410-USER-FAMILY-REAL-003"] | 2026-04-10 20:11 |
| ST-BATCH-20260410-USER-FAMILY-REAL-03 | 用户家庭详情统计真实接线 | USER-FAMILY-002 | ["USER-FAMILY-001","USER-LEDGER-001"] | 已完成 | Codex | ["USER-FAMILY-001","USER-LEDGER-001"] | ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-03-01"] | ["BUG-BATCH-20260410-USER-FAMILY-REAL-004"] | 2026-04-10 20:11 |
| ST-BATCH-20260410-USER-FAMILY-REAL-04 | 家庭闭环联调与回归 | USER-FAMILY-001 | ["USER-FAMILY-002","USER-AUTH-002"] | 已完成 | Codex | ["USER-FAMILY-001","USER-FAMILY-002"] | ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01","TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-02"] | ["BUG-BATCH-20260410-USER-FAMILY-REAL-001","BUG-BATCH-20260410-USER-FAMILY-REAL-002","BUG-BATCH-20260410-USER-FAMILY-REAL-003","BUG-BATCH-20260410-USER-FAMILY-REAL-004","BUG-BATCH-20260410-USER-FAMILY-REAL-005"] | 2026-04-10 20:11 |

## 子任务说明

### ST-BATCH-20260410-USER-FAMILY-REAL-01 后端家庭核心接口与统计口径

```yaml
subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-01
title: 后端家庭核心接口与统计口径
primary_req_id: USER-FAMILY-001
related_req_ids: ["USER-FAMILY-002","USER-LEDGER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["SYS-DATA-001","USER-AUTH-002","USER-LEDGER-001"]
implementation_notes:
  - 这是唯一允许修改后端路由接线、家庭数据访问和统计 SQL 的子任务
  - 固定补齐接口：GET /api/user/families、POST /api/user/families、POST /api/user/families/join、POST /api/user/families/join-by-link、POST /api/user/families/:familyId/leave、GET /api/user/families/:familyId、GET /api/user/families/:familyId/member-share
  - 列表响应至少包含家庭名称、创建人、家庭 ID、邀请链接、成员信息与当前 UI 所需摘要字段
  - 详情响应需返回 monthOptions、yearOptions、members 和当前详情页可直接消费的汇总结构
linked_test_case_ids: ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-01","TC-ST-BATCH-20260410-USER-FAMILY-REAL-01-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 20:11
```

Implementation Notes:

- 后端必须保证当前登录用户只能访问自己所属家庭的数据，且家庭统计口径基于该家庭成员的真实流水聚合。

### ST-BATCH-20260410-USER-FAMILY-REAL-02 用户家庭列表真实接线

```yaml
subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-02
title: 用户家庭列表真实接线
primary_req_id: USER-FAMILY-001
related_req_ids: ["USER-FAMILY-002","USER-AUTH-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-FAMILY-001"]
implementation_notes:
  - 仅负责 UserFamilies.vue 和独立家庭 API 模块
  - 移除 userFamilyMock 依赖，改为真实请求列表、创建、按 ID 加入、按邀请链接加入和退出家庭
  - 保持现有路由结构和页面信息层级，不额外扩展入口
linked_test_case_ids: ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-01","TC-ST-BATCH-20260410-USER-FAMILY-REAL-02-02"]
linked_bug_ids: ["BUG-BATCH-20260410-USER-FAMILY-REAL-003"]
last_updated_at: 2026-04-10 20:11
```

Implementation Notes:

- 家庭详情跳转仍复用现有 `/user/families/:familyId`，复制邀请链接行为保留，错误提示统一走页面消息组件。

### ST-BATCH-20260410-USER-FAMILY-REAL-03 用户家庭详情统计真实接线

```yaml
subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-03
title: 用户家庭详情统计真实接线
primary_req_id: USER-FAMILY-002
related_req_ids: ["USER-FAMILY-001","USER-LEDGER-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-FAMILY-001","USER-LEDGER-001"]
implementation_notes:
  - 仅负责 UserFamilyDetail.vue 和家庭详情/成员占比 API 调用
  - 保持月收入、月支出、年收入、年支出四个入口的点击切换交互
  - 对缺失家庭、空数据、未加入家庭等情况提供稳定兜底
linked_test_case_ids: ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-03-01"]
linked_bug_ids: ["BUG-BATCH-20260410-USER-FAMILY-REAL-004"]
last_updated_at: 2026-04-10 20:11
```

Implementation Notes:

- 家庭详情页默认展示最近月份和最近年份；成员占比图按当前选中的 periodType 和 metricType 重新请求真实数据。

### ST-BATCH-20260410-USER-FAMILY-REAL-04 家庭闭环联调与回归

```yaml
subtask_id: ST-BATCH-20260410-USER-FAMILY-REAL-04
title: 家庭闭环联调与回归
primary_req_id: USER-FAMILY-001
related_req_ids: ["USER-FAMILY-002","USER-AUTH-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-FAMILY-001","USER-FAMILY-002"]
implementation_notes:
  - 主线程负责 Playwright 回归、证据归档、Bug 台账和需求状态收口
  - 固定回归闭环：登录 -> 家庭列表 -> 创建家庭 -> 进入详情 -> 切换占比 -> 退出家庭
  - 2026-04-10 19:16 双用户联调前置账号 `e2e0410000401` 登录失败，已登记 `BUG-BATCH-20260410-USER-FAMILY-REAL-001`
  - 2026-04-10 19:31 已改用新注册并启用的联调用户 `e2efamily0410` 恢复回归，测试脚本时序异常 `BUG-BATCH-20260410-USER-FAMILY-REAL-002` 已关闭
  - 2026-04-10 19:45 Playwright 打开页面时命中 `regeneratorRuntime is not defined`，已登记 `BUG-BATCH-20260410-USER-FAMILY-REAL-003`
  - 2026-04-10 20:11 已完成“登录 -> 列表 -> 详情 -> 四种占比切换 -> 成员退出 -> 详情访问回收 -> 创建人清理”全链路回归
linked_test_case_ids: ["TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-01","TC-ST-BATCH-20260410-USER-FAMILY-REAL-04-02"]
linked_bug_ids: ["BUG-BATCH-20260410-USER-FAMILY-REAL-001","BUG-BATCH-20260410-USER-FAMILY-REAL-002","BUG-BATCH-20260410-USER-FAMILY-REAL-003","BUG-BATCH-20260410-USER-FAMILY-REAL-004","BUG-BATCH-20260410-USER-FAMILY-REAL-005"]
last_updated_at: 2026-04-10 20:11
```

Implementation Notes:

- 只要构建、浏览器控制台、网络请求或页面数据出现异常，必须先在 `bugs.md` 登记，再继续排查。
