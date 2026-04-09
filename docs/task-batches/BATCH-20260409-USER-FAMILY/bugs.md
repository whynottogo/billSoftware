# Bug 台账

## Bug 状态统计

| status | count |
| --- | --- |
| 新建 | 0 |
| 已确认 | 0 |
| 修复中 | 0 |
| 待回归 | 0 |
| 已关闭 | 2 |
| 已挂起 | 0 |

### BUG-BATCH-20260409-USER-FAMILY-001 Playwright 创建家庭回归脚本点击目标被遮罩拦截

```yaml
bug_id: BUG-BATCH-20260409-USER-FAMILY-001
primary_req_id: USER-FAMILY-001
related_subtask_id: ST-BATCH-20260409-USER-FAMILY-01
related_test_case_id: TC-ST-BATCH-20260409-USER-FAMILY-01-02
status: 已关闭
severity: 低
summary: 回归脚本使用不精确定位（按钮与表单 label 匹配不唯一），导致点击或填充操作被遮罩/同名控件拦截
environment:
  frontend: http://localhost:9000
  page: /user/families
  tool: /Users/shynin/.codex/skills/playwright/scripts/playwright_cli.sh
steps:
  - 打开 /user/families
  - 触发创建家庭弹窗
  - 脚本继续点击页面主按钮而非弹窗 footer “创建”按钮
  - 在加入家庭弹窗中，脚本用 label “家庭 ID” 命中了 radio 与输入框两个元素
expected:
  - 测试脚本应精确点击弹窗内按钮并正常提交
actual:
  - Playwright 报错 overlay intercepts pointer events，导致 timeout
  - Playwright 报错 strict mode violation，`家庭 ID` label 匹配到多个元素
evidence:
  - assets/playwright/user-families-create-ambiguous-selector-error.txt
  - assets/playwright/user-families-flow-run.txt
  - assets/playwright/user-families-flow-run-success.txt
  - assets/playwright/user-families-open-detail-run.txt
suspected_root_cause: 自动化脚本定位器不够精确且未处理弹窗遮罩状态，不属于业务页面缺陷
owner_ai: Codex-Worker-D
last_updated_at: 2026-04-09 21:57
```

Steps:

1. 打开 `http://localhost:9000/user/families`
2. 点击“创建家庭”打开弹窗
3. 继续执行未加限定的点击命令

Expected:

- 只点击弹窗内“创建”按钮，测试顺利推进

Actual:

- 点击目标被弹窗遮罩拦截，出现 `overlay intercepts pointer events`
- 判定为自动化脚本误报，已在同批次修正测试步骤并关闭

Evidence:

- `assets/playwright/user-families-create-ambiguous-selector-error.txt`
- `assets/playwright/user-families-flow-run.txt`
- `assets/playwright/user-families-flow-run-success.txt`
- `assets/playwright/user-families-open-detail-run.txt`

## 当前状态

- 当前批次累计记录 2 条缺陷，均已关闭。
- 后续仍按“发现异常先记 Bug”的流程继续执行回归。

### BUG-BATCH-20260409-USER-FAMILY-002 并行改动导致路由编译错误阻塞页面回归

```yaml
bug_id: BUG-BATCH-20260409-USER-FAMILY-002
primary_req_id: USER-FAMILY-001
related_subtask_id: ST-BATCH-20260409-USER-FAMILY-01
related_test_case_id: TC-ST-BATCH-20260409-USER-FAMILY-01-01
status: 已关闭
severity: 中
summary: 并行 worker 改动后 admin 页面文件缺失，但路由仍引用 AdminUsers.vue，导致 webpack dev server 编译 error
environment:
  frontend: http://localhost:9000
  page: /user/families
  tool: /Users/shynin/.codex/skills/playwright/scripts/playwright_cli.sh
steps:
  - 打开 /user/families
  - 观察 console error 输出
expected:
  - 路由引用的页面文件应存在，dev server 正常编译
actual:
  - 控制台出现 Module not found，无法解析 @/pages/admin/AdminUsers.vue
evidence:
  - assets/playwright/user-families-goto-console-error.txt
  - assets/playwright/user-families-goto-console-error-retest.txt
suspected_root_cause: 共享路由文件和 admin 页面文件被并行修改后未同步，产生悬挂 import
owner_ai: Codex-Worker-D
last_updated_at: 2026-04-09 21:56
```
