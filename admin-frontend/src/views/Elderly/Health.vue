<template>
  <div class="elderly-health">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="handleBack" circle>
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <h2 class="page-title">
          <el-icon><FirstAidKit /></el-icon>
          健康档案
        </h2>
        <span class="elderly-name">{{ elderlyName }}</span>
      </div>
      <div class="header-actions">
        <el-button @click="handleAddRecord">
          <el-icon><Plus /></el-icon>
          添加记录
        </el-button>
        <ExportButton
          :data="healthRecords"
          :columns="exportColumns"
          :total="total"
          filename="健康档案"
          :export-function="handleExport"
        />
      </div>
    </div>

    <!-- 健康概览卡片 -->
    <el-row :gutter="20" class="overview-row">
      <el-col :xs="12" :sm="6" v-for="stat in healthStats" :key="stat.key">
        <div class="health-stat" :class="`stat-${stat.type}`">
          <div class="stat-icon">
            <component :is="stat.icon" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 健康趋势图表 -->
    <el-card shadow="never" class="chart-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">健康趋势</span>
          <el-radio-group v-model="trendMetric" size="small">
            <el-radio-button label="bloodPressure">血压</el-radio-button>
            <el-radio-button label="bloodSugar">血糖</el-radio-button>
            <el-radio-button label="temperature">体温</el-radio-button>
            <el-radio-button label="weight">体重</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <ECharts :option="trendChartOption" height="300px" />
    </el-card>

    <!-- 健康记录列表 -->
    <el-card shadow="never" class="records-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">健康记录明细</span>
          <el-radio-group v-model="filterType" size="small">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="blood_pressure">血压</el-radio-button>
            <el-radio-button label="blood_sugar">血糖</el-radio-button>
            <el-radio-button label="temperature">体温</el-radio-button>
            <el-radio-button label="other">其他</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="filteredRecords"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="recordTime" label="记录时间" width="160" />
        <el-table-column prop="type" label="指标类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="数值" width="120">
          <template #default="{ row }">
            <span :class="{ 'text-danger': row.isAbnormal }">
              {{ row.value }} {{ row.unit }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="recordBy" label="记录人" width="100" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.isAbnormal ? 'danger' : 'success'" size="small">
              {{ row.isAbnormal ? '异常' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text type="primary" @click="handleEdit(row)">编辑</el-button>
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
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  FirstAidKit,
  Plus,
  TrendCharts
} from '@element-plus/icons-vue'
import ECharts from '@/components/Dashboard/ECharts.vue'
import ExportButton from '@/components/Export/ExportButton.vue'
import { getChartColors } from '@/composables/useECharts'
import { getHealthRecords, getElderlyList } from '@/utils/seedData'

const route = useRoute()
const router = useRouter()

const elderlyId = computed(() => route.params.id as string)
const elderlyName = ref('加载中...')
const loading = ref(false)
const filterType = ref('')
const trendMetric = ref('bloodPressure')

// 健康统计数据
const healthStats = ref([
  { key: 'records', label: '健康记录', value: 156, type: 'primary', icon: 'Document' },
  { key: 'abnormal', label: '异常次数', value: 8, type: 'danger', icon: 'Warning' },
  { key: 'recent', label: '本月测量', value: 24, type: 'success', icon: 'Calendar' },
  { key: 'score', label: '健康评分', value: 85, type: 'info', icon: 'TrendCharts' }
])

// 健康记录列表
const healthRecords = ref<any[]>([])
const total = ref(0)

const pagination = ref({
  page: 1,
  pageSize: 20
})

// 过滤后的记录
const filteredRecords = computed(() => {
  if (!filterType.value) return healthRecords.value
  return healthRecords.value.filter(r => r.type === filterType.value)
})

// 导出列配置
const exportColumns = [
  { key: 'recordTime', title: '记录时间', width: 18 },
  { key: 'type', title: '指标类型', width: 12 },
  { key: 'value', title: '数值', width: 12 },
  { key: 'unit', title: '单位', width: 8 },
  { key: 'recordBy', title: '记录人', width: 12 },
  { key: 'isAbnormal', title: '是否异常', width: 12, formatter: (row: any) => row.isAbnormal ? '是' : '否' },
  { key: 'remark', title: '备注', width: 20 }
]

// 健康趋势图表
const trendChartOption = computed(() => {
  const colors = getChartColors()
  return {
    grid: {
      top: '10%',
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: 'transparent',
      textStyle: { color: '#fff' }
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月']
    },
    yAxis: {
      type: 'value',
      name: trendMetric.value === 'bloodPressure' ? 'mmHg' :
            trendMetric.value === 'bloodSugar' ? 'mmol/L' :
            trendMetric.value === 'temperature' ? '℃' : 'kg'
    },
    series: [{
      name: '平均值',
      type: 'line',
      smooth: true,
      data: trendMetric.value === 'bloodPressure' ? [125, 128, 126, 130, 127, 124] :
            trendMetric.value === 'bloodSugar' ? [6.2, 6.5, 6.3, 6.8, 6.4, 6.1] :
            trendMetric.value === 'temperature' ? [36.5, 36.6, 36.4, 36.7, 36.5, 36.3] :
            [65, 66, 64.5, 65.5, 66, 65],
      itemStyle: { color: colors[0] },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: colors[0] + 'cc' },
            { offset: 1, color: colors[0] + '11' }
          ]
        }
      }
    }]
  }
})

