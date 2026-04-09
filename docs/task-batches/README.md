# task-batches

本目录用于按批次记录从任务池中领取的开发任务。

## 使用方式

1. 从 `docs/requirements/bill-software-requirements.md` 领取一个或多个 `req_id`
2. 创建批次目录：`docs/task-batches/<batch_id>/`
3. 可直接复制 `_template/` 作为起点
4. 补齐 `overview.md`、`subtasks.md`、`test-plan.md`、`test-cases.md`、`bugs.md`
5. 测试截图、控制台摘要、网络证据统一放到 `assets/playwright/`

## 推荐命名

- `BATCH-YYYYMMDD-<topic>`
- 示例：`BATCH-20260409-USER-AUTH`

## 注意事项

- `AGENTS.md` 负责定义工作流规则
- 本目录负责落地每个批次的执行记录
- 需求状态、批次状态、子任务状态、测试用例状态和 Bug 状态必须同步维护
