import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { VM, VMStatus, CreateVMParams } from '@/types'
import * as vmApi from '@/api/vm'

export const useVMStore = defineStore('vm', () => {
  // State
  const vmList = ref<VM[]>([])
  const currentVM = ref<VM | null>(null)
  const loading = ref(false)
  const total = ref(0)

  // Actions
  async function fetchVMList(params?: { search?: string; status?: VMStatus; page?: number; pageSize?: number }) {
    loading.value = true
    try {
      const res = await vmApi.getVMList(params)
      const data = (res as any).data || res
      vmList.value = data.list || []
      total.value = data.total || 0
      return data
    } finally {
      loading.value = false
    }
  }

  async function fetchVMDetail(id: string | number) {
    loading.value = true
    try {
      const res = await vmApi.getVMDetail(id)
      const data = (res as any).data || res
      currentVM.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  async function createVM(params: CreateVMParams) {
    loading.value = true
    try {
      const res = await vmApi.createVM(params)
      const data = (res as any).data || res
      vmList.value.unshift(data)
      return data
    } finally {
      loading.value = false
    }
  }

  async function deleteVM(id: string | number) {
    await vmApi.deleteVM(id)
    vmList.value = vmList.value.filter(vm => vm.id !== id)
    if (currentVM.value?.id === id) {
      currentVM.value = null
    }
  }

  async function startVM(id: string | number) {
    await vmApi.startVM(id)
    updateVMStatus(id, 'running')
  }

  async function stopVM(id: string | number) {
    await vmApi.stopVM(id)
    updateVMStatus(id, 'stopped')
  }

  async function restartVM(id: string | number) {
    await vmApi.restartVM(id)
    // 重启后状态可能暂时变为 creating 或维持 running
  }

  function updateVMStatus(id: string | number, status: VMStatus) {
    const vm = vmList.value.find(v => v.id === id)
    if (vm) {
      vm.status = status
    }
    if (currentVM.value?.id === id) {
      currentVM.value.status = status
    }
  }

  function clearCurrentVM() {
    currentVM.value = null
  }

  return {
    vmList,
    currentVM,
    loading,
    total,
    fetchVMList,
    fetchVMDetail,
    createVM,
    deleteVM,
    startVM,
    stopVM,
    restartVM,
    updateVMStatus,
    clearCurrentVM
  }
})
