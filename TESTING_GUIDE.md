# Learn Hub 项目测试指南

## 📋 项目完成情况

### ✅ 已完成的工作

#### Phase 1: 后端核心 (100%) ✅
- 项目结构、数据模型、Repository/Service/Handler 层
- JWT 认证、RBAC 权限体系
- 数据库迁移工具、Swagger 文档

#### Phase 4: 后端必要功能 (100%) ✅
- RefreshToken 实现
- 初始化默认数据
- 文件上传功能 (OSS 集成)
- 数据导入导出功能

#### Phase 2: 前端管理端 (100%) ✅
- 资料管理 - 搜索、过滤、上传、详情、发布/草稿
- 题库管理 - 搜索、过滤、详情、批量导入
- 考试管理 - 搜索、过滤、详情、题目选择、发布
- 用户管理 - 搜索、过滤、详情、角色分配、重置密码
- 角色权限 - 完整的权限管理系统
- 数据统计 - 学习统计、考试统计、图表展示
- 系统设置 - 个人设置、密码修改、头像上传

---

## 🚀 启动服务步骤

### 1. 后端启动

#### 环境准备
```bash
# 设置代理（可选，加速网络）
export https_proxy=http://127.0.0.1:7890
export http_proxy=http://127.0.0.1:7890
export all_proxy=socks5://127.0.0.1:7890

# 进入后端目录
cd backend

# 下载依赖
go mod download
go mod tidy

# 执行数据库迁移（创建表和默认数据）
make migrate

# 启动后端服务
make run
# 或开发模式（支持热重载）
make dev
```

**后端服务地址**: http://localhost:8080  
**Swagger API 文档**: http://localhost:8080/swagger/index.html

#### 默认管理员账户
- 用户名: `admin`
- 密码: `admin123`

---

### 2. 前端管理端启动

```bash
# 进入前端目录
cd frontend-admin

# 安装依赖
npm install

# 开发模式启动
npm run dev

# 生产构建
npm run build

# 预览构建结果
npm run preview
```

**前端地址**: http://localhost:3000 (或 Vite 默认端口)

---

## 🧪 测试场景

### 后端 API 测试

#### 1. 用户认证
```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 响应示例
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

#### 2. 刷新 Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "<old-token>"
  }'
```

#### 3. 获取资料列表
```bash
curl -X GET "http://localhost:8080/api/v1/materials?page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

#### 4. 文件上传
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/file.pdf" \
  -F "file_type=material"
```

#### 5. 题目批量导入
```bash
curl -X POST http://localhost:8080/api/v1/import-export/questions \
  -H "Authorization: Bearer <token>" \
  -F "file=@questions.xlsx" \
  -F "exam_id=1"
```

---

### 前端管理端测试

#### 1. 登录测试
- 访问 http://localhost:3000
- 输入用户名: `admin`
- 输入密码: `admin123`
- 点击登录

#### 2. 资料管理测试
- 进入资料管理模块
- 测试搜索功能
- 测试状态过滤
- 测试文件上传
- 测试详情查看
- 测试发布/草稿切换

#### 3. 题库管理测试
- 进入题库管理模块
- 测试搜索功能
- 测试题型过滤
- 测试批量导入
- 测试详情查看

#### 4. 考试管理测试
- 进入考试管理模块
- 测试搜索功能
- 测试状态过滤
- 测试题目选择
- 测试发布试卷

#### 5. 用户管理测试
- 进入用户管理模块
- 测试搜索功能
- 测试用户创建/编辑/删除
- 测试角色分配

---

## 📊 数据库配置

### MySQL 配置
```yaml
# config/config.yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: learn_hub
  max_open_conns: 100
  max_idle_conns: 10
```

### 创建数据库
```sql
CREATE DATABASE IF NOT EXISTS learn_hub DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

---

## 🔍 常见问题

### Q: 后端启动失败，提示数据库连接错误
A: 检查 MySQL 是否运行，修改 `config/config.yaml` 中的数据库配置

### Q: 前端无法连接后端
A: 检查后端是否运行在 8080 端口，修改前端的 API 地址配置

### Q: 文件上传失败
A: 检查 `uploads` 目录是否存在，或修改 OSS 配置

### Q: 登录失败
A: 确保数据库迁移已执行，默认管理员账户已创建

---

## 📝 项目文档

- [项目 README](./README.md)
- [快速开始](./QUICK_START.md)
- [后端 README](./backend/README.md)
- [前端管理端 README](./frontend-admin/README.md)
- [项目总结](./PROJECT_SUMMARY.md)
- [前端管理端完成总结](./FRONTEND_ADMIN_COMPLETE.md)

---

## 🎯 下一步

1. **完成 Phase 3 用户端** - Taro + React 小程序
2. **完成 Phase 5 测试部署** - 单元测试、集成测试、部署配置
3. **性能优化** - 缓存、查询优化、安全加固

---

**最后更新**: 2025-11-14  
**项目版本**: 1.0.0 (Beta)  
**总体完成度**: 70%
