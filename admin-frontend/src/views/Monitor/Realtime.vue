<template>
  <div class="realtime-monitor">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><VideoCamera /></el-icon>
          实时监控中心
        </h2>
        <el-badge :value="onlineCount" type="success" class="status-badge">
          <el-tag :type="connected ? 'success' : 'danger'" size="large">
            {{ connected ? '已连接' : '未连接' }}
          </el-tag>
        </el-badge>
      </div>
      <div class="header-actions">
        <el-button @click="refreshData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-switch
          v-model="autoRefresh"
          active-text="自动刷新"
          @change="toggleAutoRefresh"
        />
      </div>
    </div>

    <!-- 警告通知区域 -->
    <el-alert
      v-if="alerts.length > 0"
      :title="`当前有 ${alerts.length} 条未处理警报`"
      type="warning"
      show-icon
      closable
      class="alerts-banner"
    >
      <template #default>
        <div class="alerts-list">
          <div
            v-for="alert in alerts.slice(0, 3)"
            :key="alert.id"
            class="alert-item"
            @click="handleAlertClick(alert)"
          >
            <el-icon class="alert-icon"><WarningFilled /></el-icon>
            <span class="alert-message">{{ alert.message }}</span>
            <span class="alert-time">{{ formatTime(alert.time) }}</span>
          </div>
          <el-button v-if="alerts.length > 3" text size="small" @click="handleViewAllAlerts">
            查看全部 {{ alerts.length }} 条警报
          </el-button>
        </div>
      </template>
    </el-alert>

    <!-- 实时统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in realtimeStats" :key="stat.key">
        <div class="realtime-stat" :class="`stat-${stat.type}`">
          <div class="stat-icon">
            <component :is="stat.icon" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
          <div class="stat-indicator" :class="{ 'is-active': stat.active }"></div>
        </div>
      </el-col>
    </el-row>

    <!-- 监控面板 -->
    <el-row :gutter="20" class="monitor-panels">
      <!-- 护理任务进度 -->
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover" class="monitor-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><ListDone /></el-icon>
                护理任务进度
              </span>
              <el-progress
                :percentage="taskProgress"
                :color="getProgressColor(taskProgress)"
                :width="60"
                type="circle"
              />
            </div>
          </template>
          <div class="task-list">
            <div
              v-for="task in recentTasks"
              :key="task.id"
              class="task-item"
              :class="`task-${task.status}`"
            >
              <div class="task-info">
                <div class="task-title">{{ task.title }}</div>
                <div class="task-meta">
                  <span>{{ task.elderly }}</span>
                  <el-divider direction="vertical" />
                  <span>{{ task.nurse }}</span>
                </div>
              </div>
              <div class="task-status">
                <el-tag :type="getTaskStatusType(task.status)" size="small">
                  {{ getTaskStatusText(task.status) }}
                </el-tag>
                <span class="task-time">{{ formatTime(task.time) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 健康数据异常 -->
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover" class="monitor-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Monitor /></el-icon>
                健康数据监测
              </span>
              <el-button size="small" text @click="handleViewAllHealth">
                查看全部
              </el-button>
            </div>
          </template>
          <div class="health-list">
            <div
              v-for="health in healthData"
              :key="health.id"
              class="health-item"
              :class="`health-${health.level}`"
            >
              <div class="health-avatar">
                <el-avatar :src="health.avatar" :size="40">{{ health.elderlyName[0] }}</el-avatar>
              </div>
              <div class="health-info">
                <div class="health-name">{{ health.elderlyName }}</div>
                <div class="health-metrics">
                  <el-tag size="small" :type="getHealthTagType(health.level)">
                    {{ health.metric }}: {{ health.value }}
                  </el-tag>
                  <span class="health-trend" :class="health.trend > 0 ? 'up' : 'down'">
                    <el-icon><component :is="health.trend > 0 ? ArrowUp : ArrowDown" /></el-icon>
                    {{ Math.abs(health.trend) }}%
                  </span>
                </div>
              </div>
              <div class="health-actions">
                <el-button size="small" type="primary" link @click="handleHealthDetail(health)">
                  详情
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 设备状态监控 -->
    <el-row :gutter="20" class="device-panels">
      <el-col :span="24">
        <el-card shadow="hover" class="monitor-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Connection /></el-icon>
                设备在线状态
              </span>
              <el-radio-group v-model="deviceFilter" size="small">
                <el-radio-button label="all">全部</el-radio-button>
                <el-radio-button label="online">在线</el-radio-button>
                <el-radio-button label="offline">离线</el-radio-button>
                <el-radio-button label="warning">告警</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="device-grid">
            <div
              v-for="device in filteredDevices"
              :key="device.id"
              class="device-item"
              :class="`device-${device.status}`"
              @click="handleDeviceClick(device)"
            >
              <div class="device-icon">
                <component :is="getDeviceIcon(device.type)" />
              </div>
              <div class="device-info">
                <div class="device-name">{{ device.name }}</div>
                <div class="device-location">{{ device.location }}</div>
              </div>
              <div class="device-status">
                <el-tag :type="getDeviceStatusType(device.status)" size="small">
                  {{ getDeviceStatusText(device.status) }}
                </el-tag>
              </div>
              <div class="device-indicator" :class="`indicator-${device.status}`"></div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 考勤实时统计 -->
    <el-row :gutter="20" class="attendance-panels">
      <el-col :xs="24" :lg="16">
        <el-card shadow="hover" class="monitor-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><UserFilled /></el-icon>
                今日考勤统计
              </span>
              <el-button size="small" text @click="handleExportAttendance">
                导出报表
              </el-button>
            </div>
          </template>
          <ECharts :option="attendanceChartOption" height="300px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card shadow="hover" class="monitor-card">
          <template #header>
            <span class="card-title">
              <el-icon><Clock /></el-icon>
              排班概览
            </span>
          </template>
          <div class="schedule-list">
            <div
              v-for="shift in shifts"
              :key="shift.id"
              class="shift-item"
              :class="{ 'is-current': shift.isCurrent }"
            >
              <div class="shift-time">{{ shift.timeRange }}</div>
              <div class="shift-info">
                <div class="shift-name">{{ shift.name }}</div>
                <div class="shift-count">
                  在岗 {{ shift.onDuty }}/{{ shift.total }} 人
                </div>
              </div>
              <el-progress
                :percentage="Math.round((shift.onDuty / shift.total) * 100)"
                :stroke-width="6"
              />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  VideoCamera,
  Refresh,
  WarningFilled,
  ListDone,
  Monitor,
  Connection,
  UserFilled,
  Clock,
  ArrowUp,
  ArrowDown
} from '@element-plus/icons-vue'
import ECharts from '@/components/Dashboard/ECharts.vue'
import { useWebSocket, usePolling } from '@/composables/useWebSocket'
import { getChartColors } from '@/composables/useECharts'

// 连接状态
const connected = ref(false)
const loading = ref(false)
const autoRefresh = ref(true)
const onlineCount = ref(48)
const deviceFilter = ref('all')

// 警报数据
const alerts = ref([
  { id: 1, message: '张奶奶血压异常偏高', time: Date.now() - 300000, level: 'high' },
  { id: 2, message: '3号楼设备离线', time: Date.now() - 600000, level: 'medium' },
  { id: 3, message: '李爷爷未按时服药', time: Date.now() - 900000, level: 'low' }
])

// 实时统计数据
const realtimeStats = ref([
  { key: 'tasks', label: '今日任务', value: 1245, type: 'primary', icon: ListDone, active: true },
  { key: 'completed', label: '已完成', value: 1102, type: 'success', icon: 'Check', active: false },
  { key: 'pending', label: '进行中', value: 98, type: 'warning', icon: Clock, active: true },
  { key: 'overdue', label: '已逾期', value: 45, type: 'danger', icon: WarningFilled, active: true }
])

// 任务进度
const taskProgress = ref(88.5)

// 最近任务
const recentTasks = ref([
  { id: 1, title: '血压测量', elderly: '张奶奶', nurse: '赵护士', time: Date.now() - 120000, status: 'completed' },
  { id: 2, title: '康复训练', elderly: '王爷爷', nurse: '陈康复师', time: Date.now() - 300000, status: 'in-progress' },
  { id: 3, title: '用药管理', elderly: '李奶奶', nurse: '周护士', time: Date.now() - 480000, status: 'pending' },
  { id: 4, title: '营养配餐', elderly: '刘爷爷', nurse: '吴护工', time: Date.now() - 600000, status: 'overdue' },
  { id: 5, title: '心理疏导', elderly: '孙奶奶', nurse: '郑心理咨询师', time: Date.now() - 720000, status: 'completed' }
])

// 健康数据
const healthData = ref([
  { id: 1, elderlyName: '张奶奶', metric: '血压', value: '158/95', level: 'high', trend: 12, avatar: '' },
  { id: 2, elderlyName: '王爷爷', metric: '血糖', value: '8.2', level: 'medium', trend: 5, avatar: '' },
  { id: 3, elderlyName: '李奶奶', metric: '体温', value: '37.8', level: 'low', trend: -3, avatar: '' },
  { id: 4, elderlyName: '刘爷爷', metric: '血氧', value: '94%', level: 'normal', trend: 0, avatar: '' }
])

// 设备数据
const devices = ref([
  { id: 1, name: '智能床垫 #001', type: 'bed', location: '3号楼201', status: 'online' },
  { id: 2, name: '血压监测仪 #003', type: 'monitor', location: '2号楼护士站', status: 'online' },
  { id: 3, name: '紧急呼叫器 #012', type: 'alert', location: '1号楼305', status: 'offline' },
  { id: 4, name: '输液泵 #008', type: 'pump', location: '医务室', status: 'warning' },
  { id: 5, name: '智能床垫 #002', type: 'bed', location: '3号楼202', status: 'online' },
  { id: 6, name: '血糖仪 #005', type: 'meter', location: '2号楼护士站', status: 'online' },
  { id: 7, name: '心电监护仪 #001', type: 'ecg', location: 'ICU', status: 'online' },
  { id: 8, name: '智能手环 #045', type: 'wearable', location: '活动室', status: 'offline' }
])

// 排班数据
const shifts = ref([
  { id: 1, name: '早班', timeRange: '06:00-14:00', onDuty: 18, total: 20, isCurrent: true },
  { id: 2, name: '中班', timeRange: '14:00-22:00', onDuty: 15, total: 18, isCurrent: false },
  { id: 3, name: '夜班', timeRange: '22:00-06:00', onDuty: 8, total: 12, isCurrent: false }
])

// 过滤设备
const filteredDevices = computed(() => {
  if (deviceFilter.value === 'all') return devices.value
  return devices.value.filter(d => d.status === deviceFilter.value)
})

// 考勤图表配置
const attendanceChartOption = computed(() => {
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
    legend: {
      data: ['正常出勤', '迟到', '早退', '请假'],
      bottom: 0
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '正常出勤',
        type: 'bar',
        stack: 'total',
        data: [45, 48, 46, 47, 45, 40, 38],
        itemStyle: { color: colors[1] }
      },
      {
        name: '迟到',
        type: 'bar',
        stack: 'total',
        data: [2, 1, 3, 1, 2, 1, 0],
        itemStyle: { color: colors[3] }
      },
      {
        name: '早退',
        type: 'bar',
        stack: 'total',
        data: [1, 0, 1, 0, 1, 0, 1],
        itemStyle: { color: colors[2] }
      },
      {
        name: '请假',
        type: 'bar',
        stack: 'total',
        data: [3, 2, 2, 3, 4, 5, 6],
        itemStyle: { color: colors[4] }
      }
    ]
  }
})

