<template>
  <div class="error-boundary">
    <!-- 正常内容 -->
    <div v-if="!hasError" class="error-boundary-content">
      <slot />
    </div>

    <!-- 错误展示 -->
    <div v-else class="error-fallback">
      <div class="error-icon">
        <el-icon><WarningFilled /></el-icon>
      </div>
      <h2 class="error-title">Oops, 出错了</h2>
      <p class="error-message">{{ error?.message || '页面加载失败' }}</p>
      <p v-if="error?.type === 'network'" class="error-hint">
        请检查网络连接后重试
      </p>

      <div class="error-actions">
        <el-button type="primary" @click="handleRetry">
          <el-icon><Refresh /></el-icon>
          重试
        </el-button>
        <el-button @click="handleGoHome">
          <el-icon><HomeFilled /></el-icon>
          返回首页
        </el-button>
      </div>

      <!-- 开发环境显示错误详情 -->
      <div v-if="isDev && error?.stack" class="error-details">
        <el-collapse>
          <el-collapse-item title="错误详情（仅开发环境）" name="stack">
            <pre class="error-stack">{{ error.stack }}</pre>
          </el-collapse-item>
        </el-collapse>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onErrorCaptured, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { WarningFilled, Refresh, HomeFilled } from '@element-plus/icons-vue'

const props = defineProps<{
  // 是否捕获所有错误（包括子组件）
  captureAll?: boolean
  // 错误回调
  onError?: (error: Error) => void
}>()

const emit = defineEmits<{
  error: [error: Error]
  reset: []
}>()

const router = useRouter()
const hasError = ref(false)
const error = ref<Error & { type?: string; code?: string | number } | null>(null)

const isDev = computed(() => import.meta.env.DEV)

// 捕获子组件错误
onErrorCaptured((err, instance, info) => {
  console.error('[ErrorBoundary] Caught error:', err, info)

  hasError.value = true
  error.value = err

  // 调用错误回调
  props.onError?.(err)

  // 发送错误事件
  emit('error', err)

  // 阻止错误继续传播
  if (props.captureAll) {
    return false
  }
})

// 处理自身错误
onMounted(() => {
  window.addEventListener('error', (event) => {
    if (event.error) {
      hasError.value = true
      error.value = event.error
    }
  })
})

// 重试
const handleRetry = () => {
  hasError.value = false
  error.value = null
  emit('reset')
}

// 返回首页
const handleGoHome = () => {
  router.push('/')
}

// 暴露方法供父组件调用
defineExpose({
  hasError,
  error,
  reset: () => {
    hasError.value = false
    error.value = null
  }
})
</script>

<style scoped lang="scss">
.error-boundary {
  width: 100%;
  height: 100%;
}

.error-boundary-content {
  width: 100%;
  height: 100%;
}

.error-fallback {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 40px;
  text-align: center;

  .error-icon {
    font-size: 64px;
    color: var(--warning-color);
    margin-bottom: 20px;
  }

  .error-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 12px 0;
  }

  .error-message {
    font-size: 16px;
    color: var(--text-secondary);
    margin: 0 0 8px 0;
  }

  .error-hint {
    font-size: 14px;
    color: var(--text-tertiary);
    margin: 0 0 24px 0;
  }

  .error-actions {
    display: flex;
    gap: 12px;
  }

  .error-details {
    margin-top: 32px;
    width: 100%;
    max-width: 600px;
    text-align: left;

    .error-stack {
      background: var(--bg-tertiary);
      padding: 12px;
      border-radius: 4px;
      font-size: 12px;
      color: var(--danger-color);
      overflow-x: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
  }
}

@media (max-width: 768px) {
  .error-fallback {
    padding: 20px;

    .error-icon {
      font-size: 48px;
    }

    .error-title {
      font-size: 20px;
    }

    .error-message {
      font-size: 14px;
    }

    .error-actions {
      flex-direction: column;

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
