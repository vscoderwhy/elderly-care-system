import { ref, computed } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import type { Router } from 'vue-router'

export interface ErrorInfo {
  code?: string | number
  message: string
  type?: 'network' | 'business' | 'system' | 'validation'
  stack?: string
  details?: any
}

const globalError = ref<ErrorInfo | null>(null)
const errorCount = ref(0)

export const useErrorHandler = (router?: Router) => {
  /**
   * 处理API错误
   */
  const handleApiError = (error: any, defaultMessage: string = '操作失败') => {
    console.error('[API Error]', error)

    let errorMessage = defaultMessage
    let errorCode: string | number | undefined

    // Axios 错误
    if (error?.response) {
      const { status, data } = error.response
      errorCode = status
      errorMessage = data?.message || data?.error || `请求失败 (${status})`
    }
    // 网络错误
    else if (error?.request) {
      errorCode = 'NETWORK_ERROR'
      errorMessage = '网络连接失败，请检查网络设置'
    }
    // 其他错误
    else if (error?.message) {
      errorMessage = error.message
    }

    const errorInfo: ErrorInfo = {
      code: errorCode,
      message: errorMessage,
      type: 'network',
      details: error
    }

    globalError.value = errorInfo
    errorCount.value++

    // 显示错误消息
    ElMessage.error(errorMessage)

    return errorInfo
  }

  /**
   * 处理业务错误
   */
  const handleBusinessError = (message: string, details?: any) => {
    const errorInfo: ErrorInfo = {
      code: 'BUSINESS_ERROR',
      message,
      type: 'business',
      details
    }

    globalError.value = errorInfo
    errorCount.value++

    ElMessage.warning(message)

    return errorInfo
  }

  /**
   * 处理验证错误
   */
  const handleValidationError = (message: string, errors?: any) => {
    const errorInfo: ErrorInfo = {
      code: 'VALIDATION_ERROR',
      message,
      type: 'validation',
      details: errors
    }

    globalError.value = errorInfo

    ElMessage.warning(message)

    return errorInfo
  }

  /**
   * 处理系统错误
   */
  const handleSystemError = (error: Error, context?: string) => {
    console.error('[System Error]', error, context)

    const errorInfo: ErrorInfo = {
      code: 'SYSTEM_ERROR',
      message: error.message || '系统错误',
      type: 'system',
      stack: error.stack,
      details: context
    }

    globalError.value = errorInfo
    errorCount.value++

    // 系统错误显示通知
    ElNotification.error({
      title: '系统错误',
      message: error.message || '发生了系统错误，请稍后重试',
      duration: 5000
    })

    // 开发环境输出详细堆栈
    if (import.meta.env.DEV && error.stack) {
      console.error('Error Stack:', error.stack)
    }

    return errorInfo
  }

  /**
   * 清除错误
   */
  const clearError = () => {
    globalError.value = null
  }

  /**
   * 显示成功消息
   */
  const showSuccess = (message: string) => {
    ElMessage.success(message)
    clearError()
  }

  /**
   * 异步操作包装器
   */
  const withErrorHandling = async <T>(
    asyncFn: () => Promise<T>,
    options?: {
      errorMessage?: string
      successMessage?: string
      showSuccess?: boolean
    }
  ): Promise<T | null> => {
    try {
      const result = await asyncFn()

      if (options?.showSuccess && options?.successMessage) {
        showSuccess(options.successMessage)
      }

      return result
    } catch (error) {
      handleApiError(error, options?.errorMessage)
      return null
    }
  }

  /**
   * 页面级错误处理
   */
  const withPageErrorHandling = async <T>(
    asyncFn: () => Promise<T>,
    options?: {
      errorMessage?: string
      redirectOnError?: string
    }
  ): Promise<T | null> => {
    const result = await withErrorHandling(asyncFn, options)

    if (result === null && options?.redirectOnError && router) {
      router.push(options.redirectOnError)
    }

    return result
  }

  return {
    globalError: computed(() => globalError.value),
    errorCount: computed(() => errorCount.value),
    hasError: computed(() => globalError.value !== null),
    handleApiError,
    handleBusinessError,
    handleValidationError,
    handleSystemError,
    clearError,
    showSuccess,
    withErrorHandling,
    withPageErrorHandling
  }
}

/**
 * 全局错误处理器（在main.ts中注册）
 */
export const setupGlobalErrorHandler = (app: any, router: Router) => {
  // 处理未捕获的Promise错误
  window.addEventListener('unhandledrejection', (event) => {
    console.error('[Unhandled Rejection]', event.reason)
    event.preventDefault()

    ElNotification.error({
      title: '未处理的错误',
      message: event.reason?.message || '发生了未知错误',
      duration: 5000
    })
  })

  // 处理全局错误
  window.addEventListener('error', (event) => {
    console.error('[Global Error]', event.error)

    if (import.meta.env.PROD) {
      // 生产环境上报错误
      // TODO: 接入错误追踪服务（如Sentry）
    }
  })

  // Vue错误处理
  app.config.errorHandler = (err, instance, info) => {
    console.error('[Vue Error]', err, info)

    ElNotification.error({
      title: '组件错误',
      message: err.message || '组件渲染出错',
      duration: 5000
    })

    if (import.meta.env.DEV) {
      console.error('Error Info:', info)
    }
  }

  // 路由错误处理
  router.onError((error) => {
    console.error('[Router Error]', error)
    ElMessage.error('页面加载失败')
  })
}
