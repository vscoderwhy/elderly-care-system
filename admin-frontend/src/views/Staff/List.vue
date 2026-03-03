<template>
  <div class="staff-management">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><User /></el-icon>
          员工管理
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          添加员工
        </el-button>
        <el-button @click="handleImport">
          <el-icon><Upload /></el-icon>
          批量导入
        </el-button>
      </div>
    </div>

    <!-- 员工统计 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in staffStats" :key="stat.key">
        <div class="staff-stat" :class="`stat-${stat.type}`">
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

    <!-- 员工列表 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>员工列表</span>
          <div class="header-actions">
            <el-select v-model="filterForm.department" placeholder="全部部门" clearable style="width: 120px" @change="loadStaff">
              <el-option label="全部部门" value="" />
              <el-option label="护理部" value="nursing" />
              <el-option label="医务室" value="medical" />
              <el-option label="康复科" value="rehab" />
              <el-option label="膳食部" value="dining" />
              <el-option label="行政部" value="admin" />
            </el-select>
            <el-select v-model="filterForm.position" placeholder="全部职位" clearable style="width: 120px" @change="loadStaff">
              <el-option label="全部职位" value="" />
              <el-option label="院长" value="director" />
              <el-option label="医生" value="doctor" />
              <el-option label="护士" value="nurse" />
              <el-option label="护理员" value="caregiver" />
              <el-option label="康复师" value="therapist" />
              <el-option label="营养师" value="nutritionist" />
              <el-option label="护工" value="worker" />
            </el-select>
            <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 120px" @change="loadStaff">
              <el-option label="全部" value="" />
              <el-option label="在职" value="active" />
              <el-option label="离职" value="inactive" />
              <el-option label="休假" value="leave" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索姓名或工号"
              style="width: 200px"
              clearable
              @clear="loadStaff"
              @keyup.enter="loadStaff"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>

      <el-table :data="staffList" stripe v-loading="loading">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="employeeNo" label="工号" width="100" />
        <el-table-column label="姓名" width="100">
          <template #default="{ row }">
            <div class="name-cell">
              <el-avatar :size="32" :src="row.avatar">{{ row.name.charAt(0) }}</el-avatar>
              <span class="name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="60" />
        <el-table-column prop="department" label="部门" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.department }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="position" label="职位" width="100" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="hireDate" label="入职日期" width="110" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" text @click="handleSchedule(row)">排班</el-button>
            <el-dropdown trigger="click">
              <el-button size="small" text>
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleAttendance(row)">
                    <el-icon><Calendar /></el-icon>
                    考勤记录
                  </el-dropdown-item>
                  <el-dropdown-item @click="handlePerformance(row)">
                    <el-icon><TrendCharts /></el-icon>
                    绩效考核
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleTraining(row)">
                    <el-icon><Reading /></el-icon>
                    培训记录
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleResign(row)" v-if="row.status === 'active'">
                    <el-icon><Remove /></el-icon>
                    办理离职
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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

    <!-- 员工表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form :model="staffForm" :rules="staffRules" ref="staffFormRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="staffForm.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="gender">
              <el-radio-group v-model="staffForm.gender">
                <el-radio label="男">男</el-radio>
                <el-radio label="女">女</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="身份证号" prop="idCard">
              <el-input v-model="staffForm.idCard" placeholder="请输入身份证号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="staffForm.phone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部门" prop="department">
              <el-select v-model="staffForm.department" placeholder="请选择部门" style="width: 100%">
                <el-option label="护理部" value="nursing" />
                <el-option label="医务室" value="medical" />
                <el-option label="康复科" value="rehab" />
                <el-option label="膳食部" value="dining" />
                <el-option label="行政部" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="职位" prop="position">
              <el-select v-model="staffForm.position" placeholder="请选择职位" style="width: 100%">
                <el-option label="院长" value="director" />
                <el-option label="医生" value="doctor" />
                <el-option label="护士" value="nurse" />
                <el-option label="护理员" value="caregiver" />
                <el-option label="康复师" value="therapist" />
                <el-option label="营养师" value="nutritionist" />
                <el-option label="护工" value="worker" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="入职日期" prop="hireDate">
              <el-date-picker v-model="staffForm.hireDate" type="date" placeholder="选择日期" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学历" prop="education">
              <el-select v-model="staffForm.education" placeholder="请选择学历" style="width: 100%">
                <el-option label="高中" value="high" />
                <el-option label="大专" value="college" />
                <el-option label="本科" value="bachelor" />
                <el-option label="硕士" value="master" />
                <el-option label="博士" value="doctor" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="居住地址" prop="address">
          <el-input v-model="staffForm.address" type="textarea" :rows="2" placeholder="请输入居住地址" />
        </el-form-item>

        <el-form-item label="资格证书">
          <image-upload v-model="staffForm.certificates" :limit="5" />
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="staffForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
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
  User,
  Plus,
  Upload,
  Search,
  ArrowDown,
  Calendar,
  TrendCharts,
  Reading,
  Remove
} from '@element-plus/icons-vue'
import ImageUpload from '@/components/Upload/ImageUpload.vue'

