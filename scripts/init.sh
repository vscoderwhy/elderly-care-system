#!/bin/bash

# 养老院管理系统 - 初始化脚本

set -e

echo "=========================================="
echo "  养老院管理系统 - 初始化安装"
echo "=========================================="
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

echo "✅ Docker 环境检查通过"
echo ""

# 创建环境变量文件
if [ ! -f .env ]; then
    echo "📝 创建环境变量文件..."
    cp .env.example .env
    echo "⚠️  请编辑 .env 文件，修改必要的配置（如数据库密码）"
    echo ""
    read -p "是否现在编辑 .env 文件？(y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ${EDITOR:-vi} .env
    fi
else
    echo "✅ 环境变量文件已存在"
fi

# 创建必要的目录
echo ""
echo "📁 创建必要的目录..."
mkdir -p uploads
mkdir -p scripts
mkdir -p nginx/ssl
mkdir -p monitoring/prometheus
mkdir -p monitoring/grafana/dashboards
mkdir -p monitoring/grafana/datasources

# 生成自签名 SSL 证书（开发用）
if [ ! -f nginx/ssl/cert.pem ]; then
    echo ""
    echo "🔐 生成自签名 SSL 证书（仅用于开发）..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=ElderlyCare/CN=localhost"
    echo "✅ SSL 证书已生成"
fi

# 创建监控配置
echo ""
echo "📊 创建监控配置..."
cat > monitoring/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'backend'
    static_configs:
      - targets: ['backend:8080']
EOF

cat > monitoring/grafana/datasources/prometheus.yml <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
EOF

echo "✅ 监控配置已创建"

# 初始化数据库
echo ""
echo "🗄️  初始化数据库..."
cat > scripts/init.sql <<'EOF'
-- 创建数据库扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建触发器函数
CREATE OR REPLACE FUNCTION updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';
EOF

echo "✅ 数据库初始化脚本已创建"

# 构建并启动服务
echo ""
echo "🚀 构建并启动服务..."
if docker compose version &> /dev/null; then
    docker compose up -d --build
else
    docker-compose up -d --build
fi

# 等待服务启动
echo ""
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "📊 服务状态："
if docker compose version &> /dev/null; then
    docker compose ps
else
    docker-compose ps
fi

echo ""
echo "=========================================="
echo "  安装完成！"
echo "=========================================="
echo ""
echo "🌐 访问地址："
echo "  - 管理后台: http://localhost"
echo "  - 移动端 H5: http://mobile.localhost"
echo "  - API 文档: http://localhost:8080/swagger"
echo "  - Prometheus: http://localhost:9090"
echo "  - Grafana: http://localhost:3000 (admin/admin)"
echo ""
echo "📝 下一步："
echo "  1. 访问管理后台并登录"
echo "  2. 完成系统初始化配置"
echo "  3. 添加老人、员工等基础数据"
echo ""
echo "📖 更多信息请查看 README.md 和 DEPLOYMENT.md"
echo ""
