import request from './request'

export interface LoginParams {
  username: string
  password: string
  captchaId: string
  verifyValue: string
}

export interface LoginResult {
  token: string
  userId: number
  username: string
  isSuper?: number
  roleId?: number
  shopName?: string
}

// 构建 form-urlencoded 请求体（后端使用 c.PostForm 解析）
function toFormData(data: Record<string, any>): URLSearchParams {
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(data)) {
    params.append(key, String(value ?? ''))
  }
  return params
}

export function adminLogin(data: LoginParams) {
  return request.post('/admin/login', toFormData(data)) as Promise<{ code: number; data: LoginResult }>
}

export function merchantLogin(data: LoginParams) {
  return request.post('/merchant/login', toFormData(data)) as Promise<{ code: number; data: LoginResult }>
}

export function adminLogout() {
  return request.get('/admin/logout')
}

export function merchantLogout() {
  return request.get('/merchant/logout')
}

export function getCaptcha() {
  return request.get('/public/captcha')
}

export function getAdminCurrentUser() {
  return request.get('/admin/currentUser')
}
