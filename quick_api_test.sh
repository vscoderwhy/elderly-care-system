#!/bin/bash

# 养老院管理系统 - 快速API测试
BASE_URL="http://1.12.223.138/api"

echo "========================================="
echo "养老院管理系统 - API测试"
echo "========================================="

# 1. 登录
echo -e "\n1. 测试登录..."
LOGIN_RESP=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}')

echo "登录响应: $LOGIN_RESP"

# 提取token
TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败：无法获取token"
    exit 1
fi

echo "✅ 登录成功"
echo "Token: ${TOKEN:0:30}..."

# 2. 测试Dashboard API
echo -e "\n2. 测试Dashboard API..."
DASHBOARD_RESP=$(curl -s "${BASE_URL}/dashboard/overview" -H "Authorization: Bearer $TOKEN")
echo "Dashboard响应: ${DASHBOARD_RESP:0:200}..."

# 3. 测试统计数据API（重点）
echo -e "\n3. 测试统计数据API..."

# 入住率趋势
echo -e "\n  3.1 入住率趋势(7天)..."
OCCUPANCY_RESP=$(curl -s "${BASE_URL}/statistics/occupancy?days=7" -H "Authorization: Bearer $TOKEN")
echo "响应: ${OCCUPANCY_RESP:0:300}..."

# 护理等级分布
echo -e "\n  3.2 护理等级分布..."
CARE_LEVEL_RESP=$(curl -s "${BASE_URL}/statistics/care-levels" -H "Authorization: Bearer $TOKEN")
echo "响应: ${CARE_LEVEL_RESP:0:300}..."

# 性别分布
echo -e "\n  3.3 性别分布..."
GENDER_RESP=$(curl -s "${BASE_URL}/statistics/gender" -H "Authorization: Bearer $TOKEN")
echo "响应: ${GENDER_RESP:0:300}..."

# 年龄分布
echo -e "\n  3.4 年龄分布..."
AGE_RESP=$(curl -s "${BASE_URL}/statistics/age" -H "Authorization: Bearer $TOKEN")
echo "响应: ${AGE_RESP:0:300}..."

# 护理统计
echo -e "\n  3.5 护理统计..."
CARE_STATS_RESP=$(curl -s "${BASE_URL}/statistics/care-stats" -H "Authorization: Bearer $TOKEN")
echo "响应: ${CARE_STATS_RESP:0:300}..."

# 4. 测试老人管理API
echo -e "\n4. 测试老人管理API..."
ELDERLY_RESP=$(curl -s "${BASE_URL}/elderly?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "老人列表响应: ${ELDERLY_RESP:0:300}..."

# 5. 测试护理记录API
echo -e "\n5. 测试护理记录API..."
CARE_RESP=$(curl -s "${BASE_URL}/care?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "护理记录响应: ${CARE_RESP:0:300}..."

# 6. 测试用药管理API
echo -e "\n6. 测试用药管理API..."
MED_RESP=$(curl -s "${BASE_URL}/medications?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "用药记录响应: ${MED_RESP:0:300}..."

# 7. 测试财务管理API
echo -e "\n7. 测试财务管理API..."
BILLS_RESP=$(curl -s "${BASE_URL}/bills?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "账单列表响应: ${BILLS_RESP:0:300}..."

# 8. 测试房间管理API
echo -e "\n8. 测试房间管理API..."
ROOMS_RESP=$(curl -s "${BASE_URL}/rooms" -H "Authorization: Bearer $TOKEN")
echo "房间列表响应: ${ROOMS_RESP:0:300}..."

# 9. 测试员工管理API
echo -e "\n9. 测试员工管理API..."
USERS_RESP=$(curl -s "${BASE_URL}/users?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "员工列表响应: ${USERS_RESP:0:300}..."

# 10. 测试库存管理API
echo -e "\n10. 测试库存管理API..."
INV_RESP=$(curl -s "${BASE_URL}/inventory?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "库存列表响应: ${INV_RESP:0:300}..."

# 11. 测试探视记录API
echo -e "\n11. 测试探视记录API..."
VISITS_RESP=$(curl -s "${BASE_URL}/visits?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "探视记录响应: ${VISITS_RESP:0:300}..."

# 12. 测试告警管理API
echo -e "\n12. 测试告警管理API..."
ALERTS_RESP=$(curl -s "${BASE_URL}/alerts?page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "告警列表响应: ${ALERTS_RESP:0:300}..."

# 13. 测试排班管理API
echo -e "\n13. 测试排班管理API..."
SCHEDULE_RESP=$(curl -s "${BASE_URL}/schedule" -H "Authorization: Bearer $TOKEN")
echo "排班列表响应: ${SCHEDULE_RESP:0:300}..."

echo -e "\n========================================="
echo "API测试完成"
echo "========================================="
