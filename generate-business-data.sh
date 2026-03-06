#!/bin/bash
# 一键生成所有业务数据

echo "=== 生成所有业务数据 ==="
echo ""

cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend

# 后台运行，避免超时
nohup go run cmd/generator/business.go > /tmp/business-data-generator.log 2>&1 &
PID=$!

echo "数据生成器已在后台启动，PID: $PID"
echo "日志文件: /tmp/business-data-generator.log"
echo ""
echo "查看进度:"
echo "  tail -f /tmp/business-data-generator.log"
echo ""
echo "查看数据统计:"
echo "  docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM medications;'"
echo ""
echo "等待30秒后自动检查进度..."
sleep 30

echo ""
echo "=== 当前进度 ==="
docker exec elderly-care-db psql -U postgres -d elderly_care -c "
SELECT '员工' as 类型, COUNT(*) as 数量 FROM users WHERE role = 'caregiver'
UNION ALL
SELECT '药品', COUNT(*) FROM medications
UNION ALL  
SELECT '用药记录', COUNT(*) FROM medication_records
UNION ALL
SELECT '账单', COUNT(*) FROM bills
UNION ALL
SELECT '账单明细', COUNT(*) FROM bill_items
UNION ALL
SELECT '支付记录', COUNT(*) FROM payments
UNION ALL
SELECT '健康记录', COUNT(*) FROM health_records
UNION ALL
SELECT '考勤记录', COUNT(*) FROM attendances;
"

echo ""
echo "生成器仍在后台运行..."
echo "完成后会在日志中显示最终统计"
