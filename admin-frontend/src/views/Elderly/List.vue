<template>
  <div class="elderly-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><User /></el-icon>
          老人管理
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          添加老人
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="姓名">
          <el-input v-model="searchForm.name" placeholder="请输入姓名" clearable />
        </el-form-item>
        <el-form-item label="护理等级">
          <el-select v-model="searchForm.careLevel" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="三级护理" value="level3" />
            <el-option label="二级护理" value="level2" />
            <el-option label="一级护理" value="level1" />
            <el-option label="特级护理" value="special" />
          </el-select>
        </el-form-item>
        <el-form-item label="床位">
          <el-input v-model="searchForm.bedNumber" placeholder="请输入床位号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="在院" value="active" />
            <el-option label="请假" value="leave" />
            <el-option label="住院" value="hospital" />
            <el-option label="退住" value="discharged" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <span>共 {{ total }} 位老人</span>
            <el-divider direction="vertical" />
            <el-tag type="success" size="small">在院 {{ statistics.active }} 人</el-tag>
            <el-tag type="warning" size="small">请假 {{ statistics.leave }} 人</el-tag>
          </div>
          <div class="header-actions">
            <ExportButton
              :data="tableData"
              :columns="tableColumns"
              :total="total"
              :export-function="handleExport"
              filename="老人名单"
            />
          </div>
        </div>
      </template>

      <el-table
        ref="tableRef"
        v-loading="loading"
        :data="tableData"
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :src="row.avatar" :size="50">
              {{ row.name?.charAt(0) }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="gender" label="性别" width="60">
          <template #default="{ row }">
            <el-tag :type="row.gender === '男' ? 'primary' : 'danger'" size="small">
              {{ row.gender }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="age" label="年龄" width="70" />
        <el-table-column prop="careLevel" label="护理等级" width="100">
          <template #default="{ row }">
            <el-tag :type="getCareLevelType(row.careLevel)" size="small">
              {{ getCareLevelText(row.careLevel) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="checkInDate" label="入住日期" width="120" />
        <el-table-column prop="familyName" label="家属联系人" width="120" />
        <el-table-column prop="familyPhone" label="家属电话" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">
              查看
            </el-button>
            <el-button size="small" text @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button size="small" text @click="handleHealth(row)">
              健康
            </el-button>
            <el-dropdown trigger="click">
              <el-button size="small" text>
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleCare(row)">
                    <el-icon><Briefcase /></el-icon>
                    护理记录
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleBill(row)">
                    <el-icon><Wallet /></el-icon>
                    费用账单
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleLeave(row)">
                    <el-icon><CircleClose /></el-icon>
                    请假/退住
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="formData.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="gender">
              <el-radio-group v-model="formData.gender">
                <el-radio label="男">男</el-radio>
                <el-radio label="女">女</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="出生日期" prop="birthDate">
              <el-date-picker
                v-model="formData.birthDate"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="身份证号" prop="idCard">
              <el-input v-model="formData.idCard" placeholder="请输入身份证号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="护理等级" prop="careLevel">
              <el-select v-model="formData.careLevel" placeholder="请选择护理等级">
                <el-option label="三级护理" value="level3" />
                <el-option label="二级护理" value="level2" />
                <el-option label="一级护理" value="level1" />
                <el-option label="特级护理" value="special" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="床位号" prop="bedNumber">
              <el-input v-model="formData.bedNumber" placeholder="请输入床位号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="入住日期" prop="checkInDate">
              <el-date-picker
                v-model="formData.checkInDate"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="formData.phone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="家属联系人" prop="familyName">
          <el-input v-model="formData.familyName" placeholder="请输入家属姓名" />
        </el-form-item>

        <el-form-item label="家属电话" prop="familyPhone">
          <el-input v-model="formData.familyPhone" placeholder="请输入家属电话" />
        </el-form-item>

        <el-form-item label="与老人关系" prop="familyRelation">
          <el-input v-model="formData.familyRelation" placeholder="如：父子、母女等" />
        </el-form-item>

        <el-form-item label="健康状况" prop="healthStatus">
          <el-input
            v-model="formData.healthStatus"
            type="textarea"
            :rows="3"
            placeholder="请输入健康状况描述"
          />
        </el-form-item>

        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注"
          />
        </el-form-item>

        <el-form-item label="照片" prop="avatar">
          <ImageUpload v-model="formData.avatars" :limit="3" tip="最多上传3张照片" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  User,
  Plus,
  Search,
  RefreshLeft,
  ArrowDown,
  Briefcase,
  Wallet,
  CircleClose,
  Delete
} from '@element-plus/icons-vue'
import ImageUpload from '@/components/Upload/ImageUpload.vue'
import ExportButton from '@/components/Export/ExportButton.vue'
import { exportElderlyList } from '@/utils/export'
import { getElderlyList } from '@/utils/seedData'

const router = useRouter()

// 搜索表单
const searchForm = reactive({
  name: '',
  careLevel: '',
  bedNumber: '',
  status: ''
})

// 表格数据
const loading = ref(false)
const tableRef = ref()
const tableData = ref([])
const selectedRows = ref([])
const total = ref(0)

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20
})

// 统计数据
const statistics = ref({
  active: 248,
  leave: 12,
  hospital: 3,
  discharged: 5
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: '',
  name: '',
  gender: '男',
  birthDate: '',
  idCard: '',
  careLevel: '',
  bedNumber: '',
  checkInDate: '',
  phone: '',
  familyName: '',
  familyPhone: '',
  familyRelation: '',
  healthStatus: '',
  remark: '',
  avatars: [] as string[]
})

// 表单验证规则
const formRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  idCard: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { pattern: /^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/, message: '身份证号格式不正确', trigger: 'blur' }
  ],
  careLevel: [{ required: true, message: '请选择护理等级', trigger: 'change' }],
  bedNumber: [{ required: true, message: '请输入床位号', trigger: 'blur' }],
  checkInDate: [{ required: true, message: '请选择入住日期', trigger: 'change' }],
  phone: [{ pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }],
  familyName: [{ required: true, message: '请输入家属姓名', trigger: 'blur' }],
  familyPhone: [
    { required: true, message: '请输入家属电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ]
}

// 表格列配置
const tableColumns = [
  { key: 'name', title: '姓名', width: 12 },
  { key: 'gender', title: '性别', width: 8 },
  { key: 'age', title: '年龄', width: 8 },
  { key: 'idCard', title: '身份证号', width: 20 },
  { key: 'careLevel', title: '护理等级', width: 12 },
  { key: 'bedNumber', title: '床位号', width: 12 },
  { key: 'checkInDate', title: '入住日期', width: 15 },
  { key: 'familyName', title: '家属联系人', width: 12 },
  { key: 'familyPhone', title: '家属电话', width: 15 },
  { key: 'healthStatus', title: '健康状况', width: 15 },
  { key: 'status', title: '状态', width: 10 }
]

// 获取护理等级类型
const getCareLevelType = (level: string) => {
  const map: Record<string, any> = {
    level3: 'info',
    level2: '',
    level1: 'warning',
    special: 'danger'
  }
  return map[level] || ''
}

// 获取护理等级文本
const getCareLevelText = (level: string) => {
  const map: Record<string, string> = {
    level3: '三级护理',
    level2: '二级护理',
    level1: '一级护理',
    special: '特级护理'
  }
  return map[level] || level
}

// 获取状态类型
const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    active: 'success',
    leave: 'warning',
    hospital: 'danger',
    discharged: 'info'
  }
  return map[status] || ''
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '在院',
    leave: '请假',
    hospital: '住院',
    discharged: '退住'
  }
  return map[status] || status
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 使用种子数据
    await new Promise(resolve => setTimeout(resolve, 300))
    let data = getElderlyList()

    // 应用搜索过滤
    if (searchForm.name) {
      data = data.filter(e => e.name.includes(searchForm.name))
    }
    if (searchForm.careLevel) {
      const levelMap: Record<string, string> = {
        'level3': '三级',
        'level2': '二级',
        'level1': '一级',
        'special': '特级'
      }
      data = data.filter(e => e.careLevel === levelMap[searchForm.careLevel])
    }
    if (searchForm.bedNumber) {
      data = data.filter(e => e.bedNumber.includes(searchForm.bedNumber))
    }
    if (searchForm.status) {
      const statusMap: Record<string, string> = {
        'active': '在院',
        'leave': '外出',
        'hospital': '住院',
        'discharged': '退住'
      }
      data = data.filter(e => e.status === statusMap[searchForm.status])
    }

    // 分页
    const start = (pagination.page - 1) * pagination.pageSize
    const end = start + pagination.pageSize
    tableData.value = data.slice(start, end).map(e => ({
      ...e,
      careLevel: e.careLevel.includes('级') ? (
        e.careLevel === '特级' ? 'special' :
        e.careLevel === '一级' ? 'level1' :
        e.careLevel === '二级' ? 'level2' : 'level3'
      ) : e.careLevel,
      status: (
        e.status === '在院' ? 'active' :
        e.status === '外出' ? 'leave' :
        e.status === '住院' ? 'hospital' : 'discharged'
      ),
      checkInDate: e.admitDate || e.checkInDate,
      familyName: e.emergencyContact,
      familyPhone: e.emergencyPhone,
      healthStatus: `健康评分: ${e.healthScore}`
    }))
    total.value = data.length

    // 更新统计
    statistics.value = {
      active: data.filter(e => e.status === '在院').length,
      leave: data.filter(e => e.status === '外出').length,
      hospital: data.filter(e => e.status === '住院').length,
      discharged: data.filter(e => e.status === '退住').length
    }
  } catch (error) {
    console.error('加载数据失败', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    careLevel: '',
    bedNumber: '',
    status: ''
  })
  handleSearch()
}

