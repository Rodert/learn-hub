# 修复登录问题 - 完整指南

## 问题诊断

如果登录返回 `{"status":"error","type":"account","currentAuthority":"guest"}`，通常是密码哈希值不匹配。

## 解决方案

### 步骤1：执行 SQL 更新密码

**必须执行以下 SQL 来更新数据库中的密码哈希值：**

```sql
UPDATE sys_user 
SET password = '$2a$10$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW',
    status = 1
WHERE username = 'admin';
```

### 步骤2：执行方式

**方法1：使用 SQL 文件（推荐）**
```bash
cd backend
cat scripts/fix-password-correct.sql | docker compose exec -T mysql mysql -u root -proot123456 learn_hub
```

**方法2：直接执行 SQL 命令**
```bash
docker compose exec mysql mysql -u root -proot123456 learn_hub -e "UPDATE sys_user SET password = '\$2a\$10\$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW', status = 1 WHERE username = 'admin';"
```

**方法3：手动连接 MySQL**
```bash
docker compose exec mysql mysql -u root -proot123456 learn_hub
```
然后执行：
```sql
UPDATE sys_user 
SET password = '$2a$10$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW',
    status = 1
WHERE username = 'admin';
```

### 步骤3：验证修复

执行 SQL 后，运行测试程序：

```bash
cd backend
go run scripts/test-password.go
```

应该显示：
```
密码验证结果: true
✅ 密码验证成功！
```

### 步骤4：测试登录

```bash
curl -X POST http://localhost:8080/api/login/account \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","type":"account"}'
```

应该返回：
```json
{
  "status": "ok",
  "type": "account",
  "currentAuthority": "admin",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 重要说明

1. **密码哈希值已通过验证**：`$2a$10$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW` 对应密码 `admin123`
2. **必须执行 SQL**：数据库中的密码哈希值必须更新，否则登录会一直失败
3. **检查后端日志**：如果仍然失败，查看后端日志中的 "密码验证失败" 信息，确认密码哈希值是否正确

## 如果仍然失败

1. 检查后端日志，查看密码哈希值的前30个字符
2. 确认数据库中是否真的更新了
3. 确认用户状态是否为 1（启用）

