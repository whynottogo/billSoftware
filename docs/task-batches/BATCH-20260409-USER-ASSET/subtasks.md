# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-USER-ASSET-01 | 资产总览与账户分类列表 Make 化 | USER-ASSET-001 | ["USER-ASSET-002","USER-ASSET-003"] | 已完成 | Codex | ["USER-AUTH-002"] | ["TC-ST-BATCH-20260409-USER-ASSET-01-01"] | ["BUG-BATCH-20260409-USER-ASSET-001","BUG-BATCH-20260409-USER-ASSET-002","BUG-BATCH-20260409-USER-ASSET-004"] | 2026-04-09 22:01 |
| ST-BATCH-20260409-USER-ASSET-02 | 新增与编辑账户表单壳 Make 化 | USER-ASSET-002 | ["USER-ASSET-001"] | 已完成 | Codex | ["USER-ASSET-001"] | ["TC-ST-BATCH-20260409-USER-ASSET-02-01"] | ["BUG-BATCH-20260409-USER-ASSET-003"] | 2026-04-09 22:01 |
| ST-BATCH-20260409-USER-ASSET-03 | 账户详情与余额变动记录 Make 化 | USER-ASSET-003 | ["USER-ASSET-001","USER-ASSET-002"] | 已完成 | Codex | ["USER-ASSET-001","USER-ASSET-002"] | ["TC-ST-BATCH-20260409-USER-ASSET-03-01","TC-ST-BATCH-20260409-USER-ASSET-03-02"] | ["BUG-BATCH-20260409-USER-ASSET-005"] | 2026-04-09 22:38 |

## 子任务说明模板

### ST-BATCH-20260409-USER-ASSET-01 资产总览与账户分类列表 Make 化

```yaml
subtask_id: ST-BATCH-20260409-USER-ASSET-01
title: 资产总览与账户分类列表 Make 化
primary_req_id: USER-ASSET-001
related_req_ids: ["USER-ASSET-002","USER-ASSET-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-AUTH-002"]
implementation_notes:
  - 以 Figma Make 的资产总览页面结构为参考，重构资产首页信息密度和卡片层次
  - 展示净资产、资产、负债、资产变化提示以及分类分组账户列表
  - 每个分类展示分类总额，账户项展示余额并可进入详情
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-ASSET-01-01"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-ASSET-001","BUG-BATCH-20260409-USER-ASSET-002","BUG-BATCH-20260409-USER-ASSET-004"]
last_updated_at: 2026-04-09 22:01
```

Implementation Notes:

- 该子任务重点是完成首页可读性和分组结构，不引入真实接口。

### ST-BATCH-20260409-USER-ASSET-02 新增与编辑账户表单壳 Make 化

```yaml
subtask_id: ST-BATCH-20260409-USER-ASSET-02
title: 新增与编辑账户表单壳 Make 化
primary_req_id: USER-ASSET-002
related_req_ids: ["USER-ASSET-001"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-ASSET-001"]
implementation_notes:
  - 在资产页提供新增账户和编辑账户入口，使用表单弹窗实现交互壳
  - 表单字段包含名称、备注、余额，银行卡和信用卡强制卡号，虚拟账户强制微信/支付宝区分
  - 编辑时账户类型不可修改，保持需求定义一致
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-ASSET-02-01"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-ASSET-003"]
last_updated_at: 2026-04-09 22:01
```

Implementation Notes:

- 该子任务强调字段校验和类型差异化，不做真实保存接口。

### ST-BATCH-20260409-USER-ASSET-03 账户详情与余额变动记录 Make 化

```yaml
subtask_id: ST-BATCH-20260409-USER-ASSET-03
title: 账户详情与余额变动记录 Make 化
primary_req_id: USER-ASSET-003
related_req_ids: ["USER-ASSET-001","USER-ASSET-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-ASSET-001","USER-ASSET-002"]
implementation_notes:
  - 新增账户详情页，支持调整余额、增加余额、减少余额三类操作并生成记录
  - 余额记录支持按月筛选查看
  - 账户设置中修改余额时自动追加一条调整记录
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-ASSET-03-01","TC-ST-BATCH-20260409-USER-ASSET-03-02"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-ASSET-005"]
last_updated_at: 2026-04-09 22:38
```

Implementation Notes:

- 该子任务使用前端 mock 数据与本地状态模拟记录变更，便于后续联调替换。
