# 子任务拆分

| subtask_id | title | primary_req_id | related_req_ids | status | owner_ai | dependencies | linked_test_case_ids | linked_bug_ids | last_updated_at |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ST-BATCH-20260409-USER-PROFILE-SESSION-01 | 个人信息查看编辑与头像上传占位 | USER-PROFILE-001 | ["USER-PROFILE-002","SYS-AUTH-002"] | 已完成 | Codex | ["USER-AUTH-002"] | ["TC-ST-BATCH-20260409-USER-PROFILE-SESSION-01-01"] | ["BUG-BATCH-20260409-USER-PROFILE-SESSION-001"] | 2026-04-09 21:50 |
| ST-BATCH-20260409-USER-PROFILE-SESSION-02 | 修改密码页与重新登录承接 | USER-PROFILE-002 | ["USER-PROFILE-001","USER-AUTH-003"] | 已完成 | Codex | ["USER-PROFILE-001"] | ["TC-ST-BATCH-20260409-USER-PROFILE-SESSION-02-01"] | [] | 2026-04-09 21:50 |
| ST-BATCH-20260409-USER-PROFILE-SESSION-03 | 主动退出与被挤下线统一提示 | USER-AUTH-003 | ["SYS-AUTH-002"] | 已完成 | Codex | ["USER-AUTH-002"] | ["TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-01","TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-02"] | [] | 2026-04-09 21:50 |

## 子任务说明模板

### ST-BATCH-20260409-USER-PROFILE-SESSION-01 个人信息查看编辑与头像上传占位

```yaml
subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-01
title: 个人信息查看编辑与头像上传占位
primary_req_id: USER-PROFILE-001
related_req_ids: ["USER-PROFILE-002","SYS-AUTH-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-AUTH-002"]
implementation_notes:
  - 以 UserProfile Make 原型信息结构为参考，重做个人信息页布局和资料卡
  - 账户、手机号只读展示，昵称、邮箱可编辑
  - 头像上传采用本地预览和压缩占位提示，不接真实上传接口
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-PROFILE-SESSION-01-01"]
linked_bug_ids: ["BUG-BATCH-20260409-USER-PROFILE-SESSION-001"]
last_updated_at: 2026-04-09 21:50
```

Implementation Notes:

- 该子任务输出 profile mock 数据层，避免与账单预算 mock 混写。

### ST-BATCH-20260409-USER-PROFILE-SESSION-02 修改密码页与重新登录承接

```yaml
subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-02
title: 修改密码页与重新登录承接
primary_req_id: USER-PROFILE-002
related_req_ids: ["USER-PROFILE-001","USER-AUTH-003"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-PROFILE-001"]
implementation_notes:
  - 新增 /user/profile/password 页面并补齐旧密码、新密码、确认新密码校验
  - 使用 mock 密码校验逻辑模拟修改密码成功和失败分支
  - 修改成功后清空用户端登录态并跳回登录页提示重新登录
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-PROFILE-SESSION-02-01"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:50
```

Implementation Notes:

- 该子任务只负责前端表现层，不实现真实密码接口联调。

### ST-BATCH-20260409-USER-PROFILE-SESSION-03 主动退出与被挤下线统一提示

```yaml
subtask_id: ST-BATCH-20260409-USER-PROFILE-SESSION-03
title: 主动退出与被挤下线统一提示
primary_req_id: USER-AUTH-003
related_req_ids: ["SYS-AUTH-002"]
status: 已完成
owner_ai: Codex
dependencies: ["USER-AUTH-002"]
implementation_notes:
  - 登录页补充主动退出、会话失效、密码修改后的统一提示承接
  - 新增会话失效提示承接页，作为 SYS-AUTH-002 的前端表现层入口
  - 会话提示由本地状态和 query 参数驱动，不改真实后端互踢逻辑
linked_test_case_ids: ["TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-01","TC-ST-BATCH-20260409-USER-PROFILE-SESSION-03-02"]
linked_bug_ids: []
last_updated_at: 2026-04-09 21:50
```

Implementation Notes:

- 该子任务需要覆盖主动退出和被挤下线两个入口，提示文案保持一致。
