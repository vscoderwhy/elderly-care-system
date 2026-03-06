#!/bin/bash
# 回滚错误的修复 - 移除多余的 /api 前缀

echo "=== 回滚错误的API路径修复 ==="
echo ""

FRONTEND_DIR="/root/.openclaw/workspace-feishu-elderly/elderly-care-system/admin-frontend/src"

cd "$FRONTEND_DIR"

echo "回滚 Alerts 页面..."
sed -i "s|axios.get('/api/alerts'|axios.get('/alerts'|g" views/Alerts/Index.vue
sed -i "s|axios.post('/api/alerts|axios.post('/alerts|g" views/Alerts/Index.vue
sed -i "s|axios.put(\`/api/alerts/|axios.put(\`/alerts/|g" views/Alerts/Index.vue

echo "回滚 Medications 页面..."
sed -i "s|axios.get('/api/medications'|axios.get('/medications'|g" views/Medications/Index.vue
sed -i "s|axios.post('/api/medications|axios.post('/medications|g" views/Medications/Index.vue
sed -i "s|axios.put(\`/api/medications/|axios.put(\`/medications/|g" views/Medications/Index.vue
sed -i "s|axios.delete(\`/api/medications/|axios.delete(\`/medications/|g" views/Medications/Index.vue

echo "回滚 Statistics/Dashboard 页面..."
sed -i "s|axios.get('/api/statistics/|axios.get('/statistics/|g" views/Statistics/Dashboard.vue

echo "回滚 Statistics/Report 页面..."
sed -i "s|axios.get('/api/statistics/|axios.get('/statistics/|g" views/Statistics/Report.vue

echo "回滚 System/Menus 页面..."
sed -i "s|axios.get('/api/system/menus')|axios.get('/system/menus')|g" views/System/Menus.vue
sed -i "s|axios.post('/api/system/menus'|axios.post('/system/menus'|g" views/System/Menus.vue
sed -i "s|axios.put(\`/api/system/menus/|axios.put(\`/system/menus/|g" views/System/Menus.vue
sed -i "s|axios.delete(\`/api/system/menus/|axios.delete(\`/system/menus/|g" views/System/Menus.vue

echo "回滚 System/Roles 页面..."
sed -i "s|axios.get('/api/system/roles')|axios.get('/system/roles')|g" views/System/Roles.vue
sed -i "s|axios.post('/api/system/roles'|axios.post('/system/roles'|g" views/System/Roles.vue
sed -i "s|axios.put(\`/api/system/roles/|axios.put(\`/system/roles/|g" views/System/Roles.vue
sed -i "s|axios.delete(\`/api/system/roles/|axios.delete(\`/system/roles/|g" views/System/Roles.vue

echo ""
echo "=== 回滚完成 ==="
echo ""
echo "API baseURL已经配置为'/api'，所以代码中应该使用相对路径"
echo "例如: axios.get('/alerts') 会自动变成 '/api/alerts'"
