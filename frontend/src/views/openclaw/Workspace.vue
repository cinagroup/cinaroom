<template>
  <div class="workspace-container">
    <Header>
      <template #left>
        <span class="page-title">OpenClaw 工作空间管理</span>
      </template>
      <template #right>
        <el-button type="primary" @click="handleCreateWorkspace">
          <el-icon><Plus /></el-icon>
          新建工作空间
        </el-button>
      </template>
    </Header>
    
    <el-row :gutter="20">
      <el-col :span="8" v-for="ws in workspaces" :key="ws.id">
        <el-card class="workspace-card">
          <template #header>
            <div class="card-header">
              <span class="workspace-name">{{ ws.name }}</span>
              <el-tag size="small" :type="ws.current ? 'success' : 'info'">
                {{ ws.current ? '当前' : '历史' }}
              </el-tag>
            </div>
          </template>
          
          <div class="workspace-info">
            <div class="info-item">
              <el-icon><Folder /></el-icon>
              <span class="label">路径:</span>
              <span class="value">{{ ws.path }}</span>
            </div>
            <div class="info-item">
              <el-icon><Document /></el-icon>
              <span class="label">文件数:</span>
              <span class="value">{{ ws.fileCount }}</span>
            </div>
            <div class="info-item">
              <el-icon><Coin /></el-icon>
              <span class="label">大小:</span>
              <span class="value">{{ ws.size }}</span>
            </div>
            <div class="info-item">
              <el-icon><Clock /></el-icon>
              <span class="label">修改时间:</span>
              <span class="value">{{ ws.modifiedAt }}</span>
            </div>
          </div>
          
          <div class="workspace-actions">
            <el-button type="primary" link @click="handleOpen(ws)">进入</el-button>
            <el-button link @click="handleBackup(ws)">备份</el-button>
            <el-button type="danger" link @click="handleDelete(ws)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card style="margin-top: 20px">
      <template #header>
        <span>工作空间统计</span>
      </template>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="工作空间总数">{{ workspaces.length }}</el-descriptions-item>
        <el-descriptions-item label="总文件数">12,580</el-descriptions-item>
        <el-descriptions-item label="总大小">8.5 GB</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Header from '@/components/Header.vue'

const workspaces = ref([
  {
    id: 1,
    name: 'default',
    path: '/root/.openclaw/workspace',
    fileCount: 8520,
    size: '5.2 GB',
    modifiedAt: '2024-04-01 10:30',
    current: true
  },
  {
    id: 2,
    name: 'backup-20240328',
    path: '/root/.openclaw/workspace-backup',
    fileCount: 4060,
    size: '3.3 GB',
    modifiedAt: '2024-03-28 15:00',
    current: false
  }
])

const handleCreateWorkspace = () => {
  ElMessage.info('创建新工作空间')
}

const handleOpen = (ws) => {
  ElMessage.success(`进入工作空间：${ws.name}`)
}

const handleBackup = (ws) => {
  ElMessage.success(`开始备份工作空间：${ws.name}`)
}

const handleDelete = (ws) => {
  if (ws.current) {
    ElMessage.warning('不能删除当前工作空间')
    return
  }
  
  ElMessageBox.confirm(`确定要删除工作空间 "${ws.name}" 吗？`, '警告', { type: 'warning' })
    .then(() => ElMessage.success('删除成功'))
}
</script>

<style scoped lang="scss">
.workspace-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}

.workspace-card {
  margin-bottom: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .workspace-name {
      font-weight: bold;
      font-size: 16px;
    }
  }
  
  .workspace-info {
    .info-item {
      display: flex;
      align-items: center;
      margin-bottom: 10px;
      
      .el-icon {
        margin-right: 8px;
        color: var(--accent-color);
      }
      
      .label {
        color: #909399;
        margin-right: 10px;
      }
      
      .value {
        color: #303133;
      }
    }
  }
  
  .workspace-actions {
    margin-top: 15px;
    display: flex;
    gap: 10px;
  }
}
</style>
