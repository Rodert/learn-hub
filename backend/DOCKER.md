# Docker 使用指南

## 使用 Docker Compose 启动 MySQL

### 快速开始

1. **启动 MySQL 服务**

```bash
docker-compose up -d
```

2. **查看服务状态**

```bash
docker-compose ps
```

3. **查看日志**

```bash
docker-compose logs -f mysql
```

4. **停止服务**

```bash
docker-compose down
```

5. **停止并删除数据卷（清空数据）**

```bash
docker-compose down -v
```

### 配置说明

#### MySQL 配置

- **镜像版本**: MySQL 8.0
- **端口映射**: 3306:3306
- **Root 密码**: `root123456`
- **数据库名**: `learn_hub`
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_unicode_ci

#### 环境变量

默认配置：
- `MYSQL_ROOT_PASSWORD`: root123456
- `MYSQL_DATABASE`: learn_hub
- `MYSQL_USER`: learnhub
- `MYSQL_PASSWORD`: learnhub123

### 连接信息

#### 从宿主机连接

```bash
mysql -h 127.0.0.1 -P 3306 -u root -proot123456
```

#### 从 Docker 容器连接

```bash
docker-compose exec mysql mysql -u root -proot123456 learn_hub
```

### 后端服务配置

创建 `.env` 文件（或使用 `.env.docker`）：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root123456
DB_NAME=learn_hub
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

### 数据持久化

数据存储在 Docker volume `mysql_data` 中，即使容器删除，数据也会保留。

如需清空数据：

```bash
docker-compose down -v
```

### 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f mysql

# 进入 MySQL 容器
docker-compose exec mysql bash

# 执行 SQL
docker-compose exec mysql mysql -u root -proot123456 learn_hub -e "SHOW TABLES;"

# 备份数据库
docker-compose exec mysql mysqldump -u root -proot123456 learn_hub > backup.sql

# 恢复数据库
docker-compose exec -T mysql mysql -u root -proot123456 learn_hub < backup.sql
```

### 修改配置

如需修改 MySQL 配置，编辑 `docker-compose.yml` 文件中的环境变量，然后重启服务：

```bash
docker-compose down
docker-compose up -d
```

### 故障排查

#### 1. 端口被占用

如果 3306 端口被占用，可以修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "3307:3306"  # 改为其他端口
```

同时修改 `.env` 文件中的 `DB_PORT=3307`

#### 2. 连接失败

- 检查 MySQL 容器是否正常运行：`docker-compose ps`
- 查看日志：`docker-compose logs mysql`
- 检查防火墙设置

#### 3. 字符集问题

MySQL 8.0 默认使用 utf8mb4，如果遇到字符集问题，检查：
- 数据库字符集：`SHOW CREATE DATABASE learn_hub;`
- 表字符集：`SHOW CREATE TABLE sys_user;`

### 生产环境建议

1. **修改默认密码**
   - 修改 `MYSQL_ROOT_PASSWORD`
   - 使用强密码

2. **数据备份**
   - 定期备份数据卷
   - 使用定时任务自动备份

3. **网络安全**
   - 不要暴露 3306 端口到公网
   - 使用 Docker 网络隔离
   - 配置防火墙规则

4. **资源限制**
   - 添加内存和 CPU 限制
   - 监控资源使用情况

示例配置：

```yaml
services:
  mysql:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

