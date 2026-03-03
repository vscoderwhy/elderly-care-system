<template>
  <view class="login-page">
    <!-- 背景装饰 -->
    <view class="bg-decoration">
      <view class="bg-circle circle-1"></view>
      <view class="bg-circle circle-2"></view>
      <view class="bg-circle circle-3"></view>
    </view>

    <!-- Logo区域 -->
    <view class="logo-section">
      <image class="logo-image" src="/static/logo.png" mode="aspectFit" />
      <text class="app-name">智慧养老</text>
      <text class="app-slogan">用心守护，让爱传递</text>
    </view>

    <!-- 登录表单 -->
    <view class="login-form">
      <view class="tab-bar">
        <view
          class="tab-item"
          :class="{ active: loginType === 'wechat' }"
          @click="loginType = 'wechat'"
        >
          微信登录
        </view>
        <view
          class="tab-item"
          :class="{ active: loginType === 'phone' }"
          @click="loginType = 'phone'"
        >
          手机登录
        </view>
      </view>

      <!-- 微信登录 -->
      <view v-if="loginType === 'wechat'" class="wechat-login">
        <button class="wechat-btn" open-type="getUserInfo" @getuserinfo="handleWechatLogin">
          <uni-icons type="weixin" size="24" color="#fff" />
          <text>微信一键登录</text>
        </button>
        <view class="login-tip">
          <text>未注册的微信号将自动创建账号</text>
        </view>
      </view>

      <!-- 手机登录 -->
      <view v-else class="phone-login">
        <view class="form-item">
          <view class="input-wrapper">
            <uni-icons type="phone" size="20" color="#909399" />
            <input
              class="form-input"
              type="number"
              maxlength="11"
              v-model="phoneForm.phone"
              placeholder="请输入手机号"
              placeholder-class="input-placeholder"
            />
          </view>
        </view>

        <view class="form-item">
          <view class="input-wrapper">
            <uni-icons type="locked" size="20" color="#909399" />
            <input
              class="form-input"
              type="number"
              maxlength="6"
              v-model="phoneForm.code"
              placeholder="请输入验证码"
              placeholder-class="input-placeholder"
            />
            <button
              class="code-btn"
              :disabled="counting"
              @click="handleSendCode"
            >
              {{ counting ? `${countdown}s` : '获取验证码' }}
            </button>
          </view>
        </view>

        <button class="login-btn" :disabled="logging" @click="handlePhoneLogin">
          {{ logging ? '登录中...' : '登录' }}
        </button>

        <view class="login-agreement">
          <checkbox-group @change="handleAgreementChange">
            <label class="agreement-label">
              <checkbox value="agree" :checked="agreed" color="#409eff" />
              <text class="agreement-text">
                我已阅读并同意
                <text class="link" @click.stop="handleViewAgreement">《用户协议》</text>
                和
                <text class="link" @click.stop="handleViewPrivacy">《隐私政策》</text>
              </text>
            </label>
          </checkbox-group>
        </view>
      </view>
    </view>

    <!-- 底部信息 -->
    <view class="footer-info">
      <text class="copyright">© 2026 智慧养老 All Rights Reserved</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const loginType = ref<'wechat' | 'phone'>('wechat')
const counting = ref(false)
const countdown = ref(60)
const logging = ref(false)
const agreed = ref(false)

const phoneForm = ref({
  phone: '',
  code: ''
})

const handleWechatLogin = async (e: any) => {
  if (e.detail.userInfo) {
    logging.value = true

    // TODO: 调用后端微信登录接口
    try {
      // 模拟登录
      await new Promise(resolve => setTimeout(resolve, 1000))

      uni.showToast({
        title: '登录成功',
        icon: 'success'
      })

      // 跳转到首页
      setTimeout(() => {
        uni.switchTab({
          url: '/pages/index/index'
        })
      }, 1500)
    } catch (error) {
      uni.showToast({
        title: '登录失败',
        icon: 'none'
      })
    } finally {
      logging.value = false
    }
  }
}

const handleSendCode = () => {
  if (!/^1[3-9]\d{9}$/.test(phoneForm.value.phone)) {
    uni.showToast({
      title: '请输入正确的手机号',
      icon: 'none'
    })
    return
  }

  counting.value = true
  countdown.value = 60

  // TODO: 调用发送验证码接口

  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
      counting.value = false
    }
  }, 1000)

  uni.showToast({
    title: '验证码已发送',
    icon: 'success'
  })
}

