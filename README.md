# billSoftware

## 目录结构

- `frontend`：Vue 3 + Element Plus + Axios + Webpack + ECharts + Babel 6 前端工程
- `backend`：Golang + Gin + XORM 后端工程
- `database`：数据库初始化与后续迁移脚本
- `docs`：需求与项目文档

## 快速开始

### 前端

```bash
cd frontend
npm install
npm run dev
```

默认访问地址：

- 用户端：`http://localhost:9000/user/login`
- 管理端：`http://localhost:9000/admin/login`

### 后端

先复制一份本地配置：

```bash
cd backend
cp configs/app.example.yaml configs/app.yaml
```

按实际环境修改 `backend/configs/app.yaml` 中的数据库连接后，再执行：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

默认服务地址：

- 后端：`http://localhost:8080`

### 数据库

数据库脚本位于 `database` 目录：

- `database/init/0000_create_database.sql`
- `database/init/0001_schema.sql`
- `database/init/0002_seed_default_categories.sql`

初始化顺序：

```bash
mysql -h qiqiqu.cn -P 3306 -u root -p < database/init/0000_create_database.sql
mysql -h qiqiqu.cn -P 3306 -u root -p billSoftware < database/init/0001_schema.sql
mysql -h qiqiqu.cn -P 3306 -u root -p billSoftware < database/init/0002_seed_default_categories.sql
```
