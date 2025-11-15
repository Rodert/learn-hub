# 快速开始指南

## 前置要求

- Go 1.21 或更高版本
- MySQL 8.0 或更高版本

## 安装步骤

### 1. 安装依赖

```bash
go mod download
```

### 2. 创建数据库

```sql
CREATE DATABASE learn_hub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 配置环境变量

复制 `.env.example` 为 `.env`：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置数据库连接信息：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=learn_hub
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

### 4. 运行服务

```bash
go run main.go
```

或者使用 Makefile：

```bash
make run
```

服务将在 `http://localhost:8080` 启动。

## 测试接口

### 1. 登录

```bash
curl -X POST http://localhost:8080/api/login/account \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123",
    "type": "account"
  }'
```

响应示例：
```json
{
  "status": "ok",
  "type": "account",
  "currentAuthority": "admin",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. 获取当前用户信息

```bash
curl -X GET http://localhost:8080/api/currentUser \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 获取菜单列表

```bash
curl -X GET http://localhost:8080/api/menu/list \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. 获取用户权限

```bash
curl -X GET http://localhost:8080/api/user/permissions \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 默认账号

- **用户名**: `admin`
- **密码**: `admin123`
- **权限**: `admin`

## 前端对接

修改前端项目的 `config/proxy.ts`，将代理指向后端服务：

```typescript
dev: {
  '/api/': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    pathRewrite: { '^': '' },
  },
},
```

## 常见问题

### 1. 数据库连接失败

- 检查 MySQL 服务是否启动
- 检查 `.env` 文件中的数据库配置是否正确
- 确认数据库已创建

### 2. Token 验证失败

- 检查请求头中是否包含 `Authorization: Bearer <token>`
- 确认 token 未过期（默认24小时）
- 检查 JWT_SECRET 配置是否正确

### 3. 菜单不显示

- 确认用户已分配角色
- 确认角色已分配菜单
- 检查菜单的 `status` 字段是否为 1（启用）

## 开发建议

1. **扩展功能**: 在对应的目录下添加新文件，保持代码结构清晰
2. **数据库迁移**: 修改模型后，删除数据库重新运行会自动迁移
3. **添加新接口**: 在 `controllers/` 创建控制器，在 `routes/routes.go` 注册路由
4. **权限控制**: 使用 `middleware.AuthMiddleware()` 保护需要认证的接口

