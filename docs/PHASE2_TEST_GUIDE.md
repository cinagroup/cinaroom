# Phase 2 测试指南 - CinaToken OAuth 集成

## 📋 测试前准备

### 1. 环境检查

```bash
# 检查 Go 版本（需要 1.23+）
go version

# 检查 Node.js 版本（需要 18+）
node -v

# 检查 PostgreSQL 连接
psql -h <DB_HOST> -U <DB_USER> -d <DB_NAME> -c "SELECT 1"
```

### 2. 配置 OAuth 凭证

**在 CinaToken 管理后台**:
1. 登录 https://cinatoken.com/admin
2. 进入 "OAuth 应用管理" → "创建应用"
3. 填写应用信息：
   - **应用名称**: CinaRoom
   - **应用描述**: CinaRoom 远程管理平台
   - **回调地址**: http://localhost:3000/oauth/callback
   - **Scopes**: user:read user:email
4. 保存后获取：
   - `Client ID`: `cinaroom_xxx`
   - `Client Secret`: `sk_xxx`

**更新后端配置**:
```bash
cd backend
cp .env.example .env
vim .env

# 修改以下配置
CINATOKEN_CLIENT_ID=cinaroom_xxx
CINATOKEN_CLIENT_SECRET=sk_xxx
CINATOKEN_REDIRECT_URI=http://localhost:3000/oauth/callback
```

### 3. 执行数据库迁移

```bash
cd backend

# 方式 1：使用迁移脚本
chmod +x scripts/migrate.sh
./scripts/migrate.sh

# 方式 2：手动执行 SQL
psql -h <DB_HOST> -U <DB_USER> -d <DB_NAME> -f scripts/migrate.sql

# 验证迁移结果
psql -h <DB_HOST> -U <DB_USER> -d <DB_NAME> -c "\d users"
```

**预期输出**:
```
Table "public.users"
     Column     |            Type
----------------+------------------
 id             | integer
 cinatoken_id   | integer          ← 新增
 provider       | character varying(50)  ← 新增
 active         | boolean          ← 新增
 username       | character varying(20)
 email          | character varying(100)
 ...
```

---

## 🧪 测试步骤

### Step 1: 启动后端服务

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

**预期输出**:
```
Server starting on port 8080
```

### Step 2: 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

**预期输出**:
```
VITE v5.x.x  ready in xxx ms
➜  Local:   http://localhost:3000/
```

### Step 3: 测试 OAuth 登录流程

#### 3.1 访问登录页

浏览器打开：http://localhost:3000/login

**检查点**:
- ✅ 页面显示 "使用 CinaToken 登录" 按钮
- ✅ 页面显示 "传统账号密码登录" 区域
- ✅ 点击 "查看支持的登录方式" 显示 9+ Provider

#### 3.2 点击 OAuth 登录

点击 "使用 CinaToken 登录" 按钮

**检查点**:
- ✅ 浏览器跳转到 CinaToken 授权页
- ✅ URL 格式：`https://cinatoken.com/oauth/authorize?client_id=xxx&redirect_uri=xxx&state=xxx`
- ✅ 显示 CinaToken 支持的登录方式（GitHub/Google/Microsoft 等）

#### 3.3 授权登录

在 CinaToken 授权页选择一个 Provider（如 GitHub）并完成授权

**检查点**:
- ✅ 授权成功后跳转回 http://localhost:3000/oauth/callback?code=xxx&state=xxx
- ✅ 页面显示 "正在验证授权..."
- ✅ 页面显示 "登录成功，正在跳转..."
- ✅ 自动跳转到 /vms 页面

#### 3.4 验证登录状态

在 /vms 页面检查：

**检查点**:
- ✅ 顶部导航栏显示用户名
- ✅ 可以访问需要认证的页面
- ✅ localStorage 中存在 token

**浏览器控制台**:
```javascript
localStorage.getItem('token')
// 应该返回 JWT token 字符串
```

---

## 🔬 API 测试

### 使用 curl 测试 OAuth 接口

#### 1. 获取授权 URL

```bash
curl -X GET "http://localhost:8080/api/v1/oauth/authorize" \
  -H "Content-Type: application/json"
```

**预期响应**:
```json
{
  "code": 200,
  "data": {
    "authorize_url": "https://cinatoken.com/oauth/authorize?client_id=xxx&..."
  }
}
```

#### 2. 获取 OAuth Provider 列表

```bash
curl -X GET "http://localhost:8080/api/v1/oauth/providers"
```

**预期响应**:
```json
{
  "code": 200,
  "data": {
    "providers": [
      {"name": "github", "display_name": "GitHub", "enabled": true},
      {"name": "google", "display_name": "Google", "enabled": true},
      ...
    ],
    "note": "通过 CinaToken 统一认证，支持 9+ OAuth 提供商"
  }
}
```

