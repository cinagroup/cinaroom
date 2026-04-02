import request from './request'
import type { VM, CreateVMParams, VMStatus, VMSnapshot, VMMetrics, ApiResponse } from '@/types'

/** 获取虚拟机列表 */
export function getVMList(params?: { search?: string; status?: VMStatus; page?: number; pageSize?: number }) {
  return request.get<ApiResponse<{ list: VM[]; total: number }>>('/vms', { params })
}

/** 获取虚拟机详情 */
export function getVMDetail(id: string | number) {
  return request.get<ApiResponse<VM>>(`/vms/${id}`)
}

/** 创建虚拟机 */
export function createVM(data: CreateVMParams) {
  return request.post<ApiResponse<VM>>('/vms', data)
}

/** 更新虚拟机配置 */
export function updateVM(id: string | number, data: Partial<VM>) {
  return request.put<ApiResponse<VM>>(`/vms/${id}`, data)
}

/** 删除虚拟机 */
export function deleteVM(id: string | number) {
  return request.delete<ApiResponse<null>>(`/vms/${id}`)
}

/** 启动虚拟机 */
export function startVM(id: string | number) {
  return request.post<ApiResponse<null>>(`/vms/${id}/start`)
}

/** 停止虚拟机 */
export function stopVM(id: string | number) {
  return request.post<ApiResponse<null>>(`/vms/${id}/stop`)
}

/** 重启虚拟机 */
export function restartVM(id: string | number) {
  return request.post<ApiResponse<null>>(`/vms/${id}/restart`)
}

/** 暂停虚拟机 */
export function suspendVM(id: string | number) {
  return request.post<ApiResponse<null>>(`/vms/${id}/suspend`)
}

/** 获取虚拟机监控数据 */
export function getVMMetrics(id: string | number, params?: { start?: number; end?: number; step?: number }) {
  return request.get<ApiResponse<{ metrics: VMMetrics[] }>>(`/vms/${id}/metrics`, { params })
}

/** 获取虚拟机快照列表 */
export function getVMSnapshots(id: string | number) {
  return request.get<ApiResponse<{ snapshots: VMSnapshot[] }>>(`/vms/${id}/snapshots`)
}

/** 创建虚拟机快照 */
export function createVMSnapshot(id: string | number, data: { name: string }) {
  return request.post<ApiResponse<VMSnapshot>>(`/vms/${id}/snapshots`, data)
}

/** 恢复虚拟机快照 */
export function restoreVMSnapshot(vmId: string | number, snapshotId: string | number) {
  return request.post<ApiResponse<null>>(`/vms/${vmId}/snapshots/${snapshotId}/restore`)
}

/** 删除虚拟机快照 */
export function deleteVMSnapshot(vmId: string | number, snapshotId: string | number) {
  return request.delete<ApiResponse<null>>(`/vms/${vmId}/snapshots/${snapshotId}`)
}

/** 获取可用镜像列表 */
export function getImages() {
  return request.get<ApiResponse<{ images: Array<{ name: string; alias: string; os: string; version: string; size: string }> }>>('/images')
}

/** 获取虚拟机 WebSocket 终端地址 */
export function getShellWsUrl(id: string | number) {
  return request.get<ApiResponse<{ url: string }>>(`/vms/${id}/shell`)
}
