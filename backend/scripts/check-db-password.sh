#!/bin/bash

echo "=== 检查数据库中的密码哈希值 ==="
echo ""

docker compose exec -T mysql mysql -u root -proot123456 learn_hub << 'EOF'
SELECT 
    username,
    SUBSTRING(password, 1, 30) as password_hash_preview,
    LENGTH(password) as password_length,
    CASE 
        WHEN password = '$2a$10$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW' THEN '✅ 正确'
        ELSE '❌ 需要更新'
    END as status_check,
    status as user_status
FROM sys_user 
WHERE username = 'admin';
EOF

echo ""
echo "如果显示 '❌ 需要更新'，请执行修复 SQL"

