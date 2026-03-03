<template>
  <div class="attendance-page">
    <van-nav-bar title="考勤打卡" left-arrow @click-left="onClickLeft" />

    <!-- 当前时间 -->
    <div class="current-time">
      <div class="time">{{ currentTime }}</div>
      <div class="date">{{ currentDate }}</div>
    </div>

    <!-- 打卡区域 -->
    <div class="clock-area">
      <div class="clock-in" :class="{ disabled: clockedIn }" @click="handleClockIn">
        <div class="clock-icon">📍</div>
        <div class="clock-label">上班打卡</div>
        <div class="clock-time" v-if="clockInTime">{{ clockInTime }}</div>
        <div class="clock-status" :class="clockInStatus">{{ clockInStatusText }}</div>
      </div>

      <div class="clock-out" :class="{ disabled: !clockedIn || clockedOut }" @click="handleClockOut">
        <div class="clock-icon">🏠</div>
        <div class="clock-label">下班打卡</div>
        <div class="clock-time" v-if="clockOutTime">{{ clockOutTime }}</div>
        <div class="clock-status" :class="clockOutStatus">{{ clockOutStatusText }}</div>
      </div>
    </div>

    <!-- 位置信息 -->
    <div class="location-info">
      <van-icon name="location-o" />
      <span>{{ location || '正在获取位置...' }}</span>
    </div>

    <!-- 本月统计 -->
    <div class="stats-card">
      <div class="stat-item">
        <div class="stat-value">{{ stats.presentDays }}</div>
        <div class="stat-label">出勤天数</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ stats.lateDays }}</div>
        <div class="stat-label">迟到次数</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ stats.attendanceRate }}%</div>
        <div class="stat-label">出勤率</div>
      </div>
    </div>

    <!-- 考勤记录 -->
    <div class="records-section">
      <van-divider>本月考勤记录</van-divider>
      <van-list v-model:loading="loading" :finished="finished" @load="onLoad">
        <div class="record-item" v-for="record in records" :key="record.id">
          <div class="record-date">{{ record.date }}</div>
          <div class="record-detail">
            <div class="record-times">
              <span class="time-in">上班: {{ record.clockIn || '--' }}</span>
              <span class="time-out">下班: {{ record.clockOut || '--' }}</span>
            </div>
            <div class="record-status">
              <van-tag v-if="record.status === 'normal'" type="success">正常</van-tag>
              <van-tag v-else-if="record.status === 'late'" type="warning">迟到</van-tag>
              <van-tag v-else-if="record.status === 'early'" type="danger">早退</van-tag>
            </div>
          </div>
        </div>
      </van-list>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '@/api/index'
import { showToast, showLoadingToast, closeToast } from 'vant'

const router = useRouter()

const currentTime = ref('')
const currentDate = ref('')
const location = ref('')
const clockedIn = ref(false)
const clockedOut = ref(false)
const clockInTime = ref('')
const clockOutTime = ref('')
const clockInStatus = ref('')
const clockOutStatus = ref('')
const clockInStatusText = ref('')
const clockOutStatusText = ref('')

const stats = ref({
  presentDays: 0,
  lateDays: 0,
  attendanceRate: 100
})

const records = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
let timer: any

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  currentDate.value = now.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })
}

// 获取位置
const getLocation = () => {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords
        location.value = `${latitude.toFixed(4)}, ${longitude.toFixed(4)}`
      },
      () => {
        location.value = '无法获取位置'
      }
    )
  }
}

// 上班打卡
const handleClockIn = async () => {
  if (clockedIn.value) return

  showLoadingToast({ message: '打卡中...', forbidClick: true })

  try {
    const res = await axios.post('/attendance/clock-in', {
      latitude: 0,
      longitude: 0,
      location: location.value,
      device: navigator.userAgent
    })

    closeToast()
    showToast({ type: 'success', message: '上班打卡成功' })

    clockInTime.value = res.data.time.substring(11, 16)
    clockedIn.value = true
    clockInStatus.value = res.data.status || 'normal'
    clockInStatusText.value = getStatusText(res.data.status)

    loadTodayAttendance()
  } catch (error: any) {
    closeToast()
    showToast({ type: 'fail', message: error.response?.data?.message || '打卡失败' })
  }
}

