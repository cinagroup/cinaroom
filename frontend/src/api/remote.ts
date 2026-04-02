import request from './request'
import type { RemoteAccessConfig, WhitelistEntry, AccessLog, ApiResponse } from '@/types'

/** 获取远程访问配置 */
export function getRemoteAccessConfig() {
  return request.get<ApiResponse<RemoteAccessConfig>>('/remote/config')
}

/** 更新远程访问配置 */
export function updateRemoteAccessConfig(data: Partial<RemoteAccessConfig>) {
  return request.put<ApiResponse<null>>('/remote/config', data)
}

/** 获取 IP 白名单 */
export function getWhitelist() {
  return request.get<ApiResponse<{ list: WhitelistEntry[] }>>('/remote/whitelist')
}

/** 添加 IP 白名单 */
export function addWhitelistEntry(data: { ip: string; note?: string }) {
  return request.post<ApiResponse<WhitelistEntry>>('/remote/whitelist', data)
}

/** 批量导入 IP 白名单 */
export function batchImportWhitelist(data: { entries: Array<{ ip: string; note?: string }> }) {
  return request.post<ApiResponse<{ imported: number }>>('/remote/whitelist/batch', data)
}

/** 删除 IP 白名单 */
export function removeWhitelistEntry(ip: string) {
  return request.delete<ApiResponse<null>>(`/remote/whitelist/${encodeURIComponent(ip)}`)
}

/** 获取访问日志 */
export function getAccessLogs(params?: { start?: string; end?: string; page?: number; pageSize?: number }) {
  return request.get<ApiResponse<{ logs: AccessLog[]; total: number }>>('/remote/access-logs', { params })
}

/** 导出访问日志 */
export function exportAccessLogs(params?: { start?: string; end?: string }) {
  return request.get<ApiResponse<{ url: string }>>('/remote/access-logs/export', { params })
}
