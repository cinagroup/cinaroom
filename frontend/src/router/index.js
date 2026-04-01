import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/components/Layout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/Register.vue'),
    meta: { title: '注册', requiresAuth: false }
  },
  {
    path: '/oauth/callback',
    name: 'OAuthCallback',
    component: () => import('@/views/user/OAuthCallback.vue'),
    meta: { title: '登录验证', requiresAuth: false }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/vms',
    children: [
      // 用户中心
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人信息', icon: 'User' }
      },
      {
        path: 'security',
        name: 'Security',
        component: () => import('@/views/user/Security.vue'),
        meta: { title: '安全设置', icon: 'Lock' }
      },
      // 虚拟机管理
      {
        path: 'vms',
        name: 'VMs',
        component: () => import('@/views/vms/VMList.vue'),
        meta: { title: '虚拟机列表', icon: 'Monitor' }
      },
      {
        path: 'vms/create',
        name: 'CreateVM',
        component: () => import('@/views/vms/CreateVM.vue'),
        meta: { title: '创建虚拟机', icon: 'Plus' }
      },
      {
        path: 'vms/:id',
        name: 'VMDetail',
        component: () => import('@/views/vms/VMDetail.vue'),
        meta: { title: '虚拟机详情', icon: 'Document' }
      },
      // Web 远程管理
      {
        path: 'shell/:id',
        name: 'WebShell',
        component: () => import('@/views/shell/WebShell.vue'),
        meta: { title: 'WebShell', icon: 'Terminal' }
      },
      {
        path: 'remote/:id',
        name: 'RemoteControl',
        component: () => import('@/views/remote/RemoteControl.vue'),
        meta: { title: '远程管控', icon: 'Setting' }
      },
      {
        path: 'logs/:id',
        name: 'LogViewer',
        component: () => import('@/views/logs/LogViewer.vue'),
        meta: { title: '日志查看', icon: 'Document' }
      },
      // 目录挂载
      {
        path: 'mounts',
        name: 'Mounts',
        component: () => import('@/views/mounts/MountManager.vue'),
        meta: { title: '挂载管理', icon: 'Folder' }
      },
      // OpenClaw 管理
      {
        path: 'openclaw/deploy',
        name: 'OpenClawDeploy',
        component: () => import('@/views/openclaw/Deploy.vue'),
        meta: { title: '部署管理', icon: 'Download' }
      },
      {
        path: 'openclaw/config',
        name: 'OpenClawConfig',
        component: () => import('@/views/openclaw/Config.vue'),
        meta: { title: '管控配置', icon: 'Setting' }
      },
      {
        path: 'openclaw/monitor',
        name: 'OpenClawMonitor',
        component: () => import('@/views/openclaw/Monitor.vue'),
        meta: { title: '监控面板', icon: 'DataAnalysis' }
      },
      {
        path: 'openclaw/workspace',
        name: 'OpenClawWorkspace',
        component: () => import('@/views/openclaw/Workspace.vue'),
        meta: { title: '工作空间', icon: 'FolderOpened' }
      },
      // 远程访问管理
      {
        path: 'remote-access',
        name: 'RemoteAccess',
        component: () => import('@/views/access/RemoteAccess.vue'),
        meta: { title: '远程访问', icon: 'Link' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - CinaSeek 云端开发工作室` : 'CinaSeek 云端开发工作室'
  
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth === false) {
    next()
  } else {
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router
