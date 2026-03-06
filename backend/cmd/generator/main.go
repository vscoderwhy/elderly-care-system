package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

var (
	db          *sql.DB
	totalInsert int64
)

// 姓氏和名字库
var surnames = []string{
	"王", "李", "张", "刘", "陈", "杨", "黄", "赵", "吴", "周",
	"徐", "孙", "马", "朱", "胡", "郭", "何", "高", "林", "罗",
}

var lastNames = []string{
	"奶奶", "爷爷", "婆婆", "公公", "大爷", "阿婆",
}

var givenNames = []string{
	"秀英", "桂英", "秀珍", "玉兰", "淑珍", "淑兰", "秀兰", "玉兰", "桂兰", "秀华",
	"明", "伟", "芳", "娜", "敏", "静", "丽", "强", "磊", "军",
	"洋", "勇", "艳", "杰", "娟", "涛", "明", "超", "秀", "霞",
	"平", "刚", "桂英", "玉兰", "秀英", "淑珍", "华", "红", "海", "萍",
}

// 地址库
var provinces = []string{
	"北京市", "上海市", "广东省", "浙江省", "江苏省", "四川省", "山东省", "河南省",
}

var districts = []string{
	"朝阳区", "海淀区", "浦东新区", "黄浦区", "天河区", "西湖区", "江宁区", "武侯区",
}

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

func randomBirthDate() time.Time {
	now := time.Now()
	minAge := 60
	maxAge := 100
	randomAge := minAge + rand.Intn(maxAge-minAge+1)
	return now.AddDate(-randomAge, 0, 0)
}

func randomAdmissionDate() time.Time {
	now := time.Now()
	days := rand.Intn(3650) // 过去10年内
	return now.AddDate(-days/365, 0, -days%365)
}

func randomPhone() string {
	return fmt.Sprintf("138%08d", rand.Intn(100000000))
}

func randomIDCard() string {
	return fmt.Sprintf("%d%d%d", rand.Intn(1000000000), rand.Intn(10000), rand.Intn(10))
}

func randomAddress() string {
	province := provinces[rand.Intn(len(provinces))]
	district := districts[rand.Intn(len(districts))]
	street := fmt.Sprintf("%s街道%d号", []string{"人民", "解放", "建设", "和平", "友谊"}[rand.Intn(5)], rand.Intn(1000))
	return fmt.Sprintf("%s%s%s", province, district, street)
}

func randomName(gender string) string {
	surname := surnames[rand.Intn(len(surnames))]
	if rand.Intn(3) == 0 {
		return surname + lastNames[rand.Intn(len(lastNames))]
	}
	return surname + givenNames[rand.Intn(len(givenNames))]
}

// 获取一个随机的可用床位ID
func getRandomAvailableBed() sql.NullInt64 {
	var bedID sql.NullInt64
	
	// 查询一个可用的床位（状态为available或empty）
	err := db.QueryRow("SELECT id FROM beds WHERE status IN ('available', 'empty', NULL) ORDER BY RANDOM() LIMIT 1").Scan(&bedID.Int64)
	if err == nil {
		bedID.Valid = true
	}
	
	// 如果没有可用床位，查询任意床位
	if !bedID.Valid {
		err = db.QueryRow("SELECT id FROM beds ORDER BY RANDOM() LIMIT 1").Scan(&bedID.Int64)
		if err == nil {
			bedID.Valid = true
		}
	}
	
	return bedID
}

