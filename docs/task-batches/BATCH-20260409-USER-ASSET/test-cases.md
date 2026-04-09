# 测试用例

## 用例状态统计

| status | count |
| --- | --- |
| 未开始 | 0 |
| 执行中 | 0 |
| 已通过 | 4 |
| 已阻塞 | 0 |

### TC-ST-BATCH-20260409-USER-ASSET-01-01 资产总览与账户分类列表展示

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-ASSET-01-01
related_subtask_id: ST-BATCH-20260409-USER-ASSET-01
primary_req_id: USER-ASSET-001
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:58
evidence:
  - assets/playwright/user-assets-tc01.png
  - assets/playwright/user-assets-tc01-snapshot.md
  - assets/playwright/user-assets-tc01-url.txt
  - assets/playwright/user-assets-tc01-console-error.log
  - assets/playwright/user-assets-tc01-network.log
  - assets/playwright/user-assets-tc01-detail-url.txt
  - assets/playwright/user-assets-tc01-detail.png
```

Preconditions:

- 已具备用户端登录态，或手动写入 `bill_user_token`
- 浏览器可访问 `/user/assets`

Steps:

1. 打开 `/user/assets`
2. 检查净资产、资产、负债指标卡片
3. 检查账户分类列表、分类总额和账户余额
4. 点击任一账户进入详情页

Expected:

- 页面展示净资产、资产、负债三项指标
- 列表按分类展示账户，分类显示总额，账户显示余额
- 点击账户可进入 `/user/assets/:accountId`

Actual:

- 页面展示净资产、资产、负债指标与分类账户列表
- 账户余额和分类总额展示正常
- 点击详情后成功进入 `http://localhost:9000/user/assets/acc-cmb-6225`
- 控制台无 error

Evidence:

- `assets/playwright/user-assets-tc01.png`
- `assets/playwright/user-assets-tc01-snapshot.md`
- `assets/playwright/user-assets-tc01-url.txt`
- `assets/playwright/user-assets-tc01-console-error.log`
- `assets/playwright/user-assets-tc01-network.log`
- `assets/playwright/user-assets-tc01-detail-url.txt`
- `assets/playwright/user-assets-tc01-detail.png`

### TC-ST-BATCH-20260409-USER-ASSET-02-01 新增与编辑账户表单校验

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-ASSET-02-01
related_subtask_id: ST-BATCH-20260409-USER-ASSET-02
primary_req_id: USER-ASSET-002
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 21:59
evidence:
  - assets/playwright/user-assets-tc02-validation.log
  - assets/playwright/user-assets-tc02-type-immutable.log
  - assets/playwright/user-assets-tc02.png
  - assets/playwright/user-assets-tc02-snapshot.md
  - assets/playwright/user-assets-tc02-url.txt
  - assets/playwright/user-assets-tc02-console-error.log
  - assets/playwright/user-assets-tc02-network.log
```

Preconditions:

- 已具备用户端登录态
- 浏览器可访问 `/user/assets`

Steps:

1. 在资产页打开新增账户弹窗
2. 校验名称、备注、余额必填
3. 切换为银行卡或信用卡时校验卡号必填
4. 切换为虚拟账户时校验渠道必填
5. 进入编辑账户并确认类型不可修改

Expected:

- 新增账户必填项完整生效
- 银行卡/信用卡卡号规则生效，虚拟账户渠道规则生效
- 编辑账户时类型字段不可修改

Actual:

- 新增银行卡缺少卡号时提示：`银行卡和信用卡账户必须填写卡号`
- 补齐卡号后创建成功提示：`账户已创建（原型本地数据）`
- 新增虚拟账户缺少渠道时提示：`虚拟账户必须区分微信或支付宝`
- 编辑态下即使尝试改类型，账户类型仍保持不变（`typeImmutable: true`）
- 控制台无 error

Evidence:

- `assets/playwright/user-assets-tc02-validation.log`
- `assets/playwright/user-assets-tc02-type-immutable.log`
- `assets/playwright/user-assets-tc02.png`
- `assets/playwright/user-assets-tc02-snapshot.md`
- `assets/playwright/user-assets-tc02-url.txt`
- `assets/playwright/user-assets-tc02-console-error.log`
- `assets/playwright/user-assets-tc02-network.log`

### TC-ST-BATCH-20260409-USER-ASSET-03-01 账户详情三种余额操作与按月记录

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-ASSET-03-01
related_subtask_id: ST-BATCH-20260409-USER-ASSET-03
primary_req_id: USER-ASSET-003
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:00
evidence:
  - assets/playwright/user-assets-tc03-operations.log
  - assets/playwright/user-assets-tc03-month-switch.log
  - assets/playwright/user-assets-tc03.png
  - assets/playwright/user-assets-tc03-snapshot.md
  - assets/playwright/user-assets-tc03-url.txt
  - assets/playwright/user-assets-tc03-console-error.log
  - assets/playwright/user-assets-tc03-network.log
```

