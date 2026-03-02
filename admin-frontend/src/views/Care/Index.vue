<template>
  <div class="care-records">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>护理记录</span>
          <div class="header-actions">
            <el-select v-model="filterElderlyId" placeholder="选择老人筛选" clearable @change="fetchData" style="width: 200px; margin-right: 10px;">
              <el-option v-for="e in elderlyList" :key="e.id" :label="e.name" :value="e.id" />
            </el-select>
            <el-button type="primary" @click="showQuickRecord">快速记录</el-button>
          </div>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="elderly.name" label="老人姓名" width="120" />
        <el-table-column prop="care_item.name" label="护理项目" width="150" />
        <el-table-column prop="staff.nickname" label="护理员" width="120" />
        <el-table-column prop="notes" label="备注" />
        <el-table-column prop="recorded_at" label="记录时间" width="180" />
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

    <!-- 快速记录对话框 -->
    <el-dialog v-model="showRecordDialog" title="快速护理记录" width="500px">
      <el-form :model="recordForm" label-width="100px">
        <el-form-item label="选择老人" required>
          <el-select v-model="recordForm.elderly_id" placeholder="请选择老人" style="width: 100%">
            <el-option v-for="e in elderlyList" :key="e.id" :label="e.name" :value="e.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="护理项目" required>
          <el-select v-model="recordForm.care_item_id" placeholder="请选择护理项目" style="width: 100%">
            <el-option v-for="item in careItems" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="recordForm.notes" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRecordDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRecord" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import instance, { careApi, elderlyApi } from '@/api'

const loading = ref(false)
const tableData = ref([])
const elderlyList = ref<any[]>([])
const careItems = ref<any[]>([])
const filterElderlyId = ref<number | null>(null)
const showRecordDialog = ref(false)
const submitting = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const recordForm = reactive({
  elderly_id: null as number | null,
  care_item_id: null as number | null,
  notes: ''
})

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filterElderlyId.value) {
      params.elderly_id = filterElderlyId.value
    }
    const result = await careApi.listRecords(params)
    tableData.value = result.list || []
    pagination.total = result.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadElderlyList = async () => {
  try {
    const result = await elderlyApi.list({ page: 1, page_size: 1000 })
    elderlyList.value = result.list || []
  } catch (error) {
    console.error(error)
  }
}

const loadCareItems = async () => {
  try {
    const result = await instance.get('/care/items')
    careItems.value = result || []
  } catch (error) {
    console.error(error)
  }
}

const showQuickRecord = () => {
  recordForm.elderly_id = null
  recordForm.care_item_id = null
  recordForm.notes = ''
  showRecordDialog.value = true
}

const submitRecord = async () => {
  if (!recordForm.elderly_id || !recordForm.care_item_id) {
    ElMessage.warning('请选择老人和护理项目')
    return
  }

  submitting.value = true
  try {
    await instance.post('/care/records', recordForm)
    ElMessage.success('记录成功')
    showRecordDialog.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('记录失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchData()
  loadElderlyList()
  loadCareItems()
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
  align-items: center;
}
</style>
