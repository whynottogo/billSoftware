# 批次概览

```yaml
batch_id: BATCH-20260409-REQ-STATUS-REALIGN
title: 协同任务状态口径纠偏与真实完成度回写
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - SYS-AUTH-002
  - USER-LEDGER-001
  - USER-LEDGER-002
  - USER-LEDGER-003
  - USER-LEDGER-004
  - USER-BILL-001
  - USER-BILL-002
  - USER-BILL-003
  - USER-BUDGET-001
  - USER-BUDGET-002
  - USER-ASSET-001
  - USER-ASSET-002
  - USER-ASSET-003
  - USER-CHART-001
  - USER-CHART-002
  - USER-PROFILE-001
  - USER-PROFILE-002
  - USER-FAMILY-001
  - USER-FAMILY-002
  - ADMIN-USER-001
  - ADMIN-USER-002
  - ADMIN-USER-003
status: 已完成
owner_ai: Codex
created_at: 2026-04-09 23:00
last_updated_at: 2026-04-09 23:20
goal: 按需求验收项和真实实现覆盖重审协同任务状态，纠正“前端原型完成即回写已完成”的失真状态。
scope: both
acceptance_summary:
  - 为需求状态补充统一口径：仅当前实现满足验收项，或 Notes 明确允许 mock 时，才标记为已完成
  - 将仅完成 Figma/Make 到 Vue 表现层、但未满足真实业务验收项的需求回退到开发中
  - 保留已满足真实认证能力或 Notes 明确允许 mock 的需求为已完成，并区分用户管理联调状态
```

## 需求记录

- 需求来源：`docs/requirements/bill-software-requirements.md` 协同任务清单
- 业务背景：当前任务池中大量需求已被回写为 `已完成`，但批次记录与源码显示其中一部分仅完成前端原型翻写或本地 mock 演示，导致状态口径失真
- 本批次覆盖范围：系统/认证、用户端业务模块、管理端用户模块与后台门户相关 req 的真实完成度审计；模块总览状态修正；状态口径补充说明
- 不在本批次范围内的内容：新增业务功能、补做缺失后端接口、补写 Playwright 页面证据、变更现有前后端业务逻辑

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| 2026-04-09 23:00 | 批次 | 未开始 -> 开发中 | 创建状态纠偏批次，开始逐条核对 req 验收项与真实实现 |
| 2026-04-09 23:08 | 需求 `SYS-AUTH-002` | 开发中 -> 已完成 | 核对用户单会话互踢链路，确认后端失活旧会话、受保护接口返回失效、前端收到 401 后跳转并提示已具备 |
| 2026-04-09 23:12 | 多个用户业务需求 | 已完成 -> 开发中 | 核对到仅完成前端原型或本地 mock，未满足真实业务验收项，统一回退状态 |
| 2026-04-09 23:14 | 需求 `ADMIN-USER-001` | 已完成 -> 联调中 | 后端真实用户列表接口已存在，但前端用户列表仍使用 mock 数据，调整为联调中 |
| 2026-04-09 23:14 | 需求 `ADMIN-USER-002` | 已完成 -> 联调中 | 后端启停接口已存在且登录态校验生效，但前端启停动作仍未接线，调整为联调中 |
| 2026-04-09 23:15 | 需求 `ADMIN-USER-003` | 已完成 -> 开发中 | 真实月/年账单详情接口尚未补齐，页面仍为 mock 详情，回退为开发中 |
| 2026-04-09 23:18 | 批次 | 开发中 -> 已完成 | 需求文档、模块状态和状态口径说明已同步更新，审计用例通过 |
| 2026-04-09 23:20 | 批次 | 已完成 -> 已完成 | 补充后续开发优先级建议，便于新批次直接按顺序领取 |

## 完成标准

- [x] 所有关联 `req_id` 的子任务已完成
- [x] 相关测试用例全部通过
- [x] 未关闭 Bug 已清零
- [x] 需求文档中的状态已同步更新

