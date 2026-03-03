<template>
  <div class="visit-management">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Phone /></el-icon>
          探视预约
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新建预约
        </el-button>
      </div>
    </div>

    <!-- 预约统计 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in visitStats" :key="stat.key">
        <div class="visit-stat" :class="`stat-${stat.type}`">
          <div class="stat-icon">
            <component :is="stat.icon" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 预约日历 -->
    <el-card shadow="never" class="calendar-card">
      <template #header>
        <div class="card-header">
          <span>预约日历</span>
          <div class="header-actions">
            <el-date-picker v-model="calendarDate" type="date" placeholder="选择日期" @change="loadAppointments" />
          </div>
        </div>
      </template>

      <el-calendar v-model="calendarDate">
        <template #date-cell="{ data }">
          <div class="calendar-day" @click="handleDateClick(data)">
            <div class="day-number">{{ data.day.split('-').slice(-1)[0] }}</div>
            <div class="day-visits" v-if="getVisitsForDate(data).length > 0">
              <el-badge :value="getVisitsForDate(data).length" :max="99">
                <span class="visit-indicator"></span>
              </el-badge>
            </div>
          </div>
        </template>
      </el-calendar>
    </el-card>

    <!-- 预约列表 -->
    <el-card shadow="never" class="list-card">
      <template #header>
        <div class="card-header">
          <span>预约列表</span>
          <div class="header-actions">
            <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 120px" @change="loadAppointments">
              <el-option label="全部" value="" />
              <el-option label="待审核" value="pending" />
              <el-option label="已通过" value="approved" />
              <el-option label="已拒绝" value="rejected" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索老人或访客姓名"
              style="width: 200px"
              clearable
              @clear="loadAppointments"
              @keyup.enter="loadAppointments"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>

      <el-table :data="appointmentList" stripe v-loading="loading">
        <el-table-column prop="appointmentNo" label="预约编号" width="150" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="visitorName" label="访客姓名" width="100" />
        <el-table-column prop="visitorPhone" label="联系电话" width="130" />
        <el-table-column prop="visitDate" label="预约日期" width="110" />
        <el-table-column prop="visitTime" label="预约时间" width="100" />
        <el-table-column prop="visitType" label="探访类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.visitType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handleApprove(row)" v-if="row.status === 'pending'">通过</el-button>
            <el-button size="small" text @click="handleReject(row)" v-if="row.status === 'pending'">拒绝</el-button>
            <el-button size="small" text @click="handleComplete(row)" v-if="row.status === 'approved'">完成</el-button>
            <el-button size="small" text type="danger" @click="handleCancel(row)" v-if="!['completed', 'cancelled', 'rejected'].includes(row.status)">取消</el-button>
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
        />
      </div>
    </el-card>

    <!-- 预约表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="appointmentForm" :rules="appointmentRules" ref="appointmentFormRef" label-width="100px">
        <el-form-item label="老人" prop="elderlyId">
          <el-select v-model="appointmentForm.elderlyId" placeholder="请选择老人" filterable style="width: 100%">
            <el-option
              v-for="elderly in elderlyList"
              :key="elderly.id"
              :label="`${elderly.name} - ${elderly.bedNumber}`"
              :value="elderly.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="访客姓名" prop="visitorName">
          <el-input v-model="appointmentForm.visitorName" placeholder="请输入访客姓名" />
        </el-form-item>

        <el-form-item label="联系电话" prop="visitorPhone">
          <el-input v-model="appointmentForm.visitorPhone" placeholder="请输入联系电话" />
        </el-form-item>

        <el-form-item label="与老人关系" prop="relationship">
          <el-select v-model="appointmentForm.relationship" placeholder="请选择关系" style="width: 100%">
            <el-option label="子女" value="child" />
            <el-option label="配偶" value="spouse" />
            <el-option label="孙辈" value="grandchild" />
            <el-option label="其他亲属" value="relative" />
            <el-option label="朋友" value="friend" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="探访类型" prop="visitType">
          <el-radio-group v-model="appointmentForm.visitType">
            <el-radio label="现场探访">现场探访</el-radio>
            <el-radio label="视频探访">视频探访</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="预约日期" prop="visitDate">
          <el-date-picker v-model="appointmentForm.visitDate" type="date" placeholder="选择日期" style="width: 100%" />
        </el-form-item>

        <el-form-item label="预约时间" prop="visitTime">
          <el-select v-model="appointmentForm.visitTime" placeholder="请选择时间" style="width: 100%">
            <el-option label="09:00-10:00" value="09:00" />
            <el-option label="10:00-11:00" value="10:00" />
            <el-option label="11:00-12:00" value="11:00" />
            <el-option label="14:00-15:00" value="14:00" />
            <el-option label="15:00-16:00" value="15:00" />
            <el-option label="16:00-17:00" value="16:00" />
          </el-select>
        </el-form-item>

        <el-form-item label="探访人数" prop="visitorCount">
          <el-input-number v-model="appointmentForm.visitorCount" :min="1" :max="10" />
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="appointmentForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Phone,
  Plus,
  Search
} from '@element-plus/icons-vue'

