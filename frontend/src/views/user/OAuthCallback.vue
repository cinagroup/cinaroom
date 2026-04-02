<template>
  <div class="callback-container">
    <div class="callback-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
    </div>
    <el-card class="callback-card">
      <div class="callback-content">
        <div v-if="loading" class="loading-state">
          <el-icon class="loading-icon" :size="56"><Loading /></el-icon>
          <p class="loading-text">{{ statusText }}</p>
          <el-progress :percentage="progress" :show-text="false" :stroke-width="4" style="width: 200px; margin-top: 20px" />
        </div>
        <div v-else-if="success" class="success-state">
          <el-icon class="success-icon" :size="56" color="#67C23A"><CircleCheckFilled /></el-icon>
          <p class="success-text">登录成功</p>
          <p class="welcome-text">欢迎回来，{{ username }}！</p>
          <p class="redirect-text">{{ redirectCountdown }} 秒后自动跳转...</p>
        </div>
        <div v-else class="error-state">
          <el-icon class="error-icon" :size="56" color="#F56C6C"><CircleCloseFilled /></el-icon>
          <p class="error-text">登录失败</p>
          <p class="error-detail">{{ errorMessage }}</p>
          <el-button type="primary" @click="$router.push('/login')" style="margin-top: 20px">
            返回登录
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Loading, CircleCheckFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import { oauthCallback } from '@/api/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(true)
const success = ref(false)
const statusText = ref('正在处理登录...')
const progress = ref(0)
const errorMessage = ref('')
const username = ref('')
const redirectCountdown = ref(3)

let countdownTimer: ReturnType<typeof setInterval> | null = null
let progressTimer: ReturnType<typeof setInterval> | null = null

const startProgress = () => {
  progressTimer = setInterval(() => {
    if (progress.value < 85) {
      progress.value += Math.random() * 15
    }
  }, 300)
}

const finishProgress = () => {
  if (progressTimer) clearInterval(progressTimer)
  progress.value = 100
}

const startRedirectCountdown = () => {
  countdownTimer = setInterval(() => {
    redirectCountdown.value--
    if (redirectCountdown.value <= 0) {
      if (countdownTimer) clearInterval(countdownTimer)
      const redirect = (route.query.redirect as string) || '/vms'
      router.push(redirect)
    }
  }, 1000)
}

onMounted(async () => {
  startProgress()

  try {
    const { code, state, error, error_description } = route.query as Record<string, string>

    // Check for OAuth error response
    if (error) {
      throw new Error(error_description || error || '授权被拒绝')
    }

    if (!code) {
      throw new Error('未收到授权码，请重试')
    }

    statusText.value = '正在验证授权...'

    // Call the OAuth callback API
    const res = await oauthCallback(code, state) as any
    const data = res.data || res

    // Store token and user info
    userStore.setToken(data.token)
    userStore.setUserInfo(data.user)

    username.value = data.user.nickname || data.user.username
    success.value = true
    finishProgress()

    ElMessage.success(`欢迎回来，${username.value}！`)
    startRedirectCountdown()
  } catch (error: any) {
    loading.value = false
    finishProgress()
    errorMessage.value = error?.response?.data?.msg || error?.message || '登录验证失败，请重试'
    ElMessage.error(errorMessage.value)

    // Auto redirect to login after 5 seconds
    setTimeout(() => {
      router.push('/login')
    }, 5000)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
  if (progressTimer) clearInterval(progressTimer)
})
</script>

<style scoped lang="scss">
.callback-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.callback-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;

  .bg-shape {
    position: absolute;
    border-radius: 50%;
    opacity: 0.1;
    background: #fff;
  }
  .shape-1 { width: 300px; height: 300px; top: -80px; right: -60px; }
  .shape-2 { width: 200px; height: 200px; bottom: -40px; left: -40px; }
}

.callback-card {
  width: 420px;
  position: relative;
  z-index: 1;
  border-radius: 12px !important;
}

.callback-content {
  text-align: center;
  padding: 40px 20px;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.loading-icon {
  animation: spin 1.2s linear infinite;
  color: var(--accent-color);
}

.loading-text {
  margin-top: 20px;
  color: #606266;
  font-size: 16px;
}

.success-icon, .error-icon {
  margin-bottom: 16px;
}

.success-text, .error-text {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
}

.success-text { color: #67C23A; }
.error-text { color: #F56C6C; }

.welcome-text {
  color: #606266;
  font-size: 14px;
  margin-bottom: 8px;
}

.redirect-text {
  color: $info-color;
  font-size: 13px;
}

.error-detail {
  color: $info-color;
  font-size: 14px;
  max-width: 300px;
  word-break: break-all;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
