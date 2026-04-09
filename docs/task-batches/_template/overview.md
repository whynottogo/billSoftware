# 批次概览

```yaml
batch_id: BATCH-YYYYMMDD-TOPIC
title: 批次标题
source_task_pool: docs/requirements/bill-software-requirements.md
related_req_ids:
  - REQ-ID-001
status: 未开始
owner_ai: 待分配
created_at: YYYY-MM-DD HH:MM
last_updated_at: YYYY-MM-DD HH:MM
goal: 简要说明本批次目标
scope: frontend | backend | both
acceptance_summary:
  - 验收标准摘要 1
  - 验收标准摘要 2
```

## 需求记录

- 需求来源：
- 业务背景：
- 本批次覆盖范围：
- 不在本批次范围内的内容：

## 状态历史

| 时间 | 对象 | 状态变更 | 说明 |
| --- | --- | --- | --- |
| YYYY-MM-DD HH:MM | 批次 | 未开始 -> 开发中 | 领取任务并开始执行 |

## 完成标准

- [ ] 所有关联 `req_id` 的子任务已完成
- [ ] 相关测试用例全部通过
- [ ] 未关闭 Bug 已清零
- [ ] 需求文档中的状态已同步更新
