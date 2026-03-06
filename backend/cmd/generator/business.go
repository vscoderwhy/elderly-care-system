package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=elderly_care sslmode=disable")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	if err = db.Ping(); err != nil {
		return err
	}
	log.Println("数据库连接成功")
	return nil
}

func getElderlyIDs() ([]int, error) {
	var ids []int
	rows, err := db.Query("SELECT id FROM elderly ORDER BY id LIMIT 10000")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func getStaffIDs() ([]int, error) {
	var ids []int
	rows, err := db.Query("SELECT id FROM users WHERE role = 'caregiver' ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		ids = []int{1}
	}
	return ids, nil
}

func generateMedicationRecords(elderlyIDs []int) error {
	log.Println("生成用药记录...")

	medications := []struct {
		name              string
		genericName       string
		specification     string
		unit              string
		manufacturer      string
		usageInstructions string
	}{
		{"阿司匹林", "乙酰水杨酸", "100mg", "片", "拜耳", "口服，每日一次，每次一片"},
		{"二甲双胍", "盐酸二甲双胍", "500mg", "片", "中美施贵宝", "口服，每日两次，每次一片，餐后服用"},
		{"氨氯地平", "苯磺酸氨氯地平", "5mg", "片", "辉瑞", "口服，每日一次，每次一片"},
		{"辛伐他汀", "辛伐他汀", "20mg", "片", "默沙东", "口服，每日一次，睡前服用"},
		{"奥美拉唑", "奥美拉唑", "20mg", "胶囊", "阿斯利康", "口服，每日一次，晨起空腹服用"},
		{"格列吡嗪", "格列吡嗪", "5mg", "片", "辉瑞", "口服，每日两次，餐前30分钟"},
		{"美托洛尔", "酒石酸美托洛尔", "25mg", "片", "阿斯利康", "口服，每日两次"},
		{"氯吡格雷", "硫酸氢氯吡格雷", "75mg", "片", "赛诺菲", "口服，每日一次"},
		{"胰岛素", "生物合成人胰岛素", "400IU/ml", "支", "诺和诺德", "皮下注射，每日三次，餐前15分钟"},
		{"华法林", "华法林钠", "2.5mg", "片", "百时美", "口服，每日一次，根据INR调整剂量"},
	}

	var medicationIDs []int
	for _, med := range medications {
		var id int
		err := db.QueryRow(`
			INSERT INTO medications (name, generic_name, specification, unit, stock, min_stock, 
				expiry_date, manufacturer, usage_instructions, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`, med.name, med.genericName, med.specification, med.unit, 
			rand.Intn(500)+100, 100, 
			time.Now().AddDate(1, rand.Intn(12), 0), 
			med.manufacturer, med.usageInstructions, "active").Scan(&id)
		if err != nil {
			log.Printf("插入药品失败: %v", err)
			continue
		}
		medicationIDs = append(medicationIDs, id)
	}

	log.Printf("已创建 %d 种药品", len(medicationIDs))

	if len(medicationIDs) == 0 {
		log.Println("警告：没有创建任何药品，跳过用药记录生成")
		return nil
	}

	staffIDs, _ := getStaffIDs()
	count := 0
	for _, elderlyID := range elderlyIDs {
		numMeds := rand.Intn(4) + 2
		for i := 0; i < numMeds; i++ {
			medID := medicationIDs[rand.Intn(len(medicationIDs))]
			dosage := []string{"每次一片", "每次两片", "每次一片半", "每次2.5mg"}[rand.Intn(4)]
			frequency := []string{"每日一次", "每日两次", "每日三次", "每日四次", "隔日一次"}[rand.Intn(5)]
			
			var recordID int
			err := db.QueryRow(`
				INSERT INTO medication_records (elderly_id, medication_id, dosage, frequency, 
					prescribed_date, prescribed_by, notes, status)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id
			`, elderlyID, medID, dosage, frequency,
				time.Now().AddDate(-rand.Intn(365), 0, 0),
				staffIDs[rand.Intn(len(staffIDs))],
				"遵医嘱服用", "active").Scan(&recordID)
			
			if err == nil {
				count++
			}
		}
	}
	
	log.Printf("已生成 %d 条用药记录", count)
	return nil
}

func generateBills(elderlyIDs []int) error {
	log.Println("生成账单数据...")

	count := 0
	for _, elderlyID := range elderlyIDs {
		for month := 0; month < 12; month++ {
			billPeriodStart := time.Now().AddDate(0, -month-1, 1)
			billPeriodEnd := billPeriodStart.AddDate(0, 1, -1)
			
			baseFee := float64(3000+rand.Intn(5000))
			careLevelFee := float64(rand.Intn(5)+1) * 1000
			mealFee := float64(1500+rand.Intn(500))
			otherFee := float64(rand.Intn(1000))
			totalAmount := baseFee + careLevelFee + mealFee + otherFee
			
			var billID int
			err := db.QueryRow(`
				INSERT INTO bills (elderly_id, bill_no, total_amount, status, 
					bill_period_start, bill_period_end, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
			`, elderlyID, 
				fmt.Sprintf("BILL%d%04d%02d", elderlyID, billPeriodStart.Year(), int(billPeriodStart.Month())),
				totalAmount,
				[]string{"paid", "unpaid", "partial"}[rand.Intn(3)],
				billPeriodStart, billPeriodEnd,
				billPeriodStart.AddDate(0, 0, rand.Intn(10))).Scan(&billID)
			
			if err != nil {
				continue
			}
			
			count++
			
			db.Exec(`
				INSERT INTO bill_items (bill_id, item_name, quantity, unit_price, amount)
				VALUES 
					($1, '护理费', 1, $2, $2),
					($1, '床位费', 1, $3, $3),
					($1, '餐费', 1, $4, $4),
					($1, '水电费', 1, $5, $5)
			`, billID, careLevelFee, baseFee, mealFee, otherFee)
			
			if rand.Intn(10) < 7 {
				paidAt := billPeriodEnd.AddDate(0, 0, rand.Intn(30))
				paymentMethod := []string{"cash", "card", "transfer", "wechat", "alipay"}[rand.Intn(5)]
				db.Exec(`
					INSERT INTO payments (bill_id, amount, payment_method, paid_at, status)
					VALUES ($1, $2, $3, $4, 'completed')
				`, billID, totalAmount, paymentMethod, paidAt)
			}
		}
	}
	
	log.Printf("已生成 %d 条账单", count)
	return nil
}

func generateStaff() error {
	log.Println("生成员工数据...")

	positions := []string{"护士长", "主管护工", "护工", "护理员", "康复师", "营养师", "医生", "清洁工", "保安", "行政"}
	
	count := 0
	for i := 0; i < 50; i++ {
		phone := fmt.Sprintf("139%08d", i+1000)
		name := fmt.Sprintf("%s%d", positions[i%len(positions)], i+1)
		passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMye1AJmK5Q8zL3p6J5z5G5z5G5z5G5z5G5"
		
		var id int
		err := db.QueryRow(`
			INSERT INTO users (phone, password_hash, name, role, status, nickname)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (phone) DO NOTHING
			RETURNING id
		`, phone, passwordHash, name, "caregiver", "active", name).Scan(&id)
		
		if err == nil {
			count++
			generateAttendance(id)
		}
	}
	
	log.Printf("已生成 %d 个员工", count)
	return nil
}

func generateAttendance(staffID int) {
	for day := 0; day < 30; day++ {
		date := time.Now().AddDate(0, 0, -day)
		
		if rand.Intn(10) < 8 {
			checkIn := time.Date(date.Year(), date.Month(), date.Day(), 7+rand.Intn(2), rand.Intn(60), 0, 0, time.Local)
			checkOut := time.Date(date.Year(), date.Month(), date.Day(), 17+rand.Intn(2), rand.Intn(60), 0, 0, time.Local)
			
			db.Exec(`
				INSERT INTO attendances (staff_id, date, check_in, check_out, status, work_hours)
				VALUES ($1, $2, $3, $4, 'present', $5)
			`, staffID, date, checkIn, checkOut, 8.0)
		} else {
			status := []string{"leave", "sick", "absent"}[rand.Intn(3)]
			db.Exec(`
				INSERT INTO attendances (staff_id, date, status, work_hours)
				VALUES ($1, $2, $3, 0)
			`, staffID, date, status)
		}
	}
}

func generateHealthRecords(elderlyIDs []int) error {
	log.Println("生成健康记录...")

	staffIDs, _ := getStaffIDs()
	
	recordTypes := []struct {
		recordType string
		unit       string
		min        float64
		max        float64
	}{
		{"blood_pressure_systolic", "mmHg", 90, 160},
		{"blood_pressure_diastolic", "mmHg", 60, 100},
		{"heart_rate", "bpm", 60, 100},
		{"blood_sugar", "mmol/L", 3.9, 11.1},
		{"temperature", "℃", 36.0, 37.5},
		{"oxygen_saturation", "%", 90, 100},
		{"weight", "kg", 40, 90},
	}

	count := 0
	for _, elderlyID := range elderlyIDs {
		for day := 0; day < 30; day++ {
			recordedAt := time.Now().AddDate(0, 0, -day)
			
			if rand.Intn(10) < 8 {
				for _, rt := range recordTypes {
					value := rt.min + float64(rand.Intn(int((rt.max-rt.min)*10)))/10.0
					
					var recordID int
					err := db.QueryRow(`
						INSERT INTO health_records (elderly_id, record_type, value, unit, recorded_at, recorded_by)
						VALUES ($1, $2, $3, $4, $5, $6)
						RETURNING id
					`, elderlyID, rt.recordType, fmt.Sprintf("%.1f", value), rt.unit, recordedAt,
						staffIDs[rand.Intn(len(staffIDs))]).Scan(&recordID)
					
					if err == nil {
						count++
					}
				}
			}
		}
	}
	
	log.Printf("已生成 %d 条健康记录", count)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := initDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	startTime := time.Now()

	log.Println("\n=== 开始生成业务数据 ===\n")
	generateStaff()

	elderlyIDs, err := getElderlyIDs()
	if err != nil {
		log.Fatalf("获取老人ID失败: %v", err)
	}
	log.Printf("获取到 %d 个老人", len(elderlyIDs))

	generateMedicationRecords(elderlyIDs)
	generateBills(elderlyIDs)
	generateHealthRecords(elderlyIDs)

	elapsed := time.Since(startTime)

	var medicationCount, billCount, healthCount, attendanceCount, staffCount int
	db.QueryRow("SELECT COUNT(*) FROM medication_records").Scan(&medicationCount)
	db.QueryRow("SELECT COUNT(*) FROM bills").Scan(&billCount)
	db.QueryRow("SELECT COUNT(*) FROM health_records").Scan(&healthCount)
	db.QueryRow("SELECT COUNT(*) FROM attendances").Scan(&attendanceCount)
	db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'caregiver'").Scan(&staffCount)

	log.Println("\n=== 数据生成完成 ===")
	log.Printf("总耗时: %v", elapsed)
	log.Println("\n=== 数据库统计 ===")
	log.Printf("员工总数: %d", staffCount)
	log.Printf("用药记录: %d", medicationCount)
	log.Printf("账单总数: %d", billCount)
	log.Printf("健康记录: %d", healthCount)
	log.Printf("考勤记录: %d", attendanceCount)
}
