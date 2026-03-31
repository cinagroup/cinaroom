<template>
  <div class="vm-list-container">
    <Header>
      <template #left>
        <span class="page-title">虚拟机管理</span>
      </template>
      <template #right>
        <el-button type="primary" @click="$router.push('/vms/create')">
          <el-icon><Plus /></el-icon>
          新建虚拟机
        </el-button>
      </template>
    </Header>
    
    <el-card>
      <div class="filter-bar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索虚拟机名称或 IP"
          prefix-icon="Search"
          clearable
          style="width: 300px"
        />
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px; margin-left: 10px">
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="已暂停" value="suspended" />
        </el-select>
      </div>
      
      <el-table :data="filteredVMs" style="width: 100%; margin-top: 20px" v-loading="loading">
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址" width="140" />
        <el-table-column prop="cpu" label="CPU" width="100">
          <template #default="{ row }">
            {{ row.cpu }} 核
          </template>
        </el-table-column>
        <el-table-column prop="memory" label="内存" width="100">
          <template #default="{ row }">
            {{ row.memory }} GB
          </template>
        </el-table-column>
        <el-table-column prop="disk" label="磁盘" width="100">
          <template #default="{ row }">
            {{ row.disk }} GB
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button
              :type="row.status === 'running' ? 'warning' : 'success'"
              link
              @click="handleToggleVM(row)"
            >
              {{ row.status === 'running' ? '停止' : '启动' }}
            </el-button>
            <el-button link @click="handleRestart(row)">重启</el-button>
            <el-button link @click="$router.push(`/vms/${row.id}`)">详情</el-button>
            <el-button link @click="$router.push(`/shell/${row.id}`)">Shell</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useVMStore } from '@/stores/vm'
import { ElMessage, ElMessageBox } from 'element-plus'
import Header from '@/components/Header.vue'

const vmStore = useVMStore()
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')

const vmList = ref([
  { id: 1, name: 'ubuntu-dev', status: 'running', ip: '192.168.64.10', cpu: 2, memory: 4, disk: 50 },
  { id: 2, name: 'centos-prod', status: 'stopped', ip: '192.168.64.11', cpu: 4, memory: 8, disk: 100 },
  { id: 3, name: 'debian-test', status: 'running', ip: '192.168.64.12', cpu: 1, memory: 2, disk: 25 }
])

const filteredVMs = computed(() => {
  return vmList.value.filter(vm => {
    const matchSearch = !searchQuery.value || 
      vm.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      vm.ip.includes(searchQuery.value)
    
    const matchStatus = !statusFilter.value || vm.status === statusFilter.value
    
    return matchSearch && matchStatus
  })
})

const getStatusType = (status) => {
  const types = {
    running: 'success',
    stopped: 'info',
    suspended: 'warning'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    running: '运行中',
    stopped: '已停止',
    suspended: '已暂停'
  }
  return texts[status] || status
}

const handleToggleVM = (vm) => {
  const action = vm.status === 'running' ? '停止' : '启动'
  ElMessageBox.confirm(`确定要${action}虚拟机 "${vm.name}" 吗？`, '提示', {
    type: 'warning'
  }).then(() => {
    // TODO: 调用 API
    ElMessage.success(`已发送${action}命令`)
  })
}

const handleRestart = (vm) => {
  ElMessageBox.confirm(`确定要重启虚拟机 "${vm.name}" 吗？`, '提示', {
    type: 'warning'
  }).then(() => {
    // TODO: 调用 API
    ElMessage.success('已发送重启命令')
  })
}

const handleDelete = (vm) => {
  ElMessageBox.confirm(`确定要删除虚拟机 "${vm.name}" 吗？此操作不可恢复！`, '警告', {
    type: 'error',
    confirmButtonText: '删除',
    confirmButtonClass: 'el-button--danger'
  }).then(() => {
    // TODO: 调用 API
    ElMessage.success('删除成功')
  })
}

onMounted(() => {
  loading.value = true
  // TODO: 调用 API 获取虚拟机列表
  setTimeout(() => {
    loading.value = false
  }, 500)
})
</script>

<style scoped lang="scss">
.vm-list-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
  
  .filter-bar {
    display: flex;
    align-items: center;
  }
}
</style>
