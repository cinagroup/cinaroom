-- CinaRoom Phase 2 数据库迁移脚本
-- 用于集成 CinaToken OAuth 认证

-- ===========================================
-- 1. 添加 CinaToken OAuth 相关字段
-- ===========================================

-- 添加 cinatoken_id 字段（CinaToken 用户 ID）
ALTER TABLE users ADD COLUMN IF NOT EXISTS cinatoken_id INTEGER;

-- 添加 provider 字段（OAuth 提供商：github/google/microsoft 等）
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider VARCHAR(50);

-- 添加 active 字段（用户激活状态）
ALTER TABLE users ADD COLUMN IF NOT EXISTS active BOOLEAN DEFAULT true;

-- ===========================================
-- 2. 创建唯一索引
-- ===========================================

-- cinatoken_id 唯一索引（确保一个 CinaToken 账号只对应一个本地用户）
CREATE INDEX IF NOT EXISTS idx_users_cinatoken_id ON users(cinatoken_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_cinatoken_id_unique ON users(cinatoken_id) WHERE cinatoken_id IS NOT NULL;

-- ===========================================
-- 3. 调整现有索引
-- ===========================================

-- 如果 username 和 email 的唯一索引不存在，创建它们
-- 注意：如果已存在唯一约束，这些语句会失败，可以忽略

-- ===========================================
-- 4. 数据迁移（可选）
-- ===========================================

-- 为现有用户设置默认 active 状态
UPDATE users SET active = true WHERE active IS NULL;

-- 为现有用户设置默认 provider（空字符串表示传统账号密码登录）
UPDATE users SET provider = '' WHERE provider IS NULL;

-- ===========================================
-- 5. 验证迁移结果
-- ===========================================

-- 检查表结构
-- \d users

-- 检查索引
-- SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'users';

-- ===========================================
-- 回滚脚本（如需要）
-- ===========================================

-- DROP INDEX IF EXISTS idx_users_cinatoken_id_unique;
-- DROP INDEX IF EXISTS idx_users_cinatoken_id;
-- ALTER TABLE users DROP COLUMN IF EXISTS active;
-- ALTER TABLE users DROP COLUMN IF EXISTS provider;
-- ALTER TABLE users DROP COLUMN IF EXISTS cinatoken_id;
