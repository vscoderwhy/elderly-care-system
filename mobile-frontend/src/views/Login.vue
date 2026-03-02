<template>
  <div class="login-page">
    <div class="logo-section">
      <h1 class="title">养老院服务</h1>
      <p class="subtitle">家属端</p>
    </div>

    <van-form @submit="handleLogin" class="login-form">
      <van-cell-group inset>
        <van-field
          v-model="phone"
          name="phone"
          label="手机号"
          placeholder="请输入手机号"
          :rules="[{ required: true, message: '请输入手机号' }]"
          type="tel"
          maxlength="11"
        />
        <van-field
          v-model="password"
          name="password"
          label="密码"
          placeholder="请输入密码"
          :rules="[{ required: true, message: '请输入密码' }]"
          type="password"
        />
      </van-cell-group>

      <div class="button-group">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          登录
        </van-button>
      </div>
    </van-form>

    <div class="tips">
      <p>测试账号: 13800138000</p>
      <p>密码: 123456</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import api from '@/api'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const phone = ref('13800138000')
const password = ref('123456')
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  try {
    const result = await api.auth.login(phone.value, password.value)
    userStore.setToken(result.token)
    userStore.setUserInfo(result.user)
    showToast({ type: 'success', message: '登录成功' })
    setTimeout(() => {
      router.replace('/home')
    }, 500)
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 20px 20px;
  display: flex;
  flex-direction: column;
}

.logo-section {
  text-align: center;
  margin-bottom: 60px;
}

.title {
  font-size: 32px;
  color: #fff;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.login-form {
  margin-bottom: 30px;
}

.button-group {
  margin: 30px 16px;
}

.tips {
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  line-height: 1.8;
}
</style>
