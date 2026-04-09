# 测试计划

```yaml
batch_id: BATCH-20260409-LEDGER-CORE-REAL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - GET /api/user/ledger
  - POST /api/user/ledger
  - GET /api/user/categories
  - POST /api/user/categories
  - DELETE /api/user/categories/:id
  - /user/ledger
verification_strategy: Playwright + 接口验证
last_updated_at: 2026-04-10 00:18
```

## 测试范围

- 覆盖页面：`/user/ledger`
- 覆盖接口：真实月收支、真实新增收支、真实分类管理
- 覆盖核心流程：用户隔离、记账回显、默认分类、自定义分类管理、禁用后会话失效

## 执行结果

1. 已安装本地 Go 运行环境并生成 `backend/go.sum`。
2. 已初始化远端数据库 schema 与默认分类种子。
3. 已通过真实页面验证新增支出、新增收入和分类管理。
4. 已通过接口验证第二个用户读取不到第一个用户的月账本数据。
