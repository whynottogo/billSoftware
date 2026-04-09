# 批次概览

```yaml
batch_id: BATCH-20260409-LEDGER-CORE-REAL
title: 用户数据权限与当月收支真实后端闭环
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - SYS-DATA-001
  - USER-LEDGER-001
  - USER-LEDGER-002
  - USER-LEDGER-003
  - USER-LEDGER-004
status: 已完成
owner_ai: Codex
created_at: 2026-04-09 23:30
last_updated_at: 2026-04-10 00:18
goal: 以用户级数据隔离为前提，补齐当月收支读取、新增收支和分类管理的真实接口，并让管理端禁用用户时立即失效旧会话。
scope: both
acceptance_summary:
  - 用户端 `ledger/categories` 五个真实接口全部可用
  - `UserHome` 已切到真实 `/api/user/ledger` 与 `/api/user/categories`
  - 自定义分类可新增删除，默认分类不可删除
  - 新增支出与新增收入可即时回显到月汇总和日分组
  - 已验证用户隔离：用户 1 月账本为空，用户 2 可见自己的 2 条记录
```

## 需求记录

- 需求来源：`docs/requirements/bill-software-requirements.md` 协同任务清单
- 本批次覆盖范围：`SYS-DATA-001`、`USER-LEDGER-001/002/003/004` 的后端实现与 `/user/ledger` 前端真实接线
- 不在本批次范围：账单、预算、图表、家庭共享、管理端用户账单详情
- 备注：`SYS-DATA-001` 在需求文档中保留为 `联调中`，因为预算/资产/图表等其他业务域尚未全部接入真实权限校验；本批次覆盖的 ledger/category 范围已通过

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| 2026-04-09 23:30 | 批次 | 未开始 -> 开发中 | 领取用户数据权限与当月收支真实闭环批次并开始执行 |
| 2026-04-09 23:40 | 批次 | 开发中 -> 测试中 | 后端接口与会话失效逻辑已实现，进入验证阶段 |
| 2026-04-09 23:42 | 批次 | 测试中 -> 已阻塞 | 本机缺少 Go 运行环境且 `localhost:8080` 未启动，接口联调验证阻塞 |
| 2026-04-10 00:00 | 批次 | 已阻塞 -> 测试中 | 已补齐 Go 运行环境并初始化数据库 schema，后端恢复可用 |
| 2026-04-10 00:14 | 批次 | 测试中 -> 已完成 | 分类管理、收支回显、用户隔离与旧会话失效验证全部通过 |

## 完成标准

- [x] 所有关联 `req_id` 的子任务已完成
- [x] 相关测试用例全部通过
- [x] 未关闭 Bug 已清零
- [x] 需求文档中的状态已同步更新
