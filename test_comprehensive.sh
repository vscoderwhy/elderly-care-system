#!/bin/bash

# 养老院管理系统 - 全面API测试脚本
# 测试时间: 2026-03-04

BASE_URL="http://1.12.223.138/api"
USERNAME="13800138000"
PASSWORD="123456"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果记录
declare -a FAILED_ENDPOINTS
declare -a PASSED_ENDPOINTS

# 测试函数
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local token="$4"
    local data="$5"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -e "\n${YELLOW}测试 ${TOTAL_TESTS}: ${name}${NC}"
    echo "端点: ${method} ${BASE_URL}${endpoint}"

    if [ -z "$token" ]; then
        if [ -z "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}")
        else
            response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d "${data}")
        fi
    else
        if [ -z "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}" \
                -H "Authorization: Bearer ${token}")
        else
            response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}" \
                -H "Authorization: Bearer ${token}" \
                -H "Content-Type: application/json" \
                -d "${data}")
        fi
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ 通过${NC} (HTTP $http_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        PASSED_ENDPOINTS+=("${name} (${method} ${endpoint})")
        # 显示响应数据的前200字符
        echo "响应: $(echo "$body" | head -c 200)..."
    else
        echo -e "${RED}✗ 失败${NC} (HTTP $http_code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        FAILED_ENDPOINTS+=("${name} (${method} ${endpoint}) - HTTP $http_code")
        echo "响应: $body"
    fi
}

echo "========================================="
echo "养老院管理系统 - 全面API测试"
echo "========================================="
echo "基础URL: $BASE_URL"
echo "测试账号: $USERNAME"
echo "========================================="

# 1. 认证测试
echo -e "\n${YELLOW}========== 1. 认证模块测试 ==========${NC}"

test_endpoint "用户登录" "POST" "/auth/login" "" "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\""

# 从登录响应中提取token
login_response=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}致命错误: 无法获取Token，停止测试${NC}"
    exit 1
fi

echo -e "${GREEN}Token已获取: ${TOKEN:0:20}...${NC}"

# 2. Dashboard和工作台
echo -e "\n${YELLOW}========== 2. Dashboard测试 ==========${NC}"
test_endpoint "工作台概览" "GET" "/dashboard/overview" "$TOKEN"
test_endpoint "待办事项" "GET" "/dashboard/todos" "$TOKEN"

# 3. 统计数据（重点测试）
echo -e "\n${YELLOW}========== 3. 统计数据测试（重点） ==========${NC}"
test_endpoint "入住率趋势(7天)" "GET" "/statistics/occupancy?days=7" "$TOKEN"
test_endpoint "入住率趋势(30天)" "GET" "/statistics/occupancy?days=30" "$TOKEN"
test_endpoint "护理等级分布" "GET" "/statistics/care-levels" "$TOKEN"
test_endpoint "性别分布" "GET" "/statistics/gender" "$TOKEN"
test_endpoint "年龄分布" "GET" "/statistics/age" "$TOKEN"
test_endpoint "护理统计" "GET" "/statistics/care-stats" "$TOKEN"
test_endpoint "报表数据" "GET" "/stats/summary" "$TOKEN"

