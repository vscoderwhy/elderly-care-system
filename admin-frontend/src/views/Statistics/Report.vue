<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>月度综合报表</span>
        <el-date-picker
          v-model="reportMonth"
          type="month"
          placeholder="选择月份"
          format="YYYY年MM月"
          value-format="YYYY-MM"
          @change="fetchMonthlyReport"
        />
      </div>
    </template>

    <div v-loading="loading">
      <el-descriptions :column="3" border>
        <el-descriptions-item label="统计月份">{{ report.period?.year }}年{{ report.period?.month }}月</el-descriptions-item>
        <el-descriptions-item label="统计时间">
          {{ report.period?.start }} 至 {{ report.period?.end }}
        </el-descriptions-item>
        <el-descriptions-item label="报表生成时间">{{ new Date().toLocaleString() }}</el-descriptions-item>
      </el-descriptions>

      <el-divider content-position="left">老人统计</el-divider>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic title="在院老人" :value="report.elderly?.total || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="新入住" :value="report.elderly?.new_admitted || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="出院" :value="report.elderly?.discharged || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="去世" :value="report.elderly?.deceased || 0" />
        </el-col>
      </el-row>

      <el-divider content-position="left">护理统计</el-divider>
      <el-row :gutter="20">
        <el-col :span="8">
          <el-statistic title="护理记录总数" :value="report.care?.total_records || 0" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="完成率" :value="report.care?.completion_rate || 0" suffix="%" />
        </el-col>
      </el-row>

      <el-divider content-position="left">财务统计</el-divider>
      <el-row :gutter="20">
        <el-col :span="8">
          <el-statistic title="本月收入" :value="report.finance?.income || 0" :precision="2" prefix="¥" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="本月支出" :value="report.finance?.expense || 0" :precision="2" prefix="¥" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="待收金额" :value="report.finance?.pending_amount || 0" :precision="2" prefix="¥" />
        </el-col>
      </el-row>

      <el-divider content-position="left">健康异常统计</el-divider>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-statistic title="异常记录数" :value="report.health?.abnormal_count || 0" />
        </el-col>
        <el-col :span="12">
          <el-statistic title="紧急事件数" :value="report.health?.emergency_count || 0" />
        </el-col>
      </el-row>

      <el-divider />
      <el-button type="primary" @click="exportReport">导出报表</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import axios from '@/api/index'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const reportMonth = ref(new Date().toISOString().slice(0, 7))
const report = reactive({
  period: null,
  elderly: null,
  care: null,
  finance: null,
  health: null
})

const fetchMonthlyReport = async () => {
  loading.value = true
  try {
    const [year, month] = reportMonth.value.split('-')
    const res = await axios.get('/statistics/monthly-report', {
      params: { year: parseInt(year), month: parseInt(month) }
    })
    Object.assign(report, res.data)
  } catch (error) {
    console.error('Failed to fetch monthly report:', error)
    ElMessage.error('获取报表失败')
  } finally {
    loading.value = false
  }
}

const exportReport = () => {
  ElMessage.info('导出功能开发中...')
}

// 初始加载
fetchMonthlyReport()
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-statistic__head) {
  font-size: 14px;
  color: #909399;
}

:deep(.el-statistic__content) {
  font-size: 24px;
  font-weight: bold;
}
</style>
