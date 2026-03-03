# 种子数据使用指南

## 数据概览

本项目提供了完整的种子数据，可用于系统演示和客户展示：

### 数据统计

| 数据类型 | 数量 | 说明 |
|---------|------|------|
| 老人数据 | 150条 | 包含姓名、性别、年龄、护理等级、床位等信息 |
| 护理记录 | 300条 | 包含护理类型、内容、时间、评价等 |
| 健康数据 | 500条 | 包含血压、心率、血糖、体温等 |
| 账单数据 | 200条 | 包含账单类型、金额、支付状态等 |
| 护理任务 | 150条 | 自动生成的今日任务 |
| 员工数据 | 60条 | 包含各部门员工信息 |
| 探视预约 | 100条 | 包含预约状态、访客信息等 |

## 使用方法

### 方法一：直接导入 JSON（推荐）

#### 1. 将种子数据移动到正确位置

```bash
# 确保种子数据在正确位置
cp seed_data/*.json admin-frontend/src/seed_data/
```

#### 2. 更新前端页面使用种子数据

在 `admin-frontend/src/views/Elderly/List.vue` 等页面中，将数据获取部分替换为：

```typescript
import { getElderlyList } from '@/utils/seedData'

// 在组件中使用
const loadElderly = () => {
  loading.value = true
  setTimeout(() => {
    elderlyList.value = getElderlyList()
    total.value = elderlyList.value.length
    loading.value = false
  }, 300)
}
```

#### 3. 支持的关键页面

以下页面可以直接使用种子数据：

- **老人管理**: `views/Elderly/List.vue`
- **护理记录**: `views/Care/Records.vue`
- **财务管理**: `views/Finance/Bills.vue`
- **护理任务**: `views/Care/Tasks.vue`
- **探视管理**: `views/Visits/Appointments.vue`
- **数据看板**: `views/Dashboard/Overview.vue`

### 方法二：更新 API Handler 临时返回种子数据

修改 `backend/internal/api/handler/` 中的 handler 文件，直接返回种子数据：

```go
// 示例：更新 elderly.go
func (h *ElderlyHandler) List(c *gin.Context) {
    // 读取种子数据
    data, err := os.ReadFile("seed_data/elderly_full.json")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    var elderlyList []map[string]interface{}
    json.Unmarshal(data, &elderlyList)

    c.JSON(200, gin.H{
        "code": 0,
        "message": "获取成功",
        "data": gin.H{
            "list": elderlyList,
            "total": len(elderlyList),
        },
    })
}
```

### 方法三：批量导入到数据库

#### 1. 启动后端服务

```bash
cd backend
go run cmd/main.go
```

#### 2. 使用脚本导入数据

```bash
# 运行种子数据生成器
cd backend
go run cmd/seed/main.go
```

这会将数据保存到 `seed_data/` 目录。

#### 3. 使用数据库导入工具

```bash
# 连接到 PostgreSQL
psql -U elderly_care -d elderly_care_db

# 在 psql 中
\i seed_data/init.sql
```

## 快速集成到现有页面

### 示例 1：更新老人列表页面

编辑 `admin-frontend/src/views/Elderly/List.vue`：

```typescript
// 在 import 部分添加
import { getElderlyList, getStatistics } from '@/utils/seedData'

// 替换 loadBills 函数
const loadBills = async () => {
  loading.value = true
  try {
    // 模拟API延迟
    await new Promise(resolve => setTimeout(resolve, 300))

    // 使用种子数据
    elderlyList.value = getElderlyList()
    total.value = elderlyList.value.length

    // 应用筛选
    if (currentFilter.value && currentFilter.value !== 'all') {
      elderlyList.value = elderlyList.value.filter(e =>
        e.status === currentFilter.value ||
        e.careLevel === currentFilter.value
      )
    }
  } finally {
    loading.value = false
  }
}
```

### 示例 2：更新数据看板

编辑 `admin-frontend/src/views/Dashboard/Overview.vue`：

```typescript
import { getStatistics } from '@/utils/seedData'

// 在 onMounted 中使用
onMounted(() => {
  const statistics = getStatistics()

  stats.value = [
    { key: 'total', label: '在院老人', value: statistics.elderly.total, unit: '人', trend: 8.5 },
    { key: 'occupancy', label: '入住率', value: statistics.elderly.occupancyRate, unit: '%', trend: 5.2 },
    { key: 'health', label: '健康评分', value: statistics.elderly.avgHealthScore, unit: '分', trend: 3.1 },
    { key: 'care', label: '护理记录', value: statistics.care.totalRecords, unit: '条', trend: 12.3 }
  ]

  // 财务统计
  financeStats.value = [
    { key: 'total', label: '总账单', value: statistics.finance.totalBills, unit: '笔' },
    { key: 'amount', label: '总金额', value: statistics.finance.totalAmount, unit: '元' },
    { key: 'paid', label: '已收', value: statistics.finance.paidAmount, unit: '元' },
    { key: 'unpaid', label: '待收', value: statistics.finance.unpaidAmount, unit: '元' }
  ]
})
```

## 数据文件说明

### elderly_full.json
- **路径**: `seed_data/elderly_full.json`
- **数量**: 50条（可扩展）
- **字段**: id, name, gender, age, careLevel, bedNumber, admitDate, status, healthScore, phone, emergencyContact, emergencyPhone

### care_records_full.json
- **路径**: `seed_data/care_records_full.json`
- **数量**: 30条（可扩展）
- **字段**: id, elderlyId, elderlyName, bedNumber, careType, content, careTime, nurseName, images, evaluation, tags, remarks

### bills_full.json
- **路径**: `seed_data/bills_full.json`
- **数量**: 50条（可扩展）
- **字段**: id, billNo, elderlyId, elderlyName, bedNumber, billType, amount, period, billDate, dueDate, status, paymentDate, paymentMethod

## 扩展数据

如需生成更多数据，可以：

### 1. 使用种子生成器

```bash
cd backend
go run cmd/seed/main.go
```

### 2. 手动扩展 JSON 文件

复制现有的 JSON 条目并修改相应字段值。

### 3. 使用在线工具

- JSON Generator: https://www.json-generator.com/
- Mock Data: https://www.mockaroo.com/

## 注意事项

1. **ID唯一性**: 确保所有数据的ID字段唯一
2. **关联性**: elderlyId 等外键应与实际数据匹配
3. **日期格式**: 统一使用 `YYYY-MM-DD` 格式
4. **状态枚举**: 使用正确的状态值（paid/unpaid/overdue等）
5. **金额精度**: 金额保留两位小数

## 常见问题

### Q: 前端页面显示空白？
A: 检查浏览器控制台是否有错误，确认种子数据文件路径正确。

### Q: 数据重复或ID冲突？
A: 运行数据生成器前确保使用随机种子，或手动修改ID。

### Q: 如何快速生成更多数据？
A: 修改 `backend/cmd/seed/main.go` 中的生成数量参数，重新运行生成器。

## 客户演示建议

1. **准备环境**: 提前一天准备好所有数据和页面
2. **场景演练**: 按照典型使用流程演练一遍
3. **数据更新**: 演示当天更新一些时间相关数据（如今日日期）
4. **备选方案**: 准备多套数据场景，根据客户关注点切换
