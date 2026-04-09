# 批次概览

```yaml
batch_id: BATCH-20260409-USER-FAMILY
title: Make 原型并行替换：用户端家庭模块
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - USER-FAMILY-001
  - USER-FAMILY-002
status: 已完成
owner_ai: Codex-Worker-D
created_at: 2026-04-09 21:34
last_updated_at: 2026-04-09 21:57
goal: 以 Figma Make 的家庭页面原型为基准，完成家庭列表、创建/加入/退出入口与家庭详情统计页面替换
scope: frontend
acceptance_summary:
  - /user/families 展示家庭列表、创建家庭、按家庭 ID 加入、按邀请链接加入和退出入口
  - 新增 /user/families/:familyId 家庭详情页，支持月度与年度汇总切换
  - 家庭详情支持成员收入/支出占比展示，交互点击后切换占比饼图
  - 本批次仅使用家庭模块本地 mock 数据，不接入真实后端接口
```

## 需求记录

- 需求来源：`docs/requirements/bill-software-requirements.md` 协同任务清单
- 业务背景：当前家庭模块仍为占位页，需要按 Make 原型完成可浏览、可交互的前端原型实现
- 本批次覆盖范围：用户端家庭列表页、家庭详情页、家庭 mock 数据与占比图展示
- 不在本批次范围内的内容：家庭成员审批流、成员移除、后端真实家庭接口、跨账号真实邀请联调

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| 2026-04-09 21:34 | 批次 | 未开始 -> 开发中 | Worker D 领取家庭模块任务并开始执行 |
| 2026-04-09 21:34 | 需求 `USER-FAMILY-001` | 未开始 -> 开发中 | 开始实现家庭列表、创建/加入/退出入口 |
| 2026-04-09 21:34 | 需求 `USER-FAMILY-002` | 未开始 -> 开发中 | 开始实现家庭详情月/年统计与成员占比 |
| 2026-04-09 21:47 | 批次 | 开发中 -> 测试中 | 家庭列表页、详情页、mock 与路由已完成，开始执行 build 和 Playwright 回归 |
| 2026-04-09 21:46 | Bug `BUG-BATCH-20260409-USER-FAMILY-001` | 新建 -> 已关闭 | 记录并关闭 Playwright 选择器误报，不属于业务页面缺陷 |
| 2026-04-09 21:55 | Bug `BUG-BATCH-20260409-USER-FAMILY-002` | 新建 -> 已关闭 | 并行改动触发路由悬挂 import，复测已清除编译错误 |
| 2026-04-09 21:53 | 需求 `USER-FAMILY-001` | 开发中 -> 已完成 | 家庭列表、创建、加入、退出入口测试通过 |
| 2026-04-09 21:53 | 需求 `USER-FAMILY-002` | 开发中 -> 已完成 | 家庭详情月/年统计与成员占比测试通过 |
| 2026-04-09 21:53 | 批次 | 测试中 -> 已完成 | 所有测试用例通过且无未关闭 Bug，批次收口完成 |

## 完成标准

- [x] 所有关联 `req_id` 的子任务已完成
- [x] 相关测试用例全部通过
- [x] 未关闭 Bug 已清零
- [x] 需求文档中的状态已同步更新
