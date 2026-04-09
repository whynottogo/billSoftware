# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 3 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-01-01 新增支出后即时回显

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-01-01
related_subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-01
primary_req_id: USER-LEDGER-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:08
evidence:
  - assets/playwright/tc01-final-form-check.txt
  - assets/playwright/tc01-final-assert.txt
  - assets/playwright/tc01-final-after-save.png
  - assets/playwright/tc01-final-after-save-snapshot.yml
  - assets/playwright/tc01-final-after-save-url.txt
  - assets/playwright/tc01-final-console.log
  - assets/playwright/tc01-final-network.log
```

Preconditions:

- 前端构建已完成（`npm run build`）
- 使用独立静态服务 `http://localhost:9101` 做回归，避免并行 worker 引起的 dev-server overlay 干扰
- 已写入用户端登录态并可进入 `/user/ledger`

Steps:

1. 打开 `/user/ledger`，点击“记一笔支出”
2. 检查日期默认当天、图片输入 `multiple=false`
3. 填写分类、金额、备注并保存
4. 在当月列表确认新增记录即时可见

Expected:

- 新增支出弹层字段完整
- 日期默认当天、图片可选且最多 1 张
- 保存后无需刷新即可看到新增支出记录

Actual:

- 日期默认值为 `2026-04-09`，图片输入 `fileMultiple=false`
- 保存后页面仍在 `/user/ledger`，新增备注“测试新增支出-WorkerE-Final2”已出现在列表
- 控制台无 error

Evidence:

- `assets/playwright/tc01-final-form-check.txt`
- `assets/playwright/tc01-final-assert.txt`
- `assets/playwright/tc01-final-after-save.png`
- `assets/playwright/tc01-final-after-save-snapshot.yml`
- `assets/playwright/tc01-final-after-save-url.txt`
- `assets/playwright/tc01-final-console.log`
- `assets/playwright/tc01-final-network.log`

### TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-02-01 新增收入后即时回显

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-02-01
related_subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-02
primary_req_id: USER-LEDGER-003
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:08
evidence:
  - assets/playwright/tc02-final-assert.txt
  - assets/playwright/tc02-final-after-save.png
  - assets/playwright/tc02-final-after-save-snapshot.yml
  - assets/playwright/tc02-final-after-save-url.txt
  - assets/playwright/tc02-final-console.log
  - assets/playwright/tc02-final-network.log
```

Preconditions:

- 可访问 `http://localhost:9101/user/ledger`
- 已写入用户端登录态

Steps:

1. 点击“记一笔收入”
2. 填写收入分类、金额、备注并保存
3. 验证新增记录即时回显

Expected:

- 新增收入流程可执行
- 保存后无需刷新即可在列表看到新增收入

Actual:

- 保存后页面仍在 `/user/ledger`，新增备注“测试新增收入-WorkerE-Final2”已出现在列表
- 控制台无 error

Evidence:

- `assets/playwright/tc02-final-assert.txt`
- `assets/playwright/tc02-final-after-save.png`
- `assets/playwright/tc02-final-after-save-snapshot.yml`
- `assets/playwright/tc02-final-after-save-url.txt`
- `assets/playwright/tc02-final-console.log`
- `assets/playwright/tc02-final-network.log`

### TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-03-01 收支分类管理

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-LEDGER-ACTIONS-03-01
related_subtask_id: ST-BATCH-20260409-USER-LEDGER-ACTIONS-03
primary_req_id: USER-LEDGER-004
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:08
evidence:
  - assets/playwright/tc03-rerun-income-tab-snapshot.yml
  - assets/playwright/tc03-final-assert.txt
  - assets/playwright/tc03-final-after-add.png
  - assets/playwright/tc03-final-after-add-snapshot.yml
  - assets/playwright/tc03-final-after-add-url.txt
  - assets/playwright/tc03-final-console.log
  - assets/playwright/tc03-final-network.log
```

Preconditions:

- 可访问 `http://localhost:9101/user/ledger`
- 已写入用户端登录态

Steps:

1. 进入“分类管理”
2. 在支出分类新增“测试支出分类Final”
3. 切换收入分类并新增“测试收入分类Final”
4. 校验分类列表显示新增项

Expected:

- 支出与收入分类可分别新增
- 新增后分类列表即时显示
- 页面无异常报错

Actual:

- 支出分类新增项已出现在分类管理列表（见 `tc03-rerun-income-tab-snapshot.yml`）
- 收入分类新增项“测试收入分类Final”已显示（见 `tc03-final-after-add-snapshot.yml`）
- 控制台无 error

Evidence:

- `assets/playwright/tc03-rerun-income-tab-snapshot.yml`
- `assets/playwright/tc03-final-assert.txt`
- `assets/playwright/tc03-final-after-add.png`
- `assets/playwright/tc03-final-after-add-snapshot.yml`
- `assets/playwright/tc03-final-after-add-url.txt`
- `assets/playwright/tc03-final-console.log`
- `assets/playwright/tc03-final-network.log`
