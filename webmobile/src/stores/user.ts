import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/rpc'
import { wsClient } from '@/utils/websocket'

interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const result = await authApi.login(username, password)
    token.value = result.token
    userInfo.value = result.user
    localStorage.setItem('token', result.token)
    wsClient.connect(result.token)
    return result
  }

  async function register(username: string, password: string, nickname?: string) {
    const result = await authApi.register(username, password, nickname)
    token.value = result.token
    userInfo.value = result.user
    localStorage.setItem('token', result.token)
    wsClient.connect(result.token)
    return result
  }

  async function fetchUserInfo() {
    if (!token.value) return null
    try {
      const user = await authApi.getUserInfo()
      userInfo.value = user
      return user
    } catch {
      logout()
      return null
    }
  }

  function logout() {
    token.value = null
    userInfo.value = null
    localStorage.removeItem('token')
    wsClient.disconnect()
  }

  function initWebSocket() {
    if (token.value) {
      wsClient.connect(token.value)
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    login,
    register,
    fetchUserInfo,
    logout,
    initWebSocket,
  }
})
