// ========== 日期时间相关 ==========

/**
 * 格式化日期
 * @param date 日期对象、时间戳或日期字符串
 * @param format 格式化模板，默认 'YYYY-MM-DD HH:mm:ss'
 */
export const formatDate = (
  date: Date | string | number,
  format: string = 'YYYY-MM-DD HH:mm:ss'
): string => {
  const d = new Date(date)
  if (isNaN(d.getTime())) return ''

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 相对时间格式化（多久之前）
 */
export const formatRelativeTime = (date: Date | string | number): string => {
  const d = new Date(date)
  const now = new Date()
  const diff = now.getTime() - d.getTime()

  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  const week = 7 * day
  const month = 30 * day
  const year = 365 * day

  if (diff < minute) {
    return '刚刚'
  } else if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`
  } else if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`
  } else if (diff < week) {
    return `${Math.floor(diff / day)}天前`
  } else if (diff < month) {
    return `${Math.floor(diff / week)}周前`
  } else if (diff < year) {
    return `${Math.floor(diff / month)}个月前`
  } else {
    return `${Math.floor(diff / year)}年前`
  }
}

/**
 * 获取今天开始/结束时间戳
 */
export const getTodayRange = (): [number, number] => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const end = start + 24 * 60 * 60 * 1000 - 1
  return [start, end]
}

/**
 * 获取本周开始/结束时间戳
 */
export const getWeekRange = (): [number, number] => {
  const now = new Date()
  const day = now.getDay() || 7
  const start = new Date(now.getFullYear(), now.getMonth(), now.getDate() - day + 1).getTime()
  const end = start + 7 * 24 * 60 * 60 * 1000 - 1
  return [start, end]
}

/**
 * 获取本月开始/结束时间戳
 */
export const getMonthRange = (): [number, number] => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1).getTime()
  const end = new Date(now.getFullYear(), now.getMonth() + 1, 1).getTime() - 1
  return [start, end]
}

// ========== 数字格式化 ==========

/**
 * 格式化金额（千分位）
 */
export const formatMoney = (amount: number | string, decimals = 2): string => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  if (isNaN(num)) return '0.00'

  return num.toFixed(decimals).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/**
 * 格式化文件大小
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

// ========== 字符串处理 ==========

/**
 * 手机号脱敏
 */
export const maskPhone = (phone: string): string => {
  if (!phone || phone.length !== 11) return phone
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/**
 * 身份证号脱敏
 */
export const maskIdCard = (idCard: string): string => {
  if (!idCard || idCard.length < 18) return idCard
  return idCard.replace(/(\d{6})\d{8}(\d{4})/, '$1********$2')
}

/**
 * 姓名脱敏
 */
export const maskName = (name: string): string => {
  if (!name) return name
  if (name.length === 2) {
    return name[0] + '*'
  } else if (name.length > 2) {
    return name[0] + '*'.repeat(name.length - 2) + name[name.length - 1]
  }
  return name
}

/**
 * 截断文本
 */
export const truncateText = (text: string, maxLength: number, suffix = '...'): string => {
  if (!text || text.length <= maxLength) return text
  return text.substring(0, maxLength) + suffix
}

// ========== 数据验证 ==========

/**
 * 验证手机号
 */
export const isValidPhone = (phone: string): boolean => {
  return /^1[3-9]\d{9}$/.test(phone)
}

/**
 * 验证身份证号
 */
export const isValidIdCard = (idCard: string): boolean => {
  return /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/.test(idCard)
}

/**
 * 验证邮箱
 */
export const isValidEmail = (email: string): boolean => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

// ========== 存储相关 ==========

/**
 * 设置存储（带过期时间）
 */
export const setStorageWithExpire = (key: string, value: any, expireSeconds: number) => {
  const data = {
    value,
    expire: Date.now() + expireSeconds * 1000
  }
  uni.setStorageSync(key, JSON.stringify(data))
}

/**
 * 获取存储（检查过期）
 */
export const getStorageWithExpire = (key: string) => {
  const dataStr = uni.getStorageSync(key)
  if (!dataStr) return null

  try {
    const data = JSON.parse(dataStr)
    if (data.expire && data.expire < Date.now()) {
      uni.removeStorageSync(key)
      return null
    }
    return data.value
  } catch {
    return dataStr
  }
}

// ========== 图片处理 ==========

/**
 * 压缩图片
 */
export const compressImage = (src: string, quality = 80): Promise<string> => {
  return new Promise((resolve, reject) => {
    uni.compressImage({
      src,
      quality,
      success: (res) => resolve(res.tempFilePath),
      fail: reject
    })
  })
}

/**
 * 选择并上传图片
 */
export const chooseAndUploadImage = (count = 1, sizeType = ['compressed'], sourceType = ['album', 'camera']) => {
  return new Promise((resolve, reject) => {
    uni.chooseImage({
      count,
      sizeType,
      sourceType,
      success: async (res) => {
        try {
          const tempFilePaths = res.tempFilePaths
          // 这里可以添加上传逻辑
          resolve(tempFilePaths)
        } catch (error) {
          reject(error)
        }
      },
      fail: reject
    })
  })
}

// ========== 位置相关 ==========

/**
 * 获取当前位置
 */
export const getCurrentLocation = (): Promise<{
  latitude: number
  longitude: number
}> => {
  return new Promise((resolve, reject) => {
    uni.getLocation({
      type: 'gcj02',
      success: (res) => {
        resolve({
          latitude: res.latitude,
          longitude: res.longitude
        })
      },
      fail: reject
    })
  })
}

/**
 * 计算两点距离（单位：米）
 */
export const getDistance = (
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number => {
  const R = 6371e3 // 地球半径（米）
  const φ1 = (lat1 * Math.PI) / 180
  const φ2 = (lat2 * Math.PI) / 180
  const Δφ = ((lat2 - lat1) * Math.PI) / 180
  const Δλ = ((lon2 - lon1) * Math.PI) / 180

  const a =
    Math.sin(Δφ / 2) * Math.sin(Δφ / 2) +
    Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))

  return R * c
}

// ========== 防抖节流 ==========

/**
 * 防抖
 */
export const debounce = <T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timer: number | null = null

  return function (this: any, ...args: Parameters<T>) {
    if (timer) clearTimeout(timer)
    timer = window.setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

/**
 * 节流
 */
export const throttle = <T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let lastTime = 0

  return function (this: any, ...args: Parameters<T>) {
    const now = Date.now()
    if (now - lastTime >= delay) {
      lastTime = now
      fn.apply(this, args)
    }
  }
}

// ========== 页面导航 ==========

/**
 * 延迟跳转（防止快速点击）
 */
export const navigateToDebounce = debounce((url: string) => {
  uni.navigateTo({ url })
}, 300)

/**
 * 返回上一页或首页
 */
export const navigateBackOrHome = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.reLaunch({
      url: '/pages/index/index'
    })
  }
}
