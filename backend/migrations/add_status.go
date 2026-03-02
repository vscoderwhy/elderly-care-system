package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// 连接数据库
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=elderly_care sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查列是否存在
	var columnName string
	err = db.QueryRow("SELECT column_name FROM information_schema.columns WHERE table_name='users' AND column_name='status'").Scan(&columnName)

	if err == sql.ErrNoRows {
		// 列不存在，添加列
		fmt.Println("Adding status column to users table...")
		_, err = db.Exec("ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active'")
		if err != nil {
			log.Fatal("Failed to add column:", err)
		}

		// 创建索引
		_, err = db.Exec("CREATE INDEX idx_users_status ON users(status)")
		if err != nil {
			log.Println("Warning: Failed to create index:", err)
		}

		// 更新现有数据
		_, err = db.Exec("UPDATE users SET status = 'active' WHERE status IS NULL")
		if err != nil {
			log.Println("Warning: Failed to update existing data:", err)
		}

		fmt.Println("Status column added successfully!")
	} else if err != nil {
		log.Fatal("Failed to check column:", err)
	} else {
		fmt.Println("Status column already exists!")
	}

	fmt.Println("Migration completed!")
}
