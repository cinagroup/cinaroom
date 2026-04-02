<template>
  <div class="app-layout" :class="{ 'sidebar-collapsed': isCollapse, 'is-mobile': isMobile }">
    <!-- Header -->
    <header class="app-header">
      <div class="header-left">
        <!-- 移动端菜单按钮 -->
        <button v-if="isMobile" class="mobile-menu-btn" @click="drawerVisible = true">
          <component :is="icons.Menu" :size="20" />
        </button>
        <!-- 折叠按钮 -->
        <button v-if="!isMobile" class="collapse-btn" @click="toggleCollapse">
          <component :is="isCollapse ? icons.PanelLeftOpen : icons.PanelLeftClose" :size="20" />
        </button>
        <!-- Logo -->
        <div class="header-logo" @click="$router.push('/vms')">
          <span class="logo-text">CinaSeek</span>
        </div>
      </div>

      <div class="header-center">
        <!-- 面包屑 -->
        <el-breadcrumb separator="/" v-if="breadcrumbs.length > 0">
          <el-breadcrumb-item v-for="(item, idx) in breadcrumbs" :key="idx" :to="item.path ? { path: item.path } : undefined">
            {{ item.title }}
          </el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <div class="header-right">
        <!-- 主题切换 -->
        <el-dropdown trigger="click" @command="handleThemeChange">
          <button class="icon-btn">
            <component :is="themeIcon" :size="18" />
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="light" :class="{ 'is-active': appStore.theme === 'light' }">
                <component :is="icons.Sun" :size="14" style="vertical-align: middle; margin-right: 6px;" />
                浅色模式
              </el-dropdown-item>
              <el-dropdown-item command="dark" :class="{ 'is-active': appStore.theme === 'dark' }">
                <component :is="icons.Moon" :size="14" style="vertical-align: middle; margin-right: 6px;" />
                深色模式
              </el-dropdown-item>
              <el-dropdown-item command="auto" :class="{ 'is-active': appStore.theme === 'auto' }">
                <component :is="icons.Monitor" :size="14" style="vertical-align: middle; margin-right: 6px;" />
                跟随系统
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <!-- 用户下拉 -->
        <el-dropdown trigger="click" @command="handleUserCommand">
          <div class="user-area">
            <el-avatar :size="30" :src="userStore.avatar || undefined">
              <component :is="icons.User" :size="16" />
            </el-avatar>
            <span class="username">{{ userStore.username || '用户' }}</span>
            <component :is="icons.ChevronDown" :size="14" />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <component :is="icons.UserCircle" :size="14" style="vertical-align: middle; margin-right: 6px;" />
                个人中心
              </el-dropdown-item>
              <el-dropdown-item v-if="userStore.isAdmin" command="admin" divided>
                <component :is="icons.Shield" :size="14" style="vertical-align: middle; margin-right: 6px;" />
                管理后台
              </el-dropdown-item>
              <el-dropdown-item command="logout" divided>
                <component :is="icons.LogOut" :size="14" style="vertical-align: middle; margin-right: 6px;" />
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <!-- Sidebar (Desktop) -->
    <aside v-if="!isMobile" class="app-sidebar" :class="{ collapsed: isCollapse }">
      <nav class="sidebar-nav">
        <div class="nav-group">
          <div class="nav-group-label" v-if="!isCollapse">工作区</div>
          <router-link
            v-for="item in workspaceMenus"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
          >
            <component :is="item.icon" :size="18" />
            <span v-if="!isCollapse" class="nav-text">{{ item.title }}</span>
          </router-link>
        </div>

        <div class="nav-group" v-if="openClawMenus.length">
          <div class="nav-group-label" v-if="!isCollapse">OpenClaw</div>
          <router-link
            v-for="item in openClawMenus"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
          >
            <component :is="item.icon" :size="18" />
            <span v-if="!isCollapse" class="nav-text">{{ item.title }}</span>
          </router-link>
        </div>

        <!-- 管理员区域 -->
        <div class="nav-group" v-if="userStore.isAdmin && adminMenus.length">
          <div class="nav-group-label" v-if="!isCollapse">管理员</div>
          <router-link
            v-for="item in adminMenus"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
          >
            <component :is="item.icon" :size="18" />
            <span v-if="!isCollapse" class="nav-text">{{ item.title }}</span>
          </router-link>
        </div>
      </nav>

      <!-- 底部折叠按钮 -->
      <div class="sidebar-footer">
        <button class="collapse-toggle-btn" @click="toggleCollapse">
          <component
            :is="icons.ChevronLeft"
            :size="16"
            :style="{ transform: isCollapse ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.3s' }"
          />
          <span v-if="!isCollapse">收起侧边栏</span>
        </button>
      </div>
    </aside>

    <!-- Mobile Drawer -->
    <el-drawer
      v-if="isMobile"
      v-model="drawerVisible"
      direction="ltr"
      :size="260"
      :show-close="false"
      class="mobile-drawer"
    >
      <template #header>
        <div class="drawer-header">
          <span class="logo-text">CinaSeek</span>
        </div>
      </template>
      <nav class="sidebar-nav">
        <div class="nav-group">
          <div class="nav-group-label">工作区</div>
          <router-link
            v-for="item in workspaceMenus"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
            @click="drawerVisible = false"
          >
            <component :is="item.icon" :size="18" />
            <span class="nav-text">{{ item.title }}</span>
          </router-link>
        </div>

        <div class="nav-group" v-if="openClawMenus.length">
          <div class="nav-group-label">OpenClaw</div>
          <router-link
            v-for="item in openClawMenus"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
            @click="drawerVisible = false"
          >
            <component :is="item.icon" :size="18" />
            <span class="nav-text">{{ item.title }}</span>
          </router-link>
        </div>

        <div class="nav-group" v-if="userStore.isAdmin && adminMenus.length">
          <div class="nav-group-label">管理员</div>
          <router-link
            v-for="item in adminMenus"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :class="{ active: isActive(item.path) }"
            @click="drawerVisible = false"
          >
            <component :is="item.icon" :size="18" />
            <span class="nav-text">{{ item.title }}</span>
          </router-link>
        </div>
      </nav>
    </el-drawer>

    <!-- Main Content -->
    <main class="app-main" :class="{ 'sidebar-expanded': !isMobile && !isCollapse, 'sidebar-collapsed-main': !isMobile && isCollapse }">
      <div class="content-wrapper">
        <router-view />
      </div>
      <!-- Footer -->
      <footer class="app-footer">
        <span>© 2026 CinaGroup. All rights reserved.</span>
      </footer>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, shallowRef, type Component } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import {
  Monitor,
  Terminal,
  FolderOpen,
  Settings,
  CloudUpload,
  Activity,
  Briefcase,
  Link,
  Users,
  FileText,
  User,
  UserCircle,
  Shield,
  LogOut,
  Sun,
  Moon,
  ChevronDown,
  ChevronLeft,
  PanelLeftClose,
  PanelLeftOpen,
  Menu,
  type LucideIcon
} from 'lucide-vue-next'

