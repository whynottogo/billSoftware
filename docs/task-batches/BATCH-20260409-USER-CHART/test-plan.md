# 测试计划

```yaml
batch_id: BATCH-20260409-USER-CHART
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - /user/charts/expense 支出图表页
  - /user/charts/income 收入图表页
entry_points:
  - /user/charts/expense
  - /user/charts/income
verification_strategy: Playwright 优先
risks:
  - 共享路由文件并行改动频繁，可能导致图表路由短时不可达
  - 图表数据口径调整可能影响排行榜和折线图一致性
last_updated_at: 2026-04-09 22:21
```

## 测试范围

- 覆盖页面：`/user/charts/expense`、`/user/charts/income`
- 覆盖接口：本批次不做真实接口联调，验证本地 chart mock 渲染结果
- 覆盖核心流程：图表页访问、年份切换、排行榜展示、路由互跳

## 入口与前置条件

- 页面入口：`http://localhost:9000/user/charts/expense` 与 `http://localhost:9000/user/charts/income`
- 登录账号：通过 `bill_user_token` 与 `bill_user_profile` 注入用户端登录态
- 初始化数据：使用 `userChartMock.js` 数据
- 环境依赖：前端服务可访问，Playwright wrapper 可执行

## 执行顺序

1. 访问支出图表页并验证月/年折线与排行榜
2. 验证支出页年份切换和路由按钮
3. 访问收入图表页并验证年折线与排行榜
4. 验证收入页年份切换
5. 保存截图、URL、控制台摘要、关键网络请求

## 验证策略

- 前端流程使用 Playwright wrapper 进行真实页面验证
- 图表页每条用例都归档四件套证据
- 出现异常时先更新 `bugs.md` 再继续排查

## 风险记录

- 风险 1：图表组件参数不一致会导致渲染空白或坐标轴错位
- 风险 2：共享路由变更冲突会导致图表页短时不可达
