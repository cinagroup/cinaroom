<template>
  <div class="remote-control-container">
    <Header>
      <template #left>
        <span class="page-title">远程管控 - {{ vmName }}</span>
      </template>
    </Header>
    
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="文件管理" name="files">
        <div class="file-manager">
          <div class="file-toolbar">
            <el-button type="primary" size="small" @click="handleUpload">
              <el-icon><Upload /></el-icon>
              上传
            </el-button>
            <el-button size="small" @click="handleNewFolder">
              <el-icon><FolderAdd /></el-icon>
              新建文件夹
            </el-button>
            <el-button size="small" @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
          
          <div class="file-path">
            <el-icon><FolderOpened /></el-icon>
            <span>{{ currentPath }}</span>
          </div>
          
          <el-table :data="fileList" style="width: 100%" @row-dblclick="handleOpenFolder">
            <el-table-column prop="name" label="名称" min-width="200">
              <template #default="{ row }">
                <el-icon><component :is="row.type === 'folder' ? Folder : Document" /></el-icon>
                <span style="margin-left: 8px">{{ row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="size" label="大小" width="120" />
            <el-table-column prop="modifiedAt" label="修改时间" width="180" />
            <el-table-column label="操作" width="180">
              <template #default="{ row }">
                <el-button type="primary" link @click.stop="handleDownload(row)">下载</el-button>
                <el-button type="danger" link @click.stop="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="进程管理" name="process">
        <el-table :data="processList" style="width: 100%">
          <el-table-column prop="pid" label="PID" width="80" />
          <el-table-column prop="name" label="进程名" min-width="200" />
          <el-table-column prop="cpu" label="CPU%" width="80" />
          <el-table-column prop="memory" label="内存%" width="80" />
          <el-table-column prop="user" label="用户" width="100" />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button type="danger" link @click="handleKillProcess(row)">结束</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      
      <el-tab-pane label="服务管理" name="service">
        <el-table :data="serviceList" style="width: 100%">
          <el-table-column prop="name" label="服务名" min-width="200" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'running' ? 'success' : 'info'">
                {{ row.status === 'running' ? '运行中' : '已停止' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button 
                :type="row.status === 'running' ? 'warning' : 'success'" 
                link 
                @click="handleToggleService(row)"
              >
                {{ row.status === 'running' ? '停止' : '启动' }}
              </el-button>
              <el-button link @click="handleRestartService(row)">重启</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Header from '@/components/Header.vue'

const route = useRoute()
const activeTab = ref('files')
const vmName = ref('ubuntu-dev')
const currentPath = ref('/home/ubuntu')

const fileList = ref([
  { name: 'documents', type: 'folder', size: '-', modifiedAt: '2024-04-01 10:00' },
  { name: 'downloads', type: 'folder', size: '-', modifiedAt: '2024-04-01 09:30' },
  { name: 'project.zip', type: 'file', size: '2.5 MB', modifiedAt: '2024-04-01 08:15' },
  { name: 'config.json', type: 'file', size: '1.2 KB', modifiedAt: '2024-03-31 16:20' }
])

const processList = ref([
  { pid: 1, name: 'systemd', cpu: 0.1, memory: 0.5, user: 'root' },
  { pid: 1234, name: 'nginx', cpu: 2.3, memory: 1.2, user: 'www-data' },
  { pid: 5678, name: 'node', cpu: 15.6, memory: 8.5, user: 'ubuntu' }
])

const serviceList = ref([
  { name: 'nginx', status: 'running' },
  { name: 'mysql', status: 'running' },
  { name: 'docker', status: 'stopped' }
])

const handleUpload = () => {
  ElMessage.info('上传功能开发中')
}

const handleNewFolder = () => {
  ElMessage.info('新建文件夹功能开发中')
}

const handleRefresh = () => {
  ElMessage.success('刷新成功')
}

const handleOpenFolder = (row) => {
  if (row.type === 'folder') {
    currentPath.value += '/' + row.name
    // TODO: 加载新目录内容
  }
}

const handleDownload = (file) => {
  ElMessage.success(`开始下载：${file.name}`)
}

const handleDelete = (file) => {
  ElMessageBox.confirm(`确定要删除 "${file.name}" 吗？`, '警告', { type: 'warning' })
    .then(() => ElMessage.success('删除成功'))
}

const handleKillProcess = (process) => {
  ElMessageBox.confirm(`确定要结束进程 "${process.name}" (PID: ${process.pid}) 吗？`, '警告', { type: 'warning' })
    .then(() => ElMessage.success('进程已结束'))
}

const handleToggleService = (service) => {
  const action = service.status === 'running' ? '停止' : '启动'
  ElMessage.success(`已${action}服务：${service.name}`)
}

const handleRestartService = (service) => {
  ElMessage.success(`已重启服务：${service.name}`)
}
</script>

<style scoped lang="scss">
.remote-control-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}

.file-manager {
  .file-toolbar {
    margin-bottom: 15px;
    display: flex;
    gap: 10px;
  }
  
  .file-path {
    display: flex;
    align-items: center;
    padding: 10px 15px;
    background-color: var(--bg-secondary);
    border-radius: $border-radius;
    margin-bottom: 15px;
    
    .el-icon {
      margin-right: 8px;
      color: var(--accent-color);
    }
  }
}
</style>