const icons = {
  Monitor, Terminal, FolderOpen, Settings, CloudUpload, Activity, Briefcase, Link,
  Users, FileText, User, UserCircle, Shield, LogOut, Sun, Moon,
  ChevronDown, ChevronLeft, PanelLeftClose, PanelLeftOpen, Menu
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const isCollapse = ref(false)
const drawerVisible = ref(false)
const isMobile = ref(window.innerWidth < 768)

// 监听窗口大小
window.addEventListener('resize', () => {
  isMobile.value = window.innerWidth < 768
  if (!isMobile.value) drawerVisible.value = false
})

// 面包屑
const breadcrumbs = computed(() => {
  const items: Array<{ title: string; path?: string }> = []
  const matched = route.matched
  for (const r of matched) {
    if (r.meta?.title) {
      items.push({ title: r.meta.title as string, path: r.path })
    }
  }
  return items
})

// 主题图标
const themeIcon = computed(() => {
  if (appStore.theme === 'auto') return Monitor
  return appStore.isDark ? Moon : Sun
})

// 菜单定义
interface MenuItem {
  path: string
  title: string
  icon: LucideIcon | Component
  requireAdmin?: boolean
}

const workspaceMenus = computed<MenuItem[]>(() => {
  const items: MenuItem[] = [
    { path: '/vms', title: '我的虚拟机', icon: Monitor },
    { path: '/mounts', title: '挂载管理', icon: FolderOpen },
  ]
  // 普通用户也可以看到个人中心
  if (!userStore.isAdmin) {
    items.push({ path: '/profile', title: '个人中心', icon: UserCircle })
  }
  return items
})

const openClawMenus: MenuItem[] = [
  { path: '/openclaw/deploy', title: '部署管理', icon: CloudUpload },
  { path: '/openclaw/config', title: '管控配置', icon: Settings },
  { path: '/openclaw/monitor', title: '监控面板', icon: Activity },
  { path: '/openclaw/workspace', title: '工作空间', icon: Briefcase },
]

const adminMenus = computed<MenuItem[]>(() => {
  if (!userStore.isAdmin) return []
  return [
    { path: '/admin/users', title: '用户管理', icon: Users, requireAdmin: true },
    { path: '/admin/logs', title: '系统日志', icon: FileText, requireAdmin: true },
    { path: '/admin/settings', title: '系统设置', icon: Settings, requireAdmin: true },
  ]
})

function isActive(path: string) {
  return route.path === path || route.path.startsWith(path + '/')
}

function toggleCollapse() {
  isCollapse.value = !isCollapse.value
}

function handleThemeChange(mode: string) {
  appStore.setTheme(mode as 'light' | 'dark' | 'auto')
}

function handleUserCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'admin':
      router.push('/admin/users')
      break
    case 'logout':
      userStore.logout()
      router.push('/login')
      break
  }
}

// 初始化主题
appStore.initTheme()
</script>

<style scoped lang="scss">
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

