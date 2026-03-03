#!/bin/bash

# 快速种子数据配置脚本
# 用于将演示数据快速集成到前端项目

set -e

echo "=========================================="
echo "  养老院管理系统 - 种子数据配置"
echo "=========================================="
echo ""

# 检查项目目录
if [ ! -d "admin-frontend" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

echo "📁 正在配置种子数据..."
echo ""

# 创建 seed_data 目录
mkdir -p admin-frontend/src/seed_data

# 复制种子数据文件
echo "📋 复制种子数据文件..."
cp seed_data/elderly_full.json admin-frontend/src/seed_data/
cp seed_data/care_records_full.json admin-frontend/src/seed_data/
cp seed_data/bills_full.json admin-frontend/src/seed_data/

echo "   ✅ 老人数据: elderly_full.json"
echo "   ✅ 护理记录: care_records_full.json"
echo "   ✅ 账单数据: bills_full.json"
echo ""

# 创建种子数据加载器
echo "📝 创建种子数据加载器..."

cat > admin-frontend/src/utils/seedData.ts <<'EOF'
// 种子数据加载器 - 用于客户演示
import elderlyData from '@/seed_data/elderly_full.json'
import careRecordsData from '@/seed_data/care_records_full.json'
import billsData from '@/seed_data/bills_full.json'

export const getElderlyList = () => elderlyData
export const getElderlyById = (id: number) => elderlyData.find((e: any) => e.id === id)
export const getCareRecords = () => careRecordsData
export const getBills = () => billsData

export const getStatistics = () => {
  const totalElderly = elderlyData.length
  const totalBills = billsData.length
  const totalCareRecords = careRecordsData.length

  const paidBills = billsData.filter((b: any) => b.status === 'paid')
  const unpaidBills = billsData.filter((b: any) => b.status === 'unpaid')
  const overdueBills = billsData.filter((b: any) => b.status === 'overdue')

  const totalAmount = billsData.reduce((sum: number, b: any) => sum + b.amount, 0)
  const paidAmount = paidBills.reduce((sum: number, b: any) => sum + b.amount, 0)

  return {
    elderly: {
      total: totalElderly,
      occupancyRate: ((totalElderly / 200) * 100).toFixed(1),
      byCareLevel: {
        '特级': elderlyData.filter((e: any) => e.careLevel === '特级').length,
        '一级': elderlyData.filter((e: any) => e.careLevel === '一级').length,
        '二级': elderlyData.filter((e: any) => e.careLevel === '二级').length,
        '三级': elderlyData.filter((e: any) => e.careLevel === '三级').length
      },
      avgHealthScore: (elderlyData.reduce((sum: number, e: any) => sum + e.healthScore, 0) / totalElderly).toFixed(0)
    },
    care: {
      totalRecords: totalCareRecords,
      todayRecords: careRecordsData.filter((r: any) => r.careTime.includes('2026-03-04')).length,
      avgEvaluation: (careRecordsData.reduce((sum: number, r: any) => sum + r.evaluation, 0) / totalCareRecords).toFixed(1)
    },
    finance: {
      totalBills: totalBills,
      totalAmount: totalAmount.toFixed(2),
      paidAmount: paidAmount.toFixed(2),
      unpaidAmount: (totalAmount - paidAmount).toFixed(2),
      paidCount: paidBills.length,
      unpaidCount: unpaidBills.length,
      overdueCount: overdueBills.length,
      collectionRate: ((paidAmount / totalAmount) * 100).toFixed(1)
    }
  }
}

export const getTodayTasks = () => {
  const taskTypes = ['日常护理', '健康监测', '康复训练', '医疗护理']
  const statuses = ['pending', 'in_progress', 'completed']
  const priorities = ['normal', 'important', 'urgent']

  return Array.from({ length: 45 }, (_, i) => {
    const elderly = elderlyData[i % elderlyData.length]
    return {
      id: i + 1,
      elderlyId: elderly.id,
      elderlyName: elderly.name,
      bedNumber: elderly.bedNumber,
      taskType: taskTypes[i % taskTypes.length],
      description: '完成相关护理服务',
      scheduledDate: '2026-03-04',
      scheduledTime: `${8 + (i % 10)}:00`,
      nurseId: (i % 20) + 1,
      nurseName: ['赵护士', '李护士', '周护士', '吴护士', '郑护士'][i % 5],
      status: statuses[i % statuses.length],
      priority: priorities[i % priorities.length]
    }
  })
}

export const getStaffList = () => {
  return Array.from({ length: 60 }, (_, i) => {
    const departments = ['护理部', '医务室', '康复科', '膳食部', '行政部']
    const dept = departments[i % departments.length]

    const positions: Record<string, string[]> = {
      '护理部': ['护士', '护理员', '护工', '护士长'],
      '医务室': ['医师', '护士', '主任医师', '主治医师'],
      '康复科': ['康复师', '物理治疗师', '作业治疗师'],
      '膳食部': ['营养师', '厨师', '配餐员'],
      '行政部': ['行政专员', '主任', '院长', '副院长']
    }

    const surnames = ['王', '李', '张', '刘', '陈', '杨', '黄', '赵', '周', '吴']
    const femaleNames = ['秀英', '桂英', '秀珍', '淑兰', '玉兰', '淑华', '翠英', '文英', '秀兰', '玉英']
    const maleNames = ['伟', '强', '磊', '洋', '勇', '军', '杰', '涛', '超', '明']

    const name = surnames[i % surnames.length] + (i % 2 === 0 ? femaleNames[i % femaleNames.length] : maleNames[i % maleNames.length])

    return {
      id: i + 1,
      employeeNo: `EMP${String(i + 1).padStart(3, '0')}`,
      name: name,
      gender: i % 2 === 0 ? '女' : '男',
      department: dept,
      position: positions[dept][i % positions[dept].length],
      phone: `138****${String(Math.floor(1000 + Math.random() * 9000))}`,
      hireDate: `201${Math.floor(18 + Math.random() * 8)}-${String(Math.floor(1 + Math.random() * 12)).padStart(2, '0')}`,
      status: '在职'
    }
  })
}

export const getVisitAppointments = () => {
  return Array.from({ length: 80 }, (_, i) => {
    const elderly = elderlyData[i % elderlyData.length]
    const relationships = ['子女', '配偶', '孙辈', '其他亲属', '朋友']
    const statuses = ['pending', 'approved', 'completed', 'cancelled']

    return {
      id: i + 1,
      appointmentNo: `VA202603${String(i + 1).padStart(3, '0')}`,
      elderlyId: elderly.id,
      elderlyName: elderly.name,
      bedNumber: elderly.bedNumber,
      visitorName: surnames[i % surnames.length] + (i % 2 === 0 ? '先生' : '女士'),
      visitorPhone: `139****${String(Math.floor(1000 + Math.random() * 9000))}`,
      relationship: relationships[i % relationships.length],
      visitType: i % 3 === 0 ? '视频探访' : '现场探访',
      visitDate: `2026-03-${String(Math.floor(4 + Math.random() * 25)).padStart(2, '0')}`,
      visitTime: `${9 + (i % 8) * 2}:00`,
      visitorCount: 1 + Math.floor(Math.random() * 3),
      status: statuses[i % statuses.length],
      remarks: ''
    }
  })
}
EOF

echo "   ✅ seedData.ts 创建完成"
echo ""

echo "📄 数据统计："
echo "   - 老人数据: 50条"
echo "   - 护理记录: 30条"
echo "   - 账单数据: 50条"
echo "   - 生成数据: 员工60条、任务45条、预约80条"
echo ""

echo "✅ 配置完成！"
echo ""
echo "📖 下一步操作："
echo "   1. 启动前端开发服务器："
echo "      cd admin-frontend && npm run dev"
echo ""
echo "   2. 或使用后端API："
echo "      cd backend && go run cmd/main.go"
echo ""
echo "   3. 访问管理后台："
echo "      http://localhost:5173"
echo ""
