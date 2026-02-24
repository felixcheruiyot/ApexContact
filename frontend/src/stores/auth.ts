import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('ls_token'))
  const user = ref<User | null>(JSON.parse(localStorage.getItem('ls_user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(email: string, password: string) {
    const res = await authApi.login({ email, password })
    const { token: t, user: u } = res.data.data!
    token.value = t
    user.value = u
    localStorage.setItem('ls_token', t)
    localStorage.setItem('ls_user', JSON.stringify(u))
  }

  async function register(email: string, password: string, full_name: string) {
    const res = await authApi.register({ email, password, full_name })
    const { token: t, user: u } = res.data.data!
    token.value = t
    user.value = u
    localStorage.setItem('ls_token', t)
    localStorage.setItem('ls_user', JSON.stringify(u))
  }

  async function loginWithGoogle(code: string, redirectUri: string) {
    const res = await authApi.googleCallback({ code, redirect_uri: redirectUri })
    const { token: t, user: u } = res.data.data!
    token.value = t
    user.value = u
    localStorage.setItem('ls_token', t)
    localStorage.setItem('ls_user', JSON.stringify(u))
  }

  function logout() {
    authApi.logout().catch(() => {}) // fire and forget
    token.value = null
    user.value = null
    localStorage.removeItem('ls_token')
    localStorage.removeItem('ls_user')
  }

  return { token, user, isAuthenticated, isAdmin, login, register, loginWithGoogle, logout }
})
