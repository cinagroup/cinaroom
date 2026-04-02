<template>
  <div class="login-container" :class="{ 'dark-theme': appStore.isDark }">
    <div class="login-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>

    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <div class="brand">
            <h1>CinaSeek</h1>
            <span class="brand-sub">云端开发工作室</span>
          </div>
        </div>
      </template>
      
      <div class="login-content">
        <!-- OAuth 登录 -->
        <div class="oauth-section">
          <h3>使用 CinaToken 账号登录</h3>
          <p class="oauth-desc">支持 GitHub、Google、Microsoft 等 9+ 平台统一认证</p>
          
          <el-button
            type="primary"
            size="large"
            :loading="oauthLoading"
            @click="handleOAuthLogin"
            class="oauth-btn"
          >
            <el-icon><Connection /></el-icon>
            使用 CinaToken 登录
          </el-button>
        </div>
        
        <el-divider>
          <span class="divider-text">或使用密码登录</span>
        </el-divider>
        
        <!-- 传统登录 -->
        <el-form
          ref="formRef"
          :model="loginForm"
          :rules="rules"
          label-position="top"
          @submit.prevent="handleLogin"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名或邮箱"
              :prefix-icon="User"
              size="large"
              clearable
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              size="large"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <div class="form-options">
            <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
          </div>
          
          <el-form-item>
            <el-button
              type="default"
              size="large"
              :loading="loginLoading"
              @click="handleLogin"
              style="width: 100%"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="links">
          <router-link to="/register">没有账号？立即注册</router-link>
          <span class="divider">|</span>
          <a href="#" @click.prevent="showProviders">查看支持的登录方式</a>
        </div>
      </div>
    </el-card>
    
    <!-- 支持的登录方式对话框 -->
    <el-dialog v-model="showProvidersDialog" title="支持的登录平台" width="450px">
      <div class="provider-list">
        <div v-for="p in providers" :key="p.name" class="provider-item">
          <el-icon size="20"><CircleCheck /></el-icon>
          <span>{{ p.display_name }}</span>
        </div>
        <p class="provider-tip">所有平台通过 CinaToken 统一认证，安全便捷</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock, Connection, CircleCheck } from '@element-plus/icons-vue'
import * as authApi from '@/api/auth'
import type { OAuthProvider } from '@/types/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()
const formRef = ref<FormInstance>()
const oauthLoading = ref(false)
const loginLoading = ref(false)
const showProvidersDialog = ref(false)
const providers = ref<OAuthProvider[]>([])

const loginForm = reactive({
  username: '',
  password: '',
  remember: false
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' }
  ]
}

/** OAuth 登录 */
const handleOAuthLogin = async () => {
  oauthLoading.value = true
  try {
    const res = await authApi.getOAuthAuthorizeUrl()
    const data = (res as any).data || res
    if (data.authorize_url) {
      window.location.href = data.authorize_url
    }
  } catch (error: any) {
    ElMessage.error('获取授权地址失败：' + (error.message || '未知错误'))
  } finally {
    oauthLoading.value = false
  }
}

/** 传统登录 */
const handleLogin = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loginLoading.value = true
  try {
    await userStore.login({
      username: loginForm.username,
      password: loginForm.password,
      remember: loginForm.remember
    })
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/vms'
    router.push(redirect)
  } catch (error: any) {
    // 错误已在拦截器中处理
  } finally {
    loginLoading.value = false
  }
}

/** 显示支持的登录方式 */
const showProviders = async () => {
  try {
    const res = await authApi.getOAuthProviders()
    const data = (res as any).data || res
    providers.value = (data.providers || []).filter((p: OAuthProvider) => p.enabled)
    showProvidersDialog.value = true
  } catch {
    ElMessage.error('获取登录方式失败')
  }
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #0c1929 0%, #1a365d 50%, #234e7e 100%);
}

.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;

  .bg-shape {
    position: absolute;
    border-radius: 50%;
    opacity: 0.08;
    background: #4A90D9;
  }

  .shape-1 {
    width: 600px;
    height: 600px;
    top: -200px;
    right: -100px;
    animation: float 20s ease-in-out infinite;
  }

  .shape-2 {
    width: 400px;
    height: 400px;
    bottom: -150px;
    left: -100px;
    animation: float 15s ease-in-out infinite reverse;
  }

  .shape-3 {
    width: 200px;
    height: 200px;
    top: 40%;
    left: 20%;
    animation: float 10s ease-in-out infinite;
  }
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -30px); }
}

.login-card {
  width: 440px;
  position: relative;
  z-index: 1;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  
  .card-header {
    text-align: center;
    
    .brand {
      h1 {
        margin: 0;
        font-size: 28px;
        background: linear-gradient(135deg, #4A90D9, #667eea);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        letter-spacing: 2px;
      }
      
      .brand-sub {
        color: #909399;
        font-size: 14px;
      }
    }
  }
}

.login-content {
  .oauth-section {
    text-align: center;
    padding: 15px 0;
    
    h3 {
      margin: 0 0 8px;
      font-size: 15px;
      color: #303133;
    }
    
    .oauth-desc {
      margin: 0 0 16px;
      color: #909399;
      font-size: 13px;
    }
    
    .oauth-btn {
      width: 100%;
      height: 44px;
      font-size: 15px;
      border-radius: 8px;
    }
  }
  
  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 18px;
  }
}

.links {
  text-align: center;
  margin-top: 12px;
  
  a {
    color: #4A90D9;
    text-decoration: none;
    font-size: 13px;
    transition: color 0.2s;
    
    &:hover {
      color: #667eea;
    }
  }
  
  .divider {
    margin: 0 10px;
    color: #dcdfe6;
  }
}

.provider-list {
  .provider-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 0;
    color: #303133;
  }
  
  .provider-tip {
    margin-top: 15px;
    color: #909399;
    font-size: 13px;
    text-align: center;
  }
}
</style>
