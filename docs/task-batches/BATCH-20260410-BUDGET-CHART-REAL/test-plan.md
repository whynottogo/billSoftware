# 测试计划

```yaml
batch_id: BATCH-20260410-BUDGET-CHART-REAL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先（浏览器 MCP 初始化异常时切换本机 Playwright wrapper）
last_updated_at: 2026-04-10 12:54
```

## 测试范围

- 覆盖页面：
  - `/user/budget/month`
  - `/user/budget/year`
  - `/user/charts/expense`
  - `/user/charts/income`
- 覆盖接口：
  - `GET /api/user/budgets/month/current`
  - `PUT /api/user/budgets/month/current`
  - `GET /api/user/budgets/year/options`
  - `GET /api/user/budgets/year/:year`
  - `PUT /api/user/budgets/year/current`
  - `GET /api/user/charts/years`
  - `GET /api/user/charts/expense/:year`
  - `GET /api/user/charts/income/:year`
- 覆盖核心流程：
  - 当前月份预算设置与回显
  - 当前年份预算设置与回显
  - 支出图表月趋势、年趋势、排行榜真实渲染
  - 收入图表年趋势、排行榜真实渲染

## 入口与前置条件

- 页面入口：
  - 用户端登录页 `http://localhost:9000/user/login`
- 登录账号：
  - 复用联调用户 `e2e0410000402`
- 初始化数据：
  - 通过真实预算接口写入月预算和年预算测试数据
  - 账本数据复用当前数据库中 `2025-11` 至 `2026-04` 的真实流水
- 环境依赖：
  - 本地前端、后端服务可访问
  - 代理 `/api -> http://127.0.0.1:8080` 可用

## 执行顺序

1. 启动前端与后端服务
2. 通过接口验证预算与图表结构
3. 登录用户端并完成月预算设置回归
4. 进入年预算页验证年度预算与月度执行图
5. 进入支出图表页验证月趋势、年趋势与排行榜
6. 切换收入图表验证年趋势与排行榜
7. 保存截图、快照、控制台摘要和关键网络请求
8. 执行最终 `build`、`go test` 与健康检查

## 验证策略

- 前端和联调任务使用 Playwright 做真实页面验证；浏览器 MCP 初始化失败时，按 Bug 先登记后切换本机 wrapper 继续回归
- 后端接口先用 JSON 证据落盘，确认返回结构与页面消费字段一致
- 预算设置必须验证保存后即时回显，不接受只验证请求成功
- 失败时先补充截图、控制台摘要、关键网络请求，再同步回写 `bugs.md`

## 风险记录

- 风险 1：预算表历史为空，首次写入可能暴露 `budgets` / `budget_items` upsert 或事务问题
- 风险 2：图表页对空数据和年份切换敏感，若返回数组长度异常可能导致 SVG 渲染错位
- 风险 3：Playwright MCP 在当前桌面环境下初始化目录失败，需要 wrapper 兜底；已登记并关闭

## 执行结果

- `npm run build` 已通过，证据：`docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-frontend-build-final.txt`
- `/Users/shynin/.local/go-1.26.2/bin/go test ./...` 已通过，证据：`docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-backend-go-test-final.txt`
- `GET /api/health` 已通过，证据：`docs/task-batches/BATCH-20260410-BUDGET-CHART-REAL/assets/playwright/20260410-backend-health-final.json`
- 月预算、年预算、支出图表、收入图表真实页面回归均已通过，详细结果见 `test-cases.md`