Preconditions:

- 已具备用户端登录态
- 浏览器可访问 `/user/assets/:accountId`

Steps:

1. 打开账户详情页
2. 分别执行调整余额、增加余额、减少余额操作
3. 切换月份筛选余额变动记录
4. 检查记录区按月份展示结果

Expected:

- 页面支持三种余额操作并成功追加记录
- 页面支持按月切换查看余额变动记录
- 余额数字和记录列表随着操作同步更新

Actual:

- 执行调整/增加/减少后，最新三条记录包含三种动作
- 记录月份集合包含 `2026-04`、`2026-03`、`2026-02`
- 通过月份下拉切换后，记录行数收敛为 `2`（筛选生效）
- 控制台无 error

Evidence:

- `assets/playwright/user-assets-tc03-operations.log`
- `assets/playwright/user-assets-tc03-month-switch.log`
- `assets/playwright/user-assets-tc03.png`
- `assets/playwright/user-assets-tc03-snapshot.md`
- `assets/playwright/user-assets-tc03-url.txt`
- `assets/playwright/user-assets-tc03-console-error.log`
- `assets/playwright/user-assets-tc03-network.log`

### TC-ST-BATCH-20260409-USER-ASSET-03-02 设置修改余额自动追加调整记录

```yaml
test_case_id: TC-ST-BATCH-20260409-USER-ASSET-03-02
related_subtask_id: ST-BATCH-20260409-USER-ASSET-03
primary_req_id: USER-ASSET-003
status: 已通过
executor: Codex
last_executed_at: 2026-04-09 22:00
evidence:
  - assets/playwright/user-assets-tc04-settings.log
  - assets/playwright/user-assets-tc04.png
  - assets/playwright/user-assets-tc04-snapshot.md
  - assets/playwright/user-assets-tc04-url.txt
  - assets/playwright/user-assets-tc04-console-error.log
  - assets/playwright/user-assets-tc04-network.log
```

Preconditions:

- 已具备用户端登录态
- 浏览器可访问 `/user/assets/:accountId`

Steps:

1. 打开账户设置弹窗并修改账户余额
2. 保存设置
3. 检查余额变动记录列表顶部新增一条调整记录

Expected:

- 通过设置修改余额后，余额数值更新
- 记录列表新增一条来源为设置的调整记录

Actual:

- 设置保存后余额从 `36850.3` 更新为 `37171.3`
- 最新记录为 `调整`，来源为 `账户设置`
- `isAdjustmentFromSettings: true`
- 控制台无 error

Evidence:

- `assets/playwright/user-assets-tc04-settings.log`
- `assets/playwright/user-assets-tc04.png`
- `assets/playwright/user-assets-tc04-snapshot.md`
- `assets/playwright/user-assets-tc04-url.txt`
- `assets/playwright/user-assets-tc04-console-error.log`
- `assets/playwright/user-assets-tc04-network.log`