## 推荐优先级

### P0 先打通最小可用闭环

- `ADMIN-USER-001`、`ADMIN-USER-002`
  这两条后端接口已经存在，前端仍在用 mock；优先接到真实接口后，才能形成“用户注册 -> 管理员启用/禁用 -> 用户登录受影响”的真实闭环。
- `SYS-DATA-001`
  这是用户侧真实业务接口的安全底座，后续收支、资产、资料、家庭都要依赖它做用户隔离和管理员访问控制。
- `USER-LEDGER-001`、`USER-LEDGER-002`、`USER-LEDGER-003`、`USER-LEDGER-004`
  当月收支是用户端核心模块，也是账单、预算、图表、家庭统计的数据来源。建议把“列表读取 + 新增收支 + 分类管理”作为一个连续批次打通。

建议顺序：

1. 接通 `ADMIN-USER-001`、`ADMIN-USER-002`
2. 完成 `SYS-DATA-001`
3. 打通 `USER-LEDGER-001`、`USER-LEDGER-002`、`USER-LEDGER-003`、`USER-LEDGER-004`

### P1 基于真实收支数据补齐派生能力

- `USER-BILL-001`、`USER-BILL-002`、`USER-BILL-003`
  账单直接依赖收支汇总，且 `ADMIN-USER-003` 还要复用这套统计口径。
- `ADMIN-USER-003`
  应放在账单能力之后，一次复用用户侧月/年账单统计，不建议单独先做管理端 mock 替换。
- `USER-BUDGET-001`、`USER-BUDGET-002`
  预算依赖真实分类和收支执行数据，放在收支和账单之后性价比更高。
- `USER-CHART-001`、`USER-CHART-002`
  图表本质上是收支聚合视图，等真实收支和分类跑通后补最稳。

建议顺序：

1. `USER-BILL-001`、`USER-BILL-002`、`USER-BILL-003`
2. `ADMIN-USER-003`
3. `USER-BUDGET-001`、`USER-BUDGET-002`
4. `USER-CHART-001`、`USER-CHART-002`

### P2 补齐独立业务域

- `USER-PROFILE-001`、`USER-PROFILE-002`
  资料编辑和密码修改相对独立，不阻塞主记账闭环，可以在核心链路稳定后推进。
- `USER-ASSET-001`、`USER-ASSET-002`、`USER-ASSET-003`
  资产管家是独立数据域，且与账单是弱联动，适合在收支主线稳定后单独开批次完成。
- `USER-FAMILY-001`、`USER-FAMILY-002`
  家庭功能涉及共享范围、成员关系和汇总口径，是复杂度最高的一条，建议最后做。

建议顺序：

1. `USER-PROFILE-001`、`USER-PROFILE-002`
2. `USER-ASSET-001`、`USER-ASSET-002`、`USER-ASSET-003`
3. `USER-FAMILY-001`、`USER-FAMILY-002`

## 推荐批次拆法

- 批次 A：管理端真实用户联调
  覆盖 `ADMIN-USER-001`、`ADMIN-USER-002`
- 批次 B：用户数据权限与当月收支
  覆盖 `SYS-DATA-001`、`USER-LEDGER-001`、`USER-LEDGER-002`、`USER-LEDGER-003`、`USER-LEDGER-004`
- 批次 C：真实账单与管理端用户详情
  覆盖 `USER-BILL-001`、`USER-BILL-002`、`USER-BILL-003`、`ADMIN-USER-003`
- 批次 D：预算与图表
  覆盖 `USER-BUDGET-001`、`USER-BUDGET-002`、`USER-CHART-001`、`USER-CHART-002`
- 批次 E：个人信息
  覆盖 `USER-PROFILE-001`、`USER-PROFILE-002`
- 批次 F：资产管家
  覆盖 `USER-ASSET-001`、`USER-ASSET-002`、`USER-ASSET-003`
- 批次 G：家庭功能
  覆盖 `USER-FAMILY-001`、`USER-FAMILY-002`
