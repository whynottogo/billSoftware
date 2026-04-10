# 测试计划

```yaml
batch_id: BATCH-20260410-USER-FAMILY-REAL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - /user/families 家庭列表、创建、加入、退出
  - /user/families/:familyId 家庭详情月度年度统计与成员占比
  - /api/user/families*
entry_points:
  - http://localhost:9000/user/login
  - http://localhost:9000/user/families
verification_strategy: Playwright 优先（必要时补接口 smoke 与构建验证）
risks:
  - 家庭成员关系表和 ledger 聚合口径可能与现有数据分布不完全匹配
  - 创建/加入/退出流程涉及真实数据库写入，需要控制测试数据可回归
  - 详情页依赖真实月度年度聚合，空家庭或空流水场景需要稳定兜底
last_updated_at: 2026-04-10 20:11
```

## 测试范围

- 覆盖页面：`/user/families`、`/user/families/:familyId`
- 覆盖接口：
  - `GET /api/user/families`
  - `POST /api/user/families`
  - `POST /api/user/families/join`
  - `POST /api/user/families/join-by-link`
  - `POST /api/user/families/:familyId/leave`
  - `GET /api/user/families/:familyId`
  - `GET /api/user/families/:familyId/member-share`
- 覆盖核心流程：
  - 家庭列表真实展示
  - 创建家庭生成独立家庭 ID
  - 按家庭 ID / 邀请链接加入家庭
  - 进入详情页查看月度年度汇总并切换成员占比
  - 退出家庭后列表和详情状态收敛

## 入口与前置条件

- 页面入口：
  - 用户端登录页 `http://localhost:9000/user/login`
- 登录账号：
  - 复用现有联调用户，必要时通过真实注册/登录接口获取新的测试账号
- 初始化数据：
  - 如数据库中不存在可用家庭，允许通过真实创建接口生成本批次测试家庭
  - 成员占比数据依赖家庭成员关联到账本用户，必要时使用现有有流水账号做联调
- 环境依赖：
  - 本地前端、后端服务可访问
  - 前端代理 `/api -> http://127.0.0.1:8080` 可用

## 执行顺序

1. 启动前端与后端服务
2. 通过接口 smoke 验证家庭列表、创建、加入、退出和详情接口结构
3. 登录用户端并回归家庭列表与创建/加入/退出流程
4. 进入家庭详情页验证月度年度汇总和成员占比切换
5. 保存截图、快照、控制台摘要和关键网络请求
6. 执行最终 `npm run build`、`go test ./...` 与健康检查

## 验证策略

- 前端和联调任务优先使用 Playwright 做真实页面验证
- 后端接口先确认返回结构可直接驱动当前 UI，再做浏览器联调
- 对创建/加入/退出动作必须验证页面状态和列表数据同步变化
- 失败时先补充截图、控制台摘要、关键网络请求，再同步回写 `bugs.md`

## 风险记录

- 风险 1：现网数据库 `families` 相关表结构可能与当前前端 mock 字段不完全一致，需要做兼容映射
- 风险 2：家庭创建/加入流程可能引入重复成员、重复家庭 ID 或无效邀请链接边界
- 风险 3：详情页统计依赖家庭成员到账本用户的映射，如仅存在家庭记录但无流水，需正确返回 0 值而非报错
- 风险 4：双用户闭环依赖现成联调账号，若历史密码已被其他批次修改会直接阻塞创建/加入链路；2026-04-10 19:16 已命中该风险并登记 `BUG-BATCH-20260410-USER-FAMILY-REAL-001`，当前已切换到新注册联调用户继续回归
- 风险 5：当前前端构建链未自动注入 `regeneratorRuntime`，若页面脚本继续使用 `async/await`，真实页面会在运行时直接报错；2026-04-10 19:45 已命中并登记 `BUG-BATCH-20260410-USER-FAMILY-REAL-003`
