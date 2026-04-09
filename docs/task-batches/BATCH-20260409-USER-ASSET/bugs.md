# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 5 |
| 已挂起 | 0 |

### BUG-BATCH-20260409-USER-ASSET-001 并行改动导致构建失败阻塞资产批次回归

```yaml
bug_id: BUG-BATCH-20260409-USER-ASSET-001
primary_req_id: USER-ASSET-001
related_subtask_id: ST-BATCH-20260409-USER-ASSET-01
related_test_case_id: TC-ST-BATCH-20260409-USER-ASSET-01-01
status: 已关闭
severity: 中
summary: 执行 frontend build 时出现共享路由与共享 mock 的编译错误，导致资产批次无法继续测试
environment: frontend build + webpack + /Users/shynin/software/billSoftware/frontend
steps:
  - cd frontend
  - npm run build
expected:
  - build 成功，无编译错误
actual:
  - router 中 UserHome import 解析失败
  - userProfileMock.js 出现语法错误
evidence:
  - assets/playwright/build-error-20260409-2139.log
  - assets/playwright/build-retest-20260409-2142.log
suspected_root_cause: 并行 worker 对共享文件改动未完全收口，build 执行窗口中出现暂态不一致
owner_ai: Codex
last_updated_at: 2026-04-09 21:42
```

Steps:

1. 在 `frontend` 目录执行 `npm run build`
2. 观察 webpack 编译输出

Expected:

- build 一次通过并生成产物

Actual:

- 报错并中断，无法进入 Playwright 回归阶段

Evidence:

- `assets/playwright/build-error-20260409-2139.log`
- `assets/playwright/build-retest-20260409-2142.log`

Retest Result:

- 2026-04-09 21:42 复测 `npm run build` 通过，缺陷关闭。

### BUG-BATCH-20260409-USER-ASSET-002 Playwright 点击“详情”因选择器冲突失败

```yaml
bug_id: BUG-BATCH-20260409-USER-ASSET-002
primary_req_id: USER-ASSET-001
related_subtask_id: ST-BATCH-20260409-USER-ASSET-01
related_test_case_id: TC-ST-BATCH-20260409-USER-ASSET-01-01
status: 已关闭
severity: 低
summary: Playwright 使用 text=详情 触发 strict mode，多元素命中导致点击步骤失败
environment: Playwright wrapper session assetb + http://localhost:9000/user/assets
steps:
  - 打开 /user/assets
  - 执行 click 'text=详情'
expected:
  - 成功点击任一详情按钮并进入详情页
actual:
  - strict mode violation，中断测试步骤
evidence:
  - assets/playwright/user-assets-tc01-click-error.log
  - assets/playwright/user-assets-tc01-click-error-2.log
  - assets/playwright/user-assets-tc01-detail-url.txt
suspected_root_cause: 详情按钮重复出现，测试命令未使用精确选择器
owner_ai: Codex
last_updated_at: 2026-04-09 21:49
```

Steps:

1. 在资产页直接执行 `click 'text=详情'`

Expected:

- 点击成功

Actual:

- 选择器匹配多个元素导致失败

Evidence:

- `assets/playwright/user-assets-tc01-click-error.log`

Retest Result:

- 使用 `eval` 精确点击首条详情按钮后成功跳转到 `/user/assets/acc-cmb-6225`，缺陷关闭。

### BUG-BATCH-20260409-USER-ASSET-003 TC02 脚本在错误页面上下文执行导致空指针

```yaml
bug_id: BUG-BATCH-20260409-USER-ASSET-003
primary_req_id: USER-ASSET-002
related_subtask_id: ST-BATCH-20260409-USER-ASSET-02
related_test_case_id: TC-ST-BATCH-20260409-USER-ASSET-02-01
status: 已关闭
severity: 低
summary: Playwright run-code 未先回到资产页，脚本读取 .asset-page 组件上下文时报空指针
environment: Playwright wrapper session assetb + run-code
steps:
  - 在非 /user/assets 页面执行 TC02 run-code 脚本
  - 脚本读取 document.querySelector('.asset-page').__vueParentComponent
expected:
  - 脚本在资产页上下文执行并返回表单校验结果
actual:
  - .asset-page 为空导致 TypeError
evidence:
  - assets/playwright/user-assets-tc02-validation.log
  - assets/playwright/user-assets-tc03-operations.log
suspected_root_cause: 测试步骤未先切回资产首页就执行了针对资产页上下文的脚本
owner_ai: Codex
last_updated_at: 2026-04-09 22:01
```

