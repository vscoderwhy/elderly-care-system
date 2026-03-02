<template>
  <div class="schedule-page">
    <!-- 月度统计 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ monthlyStats.total_shifts || 0 }}</div>
          <div class="stat-label">本月排班</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card morning">
          <div class="stat-value">{{ monthlyStats.morning_shifts || 0 }}</div>
          <div class="stat-label">早班</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card afternoon">
          <div class="stat-value">{{ monthlyStats.afternoon_shifts || 0 }}</div>
          <div class="stat-label">中班</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card night">
          <div class="stat-value">{{ monthlyStats.night_shifts || 0 }}</div>
          <div class="stat-label">晚班</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 日期选择和操作栏 -->
    <el-card class="toolbar-card">
      <el-row :gutter="20" align="middle">
        <el-col :span="12">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="loadSchedules"
          />
        </el-col>
        <el-col :span="12" style="text-align: right;">
          <el-button type="primary" @click="showAddDialog">添加排班</el-button>
          <el-button @click="batchAddDialog = true">批量排班</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 排班日历视图 -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>排班日历</span>
          <el-radio-group v-model="viewMode" size="small">
            <el-radio-button label="calendar">日历视图</el-radio-button>
            <el-radio-button label="list">列表视图</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <!-- 日历视图 -->
      <div v-if="viewMode === 'calendar'" class="calendar-view">
        <el-calendar v-model="currentDate">
          <template #date-cell="{ data }">
            <div class="calendar-cell" @click="selectDate(data)">
              <div class="date-label">{{ data.day }}</div>
              <div class="schedule-dots">
                <span
                  v-for="s in getSchedulesForDate(data)"
                  :key="s.id"
                  class="dot"
                  :class="getShiftTypeClass(s.shift_type)"
                  :title="`${s.staff_name} - ${s.shift_type}`"
                ></span>
              </div>
            </div>
          </template>
        </el-calendar>
      </div>

      <!-- 列表视图 -->
      <el-table v-else :data="schedules" v-loading="loading">
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column prop="staff_name" label="员工" width="120" />
        <el-table-column prop="shift_type" label="班次" width="100">
          <template #default="{ row }">
            <el-tag :type="getShiftTagType(row.shift_type)">
              {{ row.shift_type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="150">
          <template #default="{ row }">
            {{ row.start_time }} - {{ row.end_time }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="notes" label="备注" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="editSchedule(row)">编辑</el-button>
            <el-popconfirm title="确定删除?" @confirm="deleteSchedule(row.id)">
              <template #reference>
                <el-button type="danger" size="small" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        @current-change="loadSchedules"
        layout="total, prev, pager, next"
        style="margin-top: 20px; justify-content: center"
      />
    </el-card>

    <!-- 添加/编辑排班对话框 -->
    <el-dialog v-model="addDialogVisible" :title="isEdit ? '编辑排班' : '添加排班'" width="500px">
      <el-form :model="scheduleForm" label-width="100px" :rules="formRules" ref="formRef">
        <el-form-item label="员工" prop="staff_id">
          <el-select v-model="scheduleForm.staff_id" placeholder="请选择员工" style="width: 100%">
            <el-option v-for="staff in staffList" :key="staff.id" :label="staff.nickname" :value="staff.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期" prop="date">
          <el-date-picker v-model="scheduleForm.date" type="date" placeholder="选择日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="班次" prop="shift_type">
          <el-select v-model="scheduleForm.shift_type" placeholder="请选择班次" style="width: 100%">
            <el-option label="早班 (08:00-16:00)" value="早班" />
            <el-option label="中班 (16:00-00:00)" value="中班" />
            <el-option label="晚班 (00:00-08:00)" value="晚班" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="scheduleForm.notes" type="textarea" :rows="2" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSchedule" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import instance from '@/api'

interface Schedule {
  id: number
  staff_id: number
  staff_name: string
  date: string
  shift_type: string
  start_time: string
  end_time: string
  status: string
  notes: string
}

const loading = ref(false)
const submitting = ref(false)
const viewMode = ref('calendar')
const currentDate = ref(new Date())
const dateRange = ref<[Date, Date] | null>(null)
const schedules = ref<Schedule[]>([])
const staffList = ref<any[]>([])
const addDialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()

const monthlyStats = ref({
  total_shifts: 0,
  morning_shifts: 0,
  afternoon_shifts: 0,
  night_shifts: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const scheduleForm = reactive({
  staff_id: null as number | null,
  date: null as Date | null,
  shift_type: '',
  notes: ''
})

const formRules: FormRules = {
  staff_id: [{ required: true, message: '请选择员工', trigger: 'change' }],
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  shift_type: [{ required: true, message: '请选择班次', trigger: 'change' }]
}

const getShiftTagType = (type: string) => {
  const types: Record<string, string> = {
    '早班': 'warning',
    '中班': 'primary',
    '晚班': 'info'
  }
  return types[type] || 'default'
}

const getShiftTypeClass = (type: string) => {
  const classes: Record<string, string> = {
    '早班': 'morning',
    '中班': 'afternoon',
    '晚班': 'night'
  }
  return classes[type] || ''
}

const getStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    scheduled: 'info',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    scheduled: '已排班',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const getSchedulesForDate = (date: any) => {
  const dateStr = `${date.year}-${String(date.month).padStart(2, '0')}-${String(date.day).padStart(2, '0')}`
  return schedules.value.filter(s => s.date === dateStr)
}

const selectDate = (date: any) => {
  scheduleForm.date = new Date(date.year, date.month - 1, date.day)
  addDialogVisible.value = true
}

const loadSchedules = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (dateRange.value) {
      params.start_date = formatDate(dateRange.value[0])
      params.end_date = formatDate(dateRange.value[1])
    }
    const res = await instance.get('/schedules', { params })
    schedules.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadStaffList = async () => {
  try {
    const res = await instance.get('/staff', { params: { page: 1, page_size: 100 } })
    staffList.value = res.list || []
  } catch (error) {
    console.error(error)
  }
}

const loadMonthlyStats = async () => {
  try {
    const now = new Date()
    const res = await instance.get('/schedules/stats/monthly', {
      params: { year: now.getFullYear(), month: now.getMonth() + 1 }
    })
    monthlyStats.value = res || monthlyStats.value
  } catch (error) {
    console.error(error)
  }
}

const showAddDialog = () => {
  isEdit.value = false
  editId.value = 0
  scheduleForm.staff_id = null
  scheduleForm.date = null
  scheduleForm.shift_type = ''
  scheduleForm.notes = ''
  addDialogVisible.value = true
}

const editSchedule = (row: Schedule) => {
  isEdit.value = true
  editId.value = row.id
  scheduleForm.staff_id = row.staff_id
  scheduleForm.date = new Date(row.date)
  scheduleForm.shift_type = row.shift_type
  scheduleForm.notes = row.notes
  addDialogVisible.value = true
}

const submitSchedule = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const data = {
        staff_id: scheduleForm.staff_id,
        date: formatDate(scheduleForm.date!),
        shift_type: scheduleForm.shift_type,
        notes: scheduleForm.notes,
        start_time: scheduleForm.shift_type === '早班' ? '08:00' : scheduleForm.shift_type === '中班' ? '16:00' : '00:00',
        end_time: scheduleForm.shift_type === '早班' ? '16:00' : scheduleForm.shift_type === '中班' ? '00:00' : '08:00'
      }
      await instance.post('/schedules', data)
      ElMessage.success('排班成功')
      addDialogVisible.value = false
      loadSchedules()
      loadMonthlyStats()
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const deleteSchedule = async (id: number) => {
  try {
    await instance.delete(`/schedules/${id}`)
    ElMessage.success('删除成功')
    loadSchedules()
    loadMonthlyStats()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const formatDate = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

onMounted(() => {
  loadSchedules()
  loadStaffList()
  loadMonthlyStats()
})
</script>

<style scoped>
.schedule-page {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  padding: 20px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
}

.stat-card.morning .stat-value {
  color: #e6a23c;
}

.stat-card.afternoon .stat-value {
  color: #67c23a;
}

.stat-card.night .stat-value {
  color: #909399;
}

.stat-label {
  color: #666;
  margin-top: 8px;
}

.toolbar-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.calendar-view {
  padding: 10px;
}

.calendar-cell {
  min-height: 80px;
  padding: 8px;
  cursor: pointer;
}

.date-label {
  font-size: 14px;
  font-weight: 500;
}

.schedule-dots {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 4px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.dot.morning {
  background: #e6a23c;
}

.dot.afternoon {
  background: #67c23a;
}

.dot.night {
  background: #909399;
}
</style>
