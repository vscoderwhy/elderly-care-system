<template>
  <div class="staff-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>员工管理</span>
          <div class="header-actions">
            <el-select v-model="roleFilter" placeholder="角色筛选" clearable @change="fetchData" style="width: 150px; margin-right: 10px;">
              <el-option label="护理员" value="caregiver" />
              <el-option label="医生" value="doctor" />
              <el-option label="管理员" value="admin" />
            </el-select>
            <el-button type="primary" @click="showAddDialog">添加员工</el-button>
          </div>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="nickname" label="姓名" width="120" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag v-for="role in row.roles" :key="role.id" size="small" style="margin-right: 4px;">
              {{ getRoleName(role.name) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="入职时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showEditDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除该员工?" @confirm="deleteStaff(row.id)">
              <template #reference>
                <el-button type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑员工' : '添加员工'" width="500px">
      <el-form :model="form" label-width="100px" :rules="rules" ref="formRef">
        <el-form-item label="姓名" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role_id">
          <el-select v-model="form.role_id" placeholder="请选择角色" style="width: 100%">
            <el-option label="护理员" :value="2" />
            <el-option label="医生" :value="3" />
            <el-option label="管理员" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" v-if="isEdit">
          <el-radio-group v-model="form.status">
            <el-radio label="active">在职</el-radio>
            <el-radio label="inactive">离职</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import instance from '@/api'

const loading = ref(false)
const tableData = ref([])
const roleFilter = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editId = ref(0)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const form = reactive({
  nickname: '',
  phone: '',
  password: '',
  role_id: null as number | null,
  status: 'active'
})

const rules: FormRules = {
  nickname: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role_id: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const getRoleName = (name: string) => {
  const map: Record<string, string> = {
    admin: '管理员',
    caregiver: '护理员',
    doctor: '医生',
    family: '家属'
  }
  return map[name] || name
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (roleFilter.value) {
      params.role = roleFilter.value
    }
    const result = await instance.get('/staff', { params })
    tableData.value = result.list || []
    pagination.total = result.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.nickname = ''
  form.phone = ''
  form.password = ''
  form.role_id = null
  form.status = 'active'
}

const showAddDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const showEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  form.nickname = row.nickname
  form.phone = row.phone
  form.status = row.status
  form.role_id = row.roles?.[0]?.id || null
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await instance.put(`/staff/${editId.value}`, {
          nickname: form.nickname,
          status: form.status
        })
        ElMessage.success('更新成功')
      } else {
        await instance.post('/staff', form)
        ElMessage.success('添加成功')
      }
      dialogVisible.value = false
      fetchData()
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const deleteStaff = async (id: number) => {
  try {
    await instance.delete(`/staff/${id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.staff-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}
</style>
