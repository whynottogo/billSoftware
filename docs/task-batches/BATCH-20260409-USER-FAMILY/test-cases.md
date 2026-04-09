# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 4 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-USER-FAMILY-01-01 家庭列表页展示与详情跳转

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-FAMILY-01-01
related_subtask_id: ST-BATCH-20260409-USER-FAMILY-01
primary_req_id: USER-FAMILY-001
status: 已通过
preconditions:
  - 已具备用户端登录态，或手动写入 bill_user_token
  - 浏览器可访问 /user/families
steps:
  - 打开 /user/families
  - 检查家庭名称、创建人、成员数等信息展示
  - 点击任一家庭“查看详情”，确认跳转到 /user/families/:familyId
expected:
  - 列表展示家庭名称和创建人
  - 可以进入家庭详情页
actual: 家庭列表可见且含创建人信息，点击“查看详情”后成功进入 /user/families/FAM-4821，控制台无业务 error
evidence:
  - assets/playwright/user-families-list.png
  - assets/playwright/user-families-list-snapshot.txt
  - assets/playwright/user-families-list-url.txt
  - assets/playwright/user-families-list-console.txt
  - assets/playwright/user-families-list-network.txt
  - assets/playwright/user-families-open-detail-run-fixed.txt
executor: Codex-Worker-D
last_executed_at: 2026-04-09 21:57
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/families`

Steps:

1. 打开 `/user/families`
2. 检查家庭名称、创建人、成员数等信息展示
3. 点击任一家庭“查看详情”，确认跳转到 `/user/families/:familyId`

Expected:

- 列表展示家庭名称和创建人
- 可以进入家庭详情页

Actual:

- 当前 URL：`http://localhost:9000/user/families`
- 列表页展示家庭名称、创建人、成员数和摘要统计
- 点击家庭卡片“查看详情”后跳转到 `http://localhost:9000/user/families/FAM-4821`

Evidence:

- `assets/playwright/user-families-list.png`
- `assets/playwright/user-families-list-snapshot.txt`
- `assets/playwright/user-families-list-url.txt`
- `assets/playwright/user-families-list-console.txt`
- `assets/playwright/user-families-list-network.txt`
- `assets/playwright/user-families-open-detail-run-fixed.txt`

### TC-ST-BATCH-20260409-USER-FAMILY-01-02 家庭创建加入退出入口

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-FAMILY-01-02
related_subtask_id: ST-BATCH-20260409-USER-FAMILY-01
primary_req_id: USER-FAMILY-001
status: 已通过
preconditions:
  - 已具备用户端登录态，或手动写入 bill_user_token
  - 浏览器可访问 /user/families
steps:
  - 创建一个新家庭并观察生成的家庭 ID
  - 使用家庭 ID 加入家庭
  - 使用邀请链接加入家庭
  - 触发退出家庭入口并确认状态反馈
expected:
  - 创建家庭后可看到独立家庭 ID
  - 用户可通过家庭 ID 或邀请链接加入家庭
  - 退出家庭后列表和统计即时更新
actual: 创建、按 ID 加入、按邀请链接加入、退出家庭流程都可完成；期间出现的 Playwright 选择器误报已登记并关闭
evidence:
  - assets/playwright/user-families-flow.png
  - assets/playwright/user-families-flow-snapshot.txt
  - assets/playwright/user-families-flow-url.txt
  - assets/playwright/user-families-flow-console.txt
  - assets/playwright/user-families-flow-network.txt
  - assets/playwright/user-families-create-run.txt
  - assets/playwright/user-families-join-id-run.txt
  - assets/playwright/user-families-join-link-run.txt
  - assets/playwright/user-families-leave-run.txt
  - assets/playwright/user-families-create-ambiguous-selector-error.txt
