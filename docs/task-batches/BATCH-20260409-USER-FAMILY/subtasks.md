# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-USER-FAMILY-01 | 家庭列表与创建加入退出入口 Make 化 | USER-FAMILY-001 | ["USER-FAMILY-002","USER-AUTH-002"] | 已完成 | Codex-Worker-D | ["USER-AUTH-002"] | ["TC-ST-BATCH-20260409-USER-FAMILY-01-01","TC-ST-BATCH-20260409-USER-FAMILY-01-02"] | ["BUG-BATCH-20260409-USER-FAMILY-001"] | 2026-04-09 21:53 |
| ST-BATCH-20260409-USER-FAMILY-02 | 家庭详情统计与成员占比 Make 化 | USER-FAMILY-002 | ["USER-FAMILY-001","USER-LEDGER-001"] | 已完成 | Codex-Worker-D | ["USER-FAMILY-001","USER-LEDGER-001"] | ["TC-ST-BATCH-20260409-USER-FAMILY-02-01","TC-ST-BATCH-20260409-USER-FAMILY-02-02"] | [] | 2026-04-09 21:53 |

## 子任务说明模板

### ST-BATCH-20260409-USER-FAMILY-01 家庭列表与创建加入退出入口 Make 化

```yaml
subtask_id: ST-BATCH-20260409-USER-FAMILY-01
title: 家庭列表与创建加入退出入口 Make 化
primary_req_id: USER-FAMILY-001
related_req_ids: ["USER-FAMILY-002","USER-AUTH-002"]
status: 已完成
owner_ai: Codex-Worker-D
dependencies: ["USER-AUTH-002"]
implementation_notes:
  - 以 FamilyList 与 FamilyEmpty 的信息结构为参考，重做家庭列表页
  - 支持创建家庭、按家庭 ID 加入、按邀请链接加入、主动退出的原型交互
  - 数据全部落在家庭模块独立 mock，不依赖后端接口
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-FAMILY-01-01","TC-ST-BATCH-20260409-USER-FAMILY-01-02"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-FAMILY-001"]
last_updated_at: 2026-04-09 21:53
```

Implementation Notes:

- 该子任务重点是家庭列表可用性和入口完整性，保证用户能进入下一层详情页。

### ST-BATCH-20260409-USER-FAMILY-02 家庭详情统计与成员占比 Make 化

```yaml
subtask_id: ST-BATCH-20260409-USER-FAMILY-02
title: 家庭详情统计与成员占比 Make 化
primary_req_id: USER-FAMILY-002
related_req_ids: ["USER-FAMILY-001","USER-LEDGER-001"]
status: 已完成
owner_ai: Codex-Worker-D
dependencies: ["USER-FAMILY-001","USER-LEDGER-001"]
implementation_notes:
  - 以 FamilyDetail 的信息层级为参考，实现家庭月度与年度汇总
  - 提供月收入、月支出、年收入、年支出的点击占比切换能力
  - 通过独立组件渲染成员占比饼图与排名列表
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-FAMILY-02-01","TC-ST-BATCH-20260409-USER-FAMILY-02-02"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:53
```

Implementation Notes:

- 该子任务重点是详情页的数据表达和占比交互，不展开成员流水明细。
