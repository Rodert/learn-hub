-- 诊断 SQL 脚本
-- 用于排查登录问题

-- 1. 检查用户是否存在及基本信息
SELECT 
    id, 
    username, 
    name, 
    email,
    access,
    status,
    created_at,
    updated_at
FROM sys_user 
WHERE username = 'admin';

-- 2. 检查密码哈希值（显示前30个字符和总长度）
SELECT 
    username,
    SUBSTRING(password, 1, 30) as password_hash_preview,
    LENGTH(password) as password_length,
    CASE 
        WHEN password LIKE '$2a$10$%' THEN '✅ 格式正确'
        ELSE '❌ 格式错误'
    END as hash_format_check
FROM sys_user 
WHERE username = 'admin';

-- 3. 检查用户角色关联
SELECT 
    u.id as user_id,
    u.username,
    u.access,
    r.id as role_id,
    r.code as role_code,
    r.name as role_name,
    r.status as role_status
FROM sys_user u
LEFT JOIN sys_user_role ur ON u.id = ur.user_id
LEFT JOIN sys_role r ON ur.role_id = r.id
WHERE u.username = 'admin';

-- 4. 检查角色是否存在
SELECT id, code, name, status FROM sys_role;

-- 5. 检查菜单数据
SELECT id, parent_id, name, path, component, icon, access, status 
FROM sys_menu 
ORDER BY parent_id, sort_order;

-- 6. 检查角色菜单关联
SELECT 
    r.code as role_code,
    m.name as menu_name,
    m.path as menu_path,
    m.access as menu_access
FROM sys_role r
JOIN sys_role_menu rm ON r.id = rm.role_id
JOIN sys_menu m ON rm.menu_id = m.id
WHERE r.code = 'admin'
ORDER BY m.parent_id, m.sort_order;

