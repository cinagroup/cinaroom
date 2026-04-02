import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/components/Layout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/Login.vue'),
    meta: { title: '登录', requireAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/Register.vue'),
    meta: { title: '注册', requireAuth: false }
  },
  {
    path: '/oauth/callback',
    name: 'OAuthCallback',
    component: () => import('@/views/user/OAuthCallback.vue'),
    meta: { title: '登录验证', requireAuth: false }
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/Forbidden.vue'),
    meta: { title: '无权限', requireAuth: false }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/vms',
    meta: { requireAuth: true },
    children: [
      // 用户中心
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人信息', icon: 'User', requireAuth: true }
      },
      {
        path: 'security',
        name: 'Security',
        component: () => import('@/views/user/Security.vue'),
        meta: { title: '安全设置', icon: 'Lock', requireAuth: true }
      },
      // 虚拟机管理
      {
        path: 'vms',
        name: 'VMs',
        component: () => import('@/views/vms/VMList.vue'),
        meta: { title: '虚拟机列表', icon: 'Monitor', requireAuth: true }
      },
      {
        path: 'vms/create',
        name: 'CreateVM',
        component: () => import('@/views/vms/CreateVM.vue'),
        meta: { title: '创建虚拟机', icon: 'Plus', requireAuth: true }
      },
      {
        path: 'vms/:id',
        name: 'VMDetail',
        component: () => import('@/views/vms/VMDetail.vue'),
        meta: { title: '虚拟机详情', icon: 'Document', requireAuth: true }
      },
      // Web 远程管理
      {
        path: 'shell/:id',
        name: 'WebShell',
        component: () => import('@/views/shell/WebShell.vue'),
        meta: { title: 'WebShell', icon: 'Terminal', requireAuth: true }
      },
      {
        path: 'remote/:id',
        name: 'RemoteControl',
        component: () => import('@/views/remote/RemoteControl.vue'),
        meta: { title: '远程管控', icon: 'Setting', requireAuth: true }
      },
      {
        path: 'logs/:id',
        name: 'LogViewer',
        component: () => import('@/views/logs/LogViewer.vue'),
        meta: { title: '日志查看', icon: 'Document', requireAuth: true }
      },
      // 目录挂载
      {
        path: 'mounts',
        name: 'Mounts',
        component: () => import('@/views/mounts/MountManager.vue'),
        meta: { title: '挂载管理', icon: 'Folder', requireAuth: true }
      },
      // OpenClaw 管理
      {
        path: 'openclaw/deploy',
        name: 'OpenClawDeploy',
        component: () => import('@/views/openclaw/Deploy.vue'),
        meta: { title: '部署管理', icon: 'Download', requireAuth: true }
      },
      {
        path: 'openclaw/config',
        name: 'OpenClawConfig',
        component: () => import('@/views/openclaw/Config.vue'),
        meta: { title: '管控配置', icon: 'Setting', requireAuth: true }
      },
      {
        path: 'openclaw/monitor',
        name: 'OpenClawMonitor',
        component: () => import('@/views/openclaw/Monitor.vue'),
        meta: { title: '监控面板', icon: 'DataAnalysis', requireAuth: true }
      },
      {
        path: 'openclaw/workspace',
        name: 'OpenClawWorkspace',
        component: () => import('@/views/openclaw/Workspace.vue'),
        meta: { title: '工作空间', icon: 'FolderOpened', requireAuth: true }
      },
      // 远程访问管理
      {
        path: 'remote-access',
        name: 'RemoteAccess',
        component: () => import('@/views/access/RemoteAccess.vue'),
        meta: { title: '远程访问', icon: 'Link', requireAuth: true }
      },
      // ===== 管理员路由 =====
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
        meta: { title: '用户管理', icon: 'Users', requireAuth: true, requireAdmin: true }
      },
      {
        path: 'admin/logs',
        name: 'AdminLogs',
        component: () => import('@/views/admin/Logs.vue'),
        meta: { title: '系统日志', icon: 'Document', requireAuth: true, requireAdmin: true }
      },
      {
        path: 'admin/settings',
        name: 'AdminSettings',
        component: () => import('@/views/admin/Settings.vue'),
        meta: { title: '系统设置', icon: 'Setting', requireAuth: true, requireAdmin: true }
      }
    ]
  },
  // 404 fallback
  {
    path: '/:pathMatch(.*)*',
    redirect: '/vms'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - CinaSeek 云端开发工作室` : 'CinaSeek 云端开发工作室'

  // 不需要认证的页面直接放行
  if (to.meta.requireAuth === false) {
    next()
    return
  }

  const token = localStorage.getItem('token')
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')

  // 未登录 → 跳转登录页
  if (!token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  // 需要管理员权限
  if (to.meta.requireAdmin) {
    if (!userInfo || typeof userInfo.role !== 'number' || userInfo.role < 10) {
      next('/403')
      return
    }
  }

  next()
})

export default router
