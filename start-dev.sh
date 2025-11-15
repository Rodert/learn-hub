#!/bin/bash

# Learn Hub 开发环境启动脚本

set -e

echo "========================================="
echo "Learn Hub 开发环境启动"
echo "========================================="
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装"
    echo "请安装 Docker Desktop: https://www.docker.com/products/docker-desktop"
    exit 1
fi

echo "✅ Docker 已安装"

# 检查 docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose 未安装"
    echo "请安装 Docker Desktop 或单独安装 docker-compose"
    exit 1
fi

echo "✅ docker-compose 已安装"
echo ""

# 启动 MySQL
echo "🚀 启动 MySQL 容器..."
docker-compose up -d mysql

echo "⏳ 等待 MySQL 启动完成..."
sleep 10

# 检查 MySQL 是否启动成功
if docker exec learn-hub-mysql mysqladmin ping -h localhost -u root -ppassword &> /dev/null; then
    echo "✅ MySQL 已启动成功"
else
    echo "❌ MySQL 启动失败"
    echo "查看日志: docker logs learn-hub-mysql"
    exit 1
fi

echo ""

# 进入后端目录
cd backend

# 执行数据库迁移
echo "🔄 执行数据库迁移..."
go run ./cmd/migrate/main.go

if [ $? -eq 0 ]; then
    echo "✅ 数据库迁移成功"
else
    echo "❌ 数据库迁移失败"
    exit 1
fi

echo ""
echo "========================================="
echo "✅ 开发环境启动完成！"
echo "========================================="
echo ""
echo "📝 后端服务信息:"
echo "  地址: http://localhost:8080"
echo "  Swagger: http://localhost:8080/swagger/index.html"
echo "  默认账户: admin/admin123"
echo ""
echo "🚀 启动后端服务:"
echo "  make dev     (开发模式，支持热重载)"
echo "  make run     (生产模式)"
echo ""
echo "🌐 启动前端管理端:"
echo "  cd ../frontend-admin"
echo "  npm install"
echo "  npm run dev"
echo ""
echo "📊 查看 MySQL 日志:"
echo "  docker logs -f learn-hub-mysql"
echo ""
echo "🛑 停止服务:"
echo "  docker-compose down"
echo ""
