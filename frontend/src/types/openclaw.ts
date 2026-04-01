/** OpenClaw 运行状态 */
export interface OpenClawStatus {
  running: boolean
  version: string
  uptime: string
  installPath: string
  nodeVersion: string
}

/** OpenClaw 工具 */
export interface OpenClawTool {
  name: string
  description: string
  enabled: boolean
  config?: Record<string, unknown>
}

/** OpenClaw 技能 */
export interface OpenClawSkill {
  name: string
  version: string
  enabled: boolean
  description?: string
}

/** OpenClaw 工作空间 */
export interface OpenClawWorkspace {
  id: number | string
  name: string
  path: string
  fileCount: number
  size: string
  modifiedAt: string
  current: boolean
}

/** OpenClaw 模型配置 */
export interface OpenClawModelConfig {
  model: string
  apiKey: string
  switchStrategy: 'auto' | 'manual'
}

/** 请求日志 */
export interface RequestLog {
  time: string
  type: string
  model: string
  tokens: number
  duration: number
  status: 'success' | 'failed'
}

/** 部署日志 */
export interface DeployLog {
  id: number
  time: string
  level: 'INFO' | 'WARN' | 'ERROR'
  message: string
}
