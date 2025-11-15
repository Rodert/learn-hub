#!/bin/bash

echo "=== 检查数据库中的用户信息 ==="
echo ""

echo "1. 查看用户表:"
docker compose exec -T mysql mysql -u root -proot123456 learn_hub -e "
SELECT id, username, name, access, status, 
       SUBSTRING(password, 1, 20) as password_hash_preview,
       LENGTH(password) as password_length
FROM sys_user;
" 2>/dev/null

echo ""
echo "2. 查看用户角色关联:"
docker compose exec -T mysql mysql -u root -proot123456 learn_hub -e "
SELECT u.id, u.username, r.code as role_code, r.name as role_name
FROM sys_user u
LEFT JOIN sys_user_role ur ON u.id = ur.user_id
LEFT JOIN sys_role r ON ur.role_id = r.id
WHERE u.username = 'admin';
" 2>/dev/null

echo ""
echo "3. 测试密码哈希值:"
echo "   正确的 admin123 密码哈希应该以: \$2a\$10\$ 开头"
echo "   长度应该是: 60 字符"

