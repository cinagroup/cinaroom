<template>
  <div class="mount-manager-container">
    <Header>
      <template #left>
        <span class="page-title">挂载管理</span>
      </template>
      <template #right>
        <el-button type="primary" @click="showAddDialog = true">
          <el-icon><Plus /></el-icon>
          新建挂载
        </el-button>
      </template>
    </Header>
    
    <el-card>
      <el-table :data="mountList" style="width: 100%">
        <el-table-column prop="name" label="挂载名称" min-width="150" />
        <el-table-column prop="hostPath" label="宿主机路径" min-width="200" />
        <el-table-column prop="vmPath" label="虚拟机路径" min-width="200" />
        <el-table-column prop="vmName" label="所属虚拟机" width="150" />
        <el-table-column prop="permission" label="权限" width="100">
          <template #default="{ row }">
            <el-tag :type="row.permission === 'rw' ? 'success' : 'info'">
              {{ row.permission === 'rw' ? '读写' : '只读' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'info'">
              {{ row.status ? '已挂载' : '未挂载' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="autoMount" label="自动挂载" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.autoMount" disabled />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button
              :type="row.status ? 'warning' : 'success'"
              link
              @click="handleToggleMount(row)"
            >
              {{ row.status ? '卸载' : '挂载' }}
            </el-button>
            <el-button link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="showAddDialog" title="新建挂载" width="600px">
      <el-form :model="mountForm" label-width="120px">
        <el-form-item label="挂载名称">
          <el-input v-model="mountForm.name" placeholder="请输入挂载名称" />
        </el-form-item>
        
        <el-form-item label="宿主机路径">
          <el-input v-model="mountForm.hostPath" placeholder="/path/on/host" />
        </el-form-item>
        
        <el-form-item label="虚拟机路径">
          <el-input v-model="mountForm.vmPath" placeholder="/path/on/vm" />
        </el-form-item>
        
        <el-form-item label="所属虚拟机">
          <el-select v-model="mountForm.vmId" placeholder="请选择虚拟机" style="width: 100%">
            <el-option label="ubuntu-dev" value="1" />
            <el-option label="centos-prod" value="2" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="权限">
          <el-radio-group v-model="mountForm.permission">
            <el-radio label="rw">读写</el-radio>
            <el-radio label="ro">只读</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="自动挂载">
          <el-switch v-model="mountForm.autoMount" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Header from '@/components/Header.vue'

const showAddDialog = ref(false)

const mountList = ref([
  {
    id: 1,
    name: '工作空间',
    hostPath: '/root/.openclaw/workspace',
    vmPath: '/mnt/workspace',
    vmName: 'ubuntu-dev',
    permission: 'rw',
    status: true,
    autoMount: true
  },
  {
    id: 2,
    name: '技能目录',
    hostPath: '/root/.openclaw/skills',
    vmPath: '/mnt/skills',
    vmName: 'ubuntu-dev',
    permission: 'ro',
    status: true,
    autoMount: true
  }
])

const mountForm = reactive({
  name: '',
  hostPath: '',
  vmPath: '',
  vmId: '',
  permission: 'rw',
  autoMount: true
})

const handleToggleMount = (mount) => {
  const action = mount.status ? '卸载' : '挂载'
  ElMessageBox.confirm(`确定要${action} "${mount.name}" 吗？`, '提示', { type: 'warning' })
    .then(() => {
      mount.status = !mount.status
      ElMessage.success(`${action}成功`)
    })
}

const handleEdit = (mount) => {
  mountForm.name = mount.name
  mountForm.hostPath = mount.hostPath
  mountForm.vmPath = mount.vmPath
  mountForm.vmId = '1'
  mountForm.permission = mount.permission
  mountForm.autoMount = mount.autoMount
  showAddDialog.value = true
}

const handleDelete = (mount) => {
  ElMessageBox.confirm(`确定要删除挂载 "${mount.name}" 吗？`, '警告', { type: 'warning' })
    .then(() => {
      const index = mountList.value.findIndex(m => m.id === mount.id)
      if (index > -1) mountList.value.splice(index, 1)
      ElMessage.success('删除成功')
    })
}

const handleCreate = () => {
  mountList.value.push({
    id: Date.now(),
    ...mountForm,
    vmName: 'ubuntu-dev',
    status: false
  })
  showAddDialog.value = false
  ElMessage.success('创建成功')
}
</script>

<style scoped lang="scss">
.mount-manager-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}
</style>
