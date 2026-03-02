# 养老院管理系统 - 技术设计文档

**日期**: 2026-03-02
**版本**: v1.0

---

## 1. 项目概述

### 1.1 项目背景
单体养老院数字化转型项目，解决护理管理混乱、家属缺乏信任感、数据分散管理等问题。

### 1.2 目标用户
- **家属**: 查看父母状况、护理记录、在线缴费
- **老人**: 呼叫服务、查看活动
- **护工**: 记录护理、查看排班、处理请求
- **管理员**: 全面管理养老院运营

### 1.3 核心痛点
1. 护理管理难 - 护工工作记录混乱、护理标准不统一
2. 家属不放心 - 无法及时了解父母状况、缺乏信任感
3. 数据管理混乱 - 健康档案、用药记录、费用账单等数据分散

---

## 2. 技术架构

### 2.1 技术栈

| 层级 | 技术选择 | 说明 |
|------|----------|------|
| **小程序** | 微信小程序原生 | 性能最优，用户体验好 |
| **管理后台** | Vue 3 + Element Plus + Vite | 开发效率高，组件丰富 |
| **后端** | Go 1.21+ + Gin + GORM | 高性能、开发效率高 |
| **数据库** | PostgreSQL 15+ | 关系型、支持JSON、稳定性高 |
| **缓存** | Redis 7+ | 会话、热点数据缓存 |
| **文件存储** | 腾讯云 COS | 照片、视频存储 |
| **消息推送** | 微信模板消息 + 短信 | 实时通知 |
| **部署** | Docker Compose | 简化部署运维 |

### 2.2 系统架构图

```
┌─────────────────────────────────────────┐
│           Nginx (反向代理)                │
├─────────────────────────────────────────┤
│  小程序静态资源  │  管理后台静态资源      │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│        Go 后端服务 (Gin + GORM)          │
│  ┌─────────────────────────────────┐    │
│  │  认证中间件  │  日志  │  CORS   │    │
│  ├─────────────────────────────────┤    │
│  │         业务逻辑层               │    │
│  │  ┌─────┬─────┬─────┬─────┐      │    │
│  │  │用户 │护理 │财务 │消息 │      │    │
│  │  │模块 │模块 │模块 │模块 │      │    │
│  │  └─────┴─────┴─────┴─────┘      │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
          ↓           ↓           ↓
   ┌──────────┐  ┌──────────┐  ┌──────────┐
   │PostgreSQL│  │  Redis   │  │ OSS存储  │
   │  主数据库 │  │   缓存   │  │ 文件/图片│
   └──────────┘  └──────────┘  └──────────┘
```

---

## 3. 功能模块设计

### 3.1 小程序端

#### 3.1.1 家属端

| 模块 | 功能 | 优先级 |
|------|------|--------|
| **登录/注册** | 手机号+微信授权登录 | P0 |
| **老人概况** | 基本信息、入住状态、所在区域 | P0 |
| **护理记录** | 查看护理日志（带时间轴） | P0 |
| **照片视频** | 护工上传的日常生活照/视频 | P1 |
| **健康数据** | 血压、血糖、体温趋势图 | P1 |
| **在线缴费** | 查看账单、微信支付 | P1 |
| **探视预约** | 预约探视时间 | P2 |
| **消息通知** | 护理异常、账单、活动通知 | P0 |
| **家属留言** | 给老人/护工留言 | P2 |

#### 3.1.2 老人端

| 模块 | 功能 | 优先级 |
|------|------|--------|
| **服务呼叫** | 一键呼叫护工（分类） | P0 |
| **呼叫记录** | 历史呼叫、处理状态 | P1 |
| **查看活动** | 院内活动日程、报名 | P2 |

#### 3.1.3 护工端

| 模块 | 功能 | 优先级 |
|------|------|--------|
| **今日任务** | 排班任务、待护理老人列表 | P0 |
| **护理记录** | 快速记录（打钩+拍照） | P0 |
| **老人信息** | 健康档案、用药提醒 | P1 |
| **服务请求** | 接收/处理老人呼叫 | P1 |
| **排班查看** | 我的排班表 | P1 |

### 3.2 管理后台端

| 模块 | 功能 | 优先级 |
|------|------|--------|
| **登录/权限** | 管理员登录、角色权限管理 | P0 |
| **工作台** | 数据概览、待办事项 | P1 |
| **老人档案** | 入住登记、信息管理、家属关联 | P0 |
| **房间管理** | 楼栋/房间/床位管理，可视化 | P0 |
| **员工管理** | 护工信息、排班、考核 | P0 |
| **护理标准** | 护理项目配置、频次要求 | P0 |
| **护理记录** | 查询、统计、导出 | P0 |
| **服务请求** | 呼叫记录、响应统计 | P1 |
| **财务管理** | 费用配置、账单、收款 | P1 |
| **消息中心** | 群发通知、模板消息 | P1 |
| **数据报表** | 各类统计报表 | P2 |
| **系统设置** | 角色权限、操作日志 | P1 |

