# Learn Hub 数据库启动指南

## 📋 环境要求

- Docker 或 Docker Desktop
- 或本地 MySQL 8.0+

---

## 🐳 方案 1: 使用 Docker (推荐)

### 前置条件
- 安装 Docker Desktop (https://www.docker.com/products/docker-desktop)
- 启动 Docker Desktop

### 启动 MySQL 容器

```bash
# 进入项目目录
cd /Users/xuanxuanzi/home/s/javapub/learn-hub

# 启动 MySQL 容器
docker-compose up -d mysql

# 等待 MySQL 启动完成（约 30 秒）
docker-compose logs -f mysql

# 当看到 "ready for connections" 时，MySQL 已启动成功
```

### 验证 MySQL 连接

```bash
# 进入 MySQL 容器
docker exec -it learn-hub-mysql mysql -u root -ppassword

# 查看数据库
SHOW DATABASES;

# 退出
EXIT;
```

---

## 💻 方案 2: 使用本地 MySQL

### 安装 MySQL

#### macOS (使用 Homebrew)
```bash
# 安装 Homebrew (如果未安装)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装 MySQL
brew install mysql

# 启动 MySQL 服务
brew services start mysql

# 验证安装
mysql --version
```

#### macOS (使用 MacPorts)
```bash
# 安装 MacPorts (如果未安装)
# 访问 https://www.macports.org/install.php

# 安装 MySQL
sudo port install mysql80-server

# 启动 MySQL
sudo port load mysql80-server
```

#### macOS (使用 DMG 安装包)
1. 下载 MySQL DMG: https://dev.mysql.com/downloads/mysql/
2. 双击安装包按照提示安装
3. 启动 MySQL: System Preferences > MySQL > Start MySQL Server

### 创建数据库

```bash
# 连接到 MySQL
mysql -u root -p

# 输入密码（默认为空或你设置的密码）

# 创建数据库
CREATE DATABASE IF NOT EXISTS learn_hub DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 查看数据库
SHOW DATABASES;

# 退出
EXIT;
```

### 修改后端配置

编辑 `backend/config/config.yaml`:

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  user: root
  password: ""  # 如果没有密码，留空
  dbname: learn_hub
  max_open_conns: 100
  max_idle_conns: 10
```

---

## 🚀 启动后端服务

### 1. 执行数据库迁移

```bash
cd /Users/xuanxuanzi/home/s/javapub/learn-hub/backend

# 执行迁移（创建表和默认数据）
make migrate

# 或直接运行
go run ./cmd/migrate/main.go
```

### 2. 启动后端服务

```bash
# 开发模式（支持热重载）
make dev

# 或生产模式
make run
```

**后端服务地址**: http://localhost:8080  
**Swagger API 文档**: http://localhost:8080/swagger/index.html

---

## 🧪 测试数据库连接

### 使用 curl 测试登录

```bash
# 登录测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 预期响应
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员"
    }
  }
}
```

---

## 🔧 常见问题

### Q: Docker 容器启动失败
A: 
1. 检查 Docker Desktop 是否运行
2. 检查端口 3306 是否被占用: `lsof -i :3306`
3. 删除旧容器: `docker rm learn-hub-mysql`
4. 重新启动: `docker-compose up -d mysql`

### Q: MySQL 连接超时
A:
1. 检查 MySQL 是否启动: `docker ps`
2. 查看日志: `docker logs learn-hub-mysql`
3. 等待 MySQL 完全启动（约 30 秒）

### Q: 数据库迁移失败
A:
1. 确保 MySQL 已启动
2. 检查数据库配置: `backend/config/config.yaml`
3. 检查数据库用户权限
4. 查看错误日志

### Q: 忘记 MySQL 密码
A:
```bash
# 重置 MySQL 密码
docker exec -it learn-hub-mysql mysql -u root -ppassword

# 或删除容器重新启动
docker rm -f learn-hub-mysql
docker-compose up -d mysql
```

---

## 📊 数据库初始化

运行 `make migrate` 后，数据库会自动创建以下内容：

### 表结构
- users - 用户表
- roles - 角色表
- permissions - 权限表
- menus - 菜单表
- user_roles - 用户角色关联表
- role_permissions - 角色权限关联表
- materials - 学习资料表
- questions - 题目表
- exams - 考试表
- exam_records - 考试记录表
- course_records - 学习记录表

### 默认数据
- 3 个角色: user, admin, system_admin
- 11 个权限
- 6 个菜单项
- 1 个管理员账户: admin/admin123

---

## ✅ 启动检查清单

- [ ] MySQL 已启动
- [ ] 数据库 `learn_hub` 已创建
- [ ] 后端配置文件已修改
- [ ] 数据库迁移已执行
- [ ] 后端服务已启动
- [ ] 可以访问 http://localhost:8080/swagger/index.html
- [ ] 可以用 admin/admin123 登录

---

**最后更新**: 2025-11-14  
**版本**: 1.0.0
