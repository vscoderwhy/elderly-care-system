<template>
  <div class="nursing-quality-assessment">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <h3>护理质量评估</h3>
          <div class="header-actions">
            <el-button size="small" @click="showHistory = !showHistory">
              {{ showHistory ? '隐藏历史' : '查看历史' }}
            </el-button>
            <el-button type="primary" size="small" @click="startNewAssessment">
              新建评估
            </el-button>
          </div>
        </div>
      </template>

      <!-- 评估表单 -->
      <div v-if="showForm" class="assessment-form">
        <el-form :model="assessment" :rules="rules" ref="formRef" label-width="120px">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="评估对象" prop="elderlyId">
                <el-select
                  v-model="assessment.elderlyId"
                  placeholder="请选择老人"
                  filterable
                  style="width: 100%"
                >
                  <el-option
                    v-for="elderly in elderlyList"
                    :key="elderly.id"
                    :label="`${elderly.name} - ${elderly.bedNumber}`"
                    :value="elderly.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="评估日期" prop="assessmentDate">
                <el-date-picker
                  v-model="assessment.assessmentDate"
                  type="date"
                  value-format="YYYY-MM-DD"
                  placeholder="选择日期"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="评估人" prop="assessor">
                <el-input v-model="assessment.assessor" placeholder="请输入评估人姓名" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="评估类型" prop="type">
                <el-select v-model="assessment.type" placeholder="请选择评估类型" style="width: 100%">
                  <el-option label="月度评估" value="monthly" />
                  <el-option label="季度评估" value="quarterly" />
                  <el-option label="年度评估" value="annual" />
                  <el-option label="专项评估" value="special" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 评分维度 -->
          <div class="dimension-section">
            <h4>评分维度</h4>
            <p class="section-tip">请对以下护理质量维度进行评分（1-100分）</p>

            <el-table :data="assessment.dimensions" border class="dimension-table">
              <el-table-column prop="name" label="评估维度" width="150" />
              <el-table-column prop="description" label="评估内容" min-width="200" />
              <el-table-column label="权重" width="100">
                <template #default="{ row }">
                  <el-input-number
                    v-model="row.weight"
                    :min="0"
                    :max="100"
                    :step="5"
                    size="small"
                    @change="calculateTotalScore"
                  />
                  <span style="margin-left: 4px">%</span>
                </template>
              </el-table-column>
              <el-table-column label="评分" width="150">
                <template #default="{ row }">
                  <el-slider
                    v-model="row.score"
                    :min="0"
                    :max="100"
                    :step="1"
                    show-input
                    @change="calculateTotalScore"
                  />
                </template>
              </el-table-column>
              <el-table-column label="加权得分" width="120">
                <template #default="{ row }">
                  <span class="weighted-score">{{ (row.score * row.weight) / 100 }}</span>
                </template>
              </el-table-column>
              <el-table-column label="评估说明" width="200">
                <template #default="{ row }">
                  <el-input
                    v-model="row.comment"
                    type="textarea"
                    :rows="1"
                    placeholder="选填"
                    size="small"
                  />
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- 总分展示 -->
          <div class="total-score-section">
            <el-result
              :icon="getScoreLevel(assessment.totalScore).icon"
              :title="`总分: ${assessment.totalScore}分`"
              :sub-title="getScoreLevel(assessment.totalScore).text"
            >
              <template #extra>
                <div class="score-level">
                  <el-tag :type="getScoreLevel(assessment.totalScore).type" size="large">
                    {{ getScoreLevel(assessment.totalScore).level }}
                  </el-tag>
                </div>
              </template>
            </el-result>
          </div>

          <!-- 雷达图对比 -->
          <div class="radar-section" v-if="showRadar">
            <h4>评估结果雷达图</h4>
            <div ref="radarChartRef" class="radar-chart"></div>
          </div>

          <!-- 改进建议 -->
          <el-form-item label="改进建议" prop="suggestions">
            <el-input
              v-model="assessment.suggestions"
              type="textarea"
              :rows="4"
              placeholder="请输入改进建议和措施"
            />
          </el-form-item>

          <el-form-item label="备注">
            <el-input
              v-model="assessment.remark"
              type="textarea"
              :rows="2"
              placeholder="其他备注信息"
            />
          </el-form-item>
        </el-form>

        <div class="form-actions">
          <el-button @click="showForm = false">取消</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button type="primary" @click="submitAssessment" :loading="submitting">
            提交评估
          </el-button>
        </div>
      </div>

      <!-- 历史记录 -->
      <div v-else-if="showHistory" class="history-section">
        <div class="history-filters">
          <el-form :inline="true" :model="historyFilter">
            <el-form-item label="评估对象">
              <el-select
                v-model="historyFilter.elderlyId"
                placeholder="全部"
                clearable
                style="width: 150px"
              >
                <el-option
                  v-for="elderly in elderlyList"
                  :key="elderly.id"
                  :label="elderly.name"
                  :value="elderly.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="评估类型">
              <el-select
                v-model="historyFilter.type"
                placeholder="全部"
                clearable
                style="width: 120px"
              >
                <el-option label="月度评估" value="monthly" />
                <el-option label="季度评估" value="quarterly" />
                <el-option label="年度评估" value="annual" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadHistory">查询</el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-table :data="historyList" v-loading="historyLoading" border>
          <el-table-column prop="assessmentDate" label="评估日期" width="120" />
          <el-table-column prop="elderlyName" label="评估对象" width="100" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="assessor" label="评估人" width="100" />
          <el-table-column prop="totalScore" label="总分" width="100">
            <template #default="{ row }">
              <el-tag :type="getScoreLevel(row.totalScore).type">{{ row.totalScore }}分</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="维度得分" min-width="300">
            <template #default="{ row }">
              <div class="dimension-scores">
                <span v-for="dim in row.dimensions" :key="dim.name" class="dim-score">
                  {{ dim.name }}: {{ dim.score }}
                </span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="suggestions" label="改进建议" min-width="200" show-overflow-tooltip />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button size="small" text @click="viewHistoryDetail(row)">查看</el-button>
              <el-button size="small" text type="primary" @click="compareWithHistory(row)">
                对比
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 默认空状态 -->
      <el-empty v-else description="点击"新建评估"开始护理质量评估" />
    </el-card>

    <!-- 历史对比弹窗 -->
    <el-dialog
      v-model="compareVisible"
      title="评估历史对比"
      width="900px"
      :close-on-click-modal="false"
    >
      <div class="compare-content" v-if="compareData.current && compareData.previous">
        <div class="compare-info">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="对比周期" :span="2">
              {{ compareData.previous.assessmentDate }} vs {{ compareData.current.assessmentDate }}
            </el-descriptions-item>
            <el-descriptions-item label="评估对象">
              {{ compareData.current.elderlyName }}
            </el-descriptions-item>
            <el-descriptions-item label="评估人">
              {{ compareData.current.assessor }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="compare-table-section">
          <h4>维度得分对比</h4>
          <el-table :data="compareData.dimensions" border>
            <el-table-column prop="name" label="评估维度" width="150" />
            <el-table-column label="上次得分" width="120">
              <template #default="{ row }">
                <span :class="getScoreChangeClass(row.change)">
                  {{ row.previousScore }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="本次得分" width="120">
              <template #default="{ row }">
                <span :class="getScoreChangeClass(row.change)">
                  {{ row.currentScore }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="变化" width="100">
              <template #default="{ row }">
                <el-tag v-if="row.change > 0" type="success" size="small">
                  +{{ row.change }}
                </el-tag>
                <el-tag v-else-if="row.change < 0" type="danger" size="small">
                  {{ row.change }}
                </el-tag>
                <el-tag v-else type="info" size="small">无变化</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="趋势" width="80">
              <template #default="{ row }">
                <el-icon v-if="row.change > 0" color="#67c23a"><Top /></el-icon>
                <el-icon v-else-if="row.change < 0" color="#f56c6c"><Bottom /></el-icon>
                <el-icon v-else color="#909399"><Minus /></el-icon>
              </template>
            </el-table-column>
            <el-table-column label="对比说明" min-width="200">
              <template #default="{ row }">
                <span class="compare-comment">{{ getCompareComment(row) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="compare-summary">
          <h4>总结</h4>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="summary-item">
                <div class="summary-label">总分变化</div>
                <div class="summary-value" :class="getScoreChangeClass(compareData.totalChange)">
                  {{ compareData.totalChange > 0 ? '+' : '' }}{{ compareData.totalChange }}分
                </div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="summary-item">
                <div class="summary-label">提升维度</div>
                <div class="summary-value success">
                  {{ compareData.improvedCount }}个
                </div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="summary-item">
                <div class="summary-label">下降维度</div>
                <div class="summary-value danger">
                  {{ compareData.declinedCount }}个
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>

      <template #footer>
        <el-button @click="compareVisible = false">关闭</el-button>
        <el-button type="primary" @click="exportCompare">导出对比报告</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="评估详情" width="800px">
      <div class="detail-content" v-if="currentDetail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="评估日期">
            {{ currentDetail.assessmentDate }}
          </el-descriptions-item>
          <el-descriptions-item label="评估类型">
            <el-tag size="small">{{ getTypeText(currentDetail.type) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="评估对象">
            {{ currentDetail.elderlyName }}
          </el-descriptions-item>
          <el-descriptions-item label="评估人">
            {{ currentDetail.assessor }}
          </el-descriptions-item>
          <el-descriptions-item label="总分" :span="2">
            <el-tag :type="getScoreLevel(currentDetail.totalScore).type" size="large">
              {{ currentDetail.totalScore }}分 - {{ getScoreLevel(currentDetail.totalScore).level }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div class="detail-dimensions">
          <h4>维度得分</h4>
          <el-table :data="currentDetail.dimensions" border size="small">
            <el-table-column prop="name" label="维度" width="120" />
            <el-table-column prop="score" label="得分" width="80" />
            <el-table-column prop="weight" label="权重" width="80" />
            <el-table-column prop="comment" label="说明" min-width="200" />
          </el-table>
        </div>

        <div class="detail-suggestions" v-if="currentDetail.suggestions">
          <h4>改进建议</h4>
          <p>{{ currentDetail.suggestions }}</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Top, Bottom, Minus } from '@element-plus/icons-vue'
import { getElderlyList } from '@/utils/seedData'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'

interface Dimension {
  name: string
  description: string
  weight: number
  score: number
  comment: string
}

interface Assessment {
  id?: number
  elderlyId: number | null
  elderlyName?: string
  assessmentDate: string
  assessor: string
  type: string
  dimensions: Dimension[]
  totalScore: number
  suggestions: string
  remark: string
}

const elderlyList = getElderlyList()

const showForm = ref(false)
const showHistory = ref(false)
const showRadar = ref(false)
const submitting = ref(false)
const historyLoading = ref(false)

const formRef = ref<FormInstance>()
const radarChartRef = ref<HTMLElement>()
let radarChart: ECharts | null = null

const assessment = reactive<Assessment>({
  elderlyId: null,
  assessmentDate: new Date().toISOString().split('T')[0],
  assessor: '',
  type: 'monthly',
  dimensions: [
    { name: '生活照料', description: '饮食、起居、卫生等日常护理', weight: 25, score: 80, comment: '' },
    { name: '健康监测', description: '生命体征、用药、健康数据记录', weight: 20, score: 80, comment: '' },
    { name: '康复护理', description: '康复训练、功能锻炼指导', weight: 15, score: 80, comment: '' },
    { name: '心理护理', description: '心理疏导、情感陪伴', weight: 15, score: 80, comment: '' },
    { name: '安全管理', description: '安全防护、风险评估', weight: 15, score: 80, comment: '' },
    { name: '服务质量', description: '服务态度、响应速度', weight: 10, score: 80, comment: '' }
  ],
  totalScore: 80,
  suggestions: '',
  remark: ''
})

const rules: FormRules = {
  elderlyId: [{ required: true, message: '请选择评估对象', trigger: 'change' }],
  assessmentDate: [{ required: true, message: '请选择评估日期', trigger: 'change' }],
  assessor: [{ required: true, message: '请输入评估人', trigger: 'blur' }],
  type: [{ required: true, message: '请选择评估类型', trigger: 'change' }]
}

const historyFilter = reactive({
  elderlyId: null,
  type: ''
})

const historyList = ref<Assessment[]>([])

const compareVisible = ref(false)
const compareData = reactive<any>({
  current: null,
  previous: null,
  dimensions: [],
  totalChange: 0,
  improvedCount: 0,
  declinedCount: 0
})

const detailVisible = ref(false)
const currentDetail = ref<Assessment | null>(null)

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    monthly: '月度评估',
    quarterly: '季度评估',
    annual: '年度评估',
    special: '专项评估'
  }
  return map[type] || type
}

const getScoreLevel = (score: number) => {
  if (score >= 90) return { level: '优秀', icon: 'success', type: 'success', text: '护理质量优秀' }
  if (score >= 80) return { level: '良好', icon: 'success', type: 'primary', text: '护理质量良好' }
  if (score >= 70) return { level: '中等', icon: 'warning', type: 'warning', text: '护理质量中等，需改进' }
  if (score >= 60) return { level: '合格', icon: 'warning', type: 'warning', text: '护理质量合格，需要提升' }
  return { level: '不合格', icon: 'error', type: 'danger', text: '护理质量不合格，急需改进' }
}

const getScoreChangeClass = (change: number) => {
  if (change > 0) return 'score-up'
  if (change < 0) return 'score-down'
  return ''
}

const getCompareComment = (row: any) => {
  if (row.change > 5) return '显著提升，继续保持'
  if (row.change > 0) return '有所提升'
  if (row.change < -5) return '明显下降，需要关注'
  if (row.change < 0) return '有所下降'
  return '保持稳定'
}

const calculateTotalScore = () => {
  let totalWeight = 0
  let weightedSum = 0

  assessment.dimensions.forEach(dim => {
    weightedSum += dim.score * dim.weight
    totalWeight += dim.weight
  })

  assessment.totalScore = totalWeight > 0 ? Math.round(weightedSum / totalWeight) : 0

  // 更新雷达图
  if (radarChart) {
    updateRadarChart()
  }
}

const initRadarChart = () => {
  if (!radarChartRef.value) return

  radarChart = echarts.init(radarChartRef.value)
  updateRadarChart()
}

const updateRadarChart = () => {
  if (!radarChart) return

  const option = {
    radar: {
      indicator: assessment.dimensions.map(dim => ({
        name: dim.name,
        max: 100
      }))
    },
    series: [
      {
        type: 'radar',
        data: [
          {
            value: assessment.dimensions.map(dim => dim.score),
            name: '当前评估',
            areaStyle: {
              color: 'rgba(64, 158, 255, 0.2)'
            }
          }
        ]
      }
    ]
  }

  radarChart.setOption(option)
}

const startNewAssessment = () => {
  resetForm()
  showForm.value = true
  showHistory.value = false
  showRadar.value = true

  nextTick(() => {
    initRadarChart()
  })
}

const resetForm = () => {
  Object.assign(assessment, {
    elderlyId: null,
    assessmentDate: new Date().toISOString().split('T')[0],
    assessor: '',
    type: 'monthly',
    dimensions: [
      { name: '生活照料', description: '饮食、起居、卫生等日常护理', weight: 25, score: 80, comment: '' },
      { name: '健康监测', description: '生命体征、用药、健康数据记录', weight: 20, score: 80, comment: '' },
      { name: '康复护理', description: '康复训练、功能锻炼指导', weight: 15, score: 80, comment: '' },
      { name: '心理护理', description: '心理疏导、情感陪伴', weight: 15, score: 80, comment: '' },
      { name: '安全管理', description: '安全防护、风险评估', weight: 15, score: 80, comment: '' },
      { name: '服务质量', description: '服务态度、响应速度', weight: 10, score: 80, comment: '' }
    ],
    totalScore: 80,
    suggestions: '',
    remark: ''
  })
  formRef.value?.clearValidate()
}

const submitAssessment = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    submitting.value = true

    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 添加到历史
    const elderly = elderlyList.find(e => e.id === assessment.elderlyId)
    const newAssessment = {
      ...assessment,
      id: Date.now(),
      elderlyName: elderly?.name || ''
    }
    historyList.value.unshift(newAssessment)

    ElMessage.success('评估提交成功')
    showForm.value = false
    showRadar.value = false
    showHistory.value = true
    loadHistory()
  } finally {
    submitting.value = false
  }
}

const loadHistory = async () => {
  historyLoading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))

    // 生成模拟历史数据
    if (historyList.value.length === 0) {
      const mockHistory: Assessment[] = []
      const types = ['monthly', 'quarterly', 'annual']

      for (let i = 0; i < 10; i++) {
        const elderly = elderlyList[i % elderlyList.length]
        const type = types[i % types.length]
        const dimensions = [
          { name: '生活照料', description: '饮食、起居、卫生等日常护理', weight: 25, score: 70 + Math.floor(Math.random() * 25), comment: '' },
          { name: '健康监测', description: '生命体征、用药、健康数据记录', weight: 20, score: 70 + Math.floor(Math.random() * 25), comment: '' },
          { name: '康复护理', description: '康复训练、功能锻炼指导', weight: 15, score: 70 + Math.floor(Math.random() * 25), comment: '' },
          { name: '心理护理', description: '心理疏导、情感陪伴', weight: 15, score: 70 + Math.floor(Math.random() * 25), comment: '' },
          { name: '安全管理', description: '安全防护、风险评估', weight: 15, score: 70 + Math.floor(Math.random() * 25), comment: '' },
          { name: '服务质量', description: '服务态度、响应速度', weight: 10, score: 70 + Math.floor(Math.random() * 25), comment: '' }
        ]

        let totalWeight = 0
        let weightedSum = 0
        dimensions.forEach(dim => {
          weightedSum += dim.score * dim.weight
          totalWeight += dim.weight
        })

        mockHistory.push({
          id: Date.now() + i,
          elderlyId: elderly.id,
          elderlyName: elderly.name,
          assessmentDate: new Date(Date.now() - i * 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
          assessor: ['张护士', '李护士', '王护士'][i % 3],
          type,
          dimensions,
          totalScore: Math.round(weightedSum / totalWeight),
          suggestions: '继续加强日常护理，关注老人心理健康',
          remark: ''
        })
      }

      historyList.value = mockHistory
    }
  } finally {
    historyLoading.value = false
  }
}

const viewHistoryDetail = (row: Assessment) => {
  currentDetail.value = row
  detailVisible.value = true
}

const compareWithHistory = (current: Assessment) => {
  // 找到上一次评估
  const currentIndex = historyList.value.findIndex(h => h.id === current.id)
  const previous = historyList.value[currentIndex + 1]

  if (!previous) {
    ElMessage.warning('没有找到历史评估数据进行对比')
    return
  }

  compareData.current = current
  compareData.previous = previous

  // 生成对比维度数据
  compareData.dimensions = current.dimensions.map(currDim => {
    const prevDim = previous.dimensions.find(d => d.name === currDim.name) || currDim
    return {
      name: currDim.name,
      currentScore: currDim.score,
      previousScore: prevDim.score,
      change: currDim.score - prevDim.score
    }
  })

  compareData.totalChange = current.totalScore - previous.totalScore
  compareData.improvedCount = compareData.dimensions.filter(d => d.change > 0).length
  compareData.declinedCount = compareData.dimensions.filter(d => d.change < 0).length

  compareVisible.value = true
}

const exportCompare = () => {
  ElMessage.success('对比报告导出功能开发中')
}

onMounted(() => {
  // 初始加载历史数据
  loadHistory()
})
</script>

<style scoped lang="scss">
.nursing-quality-assessment {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .assessment-form {
    .dimension-section {
      margin: 24px 0;
      padding: 20px;
      background: var(--bg-tertiary);
      border-radius: 8px;

      h4 {
        margin: 0 0 8px;
        font-size: 16px;
        font-weight: 600;
      }

      .section-tip {
        margin: 0 0 16px;
        font-size: 13px;
        color: var(--text-secondary);
      }

      .dimension-table {
        .weighted-score {
          font-weight: 600;
          color: var(--primary-color);
        }
      }
    }

    .total-score-section {
      margin: 24px 0;

      .score-level {
        margin-top: 12px;
      }
    }

    .radar-section {
      margin: 24px 0;

      h4 {
        margin: 0 0 16px;
        font-size: 16px;
        font-weight: 600;
      }

      .radar-chart {
        width: 100%;
        height: 400px;
      }
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      margin-top: 24px;
      padding-top: 24px;
      border-top: 1px solid var(--border-color-lighter);
    }
  }

  .history-section {
    .history-filters {
      margin-bottom: 16px;
    }

    .dimension-scores {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .dim-score {
        font-size: 12px;
        padding: 2px 8px;
        background: var(--bg-tertiary);
        border-radius: 4px;
      }
    }
  }

  .compare-content {
    .compare-info {
      margin-bottom: 20px;
    }

    .compare-table-section {
      margin: 20px 0;

      h4 {
        margin: 0 0 12px;
        font-size: 14px;
        font-weight: 600;
      }

      .score-up {
        color: var(--success-color);
        font-weight: 600;
      }

      .score-down {
        color: var(--danger-color);
        font-weight: 600;
      }

      .compare-comment {
        font-size: 13px;
        color: var(--text-secondary);
      }
    }

    .compare-summary {
      margin: 20px 0;
      padding: 16px;
      background: var(--bg-tertiary);
      border-radius: 8px;

      h4 {
        margin: 0 0 12px;
        font-size: 14px;
        font-weight: 600;
      }

      .summary-item {
        text-align: center;

        .summary-label {
          font-size: 12px;
          color: var(--text-secondary);
          margin-bottom: 8px;
        }

        .summary-value {
          font-size: 24px;
          font-weight: 600;

          &.success {
            color: var(--success-color);
          }

          &.danger {
            color: var(--danger-color);
          }
        }
      }
    }
  }

  .detail-content {
    .detail-dimensions {
      margin: 20px 0;

      h4 {
        margin: 0 0 12px;
        font-size: 14px;
        font-weight: 600;
      }
    }

    .detail-suggestions {
      margin: 20px 0;

      h4 {
        margin: 0 0 8px;
        font-size: 14px;
        font-weight: 600;
      }

      p {
        margin: 0;
        color: var(--text-secondary);
        line-height: 1.6;
      }
    }
  }
}
</style>
