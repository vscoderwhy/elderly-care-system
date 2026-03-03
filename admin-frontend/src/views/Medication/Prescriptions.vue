<template>
  <div class="medication-management">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><FirstAidKit /></el-icon>
          用药管理
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAddPrescription">
          <el-icon><Plus /></el-icon>
          开具处方
        </el-button>
        <el-button @click="handleMedicationPlan">
          <el-icon><Tickets /></el-icon>
          用药计划
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in medStats" :key="stat.key">
        <div class="med-stat" :class="`stat-${stat.type}`">
          <div class="stat-header">
            <span class="stat-value">{{ stat.value }}</span>
            <span class="stat-unit">{{ stat.unit }}</span>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-chart">
              <div class="chart-bar" :style="{ width: stat.percent + '%' }"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 今日用药 -->
    <el-card shadow="never" class="today-card">
      <template #header>
        <div class="card-header">
          <span>今日用药安排</span>
          <div class="header-actions">
            <el-date-picker v-model="todayDate" type="date" placeholder="选择日期" />
            <el-button size="small" @click="handleExportToday">导出今日</el-button>
          </div>
        </div>
      </template>

      <el-table :data="todayMedications" stripe>
        <el-table-column prop="elderlyName" label="老人姓名" width="120" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="medicationName" label="药品名称" width="150" />
        <el-table-column prop="dosage" label="剂量" width="100" />
        <el-table-column prop="frequency" label="频次" width="100" />
        <el-table-column prop="plannedTime" label="计划时间" width="150">
          <template #default="{ row }">
            <el-tag
              v-for="(time, index) in row.plannedTimes"
              :key="index"
              :type="getTimeStatus(time, row.status)"
              size="small"
            >
              {{ time }}
              <el-icon v-if="row.status === 'completed'" class="check-icon"><Check /></el-icon>
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
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleConfirm(row)" :disabled="row.status === 'completed'">
              确认服药
            </el-button>
            <el-button size="small" text @click="handleSkip(row)" :disabled="row.status === 'completed'">
              跳过
            </el-button>
            <el-dropdown trigger="click">
              <el-button size="small" text>
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleViewPrescription(row)">
                    <el-icon><Document /></el-icon>
                    查看处方
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleRecord(row)">
                    <el-icon><Edit /></el-icon>
                    记录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 处方列表 -->
    <el-card shadow="never" class="prescriptions-card">
      <template #header>
        <div class="card-header">
          <span>处方列表</span>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索老人姓名或药品"
            style="width: 200px"
            clearable
            @clear="loadPrescriptions"
            @keyup.enter="loadPrescriptions"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>

      <el-table :data="prescriptions" stripe v-loading="loading">
        <el-table-column prop="prescriptionNo" label="处方编号" width="150" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="doctorName" label="开方医生" width="100" />
        <el-table-column prop="prescriptionDate" label="开方日期" width="120" />
        <el-table-column prop="medications" label="药品" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag
              v-for="(med, index) in row.medicationList.slice(0, 3)"
              :key="index"
              size="small"
              class="med-tag"
            >
              {{ med.name }}
            </el-tag>
            <span v-if="row.medicationList.length > 3" class="more-tag">
              等{{ row.medicationList.length }}种
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getPrescriptionStatusType(row.status)" size="small">
              {{ getPrescriptionStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleViewPrescription(row)">查看</el-button>
            <el-button size="small" text @click="handlePrintPrescription(row)">打印</el-button>
            <el-button size="small" text type="danger" @click="handleDeletePrescription(row)" v-if="row.status === 'draft'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </el-card>

    <!-- 用药确认对话框 -->
    <el-dialog v-model="confirmDialogVisible" title="确认服药" width="500px">
      <el-form :model="confirmForm" label-width="100px">
        <el-form-item label="老人姓名">
          <el-input :value="currentMedication?.elderlyName" disabled />
        </el-form-item>
        <el-form-item label="药品名称">
          <el-input :value="currentMedication?.medicationName" disabled />
        </el-form-item>
        <el-form-item label="计划时间">
          <el-input :value="currentTime" disabled />
        </el-form-item>
        <el-form-item label="实际用量">
          <el-input-number v-model="confirmForm.actualDose" :min="0" :max="100" />
          <span class="unit-text">片/粒</span>
        </el-form-item>
        <el-form-item label="服用情况">
          <el-radio-group v-model="confirmForm.result">
            <el-radio label="normal">正常服用</el-radio>
            <el-radio label="refused">拒服</el-radio>
            <el-radio label="vomited">服药后呕吐</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="confirmForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="confirmDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmSubmit" :loading="confirming">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  FirstAidKit,
  Plus,
  Tickets,
  Search,
  Check,
  ArrowDown,
  Document,
  Edit
} from '@element-plus/icons-vue'

// 统计数据
const medStats = ref([
  { key: 'today', label: '今日用药', value: 145, unit: '次', percent: 85, type: 'primary' },
  { key: 'completed', label: '已完成', value: 123, unit: '次', percent: 90, type: 'success' },
  { key: 'pending', label: '待服用', value: 22, unit: '次', percent: 15, type: 'warning' },
  { key: 'skipped', label: '已跳过', value: 8, unit: '次', percent: 5, type: 'danger' }
])

const todayDate = ref(new Date())
const searchKeyword = ref('')
const loading = ref(false)

const todayMedications = ref([
  {
    id: 1,
    elderlyName: '张奶奶',
    bedNumber: '3号楼201',
    medicationName: '降压药',
    dosage: '1片',
    frequency: '一日3次',
    plannedTimes: ['08:00', '12:00', '18:00'],
    status: 'completed'
  },
  {
    id: 2,
    elderlyName: '王爷爷',
    bedNumber: '3号楼202',
    medicationName: '降糖药',
    dosage: '2片',
    frequency: '一日2次',
    plannedTimes: ['08:00', '20:00'],
    status: 'pending'
  },
  {
    id: 3,
    elderlyName: '李奶奶',
    bedNumber: '2号楼105',
    medicationName: '心脏病药',
    dosage: '1片',
    frequency: '一日1次',
    plannedTimes: ['08:00'],
    status: 'skipped'
  },
  {
    id: 4,
    elderlyName: '刘爷爷',
    bedNumber: '2号楼108',
    medicationName: '钙片',
    dosage: '1片',
    frequency: '一日3次',
    plannedTimes: ['08:00', '13:00', '20:00'],
    status: 'pending'
  }
])

const prescriptions = ref([
  {
    id: 1,
    prescriptionNo: 'RX2026030001',
    elderlyName: '张奶奶',
    bedNumber: '3号楼201',
    doctorName: '王医生',
    prescriptionDate: '2026-03-01',
    medicationList: [
      { name: '降压药', dosage: '1片', frequency: '一日3次', days: 30 },
      { name: '钙片', dosage: '1片', frequency: '一日1次', days: 30 },
      { name: '维生素', dosage: '1片', frequency: '一日2次', days: 30 }
    ],
    status: 'active'
  },
  {
    id: 2,
    prescriptionNo: 'RX2026030002',
    elderlyName: '王爷爷',
    bedNumber: '3号楼202',
    doctorName: '李医生',
    prescriptionDate: '2026-02-28',
    medicationList: [
      { name: '降糖药', dosage: '2片', frequency: '一日2次', days: 30 },
      { name: '阿司匹林', dosage: '1片', frequency: '一日1次', days: 90 }
    ],
    status: 'active'
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(15)

// 确认对话框
const confirmDialogVisible = ref(false)
const confirming = ref(false)
const currentMedication = ref<any>(null)
const currentTime = ref('08:00')

const confirmForm = reactive({
  actualDose: 1,
  result: 'normal',
  remark: ''
})

const getTimeStatus = (time: string, status: string) => {
  if (status === 'completed') return 'success'
  return 'info'
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    completed: 'success',
    pending: 'warning',
    skipped: 'danger'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    completed: '已完成',
    pending: '待服用',
    skipped: '已跳过'
  }
  return map[status] || status
}

const getPrescriptionStatusType = (status: string) => {
  if (status === 'active') return 'success'
  if (status === 'expired') return 'danger'
  return 'info'
}

const getPrescriptionStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '有效',
    draft: '草稿',
    expired: '已过期',
    completed: '已完成'
  }
  return map[status] || status
}

const handleAddPrescription = () => {
  console.log('开具处方')
}

const handleMedicationPlan = () => {
  console.log('用药计划')
}

const handleExportToday = () => {
  console.log('导出今日用药')
}

const handleConfirm = (row: any) => {
  currentMedication.value = row
  const now = new Date()
  currentTime.value = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
  confirmDialogVisible.value = true
}

const handleSkip = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要跳过"${row.elderlyName} - ${row.medicationName}"吗？`, '提示', {
      type: 'warning'
    })
    ElMessage.success('已标记为跳过')
  } catch {
    // 取消
  }
}

const handleConfirmSubmit = () => {
  confirming.value = true
  setTimeout(() => {
    confirming.value = false
    confirmDialogVisible.value = false
    ElMessage.success('确认成功')
    // 更新数据
  }, 1000)
}

const handleViewPrescription = (row: any) => {
  console.log('查看处方', row)
}

const handleRecord = (row: any) => {
  console.log('记录', row)
}

const handlePrintPrescription = (row: any) => {
  console.log('打印处方', row)
}

const handleDeletePrescription = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这条处方吗？', '提示', {
      type: 'warning'
    })
    ElMessage.success('删除成功')
  } catch {
    // 取消
  }
}

const loadPrescriptions = () => {
  // TODO: 加载处方列表
}

onMounted(() => {
  loadPrescriptions()
})
</script>

<style scoped lang="scss">
.medication-management {
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
    margin: 0;
  }

  .stats-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .med-stat {
    padding: 20px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);
  }

  .stat-header {
    display: flex;
    align-items: baseline;
    gap: 8px;
    margin-bottom: 12px;
  }

  .stat-value {
    font-size: 32px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .stat-unit {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .stat-body {
    position: relative;
  }

  .stat-label {
    font-size: 14px;
    color: var(--text-secondary);
    margin-bottom: 8px;
  }

  .stat-chart {
    width: 100%;
    height: 8px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    overflow: hidden;
  }

  .chart-bar {
    height: 100%;
    background: var(--primary-color);
    border-radius: 4px;
    transition: width 0.5s ease;
  }

  &.stat-primary .chart-bar { background: var(--gradient-blue); }
  &.stat-success .chart-bar { background: var(--gradient-green); }
  &.stat-warning .chart-bar { background: var(--gradient-orange); }
  &.stat-danger .chart-bar { background: var(--gradient-red); }

  .today-card,
  .prescriptions-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .med-tag {
    margin-right: 4px;
  }

  .more-tag {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .check-icon {
    margin-left: 4px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

@media (max-width: 768px) {
  .medication-management {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .stats-row {
      :deep(.el-col) {
        margin-bottom: 8px;
      }
    }

    .med-stat {
      padding: 16px;

      .stat-value {
        font-size: 24px;
      }
    }

    .header-actions {
      flex-wrap: wrap;

      * {
        flex-shrink: 0;
      }
    }
  }
}
</style>
