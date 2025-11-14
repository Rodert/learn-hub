# Learn Hub API 文档

## 基础信息

- **API 基础 URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token (JWT)
- **请求头**: `Authorization: Bearer {token}`

## 认证接口

### 1. 用户登录
```
POST /auth/login
```

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "string",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "status": "active"
    }
  }
}
```

### 2. 用户注册
```
POST /auth/register
```

**请求体**:
```json
{
  "username": "string",
  "password": "string",
  "nickname": "string"
}
```

### 3. 刷新 Token
```
POST /auth/refresh
```

---

## 用户接口

### 1. 获取用户信息
```
GET /users/{id}
```

**参数**: `id` (path) - 用户 ID

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "status": "active",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### 2. 更新用户信息
```
PUT /users/{id}
```

**请求体**:
```json
{
  "nickname": "string"
}
```

### 3. 获取当前用户信息
```
GET /users/profile/me
```

---

## 资料管理接口

### 1. 获取资料列表
```
GET /materials?page=1&limit=10&status=published
```

**查询参数**:
- `page` (int) - 页码，默认 1
- `limit` (int) - 每页数量，默认 10
- `status` (string) - 状态过滤（draft/published/archived）

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "资料标题",
        "description": "资料描述",
        "content_type": "text",
        "status": "published",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 10
  }
}
```

### 2. 获取资料详情
```
GET /materials/{id}
```

### 3. 创建资料
```
POST /materials
```

**请求体**:
```json
{
  "title": "string",
  "description": "string",
  "content_type": "text|video|file|mixed",
  "content": "string",
  "file_url": "string (可选)",
  "cover_url": "string (可选)"
}
```

### 4. 更新资料
```
PUT /materials/{id}
```

**请求体**:
```json
{
  "title": "string",
  "description": "string",
  "status": "draft|published|archived"
}
```

### 5. 删除资料
```
DELETE /materials/{id}
```

---

## 题库管理接口

### 1. 获取题目列表
```
GET /questions?page=1&limit=10&type=single_choice
```

**查询参数**:
- `page` (int) - 页码
- `limit` (int) - 每页数量
- `type` (string) - 题型过滤（single_choice/multiple_choice/fill_blank）

### 2. 获取题目详情
```
GET /questions/{id}
```

### 3. 创建题目
```
POST /questions
```

**请求体**:
```json
{
  "exam_id": 1,
  "question_type": "single_choice|multiple_choice|fill_blank",
  "content": "string",
  "answer": "string",
  "explanation": "string (可选)",
  "score": 10
}
```

### 4. 更新题目
```
PUT /questions/{id}
```

### 5. 删除题目
```
DELETE /questions/{id}
```

---

## 考试管理接口

### 1. 获取试卷列表
```
GET /exams?page=1&limit=10
```

### 2. 获取试卷详情
```
GET /exams/{id}
```

### 3. 创建试卷
```
POST /exams
```

**请求体**:
```json
{
  "title": "string",
  "description": "string",
  "total_score": 100,
  "pass_score": 60,
  "time_limit": 120
}
```

### 4. 更新试卷
```
PUT /exams/{id}
```

### 5. 删除试卷
```
DELETE /exams/{id}
```

### 6. 开始考试
```
POST /exams/{id}/start
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "exam_record_id": 1,
    "exam": { ... },
    "questions": [ ... ]
  }
}
```

### 7. 提交答卷
```
POST /exams/{id}/submit
```

**请求体**:
```json
{
  "exam_record_id": 1,
  "answers": [
    {
      "question_id": 1,
      "answer": "A"
    }
  ]
}
```

### 8. 获取考试成绩
```
GET /exams/{id}/records?page=1&limit=10
```

---

## 用户管理接口（管理员）

### 1. 获取用户列表
```
GET /admin/users?page=1&limit=10
```

### 2. 获取用户详情
```
GET /admin/users/{id}
```

### 3. 创建用户
```
POST /admin/users
```

**请求体**:
```json
{
  "username": "string",
  "password": "string",
  "nickname": "string"
}
```

### 4. 更新用户
```
PUT /admin/users/{id}
```

**请求体**:
```json
{
  "nickname": "string",
  "status": "active|inactive|banned"
}
```

### 5. 删除用户
```
DELETE /admin/users/{id}
```

---

## 角色管理接口（管理员）

### 1. 获取角色列表
```
GET /admin/roles
```

### 2. 创建角色
```
POST /admin/roles
```

**请求体**:
```json
{
  "name": "string",
  "description": "string"
}
```

### 3. 更新角色
```
PUT /admin/roles/{id}
```

### 4. 删除角色
```
DELETE /admin/roles/{id}
```

---

## 权限管理接口（管理员）

### 1. 获取权限列表
```
GET /admin/permissions
```

### 2. 创建权限
```
POST /admin/permissions
```

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "resource": "string",
  "action": "string"
}
```

---

## 菜单接口

### 1. 获取用户菜单
```
GET /menus
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "仪表盘",
      "path": "/dashboard",
      "icon": "dashboard",
      "component": "Dashboard",
      "order_num": 1,
      "children": []
    }
  ]
}
```

---

## 学习记录接口

### 1. 获取学习记录
```
GET /course-records?page=1&limit=10
```

### 2. 获取记录详情
```
GET /course-records/{id}
```

### 3. 更新学习进度
```
PUT /course-records/{id}
```

**请求体**:
```json
{
  "progress_percent": 50,
  "view_duration": 300
}
```

---

## 错误响应

所有错误响应格式如下：

```json
{
  "code": 1,
  "message": "error message",
  "data": null
}
```

### 常见错误码

- `401` - 未授权（Token 过期或无效）
- `403` - 禁止访问（权限不足）
- `404` - 资源不存在
- `500` - 服务器错误

---

## 认证流程

1. 调用 `/auth/login` 获取 Token
2. 将 Token 保存到 localStorage
3. 在所有请求的 Authorization 头中添加 Token
4. 如果收到 401 错误，调用 `/auth/refresh` 刷新 Token
5. 如果刷新失败，重定向到登录页面

---

## 示例代码

### JavaScript/TypeScript (Axios)

```typescript
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
```

### 使用示例

```typescript
// 登录
const login = async (username: string, password: string) => {
  const response = await api.post('/auth/login', { username, password })
  localStorage.setItem('token', response.data.data.token)
  return response.data.data
}

// 获取资料列表
const getMaterials = async (page = 1) => {
  const response = await api.get('/materials', {
    params: { page, limit: 10 },
  })
  return response.data.data
}

// 创建资料
const createMaterial = async (data: any) => {
  const response = await api.post('/materials', data)
  return response.data.data
}
```

---

## Swagger UI

访问以下地址查看完整的 API 文档：

```
http://localhost:8080/swagger/index.html
```

---

**最后更新**: 2025-11-14  
**API 版本**: v1
