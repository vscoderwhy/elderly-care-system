<template>
  <div class="task-scheduler">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <h3>
            <el-icon><Calendar /></el-icon>
            护理任务调度系统
          </h3>
          <div class="header-actions">
            <el-button size="small" @click="refreshSchedule">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button size="small" @click="showConflictPanel = !showConflictPanel">
              <el-icon><Warning /></el-icon>
              冲突检测 ({{ conflicts.length }})
            </el-button>
            <el-button type="primary" size="small" @click="autoAssignTasks">
              <el-icon><MagicStick /></el-icon>
              自动分配
            </el-button>
            <el-button type="success" size="small" @click="showAddTaskDialog">
              <el-icon><Plus /></el-icon>
              新建任务
            </el-button>
          </div>
        </div>
      </template>

      <!-- 视图切换和筛选 -->
      <div class="toolbar">
        <el-radio-group v-model="viewMode">
          <el-radio-button label="day">日视图</el-radio-button>
          <el-radio-button label="week">周视图</el-radio-button>
          <el-radio-button label="month">月视图</el-radio-button>
          <el-radio-button label="timeline">时间轴</el-radio-button>
        </el-radio-group>

        <div class="filter-group">
          <el-date-picker
            v-model="selectedDate"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            @change="loadSchedule"
          />
          <el-select v-model="filterNurse" placeholder="全部护工" clearable style="width: 150px">
            <el-option
              v-for="nurse in nurseList"
              :key="nurse.id"
              :label="nurse.name"
              :value="nurse.id"
            />
          </el-select>
          <el-select v-model="filterDepartment" placeholder="全部部门" clearable style="width: 150px">
            <el-option label="护理部" value="nursing" />
            <el-option label="康复部" value="rehab" />
            <el-option label="医疗部" value="medical" />
          </el-select>
        </div>
      </div>

      <!-- 冲突面板 -->
      <div v-if="showConflictPanel" class="conflict-panel">
        <el-alert
          v-if="conflicts.length === 0"
          title="未发现冲突"
          type="success"
          :closable="false"
        />
        <div v-else class="conflict-list">
          <div
            v-for="conflict in conflicts"
            :key="conflict.id"
            class="conflict-item"
            :class="`conflict-${conflict.level}`"
          >
            <div class="conflict-icon">
              <el-icon><WarningFilled /></el-icon>
            </div>
            <div class="conflict-content">
              <div class="conflict-title">{{ conflict.title }}</div>
              <div class="conflict-detail">{{ conflict.detail }}</div>
              <div class="conflict-time">{{ conflict.time }}</div>
            </div>
            <div class="conflict-actions">
              <el-button size="small" @click="resolveConflict(conflict)">解决</el-button>
              <el-button size="small" type="primary" @click="ignoreConflict(conflict)">
                忽略
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 日视图/周视图 -->
      <div v-if="viewMode === 'day' || viewMode === 'week'" class="schedule-grid">
        <div class="time-header">
          <div class="time-label">时间</div>
          <div
            v-for="hour in workingHours"
            :key="hour"
            class="hour-column"
          >
            {{ hour }}:00
          </div>
        </div>

        <div class="schedule-body">
          <div
            v-for="nurse in filteredNurses"
            :key="nurse.id"
            class="nurse-row"
          >
            <div class="nurse-header">
              <el-avatar :size="40" :src="nurse.avatar" />
              <div class="nurse-info">
                <div class="nurse-name">{{ nurse.name }}</div>
                <div class="nurse-stats">
                  <el-tag size="small" type="info">
                    {{ getNurseTaskCount(nurse.id) }}个任务
                  </el-tag>
                  <el-tag size="small" :type="getWorkloadType(nurse.id)">
                    {{ getWorkloadText(nurse.id) }}
                  </el-tag>
                </div>
              </div>
            </div>

            <div class="nurse-tasks">
              <div
                v-for="hour in workingHours"
                :key="hour"
                class="hour-cell"
                @click="handleCellClick(nurse, hour)"
                @drop="handleDrop($event, nurse, hour)"
                @dragover.prevent
              >
                <div
                  v-for="task in getTasksAtHour(nurse.id, hour)"
                  :key="task.id"
                  class="task-block"
                  :class="`task-${task.priority}`"
                  :style="{ width: getTaskWidth(task, hour) + '%' }"
                  draggable="true"
                  @dragstart="handleDragStart($event, task)"
                  @click.stop="handleTaskClick(task)"
                >
                  <div class="task-time">{{ task.startTime }} - {{ task.endTime }}</div>
                  <div class="task-title">{{ task.title }}</div>
                  <div class="task-elderly">{{ task.elderlyName }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 时间轴视图 -->
      <div v-else-if="viewMode === 'timeline'" class="timeline-view">
        <div class="timeline-axis">
          <div
            v-for="hour in workingHours"
            :key="hour"
            class="timeline-hour"
            :style="{ left: ((hour - 6) / 12 * 100) + '%' }"
          >
            {{ hour }}:00
          </div>
        </div>

        <div class="timeline-content">
          <div
            v-for="nurse in filteredNurses"
            :key="nurse.id"
            class="timeline-nurse"
          >
            <div class="timeline-nurse-header">
              <el-avatar :size="32" :src="nurse.avatar" />
              <span>{{ nurse.name }}</span>
            </div>

            <div class="timeline-track">
              <div
                v-for="task in getNurseTasks(nurse.id)"
                :key="task.id"
                class="timeline-task"
                :class="`task-${task.priority}`"
                :style="getTaskPosition(task)"
                @click="handleTaskClick(task)"
              >
                <div class="timeline-task-title">{{ task.title }}</div>
                <div class="timeline-task-time">
                  {{ task.startTime }} - {{ task.endTime }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 月视图 -->
      <div v-else-if="viewMode === 'month'" class="month-view">
        <el-calendar v-model="calendarDate">
          <template #date-cell="{ data }">
            <div class="calendar-day" @click="handleDayClick(data)">
              <div class="day-number">{{ data.day.split('-')[2] }}</div>
              <div class="day-tasks">
                <div
                  v-for="task in getTasksForDate(data.day)"
                  :key="task.id"
                  class="day-task"
                  :class="`task-${task.priority}`"
                  @click.stop="handleTaskClick(task)"
                >
                  {{ task.title }}
                </div>
              </div>
            </div>
          </template>
        </el-calendar>
      </div>
    </el-card>

    <!-- 新建/编辑任务弹窗 -->
    <el-dialog
      v-model="taskDialogVisible"
      :title="editingTask ? '编辑护理任务' : '新建护理任务'"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form :model="taskForm" :rules="taskFormRules" ref="taskFormRef" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="任务标题" prop="title">
              <el-input v-model="taskForm.title" placeholder="请输入任务标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="任务类型" prop="type">
              <el-select v-model="taskForm.type" placeholder="请选择" style="width: 100%">
                <el-option label="日常护理" value="daily" />
                <el-option label="医疗护理" value="medical" />
                <el-option label="康复训练" value="rehab" />
                <el-option label="健康监测" value="monitor" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="执行日期" prop="date">
              <el-date-picker
                v-model="taskForm.date"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="开始时间" prop="startTime">
              <el-time-picker
                v-model="taskForm.startTime"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="开始"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="结束时间" prop="endTime">
              <el-time-picker
                v-model="taskForm.endTime"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="结束"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="负责护工" prop="nurseId">
              <el-select v-model="taskForm.nurseId" placeholder="请选择" style="width: 100%">
                <el-option
                  v-for="nurse in nurseList"
                  :key="nurse.id"
                  :label="`${nurse.name} (${getNurseTaskCount(nurse.id)}个任务)`"
                  :value="nurse.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="服务对象" prop="elderlyId">
              <el-select v-model="taskForm.elderlyId" placeholder="请选择老人" filterable style="width: 100%">
                <el-option
                  v-for="elderly in elderlyList"
                  :key="elderly.id"
                  :label="`${elderly.name} - ${elderly.bedNumber}`"
                  :value="elderly.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="优先级" prop="priority">
              <el-radio-group v-model="taskForm.priority">
                <el-radio label="low">低</el-radio>
                <el-radio label="medium">中</el-radio>
                <el-radio label="high">高</el-radio>
                <el-radio label="urgent">紧急</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="重复">
              <el-select v-model="taskForm.repeat" placeholder="不重复" style="width: 100%">
                <el-option label="不重复" value="none" />
                <el-option label="每天" value="daily" />
                <el-option label="每周" value="weekly" />
                <el-option label="每月" value="monthly" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="任务描述" prop="description">
          <el-input
            v-model="taskForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入任务详情"
          />
        </el-form-item>

        <el-form-item label="注意事项">
          <el-input
            v-model="taskForm.notes"
            type="textarea"
            :rows="2"
            placeholder="特殊要求、注意事项等"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="taskDialogVisible = false">取消</el-button>
        <el-button @click="taskDialogVisible = false; checkConflicts()">
          检测冲突
        </el-button>
        <el-button type="primary" @click="saveTask" :loading="savingTask">保存</el-button>
      </template>
    </el-dialog>

    <!-- 任务详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" title="任务详情" width="600px">
      <div class="task-detail" v-if="currentTask">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务标题" :span="2">
            {{ currentTask.title }}
          </el-descriptions-item>
          <el-descriptions-item label="任务类型">
            <el-tag size="small">{{ getTypeText(currentTask.type) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentTask.priority)" size="small">
              {{ getPriorityText(currentTask.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="执行日期">
            {{ currentTask.date }}
          </el-descriptions-item>
          <el-descriptions-item label="执行时间">
            {{ currentTask.startTime }} - {{ currentTask.endTime }}
          </el-descriptions-item>
          <el-descriptions-item label="负责护工">
            {{ getNurseName(currentTask.nurseId) }}
          </el-descriptions-item>
          <el-descriptions-item label="服务对象">
            {{ currentTask.elderlyName }}
          </el-descriptions-item>
          <el-descriptions-item label="任务描述" :span="2" v-if="currentTask.description">
            {{ currentTask.description }}
          </el-descriptions-item>
          <el-descriptions-item label="注意事项" :span="2" v-if="currentTask.notes">
            {{ currentTask.notes }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="detail-actions">
          <el-button @click="editCurrentTask">编辑</el-button>
          <el-button type="primary" @click="completeTask(currentTask)">
            标记完成
          </el-button>
          <el-button type="danger" @click="deleteTask(currentTask)">
            删除
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Calendar,
  Refresh,
  Warning,
  MagicStick,
  Plus,
  WarningFilled
} from '@element-plus/icons-vue'
import { getElderlyList } from '@/utils/seedData'

interface Nurse {
  id: number
  name: string
  avatar: string
  department: string
  skills: string[]
}

interface Task {
  id: number
  title: string
  type: string
  date: string
  startTime: string
  endTime: string
  nurseId: number
  nurseName?: string
  elderlyId: number
  elderlyName: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
  description: string
  notes: string
  repeat: string
  status: 'pending' | 'in_progress' | 'completed'
}

interface Conflict {
  id: number
  level: 'warning' | 'error'
  title: string
  detail: string
  time: string
  tasks: number[]
}

const viewMode = ref('day')
const selectedDate = ref(new Date().toISOString().split('T')[0])
const calendarDate = ref(new Date())
const filterNurse = ref<number | null>(null)
const filterDepartment = ref('')
const showConflictPanel = ref(false)

const workingHours = [6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18]

const nurseList = ref<Nurse[]>([
  { id: 1, name: '张护士', avatar: '', department: 'nursing', skills: ['日常护理', '健康监测'] },
  { id: 2, name: '李护士', avatar: '', department: 'nursing', skills: ['日常护理', '医疗护理'] },
  { id: 3, name: '王护士', avatar: '', department: 'rehab', skills: ['康复训练', '按摩'] },
  { id: 4, name: '赵护士', avatar: '', department: 'nursing', skills: ['日常护理'] },
  { id: 5, name: '刘护士', avatar: '', department: 'medical', skills: ['医疗护理', '输液'] },
  { id: 6, name: '陈护士', avatar: '', department: 'nursing', skills: ['日常护理', '健康监测'] }
])

const elderlyList = getElderlyList()

const tasks = ref<Task[]>([])
const conflicts = ref<Conflict[]>([])

const taskDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const editingTask = ref<Task | null>(null)
const currentTask = ref<Task | null>(null)
const savingTask = ref(false)

const taskFormRef = ref<FormInstance>()
const draggedTask = ref<Task | null>(null)

const taskForm = reactive<Task>({
  id: 0,
  title: '',
  type: 'daily',
  date: new Date().toISOString().split('T')[0],
  startTime: '',
  endTime: '',
  nurseId: 0,
  elderlyId: 0,
  elderlyName: '',
  priority: 'medium',
  description: '',
  notes: '',
  repeat: 'none',
  status: 'pending'
})

const taskFormRules: FormRules = {
  title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  date: [{ required: true, message: '请选择执行日期', trigger: 'change' }],
  startTime: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  endTime: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  nurseId: [{ required: true, message: '请选择负责护工', trigger: 'change' }],
  elderlyId: [{ required: true, message: '请选择服务对象', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }]
}

const filteredNurses = computed(() => {
  let result = nurseList.value

  if (filterNurse.value) {
    result = result.filter(n => n.id === filterNurse.value)
  }

  if (filterDepartment.value) {
    result = result.filter(n => n.department === filterDepartment.value)
  }

  return result
})

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    daily: '日常护理',
    medical: '医疗护理',
    rehab: '康复训练',
    monitor: '健康监测',
    other: '其他'
  }
  return map[type] || type
}

const getPriorityType = (priority: string) => {
  const map: Record<string, any> = {
    low: 'info',
    medium: 'primary',
    high: 'warning',
    urgent: 'danger'
  }
  return map[priority] || ''
}

const getPriorityText = (priority: string) => {
  const map: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急'
  }
  return map[priority] || priority
}

const getNurseName = (nurseId: number) => {
  const nurse = nurseList.value.find(n => n.id === nurseId)
  return nurse?.name || ''
}

const getNurseTaskCount = (nurseId: number) => {
  return tasks.value.filter(
    t => t.nurseId === nurseId &&
    t.date === selectedDate.value &&
    t.status !== 'completed'
  ).length
}

const getNurseTasks = (nurseId: number) => {
  return tasks.value.filter(
    t => t.nurseId === nurseId && t.date === selectedDate.value
  )
}

const getWorkloadType = (nurseId: number) => {
  const count = getNurseTaskCount(nurseId)
  if (count >= 8) return 'danger'
  if (count >= 5) return 'warning'
  if (count >= 3) return 'primary'
  return 'success'
}

const getWorkloadText = (nurseId: number) => {
  const count = getNurseTaskCount(nurseId)
  if (count >= 8) return '超负荷'
  if (count >= 5) return '较重'
  if (count >= 3) return '正常'
  return '轻松'
}

const getTasksAtHour = (nurseId: number, hour: number) => {
  return tasks.value.filter(t => {
    if (t.nurseId !== nurseId || t.date !== selectedDate.value) return false
    const startHour = parseInt(t.startTime.split(':')[0])
    const endHour = parseInt(t.endTime.split(':')[0])
    return hour >= startHour && hour < endHour
  })
}

const getTaskWidth = (task: Task, hour: number) => {
  const startHour = parseInt(task.startTime.split(':')[0])
  const endHour = parseInt(task.endTime.split(':')[0])
  if (hour === startHour) {
    const duration = endHour - startHour
    return duration * 100
  }
  return 0
}

const getTaskPosition = (task: Task) => {
  const startHour = parseInt(task.startTime.split(':')[0])
  const endHour = parseInt(task.endTime.split(':')[0])
  const startMinute = parseInt(task.startTime.split(':')[1])
  const endMinute = parseInt(task.endTime.split(':')[1])

  const totalMinutes = (endHour - startHour) * 60 + (endMinute - startMinute)
  const offsetMinutes = (startHour - 6) * 60 + startMinute

  const workingDayMinutes = 12 * 60 // 6:00 - 18:00

  return {
    left: (offsetMinutes / workingDayMinutes * 100) + '%',
    width: (totalMinutes / workingDayMinutes * 100) + '%'
  }
}

const getTasksForDate = (dateStr: string) => {
  return tasks.value.filter(t => t.date === dateStr).slice(0, 3)
}

const loadSchedule = async () => {
  // 生成模拟任务数据
  tasks.value = []

  const taskTypes = ['daily', 'medical', 'rehab', 'monitor']
  const priorities: Array<'low' | 'medium' | 'high' | 'urgent'> = ['low', 'medium', 'high', 'urgent']

  for (let i = 0; i < 30; i++) {
    const nurse = nurseList.value[i % nurseList.value.length]
    const elderly = elderlyList[i % elderlyList.length]
    const startHour = 6 + Math.floor(Math.random() * 11)
    const endHour = startHour + 1

    tasks.value.push({
      id: Date.now() + i,
      title: `${['晨间护理', '测量血压', '康复训练', '健康检查', '用药提醒'][i % 5]}`,
      type: taskTypes[i % taskTypes.length],
      date: selectedDate.value,
      startTime: `${String(startHour).padStart(2, '0')}:00`,
      endTime: `${String(endHour).padStart(2, '0')}:00`,
      nurseId: nurse.id,
      elderlyId: elderly.id,
      elderlyName: elderly.name,
      priority: priorities[i % priorities.length],
      description: '任务详情',
      notes: '',
      repeat: 'none',
      status: 'pending'
    })
  }

  checkConflicts()
}

const checkConflicts = () => {
  conflicts.value = []

  // 检查同一护工的时间冲突
  const nurseTasksMap = new Map<number, Task[]>()

  tasks.value.forEach(task => {
    if (!nurseTasksMap.has(task.nurseId)) {
      nurseTasksMap.set(task.nurseId, [])
    }
    nurseTasksMap.get(task.nurseId)!.push(task)
  })

  nurseTasksMap.forEach((nurseTasks, nurseId) => {
    for (let i = 0; i < nurseTasks.length; i++) {
      for (let j = i + 1; j < nurseTasks.length; j++) {
        const t1 = nurseTasks[i]
        const t2 = nurseTasks[j]

        if (t1.date === t2.date) {
          const t1Start = parseInt(t1.startTime.replace(':', ''))
          const t1End = parseInt(t1.endTime.replace(':', ''))
          const t2Start = parseInt(t2.startTime.replace(':', ''))
          const t2End = parseInt(t2.endTime.replace(':', ''))

          if (t1Start < t2End && t2Start < t1End) {
            conflicts.value.push({
              id: Date.now() + conflicts.value.length,
              level: 'error',
              title: `${getNurseName(nurseId)} 时间冲突`,
              detail: `任务"${t1.title}"与"${t2.title}"时间重叠`,
              time: `${t1.date} ${t1.startTime}`,
              tasks: [t1.id, t2.id]
            })
          }
        }
      }
    }
  })

  // 检查超负荷
  nurseTasksMap.forEach((nurseTasks, nurseId) => {
    const dayTasks = nurseTasks.filter(t => t.date === selectedDate.value)
    if (dayTasks.length >= 8) {
      conflicts.value.push({
        id: Date.now() + conflicts.value.length,
        level: 'warning',
        title: `${getNurseName(nurseId)} 超负荷`,
        detail: `当日安排${dayTasks.length}个任务，超过合理负荷`,
        time: selectedDate.value,
        tasks: dayTasks.map(t => t.id)
      })
    }
  })

  if (conflicts.value.length > 0 && !showConflictPanel.value) {
    showConflictPanel.value = true
  }
}

const resolveConflict = (conflict: Conflict) => {
  ElMessage.info('请手动调整相关任务来解决冲突')
}

const ignoreConflict = (conflict: Conflict) => {
  const index = conflicts.value.findIndex(c => c.id === conflict.id)
  if (index !== -1) {
    conflicts.value.splice(index, 1)
  }
}

const showAddTaskDialog = () => {
  editingTask.value = null
  Object.assign(taskForm, {
    id: 0,
    title: '',
    type: 'daily',
    date: selectedDate.value,
    startTime: '',
    endTime: '',
    nurseId: 0,
    elderlyId: 0,
    elderlyName: '',
    priority: 'medium',
    description: '',
    notes: '',
    repeat: 'none',
    status: 'pending'
  })
  taskDialogVisible.value = true
}

const saveTask = async () => {
  if (!taskFormRef.value) return

  try {
    await taskFormRef.value.validate()

    savingTask.value = true

    const elderly = elderlyList.find(e => e.id === taskForm.elderlyId)

    if (editingTask.value) {
      // 更新任务
      const index = tasks.value.findIndex(t => t.id === editingTask.value!.id)
      if (index !== -1) {
        tasks.value[index] = {
          ...taskForm,
          elderlyName: elderly?.name || ''
        }
      }
      ElMessage.success('任务更新成功')
    } else {
      // 新建任务
      tasks.value.push({
        ...taskForm,
        id: Date.now(),
        elderlyName: elderly?.name || ''
      })
      ElMessage.success('任务创建成功')
    }

    taskDialogVisible.value = false
    checkConflicts()
  } finally {
    savingTask.value = false
  }
}

const handleTaskClick = (task: Task) => {
  currentTask.value = task
  detailDialogVisible.value = true
}

const handleCellClick = (nurse: Nurse, hour: number) => {
  taskForm.date = selectedDate.value
  taskForm.startTime = `${String(hour).padStart(2, '0')}:00`
  taskForm.endTime = `${String(hour + 1).padStart(2, '0')}:00`
  taskForm.nurseId = nurse.id
  showAddTaskDialog()
}

const handleDayClick = (data: any) => {
  selectedDate.value = data.day
  viewMode.value = 'day'
  loadSchedule()
}

const handleDragStart = (event: DragEvent, task: Task) => {
  draggedTask.value = task
}

const handleDrop = (event: DragEvent, nurse: Nurse, hour: number) => {
  if (!draggedTask.value) return

  const task = tasks.value.find(t => t.id === draggedTask.value!.id)
  if (task) {
    task.nurseId = nurse.id
    task.startTime = `${String(hour).padStart(2, '0')}:00`
    const duration = parseInt(task.endTime.split(':')[0]) - parseInt(task.startTime.split(':')[0])
    task.endTime = `${String(hour + duration).padStart(2, '0')}:00`
    ElMessage.success('任务已移动')
    checkConflicts()
  }

  draggedTask.value = null
}

const editCurrentTask = () => {
  if (!currentTask.value) return

  editingTask.value = currentTask.value
  Object.assign(taskForm, currentTask.value)
  detailDialogVisible.value = false
  taskDialogVisible.value = true
}

const completeTask = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确认标记此任务为已完成？', '完成任务', {
      type: 'success'
    })

    task.status = 'completed'
    detailDialogVisible.value = false
    ElMessage.success('任务已完成')
  } catch {
    // 取消
  }
}

