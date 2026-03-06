package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=elderly_care sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 测试生成100条
	genders := []string{"男", "女"}
	careLevels := []int{1, 2, 3, 4, 5}
	statuses := []string{"active", "hospitalized", "left"}
	
	success := 0
	for i := 0; i < 100; i++ {
		tx, _ := db.Begin()
		
		gender := genders[rand.Intn(len(genders))]
		
		_, err := tx.Exec(`INSERT INTO elderly (name, gender, birth_date, id_card, phone, address, emergency_contact, emergency_phone, admission_date, care_level, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			"测试老人"+fmt.Sprint(i),
			gender,
			time.Now().AddDate(-70, 0, 0),
			fmt.Sprintf("%d", rand.Intn(1000000000)),
			fmt.Sprintf("138%08d", rand.Intn(100000000)),
			"测试地址",
			"紧急联系人",
			"13900000000",
			time.Now().AddDate(-1, 0, 0),
			careLevels[rand.Intn(len(careLevels))],
			statuses[rand.Intn(len(statuses))],
		)
		
		if err != nil {
			tx.Rollback()
			fmt.Printf("第%d条失败: %v\n", i, err)
			continue
		}
		
		if err := tx.Commit(); err != nil {
			fmt.Printf("第%d条提交失败: %v\n", i, err)
			continue
		}
		
		success++
		if success%10 == 0 {
			fmt.Printf("已成功插入 %d 条\n", success)
		}
	}
	
	fmt.Printf("\n完成! 成功: %d/100\n", success)
}
