<template>
  <div class="vm-detail-container">
    <Header>
      <template #left>
        <span class="page-title">虚拟机详情</span>
      </template>
      <template #right>
        <el-button-group>
          <el-button type="success" @click="handleStart">启动</el-button>
          <el-button type="warning" @click="handleStop">停止</el-button>
          <el-button @click="handleRestart">重启</el-button>
        </el-button-group>
      </template>
    </Header>
    
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card style="margin-bottom: 20px">
          <template #header>
            <span>基本信息</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="名称">{{ vm.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(vm.status)">{{ getStatusText(vm.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="IP 地址">{{ vm.ip }}</el-descriptions-item>
            <el-descriptions-item label="镜像">{{ vm.image }}</el-descriptions-item>
            <el-descriptions-item label="CPU">{{ vm.cpu }} 核</el-descriptions-item>
            <el-descriptions-item label="内存">{{ vm.memory }} GB</el-descriptions-item>
            <el-descriptions-item label="磁盘">{{ vm.disk }} GB</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ vm.createdAt }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
        
        <el-card style="margin-bottom: 20px">
          <template #header>
            <span>资源监控</span>
          </template>
          <div ref="cpuChart" style="height: 300px; margin-bottom: 20px"></div>
          <div ref="memoryChart" style="height: 300px"></div>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card style="margin-bottom: 20px">
          <template #header>
            <span>快捷操作</span>
          </template>
          <el-space direction="vertical" style="width: 100%">
            <el-button style="width: 100%" @click="$router.push(`/shell/${vm.id}`)">
              <el-icon><Terminal /></el-icon>
              WebShell
            </el-button>
            <el-button style="width: 100%" @click="$router.push(`/remote/${vm.id}`)">
              <el-icon><Setting /></el-icon>
              远程管控
            </el-button>
            <el-button style="width: 100%" @click="$router.push(`/logs/${vm.id}`)">
              <el-icon><Document /></el-icon>
              日志查看
            </el-button>
            <el-button style="width: 100%" @click="handleSnapshot">
              <el-icon><Camera /></el-icon>
              创建快照
            </el-button>
          </el-space>
        </el-card>
        
        <el-card>
          <template #header>
            <span>快照列表</span>
          </template>
          <el-empty v-if="snapshots.length === 0" description="暂无快照" />
          <el-table v-else :data="snapshots" style="width: 100%">
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="createdAt" label="创建时间" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleRestoreSnapshot(row)">恢复</el-button>
                <el-button type="danger" link @click="handleDeleteSnapshot(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import Header from '@/components/Header.vue'

const route = useRoute()
const cpuChart = ref(null)
const memoryChart = ref(null)

const vm = reactive({
  id: route.params.id,
  name: 'ubuntu-dev',
  status: 'running',
  ip: '192.168.64.10',
  image: 'Ubuntu 22.04 LTS',
  cpu: 2,
  memory: 4,
  disk: 50,
  createdAt: '2024-01-15 10:30:00'
})

const snapshots = ref([
  { id: 1, name: 'backup-20240401', createdAt: '2024-04-01 08:00:00' },
  { id: 2, name: 'before-update', createdAt: '2024-03-28 15:30:00' }
])

let chartInterval = null

const getStatusType = (status) => {
  const types = { running: 'success', stopped: 'info', suspended: 'warning' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { running: '运行中', stopped: '已停止', suspended: '已暂停' }
  return texts[status] || status
}

const initCharts = () => {
  if (cpuChart.value) {
    const cpuInstance = echarts.init(cpuChart.value)
    cpuInstance.setOption({
      title: { text: 'CPU 使用率' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['10:00', '10:05', '10:10', '10:15', '10:20', '10:25'] },
      yAxis: { type: 'value', max: 100 },
      series: [{ data: [30, 45, 60, 40, 55, 50], type: 'line', smooth: true, color: '#409EFF' }]
    })
  }
  
  if (memoryChart.value) {
    const memoryInstance = echarts.init(memoryChart.value)
    memoryInstance.setOption({
      title: { text: '内存使用率' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['10:00', '10:05', '10:10', '10:15', '10:20', '10:25'] },
      yAxis: { type: 'value', max: 100 },
      series: [{ data: [50, 55, 60, 58, 62, 60], type: 'line', smooth: true, color: '#67C23A' }]
    })
  }
}

const handleStart = () => {
  ElMessage.success('已发送启动命令')
}

const handleStop = () => {
  ElMessageBox.confirm('确定要停止此虚拟机吗？', '提示', { type: 'warning' })
    .then(() => ElMessage.success('已发送停止命令'))
}

const handleRestart = () => {
  ElMessageBox.confirm('确定要重启此虚拟机吗？', '提示', { type: 'warning' })
    .then(() => ElMessage.success('已发送重启命令'))
}

const handleSnapshot = () => {
  ElMessage.success('快照创建成功')
}

const handleRestoreSnapshot = (snapshot) => {
  ElMessageBox.confirm(`确定要恢复到快照 "${snapshot.name}" 吗？`, '警告', { type: 'warning' })
    .then(() => ElMessage.success('快照恢复成功'))
}

const handleDeleteSnapshot = (snapshot) => {
  ElMessageBox.confirm(`确定要删除快照 "${snapshot.name}" 吗？`, '警告', { type: 'error' })
    .then(() => ElMessage.success('快照删除成功'))
}

onMounted(() => {
  initCharts()
  chartInterval = setInterval(() => {
    // TODO: 获取最新监控数据并更新图表
  }, 5000)
})

onUnmounted(() => {
  if (chartInterval) clearInterval(chartInterval)
  if (cpuChart.value) echarts.dispose(cpuChart.value)
  if (memoryChart.value) echarts.dispose(memoryChart.value)
})
</script>

<style scoped lang="scss">
.vm-detail-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}
</style>
