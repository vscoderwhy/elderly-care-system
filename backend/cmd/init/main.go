package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=elderly_care sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	// 检查并添加status字段
	var columnName string
	err = db.QueryRow("SELECT column_name FROM information_schema.columns WHERE table_name='users' AND column_name='status'").Scan(&columnName)

	if err == sql.ErrNoRows {
		fmt.Println("Adding status column to users table...")
		_, err = db.Exec("ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active'")
		if err != nil {
			log.Fatal("Failed to add column:", err)
		}
		_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)")
		fmt.Println("Status column added!")
	} else {
		fmt.Println("Status column exists!")
	}

	// 插入管理员账号
	result, err := db.Exec(`
		INSERT INTO users (id, phone, password, nickname, status) VALUES
		(1, '13800138000', '$2a$10$N9qo8uLOickgx2ZMRZoMye1IK4w7jL5Z2G5v5Y7w8X9z0C1d2E3F4', '系统管理员', 'active')
		ON CONFLICT (phone) DO UPDATE SET nickname = '系统管理员', status = 'active'
	`)
	if err != nil {
		log.Println("Warning:", err)
	} else {
		rows, _ := result.RowsAffected()
		fmt.Printf("Admin user updated: %d rows\n", rows)
	}

	// 关联管理员角色
	_, err = db.Exec(`
		INSERT INTO user_roles (user_id, role_id) VALUES (1, 4)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`)
	if err != nil {
		log.Println("Warning:", err)
	}

	fmt.Println("Database initialized successfully!")
}
