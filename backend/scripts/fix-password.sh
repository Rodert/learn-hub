#!/bin/bash

echo "=== 修复 admin 用户密码 ==="
echo ""

# 生成新的密码哈希（使用 Go 程序）
cat > /tmp/gen-password.go << 'EOF'
package main
import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)
func main() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
EOF

echo "生成新的密码哈希..."
NEW_HASH=$(cd /tmp && go run gen-password.go 2>/dev/null)

if [ -z "$NEW_HASH" ]; then
    echo "❌ 无法生成密码哈希，请手动更新"
    echo "   正确的密码哈希（admin123）:"
    echo "   \$2a\$10\$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iwK8pJ5C"
    exit 1
fi

echo "✅ 新密码哈希: ${NEW_HASH:0:30}..."
echo ""

echo "更新数据库中的密码..."
docker compose exec -T mysql mysql -u root -proot123456 learn_hub << EOF
UPDATE sys_user 
SET password = '$NEW_HASH' 
WHERE username = 'admin';
SELECT '密码更新成功' as result;
EOF

echo ""
echo "✅ 密码已更新！现在可以使用 admin/admin123 登录"

