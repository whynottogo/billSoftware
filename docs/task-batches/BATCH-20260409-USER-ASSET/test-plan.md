# 测试计划

```yaml
batch_id: BATCH-20260409-USER-ASSET
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-09 22:01
```

## 测试范围

- 覆盖页面：`/user/assets`、`/user/assets/:accountId`
- 覆盖核心流程：资产总览展示、账户分类列表、账户详情跳转、余额操作、按月记录筛选、账户新增编辑表单校验
- 覆盖编译验证：`frontend` 执行 `npm run build`

## 入口与前置条件

- 页面入口：浏览器可访问 `http://localhost:9000`
- 登录态要求：需具备用户端 token，或通过写入 `bill_user_token` 和 `bill_user_profile` 访问 `/user/*` 受保护页面
- 环境依赖：前端开发服务正常启动，Playwright wrapper 可用

## 执行顺序

1. 启动前端并确认资产模块路由可访问
2. 执行 `npm run build` 验证工程可编译
3. 使用 Playwright 访问资产总览页并验证资产指标、分类列表和详情跳转
4. 使用 Playwright 验证新增账户和编辑账户表单差异化字段校验
5. 使用 Playwright 验证账户详情页三种余额操作、设置改余额和按月筛选记录
6. 保存截图、页面快照、控制台摘要和关键网络请求证据

## 验证策略

- 前端页面改造优先使用 Playwright 进行真实页面验证
- 若验证过程中出现构建失败、路由错误、控制台异常、4xx/5xx，先写入 `bugs.md` 再继续排查
- 每个测试用例都需归档证据到 `assets/playwright/`
- 本批次已完成 build 与 4 条测试用例回归，Bug 已全部关闭

## 风险记录

- 风险 1：资产详情在本轮为本地状态模拟，刷新页面后操作记录不会持久化
- 风险 2：并行协作下若共享路由被其他 worker 同时修改，可能影响详情页跳转验证
