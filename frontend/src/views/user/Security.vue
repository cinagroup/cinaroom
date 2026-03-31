<template>
  <div class="security-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>安全设置</span>
        </div>
      </template>
      
      <el-divider content-position="left">修改密码</el-divider>
      <el-form :model="passwordForm" label-width="120px" style="max-width: 500px">
        <el-form-item label="当前密码">
          <el-input v-model="passwordForm.currentPassword" type="password" show-password />
        </el-form-item>
        
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.newPassword" type="password" show-password />
        </el-form-item>
        
        <el-form-item label="确认新密码">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleChangePassword">修改密码</el-button>
        </el-form-item>
      </el-form>
      
      <el-divider content-position="left">两步验证</el-divider>
      <div class="security-item">
        <span>启用两步验证，提高账户安全性</span>
        <el-switch v-model="twoFactorEnabled" />
      </div>
      
      <el-divider content-position="left">登录日志</el-divider>
      <el-table :data="loginLogs" style="width: 100%">
        <el-table-column prop="time" label="登录时间" width="180" />
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="location" label="登录地点" />
        <el-table-column prop="device" label="设备" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === '成功' ? 'success' : 'danger'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <el-divider content-position="left">活跃会话</el-divider>
      <el-table :data="sessions" style="width: 100%">
        <el-table-column prop="device" label="设备" />
        <el-table-column prop="location" label="地点" />
        <el-table-column prop="loginTime" label="登录时间" />
        <el-table-column prop="ip" label="IP 地址" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleLogoutSession(row)">下线</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const twoFactorEnabled = ref(false)

const loginLogs = ref([
  { time: '2024-04-01 10:30:00', ip: '192.168.1.100', location: '新加坡', device: 'Chrome / Windows', status: '成功' },
  { time: '2024-03-31 15:20:00', ip: '192.168.1.100', location: '新加坡', device: 'Chrome / Windows', status: '成功' }
])

const sessions = ref([
  { device: 'Chrome / Windows', location: '新加坡', loginTime: '2024-04-01 10:30:00', ip: '192.168.1.100', current: true },
  { device: 'Safari / iPhone', location: '新加坡', loginTime: '2024-04-01 09:00:00', ip: '192.168.1.101', current: false }
])

const handleChangePassword = () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的新密码不一致')
    return
  }
  // TODO: 调用修改密码 API
  ElMessage.success('密码修改成功')
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

const handleLogoutSession = (session) => {
  if (session.current) {
    ElMessage.warning('不能下线当前会话')
    return
  }
  
  ElMessageBox.confirm('确定要下线该会话吗？', '提示', {
    type: 'warning'
  }).then(() => {
    // TODO: 调用下线 API
    ElMessage.success('已下线该会话')
  })
}
</script>

<style scoped lang="scss">
.security-container {
  max-width: 1000px;
}

.card-header {
  font-weight: bold;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}
</style>
