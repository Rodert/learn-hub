# Learn Hub 快速开始指南

## 🚀 项目状态

| 阶段 | 状态 | 进度 |
|------|------|------|
| Phase 1 - 后端核心 | ✅ 完成 | 100% |
| Phase 2 - 管理端 | 🚀 进行中 | 70% |
| Phase 3 - 用户端 | ⏳ 待开始 | 0% |
| Phase 4 - 后端必要功能 | ✅ 完成 | 100% |
| Phase 5 - 测试部署 | ⏳ 待开始 | 0% |

## 📦 安装依赖

### 后端
```bash
cd backend
go mod download
# 新增依赖
go get github.com/xuri/excelize/v2
```

### 前端管理端
```bash
cd frontend-admin
npm install
```

## 🗄️ 数据库初始化

```bash
cd backend

# 执行迁移（创建表和默认数据）
make migrate

# 或直接运行迁移工具
go run cmd/migrate/main.go
```

**默认管理员账户**:
- 用户名: `admin`
- 密码: `admin123`

## 🏃 启动服务

### 后端 API
```bash
cd backend

# 开发模式（需要 air）
make dev

# 或直接运行
make run

# 服务将在 http://localhost:8080 启动
# Swagger 文档: http://localhost:8080/swagger/index.html
```

### 前端管理端
```bash
cd frontend-admin

# 开发模式
npm run dev

# 生产构建
npm run build

# 预览
npm run preview
```

## 📝 API 文档

### 认证接口

#### 登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

#### 刷新 Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "<old-token>"
  }'
```

### 文件上传

#### 上传文件
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/file.pdf" \
  -F "file_type=material"
```

#### 删除文件
```bash
curl -X POST http://localhost:8080/api/v1/files/delete \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "/uploads/1234567890_filename.pdf"
  }'
```

### 数据导入导出

#### 导入题目
```bash
curl -X POST http://localhost:8080/api/v1/import-export/questions \
  -H "Authorization: Bearer <token>" \
  -F "file=@questions.xlsx" \
  -F "exam_id=1"
```

#### 导入用户
```bash
curl -X POST http://localhost:8080/api/v1/import-export/users \
  -H "Authorization: Bearer <token>" \
  -F "file=@users.xlsx" \
  -F "role_id=1"
```

#### 导出考试成绩
```bash
curl -X GET "http://localhost:8080/api/v1/import-export/exam-scores?exam_id=1" \
  -H "Authorization: Bearer <token>" \
  -o exam_scores.xlsx
```

## 📂 项目结构

```
learn-hub/
├── backend/                          # Go 后端
│   ├── cmd/
│   │   ├── main.go                  # 应用入口
│   │   └── migrate/main.go          # 数据库迁移
│   ├── config/                      # 配置
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/             # 请求处理器
│   │   │   └── routes.go            # 路由定义
│   │   ├── middleware/              # 中间件
│   │   ├── model/                   # 数据模型
│   │   ├── repository/              # 数据访问层
│   │   └── service/                 # 业务逻辑
│   ├── pkg/
│   │   ├── database/                # 数据库工具
│   │   ├── oss/                     # 文件存储
│   │   └── excel/                   # Excel 处理
│   └── go.mod
├── frontend-admin/                   # React 管理端
│   ├── src/
│   │   ├── pages/                   # 页面
│   │   ├── components/              # 组件
│   │   ├── services/                # API 服务
│   │   └── App.tsx
│   └── package.json
├── README.md                         # 项目文档
├── TODO.md                           # 任务规划
└── BACKEND_PHASE4_SUMMARY.md        # Phase 4 总结
```

## ⚙️ 配置说明

### 后端配置 (config/config.yaml)

```yaml
server:
  port: 8080
  env: development
  log_level: debug

database:
  driver: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: learn_hub
  max_open_conns: 100
  max_idle_conns: 10

jwt:
  secret: your-secret-key
  expire_hours: 24
  refresh_expire_hours: 720

oss:
  provider: local              # 本地存储（开发环境）
  endpoint: ./uploads
  # 生产环境配置
  # provider: aliyun
  # endpoint: oss-cn-hangzhou.aliyuncs.com
  # access_key: your-access-key
  # secret_key: your-secret-key
  # bucket: learn-hub
  # region: cn-hangzhou
```

