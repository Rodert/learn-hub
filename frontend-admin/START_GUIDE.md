# 前端启动指南

## 正确的启动命令

### 方式1：开发环境启动（推荐）
```bash
cd frontend-admin
npm run start:dev
# 或
npm run dev
```

这会：
- 设置 `REACT_APP_ENV=dev`
- 使用 `proxy.dev` 配置
- 代理 `/api/*` 到 `http://localhost:8080`

### 方式2：普通启动
```bash
cd frontend-admin
npm start
```

这会：
- 设置 `UMI_ENV=dev`
- 默认使用 `proxy.dev` 配置（如果 `REACT_APP_ENV` 未设置）

## 检查代理是否生效

启动后，前端默认运行在 `http://localhost:8000`（或控制台显示的端口）。

### 验证代理配置

1. **检查启动日志**：应该看到类似：
   ```
   App running at:
   - Local:   http://localhost:8000
   ```

2. **检查代理**：访问 `http://localhost:8000/api/login/captcha?phone=13800138000`
   - 如果返回验证码，说明代理正常
   - 如果 404 或连接失败，说明代理未生效

## 常见问题

### 问题1：访问 8003 端口登录失败

**原因**：可能使用了错误的启动命令，或者代理配置未生效。

**解决**：
1. 停止当前前端服务
2. 使用正确的启动命令：`npm run start:dev`
3. 访问控制台显示的端口（通常是 8000）

### 问题2：代理不生效

**检查**：
1. 确认 `config/proxy.ts` 中有 `dev` 配置
2. 确认启动时设置了 `REACT_APP_ENV=dev`
3. 检查后端服务是否在 `http://localhost:8080` 运行

### 问题3：端口被占用

如果 8000 端口被占用，Umi 会自动使用其他端口（如 8001, 8002, 8003...）

**解决**：
- 查看控制台输出的实际端口
- 或者指定端口：`PORT=8000 npm run start:dev`

## 完整启动流程

```bash
# 1. 确保后端服务运行
cd backend
go run main.go
# 后端应该在 http://localhost:8080 运行

# 2. 启动前端（新终端）
cd frontend-admin
npm run start:dev
# 前端应该在 http://localhost:8000 运行

# 3. 访问前端
# 浏览器打开 http://localhost:8000
```

## 环境变量说明

- `REACT_APP_ENV=dev`：使用开发环境配置，启用 `proxy.dev`
- `MOCK=none`：禁用 Mock 数据，使用真实后端
- `UMI_ENV=dev`：Umi 开发模式

