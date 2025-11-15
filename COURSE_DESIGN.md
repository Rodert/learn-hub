# 学习管理系统设计方案

## 一、整体架构

### 1.1 系统角色
- **员工（User）**：学习课程，查看进度
- **管理员（Admin）**：管理课程，查看学习进度

### 1.2 功能模块
```
┌─────────────────────────────────────┐
│    员工端（H5移动端 - React）         │
│  - 课程列表                          │
│  - 课程详情（视频/文本）              │
│  - 学习进度追踪                      │
│  - 标记完成                          │
│  技术栈：React + Ant Design Mobile   │
└─────────────────────────────────────┘
              │
              │ API
              ▼
┌─────────────────────────────────────┐
│       后端服务（Go + GORM）          │
│  - 课程管理 API                      │
│  - 学习记录 API                      │
│  - 进度追踪 API                      │
└─────────────────────────────────────┘
              │
              │ SQL
              ▼
┌─────────────────────────────────────┐
│         MySQL 数据库                 │
│  - 课程表                            │
│  - 学习记录表                        │
└─────────────────────────────────────┘
              ▲
              │ API
              │
┌─────────────────────────────────────┐
│  管理后台（PC端 - Ant Design Pro）   │
│  - 课程管理（CRUD）                  │
│  - 学习进度查看                      │
└─────────────────────────────────────┘
```

### 1.3 技术栈选择

**员工端（H5移动端）：**
- **框架**：React 19 + Umi Max（与现有项目保持一致）
- **UI组件库**：Ant Design Mobile（移动端专用）
- **视频播放**：react-player 或 xgplayer（移动端优化）
- **样式方案**：Less + CSS Modules
- **构建工具**：Umi Max（支持移动端构建）

**管理后台（PC端）：**
- **框架**：React 19 + Umi Max（现有）
- **UI组件库**：Ant Design Pro（现有）

## 二、数据库设计

### 2.1 课程表（sys_course）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| title | string(200) | 课程标题 |
| description | text | 课程描述 |
| cover_image | string(500) | 封面图片URL |
| content_type | int | 内容类型：1-视频，2-文本，3-混合 |
| video_url | string(500) | 视频URL（可选） |
| text_content | text | 文本内容（可选） |
| duration | int | 视频时长（秒），文本为0 |
| status | int | 状态：0-草稿，1-已发布，2-已下架 |
| sort_order | int | 排序 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除 |

**索引：**
- `idx_status` (status)
- `idx_sort_order` (sort_order)

### 2.2 学习记录表（sys_course_record）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID（外键） |
| course_id | uint | 课程ID（外键） |
| progress | int | 学习进度（0-100） |
| duration | int | 已学习时长（秒） |
| is_completed | bool | 是否完成 |
| completed_at | datetime | 完成时间（可选） |
| last_study_at | datetime | 最后学习时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**索引：**
- `idx_user_course` (user_id, course_id) - 唯一索引
- `idx_user_id` (user_id)
- `idx_course_id` (course_id)

### 2.3 数据关系

```
User (1) ────< (N) CourseRecord (N) >─── (1) Course
```

## 三、后端API设计

### 3.1 课程管理API（管理员）

#### 3.1.1 获取课程列表
```
GET /api/course/list
参数：
  - current: 页码（默认1）
  - pageSize: 每页数量（默认20）
  - title: 课程标题（搜索）
  - status: 状态筛选（0-草稿，1-已发布，2-已下架）

响应：
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "课程标题",
      "description": "课程描述",
      "coverImage": "封面URL",
      "contentType": 1,
      "duration": 3600,
      "status": 1,
      "sortOrder": 1,
      "createdAt": "2025-01-01 10:00:00"
    }
  ],
  "total": 100
}
```

#### 3.1.2 创建课程
```
POST /api/course
请求体：
{
  "title": "课程标题",
  "description": "课程描述",
  "coverImage": "封面URL",
  "contentType": 1,  // 1-视频，2-文本，3-混合
  "videoUrl": "视频URL",
  "textContent": "文本内容",
  "duration": 3600,
  "status": 1,
  "sortOrder": 1
}
```

#### 3.1.3 更新课程
```
PUT /api/course/:id
请求体：同创建课程
```

#### 3.1.4 删除课程
```
DELETE /api/course/:id
```

#### 3.1.5 发布/下架课程
```
POST /api/course/:id/publish
请求体：
{
  "status": 1  // 1-发布，2-下架
}
```

