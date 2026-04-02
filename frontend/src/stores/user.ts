import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, LoginParams, RegisterParams } from '@/types/user'
import * as authApi from '@/api/auth'
import { isAdmin as checkIsAdmin, isRoot as checkIsRoot } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  // State
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo>(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  const loading = ref(false)

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value.username || '')
  const email = computed(() => userInfo.value.email || '')
  const avatar = computed(() => userInfo.value.avatar || '')
  const roles = computed(() => userInfo.value.roles || [])
  const role = computed(() => (userInfo.value as any).role ?? 0)
  const isAdmin = computed(() => checkIsAdmin(userInfo.value))
  const isRoot = computed(() => checkIsRoot(userInfo.value))

  // Actions
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUserInfo(info: UserInfo) {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  /** 从 JWT 解析 role */
  function parseRoleFromToken(tokenStr: string): number {
    try {
      const parts = tokenStr.split('.')
      if (parts.length !== 3) return 0
      const payload = JSON.parse(atob(parts[1]))
      return payload.role ?? 0
    } catch {
      return 0
    }
  }

  async function login(params: LoginParams) {
    loading.value = true
    try {
      const res = await authApi.login(params)
      const { token: newToken, user } = (res as any).data || res
      // 如果 user 中没有 role，从 JWT 解析
      if (user && !user.role) {
        user.role = parseRoleFromToken(newToken)
      }
      setToken(newToken)
      setUserInfo(user)
      return user
    } finally {
      loading.value = false
    }
  }

  async function register(params: RegisterParams) {
    loading.value = true
    try {
      await authApi.register(params)
    } finally {
      loading.value = false
    }
  }

  async function fetchUserInfo() {
    try {
      const res = await authApi.getUserInfo()
      const user = (res as any).data || res
      // 如果 user 中没有 role，从 token 解析
      if (user && !user.role && token.value) {
        user.role = parseRoleFromToken(token.value)
      }
      setUserInfo(user)
      return user
    } catch {
      return null
    }
  }

  async function updateProfile(data: Partial<UserInfo>) {
    loading.value = true
    try {
      const res = await authApi.updateUserInfo(data)
      const user = (res as any).data || res
      setUserInfo(user)
      return user
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // 忽略退出接口错误
    }
    token.value = ''
    userInfo.value = {} as UserInfo
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  return {
    token,
    userInfo,
    loading,
    isLoggedIn,
    username,
    email,
    avatar,
    roles,
    role,
    isAdmin,
    isRoot,
    setToken,
    setUserInfo,
    login,
    register,
    fetchUserInfo,
    updateProfile,
    logout
  }
})
