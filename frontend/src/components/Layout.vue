<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar">
      <div class="logo">
        <span v-if="!isCollapse">Multipass</span>
        <span v-else>MP</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/vms">
          <el-icon><Monitor /></el-icon>
          <span>虚拟机管理</span>
        </el-menu-item>
        <el-menu-item index="/shell">
          <el-icon><Terminal /></el-icon>
          <span>WebShell</span>
        </el-menu-item>
        <el-menu-item index="/mounts">
          <el-icon><Folder /></el-icon>
          <span>目录挂载</span>
        </el-menu-item>
        <el-sub-menu index="openclaw">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>OpenClaw 管理</span>
          </template>
          <el-menu-item index="/openclaw/deploy">部署管理</el-menu-item>
          <el-menu-item index="/openclaw/config">管控配置</el-menu-item>
          <el-menu-item index="/openclaw/monitor">监控面板</el-menu-item>
          <el-menu-item index="/openclaw/workspace">工作空间</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/remote-access">
          <el-icon><Link /></el-icon>
          <span>远程访问</span>
        </el-menu-item>
        <el-sub-menu index="user">
          <template #title>
            <el-icon><User /></el-icon>
            <span>用户中心</span>
          </template>
          <el-menu-item index="/profile">个人信息</el-menu-item>
          <el-menu-item index="/security">安全设置</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="toggleCollapse">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" :icon="User" />
              <span class="username">{{ userStore.username || '用户' }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/profile')">
                  <el-icon><User /></el-icon>个人信息
                </el-dropdown-item>
                <el-dropdown-item @click="$router.push('/security')">
                  <el-icon><Lock /></el-icon>安全设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Monitor, Terminal, Folder, Setting, Link, User, Lock,
  Fold, Expand, SwitchButton
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const activeMenu = computed(() => route.path)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped lang="scss">
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  
  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 18px;
    font-weight: bold;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  :deep(.el-menu) {
    border-right: none;
  }
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  
  .header-left {
    .collapse-btn {
      font-size: 20px;
      cursor: pointer;
      transition: color 0.3s;
      
      &:hover {
        color: $primary-color;
      }
    }
  }
  
  .header-right {
    .user-info {
      display: flex;
      align-items: center;
      cursor: pointer;
      
      .username {
        margin-left: 8px;
      }
    }
  }
}

.main-content {
  background-color: #f0f2f5;
  padding: $content-margin;
}
</style>