const staffStats = ref([
  { key: 'total', label: '员工总数', value: 86, type: 'primary', icon: 'User' },
  { key: 'active', label: '在职', value: 78, type: 'success', icon: 'UserFilled' },
  { key: 'leave', label: '休假', value: 5, type: 'warning', icon: 'Clock' },
  { key: 'resigned', label: '离职', value: 3, type: 'danger', icon: 'User' }
])

const filterForm = reactive({
  department: '',
  position: '',
  status: ''
})

const searchKeyword = ref('')
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const saving = ref(false)
const staffFormRef = ref<FormInstance>()

const staffList = ref([
  {
    id: 1,
    employeeNo: 'EMP001',
    name: '张医生',
    gender: '男',
    department: '医务室',
    position: '医生',
    phone: '13800138001',
    hireDate: '2023-03-15',
    status: 'active',
    avatar: '',
    education: '本科',
    idCard: '',
    address: '',
    certificates: [],
    remark: ''
  },
  {
    id: 2,
    name: '李护士',
    gender: '女',
    department: '护理部',
    position: '护士',
    phone: '13800138002',
    hireDate: '2023-05-20',
    status: 'active',
    avatar: '',
    education: '大专',
    idCard: '',
    address: '',
    certificates: [],
    remark: ''
  },
  {
    id: 3,
    name: '王康复师',
    gender: '男',
    department: '康复科',
    position: '康复师',
    phone: '13800138003',
    hireDate: '2023-06-10',
    status: 'active',
    avatar: '',
    education: '本科',
    idCard: '',
    address: '',
    certificates: [],
    remark: ''
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(86)

const staffForm = reactive({
  id: '',
  employeeNo: '',
  name: '',
  gender: '男',
  idCard: '',
  phone: '',
  department: '',
  position: '',
  hireDate: '',
  education: '',
  address: '',
  certificates: [] as string[],
  remark: ''
})

const staffRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  idCard: [{ required: true, message: '请输入身份证号', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  department: [{ required: true, message: '请选择部门', trigger: 'change' }],
  position: [{ required: true, message: '请选择职位', trigger: 'change' }],
  hireDate: [{ required: true, message: '请选择入职日期', trigger: 'change' }]
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    active: 'success',
    inactive: 'danger',
    leave: 'warning'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '在职',
    inactive: '离职',
    leave: '休假'
  }
  return map[status] || status
}

const loadStaff = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
  }, 500)
}

const handleAdd = () => {
  dialogTitle.value = '添加员工'
  resetForm()
  dialogVisible.value = true
}

const handleImport = () => {
  console.log('批量导入')
}

const handleView = (row: any) => {
  console.log('查看员工', row)
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑员工'
  Object.assign(staffForm, row)
  dialogVisible.value = true
}

const handleSchedule = (row: any) => {
  console.log('排班管理', row)
}

const handleAttendance = (row: any) => {
  console.log('考勤记录', row)
}

const handlePerformance = (row: any) => {
  console.log('绩效考核', row)
}

const handleTraining = (row: any) => {
  console.log('培训记录', row)
}

const handleResign = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要为员工"${row.name}"办理离职吗？`, '提示', {
      type: 'warning'
    })
    ElMessage.success('已标记为离职')
    loadStaff()
  } catch {
    // 取消
  }
}

const resetForm = () => {
  Object.assign(staffForm, {
    id: '',
    employeeNo: '',
    name: '',
    gender: '男',
    idCard: '',
    phone: '',
    department: '',
    position: '',
    hireDate: '',
    education: '',
    address: '',
    certificates: [],
    remark: ''
  })
  staffFormRef.value?.clearValidate()
}

const handleSave = async () => {
  const valid = await staffFormRef.value?.validate()
  if (!valid) return

  saving.value = true
  setTimeout(() => {
    saving.value = false
    dialogVisible.value = false
    ElMessage.success('保存成功')
    loadStaff()
  }, 1000)
}

onMounted(() => {
  loadStaff()
})
</script>

<style scoped lang="scss">
.staff-management {
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

  .staff-stat {
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
      .stat-danger & { background: #fef0f0; color: #f56c6c; }
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

  .table-card {
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
    flex-wrap: wrap;
  }

  .name-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .name {
      font-weight: 500;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

@media (max-width: 768px) {
  .staff-management {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .header-actions {
      flex-wrap: wrap;

      * {
        flex-shrink: 0;
      }
    }
  }
}
</style>
