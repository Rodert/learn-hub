package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	
	// 生成新的密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成失败: %v\n", err)
		return
	}
	
	hashStr := string(hash)
	fmt.Printf("密码: %s\n", password)
	fmt.Printf("生成的哈希值: %s\n", hashStr)
	fmt.Printf("哈希长度: %d\n", len(hashStr))
	
	// 验证生成的哈希
	err = bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(password))
	if err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
	} else {
		fmt.Printf("✅ 验证成功！\n")
	}
	
	// 输出 SQL 更新语句
	fmt.Printf("\n=== SQL 更新语句 ===\n")
	fmt.Printf("UPDATE sys_user SET password = '%s', status = 1 WHERE username = 'admin';\n", hashStr)
}

