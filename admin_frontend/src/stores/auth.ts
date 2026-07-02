import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getToken, setToken, removeToken, setUser, getUser } from '@/utils/token'
import { adminLogin, merchantLogin, getAdminCurrentUser, type LoginParams } from '@/api/auth'
import request from '@/api/request'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(getToken())
  const user = ref<any>(getUser())
  const userType = ref<string>(getStoredUserType())

  const isLoggedIn = ref(!!token.value)

  async function login(params: LoginParams, type: 'admin' | 'merchant') {
    try {
      const res = type === 'admin'
        ? await adminLogin(params)
        : await merchantLogin(params)

      const { token: newToken, ...userInfo } = res.data
      setToken(newToken, type)
      setUser(userInfo)
      setUserType(type)
      token.value = newToken
      user.value = userInfo
      userType.value = type
      isLoggedIn.value = true
      return true
    } catch (err) {
      return false
    }
  }

  async function fetchCurrentUser() {
    if (userType.value === 'admin') {
      try {
        const res = await getAdminCurrentUser()
        user.value = res.data
        setUser(res.data)
      } catch {
        // ignore
      }
    } else if (userType.value === 'merchant') {
      try {
        const res: any = await request.get('/merchant/currentUser')
        user.value = res.data
        setUser(res.data)
      } catch {
        // ignore
      }
    }
  }

  function logout() {
    removeToken()
    token.value = null
    user.value = null
    userType.value = ''
    isLoggedIn.value = false
  }

  return {
    token,
    user,
    userType,
    isLoggedIn,
    login,
    logout,
    fetchCurrentUser,
  }
})

// 持久化 userType 到 localStorage（防止刷新后丢失）
const USER_TYPE_KEY = 'goshop_user_type'

function getStoredUserType(): string {
  return localStorage.getItem(USER_TYPE_KEY) || ''
}

function setUserType(type: string) {
  localStorage.setItem(USER_TYPE_KEY, type)
}
