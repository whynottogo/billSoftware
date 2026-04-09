# 测试计划

```yaml
batch_id: BATCH-20260409-BILL-BUDGET
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-09 21:47
```

## 测试范围

- 覆盖页面：`/user/bills/month`、`/user/bills/month/:month`、`/user/bills/year`、`/user/budget/month`、`/user/budget/year`
- 覆盖核心流程：用户端账单导航跳转、月份详情路由进入、年视图切换、预算状态展示、页面无脚本错误
- 覆盖编译验证：`frontend` 执行 `npm run build`

## 入口与前置条件

- 页面入口：浏览器可访问 `http://localhost:9000`
- 登录态要求：需具备用户端 token，或通过写入 `bill_user_token` 访问 `/user/*` 受保护页面
- 环境依赖：前端开发服务正常启动；若需要接口代理，后端服务可作为补充，但本批次页面验证以静态原型替换为主

## 执行顺序

1. 启动前端并确认新增路由已注册
2. 执行 `npm run build` 验证工程可编译
3. 使用 Playwright 访问月账单列表并检查跳转到月详情和年账单
4. 使用 Playwright 分别访问月预算和年预算页面并检查布局、摘要和交互占位
5. 保存截图、页面快照、控制台摘要和关键网络请求证据

## 验证策略

- 前端页面改造优先使用 Codex tools 中的 Playwright 进行真实页面验证
- 对于受保护路由，使用与上一批次一致的本地用户端登录态注入方式
- 若验证过程中出现构建失败、路由错误或控制台异常，必须同步更新 `bugs.md`
- 阻塞缺陷已修复并完成回归，本批次 5 条测试用例全部通过

## 风险记录

- 风险 1：Figma Make 的页面密度较高，迁移到 Vue 单文件组件后可能出现局部排版差异
- 风险 2：当前账单和预算页面先以前端静态数据替换壳层，真实接口字段接入后可能需要二次对齐
- 风险 3：并行 worker 修改共享页面时可能引入编译异常，本批次通过补充回归和兼容修复已解除阻塞
