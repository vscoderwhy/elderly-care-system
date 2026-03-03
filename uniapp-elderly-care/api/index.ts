import { get, post, put, del } from './request'

// ========== 用户相关 ==========
// 微信登录
export const wxLogin = (code: string) => {
  return post('/auth/wx/login', { code })
}

// 手机号登录
export const phoneLogin = (phone: string, code: string) => {
  return post('/auth/phone/login', { phone, code })
}

// 发送验证码
export const sendSmsCode = (phone: string) => {
  return post('/auth/sms/send', { phone })
}

// 退出登录
export const logout = () => {
  return post('/auth/logout')
}

// 获取用户信息
export const getUserInfo = () => {
  return get('/user/info')
}

// 更新用户信息
export const updateUserInfo = (data: any) => {
  return put('/user/info', data)
}

// ========== 老人相关 ==========
// 获取老人列表（家属端）
export const getElderlyList = () => {
  return get('/elderly/my')
}

// 获取老人详情
export const getElderlyDetail = (id: string | number) => {
  return get(`/elderly/${id}`)
}

// 获取老人基本信息
export const getElderlyBaseInfo = (id: string | number) => {
  return get(`/elderly/${id}/base`)
}

// 获取老人照片列表
export const getElderlyPhotos = (id: string | number, page = 1, pageSize = 20) => {
  return get(`/elderly/${id}/photos`, { page, pageSize })
}

// ========== 护理记录相关 ==========
// 获取护理记录列表
export const getCareRecords = (params?: {
  elderlyId?: string | number
  page?: number
  pageSize?: number
}) => {
  return get('/care/records', params)
}

// 获取护理记录详情
export const getCareRecordDetail = (id: string | number) => {
  return get(`/care/records/${id}`)
}

// 提交服务评价
export const submitEvaluation = (recordId: string | number, data: {
  rating: number
  comment: string
  tags?: string[]
}) => {
  return post(`/care/records/${recordId}/evaluation`, data)
}

// ========== 健康数据相关 ==========
// 获取健康数据列表
export const getHealthData = (params?: {
  elderlyId?: string | number
  type?: string
  startDate?: string
  endDate?: string
  page?: number
  pageSize?: number
}) => {
  return get('/health/data', params)
}

// 获取健康数据详情
export const getHealthDetail = (id: string | number) => {
  return get(`/health/data/${id}`)
}

// 获取健康趋势图数据
export const getHealthTrend = (params: {
  elderlyId: string | number
  type: 'bloodPressure' | 'bloodSugar' | 'temperature' | 'weight'
  days?: number
}) => {
  return get('/health/trend', params)
}

// 获取最新健康数据
export const getLatestHealth = (elderlyId: string | number) => {
  return get(`/health/latest/${elderlyId}`)
}

// ========== 费用账单相关 ==========
// 获取账单列表
export const getBillList = (params?: {
  status?: string
  year?: number
  month?: number
  page?: number
  pageSize?: number
}) => {
  return get('/bills', params)
}

// 获取账单详情
export const getBillDetail = (id: string | number) => {
  return get(`/bills/${id}`)
}

// 创建支付订单
export const createPayment = (billId: string | number, method: 'wechat' | 'alipay') => {
  return post('/payment/create', { billId, method })
}

// 查询支付状态
export const getPaymentStatus = (orderId: string) => {
  return get(`/payment/status/${orderId}`)
}

// 下载发票
export const downloadInvoice = (billId: string | number) => {
  return get(`/bills/${billId}/invoice`)
}

// ========== 探视预约相关 ==========
// 获取预约列表
export const getVisitAppointments = (params?: {
  status?: string
  page?: number
  pageSize?: number
}) => {
  return get('/visits/appointments', params)
}

// 创建预约
export const createAppointment = (data: {
  elderlyId: string | number
  visitDate: string
  visitTime: string
  visitors: number
  contactName: string
  contactPhone: string
  remark?: string
}) => {
  return post('/visits/appointments', data)
}

// 取消预约
export const cancelAppointment = (id: string | number) => {
  return post(`/visits/appointments/${id}/cancel`)
}

// 获取预约详情
export const getAppointmentDetail = (id: string | number) => {
  return get(`/visits/appointments/${id}`)
}

// 获取可预约时段
export const getAvailableSlots = (date: string) => {
  return get('/visits/slots', { date })
}

// ========== 消息通知相关 ==========
// 获取通知列表
export const getNotifications = (params?: {
  type?: string
  page?: number
  pageSize?: number
}) => {
  return get('/notifications', params)
}

// 标记已读
export const markAsRead = (id: string | number) => {
  return post(`/notifications/${id}/read`)
}

// 全部标记已读
export const markAllAsRead = () => {
  return post('/notifications/read-all')
}

// 获取未读数量
export const getUnreadCount = () => {
  return get('/notifications/unread-count')
}

// ========== 护工端相关 ==========
// 获取今日任务列表
export const getTodayTasks = () => {
  return get('/nurse/tasks/today')
}

// 获取任务详情
export const getTaskDetail = (id: string | number) => {
  return get(`/nurse/tasks/${id}`)
}

// 更新任务状态
export const updateTaskStatus = (id: string | number, status: string, data?: any) => {
  return post(`/nurse/tasks/${id}/status`, { status, ...data })
}

// 快速记录护理
export const quickRecord = (data: {
  elderlyId: string | number
  careType: string
  description: string
  images?: string[]
  location?: string
}) => {
  return post('/nurse/quick-record', data)
}

// 上传图片
export const uploadImage = (filePath: string) => {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: import.meta.env.VITE_API_BASE_URL + '/upload/image',
      filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${uni.getStorageSync('token')}`
      },
      success: (res) => {
        const data = JSON.parse(res.data)
        if (data.code === 0) {
          resolve(data.data.url)
        } else {
          reject(data)
        }
      },
      fail: reject
    })
  })
}

// 打卡签到
export const clockIn = (type: 'in' | 'out', location?: { latitude: number; longitude: number }) => {
  return post('/nurse/clock', { type, location })
}

// 获取考勤记录
export const getAttendanceRecords = (params?: {
  year?: number
  month?: number
}) => {
  return get('/nurse/attendance', params)
}

// 获取排班信息
export const getSchedule = (params?: {
  startDate?: string
  endDate?: string
}) => {
  return get('/nurse/schedule', params)
}

// ========== 统计相关 ==========
// 获取首页统计数据
export const getHomeStats = () => {
  return get('/stats/home')
}

// 获取护理统计
export const getCareStats = (params?: {
  elderlyId?: string | number
  startDate?: string
  endDate?: string
}) => {
  return get('/stats/care', params)
}

// 获取健康统计
export const getHealthStats = (elderlyId: string | number) => {
  return get(`/stats/health/${elderlyId}`)
}
