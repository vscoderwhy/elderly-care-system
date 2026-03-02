import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<any>(null)

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUserInfo(info: any) {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  function loadUserInfo() {
    const saved = localStorage.getItem('userInfo')
    if (saved) {
      try {
        userInfo.value = JSON.parse(saved)
      } catch (e) {
        console.error('Failed to parse user info', e)
      }
    }
  }

  loadUserInfo()

  return {
    token,
    userInfo,
    setToken,
    setUserInfo,
    logout
  }
})
