# 批次概览

```yaml
batch_id: BATCH-20260409-USER-CHART
title: Make 原型第四批替换：用户端图表模块
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - USER-CHART-001
  - USER-CHART-002
status: 已完成
owner_ai: Codex
created_at: 2026-04-09 22:08
last_updated_at: 2026-04-09 22:21
goal: 以 Figma Make 的 ExpenseChart/IncomeChart 原型为参考，完成用户端支出图表和收入图表双路由落地
scope: frontend
acceptance_summary:
  - /user/charts/expense 提供月/年支出图表和支出排行榜
  - /user/charts/income 提供年收入图表和收入排行榜
  - 保持稳定 SVG 图表方案，不引入新的高风险图表依赖
  - 图表页可切换年份，页面交互和路由可达性通过 Playwright 验证
```

## 需求记录

- 需求来源：`docs/requirements/bill-software-requirements.md` 协同任务清单
- 业务背景：用户要求继续按 Figma Make 原型完成图表模块翻写，并保持与现有用户端设计语言一致
- 本批次覆盖范围：`/user/charts/expense`、`/user/charts/income`、图表 mock 数据、图表页切换与排行榜展示
- 不在本批次范围内的内容：真实统计接口联调、导出报表、图表筛选器高级条件

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| 2026-04-09 22:08 | 批次 | 未开始 -> 开发中 | Worker E 领取图表模块批次并开始执行 |
| 2026-04-09 22:08 | 需求 `USER-CHART-001` | 未开始 -> 开发中 | 开始实现支出图表页 |
| 2026-04-09 22:08 | 需求 `USER-CHART-002` | 未开始 -> 开发中 | 开始实现收入图表页 |
| 2026-04-09 22:21 | 需求 `USER-CHART-001` | 开发中 -> 已完成 | 支出图表页与年份切换、排行榜验证通过 |
| 2026-04-09 22:21 | 需求 `USER-CHART-002` | 开发中 -> 已完成 | 收入图表页与年份切换、排行榜验证通过 |
| 2026-04-09 22:21 | 批次 | 开发中 -> 已完成 | 图表双路由通过回归，批次收口完成 |

## 完成标准

- [x] 所有关联 `req_id` 的子任务已完成
- [x] 相关测试用例全部通过
- [x] 未关闭 Bug 已清零
- [x] 需求文档中的状态已同步更新
