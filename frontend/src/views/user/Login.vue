<template>
  <div class="login-page" :class="{ dark: appStore.isDark }">
    <!-- Background -->
    <div class="login-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>

    <!-- Login Card -->
    <div class="login-card">
      <!-- Logo -->
      <div class="login-brand">
        <h1 class="brand-title">CinaSeek</h1>
        <p class="brand-subtitle">云端开发工作室</p>
      </div>

      <!-- Tabs -->
      <el-tabs v-model="activeTab" class="login-tabs">
        <!-- Account Login -->
        <el-tab-pane label="账号登录" name="account">
          <el-form
            ref="formRef"
            :model="loginForm"
            :rules="rules"
            label-position="top"
            @submit.prevent="handleLogin"
            class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名或邮箱"
                size="large"
                clearable
              >
                <template #prefix>
                  <component :is="UserIcon" :size="16" style="color: var(--text-tertiary)" />
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                size="large"
                show-password
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <component :is="LockIcon" :size="16" style="color: var(--text-tertiary)" />
                </template>
              </el-input>
            </el-form-item>

            <div class="form-options">
              <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
            </div>

            <el-button
              type="primary"
              size="large"
              :loading="loginLoading"
              @click="handleLogin"
              class="login-btn"
            >
              登录
            </el-button>
          </el-form>
        </el-tab-pane>

        <!-- OAuth Login -->
        <el-tab-pane label="OAuth 登录" name="oauth">
          <div class="oauth-section">
            <p class="oauth-desc">支持 GitHub、Google、Microsoft 等 9+ 平台统一认证</p>

            <button
              class="github-btn"
              :disabled="oauthLoading"
              @click="handleOAuthLogin"
            >
              <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
              <span>{{ oauthLoading ? '跳转中...' : '使用 CinaToken 账号登录' }}</span>
            </button>

            <a href="#" class="view-providers" @click.prevent="showProviders">查看支持的登录方式</a>
          </div>
        </el-tab-pane>
      </el-tabs>

      <!-- Register Link -->
      <div class="login-footer">
        <router-link to="/register">还没有账号？立即注册</router-link>
      </div>
    </div>

    <!-- 支持的登录方式对话框 -->
    <el-dialog v-model="showProvidersDialog" title="支持的登录平台" width="450px">
      <div class="provider-list">
        <div v-for="p in providers" :key="p.name" class="provider-item">
          <component :is="CircleCheckIcon" :size="18" style="color: var(--success-color)" />
          <span>{{ p.display_name }}</span>
        </div>
        <p class="provider-tip">所有平台通过 CinaToken 统一认证，安全便捷</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, shallowRef } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User as UserIcon, Lock as LockIcon, CircleCheck as CircleCheckIcon } from 'lucide-vue-next'
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
const activeTab = ref('account')

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
.login-page {
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
  width: 420px;
  position: relative;
  z-index: 1;
  background: var(--card-bg);
  border: 1px solid var(--border-color-light);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 40px 36px 32px;
  transition: background-color 0.3s, border-color 0.3s;
}

.login-brand {
  text-align: center;
  margin-bottom: 32px;

  .brand-title {
    margin: 0;
    font-size: 32px;
    font-weight: 700;
    background: linear-gradient(135deg, #4A90D9, #667eea);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 2px;
  }

  .brand-subtitle {
    margin: 8px 0 0;
    color: var(--text-tertiary);
    font-size: 14px;
  }
}

.login-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 24px;
  }

  :deep(.el-tabs__nav-wrap::after) {
    background-color: var(--border-color-light);
  }

  :deep(.el-tabs__item) {
    font-size: 15px;
    font-weight: 500;
  }
}

.login-form {
  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 18px;
  }

  .login-btn {
    width: 100%;
    height: 44px;
    font-size: 15px;
    border-radius: 8px;
    font-weight: 600;
  }
}

.oauth-section {
  text-align: center;
  padding: 20px 0;

  .oauth-desc {
    margin: 0 0 24px;
    color: var(--text-tertiary);
    font-size: 14px;
  }

  .github-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    width: 100%;
    height: 48px;
    background: #24292e;
    color: #fff;
    border: none;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
      background: #2f363d;
    }

    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }

  .view-providers {
    display: inline-block;
    margin-top: 16px;
    color: var(--accent-color);
    font-size: 13px;
    text-decoration: none;

    &:hover {
      color: var(--accent-color-hover);
    }
  }
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color-light);

  a {
    color: var(--accent-color);
    font-size: 14px;
    text-decoration: none;
    font-weight: 500;

    &:hover {
      color: var(--accent-color-hover);
    }
  }
}

.provider-list {
  .provider-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 0;
    color: var(--text-primary);
  }

  .provider-tip {
    margin-top: 15px;
    color: var(--text-tertiary);
    font-size: 13px;
    text-align: center;
  }
}
</style>
