// API 基础配置
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

// 请求拦截器
const requestInterceptor = (config: UniApp.RequestOptions) => {
  // 添加 token
  const token = uni.getStorageSync('token')
  if (token) {
    config.header = {
      ...config.header,
      'Authorization': `Bearer ${token}`
    }
  }

  // 添加时间戳防止缓存
  if (config.url) {
    config.url += config.url.includes('?') ? `&_t=${Date.now()}` : `?_t=${Date.now()}`
  }

  console.log('[Request]', config.method?.toUpperCase(), config.url, config.data)
  return config
}

// 响应拦截器
const responseInterceptor = (response: UniApp.RequestSuccessCallbackResult) => {
  const { statusCode, data } = response

  console.log('[Response]', response)

  // HTTP 状态码检查
  if (statusCode >= 200 && statusCode < 300) {
    // 业务状态码检查
    if (data.code === 0 || data.code === 200) {
      return data.data || data
    } else {
      // 业务错误
      uni.showToast({
        title: data.message || '请求失败',
        icon: 'none',
        duration: 2000
      })
      return Promise.reject(data)
    }
  } else {
    // HTTP 错误
    let errorMsg = '网络请求失败'
    if (statusCode === 401) {
      errorMsg = '登录已过期，请重新登录'
      // 跳转到登录页
      uni.reLaunch({
        url: '/pages/login/index'
      })
    } else if (statusCode === 403) {
      errorMsg = '没有权限访问'
    } else if (statusCode === 404) {
      errorMsg = '请求的资源不存在'
    } else if (statusCode === 500) {
      errorMsg = '服务器错误'
    }

    uni.showToast({
      title: errorMsg,
      icon: 'none',
      duration: 2000
    })

    return Promise.reject({ statusCode, message: errorMsg })
  }
}

// 错误拦截器
const errorInterceptor = (error: any) => {
  console.error('[Request Error]', error)

  let errorMsg = '网络异常'
  if (error.errMsg) {
    if (error.errMsg.includes('timeout')) {
      errorMsg = '请求超时，请检查网络'
    } else if (error.errMsg.includes('fail')) {
      errorMsg = '网络连接失败'
    }
  }

  uni.showToast({
    title: errorMsg,
    icon: 'none',
    duration: 2000
  })

  return Promise.reject(error)
}

// 通用请求方法
interface RequestConfig {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  data?: any
  header?: Record<string, string>
  timeout?: number
}

export const request = <T = any>(config: RequestConfig): Promise<T> => {
  return new Promise((resolve, reject) => {
    // 应用请求拦截器
    const interceptedConfig = requestInterceptor({
      url: BASE_URL + config.url,
      method: config.method || 'GET',
      data: config.data,
      header: {
        'Content-Type': 'application/json',
        ...config.header
      },
      timeout: config.timeout || 30000
    })

    uni.request({
      ...interceptedConfig,
      success: (response) => {
        responseInterceptor(response)
          .then(resolve)
          .catch(reject)
      },
      fail: (error) => {
        errorInterceptor(error).catch(reject)
      }
    })
  })
}

// GET 请求
export const get = <T = any>(url: string, data?: any, config?: Partial<RequestConfig>): Promise<T> => {
  return request<T>({
    url,
    method: 'GET',
    data,
    ...config
  })
}

// POST 请求
export const post = <T = any>(url: string, data?: any, config?: Partial<RequestConfig>): Promise<T> => {
  return request<T>({
    url,
    method: 'POST',
    data,
    ...config
  })
}

// PUT 请求
export const put = <T = any>(url: string, data?: any, config?: Partial<RequestConfig>): Promise<T> => {
  return request<T>({
    url,
    method: 'PUT',
    data,
    ...config
  })
}

// DELETE 请求
export const del = <T = any>(url: string, data?: any, config?: Partial<RequestConfig>): Promise<T> => {
  return request<T>({
    url,
    method: 'DELETE',
    data,
    ...config
  })
}

// 文件上传
export const upload = (url: string, filePath: string, name = 'file', formData?: Record<string, string>): Promise<any> => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')

    uni.uploadFile({
      url: BASE_URL + url,
      filePath,
      name,
      formData,
      header: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (response) => {
        const data = JSON.parse(response.data)
        if (data.code === 0 || data.code === 200) {
          resolve(data.data || data)
        } else {
          uni.showToast({
            title: data.message || '上传失败',
            icon: 'none'
          })
          reject(data)
        }
      },
      fail: (error) => {
        uni.showToast({
          title: '上传失败',
          icon: 'none'
        })
        reject(error)
      }
    })
  })
}

// 下载文件
export const download = (url: string, fileName?: string): Promise<any> => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')

    uni.downloadFile({
      url: BASE_URL + url,
      header: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (response) => {
        if (response.statusCode === 200) {
          resolve(response.tempFilePath)
        } else {
          uni.showToast({
            title: '下载失败',
            icon: 'none'
          })
          reject(response)
        }
      },
      fail: reject
    })
  })
}

export default {
  request,
  get,
  post,
  put,
  del: del,
  upload,
  download
}
