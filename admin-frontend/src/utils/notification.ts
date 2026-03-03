import { ElMessage, ElNotification, ElMessageBox } from 'element-plus'
import type { MessageParams, NotificationParams, MessageBoxData } from 'element-plus'

/**
 * 通知服务 - 统一的消息提示和弹窗管理
 */
export const notify = {
  /**
   * 成功消息
   */
  success(message: string, options?: Partial<MessageParams>) {
    return ElMessage.success({
      message,
      duration: 3000,
      showClose: true,
      ...options
    })
  },

  /**
   * 警告消息
   */
  warning(message: string, options?: Partial<MessageParams>) {
    return ElMessage.warning({
      message,
      duration: 3000,
      showClose: true,
      ...options
    })
  },

  /**
   * 信息消息
   */
  info(message: string, options?: Partial<MessageParams>) {
    return ElMessage.info({
      message,
      duration: 3000,
      showClose: true,
      ...options
    })
  },

  /**
   * 错误消息
   */
  error(message: string, options?: Partial<MessageParams>) {
    return ElMessage.error({
      message,
      duration: 5000,
      showClose: true,
      ...options
    })
  },

  /**
   * 加载消息
   */
  loading(message: string = '加载中...', options?: Partial<MessageParams>) {
    return ElMessage({
      message,
      duration: 0,
      showClose: false,
      ...options
    })
  }
}

/**
 * 通知弹窗
 */
export const notification = {
  /**
   * 成功通知
   */
  success(title: string, message?: string, options?: Partial<NotificationParams>) {
    return ElNotification.success({
      title,
      message,
      duration: 3000,
      ...options
    })
  },

  /**
   * 警告通知
   */
  warning(title: string, message?: string, options?: Partial<NotificationParams>) {
    return ElNotification.warning({
      title,
      message,
      duration: 3000,
      ...options
    })
  },

  /**
   * 信息通知
   */
  info(title: string, message?: string, options?: Partial<NotificationParams>) {
    return ElNotification.info({
      title,
      message,
      duration: 3000,
      ...options
    })
  },

  /**
   * 错误通知
   */
  error(title: string, message?: string, options?: Partial<NotificationParams>) {
    return ElNotification.error({
      title,
      message,
      duration: 5000,
      ...options
    })
  }
}

/**
 * 确认对话框
 */
export const confirm = {
  /**
   * 简单确认
   */
  async confirm(
    message: string,
    title: string = '提示',
    options?: any
  ): Promise<boolean> {
    try {
      await ElMessageBox.confirm(message, title, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        ...options
      })
      return true
    } catch {
      return false
    }
  },

  /**
   * 危险操作确认
   */
  async danger(
    message: string,
    title: string = '危险操作',
    options?: any
  ): Promise<boolean> {
    try {
      await ElMessageBox.confirm(message, title, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error',
        confirmButtonClass: 'el-button--danger',
        ...options
      })
      return true
    } catch {
      return false
    }
  },

  /**
   * 删除确认
   */
  async delete(
    itemName: string = '此项目',
    options?: any
  ): Promise<boolean> {
    return this.danger(`确定要删除${itemName}吗？删除后无法恢复。`, '删除确认', options)
  },

  /**
   * 输入对话框
   */
  async prompt(
    message: string,
    title: string = '输入',
    defaultValue: string = '',
    options?: any
  ): Promise<string | null> {
    try {
      const { value } = await ElMessageBox.prompt(message, title, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputValue: defaultValue,
        ...options
      })
      return value as string
    } catch {
      return null
    }
  },

  /**
   * 警告对话框
   */
  async alert(
    message: string,
    title: string = '提示',
    type: 'success' | 'warning' | 'info' | 'error' = 'info',
    options?: any
  ): Promise<void> {
    await ElMessageBox.alert(message, title, {
      type,
      confirmButtonText: '确定',
      ...options
    })
  }
}

/**
 * Toast快捷方式
 */
export const toast = notify

/**
 * Alert快捷方式
 */
export const alert = confirm.alert

/**
 * Prompt快捷方式
 */
export const prompt = confirm.prompt
