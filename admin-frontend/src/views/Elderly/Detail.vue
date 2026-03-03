<template>
  <div class="elderly-detail">
    <!-- 返回按钮 -->
    <div class="back-header">
      <el-button @click="handleBack">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
    </div>

    <!-- 老人基本信息卡片 -->
    <el-card shadow="never" class="info-card">
      <div class="elderly-profile">
        <div class="profile-avatar">
          <el-avatar :src="elderlyInfo.avatar" :size="120">
            {{ elderlyInfo.name?.charAt(0) }}
          </el-avatar>
          <el-tag
            :type="getStatusType(elderlyInfo.status)"
            size="large"
            class="status-tag"
          >
            {{ getStatusText(elderlyInfo.status) }}
          </el-tag>
        </div>

        <div class="profile-info">
          <h2 class="elderly-name">
            {{ elderlyInfo.name }}
            <el-tag :type="getCareLevelType(elderlyInfo.careLevel)" size="large">
              {{ getCareLevelText(elderlyInfo.careLevel) }}
            </el-tag>
          </h2>

          <el-descriptions :column="3" border>
            <el-descriptions-item label="性别">
              {{ elderlyInfo.gender }}
            </el-descriptions-item>
            <el-descriptions-item label="年龄">
              {{ elderlyInfo.age }} 岁
            </el-descriptions-item>
            <el-descriptions-item label="出生日期">
              {{ elderlyInfo.birthDate }}
            </el-descriptions-item>
            <el-descriptions-item label="身份证号">
              {{ elderlyInfo.idCard }}
            </el-descriptions-item>
            <el-descriptions-item label="床位号">
              {{ elderlyInfo.bedNumber }}
            </el-descriptions-item>
            <el-descriptions-item label="入住日期">
              {{ elderlyInfo.checkInDate }}
            </el-descriptions-item>
            <el-descriptions-item label="联系电话">
              {{ elderlyInfo.phone || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="家属联系人">
              {{ elderlyInfo.familyName }}
            </el-descriptions-item>
            <el-descriptions-item label="家属电话">
              {{ elderlyInfo.familyPhone }}
            </el-descriptions-item>
            <el-descriptions-item label="与老人关系">
              {{ elderlyInfo.familyRelation || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="入住天数">
              {{ elderlyInfo.stayDays }} 天
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="profile-actions">
        <el-button type="primary" @click="handleEdit">
          <el-icon><Edit /></el-icon>
          编辑信息
        </el-button>
        <el-button @click="handleHealth">
          <el-icon><Monitor /></el-icon>
          健康档案
        </el-button>
        <el-button @click="handleCare">
          <el-icon><Briefcase /></el-icon>
          护理记录
        </el-button>
        <el-button @click="handleBill">
          <el-icon><Wallet /></el-icon>
          费用账单
        </el-button>
      </div>
    </el-card>

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="detail-tabs">
      <!-- 健康状况 -->
      <el-tab-pane label="健康状况" name="health">
        <el-card shadow="never">
          <div class="health-overview">
            <div class="health-item" v-for="item in healthData" :key="item.key">
              <div class="health-icon" :class="`health-${item.level}`">
                <component :is="item.icon" />
              </div>
              <div class="health-content">
                <div class="health-label">{{ item.label }}</div>
                <div class="health-value" :class="`text-${item.level}`">
                  {{ item.value }}
                  <span v-if="item.unit" class="health-unit">{{ item.unit }}</span>
                </div>
              </div>
              <div class="health-trend" :class="item.trend > 0 ? 'up' : 'down'">
                <el-icon><component :is="item.trend > 0 ? ArrowUp : ArrowDown" /></el-icon>
                {{ Math.abs(item.trend) }}%
              </div>
            </div>
          </div>

          <!-- 健康趋势图 -->
          <div class="health-chart">
            <ECharts :option="healthChartOption" height="300px" />
          </div>
        </el-card>
      </el-tab-pane>

      <!-- 护理记录 -->
      <el-tab-pane label="护理记录" name="care">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>最近护理记录</span>
              <el-button size="small" @click="handleViewAllCare">
                查看全部
              </el-button>
            </div>
          </template>

          <el-timeline>
            <el-timeline-item
              v-for="record in careRecords"
              :key="record.id"
              :timestamp="record.time"
              placement="top"
            >
              <el-card>
                <div class="care-record">
                  <div class="care-header">
                    <el-tag :type="getRecordType(record.type)" size="small">
                      {{ record.type }}
                    </el-tag>
                    <span class="care-nurse">{{ record.nurse }}</span>
                  </div>
                  <div class="care-content">{{ record.content }}</div>
                  <div v-if="record.images" class="care-images">
                    <el-image
                      v-for="(img, index) in record.images"
                      :key="index"
                      :src="img"
                      :preview-src-list="record.images"
                      fit="cover"
                      class="care-image"
                    />
                  </div>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-tab-pane>

      <!-- 费用账单 -->
      <el-tab-pane label="费用账单" name="bill">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>最近账单</span>
              <el-button size="small" @click="handleViewAllBills">
                查看全部
              </el-button>
            </div>
          </template>

          <el-table :data="bills" stripe>
            <el-table-column prop="billNo" label="账单号" width="150" />
            <el-table-column prop="billType" label="费用类型" width="120">
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
            <el-table-column prop="dueDate" label="应付日期" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getBillStatusType(row.status)" size="small">
                  {{ getBillStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" text @click="handlePayBill(row)">
                  支付
                </el-button>
                <el-button size="small" text @click="handleViewBill(row)">
                  详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- 家属探视 -->
      <el-tab-pane label="家属探视" name="visit">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>探视记录</span>
              <el-button size="small" @click="handleAddVisit">
                添加记录
              </el-button>
            </div>
          </template>

          <el-table :data="visits" stripe>
            <el-table-column prop="visitorName" label="访客姓名" width="120" />
            <el-table-column prop="relation" label="关系" width="100" />
            <el-table-column prop="visitDate" label="访视日期" width="120" />
            <el-table-column prop="duration" label="时长(分钟)" width="100" />
            <el-table-column prop="purpose" label="来访目的" show-overflow-tooltip />
            <el-table-column prop="remark" label="备注" show-overflow-tooltip />
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- 照片相册 -->
      <el-tab-pane label="照片相册" name="photos">
        <el-card shadow="never">
          <div class="photo-gallery">
            <div
              v-for="(photo, index) in photos"
              :key="index"
              class="photo-item"
            >
              <el-image
                :src="photo.url"
                :preview-src-list="photos.map(p => p.url)"
                :initial-index="index"
                fit="cover"
                class="photo-image"
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><PictureFilled /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="photo-info">{{ photo.date }}</div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  Edit,
  Monitor,
  Briefcase,
  Wallet,
  ArrowUp,
  ArrowDown,
  PictureFilled
} from '@element-plus/icons-vue'
import ECharts from '@/components/Dashboard/ECharts.vue'
import { getChartColors } from '@/composables/useECharts'

const route = useRoute()
const router = router

const activeTab = ref('health')
const elderlyId = computed(() => route.params.id as string)

// 老人信息
const elderlyInfo = ref({
  id: '',
  name: '张奶奶',
  gender: '女',
  age: 78,
  birthDate: '1946-01-01',
  idCard: '110101194601011234',
  careLevel: 'level2',
  bedNumber: '3号楼201',
  checkInDate: '2023-06-15',
  stayDays: 268,
  phone: '138****1234',
  familyName: '张先生',
  familyPhone: '13812345678',
  familyRelation: '父子',
  status: 'active',
  avatar: ''
})

// 健康数据
const healthData = ref([
  { key: 'bloodPressure', label: '血压', value: '128/82', unit: 'mmHg', level: 'normal', trend: -2, icon: 'Monitor' },
  { key: 'heartRate', label: '心率', value: '72', unit: '次/分', level: 'normal', trend: 0, icon: 'Heart' },
  { key: 'bloodSugar', label: '血糖', value: '6.3', unit: 'mmol/L', level: 'normal', trend: 5, icon: 'Sugar' },
  { key: 'temperature', label: '体温', value: '36.5', unit: '℃', level: 'normal', trend: 0, icon: 'Temperature' },
  { key: 'weight', label: '体重', value: '58', unit: 'kg', level: 'normal', trend: -1, icon: 'Weight' },
  { key: 'oxygen', label: '血氧', value: '97', unit: '%', level: 'normal', trend: 0, icon: 'Oxygen' }
])

// 护理记录
const careRecords = ref([
  {
    id: 1,
    type: '日常护理',
    nurse: '赵护士',
    time: '2026-03-03 14:30',
    content: '测量血压、体温，协助用餐，老人状态良好',
    images: []
  },
  {
    id: 2,
    type: '康复训练',
    nurse: '陈康复师',
    time: '2026-03-03 10:00',
    content: '上肢关节活动训练30分钟，老人配合度良好',
    images: []
  },
  {
    id: 3,
    type: '健康监测',
    nurse: '周护士',
    time: '2026-03-02 16:00',
    content: '血压130/85mmHg，心率72次/分，体温36.4℃',
    images: []
  }
])

// 账单
const bills = ref([
  { billNo: 'B202603001', billType: '床位费', amount: 3500, billDate: '2026-03-01', dueDate: '2026-03-10', status: 'unpaid' },
  { billNo: 'B202603002', billType: '护理费', amount: 1800, billDate: '2026-03-01', dueDate: '2026-03-10', status: 'unpaid' },
  { billNo: 'B202603003', billType: '餐费', amount: 900, billDate: '2026-03-01', dueDate: '2026-03-10', status: 'paid' }
])

// 探视记录
const visits = ref([
  { visitorName: '张先生', relation: '儿子', visitDate: '2026-02-28', duration: 120, purpose: '看望老人，送去生活用品', remark: '' },
  { visitorName: '李女士', relation: '女儿', visitDate: '2026-02-20', duration: 90, purpose: '周末探视', remark: '' }
])

// 照片
const photos = ref([
  { url: '', date: '2026-03-01' },
  { url: '', date: '2026-02-28' },
  { url: '', date: '2026-02-25' },
  { url: '', date: '2026-02-20' },
  { url: '', date: '2026-02-15' },
  { url: '', date: '2026-02-10' }
])

// 健康趋势图
const healthChartOption = computed(() => {
  const colors = getChartColors()
  return {
    grid: {
      top: '10%',
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['收缩压', '舒张压', '血糖'],
      bottom: 0
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '收缩压',
        type: 'line',
        smooth: true,
        data: [130, 128, 132, 126, 129, 127, 128],
        itemStyle: { color: colors[0] }
      },
      {
        name: '舒张压',
        type: 'line',
        smooth: true,
        data: [85, 82, 84, 80, 83, 81, 82],
        itemStyle: { color: colors[1] }
      },
      {
        name: '血糖',
        type: 'line',
        smooth: true,
        data: [6.5, 6.3, 6.4, 6.2, 6.5, 6.1, 6.3],
        itemStyle: { color: colors[2] }
      }
    ]
  }
})

// 辅助函数
const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    active: 'success',
    leave: 'warning',
    hospital: 'danger'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '在院',
    leave: '请假',
    hospital: '住院'
  }
  return map[status] || status
}

const getCareLevelType = (level: string) => {
  const map: Record<string, any> = {
    level3: 'info',
    level2: '',
    level1: 'warning',
    special: 'danger'
  }
  return map[level] || ''
}

const getCareLevelText = (level: string) => {
  const map: Record<string, string> = {
    level3: '三级护理',
    level2: '二级护理',
    level1: '一级护理',
    special: '特级护理'
  }
  return map[level] || level
}

const getRecordType = (type: string) => {
  const map: Record<string, any> = {
    '日常护理': 'info',
    '康复训练': 'success',
    '健康监测': 'warning',
    '医疗护理': 'danger'
  }
  return map[type] || ''
}

const getBillStatusType = (status: string) => {
  const map: Record<string, any> = {
    paid: 'success',
    unpaid: 'warning',
    overdue: 'danger'
  }
  return map[status] || ''
}

const getBillStatusText = (status: string) => {
  const map: Record<string, string> = {
    paid: '已支付',
    unpaid: '未支付',
    overdue: '已逾期'
  }
  return map[status] || status
}

// 事件处理
const handleBack = () => {
  router.back()
}

const handleEdit = () => {
  console.log('编辑')
}

const handleHealth = () => {
  activeTab.value = 'health'
}

const handleCare = () => {
  activeTab.value = 'care'
}

const handleBill = () => {
  activeTab.value = 'bill'
}

const handleViewAllCare = () => {
  router.push(`/care/records?elderlyId=${elderlyId.value}`)
}

const handleViewAllBills = () => {
  router.push(`/finance/bills?elderlyId=${elderlyId.value}`)
}

const handlePayBill = (bill: any) => {
  console.log('支付账单', bill)
}

const handleViewBill = (bill: any) => {
  console.log('查看账单详情', bill)
}

const handleAddVisit = () => {
  console.log('添加探视记录')
}

onMounted(() => {
  // 加载老人详情数据
})
</script>

<style scoped lang="scss">
.elderly-detail {
  padding: 20px;

  .back-header {
    margin-bottom: 20px;
  }

  .info-card {
    margin-bottom: 20px;
  }

  .elderly-profile {
    display: flex;
    gap: 32px;
    margin-bottom: 24px;
  }

  .profile-avatar {
    position: relative;
    flex-shrink: 0;
    text-align: center;

    .status-tag {
      position: absolute;
      bottom: -8px;
      left: 50%;
      transform: translateX(-50%);
    }
  }

  .profile-info {
    flex: 1;
  }

  .elderly-name {
    font-size: 24px;
    font-weight: 600;
    margin: 0 0 16px 0;
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .profile-actions {
    display: flex;
    gap: 12px;
    padding-top: 16px;
    border-top: 1px solid var(--border-color-lighter);
  }

  .detail-tabs {
    :deep(.el-tabs__content) {
      padding-top: 20px;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .health-overview {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  .health-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: 8px;
  }

  .health-icon {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;

    &.health-normal {
      background: #f0f9ff;
      color: #67c23a;
    }

    &.health-warning {
      background: #fdf6ec;
      color: #e6a23c;
    }

    &.health-danger {
      background: #fef0f0;
      color: #f56c6c;
    }
  }

  .health-label {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .health-value {
    font-size: 20px;
    font-weight: 600;
    color: var(--text-primary);

    &.text-normal {
      color: #67c23a;
    }

    &.text-warning {
      color: #e6a23c;
    }

    &.text-danger {
      color: #f56c6c;
    }
  }

  .health-unit {
    font-size: 12px;
    font-weight: 400;
    color: var(--text-secondary);
    margin-left: 4px;
  }

  .health-trend {
    margin-left: auto;
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: 2px;

    &.up {
      color: var(--danger-color);
    }

    &.down {
      color: var(--success-color);
    }
  }

  .care-record {
    .care-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 8px;
    }

    .care-nurse {
      font-size: 12px;
      color: var(--text-secondary);
    }

    .care-content {
      color: var(--text-primary);
      margin-bottom: 8px;
    }

    .care-images {
      display: flex;
      gap: 8px;
    }

    .care-image {
      width: 60px;
      height: 60px;
      border-radius: 4px;
    }
  }

  .amount {
    font-weight: 600;
    color: var(--danger-color);
  }

  .photo-gallery {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
  }

  .photo-item {
    text-align: center;
  }

  .photo-image {
    width: 100%;
    height: 150px;
    border-radius: 8px;
    margin-bottom: 8px;
  }

  .image-error {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    background: var(--bg-tertiary);
    color: var(--text-tertiary);
    font-size: 32px;
  }

  .photo-info {
    font-size: 12px;
    color: var(--text-secondary);
  }
}

@media (max-width: 768px) {
  .elderly-detail {
    padding: 10px;

    .elderly-profile {
      flex-direction: column;
      align-items: center;
      text-align: center;
    }

    .profile-actions {
      flex-wrap: wrap;
      justify-content: center;
    }

    .health-overview {
      grid-template-columns: 1fr;
    }
  }
}
</style>