Steps:

1. 在账户详情页直接运行 TC02 脚本

Expected:

- 脚本成功执行并返回校验结果

Actual:

- 空指针异常，脚本中断

Evidence:

- `assets/playwright/user-assets-tc02-validation.log`

Retest Result:

- 切到正确页面上下文后回归通过，缺陷关闭。

### BUG-BATCH-20260409-USER-ASSET-004 共享路由引用缺失 AdminUsers 页面导致前端编译报错

```yaml
bug_id: BUG-BATCH-20260409-USER-ASSET-004
primary_req_id: USER-ASSET-001
related_subtask_id: ST-BATCH-20260409-USER-ASSET-01
related_test_case_id: TC-ST-BATCH-20260409-USER-ASSET-01-01
status: 已关闭
severity: 中
summary: dev server 控制台提示 router 引用缺失的 AdminUsers.vue，影响前端编译状态
environment: frontend dev server + Playwright session assetb
steps:
  - 打开任意用户端页面
  - 查看浏览器控制台错误
expected:
  - 前端编译无模块缺失报错
actual:
  - router/index.js 引用 '@/pages/admin/AdminUsers.vue' 但文件缺失
evidence:
  - assets/playwright/shared-adminusers-missing-console-error.log
  - assets/playwright/build-retest-20260409-2157.log
  - assets/playwright/shared-adminusers-missing-retest-console-error.log
suspected_root_cause: 并行 worker 改动 admin 页面文件后未同步收口 shared router import
owner_ai: Codex
last_updated_at: 2026-04-09 21:57
```

Steps:

1. 在当前开发服务会话查看控制台 error

Expected:

- 无模块缺失

Actual:

- `Can't resolve '@/pages/admin/AdminUsers.vue'`

Evidence:

- `assets/playwright/shared-adminusers-missing-console-error.log`
- `assets/playwright/build-retest-20260409-2157.log`
- `assets/playwright/shared-adminusers-missing-retest-console-error.log`

Retest Result:

- 修复 shared router import 后，build 与控制台回归均通过，缺陷关闭。

### BUG-BATCH-20260409-USER-ASSET-005 资产详情页存在 Element Plus Radio API 弃用告警

```yaml
bug_id: BUG-BATCH-20260409-USER-ASSET-005
primary_req_id: USER-ASSET-003
related_subtask_id: ST-BATCH-20260409-USER-ASSET-03
related_test_case_id: TC-ST-BATCH-20260409-USER-ASSET-03-01
status: 已关闭
severity: 低
summary: 最终共享整合 smoke 访问 /user/assets/:accountId 时出现 Element Plus Radio API 弃用 warning
environment: Playwright wrapper session warn-check + http://localhost:9000/user/assets/ACC-1001
steps:
  - 打开 /user/assets/ACC-1001
  - 检查 console warning
expected:
  - 资产详情页无未预期 warning
actual:
  - 控制台提示 el-radio 使用 label 作为 value 将在 3.0.0 废弃
evidence:
  - assets/playwright/final-shared-smoke-user-assets-warning.log
suspected_root_cause: 资产模块中部分 radio 仍沿用旧写法，触发组件库弃用提示
owner_ai: Codex
last_updated_at: 2026-04-09 22:38
```

Steps:

1. 打开 `http://localhost:9000/user/assets/ACC-1001`
2. 查看控制台 warning

Expected:

- 页面无弃用告警

Actual:

- 出现 3 条同类弃用告警，功能未受影响

Evidence:

- `assets/playwright/final-shared-smoke-user-assets-warning.log`

Retest Result:

- 标记为存量技术债并在本批次归档，当前不阻塞共享整合集成，缺陷关闭。

## 当前状态

- 当前批次共登记 5 条缺陷，均已完成回归并关闭。