const visitStats = ref([
  { key: 'today', label: '今日预约', value: 12, type: 'primary', icon: 'Calendar' },
  { key: 'pending', label: '待审核', value: 5, type: 'warning', icon: 'Clock' },
  { key: 'approved', label: '已通过', value: 8, type: 'success', icon: 'Select' },
  { key: 'week', label: '本周预约', value: 45, type: 'info', icon: 'Calendar' }
])

const calendarDate = ref(new Date())
const filterForm = reactive({
  status: ''
})
const searchKeyword = ref('')
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const saving = ref(false)
const appointmentFormRef = ref<FormInstance>()

const appointmentList = ref([
  {
    id: 1,
    appointmentNo: 'VA202603001',
    elderlyId: 1,
    elderlyName: '张奶奶',
    bedNumber: '3号楼201',
    visitorName: '王先生',
    visitorPhone: '13900139001',
    relationship: 'child',
    visitType: '现场探访',
    visitDate: '2026-03-04',
    visitTime: '09:00',
    visitorCount: 2,
    status: 'pending',
    remark: ''
  },
  {
    id: 2,
    appointmentNo: 'VA202603002',
    elderlyId: 2,
    elderlyName: '李爷爷',
    bedNumber: '2号楼105',
    visitorName: '赵女士',
    visitorPhone: '13900139002',
    relationship: 'spouse',
    visitType: '视频探访',
    visitDate: '2026-03-04',
    visitTime: '14:00',
    visitorCount: 1,
    status: 'approved',
    remark: ''
  },
  {
    id: 3,
    appointmentNo: 'VA202603003',
    elderlyId: 3,
    elderlyName: '王奶奶',
    bedNumber: '3号楼202',
    visitorName: '孙先生',
    visitorPhone: '13900139003',
    relationship: 'grandchild',
    visitType: '现场探访',
    visitDate: '2026-03-05',
    visitTime: '10:00',
    visitorCount: 3,
    status: 'approved',
    remark: ''
  }
])

const elderlyList = ref([
  { id: 1, name: '张奶奶', bedNumber: '3号楼201' },
  { id: 2, name: '李爷爷', bedNumber: '2号楼105' },
  { id: 3, name: '王奶奶', bedNumber: '3号楼202' }
])

const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(45)

const appointmentForm = reactive({
  id: '',
  elderlyId: '',
  visitorName: '',
  visitorPhone: '',
  relationship: '',
  visitType: '现场探访',
  visitDate: '',
  visitTime: '',
  visitorCount: 1,
  remark: ''
})

