# Learn Hub - 学习系统 MVP

一个完整的在线学习平台，支持资料浏览、题库管理、考试系统和学习进度追踪。

[GitHub](https://github.com/Rodert/learn-hub)

## 📋 目录

- [项目概述](#项目概述)
- [技术栈](#技术栈)
- [系统架构](#系统架构)
- [核心功能](#核心功能)
- [数据结构](#数据结构)
- [API 设计](#api-设计)
- [认证与权限](#认证与权限)
- [部署指南](#部署指南)

---

## 项目概述

Learn Hub 是一个 MVP 级别的在线学习系统，提供完整的学习闭环：

**资料上传 → 用户学习 → 做题 → 自动评分 → 管理员查看**

### 核心特性

- ✅ 跨平台支持（H5 小程序、Web 浏览器）
- ✅ 完整的考试系统（单选、多选、填空）
- ✅ 自动评分机制
- ✅ 学习进度追踪
- ✅ 管理后台（资料、题库、用户数据管理）
- ✅ OSS 文件存储

---

## 技术栈

### 后端
- **语言**：Go 1.20+
- **框架**：Gin / Echo
- **ORM**：GORM
- **数据库**：MySQL 8.0+
- **认证**：JWT
- **文件存储**：阿里云 OSS / 腾讯云 COS

### 前端 - 用户端（H5 小程序）
- **框架**：Taro + React
- **UI 组件**：taro-ui / nutui
- **状态管理**：Redux / Zustand
- **HTTP 客户端**：axios

### 前端 - 管理端（Web）
- **框架**：React 18+
- **UI 框架**：Ant Design Pro
- **状态管理**：Redux / Zustand
- **HTTP 客户端**：axios

### 开发工具
- **容器化**：Docker + docker-compose
- **版本控制**：Git
- **测试**：Go testing + Jest + React Testing Library
- **API 文档**：Swagger / OpenAPI

---

## 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                     用户端 & 管理端                          │
├─────────────────────────────────────────────────────────────┤
│  H5 小程序 (Taro + React)  │  Web 后台 (Ant Design Pro)    │
└────────────────┬──────────────────────────┬─────────────────┘
                 │ HTTPS API                │ HTTPS API
                 └────────────┬─────────────┘
                              │
                    ┌─────────▼──────────┐
                    │   后端服务 (Go)     │
                    ├────────────────────┤
                    │ - API 接口         │
                    │ - JWT 认证         │
                    │ - 业务逻辑         │
                    │ - 文件处理         │
                    └─────────┬──────────┘
                              │
                    ┌─────────▼──────────┐
                    │   MySQL 数据库      │
                    └────────────────────┘
                              │
                    ┌─────────▼──────────┐
                    │   OSS 文件存储      │
                    └────────────────────┘
```

### 核心模块

| 模块 | 职责 | 技术 |
|------|------|------|
| **用户认证** | 账号密码登录、JWT token 管理 | Go + JWT |
| **内容管理** | 资料 CRUD、文件上传 | Go + GORM + OSS |
| **题库系统** | 题目管理、试卷组织 | Go + GORM |
| **考试系统** | 试卷下发、做题、自动评分 | Go + 评分引擎 |
| **学习记录** | 进度追踪、成绩记录 | Go + GORM |
| **管理后台** | 数据管理、统计查看 | React + Ant Design Pro |

---

## 核心功能

### 1. 内容学习模块
- 浏览学习资料列表
- 查看资料详情（图文、视频、文件）
- 标记学习完成状态
- 记录学习时间

### 2. 考试系统
- **题型支持**：单选题、多选题、填空题
- **试卷管理**：创建试卷、组织题目、设置分值
- **做题流程**：
  - 用户开始考试
  - 实时保存答题进度（可选）
  - 提交试卷
  - 自动评分
- **成绩查看**：用户查看个人成绩，管理员查看全部成绩

### 3. 学习记录
- 用户学习进度追踪
- 考试成绩记录
- 学习统计（完成率、平均分等）

### 4. 管理后台
- 资料上传与管理（CRUD）
- 题库管理（题目 CRUD、试卷组织）
- 用户学习数据查看
- 考试成绩统计

---

## 数据结构

### 数据库表设计

#### 用户相关

**users** - 用户表
```sql
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  nickname VARCHAR(100),
  openid VARCHAR(255),
  status ENUM('active', 'inactive', 'banned') DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  INDEX idx_username (username),
  INDEX idx_created_at (created_at)
);
```

**roles** - 角色表
```sql
CREATE TABLE roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) UNIQUE NOT NULL COMMENT '角色名称：user, admin, system_admin',
  description VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  INDEX idx_name (name)
);
```

**user_roles** - 用户角色关联表
```sql
CREATE TABLE user_roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (role_id) REFERENCES roles(id),
  UNIQUE KEY uk_user_role (user_id, role_id),
  INDEX idx_user_id (user_id),
  INDEX idx_role_id (role_id)
);
```

**permissions** - 权限表
```sql
CREATE TABLE permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) UNIQUE NOT NULL COMMENT '权限标识：materials:create, exams:view 等',
  description VARCHAR(255),
  resource VARCHAR(100) COMMENT '资源类型：materials, exams, users 等',
  action VARCHAR(50) COMMENT '操作：create, read, update, delete',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  INDEX idx_name (name),
  INDEX idx_resource (resource)
);
```

**role_permissions** - 角色权限关联表
```sql
CREATE TABLE role_permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (permission_id) REFERENCES permissions(id),
  UNIQUE KEY uk_role_permission (role_id, permission_id),
  INDEX idx_role_id (role_id),
  INDEX idx_permission_id (permission_id)
);
```

**menus** - 菜单表
```sql
CREATE TABLE menus (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '菜单名称',
  path VARCHAR(255) COMMENT '路由路径',
  icon VARCHAR(100) COMMENT '菜单图标',
  component VARCHAR(255) COMMENT '组件路径',
  parent_id BIGINT COMMENT '父菜单 ID',
  order_num INT DEFAULT 0 COMMENT '排序号',
  visible TINYINT DEFAULT 1 COMMENT '是否显示',
  type ENUM('menu', 'button') DEFAULT 'menu' COMMENT '菜单类型',
  permission VARCHAR(100) COMMENT '关联权限标识',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (parent_id) REFERENCES menus(id),
  INDEX idx_parent_id (parent_id),
  INDEX idx_order_num (order_num),
  INDEX idx_permission (permission)
);
```

**role_menus** - 角色菜单关联表
```sql
CREATE TABLE role_menus (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  role_id BIGINT NOT NULL,
  menu_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (menu_id) REFERENCES menus(id),
  UNIQUE KEY uk_role_menu (role_id, menu_id),
  INDEX idx_role_id (role_id),
  INDEX idx_menu_id (menu_id)
);
```

#### 内容相关

**materials** - 学习资料表
```sql
CREATE TABLE materials (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  content_type ENUM('text', 'video', 'file', 'mixed') DEFAULT 'text',
  content TEXT,
  file_url VARCHAR(500),
  file_size BIGINT,
  cover_url VARCHAR(500),
  order_num INT DEFAULT 0,
  status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
  created_by BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (created_by) REFERENCES users(id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at),
  INDEX idx_order_num (order_num)
);
```

#### 题库相关

**questions** - 题库表
```sql
CREATE TABLE questions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  exam_id BIGINT,
  question_type ENUM('single_choice', 'multiple_choice', 'fill_blank') NOT NULL,
  content TEXT NOT NULL,
  options JSON,
  answer VARCHAR(500) NOT NULL,
  explanation TEXT,
  score DECIMAL(10, 2) DEFAULT 1.00,
  order_num INT DEFAULT 0,
  created_by BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (exam_id) REFERENCES exams(id),
  FOREIGN KEY (created_by) REFERENCES users(id),
  INDEX idx_exam_id (exam_id),
  INDEX idx_question_type (question_type),
  INDEX idx_created_at (created_at)
);
```

**exams** - 试卷表
```sql
CREATE TABLE exams (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  total_score DECIMAL(10, 2) DEFAULT 100.00,
  pass_score DECIMAL(10, 2) DEFAULT 60.00,
  time_limit INT COMMENT '考试时间限制（分钟）',
  status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
  created_by BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (created_by) REFERENCES users(id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at)
);
```

#### 考试记录相关

**exam_records** - 用户考试记录表
```sql
CREATE TABLE exam_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  exam_id BIGINT NOT NULL,
  score DECIMAL(10, 2),
  status ENUM('in_progress', 'submitted', 'graded') DEFAULT 'in_progress',
  answers JSON COMMENT '用户答题记录',
  start_time TIMESTAMP,
  submit_time TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (exam_id) REFERENCES exams(id),
  INDEX idx_user_id (user_id),
  INDEX idx_exam_id (exam_id),
  INDEX idx_user_exam (user_id, exam_id),
  INDEX idx_created_at (created_at)
);
```

#### 学习记录相关

**course_records** - 用户学习记录表
```sql
CREATE TABLE course_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  material_id BIGINT NOT NULL,
  status ENUM('not_started', 'in_progress', 'completed') DEFAULT 'not_started',
  progress_percent INT DEFAULT 0 COMMENT '学习进度百分比',
  view_duration INT COMMENT '浏览时长（秒）',
  completed_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (material_id) REFERENCES materials(id),
  INDEX idx_user_id (user_id),
  INDEX idx_material_id (material_id),
  INDEX idx_user_material (user_id, material_id),
  INDEX idx_created_at (created_at)
);
```

#### 专题相关（可选）

**topics** - 学习专题表
```sql
CREATE TABLE topics (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  order_num INT DEFAULT 0,
  status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
  created_by BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (created_by) REFERENCES users(id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at)
);
```

**topic_materials** - 专题资料关联表
```sql
CREATE TABLE topic_materials (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  topic_id BIGINT NOT NULL,
  material_id BIGINT NOT NULL,
  order_num INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (topic_id) REFERENCES topics(id),
  FOREIGN KEY (material_id) REFERENCES materials(id),
  UNIQUE KEY uk_topic_material (topic_id, material_id),
  INDEX idx_topic_id (topic_id)
);
```

**topic_exams** - 专题考试关联表
```sql
CREATE TABLE topic_exams (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  topic_id BIGINT NOT NULL,
  exam_id BIGINT NOT NULL,
  order_num INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (topic_id) REFERENCES topics(id),
  FOREIGN KEY (exam_id) REFERENCES exams(id),
  UNIQUE KEY uk_topic_exam (topic_id, exam_id),
  INDEX idx_topic_id (topic_id)
);
```

---

## API 设计

### 认证接口

#### 登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "id": 1,
      "username": "user@example.com",
      "nickname": "张三",
      "role": "user"
    }
  }
}
```

