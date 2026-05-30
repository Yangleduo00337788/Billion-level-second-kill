import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, logout as logoutApi } from '../api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('accessToken') || '')
  const refreshToken = ref(localStorage.getItem('refreshToken') || '')
  const username = ref(localStorage.getItem('username') || '')
  const userId = ref(localStorage.getItem('userId') || '')

  function setToken(accessToken, refresh) {
    token.value = accessToken
    refreshToken.value = refresh
    localStorage.setItem('accessToken', accessToken)
    localStorage.setItem('refreshToken', refresh)
  }

  function setUsername(name) {
    username.value = name
    localStorage.setItem('username', name)
  }

  function setUserId(id) {
    userId.value = id
    localStorage.setItem('userId', id)
  }

  async function login(loginForm) {
    const res = await loginApi(loginForm)
    const tokenData = res.data.token || res.data
    setToken(tokenData.accessToken, tokenData.refreshToken)
    setUserId(res.data.userId)
    setUsername(res.data.username || loginForm.username)
    return res
  }

  async function logout() {
    try {
      await logoutApi()
    } catch (e) {
      console.error('logout error', e)
    }
    token.value = ''
    refreshToken.value = ''
    username.value = ''
    userId.value = ''
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('username')
    localStorage.removeItem('userId')
  }

  function isLoggedIn() {
    return !!token.value
  }

  return {
    token,
    refreshToken,
    username,
    userId,
    setToken,
    setUsername,
    setUserId,
    login,
    logout,
    isLoggedIn
  }
})
