<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>数据导出</span>
      </div>
    </template>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="export-card" shadow="hover" @click="exportElderly">
          <div class="export-icon">👴</div>
          <div class="export-title">老人列表</div>
          <div class="export-desc">导出所有老人信息</div>
          <el-button type="primary" plain>导出CSV</el-button>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="export-card" shadow="hover" @click="exportCareRecords">
          <div class="export-icon">📝</div>
          <div class="export-title">护理记录</div>
          <div class="export-desc">导出护理记录数据</div>
          <el-button type="primary" plain>导出CSV</el-button>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="export-card" shadow="hover" @click="exportHealthData">
          <div class="export-icon">❤️</div>
          <div class="export-title">健康数据</div>
          <div class="export-desc">导出健康记录数据</div>
          <el-button type="primary" plain>导出CSV</el-button>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="export-card" shadow="hover" @click="exportFinance">
          <div class="export-icon">💰</div>
          <div class="export-title">财务报表</div>
          <div class="export-desc">导出账单和支付数据</div>
          <el-button type="primary" plain>导出CSV</el-button>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="export-card" shadow="hover" @click="exportMedicationRecords">
          <div class="export-icon">💊</div>
          <div class="export-title">用药记录</div>
          <div class="export-desc">导出用药记录数据</div>
          <el-button type="primary" plain>导出CSV</el-button>
        </el-card>
      </el-col>
    </el-row>

    <!-- 导出对话框 -->
    <el-dialog v-model="dialogVisible" title="选择导出时间范围" width="500px">
      <el-form label-width="100px">
        <el-form-item label="开始日期">
          <el-date-picker v-model="dateRange.start" type="date" placeholder="选择开始日期" />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker v-model="dateRange.end" type="date" placeholder="选择结束日期" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="doExport">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/api/index'

const dialogVisible = ref(false)
const currentExport = ref('')
const dateRange = reactive({
  start: '',
  end: ''
})

const exportElderly = () => {
  currentExport.value = 'elderly'
  dialogVisible.value = true
}

const exportCareRecords = () => {
  currentExport.value = 'care'
  dialogVisible.value = true
}

const exportHealthData = () => {
  currentExport.value = 'health'
  dialogVisible.value = true
}

const exportFinance = () => {
  currentExport.value = 'finance'
  dialogVisible.value = true
}

const exportMedicationRecords = () => {
  currentExport.value = 'medication'
  dialogVisible.value = true
}

const doExport = async () => {
  const exportUrls: Record<string, string> = {
    elderly: '/export/elderly',
    care: '/export/care-records',
    health: '/export/health-data',
    finance: '/export/finance',
    medication: '/export/medication-records'
  }

  const url = exportUrls[currentExport.value]
  if (!url) {
    ElMessage.error('未知的导出类型')
    return
  }

  try {
    const params: any = {}
    if (dateRange.start) params.start_date = dateRange.start
    if (dateRange.end) params.end_date = dateRange.end

    const response = await axios.get(url, {
      params,
      responseType: 'blob',
      timeout: 30000
    })

    const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = window.URL.createObjectURL(blob)

    const filenames: Record<string, string> = {
      elderly: '老人列表',
      care: '护理记录',
      health: '健康数据',
      finance: '财务报表',
      medication: '用药记录'
    }

    const dateStr = dateRange.start && dateRange.end
      ? `_${dateRange.start}_${dateRange.end}`
      : `_${new Date().toISOString().split('T')[0]}`

    link.download = `${filenames[currentExport.value]}${dateStr}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(link.href)

    dialogVisible.value = false
    ElMessage.success('导出成功')
  } catch (error: any) {
    console.error('导出失败:', error)
    const errorMsg = error.response?.data?.message || error.message || '导出失败，请稍后重试'
    ElMessage.error(errorMsg)
  }
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.export-card {
  text-align: center;
  padding: 20px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.export-card:hover {
  transform: translateY(-5px);
}

.export-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.export-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 5px;
}

.export-desc {
  color: #999;
  font-size: 14px;
  margin-bottom: 15px;
}
</style>
