<template>
  <div class="deploy-container">
    <Header>
      <template #left>
        <span class="page-title">OpenClaw 部署管理</span>
      </template>
    </Header>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>部署状态</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="运行状态">
              <el-tag type="success">运行中</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="版本号">v2026.3.24</el-descriptions-item>
            <el-descriptions-item label="运行时间">15 天 8 小时 30 分钟</el-descriptions-item>
            <el-descriptions-item label="安装路径">/root/.openclaw</el-descriptions-item>
            <el-descriptions-item label="Node 版本">v22.22.1</el-descriptions-item>
          </el-descriptions>
          
          <div class="deploy-actions" style="margin-top: 20px">
            <el-button type="primary" @click="handleUpdate">检查更新</el-button>
            <el-button type="warning" @click="handleRestart">重启服务</el-button>
            <el-button type="danger" @click="handleStop">停止服务</el-button>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>版本历史</span>
          </template>
          <el-timeline>
            <el-timeline-item timestamp="2026-03-24" placement="top" type="success">
              <el-card>
                <h4>v2026.3.24</h4>
                <p>当前版本</p>
              </el-card>
            </el-timeline-item>
            <el-timeline-item timestamp="2026-03-10" placement="top">
              <el-card>
                <h4>v2026.3.10</h4>
                <p>修复已知问题</p>
              </el-card>
            </el-timeline-item>
            <el-timeline-item timestamp="2026-02-28" placement="top">
              <el-card>
                <h4>v2026.2.28</h4>
                <p>新增技能管理功能</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card style="margin-top: 20px">
      <template #header>
        <span>部署日志</span>
      </template>
      <div class="deploy-logs">
        <div v-for="log in deployLogs" :key="log.id" class="log-line">
          <span class="log-time">{{ log.time }}</span>
          <span :class="['log-level', log.level]">{{ log.level }}</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import Header from '@/components/Header.vue'

const deployLogs = ref([
  { id: 1, time: '2026-03-24 10:30:00', level: 'INFO', message: 'OpenClaw Gateway started successfully' },
  { id: 2, time: '2026-03-24 10:30:01', level: 'INFO', message: 'Loading extensions...' },
  { id: 3, time: '2026-03-24 10:30:02', level: 'INFO', message: 'Extension qqbot loaded' },
  { id: 4, time: '2026-03-24 10:30:02', level: 'INFO', message: 'Extension wecom loaded' },
  { id: 5, time: '2026-03-24 10:30:03', level: 'INFO', message: 'Gateway ready on port 3000' }
])

const handleUpdate = () => {
  ElMessage.info('检查更新中...')
  setTimeout(() => {
    ElMessage.success('已是最新版本')
  }, 2000)
}

const handleRestart = () => {
  ElMessage.success('服务重启成功')
}

const handleStop = () => {
  ElMessage.warning('服务已停止')
}
</script>

<style scoped lang="scss">
.deploy-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
  
  .deploy-actions {
    display: flex;
    gap: 10px;
  }
}

.deploy-logs {
  height: 300px;
  overflow-y: auto;
  background-color: #1e1e1e;
  border-radius: $border-radius;
  padding: 15px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  
  .log-line {
    display: flex;
    gap: 15px;
    padding: 3px 0;
    
    .log-time {
      color: #569cd6;
      min-width: 160px;
    }
    
    .log-level {
      min-width: 60px;
      font-weight: bold;
      
      &.INFO {
        color: #4ec9b0;
      }
      
      &.WARN {
        color: #dcdcaa;
      }
      
      &.ERROR {
        color: #f44747;
      }
    }
    
    .log-message {
      color: #ce9178;
    }
  }
}
</style>
