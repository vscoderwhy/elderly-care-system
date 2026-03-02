<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>探视预约</span>
        <div class="header-actions">
          <el-button @click="fetchTodayVisits">今日探视</el-button>
          <el-button type="primary" @click="showCreateDialog">新增预约</el-button>
        </div>
      </div>
    </template>

    <el-table :data="visits" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="elderly.name" label="老人姓名" width="120" />
      <el-table-column prop="visitor_name" label="访客姓名" width="120" />
      <el-table-column prop="visitor_phone" label="访客电话" width="150" />
      <el-table-column prop="relationship" label="关系" width="100" />
      <el-table-column prop="visit_date" label="探视日期" width="120">
        <template #default="{ row }">{{ formatDate(row.visit_date) }}</template>
      </el-table-column>
      <el-table-column prop="visit_time" label="探视时间" width="100" />
      <el-table-column prop="visitor_count" label="人数" width="80" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" size="small" @click="confirmVisit(row.id)">确认</el-button>
          <el-button v-if="row.status === 'pending'" size="small" @click="cancelVisit(row.id)">取消</el-button>
          <el-button v-if="row.status === 'confirmed'" size="small" @click="completeVisit(row.id)">完成</el-button>
          <el-button size="small" @click="editVisit(row)">编辑</el-button>
          <el-button v-if="row.status !== 'completed'" size="small" type="danger" @click="deleteVisit(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      @current-change="fetchVisits"
      layout="total, prev, pager, next"
    />

    <!-- 新增/编辑预约对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="老人" required>
          <el-select v-model="form.elderly_id" placeholder="选择老人" filterable>
            <el-option
              v-for="elderly in elderlyList"
              :key="elderly.id"
              :label="elderly.name"
              :value="elderly.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="访客姓名" required>
          <el-input v-model="form.visitor_name" />
        </el-form-item>
        <el-form-item label="访客电话" required>
          <el-input v-model="form.visitor_phone" />
        </el-form-item>
        <el-form-item label="关系" required>
          <el-input v-model="form.relationship" placeholder="如：父子、母女、夫妻等" />
        </el-form-item>
        <el-form-item label="探视日期" required>
          <el-date-picker v-model="form.visit_date" type="date" />
        </el-form-item>
        <el-form-item label="探视时间" required>
          <el-input v-model="form.visit_time" placeholder="如：09:00-10:00" />
        </el-form-item>
        <el-form-item label="人数">
          <el-input-number v-model="form.visitor_count" :min="1" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.notes" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveVisit">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from '@/api/index'

const visits = ref([])
const elderlyList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const form = reactive({
  id: null as number | null,
  elderly_id: null as number | null,
  visitor_name: '',
  visitor_phone: '',
  relationship: '',
  visit_date: '',
  visit_time: '',
  visitor_count: 1,
  notes: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const fetchVisits = async () => {
  loading.value = true
  try {
    const res = await axios.get('/visits', {
      params: { page: pagination.page, page_size: pagination.pageSize }
    })
    visits.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('获取探视预约列表失败')
  }
  loading.value = false
}

const fetchTodayVisits = async () => {
  loading.value = true
  try {
    const res = await axios.get('/visits/today')
    visits.value = res.data || []
    pagination.total = res.data?.length || 0
  } catch (error) {
    ElMessage.error('获取今日探视失败')
  }
  loading.value = false
}

const fetchElderlyList = async () => {
  try {
    const res = await axios.get('/elderly', { params: { page: 1, page_size: 1000 } })
    elderlyList.value = res.data.list || []
  } catch (error) {
    ElMessage.error('获取老人列表失败')
  }
}

const showCreateDialog = () => {
  dialogTitle.value = '新增探视预约'
  Object.assign(form, {
    id: null,
    elderly_id: null,
    visitor_name: '',
    visitor_phone: '',
    relationship: '',
    visit_date: '',
    visit_time: '',
    visitor_count: 1,
    notes: ''
  })
  dialogVisible.value = true
}

const editVisit = (row: any) => {
  dialogTitle.value = '编辑探视预约'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const saveVisit = async () => {
  try {
    const data = {
      ...form,
      visit_date: form.visit_date ? new Date(form.visit_date).toISOString().split('T')[0] : ''
    }
    if (form.id) {
      await axios.put(`/visits/${form.id}`, data)
      ElMessage.success('更新成功')
    } else {
      await axios.post('/visits', data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchVisits()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const confirmVisit = async (id: number) => {
  try {
    await axios.put(`/visits/${id}/confirm`)
    ElMessage.success('已确认')
    fetchVisits()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const cancelVisit = async (id: number) => {
  try {
    await axios.put(`/visits/${id}/cancel`)
    ElMessage.success('已取消')
    fetchVisits()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const completeVisit = async (id: number) => {
  try {
    await axios.put(`/visits/${id}/complete`)
    ElMessage.success('已完成')
    fetchVisits()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const deleteVisit = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除此预约吗？', '提示', {
      type: 'warning'
    })
    await axios.delete(`/visits/${id}`)
    ElMessage.success('删除成功')
    fetchVisits()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'info',
    confirmed: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchVisits()
  fetchElderlyList()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.el-pagination {
  margin-top: 20px;
}
</style>
