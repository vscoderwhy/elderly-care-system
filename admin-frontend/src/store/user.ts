import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface User {
  id: number
  phone: string
  nickname: string
  avatar: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(newUser: User) {
    user.value = newUser
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    user,
    setToken,
    setUser,
    logout
  }
})
