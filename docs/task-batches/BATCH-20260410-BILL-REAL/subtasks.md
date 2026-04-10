# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260410-BILL-REAL-01 | 后端账单统计共享内核与接口 | USER-BILL-001 | ["USER-BILL-002","USER-BILL-003","ADMIN-USER-003"] | 已完成 | Codex | ["USER-LEDGER-001"] | ["TC-ST-BATCH-20260410-BILL-REAL-01-01","TC-ST-BATCH-20260410-BILL-REAL-01-02"] | [] | 2026-04-10 02:46 |
| ST-BATCH-20260410-BILL-REAL-02 | 用户月账单列表与年账单真实接线 | USER-BILL-001 | ["USER-BILL-003"] | 已完成 | Codex | ["USER-BILL-001"] | ["TC-ST-BATCH-20260410-BILL-REAL-02-01","TC-ST-BATCH-20260410-BILL-REAL-02-02"] | [] | 2026-04-10 02:46 |
| ST-BATCH-20260410-BILL-REAL-03 | 用户月账单详情真实接线 | USER-BILL-002 | ["USER-BILL-001","USER-BILL-003"] | 已完成 | Codex | ["USER-BILL-001"] | ["TC-ST-BATCH-20260410-BILL-REAL-03-01"] | [] | 2026-04-10 02:46 |
| ST-BATCH-20260410-BILL-REAL-04 | 管理端用户账单详情真实接线 | ADMIN-USER-003 | ["USER-BILL-001","USER-BILL-002","USER-BILL-003"] | 已完成 | Codex | ["ADMIN-USER-001","USER-BILL-001","USER-BILL-002","USER-BILL-003"] | ["TC-ST-BATCH-20260410-BILL-REAL-04-01"] | [] | 2026-04-10 02:46 |
| ST-BATCH-20260410-BILL-REAL-05 | 账单闭环联调与回归 | USER-BILL-001 | ["USER-BILL-002","USER-BILL-003","ADMIN-USER-003"] | 已完成 | Codex | ["USER-BILL-001","USER-BILL-002","USER-BILL-003","ADMIN-USER-003"] | ["TC-ST-BATCH-20260410-BILL-REAL-05-01","TC-ST-BATCH-20260410-BILL-REAL-05-02"] | [] | 2026-04-10 02:46 |

## 子任务说明

### ST-BATCH-20260410-BILL-REAL-01 后端账单统计共享内核与接口

```yaml
subtask_id: ST-BATCH-20260410-BILL-REAL-01
title: 后端账单统计共享内核与接口
primary_req_id: USER-BILL-001
related_req_ids: ["USER-BILL-002","USER-BILL-003","ADMIN-USER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-LEDGER-001"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BILL-REAL-01-01","TC-ST-BATCH-20260410-BILL-REAL-01-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 02:46
```

Implementation Notes:

- 唯一允许修改后端共享统计内核、路由接线和账单聚合口径的子任务
- 需要补齐用户端 bill 接口与管理端单用户账单接口，确保用户端与管理端统计口径一致
- 统计必须按当前登录用户或管理员指定用户隔离，不可跨用户串数据
- 已交付 `GET /api/user/bills/years`、`GET /api/user/bills/year/:year`、`GET /api/user/bills/month/:month`、`GET /api/admin/users/:id/bills/overview`
- 已通过 curl 与 Playwright 双重验证用户 `e2e0410000402` 的 2025/2026 统计口径

### ST-BATCH-20260410-BILL-REAL-02 用户月账单列表与年账单真实接线

```yaml
subtask_id: ST-BATCH-20260410-BILL-REAL-02
title: 用户月账单列表与年账单真实接线
primary_req_id: USER-BILL-001
related_req_ids: ["USER-BILL-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BILL-001"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BILL-REAL-02-01","TC-ST-BATCH-20260410-BILL-REAL-02-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 02:46
```

Implementation Notes:

- 仅负责 `UserBillsMonth.vue`、`UserBillsYear.vue` 及其独立 API 模块
- 不改路由结构，不改月账单详情页和管理端详情页
- 目标是把年份切换、月份列表、年汇总与历年汇总全部改用真实接口
- 已修正降序年份数组下的左右箭头方向，左箭头切往更早年份，右箭头切回更新年份

### ST-BATCH-20260410-BILL-REAL-03 用户月账单详情真实接线

```yaml
subtask_id: ST-BATCH-20260410-BILL-REAL-03
title: 用户月账单详情真实接线
primary_req_id: USER-BILL-002
related_req_ids: ["USER-BILL-001","USER-BILL-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BILL-001"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BILL-REAL-03-01"]
linked_bug_ids: []
last_updated_at: 2026-04-10 02:46
```

Implementation Notes:

- 仅负责 `UserBillDetail.vue` 及其独立 API 模块
- 需要真实渲染月汇总、分类分布、支出排行、日趋势、近 6 个月对比和记账成就
- 必须复用后端统一统计口径，不允许前端再次手工聚合 ledger 列表
- 已修正错误态占位数据不再携带伪造的 `previousKey`，接口失败时“上一月”按钮保持禁用

### ST-BATCH-20260410-BILL-REAL-04 管理端用户账单详情真实接线

```yaml
subtask_id: ST-BATCH-20260410-BILL-REAL-04
title: 管理端用户账单详情真实接线
primary_req_id: ADMIN-USER-003
related_req_ids: ["USER-BILL-001","USER-BILL-002","USER-BILL-003"]
status: 已完成
owner_ai: Codex
dependencies: ["ADMIN-USER-001","USER-BILL-001","USER-BILL-002","USER-BILL-003"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BILL-REAL-04-01"]
linked_bug_ids: []
last_updated_at: 2026-04-10 02:46
```

Implementation Notes:

- 仅负责 `AdminUserDetail.vue` 及其独立 API 模块
- 管理端页面只读，不开放新增或编辑账单
- 顶部汇总与用户端对应月/年统计必须使用同一口径
- 已兼容后端返回的 `YYYY-MM-DD HH:mm` 注册时间格式，避免浏览器日期解析差异

### ST-BATCH-20260410-BILL-REAL-05 账单闭环联调与回归

```yaml
subtask_id: ST-BATCH-20260410-BILL-REAL-05
title: 账单闭环联调与回归
primary_req_id: USER-BILL-001
related_req_ids: ["USER-BILL-002","USER-BILL-003","ADMIN-USER-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-BILL-001","USER-BILL-002","USER-BILL-003","ADMIN-USER-003"]
linked_test_case_ids: ["TC-ST-BATCH-20260410-BILL-REAL-05-01","TC-ST-BATCH-20260410-BILL-REAL-05-02"]
linked_bug_ids: []
last_updated_at: 2026-04-10 02:46
```

Implementation Notes:

- 持续维护 Playwright 回归、接口证据与批次文档状态
- 用例重点覆盖用户端年/月账单页面、管理端用户详情页与口径一致性
- 发现异常必须先落证据，再登记 Bug
- 本轮 8 条用例全部通过，未触发业务 Bug 新增
