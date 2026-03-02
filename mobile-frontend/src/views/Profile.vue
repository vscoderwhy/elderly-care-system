<template>
  <div class="profile-page">
    <van-nav-bar title="我的" />

    <div class="header">
      <van-image
        round
        width="60"
        height="60"
        :src="userStore.userInfo?.avatar || 'https://fastly.jsdelivr.net/npm/@vant/assets/user-active.png'"
      />
      <div class="info">
        <div class="name">{{ userStore.userInfo?.nickname || '用户' }}</div>
        <div class="phone">{{ userStore.userInfo?.phone }}</div>
      </div>
    </div>

    <van-cell-group inset>
      <van-cell title="切换账号" icon="exchange" is-link @click="handleLogout" />
      <van-cell title="关于我们" icon="info-o" is-link />
    </van-cell-group>

    <div class="logout-btn">
      <van-button type="danger" block round @click="handleLogout">
        退出登录
      </van-button>
    </div>

    <div class="version">
      <span>版本 1.0.0</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { showConfirmDialog } from 'vant'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

function handleLogout() {
  showConfirmDialog({
    title: '确认退出',
    message: '确定要退出登录吗？'
  }).then(() => {
    userStore.logout()
    router.replace('/login')
  }).catch(() => {
    // 取消
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.header .name {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.header .phone {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4px;
}

.van-cell-group {
  margin: 15px;
}

.logout-btn {
  padding: 0 15px;
  margin-top: 30px;
}

.version {
  text-align: center;
  padding: 30px 0;
  font-size: 12px;
  color: #999;
}
</style>
