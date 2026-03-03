<template>
  <div class="advanced-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">高级数据分析</h2>
      <div class="header-actions">
        <el-button-group>
          <el-button :type="viewMode === 'grid' ? 'primary' : ''" @click="viewMode = 'grid'">
            <el-icon><Grid /></el-icon>
            网格视图
          </el-button>
          <el-button :type="viewMode === 'list' ? 'primary' : ''" @click="viewMode = 'list'">
            <el-icon><List /></el-icon>
            列表视图
          </el-button>
        </el-button-group>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出报告
        </el-button>
      </div>
    </div>

    <!-- 过滤器 -->
    <el-card shadow="never" class="filter-card">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="6" :md="4">
          <el-select v-model="filters.department" placeholder="部门" clearable>
            <el-option label="全部部门" value="" />
            <el-option label="护理部" value="nursing" />
            <el-option label="医务部" value="medical" />
            <el-option label="后勤部" value="logistics" />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="6" :md="4">
          <el-select v-model="filters.timeRange" placeholder="时间范围" clearable>
            <el-option label="最近7天" value="7d" />
            <el-option label="最近30天" value="30d" />
            <el-option label="最近90天" value="90d" />
            <el-option label="最近一年" value="1y" />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 100%"
          />
        </el-col>
        <el-col :xs="24" :sm="12" :md="4">
          <el-button type="primary" @click="handleFilter">应用筛选</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 图表网格 -->
    <div class="charts-grid" :class="`view-${viewMode}`">
      <!-- 护理质量雷达图 -->
      <el-card shadow="hover" class="chart-item">
        <template #header>
          <div class="card-header">
            <span class="card-title">护理质量评分</span>
            <el-button size="small" text>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
          </div>
        </template>
        <RadarChart
          :data="radarData"
          :indicator="radarIndicator"
          height="350px"
          title="护理质量多维度评估"
        />
      </el-card>

      <!-- 老人流转桑基图 -->
      <el-card shadow="hover" class="chart-item chart-item--large">
        <template #header>
          <div class="card-header">
            <span class="card-title">老人流转路径</span>
            <el-button size="small" text>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
          </div>
        </template>
        <SankeyChart
          :nodes="sankeyNodes"
          :links="sankeyLinks"
          height="400px"
          title="老人状态流转分析"
        />
      </el-card>

      <!-- 费用支付漏斗图 -->
      <el-card shadow="hover" class="chart-item">
        <template #header>
          <div class="card-header">
            <span class="card-title">费用支付转化</span>
            <el-button size="small" text>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
          </div>
        </template>
        <FunnelChart
          :data="funnelData"
          height="350px"
          title="费用支付流程转化率"
        />
      </el-card>

      <!-- 关系网络图 -->
      <el-card shadow="hover" class="chart-item chart-item--large">
        <template #header>
          <div class="card-header">
            <span class="card-title">老人关系网络</span>
            <el-button size="small" text>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
          </div>
        </template>
        <GraphChart
          :nodes="graphNodes"
          :links="graphLinks"
          :categories="graphCategories"
          height="450px"
          title="老人-家属-护工关系网络"
          layout="force"
        />
      </el-card>

      <!-- 护理记录词云 -->
      <el-card shadow="hover" class="chart-item">
        <template #header>
          <div class="card-header">
            <span class="card-title">护理记录关键词</span>
            <el-button size="small" text>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
          </div>
        </template>
        <WordCloud
          :data="wordCloudData"
          height="350px"
          title="护理记录高频词分析"
        />
      </el-card>

      <!-- 健康数据趋势 -->
      <el-card shadow="hover" class="chart-item chart-item--wide">
        <template #header>
          <div class="card-header">
            <span class="card-title">健康数据趋势</span>
            <el-radio-group v-model="healthMetric" size="small">
              <el-radio-button label="bloodPressure">血压</el-radio-button>
              <el-radio-button label="bloodSugar">血糖</el-radio-button>
              <el-radio-button label="temperature">体温</el-radio-button>
            </el-radio-group>
          </div>
        </template>
        <ECharts :option="healthTrendOption" height="300px" />
      </el-card>
    </div>

    <!-- 数据表格 -->
    <el-card shadow="hover" class="table-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">详细数据</span>
          <el-button size="small" @click="handleExportTable">
            <el-icon><Download /></el-icon>
            导出Excel
          </el-button>
        </div>
      </template>
      <el-table :data="tableData" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="category" label="类别" width="120" />
        <el-table-column prop="indicator" label="指标" width="150" />
        <el-table-column prop="currentValue" label="当前值" width="120" />
        <el-table-column prop="previousValue" label="上期值" width="120" />
        <el-table-column prop="change" label="变化" width="100">
          <template #default="scope">
            <span :class="scope.row.change >= 0 ? 'text-success' : 'text-danger'">
              {{ scope.row.change >= 0 ? '+' : '' }}{{ scope.row.change }}%
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="trend" label="趋势" width="100">
          <template #default="scope">
            <el-tag :type="getTrendType(scope.row.trend)" size="small">
              {{ scope.row.trend }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updateTime" label="更新时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button size="small" text @click="handleViewDetail(scope.row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Grid, List, Download, MoreFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import RadarChart from '@/components/Charts/RadarChart.vue'
import SankeyChart from '@/components/Charts/SankeyChart.vue'
import FunnelChart from '@/components/Charts/FunnelChart.vue'
import GraphChart from '@/components/Charts/GraphChart.vue'
import WordCloud from '@/components/Charts/WordCloud.vue'
import ECharts from '@/components/Dashboard/ECharts.vue'
import { getChartColors } from '@/composables/useECharts'
import { exportToExcel, exportToJSON } from '@/utils/export'
import { getStatistics, getElderlyList, getCareRecords } from '@/utils/seedData'

// 视图模式
const viewMode = ref<'grid' | 'list'>('grid')
const loading = ref(false)
const healthMetric = ref('bloodPressure')

// 过滤器
const filters = ref({
  department: '',
  timeRange: '30d',
  dateRange: []
})

// 雷达图数据
const radarIndicator = [
  { name: '服务态度', max: 100 },
  { name: '专业技能', max: 100 },
  { name: '响应速度', max: 100 },
  { name: '沟通能力', max: 100 },
  { name: '责任心', max: 100 },
  { name: '团队协作', max: 100 }
]

const radarData = [
  { name: '本月', value: [85, 90, 78, 88, 92, 80] },
  { name: '上月', value: [80, 85, 75, 82, 88, 78] }
]

// 桑基图数据
const sankeyNodes = [
  { name: '自理' },
  { name: '半自理' },
  { name: '不能自理' },
  { name: '一级护理' },
  { name: '二级护理' },
  { name: '三级护理' },
  { name: '特级护理' },
  { name: '出院' }
]

const sankeyLinks = [
  { source: '自理', target: '三级护理', value: 45 },
  { source: '自理', target: '二级护理', value: 12 },
  { source: '半自理', target: '二级护理', value: 38 },
  { source: '半自理', target: '一级护理', value: 25 },
  { source: '不能自理', target: '一级护理', value: 18 },
  { source: '不能自理', target: '特级护理', value: 32 },
  { source: '一级护理', target: '出院', value: 8 },
  { source: '二级护理', target: '出院', value: 15 },
  { source: '三级护理', target: '出院', value: 22 }
]

// 漏斗图数据
const funnelData = [
  { name: '账单生成', value: 1000 },
  { name: '发送通知', value: 950 },
  { name: '查看账单', value: 780 },
  { name: '发起支付', value: 520 },
  { name: '支付完成', value: 485 }
]

// 关系图数据
const graphCategories = [
  { name: '老人' },
  { name: '家属' },
  { name: '护工' }
]

const graphNodes = [
  { id: '1', name: '张奶奶', category: 0, symbolSize: 20, value: 85 },
  { id: '2', name: '王爷爷', category: 0, symbolSize: 18, value: 72 },
  { id: '3', name: '李奶奶', category: 0, symbolSize: 22, value: 90 },
  { id: '4', name: '张先生（家属）', category: 1, symbolSize: 15 },
  { id: '5', name: '王女士（家属）', category: 1, symbolSize: 15 },
  { id: '6', name: '李先生（家属）', category: 1, symbolSize: 15 },
  { id: '7', name: '赵护士', category: 2, symbolSize: 16 },
  { id: '8', name: '陈康复师', category: 2, symbolSize: 16 },
  { id: '9', name: '周护工', category: 2, symbolSize: 16 }
]

const graphLinks = [
  { source: '1', target: '4', value: 1 },
  { source: '2', target: '5', value: 1 },
  { source: '3', target: '6', value: 1 },
  { source: '1', target: '7', value: 2 },
  { source: '2', target: '8', value: 2 },
  { source: '3', target: '9', value: 2 },
  { source: '7', target: '8', value: 1 },
  { source: '8', target: '9', value: 1 }
]

// 词云数据
const wordCloudData = [
  { name: '血压监测', value: 156 },
  { name: '康复训练', value: 142 },
  { name: '营养配餐', value: 128 },
  { name: '心理疏导', value: 115 },
  { name: '健康宣教', value: 108 },
  { name: '用药管理', value: 95 },
  { name: '生活照料', value: 88 },
  { name: '社交活动', value: 76 },
  { name: '紧急救护', value: 52 },
  { name: '慢病管理', value: 48 },
  { name: '中医理疗', value: 42 },
  { name: '文娱活动', value: 38 },
  { name: '家属沟通', value: 35 },
  { name: '定期巡房', value: 32 },
  { name: '协助用餐', value: 28 }
]

// 健康趋势图配置
const healthTrendOption = computed(() => {
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
      textStyle: {
        color: '#fff'
      }
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value',
      name: healthMetric.value === 'bloodPressure' ? 'mmHg' :
            healthMetric.value === 'bloodSugar' ? 'mmol/L' : '℃'
    },
    series: [
      {
        name: '平均值',
        type: 'line',
        smooth: true,
        data: healthMetric.value === 'bloodPressure' ? [125, 128, 126, 130, 127, 124, 126] :
              healthMetric.value === 'bloodSugar' ? [6.2, 6.5, 6.3, 6.8, 6.4, 6.1, 6.3] :
              [36.5, 36.6, 36.4, 36.7, 36.5, 36.3, 36.5],
        itemStyle: { color: colors[0] },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: colors[0] + 'cc' },
              { offset: 1, color: colors[0] + '11' }
            ]
          }
        }
      }
    ]
  }
})

