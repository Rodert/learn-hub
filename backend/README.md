# Learn Hub Backend

基于 Golang + GORM 的后端服务，提供用户认证和菜单权限管理功能。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **认证**: JWT
- **密码加密**: bcrypt

## 项目结构

```
backend/
├── config/          # 配置管理
├── controllers/     # 控制器层
├── database/        # 数据库初始化
├── middleware/      # 中间件（JWT认证等）
├── models/          # 数据模型
├── routes/          # 路由配置
├── services/        # 业务逻辑层
├── utils/           # 工具函数
├── main.go          # 入口文件
└── go.mod           # 依赖管理
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置数据库连接等信息。

### 3. 创建数据库

```sql
CREATE DATABASE learn_hub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 认证接口

#### 用户登录
```
POST /api/login/account
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123",
  "type": "account"
}
```

#### 获取当前用户
```
GET /api/currentUser
Headers: Authorization: Bearer {token}
```

#### 退出登录
```
POST /api/login/outLogin
Headers: Authorization: Bearer {token}
```

#### 获取验证码
```
GET /api/login/captcha?phone=13800138000
```

### 菜单权限接口

#### 获取菜单列表
```
GET /api/menu/list
Headers: Authorization: Bearer {token}
```

#### 获取用户权限
```
GET /api/user/permissions
Headers: Authorization: Bearer {token}
```

## 默认数据

系统初始化时会自动创建：

- **默认管理员账号**
  - 用户名: `admin`
  - 密码: `admin123`
  - 权限: `admin`

- **默认角色**
  - `admin`: 管理员
  - `user`: 普通用户

- **默认菜单**
  - 欢迎页
  - 管理员页面（需要admin权限）
  - 列表页面

## 开发说明

### 扩展功能

1. **添加新的模型**: 在 `models/` 目录下创建新文件
2. **添加新的服务**: 在 `services/` 目录下创建新文件
3. **添加新的控制器**: 在 `controllers/` 目录下创建新文件，并在 `routes/routes.go` 中注册路由
4. **添加新的中间件**: 在 `middleware/` 目录下创建新文件

### 数据库迁移

使用 GORM 的 AutoMigrate 功能，在 `database/database.go` 的 `AutoMigrate()` 函数中添加新模型。

### 权限扩展

- 在 `models/menu.go` 中可以为菜单添加 `access` 字段
- 在 `services/user_service.go` 的 `GetUserMenus()` 中实现权限过滤逻辑
- 前端通过 `access.ts` 中的权限标识进行权限控制

## 注意事项

1. **生产环境**: 
   - 修改 `JWT_SECRET` 为强随机字符串
   - 使用 HTTPS
   - 配置合适的 CORS 策略

2. **密码安全**: 
   - 默认密码已使用 bcrypt 加密
   - 生产环境应强制用户修改默认密码

3. **Token 管理**: 
   - 当前实现为简单 JWT，生产环境建议添加 token 刷新机制和黑名单

## License

MIT

