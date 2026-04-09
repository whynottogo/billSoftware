# 测试计划

```yaml
batch_id: BATCH-20260409-AUTH-SHELL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-09 18:17
```

## 测试范围

- 覆盖页面：`/user/login`、`/user/register`、`/user/ledger`、`/admin/login`、`/admin/users`
- 覆盖接口：用户登录、用户注册、管理员登录相关前端调用与错误提示链路
- 覆盖核心流程：双端登录入口隔离、用户端导航、管理员端用户列表框架、退出登录跳转

## 入口与前置条件

- 页面入口：浏览器分别打开 `http://localhost:9000/user/login` 与 `http://localhost:9000/admin/login`
- 登录账号：优先使用现有联调账号；若后端未准备好，则验证前端回退提示和本地 mock token 流程
- 初始化数据：至少保留一条用户列表展示数据，或允许使用前端占位数据验证布局
- 环境依赖：前端开发服务正常启动，若涉及真实登录则后端服务和代理可用

## 执行顺序

1. 启动前端并确认路由可访问
2. 用 Playwright 打开用户端登录与注册页，检查表单、文案和入口切换
3. 进入用户端首页，检查导航、摘要区和退出登录交互
4. 打开管理员登录与用户列表页，检查双端隔离和后台壳样式
5. 保存截图、控制台摘要和关键网络请求证据

## 验证策略

- 前端页面改造优先使用 Codex tools 中的 Playwright 进行真实页面验证
- 同步执行 `npm run build` 验证 Vue 工程可编译
- 失败时必须补充截图、控制台摘要、关键网络请求到 `assets/playwright/`

## 风险记录

- 风险 1：Figma Make 原型基于 React/Tailwind，迁移到 Vue + Element Plus 后可能出现局部视觉偏差
- 风险 2：当前登录接口可能仍处于脚手架或联调初期，真实鉴权结果与页面替换进度并不完全同步