---

## 4. 数据库设计

### 4.1 核心表结构

```sql
-- 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE,
    password VARCHAR(255),
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    openid VARCHAR(100),        -- 微信openid
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 角色表
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE,    -- 家属/老人/护工/管理员
    description TEXT
);

-- 用户角色关联
CREATE TABLE user_roles (
    user_id BIGINT REFERENCES users(id),
    role_id INT REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- 老人档案
CREATE TABLE elderly (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    gender VARCHAR(10),
    birth_date DATE,
    id_card VARCHAR(20),
    phone VARCHAR(20),
    emergency_contact VARCHAR(50),
    emergency_phone VARCHAR(20),
    admission_date DATE,
    bed_id BIGINT,
    health_status TEXT,         -- JSON: 健康状况
    care_level INT,             -- 护理等级
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 楼栋
CREATE TABLE buildings (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    floors_count INT DEFAULT 1
);

-- 楼层
CREATE TABLE floors (
    id SERIAL PRIMARY KEY,
    building_id INT REFERENCES buildings(id),
    name VARCHAR(50) NOT NULL,
    sort_order INT DEFAULT 0
);

-- 房间
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    floor_id INT REFERENCES floors(id),
    name VARCHAR(50) NOT NULL,
    bed_capacity INT DEFAULT 1,
    sort_order INT DEFAULT 0
);

-- 床位
CREATE TABLE beds (
    id BIGSERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(id),
    name VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'empty', -- empty/occupied/reserved
    elderly_id BIGINT REFERENCES elderly(id),
    sort_order INT DEFAULT 0
);

-- 护理项目
CREATE TABLE care_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    category VARCHAR(30),       -- 喂饭/翻身/清洁/用药/其他
    unit VARCHAR(20),           -- 次/分钟
    description TEXT
);

-- 护理标准（老人-项目-频次）
CREATE TABLE care_standards (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    care_item_id INT REFERENCES care_items(id),
    frequency INT,              -- 每天次数
    priority INT DEFAULT 0
);

-- 护理记录
CREATE TABLE care_records (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    care_item_id INT REFERENCES care_items(id),
    staff_id BIGINT REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'completed',
    notes TEXT,
    images TEXT,                -- JSON: 图片URL数组
    recorded_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 服务请求
CREATE TABLE service_requests (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    type VARCHAR(30),           -- 护理/送餐/打扫/其他
    status VARCHAR(20) DEFAULT 'pending', -- pending/processing/completed
    requester_id BIGINT REFERENCES users(id),
    handler_id BIGINT REFERENCES users(id),
    notes TEXT,
    requested_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- 费用项目
CREATE TABLE fee_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    type VARCHAR(30),           -- 床位费/护理费/餐费/其他
    cycle VARCHAR(20)           -- daily/monthly/once
);

-- 账单
CREATE TABLE bills (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    bill_no VARCHAR(50) UNIQUE,
    total_amount DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'unpaid', -- unpaid/paid/overdue
    bill_period_start DATE,
    bill_period_end DATE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 账单明细
CREATE TABLE bill_items (
    id BIGSERIAL PRIMARY KEY,
    bill_id BIGINT REFERENCES bills(id),
    fee_item_id INT REFERENCES fee_items(id),
    quantity INT DEFAULT 1,
    amount DECIMAL(10,2),
    description TEXT
);

-- 支付记录
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    bill_id BIGINT REFERENCES bills(id),
    amount DECIMAL(10,2),
    method VARCHAR(30),         -- wechat/cash/transfer
    transaction_no VARCHAR(100),
    paid_at TIMESTAMP DEFAULT NOW()
);

-- 健康记录
CREATE TABLE health_records (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    type VARCHAR(30),           -- blood_pressure/blood_sugar/temperature/weight
    value VARCHAR(50),
    unit VARCHAR(20),
    recorded_by BIGINT REFERENCES users(id),
    recorded_at TIMESTAMP DEFAULT NOW()
);

-- 用药记录
CREATE TABLE medication_records (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    medicine_name VARCHAR(100),
    dosage VARCHAR(50),
    frequency VARCHAR(50),
    taken_at TIMESTAMP,
    recorded_by BIGINT REFERENCES users(id)
);

-- 员工信息
CREATE TABLE staff (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name VARCHAR(50) NOT NULL,
    position VARCHAR(30),       -- 护工/护士/医生/其他
    department VARCHAR(50),
    hire_date DATE,
    status VARCHAR(20) DEFAULT 'active'
);

-- 排班
CREATE TABLE schedules (
    id BIGSERIAL PRIMARY KEY,
    staff_id BIGINT REFERENCES staff(id),
    date DATE,
    shift_type VARCHAR(20),     -- morning/afternoon/night
    created_at TIMESTAMP DEFAULT NOW()
);

-- 消息通知
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    type VARCHAR(30),
    title VARCHAR(100),
    content TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 家属关联
CREATE TABLE elderly_family (
    id BIGSERIAL PRIMARY KEY,
    elderly_id BIGINT REFERENCES elderly(id),
    user_id BIGINT REFERENCES users(id),
    relation VARCHAR(30),       -- 子女/配偶/亲属
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 操作日志
CREATE TABLE operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    action VARCHAR(50),
    resource_type VARCHAR(50),
    resource_id BIGINT,
    details TEXT,
    ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 5. 接口设计

### 5.1 认证相关

```
POST   /api/auth/register          # 注册
POST   /api/auth/login             # 登录
POST   /api/auth/logout            # 登出
POST   /api/auth/refresh           # 刷新token
POST   /api/auth/wechat-login      # 微信登录
```

### 5.2 用户相关

```
GET    /api/user/profile           # 获取个人信息
PUT    /api/user/profile           # 更新个人信息
GET    /api/user/elderly-list      # 家属获取关联老人列表
```

### 5.3 老人相关

```
GET    /api/elderly                # 老人列表
GET    /api/elderly/:id            # 老人详情
POST   /api/elderly                # 创建老人档案
PUT    /api/elderly/:id            # 更新老人档案
DELETE /api/elderly/:id            # 删除老人档案
GET    /api/elderly/:id/health     # 健康记录
```

### 5.4 护理相关

```
GET    /api/care/items             # 护理项目列表
GET    /api/care/records           # 护理记录列表
POST   /api/care/records           # 创建护理记录
GET    /api/care/records/:id       # 护理记录详情
GET    /api/care/my-tasks          # 护工今日任务
```

### 5.5 服务请求

```
GET    /api/service/requests       # 服务请求列表
POST   /api/service/requests       # 创建服务请求
PUT    /api/service/requests/:id   # 处理服务请求
```

### 5.6 财务相关

```
GET    /api/bills                  # 账单列表
GET    /api/bills/:id              # 账单详情
POST   /api/bills/:id/pay          # 支付账单
GET    /api/payments               # 支付记录
```

### 5.7 房间管理

```
GET    /api/rooms                  # 房间列表
GET    /api/beds                   # 床位列表
PUT    /api/beds/:id               # 更新床位状态
```

---

## 6. 项目结构

```
elderly-care-system/
├── backend/                       # Go 后端
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/               # 配置
│   │   ├── middleware/           # 中间件
│   │   ├── model/                # 数据模型
│   │   ├── repository/           # 数据访问层
│   │   ├── service/              # 业务逻辑层
│   │   └── handler/              # HTTP 处理器
│   ├── pkg/                      # 公共包
│   ├── migrations/               # 数据库迁移
│   ├── go.mod
│   └── go.sum
│
├── admin-frontend/               # 管理后台 Vue3
│   ├── src/
│   │   ├── api/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── router/
│   │   ├── store/
│   │   ├── views/
│   │   ├── App.vue
│   │   └── main.ts
│   ├── package.json
│   └── vite.config.ts
│
├── miniprogram/                  # 微信小程序
│   ├── pages/
│   ├── components/
│   ├── utils/
│   ├── app.ts
│   └── app.json
│
├── docs/                         # 文档
│   └── plans/
│
└── docker-compose.yml            # 本地开发环境
```

---

## 7. 安全设计

1. **认证**: JWT Token + Redis 存储黑名单
2. **密码**: bcrypt 加密
3. **权限**: RBAC 角色权限控制
4. **API 限流**: Redis + 滑动窗口
5. **数据脱敏**: 敏感字段加密存储
6. **SQL 注入防护**: 使用 GORM 参数化查询
7. **XSS 防护**: 前端输入过滤 + CSP 头
8. **操作日志**: 记录关键操作审计

---

## 8. 开发计划

### 第一阶段（MVP - 4周）

**Week 1-2: 基础框架**
- 项目初始化
- 数据库设计 + 迁移脚本
- 用户认证 + JWT
- 管理后台框架搭建

**Week 3: 核心功能**
- 老人档案管理
- 房间床位管理
- 护理记录功能
- 家属端查看功能

**Week 4: 小程序开发**
- 家属端小程序
- 护工端小程序
- 基础消息推送

### 第二阶段（2周）

- 服务呼叫功能
- 在线缴费
- 护工排班
- 健康记录

### 第三阶段（2周）

- 数据报表
- 探视预约
- 系统优化
- 测试上线

---

## 9. 部署方案

### 9.1 云服务选型

- **服务器**: 腾讯云 CVM 2核4G 起步
- **数据库**: 云 PostgreSQL 或自建
- **缓存**: 云 Redis
- **存储**: 腾讯云 COS
- **CDN**: 腾讯云 CDN

### 9.2 本地开发

```bash
docker-compose up -d
```

启动 PostgreSQL + Redis + 后端服务

---

**文档版本**: v1.0
**最后更新**: 2026-03-02
