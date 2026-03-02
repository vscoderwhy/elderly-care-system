<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>智能预警</span>
        <el-button type="primary" @click="checkAlerts">检查预警</el-button>
      </div>
    </template>

    <!-- 预警统计 -->
    <el-row :gutter="20" class="summary">
      <el-col :span="6">
        <el-card class="summary-card critical">
          <div class="number">{{ summary.critical || 0 }}</div>
          <div class="label">严重</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="summary-card warning">
          <div class="number">{{ summary.warning || 0 }}</div>
          <div class="label">警告</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="summary-card info">
          <div class="number">{{ summary.info || 0 }}</div>
          <div class="label">提示</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="summary-card total">
          <div class="number">{{ summary.total_active || 0 }}</div>
          <div class="label">总计</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 预警列表 -->
    <el-table :data="alerts" style="width: 100%; margin-top: 20px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="type" label="类型" width="150">
        <template #default="{ row }">
          {{ getAlertTypeText(row.type) }}
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" width="200" />
      <el-table-column prop="content" label="内容" show-overflow-tooltip />
      <el-table-column prop="level" label="级别" width="100">
        <template #default="{ row }">
          <el-tag :type="getLevelType(row.level)" size="small">
            {{ getLevelText(row.level) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button v-if="row.status === 'active'" size="small" @click="acknowledgeAlert(row.id)">确认</el-button>
          <el-button v-if="row.status === 'active' || row.status === 'acknowledged'" size="small" @click="resolveAlert(row.id)">解决</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      @current-change="fetchAlerts"
      layout="total, prev, pager, next"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/api/index'

const alerts = ref([])
const loading = ref(false)
const summary = reactive({
  critical: 0,
  warning: 0,
  info: 0,
  total_active: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const fetchAlerts = async () => {
  loading.value = true
  try {
    const res = await axios.get('/alerts', {
      params: { page: pagination.page, page_size: pagination.pageSize, status: 'active' }
    })
    alerts.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('获取预警列表失败')
  }
  loading.value = false
}

const fetchSummary = async () => {
  try {
    const res = await axios.get('/alerts/summary')
    Object.assign(summary, res.data)
  } catch (error) {
    ElMessage.error('获取预警统计失败')
  }
}

const checkAlerts = async () => {
  try {
    await axios.post('/alerts/check')
    ElMessage.success('预警检查完成')
    fetchAlerts()
    fetchSummary()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '检查失败')
  }
}

const acknowledgeAlert = async (id: number) => {
  try {
    await axios.put(`/alerts/${id}/acknowledge`)
    ElMessage.success('已确认')
    fetchAlerts()
    fetchSummary()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const resolveAlert = async (id: number) => {
  try {
    await axios.put(`/alerts/${id}/resolve`)
    ElMessage.success('已标记为已解决')
    fetchAlerts()
    fetchSummary()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const getAlertTypeText = (type: string) => {
  const map: Record<string, string> = {
    medication_low: '药品库存不足',
    medication_expiry: '药品即将过期',
    health_abnormal: '健康数据异常',
    bill_overdue: '账单逾期',
    bed_available: '床位可用',
    bed_full: '床位紧张'
  }
  return map[type] || type
}

const getLevelType = (level: string) => {
  const map: Record<string, string> = {
    critical: 'danger',
    warning: 'warning',
    info: 'info'
  }
  return map[level] || 'info'
}

const getLevelText = (level: string) => {
  const map: Record<string, string> = {
    critical: '严重',
    warning: '警告',
    info: '提示'
  }
  return map[level] || level
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'danger',
    acknowledged: 'warning',
    resolved: 'success'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '活跃',
    acknowledged: '已确认',
    resolved: '已解决'
  }
  return map[status] || status
}

const formatDateTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchAlerts()
  fetchSummary()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary {
  margin-bottom: 20px;
}

.summary-card {
  text-align: center;
}

.summary-card .number {
  font-size: 32px;
  font-weight: bold;
}

.summary-card .label {
  color: #999;
  margin-top: 10px;
}

.summary-card.critical { border-left: 4px solid #f56c6c; }
.summary-card.warning { border-left: 4px solid #e6a23c; }
.summary-card.info { border-left: 4px solid #909399; }
.summary-card.total { border-left: 4px solid #409eff; }

.el-pagination {
  margin-top: 20px;
}
</style>