#### 刷新 Token
```
POST /api/v1/auth/refresh
Authorization: Bearer {token}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGc..."
  }
}
```

### 内容管理接口

#### 获取资料列表
```
GET /api/v1/materials?page=1&limit=10&status=published

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Go 基础教程",
        "description": "...",
        "content_type": "video",
        "cover_url": "https://oss.example.com/...",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 10
  }
}
```

#### 获取资料详情
```
GET /api/v1/materials/{id}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Go 基础教程",
    "content": "...",
    "file_url": "https://oss.example.com/...",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 上传资料（管理员）
```
POST /api/v1/admin/materials
Authorization: Bearer {admin_token}
Content-Type: multipart/form-data

{
  "title": "Go 基础教程",
  "description": "...",
  "content_type": "video",
  "file": <binary>
}

Response 201:
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "file_url": "https://oss.example.com/..."
  }
}
```

### 考试接口

#### 获取试卷列表
```
GET /api/v1/exams?page=1&limit=10

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Go 基础测试",
        "total_score": 100,
        "pass_score": 60,
        "time_limit": 60
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 10
  }
}
```

#### 开始考试
```
POST /api/v1/exams/{id}/start
Authorization: Bearer {user_token}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "exam_record_id": 123,
    "exam": {
      "id": 1,
      "title": "Go 基础测试",
      "time_limit": 60
    },
    "questions": [
      {
        "id": 1,
        "question_type": "single_choice",
        "content": "Go 的并发模型是什么？",
        "options": ["A. 进程", "B. 线程", "C. Goroutine", "D. 协程"],
        "score": 5
      }
    ]
  }
}
```

#### 提交答卷
```
POST /api/v1/exams/{id}/submit
Authorization: Bearer {user_token}
Content-Type: application/json

