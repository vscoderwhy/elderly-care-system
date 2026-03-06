#!/bin/bash

# 简单的 Go 程序来检查性别分布问题

cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend

cat > /tmp/check_gender.go << 'EOF'
package main

import (
	"database/sql"
	"encoding/json"
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

	rows, err := db.Query("SELECT id, name, gender FROM elderly")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	genderDist := make(map[string]int)
	
	fmt.Println("=== 老人列表 ===")
	for rows.Next() {
		var id int
		var name, gender string
		if err := rows.Scan(&id, &name, &gender); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Gender: '%s' (len=%d)\n", id, name, gender, len(gender))
		genderDist[gender]++
	}
	
	fmt.Println("\n=== 性别分布 ===")
	jsonData, _ := json.MarshalIndent(genderDist, "", "  ")
	fmt.Println(string(jsonData))
}
EOF

go run /tmp/check_gender.go
