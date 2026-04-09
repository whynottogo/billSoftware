# 测试计划

```yaml
batch_id: BATCH-20260409-CLOSED-LOOP-QA
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - /user/register
  - /user/login
  - /admin/login
  - /admin/users
entry_points:
  - http://localhost:9000/user/register
  - http://localhost:9000/user/login
  - http://localhost:9000/admin/login
  - http://localhost:9000/admin/users
verification_strategy: Playwright 优先
risks:
  - Playwright CLI 对阻塞性 alert 的处理有限，因此旧会话失效以“前端弹窗 + 后端 401 日志”双证据确认
last_updated_at: 2026-04-10 00:18
```

## 测试范围

- 覆盖页面：注册、用户登录、管理员登录、管理端用户列表
- 覆盖核心流程：注册默认禁用、启用可登录、禁用不可登录、旧会话失效回跳

## 本轮执行结果

1. 已完成真实闭环回归。
2. 首次阻塞 504 的根因是后端环境未拉起；第二次阻塞 500 的根因是远端数据库未初始化 schema。
3. 修复后，五条固定用例全部通过。
