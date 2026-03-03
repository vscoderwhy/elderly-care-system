import { defineStore } from 'pinia'
import { ref } from 'vue'
import { request } from '@/api'

interface CaregiverUser {
  id: string
  name: string
  phone: string
  avatar: string
  role: 'caregiver'
  employeeNo: string
  department: string
  position: string
  token: string
}

export const useCaregiverStore = defineStore('caregiver', () => {
  const user = ref<CaregiverUser | null>(null)
  const token = ref('')
  const currentElderly = ref<any>(null)

  // 登录
  const login = async (loginData: { phone: string; code: string }) => {
    try {
      const res = await request({
        url: '/auth/caregiver/login',
        method: 'POST',
        data: loginData
      })

      user.value = {
        id: res.data.id,
        name: res.data.name,
        phone: res.data.phone,
        avatar: res.data.avatar,
        role: 'caregiver',
        employeeNo: res.data.employeeNo,
        department: res.data.department,
        position: res.data.position,
        token: res.data.token
      }
      token.value = res.data.token

      // 持久化存储
      uni.setStorageSync('caregiver_token', res.data.token)
      uni.setStorageSync('caregiver_user', res.data)

      return res
    } catch (error) {
      throw error
    }
  }

  // 微信登录
  const wechatLogin = async (wechatData: any) => {
    try {
      const res = await request({
        url: '/auth/caregiver/wechat-login',
        method: 'POST',
        data: wechatData
      })

      user.value = {
        id: res.data.id,
        name: res.data.name,
        phone: res.data.phone,
        avatar: res.data.avatar,
        role: 'caregiver',
        employeeNo: res.data.employeeNo,
        department: res.data.department,
        position: res.data.position,
        token: res.data.token
      }
      token.value = res.data.token

      uni.setStorageSync('caregiver_token', res.data.token)
      uni.setStorageSync('caregiver_user', res.data)

      return res
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logout = async () => {
    try {
      await request({
        url: '/auth/caregiver/logout',
        method: 'POST'
      })
    } catch (error) {
      console.error('登出失败:', error)
    } finally {
      user.value = null
      token.value = ''
      currentElderly.value = null
      uni.removeStorageSync('caregiver_token')
      uni.removeStorageSync('caregiver_user')
    }
  }

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const res = await request({
        url: '/caregiver/info',
        method: 'GET'
      })

      user.value = {
        ...user.value!,
        ...res.data
      }

      uni.setStorageSync('caregiver_user', user.value)

      return res
    } catch (error) {
      throw error
    }
  }

  // 选择当前服务的老人
  const selectElderly = (elderly: any) => {
    currentElderly.value = elderly
    uni.setStorageSync('current_elderly', elderly)
  }

  // 从本地存储恢复状态
  const restoreState = () => {
    const savedToken = uni.getStorageSync('caregiver_token')
    const savedUser = uni.getStorageSync('caregiver_user')
    const savedElderly = uni.getStorageSync('current_elderly')

    if (savedToken && savedUser) {
      token.value = savedToken
      user.value = savedUser
    }

    if (savedElderly) {
      currentElderly.value = savedElderly
    }
  }

  return {
    user,
    token,
    currentElderly,
    login,
    wechatLogin,
    logout,
    getUserInfo,
    selectElderly,
    restoreState
  }
})
