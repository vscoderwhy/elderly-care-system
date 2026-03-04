<template>
  <div class="care-assessment">
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><Document /></el-icon>
        护理评估
      </h2>
      <el-button type="primary" @click="handleNewAssessment">
        <el-icon><Plus /></el-icon>
        新建评估
      </el-button>
    </div>

    <!-- 评估列表 -->
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>评估记录列表</span>
          <div class="header-actions">
            <el-select v-model="filterForm.elderlyId" placeholder="选择老人" clearable style="width: 150px">
              <el-option
                v-for="e in elderlyList"
                :key="e.id"
                :label="e.name"
                :value="e.id"
              />
            </el-select>
            <el-select v-model="filterForm.type" placeholder="评估类型" clearable style="width: 120px">
              <el-option label="全部" value="" />
              <el-option label="入院评估" value="admission" />
              <el-option label="定期评估" value="periodic" />
              <el-option label="变更评估" value="change" />
            </el-select>
            <el-button @click="handleFilter">查询</el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredAssessments" stripe v-loading="loading">
        <el-table-column prop="assessmentNo" label="评估编号" width="150" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="assessmentType" label="评估类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getAssessmentTypeText(row.assessmentType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="assessmentDate" label="评估日期" width="120" />
        <el-table-column prop="nurseName" label="评估人" width="100" />
        <el-table-column prop="overallScore" label="综合评分" width="100">
          <template #default="{ row }">
            <el-tag :type="getScoreType(row.overallScore)" size="small">
              {{ row.overallScore }} 分
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="careLevel" label="建议护理等级" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'completed' ? 'success' : 'warning'" size="small">
              {{ row.status === 'completed' ? '已完成' : '待审核' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" text @click="handlePrint(row)">打印</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next"
        />
      </div>
    </el-card>

    <!-- 评估表单对话框 -->
    <el-dialog
      v-model="assessmentDialogVisible"
      :title="dialogTitle"
      width="900px"
      @close="handleDialogClose"
    >
      <el-form
        ref="assessmentFormRef"
        :model="assessmentForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="老人" prop="elderlyId">
              <el-select
                v-model="assessmentForm.elderlyId"
                placeholder="请选择老人"
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="e in elderlyList"
                  :key="e.id"
                  :label="`${e.name} (${e.bedNumber})`"
                  :value="e.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="评估类型" prop="assessmentType">
              <el-select v-model="assessmentForm.assessmentType" placeholder="请选择" style="width: 100%">
                <el-option label="入院评估" value="admission" />
                <el-option label="定期评估" value="periodic" />
                <el-option label="变更评估" value="change" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="评估日期" prop="assessmentDate">
              <el-date-picker
                v-model="assessmentForm.assessmentDate"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="评估人" prop="nurseId">
              <el-select v-model="assessmentForm.nurseId" placeholder="请选择" style="width: 100%">
                <el-option
                  v-for="n in nurseList"
                  :key="n.id"
                  :label="n.name"
                  :value="n.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 评估维度 -->
        <div class="assessment-section">
          <h4 class="section-title">护理质量评估（{{ getAssessmentScore() }}分）</h4>

          <el-table :data="assessmentForm.dimensions" border class="dimensions-table">
            <el-table-column prop="name" label="评估维度" width="150" />
            <el-table-column prop="description" label="评估内容" min-width="200" />
            <el-table-column label="评分" width="180">
              <template #default="{ row, $index }">
                <el-slider
                  v-model="row.score"
                  :max="100"
                  :step="5"
                  :marks="{ 0: '0', 50: '50', 100: '100' }"
                  show-input
                />
              </template>
            </el-table-column>
            <el-table-column prop="weight" label="权重" width="80" />
          </el-table>
        </div>

        <div class="assessment-section">
          <h4 class="section-title">日常生活能力（ADL）评估</h4>

          <el-form-item label=" Barthel指数">
            <el-radio-group v-model="assessmentForm.adlScore">
              <el-radio :label="100">完全自理（100分）</el-radio>
              <el-radio :label="60">轻度依赖（61-99分）</el-radio>
              <el-radio :label="40">中度依赖（41-60分）</el-radio>
              <el-radio :label="20">重度依赖（≤40分）</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="能力评估">
            <el-checkbox-group v-model="assessmentForm.adlItems">
              <el-checkbox label="eating">进食</el-checkbox>
              <el-checkbox label="bathing">洗澡</el-checkbox>
              <el-checkbox label="dressing">穿衣</el-checkbox>
              <el-checkbox label="toileting">如厕</el-checkbox>
              <el-checkbox label="moving">活动</el-checkbox>
              <el-checkbox label="controlling"排便控制</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </div>

        <div class="assessment-section">
          <h4 class="section-title">认知功能评估</h4>

          <el-form-item label="认知状态">
            <el-radio-group v-model="assessmentForm.cognitiveStatus">
              <el-radio label="normal">正常</el-radio>
              <el-radio label="mild">轻度障碍</el-radio>
              <el-radio label="moderate">中度障碍</el-radio>
              <el-radio label="severe">重度障碍</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="精神状态">
            <el-input
              v-model="assessmentForm.mentalStatus"
              type="textarea"
              :rows="3"
              placeholder="请描述老人精神状态、情绪表现等"
            />
          </el-form-item>
        </div>

        <el-form-item label="建议护理等级" prop="careLevel">
          <el-radio-group v-model="assessmentForm.careLevel">
            <el-radio label="三级">三级护理</el-radio>
            <el-radio label="二级">二级护理</el-radio>
            <el-radio label="一级">一级护理</el-radio>
            <el-radio label="special">特级护理</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="综合建议">
          <el-input
            v-model="assessmentForm.suggestions"
            type="textarea"
            :rows="4"
            placeholder="请输入综合护理建议和注意事项"
          />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="assessmentForm.remark"
            type="textarea"
            :rows="2"
            placeholder="其他需要说明的事项"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="assessmentDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAssessment" :loading="submitting">
          {{ submitting ? '保存中...' : '保存' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="评估详情"
      width="800px"
    >
      <div v-if="currentAssessment" class="assessment-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="评估编号">
            {{ currentAssessment.assessmentNo }}
          </el-descriptions-item>
          <el-descriptions-item label="老人姓名">
            {{ currentAssessment.elderlyName }}
          </el-descriptions-item>
          <el-descriptions-item label="评估类型">
            {{ getAssessmentTypeText(currentAssessment.assessmentType) }}
          </el-descriptions-item>
          <el-descriptions-item label="评估日期">
            {{ currentAssessment.assessmentDate }}
          </el-descriptions-item>
          <el-descriptions-item label="评估人">
            {{ currentAssessment.nurseName }}
          </el-descriptions-item>
          <el-descriptions-item label="综合评分">
            <el-tag :type="getScoreType(currentAssessment.overallScore)">
              {{ currentAssessment.overallScore }} 分
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="建议护理等级" :span="2">
            {{ currentAssessment.careLevel }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentAssessment.suggestions" class="detail-section">
          <h4>综合建议</h4>
          <p>{{ currentAssessment.suggestions }}</p>
        </div>
      </div>

      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handlePrint(currentAssessment)">打印</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Document, Plus } from '@element-plus/icons-vue'
import { getElderlyList, getCareRecords } from '@/utils/seedData'

const loading = ref(false)
const submitting = ref(false)
const assessmentDialogVisible = ref(false)
const viewDialogVisible = ref(false)
const dialogTitle = ref('新建评估')
const assessmentFormRef = ref<FormInstance>()
const currentAssessment = ref<any>(null)

// 筛选表单
const filterForm = reactive({
  elderlyId: '',
  type: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(0)

// 老人和护理员列表
const elderlyList = ref(getElderlyList())
const nurseList = ref([
  { id: 1, name: '赵护士' },
  { id: 2, name: '钱护士' },
  { id: 3, name: '孙护士' },
  { id: 4, name: '李护士' }
])

// 评估列表
const assessmentList = ref([
  {
    id: 1,
    assessmentNo: 'AS202603001',
    elderlyId: 1,
    elderlyName: '张奶奶',
    assessmentType: 'admission',
    assessmentDate: '2026-03-01',
    nurseId: 1,
    nurseName: '赵护士',
    overallScore: 85,
    careLevel: '二级',
    status: 'completed',
    adlScore: 60,
    cognitiveStatus: 'mild',
    suggestions: '建议加强日常活动能力训练，注意防跌倒',
    dimensions: [
      { name: '生活照料', description: '日常护理质量', score: 85, weight: '25%' },
      { name: '健康监测', description: '健康数据监测', score: 90, weight: '20%' },
      { name: '康复护理', description: '康复训练', score: 80, weight: '20%' },
      { name: '心理护理', description: '心理疏导', score: 88, weight: '15%' },
      { name: '安全管理', description: '安全防护', score: 82, weight: '20%' }
    ]
  },
  {
    id: 2,
    assessmentNo: 'AS202603002',
    elderlyId: 2,
    elderlyName: '王爷爷',
    assessmentType: 'periodic',
    assessmentDate: '2026-02-28',
    nurseId: 2,
    nurseName: '钱护士',
    overallScore: 72,
    careLevel: '一级',
    status: 'completed',
    adlScore: 40,
    cognitiveStatus: 'moderate',
    suggestions: '需要较多协助，建议增加护理时间',
    dimensions: []
  },
  {
    id: 3,
    assessmentNo: 'AS202603003',
    elderlyId: 3,
    elderlyName: '李奶奶',
    assessmentType: 'admission',
    assessmentDate: '2026-03-02',
    nurseId: 1,
    nurseName: '赵护士',
    overallScore: 68,
    careLevel: '一级',
    status: 'pending',
    adlScore: 40,
    cognitiveStatus: 'moderate',
    suggestions: '',
    dimensions: []
  }
])

// 评估表单
const assessmentForm = reactive({
  id: undefined,
  elderlyId: '',
  assessmentType: 'admission',
  assessmentDate: '',
  nurseId: '',
  careLevel: '二级',
  adlScore: 100,
  adlItems: [],
  cognitiveStatus: 'normal',
  mentalStatus: '',
  suggestions: '',
  remark: '',
  dimensions: [
    { name: '生活照料', description: '日常护理质量', score: 80, weight: '25%' },
    { name: '健康监测', description: '健康数据监测', score: 80, weight: '20%' },
    { name: '康复护理', description: '康复训练', score: 80, weight: '20%' },
    { name: '心理护理', description: '心理疏导', score: 80, weight: '15%' },
    { name: '安全管理', description: '安全防护', score: 80, weight: '20%' }
  ]
})

// 表单验证规则
const formRules: FormRules = {
  elderlyId: [{ required: true, message: '请选择老人', trigger: 'change' }],
  assessmentType: [{ required: true, message: '请选择评估类型', trigger: 'change' }],
  assessmentDate: [{ required: true, message: '请选择评估日期', trigger: 'change' }],
  nurseId: [{ required: true, message: '请选择评估人', trigger: 'change' }],
  careLevel: [{ required: true, message: '请选择护理等级', trigger: 'change' }]
}

// 计算评估分数
const getAssessmentScore = () => {
  const { dimensions } = assessmentForm
  let totalScore = 0
  dimensions.forEach((d: any) => {
    totalScore += d.score * (parseInt(d.weight) / 100)
  })
  return Math.round(totalScore)
}

// 过滤后的评估列表
const filteredAssessments = computed(() => {
  let list = assessmentList.value

  if (filterForm.elderlyId) {
    list = list.filter(a => a.elderlyId === Number(filterForm.elderlyId))
  }

  if (filterForm.type) {
    list = list.filter(a => a.assessmentType === filterForm.type)
  }

  return list
})

// 获取评估类型文本
const getAssessmentTypeText = (type: string) => {
  const map: Record<string, string> = {
    admission: '入院评估',
    periodic: '定期评估',
    change: '变更评估'
  }
  return map[type] || type
}

// 获取分数类型
const getScoreType = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

// 筛选
const handleFilter = () => {
  console.log('筛选', filterForm)
}

// 新建评估
const handleNewAssessment = () => {
  dialogTitle.value = '新建评估'
  resetForm()
  assessmentForm.assessmentDate = new Date().toISOString().slice(0, 10)
  assessmentDialogVisible.value = true
}

// 编辑评估
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑评估'
  Object.assign(assessmentForm, {
    ...row,
    dimensions: row.dimensions ? [...row.dimensions] : assessmentForm.dimensions
  })
  assessmentDialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  Object.assign(assessmentForm, {
    id: undefined,
    elderlyId: '',
    assessmentType: 'admission',
    assessmentDate: '',
    nurseId: '',
    careLevel: '二级',
    adlScore: 100,
    adlItems: [],
    cognitiveStatus: 'normal',
    mentalStatus: '',
    suggestions: '',
    remark: '',
    dimensions: [
      { name: '生活照料', description: '日常护理质量', score: 80, weight: '25%' },
      { name: '健康监测', description: '健康数据监测', score: 80, weight: '20%' },
      { name: '康复护理', description: '康复训练', score: 80, weight: '20%' },
      { name: '心理护理', description: '心理疏导', score: 80, weight: '15%' },
      { name: '安全管理', description: '安全防护', score: 80, weight: '20%' }
    ]
  })
  assessmentFormRef.value?.clearValidate()
}

// 关闭对话框
const handleDialogClose = () => {
  resetForm()
}

// 提交评估
const handleSubmitAssessment = async () => {
  if (!assessmentFormRef.value) return

  try {
    await assessmentFormRef.value.validate()

    submitting.value = true

    const score = getAssessmentScore()

    // TODO: API 请求
    await new Promise(resolve => setTimeout(resolve, 1000))

    if (assessmentForm.id) {
      // 更新
      const index = assessmentList.value.findIndex(a => a.id === assessmentForm.id)
      if (index > -1) {
        Object.assign(assessmentList.value[index], {
          ...assessmentForm,
          overallScore: score,
          nurseName: nurseList.value.find(n => n.id === assessmentForm.nurseId)?.name,
          elderlyName: elderlyList.value.find(e => e.id === Number(assessmentForm.elderlyId))?.name
        })
      }
      ElMessage.success('评估更新成功')
    } else {
      // 新建
      const newAssessment = {
        id: Date.now(),
        assessmentNo: `AS${Date.now()}`,
        ...assessmentForm,
        overallScore: score,
        nurseName: nurseList.value.find(n => n.id === assessmentForm.nurseId)?.name,
        elderlyName: elderlyList.value.find(e => e.id === Number(assessmentForm.elderlyId))?.name,
        status: 'pending'
      }
      assessmentList.value.unshift(newAssessment)
      ElMessage.success('评估创建成功')
    }

    assessmentDialogVisible.value = false
  } catch (error) {
    console.error('表单验证失败', error)
  } finally {
    submitting.value = false
  }
}

// 查看详情
const handleView = (row: any) => {
  currentAssessment.value = row
  viewDialogVisible.value = true
}

// 打印
const handlePrint = (row: any) => {
  ElMessage.info('打印功能开发中')
}

onMounted(() => {
  total.value = assessmentList.value.length
})
</script>

<style scoped lang="scss">
.care-assessment {
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

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;

    .header-actions {
      display: flex;
      gap: 8px;
      align-items: center;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .assessment-section {
    margin: 24px 0;
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: 8px;

    .section-title {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
      padding-bottom: 8px;
      border-bottom: 2px solid var(--border-color);
    }
  }

  .dimensions-table {
    background: var(--card-bg);

    :deep(.el-input__inner) {
      width: 80px;
    }

    :deep(.el-slider__runway) {
      margin-right: 10px;
    }
  }

  .assessment-detail {
    .detail-section {
      margin-top: 24px;

      h4 {
        font-size: 14px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 12px 0;
      }

      p {
        color: var(--text-secondary);
        line-height: 1.6;
        margin: 0;
      }
    }
  }

  :deep(.el-checkbox-group) {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;

    .el-checkbox {
      margin-right: 0;
    }
  }

  :deep(.el-radio-group) {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;

    .el-radio {
      margin-right: 0;
    }
  }
}

@media (max-width: 768px) {
  .care-assessment {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .card-header {
      flex-direction: column;
      align-items: flex-start;

      .header-actions {
        width: 100%;
        flex-direction: column;

        .el-select {
          width: 100% !important;
        }
      }
    }

    :deep(.el-dialog) {
      width: 95% !important;
    }
  }
}
</style>
