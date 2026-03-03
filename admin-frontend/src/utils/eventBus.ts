/**
 * 事件总线 - 用于跨组件通信
 */
type EventBusCallback = (...args: any[]) => void

class EventBus {
  private events: Record<string, EventBusCallback[]> = {}

  /**
   * 订阅事件
   */
  on(event: string, callback: EventBusCallback) {
    if (!this.events[event]) {
      this.events[event] = []
    }
    this.events[event].push(callback)
  }

  /**
   * 取消订阅
   */
  off(event: string, callback: EventBusCallback) {
    if (!this.events[event]) return

    const index = this.events[event].indexOf(callback)
    if (index > -1) {
      this.events[event].splice(index, 1)
    }
  }

  /**
   * 触发事件
   */
  emit(event: string, ...args: any[]) {
    if (!this.events[event]) return

    this.events[event].forEach(callback => {
      try {
        callback(...args)
      } catch (error) {
        console.error(`[EventBus] Error in event handler for "${event}":`, error)
      }
    })
  }

  /**
   * 订阅一次性事件
   */
  once(event: string, callback: EventBusCallback) {
    const onceCallback = (...args: any[]) => {
      callback(...args)
      this.off(event, onceCallback)
    }
    this.on(event, onceCallback)
  }

  /**
   * 清除所有事件
   */
  clear() {
    this.events = {}
  }

  /**
   * 清除指定事件的所有监听器
   */
  clearEvent(event: string) {
    delete this.events[event]
  }
}

// 导出单例
export const eventBus = new EventBus()

// 事件名称常量
export const Events = {
  // 主题相关
  THEME_CHANGE: 'theme:change',
  THEME_DARK: 'theme:dark',
  THEME_LIGHT: 'theme:light',

  // 用户相关
  USER_LOGIN: 'user:login',
  USER_LOGOUT: 'user:logout',
  USER_UPDATE: 'user:update',

  // 数据相关
  DATA_REFRESH: 'data:refresh',
  DATA_UPDATE: 'data:update',
  DATA_DELETE: 'data:delete',

  // 护理相关
  CARE_TASK_CREATE: 'care:task:create',
  CARE_TASK_UPDATE: 'care:task:update',
  CARE_TASK_COMPLETE: 'care:task:complete',

  // 健康相关
  HEALTH_ALERT: 'health:alert',
  HEALTH_DATA_UPDATE: 'health:data:update',

  // 财务相关
  BILL_CREATE: 'bill:create',
  BILL_PAID: 'bill:paid',
  BILL_OVERDUE: 'bill:overdue',

  // 系统相关
  NOTIFICATION: 'system:notification',
  ERROR: 'system:error',
  WARNING: 'system:warning'
}
