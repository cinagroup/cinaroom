# CinaRoom Backend API 文档

## 基础信息

- **Base URL**: `/api/v1`
- **认证方式**: JWT Bearer Token
- **请求格式**: `application/json`
- **响应格式**: `application/json`

## 认证说明

除了 `/auth/register`、`/auth/login` 和 `/auth/reset-pwd` 接口外，所有接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <your_token>
```

## 统一响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误描述",
  "data": null
}
```

### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 1. 认证模块 (Auth)

### 1.1 用户注册

**POST** `/auth/register`

**请求体**:
```json
{
  "username": "admin",
  "email": "admin@example.com",
  "password": "Admin123",
  "confirm_password": "Admin123"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    }
  }
}
```

### 1.2 用户登录

**POST** `/auth/login`

**请求体**:
```json
{
  "username": "admin",
  "password": "Admin123",
  "remember": true
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "nickname": "管理员",
      "avatar": "https://..."
    }
  }
}
```

### 1.3 重置密码

**POST** `/auth/reset-pwd`

**请求体**:
```json
{
  "email": "admin@example.com",
  "new_password": "NewAdmin123",
  "code": "123456"
}
```

### 1.4 用户登出

**POST** `/auth/logout`

**认证**: 需要

### 1.5 获取用户信息

**GET** `/auth/user-info`

**认证**: 需要

### 1.6 更新用户信息

**PUT** `/auth/user-info`

**请求体**:
```json
{
  "email": "newemail@example.com",
  "nickname": "新昵称",
  "phone": "13800138000",
  "avatar": "https://..."
}
```

### 1.7 修改密码

**PUT** `/auth/user-pwd`

**请求体**:
```json
{
  "current_password": "OldPassword123",
  "new_password": "NewPassword123"
}
```

### 1.8 获取登录日志

**GET** `/auth/login-logs`

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "login_time": "2024-01-01T10:00:00Z",
      "ip": "192.168.1.1",
      "location": "北京市",
      "device": "Chrome 120.0.0.0"
    }
  ]
}
```

### 1.9 获取活跃会话

**GET** `/auth/sessions`

### 1.10 撤销会话

**POST** `/auth/sessions/revoke`

**请求体**:
```json
{
  "session_id": 1
}
```

---

## 2. 虚拟机管理模块 (VM)

### 2.1 获取虚拟机列表

**GET** `/vm/list`

**查询参数**:
- `name`: 虚拟机名称（模糊搜索）
- `status`: 状态筛选（running/stopped/paused）
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10）

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "name": "my-vm",
        "status": "running",
        "ip": "192.168.1.100",
        "image": "ubuntu:22.04",
        "cpu": 2,
        "memory": 4,
        "disk": 50,
        "network_type": "nat",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 10
  }
}
```

### 2.2 获取虚拟机详情

**GET** `/vm/detail/:id`

### 2.3 创建虚拟机

**POST** `/vm/create`

**请求体**:
```json
{
  "name": "my-vm",
  "image": "ubuntu:22.04",
  "cpu": 2,
  "memory": 4,
  "disk": 50,
  "network_type": "nat",
  "ssh_key": "ssh-rsa AAAA...",
  "init_script": "#!/bin/bash\necho 'Hello'"
}
```

### 2.4 操作虚拟机

**POST** `/vm/operate/:id`

**请求体**:
```json
{
  "operation": "start"
}
```

**操作类型**:
- `start`: 启动
- `stop`: 停止
- `restart`: 重启
- `pause`: 暂停
- `resume`: 恢复
- `delete`: 删除

### 2.5 更新虚拟机配置

**PUT** `/vm/update-config/:id`

**请求体**:
```json
{
  "cpu": 4,
  "memory": 8,
  "disk": 100
}
```

### 2.6 获取快照列表

**GET** `/vm/snapshots/:id`

### 2.7 创建快照

**POST** `/vm/snapshot/:id`

**请求体**:
```json
{
  "name": "backup-2024-01-01"
}
```

### 2.8 恢复快照

**POST** `/vm/snapshot/:id/restore`

**请求体**:
```json
{
  "snapshot_id": 1
}
```

### 2.9 删除快照

**DELETE** `/vm/snapshot/:id/:snapshot_id`

### 2.10 获取操作日志

**GET** `/vm/logs/:id`

### 2.11 获取监控指标

**GET** `/vm/metrics/:id`

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "vm_id": 1,
      "cpu_usage": 25.5,
      "memory_usage": 512.3,
      "disk_io": 1024.5,
      "network_rx": 2048.0,
      "network_tx": 1024.0,
      "timestamp": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

## 3. 目录挂载模块 (Mount)

### 3.1 获取挂载列表

