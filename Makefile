.PHONY: help build run test clean docker-up docker-down docker-logs migrate

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make build         - 编译后端"
	@echo "  make run           - 运行后端服务"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理编译文件"
	@echo "  make migrate       - 执行数据库迁移"
	@echo "  make migrate-rollback - 回滚数据库迁移"
	@echo "  make docker-up     - 启动 Docker 服务"
	@echo "  make docker-down   - 停止 Docker 服务"
	@echo "  make docker-logs   - 查看 Docker 日志"
	@echo "  make install-deps  - 安装依赖"

# 编译后端
build:
	@echo "编译后端..."
	cd backend && go build -o bin/server ./cmd/server

# 运行后端
run:
	@echo "启动后端服务..."
	cd backend && go run cmd/server/main.go

# 运行测试
test:
	@echo "运行测试..."
	cd backend && go test -v ./...

# 清理
clean:
	@echo "清理编译文件..."
	cd backend && rm -rf bin/
	cd admin-frontend && rm -rf dist/

# 安装依赖
install-deps:
	@echo "安装后端依赖..."
	cd backend && go mod download
	@echo "安装前端依赖..."
	cd admin-frontend && npm install

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	@docker exec -i elderly-care-db psql -U postgres -d elderly_care < backend/migrations/000001_init_schema.up.sql

# 启动 Docker 服务
docker-up:
	@echo "启动 Docker 服务..."
	docker-compose up -d

# 停止 Docker 服务
docker-down:
	@echo "停止 Docker 服务..."
	docker-compose down

# 查看 Docker 日志
docker-logs:
	docker-compose logs -f

# 前端开发服务器
dev-frontend:
	@echo "启动前端开发服务器..."
	cd admin-frontend && npm run dev

# 构建前端
build-frontend:
	@echo "构建前端..."
	cd admin-frontend && npm run build

# 初始化项目
init: install-deps docker-up
	@echo "等待数据库启动..."
	@sleep 5
	@echo "执行数据库迁移..."
	$(MAKE) migrate
	@echo "初始化完成！"
	@echo "后端: http://localhost:8080"
	@echo "前端: http://localhost:3000"