// 选择变化
const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadData()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加老人'
  resetForm()
  dialogVisible.value = true
}

// 查看
const handleView = (row: any) => {
  router.push(`/elderly/detail/${row.id}`)
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑老人信息'
  Object.assign(formData, {
    ...row,
    avatars: row.avatar ? [row.avatar] : []
  })
  dialogVisible.value = true
}

// 健康档案
const handleHealth = (row: any) => {
  router.push(`/elderly/health/${row.id}`)
}

// 护理记录
const handleCare = (row: any) => {
  router.push(`/care/records?elderlyId=${row.id}`)
}

// 费用账单
const handleBill = (row: any) => {
  router.push(`/finance/bills?elderlyId=${row.id}`)
}

// 请假/退住
const handleLeave = (row: any) => {
  ElMessageBox.confirm(
    `确定要为"${row.name}"办理请假/退住手续吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('操作成功')
    loadData()
  })
}

// 删除
const handleDelete = (row: any) => {
  ElMessageBox.confirm(
    `确定要删除"${row.name}"的信息吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('删除成功')
    loadData()
  })
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: '',
    name: '',
    gender: '男',
    birthDate: '',
    idCard: '',
    careLevel: '',
    bedNumber: '',
    checkInDate: '',
    phone: '',
    familyName: '',
    familyPhone: '',
    familyRelation: '',
    healthStatus: '',
    remark: '',
    avatars: []
  })
  formRef.value?.clearValidate()
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    // TODO: 实际 API 请求
    await new Promise(resolve => setTimeout(resolve, 1000))

    ElMessage.success(formData.id ? '修改成功' : '添加成功')
    dialogVisible.value = false
    loadData()
  } catch (error) {
    console.error('提交失败', error)
  } finally {
    submitting.value = false
  }
}

// 导出
const handleExport = async (params: any) => {
  if (params.type === 'all') {
    // 获取全部数据
    const allData = [] // await getAllElderlyList()
    exportElderlyList(allData)
  } else {
    exportElderlyList(params.data)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.elderly-list {
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
    color: var(--text-primary);
    margin: 0;
  }

  .search-card {
    margin-bottom: 20px;
  }

  .table-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .header-title {
      display: flex;
      align-items: center;
      gap: 12px;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