#### 3. 测试 OAuth 回调（需要真实授权码）

```bash
# 注意：这个测试需要先从 CinaToken 获取真实的授权码
curl -X GET "http://localhost:8080/api/v1/oauth/callback?code=AUTH_CODE&state=STATE"
```

---

## 🧪 单元测试

### 运行后端测试

```bash
cd backend
go test -v ./tests/oauth_test.go -run TestCinaTokenClient
```

**预期输出**:
```
=== RUN   TestCinaTokenClient_GetAuthorizationURL
--- PASS: TestCinaTokenClient_GetAuthorizationURL (0.00s)
PASS
ok      multipass-backend/tests    0.003s
```

### 运行所有 OAuth 测试

```bash
go test -v ./tests/... -run OAuth
```

**预期输出**:
```
=== RUN   TestOAuthHandler_Providers
--- PASS: TestOAuthHandler_Providers (0.00s)
=== RUN   TestOAuthConfig_EnvironmentVariables
--- PASS: TestOAuthConfig_EnvironmentVariables (0.00s)
PASS
```

---

## 🐛 常见问题排查

### 1. "CinaToken Client ID 无效"

**原因**: 配置的 Client ID 不正确或已过期

**解决**:
1. 检查 `.env` 文件中的 `CINATOKEN_CLIENT_ID`
2. 在 CinaToken 管理后台重新生成 Client ID
3. 重启后端服务

### 2. "State 验证失败"

**原因**: CSRF state 不匹配

**解决**:
1. 检查浏览器是否启用了 Cookie
2. 清除浏览器缓存后重试
3. 检查 `oauth_state` cookie 是否正确设置

### 3. "回调地址不匹配"

**原因**: CinaToken 后台配置的回调地址与实际不符

**解决**:
1. 在 CinaToken 管理后台检查回调地址配置
2. 确保配置为 `http://localhost:3000/oauth/callback`
3. 保存后等待 1-2 分钟生效

### 4. "数据库字段不存在"

**原因**: 未执行数据库迁移

**解决**:
```bash
cd backend
./scripts/migrate.sh
```

### 5. "CORS 错误"

**原因**: 跨域配置问题

**解决**:
```bash
# 检查 .env 中的 CORS 配置
CORS_ALLOW_ORIGINS=http://localhost:3000

# 或临时允许所有来源
CORS_ALLOW_ORIGINS=*
```

---

## ✅ 验收标准

完成以下检查清单：

- [ ] 后端服务启动成功（端口 8080）
- [ ] 前端服务启动成功（端口 3000）
- [ ] 登录页面显示 OAuth 按钮
- [ ] 点击 OAuth 按钮跳转到 CinaToken
- [ ] 授权成功后跳转回回调页面
- [ ] 回调页面显示 "登录成功"
- [ ] 自动跳转到 /vms 页面
- [ ] 顶部导航栏显示用户名
- [ ] localStorage 中存在 token
- [ ] 单元测试全部通过
- [ ] 数据库迁移成功（users 表有 cinatoken_id 字段）

---

## 📊 性能指标

**预期性能**:
- OAuth 授权跳转延迟：< 200ms
- Token 交换延迟：< 500ms
- 用户信息获取延迟：< 300ms
- 整体登录流程：< 2 秒

**测试方法**:
```bash
# 使用 curl 测试响应时间
curl -w "@curl-format.txt" -o /dev/null -s "http://localhost:8080/api/v1/oauth/providers"
```

**curl-format.txt**:
```
time_namelookup:  %{time_namelookup}\n
time_connect:     %{time_connect}\n
time_starttransfer: %{time_starttransfer}\n
time_total:       %{time_total}\n
```

---

## 📝 测试报告模板

完成测试后，填写以下报告：

```markdown
## Phase 2 测试报告

**测试时间**: 2026-04-01
**测试人员**: [姓名]
**测试环境**: 
- Go: [版本]
- Node.js: [版本]
- PostgreSQL: [版本]

### 测试结果

| 测试项 | 状态 | 备注 |
|--------|------|------|
| OAuth 授权跳转 | ✅/❌ | |
| Token 交换 | ✅/❌ | |
| 用户信息获取 | ✅/❌ | |
| 本地用户创建 | ✅/❌ | |
| JWT Token 生成 | ✅/❌ | |
| 前端回调处理 | ✅/❌ | |
| 单元测试 | ✅/❌ | [通过率] |

### 发现的问题

1. [问题描述]
   - 严重程度：高/中/低
   - 解决方案：[描述]

### 结论

Phase 2 OAuth 集成测试 [通过/未通过]，[可以/不可以] 进入 Phase 3。
```

---

**测试完成后，将报告保存到 `docs/PHASE2_TEST_REPORT.md`**
