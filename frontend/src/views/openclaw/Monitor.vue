<template>
  <div class="monitor-container">
    <Header>
      <template #left>
        <span class="page-title">OpenClaw 监控面板</span>
      </template>
    </Header>
    
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #409EFF">
            <el-icon><Cpu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">45%</div>
            <div class="stat-label">CPU 使用率</div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #67C23A">
            <el-icon><Memory /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">1.2GB</div>
            <div class="stat-label">内存使用</div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #E6A23C">
            <el-icon><Folder /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">2.8GB</div>
            <div class="stat-label">磁盘使用</div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #F56C6C">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">12</div>
            <div class="stat-label">活跃会话</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>请求统计</span>
          </template>
          <div ref="requestChart" style="height: 300px"></div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>响应时间</span>
          </template>
          <div ref="responseChart" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card style="margin-top: 20px">
      <template #header>
        <span>今日请求明细</span>
      </template>
      <el-table :data="requestLogs" style="width: 100%">
        <el-table-column prop="time" label="时间" width="160" />
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column prop="model" label="模型" min-width="150" />
        <el-table-column prop="tokens" label="Token 数" width="100" />
        <el-table-column prop="duration" label="耗时" width="100">
          <template #default="{ row }">
            {{ row.duration }}ms
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import Header from '@/components/Header.vue'

const requestChart = ref(null)
const responseChart = ref(null)

const requestLogs = ref([
  { time: '2024-04-01 10:30:00', type: 'chat', model: 'qwencode/qwen3.5-plus', tokens: 1250, duration: 850, status: 'success' },
  { time: '2024-04-01 10:28:15', type: 'tool', model: 'browser', tokens: 0, duration: 1200, status: 'success' },
  { time: '2024-04-01 10:25:30', type: 'chat', model: 'qwencode/qwen3.5-plus', tokens: 2100, duration: 1500, status: 'success' },
  { time: '2024-04-01 10:22:00', type: 'chat', model: 'claude-3-sonnet', tokens: 1800, duration: 2000, status: 'success' },
  { time: '2024-04-01 10:20:10', type: 'tool', model: 'exec', tokens: 0, duration: 300, status: 'success' }
])

const initCharts = () => {
  if (requestChart.value) {
    const requestInstance = echarts.init(requestChart.value)
    requestInstance.setOption({
      title: { text: '每小时请求数' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['08:00', '09:00', '10:00', '11:00', '12:00', '13:00'] },
      yAxis: { type: 'value' },
      series: [{ data: [45, 62, 78, 95, 88, 72], type: 'bar', color: '#409EFF' }]
    })
  }
  
  if (responseChart.value) {
    const responseInstance = echarts.init(responseChart.value)
    responseInstance.setOption({
      title: { text: '平均响应时间 (ms)' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['08:00', '09:00', '10:00', '11:00', '12:00', '13:00'] },
      yAxis: { type: 'value' },
      series: [{ data: [920, 880, 950, 1100, 980, 900], type: 'line', smooth: true, color: '#67C23A' }]
    })
  }
}

onMounted(() => {
  initCharts()
})
</script>

<style scoped lang="scss">
.monitor-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  
  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 20px;
    
    .el-icon {
      font-size: 30px;
      color: #fff;
    }
  }
  
  .stat-info {
    .stat-value {
      font-size: 28px;
      font-weight: bold;
      color: #303133;
    }
    
    .stat-label {
      font-size: 14px;
      color: #909399;
      margin-top: 5px;
    }
  }
}
</style>