// 下班打卡
const handleClockOut = async () => {
  if (!clockedIn.value || clockedOut.value) return

  showLoadingToast({ message: '打卡中...', forbidClick: true })

  try {
    const res = await axios.post('/attendance/clock-out', {
      latitude: 0,
      longitude: 0,
      location: location.value,
      device: navigator.userAgent
    })

    closeToast()
    showToast({ type: 'success', message: '下班打卡成功' })

    clockOutTime.value = res.data.time.substring(11, 16)
    clockedOut.value = true
    clockOutStatus.value = res.data.status || 'normal'
    clockOutStatusText.value = getStatusText(res.data.status)
  } catch (error: any) {
    closeToast()
    showToast({ type: 'fail', message: error.response?.data?.message || '打卡失败' })
  }
}

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'normal': '正常',
    'late': '迟到',
    'early': '早退'
  }
  return statusMap[status] || '正常'
}

// 加载今日打卡状态
const loadTodayAttendance = async () => {
  try {
    const res = await axios.get('/attendance/today')
    if (res.data.clock_in) {
      clockedIn.value = true
      clockInTime.value = res.data.clock_in.time.substring(11, 16)
      clockInStatus.value = res.data.clock_in.status
      clockInStatusText.value = getStatusText(res.data.clock_in.status)
    }
    if (res.data.clock_out) {
      clockedOut.value = true
      clockOutTime.value = res.data.clock_out.time.substring(11, 16)
      clockOutStatus.value = res.data.clock_out.status
      clockOutStatusText.value = getStatusText(res.data.clock_out.status)
    }
  } catch (error) {
    // 忽略错误
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    const now = new Date()
    const res = await axios.get('/attendance/stats', {
      params: {
        year: now.getFullYear(),
        month: now.getMonth() + 1
      }
    })
    stats.value = {
      presentDays: res.data.present_days || 0,
      lateDays: res.data.late_days || 0,
      attendanceRate: Math.round(res.data.attendance_rate || 100)
    }
  } catch (error) {
    // 忽略错误
  }
}

// 加载考勤记录
const onLoad = async () => {
  loading.value = true
  // TODO: 加载考勤记录
  loading.value = false
  finished.value = true
}

const onClickLeft = () => {
  router.back()
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
  getLocation()
  loadTodayAttendance()
  loadStats()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.attendance-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 20px;
}

.current-time {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 30px 20px;
  text-align: center;
}

.time {
  font-size: 48px;
  font-weight: bold;
  font-family: 'SF Mono', 'Monaco', monospace;
}

.date {
  font-size: 16px;
  opacity: 0.9;
  margin-top: 5px;
}

.clock-area {
  display: flex;
  justify-content: space-around;
  padding: 30px 20px;
  background: white;
  margin: -20px 20px 20px;
  border-radius: 20px 20px 0 0;
}

.clock-in, .clock-out {
  flex: 1;
  text-align: center;
  padding: 20px;
  margin: 0 5px;
  border-radius: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.clock-in {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.clock-in.disabled {
  background: #ccc;
  cursor: not-allowed;
}

.clock-out {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.clock-out.disabled {
  background: #ccc;
  cursor: not-allowed;
}

.clock-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.clock-label {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.clock-time {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 5px;
}

.clock-status.normal {
  color: #07c160;
}

.clock-status.late,
.clock-status.early {
  color: #ff976a;
}

.location-info {
  background: white;
  padding: 15px 20px;
  margin: 0 20px 20px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #666;
}

.stats-card {
  background: white;
  padding: 20px;
  margin: 0 20px 20px;
  border-radius: 15px;
  display: flex;
  justify-content: space-around;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #667eea;
}

.stat-label {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.records-section {
  background: white;
  margin: 0 20px;
  border-radius: 15px;
  padding: 10px 20px 20px;
}

.record-item {
  padding: 15px 0;
  border-bottom: 1px solid #eee;
}

.record-date {
  font-weight: bold;
  margin-bottom: 8px;
}

.record-detail {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.record-times {
  font-size: 14px;
  color: #666;
}

.record-times span {
  margin-right: 15px;
}
</style>
