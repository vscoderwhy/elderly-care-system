<template>
  <div class="finance-bills">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Wallet /></el-icon>
          费用账单
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          生成账单
        </el-button>
        <ExportButton
          :data="billList"
          :columns="exportColumns"
          :total="total"
          filename="财务报表"
          :export-function="handleExport"
        />
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in billStats" :key="stat.key">
        <div class="bill-stat" :class="`stat-${stat.type}`">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value">¥{{ stat.value.toLocaleString() }}</div>
          <div class="stat-trend" :class="stat.trend >= 0 ? 'up' : 'down'">
            <el-icon><component :is="stat.trend >= 0 ? ArrowUp : ArrowDown" /></el-icon>
            {{ Math.abs(stat.trend) }}%
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 筛选表单 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="老人">
          <el-select v-model="filterForm.elderlyId" placeholder="请选择" clearable filterable>
            <el-option
              v-for="e in elderlyList"
              :key="e.id"
              :label="`${e.name} (${e.bedNumber})`"
              :value="e.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="账单类型">
          <el-select v-model="filterForm.billType" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="床位费" value="床位费" />
            <el-option label="护理费" value="护理费" />
            <el-option label="伙食费" value="伙食费" />
            <el-option label="医疗费" value="医疗费" />
            <el-option label="其他费用" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="支付状态">
          <el-select v-model="filterForm.status" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="未支付" value="unpaid" />
            <el-option label="已支付" value="paid" />
            <el-option label="已逾期" value="overdue" />
          </el-select>
        </el-form-item>
        <el-form-item label="账单月份">
          <el-date-picker
            v-model="filterForm.month"
            type="month"
            placeholder="选择月份"
            value-format="YYYY-MM"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 账单列表 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>账单明细</span>
          <div class="header-actions">
            <el-checkbox v-model="showUnpaidOnly">仅显示未支付</el-checkbox>
            <el-button size="small" @click="handleBatchPay" :disabled="!hasSelected">
              批量支付 ({{ selectedCount }})
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        ref="tableRef"
        v-loading="loading"
        :data="filteredBillList"
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="billNo" label="账单号" width="150" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="billType" label="费用类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.billType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="billDate" label="账单日期" width="120" />
        <el-table-column prop="dueDate" label="应付日期" width="120">
          <template #default="{ row }">
            <span :class="{ 'text-danger': isOverdue(row.dueDate, row.status) }">
              {{ row.dueDate }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button
              v-if="row.status === 'unpaid'"
              size="small"
              type="primary"
              text
              @click="handlePay(row)"
            >
              支付
            </el-button>
            <el-button size="small" text @click="handleReminder(row)" v-if="row.status === 'unpaid'">
              催缴
            </el-button>
            <el-dropdown trigger="click">
              <el-button size="small" text>
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handlePrint(row)">
                    <el-icon><Printer /></el-icon>
                    打印
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleInvoice(row)">
                    <el-icon><Ticket /></el-icon>
                    开票
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">
                    <el-icon><Delete /></el-icon>
                    作废
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 支付对话框 -->
    <el-dialog v-model="payDialogVisible" title="在线支付" width="500px">
      <div class="pay-dialog">
        <div class="pay-info">
          <div class="info-item">
            <span class="label">账单号：</span>
            <span class="value">{{ currentBill?.billNo }}</span>
          </div>
          <div class="info-item">
            <span class="label">应付金额：</span>
            <span class="value amount-value">¥{{ currentBill?.amount.toFixed(2) }}</span>
          </div>
        </div>

        <el-divider />

        <div class="pay-methods">
          <div class="section-title">选择支付方式</div>
          <el-radio-group v-model="payForm.paymentMethod" class="method-list">
            <el-radio label="wechat" class="method-item">
              <div class="method-content">
                <el-icon class="method-icon"><WechatFilled /></el-icon>
                <span>微信支付</span>
              </div>
            </el-radio>
            <el-radio label="alipay" class="method-item">
              <div class="method-content">
                <el-icon class="method-icon"><Alipay /></el-icon>
                <span>支付宝</span>
              </div>
            </el-radio>
            <el-radio label="cash" class="method-item">
              <div class="method-content">
                <el-icon class="method-icon"><Wallet /></el-icon>
                <span>现金/POS机</span>
              </div>
            </el-radio>
          </el-radio-group>
        </div>

        <el-form v-if="payForm.paymentMethod === 'cash'" :model="payForm" class="pay-form">
          <el-form-item label="支付方式">
            <el-select v-model="payForm.cashType">
              <el-option label="现金" value="cash" />
              <el-option label="POS机刷卡" value="pos" />
              <el-option label="银行转账" value="transfer" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmPay" :loading="paying">
          确认支付 ¥{{ currentBill?.amount.toFixed(2) }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Wallet,
  Plus,
  ArrowUp,
  ArrowDown,
  ArrowDown as ArrowDropDown,
  Printer,
  Ticket,
  Delete,
  WechatFilled,
  Alipay
} from '@element-plus/icons-vue'
import ExportButton from '@/components/Export/ExportButton.vue'
import { exportFinancialReport } from '@/utils/export'
import { getBills, getElderlyList, getStatistics } from '@/utils/seedData'

// 统计数据
const billStats = ref([
  { key: 'total', label: '本月应收', value: 89500, trend: 5.2, type: 'primary' },
  { key: 'paid', label: '本月已收', value: 76800, trend: 8.1, type: 'success' },
  { key: 'unpaid', label: '待收款', value: 12700, trend: -12.5, type: 'warning' },
  { key: 'overdue', label: '逾期金额', value: 3200, trend: -25.0, type: 'danger' }
])

const loading = ref(false)
const tableRef = ref()
const billList = ref([])
const selectedRows = ref([])
const total = ref(0)
const showUnpaidOnly = ref(false)

const filterForm = reactive({
  elderlyId: '',
  billType: '',
  status: '',
  month: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20
})

const elderlyList = ref(getElderlyList().map(e => ({
  id: e.id,
  name: e.name,
  bedNumber: e.bedNumber
})))

const exportColumns = [
  { key: 'elderlyName', title: '老人姓名', width: 12 },
  { key: 'billNo', title: '账单号', width: 18 },
  { key: 'billType', title: '费用类型', width: 12 },
  { key: 'amount', title: '金额(元)', width: 12 },
  { key: 'billDate', title: '账单日期', width: 15 },
  { key: 'dueDate', title: '应付日期', width: 15 },
  { key: 'paidDate', title: '支付日期', width: 15 },
  { key: 'paymentMethod', title: '支付方式', width: 12 },
  { key: 'status', title: '状态', width: 10 },
  { key: 'remark', title: '备注', width: 20 }
]

// 支付对话框
const payDialogVisible = ref(false)
const paying = ref(false)
const currentBill = ref<any>(null)
const payForm = reactive({
  paymentMethod: 'wechat',
  cashType: 'cash'
})

// 是否有选中
const hasSelected = computed(() => selectedRows.value.length > 0)
const selectedCount = computed(() => selectedRows.value.length)

// 过滤后的账单列表
const filteredBillList = computed(() => {
  let list = billList.value
  if (showUnpaidOnly.value) {
    list = list.filter((b: any) => b.status === 'unpaid')
  }
  return list
})

// 判断是否逾期
const isOverdue = (dueDate: string, status: string) => {
  if (status === 'paid') return false
  return new Date(dueDate) < new Date()
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    paid: 'success',
    unpaid: 'warning',
    overdue: 'danger'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    paid: '已支付',
    unpaid: '未支付',
    overdue: '已逾期'
  }
  return map[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    // 使用种子数据
    await new Promise(resolve => setTimeout(resolve, 300))
    let data = getBills()

    // 应用筛选
    if (filterForm.elderlyId) {
      data = data.filter(b => b.elderlyId === filterForm.elderlyId)
    }
    if (filterForm.billType) {
      data = data.filter(b => b.billType === filterForm.billType)
    }
    if (filterForm.status) {
      data = data.filter(b => b.status === filterForm.status)
    }

    // 分页
    const start = (pagination.page - 1) * pagination.pageSize
    const end = start + pagination.pageSize
    billList.value = data.slice(start, end)

    total.value = data.length

    // 更新统计数据
    const stats = getStatistics()
    billStats.value = [
      { key: 'total', label: '本月应收', value: Number(stats.finance.totalAmount), trend: 5.2, type: 'primary' },
      { key: 'paid', label: '本月已收', value: Number(stats.finance.paidAmount), trend: 8.1, type: 'success' },
      { key: 'unpaid', label: '待收款', value: Number(stats.finance.unpaidAmount), trend: -12.5, type: 'warning' },
      { key: 'overdue', label: '逾期金额', value: 3200, trend: -25.0, type: 'danger' }
    ]
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  Object.assign(filterForm, {
    elderlyId: '',
    billType: '',
    status: '',
    month: ''
  })
  handleFilter()
}

const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadData()
}

