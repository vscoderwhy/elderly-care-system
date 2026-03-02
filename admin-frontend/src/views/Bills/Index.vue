<template>
  <div class="bills">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>财务管理</span>
          <el-radio-group v-model="statusFilter" size="small" @change="fetchData">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="unpaid">待支付</el-radio-button>
            <el-radio-button label="paid">已支付</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="bill_no" label="账单号" width="180" />
        <el-table-column prop="elderly.name" label="老人姓名" width="120" />
        <el-table-column prop="total_amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.total_amount }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'paid' ? 'success' : 'warning'">
              {{ row.status === 'paid' ? '已支付' : '待支付' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showPayDialog(row)" v-if="row.status !== 'paid'">
              缴费
            </el-button>
            <el-button type="info" size="small" @click="showDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
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

    <!-- 缴费对话框 -->
    <el-dialog v-model="payDialogVisible" title="在线缴费" width="500px">
      <div class="pay-info" v-if="selectedBill">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="账单号">{{ selectedBill.bill_no }}</el-descriptions-item>
          <el-descriptions-item label="老人姓名">{{ selectedBill.elderly?.name }}</el-descriptions-item>
          <el-descriptions-item label="账单金额">
            <span class="amount">¥{{ selectedBill.total_amount }}</span>
          </el-descriptions-item>
        </el-descriptions>

        <el-divider>选择支付方式</el-divider>

        <div class="pay-methods">
          <div
            class="pay-method"
            :class="{ active: payMethod === 'wechat' }"
            @click="payMethod = 'wechat'"
          >
            <el-icon size="32" color="#07C160"><Wallet /></el-icon>
            <span>微信支付</span>
          </div>
          <div
            class="pay-method"
            :class="{ active: payMethod === 'alipay' }"
            @click="payMethod = 'alipay'"
          >
            <el-icon size="32" color="#1677FF"><CreditCard /></el-icon>
            <span>支付宝</span>
          </div>
          <div
            class="pay-method"
            :class="{ active: payMethod === 'cash' }"
            @click="payMethod = 'cash'"
          >
            <el-icon size="32" color="#67C23A"><Money /></el-icon>
            <span>现金支付</span>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmPay" :loading="paying">
          确认支付 ¥{{ selectedBill?.total_amount || 0 }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 账单详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="账单详情" width="600px">
      <div v-if="selectedBill">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="账单号">{{ selectedBill.bill_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="selectedBill.status === 'paid' ? 'success' : 'warning'">
              {{ selectedBill.status === 'paid' ? '已支付' : '待支付' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="老人姓名">{{ selectedBill.elderly?.name }}</el-descriptions-item>
          <el-descriptions-item label="账单金额">¥{{ selectedBill.total_amount }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ selectedBill.created_at }}</el-descriptions-item>
          <el-descriptions-item label="支付时间" v-if="selectedBill.paid_at">{{ selectedBill.paid_at }}</el-descriptions-item>
        </el-descriptions>

        <el-divider>费用明细</el-divider>

        <el-table :data="selectedBill.items" border size="small">
          <el-table-column prop="name" label="项目名称" />
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column prop="unit_price" label="单价" width="100">
            <template #default="{ row }">¥{{ row.unit_price }}</template>
          </el-table-column>
          <el-table-column prop="amount" label="金额" width="100">
            <template #default="{ row }">¥{{ row.amount }}</template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Wallet, CreditCard, Money } from '@element-plus/icons-vue'
import instance, { billApi } from '@/api'

const loading = ref(false)
const tableData = ref([])
const statusFilter = ref('')
const payDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const selectedBill = ref<any>(null)
const payMethod = ref('wechat')
const paying = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    const result = await billApi.list(params)
    tableData.value = result.list || []
    pagination.total = result.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const showPayDialog = (bill: any) => {
  selectedBill.value = bill
  payMethod.value = 'wechat'
  payDialogVisible.value = true
}

const showDetail = async (bill: any) => {
  try {
    const result = await billApi.get(bill.id)
    selectedBill.value = result
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取详情失败')
  }
}

const confirmPay = async () => {
  if (!selectedBill.value) return

  paying.value = true
  try {
    await instance.post(`/bills/${selectedBill.value.id}/pay`, {
      amount: selectedBill.value.total_amount,
      method: payMethod.value
    })
    ElMessage.success('支付成功')
    payDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('支付失败')
  } finally {
    paying.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pay-info {
  padding: 10px 0;
}

.amount {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

.pay-methods {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}

.pay-method {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 30px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.pay-method:hover {
  border-color: #409EFF;
}

.pay-method.active {
  border-color: #409EFF;
  background: #ecf5ff;
}

.pay-method span {
  margin-top: 8px;
  font-size: 14px;
  color: #606266;
}
</style>
