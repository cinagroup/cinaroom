import request from './request'
import type { ApiResponse } from '@/types/user'

/** 系统信息 */
export interface SystemInfo {
  hostname: string
  os: string
  kernel: string
  uptime: number
  cpuModel: string
  cpuCores: number
  totalMemory: number
  usedMemory: number
  totalDisk: number
  usedDisk: number
  loadAvg: [number, number, number]
}

/** 系统监控数据 */
export interface SystemMetrics {
  timestamp: number
  cpuUsage: number
  memoryUsage: number
  diskUsage: number
  networkIn: number
  networkOut: number
  loadAvg: number
  processes: number
}

/** 获取系统信息 */
export function getSystemInfo() {
  return request.get<ApiResponse<SystemInfo>>('/system/info')
}

/** 获取系统监控数据 */
export function getSystemMetrics(params?: { start?: number; end?: number; step?: number }) {
  return request.get<ApiResponse<{ metrics: SystemMetrics[] }>>('/system/metrics', { params })
}

/** 获取虚拟机日志 */
export function getVMLogs(vmId: string | number, params?: {
  type?: string
  level?: string
  keyword?: string
  start?: string
  end?: string
  limit?: number
}) {
  return request.get<ApiResponse<{ logs: Array<{ id: number; timestamp: string; level: string; message: string }> }>>(`/vms/${vmId}/logs`, { params })
}

/** 导出虚拟机日志 */
export function exportVMLogs(vmId: string | number, params?: { type?: string; start?: string; end?: string }) {
  return request.get<ApiResponse<{ url: string }>>(`/vms/${vmId}/logs/export`, { params })
}

/** 获取虚拟机进程列表 */
export function getVMProcesses(vmId: string | number) {
  return request.get<ApiResponse<{ processes: Array<{ pid: number; name: string; cpu: number; memory: number; user: string }> }>>(`/vms/${vmId}/processes`)
}

/** 结束虚拟机进程 */
export function killVMProcess(vmId: string | number, pid: number) {
  return request.delete<ApiResponse<null>>(`/vms/${vmId}/processes/${pid}`)
}

/** 获取虚拟机服务列表 */
export function getVMServices(vmId: string | number) {
  return request.get<ApiResponse<{ services: Array<{ name: string; status: string; description?: string }> }>>(`/vms/${vmId}/services`)
}

/** 操作虚拟机服务 */
export function controlVMService(vmId: string | number, service: string, action: 'start' | 'stop' | 'restart') {
  return request.post<ApiResponse<null>>(`/vms/${vmId}/services/${service}/${action}`)
}

/** 获取虚拟机文件列表 */
export function getVMFiles(vmId: string | number, params?: { path?: string }) {
  return request.get<ApiResponse<{
    path: string
    files: Array<{ name: string; type: 'file' | 'folder'; size: string; modifiedAt: string }>
  }>>(`/vms/${vmId}/files`, { params })
}
