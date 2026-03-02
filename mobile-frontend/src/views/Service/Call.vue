<template>
  <div class="service-call-page">
    <van-nav-bar title="服务呼叫" left-arrow @click-left="$router.back()" />

    <div class="call-types">
      <div class="type-title">请选择呼叫类型</div>

      <van-grid :column-num="2" :gutter="12">
        <van-grid-item
          v-for="item in callTypes"
          :key="item.type"
          :icon="item.icon"
          :text="item.label"
          @click="selectedType = item.type"
          :class="{ active: selectedType === item.type }"
        />
      </van-grid>
    </div>

    <div class="notes-section">
      <van-field
        v-model="notes"
        rows="2"
        autosize
        label="备注"
        type="textarea"
        placeholder="请输入需要说明的情况（选填）"
      />
    </div>

    <div class="submit-btn">
      <van-button
        type="primary"
        block
        round
        size="large"
        :loading="loading"
        :disabled="!selectedType"
        @click="submitCall"
      >
        立即呼叫
      </van-button>
    </div>

    <!-- 呼叫记录 -->
    <div class="history-section">
      <van-cell-group title="最近呼叫记录">
        <van-empty v-if="callHistory.length === 0" description="暂无呼叫记录" />
        <van-cell
          v-for="record in callHistory"
          :key="record.id"
          :title="record.type"
          :label="formatTime(record.requested_at)"
        >
          <template #value>
            <van-tag :type="getStatusType(record.status)">
              {{ getStatusText(record.status) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showSuccessToast } from 'vant'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

const elderlyId = ref(Number(route.params.elderlyId) || 0)
const selectedType = ref('')
const notes = ref('')
const loading = ref(false)
const callHistory = ref<any[]>([])

const callTypes = [
  { type: '护理', label: '护理协助', icon: 'user-o' },
  { type: '送餐', label: '送餐服务', icon: 'food' },
  { type: '打扫', label: '房间打扫', icon: 'clean' },
  { type: '紧急', label: '紧急呼叫', icon: 'warning-o' },
  { type: '医疗', label: '医疗服务', icon: 'medical' },
  { type: '其他', label: '其他服务', icon: 'ellipsis' }
]

const token = localStorage.getItem('token')

const submitCall = async () => {
  if (!selectedType.value) {
    showToast('请选择呼叫类型')
    return
  }

  loading.value = true
  try {
    await axios.post('http://1.12.223.138:8080/api/service/requests', {
      elderly_id: elderlyId.value,
      type: selectedType.value,
      notes: notes.value
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })

    showSuccessToast('呼叫已发送，请耐心等待')
    notes.value = ''
    selectedType.value = ''
    loadHistory()
  } catch (error: any) {
    showToast(error.response?.data?.message || '发送失败')
  } finally {
    loading.value = false
  }
}

const loadHistory = async () => {
  try {
    const res = await axios.get('http://1.12.223.138:8080/api/service/requests', {
      headers: { Authorization: `Bearer ${token}` },
      params: { page: 1, page_size: 10 }
    })
    callHistory.value = res.data.list || []
  } catch (error) {
    console.error('加载历史记录失败', error)
  }
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success'
  }
  return map[status] || 'default'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成'
  }
  return map[status] || status
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.service-call-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.call-types {
  background: #fff;
  padding: 16px;
  margin-bottom: 12px;
}

.type-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 12px;
}

.call-types :deep(.van-grid-item.active) {
  background: #e8f4ff;
}

.call-types :deep(.van-grid-item.active .van-grid-item__text) {
  color: #1989fa;
}

.notes-section {
  margin-bottom: 12px;
}

.submit-btn {
  padding: 24px 16px;
}

.history-section {
  margin-top: 12px;
}
</style>
