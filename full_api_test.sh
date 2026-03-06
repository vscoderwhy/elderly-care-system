#!/bin/bash

# 养老院管理系统 - 全面功能测试（更新版）
BASE_URL="http://1.12.223.138/api"

echo "========================================="
echo "养老院管理系统 - 全面API功能测试"
echo "测试时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果记录
declare -a FAILED_ENDPOINTS
declare -a PASSED_ENDPOINTS
declare -a ERROR_RESPONSES

# 测试函数
test_api() {
    local name="$1"
    local endpoint="$2"
    local method="${3:-GET}"
    local is_critical="${4:-false}"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -e "\n${BLUE}[${TOTAL_TESTS}]${NC} 测试: ${name}"
    echo "端点: ${method} ${endpoint}"

    response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}" \
        -H "Authorization: Bearer ${TOKEN}")

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "200" ]; then
        # 检查响应是否包含 "code":0 (成功)
        if echo "$body" | grep -q '"code":0'; then
            echo -e "${GREEN}✓ 通过${NC} (HTTP 200)"
            PASSED_TESTS=$((PASSED_TESTS + 1))
            PASSED_ENDPOINTS+=("${name}")
            # 显示响应数据的前150字符
            echo "响应: $(echo "$body" | head -c 150)..."
        else
            echo -e "${YELLOW}⚠ 警告${NC} (HTTP 200 但业务失败)"
            echo "响应: $body"
            FAILED_TESTS=$((FAILED_TESTS + 1))
            FAILED_ENDPOINTS+=("${name} - 业务错误")
            ERROR_RESPONSES+=("${name}: $body")
        fi
    else
        echo -e "${RED}✗ 失败${NC} (HTTP $http_code)"
        echo "响应: $body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        FAILED_ENDPOINTS+=("${name} - HTTP $http_code")
        ERROR_RESPONSES+=("${name}: $body")

        if [ "$is_critical" = "true" ]; then
            echo -e "${RED}关键测试失败，停止测试${NC}"
            break
        fi
    fi
}

# ============ 开始测试 ============

# 1. 登录
echo -e "\n${YELLOW}========== 1. 认证测试 ==========${NC}"
LOGIN_RESP=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}')

TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 登录失败：无法获取token${NC}"
    echo "响应: $LOGIN_RESP"
    exit 1
fi

echo -e "${GREEN}✓ 登录成功${NC}"
echo "Token: ${TOKEN:0:30}..."
TOTAL_TESTS=$((TOTAL_TESTS + 1))
PASSED_TESTS=$((PASSED_TESTS + 1))
PASSED_ENDPOINTS+=("用户登录")

# 2. 统计数据（重点测试 - 之前有图表渲染问题）
echo -e "\n${YELLOW}========== 2. 统计数据测试（重点） ==========${NC}"

echo -e "\n${BLUE}[重要]${NC} 这些API对应之前修复的Dashboard图表"

test_api "统计Dashboard" "/statistics/dashboard"
test_api "入住率趋势" "/statistics/occupancy-trend?days=7"
test_api "健康趋势" "/statistics/health-trend?days=7"
test_api "财务统计" "/statistics/finance"
test_api "护理统计" "/statistics/care"
test_api "年龄分布" "/statistics/age-distribution"
test_api "月度报表" "/statistics/monthly-report?year=2026&month=3"

# 3. Dashboard概览
echo -e "\n${YELLOW}========== 3. Dashboard概览 ==========${NC}"
test_api "Dashboard统计" "/stats/dashboard"
test_api "床位入住率" "/stats/occupancy"
test_api "护理统计" "/stats/care"
test_api "财务统计" "/stats/finance"

# 4. 老人管理
echo -e "\n${YELLOW}========== 4. 老人管理 ==========${NC}"
test_api "老人列表" "/elderly?page=1&page_size=10"

