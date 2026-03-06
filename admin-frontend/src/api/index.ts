import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'

const instance: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// Request interceptor
instance.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    // 如果是文件下载（blob类型），直接返回响应
    if (response.config.responseType === 'blob') {
      return response
    }
    
    // 普通JSON响应
    const { code, message, data } = response.data

    if (code === 0) {
      return data
    } else {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message))
    }
  },
  (error) => {
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      window.location.href = '/login'
    }
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default instance

// API functions
export const authApi = {
  login: (phone: string, password: string) =>
    instance.post('/auth/login', { phone, password }),

  logout: () => instance.post('/auth/logout')
}

export const elderlyApi = {
  list: (params: any) => instance.get('/elderly', { params }),
  get: (id: number) => instance.get(`/elderly/${id}`),
  create: (data: any) => instance.post('/elderly', data),
  update: (id: number, data: any) => instance.put(`/elderly/${id}`, data),
  delete: (id: number) => instance.delete(`/elderly/${id}`)
}

export const careApi = {
  listRecords: (params: any) => instance.get('/care/records', { params }),
  createRecord: (data: any) => instance.post('/care/records', data),
  listItems: () => instance.get('/care/items')
}

export const billApi = {
  list: (params: any) => instance.get('/bills', { params }),
  get: (id: number) => instance.get(`/bills/${id}`),
  pay: (id: number, data: any) => instance.post(`/bills/${id}/pay`, data)
}
