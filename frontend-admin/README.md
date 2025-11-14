# Learn Hub 管理后台

基于 React 18 + Ant Design Pro 的管理后台前端项目。

## 项目结构

```
frontend-admin/
├── src/
│   ├── pages/              # 页面组件
│   │   ├── Login.tsx       # 登录页
│   │   ├── Dashboard.tsx   # 仪表盘
│   │   ├── Materials.tsx   # 资料管理
│   │   ├── Questions.tsx   # 题库管理
│   │   ├── Exams.tsx       # 考试管理
│   │   ├── Users.tsx       # 用户管理
│   │   └── Roles.tsx       # 角色权限
│   ├── layouts/            # 布局组件
│   │   └── Layout.tsx      # 主布局
│   ├── services/           # 服务层
│   │   └── api.ts          # API 客户端
│   ├── App.tsx             # 应用入口
│   ├── main.tsx            # 主文件
│   └── index.css           # 全局样式
├── index.html              # HTML 模板
├── vite.config.ts          # Vite 配置
├── tsconfig.json           # TypeScript 配置
├── package.json            # 项目依赖
└── README.md               # 本文件
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

访问 `http://localhost:3000`

### 生产构建

```bash
npm run build
```

### 预览构建结果

```bash
npm run preview
```

## 技术栈

- **React 18** - UI 框架
- **TypeScript** - 类型检查
- **Ant Design 5** - UI 组件库
- **React Router 6** - 路由管理
- **Zustand** - 状态管理
- **Axios** - HTTP 客户端
- **Vite** - 构建工具
- **Tailwind CSS** - 样式框架

## 功能模块

- ✅ 登录/注册
- ✅ 仪表盘
- ⏳ 资料管理
- ⏳ 题库管理
- ⏳ 考试管理
- ⏳ 用户管理
- ⏳ 角色权限

## 开发规范

- 使用 TypeScript 进行类型检查
- 遵循 ESLint 规则
- 使用 Prettier 格式化代码
- 组件使用函数式组件 + Hooks

## 环境变量

创建 `.env` 文件：

```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 常见问题

### Q: 如何连接后端 API？

A: 在 `vite.config.ts` 中配置代理，或在 `.env` 中设置 API 地址。

### Q: 如何添加新页面？

A: 
1. 在 `src/pages` 中创建新组件
2. 在 `src/layouts/Layout.tsx` 中添加路由
3. 在菜单中添加对应项

## 部署

### Docker 部署

```bash
docker build -t learn-hub-admin:latest .
docker run -p 3000:3000 learn-hub-admin:latest
```

### Nginx 部署

```nginx
server {
    listen 80;
    server_name admin.example.com;

    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8080;
    }
}
```

## 许可证

MIT