// 方法
const formatTime = (timestamp: number) => {
  const now = Date.now()
  const diff = now - timestamp
  const minutes = Math.floor(diff / 60000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

const getProgressColor = (percentage: number) => {
  if (percentage >= 80) return '#67c23a'
  if (percentage >= 60) return '#409eff'
  return '#f56c6c'
}

const getTaskStatusType = (status: string) => {
  const map: Record<string, any> = {
    completed: 'success',
    'in-progress': 'primary',
    pending: 'info',
    overdue: 'danger'
  }
  return map[status] || 'info'
}

const getTaskStatusText = (status: string) => {
  const map: Record<string, string> = {
    completed: '已完成',
    'in-progress': '进行中',
    pending: '待处理',
    overdue: '已逾期'
  }
  return map[status] || status
}

const getHealthTagType = (level: string) => {
  const map: Record<string, any> = {
    high: 'danger',
    medium: 'warning',
    low: 'info',
    normal: 'success'
  }
  return map[level] || 'info'
}

const getDeviceIcon = (type: string) => {
  return Monitor // 简化处理，实际应根据类型返回不同图标
}

const getDeviceStatusType = (status: string) => {
  const map: Record<string, any> = {
    online: 'success',
    offline: 'info',
    warning: 'danger'
  }
  return map[status] || 'info'
}

const getDeviceStatusText = (status: string) => {
  const map: Record<string, string> = {
    online: '在线',
    offline: '离线',
    warning: '告警'
  }
  return map[status] || status
}

const refreshData = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
  }, 1000)
}

