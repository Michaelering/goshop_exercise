const TOKEN_KEY = 'goshop_admin_token'
const USER_KEY = 'goshop_admin_user'
const USER_TYPE_KEY = 'goshop_user_type'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string, type?: string) {
  localStorage.setItem(TOKEN_KEY, token)
  if (type) {
    localStorage.setItem(USER_TYPE_KEY, type)
  }
}

export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.removeItem(USER_TYPE_KEY)
}

export function getUser(): any | null {
  const userStr = localStorage.getItem(USER_KEY)
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch {
      return null
    }
  }
  return null
}

export function setUser(user: any) {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function getUserType(): string {
  return localStorage.getItem(USER_TYPE_KEY) || ''
}
