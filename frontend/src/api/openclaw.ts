import request from './request'
import type { OpenClawStatus, OpenClawTool, OpenClawSkill, OpenClawModelConfig, DeployLog, ApiResponse } from '@/types'

/** 获取 OpenClaw 运行状态 */
export function getOpenClawStatus() {
  return request.get<ApiResponse<OpenClawStatus>>('/openclaw/status')
}

/** 检查更新 */
export function checkUpdate() {
  return request.get<ApiResponse<{ hasUpdate: boolean; version?: string; changelog?: string }>>('/openclaw/check-update')
}

/** 执行更新 */
export function performUpdate() {
  return request.post<ApiResponse<null>>('/openclaw/update')
}

/** 重启 OpenClaw 服务 */
export function restartOpenClaw() {
  return request.post<ApiResponse<null>>('/openclaw/restart')
}

/** 停止 OpenClaw 服务 */
export function stopOpenClaw() {
  return request.post<ApiResponse<null>>('/openclaw/stop')
}

/** 获取版本历史 */
export function getVersionHistory() {
  return request.get<ApiResponse<{ versions: Array<{ version: string; date: string; changelog: string }> }>>('/openclaw/versions')
}

/** 获取部署日志 */
export function getDeployLogs(params?: { limit?: number; level?: string }) {
  return request.get<ApiResponse<{ logs: DeployLog[] }>>('/openclaw/logs', { params })
}

/** 一键部署 OpenClaw 到虚拟机 */
export function deployToVM(vmId: string | number) {
  return request.post<ApiResponse<{ taskId: string }>>(`/openclaw/deploy/${vmId}`)
}

/** 获取模型配置 */
export function getModelConfig() {
  return request.get<ApiResponse<OpenClawModelConfig>>('/openclaw/config/model')
}

/** 更新模型配置 */
export function updateModelConfig(data: OpenClawModelConfig) {
  return request.put<ApiResponse<null>>('/openclaw/config/model', data)
}

/** 获取工具列表 */
export function getTools() {
  return request.get<ApiResponse<{ tools: OpenClawTool[] }>>('/openclaw/tools')
}

/** 更新工具状态 */
export function updateTool(name: string, enabled: boolean) {
  return request.put<ApiResponse<null>>(`/openclaw/tools/${name}`, { enabled })
}

/** 获取技能列表 */
export function getSkills() {
  return request.get<ApiResponse<{ skills: OpenClawSkill[] }>>('/openclaw/skills')
}

/** 安装技能 */
export function installSkill(name: string) {
  return request.post<ApiResponse<null>>(`/openclaw/skills/${name}/install`)
}

/** 卸载技能 */
export function uninstallSkill(name: string) {
  return request.delete<ApiResponse<null>>(`/openclaw/skills/${name}`)
}

/** 更新技能 */
export function updateSkill(name: string) {
  return request.post<ApiResponse<null>>(`/openclaw/skills/${name}/update`)
}

/** 切换技能启用状态 */
export function toggleSkill(name: string, enabled: boolean) {
  return request.put<ApiResponse<null>>(`/openclaw/skills/${name}`, { enabled })
}
