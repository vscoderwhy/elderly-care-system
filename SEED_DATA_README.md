# 种子数据完整说明

## 📊 数据文件清单

| 文件名 | 路径 | 数据量 | 说明 |
|--------|------|--------|------|
| 老人数据 | `seed_data/elderly_full.json` | 50条 | 包含姓名、护理等级、床位等完整信息 |
| 护理记录 | `seed_data/care_records_full.json` | 30条 | 包含护理类型、内容、评价等 |
| 账单数据 | `seed_data/bills_full.json` | 50条 | 包含账单类型、金额、状态等 |
| SQL脚本 | `scripts/seed_data.sql` | - | 可直接导入数据库的SQL脚本 |

## 🚀 快速使用

### 方法一：自动化脚本（推荐）

```bash
# 在项目根目录执行
./scripts/seed_data.sh
```

这将自动完成以下操作：
1. 复制数据文件到前端项目
2. 创建种子数据加载器
3. 配置好所有导入路径

### 方法二：手动配置

#### 步骤 1：复制数据文件

```bash
mkdir -p admin-frontend/src/seed_data
cp seed_data/*.json admin-frontend/src/seed_data/
```

#### 步骤 2：创建数据加载器

将 `utils/seedData.ts` 复制到项目中。

#### 步骤 3：在页面中使用

```typescript
import { getElderlyList } from '@/utils/seedData'

// 在组件中
const loadData = () => {
  elderlyList.value = getElderlyList()
}
```

## 📋 数据详情

### 1. 老人数据 (elderly_full.json)

```json
{
  "id": 1,
  "name": "王秀英",
  "gender": "女",
  "age": 78,
  "careLevel": "二级",
  "bedNumber": "1号楼101",
  "admitDate": "2023-03-15",
  "status": "在院",
  "healthScore": 85,
  "phone": "138****1234",
  "emergencyContact": "王强",
  "emergencyPhone": "139****5678"
}
```

**字段说明**：
- `id`: 唯一标识
- `name`: 姓名
- `gender`: 性别（男/女）
- `age`: 年龄
- `careLevel`: 护理等级（特级/一级/二级/三级）
- `bedNumber`: 床位号（楼栋+楼层+房间号）
- `admitDate`: 入住日期
- `status`: 状态（在院/外出/退住）
- `healthScore`: 健康评分（0-100）
- `phone`: 联系电话
- `emergencyContact`: 紧急联系人
- `emergencyPhone`: 紧急联系电话

**数据分布**：
- 性别：女26人，男24人
- 年龄：60-90岁，平均78岁
- 护理等级：特级8人、一级18人、二级16人、三级8人
- 楼栋：1-4号楼，均匀分布

### 2. 护理记录 (care_records_full.json)

```json
{
  "id": 1,
  "elderlyId": 1,
  "elderlyName": "王秀英",
  "bedNumber": "1号楼101",
  "careType": "日常护理",
  "content": "协助老人起床洗漱，测量血压128/82mmHg...",
  "careTime": "2026-03-04 08:30",
  "nurseName": "赵护士",
  "images": ["/uploads/care1.jpg"],
  "evaluation": 5,
  "tags": ["服务热情", "专业细致"],
  "remarks": "老人今日情绪稳定"
}
```

**护理类型分布**：
- 日常护理：30%
- 健康监测：25%
- 康复训练：20%
- 医疗护理：15%
- 心理疏导：10%

**评价分布**：
- 5分：40%
- 4分：50%
- 3分：10%

### 3. 账单数据 (bills_full.json)

```json
{
  "id": 1,
  "billNo": "B202603001",
  "elderlyId": 1,
  "elderlyName": "王秀英",
  "bedNumber": "1号楼101",
  "billType": "床位费",
  "amount": 2800,
  "period": "2026年3月",
  "billDate": "2026-03-01",
  "dueDate": "2026-03-10",
  "status": "paid",
  "paymentDate": "2026-03-05",
  "paymentMethod": "微信支付"
}
```

**账单类型分布**：
- 床位费：40%
- 护理费：30%
- 伙食费：15%
- 医疗费：10%
- 其他：5%

**状态分布**：
- 已支付：50%
- 未支付：30%
- 已逾期：20%

**金额范围**：
- 床位费：2500-3500元
- 护理费：500-2000元
- 伙食费：600元/月
- 医疗费：100-5000元

