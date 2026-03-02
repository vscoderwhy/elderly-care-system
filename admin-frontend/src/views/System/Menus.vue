<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>菜单管理</span>
        <el-button type="primary" @click="showCreateDialog">新增菜单</el-button>
      </div>
    </template>

    <el-table
      :data="menuTree"
      style="width: 100%"
      row-key="id"
      :tree-props="{ children: 'children' }"
      v-loading="loading"
    >
      <el-table-column prop="name" label="菜单名称" width="200" />
      <el-table-column prop="path" label="路由路径" width="200" />
      <el-table-column prop="component" label="组件路径" />
      <el-table-column prop="icon" label="图标" width="120" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.type === 'menu' ? 'primary' : 'info'" size="small">
            {{ row.type === 'menu' ? '菜单' : row.type === 'directory' ? '目录' : '按钮' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column prop="visible" label="可见" width="80">
        <template #default="{ row }">
          <el-tag :type="row.visible ? 'success' : 'info'" size="small">
            {{ row.visible ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="editMenu(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteMenu(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑菜单对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parent_id"
            :data="menuTree"
            :props="{ label: 'name', value: 'id' }"
            placeholder="选择上级菜单（可不选）"
            check-strictly
            clearable
          />
        </el-form-item>
        <el-form-item label="菜单名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="路由路径">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item label="组件路径">
          <el-input v-model="form.component" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" />
        </el-form-item>
        <el-form-item label="菜单类型">
          <el-select v-model="form.type" placeholder="选择类型">
            <el-option label="菜单" value="menu" />
            <el-option label="目录" value="directory" />
            <el-option label="按钮" value="button" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限标识">
          <el-input v-model="form.permission" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="是否可见">
          <el-switch v-model="form.visible" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveMenu">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from '@/api/index'

interface Menu {
  id: number
  parent_id?: number
  name: string
  path: string
  component: string
  icon: string
  sort: number
  type: string
  permission: string
  visible: boolean
  children?: Menu[]
}

const menuTree = ref<Menu[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const form = reactive<{
  id: number | null
  parent_id: number | null
  name: string
  path: string
  component: string
  icon: string
  sort: number
  type: string
  permission: string
  visible: boolean
}>({
  id: null,
  parent_id: null,
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  type: 'menu',
  permission: '',
  visible: true
})

const fetchMenus = async () => {
  loading.value = true
  try {
    const res = await axios.get('/system/menus')
    menuTree.value = res.data || []
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
  }
  loading.value = false
}

const showCreateDialog = () => {
  dialogTitle.value = '新增菜单'
  Object.assign(form, {
    id: null,
    parent_id: null,
    name: '',
    path: '',
    component: '',
    icon: '',
    sort: 0,
    type: 'menu',
    permission: '',
    visible: true
  })
  dialogVisible.value = true
}

const editMenu = (row: Menu) => {
  dialogTitle.value = '编辑菜单'
  Object.assign(form, {
    id: row.id,
    parent_id: row.parent_id || null,
    name: row.name,
    path: row.path,
    component: row.component,
    icon: row.icon,
    sort: row.sort,
    type: row.type,
    permission: row.permission,
    visible: row.visible
  })
  dialogVisible.value = true
}

const saveMenu = async () => {
  try {
    const data = { ...form }
    if (form.id) {
      await axios.put(`/system/menus/${form.id}`, data)
      ElMessage.success('更新成功')
    } else {
      await axios.post('/system/menus', data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchMenus()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const deleteMenu = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除此菜单吗？子菜单也将被删除。', '提示', {
      type: 'warning'
    })
    await axios.delete(`/system/menus/${id}`)
    ElMessage.success('删除成功')
    fetchMenus()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

onMounted(() => {
  fetchMenus()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