executor: Codex-Worker-D
last_executed_at: 2026-04-09 21:57
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/families`

Steps:

1. 创建一个新家庭并观察生成的家庭 ID
2. 使用家庭 ID 加入家庭
3. 使用邀请链接加入家庭
4. 触发退出家庭入口并确认状态反馈

Expected:

- 创建家庭后可看到独立家庭 ID
- 用户可通过家庭 ID 或邀请链接加入家庭
- 退出家庭后列表和统计即时更新

Actual:

- 创建家庭成功，列表新增记录并包含 `FAM-xxxx` 形式 ID
- 按家庭 ID 加入和按邀请链接加入均成功
- 退出家庭流程可执行并即时提示结果
- Playwright 过程中出现选择器误报，已按规则登记 `BUG-BATCH-20260409-USER-FAMILY-001` 并关闭

Evidence:

- `assets/playwright/user-families-flow.png`
- `assets/playwright/user-families-flow-snapshot.txt`
- `assets/playwright/user-families-flow-url.txt`
- `assets/playwright/user-families-flow-console.txt`
- `assets/playwright/user-families-flow-network.txt`
- `assets/playwright/user-families-create-run.txt`
- `assets/playwright/user-families-join-id-run.txt`
- `assets/playwright/user-families-join-link-run.txt`
- `assets/playwright/user-families-leave-run.txt`
- `assets/playwright/user-families-create-ambiguous-selector-error.txt`

### TC-ST-BATCH-20260409-USER-FAMILY-02-01 家庭详情月度与年度汇总

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-FAMILY-02-01
related_subtask_id: ST-BATCH-20260409-USER-FAMILY-02
primary_req_id: USER-FAMILY-002
status: 已通过
preconditions:
  - 已具备用户端登录态，或手动写入 bill_user_token
  - 浏览器可访问 /user/families/:familyId
steps:
  - 打开任一家庭详情页
  - 检查月收入、月支出、月结余
  - 检查年收入、年支出、年结余
expected:
  - 可查看指定月份和年份的汇总统计
  - 月度和年度统计都可切换并更新
actual: 家庭详情页可显示月度与年度收入/支出/结余汇总，页面加载和切换无控制台 error
evidence:
  - assets/playwright/user-family-detail.png
  - assets/playwright/user-family-detail-snapshot.txt
  - assets/playwright/user-family-detail-url.txt
  - assets/playwright/user-family-detail-console.txt
  - assets/playwright/user-family-detail-network.txt
executor: Codex-Worker-D
last_executed_at: 2026-04-09 21:57
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/families/:familyId`

Steps:

1. 打开任一家庭详情页
2. 检查月收入、月支出、月结余
3. 检查年收入、年支出、年结余

Expected:

- 可查看指定月份和年份的汇总统计
- 月度和年度统计都可切换并更新

Actual:

- 当前 URL：`http://localhost:9000/user/families/FAM-4821`
- 月度与年度统计卡片均展示收入、支出与结余
- 控制台 error 数量为 0

Evidence:

- `assets/playwright/user-family-detail.png`
- `assets/playwright/user-family-detail-snapshot.txt`
- `assets/playwright/user-family-detail-url.txt`
- `assets/playwright/user-family-detail-console.txt`
- `assets/playwright/user-family-detail-network.txt`

### TC-ST-BATCH-20260409-USER-FAMILY-02-02 家庭成员占比点击切换

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-FAMILY-02-02
related_subtask_id: ST-BATCH-20260409-USER-FAMILY-02
primary_req_id: USER-FAMILY-002
status: 已通过
preconditions:
  - 已具备用户端登录态，或手动写入 bill_user_token
  - 浏览器可访问 /user/families/:familyId
steps:
  - 在家庭详情页点击月收入、月支出、年收入、年支出四个入口
  - 观察成员占比图和列表的标题与数据变化
expected:
  - 点击不同入口后占比图切换到对应口径
  - 成员占比列表和总额随入口同步更新
actual: 四个入口点击后成员占比图和列表可切换，标题与总额按口径同步变化
evidence:
  - assets/playwright/user-family-detail-share-switch-run.txt
  - assets/playwright/user-family-detail.png
  - assets/playwright/user-family-detail-snapshot.txt
executor: Codex-Worker-D
last_executed_at: 2026-04-09 21:57
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/families/:familyId`

Steps:

1. 在家庭详情页点击月收入、月支出、年收入、年支出四个入口
2. 观察成员占比图和列表的标题与数据变化

Expected:

- 点击不同入口后占比图切换到对应口径
- 成员占比列表和总额随入口同步更新

Actual:

- 在详情页中依次点击四个入口，页面交互正常
- 成员占比组件标题、总额与成员占比列表同步更新

Evidence:

- `assets/playwright/user-family-detail-share-switch-run.txt`
- `assets/playwright/user-family-detail.png`
- `assets/playwright/user-family-detail-snapshot.txt`
