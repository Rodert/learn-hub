-- 使用正确的密码哈希值更新 admin 用户
-- 这个哈希值对应密码: admin123
-- 已通过 bcrypt 验证，确保正确

UPDATE sys_user 
SET password = '$2a$10$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW',
    status = 1
WHERE username = 'admin';

-- 验证更新结果
SELECT 
    username,
    SUBSTRING(password, 1, 30) as password_hash_preview,
    LENGTH(password) as password_length,
    status
FROM sys_user 
WHERE username = 'admin';