**GET** `/mount/list`

**查询参数**:
- `vm_id`: 虚拟机 ID

### 3.2 添加挂载

**POST** `/mount/add`

**请求体**:
```json
{
  "vm_id": 1,
  "name": "workspace",
  "host_path": "/home/user/workspace",
  "vm_path": "/root/workspace",
  "permission": "rw",
  "auto_mount": true
}
```

### 3.3 操作挂载

**POST** `/mount/operate/:id`

**请求体**:
```json
{
  "operation": "mount",
  "name": "新名称",
  "vm_path": "/new/path",
  "permission": "ro",
  "auto_mount": false
}
```

**操作类型**:
- `mount`: 挂载
- `unmount`: 卸载
- `edit`: 编辑
- `delete`: 删除

### 3.4 获取 OpenClaw 配置

**GET** `/mount/openclaw-config?vm_id=1`

### 3.5 配置 OpenClaw 挂载

**POST** `/mount/openclaw-config`

**请求体**:
```json
{
  "vm_id": 1,
  "workspace_path": "/home/user/openclaw/workspace",
  "skills_path": "/home/user/openclaw/skills",
  "sync_openclaw_json": true,
  "sync_tool_configs": true
}
```

---

## 4. OpenClaw 管理模块

### 4.1 获取状态

**GET** `/openclaw/status?vm_id=1`

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "running",
    "version": "2026.3.24",
    "running_time": 86400,
    "last_deployed_at": "2024-01-01T10:00:00Z"
  }
}
```

### 4.2 部署 OpenClaw

**POST** `/openclaw/deploy`

**请求体**:
```json
{
  "vm_id": 1,
  "version": "latest",
  "api_key": "sk-...",
  "default_model": "qwencode/qwen3.5-plus"
}
```

### 4.3 操作 OpenClaw

**POST** `/openclaw/operate/:id`

**请求体**:
```json
{
  "operation": "start",
  "version": "2026.3.24"
}
```

**操作类型**:
- `start`: 启动
- `stop`: 停止
- `restart`: 重启
- `update`: 更新

### 4.4 获取日志

**GET** `/openclaw/log/:id`

### 4.5 更新配置

**PUT** `/openclaw/config/:id`

**请求体**:
```json
{
  "default_model": "qwencode/qwen3.5-plus",
  "api_key": "sk-...",
  "enabled_tools": ["browser", "web_search"],
  "enabled_skills": ["feishu-doc", "wecom-doc"]
}
```

### 4.6 获取监控数据

**GET** `/openclaw/monitor?vm_id=1`

### 4.7 获取工作空间

**GET** `/openclaw/workspace?vm_id=1`

---

## 5. 远程访问模块 (Remote)

### 5.1 获取状态

**GET** `/remote/status?vm_id=1`

### 5.2 切换开关

**PUT** `/remote/switch/:id`

**请求体**:
```json
{
  "enabled": true
}
```

### 5.3 获取 IP 白名单

**GET** `/remote/ip-whitelist?vm_id=1`

### 5.4 添加白名单

**POST** `/remote/ip-whitelist`

**请求体**:
```json
{
  "vm_id": 1,
  "ip": "192.168.1.0/24",
  "note": "办公室网络"
}
```

### 5.5 删除白名单

**DELETE** `/remote/ip-whitelist/:id/:whitelist_id`

### 5.6 获取访问日志

**GET** `/remote/log/:id`

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `ip`: IP 筛选
- `status`: 状态码筛选

---

## 6. 系统模块 (System)

### 6.1 获取设置

**GET** `/system/setting?key=xxx`

### 6.2 更新设置

**PUT** `/system/setting`

**请求体**:
```json
{
  "key": "system_name",
  "value": "CinaRoom 管理平台"
}
```

### 6.3 获取版本

**GET** `/system/version`

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": "1.0.0",
    "build": "20260401",
    "api_version": "v1",
    "go_version": "1.21.5",
    "features": [...]
  }
}
```

### 6.4 获取仪表盘

**GET** `/system/dashboard`

### 6.5 获取统计

**GET** `/system/statistics`

### 6.6 搜索虚拟机

**GET** `/system/search?keyword=xxx`

### 6.7 批量操作虚拟机

**POST** `/system/batch-vm`

**请求体**:
```json
{
  "vm_ids": [1, 2, 3],
  "operation": "start"
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（Token 无效或过期） |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 限流说明

API 接口有限流保护，默认限制：
- 每秒 10 个请求
- 突发允许 20 个请求

超过限制将返回 `429 Too Many Requests` 错误。

## 版本历史

- v1.0.0 (2026-04-01): 初始版本，包含完整的用户管理、虚拟机管理、OpenClaw 集成等功能
