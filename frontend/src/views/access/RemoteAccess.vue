<template>
  <div class="remote-access-container">
    <Header>
      <template #left>
        <span class="page-title">远程访问管理</span>
      </template>
    </Header>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>访问状态</span>
          </template>
          
          <div class="status-section">
            <div class="status-item">
              <span class="label">远程访问</span>
              <el-switch v-model="accessEnabled" size="large" @change="handleToggleAccess" />
            </div>
            
            <div class="status-item" v-if="accessEnabled">
              <span class="label">访问地址</span>
              <el-input value="https://cinaseek.run" readonly>
                <template #append>
                  <el-button @click="handleCopy">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </template>
              </el-input>
            </div>
            
            <div class="status-item" v-if="accessEnabled">
              <span class="label">移动端访问</span>
              <div class="qrcode">
                <el-qrcode :value="accessUrl" :size="128" />
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>IP 白名单</span>
          </template>
          
          <div class="whitelist-toolbar">
            <el-button type="primary" size="small" @click="showAddIpDialog = true">
              <el-icon><Plus /></el-icon>
              添加 IP
            </el-button>
            <el-button size="small" @click="handleBatchImport">
              <el-icon><Upload /></el-icon>
              批量导入
            </el-button>
          </div>
          
          <el-table :data="whitelist" style="width: 100%; margin-top: 15px">
            <el-table-column prop="ip" label="IP 地址/CIDR" min-width="150" />
            <el-table-column prop="note" label="备注" min-width="150" />
            <el-table-column prop="addedAt" label="添加时间" width="160" />
            <el-table-column label="操作" width="80">
              <template #default="{ row }">
                <el-button type="danger" link @click="handleRemoveIp(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card style="margin-top: 20px">
      <template #header>
        <span>访问日志</span>
      </template>
      
      <div class="log-toolbar">
        <el-date-picker
          v-model="logDateRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          style="width: 350px"
        />
        <el-button type="primary" @click="handleExportLogs" style="margin-left: 10px">导出日志</el-button>
      </div>
      
      <el-table :data="accessLogs" style="width: 100%; margin-top: 15px">
        <el-table-column prop="time" label="访问时间" width="160" />
        <el-table-column prop="ip" label="访问 IP" width="140" />
        <el-table-column prop="path" label="访问路径" min-width="200" />
        <el-table-column prop="userAgent" label="用户代理" min-width="250" />
        <el-table-column prop="status" label="状态码" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 200 ? 'success' : 'danger'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="showAddIpDialog" title="添加 IP 白名单" width="500px">
      <el-form :model="ipForm" label-width="100px">
        <el-form-item label="IP 地址/CIDR">
          <el-input v-model="ipForm.ip" placeholder="192.168.1.0/24" />
        </el-form-item>
        
        <el-form-item label="备注">
          <el-input v-model="ipForm.note" placeholder="可选" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showAddIpDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddIp">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Header from '@/components/Header.vue'

const accessEnabled = ref(true)
const accessUrl = ref('https://cinaseek.run')
const showAddIpDialog = ref(false)
const logDateRange = ref([])

const whitelist = ref([
  { ip: '192.168.1.0/24', note: '办公室网络', addedAt: '2024-04-01' },
  { ip: '10.0.0.0/8', note: '内网', addedAt: '2024-03-28' },
  { ip: '203.0.113.50', note: '家庭 IP', addedAt: '2024-03-25' }
])

const accessLogs = ref([
  { time: '2024-04-01 10:30:15', ip: '192.168.1.100', path: '/vms', userAgent: 'Mozilla/5.0...', status: 200 },
  { time: '2024-04-01 10:28:30', ip: '192.168.1.101', path: '/shell/1', userAgent: 'Mozilla/5.0...', status: 200 },
  { time: '2024-04-01 10:25:00', ip: '203.0.113.50', path: '/login', userAgent: 'Mozilla/5.0...', status: 200 },
  { time: '2024-04-01 10:20:10', ip: '198.51.100.23', path: '/admin', userAgent: 'curl/7.68.0', status: 403 }
])

const ipForm = reactive({
  ip: '',
  note: ''
})

const handleToggleAccess = () => {
  const status = accessEnabled.value ? '启用' : '禁用'
  ElMessage.success(`远程访问已${status}`)
}

const handleCopy = () => {
  navigator.clipboard.writeText(accessUrl.value)
  ElMessage.success('地址已复制')
}

const handleAddIp = () => {
  whitelist.value.push({
    ip: ipForm.ip,
    note: ipForm.note,
    addedAt: new Date().toLocaleDateString()
  })
  showAddIpDialog.value = false
  ElMessage.success('添加成功')
  ipForm.ip = ''
  ipForm.note = ''
}

const handleRemoveIp = (item) => {
  ElMessageBox.confirm(`确定要删除 IP "${item.ip}" 吗？`, '提示', { type: 'warning' })
    .then(() => {
      const index = whitelist.value.findIndex(w => w.ip === item.ip)
      if (index > -1) whitelist.value.splice(index, 1)
      ElMessage.success('删除成功')
    })
}

const handleBatchImport = () => {
  ElMessage.info('批量导入功能开发中')
}

const handleExportLogs = () => {
  ElMessage.success('日志导出成功')
}
</script>

<style scoped lang="scss">
.remote-access-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}

.status-section {
  .status-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 0;
    border-bottom: 1px solid #ebeef5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .label {
      font-weight: bold;
    }
    
    .qrcode {
      display: flex;
      justify-content: center;
    }
  }
}

.whitelist-toolbar {
  display: flex;
  gap: 10px;
}

.log-toolbar {
  display: flex;
  align-items: center;
}
</style>