{
  "exam_record_id": 123,
  "answers": [
    {
      "question_id": 1,
      "answer": "C"
    },
    {
      "question_id": 2,
      "answer": "答案内容"
    }
  ]
}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "score": 85,
    "pass": true,
    "details": [
      {
        "question_id": 1,
        "user_answer": "C",
        "correct_answer": "C",
        "is_correct": true,
        "score": 5
      }
    ]
  }
}
```

#### 获取考试成绩
```
GET /api/v1/exams/{id}/records
Authorization: Bearer {user_token}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 123,
        "exam_id": 1,
        "score": 85,
        "status": "graded",
        "submit_time": "2024-01-01T10:30:00Z"
      }
    ]
  }
}
```

### 学习记录接口

#### 获取学习进度
```
GET /api/v1/course-records?material_id={id}
Authorization: Bearer {user_token}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "material_id": 1,
    "status": "in_progress",
    "progress_percent": 50,
    "view_duration": 1800,
    "completed_at": null
  }
}
```

#### 更新学习进度
```
PUT /api/v1/course-records/{id}
Authorization: Bearer {user_token}
Content-Type: application/json

{
  "status": "completed",
  "progress_percent": 100,
  "view_duration": 3600
}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "status": "completed",
    "completed_at": "2024-01-01T11:00:00Z"
  }
}
```

---

## 认证与权限

### JWT Token 结构

```json
{
  "sub": "user_id",
  "username": "user@example.com",
  "roles": ["user", "admin"],
  "permissions": ["materials:view", "exams:view", "exams:submit"],
  "iat": 1704110400,
  "exp": 1704196800
}
```

### RBAC 权限体系

#### 角色定义

| 角色 | 描述 |
|------|------|
| **user** | 普通用户，可浏览资料、做题、查看个人成绩 |
| **admin** | 管理员，可管理资料、题库、查看用户数据 |
| **system_admin** | 系统管理员，可管理所有内容和用户 |

#### 权限定义

| 权限标识 | 资源 | 操作 | 描述 |
|---------|------|------|------|
| `materials:view` | materials | read | 浏览学习资料 |
| `materials:create` | materials | create | 上传学习资料 |
| `materials:update` | materials | update | 编辑学习资料 |
| `materials:delete` | materials | delete | 删除学习资料 |
| `exams:view` | exams | read | 查看试卷 |
| `exams:submit` | exams | create | 提交答卷 |
| `exams:manage` | exams | update | 管理试卷 |
| `questions:manage` | questions | update | 管理题库 |
| `users:view` | users | read | 查看用户数据 |
| `users:manage` | users | update | 管理用户 |
| `roles:manage` | roles | update | 管理角色 |

#### 角色权限映射

**user 角色权限：**
```
- materials:view
- exams:view
- exams:submit
```

**admin 角色权限：**
```
- materials:view
- materials:create
- materials:update
- materials:delete
- exams:view
- exams:manage
- questions:manage
- users:view
```

**system_admin 角色权限：**
```
- 所有权限
```

### 权限验证中间件

后端在每个受保护的接口前添加权限验证中间件：

```go
// 验证 JWT token
AuthMiddleware()

