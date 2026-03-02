<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>用户管理</span>
        <el-button type="primary" @click="showCreateDialog">新增用户</el-button>
      </div>
    </template>

    <el-table :data="users" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="phone" label="手机号" width="150" />
      <el-table-column prop="nickname" label="姓名" width="150" />
      <el-table-column label="角色" width="200">
        <template #default="{ row }">
          <el-tag v-for="role in row.roles" :key="role.id" size="small">{{ role.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="editUser(row)">编辑</el-button>
          <el-button size="small" @click="editRoles(row)">角色</el-button>
          <el-button size="small" type="danger" @click="deleteUser(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      @current-change="fetchUsers"
      @size-change="fetchUsers"
      layout="total, sizes, prev, pager, next, jumper"
    />

    <!-- 新增/编辑用户对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="手机号">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="密码" v-if="!form.id">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUser">确定</el-button>
      </template>
    </el-dialog>

    <!-- 角色编辑对话框 -->
    <el-dialog v-model="roleDialogVisible" title="分配角色" width="400px">
      <el-checkbox-group v-model="selectedRoles">
        <el-checkbox v-for="role in allRoles" :key="role.id" :label="role.id">{{ role.name }}</el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRoles">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from '@/api/index'

const users = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const roleDialogVisible = ref(false)
const dialogTitle = ref('')
const form = reactive({
  id: null as number | null,
  phone: '',
  nickname: '',
  password: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const allRoles = ref([])
const selectedRoles = ref<number[]>([])
const currentUserId = ref<number | null>(null)

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await axios.get('/system/users', {
      params: { page: pagination.page, page_size: pagination.pageSize }
    })
    users.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  }
  loading.value = false
}

const fetchRoles = async () => {
  try {
    const res = await axios.get('/system/roles')
    allRoles.value = res.data || []
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  }
}

const showCreateDialog = () => {
  dialogTitle.value = '新增用户'
  Object.assign(form, { id: null, phone: '', nickname: '', password: '' })
  dialogVisible.value = true
}

const editUser = (row: any) => {
  dialogTitle.value = '编辑用户'
  Object.assign(form, { id: row.id, phone: row.phone, nickname: row.nickname, password: '' })
  dialogVisible.value = true
}

const saveUser = async () => {
  try {
    if (form.id) {
      await axios.put(`/system/users/${form.id}`, form)
      ElMessage.success('更新成功')
    } else {
      await axios.post('/system/users', form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchUsers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const editRoles = async (row: any) => {
  currentUserId.value = row.id
  selectedRoles.value = row.roles?.map((r: any) => r.id) || []
  roleDialogVisible.value = true
}

const saveRoles = async () => {
  try {
    await axios.put(`/system/users/${currentUserId.value}/roles`, { role_ids: selectedRoles.value })
    ElMessage.success('角色分配成功')
    roleDialogVisible.value = false
    fetchUsers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const deleteUser = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除此用户吗？', '提示', {
      type: 'warning'
    })
    await axios.delete(`/system/users/${id}`)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

onMounted(() => {
  fetchUsers()
  fetchRoles()
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
