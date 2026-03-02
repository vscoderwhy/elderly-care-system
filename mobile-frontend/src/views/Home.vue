<template>
  <div class="home-page">
    <!-- 头部用户信息 -->
    <div class="header">
      <div class="user-info">
        <van-image
          round
          width="50"
          height="50"
          :src="userStore.userInfo?.avatar || 'https://fastly.jsdelivr.net/npm/@vant/assets/user-active.png'"
        />
        <div class="user-detail">
          <div class="name">{{ userStore.userInfo?.nickname || '用户' }}</div>
          <div class="role">家属</div>
        </div>
      </div>
    </div>

    <!-- 关联老人 -->
    <div class="section" v-if="elderlyList.length > 0">
      <div class="section-title">我的家人</div>
      <div class="elderly-cards">
        <div
          class="elderly-card"
          v-for="item in elderlyList"
          :key="item.id"
          @click="goToElderly(item.id)"
        >
          <div class="elderly-info">
            <van-icon name="user-o" size="24" />
            <div class="detail">
              <div class="name">{{ item.name }}</div>
              <div class="room">{{ item.bed?.name || '未分配床位' }}</div>
            </div>
          </div>
          <van-icon name="arrow" />
        </div>
      </div>
    </div>

    <!-- 快捷功能 -->
    <div class="section">
      <div class="section-title">快捷功能</div>
      <div class="quick-actions">
        <div class="action-item" @click="router.push('/care')">
          <van-icon name="todo-list-o" size="32" color="#667eea" />
          <span>护理记录</span>
        </div>
        <div class="action-item" @click="router.push('/bills')">
          <van-icon name="balance-pay-o" size="32" color="#764ba2" />
          <span>费用账单</span>
        </div>
        <div class="action-item" @click="makeCall">
          <van-icon name="phone-o" size="32" color="#f56c6c" />
          <span>联系客服</span>
        </div>
      </div>
    </div>

    <!-- 最新护理 -->
    <div class="section" v-if="recentCare.length > 0">
      <div class="section-title">最新护理</div>
      <div class="care-list">
        <div class="care-item" v-for="item in recentCare" :key="item.id">
          <div class="care-time">{{ formatTime(item.recorded_at) }}</div>
          <div class="care-content">
            <div class="care-name">{{ item.care_item?.name }}</div>
            <div class="care-staff">{{ item.staff?.nickname }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showLoadingToast, closeToast, showEmptyImage } from 'vant'
import api from '@/api'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const elderlyList = ref<any[]>([])
const recentCare = ref<any[]>([])

onMounted(() => {
  loadData()
})

function loadData() {
  showLoadingToast({ message: '加载中...', forbidClick: true })
  api.elderly.familyList().then((data: any) => {
    elderlyList.value = data || []
    if (data && data.length > 0) {
      loadRecentCare(data[0].id)
    }
    closeToast()
  }).catch(() => {
    closeToast()
  })
}

function loadRecentCare(elderlyId: number) {
  api.care.records({ elderly_id: elderlyId, page: 1, page_size: 5 }).then((data: any) => {
    recentCare.value = data.list || []
  })
}

function goToElderly(id: number) {
  router.push(`/elderly/${id}`)
}

function formatTime(time: string) {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const hours = Math.floor(diff / (1000 * 60 * 60))
  if (hours < 1) {
    return '刚刚'
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    return `${Math.floor(hours / 24)}天前`
  }
}

function makeCall() {
  window.location.href = 'tel:400-123-4567'
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-detail .name {
  font-size: 18px;
  color: #fff;
  font-weight: 600;
}

.user-detail .role {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4px;
}

.section {
  margin: 15px;
  background: #fff;
  border-radius: 12px;
  padding: 15px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.elderly-cards {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.elderly-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.elderly-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.elderly-info .name {
  font-size: 15px;
  color: #333;
  font-weight: 500;
}

.elderly-info .room {
  font-size: 13px;
  color: #999;
  margin-top: 2px;
}

.quick-actions {
  display: flex;
  justify-content: space-around;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.action-item span {
  font-size: 13px;
  color: #666;
}

.care-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.care-item {
  padding: 12px 0;
  border-bottom: 1px solid #eee;
}

.care-item:last-child {
  border-bottom: none;
}

.care-time {
  font-size: 12px;
  color: #999;
  margin-bottom: 6px;
}

.care-name {
  font-size: 14px;
  color: #333;
  margin-bottom: 4px;
}

.care-staff {
  font-size: 12px;
  color: #999;
}
</style>
