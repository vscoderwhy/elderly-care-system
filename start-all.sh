#!/bin/bash
# 养老院管理系统 - 生产环境启动脚本
# 使用 nohup 在后台运行，即使 SSH 断开也继续运行

set -e

WORKDIR="/root/.openclaw/workspace-feishu-elderly/elderly-care-system"
cd "$WORKDIR"

echo "=== 启动养老院管理系统 ==="
echo "$(date '+%Y-%m-%d %H:%M:%S') - 启动服务..."

# 停止旧进程
echo "1. 清理旧进程..."
pkill -9 -f "go run cmd/server/main.go" || true
pkill -9 -f "vite" || true
sleep 2

# 启动后端
echo "2. 启动后端服务..."
cd backend
nohup go run cmd/server/main.go > /tmp/elderly-backend.log 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > /tmp/elderly-backend.pid
echo "   后端 PID: $BACKEND_PID"
cd ..

# 等待后端启动
echo "   等待后端就绪..."
sleep 3

# 启动前端
echo "3. 启动前端服务..."
cd admin-frontend
nohup npm run dev > /tmp/elderly-frontend.log 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > /tmp/elderly-frontend.pid
echo "   前端 PID: $FRONTEND_PID"
cd ..

# 重启 Nginx
echo "4. 重启 Nginx..."
systemctl restart nginx
sleep 2

# 验证服务
echo ""
echo "=== 服务验证 ==="
echo ""

# 检查端口
echo "端口监听:"
netstat -tlnp 2>/dev/null | grep -E "80|3001|8080" | awk '{printf "   %-6s %s\n", $4, $7}'

echo ""
echo "服务测试:"
echo -n "   前端 (80): "
if curl -s -o /dev/null -w "%{http_code}" http://localhost/ | grep -q "200"; then
    echo "✅ OK"
else
    echo "❌ FAIL"
fi

echo -n "   后端 (8080): "
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/stats/dashboard | grep -q "401\|200"; then
    echo "✅ OK"
else
    echo "❌ FAIL"
fi

echo ""
echo "=== 启动完成 ==="
echo ""
echo "访问地址: http://1.12.223.138"
echo "登录账号: 13800138000"
echo "登录密码: 123456"
echo ""
echo "进程 ID:"
echo "  后端: $(cat /tmp/elderly-backend.pid)"
echo "  前端: $(cat /tmp/elderly-frontend.pid)"
echo ""
echo "日志文件:"
echo "  后端: tail -f /tmp/elderly-backend.log"
echo "  前端: tail -f /tmp/elderly-frontend.log"
echo ""
