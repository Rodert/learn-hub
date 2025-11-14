# Learn Hub Backend

Go 后端服务，提供完整的 API 接口支持。

## 项目结构

```
backend/
├── cmd/
│   ├── main.go              # 应用入口
│   └── migrate/
│       └── main.go          # 数据库迁移工具
├── config/
│   ├── config.go            # 配置加载
│   └── config.yaml          # 配置文件
├── internal/
│   ├── api/
│   │   ├── routes.go        # 路由定义
│   │   └── handler/         # 请求处理器
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型
│   ├── service/             # 业务逻辑（可选）
│   └── repository/          # 数据访问层（可选）
├── pkg/
│   └── database/            # 数据库工具
├── docs/                    # Swagger 文档（自动生成）
├── Makefile                 # 构建脚本
├── Dockerfile               # Docker 配置
├── go.mod                   # Go 模块文件
└── README.md                # 本文件
```

## 快速开始

### 前置条件

- Go 1.20+
- MySQL 8.0+
- Docker & Docker Compose（可选）

### 本地开发

1. **安装依赖**

```bash
make deps
```

2. **配置数据库**

编辑 `config/config.yaml`，配置数据库连接信息。

3. **执行迁移**

```bash
make migrate
```

4. **生成 Swagger 文档**

```bash
make swagger
```

5. **运行应用**

```bash
# 开发模式（支持热重载）
make dev

# 或直接运行
make run
```

应用将在 `http://localhost:8080` 启动。

### 访问 Swagger 文档

启动应用后，访问 `http://localhost:8080/swagger/index.html` 查看 API 文档。

## 常用命令

```bash
# 构建二进制文件
make build

# 运行应用
make run

# 开发模式（需要 air）
make dev

# 运行单元测试
make test

# 生成测试覆盖率报告
make test-coverage

# 执行数据库迁移
make migrate

# 生成 Swagger 文档
make swagger

# 代码检查
make lint

# 代码格式化
make fmt

# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 清理构建文件
make clean
```

## 配置说明

### config/config.yaml

```yaml
server:
  port: 8080              # 服务端口
  env: development        # 环境：development/production
  log_level: debug        # 日志级别

database:
  driver: mysql           # 数据库驱动
  host: localhost         # 数据库主机
  port: 3306              # 数据库端口
  user: root              # 数据库用户
  password: password       # 数据库密码
  dbname: learn_hub       # 数据库名称
  max_open_conns: 100     # 最大连接数
  max_idle_conns: 10      # 最大空闲连接数
  conn_max_lifetime: 3600 # 连接最大生命周期（秒）

jwt:
  secret: your-secret-key # JWT 密钥
  expire_hours: 24        # Token 过期时间（小时）
  refresh_expire_hours: 720 # 刷新 Token 过期时间（小时）

oss:
  provider: aliyun        # OSS 提供商：aliyun/tencent
  endpoint: ...           # OSS 端点
  access_key: ...         # 访问密钥
  secret_key: ...         # 密钥
  bucket: learn-hub       # 存储桶名称
  region: cn-hangzhou     # 区域

log:
  level: debug            # 日志级别
  format: json            # 日志格式
  output: stdout          # 输出方式：stdout/file
  file_path: ./logs/app.log # 日志文件路径
```

## API 文档

### 认证接口

#### 登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}
```

#### 注册
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}
```

### 资料接口

#### 获取资料列表
```
GET /api/v1/materials
Authorization: Bearer {token}
```

#### 上传资料
```
POST /api/v1/materials
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

### 考试接口

#### 获取试卷列表
```
GET /api/v1/exams
Authorization: Bearer {token}
```

#### 开始考试
```
POST /api/v1/exams/{id}/start
Authorization: Bearer {token}
```

#### 提交答卷
```
POST /api/v1/exams/{id}/submit
Authorization: Bearer {token}
Content-Type: application/json
```

更多 API 详情，请查看 Swagger 文档。

## 数据库迁移

### 自动迁移

运行迁移工具会自动创建所有表：

```bash
make migrate
```

### 手动迁移

如需手动执行 SQL，可参考 README.md 中的数据库设计部分。

## 开发规范

### 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go)
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查

### 提交规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：feat, fix, docs, style, refactor, test, chore

### 测试要求

- 单元测试覆盖率 ≥ 70%
- 关键业务逻辑必须有测试

## 部署

### Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### 生产环境

1. 修改 `config/config.yaml` 中的配置
2. 设置环境变量 `SERVER_ENV=production`
3. 使用 Docker 或直接运行二进制文件

## 常见问题

### Q: 如何修改数据库连接？
A: 编辑 `config/config.yaml` 中的 database 配置。

### Q: 如何生成 Swagger 文档？
A: 运行 `make swagger` 命令。

### Q: 如何添加新的 API 接口？
A: 
1. 在 `internal/api/handler` 中创建处理器
2. 在 `internal/api/routes.go` 中注册路由
3. 添加 Swagger 注释
4. 运行 `make swagger` 生成文档

### Q: 如何运行测试？
A: 运行 `make test` 命令。

## 许可证

MIT

## 联系方式

如有问题，请提交 Issue。
