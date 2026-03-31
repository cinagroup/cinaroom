<template>
  <div class="log-viewer-container">
    <Header>
      <template #left>
        <span class="page-title">日志查看 - {{ vmName }}</span>
      </template>
    </Header>
    
    <el-card>
      <div class="filter-bar">
        <el-select v-model="logType" placeholder="日志类型" style="width: 150px">
          <el-option label="系统日志" value="system" />
          <el-option label="应用日志" value="application" />
          <el-option label="安全日志" value="security" />
        </el-select>
        
        <el-date-picker
          v-model="dateRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          style="margin: 0 10px"
        />
        
        <el-input
          v-model="keyword"
          placeholder="搜索关键词"
          clearable
          style="width: 200px"
        />
        
        <el-select v-model="logLevel" placeholder="日志级别" clearable style="width: 120px; margin-left: 10px">
          <el-option label="INFO" value="INFO" />
          <el-option label="WARN" value="WARN" />
          <el-option label="ERROR" value="ERROR" />
        </el-select>
        
        <el-button type="primary" @click="handleSearch" style="margin-left: 10px">查询</el-button>
        <el-button @click="handleExport">导出</el-button>
        <el-button :type="autoRefresh ? 'success' : ''" @click="toggleAutoRefresh">
          {{ autoRefresh ? '停止刷新' : '自动刷新' }}
        </el-button>
      </div>
      
      <div class="log-content" ref="logContainer">
        <div v-for="log in filteredLogs" :key="log.id" class="log-item" :class="log.level.toLowerCase()">
          <span class="log-time">{{ log.timestamp }}</span>
          <span class="log-level">{{ log.level }}</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import Header from '@/components/Header.vue'

const route = useRoute()
const logContainer = ref(null)

const vmName = ref('ubuntu-dev')
const logType = ref('system')
const dateRange = ref([])
const keyword = ref('')
const logLevel = ref('')
const autoRefresh = ref(true)

const logs = ref([
  { id: 1, timestamp: '2024-04-01 10:30:15', level: 'INFO', message: 'System started successfully' },
  { id: 2, timestamp: '2024-04-01 10:30:20', level: 'INFO', message: 'Network interface eth0 up' },
  { id: 3, timestamp: '2024-04-01 10:31:05', level: 'WARN', message: 'High memory usage detected: 85%' },
  { id: 4, timestamp: '2024-04-01 10:32:10', level: 'ERROR', message: 'Failed to connect to database' },
  { id: 5, timestamp: '2024-04-01 10:33:00', level: 'INFO', message: 'User login: admin' },
  { id: 6, timestamp: '2024-04-01 10:34:15', level: 'INFO', message: 'Service nginx restarted' },
  { id: 7, timestamp: '2024-04-01 10:35:20', level: 'WARN', message: 'Disk usage above 80%' },
  { id: 8, timestamp: '2024-04-01 10:36:00', level: 'ERROR', message: 'Permission denied: /etc/config' }
])

const filteredLogs = computed(() => {
  return logs.value.filter(log => {
    const matchKeyword = !keyword.value || log.message.toLowerCase().includes(keyword.value.toLowerCase())
    const matchLevel = !logLevel.value || log.level === logLevel.value
    return matchKeyword && matchLevel
  })
})

let refreshInterval = null

const handleSearch = () => {
  ElMessage.success('查询完成')
  scrollToBottom()
}

const handleExport = () => {
  ElMessage.success('日志导出成功')
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

const startAutoRefresh = () => {
  refreshInterval = setInterval(() => {
    // TODO: 获取最新日志
    logs.value.push({
      id: logs.value.length + 1,
      timestamp: new Date().toLocaleString(),
      level: 'INFO',
      message: 'Auto-refresh: New log entry'
    })
    scrollToBottom()
  }, 5000)
}

const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}

const scrollToBottom = () => {
  setTimeout(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  }, 100)
}

onMounted(() => {
  if (autoRefresh.value) {
    startAutoRefresh()
  }
  scrollToBottom()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped lang="scss">
.log-viewer-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
  
  .filter-bar {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 10px;
  }
}

.log-content {
  height: 600px;
  overflow-y: auto;
  background-color: #1e1e1e;
  border-radius: $border-radius;
  padding: 15px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  
  .log-item {
    display: flex;
    gap: 15px;
    padding: 5px 0;
    border-bottom: 1px solid #333;
    
    .log-time {
      color: #569cd6;
      min-width: 160px;
    }
    
    .log-level {
      min-width: 60px;
      font-weight: bold;
      
      &.info {
        color: #4ec9b0;
      }
      
      &.warn {
        color: #dcdcaa;
      }
      
      &.error {
        color: #f44747;
      }
    }
    
    .log-message {
      color: #ce9178;
      word-break: break-all;
    }
  }
}
</style>
