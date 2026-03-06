#!/bin/bash
# 百万级数据生成器 - 一键启动脚本

set -e

echo "=== 智慧养老系统 - 百万级数据生成器 ==="
echo ""

cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend

# 检查数据库连接
echo "🔍 检查数据库连接..."
if ! docker exec elderly-care-db pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ 数据库未运行"
    echo "请先启动服务: bash start-all.sh"
    exit 1
fi
echo "✅ 数据库连接正常"
echo ""

# 显示菜单
echo "请选择要生成的数据量:"
echo ""
echo "1️⃣  测试模式 (100条老人数据) - 验证功能是否正常"
echo "2️⃣  小批量 (1万条老人数据) - 快速测试"
echo "3️⃣  中批量 (10万条老人数据) - 中等规模测试"
echo "4️⃣  大批量 (100万条老人数据) - 完整压力测试"
echo "5️⃣  仅生成护理记录 (需要先有老人数据)"
echo "6️⃣  仅创建设施数据 (楼栋/房间/床位)"
echo ""
read -p "请输入选项 [1-6]: " choice

case $choice in
    1)
        echo ""
        echo "🚀 开始生成测试数据 (100条老人)..."
        go run cmd/generator/main.go test
        ;;
    2)
        echo ""
        echo "🚀 开始生成1万条老人数据..."
        cat > main_small.go << 'EOF'
package main
import "fmt"
func main() {
    fmt.Println("生成1万条数据...")
    // 临时修改 elderlyCount 为 10000
    // 可以通过修改源代码中的参数实现
}
EOF
        echo "提示：请修改代码中的 elderlyCount 参数为 10000"
        echo "然后运行: go run cmd/generator/main.go elderly"
        ;;
    3)
        echo ""
        echo "🚀 开始生成10万条老人数据..."
        echo "提示：请修改代码中的 elderlyCount 参数为 100000"
        echo "然后运行: go run cmd/generator/main.go elderly"
        ;;
    4)
        echo ""
        echo "🚀 开始生成100万条老人数据 + 护理记录..."
        echo "⏱️  预计耗时: 15-30分钟"
        echo ""
        read -p "确认继续? (y/n): " confirm
        if [ "$confirm" = "y" ]; then
            go run cmd/generator/main.go all
        else
            echo "已取消"
        fi
        ;;
    5)
        echo ""
        echo "🚀 开始生成护理记录..."
        go run cmd/generator/main.go records
        ;;
    6)
        echo ""
        echo "🚀 开始创建设施数据..."
        go run cmd/generator/main.go facility
        ;;
    *)
        echo "❌ 无效选项"
        exit 1
        ;;
esac

echo ""
echo "=== 完成 ==="
echo ""
echo "📊 查看数据统计:"
echo "  老人总数: docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM elderly;'"
echo "  护理记录: docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM care_records;'"
echo ""
echo "🌐 访问系统: http://1.12.223.138"
echo ""
