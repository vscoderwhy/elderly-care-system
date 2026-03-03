<template>
  <div class="overview-dashboard">
    <!-- 统计卡片行 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6" v-for="stat in stats" :key="stat.key">
        <StatCard
          :label="stat.label"
          :value="stat.value"
          :unit="stat.unit"
          :color="stat.color"
          :trend="stat.trend"
          :show-sparkline="stat.showSparkline"
          :sparkline-data="stat.sparklineData"
        />
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <!-- 入住趋势图 -->
      <el-col :xs="24" :lg="16">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">入住趋势分析</span>
              <el-radio-group v-model="trendPeriod" size="small">
                <el-radio-button label="week">周</el-radio-button>
                <el-radio-button label="month">月</el-radio-button>
                <el-radio-button label="year">年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <ECharts :option="trendChartOption" height="350px" />
        </el-card>
      </el-col>

      <!-- 护理等级分布 -->
      <el-col :xs="24" :lg="8">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <span class="card-title">护理等级分布</span>
          </template>
          <ECharts :option="careLevelChartOption" height="350px" />
        </el-card>
      </el-col>
    </el-row>

    <!-- KPI 指标 -->
    <el-row :gutter="20" class="kpi-row">
      <el-col :xs="24" :sm="8" v-for="kpi in kpiData" :key="kpi.key">
        <el-card shadow="hover" class="kpi-card">
          <KPIGauge
            :value="kpi.value"
            :min="kpi.min"
            :max="kpi.max"
            :title="kpi.title"
            :unit="kpi.unit"
            :ranges="kpi.ranges"
            height="250px"
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细数据表格 -->
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="hover" class="table-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">近期护理记录</span>
              <el-button type="primary" size="small" @click="handleViewAll">
                查看全部
              </el-button>
            </div>
          </template>
          <el-table :data="recentRecords" stripe style="width: 100%">
            <el-table-column prop="elderlyName" label="老人姓名" width="120" />
            <el-table-column prop="careType" label="护理类型" width="120" />
            <el-table-column prop="nurseName" label="护理员" width="120" />
            <el-table-column prop="recordTime" label="记录时间" width="180" />
            <el-table-column prop="description" label="护理内容" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'completed' ? 'success' : 'warning'" size="small">
                  {{ scope.row.status === 'completed' ? '已完成' : '进行中' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import StatCard from '@/components/Dashboard/StatCard.vue'
import KPIGauge from '@/components/Dashboard/KPIGauge.vue'
import ECharts from '@/components/Dashboard/ECharts.vue'
import { getChartColors } from '@/composables/useECharts'
import { getStatistics, getCareRecords, getElderlyList } from '@/utils/seedData'

// 统计卡片数据
const stats = ref([
  {
    key: 'total',
    label: '在院老人',
    value: 0,
    unit: '人',
    color: '#409eff',
    trend: 8.5,
    showSparkline: true,
    sparklineData: [220, 225, 230, 235, 240, 245, 248]
  },
  {
    key: 'occupancy',
    label: '入住率',
    value: 0,
    unit: '%',
    color: '#67c23a',
    trend: 5.2,
    showSparkline: true,
    sparklineData: [75, 76, 78, 79, 80, 81, 82.5]
  },
  {
    key: 'nurses',
    label: '护理员',
    value: 56,
    unit: '人',
    color: '#e6a23c',
    trend: 3.1,
    showSparkline: true,
    sparklineData: [48, 50, 52, 53, 54, 55, 56]
  },
  {
    key: 'tasks',
    label: '今日护理任务',
    value: 0,
    unit: '次',
    color: '#f56c6c',
    trend: -2.3,
    showSparkline: true,
    sparklineData: [1300, 1280, 1270, 1260, 1250, 1248, 1245]
  }
])

// 近期护理记录
const recentRecords = ref([])

// 加载数据
onMounted(() => {
  const statistics = getStatistics()
  const colors = getChartColors()

  stats.value[0].value = Number(statistics.elderly.total)
  stats.value[1].value = Number(statistics.elderly.occupancyRate)
  stats.value[3].value = statistics.care.totalRecords

  // 获取护理等级分布数据
  const byCareLevel = statistics.elderly.byCareLevel
  careLevelChartOption.value.series[0].data = [
    { value: byCareLevel['三级'], name: '三级护理', itemStyle: { color: colors[0] } },
    { value: byCareLevel['二级'], name: '二级护理', itemStyle: { color: colors[1] } },
    { value: byCareLevel['一级'], name: '一级护理', itemStyle: { color: colors[2] } },
    { value: byCareLevel['特级'], name: '特级护理', itemStyle: { color: colors[3] } }
  ]

  // 加载近期护理记录
  const records = getCareRecords()
  recentRecords.value = records.slice(0, 5).map(r => ({
    elderlyName: r.elderlyName,
    careType: r.careType,
    nurseName: r.nurseName,
    recordTime: r.careTime,
    description: r.content,
    status: r.evaluation >= 4 ? 'completed' : 'in-progress'
  }))
})

// KPI 数据
const kpiData = ref([
  {
    key: 'satisfaction',
    title: '家属满意度',
    value: 92.5,
    min: 0,
    max: 100,
    unit: '%',
    ranges: [
      { start: 0, end: 60, color: '#f56c6c' },
      { start: 60, end: 80, color: '#e6a23c' },
      { start: 80, end: 100, color: '#67c23a' }
    ]
  },
  {
    key: 'completion',
    title: '任务完成率',
    value: 88,
    min: 0,
    max: 100,
    unit: '%',
    ranges: [
      { start: 0, end: 70, color: '#f56c6c' },
      { start: 70, end: 85, color: '#e6a23c' },
      { start: 85, end: 100, color: '#67c23a' }
    ]
  },
  {
    key: 'health',
    title: '健康指标达标率',
    value: 78,
    min: 0,
    max: 100,
    unit: '%',
    ranges: [
      { start: 0, end: 60, color: '#f56c6c' },
      { start: 60, end: 80, color: '#e6a23c' },
      { start: 80, end: 100, color: '#67c23a' }
    ]
  }
])

// 趋势图时间段
const trendPeriod = ref('month')


// 入住趋势图配置
const trendChartOption = computed(() => {
  const colors = getChartColors()
  return {
    grid: {
      top: '15%',
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
    legend: {
      data: ['入住人数', '退住人数'],
      bottom: 0
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['1月', '2月', '3月', '4月', '5月', '6月']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '入住人数',
        type: 'line',
        smooth: true,
        data: [180, 195, 210, 225, 235, 248],
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
        },
        itemStyle: { color: colors[0] }
      },
      {
        name: '退住人数',
        type: 'line',
        smooth: true,
        data: [12, 8, 15, 10, 8, 5],
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: colors[3] + 'cc' },
              { offset: 1, color: colors[3] + '11' }
            ]
          }
        },
        itemStyle: { color: colors[3] }
      }
    ]
  }
})

// 护理等级分布图配置
const careLevelChartOption = ref({
  tooltip: {
    trigger: 'item',
    backgroundColor: 'rgba(0, 0, 0, 0.8)',
    borderColor: 'transparent',
    textStyle: {
      color: '#fff'
    },
    formatter: '{b}: {c}人 ({d}%)'
  },
  legend: {
    orient: 'vertical',
    right: '10%',
    top: 'center'
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['35%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: 'inherit',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 20,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: []
    }
  ]
})

const handleViewAll = () => {
  console.log('查看全部护理记录')
}
</script>

<style scoped lang="scss">
.overview-dashboard {
  padding: 20px;

  .stats-row {
    margin-bottom: 20px;
  }

  .charts-row {
    margin-bottom: 20px;
  }

  .kpi-row {
    margin-bottom: 20px;
  }

  .chart-card,
  .table-card {
    height: 100%;

    :deep(.el-card__body) {
      padding: 20px;
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
}

// 响应式适配
@media (max-width: 768px) {
  .overview-dashboard {
    padding: 10px;

    .stats-row,
    .charts-row,
    .kpi-row {
      :deep(.el-col) {
        margin-bottom: 10px;
      }
    }
  }
}
</style>
