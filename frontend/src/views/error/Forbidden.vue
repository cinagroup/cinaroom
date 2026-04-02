<template>
  <div class="forbidden-page">
    <div class="forbidden-card">
      <div class="error-code">403</div>
      <h2 class="error-title">无权限访问</h2>
      <p class="error-desc">抱歉，您没有权限访问此页面。请联系管理员获取权限。</p>
      <el-button type="primary" size="large" @click="goHome" class="back-btn">
        <el-icon style="margin-right: 6px"><HomeFilled /></el-icon>
        返回首页
      </el-button>
    </div>

    <!-- Background decorations -->
    <div class="bg-decor">
      <div class="circle c1"></div>
      <div class="circle c2"></div>
      <div class="circle c3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled } from '@element-plus/icons-vue'

const router = useRouter()

function goHome() {
  router.push('/vms')
}
</script>

<style scoped lang="scss">
.forbidden-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--bg-secondary);
}

.bg-decor {
  position: absolute;
  inset: 0;
  pointer-events: none;

  .circle {
    position: absolute;
    border-radius: 50%;
    opacity: 0.06;
    background: var(--accent-color);
  }

  .c1 {
    width: 500px;
    height: 500px;
    top: -150px;
    right: -100px;
    animation: drift 18s ease-in-out infinite;
  }

  .c2 {
    width: 350px;
    height: 350px;
    bottom: -100px;
    left: -80px;
    animation: drift 14s ease-in-out infinite reverse;
  }

  .c3 {
    width: 200px;
    height: 200px;
    top: 50%;
    left: 60%;
    animation: drift 10s ease-in-out infinite;
  }
}

@keyframes drift {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, -20px); }
}

.forbidden-card {
  position: relative;
  z-index: 1;
  text-align: center;
  background: var(--card-bg);
  border: 1px solid var(--border-color-light);
  border-radius: 16px;
  box-shadow: var(--shadow-lg);
  padding: 60px 48px;
  max-width: 440px;
  width: 90%;
  transition: background-color 0.3s, border-color 0.3s;

  .error-code {
    font-size: 96px;
    font-weight: 800;
    line-height: 1;
    background: linear-gradient(135deg, #4A90D9, #667eea);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-bottom: 16px;
  }

  .error-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 12px;
  }

  .error-desc {
    font-size: 15px;
    color: var(--text-tertiary);
    margin: 0 0 32px;
    line-height: 1.6;
  }

  .back-btn {
    min-width: 160px;
    height: 44px;
    font-size: 15px;
    border-radius: 8px;
    font-weight: 600;
  }
}
</style>