const handleCreate = () => {
  console.log('生成账单')
}

const handleView = (row: any) => {
  console.log('查看账单', row)
}

const handlePay = (row: any) => {
  currentBill.value = row
  payDialogVisible.value = true
}

const handleConfirmPay = async () => {
  paying.value = true
  try {
    // TODO: API 请求
    await new Promise(resolve => setTimeout(resolve, 1500))
    ElMessage.success('支付成功')
    payDialogVisible.value = false
    loadData()
  } finally {
    paying.value = false
  }
}

const handleBatchPay = () => {
  console.log('批量支付', selectedRows.value)
}

const handleReminder = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要向家属发送催缴通知吗？`, '提示', {
      type: 'warning'
    })
    ElMessage.success('催缴通知已发送')
  } catch {
    // 取消
  }
}

const handlePrint = (row: any) => {
  console.log('打印账单', row)
}

const handleInvoice = (row: any) => {
  console.log('开具发票', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要作废此账单吗？', '提示', {
      type: 'warning'
    })
    ElMessage.success('账单已作废')
    loadData()
  } catch {
    // 取消
  }
}

const handleExport = async (params: any) => {
  if (params.type === 'all') {
    exportFinancialReport([])
  } else {
    exportFinancialReport(params.data)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.finance-bills {
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

  .bill-stat {
    padding: 20px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);

    .stat-label {
      font-size: 12px;
      color: var(--text-secondary);
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 8px;
    }

    .stat-trend {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;

      &.up {
        color: var(--success-color);
      }

      &.down {
        color: var(--danger-color);
      }
    }

    &.stat-primary .stat-value {
      color: var(--primary-color);
    }

    &.stat-success .stat-value {
      color: var(--success-color);
    }

    &.stat-warning .stat-value {
      color: var(--warning-color);
    }

    &.stat-danger .stat-value {
      color: var(--danger-color);
    }
  }

  .filter-card {
    margin-bottom: 20px;
  }

  .table-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .header-actions {
      display: flex;
      align-items: center;
      gap: 16px;
    }
  }

  .amount {
    font-weight: 600;
    color: var(--danger-color);
  }

  .text-danger {
    color: var(--danger-color);
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .pay-dialog {
    .pay-info {
      .info-item {
        display: flex;
        justify-content: space-between;
        padding: 8px 0;
        font-size: 14px;

        .label {
          color: var(--text-secondary);
        }

        .value {
          color: var(--text-primary);
        }

        .amount-value {
          font-size: 20px;
          font-weight: 600;
          color: var(--danger-color);
        }
      }
    }

    .pay-methods {
      margin-top: 20px;
    }

    .section-title {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 12px;
    }

    .method-list {
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    .method-item {
      display: block;
      padding: 12px 16px;
      border: 1px solid var(--border-color);
      border-radius: 8px;
      transition: all 0.3s;

      &:hover {
        border-color: var(--primary-color);
      }

      &.is-checked {
        border-color: var(--primary-color);
        background: var(--el-fill-color-light);
      }
    }

    .method-content {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .method-icon {
      font-size: 24px;
    }
  }
}

@media (max-width: 768px) {
  .finance-bills {
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

    .bill-stat {
      padding: 16px;

      .stat-value {
        font-size: 20px;
      }
    }
  }
}
</style>
