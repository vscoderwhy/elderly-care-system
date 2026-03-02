import axios from 'axios'
import type { AxiosInstance } from 'axios'
import { showToast, closeToast } from 'vant'
import type { ToastInstance } from 'vant'
import { mockApi } from './mock'

// 开发模式使用 Mock 数据
const USE_MOCK = false

const baseURL = import.meta.env.MODE === 'production'
  ? '/api'
  : 'http://1.12.223.138:8080/api'

const instance: AxiosInstance = axios.create({
  baseURL,
  timeout: 10000
})

let loadingToast: ToastInstance | null = null

instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    loadingToast = showToast({
      message: '加载中...',
      forbidClick: true,
      duration: 0
    })
    return config
  }
)

instance.interceptors.response.use(
  (response) => {
    if (loadingToast) {
      closeToast()
    }
    const { code, message, data } = response.data
    if (code === 0) {
      return data
    } else {
      showToast({ message: message || '请求失败' })
      return Promise.reject(new Error(message))
    }
  },
  (error) => {
    if (loadingToast) {
      closeToast()
    }
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    showToast({ message: error.message || '网络错误' })
    return Promise.reject(error)
  }
)

// 导出 API - 开发模式使用 Mock，生产模式使用真实 API
export default USE_MOCK ? mockApi : {
  auth: {
    login: (phone: string, password: string) =>
      instance.post('/auth/login', { phone, password }),
    wechatLogin: (code: string, userInfo: any) =>
      instance.post('/auth/wechat-login', { code, ...userInfo })
  },
  elderly: {
    list: (params: any) => instance.get('/elderly', { params }),
    get: (id: number) => instance.get(`/elderly/${id}`),
    familyList: () => instance.get('/user/elderly-list')
  },
  care: {
    records: (params: any) => instance.get('/care/records', { params }),
    items: () => instance.get('/care/items'),
    myTasks: () => instance.get('/care/my-tasks'),
    createRecord: (data: any) => instance.post('/care/records', data)
  },
  health: {
    records: (params: any) => instance.get('/health/records', { params }),
    create: (data: any) => instance.post('/health/records', data)
  },
  service: {
    requests: (params: any) => instance.get('/service/requests', { params }),
    create: (data: any) => instance.post('/service/requests', data),
    handle: (id: number, data: any) => instance.put(`/service/requests/${id}`, data)
  },
  bill: {
    list: (params: any) => instance.get('/bills', { params }),
    get: (id: number) => instance.get(`/bills/${id}`)
  }
}
