<template>
  <div class="schedule-management">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Calendar /></el-icon>
          排班管理
        </h2>
      </div>
      <div class="header-actions">
        <el-date-picker v-model="currentWeek" type="week" format="YYYY 第 ww 周" placeholder="选择周" @change="loadSchedule" />
        <el-button type="primary" @click="handleAutoSchedule">
          <el-icon><MagicStick /></el-icon>
          自动排班
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
    </div>

    <!-- 班次说明 -->
    <el-card shadow="never" class="legend-card">
      <div class="shift-legend">
        <div class="legend-item" v-for="shift in shiftTypes" :key="shift.key">
          <span class="legend-color" :style="{ background: shift.color }"></span>
          <span class="legend-label">{{ shift.label }} ({{ shift.time }})</span>
        </div>
      </div>
    </el-card>

    <!-- 排班日历 -->
    <el-card shadow="never" class="calendar-card">
      <div class="calendar-header">
        <div class="week-info">
          <el-button circle @click="prevWeek">
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
          <span class="week-text">{{ weekRangeText }}</span>
          <el-button circle @click="nextWeek">
            <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
        <div class="staff-filter">
          <el-select v-model="filterDepartment" placeholder="全部部门" clearable style="width: 150px" @change="loadSchedule">
            <el-option label="全部部门" value="" />
            <el-option label="护理部" value="nursing" />
            <el-option label="医务室" value="medical" />
            <el-option label="康复科" value="rehab" />
            <el-option label="膳食部" value="dining" />
          </el-select>
        </div>
      </div>

      <div class="schedule-table">
        <div class="table-header">
          <div class="header-cell name-col">员工</div>
          <div class="header-cell day-col" v-for="day in weekDays" :key="day.value">
            <div class="day-name">{{ day.label }}</div>
            <div class="day-date">{{ day.date }}</div>
          </div>
        </div>

        <div class="table-body">
          <div class="table-row" v-for="staff in scheduleList" :key="staff.id">
            <div class="body-cell name-col">
              <div class="staff-info">
                <el-avatar :size="32" :src="staff.avatar">{{ staff.name.charAt(0) }}</el-avatar>
                <div class="staff-detail">
                  <div class="staff-name">{{ staff.name }}</div>
                  <div class="staff-position">{{ staff.position }}</div>
                </div>
              </div>
            </div>
            <div
              class="body-cell day-col"
              v-for="day in weekDays"
              :key="day.value"
              @click="handleCellClick(staff, day)"
            >
              <div
                v-if="getShift(staff, day)"
                class="shift-badge"
                :style="{ background: getShiftColor(staff, day) }"
              >
                {{ getShiftText(staff, day) }}
              </div>
              <div v-else class="empty-cell">+</div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 排班编辑对话框 -->
    <el-dialog v-model="dialogVisible" title="编辑排班" width="500px">
      <el-form :model="scheduleForm" label-width="100px">
        <el-form-item label="员工">
          <el-input :value="currentStaff?.name" disabled />
        </el-form-item>
        <el-form-item label="日期">
          <el-input :value="currentDay?.label" disabled />
        </el-form-item>
        <el-form-item label="班次" prop="shift">
          <el-radio-group v-model="scheduleForm.shift">
            <el-radio label="morning">早班 (06:00-14:00)</el-radio>
            <el-radio label="day">白班 (08:00-17:00)</el-radio>
            <el-radio label="evening">晚班 (14:00-22:00)</el-radio>
            <el-radio label="night">夜班 (22:00-06:00)</el-radio>
            <el-radio label="rest">休息</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="scheduleForm.remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleClearShift">清除排班</el-button>
        <el-button type="primary" @click="handleSaveShift">保存</el-button>
      </template>
    </el-dialog>

    <!-- 自动排班对话框 -->
    <el-dialog v-model="autoDialogVisible" title="自动排班配置" width="600px">
      <el-form :model="autoForm" label-width="120px">
        <el-form-item label="排班周期">
          <el-radio-group v-model="autoForm.period">
            <el-radio label="week">按周</el-radio>
            <el-radio label="month">按月</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="早班人数">
          <el-input-number v-model="autoForm.morningCount" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="白班人数">
          <el-input-number v-model="autoForm.dayCount" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="晚班人数">
          <el-input-number v-model="autoForm.eveningCount" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="夜班人数">
          <el-input-number v-model="autoForm.nightCount" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="公平原则">
          <el-checkbox-group v-model="autoForm.fairness">
            <el-checkbox label="balance">班次均衡</el-checkbox>
            <el-checkbox label="weekend">周末轮休</el-checkbox>
            <el-checkbox label="night">夜班轮换</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="autoDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRunAutoSchedule" :loading="autoScheduling">
          开始排班
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Calendar,
  MagicStick,
  Download,
  ArrowLeft,
  ArrowRight
} from '@element-plus/icons-vue'

