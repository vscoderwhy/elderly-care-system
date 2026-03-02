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

	var columnName string
	err = db.QueryRow("SELECT column_name FROM information_schema.columns WHERE table_name='users' AND column_name='status'").Scan(&columnName)

	if err == sql.ErrNoRows {
		_, err = db.Exec("ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active'")
		if err != nil {
			log.Fatal("Failed to add column:", err)
		}
		_, err = db.Exec("CREATE INDEX idx_users_status ON users(status)")
		_, err = db.Exec("UPDATE users SET status = 'active'")
		fmt.Println("Status column added!")
	} else {
		fmt.Println("Status column already exists!")
	}
}
