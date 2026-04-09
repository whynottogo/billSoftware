# 测试计划

```yaml
batch_id: BATCH-20260409-USER-PROFILE-SESSION
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-09 21:50
```

## 测试范围

- 覆盖页面：`/user/profile`、`/user/profile/password`、`/user/session-kickout`、`/user/login`
- 覆盖核心流程：资料编辑保存、头像上传占位、密码修改校验与跳转、主动退出提示、会话失效提示承接
- 覆盖编译验证：`frontend` 执行 `npm run build`

## 入口与前置条件

- 页面入口：浏览器可访问 `http://localhost:9000`
- 登录态要求：访问 `/user/*` 受保护页面前写入 `bill_user_token` 与 `bill_user_profile`
- 环境依赖：前端开发服务已启动；后端仅作代理补充，本批次验证以前端表现层为主

## 执行顺序

1. 初始化用户端登录态并打开 `/user/profile`
2. 验证资料编辑、头像上传占位和保存提示
3. 进入 `/user/profile/password` 验证校验分支与成功后的跳转
4. 执行主动退出流程并验证登录页提示
5. 打开 `/user/session-kickout` 验证统一提示与登录承接
6. 保存截图、URL、控制台摘要和关键网络请求证据
7. 执行 `npm run build` 验证构建稳定

## 验证策略

- 前端任务优先使用 Playwright wrapper 进行真实页面验证
- 受保护路由通过写入本地登录态访问，不依赖真实用户接口
- 若出现构建失败、路由错误、控制台异常，先记录 `bugs.md` 再排查修复
- 首轮构建失败已记录 `BUG-BATCH-20260409-USER-PROFILE-SESSION-001`，修复后完成 build 回归和 4 个 Playwright 用例验证

## 风险记录

- 风险 1：并行开发期间共享路由文件可能存在冲突，需在回写时核对本批次相关路由是否保留
- 风险 2：当前会话提示使用前端 mock 承接，后续接入真实后端状态码时可能需要调整触发逻辑