### 3.2 认证API（员工端和管理端共用）

#### 3.2.1 用户登录
```
POST /api/login/account
请求体：
{
  "username": "用户名",
  "password": "密码",
  "type": "account",
  "autoLogin": true
}

响应：
{
  "status": "ok",
  "type": "account",
  "currentAuthority": "user",  // 或 "admin"
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 3.2.2 获取当前用户信息
```
GET /api/currentUser
Header: Authorization: Bearer <token>

响应：
{
  "success": true,
  "data": {
    "name": "用户姓名",
    "avatar": "头像URL",
    "userid": "00000001",
    "email": "email@example.com",
    "access": "user",  // 或 "admin"
    ...
  }
}
```

#### 3.2.3 退出登录
```
POST /api/login/outLogin
Header: Authorization: Bearer <token>

响应：
{
  "success": true,
  "data": "ok"
}
```

### 3.3 学习API（员工端，需要认证）

#### 3.3.1 获取可学习课程列表
```
GET /api/learn/courses
参数：
  - current: 页码
  - pageSize: 每页数量
  - title: 搜索标题

响应：
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "课程标题",
      "description": "课程描述",
      "coverImage": "封面URL",
      "contentType": 1,
      "duration": 3600,
      "progress": 50,  // 当前用户的学习进度
      "isCompleted": false,
      "lastStudyAt": "2025-01-01 10:00:00"
    }
  ],
  "total": 100
}
```

#### 3.3.2 获取课程详情
```
GET /api/learn/course/:id

响应：
{
  "success": true,
  "data": {
    "id": 1,
    "title": "课程标题",
    "description": "课程描述",
    "coverImage": "封面URL",
    "contentType": 1,
    "videoUrl": "视频URL",
    "textContent": "文本内容",
    "duration": 3600,
    "progress": 50,
    "duration": 1800,  // 已学习时长
    "isCompleted": false,
    "lastStudyAt": "2025-01-01 10:00:00"
  }
}
```

#### 3.3.3 更新学习进度
```
POST /api/learn/course/:id/progress
请求体：
{
  "progress": 75,  // 0-100
  "duration": 2700  // 已学习时长（秒）
}
```

#### 3.3.4 标记课程完成
```
POST /api/learn/course/:id/complete
请求体：
{
  "progress": 100,
  "duration": 3600
}
```

### 3.4 学习进度查看API（管理员，需要认证）

#### 3.4.1 查看课程学习进度
```
GET /api/admin/course/:id/progress
参数：
  - current: 页码
  - pageSize: 每页数量
  - username: 用户名搜索

响应：
{
  "success": true,
  "data": [
    {
      "userId": 1,
      "username": "user1",
      "name": "用户1",
      "progress": 100,
      "duration": 3600,
      "isCompleted": true,
      "completedAt": "2025-01-01 12:00:00",
      "lastStudyAt": "2025-01-01 12:00:00"
    }
  ],
  "total": 50
}
```

#### 3.4.2 查看用户学习进度
```
GET /api/admin/user/:id/progress
参数：
  - current: 页码
  - pageSize: 每页数量

