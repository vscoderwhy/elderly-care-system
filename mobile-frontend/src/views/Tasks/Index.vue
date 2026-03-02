<template>
  <div class="tasks-page">
    <!-- 日期选择 -->
    <div class="date-header">
      <van-icon name="arrow-left" @click="prevDay" />
      <div class="current-date" @click="showDatePicker = true">
        {{ formatDate(currentDate) }}
      </div>
      <van-icon name="arrow" @click="nextDay" />
    </div>

    <!-- 任务统计 -->
    <div class="stats-card">
      <div class="stat-item">
        <div class="stat-value">{{ taskStats.total }}</div>
        <div class="stat-label">总任务</div>
      </div>
      <div class="stat-item">
        <div class="stat-value completed">{{ taskStats.completed }}</div>
        <div class="stat-label">已完成</div>
      </div>
      <div class="stat-item">
        <div class="stat-value pending">{{ taskStats.pending }}</div>
        <div class="stat-label">待完成</div>
      </div>
    </div>

    <!-- 任务列表 -->
    <div class="task-list">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-list
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
        >
          <div
            class="task-card"
            v-for="task in tasks"
            :key="task.id"
            :class="{ completed: task.status === 'completed' }"
          >
            <div class="task-header">
              <div class="task-type">
                <van-tag :type="getTaskTypeColor(task.care_item?.category)">
                  {{ task.care_item?.category || '护理' }}
                </van-tag>
              </div>
              <div class="task-time">
                {{ formatTime(task.recorded_at) }}
              </div>
            </div>
            <div class="task-content">
              <div class="elderly-name">
                <van-icon name="user-o" />
                {{ task.elderly?.name }}
              </div>
              <div class="task-name">{{ task.care_item?.name }}</div>
              <div class="task-notes" v-if="task.notes">
                备注: {{ task.notes }}
              </div>
            </div>
            <div class="task-footer">
              <van-button
                v-if="task.status !== 'completed'"
                type="primary"
                size="small"
                @click="completeTask(task)"
              >
                完成任务
              </van-button>
              <van-tag v-else type="success">已完成</van-tag>
            </div>
          </div>

          <van-empty v-if="!loading && tasks.length === 0" description="今日暂无任务" />
        </van-list>
      </van-pull-refresh>
    </div>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="selectedDate"
        title="选择日期"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { showToast } from 'vant'
import instance from '@/api'

interface Task {
  id: number
  elderly?: { id: number; name: string }
  care_item?: { id: number; name: string; category: string }
  status: string
  notes: string
  recorded_at: string
}

const loading = ref(false)
const refreshing = ref(false)
const finished = ref(true)
const tasks = ref<Task[]>([])
const showDatePicker = ref(false)
const currentDate = ref(new Date())
const selectedDate = ref([
  new Date().getFullYear().toString(),
  (new Date().getMonth() + 1).toString().padStart(2, '0'),
  new Date().getDate().toString().padStart(2, '0')
])

const taskStats = computed(() => {
  const total = tasks.value.length
  const completed = tasks.value.filter(t => t.status === 'completed').length
  return { total, completed, pending: total - completed }
})

const formatDate = (date: Date) => {
  const today = new Date()
  const isToday = date.toDateString() === today.toDateString()
  if (isToday) return '今天'

  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日`
}

const formatTime = (time: string) => {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

const getTaskTypeColor = (category?: string) => {
  const colors: Record<string, string> = {
    '喂饭': 'primary',
    '翻身': 'success',
    '清洁': 'warning',
    '用药': 'danger',
    '其他': 'default'
  }
  return colors[category || '其他'] || 'default'
}

const loadTasks = async () => {
  loading.value = true
  try {
    const res = await instance.get('/care/my-tasks')
    tasks.value = res || []
  } catch (error) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const onRefresh = () => {
  loadTasks()
}

const prevDay = () => {
  const date = new Date(currentDate.value)
  date.setDate(date.getDate() - 1)
  currentDate.value = date
  loadTasks()
}

const nextDay = () => {
  const date = new Date(currentDate.value)
  date.setDate(date.getDate() + 1)
  currentDate.value = date
  loadTasks()
}

const onDateConfirm = ({ selectedValues }: any) => {
  currentDate.value = new Date(
    parseInt(selectedValues[0]),
    parseInt(selectedValues[1]) - 1,
    parseInt(selectedValues[2])
  )
  showDatePicker.value = false
  loadTasks()
}

const completeTask = async (task: Task) => {
  try {
    // 这里可以调用API更新任务状态
    task.status = 'completed'
    showToast('任务已完成')
  } catch (error) {
    showToast('操作失败')
  }
}

onMounted(() => {
  loadTasks()
})
</script>

<style scoped>
.tasks-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.date-header {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  gap: 20px;
}

.current-date {
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
}

.stats-card {
  display: flex;
  background: #fff;
  margin: 15px;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-value.completed {
  color: #07c160;
}

.stat-value.pending {
  color: #ff976a;
}

.stat-label {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.task-list {
  padding: 0 15px;
}

.task-card {
  background: #fff;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.task-card.completed {
  opacity: 0.7;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.task-time {
  font-size: 12px;
  color: #999;
}

.task-content {
  margin-bottom: 12px;
}

.elderly-name {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 6px;
}

.task-name {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.task-notes {
  font-size: 12px;
  color: #999;
}

.task-footer {
  display: flex;
  justify-content: flex-end;
}
</style>
