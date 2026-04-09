# Migration Convention

后续数据库变更统一放在本目录，命名建议如下：

- `0003_add_xxx.sql`
- `0004_alter_xxx.sql`

执行原则：

1. `database/init` 只存放首版初始化脚本。
2. 新增表、字段、索引或数据修复脚本都放到 `database/migrations`。
3. 每个脚本只做一类变更，方便追踪和回滚。
4. 脚本文件名必须递增，避免多人协作时冲突。

