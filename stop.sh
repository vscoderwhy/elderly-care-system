#!/bin/bash

# 养老院管理系统停止脚本

echo "停止养老院管理系统..."

# 停止后端
if [ -f "/tmp/elderly-backend.pid" ]; then
    PID=$(cat /tmp/elderly-backend.pid)
    if ps -p $PID > /dev/null 2>&1; then
        kill $PID
        echo "✓ 后端已停止"
    fi
    rm /tmp/elderly-backend.pid
fi

# 停止移动端
if [ -f "/tmp/elderly-mobile.pid" ]; then
    PID=$(cat /tmp/elderly-mobile.pid)
    if ps -p $PID > /dev/null 2>&1; then
        kill $PID
        echo "✓ 移动端已停止"
    fi
    rm /tmp/elderly-mobile.pid
fi

# 杀死所有相关进程
pkill -f "go run cmd/server/main.go" 2>/dev/null && echo "✓ Go 后端进程已清理" || true
pkill -f "vite.*3001" 2>/dev/null && echo "✓ Vite 移动端进程已清理" || true

echo "系统已停止"
