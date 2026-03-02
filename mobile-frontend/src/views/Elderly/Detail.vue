<template>
  <div class="elderly-detail-page">
    <van-nav-bar title="老人详情" left-arrow @click-left="router.back()" />

    <div class="basic-info">
      <div class="avatar-section">
        <van-image
          round
          width="80"
          height="80"
          :src="elderly.avatar || 'https://fastly.jsdelivr.net/npm/@vant/assets/user-active.png'"
        />
        <div class="name">{{ elderly.name }}</div>
        <div class="tags">
          <van-tag>{{ elderly.gender }}</van-tag>
          <van-tag type="primary">{{ elderly.care_level }}级护理</van-tag>
        </div>
      </div>
    </div>

    <van-cell-group title="基本信息" inset>
      <van-cell title="床位" :value="elderly.bed?.name || '未分配'" />
      <van-cell title="联系电话" :value="elderly.phone" />
      <van-cell title="紧急联系人" :value="elderly.emergency_contact" />
    </van-cell-group>

    <!-- 服务呼叫按钮 -->
    <div class="call-button-section">
      <van-button
        type="danger"
        block
        round
        icon="phone-o"
        @click="router.push(`/service/call/${elderlyId}`)"
      >
        服务呼叫
      </van-button>
    </div>

    <van-tabs v-model:active="activeTab" sticky offset-top="46">
      <van-tab title="护理记录" name="care">
        <van-list
          v-model:loading="careLoading"
          :finished="careFinished"
          @load="loadCareRecords"
        >
          <div class="care-list">
            <div class="care-item" v-for="item in careRecords" :key="item.id">
              <div class="care-header">
                <van-tag type="primary">{{ item.care_item?.name }}</van-tag>
                <span class="time">{{ formatTime(item.recorded_at) }}</span>
              </div>
              <div class="care-content">
                <div class="staff">护理员: {{ item.staff?.nickname }}</div>
                <div class="notes" v-if="item.notes">{{ item.notes }}</div>
              </div>
            </div>
          </div>
        </van-list>
        <van-empty v-if="careRecords.length === 0 && !careLoading" description="暂无护理记录" />
      </van-tab>

      <van-tab title="健康数据" name="health">
        <div class="health-list">
          <div class="health-item">
            <van-icon name="activity" color="#f56c6c" />
            <div class="info">
              <div class="label">血压</div>
              <div class="value">120/80 mmHg</div>
            </div>
          </div>
          <div class="health-item">
            <van-icon name="flower-o" color="#67c23a" />
            <div class="info">
              <div class="label">血糖</div>
              <div class="value">5.6 mmol/L</div>
            </div>
          </div>
          <div class="health-item">
            <van-icon name="fire-o" color="#e6a23c" />
            <div class="info">
              <div class="label">体温</div>
              <div class="value">36.5℃</div>
            </div>
          </div>
        </div>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import api from '@/api'

const router = useRouter()
const route = useRoute()

const elderly = ref<any>({})
const activeTab = ref('care')
const careRecords = ref<any[]>([])
const careLoading = ref(false)
const careFinished = ref(false)

const elderlyId = Number(route.params.id)

onMounted(() => {
  loadElderly()
})

function loadElderly() {
  showToast({ message: '加载中...', forbidClick: true })
  api.elderly.get(elderlyId).then((data: any) => {
    elderly.value = data
    showToast.clear()
  }).catch(() => {
    showToast.clear()
  })
}

function loadCareRecords() {
  careLoading.value = true
  api.care.records({ elderly_id: elderlyId, page: 1, page_size: 20 }).then((data: any) => {
    careRecords.value = data.list || []
    careLoading.value = false
    careFinished.value = true
  }).catch(() => {
    careLoading.value = false
  })
}

function formatTime(time: string) {
  const date = new Date(time)
  return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}
</script>

<style scoped>
.elderly-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.basic-info {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
}

.avatar-section {
  text-align: center;
}

.name {
  font-size: 22px;
  font-weight: 600;
  color: #fff;
  margin-top: 15px;
}

.tags {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 10px;
}

.van-cell-group {
  margin: 15px 0;
}

.call-button-section {
  padding: 0 16px 16px;
}

.care-list {
  padding: 15px;
}

.care-item {
  background: #fff;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
}

.care-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.time {
  font-size: 12px;
  color: #999;
}

.care-content .staff {
  font-size: 14px;
  color: #666;
  margin-bottom: 6px;
}

.care-content .notes {
  font-size: 13px;
  color: #999;
}

.health-list {
  padding: 15px;
}

.health-item {
  display: flex;
  align-items: center;
  gap: 15px;
  background: #fff;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 12px;
}

.health-item .info {
  flex: 1;
}

.health-item .label {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.health-item .value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}
</style>
