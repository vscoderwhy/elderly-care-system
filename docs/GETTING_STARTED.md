# 快速开始指南

## 系统要求

- Go 1.21+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (可选)

## 方式一：Docker 快速启动（推荐）

### 1. 克隆项目

```bash
cd ~/elderly-care-system
```

### 2. 启动服务

```bash
docker-compose up -d
```

### 3. 执行数据库迁移

```bash
make migrate
```

### 4. 访问系统

- 管理后台: http://localhost:3000
- 后端 API: http://localhost:8080

### 5. 默认账号

- 手机号: `13800138000`
- 密码: `123456`

---

## 方式二：手动启动

### 1. 安装 Go

```bash
# OpenCloudOS/CentOS
sudo yum install -y golang

# 或手动安装
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. 启动数据库

```bash
# PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=elderly_care \
  -p 5432:5432 \
  postgres:15-alpine

# Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

### 3. 启动后端

```bash
cd backend

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 运行
go run cmd/server/main.go
```

### 4. 启动前端

```bash
cd admin-frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

---

## 使用 Makefile

```bash
# 查看所有命令
make help

# 一键初始化
make init

# 编译后端
make build

# 运行后端
make run

# 运行测试
make test

# 数据库迁移
make migrate

# 启动 Docker 服务
make docker-up

# 停止 Docker 服务
make docker-down

# 查看日志
make docker-logs
```

---

## 项目结构

```
elderly-care-system/
├── backend/              # Go 后端
│   ├── cmd/             # 入口文件
│   ├── internal/        # 业务代码
│   ├── pkg/             # 公共包
│   └── migrations/      # 数据库迁移
├── admin-frontend/      # Vue 3 管理后台
├── miniprogram/         # 微信小程序
├── docs/               # 文档
├── scripts/            # 脚本
├── docker-compose.yml  # Docker 编排
├── Makefile           # 便捷命令
└── README.md          # 项目说明
```

---

## 开发指南

### 添加新 API

1. 在 `backend/internal/model/` 添加数据模型
2. 在 `backend/internal/repository/` 添加数据访问层
3. 在 `backend/internal/service/` 添加业务逻辑层
4. 在 `backend/internal/handler/` 添加 HTTP 处理器
5. 在 `backend/cmd/server/main.go` 注册路由

### 添加新页面

1. 在 `admin-frontend/src/views/` 创建页面组件
2. 在 `admin-frontend/src/router/index.ts` 添加路由
3. 在 `admin-frontend/src/api/` 添加 API 调用

---

## 常见问题

### Q: 数据库连接失败？
A: 检查 PostgreSQL 是否启动，端口是否正确

### Q: Redis 连接失败？
A: 检查 Redis 是否启动，端口是否正确

### Q: 前端无法调用后端 API？
A: 检查 `vite.config.ts` 中的代理配置

### Q: 登录失败？
A: 使用默认账号 `13800138000` / `123456`，或先注册新用户

---

## 下一步

- [ ] 查看完整 [API 文档](./API.md)
- [ ] 查看 [设计文档](./plans/2026-03-02-elderly-care-system-design.md)
- [ ] 开始开发新功能

---

## 获取帮助

- 提交 Issue
- 查看项目文档
- 联系开发团队
