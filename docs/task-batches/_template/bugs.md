# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 0 |
| 已挂起 | 0 |

## Bug 模板

### BUG-BATCH-001-01 示例 Bug

```yaml
bug_id: BUG-BATCH-001-01
primary_req_id: REQ-ID-001
related_subtask_id: ST-BATCH-001-01
related_test_case_id: TC-ST-BATCH-001-01-01
status: 新建
severity: medium
summary: 简要描述问题
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
owner_ai: 待分配
last_updated_at: YYYY-MM-DD HH:MM
```

Steps:

1. 触发步骤 1
2. 触发步骤 2

Expected:

- 预期行为

Actual:

- 实际行为

Evidence:

- `assets/playwright/example.png`
- 当前 URL：
- 控制台摘要：
- 关键请求：

Suspected Root Cause:

- 待分析

Retest Result:

- 待回归
