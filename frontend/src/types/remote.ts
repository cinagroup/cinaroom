/** 远程访问配置 */
export interface RemoteAccessConfig {
  enabled: boolean
  url: string
  sslEnabled: boolean
}

/** IP 白名单条目 */
export interface WhitelistEntry {
  ip: string
  note: string
  addedAt: string
}

/** 访问日志 */
export interface AccessLog {
  time: string
  ip: string
  path: string
  userAgent: string
  status: number
}

/** 登录日志 */
export interface LoginLog {
  time: string
  ip: string
  location: string
  device: string
  status: string
}

/** 活跃会话 */
export interface ActiveSession {
  device: string
  location: string
  loginTime: string
  ip: string
  current: boolean
}
