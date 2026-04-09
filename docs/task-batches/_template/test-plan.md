# 测试计划

```yaml
batch_id: BATCH-YYYYMMDD-TOPIC
status: 未开始
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: YYYY-MM-DD HH:MM
```

## 测试范围

- 覆盖页面：
- 覆盖接口：
- 覆盖核心流程：

## 入口与前置条件

- 页面入口：
- 登录账号：
- 初始化数据：
- 环境依赖：

## 执行顺序

1. 准备数据
2. 进入页面或接口入口
3. 执行核心流程
4. 校验结果
5. 保存证据

## 验证策略

- 前端和联调任务使用 Codex tools 中的 Playwright 进行真实页面验证
- 纯后端无页面入口任务可做接口验证
- 失败时必须补充截图、控制台摘要、关键网络请求

## 风险记录

- 风险 1：
- 风险 2：
