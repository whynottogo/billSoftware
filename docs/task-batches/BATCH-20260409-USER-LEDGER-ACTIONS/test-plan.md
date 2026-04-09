# 测试计划

```yaml
batch_id: BATCH-20260409-USER-LEDGER-ACTIONS
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - /user/ledger 新增支出动作壳
  - /user/ledger 新增收入动作壳
  - /user/ledger 收支分类管理壳
entry_points:
  - /user/ledger
verification_strategy: Playwright 优先
risks:
  - 共享页面 UserHome.vue 结构较大，新增弹层后可能影响现有账单列表展示
  - 新增 mock 数据写入逻辑可能导致金额汇总与分组展示不一致
last_updated_at: 2026-04-09 22:08
```

## 测试范围

- 覆盖页面：`/user/ledger`
- 覆盖接口：本批次不做真实接口联调，重点验证本地 mock 写入行为
- 覆盖核心流程：新增支出、新增收入、分类管理、新增后回显

## 入口与前置条件

- 页面入口：`http://localhost:9000/user/ledger`
- 登录账号：通过 `bill_user_token` 与 `bill_user_profile` 注入用户端登录态
- 初始化数据：使用 `userLedgerMock.js` 默认月份和默认分类
- 环境依赖：前端 dev 服务可访问，Playwright wrapper 可执行

## 执行顺序

1. 注入用户登录态并打开 `/user/ledger`
2. 执行新增支出流程并验证列表即时回显
3. 执行新增收入流程并验证列表即时回显
4. 执行分类管理流程并验证默认分类与新增分类
5. 记录截图、URL、控制台摘要和关键网络请求

## 验证策略

- 前端流程使用 Playwright wrapper 完成真实页面交互验证
- 每个核心流程都保存完整证据到 `assets/playwright/`
- 若出现异常，按“异常先登记后排查”规则立即记录 `bugs.md`

## 风险记录

- 风险 1：日期和金额解析异常可能导致新增记录显示在错误日期分组
- 风险 2：分类新增/删除后，旧分类记录展示需要保持可读且不崩溃
