#!/bin/bash
# 批量修复前端API路径 - 添加 /api 前缀

echo "=== 批量修复前端API路径 ==="
echo ""

FRONTEND_DIR="/root/.openclaw/workspace-feishu-elderly/elderly-care-system/admin-frontend/src"

cd "$FRONTEND_DIR"

echo "修复 Alerts 页面..."
sed -i "s|axios.get('/alerts'|axios.get('/api/alerts'|g" views/Alerts/Index.vue
sed -i "s|axios.post('/alerts|axios.post('/api/alerts|g" views/Alerts/Index.vue
sed -i "s|axios.put(\`/alerts/|axios.put(\`/api/alerts/|g" views/Alerts/Index.vue

echo "修复 Medications 页面..."
sed -i "s|axios.get('/medications'|axios.get('/api/medications'|g" views/Medications/Index.vue
sed -i "s|axios.post('/medications|axios.post('/api/medications|g" views/Medications/Index.vue
sed -i "s|axios.put(\`/medications/|axios.put(\`/api/medications/|g" views/Medications/Index.vue
sed -i "s|axios.delete(\`/medications/|axios.delete(\`/api/medications/|g" views/Medications/Index.vue

echo "修复 Statistics/Dashboard 页面..."
sed -i "s|axios.get('/statistics/|axios.get('/api/statistics/|g" views/Statistics/Dashboard.vue

echo "修复 Statistics/Report 页面..."
sed -i "s|axios.get('/statistics/|axios.get('/api/statistics/|g" views/Statistics/Report.vue

echo "修复 System/Menus 页面..."
sed -i "s|axios.get('/system/menus')|axios.get('/api/system/menus')|g" views/System/Menus.vue
sed -i "s|axios.post('/system/menus'|axios.post('/api/system/menus'|g" views/System/Menus.vue
sed -i "s|axios.put(\`/system/menus/|axios.put(\`/api/system/menus/|g" views/System/Menus.vue
sed -i "s|axios.delete(\`/system/menus/|axios.delete(\`/api/system/menus/|g" views/System/Menus.vue

echo "修复 System/Roles 页面..."
sed -i "s|axios.get('/system/roles')|axios.get('/api/system/roles')|g" views/System/Roles.vue
sed -i "s|axios.post('/system/roles'|axios.post('/api/system/roles'|g" views/System/Roles.vue
sed -i "s|axios.put(\`/system/roles/|axios.put(\`/api/system/roles/|g" views/System/Roles.vue
sed -i "s|axios.delete(\`/system/roles/|axios.delete(\`/api/system/roles/|g" views/System/Roles.vue

echo ""
echo "=== 修复完成 ==="
echo ""
echo "修改的文件:"
git diff --name-only 2>/dev/null || echo "已修复所有文件"
echo ""
echo "请重启前端服务以应用更改:"
echo "  cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system"
echo "  bash restart-frontend.sh"
