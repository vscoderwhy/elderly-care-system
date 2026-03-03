<template>
  <div class="not-found">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <div class="error-title">页面未找到</div>
      <div class="error-desc">
        抱歉，您访问的页面不存在或已被移除
      </div>
      <div class="error-actions">
        <el-button type="primary" @click="goBack">
          <el-icon><Back /></el-icon>
          返回上一页
        </el-button>
        <el-button @click="goHome">
          <el-icon><HomeFilled /></el-icon>
          回到首页
        </el-button>
      </div>
    </div>

    <!-- 装饰插图 -->
    <div class="not-found-illustration">
      <svg viewBox="0 0 400 400" xmlns="http://www.w3.org/2000/svg">
        <circle cx="200" cy="200" r="150" fill="var(--bg-tertiary)" opacity="0.5"/>
        <circle cx="200" cy="200" r="100" fill="var(--bg-secondary)" opacity="0.5"/>
        <text x="200" y="220" text-anchor="middle" font-size="80" fill="var(--text-tertiary)" font-weight="bold">404</text>
        <path d="M150 280 Q200 320 250 280" stroke="var(--border-color)" stroke-width="3" fill="none"/>
        <circle cx="170" cy="260" r="8" fill="var(--text-tertiary)"/>
        <circle cx="230" cy="260" r="8" fill="var(--text-tertiary)"/>
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Back, HomeFilled } from '@element-plus/icons-vue'

const router = useRouter()

const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/dashboard')
}
</script>

<style scoped lang="scss">
.not-found {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-secondary);
}

.not-found-content {
  text-align: center;
  margin-bottom: 40px;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  line-height: 1;
  background: linear-gradient(135deg, var(--primary-color) 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 16px;
}

.error-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.error-desc {
  font-size: 16px;
  color: var(--text-secondary);
  margin-bottom: 32px;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.not-found-illustration {
  width: 400px;
  max-width: 100%;

  svg {
    width: 100%;
    height: auto;
  }
}

@media (max-width: 768px) {
  .error-code {
    font-size: 80px;
  }

  .error-title {
    font-size: 24px;
  }

  .error-desc {
    font-size: 14px;
  }

  .error-actions {
    flex-direction: column;
  }

  .not-found-illustration {
    width: 280px;
  }
}
</style>