const toggleAutoRefresh = (enabled: boolean) => {
  if (enabled) {
    // 启动自动刷新
  } else {
    // 停止自动刷新
  }
}

const handleAlertClick = (alert: any) => {
  console.log('处理警报', alert)
}

const handleViewAllAlerts = () => {
  console.log('查看全部警报')
}

const handleViewAllHealth = () => {
  console.log('查看全部健康数据')
}

const handleHealthDetail = (health: any) => {
  console.log('健康详情', health)
}

const handleDeviceClick = (device: any) => {
  console.log('设备详情', device)
}

const handleExportAttendance = () => {
  console.log('导出考勤报表')
}

// WebSocket 连接
// const ws = useWebSocket('ws://localhost:8080/monitor', {
//   onMessage: (message) => {
//     if (message.type === 'alert') {
//       alerts.value.unshift(message.data)
//     } else if (message.type === 'task_update') {
//       // 更新任务
//     } else if (message.type === 'health_data') {
//       // 更新健康数据
//     }
//   },
//   onOpen: () => {
//     connected.value = true
//   },
//   onClose: () => {
//     connected.value = false
//   }
// })

onMounted(() => {
  // 初始化连接
  // ws.connect()

  // 如果 WebSocket 不可用，使用轮询
  // if (!connected.value) {
  //   const polling = usePolling(async () => {
  //     await refreshData()
  //   }, 5000)
  //   polling.start()
  // }
})