# 获取第一个老人的ID
ELDERLY_RESP=$(curl -s "${BASE_URL}/elderly?page=1&page_size=1" -H "Authorization: Bearer ${TOKEN}")
ELDERLY_ID=$(echo "$ELDERLY_RESP" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')

if [ ! -z "$ELDERLY_ID" ]; then
    echo "找到老人ID: $ELDERLY_ID"
    test_api "老人详情" "/elderly/$ELDERLY_ID"
fi

# 5. 护理记录
echo -e "\n${YELLOW}========== 5. 护理记录 ==========${NC}"
test_api "护理记录列表" "/care/records?page=1&page_size=10"
test_api "我的护理任务" "/care/my-tasks"
test_api "护理项目列表" "/care/items"

# 6. 健康记录
echo -e "\n${YELLOW}========== 6. 健康记录 ==========${NC}"
test_api "健康记录列表" "/health/records?page=1&page_size=10"
if [ ! -z "$ELDERLY_ID" ]; then
    test_api "老人最新健康记录" "/health/records/latest/$ELDERLY_ID"
fi

# 7. 服务请求
echo -e "\n${YELLOW}========== 7. 服务请求 ==========${NC}"
test_api "服务请求列表" "/service/requests?page=1&page_size=10"

# 8. 用药管理
echo -e "\n${YELLOW}========== 8. 用药管理 ==========${NC}"
test_api "用药记录列表" "/medications?page=1&page_size=10"
test_api "用药预警" "/medications/alerts"
if [ ! -z "$ELDERLY_ID" ]; then
    test_api "老人用药记录" "/elderly/$ELDERLY_ID/medications"
    test_api "老人今日用药" "/elderly/$ELDERLY_ID/medications/today"
fi

# 9. 财务管理
echo -e "\n${YELLOW}========== 9. 财务管理 ==========${NC}"
test_api "账单列表" "/bills?page=1&page_size=10"

# 10. 房间管理
echo -e "\n${YELLOW}========== 10. 房间管理 ==========${NC}"
test_api "楼栋列表" "/rooms/buildings"
test_api "房间列表" "/rooms"
test_api "床位统计" "/rooms/stats"

# 11. 排班管理
echo -e "\n${YELLOW}========== 11. 排班管理 ==========${NC}"
test_api "排班列表" "/schedules?page=1&page_size=10"
test_api "我的排班" "/schedules/my"
test_api "排班月度统计" "/schedules/stats/monthly?year=2026&month=3"

# 12. 探视预约
echo -e "\n${YELLOW}========== 12. 探视预约 ==========${NC}"
test_api "探视记录列表" "/visits?page=1&page_size=10"
test_api "今日探视" "/visits/today"
test_api "即将到来的探视" "/visits/upcoming"

# 13. 告警管理
echo -e "\n${YELLOW}========== 13. 告警管理 ==========${NC}"
test_api "告警列表" "/alerts?page=1&page_size=10"
test_api "告警汇总" "/alerts/summary"
test_api "活跃告警" "/alerts/active"

# 14. 考勤管理
echo -e "\n${YELLOW}========== 14. 考勤管理 ==========${NC}"
test_api "今日考勤" "/attendance/today"
test_api "考勤统计" "/attendance/stats"
test_api "我的排班" "/attendance/schedule/my"

# 15. 库存管理
echo -e "\n${YELLOW}========== 15. 库存管理 ==========${NC}"
test_api "库存分类" "/inventory/categories"
test_api "库存物品列表" "/inventory/items?page=1&page_size=10"
test_api "低库存预警" "/inventory/low-stock"
test_api "库存统计" "/inventory/stats"

# 16. 数据导出
echo -e "\n${YELLOW}========== 16. 数据导出 ==========${NC}"
test_api "导出老人列表" "/export/elderly"
test_api "导出护理记录" "/export/care-records"
test_api "导出健康数据" "/export/health-data"
test_api "导出财务数据" "/export/finance"

# 17. 系统管理
echo -e "\n${YELLOW}========== 17. 系统管理 ==========${NC}"
test_api "用户列表" "/system/users"
test_api "角色列表" "/system/roles"
test_api "菜单列表" "/system/menus"
test_api "权限列表" "/system/permissions"

# 18. 员工管理
echo -e "\n${YELLOW}========== 18. 员工管理 ==========${NC}"
test_api "员工列表" "/staff?page=1&page_size=10"

# 19. 用户个人信息
echo -e "\n${YELLOW}========== 19. 用户信息 ==========${NC}"
test_api "用户个人信息" "/user/profile"
test_api "用户菜单" "/user/menus"
test_api "用户权限" "/user/permissions"

# ============ 测试结果汇总 ============
echo -e "\n========================================="
echo -e "测试结果汇总"
echo "========================================="
echo -e "总测试数: ${YELLOW}$TOTAL_TESTS${NC}"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"
echo -e "通过率: $(awk "BEGIN {printf \"%.2f%%\", ($PASSED_TESTS/$TOTAL_TESTS)*100}")"

if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "\n${RED}失败的端点 ($FAILED_TESTS):${NC}"
    for endpoint in "${FAILED_ENDPOINTS[@]}"; do
        echo "  ✗ $endpoint"
    done
fi

echo -e "\n${GREEN}通过的端点 ($PASSED_TESTS):${NC}"
for endpoint in "${PASSED_ENDPOINTS[@]}"; do
    echo "  ✓ $endpoint"
done

# 保存结果
REPORT_FILE="/root/.openclaw/workspace-feishu-elderly/elderly-care-system/FULL_TEST_REPORT_$(date +%Y%m%d_%H%M%S).md"

{
    echo "# 养老院管理系统 - 全面API测试报告"
    echo ""
    echo "**测试时间**: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "**测试环境**: http://1.12.223.138"
    echo ""
    echo "## 测试统计"
    echo ""
    echo "| 指标 | 数值 |"
    echo "|------|------|"
    echo "| 总测试数 | $TOTAL_TESTS |"
    echo "| 通过 | $PASSED_TESTS |"
    echo "| 失败 | $FAILED_TESTS |"
    echo "| 通过率 | $(awk "BEGIN {printf \"%.2f%%\", ($PASSED_TESTS/$TOTAL_TESTS)*100}") |"
    echo ""

    if [ $FAILED_TESTS -gt 0 ]; then
        echo "## ❌ 失败的API端点"
        echo ""
        for endpoint in "${FAILED_ENDPOINTS[@]}"; do
            echo "- $endpoint"
        done
        echo ""
    fi

    echo "## ✅ 通过的API端点"
    echo ""
    for endpoint in "${PASSED_ENDPOINTS[@]}"; do
        echo "- $endpoint"
    done
    echo ""

    if [ $FAILED_TESTS -gt 0 ]; then
        echo "## 错误响应详情"
        echo ""
        for err in "${ERROR_RESPONSES[@]}"; do
            echo "### ${err%%:*}"
            echo ""
            echo '```json'
            echo "${err#*: }"
            echo '```'
            echo ""
        done
    fi

    echo "## 测试覆盖范围"
    echo ""
    echo "- ✅ 认证模块"
    echo "- ✅ 统计数据（Dashboard图表）"
    echo "- ✅ 老人管理"
    echo "- ✅ 护理记录"
    echo "- ✅ 健康记录"
    echo "- ✅ 用药管理"
    echo "- ✅ 财务管理"
    echo "- ✅ 房间管理"
    echo "- ✅ 排班管理"
    echo "- ✅ 探视预约"
    echo "- ✅ 告警管理"
    echo "- ✅ 考勤管理"
    echo "- ✅ 库存管理"
    echo "- ✅ 数据导出"
    echo "- ✅ 系统管理"
    echo "- ✅ 员工管理"
    echo ""

} > "$REPORT_FILE"

echo -e "\n详细报告已保存到: $REPORT_FILE"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️  有 $FAILED_TESTS 个测试失败${NC}"
    exit 1
fi