响应：
{
  "success": true,
  "data": [
    {
      "courseId": 1,
      "courseTitle": "课程标题",
      "progress": 50,
      "duration": 1800,
      "isCompleted": false,
      "lastStudyAt": "2025-01-01 10:00:00"
    }
  ],
  "total": 20
}
```

## 四、前端页面设计

### 4.1 员工端页面（H5移动端）

#### 4.1.1 课程列表页（/learn/courses）
- **功能**：展示所有可学习的课程
- **移动端适配**：
  - 全屏布局，无侧边栏
  - 底部导航栏（首页、我的学习）
  - 顶部搜索栏（固定）
- **组件**（使用 Ant Design Mobile）：
  - `SearchBar` - 搜索框
  - `Card` - 课程卡片
    - 封面图片（16:9比例，圆角）
    - 课程标题（加粗，最多2行）
    - 课程描述（摘要，最多1行，灰色）
    - `Progress` - 进度条（底部）
    - `Tag` - 完成状态标签（右上角）
- **交互**：
  - 下拉刷新
  - 上拉加载更多
  - 点击课程卡片进入详情页
  - 长按显示操作菜单（可选）

#### 4.1.2 课程详情页（/learn/course/:id）
- **功能**：展示课程内容，播放视频或显示文本
- **移动端适配**：
  - 全屏播放（视频）
  - 沉浸式阅读（文本）
  - 底部操作栏（固定）
- **组件**：
  - **视频模式**：
    - `VideoPlayer` - 全屏视频播放器
    - 播放控制栏（播放/暂停、进度、音量、全屏）
    - 课程信息卡片（可收起）
  - **文本模式**：
    - 课程标题（顶部固定）
    - 文本内容区（可滚动）
    - 阅读进度指示器（顶部）
  - **通用组件**：
    - `Progress` - 学习进度条
    - `Button` - "标记完成"按钮（底部固定）
- **交互**：
  - 视频播放时自动更新进度（每10秒）
  - 文本滚动时更新进度（滚动到底部自动完成）
  - 点击"标记完成"按钮完成课程
  - 返回按钮（左上角）

#### 4.1.3 我的学习页（/learn/my-courses）
- **功能**：查看个人学习进度
- **移动端适配**：
  - 顶部统计卡片（横向滚动）
  - 课程列表（下拉刷新）
- **组件**：
  - `Statistic` - 统计卡片
    - 总课程数
    - 已完成数
    - 进行中数
    - 学习时长
  - `Tabs` - 分类标签（全部/进行中/已完成）
  - `Card` - 课程列表（带进度条）
- **交互**：
  - 切换标签筛选
  - 点击课程进入详情页

#### 4.1.4 登录页（/user/login）
- **功能**：员工登录，获取JWT Token
- **移动端适配**：
  - 全屏居中布局
  - 大按钮（易于点击，最小44x44px）
  - 输入框（移动端键盘优化，自动聚焦）
  - 记住密码选项
- **组件**（使用 Ant Design Mobile）：
  - `Form` - 登录表单
  - `Input` - 用户名/密码输入（带图标）
  - `Button` - 登录按钮（全宽，主色调）
  - `Checkbox` - 记住密码
- **交互流程**：
  1. 用户输入用户名和密码
  2. 点击登录按钮
  3. 调用 `/api/login/account` 接口
  4. 成功后保存 Token 到 localStorage
  5. 跳转到课程列表页
  6. 如果未登录访问其他页面，自动跳转到登录页
- **认证机制**：
  - 使用 JWT Token（与管理端共用）
  - Token 存储在 localStorage
  - 请求时自动在 Header 中添加 `Authorization: Bearer <token>`
  - Token 过期后自动跳转登录页

### 4.2 管理后台页面（PC端）

#### 4.2.1 课程管理页（/admin/course）
- **功能**：课程的增删改查
- **组件**：
  - 搜索和筛选
  - 课程列表表格
    - 标题
    - 类型
    - 状态
    - 创建时间
    - 操作（编辑、删除、发布/下架）
  - 新建/编辑表单
    - 标题、描述
    - 封面图片上传
    - 内容类型选择
    - 视频URL或文本内容
    - 状态设置

#### 4.2.2 学习进度查看页（/admin/course/progress）
- **功能**：查看所有课程的学习进度
- **组件**：
  - 课程选择下拉框
  - 学习记录表格
    - 用户名
    - 学习进度
    - 已学习时长
    - 完成状态
    - 最后学习时间

## 五、技术实现要点

### 5.1 认证系统（员工端和管理端共用）

#### 5.1.1 登录流程
```typescript
// 员工端登录
const handleLogin = async (values: LoginParams) => {
  try {
    const response = await login(values);
    if (response.status === 'ok' && response.token) {
      // 保存 token
      localStorage.setItem('token', response.token);
      // 获取用户信息
      const userInfo = await queryCurrentUser();
      // 跳转到课程列表
      history.push('/learn/courses');
    }
  } catch (error) {
    message.error('登录失败，请检查用户名和密码');
  }
};
```

#### 5.1.2 Token 存储和管理
- **存储位置**：localStorage（移动端和PC端都支持）
- **Token 格式**：JWT Token
- **过期处理**：后端返回 401 时，清除 token 并跳转登录页
- **自动刷新**：可选，实现 token 刷新机制

#### 5.1.3 请求拦截器
```typescript
// 自动添加 Authorization Header
requestInterceptors: [
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
],

