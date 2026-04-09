# 测试计划

```yaml
batch_id: BATCH-20260409-ADMIN-USER
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-09 22:04
```

## 测试范围

- 覆盖页面：`/admin/users`、`/admin/users/:userId`
- 覆盖核心流程：用户启用/禁用状态切换、用户详情跳转、详情页月/年 tab 切换
- 覆盖编译验证：`frontend` 执行 `npm run build`

## 入口与前置条件

- 页面入口：浏览器可访问 `http://localhost:9000`
- 登录态要求：需具备管理员 token，或手动写入 `bill_admin_token`
- 环境依赖：前端开发服务正常启动，管理端路由可访问

## 执行顺序

1. 启动前端并确认管理端新增路由已注册
2. 执行 `npm run build` 验证工程可编译
3. 使用 Playwright 验证用户列表的启用/禁用交互和统计回写
4. 使用 Playwright 验证用户详情页月/年 tab 切换与汇总展示
5. 保存截图、页面快照、URL、控制台摘要和关键网络请求证据

## 验证策略

- 采用 Playwright wrapper 进行真实页面验证，不依赖静态代码判断
- 一旦发现控制台 error、路由异常、4xx/5xx、空白页，先写 `bugs.md` 再排查
- 每条测试用例必须绑定可追溯证据路径

## 风险记录

- 风险 1：并行 worker 可能同时修改共享路由或布局，导致管理端入口冲突（已通过回归修复并关闭缺陷）
- 风险 2：管理端详情页统计与用户端口径存在偏差风险，需要在 mock 汇总层显式对齐
