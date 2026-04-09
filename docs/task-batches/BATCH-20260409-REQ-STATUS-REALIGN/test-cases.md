# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 3 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-01-01 系统与认证需求状态审计

```yaml
test_case_id: TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-01-01
related_subtask_id: ST-BATCH-20260409-REQ-STATUS-REALIGN-01
primary_req_id: SYS-AUTH-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 23:18
evidence:
  - docs/requirements/bill-software-requirements.md
  - frontend/src/router/index.js
  - frontend/src/utils/request.js
  - frontend/src/api/userAuth.js
  - backend/internal/router/router.go
  - backend/internal/handler/user_auth_handler.go
  - backend/internal/middleware/user_auth.go
```

Preconditions:

- 可读取前后端源码与需求任务池
- 认证类 req 的 Acceptance Criteria 已明确

Steps:

1. 核对 `/user` 与 `/admin` 路由前缀、登录态存储 key 和守卫逻辑
2. 核对用户注册、用户登录、管理员登录 API 封装是否接入真实接口
3. 核对用户单会话逻辑是否在登录时失活旧会话，并在受保护接口校验中生效
4. 按验收项判断系统与登录注册相关 req 是否满足已完成条件

Expected:

- `SYS-ARCH-001`、`SYS-AUTH-001`、`SYS-AUTH-002`、`SYS-AUTH-003`、`USER-AUTH-001`、`USER-AUTH-002`、`USER-AUTH-003`、`ADMIN-AUTH-001` 保持或回写为符合真实实现的状态
- `SYS-DATA-001` 维持开发中，不因认证底座已存在而误标为已完成

Actual:

- 路由隔离、token 隔离、注册登录、管理员登录、用户单会话互踢链路均有真实实现支撑
- `SYS-AUTH-002` 从 `开发中` 校正为 `已完成`
- `SYS-DATA-001` 因用户业务数据权限范围未完整覆盖，维持 `开发中`

Evidence:

- `frontend/src/router/index.js`
- `frontend/src/utils/request.js`
- `backend/internal/handler/user_auth_handler.go`
- `backend/internal/middleware/user_auth.go`

### TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-02-01 用户端业务需求状态审计

```yaml
test_case_id: TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-02-01
related_subtask_id: ST-BATCH-20260409-REQ-STATUS-REALIGN-02
primary_req_id: USER-LEDGER-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 23:18
evidence:
  - docs/requirements/bill-software-requirements.md
  - docs/task-batches/BATCH-20260409-USER-LEDGER-ACTIONS/overview.md
  - docs/task-batches/BATCH-20260409-BILL-BUDGET/overview.md
  - docs/task-batches/BATCH-20260409-USER-ASSET/overview.md
  - docs/task-batches/BATCH-20260409-USER-CHART/overview.md
  - docs/task-batches/BATCH-20260409-USER-PROFILE-SESSION/overview.md
  - docs/task-batches/BATCH-20260409-USER-FAMILY/overview.md
  - frontend/src/utils/userLedgerMock.js
  - frontend/src/utils/userFinanceMock.js
  - frontend/src/utils/userAssetMock.js
  - frontend/src/utils/userChartMock.js
  - frontend/src/utils/userFamilyMock.js
  - frontend/src/utils/userProfileMock.js
```

Preconditions:

- 已确认用户端大部分业务页面来自 Figma/Make 翻写批次
- 已确认源码中存在独立 mock util

Steps:

1. 阅读各用户业务批次 `overview.md`，确认其范围是否明确为 frontend / mock / 原型级
2. 搜索用户业务页面对 `*Mock.js`、`localStorage` 和本地状态的依赖
3. 对照需求验收项，判断真实业务能力是否已具备
4. 将仅完成前端原型或本地 mock 的 req 统一回退到 `开发中`

Expected:

- 用户端收支、账单、预算、资产、图表、个人信息、家庭等业务 req 不再保持 `已完成`
- 纠偏后状态与“真实实现”而不是“原型页完成”一致

Actual:

- 上述模块均存在明显 mock 数据或本地状态驱动，且批次记录也明确排除了真实接口联调
- `USER-LEDGER-*`、`USER-BILL-*`、`USER-BUDGET-*`、`USER-ASSET-*`、`USER-CHART-*`、`USER-PROFILE-001/002`、`USER-FAMILY-*` 全部回写为 `开发中`

Evidence:

- `frontend/src/utils/userLedgerMock.js`
- `frontend/src/utils/userFinanceMock.js`
- `frontend/src/utils/userAssetMock.js`
- `frontend/src/utils/userChartMock.js`
- `frontend/src/utils/userFamilyMock.js`
- `frontend/src/utils/userProfileMock.js`

### TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-03-01 管理端模块状态审计

```yaml
test_case_id: TC-ST-BATCH-20260409-REQ-STATUS-REALIGN-03-01
related_subtask_id: ST-BATCH-20260409-REQ-STATUS-REALIGN-03
primary_req_id: ADMIN-USER-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 23:18
evidence:
  - docs/requirements/bill-software-requirements.md
  - docs/task-batches/BATCH-20260409-ADMIN-PORTAL/overview.md
  - docs/task-batches/BATCH-20260409-ADMIN-USER/overview.md
  - frontend/src/pages/admin/AdminUsers.vue
  - frontend/src/utils/adminUserMock.js
  - backend/internal/router/router.go
  - backend/internal/handler/user_handler.go
```

Preconditions:

- 已确认后台门户 req 的 Notes 明确允许 mock 或原型级页面
- 已确认管理端用户模块同时存在前端 mock 页面和后端真实接口

Steps:

1. 对照 `ADMIN-DASHBOARD-001`、`ADMIN-APPROVAL-001`、`ADMIN-FAMILY-001` 的 Notes，确认 mock 是否允许作为完成依据
2. 核对 `ADMIN-USER-001/002/003` 是否仍由前端 mock 驱动，后端接口覆盖到何种程度
3. 按“允许 mock / 已有真实接口但未接线 / 真实能力尚未具备”三类重新回写状态

Expected:

- 后台门户三条 req 保持 `已完成`
- 用户列表和用户启停因后端已有真实接口但前端未接线，回写为 `联调中`
- 用户账单详情因真实详情接口未具备，回写为 `开发中`

Actual:

- 后台门户三条 req 维持 `已完成`
- `ADMIN-USER-001`、`ADMIN-USER-002` 调整为 `联调中`
- `ADMIN-USER-003` 调整为 `开发中`

Evidence:

- `frontend/src/pages/admin/AdminUsers.vue`
- `frontend/src/utils/adminUserMock.js`
- `backend/internal/router/router.go`
- `backend/internal/handler/user_handler.go`
