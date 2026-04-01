/** 挂载项 */
export interface Mount {
  id: number | string
  name: string
  hostPath: string
  vmPath: string
  vmId: number | string
  vmName: string
  permission: 'rw' | 'ro'
  status: boolean
  autoMount: boolean
}

/** 创建挂载参数 */
export interface CreateMountParams {
  name: string
  hostPath: string
  vmPath: string
  vmId: string
  permission: 'rw' | 'ro'
  autoMount: boolean
}
