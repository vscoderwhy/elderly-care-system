# 养老院管理系统

<div align="center">

![Version](https://img.shields.io/badge/version-v1.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql)
![License](https://img.shields.io/badge/license-MIT-green)

一个现代化的养老院管理系统，支持管理后台、微信小程序（家属端/护工端）

[快速开始](#快速开始) • [功能特性](#功能特性) • [技术栈](#技术栈) • [文档](#文档)

</div>

---

## 项目简介

养老院管理系统是基于 **Go + Vue 3 + 微信小程序** 开发的综合性养老院管理平台，旨在解决养老院护理管理混乱、家属缺乏信任感、数据分散管理等问题。

### 核心价值

- 📊 **护理数字化** - 规范护理流程，记录可追溯
- 👨‍👩‍👧 **家属透明化** - 实时了解父母状况，增强信任
- 💰 **财务管理** - 账单清晰，在线缴费
- 📱 **移动便捷** - 小程序随时随地查看

---

## 功能特性

### 🖥️ 管理后台

| 模块 | 功能 |
|------|------|
| **工作台** | 数据概览、待办事项、快捷操作 |
| **老人管理** | 档案管理、入住登记、床位分配 |
| **房间管理** | 楼栋/楼层/房间/床位可视化 |
| **护理记录** | 护理项目配置、记录查询、数据统计 |
| **财务管理** | 费用配置、账单生成、收款记录 |
| **员工管理** | 护工信息、排班管理、绩效考核 |
| **数据报表** | 入住率、护理工作量、费用统计 |

### 📱 家属端小程序

- 👀 查看老人基本信息和状况
- 📝 查看护理记录（含照片/视频）
- 📊 查看健康数据趋势
- 💳 在线缴费账单
- 📅 预约探视时间
- 💬 与护工/管理员沟通

### 👨‍⚕️ 护工端小程序

- ✅ 查看今日护理任务
- 📝 快速记录护理（打钩+拍照）
- 📋 查看老人健康档案
- 🔔 处理老人服务呼叫
- 📅 查看个人排班

---

## 技术栈

### 后端技术

```
Go 1.21+          - 后端语言
Gin               - Web 框架
GORM              - ORM 框架
PostgreSQL 15+    - 关系型数据库
Redis 7+          - 缓存
JWT               - 认证
Docker            - 容器化
```

### 前端技术

```
Vue 3             - 前端框架
TypeScript        - 类型系统
Vite              - 构建工具
Element Plus      - UI 组件库
Pinia             - 状态管理
Axios             - HTTP 客户端
```

### 小程序技术

```
微信小程序原生    - 小程序框架
TypeScript        - 类型系统
```

---

## 快速开始

### 使用 Docker（推荐）

```bash
# 1. 进入项目目录
cd ~/elderly-care-system

# 2. 启动所有服务
docker-compose up -d

# 3. 执行数据库迁移
make migrate

# 4. 访问系统
# 管理后台: http://localhost:3000
# 后端 API: http://localhost:8080
```

### 使用 Makefile

```bash
# 查看所有命令
make help

# 一键初始化
make init

# 启动开发环境
make dev
```

### 默认账号

```
手机号: 13800138000
密码: 123456
```

---

## 项目结构

```
elderly-care-system/
├── backend/                # Go 后端服务
│   ├── cmd/               # 应用入口
│   ├── internal/          # 内部代码
│   │   ├── handler/       # HTTP 处理器
│   │   ├── service/       # 业务逻辑
│   │   ├── repository/    # 数据访问
│   │   ├── model/         # 数据模型
│   │   ├── middleware/    # 中间件
│   │   └── config/        # 配置
│   ├── pkg/               # 公共包
│   ├── migrations/        # 数据库迁移
│   └── Dockerfile
├── admin-frontend/        # Vue 3 管理后台
│   ├── src/
│   │   ├── views/        # 页面组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理
│   │   └── api/          # API 封装
│   └── Dockerfile
├── miniprogram/           # 微信小程序
│   ├── pages/            # 页面
│   ├── components/       # 组件
│   └── utils/            # 工具
├── docs/                 # 项目文档
├── scripts/              # 部署脚本
├── docker-compose.yml    # Docker 编排
├── Makefile             # 便捷命令
└── README.md            # 项目说明
```

---

## 文档

- [快速开始指南](./docs/GETTING_STARTED.md)
- [API 文档](./docs/API.md)
- [技术设计文档](./docs/plans/2026-03-02-elderly-care-system-design.md)

---

## 数据库设计

### 核心表

```
users              # 用户表
elderly            # 老人档案
elderly_family     # 家属关联
care_records       # 护理记录
care_items         # 护理项目
service_requests   # 服务请求
bills              # 账单
payments           # 支付记录
health_records     # 健康记录
buildings          # 楼栋
floors             # 楼层
rooms              # 房间
beds               # 床位
staff              # 员工
schedules          # 排班
notifications      # 消息通知
operation_logs     # 操作日志
```

---

## 开发计划

### ✅ 第一阶段 (MVP)

- [x] 项目架构设计
- [x] 数据库设计
- [x] 用户认证
- [x] 管理后台框架
- [ ] 老人档案管理
- [ ] 护理记录功能
- [ ] 家属端小程序

### 🔄 第二阶段

- [ ] 服务呼叫功能
- [ ] 在线缴费
- [ ] 护工排班
- [ ] 健康记录

### 📋 第三阶段

- [ ] 数据报表
- [ ] 消息推送
- [ ] 系统优化

---

## 贡献指南

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 许可证

[MIT](LICENSE)

---

## 联系方式

- 项目地址: [GitHub](https://github.com/your-username/elderly-care-system)
- 问题反馈: [Issues](https://github.com/your-username/elderly-care-system/issues)

---

<div align="center">

**如果这个项目对你有帮助，请给一个 ⭐️**

Made with ❤️ by Claude Code

</div>
