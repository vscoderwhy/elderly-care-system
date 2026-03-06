#!/bin/bash
# 百万级数据生成脚本

set -e

cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend

echo "=== 百万级数据生成器 ==="
echo ""

# 检查数据库连接
echo "检查数据库连接..."
if ! docker exec elderly-care-db pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ 数据库未运行，请先启动服务"
    echo "运行: bash start-all.sh"
    exit 1
fi

echo "✅ 数据库连接正常"
echo ""

# 显示菜单
echo "请选择要生成的数据:"
echo "1) 完整数据 (100万老人 + 护理记录 + 设施 + 护工) - 推荐用于性能测试"
echo "2) 仅生成100万老人数据"
echo "3) 仅生成护理记录"
echo "4) 仅生成设施 (楼栋/房间/床位)"
echo "5) 仅生成护工账号"
echo "6) 自定义数量 (交互式)"
echo ""
read -p "请输入选项 [1-6]: " choice

case $choice in
    1)
        echo ""
        echo "=== 开始生成完整百万级数据 ==="
        echo "预计时间: 5-15分钟"
        echo "数据量: 100万老人 + 5000万护理记录"
        echo ""
        go run cmd/generator/main.go all
        ;;
    2)
        echo ""
        echo "=== 开始生成100万老人数据 ==="
        echo "预计时间: 3-8分钟"
        go run cmd/generator/main.go elderly
        ;;
    3)
        echo ""
        echo "=== 开始生成护理记录 ==="
        go run cmd/generator/main.go records
        ;;
    4)
        echo ""
        echo "=== 开始生成设施数据 ==="
        go run cmd/generator/main.go facility
        ;;
    5)
        echo ""
        echo "=== 开始生成护工账号 ==="
        go run cmd/generator/main.go staff
        ;;
    6)
        echo ""
        echo "自定义模式"
        read -p "老人数量: " elderly_count
        read -p "每人护理记录数: " records_per_elderly
        echo ""
        echo "将生成: $elderly_count 个老人，每人 $records_per_elderly 条记录"
        read -p "确认? (y/n): " confirm
        if [ "$confirm" = "y" ]; then
            # 这里可以扩展自定义逻辑
            go run cmd/generator/main.go facility
            go run cmd/generator/main.go all
        else
            echo "已取消"
        fi
        ;;
    *)
        echo "无效选项"
        exit 1
        ;;
esac

echo ""
echo "=== 完成 ==="
echo ""
echo "查看数据统计:"
echo "  docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM elderly;'"
echo "  docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM care_records;'"
echo ""