## 🔑 核心功能

### ✅ 已实现

#### 后端 (Phase 1 + Phase 4)
- [x] 用户认证 (登录、注册、JWT)
- [x] Token 刷新 (RefreshToken)
- [x] RBAC 权限体系
- [x] 资料管理 (CRUD)
- [x] 题库管理 (CRUD)
- [x] 考试系统 (开始、提交、自动评分)
- [x] 学习记录追踪
- [x] 文件上传 (支持多个 OSS 提供商)
- [x] 数据导入导出 (Excel)
- [x] 默认数据初始化

#### 前端管理端 (Phase 2 - 70%)
- [x] 登录页
- [x] 主布局和菜单
- [x] 资料管理 (列表、创建、编辑、删除)
- [x] 题库管理 (列表、创建、编辑、删除)
- [x] 考试管理 (列表、创建、编辑、删除)
- [x] 用户管理 (列表、创建、编辑、删除)
- [x] 角色权限管理
- [ ] 数据统计
- [ ] 系统设置
- [ ] 单元测试

### ⏳ 待实现

#### 后端
- [ ] 集成阿里云/腾讯云 OSS SDK
- [ ] 单元测试
- [ ] 性能优化 (缓存、查询优化)
- [ ] 安全加固 (速率限制、SQL 注入防护)

#### 前端管理端
- [ ] 搜索功能
- [ ] 文件上传集成
- [ ] 数据统计模块
- [ ] 系统设置模块
- [ ] 单元测试

#### 前端用户端 (Phase 3)
- [ ] 创建 Taro + React 项目
- [ ] 基础框架
- [ ] 学习模块
- [ ] 考试模块
- [ ] 个人中心

## 🧪 测试

### 后端测试
```bash
cd backend

# 运行所有测试
make test

# 生成覆盖率报告
make test-coverage
```

### 前端测试
```bash
cd frontend-admin

# 运行测试
npm run test

# 生成覆盖率报告
npm run test:coverage
```

## 📊 默认数据

### 角色
- `user` - 普通用户
- `admin` - 管理员
- `system_admin` - 系统管理员

### 权限
- `materials:view` - 浏览资料
- `materials:create` - 创建资料
- `materials:update` - 编辑资料
- `materials:delete` - 删除资料
- `exams:view` - 查看试卷
- `exams:submit` - 提交答卷
- `exams:manage` - 管理试卷
- `questions:manage` - 管理题库
- `users:view` - 查看用户
- `users:manage` - 管理用户
- `roles:manage` - 管理角色

### 菜单
- 仪表盘
- 学习资料
- 题库管理
- 考试管理
- 用户管理
- 角色权限

## 🐛 常见问题

### Q: 如何修改数据库连接？
A: 编辑 `backend/config/config.yaml` 中的 database 配置。

### Q: 如何使用阿里云 OSS？
A: 
1. 在 `config/config.yaml` 中配置 OSS
2. 运行 `go get github.com/aliyun/aliyun-oss-go-sdk`
3. 在 `backend/pkg/oss/oss.go` 中实现阿里云 SDK 集成

### Q: 如何导入题目？
A: 
1. 准备 Excel 文件，格式为: 题型 | 题目内容 | 选项 | 答案 | 分数
2. 调用 `/api/v1/import-export/questions` 接口
3. 查看导入结果

### Q: 如何导出考试成绩？
A: 调用 `/api/v1/import-export/exam-scores?exam_id=1` 接口，返回 Excel 文件。

## 📞 相关文档

- [项目 README](./README.md)
- [后端 README](./backend/README.md)
- [前端管理端 README](./frontend-admin/README.md)
- [项目 TODO](./TODO.md)
- [Phase 4 总结](./BACKEND_PHASE4_SUMMARY.md)

## 🔗 有用的链接

- [Swagger API 文档](http://localhost:8080/swagger/index.html) - 启动后端后访问
- [GitHub 项目](https://github.com/Rodert/learn-hub)
- [Go 官方文档](https://golang.org/doc/)
- [React 官方文档](https://react.dev/)
- [Ant Design Pro](https://pro.ant.design/)

---

**最后更新**: 2025-11-14  
**版本**: 1.0.0
