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
        <el-form-item label="支付状态">
          <el-select v-model="filterForm.status" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="已支付" value="paid" />
            <el-option label="已退款" value="refunded" />
            <el-option label="退款中" value="refunding" />
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
        <el-table-column prop="paymentNo" label="支付编号" width="160" />
        <el-table-column prop="billNo" label="账单号" width="140" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="amount" label="支付金额" width="120">
          <template #default="{ row }">
            <span class="amount" :class="{ refunded: row.status === 'refunded' }">
              ¥{{ row.amount.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="paymentMethod" label="支付方式" width="100">
          <template #default="{ row }">
            <div class="payment-method">
              <el-icon :class="`method-${row.paymentMethod}`">
                <component :is="getMethodIcon(row.paymentMethod)" />
              </el-icon>
              <el-tag :type="getMethodType(row.paymentMethod)" size="small">
                {{ getMethodText(row.paymentMethod) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="paymentTime" label="支付时间" width="170" />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button
              size="small"
              text
              type="warning"
              @click="handleRefund(row)"
              v-if="row.canRefund"
            >
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

    <!-- 支付详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="支付详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="payment-detail" v-if="currentPayment">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="支付编号" :span="2">
            {{ currentPayment.paymentNo }}
          </el-descriptions-item>
          <el-descriptions-item label="账单号" :span="2">
            {{ currentPayment.billNo }}
          </el-descriptions-item>
          <el-descriptions-item label="老人姓名">
            {{ currentPayment.elderlyName }}
          </el-descriptions-item>
          <el-descriptions-item label="床位">
            {{ currentPayment.bedNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="支付金额">
            <span class="detail-amount">¥{{ currentPayment.amount.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="支付方式">
            <el-tag :type="getMethodType(currentPayment.paymentMethod)" size="small">
              {{ getMethodText(currentPayment.paymentMethod) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="支付状态">
            <el-tag :type="getStatusType(currentPayment.status)" size="small">
              {{ getStatusText(currentPayment.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="支付时间" :span="2">
            {{ currentPayment.paymentTime }}
          </el-descriptions-item>
          <el-descriptions-item label="操作人">
            {{ currentPayment.operator }}
          </el-descriptions-item>
          <el-descriptions-item label="退款时间" v-if="currentPayment.refundTime">
            {{ currentPayment.refundTime }}
          </el-descriptions-item>
          <el-descriptions-item label="备注" :span="2" v-if="currentPayment.remark">
            {{ currentPayment.remark }}
          </el-descriptions-item>
          <el-descriptions-item label="退款原因" :span="2" v-if="currentPayment.refundReason">
            {{ currentPayment.refundReason }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 账单明细 -->
        <div class="bill-items" v-if="currentPayment.billItems && currentPayment.billItems.length">
          <h4>账单明细</h4>
          <el-table :data="currentPayment.billItems" size="small" max-height="200">
            <el-table-column prop="type" label="费用类型" width="120" />
            <el-table-column prop="description" label="说明" />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                ¥{{ row.amount.toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 退款确认弹窗 -->
    <el-dialog
      v-model="refundVisible"
      title="退款确认"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="refundForm" :rules="refundRules" ref="refundFormRef" label-width="100px">
        <el-alert
          title="退款确认"
          type="warning"
          :closable="false"
          style="margin-bottom: 20px"
        >
          <p>确认要对 <strong>{{ refundForm.elderlyName }}</strong> 的支付记录退款吗？</p>
          <p>退款金额：<strong class="refund-amount">¥{{ refundForm.amount?.toFixed(2) }}</strong></p>
        </el-alert>

        <el-form-item label="退款金额" prop="amount">
          <el-input-number
            v-model="refundForm.refundAmount"
            :min="0.01"
            :max="refundForm.amount || 0"
            :precision="2"
            :step="100"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="退款原因" prop="reason">
          <el-select v-model="refundForm.reason" placeholder="请选择退款原因" style="width: 100%">
            <el-option label="重复支付" value="duplicate" />
            <el-option label="多收费用" value="overcharge" />
            <el-option label="服务取消" value="service_cancelled" />
            <el-option label="其他原因" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="详细说明" prop="remark">
          <el-input
            v-model="refundForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入详细说明"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="refundVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmRefund" :loading="refundLoading">
          确认退款
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Wallet,
  Money,
  TrendCharts,
  RefreshLeft,
  Coin,
  CreditCard,
  WalletFilled,
  Pointer
} from '@element-plus/icons-vue'
import ExportButton from '@/components/Export/ExportButton.vue'
import { getBills } from '@/utils/seedData'
import { getElderlyList } from '@/utils/seedData'

interface Payment {
  id: number
  paymentNo: string
  billNo: string
  elderlyId: number
  elderlyName: string
  bedNumber: string
  amount: number
  paymentMethod: string
  status: 'paid' | 'refunded' | 'refunding'
  paymentTime: string
  operator: string
  remark?: string
  refundTime?: string
  refundReason?: string
  canRefund: boolean
  billItems?: Array<{
    type: string
    description: string
    amount: number
  }>
}

const loading = ref(false)
const paymentList = ref<Payment[]>([])
const total = ref(0)
const allPayments = ref<Payment[]>([])

const detailVisible = ref(false)
const currentPayment = ref<Payment | null>(null)

const refundVisible = ref(false)
const refundLoading = ref(false)
const refundFormRef = ref<FormInstance>()

const refundForm = reactive({
  paymentId: 0,
  elderlyName: '',
  amount: 0,
  refundAmount: 0,
  reason: '',
  remark: ''
})

const refundRules: FormRules = {
  refundAmount: [{ required: true, message: '请输入退款金额', trigger: 'blur' }],
  reason: [{ required: true, message: '请选择退款原因', trigger: 'change' }],
  remark: [{ required: true, message: '请输入详细说明', trigger: 'blur' }]
}

const filterForm = reactive({
  paymentMethod: '',
  status: '',
  dateRange: [] as string[]
})

const pagination = reactive({
  page: 1,
  pageSize: 20
})

// 生成支付数据
const generatePayments = (): Payment[] => {
  const bills = getBills()
  const elderlyList = getElderlyList()
  const payments: Payment[] = []

  const methods = ['wechat', 'alipay', 'cash', 'pos', 'transfer']
  const statuses: Array<'paid' | 'refunded' | 'refunding'> = ['paid', 'paid', 'paid', 'paid', 'refunded']

  bills.forEach((bill, index) => {
    const elderly = elderlyList.find(e => e.id === bill.elderlyId)
    if (!elderly) return

    const paymentMethod = methods[index % methods.length]
    const status = statuses[index % statuses.length]
    const isPaid = status === 'paid'
    const isRefunded = status === 'refunded'

    const payment: Payment = {
      id: index + 1,
      paymentNo: `PAY${String(bill.id).padStart(10, '0')}`,
      billNo: bill.billNo,
      elderlyId: bill.elderlyId,
      elderlyName: elderly.name,
      bedNumber: elderly.bedNumber,
      amount: bill.totalAmount - (bill.paidAmount || 0),
      paymentMethod,
      status,
      paymentTime: isPaid
        ? bill.billingDate
        : new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toLocaleString('zh-CN'),
      operator: `财务-${['张三', '李四', '王五', '赵六'][index % 4]}`,
      remark: Math.random() > 0.7 ? '按月支付' : '',
      canRefund: isPaid && Math.random() > 0.5,
      billItems: [
        { type: '基础护理费', description: '本月护理费用', amount: bill.totalAmount * 0.6 },
        { type: '住宿费', description: '床位及住宿', amount: bill.totalAmount * 0.3 },
        { type: '餐饮费', description: '本月餐饮', amount: bill.totalAmount * 0.1 }
      ]
    }

    if (isRefunded) {
      payment.refundTime = new Date(
        new Date(payment.paymentTime).getTime() + 2 * 24 * 60 * 60 * 1000
      ).toLocaleString('zh-CN')
      payment.refundReason = ['重复支付', '多收费用', '服务取消'][index % 3]
      payment.canRefund = false
    }

    payments.push(payment)
  })

  return payments.sort((a, b) => new Date(b.paymentTime).getTime() - new Date(a.paymentTime).getTime())
}

// 计算统计数据
const stats = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const weekAgo = new Date(today)
  weekAgo.setDate(weekAgo.getDate() - 7)

  const monthAgo = new Date(today)
  monthAgo.setMonth(monthAgo.getMonth() - 1)

  let todayTotal = 0
  let weekTotal = 0
  let monthTotal = 0
  let refundTotal = 0

  allPayments.value.forEach(p => {
    const paymentDate = new Date(p.paymentTime)

    if (p.status === 'paid') {
      if (paymentDate >= today) todayTotal += p.amount
      if (paymentDate >= weekAgo) weekTotal += p.amount
      if (paymentDate >= monthAgo) monthTotal += p.amount
    } else if (p.status === 'refunded') {
      refundTotal += p.amount
    }
  })

  return [
    { key: 'today', label: '今日收款', value: todayTotal, type: 'success', icon: Money },
    { key: 'week', label: '本周收款', value: weekTotal, type: 'primary', icon: Wallet },
    { key: 'month', label: '本月收款', value: monthTotal, type: 'primary', icon: TrendCharts },
    { key: 'refund', label: '退款金额', value: refundTotal, type: 'danger', icon: RefreshLeft }
  ]
})

const exportColumns = [
  { key: 'paymentNo', title: '支付编号', width: 18 },
  { key: 'billNo', title: '账单号', width: 18 },
  { key: 'elderlyName', title: '老人姓名', width: 12 },
  { key: 'amount', title: '金额(元)', width: 12 },
  { key: 'paymentMethod', title: '支付方式', width: 12 },
  { key: 'status', title: '状态', width: 10 },
  { key: 'paymentTime', title: '支付时间', width: 18 },
  { key: 'operator', title: '操作人', width: 12 }
]

const getMethodIcon = (method: string) => {
  const map: Record<string, any> = {
    wechat: WalletFilled,
    alipay: Coin,
    cash: Money,
    pos: CreditCard,
    transfer: Pointer
  }
  return map[method] || Money
}

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

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    paid: 'success',
    refunded: 'danger',
    refunding: 'warning'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    paid: '已支付',
    refunded: '已退款',
    refunding: '退款中'
  }
  return map[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    // 生成支付数据
    allPayments.value = generatePayments()

    // 应用筛选
    let filtered = [...allPayments.value]

    if (filterForm.paymentMethod) {
      filtered = filtered.filter(p => p.paymentMethod === filterForm.paymentMethod)
    }

    if (filterForm.status) {
      filtered = filtered.filter(p => p.status === filterForm.status)
    }

    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      const [start, end] = filterForm.dateRange
      filtered = filtered.filter(p => {
        const date = p.paymentTime.split(' ')[0]
        return date >= start && date <= end
      })
    }

    // 分页
    total.value = filtered.length
    const start = (pagination.page - 1) * pagination.pageSize
    const end = start + pagination.pageSize
    paymentList.value = filtered.slice(start, end)
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
    status: '',
    dateRange: []
  })
  handleFilter()
}

const handleView = (row: Payment) => {
  currentPayment.value = row
  detailVisible.value = true
}

const handleRefund = (row: Payment) => {
  refundForm.paymentId = row.id
  refundForm.elderlyName = row.elderlyName
  refundForm.amount = row.amount
  refundForm.refundAmount = row.amount
  refundForm.reason = ''
  refundForm.remark = ''
  refundVisible.value = true
}

const confirmRefund = async () => {
  if (!refundFormRef.value) return

  try {
    await refundFormRef.value.validate()

    await ElMessageBox.confirm(
      `确认退款 ¥${refundForm.refundAmount.toFixed(2)}？此操作不可撤销！`,
      '退款确认',
      {
        type: 'warning',
        confirmButtonText: '确认退款',
        cancelButtonText: '取消'
      }
    )

    refundLoading.value = true

    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 更新本地数据
    const payment = allPayments.value.find(p => p.id === refundForm.paymentId)
    if (payment) {
      payment.status = 'refunding'
      payment.refundReason = refundForm.reason
      payment.refundTime = new Date().toLocaleString('zh-CN')
      payment.canRefund = false
    }

    ElMessage.success('退款申请已提交，将在1-3个工作日内到账')
    refundVisible.value = false
    loadData()
  } catch {
    // 用户取消
  } finally {
    refundLoading.value = false
  }
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
    transition: transform 0.2s;

    &:hover {
      transform: translateY(-2px);
    }

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

  .payment-method {
    display: flex;
    align-items: center;
    gap: 6px;

    .el-icon {
      font-size: 18px;

      &.method-wechat {
        color: #07c160;
      }

      &.method-alipay {
        color: #1677ff;
      }

      &.method-cash {
        color: #f59e0b;
      }

      &.method-pos {
        color: #8b5cf6;
      }

      &.method-transfer {
        color: #10b981;
      }
    }
  }

  .amount {
    font-weight: 600;
    color: var(--success-color);

    &.refunded {
      color: var(--danger-color);
      text-decoration: line-through;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .payment-detail {
    .detail-amount {
      font-size: 18px;
      font-weight: 600;
      color: var(--success-color);
    }

    .bill-items {
      margin-top: 20px;

      h4 {
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 12px;
        color: var(--text-primary);
      }
    }
  }

  .refund-amount {
    font-size: 18px;
    font-weight: 600;
    color: var(--danger-color);
  }
}
</style>