onUnmounted(() => {
  // ws.disconnect()
})
</script>

<style scoped lang="scss">
.realtime-monitor {
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
    gap: 16px;
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

  .title-icon {
    font-size: 28px;
  }

  .status-badge {
    :deep(.el-badge__content) {
      background-color: var(--success-color);
    }
  }

  .header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .alerts-banner {
    margin-bottom: 20px;
  }

  .alerts-list {
    margin-top: 8px;
  }

  .alert-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    margin-bottom: 4px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    cursor: pointer;
    transition: var(--transition-base);

    &:hover {
      background: var(--bg-secondary);
    }
  }

  .alert-icon {
    color: var(--warning-color);
  }

  .alert-message {
    flex: 1;
    font-size: 14px;
  }

  .alert-time {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .stats-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .realtime-stat {
    position: relative;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);
    overflow: hidden;

    &:hover {
      box-shadow: var(--card-shadow-hover);
    }

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;
    }

    .stat-content {
      flex: 1;
    }

    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: var(--text-primary);
    }

    .stat-label {
      font-size: 12px;
      color: var(--text-secondary);
    }

    .stat-indicator {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--border-color);

      &.is-active {
        background: var(--success-color);
        animation: pulse 2s infinite;
      }
    }

    &.stat-primary .stat-icon {
      background: var(--gradient-blue);
      color: #fff;
    }

    &.stat-success .stat-icon {
      background: var(--gradient-green);
      color: #fff;
    }

    &.stat-warning .stat-icon {
      background: var(--gradient-orange);
      color: #fff;
    }

    &.stat-danger .stat-icon {
      background: var(--gradient-red);
      color: #fff;
    }
  }

  @keyframes pulse {
    0%, 100% {
      opacity: 1;
    }
    50% {
      opacity: 0.5;
    }
  }

  .monitor-panels,
  .device-panels,
  .attendance-panels {
    margin-bottom: 20px;
  }

  .monitor-card {
    height: 100%;
    min-height: 400px;

    :deep(.el-card__body) {
      padding: 16px;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .card-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .task-list {
    max-height: 400px;
    overflow-y: auto;
  }

  .task-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    margin-bottom: 8px;
    background: var(--bg-tertiary);
    border-radius: 8px;
    border-left: 3px solid transparent;

    &.task-completed {
      border-left-color: var(--success-color);
    }

    &.task-in-progress {
      border-left-color: var(--primary-color);
    }

    &.task-pending {
      border-left-color: var(--info-color);
    }

    &.task-overdue {
      border-left-color: var(--danger-color);
    }
  }

  .task-title {
    font-weight: 500;
    margin-bottom: 4px;
  }

  .task-meta {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .task-status {
    text-align: right;
  }

  .task-time {
    display: block;
    font-size: 12px;
    color: var(--text-secondary);
    margin-top: 4px;
  }

  .health-list {
    max-height: 400px;
    overflow-y: auto;
  }

  .health-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    margin-bottom: 8px;
    background: var(--bg-tertiary);
    border-radius: 8px;
  }

  .health-info {
    flex: 1;
  }

  .health-name {
    font-weight: 500;
    margin-bottom: 4px;
  }

  .health-metrics {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .health-trend {
    display: flex;
    align-items: center;
    font-size: 12px;
    gap: 2px;

    &.up {
      color: var(--danger-color);
    }

    &.down {
      color: var(--success-color);
    }
  }

  .device-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 12px;
  }

  .device-item {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 12px;
    background: var(--bg-tertiary);
    border-radius: 8px;
    cursor: pointer;
    transition: var(--transition-base);

    &:hover {
      background: var(--bg-secondary);
    }

    .device-icon {
      width: 32px;
      height: 32px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 18px;
      background: var(--primary-color);
      color: #fff;
    }

    .device-name {
      font-weight: 500;
    }

    .device-location {
      font-size: 12px;
      color: var(--text-secondary);
    }

    .device-indicator {
      position: absolute;
      top: 12px;
      right: 12px;
      width: 8px;
      height: 8px;
      border-radius: 50%;

      &.indicator-online {
        background: var(--success-color);
      }

      &.indicator-offline {
        background: var(--text-secondary);
      }

      &.indicator-warning {
        background: var(--danger-color);
        animation: pulse 2s infinite;
      }
    }
  }

  .schedule-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .shift-item {
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: 8px;

    &.is-current {
      background: var(--gradient-blue);
      color: #fff;
    }
  }

  .shift-time {
    font-size: 12px;
    opacity: 0.8;
    margin-bottom: 4px;
  }

  .shift-name {
    font-weight: 600;
    margin-bottom: 4px;
  }

  .shift-count {
    font-size: 12px;
    opacity: 0.8;
  }
}

// 响应式适配
@media (max-width: 768px) {
  .realtime-monitor {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .device-grid {
      grid-template-columns: 1fr;
    }

    .realtime-stat {
      .stat-value {
        font-size: 20px;
      }

      .stat-icon {
        width: 40px;
        height: 40px;
        font-size: 20px;
      }
    }
  }
}
</style>