// 表格数据
const tableData = ref([
  {
    category: '护理质量',
    indicator: '满意度评分',
    currentValue: 92.5,
    previousValue: 90.2,
    change: 2.55,
    trend: '上升',
    updateTime: '2026-03-03 16:00'
  },
  {
    category: '护理质量',
    indicator: '任务完成率',
    currentValue: 88,
    previousValue: 85,
    change: 3.53,
    trend: '上升',
    updateTime: '2026-03-03 16:00'
  },
  {
    category: '运营指标',
    indicator: '入住率',
    currentValue: 82.5,
    previousValue: 79.8,
    change: 3.38,
    trend: '上升',
    updateTime: '2026-03-03 16:00'
  },
  {
    category: '运营指标',
    indicator: '床位周转率',
    currentValue: 15.2,
    previousValue: 16.5,
    change: -7.88,
    trend: '下降',
    updateTime: '2026-03-03 16:00'
  },
  {
    category: '健康指标',
    indicator: '健康达标率',
    currentValue: 78,
    previousValue: 75,
    change: 4.0,
    trend: '上升',
    updateTime: '2026-03-03 16:00'
  }
])

// 方法
const handleExport = () => {
  // 导出完整分析报告
  const exportData = tableData.value.map(item => ({
    类别: item.category,
    指标: item.indicator,
    当前值: item.currentValue,
    上期值: item.previousValue,
    变化: item.change,
    趋势: item.trend,
    更新时间: item.updateTime
  }))

  exportToExcel(
    exportData,
    [{ key: 'category', title: '类别', width: 12 },
     { key: 'indicator', title: '指标', width: 15 },
     { key: 'currentValue', title: '当前值', width: 12 },
     { key: 'previousValue', title: '上期值', width: 12 },
     { key: 'change', title: '变化%', width: 10 },
     { key: 'trend', title: '趋势', width: 10 },
     { key: 'updateTime', title: '更新时间', width: 18 }],
    '高级数据分析报告'
  )
  ElMessage.success('分析报告导出成功')
}