// 验证权限
RequirePermission("materials:create")

// 验证数据所有权
CheckOwnership(resourceId, userId)
```

### 菜单配置示例

**admin 角色菜单：**
```json
[
  {
    "id": 1,
    "name": "学习资料",
    "path": "/materials",
    "icon": "FileText",
    "permission": "materials:view",
    "children": [
      {
        "id": 2,
        "name": "上传资料",
        "path": "/materials/create",
        "type": "button",
        "permission": "materials:create"
      }
    ]
  },
  {
    "id": 3,
    "name": "题库管理",
    "path": "/questions",
    "icon": "FileText",
    "permission": "questions:manage"
  },
  {
    "id": 4,
    "name": "用户管理",
    "path": "/users",
    "icon": "Users",
    "permission": "users:view"
  }
]
```

---

## 部署指南

### 本地开发环境

#### 前置条件
- Docker & Docker Compose
- Go 1.20+
- Node.js 16+
- MySQL 8.0+

#### 快速启动

```bash
# 克隆项目
git clone https://github.com/your-org/learn-hub.git
cd learn-hub

# 启动 Docker 容器（MySQL）
docker-compose up -d

# 后端开发
cd backend
go mod download
go run main.go

# 前端开发 - 管理端
cd frontend-admin
npm install
npm start

