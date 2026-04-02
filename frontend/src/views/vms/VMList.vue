<template>
  <div class="vm-list-container page-container">
    <!-- Page Header -->
    <div class="page-header">
      <h2>虚拟机管理</h2>
      <el-button type="primary" @click="$router.push('/vms/create')">
        <el-icon><Plus /></el-icon>新建虚拟机
      </el-button>
    </div>

    <el-card>
      <!-- Filter Bar -->
      <div class="filter-bar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索名称或 IP 地址"
          :prefix-icon="Search"
          clearable
          style="width: 280px"
          @clear="fetchList"
          @keyup.enter="fetchList"
        />
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 140px" @change="fetchList">
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="已暂停" value="suspended" />
          <el-option label="创建中" value="creating" />
          <el-option label="错误" value="error" />
        </el-select>
        <el-button @click="fetchList">
          <el-icon><Refresh /></el-icon>刷新
        </el-button>
      </div>

      <!-- Table -->
      <el-table
        :data="vmStore.vmList"
        style="width: 100%"
        v-loading="vmStore.loading"
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="name" label="名称" min-width="150" sortable show-overflow-tooltip>
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/vms/${row.id}`)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110" sortable>
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small" effect="dark">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址" width="140">
          <template #default="{ row }">
            <span class="mono">{{ row.ip || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="image" label="镜像" min-width="140" show-overflow-tooltip />
        <el-table-column prop="cpu" label="CPU" width="80" align="center">
          <template #default="{ row }">{{ row.cpu }} 核</template>
        </el-table-column>
        <el-table-column prop="memory" label="内存" width="80" align="center">
          <template #default="{ row }">{{ row.memory }} GB</template>
        </el-table-column>
        <el-table-column prop="disk" label="磁盘" width="80" align="center">
          <template #default="{ row }">{{ row.disk }} GB</template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170" sortable>
          <template #default="{ row }">
            <span class="mono">{{ formatTime(row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button-group size="small">
              <el-button
                :type="row.status === 'running' ? 'warning' : 'success'"
                :loading="operationLoading[row.id]"
                @click="handleToggle(row)"
              >
                {{ row.status === 'running' ? '停止' : '启动' }}
              </el-button>
              <el-button
                :disabled="row.status !== 'running'"
                @click="handleRestart(row)"
              >
                重启
              </el-button>
              <el-button @click="$router.push(`/vms/${row.id}`)">详情</el-button>
              <el-button @click="$router.push(`/shell/${row.id}`)">Shell</el-button>
              <el-button type="danger" @click="handleDelete(row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="vmStore.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useVMStore } from '@/stores/vm'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import type { VM, VMStatus } from '@/types'

const router = useRouter()
const vmStore = useVMStore()

const searchQuery = ref('')
const statusFilter = ref<VMStatus | ''>('')
const currentPage = ref(1)
const pageSize = ref(10)
const operationLoading = reactive<Record<string | number, boolean>>({})

const statusTypeMap: Record<string, string> = {
  running: 'success',
  stopped: 'info',
  suspended: 'warning',
  creating: '',
  error: 'danger'
}

const statusTextMap: Record<string, string> = {
  running: '运行中',
  stopped: '已停止',
  suspended: '已暂停',
  creating: '创建中',
  error: '错误'
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const fetchList = () => {
  vmStore.fetchVMList({
    search: searchQuery.value || undefined,
    status: statusFilter.value || undefined,
    page: currentPage.value,
    pageSize: pageSize.value
  })
}

const handleSortChange = ({ prop, order }: { prop: string; order: string | null }) => {
  // Could pass sort params to API later
  console.log('Sort:', prop, order)
}

const handleToggle = async (vm: VM) => {
  const action = vm.status === 'running' ? '停止' : '启动'
  try {
    await ElMessageBox.confirm(`确定要${action}虚拟机 "${vm.name}" 吗？`, '提示', { type: 'warning' })
    operationLoading[vm.id] = true
    if (vm.status === 'running') {
      await vmStore.stopVM(vm.id)
    } else {
      await vmStore.startVM(vm.id)
    }
    ElMessage.success(`已发送${action}命令`)
  } catch {
    // Cancelled
  } finally {
    operationLoading[vm.id] = false
  }
}

const handleRestart = async (vm: VM) => {
  try {
    await ElMessageBox.confirm(`确定要重启虚拟机 "${vm.name}" 吗？`, '提示', { type: 'warning' })
    operationLoading[vm.id] = true
    await vmStore.restartVM(vm.id)
    ElMessage.success('已发送重启命令')
  } catch {
    // Cancelled
  } finally {
    operationLoading[vm.id] = false
  }
}

const handleDelete = async (vm: VM) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除虚拟机 "${vm.name}" 吗？此操作不可恢复！`,
      '警告',
      { type: 'error', confirmButtonText: '确认删除', confirmButtonClass: 'el-button--danger' }
    )
    await vmStore.deleteVM(vm.id)
    ElMessage.success('虚拟机已删除')
  } catch {
    // Cancelled
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.vm-list-container {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      font-size: 22px;
      font-weight: 600;
      color: #303133;
    }
  }

  .filter-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .mono {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
  }
}
</style>
