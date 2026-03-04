<template>
  <div class="visit-appointments">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <h3>
            <el-icon><Calendar /></el-icon>
            探视预约管理
          </h3>
          <div class="header-actions">
            <el-button type="primary" size="small" @click="showAddDialog">
              <el-icon><Plus /></el-icon>
              新建预约
            </el-button>
          </div>
        </div>
      </template>

      <div class="list-view">
        <div class="filter-bar">
          <el-form :inline="true" :model="filterForm">
            <el-form-item label="状态">
              <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 120px">
                <el-option label="待审核" value="pending" />
                <el-option label="已确认" value="approved" />
                <el-option label="已完成" value="completed" />
                <el-option label="已取消" value="cancelled" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadAppointments">查询</el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-table :data="appointments" border v-loading="loading">
          <el-table-column prop="appointmentNo" label="预约编号" width="140" />
          <el-table-column prop="elderlyName" label="老人姓名" width="100" />
          <el-table-column prop="visitorName" label="访客姓名" width="100" />
          <el-table-column prop="visitorPhone" label="联系电话" width="130" />
          <el-table-column prop="appointmentDate" label="预约日期" width="120" />
          <el-table-column prop="timeSlot" label="时间段" width="100" />
          <el-table-column prop="status" label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button size="small" text @click="viewDetail(row)">查看</el-button>
              <el-button size="small" text type="primary" @click="approveAppointment(row)" v-if="row.status === 'pending'">审核</el-button>
              <el-button size="small" text type="danger" @click="cancelAppointment(row)" v-if="row.status !== 'completed'">取消</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <el-dialog v-model="appointmentDialogVisible" title="新建预约" width="700px">
      <el-form :model="appointmentForm" :rules="appointmentFormRules" ref="appointmentFormRef" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="老人姓名" prop="elderlyId">
              <el-select v-model="appointmentForm.elderlyId" placeholder="请选择" filterable style="width: 100%">
                <el-option v-for="elderly in elderlyList" :key="elderly.id" :label="elderly.name" :value="elderly.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="预约日期" prop="appointmentDate">
              <el-date-picker v-model="appointmentForm.appointmentDate" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="访客姓名" prop="visitorName">
              <el-input v-model="appointmentForm.visitorName" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="visitorPhone">
              <el-input v-model="appointmentForm.visitorPhone" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input v-model="appointmentForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="appointmentDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAppointment">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="预约详情" width="600px">
      <el-descriptions :column="2" border v-if="currentAppointment">
        <el-descriptions-item label="预约编号">{{ currentAppointment.appointmentNo }}</el-descriptions-item>
        <el-descriptions-item label="老人姓名">{{ currentAppointment.elderlyName }}</el-descriptions-item>
        <el-descriptions-item label="访客姓名">{{ currentAppointment.visitorName }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ currentAppointment.visitorPhone }}</el-descriptions-item>
        <el-descriptions-item label="预约日期">{{ currentAppointment.appointmentDate }}</el-descriptions-item>
        <el-descriptions-item label="时间段">{{ currentAppointment.timeSlot }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Calendar, Plus } from '@element-plus/icons-vue'
import { getElderlyList } from '@/utils/seedData'

interface Appointment {
  id: number
  appointmentNo: string
  elderlyId: number
  elderlyName: string
  visitorName: string
  visitorPhone: string
  appointmentDate: string
  timeSlot: string
  status: string
  notes: string
}

const elderlyList = getElderlyList()
const loading = ref(false)
const appointments = ref<Appointment[]>([])
const appointmentDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentAppointment = ref<Appointment | null>(null)
const appointmentFormRef = ref<FormInstance>()

const filterForm = reactive({ status: '' })

const appointmentForm = reactive<Appointment>({
  id: 0, appointmentNo: '', elderlyId: 0, elderlyName: '', visitorName: '', 
  visitorPhone: '', appointmentDate: '', timeSlot: '09:00-11:00', status: 'pending', notes: ''
})

const appointmentFormRules: FormRules = {
  elderlyId: [{ required: true, message: '请选择老人', trigger: 'change' }],
  visitorName: [{ required: true, message: '请输入访客姓名', trigger: 'blur' }],
  visitorPhone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }]
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = { pending: 'warning', approved: 'primary', completed: 'success', cancelled: 'info' }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待审核', approved: '已确认', completed: '已完成', cancelled: '已取消' }
  return map[status] || status
}

const loadAppointments = async () => {
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 500))
  appointments.value = Array.from({ length: 20 }, (_, i) => ({
    id: i + 1,
    appointmentNo: `APT${Date.now()}${i}`,
    elderlyId: elderlyList[i % elderlyList.length].id,
    elderlyName: elderlyList[i % elderlyList.length].name,
    visitorName: `访客${i + 1}`,
    visitorPhone: '1380013800' + String(i).padStart(2, '0'),
    appointmentDate: new Date().toISOString().split('T')[0],
    timeSlot: ['09:00-11:00', '14:00-16:00'][i % 2],
    status: ['pending', 'approved', 'completed'][i % 3],
    notes: ''
  }))
  loading.value = false
}

const showAddDialog = () => {
  Object.assign(appointmentForm, { id: 0, elderlyId: 0, visitorName: '', visitorPhone: '', appointmentDate: '', timeSlot: '09:00-11:00', status: 'pending', notes: '' })
  appointmentDialogVisible.value = true
}

const saveAppointment = async () => {
  if (!appointmentFormRef.value) return
  await appointmentFormRef.value.validate()
  const elderly = elderlyList.find(e => e.id === appointmentForm.elderlyId)
  appointments.value.unshift({
    ...appointmentForm,
    id: Date.now(),
    appointmentNo: `APT${Date.now()}`,
    elderlyName: elderly?.name || ''
  })
  ElMessage.success('预约创建成功')
  appointmentDialogVisible.value = false
}

const viewDetail = (apt: Appointment) => {
  currentAppointment.value = apt
  detailDialogVisible.value = true
}

const approveAppointment = (apt: Appointment) => { apt.status = 'approved'; ElMessage.success('已通过') }
const cancelAppointment = (apt: Appointment) => { apt.status = 'cancelled'; ElMessage.success('已取消') }

onMounted(() => { loadAppointments() })
</script>

<style scoped lang="scss">
.visit-appointments { .filter-bar { margin-bottom: 16px; } }
</style>