// ===== Header =====
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: $header-height;
  background: var(--bg-header);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
  transition: background-color 0.3s, border-color 0.3s;

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;

    .collapse-btn,
    .mobile-menu-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 36px;
      height: 36px;
      border: none;
      background: transparent;
      color: var(--text-secondary);
      cursor: pointer;
      border-radius: 6px;
      transition: all 0.2s;

      &:hover {
        background: var(--accent-color-light);
        color: var(--accent-color);
      }
    }

    .header-logo {
      cursor: pointer;
      display: flex;
      align-items: center;

      .logo-text {
        font-size: 18px;
        font-weight: 700;
        background: linear-gradient(135deg, #4A90D9, #667eea);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        letter-spacing: 1px;
      }
    }
  }

  .header-center {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;

    .icon-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 34px;
      height: 34px;
      border: none;
      background: var(--bg-tertiary);
      color: var(--text-secondary);
      cursor: pointer;
      border-radius: 50%;
      transition: all 0.2s;

      &:hover {
        background: var(--accent-color-light);
        color: var(--accent-color);
      }
    }

    .user-area {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 4px 8px;
      border-radius: 8px;
      transition: background-color 0.2s;

      &:hover {
        background: var(--accent-color-light);
      }

      .username {
        font-size: 14px;
        color: var(--text-primary);
        font-weight: 500;
      }
    }
  }
}

// ===== Sidebar =====
.app-sidebar {
  position: fixed;
  top: $header-height;
  left: 0;
  bottom: 0;
  width: $sidebar-width;
  background: var(--bg-sidebar);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  z-index: 99;
  overflow: hidden;

  &.collapsed {
    width: $sidebar-collapsed-width;
  }

  .sidebar-nav {
    flex: 1;
    overflow-y: auto;
    padding: 12px 0;

    &::-webkit-scrollbar {
      width: 0;
    }
  }

  .nav-group {
    margin-bottom: 8px;

    .nav-group-label {
      padding: 8px 20px 4px;
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      color: var(--text-sidebar-group);
      font-weight: 600;
      white-space: nowrap;
      overflow: hidden;
    }
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 9px 20px;
    color: var(--text-sidebar);
    text-decoration: none;
    transition: all 0.2s;
    white-space: nowrap;
    margin: 1px 8px;
    border-radius: 6px;

    &:hover {
      background: var(--bg-sidebar-hover);
      color: var(--text-sidebar-active);
    }

    &.active {
      background: var(--bg-sidebar-active);
      color: var(--accent-color);

      :deep(svg) {
        color: var(--accent-color);
      }
    }

    .nav-text {
      font-size: 14px;
      font-weight: 500;
    }
  }

  .collapsed .nav-item {
    justify-content: center;
    padding: 10px 0;
  }

  .sidebar-footer {
    padding: 12px;
    border-top: 1px solid rgba(255, 255, 255, 0.06);

    .collapse-toggle-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      width: 100%;
      padding: 6px 12px;
      border: 1px solid rgba(255, 255, 255, 0.1);
      background: transparent;
      color: var(--text-sidebar);
      cursor: pointer;
      border-radius: 6px;
      font-size: 13px;
      transition: all 0.2s;

      &:hover {
        background: var(--bg-sidebar-hover);
      }
    }
  }

  &.collapsed .sidebar-footer .collapse-toggle-btn {
    justify-content: center;
    padding: 6px;
  }
}

// ===== Main Content =====
.app-main {
  margin-top: $header-height;
  transition: margin-left 0.3s ease;
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - #{$header-height});

  &.sidebar-expanded {
    margin-left: $sidebar-width;
  }

  &.sidebar-collapsed-main {
    margin-left: $sidebar-collapsed-width;
  }

  .content-wrapper {
    flex: 1;
    padding: $content-margin;
  }

  .app-footer {
    text-align: center;
    padding: 16px;
    color: var(--text-tertiary);
    font-size: 13px;
    border-top: 1px solid var(--border-color-light);
  }
}

// ===== Mobile =====
.is-mobile {
  .app-main {
    margin-left: 0 !important;
  }
}

.mobile-drawer {
  :deep(.el-drawer__header) {
    margin-bottom: 0;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  :deep(.el-drawer__body) {
    padding: 0;
    background: var(--bg-sidebar);
  }

  .drawer-header {
    .logo-text {
      font-size: 18px;
      font-weight: 700;
      background: linear-gradient(135deg, #4A90D9, #667eea);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }

  .sidebar-nav {
    padding: 12px 0;

    .nav-group {
      margin-bottom: 8px;

      .nav-group-label {
        padding: 8px 20px 4px;
        font-size: 11px;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        color: var(--text-sidebar-group);
        font-weight: 600;
      }
    }

    .nav-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 10px 20px;
      color: var(--text-sidebar);
      text-decoration: none;
      transition: all 0.2s;
      margin: 1px 8px;
      border-radius: 6px;

      &:hover,
      &.active {
        background: var(--bg-sidebar-active);
        color: var(--accent-color);
      }

      .nav-text {
        font-size: 14px;
        font-weight: 500;
      }
    }
  }
}
</style>
