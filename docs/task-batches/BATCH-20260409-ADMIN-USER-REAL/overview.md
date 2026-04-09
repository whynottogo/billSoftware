# 批次概览

```yaml
batch_id: BATCH-20260409-ADMIN-USER-REAL
title: 管理端真实用户列表与启停联调
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - ADMIN-USER-001
  - ADMIN-USER-002
status: 已完成
owner_ai: Codex
created_at: 2026-04-09 23:30
last_updated_at: 2026-04-10 00:18
goal: 将管理端用户列表从 mock 切换为真实接口，并打通启用/禁用对用户登录态的真实影响。
scope: frontend
acceptance_summary:
  - `frontend` 执行 `npm run build` 成功（仅 webpack 体积告警）
  - 管理员真实登录后可进入 `/admin/users` 查看真实用户列表，且列表不暴露密码字段
  - 新注册用户默认显示为禁用，管理员可真实执行启用与禁用
  - 启用后用户可登录，禁用后旧会话触发失效提示，重新登录返回 `user is disabled`
```

## 需求记录

- 需求来源：`docs/requirements/bill-software-requirements.md` 协同任务清单
- 本批次覆盖范围：`ADMIN-USER-001`、`ADMIN-USER-002` 的前端接线、联调与回归
- 不在本批次范围：`ADMIN-USER-003` 详情真实化、账单聚合、路由结构改造

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| 2026-04-09 23:30 | 批次 | 未开始 -> 开发中 | 领取管理端真实用户联调批次并开始执行 |
| 2026-04-09 23:40 | 批次 | 开发中 -> 测试中 | 完成代码复核并进入验证阶段 |
| 2026-04-09 23:46 | 批次 | 测试中 -> 已阻塞 | 后端健康检查不可达，真实接口联调阻塞 |
| 2026-04-10 00:00 | 批次 | 已阻塞 -> 测试中 | 已补齐本地 Go 运行环境并恢复 `/api/health` |
| 2026-04-10 00:07 | 批次 | 测试中 -> 已完成 | 真实用户列表、启停联动与用户端登录影响验证全部通过 |

## 完成标准

- [x] 所有关联 `req_id` 的子任务已完成
- [x] 相关测试用例全部通过
- [x] 未关闭 Bug 已清零
- [x] 需求文档中的状态已同步更新
