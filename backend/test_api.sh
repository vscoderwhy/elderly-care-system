#!/bin/bash

# API测试脚本
BASE_URL="http://localhost:8080"
TOKEN=""

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "后端 API 测试报告"
echo "测试时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================="
echo ""

# 首先尝试登录获取token
echo -e "${YELLOW}1. 尝试登录获取认证令牌...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token（如果登录成功）
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo -e "${RED}登录失败或token为空，将尝试无认证访问${NC}"
  TOKEN=""
else
  echo -e "${GREEN}登录成功，获取到token${NC}"
fi

echo ""
echo "========================================="
echo "开始测试统计API端点"
echo "========================================="
echo ""

# 测试函数
test_api() {
  local endpoint=$1
  local description=$2
  
  echo -e "${YELLOW}测试: $description${NC}"
  echo "端点: GET $endpoint"
  
  if [ -n "$TOKEN" ]; then
    RESPONSE=$(curl -s -X GET "${BASE_URL}${endpoint}" \
      -H "Authorization: Bearer ${TOKEN}" \
      -H "Content-Type: application/json")
  else
    RESPONSE=$(curl -s -X GET "${BASE_URL}${endpoint}" \
      -H "Content-Type: application/json")
  fi
  
  echo "响应状态: $(echo $RESPONSE | grep -o '"status":"[^"]*' | cut -d'"' -f4 || echo '无status字段')"
  echo "响应内容:"
  echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
  echo ""
  
  # 保存响应到文件
  echo "$RESPONSE" > "/tmp/api_test_$(basename $endpoint | tr '/' '_').json"
}

# 测试各个端点
test_api "/api/statistics/dashboard" "Dashboard统计数据"

test_api "/api/statistics/occupancy-trend?days=30" "入住率趋势（30天）"

test_api "/api/statistics/age-distribution" "老人年龄分布"

test_api "/api/statistics/care" "护理统计数据"

echo "========================================="
echo "测试完成！"
echo "========================================="
echo ""
echo "响应文件已保存到 /tmp/ 目录"
