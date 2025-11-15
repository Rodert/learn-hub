#!/bin/bash

echo "=== 测试后端 API 接口 ==="
echo ""

echo "1. 测试登录接口 (admin/admin123):"
TOKEN=$(curl -s -X POST http://localhost:8080/api/login/account \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","type":"account"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "登录失败，返回结果:"
  curl -s -X POST http://localhost:8080/api/login/account \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123","type":"account"}'
  echo ""
  echo ""
  echo "可能原因:"
  echo "  - 数据库未初始化用户"
  echo "  - 密码不正确"
  echo "  - 用户不存在"
else
  echo "登录成功! Token: ${TOKEN:0:50}..."
  echo ""
  echo "2. 测试获取当前用户:"
  curl -s -X GET http://localhost:8080/api/currentUser \
    -H "Authorization: Bearer $TOKEN" | head -20
  echo ""
  echo ""
  echo "3. 测试获取菜单列表:"
  curl -s -X GET http://localhost:8080/api/menu/list \
    -H "Authorization: Bearer $TOKEN" | head -30
fi

echo ""
echo "4. 测试验证码接口:"
curl -s "http://localhost:8080/api/login/captcha?phone=13800138000"
echo ""

