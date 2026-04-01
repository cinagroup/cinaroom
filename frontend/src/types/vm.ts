/** 虚拟机状态 */
export type VMStatus = 'running' | 'stopped' | 'suspended' | 'creating' | 'error'

/** 虚拟机实例 */
export interface VM {
  id: number | string
  name: string
  status: VMStatus
  ip: string
  image: string
  cpu: number
  memory: number
  disk: number
  network: 'nat' | 'bridge'
  createdAt: string
  updatedAt?: string
}

/** 创建虚拟机参数 */
export interface CreateVMParams {
  name: string
  image: string
  cpu: number
  memory: number
  disk: number
  network: 'nat' | 'bridge'
  sshKey?: string
  initScript?: string
}

/** 虚拟机快照 */
export interface VMSnapshot {
  id: number | string
  name: string
  vmId: number | string
  createdAt: string
  size?: string
}

/** 虚拟机监控数据 */
export interface VMMetrics {
  timestamp: number
  cpuUsage: number
  memoryUsage: number
  diskUsage: number
  networkIn?: number
  networkOut?: number
}
