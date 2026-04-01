/** 用户信息 */
export interface UserInfo {
  id?: number | string
  username: string
  email: string
  nickname?: string
  phone?: string
  avatar?: string
  createdAt?: string
  lastLoginAt?: string
  roles?: string[]
}

/** 登录参数 */
export interface LoginParams {
  username: string
  password: string
  remember?: boolean
}

/** 注册参数 */
export interface RegisterParams {
  username: string
  email: string
  password: string
  confirmPassword: string
}

/** API 响应 */
export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

/** OAuth 授权信息 */
export interface OAuthInfo {
  authorize_url: string
  provider: string
}

/** OAuth 回调响应 */
export interface OAuthCallbackResponse {
  token: string
  user: UserInfo
  oauth?: {
    provider: string
    providerUserId: string
  }
}

/** OAuth 提供商 */
export interface OAuthProvider {
  name: string
  display_name: string
  enabled: boolean
}
