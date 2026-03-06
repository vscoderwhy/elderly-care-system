#!/bin/bash

# 添加调试日志到 statistics.go

cat > /tmp/test_elderly_list.go << 'EOF'
package main

import (
	"elderly-care-system/internal/config"
	"elderly-care-system/internal/repository"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()
	
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	elderlyRepo := repository.NewElderlyRepository(db)
	
	elderly, total, err := elderlyRepo.List(1, 100)
	if err != nil {
		log.Fatal("Error:", err)
	}
	
	fmt.Printf("=== 总数: %d ===\n", total)
	for _, e := range elderly {
		fmt.Printf("ID: %d, Name: %s, Gender: '%s', CareLevel: %d\n", 
			e.ID, e.Name, e.Gender, e.CareLevel)
	}
	
	genderDist := make(map[string]int)
	for _, e := range elderly {
		genderDist[e.Gender]++
	}
	
	fmt.Printf("\n=== 性别分布 ===\n")
	for k, v := range genderDist {
		fmt.Printf("  %s: %d\n", k, v)
	}
}
EOF

go run /tmp/test_elderly_list.go
