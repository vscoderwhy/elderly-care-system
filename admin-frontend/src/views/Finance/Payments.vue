<template>
  <div class="finance-payments">
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><Wallet /></el-icon>
        支付记录
      </h2>
      <ExportButton
        :data="paymentList"
        :columns="exportColumns"
        :total="total"
        filename="支付记录"
      />
    </div>

    <!-- 统计 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in stats" :key="stat.key">
        <div class="payment-stat">
          <div class="stat-icon" :class="`stat-${stat.type}`">
            <component :is="stat.icon" />
          </div>
          <div class="stat-content">
            <div class="stat-value">¥{{ stat.value.toLocaleString() }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 筛选 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="支付方式">
          <el-select v-model="filterForm.paymentMethod" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="微信支付" value="wechat" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="现金" value="cash" />
            <el-option label="POS机" value="pos" />
            <el-option label="银行转账" value="transfer" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 支付记录列表 -->
    <el-card shadow="never">
      <el-table :data="paymentList" stripe v-loading="loading">
        <el-table-column prop="paymentNo" label="支付编号" width="150" />
        <el-table-column prop="billNo" label="账单号" width="150" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="amount" label="支付金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="paymentMethod" label="支付方式" width="100">
          <template #default="{ row }">
            <el-tag :type="getMethodType(row.paymentMethod)" size="small">
              {{ getMethodText(row.paymentMethod) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="paymentTime" label="支付时间" width="180" />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handleRefund(row)" v-if="row.canRefund">
              退款
            </el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Wallet } from '@element-plus/icons-vue'
import ExportButton from '@/components/Export/ExportButton.vue'

const loading = ref(false)
const paymentList = ref([])
const total = ref(0)

const filterForm = reactive({
  paymentMethod: '',
  dateRange: []
})

const pagination = reactive({
  page: 1,
  pageSize: 20
})

const stats = ref([
  { key: 'today', label: '今日收款', value: 12500, type: 'success', icon: 'Money' },
  { key: 'week', label: '本周收款', value: 76800, type: 'primary', icon: 'Wallet' },
  { key: 'month', label: '本月收款', value: 289500, type: 'primary', icon: 'TrendCharts' },
  { key: 'refund', label: '退款金额', value: 1200, type: 'danger', icon: 'RefreshLeft' }
])

const exportColumns = [
  { key: 'paymentNo', title: '支付编号', width: 18 },
  { key: 'billNo', title: '账单号', width: 18 },
  { key: 'elderlyName', title: '老人姓名', width: 12 },
  { key: 'amount', title: '金额(元)', width: 12 },
  { key: 'paymentMethod', title: '支付方式', width: 12 },
  { key: 'paymentTime', title: '支付时间', width: 18 },
  { key: 'operator', title: '操作人', width: 12 }
]

const getMethodType = (method: string) => {
  const map: Record<string, any> = {
    wechat: 'success',
    alipay: 'primary',
    cash: 'warning',
    pos: 'info',
    transfer: ''
  }
  return map[method] || ''
}

const getMethodText = (method: string) => {
  const map: Record<string, string> = {
    wechat: '微信',
    alipay: '支付宝',
    cash: '现金',
    pos: 'POS机',
    transfer: '转账'
  }
  return map[method] || method
}

const loadData = async () => {
  loading.value = true
  try {
    // TODO: API 请求
    await new Promise(resolve => setTimeout(resolve, 500))
    paymentList.value = [
      {
        id: 1,
        paymentNo: 'PAY20260305001',
        billNo: 'B202603001',
        elderlyName: '张奶奶',
        bedNumber: '3号楼201',
        amount: 5300,
        paymentMethod: 'wechat',
        paymentTime: '2026-03-05 14:30:25',
        operator: '财务-张三',
        remark: '',
        canRefund: true
      },
      {
        id: 2,
        paymentNo: 'PAY20260305002',
        billNo: 'B202603002',
        elderlyName: '王爷爷',
        bedNumber: '3号楼202',
        amount: 3500,
        paymentMethod: 'alipay',
        paymentTime: '2026-03-05 10:15:30',
        operator: '财务-李四',
        remark: '',
        canRefund: true
      }
    ]
    total.value = 1256
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
    paymentMethod: '',
    dateRange: []
  })
  handleFilter()
}

const handleView = (row: any) => {
  console.log('查看详情', row)
}

const handleRefund = () => {
  ElMessage.success('退款申请已提交')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.finance-payments {
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

  .payment-stat {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;

      &.stat-success {
        background: var(--gradient-green);
        color: #fff;
      }

      &.stat-primary {
        background: var(--gradient-blue);
        color: #fff;
      }

      &.stat-danger {
        background: var(--gradient-red);
        color: #fff;
      }
    }

    .stat-value {
      font-size: 22px;
      font-weight: 600;
      color: var(--text-primary);
    }

    .stat-label {
      font-size: 12px;
      color: var(--text-secondary);
    }
  }

  .filter-card {
    margin-bottom: 20px;
  }

  .amount {
    font-weight: 600;
    color: var(--success-color);
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