const handlePhoneLogin = async () => {
  if (!agreed.value) {
    uni.showToast({
      title: '请先同意用户协议',
      icon: 'none'
    })
    return
  }

  if (!/^1[3-9]\d{9}$/.test(phoneForm.value.phone)) {
    uni.showToast({
      title: '请输入正确的手机号',
      icon: 'none'
    })
    return
  }

  if (phoneForm.value.code.length !== 6) {
    uni.showToast({
      title: '请输入验证码',
      icon: 'none'
    })
    return
  }

  logging.value = true

  try {
    // TODO: 调用后端手机登录接口
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 保存登录状态
    await userStore.login({
      id: '1',
      name: '家属用户',
      phone: phoneForm.value.phone,
      avatar: '',
      role: 'family',
      token: 'mock_token_123456'
    })

    uni.showToast({
      title: '登录成功',
      icon: 'success'
    })

    setTimeout(() => {
      uni.switchTab({
        url: '/pages/index/index'
      })
    }, 1500)
  } catch (error) {
    uni.showToast({
      title: '登录失败',
      icon: 'none'
    })
  } finally {
    logging.value = false
  }
}

const handleAgreementChange = (e: any) => {
  agreed.value = e.detail.value.includes('agree')
}

const handleViewAgreement = () => {
  uni.navigateTo({
    url: '/pages/common/agreement'
  })
}

const handleViewPrivacy = () => {
  uni.navigateTo({
    url: '/pages/common/privacy'
  })
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.bg-decoration {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  z-index: 0;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.circle-1 {
  width: 400rpx;
  height: 400rpx;
  top: -100rpx;
  right: -100rpx;
}

.circle-2 {
  width: 300rpx;
  height: 300rpx;
  bottom: 200rpx;
  left: -80rpx;
}

.circle-3 {
  width: 200rpx;
  height: 200rpx;
  bottom: -50rpx;
  right: 100rpx;
}

.logo-section {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 0 80rpx;
}

.logo-image {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 30rpx;
}

.app-name {
  font-size: 48rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 16rpx;
}

.app-slogan {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.login-form {
  position: relative;
  z-index: 1;
  padding: 0 60rpx;
}

.tab-bar {
  display: flex;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50rpx;
  padding: 8rpx;
  margin-bottom: 40rpx;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 20rpx 0;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.6);
  border-radius: 40rpx;
  transition: all 0.3s;

  &.active {
    background: #fff;
    color: #667eea;
    font-weight: 500;
  }
}

.wechat-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 30rpx;
}

.wechat-btn {
  width: 100%;
  height: 88rpx;
  background: #07c160;
  color: #fff;
  border: none;
  border-radius: 44rpx;
  font-size: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
}

.login-tip {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
}

.phone-login {
  display: flex;
  flex-direction: column;
  gap: 30rpx;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 0 24rpx;
  height: 88rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 44rpx;
}

.form-input {
  flex: 1;
  font-size: 28rpx;
  color: #fff;
}

.input-placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.code-btn {
  padding: 0 24rpx;
  height: 60rpx;
  background: #fff;
  color: #667eea;
  border: none;
  border-radius: 30rpx;
  font-size: 24rpx;
  line-height: 60rpx;

  &:disabled {
    opacity: 0.6;
  }
}

.login-btn {
  width: 100%;
  height: 88rpx;
  background: #fff;
  color: #667eea;
  border: none;
  border-radius: 44rpx;
  font-size: 32rpx;
  font-weight: 500;
  margin-top: 20rpx;

  &:disabled {
    opacity: 0.6;
  }
}

.login-agreement {
  margin-top: 20rpx;
}

.agreement-label {
  display: flex;
  align-items: flex-start;
  gap: 8rpx;
}

.agreement-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.5;
}

.link {
  color: #fff;
  text-decoration: underline;
}

.footer-info {
  position: absolute;
  bottom: 40rpx;
  left: 0;
  right: 0;
  text-align: center;
  z-index: 1;
}

.copyright {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
}
</style>