const deleteTask = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确认删除此任务？', '删除任务', {
      type: 'warning'
    })

    const index = tasks.value.findIndex(t => t.id === task.id)
    if (index !== -1) {
      tasks.value.splice(index, 1)
    }

    detailDialogVisible.value = false
    checkConflicts()
    ElMessage.success('任务已删除')
  } catch {
    // 取消
  }
}

const autoAssignTasks = () => {
  ElMessage.info('自动分配功能正在开发中')
  // TODO: 实现智能分配算法
  // 考虑因素：护工技能匹配、工作负载均衡、时间冲突避免
}

const refreshSchedule = () => {
  loadSchedule()
  ElMessage.success('已刷新')
}

onMounted(() => {
  loadSchedule()
})
</script>

<style scoped lang="scss">
.task-scheduler {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: 8px;

    .filter-group {
      display: flex;
      gap: 12px;
    }
  }

  .conflict-panel {
    margin-bottom: 20px;
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: 8px;

    .conflict-list {
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    .conflict-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      background: var(--card-bg);
      border-radius: 6px;
      border-left: 4px solid;

      &.conflict-error {
        border-left-color: var(--danger-color);
      }

      &.conflict-warning {
        border-left-color: var(--warning-color);
      }

      .conflict-icon {
        font-size: 24px;
        color: var(--warning-color);
      }

      .conflict-content {
        flex: 1;

        .conflict-title {
          font-weight: 600;
          margin-bottom: 4px;
        }

        .conflict-detail {
          font-size: 13px;
          color: var(--text-secondary);
          margin-bottom: 4px;
        }

        .conflict-time {
          font-size: 12px;
          color: var(--text-secondary);
        }
      }
    }
  }

  .schedule-grid {
    .time-header {
      display: flex;
      border-bottom: 1px solid var(--border-color);
      background: var(--bg-tertiary);

      .time-label {
        width: 150px;
        padding: 12px;
        font-weight: 600;
        border-right: 1px solid var(--border-color);
      }

      .hour-column {
        flex: 1;
        padding: 12px;
        text-align: center;
        border-right: 1px solid var(--border-color-lighter);
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .schedule-body {
      .nurse-row {
        display: flex;
        border-bottom: 1px solid var(--border-color-lighter);

        &:hover {
          background: var(--bg-tertiary);
        }

        .nurse-header {
          width: 150px;
          padding: 12px;
          border-right: 1px solid var(--border-color);
          display: flex;
          align-items: center;
          gap: 12px;

          .nurse-info {
            .nurse-name {
              font-weight: 600;
              margin-bottom: 6px;
            }

            .nurse-stats {
              display: flex;
              gap: 4px;
            }
          }
        }

        .nurse-tasks {
          flex: 1;
          display: flex;

          .hour-cell {
            flex: 1;
            min-height: 60px;
            border-right: 1px solid var(--border-color-lighter);
            padding: 4px;
            cursor: pointer;
            position: relative;

            &:hover {
              background: var(--bg-tertiary);
            }

            .task-block {
              position: absolute;
              top: 4px;
              left: 4px;
              right: 4px;
              padding: 6px 8px;
              border-radius: 4px;
              font-size: 12px;
              cursor: pointer;
              z-index: 1;

              &.task-low {
                background: #e1f3d8;
                border-left: 3px solid #67c23a;
              }

              &.task-medium {
                background: #d9ecff;
                border-left: 3px solid #409eff;
              }

              &.task-high {
                background: #fdf6ec;
                border-left: 3px solid #e6a23c;
              }

              &.task-urgent {
                background: #fef0f0;
                border-left: 3px solid #f56c6c;
              }

              .task-time {
                font-size: 11px;
                color: var(--text-secondary);
                margin-bottom: 2px;
              }

              .task-title {
                font-weight: 600;
                margin-bottom: 2px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
              }

              .task-elderly {
                font-size: 11px;
                color: var(--text-secondary);
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
              }
            }
          }
        }
      }
    }
  }

  .timeline-view {
    .timeline-axis {
      position: relative;
      height: 30px;
      margin-bottom: 16px;
      border-bottom: 1px solid var(--border-color);

      .timeline-hour {
        position: absolute;
        transform: translateX(-50%);
        font-size: 11px;
        color: var(--text-secondary);
      }
    }

    .timeline-content {
      .timeline-nurse {
        display: flex;
        align-items: center;
        margin-bottom: 12px;

        .timeline-nurse-header {
          width: 120px;
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 14px;
          font-weight: 600;
        }

        .timeline-track {
          flex: 1;
          position: relative;
          height: 50px;
          background: var(--bg-tertiary);
          border-radius: 4px;

          .timeline-task {
            position: absolute;
            top: 8px;
            height: 34px;
            padding: 6px 10px;
            border-radius: 4px;
            cursor: pointer;
            overflow: hidden;

            &.task-low {
              background: #e1f3d8;
              border-left: 3px solid #67c23a;
            }

            &.task-medium {
              background: #d9ecff;
              border-left: 3px solid #409eff;
            }

            &.task-high {
              background: #fdf6ec;
              border-left: 3px solid #e6a23c;
            }

            &.task-urgent {
              background: #fef0f0;
              border-left: 3px solid #f56c6c;
            }

            .timeline-task-title {
              font-size: 12px;
              font-weight: 600;
              margin-bottom: 2px;
            }

            .timeline-task-time {
              font-size: 11px;
              color: var(--text-secondary);
            }
          }
        }
      }
    }
  }

  .month-view {
    :deep(.el-calendar-table) {
      .el-calendar-day {
        height: 100px;
        padding: 4px;
      }

      .calendar-day {
        height: 100%;
        cursor: pointer;

        &:hover {
          background: var(--bg-tertiary);
        }

        .day-number {
          font-weight: 600;
          margin-bottom: 6px;
        }

        .day-tasks {
          display: flex;
          flex-direction: column;
          gap: 2px;

          .day-task {
            padding: 2px 6px;
            font-size: 11px;
            border-radius: 2px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;

            &.task-low {
              background: #e1f3d8;
              color: #67c23a;
            }

            &.task-medium {
              background: #d9ecff;
              color: #409eff;
            }

            &.task-high {
              background: #fdf6ec;
              color: #e6a23c;
            }

            &.task-urgent {
              background: #fef0f0;
              color: #f56c6c;
            }
          }
        }
      }
    }
  }

  .detail-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px solid var(--border-color-lighter);
  }
}
</style>
