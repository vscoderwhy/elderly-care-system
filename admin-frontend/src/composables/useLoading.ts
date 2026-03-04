import { ref, computed } from 'vue'

export interface LoadingState {
  loading: boolean
  text?: string
  background?: string
}

const globalLoading = ref<LoadingState>({
  loading: false
})

const loadingStack = ref<Set<string>>(new Set())

/**
 * 全局Loading管理
 */
export const useLoading = () => {
  /**
   * 显示loading
   */
  const show = (text?: string, background?: string) => {
    globalLoading.value = {
      loading: true,
      text,
      background
    }
  }

  /**
   * 隐藏loading
   */
  const hide = () => {
    globalLoading.value = {
      loading: false
    }
  }

  /**
   * 带标识的loading（用于多个并发请求）
   */
  const showWithKey = (key: string, text?: string) => {
    loadingStack.value.add(key)
    updateGlobalLoading()
  }

  const hideWithKey = (key: string) => {
    loadingStack.value.delete(key)
    updateGlobalLoading()
  }

  const updateGlobalLoading = () => {
    globalLoading.value = {
      loading: loadingStack.value.size > 0
    }
  }

  /**
   * 异步操作包装器
   */
  const withLoading = async <T>(
    asyncFn: () => Promise<T>,
    loadingText?: string
  ): Promise<T> => {
    show(loadingText)
    try {
      return await asyncFn()
    } finally {
      hide()
    }
  }

  /**
   * 页面loading包装器
   */
  const withPageLoading = async <T>(
    asyncFn: () => Promise<T>,
    loadingText = '加载中...'
  ): Promise<T> => {
    return withLoading(asyncFn, loadingText)
  }

  return {
    globalLoading: computed(() => globalLoading.value),
    isLoading: computed(() => globalLoading.value.loading),
    loadingText: computed(() => globalLoading.value.text),
    show,
    hide,
    showWithKey,
    hideWithKey,
    withLoading,
    withPageLoading
  }
}

/**
 * 创建实例级别的loading（用于组件内）
 */
export const useInstanceLoading = () => {
  const loading = ref(false)
  const text = ref('')

  return {
    loading,
    text,
    show: (msg?: string) => {
      loading.value = true
      if (msg) text.value = msg
    },
    hide: () => {
      loading.value = false
      text.value = ''
    },
    withLoading: async <T>(asyncFn: () => Promise<T>, loadingMsg?: string): Promise<T> => {
      loading.value = true
      if (loadingMsg) text.value = loadingMsg
      try {
        return await asyncFn()
      } finally {
        loading.value = false
        text.value = ''
      }
    }
  }
}
