# 测试计划

```yaml
batch_id: BATCH-20260409-ADMIN-USER-REAL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - /admin/users
  - GET /api/admin/users
  - PUT /api/admin/users/:id/status
verification_strategy: Playwright 优先
last_updated_at: 2026-04-10 00:18
```

## 测试范围

- 覆盖页面：`/admin/users`
- 覆盖接口：真实用户列表、真实启停接口
- 覆盖核心流程：管理员登录 -> 用户列表 -> 启停 -> 用户端登录影响

## 入口与前置条件

- 页面入口：`/admin/login`、`/admin/users`、`/user/login`
- 初始化数据：新注册测试用户 `e2e0410000402`、管理员账号 `admin`
- 环境依赖：前后端可启动，管理员账号可用

## 执行顺序

1. 管理员登录并进入用户列表
2. 校验真实列表展示与统计卡
3. 执行启用操作并验证用户端可登录
4. 执行禁用操作并验证旧会话失效与重新登录失败
5. 归档截图、YAML 快照与后端日志证据

## 验证策略

- 使用 Playwright CLI 执行真实页面联调
- 同时结合后端日志核对 `GET /api/admin/users` 与 `PUT /api/admin/users/:id/status`
- 已完成静态契约复核 + `npm run build` + 真实页面回归
