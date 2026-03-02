<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>角色管理</span>
        <el-button type="primary" @click="initializeSystem">初始化系统数据</el-button>
      </div>
    </template>

    <el-table :data="roles" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="角色名称" width="200" />
      <el-table-column prop="description" label="描述" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="viewRole(row)">查看详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 角色详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="角色详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ currentRole.id }}</el-descriptions-item>
        <el-descriptions-item label="角色名称">{{ currentRole.name }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentRole.description }}</el-descriptions-item>
      </el-descriptions>

      <el-divider>权限</el-divider>
      <div class="permissions-list">
        <el-tag v-for="perm in permissions" :key="perm.id" size="small" style="margin: 5px">
          {{ perm.name }}
        </el-tag>
        <span v-if="permissions.length === 0" style="color: #999">暂无权限</span>
      </div>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/api/index'

const roles = ref([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const currentRole = reactive({
  id: 0,
  name: '',
  description: ''
})
const permissions = ref([])

const fetchRoles = async () => {
  loading.value = true
  try {
    const res = await axios.get('/system/roles')
    roles.value = res.data || []
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  }
  loading.value = false
}

const viewRole = async (row: any) => {
  try {
    const res = await axios.get(`/system/roles/${row.id}`)
    Object.assign(currentRole, res.data)
    permissions.value = res.data.permissions || []
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取角色详情失败')
  }
}

const initializeSystem = async () => {
  try {
    await ElMessageBox.confirm('确定初始化系统数据吗？这将创建默认权限、菜单和管理员角色。', '提示', {
      type: 'warning'
    })
    await axios.post('/system/init')
    ElMessage.success('系统初始化成功')
    fetchRoles()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '初始化失败')
    }
  }
}

onMounted(() => {
  fetchRoles()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.permissions-list {
  min-height: 50px;
}
</style>