const handleFilter = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    ElMessage.success('筛选已应用')
  }, 500)
}

const handleExportTable = () => {
  const columns = [
    { key: 'category', title: '类别', width: 12 },
    { key: 'indicator', title: '指标', width: 15 },
    { key: 'currentValue', title: '当前值', width: 12 },
    { key: 'previousValue', title: '上期值', width: 12 },
    { key: 'change', title: '变化', width: 10 },
    { key: 'trend', title: '趋势', width: 10 },
    { key: 'updateTime', title: '更新时间', width: 18 }
  ]
  exportToExcel(tableData.value, columns, '详细数据表')
  ElMessage.success('Excel导出成功')
}

const handleViewDetail = (row: any) => {
  console.log('查看详情', row)
}

const getTrendType = (trend: string) => {
  return trend === '上升' ? 'success' : trend === '下降' ? 'danger' : 'info'
}
</script>

<style scoped lang="scss">
.advanced-dashboard {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 16px;
  }

  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }

  .filter-card {
    margin-bottom: 20px;

    :deep(.el-card__body) {
      padding: 16px;
    }
  }

  .charts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: 20px;
    margin-bottom: 20px;

    &.view-list {
      display: flex;
      flex-direction: column;
    }

    .chart-item {
      min-height: 400px;

      &--large {
        grid-column: span 2;
      }

      &--wide {
        grid-column: span 2;
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .card-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .table-card {
    margin-top: 20px;
  }

  .text-success {
    color: var(--success-color);
  }

  .text-danger {
    color: var(--danger-color);
  }
}

// 响应式适配
@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: 1fr !important;

    .chart-item--large,
    .chart-item--wide {
      grid-column: span 1;
    }
  }
}

@media (max-width: 768px) {
  .advanced-dashboard {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .charts-grid {
      grid-template-columns: 1fr;
      gap: 10px;
    }
  }
}
</style>