// 方法
const getTypeColor = (type: string) => {
  const map: Record<string, any> = {
    blood_pressure: 'danger',
    blood_sugar: 'warning',
    temperature: 'info',
    weight: 'success'
  }
  return map[type] || ''
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    blood_pressure: '血压',
    blood_sugar: '血糖',
    temperature: '体温',
    weight: '体重'
  }
  return map[type] || type
}

const loadData = async () => {
  loading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    // 获取老人信息
    const elderly = getElderlyList().find(e => e.id === elderlyId.value)
    if (elderly) {
      elderlyName.value = elderly.name
    }

    // 模拟健康数据
    healthRecords.value = [
      {
        id: 1,
        type: 'blood_pressure',
        value: '128/82',
        unit: 'mmHg',
        recordTime: '2026-03-03 08:30',
        recordBy: '赵护士',
        isAbnormal: false,
        remark: '血压正常'
      },
      {
        id: 2,
        type: 'blood_sugar',
        value: '6.8',
        unit: 'mmol/L',
        recordTime: '2026-03-03 08:30',
        recordBy: '赵护士',
        isAbnormal: true,
        remark: '空腹血糖略高'
      },
      {
        id: 3,
        type: 'temperature',
        value: '36.5',
        unit: '℃',
        recordTime: '2026-03-03 08:30',
        recordBy: '赵护士',
        isAbnormal: false,
        remark: '体温正常'
      }
    ]
    total.value = 156
  } finally {
    loading.value = false
  }
}

const handleBack = () => {
  router.back()
}

const handleAddRecord = () => {
  console.log('添加健康记录')
}

const handleView = (row: any) => {
  console.log('查看记录', row)
}

const handleEdit = (row: any) => {
  console.log('编辑记录', row)
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadData()
}

const handleExport = async (params: any) => {
  console.log('导出健康数据', params)
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.elderly-health {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 16px;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .page-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .elderly-name {
    font-size: 14px;
    color: var(--text-secondary);
    background: var(--bg-tertiary);
    padding: 4px 12px;
    border-radius: 12px;
  }

  .overview-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .health-stat {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);

    .stat-icon {
      width: 40px;
      height: 40px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 20px;
    }

    &.stat-primary .stat-icon {
      background: var(--gradient-blue);
      color: #fff;
    }

    &.stat-danger .stat-icon {
      background: var(--gradient-red);
      color: #fff;
    }

    &.stat-success .stat-icon {
      background: var(--gradient-green);
      color: #fff;
    }

    &.stat-info .stat-icon {
      background: linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%);
      color: #fff;
    }

    .stat-value {
      font-size: 20px;
      font-weight: 600;
      color: var(--text-primary);
    }

    .stat-label {
      font-size: 12px;
      color: var(--text-secondary);
    }
  }

  .chart-card,
  .records-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
  }

  .card-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .text-danger {
    color: var(--danger-color);
    font-weight: 600;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

@media (max-width: 768px) {
  .elderly-health {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .overview-row {
      :deep(.el-col) {
        margin-bottom: 8px;
      }
    }

    .health-stat {
      padding: 12px;

      .stat-value {
        font-size: 18px;
      }
    }
  }
}
</style>
