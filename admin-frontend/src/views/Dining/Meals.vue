<template>
  <div class="dining-management">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Food /></el-icon>
          膳食管理
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreateMeal">
          <el-icon><Plus /></el-icon>
          新建膳食
        </el-button>
        <el-button @click="handleMealPlan">
          <el-icon><Calendar /></el-icon>
          膳食计划
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in mealStats" :key="stat.key">
        <div class="meal-stat" :class="`stat-${stat.type}`">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-icon">
            <component :is="stat.icon" />
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 今日膳食 -->
    <el-card shadow="never" class="today-card">
      <template #header>
        <div class="card-header">
          <span>今日膳食安排</span>
          <el-date-picker v-model="todayDate" type="date" placeholder="选择日期" @change="loadTodayMeals" />
        </div>
      </template>

      <el-table :data="todayMeals" stripe>
        <el-table-column prop="mealType" label="餐别" width="100">
          <template #default="{ row }">
            <el-tag :type="getMealTypeColor(row.mealType)" size="small">
              {{ row.mealType }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="mealTime" label="用餐时间" width="120" />
        <el-table-column prop="menu" label="菜单" min-width="300">
          <template #default="{ row }">
            <el-tag
              v-for="dish in row.dishes"
              :key="dish"
              size="small"
              class="dish-tag"
            >
              {{ dish }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="calories" label="热量(千卡)" width="120" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleViewDetail(row)">详情</el-button>
            <el-button size="small" text @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" text @click="handleCopy(row)">复制到</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 膳食记录列表 -->
    <el-card shadow="never" class="records-card">
      <template #header>
        <div class="card-header">
          <span>膳食记录</span>
          <div class="header-actions">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              @change="loadRecords"
            />
          </div>
        </div>
      </template>

      <el-table :data="mealRecords" stripe v-loading="loading">
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column prop="mealType" label="餐别" width="100" />
        <el-table-column prop="menu" label="菜单" min-width="300" show-overflow-tooltip />
        <el-table-column prop="calories" label="热量" width="100" />
        <el-table-column prop="protein" label="蛋白质" width="100" />
        <el-table-column prop="fat" label="脂肪" width="100" />
        <el-table-column prop="carbs" label="碳水" width="100" />
        <el-table-column prop="creator" label="创建人" width="100" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" text type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </el-card>

    <!-- 膳食对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px">
      <el-form :model="mealForm" :rules="mealRules" label-width="120px">
        <el-form-item label="日期" prop="date">
          <el-date-picker v-model="mealForm.date" type="date" placeholder="选择日期" />
        </el-form-item>

        <el-form-item label="餐别" prop="mealType">
          <el-radio-group v-model="mealForm.mealType">
            <el-radio label="早餐">早餐</el-radio>
            <el-radio label="午餐">午餐</el-radio>
            <el-radio label="晚餐">晚餐</el-radio>
            <el-radio label="加餐">加餐</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="用餐时间" prop="mealTime">
          <el-time-picker v-model="mealForm.mealTime" placeholder="选择时间" format="HH:mm" />
        </el-form-item>

        <el-form-item label="菜单" prop="dishes">
          <div class="dishes-input">
            <div v-for="(dish, index) in mealForm.dishes" :key="index" class="dish-item">
              <el-input v-model="mealForm.dishes[index]" placeholder="请输入菜品名称" />
              <el-button size="small" @click="removeDish(index)" :disabled="mealForm.dishes.length <= 1">
                删除
              </el-button>
            </div>
            <el-button size="small" @click="addDish">添加菜品</el-button>
          </div>
        </el-form-item>

        <el-form-item label="营养信息">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-input v-model="mealForm.calories" placeholder="热量" type="number">
                <template #prepend>千卡</template>
              </el-input>
            </el-col>
            <el-col :span="6">
              <el-input v-model="mealForm.protein" placeholder="蛋白质" type="number">
                <template #prepend>g</template>
              </el-input>
            </el-col>
            <el-col :span="6">
              <el-input v-model="mealForm.fat" placeholder="脂肪" type="number">
                <template #prepend>g</template>
              </el-input>
            </el-col>
            <el-col :span="6">
              <el-input v-model="mealForm.carbs" placeholder="碳水" type="number">
                <template #prepend>g</template>
              </el-input>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item label="适用人群" prop="targetGroup">
          <el-select v-model="mealForm.targetGroup" placeholder="请选择">
            <el-option label="全部老人" value="all" />
            <el-option label="糖尿病老人" value="diabetes" />
            <el-option label="高血压老人" value="hypertension" />
            <el-option label="素食老人" value="vegetarian" />
            <el-option label="软食老人" value="soft" />
          </el-select>
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="mealForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Food, Plus, Calendar } from '@element-plus/icons-vue'

// 统计数据
const mealStats = ref([
  { key: 'today', label: '今日膳食', value: '4餐', type: 'primary', icon: 'Calendar' },
  { key: 'average', label: '人均热量', value: '1800千卡', type: 'success', icon: 'TrendCharts' },
  { key: 'special', label: '特殊膳食', value: '15份', type: 'warning', icon: 'Warning' },
  { key: 'variety', label: '菜品种类', value: '48种', type: 'info', icon: 'Menu' }
])

const todayDate = ref(new Date())
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const saving = ref(false)

const todayMeals = ref([
  {
    id: 1,
    mealType: '早餐',
    mealTime: '07:00',
    dishes: ['小米粥', '鸡蛋', '馒头', '咸菜', '牛奶'],
    calories: 450,
    protein: 18,
    fat: 12,
    carbs: 65
  },
  {
    id: 2,
    mealType: '午餐',
    mealTime: '11:30',
    dishes: ['米饭', '红烧肉', '清炒时蔬', '紫菜蛋花汤'],
    calories: 750,
    protein: 35,
    fat: 22,
    carbs: 95
  },
  {
    id: 3,
    mealType: '晚餐',
    mealTime: '17:30',
    dishes: ['杂粮饭', '清蒸鱼', '炒青菜', '豆腐汤'],
    calories: 620,
    protein: 28,
    fat: 15,
    carbs: 78
  },
  {
    id: 4,
    mealType: '加餐',
    mealTime: '15:00',
    dishes: ['水果', '酸奶', '小点心'],
    calories: 180,
    protein: 5,
    fat: 3,
    carbs: 32
  }
])

const dateRange = ref([])
const mealRecords = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(0)

const mealForm = reactive({
  id: '',
  date: '',
  mealType: '早餐',
  mealTime: '',
  dishes: [''],
  calories: '',
  protein: '',
  fat: '',
  carbs: '',
  targetGroup: 'all',
  remark: ''
})

const mealRules: FormRules = {
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  mealType: [{ required: true, message: '请选择餐别', trigger: 'change' }],
  mealTime: [{ required: true, message: '请选择时间', trigger: 'change' }],
  dishes: [{ required: true, message: '请至少添加一道菜品', trigger: 'change' }]
}

const getMealTypeColor = (type: string) => {
  const map: Record<string, any> = {
    '早餐': 'success',
    '午餐': 'primary',
    '晚餐': 'warning',
    '加餐': 'info'
  }
  return map[type] || ''
}

const loadTodayMeals = () => {
  // TODO: 加载今日膳食
  console.log('加载今日膳食', todayDate.value)
}

const loadRecords = () => {
  // TODO: 加载记录列表
  console.log('加载记录', dateRange.value)
}

const handleCreateMeal = () => {
  dialogTitle.value = '新建膳食'
  resetForm()
  dialogVisible.value = true
}

const handleMealPlan = () => {
  console.log('膳食计划')
}

const handleViewDetail = (row: any) => {
  console.log('查看详情', row)
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑膳食'
  Object.assign(mealForm, row)
  dialogVisible.value = true
}

const handleCopy = async (row: any) => {
  try {
    const { value } = await ElMessageBox.prompt(
      '选择要复制到的日期',
      '复制膳食',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /.+/,
        inputErrorMessage: '请输入日期'
      }
    )

    ElMessage.success('复制成功')
  } catch {
    // 用户取消
  }
}

const handleView = (row: any) => {
  console.log('查看', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这条膳食记录吗？', '提示', {
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadRecords()
  } catch {
    // 取消
  }
}

const addDish = () => {
  mealForm.dishes.push('')
}

const removeDish = (index: number) => {
  mealForm.dishes.splice(index, 1)
}

const resetForm = () => {
  Object.assign(mealForm, {
    id: '',
    date: '',
    mealType: '早餐',
    mealTime: '',
    dishes: [''],
    calories: '',
    protein: '',
    fat: '',
    carbs: '',
    targetGroup: 'all',
    remark: ''
  })
}

const handleSave = async () => {
  // TODO: 保存逻辑
  saving.value = true
  setTimeout(() => {
    saving.value = false
    dialogVisible.value = false
    ElMessage.success('保存成功')
  }, 1000)
}

onMounted(() => {
  loadTodayMeals()
  loadRecords()
})
</script>

<style scoped lang="scss">
.dining-management {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .page-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 24px;
    font-weight: 600;
    margin: 0;
  }

  .stats-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .meal-stat {
    position: relative;
    padding: 20px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);
    overflow: hidden;

    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 8px;
    }

    .stat-label {
      font-size: 14px;
      color: var(--text-secondary);
      margin-bottom: 12px;
    }

    .stat-icon {
      position: absolute;
      top: 20px;
      right: 20px;
      font-size: 48px;
      opacity: 0.1;
    }

    &.stat-primary .stat-icon { color: var(--primary-color); }
    &.stat-success .stat-icon { color: var(--success-color); }
    &.stat-warning .stat-icon { color: var(--warning-color); }
    &.stat-info .stat-icon { color: var(--info-color); }
  }

  .today-card,
  .records-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .dish-tag {
    margin-right: 4px;
    margin-bottom: 4px;
  }

  .dishes-input {
    width: 100%;

    .dish-item {
      display: flex;
      gap: 8px;
      margin-bottom: 8px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

@media (max-width: 768px) {
  .dining-management {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }
  }
}
</style>