const shiftTypes = ref([
  { key: 'morning', label: '早班', time: '06:00-14:00', color: '#67c23a' },
  { key: 'day', label: '白班', time: '08:00-17:00', color: '#409eff' },
  { key: 'evening', label: '晚班', time: '14:00-22:00', color: '#e6a23c' },
  { key: 'night', label: '夜班', time: '22:00-06:00', color: '#909399' },
  { key: 'rest', label: '休息', time: '', color: '#f4f4f5' }
])

const currentWeek = ref(new Date())
const filterDepartment = ref('')
const dialogVisible = ref(false)
const autoDialogVisible = ref(false)
const autoScheduling = ref(false)

const currentStaff = ref<any>(null)
const currentDay = ref<any>(null)

const scheduleForm = reactive({
  shift: 'day',
  remark: ''
})

const autoForm = reactive({
  period: 'week',
  morningCount: 3,
  dayCount: 5,
  eveningCount: 4,
  nightCount: 2,
  fairness: ['balance', 'weekend', 'night']
})

const scheduleList = ref([
  {
    id: 1,
    name: '张护士',
    position: '护士',
    department: 'nursing',
    avatar: '',
    schedules: {
      '2026-03-03': { shift: 'morning' },
      '2026-03-04': { shift: 'day' },
      '2026-03-05': { shift: 'evening' },
      '2026-03-06': { shift: 'rest' },
      '2026-03-07': { shift: 'day' },
      '2026-03-08': { shift: 'day' },
      '2026-03-09': { shift: 'rest' }
    }
  },
  {
    id: 2,
    name: '李护士',
    position: '护士',
    department: 'nursing',
    avatar: '',
    schedules: {
      '2026-03-03': { shift: 'day' },
      '2026-03-04': { shift: 'evening' },
      '2026-03-05': { shift: 'night' },
      '2026-03-06': { shift: 'night' },
      '2026-03-07': { shift: 'rest' },
      '2026-03-08': { shift: 'day' },
      '2026-03-09': { shift: 'evening' }
    }
  },
  {
    id: 3,
    name: '王护理员',
    position: '护理员',
    department: 'nursing',
    avatar: '',
    schedules: {
      '2026-03-03': { shift: 'evening' },
      '2026-03-04': { shift: 'rest' },
      '2026-03-05': { shift: 'day' },
      '2026-03-06': { shift: 'day' },
      '2026-03-07': { shift: 'morning' },
      '2026-03-08': { shift: 'rest' },
      '2026-03-09': { shift: 'day' }
    }
  },
  {
    id: 4,
    name: '赵康复师',
    position: '康复师',
    department: 'rehab',
    avatar: '',
    schedules: {
      '2026-03-03': { shift: 'day' },
      '2026-03-04': { shift: 'day' },
      '2026-03-05': { shift: 'rest' },
      '2026-03-06': { shift: 'day' },
      '2026-03-07': { shift: 'day' },
      '2026-03-08': { shift: 'rest' },
      '2026-03-09': { shift: 'day' }
    }
  }
])

const weekDays = computed(() => {
  const days = []
  const startOfWeek = getStartOfWeek(currentWeek.value)

  for (let i = 0; i < 7; i++) {
    const date = new Date(startOfWeek)
    date.setDate(startOfWeek.getDate() + i)

    const dayNames = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    days.push({
      value: date.toISOString().split('T')[0],
      label: dayNames[date.getDay()],
      date: `${date.getMonth() + 1}/${date.getDate()}`
    })
  }

  return days
})

const weekRangeText = computed(() => {
  const start = getStartOfWeek(currentWeek.value)
  const end = new Date(start)
  end.setDate(start.getDate() + 6)

  const formatDate = (d: Date) => `${d.getMonth() + 1}月${d.getDate()}日`
  return `${formatDate(start)} - ${formatDate(end)}`
})

const getStartOfWeek = (date: Date) => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day
  return new Date(d.setDate(diff))
}

const getShift = (staff: any, day: any) => {
  return staff.schedules[day.value]?.shift
}

const getShiftText = (staff: any, day: any) => {
  const shift = getShift(staff, day)
  const shiftType = shiftTypes.value.find(s => s.key === shift)
  return shiftType?.label || ''
}

