<template>
  <div class="security-container page-container">
    <!-- Change Password -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="header-title">安全设置</span>
        </div>
      </template>

      <el-divider content-position="left">
        <el-icon><Lock /></el-icon>修改密码
      </el-divider>
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="120px"
        style="max-width: 520px"
      >
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input v-model="passwordForm.currentPassword" type="password" show-password placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="8-20 位，含字母和数字" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="changingPassword" @click="handleChangePassword">修改密码</el-button>
        </el-form-item>
      </el-form>

      <!-- Two-Factor Auth -->
      <el-divider content-position="left">
        <el-icon><Key /></el-icon>两步验证
      </el-divider>
      <div class="security-item">
        <div>
          <p class="item-title">启用两步验证</p>
          <p class="item-desc">通过 TOTP 应用生成验证码，提高账户安全性</p>
        </div>
        <el-switch v-model="twoFactorEnabled" :loading="toggling2FA" @change="handleToggle2FA" />
      </div>
      <div v-if="twoFactorSecret" class="secret-display">
        <el-alert type="warning" :closable="false" show-icon>
          请将以下密钥添加到您的 TOTP 应用中
        </el-alert>
        <el-input :model-value="twoFactorSecret" readonly style="margin-top: 10px">
          <template #append>
            <el-button @click="copySecret">复制</el-button>
          </template>
        </el-input>
      </div>

      <!-- Login Logs -->
      <el-divider content-position="left">
        <el-icon><Document /></el-icon>登录日志
      </el-divider>
      <el-table :data="loginLogs" style="width: 100%" v-loading="loadingLogs" stripe>
        <el-table-column prop="time" label="登录时间" width="180" />
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="location" label="登录地点" min-width="120" />
        <el-table-column prop="device" label="设备" min-width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === '成功' ? 'success' : 'danger'" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>

      <!-- Active Sessions -->
      <el-divider content-position="left">
        <el-icon><Connection /></el-icon>活跃会话
      </el-divider>
      <el-table :data="sessions" style="width: 100%" v-loading="loadingSessions" stripe>
        <el-table-column prop="device" label="设备" min-width="180" />
        <el-table-column prop="location" label="地点" min-width="120" />
        <el-table-column prop="loginTime" label="登录时间" width="180" />
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.current" type="success" size="small">当前</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              v-if="!row.current"
              type="danger"
              link
              size="small"
              @click="handleKillSession(row)"
            >
              下线
            </el-button>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Lock, Key, Document, Connection } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import {
  changePassword,
  toggleTwoFactor,
  getLoginLogs,
  getActiveSessions,
  killSession
} from '@/api/auth'
import type { LoginLog, ActiveSession } from '@/types/remote'

const passwordFormRef = ref<FormInstance>()
const changingPassword = ref(false)
const toggling2FA = ref(false)
const twoFactorEnabled = ref(false)
const twoFactorSecret = ref('')
const loadingLogs = ref(false)
const loadingSessions = ref(false)

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (_rule: unknown, value: string, callback: (err?: Error) => void) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的新密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = reactive<FormRules>({
  currentPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, max: 20, message: '密码长度 8-20 位', trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: '密码需包含字母和数字', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
})

const loginLogs = ref<LoginLog[]>([])
const sessions = ref<ActiveSession[]>([])

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  const valid = await passwordFormRef.value.validate().catch(() => false)
  if (!valid) return

  changingPassword.value = true
  try {
    await changePassword({
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword
    })
    ElMessage.success('密码修改成功，请重新登录')
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error: any) {
    ElMessage.error(error?.message || '密码修改失败')
  } finally {
    changingPassword.value = false
  }
}

const handleToggle2FA = async (enabled: boolean) => {
  toggling2FA.value = true
  try {
    const res = await toggleTwoFactor(enabled) as any
    const data = res.data || res
    if (enabled && data.secret) {
      twoFactorSecret.value = data.secret
      ElMessage.success('两步验证已启用，请保存密钥')
    } else {
      twoFactorSecret.value = ''
      ElMessage.success(enabled ? '两步验证已启用' : '两步验证已禁用')
    }
  } catch (error: any) {
    twoFactorEnabled.value = !enabled
    ElMessage.error(error?.message || '操作失败')
  } finally {
    toggling2FA.value = false
  }
}

const copySecret = () => {
  navigator.clipboard.writeText(twoFactorSecret.value)
  ElMessage.success('密钥已复制到剪贴板')
}

const handleKillSession = async (session: ActiveSession) => {
  try {
    await ElMessageBox.confirm('确定要下线该会话吗？此操作不可撤回。', '提示', { type: 'warning' })
    // Extract session id from the session data (assuming device as identifier)
    await killSession(session.device)
    ElMessage.success('已下线该会话')
    await fetchSessions()
  } catch {
    // Cancelled or API error
  }
}

const fetchLoginLogs = async () => {
  loadingLogs.value = true
  try {
    const res = await getLoginLogs() as any
    const data = res.data || res
    loginLogs.value = data.logs || []
  } catch {
    loginLogs.value = []
  } finally {
    loadingLogs.value = false
  }
}

const fetchSessions = async () => {
  loadingSessions.value = true
  try {
    const res = await getActiveSessions() as any
    const data = res.data || res
    sessions.value = data.sessions || []
  } catch {
    sessions.value = []
  } finally {
    loadingSessions.value = false
  }
}

onMounted(() => {
  fetchLoginLogs()
  fetchSessions()
})
</script>

<style scoped lang="scss">
.security-container {
  max-width: 1000px;
}

.card-header {
  .header-title {
    font-size: 18px;
    font-weight: 600;
  }
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;

  .item-title {
    font-weight: 600;
    font-size: 15px;
    margin-bottom: 4px;
  }
  .item-desc {
    font-size: 13px;
    color: $info-color;
  }
}

.secret-display {
  margin-top: 12px;
  max-width: 520px;
}

.text-muted {
  color: #c0c4cc;
}
</style>
