<template>
  <div class="health-alert-system">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <h3>
            <el-icon><Warning /></el-icon>
            健康数据预警系统
          </h3>
          <div class="header-actions">
            <el-button size="small" @click="showRulesPanel = !showRulesPanel">
              {{ showRulesPanel ? '隐藏规则' : '配置规则' }}
            </el-button>
            <el-button type="primary" size="small" @click="showAddRuleDialog">
              <el-icon><Plus /></el-icon>
              新建规则
            </el-button>
          </div>
        </div>
      </template>

      <!-- 规则配置面板 -->
      <div v-if="showRulesPanel" class="rules-panel">
        <el-tabs v-model="activeRuleTab">
          <el-tab-pane label="预警规则" name="rules">
            <el-table :data="alertRules" border>
              <el-table-column prop="name" label="规则名称" width="150" />
              <el-table-column prop="metric" label="监测指标" width="120">
                <template #default="{ row }">
                  <el-tag size="small">{{ getMetricText(row.metric) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="预警阈值" width="200">
                <template #default="{ row }">
                  <div class="threshold-display">
                    <span v-if="row.condition === 'between'">
                      {{ row.minValue }} - {{ row.maxValue }}
                    </span>
                    <span v-else-if="row.condition === 'greater_than'">
                      > {{ row.maxValue }}
                    </span>
                    <span v-else-if="row.condition === 'less_than'">
                      < {{ row.minValue }}
                    </span>
                    <span class="unit">{{ getUnit(row.metric) }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="level" label="预警等级" width="100">
                <template #default="{ row }">
                  <el-tag :type="getLevelType(row.level)" size="small">
                    {{ getLevelText(row.level) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="enabled" label="状态" width="80">
                <template #default="{ row }">
                  <el-switch
                    v-model="row.enabled"
                    @change="toggleRule(row)"
                  />
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" text @click="editRule(row)">编辑</el-button>
                  <el-button size="small" text type="danger" @click="deleteRule(row)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="预警记录" name="history">
            <el-table :data="alertHistory" border v-loading="historyLoading">
              <el-table-column prop="time" label="触发时间" width="170" />
              <el-table-column prop="elderlyName" label="老人姓名" width="100" />
              <el-table-column prop="metric" label="指标" width="100">
                <template #default="{ row }">
                  {{ getMetricText(row.metric) }}
                </template>
              </el-table-column>
              <el-table-column label="触发值" width="120">
                <template #default="{ row }">
                  <span class="alert-value">{{ row.value }}{{ getUnit(row.metric) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="level" label="等级" width="80">
                <template #default="{ row }">
                  <el-tag :type="getLevelType(row.level)" size="small">
                    {{ getLevelText(row.level) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="ruleName" label="触发规则" width="150" />
              <el-table-column prop="message" label="预警信息" min-width="200" />
              <el-table-column prop="status" label="状态" width="80">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'resolved' ? 'success' : 'warning'" size="small">
                    {{ row.status === 'resolved' ? '已处理' : '待处理' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button
                    size="small"
                    text
                    type="primary"
                    @click="handleAlert(row)"
                    v-if="row.status === 'pending'"
                  >
                    处理
                  </el-button>
                  <el-button size="small" text @click="viewAlertDetail(row)">
                    详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="统计报表" name="stats">
            <el-row :gutter="20">
              <el-col :span="6" v-for="stat in alertStats" :key="stat.key">
                <div class="stat-card">
                  <div class="stat-icon" :class="`stat-${stat.type}`">
                    <component :is="stat.icon" />
                  </div>
                  <div class="stat-content">
                    <div class="stat-value">{{ stat.value }}</div>
                    <div class="stat-label">{{ stat.label }}</div>
                  </div>
                </div>
              </el-col>
            </el-row>

            <div class="chart-section">
              <h4>预警趋势</h4>
              <div ref="alertTrendChartRef" class="trend-chart"></div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 实时预警列表 -->
      <div v-else class="alerts-dashboard">
        <el-row :gutter="20">
          <el-col :span="18">
            <div class="active-alerts">
              <div class="section-header">
                <h4>实时预警</h4>
                <el-tag :type="pendingAlerts.length > 0 ? 'danger' : 'success'">
                  {{ pendingAlerts.length }}条待处理
                </el-tag>
              </div>

              <el-empty v-if="pendingAlerts.length === 0" description="暂无待处理预警" />

              <div v-else class="alert-list">
                <div
                  v-for="alert in pendingAlerts"
                  :key="alert.id"
                  class="alert-item"
                  :class="`alert-${alert.level}`"
                >
                  <div class="alert-icon">
                    <el-icon v-if="alert.level === 'critical'"><WarningFilled /></el-icon>
                    <el-icon v-else><Warning /></el-icon>
                  </div>
                  <div class="alert-content">
                    <div class="alert-header">
                      <span class="alert-title">{{ alert.elderlyName }}</span>
                      <el-tag :type="getLevelType(alert.level)" size="small">
                        {{ getLevelText(alert.level) }}
                      </el-tag>
                    </div>
                    <div class="alert-message">{{ alert.message }}</div>
                    <div class="alert-meta">
                      <span>{{ alert.metricName }}: {{ alert.value }}{{ alert.unit }}</span>
                      <span>{{ alert.time }}</span>
                    </div>
                  </div>
                  <div class="alert-actions">
                    <el-button size="small" type="primary" @click="handleAlert(alert)">
                      立即处理
                    </el-button>
                    <el-button size="small" @click="viewAlertDetail(alert)">详情</el-button>
                  </div>
                </div>
              </div>
            </div>
          </el-col>

          <el-col :span="6">
            <div class="alert-summary">
              <h4>预警统计</h4>
              <div class="summary-list">
                <div class="summary-item critical">
                  <div class="summary-label">紧急</div>
                  <div class="summary-value">{{ alertSummary.critical }}</div>
                </div>
                <div class="summary-item high">
                  <div class="summary-label">高危</div>
                  <div class="summary-value">{{ alertSummary.high }}</div>
                </div>
                <div class="summary-item medium">
                  <div class="summary-label">中等</div>
                  <div class="summary-value">{{ alertSummary.medium }}</div>
                </div>
                <div class="summary-item low">
                  <div class="summary-label">低危</div>
                  <div class="summary-value">{{ alertSummary.low }}</div>
                </div>
              </div>
            </div>

            <div class="quick-actions">
              <h4>快捷操作</h4>
              <el-button size="small" @click="simulateAlert" style="width: 100%; margin-bottom: 8px">
                <el-icon><Notification /></el-icon>
                模拟预警
              </el-button>
              <el-button size="small" @click="clearAllAlerts" style="width: 100%; margin-bottom: 8px">
                <el-icon><Delete /></el-icon>
                清除所有
              </el-button>
              <el-button size="small" @click="exportAlerts" style="width: 100%">
                <el-icon><Download /></el-icon>
                导出报告
              </el-button>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- 添加/编辑规则弹窗 -->
    <el-dialog
      v-model="ruleDialogVisible"
      :title="editingRule ? '编辑预警规则' : '新建预警规则'"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form :model="ruleForm" :rules="ruleFormRules" ref="ruleFormRef" label-width="120px">
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="ruleForm.name" placeholder="例如：高血压预警" />
        </el-form-item>

        <el-form-item label="监测指标" prop="metric">
          <el-select v-model="ruleForm.metric" placeholder="请选择监测指标" style="width: 100%">
            <el-option
              v-for="item in metricOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="预警条件" prop="condition">
          <el-select v-model="ruleForm.condition" placeholder="请选择条件" style="width: 100%">
            <el-option label="超出范围" value="between" />
            <el-option label="大于阈值" value="greater_than" />
            <el-option label="小于阈值" value="less_than" />
          </el-select>
        </el-form-item>

        <el-row :gutter="20" v-if="ruleForm.condition === 'between'">
          <el-col :span="12">
            <el-form-item label="最小值" prop="minValue">
              <el-input-number v-model="ruleForm.minValue" :precision="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最大值" prop="maxValue">
              <el-input-number v-model="ruleForm.maxValue" :precision="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="阈值" prop="maxValue" v-if="ruleForm.condition === 'greater_than'">
          <el-input-number v-model="ruleForm.maxValue" :precision="1" style="width: 100%" />
        </el-form-item>

        <el-form-item label="阈值" prop="minValue" v-if="ruleForm.condition === 'less_than'">
          <el-input-number v-model="ruleForm.minValue" :precision="1" style="width: 100%" />
        </el-form-item>

        <el-form-item label="预警等级" prop="level">
          <el-radio-group v-model="ruleForm.level">
            <el-radio label="low">低危</el-radio>
            <el-radio label="medium">中等</el-radio>
            <el-radio label="high">高危</el-radio>
            <el-radio label="critical">紧急</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="预警消息" prop="messageTemplate">
          <el-input
            v-model="ruleForm.messageTemplate"
            type="textarea"
            :rows="2"
            placeholder="支持变量: {name} {value} {unit} {elderly}"
          />
          <div class="form-tip">
            示例: {elderly}的{name}达到{value}{unit}，请立即处理！
          </div>
        </el-form-item>

        <el-form-item label="通知方式" prop="notificationMethods">
          <el-checkbox-group v-model="ruleForm.notificationMethods">
            <el-checkbox label="system">系统通知</el-checkbox>
            <el-checkbox label="sms">短信通知</el-checkbox>
            <el-checkbox label="email">邮件通知</el-checkbox>
            <el-checkbox label="wechat">微信通知</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="ruleForm.enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="ruleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRule" :loading="savingRule">保存</el-button>
      </template>
    </el-dialog>

    <!-- 预警处理弹窗 -->
    <el-dialog v-model="handleDialogVisible" title="处理预警" width="600px">
      <el-form :model="handleForm" label-width="100px">
        <el-form-item label="预警信息">
          <div class="handle-alert-info">
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="老人姓名">
                {{ handleForm.elderlyName }}
              </el-descriptions-item>
              <el-descriptions-item label="预警等级">
                <el-tag :type="getLevelType(handleForm.level)" size="small">
                  {{ getLevelText(handleForm.level) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="触发指标">
                {{ getMetricText(handleForm.metric) }}
              </el-descriptions-item>
              <el-descriptions-item label="触发值">
                {{ handleForm.value }}{{ getUnit(handleForm.metric) }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </el-form-item>

        <el-form-item label="处理措施" prop="action">
          <el-select v-model="handleForm.action" placeholder="请选择处理措施" style="width: 100%">
            <el-option label="联系医生" value="doctor" />
            <el-option label="通知家属" value="family" />
            <el-option label="调整护理计划" value="plan" />
            <el-option label="送医检查" value="hospital" />
            <el-option label="持续观察" value="observe" />
          </el-select>
        </el-form-item>

        <el-form-item label="处理说明" prop="note">
          <el-input
            v-model="handleForm.note"
            type="textarea"
            :rows="4"
            placeholder="请详细说明处理情况"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitHandle" :loading="submittingHandle">
          确认处理
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Warning,
  WarningFilled,
  Plus,
  Notification,
  Delete,
  Download,
  Bell,
  TrendCharts,
  CircleCheck,
  CircleClose
} from '@element-plus/icons-vue'
import { getElderlyList } from '@/utils/seedData'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'

interface AlertRule {
  id: number
  name: string
  metric: string
  condition: 'between' | 'greater_than' | 'less_than'
  minValue?: number
  maxValue?: number
  level: 'low' | 'medium' | 'high' | 'critical'
  messageTemplate: string
  notificationMethods: string[]
  enabled: boolean
}

interface Alert {
  id: number
  elderlyId: number
  elderlyName: string
  metric: string
  metricName: string
  value: number
  unit: string
  level: string
  ruleId: number
  ruleName: string
  message: string
  time: string
  status: 'pending' | 'resolved'
}

const showRulesPanel = ref(false)
const activeRuleTab = ref('rules')
const ruleDialogVisible = ref(false)
const handleDialogVisible = ref(false)
const editingRule = ref<AlertRule | null>(null)
const savingRule = ref(false)
const historyLoading = ref(false)
const submittingHandle = ref(false)

const ruleFormRef = ref<FormInstance>()
const alertTrendChartRef = ref<HTMLElement>()
let alertTrendChart: ECharts | null = null

const elderlyList = getElderlyList()

const metricOptions = [
  { label: '收缩压 (mmHg)', value: 'systolic_bp' },
  { label: '舒张压 (mmHg)', value: 'diastolic_bp' },
  { label: '心率 (次/分)', value: 'heart_rate' },
  { label: '血糖 (mmol/L)', value: 'blood_sugar' },
  { label: '体温 (℃)', value: 'temperature' },
  { label: '血氧饱和度 (%)', value: 'oxygen_saturation' },
  { label: '呼吸频率 (次/分)', value: 'respiration_rate' },
  { label: '体重 (kg)', value: 'weight' }
]

const alertRules = ref<AlertRule[]>([
  {
    id: 1,
    name: '高血压预警',
    metric: 'systolic_bp',
    condition: 'greater_than',
    maxValue: 140,
    level: 'high',
    messageTemplate: '{elderly}的收缩压达到{value}mmHg，超出正常范围！',
    notificationMethods: ['system', 'sms'],
    enabled: true
  },
  {
    id: 2,
    name: '低血糖预警',
    metric: 'blood_sugar',
    condition: 'less_than',
    minValue: 3.9,
    level: 'critical',
    messageTemplate: '{elderly}血糖过低({value}mmol/L)，请立即处理！',
    notificationMethods: ['system', 'sms', 'wechat'],
    enabled: true
  },
  {
    id: 3,
    name: '发热预警',
    metric: 'temperature',
    condition: 'greater_than',
    maxValue: 37.3,
    level: 'medium',
    messageTemplate: '{elderly}体温异常({value}℃)，请关注！',
    notificationMethods: ['system'],
    enabled: true
  },
  {
    id: 4,
    name: '心率异常预警',
    metric: 'heart_rate',
    condition: 'between',
    minValue: 60,
    maxValue: 100,
    level: 'high',
    messageTemplate: '{elderly}心率异常({value}次/分)！',
    notificationMethods: ['system', 'sms'],
    enabled: true
  },
  {
    id: 5,
    name: '血氧偏低预警',
    metric: 'oxygen_saturation',
    condition: 'less_than',
    minValue: 95,
    level: 'critical',
    messageTemplate: '{elderly}血氧饱和度过低({value}%)，需要吸氧！',
    notificationMethods: ['system', 'sms', 'wechat'],
    enabled: true
  }
])

const pendingAlerts = ref<Alert[]>([])

const alertHistory = ref<Alert[]>([])
const alertStats = ref([
  { key: 'today', label: '今日预警', value: 12, type: 'danger', icon: Warning },
  { key: 'handled', label: '已处理', value: 45, type: 'success', icon: CircleCheck },
  { key: 'pending', label: '待处理', value: 8, type: 'warning', icon: Bell },
  { key: 'resolved', label: '解决率', value: '85%', type: 'primary', icon: TrendCharts }
])

const alertSummary = reactive({
  critical: 2,
  high: 5,
  medium: 8,
  low: 3
})

const ruleForm = reactive<AlertRule>({
  id: 0,
  name: '',
  metric: '',
  condition: 'between',
  minValue: 0,
  maxValue: 100,
  level: 'medium',
  messageTemplate: '',
  notificationMethods: ['system'],
  enabled: true
})

const ruleFormRules: FormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  metric: [{ required: true, message: '请选择监测指标', trigger: 'change' }],
  condition: [{ required: true, message: '请选择预警条件', trigger: 'change' }],
  level: [{ required: true, message: '请选择预警等级', trigger: 'change' }],
  messageTemplate: [{ required: true, message: '请输入预警消息模板', trigger: 'blur' }]
}

const handleForm = reactive({
  id: 0,
  elderlyName: '',
  metric: '',
  level: '',
  value: 0,
  action: '',
  note: ''
})

const getMetricText = (metric: string) => {
  const option = metricOptions.find(o => o.value === metric)
  return option?.label || metric
}

const getUnit = (metric: string) => {
  const units: Record<string, string> = {
    systolic_bp: 'mmHg',
    diastolic_bp: 'mmHg',
    heart_rate: '次/分',
    blood_sugar: 'mmol/L',
    temperature: '℃',
    oxygen_saturation: '%',
    respiration_rate: '次/分',
    weight: 'kg'
  }
  return units[metric] || ''
}

const getLevelType = (level: string) => {
  const map: Record<string, any> = {
    critical: 'danger',
    high: 'warning',
    medium: 'primary',
    low: 'info'
  }
  return map[level] || ''
}

const getLevelText = (level: string) => {
  const map: Record<string, string> = {
    critical: '紧急',
    high: '高危',
    medium: '中等',
    low: '低危'
  }
  return map[level] || level
}

const showAddRuleDialog = () => {
  editingRule.value = null
  Object.assign(ruleForm, {
    id: 0,
    name: '',
    metric: '',
    condition: 'between',
    minValue: 0,
    maxValue: 100,
    level: 'medium',
    messageTemplate: '',
    notificationMethods: ['system'],
    enabled: true
  })
  ruleDialogVisible.value = true
}

const editRule = (rule: AlertRule) => {
  editingRule.value = rule
  Object.assign(ruleForm, {
    ...rule
  })
  ruleDialogVisible.value = true
}

const saveRule = async () => {
  if (!ruleFormRef.value) return

  try {
    await ruleFormRef.value.validate()

    savingRule.value = true

    await new Promise(resolve => setTimeout(resolve, 500))

    if (editingRule.value) {
      // 更新现有规则
      const index = alertRules.value.findIndex(r => r.id === editingRule.value!.id)
      if (index !== -1) {
        alertRules.value[index] = { ...ruleForm, id: editingRule.value!.id }
      }
      ElMessage.success('规则更新成功')
    } else {
      // 添加新规则
      alertRules.value.push({
        ...ruleForm,
        id: Date.now()
      })
      ElMessage.success('规则创建成功')
    }

    ruleDialogVisible.value = false
  } finally {
    savingRule.value = false
  }
}

const deleteRule = async (rule: AlertRule) => {
  try {
    await ElMessageBox.confirm(`确认删除规则"${rule.name}"吗？`, '删除确认', {
      type: 'warning'
    })

    const index = alertRules.value.findIndex(r => r.id === rule.id)
    if (index !== -1) {
      alertRules.value.splice(index, 1)
      ElMessage.success('规则已删除')
    }
  } catch {
    // 取消删除
  }
}

const toggleRule = (rule: AlertRule) => {
  ElMessage.success(rule.enabled ? '规则已启用' : '规则已禁用')
}

const handleAlert = (alert: Alert) => {
  handleForm.id = alert.id
  handleForm.elderlyName = alert.elderlyName
  handleForm.metric = alert.metric
  handleForm.level = alert.level
  handleForm.value = alert.value
  handleForm.action = ''
  handleForm.note = ''
  handleDialogVisible.value = true
}

const submitHandle = async () => {
  if (!handleForm.action) {
    ElMessage.warning('请选择处理措施')
    return
  }

  submittingHandle.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 更新预警状态
    const alert = pendingAlerts.value.find(a => a.id === handleForm.id)
    if (alert) {
      alert.status = 'resolved'
    }

    // 从待处理列表移除
    const index = pendingAlerts.value.findIndex(a => a.id === handleForm.id)
    if (index !== -1) {
      pendingAlerts.value.splice(index, 1)
    }

    // 添加到历史
    alertHistory.value.unshift({
      ...alert!,
      status: 'resolved'
    })

    ElMessage.success('预警已处理')
    handleDialogVisible.value = false
    updateSummary()
  } finally {
    submittingHandle.value = false
  }
}

const viewAlertDetail = (alert: Alert) => {
  ElMessage.info('查看预警详情功能开发中')
}

const simulateAlert = () => {
  const elderly = elderlyList[Math.floor(Math.random() * elderlyList.length)]
  const rule = alertRules.value.filter(r => r.enabled)[
    Math.floor(Math.random() * alertRules.value.filter(r => r.enabled).length)
  ]

  if (!rule) {
    ElMessage.warning('没有启用的预警规则')
    return
  }

  const alert: Alert = {
    id: Date.now(),
    elderlyId: elderly.id,
    elderlyName: elderly.name,
    metric: rule.metric,
    metricName: getMetricText(rule.metric),
    value: rule.maxValue ? rule.maxValue + Math.random() * 10 : Math.random() * 100,
    unit: getUnit(rule.metric),
    level: rule.level,
    ruleId: rule.id,
    ruleName: rule.name,
    message: rule.messageTemplate
      .replace('{elderly}', elderly.name)
      .replace('{name}', getMetricText(rule.metric))
      .replace('{value}', Math.round((rule.maxValue || 0) * 10) / 10)
      .replace('{unit}', getUnit(rule.metric)),
    time: new Date().toLocaleString('zh-CN'),
    status: 'pending'
  }

  pendingAlerts.value.unshift(alert)
  alertHistory.value.unshift(alert)
  updateSummary()

  ElNotification({
    title: '健康预警',
    message: alert.message,
    type: rule.level === 'critical' ? 'error' : 'warning',
    duration: 0
  })
}

const clearAllAlerts = async () => {
  if (pendingAlerts.value.length === 0) {
    ElMessage.info('没有待处理的预警')
    return
  }

  try {
    await ElMessageBox.confirm(`确认清除所有 ${pendingAlerts.value.length} 条预警吗？`, '清除确认', {
      type: 'warning'
    })

    pendingAlerts.value = []
    ElMessage.success('已清除所有预警')
  } catch {
    // 取消
  }
}

const exportAlerts = () => {
  ElMessage.success('预警报告导出功能开发中')
}

const updateSummary = () => {
  alertSummary.critical = pendingAlerts.value.filter(a => a.level === 'critical').length
  alertSummary.high = pendingAlerts.value.filter(a => a.level === 'high').length
  alertSummary.medium = pendingAlerts.value.filter(a => a.level === 'medium').length
  alertSummary.low = pendingAlerts.value.filter(a => a.level === 'low').length
}

const initTrendChart = () => {
  if (!alertTrendChartRef.value) return

  alertTrendChart = echarts.init(alertTrendChartRef.value)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['预警次数', '处理次数']
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '预警次数',
        type: 'line',
        smooth: true,
        data: [12, 15, 10, 18, 14, 16],
        itemStyle: { color: '#f56c6c' }
      },
      {
        name: '处理次数',
        type: 'line',
        smooth: true,
        data: [11, 14, 10, 17, 14, 15],
        itemStyle: { color: '#67c23a' }
      }
    ]
  }

  alertTrendChart.setOption(option)
}

onMounted(() => {
  // 初始化一些模拟预警数据
  for (let i = 0; i < 5; i++) {
    const elderly = elderlyList[i]
    const rule = alertRules.value[i % alertRules.value.length]

    pendingAlerts.value.push({
      id: Date.now() + i,
      elderlyId: elderly.id,
      elderlyName: elderly.name,
      metric: rule.metric,
      metricName: getMetricText(rule.metric),
      value: rule.maxValue || 0,
      unit: getUnit(rule.metric),
      level: rule.level,
      ruleId: rule.id,
      ruleName: rule.name,
      message: rule.messageTemplate
        .replace('{elderly}', elderly.name)
        .replace('{name}', getMetricText(rule.metric))
        .replace('{value}', String(rule.maxValue || 0))
        .replace('{unit}', getUnit(rule.metric)),
      time: new Date(Date.now() - i * 30 * 60 * 1000).toLocaleString('zh-CN'),
      status: 'pending'
    })
  }

  alertHistory.value = [...pendingAlerts.value]

  updateSummary()

  nextTick(() => {
    initTrendChart()
  })
})
</script>

<style scoped lang="scss">
.health-alert-system {
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

  .rules-panel {
    .threshold-display {
      display: flex;
      align-items: center;
      gap: 4px;

      .unit {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .stat-card {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px;
      background: var(--card-bg);
      border-radius: 8px;
      margin-bottom: 16px;

      .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 24px;

        &.stat-danger {
          background: var(--gradient-red);
          color: #fff;
        }

        &.stat-success {
          background: var(--gradient-green);
          color: #fff;
        }

        &.stat-warning {
          background: var(--gradient-orange);
          color: #fff;
        }

        &.stat-primary {
          background: var(--gradient-blue);
          color: #fff;
        }
      }

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: var(--text-primary);
      }

      .stat-label {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .chart-section {
      margin-top: 24px;

      h4 {
        margin: 0 0 16px;
        font-size: 14px;
        font-weight: 600;
      }

      .trend-chart {
        width: 100%;
        height: 300px;
      }
    }
  }

  .alerts-dashboard {
    .active-alerts {
      .section-header {
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

      .alert-list {
        display: flex;
        flex-direction: column;
        gap: 12px;
      }

      .alert-item {
        display: flex;
        align-items: flex-start;
        gap: 12px;
        padding: 16px;
        background: var(--card-bg);
        border-radius: 8px;
        border-left: 4px solid;
        transition: all 0.2s;

        &:hover {
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }

        &.alert-critical {
          border-left-color: #f56c6c;
          background: linear-gradient(to right, rgba(245, 108, 108, 0.05), transparent);
        }

        &.alert-high {
          border-left-color: #e6a23c;
          background: linear-gradient(to right, rgba(230, 162, 60, 0.05), transparent);
        }

        &.alert-medium {
          border-left-color: #409eff;
          background: linear-gradient(to right, rgba(64, 158, 255, 0.05), transparent);
        }

        &.alert-low {
          border-left-color: #909399;
          background: linear-gradient(to right, rgba(144, 147, 153, 0.05), transparent);
        }

        .alert-icon {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 20px;
          flex-shrink: 0;
        }

        .alert-critical & .alert-icon {
          background: rgba(245, 108, 108, 0.1);
          color: #f56c6c;
        }

        .alert-high & .alert-icon {
          background: rgba(230, 162, 60, 0.1);
          color: #e6a23c;
        }

        .alert-medium & .alert-icon {
          background: rgba(64, 158, 255, 0.1);
          color: #409eff;
        }

        .alert-low & .alert-icon {
          background: rgba(144, 147, 153, 0.1);
          color: #909399;
        }

        .alert-content {
          flex: 1;
          min-width: 0;

          .alert-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 6px;

            .alert-title {
              font-weight: 600;
              font-size: 15px;
            }
          }

          .alert-message {
            font-size: 14px;
            color: var(--text-primary);
            margin-bottom: 8px;
          }

          .alert-meta {
            display: flex;
            gap: 16px;
            font-size: 12px;
            color: var(--text-secondary);
          }
        }

        .alert-actions {
          display: flex;
          flex-direction: column;
          gap: 6px;
          flex-shrink: 0;
        }
      }
    }

    .alert-summary {
      padding: 16px;
      background: var(--card-bg);
      border-radius: 8px;
      margin-bottom: 16px;

      h4 {
        margin: 0 0 12px;
        font-size: 14px;
        font-weight: 600;
      }

      .summary-list {
        display: flex;
        flex-direction: column;
        gap: 8px;

        .summary-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 10px 12px;
          background: var(--bg-tertiary);
          border-radius: 6px;

          &.critical {
            border-left: 3px solid #f56c6c;
          }

          &.high {
            border-left: 3px solid #e6a23c;
          }

          &.medium {
            border-left: 3px solid #409eff;
          }

          &.low {
            border-left: 3px solid #909399;
          }

          .summary-label {
            font-size: 13px;
            color: var(--text-secondary);
          }

          .summary-value {
            font-size: 20px;
            font-weight: 600;
            color: var(--text-primary);
          }
        }
      }
    }

    .quick-actions {
      padding: 16px;
      background: var(--card-bg);
      border-radius: 8px;

      h4 {
        margin: 0 0 12px;
        font-size: 14px;
        font-weight: 600;
      }
    }
  }

  .form-tip {
    font-size: 12px;
    color: var(--text-secondary);
    margin-top: 4px;
  }

  .handle-alert-info {
    margin-bottom: 16px;
  }
}
</style>
