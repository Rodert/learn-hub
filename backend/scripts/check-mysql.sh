#!/bin/bash

# MySQL 连接检查脚本

echo "检查 MySQL 容器状态..."
docker-compose ps mysql

echo ""
echo "检查 MySQL 日志（最后 20 行）..."
docker-compose logs --tail=20 mysql

echo ""
echo "尝试连接 MySQL..."
docker-compose exec -T mysql mysql -u root -proot123456 -e "SELECT 1;" 2>&1

echo ""
echo "检查数据库是否存在..."
docker-compose exec -T mysql mysql -u root -proot123456 -e "SHOW DATABASES;" 2>&1

echo ""
echo "如果看到 learn_hub 数据库，说明 MySQL 已就绪"

