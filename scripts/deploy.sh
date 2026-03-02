#!/bin/bash

# 部署脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始部署养老院管理系统...${NC}"

# 检查环境
if [ ! -f .env ]; then
    echo -e "${YELLOW}警告: .env 文件不存在，使用默认配置${NC}"
fi

# 构建后端
echo -e "${GREEN}构建后端...${NC}"
cd backend
go mod download
go build -o bin/server ./cmd/server
cd ..

# 构建前端
echo -e "${GREEN}构建前端...${NC}"
cd admin-frontend
npm install
npm run build
cd ..

# 停止旧容器
echo -e "${GREEN}停止旧容器...${NC}"
docker-compose down

# 启动新容器
echo -e "${GREEN}启动新容器...${NC}"
docker-compose up -d

# 等待服务启动
echo -e "${YELLOW}等待服务启动...${NC}"
sleep 10

# 执行数据库迁移
echo -e "${GREEN}执行数据库迁移...${NC}"
docker exec -i elderly-care-db psql -U postgres -d elderly_care < backend/migrations/000001_init_schema.up.sql || true

echo -e "${GREEN}部署完成！${NC}"
echo -e "后端地址: http://localhost:8080"
echo -e "前端地址: http://localhost:3000"
echo -e ""
echo -e "默认账号: 13800138000 / 123456"
