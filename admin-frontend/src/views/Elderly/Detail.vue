<template>
  <div class="elderly-detail">
    <el-page-header @back="goBack" title="返回" content="老人详情" />

    <el-card style="margin-top: 20px" v-loading="loading">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="姓名">{{ elderly.name }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ elderly.gender }}</el-descriptions-item>
        <el-descriptions-item label="身份证号">{{ elderly.id_card }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ elderly.phone }}</el-descriptions-item>
        <el-descriptions-item label="紧急联系人">{{ elderly.emergency_contact }}</el-descriptions-item>
        <el-descriptions-item label="紧急联系电话">{{ elderly.emergency_phone }}</el-descriptions-item>
        <el-descriptions-item label="护理等级">{{ careLevelText }}</el-descriptions-item>
        <el-descriptions-item label="床位">{{ elderly.bed?.name || '未分配' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 健康数据概览 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>健康数据</span>
          <el-button type="primary" size="small" @click="showAddHealthDialog">添加记录</el-button>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :span="4" v-for="record in latestHealth" :key="record.record_type">
          <div class="health-card" @click="showHealthHistory(record.record_type)">
            <div class="health-icon">
              <el-icon size="24"><component :is="getHealthIcon(record.record_type)" /></el-icon>
            </div>
            <div class="health-value">{{ formatHealthValue(record) }}</div>
            <div class="health-label">{{ getHealthLabel(record.record_type) }}</div>
          </div>
        </el-col>
        <el-col :span="4" v-for="type in missingTypes" :key="type">
          <div class="health-card empty" @click="showAddHealthDialog(type)">
            <div class="health-icon">
              <el-icon size="24"><Plus /></el-icon>
            </div>
            <div class="health-label">{{ getHealthLabel(type) }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 护理记录 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>护理记录</span>
      </template>
      <el-table :data="careRecords" v-loading="careLoading" size="small">
        <el-table-column prop="care_item.name" label="护理项目" width="150" />
        <el-table-column prop="staff.nickname" label="护理员" width="100" />
        <el-table-column prop="notes" label="备注" />
        <el-table-column prop="recorded_at" label="时间" width="180" />
      </el-table>
    </el-card>

    <!-- 添加健康记录对话框 -->
    <el-dialog v-model="healthDialogVisible" title="添加健康记录" width="500px">
      <el-form :model="healthForm" label-width="100px">
        <el-form-item label="记录类型" required>
          <el-select v-model="healthForm.record_type" placeholder="请选择" style="width: 100%">
            <el-option label="血压" value="blood_pressure" />
            <el-option label="血糖" value="blood_sugar" />
            <el-option label="体温" value="temperature" />
            <el-option label="体重" value="weight" />
            <el-option label="心率" value="heart_rate" />
          </el-select>
        </el-form-item>
        <el-form-item label="数值" required>
          <div style="display: flex; gap: 10px; align-items: center;">
            <el-input v-model="healthForm.value" placeholder="请输入数值" style="width: 120px;" />
            <el-input v-model="healthForm.value2" placeholder="低压(血压)" style="width: 120px;" v-if="healthForm.record_type === 'blood_pressure'" />
            <span>{{ healthForm.unit || getUnit(healthForm.record_type) }}</span>
          </div>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="healthForm.notes" type="textarea" :rows="2" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="healthDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitHealthRecord" :loading="healthSubmitting">提交</el-button>
      </template>
    </el-dialog>

    <!-- 健康记录历史对话框 -->
    <el-dialog v-model="historyDialogVisible" :title="getHealthLabel(selectedType) + '历史记录'" width="700px">
      <el-table :data="healthHistory" v-loading="historyLoading" size="small">
        <el-table-column label="数值">
          <template #default="{ row }">
            {{ formatHealthValue(row) }}
          </template>
        </el-table-column>
        <el-table-column prop="notes" label="备注" />
        <el-table-column prop="recorded_at" label="记录时间" width="180" />
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button type="danger" size="small" link @click="deleteHealthRecord(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Heart, Timer, Sunny, Scale, DataLine } from '@element-plus/icons-vue'
import instance, { elderlyApi } from '@/api'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const careLoading = ref(false)
const elderly: any = ref({})
const careRecords = ref([])
const latestHealth = ref<any[]>([])

const healthDialogVisible = ref(false)
const healthSubmitting = ref(false)
const healthForm = ref({
  record_type: '',
  value: '',
  value2: '',
  unit: '',
  notes: ''
})

const historyDialogVisible = ref(false)
const historyLoading = ref(false)
const healthHistory = ref([])
const selectedType = ref('')

const careLevelText = computed(() => {
  const levels = ['', '一级护理', '二级护理', '三级护理']
  return levels[elderly.value.care_level] || '-'
})

const allTypes = ['blood_pressure', 'blood_sugar', 'temperature', 'weight', 'heart_rate']
const missingTypes = computed(() => {
  const existing = latestHealth.value.map(r => r.record_type)
  return allTypes.filter(t => !existing.includes(t))
})

const getHealthIcon = (type: string) => {
  const icons: Record<string, any> = {
    blood_pressure: Heart,
    blood_sugar: Timer,
    temperature: Sunny,
    weight: Scale,
    heart_rate: DataLine
  }
  return icons[type] || Plus
}

const getHealthLabel = (type: string) => {
  const labels: Record<string, string> = {
    blood_pressure: '血压',
    blood_sugar: '血糖',
    temperature: '体温',
    weight: '体重',
    heart_rate: '心率'
  }
  return labels[type] || type
}

const getUnit = (type: string) => {
  const units: Record<string, string> = {
    blood_pressure: 'mmHg',
    blood_sugar: 'mmol/L',
    temperature: '℃',
    weight: 'kg',
    heart_rate: '次/分'
  }
  return units[type] || ''
}

const formatHealthValue = (record: any) => {
  if (record.record_type === 'blood_pressure' && record.value2) {
    return `${record.value}/${record.value2}`
  }
  return record.value
}

const fetchData = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    elderly.value = await elderlyApi.get(id)
    loadCareRecords(id)
    loadLatestHealth(id)
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadCareRecords = async (elderlyId: number) => {
  careLoading.value = true
  try {
    const res = await instance.get('/care/records', {
      params: { elderly_id: elderlyId, page: 1, page_size: 10 }
    })
    careRecords.value = res.list || []
  } catch (error) {
    console.error(error)
  } finally {
    careLoading.value = false
  }
}

const loadLatestHealth = async (elderlyId: number) => {
  try {
    const res = await instance.get(`/health/records/latest/${elderlyId}`)
    latestHealth.value = res || []
  } catch (error) {
    console.error(error)
  }
}

const showAddHealthDialog = (type?: string) => {
  healthForm.value = {
    record_type: type || '',
    value: '',
    value2: '',
    unit: '',
    notes: ''
  }
  healthDialogVisible.value = true
}

const submitHealthRecord = async () => {
  if (!healthForm.value.record_type || !healthForm.value.value) {
    ElMessage.warning('请填写完整信息')
    return
  }

  healthSubmitting.value = true
  try {
    await instance.post('/health/records', {
      elderly_id: Number(route.params.id),
      record_type: healthForm.value.record_type,
      value: healthForm.value.value,
      value2: healthForm.value.value2,
      unit: getUnit(healthForm.value.record_type),
      notes: healthForm.value.notes
    })
    ElMessage.success('记录成功')
    healthDialogVisible.value = false
    loadLatestHealth(Number(route.params.id))
  } catch (error) {
    ElMessage.error('记录失败')
  } finally {
    healthSubmitting.value = false
  }
}

const showHealthHistory = async (type: string) => {
  selectedType.value = type
  historyDialogVisible.value = true
  historyLoading.value = true
  try {
    const res = await instance.get('/health/records', {
      params: {
        elderly_id: Number(route.params.id),
        record_type: type,
        page: 1,
        page_size: 20
      }
    })
    healthHistory.value = res.list || []
  } catch (error) {
    console.error(error)
  } finally {
    historyLoading.value = false
  }
}

const deleteHealthRecord = async (id: number) => {
  try {
    await instance.delete(`/health/records/${id}`)
    ElMessage.success('删除成功')
    showHealthHistory(selectedType.value)
    loadLatestHealth(Number(route.params.id))
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const goBack = () => {
  router.push('/elderly')
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.elderly-detail {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.health-card {
  text-align: center;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.health-card:hover {
  background: #e6f0fa;
  transform: translateY(-2px);
}

.health-card.empty {
  opacity: 0.6;
}

.health-icon {
  color: #409EFF;
  margin-bottom: 8px;
}

.health-value {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.health-label {
  font-size: 12px;
  color: #909399;
}
</style>
