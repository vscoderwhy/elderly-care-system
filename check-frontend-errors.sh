#!/bin/bash
# 前端API路径修复脚本

echo "=== 检查并修复前端API路径问题 ==="
echo ""

FRONTEND_DIR="/root/.openclaw/workspace-feishu-elderly/elderly-care-system/admin-frontend/src"

echo "1. 检查API调用是否缺少 /api 前缀..."
echo ""

# 查找所有可能缺少 /api 前缀的API调用
echo "查找可疑的API调用..."
grep -r "axios.get\|axios.post\|axios.put\|axios.delete" "$FRONTEND_DIR" --include="*.vue" --include="*.ts" --include="*.tsx" | grep -v "/api/" | head -20

echo ""
echo "2. 检查常见的前端错误模式..."
echo ""

# 检查是否有未处理的Promise
echo "检查未处理的异步错误..."
grep -r "await.*\.get\|await.*\.post" "$FRONTEND_DIR" --include="*.vue" --include="*.ts" -A 2 | grep -v "try\|catch" | head -10

echo ""
echo "3. 检查是否有未定义的变量或函数..."
echo ""

# 检查ref/reactive使用
echo "检查Vue响应式变量..."
grep -r "ref\|reactive" "$FRONTEND_DIR" --include="*.vue" | wc -l

echo ""
echo "4. 检查是否有TypeScript类型错误..."
echo ""

# 检查any类型使用
echo "检查 TypeScript any 类型使用..."
grep -r ": any" "$FRONTEND_DIR" --include="*.vue" --include="*.ts" | wc -l

echo ""
echo "=== 检查完成 ==="
echo ""
echo "建议："
echo "1. 查看浏览器控制台(F12)的具体错误信息"
echo "2. 检查Network标签，看哪些API请求失败"
echo "3. 将具体的错误信息发给我，我可以精准修复"