const getShiftColor = (staff: any, day: any) => {
  const shift = getShift(staff, day)
  const shiftType = shiftTypes.value.find(s => s.key === shift)
  return shiftType?.color || '#f4f4f5'
}

const prevWeek = () => {
  const newDate = new Date(currentWeek.value)
  newDate.setDate(newDate.getDate() - 7)
  currentWeek.value = newDate
  loadSchedule()
}

const nextWeek = () => {
  const newDate = new Date(currentWeek.value)
  newDate.setDate(newDate.getDate() + 7)
  currentWeek.value = newDate
  loadSchedule()
}

const loadSchedule = () => {
  // TODO: 加载排班数据
}

const handleCellClick = (staff: any, day: any) => {
  currentStaff.value = staff
  currentDay.value = day
  const existingShift = staff.schedules[day.value]?.shift
  scheduleForm.shift = existingShift || 'day'
  dialogVisible.value = true
}

const handleSaveShift = () => {
  if (!currentStaff.value || !currentDay.value) return

  const staff = scheduleList.value.find(s => s.id === currentStaff.value.id)
  if (staff) {
    staff.schedules[currentDay.value.value] = {
      shift: scheduleForm.shift,
      remark: scheduleForm.remark
    }
  }

  dialogVisible.value = false
  ElMessage.success('排班保存成功')
}

const handleClearShift = () => {
  if (!currentStaff.value || !currentDay.value) return

  const staff = scheduleList.value.find(s => s.id === currentStaff.value.id)
  if (staff && staff.schedules[currentDay.value.value]) {
    delete staff.schedules[currentDay.value.value]
  }

  dialogVisible.value = false
  ElMessage.success('已清除排班')
}

const handleAutoSchedule = () => {
  autoDialogVisible.value = true
}

const handleRunAutoSchedule = () => {
  autoScheduling.value = true
  setTimeout(() => {
    autoScheduling.value = false
    autoDialogVisible.value = false
    ElMessage.success('自动排班完成')
    loadSchedule()
  }, 2000)
}

const handleExport = () => {
  ElMessage.success('导出成功')
}

onMounted(() => {
  loadSchedule()
})
</script>

<style scoped lang="scss">
.schedule-management {
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

  .header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .legend-card {
    margin-bottom: 20px;
  }

  .shift-legend {
    display: flex;
    gap: 24px;
    flex-wrap: wrap;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .legend-color {
    width: 16px;
    height: 16px;
    border-radius: 4px;
  }

  .legend-label {
    font-size: 14px;
    color: var(--text-primary);
  }

  .calendar-card {
    .calendar-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
    }

    .week-info {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .week-text {
      font-size: 16px;
      font-weight: 500;
      color: var(--text-primary);
      min-width: 200px;
      text-align: center;
    }
  }

  .schedule-table {
    border: 1px solid var(--border-color-lighter);
    border-radius: 4px;
    overflow: hidden;
  }

  .table-header {
    display: flex;
    background: var(--bg-secondary);
  }

  .header-cell {
    padding: 12px 8px;
    text-align: center;
    border-right: 1px solid var(--border-color-lighter);
    font-weight: 500;
    color: var(--text-primary);

    &:last-child {
      border-right: none;
    }

    &.name-col {
      width: 180px;
      min-width: 180px;
    }

    &.day-col {
      flex: 1;
      min-width: 100px;
    }
  }

  .day-name {
    font-size: 14px;
    margin-bottom: 4px;
  }

  .day-date {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .table-body {
    .table-row {
      display: flex;
      border-top: 1px solid var(--border-color-lighter);

      &:hover {
        background: var(--bg-secondary);
      }
    }
  }

  .body-cell {
    padding: 8px;
    text-align: center;
    border-right: 1px solid var(--border-color-lighter);
    min-height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;

    &:last-child {
      border-right: none;
    }

    &:hover {
      background: var(--fill-color-light);
    }

    &.name-col {
      justify-content: flex-start;
    }
  }

  .staff-info {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
  }

  .staff-detail {
    text-align: left;
  }

  .staff-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .staff-position {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .shift-badge {
    padding: 6px 12px;
    border-radius: 4px;
    font-size: 12px;
    color: #fff;
    font-weight: 500;
    white-space: nowrap;
  }

  .empty-cell {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--fill-color-light);
    color: var(--text-secondary);
    font-size: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

@media (max-width: 1024px) {
  .schedule-table {
    overflow-x: auto;

    .table-header,
    .table-row {
      min-width: 900px;
    }
  }
}
</style>
