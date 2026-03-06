#!/bin/bash
# 测试导出功能

echo "=== 测试导出功能 ==="
echo ""

# 1. 登录获取token
echo "1. 登录..."
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}')

TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ 登录失败"
  echo "响应: $RESPONSE"
  exit 1
fi

echo "✅ 登录成功，token: ${TOKEN:0:20}..."
echo ""

# 2. 测试导出接口的响应头
echo "2. 检查导出接口响应头..."
curl -I -X GET "http://localhost:8080/api/export/elderly" \
  -H "Authorization: Bearer $TOKEN" \
  2>&1 | grep -E "HTTP|Content-Type|Content-Disposition"

echo ""

# 3. 尝试下载文件
echo "3. 下载导出文件..."
curl -X GET "http://localhost:8080/api/export/elderly" \
  -H "Authorization: Bearer $TOKEN" \
  -o /tmp/elderly_export_test.csv \
  -w "HTTP状态码: %{http_code}\n内容类型: %{content_type}\n下载大小: %{size_download} bytes\n"

echo ""

# 4. 检查文件内容
if [ -f /tmp/elderly_export_test.csv ]; then
  echo "4. 文件内容预览:"
  head -5 /tmp/elderly_export_test.csv
  echo ""
  echo "文件大小: $(wc -l < /tmp/elderly_export_test.csv) 行"
  echo "文件类型: $(file /tmp/elderly_export_test.csv)"
else
  echo "❌ 文件下载失败"
fi

echo ""
echo "=== 测试完成 ==="
