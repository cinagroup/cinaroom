import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useVMStore = defineStore('vm', () => {
  const vmList = ref([])
  const currentVM = ref(null)
  const loading = ref(false)
  
  function setVMList(list) {
    vmList.value = list
  }
  
  function setCurrentVM(vm) {
    currentVM.value = vm
  }
  
  function updateVMStatus(id, status) {
    const vm = vmList.value.find(v => v.id === id)
    if (vm) {
      vm.status = status
    }
    if (currentVM.value && currentVM.value.id === id) {
      currentVM.value.status = status
    }
  }
  
  return {
    vmList,
    currentVM,
    loading,
    setVMList,
    setCurrentVM,
    updateVMStatus
  }
})
