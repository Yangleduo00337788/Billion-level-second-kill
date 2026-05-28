import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api/user'
import { post } from '@/api/request'
import type { User } from '@/types/api'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  async function login(data: { email: string; password: string }) {
    const res = await userApi.login(data)
    token.value = res.data.token
    user.value = res.data.user
    localStorage.setItem('token', res.data.token)
  }

  async function register(data: { username: string; email: string; password: string }) {
    await userApi.register(data)
  }

  async function logout() {
    // Record logout on server
    try {
      await post('/auth/logout')
    } catch {}
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await userApi.getProfile()
      user.value = res.data
    } catch {
      logout()
    }
  }

  async function updateProfile(data: Partial<User>) {
    const res = await userApi.updateProfile(data)
    user.value = res.data
  }

  if (token.value) {
    fetchProfile()
  }

  return { token, user, isAuthenticated, login, register, logout, fetchProfile, updateProfile }
})