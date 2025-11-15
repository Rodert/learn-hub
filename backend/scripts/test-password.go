package main

import (
	"fmt"
	"learn-hub-backend/utils"
)

func main() {
	// 测试密码验证
	password := "admin123"

	// 这是数据库中应该存储的哈希值（已通过验证）
	correctHash := "$2a$10$gIOMvFijGSZs7IHr/r38l.E0F4gOXWIXqovZESLh9SAmgwQEr9eIW"

	// 测试验证
	result := utils.CheckPassword(password, correctHash)
	fmt.Printf("密码验证结果: %v\n", result)

	if result {
		fmt.Println("✅ 密码验证成功！")
	} else {
		fmt.Println("❌ 密码验证失败！")

		// 生成新的哈希值用于对比
		newHash, err := utils.HashPassword(password)
		if err != nil {
			fmt.Printf("生成新哈希失败: %v\n", err)
			return
		}
		fmt.Printf("新生成的哈希值: %s\n", newHash)
		fmt.Printf("数据库中应该使用的哈希值: %s\n", correctHash)
	}
}
