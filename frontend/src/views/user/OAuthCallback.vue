<template>
  <div class="callback-container">
    <el-card class="callback-card">
      <div class="callback-content">
        <el-icon class="loading-icon" :size="50"><Loading /></el-icon>
        <p class="loading-text">{{ statusText }}</p>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const statusText = ref('正在完成登录...')

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

onMounted(async () => {
  try {
    const { code, state, error } = route.query
    
    // 检查是否有错误
    if (error) {
      throw new Error(`授权失败：${error}`)
    }
    
    if (!code) {
      throw new Error('未收到授权码')
    }
    
    statusText.value = '正在验证授权...'
    
    // 调用后端回调接口
    const response = await axios.get(`${API_BASE}/oauth/callback`, {
      params: { code, state }
    })
    
    const { token, user, oauth } = response.data
    
    // 保存 Token 和用户信息
    userStore.setToken(token)
    userStore.setUserInfo(user)
    
    statusText.value = '登录成功，正在跳转...'
    ElMessage.success(`欢迎 ${user.nickname || user.username}！`)
    
    // 延迟跳转，让用户看到成功提示
    setTimeout(() => {
      router.push('/vms')
    }, 1000)
    
  } catch (error) {
    console.error('OAuth 回调失败:', error)
    statusText.value = '登录失败'
    
    const errorMsg = error.response?.data?.msg || error.message || '未知错误'
    ElMessage.error(`登录失败：${errorMsg}`)
    
    // 3 秒后跳转回登录页
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  }
})
</script>

<style scoped lang="scss">
.callback-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.callback-card {
  width: 400px;
  
  .callback-content {
    text-align: center;
    padding: 40px 20px;
    
    .loading-icon {
      animation: rotate 1.5s linear infinite;
      color: $primary-color;
    }
    
    .loading-text {
      margin-top: 20px;
      color: #606266;
      font-size: 16px;
    }
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
