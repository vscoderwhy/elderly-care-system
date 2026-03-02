<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>用药管理</span>
        <el-button type="primary" @click="showCreateDialog">新增药品</el-button>
      </div>
    </template>

    <el-table :data="medications" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="药品名称" width="200" />
      <el-table-column prop="specification" label="规格" width="120" />
      <el-table-column prop="manufacturer" label="生产厂家" width="200" />
      <el-table-column prop="stock" label="库存" width="100">
        <template #default="{ row }">
          <el-tag :type="row.stock <= row.min_stock ? 'danger' : 'success'">{{ row.stock }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="unit" label="单位" width="80" />
      <el-table-column prop="price" label="单价" width="100">
        <template #default="{ row }">¥{{ row.price }}</template>
      </el-table-column>
      <el-table-column prop="expiry_date" label="有效期" width="120">
        <template #default="{ row }">
          {{ formatDate(row.expiry_date) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="editMedication(row)">编辑</el-button>
          <el-button size="small" @click="viewRecords(row)">用药记录</el-button>
          <el-button size="small" type="danger" @click="deleteMedication(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      @current-change="fetchMedications"
      layout="total, prev, pager, next"
    />

    <!-- 新增/编辑药品对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="药品名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="规格">
          <el-input v-model="form.specification" />
        </el-form-item>
        <el-form-item label="生产厂家">
          <el-input v-model="form.manufacturer" />
        </el-form-item>
        <el-form-item label="库存" required>
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
        <el-form-item label="最低库存">
          <el-input-number v-model="form.min_stock" :min="0" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="form.unit" />
        </el-form-item>
        <el-form-item label="单价">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="有效期">
          <el-date-picker v-model="form.expiry_date" type="date" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveMedication">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from '@/api/index'

const medications = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const form = reactive({
  id: null as number | null,
  name: '',
  specification: '',
  manufacturer: '',
  stock: 0,
  min_stock: 10,
  unit: '',
  price: 0,
  expiry_date: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const fetchMedications = async () => {
  loading.value = true
  try {
    const res = await axios.get('/medications', {
      params: { page: pagination.page, page_size: pagination.pageSize }
    })
    medications.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('获取药品列表失败')
  }
  loading.value = false
}

const showCreateDialog = () => {
  dialogTitle.value = '新增药品'
  Object.assign(form, {
    id: null,
    name: '',
    specification: '',
    manufacturer: '',
    stock: 0,
    min_stock: 10,
    unit: '',
    price: 0,
    expiry_date: ''
  })
  dialogVisible.value = true
}

const editMedication = (row: any) => {
  dialogTitle.value = '编辑药品'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const saveMedication = async () => {
  try {
    if (form.id) {
      await axios.put(`/medications/${form.id}`, form)
      ElMessage.success('更新成功')
    } else {
      await axios.post('/medications', form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchMedications()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const deleteMedication = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除此药品吗？', '提示', {
      type: 'warning'
    })
    await axios.delete(`/medications/${id}`)
    ElMessage.success('删除成功')
    fetchMedications()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

const viewRecords = (row: any) => {
  ElMessage.info('用药记录功能开发中')
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchMedications()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-pagination {
  margin-top: 20px;
}
</style>
