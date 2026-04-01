#!/bin/bash
set -e

echo "🚀 CinaRoom 数据库迁移脚本 (Phase 2)"
echo "======================================"
echo ""

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# 从环境变量读取数据库配置
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-cinaroom}
DB_NAME=${DB_NAME:-cinatoken}
DB_SCHEMA=${DB_SCHEMA:-cinaroom}

echo "📊 数据库连接信息:"
echo "   主机：$DB_HOST:$DB_PORT"
echo "   用户：$DB_USER"
echo "   数据库：$DB_NAME"
echo "   Schema: $DB_SCHEMA"
echo ""

# 检查 psql 是否安装
if ! command -v psql &> /dev/null; then
    echo "❌ psql 未安装，请先安装 PostgreSQL 客户端"
    exit 1
fi

# 提示用户输入密码
echo "🔐 请输入数据库密码："
read -s PGPASSWORD
export PGPASSWORD

echo ""
echo "📝 开始执行迁移脚本..."
echo ""

# 执行迁移
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/migrate.sql

# 检查执行结果
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 数据库迁移完成！"
    echo ""
    echo "📊 验证迁移结果:"
    echo "   运行以下 SQL 检查表结构:"
    echo "   \\d users"
    echo ""
    echo "   或执行以下查询查看新增字段:"
    echo "   SELECT column_name, data_type, is_nullable"
    echo "   FROM information_schema.columns"
    echo "   WHERE table_name = 'users'"
    echo "   AND column_name IN ('cinatoken_id', 'provider', 'active');"
else
    echo ""
    echo "❌ 数据库迁移失败！"
    exit 1
fi

# 清除密码
unset PGPASSWORD