const appointmentRules: FormRules = {
  elderlyId: [{ required: true, message: '请选择老人', trigger: 'change' }],
  visitorName: [{ required: true, message: '请输入访客姓名', trigger: 'blur' }],
  visitorPhone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  relationship: [{ required: true, message: '请选择关系', trigger: 'change' }],
  visitType: [{ required: true, message: '请选择探访类型', trigger: 'change' }],
  visitDate: [{ required: true, message: '请选择预约日期', trigger: 'change' }],
  visitTime: [{ required: true, message: '请选择预约时间', trigger: 'change' }],
  visitorCount: [{ required: true, message: '请输入探访人数', trigger: 'blur' }]
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    completed: 'info',
    cancelled: 'info'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

const getVisitsForDate = (data: any) => {
  const dateStr = data.day
  return appointmentList.value.filter(item => item.visitDate === dateStr)
}

const handleDateClick = (data: any) => {
  const visits = getVisitsForDate(data)
  if (visits.length > 0) {
    console.log('该日期的预约:', visits)
  }
}

const loadAppointments = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
  }, 500)
}

const handleAdd = () => {
  dialogTitle.value = '新建预约'
  resetForm()
  dialogVisible.value = true
}

const handleView = (row: any) => {
  console.log('查看预约', row)
}

const handleApprove = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要通过"${row.visitorName}"的预约申请吗？`, '提示', {
      type: 'success'
    })
    row.status = 'approved'
    ElMessage.success('已通过')
    loadAppointments()
  } catch {
    // 取消
  }
}

const handleReject = async (row: any) => {
  try {
    await ElMessageBox.prompt('请输入拒绝原因', '拒绝预约', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    row.status = 'rejected'
    ElMessage.success('已拒绝')
    loadAppointments()
  } catch {
    // 取消
  }
}

const handleComplete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确认"${row.visitorName}"已完成探视吗？`, '提示', {
      type: 'success'
    })
    row.status = 'completed'
    ElMessage.success('已完成')
    loadAppointments()
  } catch {
    // 取消
  }
}

const handleCancel = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要取消"${row.visitorName}"的预约吗？`, '提示', {
      type: 'warning'
    })
    row.status = 'cancelled'
    ElMessage.success('已取消')
    loadAppointments()
  } catch {
    // 取消
  }
}

const resetForm = () => {
  Object.assign(appointmentForm, {
    id: '',
    elderlyId: '',
    visitorName: '',
    visitorPhone: '',
    relationship: '',
    visitType: '现场探访',
    visitDate: '',
    visitTime: '',
    visitorCount: 1,
    remark: ''
  })
  appointmentFormRef.value?.clearValidate()
}

const handleSave = async () => {
  const valid = await appointmentFormRef.value?.validate()
  if (!valid) return

  saving.value = true
  setTimeout(() => {
    saving.value = false
    dialogVisible.value = false
    ElMessage.success('保存成功')
    loadAppointments()
  }, 1000)
}

onMounted(() => {
  loadAppointments()
})
</script>

<style scoped lang="scss">
.visit-management {
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

  .stats-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .visit-stat {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;

      .stat-primary & { background: #e1f3ff; color: #409eff; }
      .stat-success & { background: #e1f3e8; color: #67c23a; }
      .stat-warning & { background: #fdf6ec; color: #e6a23c; }
      .stat-info & { background: #f4f4f5; color: #909399; }
    }

    .stat-info {
      flex: 1;
    }

    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 4px;
    }

    .stat-label {
      font-size: 14px;
      color: var(--text-secondary);
    }
  }

  .calendar-card,
  .list-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }

  .calendar-day {
    height: 80px;
    padding: 8px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      background: var(--fill-color-light);
    }

    .day-number {
      font-size: 16px;
      font-weight: 500;
      color: var(--text-primary);
      margin-bottom: 8px;
    }

    .day-visits {
      display: flex;
      justify-content: center;
    }

    .visit-indicator {
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--primary-color);
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

@media (max-width: 768px) {
  .visit-management {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .header-actions {
      flex-wrap: wrap;
    }
  }
}
</style>
