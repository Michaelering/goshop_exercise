import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getToken, setToken, removeToken, setUser, getUser } from '@/utils/token'
import { adminLogin, merchantLogin, getAdminCurrentUser, type LoginParams } from '@/api/auth'
import { ElMessage } from 'element-plus'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(getToken())
  const user = ref<any>(getUser())
  const userType = ref<string>('') // 'admin' | 'merchant'

  const isLoggedIn = ref(!!token.value)

  async function login(params: LoginParams, type: 'admin' | 'merchant') {
    try {
      const res = type === 'admin'
        ? await adminLogin(params)
        : await merchantLogin(params)

      const { token: newToken, ...userInfo } = res.data
      setToken(newToken)
      setUser(userInfo)
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
    try {
      const res = await getAdminCurrentUser()
      user.value = res.data
      setUser(res.data)
    } catch {
      // ignore
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
