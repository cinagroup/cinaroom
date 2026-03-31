# Multipass 远程管理平台

基于 Vue3 + Vite + Element Plus 的 Multipass 虚拟机远程管理平台。

## 技术栈

- **前端框架**: Vue 3.4+
- **构建工具**: Vite 5+
- **UI 组件库**: Element Plus 2.5+
- **状态管理**: Pinia 2.1+
- **路由管理**: Vue Router 4.3+
- **HTTP 客户端**: Axios 1.6+
- **终端组件**: xterm.js 5.3+
- **图表库**: ECharts 5.4+

## 功能模块

1. **用户中心**: 注册/登录/个人信息/安全设置
2. **虚拟机管理**: 列表/创建/操作/详情/快照
3. **Web 远程管理**: WebShell/远程管控/日志查看
4. **目录挂载**: 挂载管理/OpenClaw 专属配置
5. **OpenClaw 管理**: 部署/管控/监控/工作空间
6. **远程访问管理**: 状态/开关/IP 白名单/日志

## 快速开始

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 项目结构

```
multipass-frontend/
├── src/
│   ├── api/              # API 接口
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   │   ├── Layout.vue    # 布局组件
│   │   └── Header.vue    # 头部组件
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia 状态管理
│   ├── styles/           # 全局样式
│   ├── utils/            # 工具函数
│   ├── views/            # 页面组件
│   │   ├── user/         # 用户中心
│   │   ├── vms/          # 虚拟机管理
│   │   ├── shell/        # WebShell
│   │   ├── remote/       # 远程管控
│   │   ├── logs/         # 日志查看
│   │   ├── mounts/       # 目录挂载
│   │   ├── openclaw/     # OpenClaw 管理
│   │   └── access/       # 远程访问
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── index.html
├── package.json
├── vite.config.js
└── README.md
```

## 页面清单

| 页面 | 路由 | 说明 |
|------|------|------|
| 登录 | /login | 用户登录 |
| 注册 | /register | 用户注册 |
| 个人信息 | /profile | 个人信息管理 |
| 安全设置 | /security | 安全配置 |
| 虚拟机列表 | /vms | 虚拟机管理 |
| 创建虚拟机 | /vms/create | 新建虚拟机 |
| 虚拟机详情 | /vms/:id | 虚拟机详情与监控 |
| WebShell | /shell/:id | Web 终端 |
| 远程管控 | /remote/:id | 文件/进程/服务管理 |
| 日志查看 | /logs/:id | 系统日志查看 |
| 挂载管理 | /mounts | 目录挂载配置 |
| OpenClaw 部署 | /openclaw/deploy | 部署管理 |
| OpenClaw 配置 | /openclaw/config | 管控配置 |
| OpenClaw 监控 | /openclaw/monitor | 监控面板 |
| OpenClaw 工作空间 | /openclaw/workspace | 工作空间管理 |
| 远程访问 | /remote-access | 远程访问配置 |

## 开发规范

- 使用 ESLint + Prettier 保持代码风格
- 组件采用 Composition API (setup 语法糖)
- 状态管理使用 Pinia
- 样式使用 SCSS

## License

MIT
