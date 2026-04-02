import request from './request'
import type { LoginParams, RegisterParams, UserInfo, OAuthInfo, OAuthCallbackResponse, OAuthProvider, ApiResponse } from '@/types/user'

/** 用户登录 */
export function login(data: LoginParams) {
  return request.post<ApiResponse<{ token: string; user: UserInfo }>>('/auth/login', data)
}

/** 用户注册 */
export function register(data: RegisterParams) {
  return request.post<ApiResponse<null>>('/auth/register', data)
}

/** 获取当前用户信息 */
export function getUserInfo() {
  return request.get<ApiResponse<UserInfo>>('/auth/me')
}

/** 更新用户信息 */
export function updateUserInfo(data: Partial<UserInfo>) {
  return request.put<ApiResponse<UserInfo>>('/auth/profile', data)
}

/** 修改密码 */
export function changePassword(data: { currentPassword: string; newPassword: string }) {
  return request.post<ApiResponse<null>>('/auth/change-password', data)
}

/** 刷新 Token */
export function refreshToken() {
  return request.post<ApiResponse<{ token: string }>>('/auth/refresh')
}

/** 退出登录 */
export function logout() {
  return request.post<ApiResponse<null>>('/auth/logout')
}

/** 获取 OAuth 授权 URL */
export function getOAuthAuthorizeUrl() {
  return request.get<ApiResponse<OAuthInfo>>('/oauth/authorize')
}

/** OAuth 回调处理 */
export function oauthCallback(code: string, state?: string) {
  return request.get<ApiResponse<OAuthCallbackResponse>>('/oauth/callback', {
    params: { code, state }
  })
}

/** 获取 OAuth 提供商列表 */
export function getOAuthProviders() {
  return request.get<ApiResponse<{ providers: OAuthProvider[] }>>('/oauth/providers')
}

/** 获取登录日志 */
export function getLoginLogs() {
  return request.get<ApiResponse<{ logs: import('@/types/remote').LoginLog[] }>>('/auth/login-logs')
}

/** 获取活跃会话 */
export function getActiveSessions() {
  return request.get<ApiResponse<{ sessions: import('@/types/remote').ActiveSession[] }>>('/auth/sessions')
}

/** 下线指定会话 */
export function killSession(sessionId: string) {
  return request.delete<ApiResponse<null>>(`/auth/sessions/${sessionId}`)
}

/** 启用/禁用两步验证 */
export function toggleTwoFactor(enabled: boolean) {
  return request.post<ApiResponse<{ secret?: string }>>('/auth/two-factor', { enabled })
}

/** 上传头像 */
export function uploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('avatar', file)
  return request.post<ApiResponse<{ url: string }>>('/auth/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
