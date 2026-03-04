<template>
  <div class="medication-management">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <h3>
            <el-icon><FirstAidKit /></el-icon>
            药品管理
          </h3>
          <div class="header-actions">
            <el-button size="small" @click="showRemindersPanel = !showRemindersPanel">
              <el-icon><Bell /></el-icon>
              今日提醒 ({{ todayReminders.length }})
            </el-button>
            <el-button size="small" @click="showInventoryPanel = !showInventoryPanel">
              <el-icon><Box /></el-icon>
              库存预警 ({{ lowStockCount }})
            </el-button>
            <el-button type="primary" size="small" @click="showAddMedicationDialog">
              <el-icon><Plus /></el-icon>
              新建药品档案
            </el-button>
          </div>
        </div>
      </template>

      <!-- 提醒面板 -->
      <div v-if="showRemindersPanel" class="reminders-panel">
        <div class="panel-header">
          <h4>今日用药提醒</h4>
          <el-tag :type="pendingReminders.length > 0 ? 'warning' : 'success'">
            {{ pendingReminders.length }}条待执行
          </el-tag>
        </div>

        <div class="reminder-list">
          <div
            v-for="reminder in todayReminders"
            :key="reminder.id"
            class="reminder-item"
            :class="{ completed: reminder.status === 'completed', overdue: isOverdue(reminder) }"
          >
            <div class="reminder-time">
              <div class="time">{{ reminder.time }}</div>
              <el-tag :type="getTimeTagType(reminder)" size="small">
                {{ getTimeText(reminder) }}
              </el-tag>
            </div>
            <div class="reminder-content">
              <div class="elderly-name">{{ reminder.elderlyName }}</div>
              <div class="medication-info">
                <span class="medication-name">{{ reminder.medicationName }}</span>
                <span class="dosage">{{ reminder.dosage }}</span>
              </div>
              <div class="medication-notes" v-if="reminder.notes">
                {{ reminder.notes }}
              </div>
            </div>
            <div class="reminder-actions">
              <el-button
                v-if="reminder.status === 'pending'"
                size="small"
                type="primary"
                @click="administerMedication(reminder)"
              >
                确认服药
              </el-button>
              <el-button size="small" @click="viewReminderDetail(reminder)">
                详情
              </el-button>
            </div>
          </div>
        </div>

        <el-empty v-if="todayReminders.length === 0" description="今日无用药提醒" />
      </div>

      <!-- 库存预警面板 -->
      <div v-else-if="showInventoryPanel" class="inventory-panel">
        <div class="panel-header">
          <h4>库存预警</h4>
        </div>

        <el-table :data="lowStockMedications" border>
          <el-table-column prop="name" label="药品名称" width="200" />
          <el-table-column prop="specification" label="规格" width="120" />
          <el-table-column prop="currentStock" label="当前库存" width="100">
            <template #default="{ row }">
              <span :class="{ 'low-stock': row.currentStock <= row.minStock }">
                {{ row.currentStock }} {{ row.unit }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="minStock" label="最小库存" width="100">
            <template #default="{ row }">
              {{ row.minStock }} {{ row.unit }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.currentStock <= row.minStock" type="danger" size="small">
                库存不足
              </el-tag>
              <el-tag v-else-if="row.currentStock <= row.minStock * 2" type="warning" size="small">
                库存偏低
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="supplier" label="供应商" min-width="150" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button size="small" text type="primary" @click="restock(row)">
                补货
              </el-button>
              <el-button size="small" text @click="viewInventoryDetail(row)">
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 药品档案列表 -->
      <div v-else class="medications-list">
        <div class="filter-bar">
          <el-form :inline="true" :model="filterForm">
            <el-form-item label="药品名称">
              <el-input v-model="filterForm.keyword" placeholder="搜索" clearable style="width: 200px" />
            </el-form-item>
            <el-form-item label="药品类型">
              <el-select v-model="filterForm.type" placeholder="全部" clearable style="width: 120px">
                <el-option label="西药" value="western" />
                <el-option label="中成药" value="chinese" />
                <el-option label="保健品" value="supplement" />
                <el-option label="外用药" value="external" />
              </el-select>
            </el-form-item>
            <el-form-item label="处方类型">
              <el-select v-model="filterForm.prescription" placeholder="全部" clearable style="width: 120px">
                <el-option label="处方药" value="prescription" />
                <el-option label="非处方药" value="otc" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadMedications">查询</el-button>
              <el-button @click="resetFilter">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-table :data="medications" border v-loading="loading">
          <el-table-column prop="name" label="药品名称" width="200" />
          <el-table-column prop="genericName" label="通用名" width="180" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="specification" label="规格" width="120" />
          <el-table-column prop="currentStock" label="库存" width="100">
            <template #default="{ row }">
              <span :class="{ 'low-stock': row.currentStock <= row.minStock }">
                {{ row.currentStock }} {{ row.unit }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="unitPrice" label="单价" width="100">
            <template #default="{ row }">
              ¥{{ row.unitPrice.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="manufacturer" label="生产厂家" min-width="150" show-overflow-tooltip />
          <el-table-column prop="expiryDate" label="有效期至" width="120">
            <template #default="{ row }">
              <span :class="{ 'expiring-soon': isExpiringSoon(row.expiryDate) }">
                {{ row.expiryDate }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button size="small" text @click="viewMedication(row)">查看</el-button>
              <el-button size="small" text type="primary" @click="editMedication(row)">
                编辑
              </el-button>
              <el-button size="small" text type="success" @click="adjustStock(row)">
                调整库存
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
          />
        </div>
      </div>
    </el-card>

    <!-- 新建/编辑药品弹窗 -->
    <el-dialog
      v-model="medicationDialogVisible"
      :title="editingMedication ? '编辑药品档案' : '新建药品档案'"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form :model="medicationForm" :rules="medicationFormRules" ref="medicationFormRef" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="药品名称" prop="name">
              <el-input v-model="medicationForm.name" placeholder="请输入商品名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通用名" prop="genericName">
              <el-input v-model="medicationForm.genericName" placeholder="请输入通用名" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="药品类型" prop="type">
              <el-select v-model="medicationForm.type" placeholder="请选择" style="width: 100%">
                <el-option label="西药" value="western" />
                <el-option label="中成药" value="chinese" />
                <el-option label="保健品" value="supplement" />
                <el-option label="外用药" value="external" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="处方类型" prop="prescription">
              <el-radio-group v-model="medicationForm.prescription">
                <el-radio label="prescription">处方药</el-radio>
                <el-radio label="otc">非处方药</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="规格" prop="specification">
              <el-input v-model="medicationForm.specification" placeholder="例如：100mg*30片" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="单位" prop="unit">
              <el-select v-model="medicationForm.unit" placeholder="单位" style="width: 100%">
                <el-option label="盒" value="盒" />
                <el-option label="瓶" value="瓶" />
                <el-option label="袋" value="袋" />
                <el-option label="支" value="支" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="单价" prop="unitPrice">
              <el-input-number v-model="medicationForm.unitPrice" :precision="2" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="当前库存" prop="currentStock">
              <el-input-number v-model="medicationForm.currentStock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最小库存" prop="minStock">
              <el-input-number v-model="medicationForm.minStock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="生产厂家" prop="manufacturer">
          <el-input v-model="medicationForm.manufacturer" placeholder="请输入生产厂家" />
        </el-form-item>

        <el-form-item label="供应商" prop="supplier">
          <el-input v-model="medicationForm.supplier" placeholder="请输入供应商" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="批号" prop="batchNumber">
              <el-input v-model="medicationForm.batchNumber" placeholder="请输入批号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="有效期至" prop="expiryDate">
              <el-date-picker
                v-model="medicationForm.expiryDate"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="适应症" prop="indications">
          <el-input
            v-model="medicationForm.indications"
            type="textarea"
            :rows="2"
            placeholder="请输入适应症"
          />
        </el-form-item>

        <el-form-item label="用法用量" prop="dosageInstructions">
          <el-input
            v-model="medicationForm.dosageInstructions"
            type="textarea"
            :rows="2"
            placeholder="例如：口服，一次1片，一日3次"
          />
        </el-form-item>

        <el-form-item label="注意事项" prop="precautions">
          <el-input
            v-model="medicationForm.precautions"
            type="textarea"
            :rows="2"
            placeholder="请输入注意事项"
          />
        </el-form-item>

        <el-form-item label="不良反应" prop="adverseReactions">
          <el-input
            v-model="medicationForm.adverseReactions"
            type="textarea"
            :rows="2"
            placeholder="请输入可能的不良反应"
          />
        </el-form-item>

        <el-form-item label="禁忌" prop="contraindications">
          <el-input
            v-model="medicationForm.contraindications"
            type="textarea"
            :rows="2"
            placeholder="请输入禁忌症"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="medicationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveMedication" :loading="savingMedication">
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- 用药计划弹窗 -->
    <el-dialog
      v-model="planDialogVisible"
      title="设置用药计划"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form :model="planForm" :rules="planFormRules" ref="planFormRef" label-width="120px">
        <el-form-item label="服务对象" prop="elderlyId">
          <el-select v-model="planForm.elderlyId" placeholder="请选择老人" filterable style="width: 100%">
            <el-option
              v-for="elderly in elderlyList"
              :key="elderly.id"
              :label="`${elderly.name} - ${elderly.bedNumber}`"
              :value="elderly.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="选择药品" prop="medicationId">
          <el-select v-model="planForm.medicationId" placeholder="请选择药品" filterable style="width: 100%">
            <el-option
              v-for="med in medications"
              :key="med.id"
              :label="`${med.name} (${med.specification})`"
              :value="med.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="用药剂量" prop="dosage">
          <el-input v-model="planForm.dosage" placeholder="例如：1片、5ml" />
        </el-form-item>

        <el-form-item label="用药频次" prop="frequency">
          <el-radio-group v-model="planForm.frequency">
            <el-radio label="once">一次</el-radio>
            <el-radio label="qd">每日一次</el-radio>
            <el-radio label="bid">每日两次</el-radio>
            <el-radio label="tid">每日三次</el-radio>
            <el-radio label="qid">每日四次</el-radio>
            <el-radio label="custom">自定义</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="用药时间" prop="times">
          <div class="time-inputs">
            <el-time-picker
              v-model="planForm.times"
              format="HH:mm"
              value-format="HH:mm"
              is-range
              range-separator="至"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              style="width: 100%"
            />
          </div>
          <div class="time-tip">
            提示：根据用药频次设置具体时间，如每日三次可设置为 08:00 - 20:00，系统将自动分配为 08:00、14:00、20:00
          </div>
        </el-form-item>

        <el-form-item label="给药途径" prop="route">
          <el-select v-model="planForm.route" placeholder="请选择" style="width: 100%">
            <el-option label="口服" value="oral" />
            <el-option label="注射" value="injection" />
            <el-option label="外用" value="topical" />
            <el-option label="吸入" value="inhalation" />
            <el-option label="滴眼/滴鼻" value="drops" />
          </el-select>
        </el-form-item>

        <el-form-item label="开始日期" prop="startDate">
          <el-date-picker
            v-model="planForm.startDate"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="结束日期">
          <el-date-picker
            v-model="planForm.endDate"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="planForm.notes"
            type="textarea"
            :rows="2"
            placeholder="特殊要求、注意事项等"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="planDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePlan" :loading="savingPlan">保存</el-button>
      </template>
    </el-dialog>

    <!-- 库存调整弹窗 -->
    <el-dialog v-model="stockDialogVisible" title="调整库存" width="500px">
      <el-form :model="stockForm" label-width="120px">
        <el-form-item label="药品名称">
          <el-input :value="stockForm.medicationName" disabled />
        </el-form-item>
        <el-form-item label="当前库存">
          <el-input :value="`${stockForm.currentStock} ${stockForm.unit}`" disabled />
        </el-form-item>
        <el-form-item label="调整类型">
          <el-radio-group v-model="stockForm.adjustType">
            <el-radio label="in">入库</el-radio>
            <el-radio label="out">出库</el-radio>
            <el-radio label="loss">损耗</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整数量">
          <el-input-number v-model="stockForm.adjustQuantity" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="调整后库存">
          <el-input :value="getNewStock()" disabled />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="stockForm.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入调整原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stockDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmStockAdjust">确认</el-button>
      </template>
    </el-dialog>

    <!-- 用药确认弹窗 -->
    <el-dialog v-model="administerDialogVisible" title="确认用药" width="600px">
      <el-form :model="administerForm" label-width="120px">
        <el-form-item label="老人姓名">
          <el-input :value="administerForm.elderlyName" disabled />
        </el-form-item>
        <el-form-item label="药品名称">
          <el-input :value="administerForm.medicationName" disabled />
        </el-form-item>
        <el-form-item label="用药剂量">
          <el-input :value="administerForm.dosage" disabled />
        </el-form-item>
        <el-form-item label="实际服用人">
          <el-input v-model="administerForm.administeredBy" placeholder="请输入护理员姓名" />
        </el-form-item>
        <el-form-item label="服用情况">
          <el-radio-group v-model="administerForm.status">
            <el-radio label="taken">正常服用</el-radio>
            <el-radio label="refused">拒绝服用</el-radio>
            <el-radio label="vomited">服药后呕吐</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="administerForm.notes"
            type="textarea"
            :rows="3"
            placeholder="请记录用药情况、老人反应等"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="administerDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAdminister">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, ElNotification } from 'element-plus'
import {
  FirstAidKit,
  Bell,
  Box,
  Plus
} from '@element-plus/icons-vue'
import { getElderlyList } from '@/utils/seedData'

interface Medication {
  id: number
  name: string
  genericName: string
  type: string
  prescription: string
  specification: string
  unit: string
  unitPrice: number
  currentStock: number
  minStock: number
  manufacturer: string
  supplier: string
  batchNumber: string
  expiryDate: string
  indications: string
  dosageInstructions: string
  precautions: string
  adverseReactions: string
  contraindications: string
}

interface MedicationReminder {
  id: number
  elderlyId: number
  elderlyName: string
  medicationId: number
  medicationName: string
  dosage: string
  time: string
  date: string
  notes: string
  status: 'pending' | 'completed' | 'missed'
}

const elderlyList = getElderlyList()

const showRemindersPanel = ref(false)
const showInventoryPanel = ref(false)
const loading = ref(false)

const medications = ref<Medication[]>([])
const reminders = ref<MedicationReminder[]>([])
const total = ref(0)

const medicationDialogVisible = ref(false)
const planDialogVisible = ref(false)
const stockDialogVisible = ref(false)
const administerDialogVisible = ref(false)

const editingMedication = ref<Medication | null>(null)
const savingMedication = ref(false)
const savingPlan = ref(false)

const medicationFormRef = ref<FormInstance>()
const planFormRef = ref<FormInstance>()

const filterForm = reactive({
  keyword: '',
  type: '',
  prescription: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20
})

const medicationForm = reactive<Medication>({
  id: 0,
  name: '',
  genericName: '',
  type: 'western',
  prescription: 'prescription',
  specification: '',
  unit: '盒',
  unitPrice: 0,
  currentStock: 0,
  minStock: 10,
  manufacturer: '',
  supplier: '',
  batchNumber: '',
  expiryDate: '',
  indications: '',
  dosageInstructions: '',
  precautions: '',
  adverseReactions: '',
  contraindications: ''
})

const medicationFormRules: FormRules = {
  name: [{ required: true, message: '请输入药品名称', trigger: 'blur' }],
  genericName: [{ required: true, message: '请输入通用名', trigger: 'blur' }],
  type: [{ required: true, message: '请选择药品类型', trigger: 'change' }],
  prescription: [{ required: true, message: '请选择处方类型', trigger: 'change' }],
  specification: [{ required: true, message: '请输入规格', trigger: 'blur' }],
  unit: [{ required: true, message: '请选择单位', trigger: 'change' }]
}

const planForm = reactive({
  elderlyId: 0,
  medicationId: 0,
  dosage: '',
  frequency: 'tid',
  times: [],
  route: 'oral',
  startDate: new Date().toISOString().split('T')[0],
  endDate: '',
  notes: ''
})

const planFormRules: FormRules = {
  elderlyId: [{ required: true, message: '请选择服务对象', trigger: 'change' }],
  medicationId: [{ required: true, message: '请选择药品', trigger: 'change' }],
  dosage: [{ required: true, message: '请输入用药剂量', trigger: 'blur' }],
  frequency: [{ required: true, message: '请选择用药频次', trigger: 'change' }]
}

const stockForm = reactive({
  medicationId: 0,
  medicationName: '',
  currentStock: 0,
  unit: '',
  adjustType: 'in',
  adjustQuantity: 1,
  remark: ''
})

const administerForm = reactive({
  reminderId: 0,
  elderlyName: '',
  medicationName: '',
  dosage: '',
  administeredBy: '',
  status: 'taken',
  notes: ''
})

const todayReminders = computed(() => {
  const today = new Date().toISOString().split('T')[0]
  return reminders.value.filter(r => r.date === today).sort((a, b) => a.time.localeCompare(b.time))
})

const pendingReminders = computed(() => {
  return todayReminders.value.filter(r => r.status === 'pending')
})

const lowStockMedications = computed(() => {
  return medications.value.filter(m => m.currentStock <= m.minStock * 2)
})

const lowStockCount = computed(() => {
  return medications.value.filter(m => m.currentStock <= m.minStock).length
})

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    western: '西药',
    chinese: '中成药',
    supplement: '保健品',
    external: '外用药'
  }
  return map[type] || type
}

const isExpiringSoon = (expiryDate: string) => {
  const today = new Date()
  const expiry = new Date(expiryDate)
  const daysLeft = Math.floor((expiry.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))
  return daysLeft <= 90
}

const isOverdue = (reminder: MedicationReminder) => {
  const now = new Date()
  const [hours, minutes] = reminder.time.split(':').map(Number)
  const reminderTime = new Date()
  reminderTime.setHours(hours, minutes, 0, 0)
  return reminder.status === 'pending' && now > reminderTime
}

const getTimeTagType = (reminder: MedicationReminder) => {
  if (reminder.status === 'completed') return 'success'
  if (isOverdue(reminder)) return 'danger'
  return 'info'
}

const getTimeText = (reminder: MedicationReminder) => {
  if (reminder.status === 'completed') return '已完成'
  if (isOverdue(reminder)) return '已超时'
  return '待执行'
}

const loadMedications = async () => {
  loading.value = true
  try {
    // 模拟API请求
    await new Promise(resolve => setTimeout(resolve, 500))

    // 生成模拟数据
    medications.value = []
    const types = ['western', 'chinese', 'supplement', 'external']
    const prescriptions = ['prescription', 'otc']
    const units = ['盒', '瓶', '袋', '支']

    for (let i = 0; i < 30; i++) {
      const type = types[i % types.length]
      const minStock = Math.floor(Math.random() * 20) + 5
      const currentStock = minStock + Math.floor(Math.random() * 30) - 10

      medications.value.push({
        id: Date.now() + i,
        name: `药品${i + 1}`,
        genericName: `Generic Name ${i + 1}`,
        type,
        prescription: prescriptions[i % prescriptions.length],
        specification: '100mg*30片',
        unit: units[i % units.length],
        unitPrice: Math.random() * 100 + 10,
        currentStock: Math.max(0, currentStock),
        minStock,
        manufacturer: `制药公司${i + 1}`,
        supplier: `医药供应商${i + 1}`,
        batchNumber: `BATCH${Date.now() + i}`,
        expiryDate: new Date(Date.now() + (Math.random() * 365 - 30) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        indications: '用于治疗...',
        dosageInstructions: '口服，一次1片，一日3次',
        precautions: '请遵医嘱',
        adverseReactions: '偶见恶心、呕吐',
        contraindications: '对本品过敏者禁用'
      })
    }

    total.value = medications.value.length
  } finally {
    loading.value = false
  }
}

const loadReminders = async () => {
  // 生成今日用药提醒
  reminders.value = []
  const today = new Date().toISOString().split('T')[0]
  const times = ['08:00', '12:00', '18:00']

  for (let i = 0; i < 20; i++) {
    const elderly = elderlyList[i % elderlyList.length]
    const medication = medications.value[i % medications.value.length]

    reminders.value.push({
      id: Date.now() + i,
      elderlyId: elderly.id,
      elderlyName: elderly.name,
      medicationId: medication.id,
      medicationName: medication.name,
      dosage: '1片',
      time: times[i % times.length],
      date: today,
      notes: '饭后服用',
      status: i < 12 ? 'completed' : 'pending'
    })
  }
}

const showAddMedicationDialog = () => {
  editingMedication.value = null
  Object.assign(medicationForm, {
    id: 0,
    name: '',
    genericName: '',
    type: 'western',
    prescription: 'prescription',
    specification: '',
    unit: '盒',
    unitPrice: 0,
    currentStock: 0,
    minStock: 10,
    manufacturer: '',
    supplier: '',
    batchNumber: '',
    expiryDate: '',
    indications: '',
    dosageInstructions: '',
    precautions: '',
    adverseReactions: '',
    contraindications: ''
  })
  medicationDialogVisible.value = true
}

const saveMedication = async () => {
  if (!medicationFormRef.value) return

  try {
    await medicationFormRef.value.validate()

    savingMedication.value = true

    await new Promise(resolve => setTimeout(resolve, 500))

    if (editingMedication.value) {
      const index = medications.value.findIndex(m => m.id === editingMedication.value!.id)
      if (index !== -1) {
        medications.value[index] = { ...medicationForm }
      }
      ElMessage.success('药品档案更新成功')
    } else {
      medications.value.unshift({
        ...medicationForm,
        id: Date.now()
      })
      ElMessage.success('药品档案创建成功')
    }

    medicationDialogVisible.value = false
  } finally {
    savingMedication.value = false
  }
}

const editMedication = (medication: Medication) => {
  editingMedication.value = medication
  Object.assign(medicationForm, medication)
  medicationDialogVisible.value = true
}

const viewMedication = (medication: Medication) => {
  ElMessage.info('查看药品详情功能开发中')
}

const adjustStock = (medication: Medication) => {
  stockForm.medicationId = medication.id
  stockForm.medicationName = medication.name
  stockForm.currentStock = medication.currentStock
  stockForm.unit = medication.unit
  stockForm.adjustType = 'in'
  stockForm.adjustQuantity = 1
  stockForm.remark = ''
  stockDialogVisible.value = true
}

const getNewStock = () => {
  let newStock = stockForm.currentStock
  if (stockForm.adjustType === 'in') {
    newStock += stockForm.adjustQuantity
  } else {
    newStock -= stockForm.adjustQuantity
  }
  return `${Math.max(0, newStock)} ${stockForm.unit}`
}

const confirmStockAdjust = async () => {
  const medication = medications.value.find(m => m.id === stockForm.medicationId)
  if (!medication) return

  if (stockForm.adjustType === 'in') {
    medication.currentStock += stockForm.adjustQuantity
  } else {
    medication.currentStock = Math.max(0, medication.currentStock - stockForm.adjustQuantity)
  }

  ElMessage.success('库存调整成功')
  stockDialogVisible.value = false
}

const restock = (medication: Medication) => {
  adjustStock(medication)
}

const viewInventoryDetail = (medication: Medication) => {
  viewMedication(medication)
}

const administerMedication = (reminder: MedicationReminder) => {
  administerForm.reminderId = reminder.id
  administerForm.elderlyName = reminder.elderlyName
  administerForm.medicationName = reminder.medicationName
  administerForm.dosage = reminder.dosage
  administerForm.administeredBy = ''
  administerForm.status = 'taken'
  administerForm.notes = ''
  administerDialogVisible.value = true
}

const confirmAdminister = async () => {
  if (!administerForm.administeredBy) {
    ElMessage.warning('请输入实际服用人')
    return
  }

  const reminder = reminders.value.find(r => r.id === administerForm.reminderId)
  if (reminder) {
    reminder.status = 'completed'
  }

  ElNotification({
    title: '用药已记录',
    message: `${administerForm.elderlyName} 已服用 ${administerForm.medicationName}`,
    type: 'success'
  })

  administerDialogVisible.value = false
}

const viewReminderDetail = (reminder: MedicationReminder) => {
  ElMessage.info('查看提醒详情功能开发中')
}

const resetFilter = () => {
  Object.assign(filterForm, {
    keyword: '',
    type: '',
    prescription: ''
  })
  loadMedications()
}

onMounted(() => {
  loadMedications()
  loadReminders()
})
</script>

<style scoped lang="scss">
.medication-management {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .reminders-panel,
  .inventory-panel {
    .panel-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      h4 {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
      }
    }
  }

  .reminder-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .reminder-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: var(--card-bg);
    border-radius: 8px;
    border-left: 4px solid var(--primary-color);
    transition: all 0.2s;

    &.completed {
      border-left-color: var(--success-color);
      opacity: 0.7;
    }

    &.overdue {
      border-left-color: var(--danger-color);
      background: var(--danger-color-lighter);
    }

    .reminder-time {
      text-align: center;
      min-width: 80px;

      .time {
        font-size: 18px;
        font-weight: 600;
        margin-bottom: 6px;
      }
    }

    .reminder-content {
      flex: 1;

      .elderly-name {
        font-weight: 600;
        margin-bottom: 6px;
      }

      .medication-info {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 4px;

        .medication-name {
          font-weight: 500;
        }

        .dosage {
          font-size: 12px;
          padding: 2px 6px;
          background: var(--bg-tertiary);
          border-radius: 4px;
        }
      }

      .medication-notes {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }
  }

  .filter-bar {
    margin-bottom: 16px;
  }

  .pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }

  .low-stock {
    color: var(--danger-color);
    font-weight: 600;
  }

  .expiring-soon {
    color: var(--warning-color);
    font-weight: 600;
  }

  .time-tip {
    font-size: 12px;
    color: var(--text-secondary);
    margin-top: 6px;
  }
}
</style>
