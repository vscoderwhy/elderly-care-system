#!/bin/bash

# 养老院管理系统启动脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  养老院管理系统 - 启动脚本${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 检查 Go
echo -e "${YELLOW}检查 Go 环境...${NC}"
if [ -f "$HOME/go-install/go/bin/go" ]; then
    export PATH=$PATH:$HOME/go-install/go/bin
    echo -e "${GREEN}✓ Go 已安装${NC}"
    go version
else
    echo -e "${RED}✗ Go 未安装${NC}"
    exit 1
fi

# 获取服务器 IP
echo ""
echo -e "${YELLOW}获取服务器 IP 地址...${NC}"
SERVER_IP=$(hostname -I | awk '{print $1}')
if [ -z "$SERVER_IP" ]; then
    SERVER_IP="localhost"
fi
echo -e "${GREEN}✓ 服务器 IP: ${SERVER_IP}${NC}"

# 启动 PostgreSQL
echo ""
echo -e "${YELLOW}启动 PostgreSQL...${NC}"
if ! docker ps | grep -q "postgres"; then
    docker run -d \
        --name elderly-care-db \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB=elderly_care \
        -p 5432:5432 \
        postgres:15-alpine 2>/dev/null || echo "PostgreSQL 可能已运行"
    echo -e "${GREEN}✓ PostgreSQL 已启动${NC}"
    sleep 3
else
    echo -e "${GREEN}✓ PostgreSQL 已在运行${NC}"
fi

# 启动 Redis
echo ""
echo -e "${YELLOW}启动 Redis...${NC}"
if ! docker ps | grep -q "redis"; then
    docker run -d \
        --name elderly-care-redis \
        -p 6379:6379 \
        redis:7-alpine 2>/dev/null || echo "Redis 可能已运行"
    echo -e "${GREEN}✓ Redis 已启动${NC}"
else
    echo -e "${GREEN}✓ Redis 已在运行${NC}"
fi

# 等待数据库就绪
echo ""
echo -e "${YELLOW}等待数据库就绪...${NC}"
sleep 2

# 执行数据库迁移
echo ""
echo -e "${YELLOW}执行数据库迁移...${NC}"
docker exec -i elderly-care-db psql -U postgres -d elderly_care < backend/migrations/000001_init_schema.up.sql 2>/dev/null && echo -e "${GREEN}✓ 数据库迁移完成${NC}" || echo -e "${YELLOW}! 数据库可能已初始化${NC}"

# 启动后端
echo ""
echo -e "${YELLOW}启动后端服务...${NC}"
cd backend
nohup go run cmd/server/main.go > /tmp/elderly-backend.log 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > /tmp/elderly-backend.pid
echo -e "${GREEN}✓ 后端已启动 (PID: $BACKEND_PID)${NC}"
cd ..

# 等待后端启动
echo -e "${YELLOW}等待后端服务就绪...${NC}"
sleep 3

# 启动移动端
echo ""
echo -e "${YELLOW}启动移动端 H5...${NC}"
cd mobile-frontend
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}安装依赖...${NC}"
    npm install --silent
fi
nohup npm run dev > /tmp/elderly-mobile.log 2>&1 &
MOBILE_PID=$!
echo $MOBILE_PID > /tmp/elderly-mobile.pid
echo -e "${GREEN}✓ 移动端已启动 (PID: $MOBILE_PID)${NC}"
cd ..

# 等待移动端启动
echo -e "${YELLOW}等待移动端服务就绪...${NC}"
sleep 3

# 启动 Nginx
echo ""
echo -e "${YELLOW}启动 Nginx...${NC}"
if command -v nginx &> /dev/null; then
    sudo cp nginx.conf /etc/nginx/conf.d/elderly-care.conf 2>/dev/null || sudo cp nginx.conf /etc/nginx/nginx.conf
    sudo nginx -s reload 2>/dev/null || sudo nginx
    echo -e "${GREEN}✓ Nginx 已启动${NC}"
else
    echo -e "${YELLOW}! Nginx 未安装，请访问以下地址：${NC}"
fi

# 完成
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  系统启动完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${GREEN}访问地址：${NC}"
echo -e "  📱 手机访问: ${GREEN}http://${SERVER_IP}${NC}"
echo -e "  💻 管理后台: ${GREEN}http://${SERVER_IP}:3000${NC}"
echo -e "  🔧 后端 API: ${GREEN}http://${SERVER_IP}:8080${NC}"
echo ""
echo -e "${GREEN}默认账号：${NC}"
echo -e "  手机号: ${YELLOW}13800138000${NC}"
echo -e "  密码: ${YELLOW}123456${NC}"
echo ""
echo -e "${YELLOW}日志文件：${NC}"
echo -e "  后端: ${GREEN}tail -f /tmp/elderly-backend.log${NC}"
echo -e "  移动端: ${GREEN}tail -f /tmp/elderly-mobile.log${NC}"
echo ""
echo -e "${YELLOW}停止服务：${NC}"
echo -e "  ${GREEN}./stop.sh${NC}"
echo ""

# 保存 PID
echo $BACKEND_PID > /tmp/elderly-backend.pid
echo $MOBILE_PID > /tmp/elderly-mobile.pid
