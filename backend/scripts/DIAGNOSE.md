# 诊断指南

## 执行诊断 SQL

请执行以下 SQL 语句并将结果发给我：

```bash
# 方法1: 使用 docker compose
cd backend
cat scripts/diagnose.sql | docker compose exec -T mysql mysql -u root -proot123456 learn_hub

# 方法2: 直接连接 MySQL
mysql -h localhost -P 3306 -u root -proot123456 learn_hub < scripts/diagnose.sql
```

## 或者手动执行以下 SQL

### 1. 检查用户基本信息
```sql
SELECT id, username, name, access, status 
FROM sys_user 
WHERE username = 'admin';
```

### 2. 检查密码哈希值（关键）
```sql
SELECT 
    username,
    SUBSTRING(password, 1, 30) as password_hash_preview,
    LENGTH(password) as password_length,
    CASE 
        WHEN password LIKE '$2a$10$%' THEN '格式正确'
        ELSE '格式错误'
    END as format_check
FROM sys_user 
WHERE username = 'admin';
```

### 3. 检查用户角色关联
```sql
SELECT 
    u.username,
    u.access,
    r.code as role_code,
    r.name as role_name
FROM sys_user u
LEFT JOIN sys_user_role ur ON u.id = ur.user_id
LEFT JOIN sys_role r ON ur.role_id = r.id
WHERE u.username = 'admin';
```

## 测试密码验证

也可以运行 Go 程序测试密码验证逻辑：

```bash
cd backend
go run scripts/test-password.go
```

## 常见问题

### 问题1: 密码哈希值不正确
如果密码哈希值不是以 `$2a$10$` 开头或长度不是 60，需要更新：

```sql
UPDATE sys_user 
SET password = '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iwK8pJ5C'
WHERE username = 'admin';
```

### 问题2: 用户状态为 0（禁用）
```sql
UPDATE sys_user 
SET status = 1 
WHERE username = 'admin';
```

### 问题3: 用户不存在
执行初始化脚本：
```bash
cat scripts/init-user.sql | docker compose exec -T mysql mysql -u root -proot123456 learn_hub
```

