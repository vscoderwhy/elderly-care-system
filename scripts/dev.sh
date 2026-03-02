#!/bin/bash

# 开发环境启动脚本

set -e

echo "启动开发环境..."

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

# 启动数据库服务
echo "启动数据库服务..."
docker-compose up -d postgres redis

# 等待数据库就绪
echo "等待数据库就绪..."
sleep 5

# 执行数据库迁移
echo "执行数据库迁移..."
docker exec -i elderly-care-db psql -U postgres -d elderly_care < backend/migrations/000001_init_schema.up.sql || true

# 启动后端（后台）
echo "启动后端服务..."
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

# 启动前端（后台）
echo "启动前端服务..."
cd admin-frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "开发环境启动完成！"
echo "后端: http://localhost:8080"
echo "前端: http://localhost:3000"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 捕获退出信号
trap "echo '停止服务...'; kill $BACKEND_PID $FRONTEND_PID; docker-compose down; exit" INT TERM

# 等待
wait
