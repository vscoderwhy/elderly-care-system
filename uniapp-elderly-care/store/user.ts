import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserInfo, wxLogin, phoneLogin, logout as apiLogout } from '@/api'

interface UserInfo {
  id: string | number
  name: string
  avatar?: string
  phone?: string
  role: 'family' | 'nurse'
  elderlyList?: Array<{
    id: string | number
    name: string
    relation: string
    avatar?: string
  }>
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)
  const isLoggedIn = ref(false)

  // 初始化（从本地存储恢复）
  const init = () => {
    const savedToken = uni.getStorageSync('token')
    const savedUserInfo = uni.getStorageSync('userInfo')

    if (savedToken) {
      token.value = savedToken
      isLoggedIn.value = true
    }

    if (savedUserInfo) {
      userInfo.value = savedUserInfo
    }
  }

  // 微信登录
  const wxLoginAction = async (code: string) => {
    try {
      const res = await wxLogin(code)
      token.value = res.token
      userInfo.value = res.userInfo
      isLoggedIn.value = true

      // 持久化
      uni.setStorageSync('token', res.token)
      uni.setStorageSync('userInfo', res.userInfo)

      return res
    } catch (error) {
      console.error('微信登录失败', error)
      throw error
    }
  }

  // 手机号登录
  const phoneLoginAction = async (phone: string, smsCode: string) => {
    try {
      const res = await phoneLogin(phone, smsCode)
      token.value = res.token
      userInfo.value = res.userInfo
      isLoggedIn.value = true

      // 持久化
      uni.setStorageSync('token', res.token)
      uni.setStorageSync('userInfo', res.userInfo)

      return res
    } catch (error) {
      console.error('登录失败', error)
      throw error
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const res = await getUserInfo()
      userInfo.value = res
      uni.setStorageSync('userInfo', res)
      return res
    } catch (error) {
      console.error('获取用户信息失败', error)
      throw error
    }
  }

  // 更新用户信息
  const updateUserInfo = (info: Partial<UserInfo>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...info }
      uni.setStorageSync('userInfo', userInfo.value)
    }
  }

  // 退出登录
  const logout = async () => {
    try {
      await apiLogout()
    } catch (error) {
      console.error('退出登录请求失败', error)
    } finally {
      // 清除本地数据
      token.value = ''
      userInfo.value = null
      isLoggedIn.value = false
      uni.removeStorageSync('token')
      uni.removeStorageSync('userInfo')

      // 跳转到登录页
      uni.reLaunch({
        url: '/pages/login/index'
      })
    }
  }

  // 检查登录状态
  const checkLogin = () => {
    if (!isLoggedIn.value || !token.value) {
      uni.showToast({
        title: '请先登录',
        icon: 'none'
      })
      uni.navigateTo({
        url: '/pages/login/index'
      })
      return false
    }
    return true
  }

  // 选择关联老人
  const selectedElderlyId = ref<string | number | null>(null)

  const setSelectedElderly = (id: string | number) => {
    selectedElderlyId.value = id
    uni.setStorageSync('selectedElderlyId', id)
  }

  const getSelectedElderly = () => {
    if (!selectedElderlyId.value) {
      selectedElderlyId.value = uni.getStorageSync('selectedElderlyId')
    }
    return selectedElderlyId.value
  }

  return {
    // 状态
    token,
    userInfo,
    isLoggedIn,
    selectedElderlyId,

    // 方法
    init,
    wxLoginAction,
    phoneLoginAction,
    fetchUserInfo,
    updateUserInfo,
    logout,
    checkLogin,
    setSelectedElderly,
    getSelectedElderly
  }
})
