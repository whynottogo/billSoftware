# 测试计划

```yaml
batch_id: BATCH-20260409-ADMIN-PORTAL
status: 已完成
environment:
  frontend: http://localhost:9000
  backend: http://localhost:8080
  database: billSoftware@qiqiqu.cn:3306
verification_strategy: Playwright 优先
last_updated_at: 2026-04-09 21:46
```

## 测试范围

- 覆盖页面：`/admin/dashboard`、`/admin/approvals`、`/admin/families`
- 覆盖核心流程：管理端登录态校验、三页渲染稳定性、图表/列表区块可见、页面无脚本错误
- 覆盖编译验证：`frontend` 执行 `npm run build`

## 入口与前置条件

- 页面入口：浏览器可访问 `http://localhost:9000`
- 登录账号：通过本地存储写入 `bill_admin_token` 作为受保护路由访问凭证
- 初始化数据：页面使用 `adminPortalMock` 静态数据，不依赖后端返回
- 环境依赖：前端开发服务正常启动；Playwright wrapper 可用

## 执行顺序

1. 先运行 `npm run build`，确认新增页面可编译
2. 写入管理端登录态并访问 `/admin/dashboard`
3. 依次访问 `/admin/approvals` 与 `/admin/families`
4. 每页采集截图、URL、控制台摘要、关键网络请求
5. 回写测试用例、Bug 台账和批次状态

## 验证策略

- 前端页面验证使用 Playwright wrapper 进行真实浏览器实测
- 若控制台出现 error、路由异常、白屏或网络失败，先写 `bugs.md` 再排查
- 页面校验通过后再推进需求状态收口

## 风险记录

- 风险 1：后台三页当前为原型级数据，后续真实接口接入时字段命名可能变更
- 风险 2：共享壳层尚未做最终并线合并，当前使用临时路由入口验证新页面
