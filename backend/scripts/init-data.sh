#!/bin/bash

# 初始化数据库用户和基础数据

echo "=== 初始化数据库用户和基础数据 ==="
echo ""

# 检查 MySQL 容器是否运行
if ! docker compose ps mysql | grep -q "Up"; then
    echo "❌ MySQL 容器未运行，请先启动: docker compose up -d"
    exit 1
fi

echo "1. 检查数据库连接..."
docker compose exec -T mysql mysql -u root -proot123456 -e "SELECT 1;" learn_hub > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ 数据库连接失败"
    exit 1
fi
echo "✅ 数据库连接正常"
echo ""

echo "2. 检查用户是否存在..."
USER_COUNT=$(docker compose exec -T mysql mysql -u root -proot123456 learn_hub -sN -e "SELECT COUNT(*) FROM sys_user WHERE username='admin';" 2>/dev/null)

if [ "$USER_COUNT" -gt 0 ]; then
    echo "⚠️  用户 admin 已存在"
    read -p "是否要重新初始化？这将删除现有数据 (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "已取消"
        exit 0
    fi
    echo "删除现有数据..."
    docker compose exec -T mysql mysql -u root -proot123456 learn_hub -e "
        DELETE FROM sys_user_role;
        DELETE FROM sys_role_menu;
        DELETE FROM sys_user;
        DELETE FROM sys_role;
        DELETE FROM sys_menu;
    " 2>/dev/null
fi

echo ""
echo "3. 执行初始化脚本..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker compose exec -T mysql mysql -u root -proot123456 learn_hub < "$SCRIPT_DIR/init-user.sql"

if [ $? -eq 0 ]; then
    echo "✅ 初始化成功！"
    echo ""
    echo "默认账号信息:"
    echo "  用户名: admin"
    echo "  密码: admin123"
    echo ""
    echo "现在可以使用以下命令测试登录:"
    echo "  curl -X POST http://localhost:8080/api/login/account \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"admin\",\"password\":\"admin123\",\"type\":\"account\"}'"
else
    echo "❌ 初始化失败"
    exit 1
fi

