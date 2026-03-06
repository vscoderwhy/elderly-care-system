#!/bin/bash

# 调试老人详情API的问题

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJwaG9uZSI6IjEzODAwMTM4MDAwIiwiZXhwIjoxNzczMjI1MTc3LCJuYmYiOjE3NzI2MjAzNzcsImlhdCI6MTc3MjYyMDM3N30.XAubynU4M6VbBrWMIqSTvHBxFPRHclA35mQ026nBsUg"

echo "========================================="
echo "调试老人详情API问题"
echo "========================================="

echo -e "\n1. 测试老人列表（获取ID）："
curl -s "http://1.12.223.138/api/elderly?page=1&page_size=2" -H "Authorization: Bearer $TOKEN" | jq '.data.list[] | {id, name, bed_id}'

echo -e "\n2. 测试老人详情 ID=1："
curl -s "http://1.12.223.138/api/elderly/1" -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n3. 测试老人详情 ID=2："
curl -s "http://1.12.223.138/api/elderly/2" -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n4. 测试老人详情 ID=999（不存在）："
curl -s "http://1.12.223.138/api/elderly/999" -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n========================================="
echo "分析："
echo "如果ID=1和ID=2都返回404，但列表中有数据，"
echo "可能是Preload关联查询导致的错误。"
echo "需要检查："
echo "  1. Bed表是否存在"
echo "  2. Bed.Room外键是否正确"
echo "  3. Families表是否存在"
echo "========================================="
