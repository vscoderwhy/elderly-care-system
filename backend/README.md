# 养老院管理系统 - 后端 API

基于 Go + Gin + GORM + PostgreSQL 的养老院管理系统后端服务。

## 技术栈

- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT
- **WebSocket**: Gorilla WebSocket
- **配置**: Viper

## 项目结构

```
backend/
├── cmd/
│   └── main.go              # 应用入口
├── internal/
│   ├── api/
│   │   ├── setup.go          # 路由设置
│   │   └── handler/
│   │       └── handler.go   # 请求处理
│   ├── config/
│   │   └── config.go         # 配置管理
│   ├── middleware/
│   │   └── middleware.go     # 中间件
│   ├── models/
│   │   └── models.go         # 数据模型
│   ├── service/             # 业务逻辑层
│   └── handler/             # 处理器层
├── pkg/                     # 公共包
├── config.yaml              # 配置文件
├── go.mod                   # Go 模块文件
└── README.md                # 项目文档
```

## 功能特性

### 已实现功能

- ✅ 项目结构搭建
- ✅ 配置管理
- ✅ 数据库模型定义
- ✅ JWT 认证中间件
- ✅ CORS 跨域中间件
- ✅ 日志中间件
- ✅ 恢复中间件
- ✅ API 路由框架
- ✅ Handler 基础结构

### 待实现功能

- ⏳ 用户认证逻辑
- ⏳ 业务层实现
- ⏳ WebSocket 实时通信
- ⏳ 文件上传功能
- ⏳ 数据导出功能
- ⏳ 单元测试
- ⏳ API 文档

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+ (可选)

### 安装依赖

```bash
go mod download
```

### 配置数据库

1. 创建数据库:
```sql
CREATE DATABASE elderly_care;
```

2. 修改 `config.yaml` 中的数据库配置

### 运行服务

```bash
# 开发模式
go run cmd/main.go

# 编译后运行
go build -o elderly-care cmd/main.go
./elderly-care
```

服务将在 `http://localhost:8080` 启动

## API 文档

### 认证接口

#### 登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

#### 登出
```
POST /api/v1/auth/logout
Authorization: Bearer {token}
```

### 用户管理

#### 获取当前用户信息
```
GET /api/v1/users/me
Authorization: Bearer {token}
```

#### 获取用户列表 (管理员)
```
GET /api/v1/users
Authorization: Bearer {token}
```

### 老人管理

#### 获取老人列表
```
GET /api/v1/elderly?page=1&pageSize=20&name=&careLevel=
Authorization: Bearer {token}
```

#### 创建老人
```
POST /api/v1/elderly
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "张奶奶",
  "gender": "女",
  "birthDate": "1946-01-01T00:00:00Z",
  "idCard": "110101194601011234",
  "careLevel": "level2",
  "bedNumber": "3号楼201",
  "checkInDate": "2023-06-15T00:00:00Z"
}
```

### 护理管理

#### 获取护理任务列表
```
GET /api/v1/care/tasks?status=&date=
Authorization: Bearer {token}
```

#### 开始任务
```
POST /api/v1/care/tasks/:id/start
Authorization: Bearer {token}
```

#### 完成任务
```
POST /api/v1/care/tasks/:id/complete
Authorization: Bearer {token}
Content-Type: application/json

{
  "result": "任务完成，老人状态良好"
}
```

### 健康数据

#### 获取最新健康数据
```
GET /api/v1/health/latest/:elderlyId
Authorization: Bearer {token}
```

#### 获取健康趋势
```
GET /api/v1/health/trend?elderlyId=1&type=bloodPressure&days=7
Authorization: Bearer {token}
```

### 财务管理

#### 获取账单列表
```
GET /api/v1/finance/bills?elderlyId=&status=&month=
Authorization: Bearer {token}
```

#### 支付账单
```
POST /api/v1/finance/bills/:id/pay
Authorization: Bearer {token}
Content-Type: application/json

{
  "paymentMethod": "wechat"
}
```

## 数据模型

### User (用户)
- 用户名、密码、姓名、头像
- 角色: admin, nurse, family
- 状态: active, inactive

### Elderly (老人)
- 基本信息、护理等级、床位号
- 入住日期、健康状况
- 状态: active, leave, hospital, discharged

### CareRecord (护理记录)
- 护理类型、内容、结果
- 评价、图片、执行人

### CareTask (护理任务)
- 任务类型、标题、内容
- 计划时间、执行时间、优先级
- 状态、进度

### HealthData (健康数据)
- 血压、血糖、心率、体温等
- 记录时间、记录人

### Bill (账单)
- 账单类型、金额
- 账单日期、应付日期
- 支付状态、支付方式

### VisitAppointment (探视预约)
- 访客信息、预约时间
- 状态: pending, approved, completed, cancelled

### Notification (消息通知)
- 标题、内容、类型
- 已读状态

## 开发规范

### 错误处理

所有 API 响应遵循统一格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- `code`: 0 表示成功，非 0 表示错误
- `message`: 描述信息
- `data`: 返回数据

### 分页参数

```
page: 页码，从 1 开始
pageSize: 每页数量
```

### 状态码

- 200: 成功
- 400: 参数错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器错误

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t elderly-care-backend .

# 运行容器
docker run -p 8080:8080 \
  -e DB_HOST=database \
  -e DB_PASSWORD=password \
  elderly-care-backend
```

### 环境变量

支持通过环境变量覆盖配置：

```bash
export SERVER_PORT=8080
export DATABASE_HOST=localhost
export DATABASE_PASSWORD=password
export JWT_SECRET=your-secret-key
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...
```

## License

MIT License