# 4. 老人管理
echo -e "\n${YELLOW}========== 4. 老人管理测试 ==========${NC}"
test_endpoint "老人列表" "GET" "/elderly?page=1&page_size=10" "$TOKEN"
test_endpoint "老人统计" "GET" "/elderly/stats" "$TOKEN"
# 尝试获取第一个老人的ID
elderly_list=$(curl -s "${BASE_URL}/elderly?page=1&page_size=1" -H "Authorization: Bearer ${TOKEN}")
ELDERLY_ID=$(echo "$elderly_list" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
if [ ! -z "$ELDERLY_ID" ]; then
    test_endpoint "老人详情" "GET" "/elderly/$ELDERLY_ID" "$TOKEN"
fi

# 5. 护理记录
echo -e "\n${YELLOW}========== 5. 护理记录测试 ==========${NC}"
test_endpoint "护理记录列表" "GET" "/care?page=1&page_size=10" "$TOKEN"
test_endpoint "护理统计" "GET" "/care/stats" "$TOKEN"
if [ ! -z "$ELDERLY_ID" ]; then
    test_endpoint "老人的护理记录" "GET" "/care/elderly/$ELDERLY_ID?page=1&page_size=10" "$TOKEN"
fi

# 6. 用药管理
echo -e "\n${YELLOW}========== 6. 用药管理测试 ==========${NC}"
test_endpoint "用药记录列表" "GET" "/medications?page=1&page_size=10" "$TOKEN"
if [ ! -z "$ELDERLY_ID" ]; then
    test_endpoint "老人的用药记录" "GET" "/medications/elderly/$ELDERLY_ID" "$TOKEN"
fi

# 7. 财务管理
echo -e "\n${YELLOW}========== 7. 财务管理测试 ==========${NC}"
test_endpoint "账单列表" "GET" "/bills?page=1&page_size=10" "$TOKEN"
test_endpoint "财务统计" "GET" "/bills/stats" "$TOKEN"
if [ ! -z "$ELDERLY_ID" ]; then
    test_endpoint "老人的账单" "GET" "/bills/elderly/$ELDERLY_ID" "$TOKEN"
fi

# 8. 房间管理
echo -e "\n${YELLOW}========== 8. 房间管理测试 ==========${NC}"
test_endpoint "房间列表" "GET" "/rooms" "$TOKEN"
test_endpoint "房间统计" "GET" "/rooms/stats" "$TOKEN"

# 9. 员工管理
echo -e "\n${YELLOW}========== 9. 员工管理测试 ==========${NC}"
test_endpoint "员工列表" "GET" "/users?page=1&page_size=10" "$TOKEN"

# 10. 库存管理
echo -e "\n${YELLOW}========== 10. 库存管理测试 ==========${NC}"
test_endpoint "库存列表" "GET" "/inventory?page=1&page_size=10" "$TOKEN"

# 11. 探视记录
echo -e "\n${YELLOW}========== 11. 探视记录测试 ==========${NC}"
test_endpoint "探视记录列表" "GET" "/visits?page=1&page_size=10" "$TOKEN"
if [ ! -z "$ELDERLY_ID" ]; then
    test_endpoint "老人的探视记录" "GET" "/visits/elderly/$ELDERLY_ID" "$TOKEN"
fi

# 12. 告警管理
echo -e "\n${YELLOW}========== 12. 告警管理测试 ==========${NC}"
test_endpoint "告警列表" "GET" "/alerts?page=1&page_size=10" "$TOKEN"

# 13. 排班管理
echo -e "\n${YELLOW}========== 13. 排班管理测试 ==========${NC}"
test_endpoint "排班列表" "GET" "/schedule" "$TOKEN"

# 14. 导出功能
echo -e "\n${YELLOW}========== 14. 导出功能测试 ==========${NC}"
test_endpoint "导出老人数据" "GET" "/export/elderly" "$TOKEN"

# 15. 系统管理
echo -e "\n${YELLOW}========== 15. 系统管理测试 ==========${NC}"
test_endpoint "用户列表" "GET" "/users" "$TOKEN"
test_endpoint "角色列表" "GET" "/roles" "$TOKEN"

# 测试结果汇总
echo -e "\n========================================="
echo -e "测试结果汇总"
echo "========================================="
echo -e "总测试数: ${YELLOW}$TOTAL_TESTS${NC}"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"
echo "通过率: $(awk "BEGIN {printf \"%.2f%%\", ($PASSED_TESTS/$TOTAL_TESTS)*100}")"

if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "\n${RED}失败的端点:${NC}"
    for endpoint in "${FAILED_ENDPOINTS[@]}"; do
        echo "  ✗ $endpoint"
    done
fi

echo -e "\n${GREEN}通过的端点:${NC}"
for endpoint in "${PASSED_ENDPOINTS[@]}"; do
    echo "  ✓ $endpoint"
done

# 保存结果到文件
{
    echo "# API测试结果 - $(date)"
    echo ""
    echo "## 测试统计"
    echo "- 总测试数: $TOTAL_TESTS"
    echo "- 通过: $PASSED_TESTS"
    echo "- 失败: $FAILED_TESTS"
    echo "- 通过率: $(awk "BEGIN {printf \"%.2f%%\", ($PASSED_TESTS/$TOTAL_TESTS)*100}")"
    echo ""
    if [ $FAILED_TESTS -gt 0 ]; then
        echo "## 失败的端点"
        for endpoint in "${FAILED_ENDPOINTS[@]}"; do
            echo "- $endpoint"
        done
        echo ""
    fi
    echo "## 通过的端点"
    for endpoint in "${PASSED_ENDPOINTS[@]}"; do
        echo "- $endpoint"
    done
} > /root/.openclaw/workspace-feishu-elderly/elderly-care-system/api_test_results_$(date +%Y%m%d_%H%M%S).txt

echo -e "\n详细结果已保存到: api_test_results_$(date +%Y%m%d_%H%M%S).txt"

exit $FAILED_TESTS
