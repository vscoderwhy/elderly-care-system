import { ref, onUnmounted } from 'vue'

export interface WebSocketMessage {
  type: string
  data: any
  timestamp: number
}

export interface WebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (error: Event) => void
  reconnect?: boolean
  reconnectInterval?: number
}

export const useWebSocket = (url: string, options: WebSocketOptions = {}) => {
  const {
    onMessage,
    onOpen,
    onClose,
    onError,
    reconnect = true,
    reconnectInterval = 3000
  } = options

  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const reconnectTimer = ref<number>()

  const connect = () => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      return
    }

    try {
      ws.value = new WebSocket(url)

      ws.value.onopen = () => {
        connected.value = true
        console.log('[WebSocket] Connected')
        onOpen?.()
        // 清除重连定时器
        if (reconnectTimer.value) {
          clearTimeout(reconnectTimer.value)
          reconnectTimer.value = undefined
        }
      }

      ws.value.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          onMessage?.(message)
        } catch (error) {
          console.error('[WebSocket] Failed to parse message:', error)
        }
      }

      ws.value.onclose = () => {
        connected.value = false
        console.log('[WebSocket] Closed')
        onClose?.()

        // 自动重连
        if (reconnect && !reconnectTimer.value) {
          reconnectTimer.value = window.setTimeout(() => {
            console.log('[WebSocket] Reconnecting...')
            connect()
          }, reconnectInterval)
        }
      }

      ws.value.onerror = (error) => {
        console.error('[WebSocket] Error:', error)
        onError?.(error)
      }
    } catch (error) {
      console.error('[WebSocket] Failed to connect:', error)
    }
  }

  const disconnect = () => {
    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value)
      reconnectTimer.value = undefined
    }
    ws.value?.close()
    ws.value = null
    connected.value = false
  }

  const send = (data: any) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(data))
    } else {
      console.warn('[WebSocket] Cannot send message: not connected')
    }
  }

  // 组件卸载时断开连接
  onUnmounted(() => {
    disconnect()
  })

  return {
    ws,
    connected,
    connect,
    disconnect,
    send
  }
}

// 轮询降级方案
export const usePolling = (callback: () => Promise<void>, interval: number = 5000) => {
  const timer = ref<number>()
  const isActive = ref(false)

  const start = () => {
    if (isActive.value) return
    isActive.value = true

    const poll = async () => {
      try {
        await callback()
      } catch (error) {
        console.error('[Polling] Error:', error)
      }
    }

    // 立即执行一次
    poll()
    // 定时执行
    timer.value = window.setInterval(poll, interval)
  }

  const stop = () => {
    if (timer.value) {
      clearInterval(timer.value)
      timer.value = undefined
    }
    isActive.value = false
  }

  // 组件卸载时停止轮询
  onUnmounted(() => {
    stop()
  })

  return {
    isActive,
    start,
    stop
  }
}
