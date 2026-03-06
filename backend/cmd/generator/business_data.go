package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

var (
	db          *sql.DB
	totalInsert int64
)

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
		ids = []int{1} // 默认管理员
	}
	return ids, nil
}

// 生成用药记录
func generateMedicationRecords(elderlyIDs []int) error {
	log.Println("生成用药记录...")

	// 先创建一些药品

func generateMedicationLogs(recordID int, staffIDs []int) {
	// 过去30天内的用药记录
	numLogs := rand.Intn(30) + 1
	for i := 0; i < numLogs; i++ {
		takenAt := time.Now().AddDate(0, 0, -rand.Intn(30))
		status := []string{"taken", "taken", "taken", "missed", "refused"}[rand.Intn(5)]
		
		_, err := db.Exec(`
			INSERT INTO medication_logs (medication_record_id, staff_id, taken_at, status, notes)
			VALUES ($1, $2, $3, $4, $5)
		`, recordID, staffIDs[rand.Intn(len(staffIDs))], takenAt, status, "")
		
		if err != nil {
			// 忽略错误，可能表不存在
		}
	}
}

// 生成账单数据
func generateBills(elderlyIDs []int) error {
	log.Println("生成账单数据...")

	staffIDs, _ := getStaffIDs()
	
	count := 0
	for _, elderlyID := range elderlyIDs {
		// 每个老人过去12个月的账单
		for month := 0; month < 12; month++ {
			billPeriodStart := time.Now().AddDate(0, -month-1, 1)
			billPeriodEnd := billPeriodStart.AddDate(0, 1, -1)
			
			// 基础费用
			baseFee := float64(3000+rand.Intn(5000))
			
			// 护理等级费用
			careLevelFee := float64(rand.Intn(5)+1) * 1000
			
			// 餐费
			mealFee := float64(1500+rand.Intn(500))
			
			// 其他费用
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
				log.Printf("插入账单失败: %v", err)
				continue
			}
			
			count++
			
			// 生成账单明细
			db.Exec(`
				INSERT INTO bill_items (bill_id, item_name, quantity, unit_price, amount)
				VALUES 
					($1, '护理费', 1, $2, $2),
					($1, '床位费', 1, $3, $3),
					($1, '餐费', 1, $4, $4),
					($1, '水电费', 1, $5, $5)
			`, billID, careLevelFee, baseFee, mealFee, otherFee)
			
			// 如果已支付，生成支付记录
			if rand.Intn(10) < 7 { // 70%已支付
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

// 生成员工数据
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
			// 生成考勤记录
			generateAttendance(id)
		}
	}
	
	log.Printf("已生成 %d 个员工", count)
	return nil
}

func generateAttendance(staffID int) {
	// 过去30天的考勤
	for day := 0; day < 30; day++ {
		date := time.Now().AddDate(0, 0, -day)
		
		// 80%出勤率
		if rand.Intn(10) < 8 {
			checkIn := time.Date(date.Year(), date.Month(), date.Day(), 7+rand.Intn(2), rand.Intn(60), 0, 0, time.Local)
			checkOut := time.Date(date.Year(), date.Month(), date.Day(), 17+rand.Intn(2), rand.Intn(60), 0, 0, time.Local)
			
			db.Exec(`
				INSERT INTO attendances (staff_id, date, check_in, check_out, status, work_hours)
				VALUES ($1, $2, $3, $4, 'present', $5)
			`, staffID, date, checkIn, checkOut, 8.0)
		} else {
			// 请假或缺勤
			status := []string{"leave", "sick", "absent"}[rand.Intn(3)]
			db.Exec(`
				INSERT INTO attendances (staff_id, date, status, work_hours)
				VALUES ($1, $2, $3, 0)
			`, staffID, date, status)
		}
	}
}

// 生成健康记录
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
		// 每个老人过去30天的健康记录
		for day := 0; day < 30; day++ {
			recordedAt := time.Now().AddDate(0, 0, -day)
			
			// 80%概率有记录
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

	// 1. 生成员工
	log.Println("\n=== 开始生成业务数据 ===\n")
	generateStaff()

	// 2. 获取老人ID
	elderlyIDs, err := getElderlyIDs()
	if err != nil {
		log.Fatalf("获取老人ID失败: %v", err)
	}
	log.Printf("获取到 %d 个老人", len(elderlyIDs))

	// 3. 生成用药记录
	generateMedicationRecords(elderlyIDs)

	// 4. 生成账单
	generateBills(elderlyIDs)

	// 5. 生成健康记录
	generateHealthRecords(elderlyIDs)

	elapsed := time.Since(startTime)

	// 统计
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
