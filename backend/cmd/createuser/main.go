package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

func main() {
	// 连接数据库
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=elderly_care sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	phone := "13800138000"
	password := "123456"
	name := "测试管理员"
	role := "admin"

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("密码加密失败:", err)
	}

	// 检查用户是否存在
	var existingID int
	err = db.QueryRow("SELECT id FROM users WHERE phone = $1", phone).Scan(&existingID)

	if err == sql.ErrNoRows {
		// 用户不存在，创建新用户
		query := `INSERT INTO users (phone, password_hash, name, role, status) 
		          VALUES ($1, $2, $3, $4, 'active') RETURNING id`
		var userID int
		err = db.QueryRow(query, phone, string(hashedPassword), name, role).Scan(&userID)
		if err != nil {
			log.Fatal("创建用户失败:", err)
		}
		fmt.Printf("✅ 用户创建成功！\n")
		fmt.Printf("   ID: %d\n", userID)
		fmt.Printf("   手机号: %s\n", phone)
		fmt.Printf("   密码: %s\n", password)
		fmt.Printf("   姓名: %s\n", name)
		fmt.Printf("   角色: %s\n", role)
	} else if err != nil {
		log.Fatal("查询用户失败:", err)
	} else {
		// 用户存在，更新密码
		query := `UPDATE users SET password_hash = $1, name = $2, role = $3 WHERE phone = $4`
		_, err = db.Exec(query, string(hashedPassword), name, role, phone)
		if err != nil {
			log.Fatal("更新用户失败:", err)
		}
		fmt.Printf("✅ 用户密码已更新！\n")
		fmt.Printf("   ID: %d\n", existingID)
		fmt.Printf("   手机号: %s\n", phone)
		fmt.Printf("   新密码: %s\n", password)
	}
}
