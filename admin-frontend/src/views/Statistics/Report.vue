<template>
  <div class="report-center">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Document /></el-icon>
          报表中心
        </h2>
      </div>
      <div class="header-actions">
        <el-button @click="handleCreateReport">
          <el-icon><Plus /></el-icon>
          新建报表
        </el-button>
      </div>
    </div>

    <!-- 报表分类 -->
    <el-tabs v-model="activeCategory" class="report-tabs">
      <el-tab-pane label="全部报表" name="all">
        <ReportList :reports="allReports" @view="handleViewReport" @export="handleExportReport" />
      </el-tab-pane>
      <el-tab-pane label="运营报表" name="operation">
        <ReportList :reports="operationReports" @view="handleViewReport" @export="handleExportReport" />
      </el-tab-pane>
      <el-tab-pane label="护理报表" name="nursing">
        <ReportList :reports="nursingReports" @view="handleViewReport" @export="handleExportReport" />
      </el-tab-pane>
      <el-tab-pane label="财务报表" name="finance">
        <ReportList :reports="financeReports" @view="handleViewReport" @export="handleExportReport" />
      </el-tab-pane>
      <el-tab-pane label="健康报表" name="health">
        <ReportList :reports="healthReports" @view="handleViewReport" @export="handleExportReport" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Document, Plus } from '@element-plus/icons-vue'
import ReportList from '@/components/Report/ReportList.vue'

interface Report {
  id: string | number
  name: string
  category: string
  description: string
  updateAt: string
  schedule?: string
}

const activeCategory = ref('all')

// 报表数据
const reports: Report[] = [
  {
    id: 1,
    name: '月度运营报表',
    category: 'operation',
    description: '包含入住率、床位周转、营收等核心运营指标',
    updateAt: '2026-03-01',
    schedule: '每月1日'
  },
  {
    id: 2,
    name: '老人统计报表',
    category: 'operation',
    description: '按年龄、护理等级、入住时长等维度统计分析',
    updateAt: '2026-03-01',
    schedule: '每月1日'
  },
  {
    id: 3,
    name: '护理质量报表',
    category: 'nursing',
    description: '护理任务完成率、满意度评分、质量评估等',
    updateAt: '2026-03-01',
    schedule: '每周一'
  },
  {
    id: 4,
    name: '护理员工作量报表',
    category: 'nursing',
    description: '各护理员工作时长、任务数量、绩效统计',
    updateAt: '2026-03-01',
    schedule: '每周一'
  },
  {
    id: 5,
    name: '财务收支报表',
    category: 'finance',
    description: '月度收入、支出、利润统计及趋势分析',
    updateAt: '2026-03-01',
    schedule: '每月1日'
  },
  {
    id: 6,
    name: '费用欠款报表',
    category: 'finance',
    description: '欠款明细、账龄分析、催收记录',
    updateAt: '2026-03-01',
    schedule: '每日'
  },
  {
    id: 7,
    name: '健康指标报表',
    category: 'health',
    description: '血压、血糖等健康指标的统计分析',
    updateAt: '2026-03-01',
    schedule: '每周一'
  },
  {
    id: 8,
    name: '疾病统计报表',
    category: 'health',
    description: '按疾病类型统计患病率、趋势分析',
    updateAt: '2026-03-01',
    schedule: '每月1日'
  }
]

// 分类报表
const allReports = computed(() => reports)

const operationReports = computed(() =>
  reports.filter(r => r.category === 'operation')
)

const nursingReports = computed(() =>
  reports.filter(r => r.category === 'nursing')
)

const financeReports = computed(() =>
  reports.filter(r => r.category === 'finance')
)

const healthReports = computed(() =>
  reports.filter(r => r.category === 'health')
)

// 事件处理
const handleCreateReport = () => {
  console.log('创建报表')
}

const handleViewReport = (report: Report) => {
  console.log('查看报表', report)
}

const handleExportReport = (report: Report) => {
  console.log('导出报表', report)
}
</script>

<style scoped lang="scss">
.report-center {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .page-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .report-tabs {
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    padding: 20px;
    box-shadow: var(--card-shadow);

    :deep(.el-tabs__content) {
      padding: 0;
    }
  }
}

@media (max-width: 768px) {
  .report-center {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }
  }
}
</style>
