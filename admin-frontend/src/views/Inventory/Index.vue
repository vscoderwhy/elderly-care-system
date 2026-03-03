<template>
  <div class="inventory-page">
    <el-card class="stats-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic title="总物资数" :value="stats.total_items || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="低库存" :value="stats.low_stock_count || 0">
            <template #suffix>
              <el-icon color="#f56c6c"><Warning /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-statistic title="缺货" :value="stats.out_of_stock_count || 0">
            <template #suffix>
              <el-icon color="#f56c6c"><CircleClose /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="showPurchaseDialog">新建采购单</el-button>
        </el-col>
      </el-row>
    </el-card>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="分类">
          <el-select v-model="filterForm.category_id" placeholder="全部" clearable style="width: 150px">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filterForm.keyword" placeholder="名称/编码" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="正常" value="normal" />
            <el-option label="低库存" value="low_stock" />
            <el-option label="缺货" value="out_of_stock" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadInventories">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table :data="inventories" stripe>
        <el-table-column prop="code" label="编码" width="120" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column prop="spec" label="规格" width="100" />
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="quantity" label="库存" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'out_of_stock'" type="danger">{{ row.quantity }}</el-tag>
            <el-tag v-else-if="row.status === 'low_stock'" type="warning">{{ row.quantity }}</el-tag>
            <span v-else>{{ row.quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="min_quantity" label="最小库存" width="100" />
        <el-table-column prop="cost_price" label="成本价" width="100">
          <template #default="{ row }">¥{{ row.cost_price }}</template>
        </el-table-column>
        <el-table-column prop="location" label="存放位置" width="120" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleStockIn(row)">入库</el-button>
            <el-button size="small" @click="handleStockOut(row)">出库</el-button>
            <el-button size="small" @click="handleAdjust(row)">盘点</el-button>
            <el-button size="small" type="info" @click="showLogs(row)">记录</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 入库对话框 -->
    <el-dialog v-model="stockInVisible" title="入库" width="500px">
      <el-form :model="stockInForm" label-width="100px">
        <el-form-item label="物资名称">
          <el-input :value="currentItem?.name" disabled />
        </el-form-item>
        <el-form-item label="入库数量" required>
          <el-input-number v-model="stockInForm.quantity" :min="0.1" :step="1" :precision="2" />
        </el-form-item>
        <el-form-item label="批次号">
          <el-input v-model="stockInForm.batch_no" placeholder="选填" />
        </el-form-item>
        <el-form-item label="成本价">
          <el-input-number v-model="stockInForm.cost_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="stockInForm.reason" placeholder="采购入库/退货入库等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stockInVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmStockIn">确认</el-button>
      </template>
    </el-dialog>

    <!-- 出库对话框 -->
    <el-dialog v-model="stockOutVisible" title="出库" width="500px">
      <el-form :model="stockOutForm" label-width="100px">
        <el-form-item label="物资名称">
          <el-input :value="currentItem?.name" disabled />
        </el-form-item>
        <el-form-item label="当前库存">
          <el-input :value="currentItem?.quantity" disabled />
        </el-form-item>
        <el-form-item label="出库数量" required>
          <el-input-number v-model="stockOutForm.quantity" :min="0.1" :max="currentItem?.quantity || 0" :step="1" :precision="2" />
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="stockOutForm.reason" placeholder="领用/消耗/报废等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stockOutVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmStockOut">确认</el-button>
      </template>
    </el-dialog>

    <!-- 盘点对话框 -->
    <el-dialog v-model="adjustVisible" title="库存盘点" width="500px">
      <el-form :model="adjustForm" label-width="100px">
        <el-form-item label="物资名称">
          <el-input :value="currentItem?.name" disabled />
        </el-form-item>
        <el-form-item label="当前库存">
          <el-input :value="currentItem?.quantity" disabled />
        </el-form-item>
        <el-form-item label="实际数量" required>
          <el-input-number v-model="adjustForm.quantity" :min="0" :step="1" :precision="2" />
        </el-form-item>
        <el-form-item label="盘点原因" required>
          <el-input v-model="adjustForm.reason" type="textarea" placeholder="请输入盘点原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAdjust">确认</el-button>
      </template>
    </el-dialog>

    <!-- 采购单对话框 -->
    <el-dialog v-model="purchaseVisible" title="新建采购单" width="700px">
      <el-form :model="purchaseForm" label-width="100px">
        <el-form-item label="供应商" required>
          <el-input v-model="purchaseForm.supplier" />
        </el-form-item>
        <el-form-item label="采购明细">
          <el-button size="small" @click="addPurchaseItem">添加物资</el-button>
          <el-table :data="purchaseForm.items" style="margin-top: 10px">
            <el-table-column prop="name" label="物资" width="150">
              <template #default="{ row }">
                <el-select v-model="row.inventory_id" placeholder="选择物资">
                  <el-option v-for="item in inventories" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="数量" width="120">
              <template #default="{ row }">
                <el-input-number v-model="row.quantity" :min="1" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="单价" width="120">
              <template #default="{ row }">
                <el-input-number v-model="row.cost_price" :min="0" :precision="2" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="批次号" width="120">
              <template #default="{ row }">
                <el-input v-model="row.batch_no" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="60">
              <template #default="{ $index }">
                <el-button size="small" type="danger" link @click="removePurchaseItem($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="purchaseForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="purchaseVisible = false">取消</el-button>
        <el-button type="primary" @click="createPurchase">提交</el-button>
      </template>
    </el-dialog>

    <!-- 变动记录对话框 -->
    <el-dialog v-model="logsVisible" title="库存变动记录" width="800px">
      <el-table :data="logs" stripe>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.type === 'in'" type="success">入库</el-tag>
            <el-tag v-else-if="row.type === 'out'" type="warning">出库</el-tag>
            <el-tag v-else type="info">盘点</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="before_qty" label="变动前" width="100" />
        <el-table-column prop="after_qty" label="变动后" width="100" />
        <el-table-column prop="reason" label="原因" />
        <el-table-column prop="operator.username" label="操作人" width="100" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Warning, CircleClose } from '@element-plus/icons-vue'
import axios from 'axios'

const stats = ref<any>({})
const categories = ref<any[]>([])
const inventories = ref<any[]>([])
const logs = ref<any[]>([])
const currentItem = ref<any>(null)

const filterForm = ref({
  category_id: null as number | null,
  keyword: '',
  status: ''
})

// 入库
const stockInVisible = ref(false)
const stockInForm = ref({
  inventory_id: 0,
  quantity: 1,
  batch_no: '',
  cost_price: 0,
  reason: ''
})

// 出库
const stockOutVisible = ref(false)
const stockOutForm = ref({
  inventory_id: 0,
  quantity: 1,
  reason: ''
})

// 盘点
const adjustVisible = ref(false)
const adjustForm = ref({
  inventory_id: 0,
  quantity: 0,
  reason: ''
})

// 采购单
const purchaseVisible = ref(false)
const purchaseForm = ref({
  supplier: '',
  items: [] as any[],
  remark: ''
})

const logsVisible = ref(false)

const loadStats = async () => {
  const { data } = await axios.get('/api/inventory/stats')
  stats.value = data.data
}

const loadCategories = async () => {
  const { data } = await axios.get('/api/inventory/categories')
  categories.value = data.data
}

const loadInventories = async () => {
  const params: any = {}
  if (filterForm.value.category_id) params.category_id = filterForm.value.category_id
  if (filterForm.value.keyword) params.keyword = filterForm.value.keyword
  if (filterForm.value.status) params.status = filterForm.value.status

  const { data } = await axios.get('/api/inventory/items', { params })
  inventories.value = data.data
}

const resetFilter = () => {
  filterForm.value = { category_id: null, keyword: '', status: '' }
  loadInventories()
}

const handleStockIn = (item: any) => {
  currentItem.value = item
  stockInForm.value = {
    inventory_id: item.id,
    quantity: 1,
    batch_no: '',
    cost_price: item.cost_price,
    reason: '采购入库'
  }
  stockInVisible.value = true
}

const confirmStockIn = async () => {
  await axios.post('/api/inventory/stock-in', stockInForm.value)
  ElMessage.success('入库成功')
  stockInVisible.value = false
  loadInventories()
  loadStats()
}

const handleStockOut = (item: any) => {
  currentItem.value = item
  stockOutForm.value = {
    inventory_id: item.id,
    quantity: 1,
    reason: '领用'
  }
  stockOutVisible.value = true
}

const confirmStockOut = async () => {
  await axios.post('/api/inventory/stock-out', stockOutForm.value)
  ElMessage.success('出库成功')
  stockOutVisible.value = false
  loadInventories()
  loadStats()
}

const handleAdjust = (item: any) => {
  currentItem.value = item
  adjustForm.value = {
    inventory_id: item.id,
    quantity: item.quantity,
    reason: ''
  }
  adjustVisible.value = true
}

const confirmAdjust = async () => {
  await axios.post('/api/inventory/adjust', adjustForm.value)
  ElMessage.success('盘点成功')
  adjustVisible.value = false
  loadInventories()
  loadStats()
}

const showLogs = async (item: any) => {
  const { data } = await axios.get(`/api/inventory/items/${item.id}/logs`)
  logs.value = data.data
  logsVisible.value = true
}

const showPurchaseDialog = () => {
  purchaseForm.value = {
    supplier: '',
    items: [],
    remark: ''
  }
  purchaseVisible.value = true
}

const addPurchaseItem = () => {
  purchaseForm.value.items.push({
    inventory_id: null,
    quantity: 1,
    cost_price: 0,
    batch_no: ''
  })
}

const removePurchaseItem = (index: number) => {
  purchaseForm.value.items.splice(index, 1)
}

const createPurchase = async () => {
  await axios.post('/api/inventory/purchases', purchaseForm.value)
  ElMessage.success('采购单创建成功')
  purchaseVisible.value = false
}

onMounted(() => {
  loadStats()
  loadCategories()
  loadInventories()
})
</script>

<style scoped>
.inventory-page {
  padding: 20px;
}
.stats-card, .filter-card {
  margin-bottom: 20px;
}
</style>
