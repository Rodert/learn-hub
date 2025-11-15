# 故障排查指南

## 登录失败问题

如果登录接口返回 `{"status":"error","type":"account","currentAuthority":"guest"}`，可能的原因：

### 1. 检查数据库是否已初始化

查看后端启动日志，应该看到：
```
数据库连接成功
数据库迁移完成
基础数据初始化完成
```

### 2. 检查数据库中的用户

连接到 MySQL 数据库：

```bash
# 使用 docker compose
docker compose exec mysql mysql -u root -proot123456 learn_hub

# 或者使用本地 MySQL 客户端
mysql -h localhost -P 3306 -u root -proot123456 learn_hub
```

然后执行：

```sql
-- 查看用户表
SELECT id, username, name, access, status FROM sys_user;

-- 查看用户密码（加密后的）
SELECT id, username, password FROM sys_user WHERE username = 'admin';

-- 如果用户不存在，手动创建（密码: admin123）
-- 注意：密码需要使用 bcrypt 加密
```

### 3. 手动初始化数据

如果数据库中没有用户，可以：

**方法1：删除数据库重新初始化**

```bash
# 停止后端服务
# 删除数据库
docker compose exec mysql mysql -u root -proot123456 -e "DROP DATABASE IF EXISTS learn_hub;"
docker compose exec mysql mysql -u root -proot123456 -e "CREATE DATABASE learn_hub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 重新启动后端服务，会自动初始化数据
go run main.go
```

**方法2：手动插入用户**

```sql
-- 插入管理员用户（密码: admin123）
-- 这个密码哈希值对应 "admin123"
INSERT INTO sys_user (username, password, name, email, userid, access, status, created_at, updated_at)
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iwK8pJ5C', '管理员', 'admin@example.com', '00000001', 'admin', 1, NOW(), NOW());
```

### 4. 检查密码哈希

默认密码 `admin123` 的 bcrypt 哈希值应该是：
```
$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iwK8pJ5C
```

如果数据库中的哈希值不同，需要重新生成或更新。

### 5. 检查用户状态

确保用户状态为 1（启用）：

```sql
UPDATE sys_user SET status = 1 WHERE username = 'admin';
```

## 前端连接后端

### 1. 确认后端服务运行

```bash
# 检查后端是否在运行
curl http://localhost:8080/api/login/captcha?phone=13800138000
```

应该返回：
```json
{"success":true,"data":{"captcha":"1234","code":200,"status":"ok"}}
```

### 2. 前端代理配置

已配置 `frontend-admin/config/proxy.ts`：

```typescript
dev: {
  '/api/': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    pathRewrite: { '^': '' },
  },
},
```

### 3. 启动前端

```bash
cd frontend-admin
npm install  # 或 pnpm install
npm start    # 或 pnpm start
```

前端会在 `http://localhost:8000` 启动，所有 `/api/*` 请求会自动代理到 `http://localhost:8080`

## 常见错误

### 错误1: `driver: bad connection`

**原因**: MySQL 容器未完全启动

**解决**: 
```bash
# 等待 MySQL 启动（查看日志）
docker compose logs -f mysql

# 看到 "Ready for start up" 后再启动后端
```

### 错误2: `Access denied for user`

**原因**: 数据库密码配置不正确

**解决**: 检查 `.env` 文件中的 `DB_PASSWORD` 是否与 docker-compose.yml 中的 `MYSQL_ROOT_PASSWORD` 一致

### 错误3: 登录返回 error

**原因**: 
- 用户不存在
- 密码错误
- 用户被禁用

**解决**: 按照上面的步骤检查数据库中的用户数据

## 测试脚本

使用测试脚本快速检查：

```bash
# 测试 API 接口
./backend/scripts/test-api.sh

# 检查 MySQL 连接
./backend/scripts/check-mysql.sh
```

