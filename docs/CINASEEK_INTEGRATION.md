# CinaSeek 集成指南

## 认证集成（CinaToken OAuth）

### 当前状态
- ✅ CinaToken v0.12.0-oauth.1 已发布（9 个 OAuth Provider、39+ LLM 渠道）
- ✅ CinaToken 已并入 CinaSeek 主线（cinagroup/cinaseek main 分支）
- ⏳ CinaRoom 需集成 CinaToken OAuth 实现 SSO

### 集成步骤

#### 1. 后端认证模块调整

**移除独立用户认证：**
```go
// 移除 internal/handler/auth.go 中的注册/登录接口
// 保留 JWT Token 验证中间件
```

**添加 CinaToken OAuth 客户端：**
```go
// internal/oauth/cinatoken.go
package oauth

type CinaTokenClient struct {
    BaseURL    string // https://cinatoken.com
    ClientID   string
    ClientSecret string
}

// 验证 Token 并获取用户信息
func (c *CinaTokenClient) ValidateToken(token string) (*UserInfo, error) {
    // 调用 CinaToken /oauth/userinfo 接口
}
```

#### 2. 前端登录页面调整

**原方案：** 账号密码登录
```vue
<!-- 移除 Login.vue 中的账号密码表单 -->
```

**新方案：** CinaToken OAuth 跳转
```vue
<template>
  <div class="login-container">
    <el-button type="primary" @click="loginWithCinaToken">
      使用 CinaToken 账号登录
    </el-button>
  </div>
</template>

<script setup>
const loginWithCinaToken = () => {
  const redirectUri = encodeURIComponent(window.location.origin + '/callback')
  const oauthUrl = `https://cinatoken.com/oauth/authorize?client_id=${CLIENT_ID}&redirect_uri=${redirectUri}&response_type=code`
  window.location.href = oauthUrl
}
</script>
```

#### 3. 数据库 Schema 调整

**创建独立 schema：**
```sql
-- 在 PostgreSQL 主库（服务器 A）执行
CREATE SCHEMA IF NOT EXISTS cinaroom;

-- 迁移数据表到 cinaroom schema
-- 移除 users 表（使用 CinaToken 用户体系）
-- 保留 vms, vm_snapshots, mounts, openclaw_configs 等业务表
```

**调整后的表结构：**
```
cinaroom.vms
cinaroom.vm_snapshots
cinaroom.vm_logs
cinaroom.vm_metrics
cinaroom.mounts
cinaroom.openclaw_configs
cinaroom.openclaw_logs
cinaroom.remote_access
cinaroom.remote_logs
cinaroom.ip_whitelists
cinaroom.system_settings
```

#### 4. API 路径调整

**原路径：** `/api/v1/auth/login`
**新路径：** `/api/cinaroom/v1/...`（移除认证相关接口）

**保留接口：**
- `GET /api/cinaroom/v1/vm/list` - 虚拟机列表
- `POST /api/cinaroom/v1/vm/create` - 创建虚拟机
- `POST /api/cinaroom/v1/vm/operate` - 虚拟机操作
- `GET /api/cinaroom/v1/mount/list` - 挂载列表
- `POST /api/cinaroom/v1/openclaw/deploy` - OpenClaw 部署
- `GET /api/cinaroom/v1/remote/status` - 远程访问状态
- ...（共 40 个业务接口）

**移除接口：**
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/reset-pwd` - 重置密码
- `POST /api/v1/auth/logout` - 登出
- `POST /api/v1/user/update-pwd` - 修改密码

#### 5. JWT Token 验证

**中间件调整：**
```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"code": 401, "msg": "未授权"})
            c.Abort()
            return
        }
        
        // 调用 CinaToken 验证 Token
        userInfo, err := cinatoken.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"code": 401, "msg": "Token 无效"})
            c.Abort()
            return
        }
        
        // 将用户信息存入上下文
        c.Set("userID", userInfo.ID)
        c.Set("userEmail", userInfo.Email)
        c.Next()
    }
}
```

## 数据库集成（PostgreSQL 主从）

### 连接配置

```go
// internal/database/postgres.go
dsn := "host=43.156.66.122 port=5432 user=cinaroom password=xxx dbname=cinatoken sslmode=require"
// 使用 cinaroom schema
db.Exec("SET search_path TO cinaroom, public")
```

### 主从复制

- **写操作** → 服务器 A（主节点 43.156.66.122）
- **读操作** → 服务器 B（从节点 101.32.108.223）

## 部署集成（Kubernetes）

### Namespace 配置

```yaml
# deploy/k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cinaroom
  labels:
    app: cinaroom
    team: cinagroup
```

### Deployment 配置

```yaml
# deploy/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cinaroom-backend
  namespace: cinaroom
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: backend
        image: cinagroup/cinaroom:latest
        env:
        - name: DATABASE_HOST
          value: "postgresql-master.cinatoken.svc.cluster.local"
        - name: CINATOKEN_OAUTH_URL
          value: "https://cinatoken.com/oauth"
```

## 监控集成（Uptime Kuma）

### 添加监控面板

1. 登录 Uptime Kuma（现有实例）
2. 添加新监控：
   - **名称**: CinaRoom API
   - **URL**: https://api.cinaroom.run/health
   - **检测频率**: 5 分钟（与 CinaToken 一致）
   - **告警渠道**: 企业微信 + 邮件（复用现有配置）

## 下一步

1. ✅ 创建 GitHub 仓库（已完成）
2. ⏳ 集成 CinaToken OAuth（后端）
3. ⏳ 调整前端登录页面
4. ⏳ 创建 K8s 部署配置
5. ⏳ 注册域名 cinaroom.run
6. ⏳ 配置 Cloudflare SSL

---

*最后更新：2026-04-01*
