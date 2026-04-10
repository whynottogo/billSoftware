# 测试计划

```yaml
batch_id: BATCH-20260410-ASSET-PROFILE-REAL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-10 16:13
```

## 测试范围

- 覆盖页面：
  - `/user/assets`
  - `/user/assets/:accountId`
  - `/user/profile`
  - `/user/profile/password`
- 覆盖接口：
  - `GET /api/user/assets`
  - `POST /api/user/assets`
  - `PUT /api/user/assets/:id`
  - `GET /api/user/assets/:id`
  - `POST /api/user/assets/:id/operations`
  - `GET /api/user/profile`
  - `PUT /api/user/profile`
  - `PUT /api/user/profile/password`
- 覆盖核心流程：
  - 资产总览真实展示与分类筛选
  - 账户新增、编辑、详情和余额变动记录
  - 个人资料查看、头像上传、昵称/邮箱修改
  - 修改密码后重新登录与旧会话失效

## 入口与前置条件

- 页面入口：
  - 用户端登录页 `http://localhost:9000/user/login`
- 登录账号：
  - 复用已启用联调用户
- 环境依赖：
  - 本地前端、后端服务可访问
  - 数据库可写入资产账户测试数据

## 执行顺序

1. 启动前端与后端服务
2. 先验证资产与个人信息 API 结构
3. 登录用户端并完成资产总览、账户新增/编辑、详情页回归
4. 进入个人信息页验证资料查看与修改
5. 进入修改密码页验证原密码校验、成功修改与重新登录
6. 抓取控制台、关键网络请求和页面截图

## 验证策略

- 页面与联调优先使用 Playwright 实测
- 后端接口先用 JSON 证据确认页面消费字段
- 修改密码必须验证旧登录态失效与重新登录，不接受只验证接口返回成功

## 风险记录

- 风险 1：资产账户与余额变动记录表可能尚无后端模型，需要一次性补齐模型、路由和聚合逻辑
- 风险 2：头像上传若先以 Base64 直传实现，需确保页面和接口字段保持一致
- 风险 3：改密后会话失效会触碰用户鉴权共享逻辑，只允许单后端 owner 修改
- 风险 4：资料页最初出现的昵称/邮箱串值已通过精确 selector 复核为 Playwright wrapper 定位偏差，`BUG-BATCH-20260410-ASSET-PROFILE-REAL-003` 已关闭