// 批量插入老人数据
func batchInsertElderly(count int, wg *sync.WaitGroup) {
	defer wg.Done()

	genders := []string{"男", "女"}
	careLevels := []int{1, 2, 3, 4, 5}
	statuses := []string{"active", "hospitalized", "left"}

	for i := 0; i < count; i++ {
		tx, err := db.Begin()
		if err != nil {
			log.Printf("开始事务失败: %v", err)
			continue
		}

		gender := genders[rand.Intn(len(genders))]
		name := randomName(gender)
		birthDate := randomBirthDate()
		
		// 获取一个随机床位（可能为空）
		bedID := getRandomAvailableBed()

		_, err = tx.Exec(`
			INSERT INTO elderly (name, gender, birth_date, id_card, phone, address,
				emergency_contact, emergency_phone, admission_date, bed_id, care_level, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			name,
			gender,
			birthDate,
			randomIDCard(),
			randomPhone(),
			randomAddress(),
			randomName(gender),
			randomPhone(),
			randomAdmissionDate(),
			bedID,
			careLevels[rand.Intn(len(careLevels))],
			statuses[rand.Intn(len(statuses))],
		)
		
		if err != nil {
			tx.Rollback()
			log.Printf("插入失败: %v", err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("提交事务失败: %v", err)
			continue
		}

		atomic.AddInt64(&totalInsert, 1)
		
		// 每1000条显示一次进度
		if (i+1)%1000 == 0 {
			log.Printf("已插入: %d 条", i+1)
		}
	}
}

func batchInsertCareRecords(elderlyCount int, recordPerElderly int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取护工ID
	var users []int
	rows, err := db.Query("SELECT id FROM users WHERE role = 'caregiver' ORDER BY id")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			rows.Scan(&id)
			users = append(users, id)
		}
	}
	
	if len(users) == 0 {
		users = []int{1}
	}

	statuses := []string{"completed", "pending", "skipped"}
	notes := []string{
		"血压正常", "血糖正常", "已服药", "洗澡完毕", "康复训练完成",
		"饮食配送完成", "状态良好", "需要关注", "情绪稳定", "生命体征平稳",
	}

	now := time.Now()
	inserted := 0
	
	for elderlyID := 1; elderlyID <= elderlyCount; elderlyID++ {
		for j := 0; j < recordPerElderly; j++ {
			tx, err := db.Begin()
			if err != nil {
				continue
			}

			days := rand.Intn(365)
			recordedAt := now.AddDate(0, 0, -days)

			_, err = tx.Exec(`
				INSERT INTO care_records (elderly_id, staff_id, care_item_id, recorded_at, status, notes)
				VALUES ($1, $2, $3, $4, $5, $6)
			`,
				elderlyID,
				users[rand.Intn(len(users))],
				rand.Intn(6)+1,
				recordedAt,
				statuses[rand.Intn(len(statuses))],
				notes[rand.Intn(len(notes))],
			)
			
			if err != nil {
				tx.Rollback()
				log.Printf("插入护理记录失败: %v", err)
				continue
			}

			if err := tx.Commit(); err != nil {
				log.Printf("提交护理记录失败: %v", err)
				continue
			}

			atomic.AddInt64(&totalInsert, 1)
			inserted++
			
			if inserted%10000 == 0 {
				log.Printf("已插入护理记录: %d 条", inserted)
			}
		}
	}
}

func batchInsertBuildingsFloorsRoomsBeds() {
	log.Println("创建楼栋、楼层、房间、床位...")

	buildingCount := 5
	floorsPerBuilding := 6
	roomsPerFloor := 20
	bedsPerRoom := 4

	var buildingIDs []int

	// 创建楼栋
	for i := 1; i <= buildingCount; i++ {
		var id int
		err := db.QueryRow("INSERT INTO buildings (name, floors_count) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING RETURNING id",
			fmt.Sprintf("%c栋", 'A'+i-1), floorsPerBuilding).Scan(&id)
		if err != nil {
			// 可能已存在，尝试查询
			err = db.QueryRow("SELECT id FROM buildings WHERE name = $1", fmt.Sprintf("%c栋", 'A'+i-1)).Scan(&id)
			if err != nil {
				log.Printf("获取楼栋失败: %v", err)
				continue
			}
		}
		buildingIDs = append(buildingIDs, id)
	}

	// 创建楼层、房间、床位
	for _, buildingID := range buildingIDs {
		for floor := 1; floor <= floorsPerBuilding; floor++ {
			var floorID int
			floorName := fmt.Sprintf("%d楼", floor)
			
			// 插入楼层
			err := db.QueryRow("INSERT INTO floors (building_id, name, sort_order) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id",
				buildingID, floorName, floor).Scan(&floorID)
			if err != nil {
				// 可能已存在，查询
				db.QueryRow("SELECT id FROM floors WHERE building_id = $1 AND name = $2", buildingID, floorName).Scan(&floorID)
			}

			// 创建房间
			for room := 1; room <= roomsPerFloor; room++ {
				roomNumber := fmt.Sprintf("%d%02d", floor, room)
				var roomID int
				
				err := db.QueryRow(`
					INSERT INTO rooms (building, floor, room_number, bed_count, floor_id, name, bed_capacity, sort_order)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING RETURNING id`,
					fmt.Sprintf("%c栋", 'A'+buildingID-1), floor, roomNumber, bedsPerRoom,
					floorID, roomNumber+"室", bedsPerRoom, (floor-1)*roomsPerFloor+room).Scan(&roomID)
				if err != nil {
					// 可能已存在，查询
					db.QueryRow("SELECT id FROM rooms WHERE room_number = $1", roomNumber).Scan(&roomID)
				}

				// 创建床位
				for bed := 1; bed <= bedsPerRoom; bed++ {
					bedNumber := fmt.Sprintf("%d", bed)
					bedName := fmt.Sprintf("%d号床", bed)
					status := "available"
					if rand.Intn(10) < 6 {
						status = "occupied"
					}
					
					_, err := db.Exec("INSERT INTO beds (room_id, bed_number, status, name) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
						roomID, bedNumber, status, bedName)
					if err != nil {
						log.Printf("插入床位失败: %v", err)
					}
				}
			}
		}
	}

	log.Printf("已创建: %d栋楼, %d层/栋, %d房间/层, %d床位/房间",
		buildingCount, floorsPerBuilding, roomsPerFloor, bedsPerRoom)
}

func insertStaff() {
	log.Println("创建护工账号...")

	positions := []string{"护士长", "主管护工", "护工", "护理员"}
	phones := []string{
		"13800138001", "13800138002", "13800138003", "13800138004",
		"13800138005", "13800138006", "13800138007", "13800138008",
		"13800138009", "13800138010",
	}

	success := 0
	for i, phone := range phones {
		name := fmt.Sprintf("%s%d", positions[i%len(positions)], i+1)
		passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMye1AJmK5Q8zL3p6J5z5G5z5G5z5G5z5G5"
		_, err := db.Exec("INSERT INTO users (phone, password_hash, name, role, status) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (phone) DO NOTHING",
			phone, passwordHash, name, "caregiver", "active")
		if err != nil {
			continue
		}
		success++
	}

	log.Printf("已创建 %d 个护工账号", success)
}

func showProgress(done chan bool) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	lastCount := int64(0)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			current := atomic.LoadInt64(&totalInsert)
			if current > lastCount {
				speed := (current - lastCount) / 10
				log.Printf("总进度: %d 条 (速度: %d 条/秒)", current, speed)
				lastCount = current
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: generator <数据类型>")
		fmt.Println("数据类型:")
		fmt.Println("  all       - 生成所有数据 (100万老人 + 护理记录)")
		fmt.Println("  elderly   - 只生成老人数据 (100万)")
		fmt.Println("  records   - 只生成护理记录")
		fmt.Println("  facility  - 只生成楼栋房间床位")
		fmt.Println("  staff     - 只生成护工账号")
		fmt.Println("  test      - 测试生成100条老人数据")
		os.Exit(1)
	}

	command := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	if err := initDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	done := make(chan bool)
	go showProgress(done)
	defer close(done)

	startTime := time.Now()

	switch command {
	case "facility":
		batchInsertBuildingsFloorsRoomsBeds()
		insertStaff()

	case "staff":
		insertStaff()

	case "test":
		log.Println("测试模式：生成100条老人数据...")
		var wg sync.WaitGroup
		wg.Add(1)
		batchInsertElderly(100, &wg)
		wg.Wait()

	case "elderly":
		batchInsertBuildingsFloorsRoomsBeds()

		elderlyCount := 1000000
		batchSize := 10000
		batchCount := elderlyCount / batchSize

		log.Printf("开始生成 %d 条老人数据 (每批 %d 条)", elderlyCount, batchSize)

		var wg sync.WaitGroup
		for i := 0; i < batchCount; i++ {
			wg.Add(1)
			go batchInsertElderly(batchSize, &wg)
			time.Sleep(20 * time.Millisecond)
		}
		wg.Wait()

	case "records":
		recordPerElderly := 50

		var elderlyCount int
		db.QueryRow("SELECT COUNT(*) FROM elderly").Scan(&elderlyCount)

		log.Printf("开始为 %d 个老人生成护理记录 (每人 %d 条)", elderlyCount, recordPerElderly)

		var wg sync.WaitGroup
		wg.Add(1)
		batchInsertCareRecords(elderlyCount, recordPerElderly, &wg)
		wg.Wait()

	case "all":
		log.Println("=== 开始生成完整的百万级数据 ===")

		log.Println("\n[1/4] 创建楼栋、楼层、房间、床位...")
		batchInsertBuildingsFloorsRoomsBeds()

		log.Println("\n[2/4] 创建护工账号...")
		insertStaff()

		log.Println("\n[3/4] 生成100万老人数据...")
		elderlyCount := 1000000
		batchSize := 10000
		batchCount := elderlyCount / batchSize

		var wg sync.WaitGroup
		for i := 0; i < batchCount; i++ {
			wg.Add(1)
			go batchInsertElderly(batchSize, &wg)
			time.Sleep(20 * time.Millisecond)
		}
		wg.Wait()

		log.Println("\n[4/4] 生成护理记录...")
		recordPerElderly := 50
		wg.Add(1)
		batchInsertCareRecords(elderlyCount, recordPerElderly, &wg)
		wg.Wait()

	default:
		fmt.Printf("未知命令: %s\n", command)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)

	log.Println("\n=== 数据生成完成 ===")
	log.Printf("总耗时: %v", elapsed)
	log.Printf("总插入: %d 条", atomic.LoadInt64(&totalInsert))

	var elderlyCount, buildingCount, roomCount, bedCount, recordCount int
	db.QueryRow("SELECT COUNT(*) FROM elderly").Scan(&elderlyCount)
	db.QueryRow("SELECT COUNT(*) FROM buildings").Scan(&buildingCount)
	db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&roomCount)
	db.QueryRow("SELECT COUNT(*) FROM beds").Scan(&bedCount)
	db.QueryRow("SELECT COUNT(*) FROM care_records").Scan(&recordCount)

	log.Println("\n=== 数据库统计 ===")
	log.Printf("老人总数: %d", elderlyCount)
	log.Printf("楼栋总数: %d", buildingCount)
	log.Printf("房间总数: %d", roomCount)
	log.Printf("床位总数: %d", bedCount)
	log.Printf("护理记录: %d", recordCount)
}
