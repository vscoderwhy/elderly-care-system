package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456"
	hash := "$2a$10$.sHb3kFkv7CCS3vVHsu/2uCfJIftnuc5iRJuhRioylIVIyZ4Mln4y"
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("❌ 密码验证失败:", err)
	} else {
		fmt.Println("✅ 密码验证成功!")
	}
	
	// 生成新的哈希看看是否不同
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("生成哈希失败:", err)
	} else {
		fmt.Println("新哈希:", string(newHash))
		
		// 测试新哈希
		err = bcrypt.CompareHashAndPassword(newHash, []byte(password))
		if err != nil {
			fmt.Println("❌ 新哈希验证失败")
		} else {
			fmt.Println("✅ 新哈希验证成功!")
		}
	}
}
