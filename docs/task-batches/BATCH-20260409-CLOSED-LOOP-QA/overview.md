# 批次概览

```yaml
batch_id: BATCH-20260409-CLOSED-LOOP-QA
title: 注册启停登录闭环回归基线
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - USER-AUTH-001
  - USER-AUTH-002
  - USER-AUTH-003
  - SYS-AUTH-002
  - ADMIN-USER-001
  - ADMIN-USER-002
status: 已完成
owner_ai: Codex
created_at: 2026-04-09 23:30
last_updated_at: 2026-04-10 00:18
goal: 建立真实闭环回归基线，持续验证注册、启用、登录、禁用和会话失效提示不被多线开发破坏。
scope: both
acceptance_summary:
  - 注册新用户成功，页面提示“等待管理员启用”
  - 管理端真实列表可见新用户且默认禁用
  - 启用后用户可登录并进入 `/user/ledger`
  - 禁用后重新登录返回 `user is disabled`
  - 旧会话访问用户接口时触发 401 和统一提示文案
status_history:
  - 2026-04-09 23:30 批次从未开始进入开发中
  - 2026-04-09 23:47 批次进入测试中并开始执行闭环回归
  - 2026-04-09 23:49 注册提交阶段触发 504，批次转为已阻塞
  - 2026-04-10 00:02 注册提交阶段触发新的 500（`check username failed`），批次继续保持已阻塞
  - 2026-04-10 00:07 起 修复环境与数据库问题后，五条闭环用例全部通过
```

## 需求记录

- 需求来源：`docs/requirements/bill-software-requirements.md` 协同任务清单
- 本批次覆盖范围：闭环 Playwright 回归、证据归档、需求状态收口辅助
- 不在本批次范围：业务功能开发和页面结构改造

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| 2026-04-09 23:30 | 批次 | 未开始 -> 开发中 | 创建闭环 QA 批次并开始维护真实回归基线 |
| 2026-04-09 23:47 | 批次 | 开发中 -> 测试中 | 执行闭环回归 |
| 2026-04-09 23:49 | 批次 | 测试中 -> 已阻塞 | 注册接口返回 504，闭环阻塞 |
| 2026-04-10 00:02 | 批次 | 已阻塞 -> 已阻塞 | 注册接口改为 500，提示 `check username failed` |
| 2026-04-10 00:14 | 批次 | 已阻塞 -> 已完成 | 环境与数据库问题修复后，闭环 5 条用例全部通过 |

## 完成标准

- [x] 所有关联 `req_id` 的子任务已完成
- [x] 相关测试用例全部通过
- [x] 未关闭 Bug 已清零
- [x] 需求文档中的状态已同步更新
