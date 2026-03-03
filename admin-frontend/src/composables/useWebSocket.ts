import { ref, onUnmounted } from 'vue'

export interface WebSocketOptions {
  onMessage?: (data: any) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (error: Event) => void
  reconnect?: boolean
  reconnectInterval?: number
  maxReconnectAttempts?: number
}

export const useWebSocket = (url: string, options: WebSocketOptions = {}) => {
  const {
    onMessage,
    onOpen,
    onClose,
    onError,
    reconnect = true,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5
  } = options

  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const reconnectAttempts = ref(0)
  const data = ref<any>(null)

  const connect = () => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      return
    }

    try {
      ws.value = new WebSocket(url)

      ws.value.onopen = () => {
        connected.value = true
        reconnectAttempts.value = 0
        onOpen?.()
      }

      ws.value.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          data.value = message
          onMessage?.(message)
        } catch (error) {
          // 非JSON消息，直接传递文本
          onMessage?.(event.data)
        }
      }

      ws.value.onclose = () => {
        connected.value = false
        onClose?.()

        // 自动重连
        if (reconnect && reconnectAttempts.value < maxReconnectAttempts) {
          reconnectAttempts.value++
          setTimeout(() => {
            connect()
          }, reconnectInterval)
        }
      }

      ws.value.onerror = (error) => {
        onError?.(error)
      }
    } catch (error) {
      console.error('[WebSocket] Connection error:', error)
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
      connected.value = false
    }
  }

  const send = (message: any) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(typeof message === 'string' ? message : JSON.stringify(message))
    } else {
      console.warn('[WebSocket] Cannot send message: connection not open')
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    data,
    connect,
    disconnect,
    send
  }
}

/**
 * 轮询Hook（当WebSocket不可用时使用）
 */
export const usePolling = (callback: () => Promise<void>, interval: number = 5000) => {
  let timer: number | null = null
  const isRunning = ref(false)

  const start = () => {
    if (isRunning.value) return

    isRunning.value = true

    // 立即执行一次
    callback()

    // 设置定时器
    timer = window.setInterval(() => {
      callback()
    }, interval)
  }

  const stop = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    isRunning.value = false
  }

  // 组件卸载时自动停止
  onUnmounted(() => {
    stop()
  })

  return {
    isRunning,
    start,
    stop
  }
}
