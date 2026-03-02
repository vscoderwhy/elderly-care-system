<template>
  <div class="service-requests-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>服务呼叫管理</span>
          <el-radio-group v-model="statusFilter" size="small" @change="loadData">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="pending">待处理</el-radio-button>
            <el-radio-button label="processing">处理中</el-radio-button>
            <el-radio-button label="completed">已完成</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="requests" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column label="老人" width="100">
          <template #default="{ row }">
            {{ row.elderly?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="notes" label="备注" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="请求时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.requested_at) }}
          </template>
        </el-table-column>
        <el-table-column label="处理人" width="100">
          <template #default="{ row }">
            {{ row.handler?.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 'pending'">
              <el-button type="primary" size="small" @click="handleRequest(row.id, 'processing')">
                接单
              </el-button>
            </template>
            <template v-else-if="row.status === 'processing'">
              <el-button type="success" size="small" @click="handleRequest(row.id, 'completed')">
                完成
              </el-button>
            </template>
            <span v-else class="text-muted">已处理</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import instance, { careApi } from '@/api'

const requests = ref<any[]>([])
const loading = ref(false)
const statusFilter = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)

onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await instance.get('/service/requests', {
      params: {
        status: statusFilter.value,
        page: page.value,
        page_size: pageSize
      }
    })
    requests.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('加载失败', error)
  } finally {
    loading.value = false
  }
}

const handleRequest = async (id: number, status: string) => {
  try {
    await instance.put(`/service/requests/${id}`, { status })
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    '紧急': 'danger',
    '护理': 'primary',
    '医疗': 'warning',
    '送餐': 'success'
  }
  return map[type] || 'info'
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}
</script>

<style scoped>
.service-requests-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text-muted {
  color: #999;
  font-size: 12px;
}
</style>
