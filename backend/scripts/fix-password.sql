-- 修复 admin 用户密码为 admin123
-- 使用正确的 bcrypt 哈希值

UPDATE sys_user 
SET password = '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iwK8pJ5C'
WHERE username = 'admin';

-- 确保用户状态为启用
UPDATE sys_user 
SET status = 1 
WHERE username = 'admin';

-- 验证更新结果
SELECT id, username, name, access, status, 
       SUBSTRING(password, 1, 20) as password_hash_preview
FROM sys_user 
WHERE username = 'admin';

