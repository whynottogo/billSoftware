# 测试计划

```yaml
batch_id: BATCH-20260409-USER-FAMILY
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
test_scope:
  - /user/families 家庭列表与创建加入退出入口
  - /user/families/:familyId 家庭详情月度年度统计与成员占比
entry_points:
  - http://localhost:9000/user/families
  - http://localhost:9000/user/families/:familyId
verification_strategy: Playwright 优先
risks:
  - 共享路由文件由多个 worker 并行修改，可能出现路由冲突
  - 家庭详情统计使用前端 mock，后续联调字段可能需要二次适配
last_updated_at: 2026-04-09 21:57
```

## 测试范围

- 覆盖页面：`/user/families`、`/user/families/:familyId`
- 覆盖流程：家庭列表展示、创建家庭、按 ID 加入、按邀请链接加入、主动退出、进入家庭详情
- 覆盖交互：月收入/月支出/年收入/年支出点击后成员占比切换
- 覆盖编译验证：`frontend` 执行 `npm run build`

## 入口与前置条件

- 页面入口：浏览器访问 `http://localhost:9000/user/families`
- 登录态要求：手动写入 `bill_user_token` 与 `bill_user_profile` 后访问 `/user/*` 受保护路由
- 环境依赖：前端开发服务启动正常；本批次不依赖后端接口可用性

## 执行顺序

1. 完成家庭模块页面和 mock 数据实现
2. 执行 `npm run build` 验证工程可编译
3. 使用 Playwright 访问家庭列表页并执行创建、加入、退出与详情跳转
4. 使用 Playwright 访问家庭详情页并验证月/年统计与占比切换
5. 保存截图、页面 URL、控制台摘要、关键网络请求证据

## 验证策略

- 前端页面改造优先使用 Playwright 真实页面验证
- 每个测试用例都归档四件套证据：截图、当前 URL、控制台摘要、关键网络请求
- 任何异常先登记 `bugs.md`，再开始排查和修复

## 风险记录

- 风险 1：家庭列表交互包含前端可变状态，测试过程需控制数据可重复性
- 风险 2：详情占比交互依赖路由参数，若共享路由冲突会影响整批验证

## 执行结果

- 已完成 `frontend npm run build`，构建通过（仅有 bundle 体积警告，无编译错误）
- 已完成家庭列表与家庭详情的 Playwright 回归，页面无控制台 error
- 过程出现两次 Playwright 选择器误报，已登记并关闭为工具脚本问题，不影响业务功能
