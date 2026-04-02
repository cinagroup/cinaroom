import request from './request'
import type { Mount, CreateMountParams, ApiResponse } from '@/types'

/** 获取挂载列表 */
export function getMountList(params?: { vmId?: string }) {
  return request.get<ApiResponse<{ list: Mount[] }>>('/mounts', { params })
}

/** 获取挂载详情 */
export function getMountDetail(id: string | number) {
  return request.get<ApiResponse<Mount>>(`/mounts/${id}`)
}

/** 创建挂载 */
export function createMount(data: CreateMountParams) {
  return request.post<ApiResponse<Mount>>('/mounts', data)
}

/** 更新挂载 */
export function updateMount(id: string | number, data: Partial<CreateMountParams>) {
  return request.put<ApiResponse<Mount>>(`/mounts/${id}`, data)
}

/** 删除挂载 */
export function deleteMount(id: string | number) {
  return request.delete<ApiResponse<null>>(`/mounts/${id}`)
}

/** 挂载/卸载 */
export function toggleMount(id: string | number, mounted: boolean) {
  return request.post<ApiResponse<null>>(`/mounts/${id}/${mounted ? 'unmount' : 'mount'}`)
}
