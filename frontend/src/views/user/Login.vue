<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>CinaRoom</h2>
          <p>你的云端开发工作室</p>
        </div>
      </template>
      
      <div class="login-content">
        <div class="oauth-section">
          <h3>使用 CinaToken 账号登录</h3>
          <p class="oauth-desc">支持 GitHub、Google、Microsoft 等 9+ 平台</p>
          
          <el-button
            type="primary"
            :loading="loading"
            @click="handleOAuthLogin"
            class="oauth-btn"
          >
            <el-icon><User /></el-icon>
            使用 CinaToken 登录
          </el-button>
        </div>
        
        <el-divider>
          <span class="divider-text">或</span>
        </el-divider>
        
        <div class="traditional-login">
          <p class="tip">传统账号密码登录（兼容模式）</p>
          <el-form
            ref="formRef"
            :model="loginForm"
            :rules="rules"
            label-width="80px"
          >
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名或邮箱"
                prefix-icon="User"
              />
            </el-form-item>
            
            <el-form-item label="密码" prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            
            <el-form-item>
              <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
            </el-form-item>
            
            <el-form-item>
              <el-button type="default" :loading="loading" @click="handleLogin" style="width: 100%">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <div class="links">
          <router-link to="/register">没有账号？立即注册</router-link>
          <span class="divider">|</span>
          <a href="#" @click.prevent="showProviders">查看支持的登录方式</a>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

const loginForm = reactive({
  username: '',
  password: '',
  remember: false
})

const rules = {
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' }
  ]
}

// OAuth 登录
const handleOAuthLogin = async () => {
  loading.value = true
  try {
    // 获取授权 URL
    const response = await axios.get(`${API_BASE}/oauth/authorize`)
    const authorizeUrl = response.data.authorize_url
    
    // 跳转到 CinaToken 授权页
    window.location.href = authorizeUrl
  } catch (error) {
    ElMessage.error('获取授权地址失败：' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

// 传统登录
const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await axios.post(`${API_BASE}/auth/login`, loginForm)
        
        const { token, user } = response.data
        userStore.setToken(token)
        userStore.setUserInfo(user)
        
        ElMessage.success('登录成功')
        router.push('/vms')
      } catch (error) {
        ElMessage.error('登录失败：' + (error.response?.data?.msg || error.message))
      } finally {
        loading.value = false
      }
    }
  })
}

// 显示支持的登录方式
const showProviders = async () => {
  try {
    const response = await axios.get(`${API_BASE}/oauth/providers`)
    const providers = response.data.providers
    
    const providerList = providers
      .filter(p => p.enabled)
      .map(p => `• ${p.display_name}`)
      .join('\n')
    
    ElMessageBox.alert(
      `支持以下 ${providers.filter(p => p.enabled).length} 种登录方式：\n\n${providerList}\n\n所有平台通过 CinaToken 统一认证`,
      '支持的登录平台',
      { confirmButtonText: '知道了' }
    )
  } catch (error) {
    ElMessage.error('获取登录方式失败：' + error.message)
  }
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 450px;
  
  .card-header {
    text-align: center;
    
    h2 {
      margin: 0 0 10px;
      color: $primary-color;
      font-size: 24px;
    }
    
    p {
      margin: 0;
      color: $info-color;
      font-size: 14px;
    }
  }
}

.login-content {
  .oauth-section {
    text-align: center;
    padding: 20px 0;
    
    h3 {
      margin: 0 0 10px;
      font-size: 16px;
      color: #303133;
    }
    
    .oauth-desc {
      margin: 0 0 20px;
      color: #909399;
      font-size: 13px;
    }
    
    .oauth-btn {
      width: 100%;
      height: 45px;
      font-size: 16px;
    }
  }
  
  .traditional-login {
    padding: 10px 0;
    
    .tip {
      text-align: center;
      color: #909399;
      font-size: 13px;
      margin-bottom: 15px;
    }
  }
}

.links {
  text-align: center;
  margin-top: 15px;
  
  a {
    color: $primary-color;
    text-decoration: none;
    font-size: 14px;
    
    &:hover {
      text-decoration: underline;
    }
  }
  
  .divider {
    margin: 0 10px;
    color: #dcdfe6;
  }
}
</style>
