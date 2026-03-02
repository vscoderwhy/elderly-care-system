<template>
  <div class="elderly-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>老人管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增老人
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="gender" label="性别" width="80" />
        <el-table-column prop="phone" label="联系电话" width="150" />
        <el-table-column prop="emergency_contact" label="紧急联系人" width="120" />
        <el-table-column prop="emergency_phone" label="紧急联系电话" width="150" />
        <el-table-column prop="care_level" label="护理等级" width="100" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="warning" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        @current-change="fetchData"
        layout="total, prev, pager, next"
        style="margin-top: 20px; justify-content: center"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { elderlyApi } from '@/api'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const result = await elderlyApi.list({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    tableData.value = result.list || []
    pagination.total = result.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  router.push('/elderly/create')
}

const handleView = (row: any) => {
  router.push(`/elderly/${row.id}`)
}

const handleEdit = (row: any) => {
  router.push(`/elderly/${row.id}`)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确认删除该老人档案吗？', '提示', {
      type: 'warning'
    })
    await elderlyApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    // User cancelled
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