# 前端开发 - 用户端
cd frontend-user
npm install
npm run dev
```

### 生产部署

#### Docker 部署

```bash
# 构建后端镜像
docker build -t learn-hub-backend:latest ./backend

# 构建管理端镜像
docker build -t learn-hub-admin:latest ./frontend-admin

# 构建用户端镜像
docker build -t learn-hub-user:latest ./frontend-user

# 使用 docker-compose 启动
docker-compose -f docker-compose.prod.yml up -d
```

#### 环境变量配置

后端 `.env` 文件：
```
# 数据库
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=learn_hub

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE=24h

# OSS
OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
OSS_ACCESS_KEY=your-access-key
OSS_SECRET_KEY=your-secret-key
OSS_BUCKET=learn-hub

# 服务
SERVER_PORT=8080
SERVER_ENV=production
```

### 数据库初始化

```bash
# 运行迁移脚本
go run cmd/migrate/main.go

# 或手动执行 SQL
mysql -u root -p learn_hub < schema.sql
```

---

## 项目结构

```
learn-hub/
├── backend/                    # Go 后端
│   ├── cmd/
│   │   ├── main.go
│   │   └── migrate/
│   ├── internal/
│   │   ├── api/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   └── middleware/
│   ├── config/
│   ├── go.mod
│   └── Dockerfile
├── frontend-admin/             # React + Ant Design Pro 管理端
│   ├── src/
│   │   ├── pages/
│   │   ├── components/
│   │   ├── services/
│   │   └── App.tsx
│   ├── package.json
│   └── Dockerfile
├── frontend-user/              # Taro + React 用户端
│   ├── src/
│   │   ├── pages/
│   │   ├── components/
│   │   └── app.tsx
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml
├── docker-compose.prod.yml
├── README.md
└── DESIGN.md
```

---

## 开发规范

### 代码风格

- **Go**：遵循 [Effective Go](https://golang.org/doc/effective_go)
- **JavaScript/TypeScript**：遵循 [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- **SQL**：使用小写关键字，表名使用 snake_case

### Git 提交规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：feat, fix, docs, style, refactor, test, chore

示例：
```
feat(exam): add auto-grading for fill-blank questions

- Implement fuzzy matching for answers
- Support case-insensitive comparison
- Add explanation display

Closes #123
```

### 测试要求

- 后端：单元测试覆盖率 ≥ 70%
- 前端：关键业务逻辑测试覆盖率 ≥ 60%
- 集成测试：关键流程（登录、做题、提交）

---

## 常见问题

### Q: 如何处理大文件上传？
A: 使用分片上传，前端分片 + 后端合并，配合 OSS 的断点续传功能。

### Q: 如何防止考试作弊？
A: 
- 前端：禁用复制粘贴、截屏
- 后端：记录 IP、设备指纹、异常答题速度检测
- 后续可添加人脸识别验证

### Q: 如何处理考试超时？
A: 后端在提交时校验 `submit_time - start_time` 是否超过 `time_limit`，超时则拒绝提交。

### Q: 支持离线做题吗？
A: MVP 阶段不支持，后续可通过 Service Worker + IndexedDB 实现离线功能。

---

## 许可证

MIT

---

## 联系方式

如有问题，请提交 Issue 或联系开发团队。