// 处理 401 错误
responseInterceptors: [
  (response) => {
    if (response.status === 401) {
      localStorage.removeItem('token');
      history.push('/user/login');
    }
    return response;
  },
],
```

#### 5.1.4 路由守卫
- 未登录用户访问受保护页面时，自动跳转到登录页
- 登录成功后，跳转到之前访问的页面（可选）

### 5.2 移动端H5适配

#### 5.2.1 Viewport 设置
```html
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
```

#### 5.2.2 响应式单位
- 使用 `rem` 单位（基于 375px 设计稿）
- 1rem = 37.5px（375px / 10）
- 使用 `postcss-pxtorem` 自动转换

#### 5.2.3 触摸优化
- 按钮最小点击区域：44x44px
- 使用 `Touchable` 组件优化触摸反馈
- 防止误触（防抖处理）

#### 5.2.4 移动端视频播放
```typescript
// 使用原生 video 标签或移动端播放器
import ReactPlayer from 'react-player';

<ReactPlayer
  url={videoUrl}
  width="100%"
  height="auto"
  controls
  playing={isPlaying}
  onProgress={handleProgress}
  config={{
    file: {
      attributes: {
        playsInline: true,  // iOS 内联播放
        'webkit-playsinline': true,
      }
    }
  }}
/>
```

### 5.2 视频播放进度追踪

**方案1：前端定时上报**
```javascript
// 每10秒上报一次进度
setInterval(() => {
  const currentTime = videoPlayer.currentTime;
  const duration = videoPlayer.duration;
  const progress = Math.floor((currentTime / duration) * 100);
  
  updateProgress(courseId, progress, currentTime);
}, 10000);
```

**方案2：视频播放事件触发**
```javascript
videoPlayer.addEventListener('timeupdate', () => {
  // 节流处理，避免频繁请求
  throttle(() => {
    updateProgress(courseId, progress, currentTime);
  }, 5000);
});
```

### 5.3 文本阅读进度追踪

**方案1：滚动位置计算**
```javascript
const handleScroll = () => {
  const scrollTop = window.scrollY;
  const documentHeight = document.documentElement.scrollHeight;
  const windowHeight = window.innerHeight;
  const progress = Math.floor((scrollTop / (documentHeight - windowHeight)) * 100);
  
  updateProgress(courseId, progress);
};
```

**方案2：阅读时间计算**
```javascript
// 根据页面停留时间估算进度
const startTime = Date.now();
setInterval(() => {
  const elapsed = Date.now() - startTime;
  const estimatedProgress = Math.min(100, Math.floor((elapsed / estimatedReadTime) * 100));
  updateProgress(courseId, estimatedProgress);
}, 30000);
```

### 5.4 完成状态判断

**自动完成：**
- 视频：进度 >= 95% 且观看时长 >= 视频时长的90%
- 文本：进度 >= 100%

**手动完成：**
- 用户点击"标记完成"按钮

### 5.5 性能优化

1. **课程列表分页加载**
2. **视频使用CDN加速**
3. **图片懒加载**
4. **学习进度批量更新（防抖）**
5. **缓存课程详情**
6. **路由懒加载**
7. **虚拟滚动（长列表）**

### 5.6 移动端特殊处理

1. **iOS Safari 全屏播放**：需要特殊处理
2. **Android 视频播放**：使用原生 video 标签
3. **下拉刷新**：使用 Ant Design Mobile 的 `PullToRefresh`
4. **上拉加载**：使用 `InfiniteScroll`
5. **返回按钮**：使用 `NavBar` 组件

## 六、扩展功能（后续）

1. **课程分类/标签**
2. **学习计划/任务**
3. **学习统计报表**
4. **课程评论/评分**
5. **学习证书**
6. **推送通知**
7. **离线下载**（视频缓存）

## 七、开发优先级

### Phase 1（MVP）
1. ✅ 课程管理（CRUD）- 管理端
2. ✅ 课程列表展示 - 员工端H5
3. ✅ 课程详情展示 - 员工端H5
4. ✅ 学习进度追踪 - 员工端H5
5. ✅ 标记完成 - 员工端H5

### Phase 2
1. 学习进度查看（管理员）
2. 视频播放进度自动更新
3. 文本阅读进度追踪
4. 移动端优化（下拉刷新、上拉加载）

### Phase 3
1. 学习统计
2. 课程分类
3. 搜索优化
4. 离线功能

## 八、项目结构

```
learn-hub-v2/
├── backend/              # 后端服务（Go）
├── frontend-admin/       # 管理后台（PC端）
└── frontend-mobile/      # 员工端（H5移动端）
    ├── config/
    ├── src/
    │   ├── pages/
    │   │   ├── learn/
    │   │   └── user/
    │   ├── components/
    │   ├── services/
    │   └── utils/
    └── package.json
```
