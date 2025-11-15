#!/bin/bash

echo "=== 快速修复登录问题 ==="
echo ""

# 检查 MySQL 容器
if ! docker compose ps mysql 2>/dev/null | grep -q "Up"; then
    echo "❌ MySQL 容器未运行"
    exit 1
fi

echo "1. 更新 admin 用户密码..."
docker compose exec -T mysql mysql -u root -proot123456 learn_hub << 'EOF'
UPDATE sys_user 
SET password = '$2a$10$lZiuaxzQSb.5cuKSbV/F1.5dawptkLSm3p42zwEfY4wWuJDlm.qj2',
    status = 1
WHERE username = 'admin';

SELECT '密码更新完成' as result, username, status FROM sys_user WHERE username = 'admin';
EOF

echo ""
echo "2. 验证密码哈希..."
cd /Users/xuanxuanzi/home/s/javapub/learn-hub-v2/backend
go run scripts/test-password.go

echo ""
echo "✅ 修复完成！现在可以使用 admin/admin123 登录"

