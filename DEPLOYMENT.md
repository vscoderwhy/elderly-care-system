# 养老院管理系统部署文档

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                         负载均衡                             │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│  Nginx (前端)   │   │  Nginx (前端)   │   │  Nginx (移动端) │
│  管理后台       │   │  管理后台       │   │  uniapp H5     │
└────────────────┘   └────────────────┘   └────────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│  Go API (后端)  │   │  Go API (后端)  │   │  Go API (后端)  │
│  端口: 8080     │   │  端口: 8080     │   │  端口: 8080     │
└────────────────┘   └────────────────┘   └────────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
      ┌───────▼──────┐ ┌─────▼─────┐ ┌──────▼──────┐
      │  PostgreSQL  │ │   Redis   │ │ MinIO/OSS  │
      │   端口:5432  │ │ 端口:6379 │ │  文件存储   │
      └──────────────┘ └───────────┘ └─────────────┘
```

## 部署方式

### 方式一：Docker Compose (推荐用于开发/测试)

#### 1. 环境要求
- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ 内存
- 20GB+ 磁盘空间

#### 2. 部署步骤

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/elderly-care-system.git
cd elderly-care-system

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，修改数据库密码等配置

# 3. 启动所有服务
docker-compose up -d

# 4. 查看服务状态
docker-compose ps

# 5. 查看日志
docker-compose logs -f
```

#### 3. 服务访问
- 管理后台: http://localhost:8080
- 移动端H5: http://localhost:8081
- API文档: http://localhost:8080/swagger
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

### 方式二：Kubernetes (推荐用于生产环境)

#### 1. 创建命名空间
```bash
kubectl create namespace elderly-care
```

#### 2. 部署数据库
```bash
kubectl apply -f k8s/postgresql.yaml
kubectl apply -f k8s/redis.yaml
```

#### 3. 部署后端服务
```bash
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/backend-service.yaml
```

#### 4. 部署前端服务
```bash
kubectl apply -f k8s/admin-frontend-deployment.yaml
kubectl apply -f k8s/admin-frontend-service.yaml
kubectl apply -f k8s/mobile-frontend-deployment.yaml
```

#### 5. 配置 Ingress
```bash
kubectl apply -f k8s/ingress.yaml
```

### 方式三：传统部署

#### 后端部署 (Go)

```bash
# 1. 编译后端
cd backend
go build -o elderly-care-api ./cmd/main.go

# 2. 配置 Systemd 服务
sudo cp scripts/elderly-care-api.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable elderly-care-api
sudo systemctl start elderly-care-api

# 3. 配置 Nginx 反向代理
sudo cp config/nginx.conf /etc/nginx/sites-available/elderly-care-api
sudo ln -s /etc/nginx/sites-available/elderly-care-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

#### 前端部署 (Vue + uniapp)

```bash
# 1. 构建管理后台
cd admin-frontend
npm run build

# 2. 部署到 Nginx
sudo cp -r dist/* /var/www/admin-elderly-care/

# 3. 构建移动端 H5
cd ../uniapp-elderly-care
npm run build:h5

# 4. 部署到 Nginx
sudo cp -r dist/build/h5/* /var/www/mobile-elderly-care/
```

## 配置说明

### 后端配置 (config.yaml)

```yaml
server:
  port: 8080
  mode: production

database:
  host: localhost
  port: 5432
  user: elderly_care
  password: your_password
  dbname: elderly_care_db
  sslmode: disable

redis:
  addr: localhost:6379
  password: ""
  db: 0

jwt:
  secret: your_jwt_secret_key
  expireHours: 168

upload:
  maxFileSize: 10485760
  allowedTypes:
    - image/jpeg
    - image/png
    - application/pdf
```

### 环境变量 (.env)

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=elderly_care
DB_PASSWORD=your_password
DB_NAME=elderly_care_db

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your_jwt_secret_key

# 上传
UPLOAD_PATH=/uploads
MAX_FILE_SIZE=10485760

# 微信小程序 (用于移动端)
WECHAT_APP_ID=your_app_id
WECHAT_APP_SECRET=your_app_secret

# 支付配置
WECHAT_PAY_MCH_ID=your_mch_id
WECHAT_PAY_API_KEY=your_api_key
ALIPAY_APP_ID=your_app_id
ALIPAY_PRIVATE_KEY=your_private_key
```

## 数据库初始化

```bash
# 1. 连接到 PostgreSQL
psql -U postgres

# 2. 创建数据库和用户
CREATE DATABASE elderly_care_db;
CREATE USER elderly_care WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE elderly_care_db TO elderly_care;

# 3. 运行迁移脚本
psql -U elderly_care -d elderly_care_db -f scripts/init.sql

# 4. 初始化数据
psql -U elderly_care -d elderly_care_db -f scripts/seed.sql
```

## 监控和日志

### 日志查看

```bash
# Docker
docker-compose logs -f backend
docker-compose logs -f admin-frontend

# Kubernetes
kubectl logs -f deployment/backend

# Systemd
journalctl -u elderly-care-api -f
```

### 性能监控

- Prometheus: http://your-server:9090
- Grafana: http://your-server:3000

默认 Grafana 登录:
- 用户名: admin
- 密码: admin

## 备份和恢复

### 数据库备份

```bash
# 备份
pg_dump -U elderly_care elderly_care_db > backup_$(date +%Y%m%d).sql

# 恢复
psql -U elderly_care elderly_care_db < backup_20260303.sql
```

### 文件备份

```bash
# 备份上传文件
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz uploads/
```

## 常见问题

### 1. 后端启动失败
- 检查数据库连接是否正常
- 确认配置文件中的端口未被占用
- 查看日志: `docker-compose logs backend`

### 2. 前端页面空白
- 确认后端 API 是否可访问
- 检查浏览器控制台错误信息
- 确认 API 地址配置正确

### 3. 文件上传失败
- 检查上传目录权限
- 确认文件大小限制
- 检查磁盘空间

### 4. 数据库连接失败
- 确认数据库服务已启动
- 检查连接参数是否正确
- 查看防火墙规则

## 安全建议

1. **修改默认密码**: 部署后立即修改数据库、Redis、JWT等默认密码
2. **启用 HTTPS**: 生产环境必须使用 HTTPS
3. **限制访问**: 配置防火墙，只开放必要端口
4. **定期备份**: 每日自动备份数据库和文件
5. **更新补丁**: 定期更新系统和依赖包
6. **日志审计**: 保存并定期审查访问日志

## 性能优化

1. **数据库优化**
   - 创建适当索引
   - 定期 VACUUM
   - 配置连接池

2. **缓存优化**
   - 启用 Redis 缓存
   - 配置缓存过期时间
   - 使用连接池

3. **前端优化**
   - 启用 Gzip 压缩
   - 配置 CDN
   - 使用懒加载

## 扩展部署

### 水平扩展

```bash
# 增加后端实例
kubectl scale deployment/backend --replicas=3

# 增加前端实例
kubectl scale deployment/admin-frontend --replicas=2
```

### 负载均衡

配置 Nginx 或云负载均衡器，实现多实例负载分发。

## 联系支持

- 技术支持邮箱: support@elderly-care.com
- 问题反馈: https://github.com/your-repo/elderly-care-system/issues
