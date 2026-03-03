<template>
  <div class="command-center">
    <!-- 顶部标题栏 -->
    <div class="header">
      <div class="title-left">
        <div class="title-icon"></div>
        <h1 class="title">养老院智能指挥中心</h1>
      </div>
      <div class="title-center">
        <div class="current-time">{{ currentTime }}</div>
        <div class="date">{{ currentDate }}</div>
      </div>
      <div class="title-right">
        <div class="weather">
          <span class="weather-icon">☀️</span>
          <span class="weather-text">23°C 晴</span>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧列 -->
      <div class="left-column">
        <!-- 在院老人统计 -->
        <div class="panel">
          <div class="panel-title">在院老人概况</div>
          <div class="panel-content">
            <div class="stat-grid">
              <div class="stat-item" v-for="item in elderlyStats" :key="item.key">
                <div class="stat-value" :class="`stat-${item.type}`">{{ item.value }}</div>
                <div class="stat-label">{{ item.label }}</div>
                <div class="stat-trend" :class="item.trend > 0 ? 'up' : 'down'">
                  <span>{{ item.trend > 0 ? '↑' : '↓' }}</span>
                  <span>{{ Math.abs(item.trend) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 护理等级分布 -->
        <div class="panel">
          <div class="panel-title">护理等级分布</div>
          <div class="panel-content">
            <div ref="careLevelChart" class="chart"></div>
          </div>
        </div>

        <!-- 实时告警 -->
        <div class="panel">
          <div class="panel-title">
            <span>实时告警</span>
            <span class="alert-count">{{ alerts.length }}</span>
          </div>
          <div class="panel-content">
            <div class="alert-list">
              <div
                v-for="alert in alerts"
                :key="alert.id"
                class="alert-item"
                :class="`alert-${alert.level}`"
              >
                <div class="alert-icon">⚠️</div>
                <div class="alert-content">
                  <div class="alert-title">{{ alert.title }}</div>
                  <div class="alert-time">{{ alert.time }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 中间列 -->
      <div class="center-column">
        <!-- 核心指标数字 -->
        <div class="core-metrics">
          <div class="metric-card" v-for="metric in coreMetrics" :key="metric.key">
            <div class="metric-icon">{{ metric.icon }}</div>
            <div class="metric-value">{{ metric.value }}</div>
            <div class="metric-label">{{ metric.label }}</div>
          </div>
        </div>

        <!-- 建筑监控区域 -->
        <div class="building-panel">
          <div class="building-header">
            <span class="building-title">院区实时监控</span>
            <div class="building-tabs">
              <span
                v-for="tab in buildingTabs"
                :key="tab.key"
                :class="['tab-item', { active: currentTab === tab.key }]"
                @click="currentTab = tab.key"
              >
                {{ tab.label }}
              </span>
            </div>
          </div>
          <div class="building-content">
            <div ref="buildingChart" class="building-chart"></div>
            <div class="building-overlay">
              <div class="overlay-item" v-for="item in overlayItems" :key="item.id">
                <div class="overlay-dot" :class="`dot-${item.status}`"></div>
                <div class="overlay-label">{{ item.label }}</div>
                <div class="overlay-value">{{ item.value }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 实时任务进度 -->
        <div class="panel">
          <div class="panel-title">今日任务完成进度</div>
          <div class="panel-content">
            <div class="task-progress">
              <div class="progress-item" v-for="task in taskProgress" :key="task.type">
                <div class="progress-header">
                  <span class="progress-label">{{ task.type }}</span>
                  <span class="progress-value">{{ task.completed }}/{{ task.total }}</span>
                </div>
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{ width: `${(task.completed / task.total) * 100}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧列 -->
      <div class="right-column">
        <!-- 财务概览 -->
        <div class="panel">
          <div class="panel-title">财务概览</div>
          <div class="panel-content">
            <div ref="financeChart" class="chart"></div>
          </div>
        </div>

        <!-- 健康数据趋势 -->
        <div class="panel">
          <div class="panel-title">健康数据监测</div>
          <div class="panel-content">
            <div ref="healthChart" class="chart"></div>
          </div>
        </div>

        <!-- 护理质量雷达图 -->
        <div class="panel">
          <div class="panel-title">护理质量评估</div>
          <div class="panel-content">
            <div ref="qualityChart" class="chart"></div>
          </div>
        </div>

        <!-- 员工工作状态 -->
        <div class="panel">
          <div class="panel-title">员工工作状态</div>
          <div class="panel-content">
            <div class="staff-grid">
              <div
                v-for="staff in staffStatus"
                :key="staff.id"
                class="staff-item"
                :class="`staff-${staff.status}`"
              >
                <div class="staff-avatar">{{ staff.name[0] }}</div>
                <div class="staff-info">
                  <div class="staff-name">{{ staff.name }}</div>
                  <div class="staff-dept">{{ staff.dept }}</div>
                </div>
                <div class="staff-status">{{ statusText[staff.status] }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts/core'
import {
  CanvasRenderer
} from 'echarts/renderers'
import {
  PieChart,
  BarChart,
  LineChart,
  RadarChart
} from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'

echarts.use([
  CanvasRenderer,
  PieChart,
  BarChart,
  LineChart,
  RadarChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

// 时间更新
const currentTime = ref('')
const currentDate = ref('')
const currentTab = ref('overview')

let timeInterval: any

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
  currentDate.value = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    weekday: 'long'
  })
}

// 图表refs
const careLevelChart = ref<HTMLElement>()
const buildingChart = ref<HTMLElement>()
const financeChart = ref<HTMLElement>()
const healthChart = ref<HTMLElement>()
const qualityChart = ref<HTMLElement>()

let chartInstances: any[] = []

onMounted(async () => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)

  await nextTick()
  initCharts()
})

onUnmounted(() => {
  clearInterval(timeInterval)
  chartInstances.forEach(chart => chart?.dispose())
})

const initCharts = () => {
  // 护理等级饼图
  if (careLevelChart.value) {
    const chart = echarts.init(careLevelChart.value)
    chart.setOption({
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        borderColor: 'transparent',
        textStyle: { color: '#fff' },
        formatter: '{b}: {c}人 ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center',
        textStyle: { color: '#fff', fontSize: 12 }
      },
      series: [{
        type: 'pie',
        radius: ['45%', '75%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#0d1b2a',
          borderWidth: 3
        },
        label: { show: false },
        emphasis: {
          label: {
            show: true,
            fontSize: 24,
            fontWeight: 'bold',
            color: '#fff'
          }
        },
        labelLine: { show: false },
        data: [
          { value: 8, name: '特级护理', itemStyle: { color: '#ff6b6b' } },
          { value: 18, name: '一级护理', itemStyle: { color: '#ffd93d' } },
          { value: 16, name: '二级护理', itemStyle: { color: '#6bcb77' } },
          { value: 8, name: '三级护理', itemStyle: { color: '#4d96ff' } }
        ]
      }]
    })
    chartInstances.push(chart)
  }

  // 建筑柱状图
  if (buildingChart.value) {
    const chart = echarts.init(buildingChart.value)
    chart.setOption({
      grid: { top: 30, right: 30, bottom: 30, left: 30 },
      xAxis: {
        type: 'category',
        data: ['1号楼', '2号楼', '3号楼', '4号楼'],
        axisLine: { lineStyle: { color: '#1e3a5f' } },
        axisLabel: { color: '#4a90a4', fontSize: 14 }
      },
      yAxis: {
        type: 'value',
        max: 100,
        axisLine: { lineStyle: { color: '#1e3a5f' } },
        axisLabel: { color: '#4a90a4', fontSize: 12 },
        splitLine: { lineStyle: { color: '#1e3a5f', type: 'dashed' } }
      },
      series: [{
        type: 'bar',
        barWidth: '50%',
        data: [
          { value: 85, itemStyle: { color: '#00d4ff' } },
          { value: 72, itemStyle: { color: '#00d4ff' } },
          { value: 90, itemStyle: { color: '#ff6b6b' } },
          { value: 68, itemStyle: { color: '#00d4ff' } }
        ],
        label: {
          show: true,
          position: 'top',
          color: '#00d4ff',
          fontSize: 16,
          fontWeight: 'bold'
        }
      }]
    })
    chartInstances.push(chart)
  }

  // 财务趋势图
  if (financeChart.value) {
    const chart = echarts.init(financeChart.value)
    chart.setOption({
      grid: { top: 30, right: 20, bottom: 30, left: 50 },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月'],
        axisLine: { lineStyle: { color: '#1e3a5f' } },
        axisLabel: { color: '#4a90a4', fontSize: 11 }
      },
      yAxis: {
        type: 'value',
        axisLine: { lineStyle: { color: '#1e3a5f' } },
        axisLabel: { color: '#4a90a4', fontSize: 11 },
        splitLine: { lineStyle: { color: '#1e3a5f', type: 'dashed' } }
      },
      series: [
        {
          name: '收入',
          type: 'line',
          smooth: true,
          data: [180, 220, 250, 280, 310, 350],
          itemStyle: { color: '#00d4ff' },
          areaStyle: { color: 'rgba(0, 212, 255, 0.2)' }
        },
        {
          name: '支出',
          type: 'line',
          smooth: true,
          data: [120, 140, 160, 175, 190, 210],
          itemStyle: { color: '#ff6b6b' },
          areaStyle: { color: 'rgba(255, 107, 107, 0.2)' }
        }
      ],
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        borderColor: 'transparent',
        textStyle: { color: '#fff' }
      }
    })
    chartInstances.push(chart)
  }

  // 健康趋势图
  if (healthChart.value) {
    const chart = echarts.init(healthChart.value)
    chart.setOption({
      grid: { top: 20, right: 20, bottom: 30, left: 40 },
      xAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
        axisLine: { lineStyle: { color: '#1e3a5f' } },
        axisLabel: { color: '#4a90a4', fontSize: 10 }
      },
      yAxis: {
        type: 'value',
        max: 100,
        axisLine: { lineStyle: { color: '#1e3a5f' } },
        axisLabel: { color: '#4a90a4', fontSize: 10 },
        splitLine: { lineStyle: { color: '#1e3a5f', type: 'dashed' } }
      },
      series: [
        { name: '血压', type: 'line', smooth: true, data: [92, 94, 91, 95, 93, 96, 94], itemStyle: { color: '#00d4ff' } },
        { name: '血糖', type: 'line', smooth: true, data: [88, 90, 89, 92, 91, 93, 91], itemStyle: { color: '#6bcb77' } },
        { name: '心率', type: 'line', smooth: true, data: [95, 96, 94, 97, 95, 98, 96], itemStyle: { color: '#ffd93d' } }
      ],
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        borderColor: 'transparent',
        textStyle: { color: '#fff' }
      }
    })
    chartInstances.push(chart)
  }

  // 雷达图
  if (qualityChart.value) {
    const chart = echarts.init(qualityChart.value)
    chart.setOption({
      radar: {
        center: ['50%', '55%'],
        radius: '65%',
        indicator: [
          { name: '服务态度', max: 100 },
          { name: '专业技能', max: 100 },
          { name: '响应速度', max: 100 },
          { name: '卫生环境', max: 100 },
          { name: '安全管理', max: 100 },
          { name: '记录规范', max: 100 }
        ],
        axisName: { color: '#4a90a4', fontSize: 12 },
        splitLine: { lineStyle: { color: '#1e3a5f' } },
        splitArea: { show: true, areaStyle: { color: ['rgba(0, 212, 255, 0.05)', 'rgba(0, 212, 255, 0.1)'] } },
        axisLine: { lineStyle: { color: '#1e3a5f' } }
      },
      series: [{
        type: 'radar',
        data: [
          {
            value: [95, 88, 92, 90, 96, 89],
            name: '本月',
            itemStyle: { color: '#00d4ff' },
            areaStyle: { color: 'rgba(0, 212, 255, 0.3)' }
          },
          {
            value: [90, 85, 88, 87, 93, 86],
            name: '上月',
            itemStyle: { color: '#ff6b6b' },
            areaStyle: { color: 'rgba(255, 107, 107, 0.2)' }
          }
        ]
      }],
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        borderColor: 'transparent',
        textStyle: { color: '#fff' }
      },
      legend: {
        data: ['本月', '上月'],
        bottom: 5,
        textStyle: { color: '#4a90a4', fontSize: 12 }
      }
    })
    chartInstances.push(chart)
  }
}

// 左侧数据
const elderlyStats = ref([
  { key: 'total', label: '在院老人', value: 50, type: 'primary', trend: 8.5 },
  { key: 'today', label: '今日入住', value: 3, type: 'success', trend: 50 },
  { key: 'leave', label: '请假外出', value: 5, type: 'warning', trend: -20 },
  { key: 'emergency', label: '紧急关注', value: 2, type: 'danger', trend: 0 }
])

const alerts = ref([
  { id: 1, title: '张奶奶血压异常', time: '14:30', level: 'high' },
  { id: 2, title: '3号楼201呼叫响应', time: '14:25', level: 'medium' },
  { id: 3, title: '李护士任务超时', time: '14:20', level: 'low' },
  { id: 4, title: '王爷爷心率不齐', time: '14:15', level: 'high' }
])

// 中间数据
const coreMetrics = ref([
  { key: 'occupancy', label: '入住率', value: '75%', icon: '📊' },
  { key: 'satisfaction', label: '满意度', value: '96%', icon: '⭐' },
  { key: 'tasks', label: '今日任务', value: '124', icon: '✅' },
  { key: 'health', label: '健康达标', value: '89%', icon: '💚' }
])

const buildingTabs = [
  { key: 'overview', label: '全景' },
  { key: 'building1', label: '1号楼' },
  { key: 'building2', label: '2号楼' },
  { key: 'building3', label: '3号楼' },
  { key: 'building4', label: '4号楼' }
]

const overlayItems = ref([
  { id: 1, label: '1号楼', value: '85%', status: 'normal' },
  { id: 2, label: '2号楼', value: '72%', status: 'normal' },
  { id: 3, label: '3号楼', value: '90%', status: 'busy' },
  { id: 4, label: '4号楼', value: '68%', status: 'normal' }
])

const taskProgress = ref([
  { type: '日常护理', completed: 112, total: 124 },
  { type: '健康监测', completed: 45, total: 50 },
  { type: '康复训练', completed: 28, total: 30 },
  { type: '医疗护理', completed: 18, total: 20 }
])

// 右侧数据
const staffStatus = ref([
  { id: 1, name: '赵护士', dept: '护理部', status: 'working' },
  { id: 2, name: '李护士', dept: '护理部', status: 'break' },
  { id: 3, name: '周护士', dept: '护理部', status: 'working' },
  { id: 4, name: '吴护士', dept: '护理部', status: 'off' },
  { id: 5, name: '郑康复师', dept: '康复科', status: 'working' }
])

const statusText: Record<string, string> = {
  working: '工作中',
  break: '休息中',
  off: '下班'
}
</script>

<style scoped lang="scss">
.command-center {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0e27 0%, #1a1f3a 50%, #0d1b2a 100%);
  color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
  overflow-x: hidden;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 40px;
    background: linear-gradient(180deg, rgba(0, 212, 255, 0.1) 0%, transparent 100%);
    border-bottom: 1px solid rgba(0, 212, 255, 0.3);

    .title-left {
      display: flex;
      align-items: center;
      gap: 15px;

      .title-icon {
        width: 40px;
        height: 40px;
        background: linear-gradient(135deg, #00d4ff 0%, #0099cc 100%);
        border-radius: 8px;
        box-shadow: 0 0 20px rgba(0, 212, 255, 0.5);
        animation: pulse 2s ease-in-out infinite;
      }

      .title {
        font-size: 32px;
        font-weight: bold;
        background: linear-gradient(90deg, #00d4ff 0%, #fff 50%, #00d4ff 100%);
        background-size: 200% auto;
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;
        animation: shimmer 3s linear infinite;
      }
    }

    .title-center {
      text-align: center;

      .current-time {
        font-size: 36px;
        font-weight: bold;
        color: #00d4ff;
        text-shadow: 0 0 20px rgba(0, 212, 255, 0.5);
        font-family: 'Courier New', monospace;
      }

      .date {
        font-size: 14px;
        color: #4a90a4;
        margin-top: 5px;
      }
    }

    .title-right {
      .weather {
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 18px;
        color: #ffd93d;

        .weather-icon {
          font-size: 28px;
          animation: rotate 10s linear infinite;
        }
      }
    }
  }

  .main-content {
    display: grid;
    grid-template-columns: 28% 44% 28%;
    gap: 20px;
    padding: 20px;
  }

  .panel {
    background: linear-gradient(135deg, rgba(0, 212, 255, 0.05) 0%, rgba(0, 153, 204, 0.05) 100%);
    border: 1px solid rgba(0, 212, 255, 0.3);
    border-radius: 12px;
    margin-bottom: 20px;
    overflow: hidden;
    box-shadow: 0 0 30px rgba(0, 212, 255, 0.1),
                inset 0 0 30px rgba(0, 212, 255, 0.02);
    animation: fadeIn 0.5s ease-out;

    .panel-title {
      padding: 15px 20px;
      font-size: 16px;
      font-weight: bold;
      color: #00d4ff;
      border-bottom: 1px solid rgba(0, 212, 255, 0.2);
      display: flex;
      justify-content: space-between;
      align-items: center;

      .alert-count {
        background: #ff6b6b;
        color: #fff;
        padding: 2px 8px;
        border-radius: 10px;
        font-size: 12px;
        animation: pulse 2s ease-in-out infinite;
      }
    }

    .panel-content {
      padding: 20px;

      .chart {
        height: 200px;
      }
    }
  }

  .stat-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;

    .stat-item {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.1) 0%, rgba(0, 153, 204, 0.05) 100%);
      border: 1px solid rgba(0, 212, 255, 0.2);
      border-radius: 10px;
      padding: 15px;
      text-align: center;

      .stat-value {
        font-size: 32px;
        font-weight: bold;
        margin-bottom: 8px;

        &.stat-primary { color: #00d4ff; text-shadow: 0 0 10px rgba(0, 212, 255, 0.5); }
        &.stat-success { color: #6bcb77; text-shadow: 0 0 10px rgba(107, 203, 119, 0.5); }
        &.stat-warning { color: #ffd93d; text-shadow: 0 0 10px rgba(255, 217, 61, 0.5); }
        &.stat-danger { color: #ff6b6b; text-shadow: 0 0 10px rgba(255, 107, 107, 0.5); }
      }

      .stat-label {
        font-size: 12px;
        color: #4a90a4;
        margin-bottom: 5px;
      }

      .stat-trend {
        font-size: 11px;
        font-weight: bold;

        &.up { color: #6bcb77; }
        &.down { color: #ff6b6b; }
      }
    }
  }

  .alert-list {
    max-height: 200px;
    overflow-y: auto;

    .alert-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      background: rgba(255, 107, 107, 0.05);
      border-left: 3px solid;
      margin-bottom: 10px;
      border-radius: 0 8px 8px 0;
      animation: slideIn 0.3s ease-out;

      &.alert-high { border-left-color: #ff6b6b; background: rgba(255, 107, 107, 0.1); }
      &.alert-medium { border-left-color: #ffd93d; }
      &.alert-low { border-left-color: #6bcb77; }

      .alert-icon { font-size: 20px; }

      .alert-content {
        flex: 1;

        .alert-title {
          font-size: 13px;
          color: #fff;
          margin-bottom: 3px;
        }

        .alert-time {
          font-size: 11px;
          color: #4a90a4;
        }
      }
    }
  }

  .core-metrics {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 15px;
    margin-bottom: 20px;

    .metric-card {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.15) 0%, rgba(0, 153, 204, 0.05) 100%);
      border: 1px solid rgba(0, 212, 255, 0.4);
      border-radius: 12px;
      padding: 20px 15px;
      text-align: center;
      box-shadow: 0 0 30px rgba(0, 212, 255, 0.2);
      transition: all 0.3s;

      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 0 40px rgba(0, 212, 255, 0.4);
      }

      .metric-icon {
        font-size: 28px;
        margin-bottom: 10px;
      }

      .metric-value {
        font-size: 28px;
        font-weight: bold;
        color: #00d4ff;
        text-shadow: 0 0 15px rgba(0, 212, 255, 0.5);
        margin-bottom: 5px;
      }

      .metric-label {
        font-size: 12px;
        color: #4a90a4;
      }
    }
  }

  .building-panel {
    background: linear-gradient(135deg, rgba(0, 212, 255, 0.05) 0%, rgba(0, 153, 204, 0.05) 100%);
    border: 1px solid rgba(0, 212, 255, 0.3);
    border-radius: 12px;
    margin-bottom: 20px;
    overflow: hidden;

    .building-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 15px 20px;
      border-bottom: 1px solid rgba(0, 212, 255, 0.2);

      .building-title {
        font-size: 16px;
        font-weight: bold;
        color: #00d4ff;
      }

      .building-tabs {
        display: flex;
        gap: 5px;

        .tab-item {
          padding: 6px 12px;
          border-radius: 6px;
          font-size: 12px;
          color: #4a90a4;
          cursor: pointer;
          transition: all 0.3s;

          &:hover {
            background: rgba(0, 212, 255, 0.1);
            color: #00d4ff;
          }

          &.active {
            background: linear-gradient(135deg, #00d4ff 0%, #0099cc 100%);
            color: #0a0e27;
            font-weight: bold;
          }
        }
      }
    }

    .building-content {
      position: relative;
      padding: 20px;
      height: 280px;

      .building-chart {
        height: 100%;
      }

      .building-overlay {
        position: absolute;
        top: 40px;
        right: 40px;

        .overlay-item {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 15px;
          font-size: 12px;

          .overlay-dot {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            animation: pulse 2s ease-in-out infinite;

            &.dot-normal { background: #00d4ff; box-shadow: 0 0 10px #00d4ff; }
            &.dot-busy { background: #ff6b6b; box-shadow: 0 0 10px #ff6b6b; }
          }

          .overlay-label {
            color: #4a90a4;
            width: 50px;
          }

          .overlay-value {
            color: #00d4ff;
            font-weight: bold;
            width: 40px;
          }
        }
      }
    }
  }

  .task-progress {
    .progress-item {
      margin-bottom: 15px;

      &:last-child { margin-bottom: 0; }

      .progress-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;
        font-size: 13px;

        .progress-label { color: #4a90a4; }
        .progress-value { color: #00d4ff; font-weight: bold; }
      }

      .progress-bar {
        height: 8px;
        background: rgba(0, 212, 255, 0.1);
        border-radius: 4px;
        overflow: hidden;

        .progress-fill {
          height: 100%;
          background: linear-gradient(90deg, #00d4ff 0%, #0099cc 100%);
          border-radius: 4px;
          transition: width 1s ease-out;
          box-shadow: 0 0 10px rgba(0, 212, 255, 0.5);
        }
      }
    }
  }

  .staff-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
    max-height: 200px;
    overflow-y: auto;

    .staff-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 10px;
      background: rgba(0, 212, 255, 0.05);
      border-radius: 8px;
      border-left: 3px solid;

      &.staff-working { border-left-color: #6bcb77; }
      &.staff-break { border-left-color: #ffd93d; }
      &.staff-off { border-left-color: #4a90a4; }

      .staff-avatar {
        width: 36px;
        height: 36px;
        border-radius: 50%;
        background: linear-gradient(135deg, #00d4ff 0%, #0099cc 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: bold;
        color: #0a0e27;
      }

      .staff-info {
        flex: 1;

        .staff-name {
          font-size: 12px;
          color: #fff;
          margin-bottom: 2px;
        }

        .staff-dept {
          font-size: 10px;
          color: #4a90a4;
        }
      }

      .staff-status {
        font-size: 10px;
        color: #4a90a4;
      }
    }
  }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

@keyframes shimmer {
  0% { background-position: 0% center; }
  100% { background-position: 200% center; }
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideIn {
  from { opacity: 0; transform: translateX(-10px); }
  to { opacity: 1; transform: translateX(0); }
}

::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 212, 255, 0.05);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 212, 255, 0.3);
  border-radius: 3px;

  &:hover {
    background: rgba(0, 212, 255, 0.5);
  }
}

@media (max-width: 1600px) {
  .command-center .main-content {
    grid-template-columns: 30% 40% 30%;
  }
}

@media (max-width: 1200px) {
  .command-center .main-content {
    grid-template-columns: 1fr;
  }
}
</style>
