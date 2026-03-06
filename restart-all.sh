#!/bin/bash
# 养老院管理系统 - 服务重启脚本

echo "=== 养老院管理系统 - 服务重启 ==="
echo ""

# 停止旧进程
echo "1. 停止旧进程..."
pkill -f "go run cmd/server/main.go"
pkill -f "vite"
sleep 2

# 启动后端
echo "2. 启动后端服务..."
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend
nohup go run cmd/server/main.go > /tmp/elderly-backend.log 2>&1 &
sleep 3

# 启动前端
echo "3. 启动前端服务..."
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/admin-frontend
nohup npm run dev > /tmp/elderly-frontend.log 2>&1 &
sleep 3

# 重启 Nginx
echo "4. 重启 Nginx..."
systemctl restart nginx
sleep 2

# 检查服务状态
echo ""
echo "=== 服务状态检查 ==="
echo ""

# 检查进程
echo "进程状态："
ps aux | grep -E "vite|go run" | grep -v grep | awk '{print "  ✓ " $11 " (PID: " $2 ")"}'

echo ""
echo "端口监听："
netstat -tlnp 2>/dev/null | grep -E "80|3001|8080" | awk '{print "  ✓ " $4 " -> " $7}'

echo ""
echo "服务测试："
echo -n "  前端 (80): "
curl -s -o /dev/null -w "%{http_code}" http://localhost/
echo ""
echo -n "  后端 (8080): "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/stats/dashboard
echo ""

echo ""
echo "=== 完成！==="
echo ""
echo "访问地址: http://1.12.223.138"
echo "登录账号: 13800138000"
echo "登录密码: 123456"
echo ""
echo "日志查看:"
echo "  后端: tail -f /tmp/elderly-backend.log"
echo "  前端: tail -f /tmp/elderly-frontend.log"
echo "  Nginx: tail -f /var/log/nginx/error.log"