## 🔧 后端使用

### 导入到 PostgreSQL

```bash
# 连接到数据库
psql -U elderly_care -d elderly_care_db

# 执行SQL脚本
\i scripts/seed_data.sql
```

### 或使用Go程序

```bash
cd backend
go run cmd/seed/main.go
```

这将生成：
- 150条老人数据
- 300条护理记录
- 500条健康数据
- 200条账单数据
- 150条护理任务
- 100条探视预约

## 📱 前端使用

### 老人列表页面

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getElderlyList } from '@/utils/seedData'

const elderlyList = ref([])

onMounted(() => {
  elderlyList.value = getElderlyList()
})
</script>
```

### 数据看板

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getStatistics } from '@/utils/seedData'

const stats = ref([])

onMounted(() => {
  const statistics = getStatistics()

  stats.value = [
    {
      key: 'total',
      label: '在院老人',
      value: statistics.elderly.total,
      unit: '人',
      trend: 8.5
    },
    {
      key: 'occupancy',
      label: '入住率',
      value: statistics.elderly.occupancyRate,
      unit: '%',
      trend: 5.2
    }
  ]
})
</script>
```

## 🎯 演示场景建议

### 场景一：全面展示

展示所有模块的完整数据：
- 老人管理：50条记录，支持搜索、筛选
- 护理记录：30条记录，展示护理过程
- 健康数据：每条老人对应10条健康记录
- 账单管理：50条记录，展示收费情况

### 场景二：数据看板演示

重点展示统计图表：
- 入住率：75%
- 护理完成率：91%
- 费用收回率：85%
- 健康评分分布

### 场景三：功能演示

演示特定功能：
- 支付流程：选中未支付账单，展示支付弹窗
- 任务分配：查看今日任务，分配给员工
- 数据导出：导出老人名单、护理记录等

## ⚙️ 自定义数据

### 修改现有数据

直接编辑JSON文件，注意：
1. 保持ID唯一性
2. 维护外键关联
3. 使用正确的枚举值

### 添加更多数据

复制现有条目，修改：
1. ID（确保唯一）
2. 姓名、床位等标识信息
3. 日期、时间等时间字段

### 生成随机数据

使用在线工具：
- Mockaroo: https://www.mockaroo.com/
- JSON Generator: https://www.json-generator.com/

## 📊 数据统计

### 老人统计
- 总数：50人
- 平均年龄：78岁
- 性别比例：女52% : 男48%
- 护理等级：二级占比最高（32%）

### 护理记录统计
- 总记录：30条
- 平均评价：4.7分
- 完成率：90%
- 今日记录：10条

### 账单统计
- 总账单：50笔
- 总金额：¥124,500
- 已支付：¥62,200
- 未支付：¥62,300
- 支付率：50%

## 🔍 数据验证

### 检查数据完整性

```bash
# 检查JSON格式
cat seed_data/elderly_full.json | jq '. | length'
cat seed_data/care_records_full.json | jq '. | length'
cat seed_data/bills_full.json | jq '. | length'
```

### 检查ID唯一性

```bash
# 检查老人ID
cat seed_data/elderly_full.json | jq '[.[].id' | sort -u
```

### 检查外键关联

```bash
# 检查护理记录的老人ID是否存在
cat seed_data/care_records_full.json | jq '[].[].elderlyId' | sort -u
```

## 📝 常见问题

### Q1: 前端显示 "Cannot find module '@/seedData'"
**A**: 检查文件是否在正确的位置，应该是 `admin-frontend/src/seed_data/`

### Q2: 数据加载后页面空白
**A**: 检查JSON格式是否正确，可以使用在线JSON验证工具

### Q3: 如何快速生成更多数据？
**A**:
1. 运行 `go run cmd/seed/main.go` 生成新数据
2. 或复制现有JSON条目并修改ID

### Q4: 数据重复怎么处理？
**A**: 使用Set或filter去重，检查是否有相同ID

### Q5: 如何更新日期到今天？
**A**:
```bash
# 使用sed批量替换日期
sed -i 's/2026-03-04/$(date +%Y-%m-%d)/g' seed_data/*.json
```

## 📞 技术支持

如需帮助，请查看：
- 项目README.md
- DEPLOYMENT.md
- SEED_DATA_GUIDE.md

或联系技术支持：support@elderly-care.com
