<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="login-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- Logo 和标题 -->
      <div class="login-header">
        <img src="/logo.svg" alt="Logo" class="login-logo" />
        <h1 class="login-title">养老院管理系统</h1>
        <p class="login-subtitle">Elderly Care Management System</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            size="large"
            placeholder="请输入用户名"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            size="large"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item v-if="showCaptcha">
          <div class="captcha-wrapper">
            <el-input
              v-model="loginForm.captcha"
              size="large"
              placeholder="请输入验证码"
              :prefix-icon="Key"
              clearable
            />
            <div class="captcha-img" @click="refreshCaptcha">
              <img :src="captchaUrl" alt="验证码" />
            </div>
          </div>
        </el-form-item>

        <el-form-item>
          <div class="login-options">
            <el-checkbox v-model="loginForm.rememberMe">记住密码</el-checkbox>
            <el-link type="primary" :underline="false">忘记密码？</el-link>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 底部链接 -->
      <div class="login-footer">
        <span>还没有账号？</span>
        <el-link type="primary" :underline="false">联系管理员</el-link>
      </div>
    </div>

    <!-- 系统信息 -->
    <div class="system-info">
      <p>© 2026 养老院管理系统 v1.0.0</p>
      <p>技术支持：开发团队</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock, Key } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

// 表单引用
const loginFormRef = ref<FormInstance>()

// 加载状态
const loading = ref(false)
const showCaptcha = ref(false)
const captchaUrl = ref('')

// 登录表单
const loginForm = reactive({
  username: '',
  password: '',
  captcha: '',
  rememberMe: false
})

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 4, message: '验证码为 4 位', trigger: 'blur' }
  ]
}

// 刷新验证码
const refreshCaptcha = () => {
  captchaUrl.value = `/api/captcha?t=${Date.now()}`
}

// 登录处理
const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    const valid = await loginFormRef.value.validate()
    if (!valid) return

    loading.value = true

    // 模拟登录请求
    await new Promise((resolve) => setTimeout(resolve, 1000))

    // TODO: 实际登录请求
    // const res = await loginApi({
    //   username: loginForm.username,
    //   password: loginForm.password,
    //   captcha: loginForm.captcha,
    //   rememberMe: loginForm.rememberMe
    // })

    // 模拟登录成功
    const token = 'mock-token-' + Date.now()
    const userInfo = {
      id: 1,
      username: loginForm.username,
      name: '管理员',
      avatar: '',
      role: 'admin',
      permissions: ['*']
    }

    // 保存登录信息
    localStorage.setItem('token', token)
    localStorage.setItem('userInfo', JSON.stringify(userInfo))

    // 记住密码
    if (loginForm.rememberMe) {
      localStorage.setItem('savedUsername', loginForm.username)
    } else {
      localStorage.removeItem('savedUsername')
    }

    ElMessage.success('登录成功')

    // 跳转到首页或重定向页面
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (error) {
    console.error('登录失败', error)
    if (showCaptcha.value) {
      refreshCaptcha()
      loginForm.captcha = ''
    }
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  // 检查是否有保存的用户名
  const savedUsername = localStorage.getItem('savedUsername')
  if (savedUsername) {
    loginForm.username = savedUsername
    loginForm.rememberMe = true
  }

  // 检查是否需要验证码（登录失败3次后显示）
  const failCount = sessionStorage.getItem('loginFailCount') || '0'
  if (parseInt(failCount) >= 3) {
    showCaptcha.value = true
    refreshCaptcha()
  }
})
</script>

<style scoped lang="scss">
.login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;

  .bg-shape {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    animation: float 6s ease-in-out infinite;

    &.shape-1 {
      width: 300px;
      height: 300px;
      top: -150px;
      left: -150px;
      animation-delay: 0s;
    }

    &.shape-2 {
      width: 200px;
      height: 200px;
      top: 50%;
      right: -100px;
      animation-delay: 2s;
    }

    &.shape-3 {
      width: 150px;
      height: 150px;
      bottom: -75px;
      left: 30%;
      animation-delay: 4s;
    }
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

.login-card {
  position: relative;
  z-index: 1;
  width: 420px;
  max-width: 90%;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

.login-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.login-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.login-form {
  margin-top: 32px;
}

.captcha-wrapper {
  display: flex;
  gap: 12px;

  .el-input {
    flex: 1;
  }

  .captcha-img {
    width: 120px;
    height: 40px;
    cursor: pointer;
    border-radius: 4px;
    overflow: hidden;
    background: #f0f2f5;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    &:hover {
      opacity: 0.8;
    }
  }
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.login-footer {
  margin-top: 24px;
  text-align: center;
  font-size: 14px;
  color: #606266;

  .el-link {
    margin-left: 8px;
  }
}

.system-info {
  position: absolute;
  bottom: 20px;
  left: 0;
  right: 0;
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 12px;
  line-height: 1.8;

  p {
    margin: 0;
  }
}

// 响应式适配
@media (max-width: 768px) {
  .login-card {
    width: 90%;
    padding: 24px;
  }

  .login-title {
    font-size: 20px;
  }

  .login-subtitle {
    font-size: 12px;
  }

  .system-info {
    position: relative;
    bottom: auto;
    margin-top: 20px;
  }
}

// 暗黑模式适配
[data-theme='dark'] {
  .login-card {
    background: rgba(26, 26, 26, 0.95);
  }

  .login-title {
    color: var(--text-primary);
  }

  .login-footer {
    color: var(--text-secondary);
  }
}
</style>
