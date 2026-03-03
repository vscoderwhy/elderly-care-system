package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 数据模型定义
type Elderly struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Gender          string    `json:"gender"`
	Age             int       `json:"age"`
	IDCard          string    `json:"idCard"`
	Birthday        string    `json:"birthday"`
	Phone           string    `json:"phone"`
	Address         string    `json:"address"`
	CareLevel       string    `json:"careLevel"`
	BedNumber       string    `json:"bedNumber"`
	AdmitDate       string    `json:"admitDate"`
	Status          string    `json:"status"`
	Photo           string    `json:"photo"`
	EmergencyContact string   `json:"emergencyContact"`
	EmergencyPhone   string   `json:"emergencyPhone"`
	HealthScore     int       `json:"healthScore"`
	Remarks          string    `json:"remarks"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type Staff struct {
	ID           uint      `json:"id"`
	EmployeeNo   string    `json:"employeeNo"`
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	Phone        string    `json:"phone"`
	IDCard       string    `json:"idCard"`
	Department   string    `json:"department"`
	Position     string    `json:"position"`
	Education    string    `json:"education"`
	HireDate     string    `json:"hireDate"`
	Status       string    `json:"status"`
	Avatar       string    `json:"avatar"`
	Address      string    `json:"address"`
	Remarks      string    `json:"remarks"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CareRecord struct {
	ID          uint      `json:"id"`
	ElderlyID   uint      `json:"elderlyId"`
	CareType    string    `json:"careType"`
	Content     string    `json:"content"`
	CareTime    time.Time `json:"careTime"`
	NurseID     uint      `json:"nurseId"`
	NurseName   string    `json:"nurseName"`
	Images      string    `json:"images"`
	Evaluation  int       `json:"evaluation"`
	Tags        string    `json:"tags"`
	Remarks     string    `json:"remarks"`
	CreatedAt   time.Time `json:"createdAt"`
}

type HealthData struct {
	ID            uint      `json:"id"`
	ElderlyID     uint      `json:"elderlyId"`
	RecordDate    string    `json:"recordDate"`
	RecordTime    string    `json:"recordTime"`
	BloodPressure string    `json:"bloodPressure"`
	HeartRate     int       `json:"heartRate"`
	Temperature   float64   `json:"temperature"`
	BloodSugar    float64   `json:"bloodSugar"`
	Weight        float64   `json:"weight"`
	SPO2          int       `json:"spo2"`
	Recorder      string    `json:"recorder"`
	Remarks       string    `json:"remarks"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Bill struct {
	ID          uint      `json:"id"`
	BillNo      string    `json:"billNo"`
	ElderlyID   uint      `json:"elderlyId"`
	BillType    string    `json:"billType"`
	Amount      float64   `json:"amount"`
	Period      string    `json:"period"`
	BillDate    string    `json:"billDate"`
	DueDate     string    `json:"dueDate"`
	Status      string    `json:"status"`
	PaymentDate string    `json:"paymentDate"`
	PaymentMethod string  `json:"paymentMethod"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CareTask struct {
	ID            uint      `json:"id"`
	ElderlyID     uint      `json:"elderlyId"`
	TaskType      string    `json:"taskType"`
	Description   string    `json:"description"`
	ScheduledDate string    `json:"scheduledDate"`
	ScheduledTime string    `json:"scheduledTime"`
	NurseID       uint      `json:"nurseId"`
	Status        string    `json:"status"`
	Priority      string    `json:"priority"`
	CreatedAt     time.Time `json:"createdAt"`
}

type VisitAppointment struct {
	ID              uint      `json:"id"`
	AppointmentNo   string    `json:"appointmentNo"`
	ElderlyID       uint      `json:"elderlyId"`
	VisitorName     string    `json:"visitorName"`
	VisitorPhone    string    `json:"visitorPhone"`
	Relationship    string    `json:"relationship"`
	VisitType       string    `json:"visitType"`
	VisitDate       string    `json:"visitDate"`
	VisitTime       string    `json:"visitTime"`
	VisitorCount    int       `json:"visitorCount"`
	Status          string    `json:"status"`
	Remarks          string    `json:"remarks"`
	CreatedAt        time.Time `json:"createdAt"`
}

var (
	surnames = []string{"王", "李", "张", "刘", "陈", "杨", "黄", "赵", "周", "吴", "徐", "孙", "马", "胡", "朱", "郭", "何", "罗", "高", "林", "郑"}
	femaleNames = []string{"秀英", "桂英", "秀珍", "淑兰", "玉兰", "淑华", "翠英", "文英", "秀兰", "玉英", "玉珍", "淑珍", "秀华", "英", "玉兰", "淑芬", "秀芬", "淑琴", "玉芬", "玉华"}
	maleNames = []string{"伟", "强", "磊", "洋", "勇", "军", "杰", "涛", "超", "明", "刚", "平", "辉", "鹏", "华", "飞", "鑫", "波", "斌", "宇", "浩"}

	careTypes   = []string{"日常护理", "康复训练", "健康监测", "医疗护理", "心理疏导", "营养指导"}
	careLevels  = []string{"特级", "一级", "二级", "三级"}
	departments = []string{"护理部", "医务室", "康复科", "膳食部", "行政部"}
	positions   = map[string][]string{
		"护理部":  {"护士长", "护士", "护理员", "护工"},
		"医务室":  {"主任医师", "副主任医师", "主治医师", "医师", "护士"},
		"康复科":  {"康复师", "物理治疗师", "作业治疗师"},
		"膳食部":  {"营养师", "厨师", "配餐员"},
		"行政部":  {"院长", "副院长", "主任", "行政专员"},
	}
	billTypes    = []string{"床位费", "护理费", "伙食费", "医疗费", "康复费", "用品费", "其他费用"}
	taskTypes    = []string{"日常护理", "健康监测", "康复训练", "医疗护理", "协助用餐", "协助洗漱", "测量体征", "翻身护理"}
	relationships = []string{"子女", "配偶", "孙辈", "其他亲属", "朋友", "其他"}

	buildings   = []string{"1号楼", "2号楼", "3号楼", "4号楼"}
	floors      = []string{"1楼", "2楼", "3楼", "4楼", "5楼"}
	roomNumbers = []string{"101", "102", "103", "104", "105", "201", "202", "203", "204", "205"}

	careContents = map[string][]string{
		"日常护理": {"协助起床洗漱", "协助用餐", "协助如厕", "更换衣物", "整理床铺", "室内通风", " personal卫生清洁"},
		"康复训练": {"上肢关节活动训练", "下肢功能训练", "平衡能力训练", "步行训练", "手部精细动作训练", "语言康复训练"},
		"健康监测": {"测量血压", "测量心率", "测量体温", "测量血糖", "测量血氧饱和度", "测量体重", "健康询问"},
		"医疗护理": {"伤口换药", "静脉输液", "肌肉注射", "雾化吸入", "吸氧护理", "导尿护理", "胃管护理"},
		"心理疏导": {"心理支持", "情绪疏导", "陪伴聊天", "音乐疗法", "回忆疗法", "认知训练"},
		"营养指导": {"营养评估", "膳食建议", "进食指导", "吞咽训练", "营养宣教"},
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 连接数据库
	dsn := "host=localhost user=elderly_care password=elderly_care_password dbname=elderly_care_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}

	fmt.Println("==========================================")
	fmt.Println("  养老院管理系统 - 数据种子生成器")
	fmt.Println("==========================================")
	fmt.Println()

	// 生成数据
	generateElderly(db, 150)
	generateStaff(db, 60)
	generateCareRecords(db, 300)
	generateHealthData(db, 500)
	generateBills(db, 200)
	generateCareTasks(db, 150)
	generateVisitAppointments(db, 100)

	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("  数据生成完成！")
	fmt.Println("==========================================")
}

func generateElderly(db *gorm.DB, count int) {
	fmt.Printf("📝 正在生成 %d 条老人数据...\n", count)

	var elderlyList []Elderly
	now := time.Now()

	for i := 1; i <= count; i++ {
		gender := rand.Intn(2)
		isMale := gender == 1

		name := randomName(isMale)
		age := 60 + rand.Intn(40) // 60-100岁
		birthdayYear := now.Year() - age

		elderly := Elderly{
			ID:              uint(i),
			Name:            name,
			Gender:          map[bool]string{true: "男", false: "女"}[isMale],
			Age:             age,
			IDCard:          generateIDCard(birthdayYear, isMale),
			Birthday:        fmt.Sprintf("%d-%02d-%02d", birthdayYear, 1+rand.Intn(12), 1+rand.Intn(28)),
			Phone:           generatePhone(),
			Address:         randomAddress(),
			CareLevel:       careLevels[rand.Intn(len(careLevels))],
			BedNumber:       randomBedNumber(),
			AdmitDate:       randomDate(now.AddDate(-3, 0, 0), now),
			Status:          "在院",
			EmergencyContact: randomName(true),
			EmergencyPhone:   generatePhone(),
			HealthScore:     60 + rand.Intn(40),
			Remarks:         randomRemarks(),
			CreatedAt:        randomTime(now.AddDate(-1, 0, 0), now),
			UpdatedAt:        randomTime(now.AddDate(-1, 0, 0), now),
		}

		elderlyList = append(elderlyList, elderly)
	}

	// 输出JSON到文件
	data, _ := json.MarshalIndent(elderlyList, "", "  ")
	os.MkdirAll("seed_data", 0755)
	os.WriteFile("seed_data/elderly.json", data, 0644)

	fmt.Printf("   ✅ 老人数据已生成到 seed_data/elderly.json\n")
}

func generateStaff(db *gorm.DB, count int) {
	fmt.Printf("👨‍⚕️  正在生成 %d 条员工数据...\n", count)

	var staffList []Staff
	now := time.Now()

	for i := 1; i <= count; i++ {
		department := departments[rand.Intn(len(departments))]
		position := positions[department][rand.Intn(len(positions[department]))]

		staff := Staff{
			ID:           uint(i),
			EmployeeNo:   fmt.Sprintf("EMP%03d", i),
			Name:         randomName(rand.Intn(2) == 1),
			Gender:       map[bool]string{true: "男", false: "女"}[rand.Intn(2) == 1],
			Phone:        generatePhone(),
			IDCard:       generateIDCard(now.Year()-25-rand.Intn(40), rand.Intn(2) == 1),
			Department:   department,
			Position:     position,
			Education:    []string{"高中", "大专", "本科", "硕士"}[rand.Intn(4)],
			HireDate:     randomDate(now.AddDate(-5, 0, 0), now),
			Status:       "在职",
			Address:      randomAddress(),
			Remarks:     "",
			CreatedAt:    randomTime(now.AddDate(-2, 0, 0), now),
			UpdatedAt:    randomTime(now.AddDate(-2, 0, 0), now),
		}

		staffList = append(staffList, staff)
	}

	data, _ := json.MarshalIndent(staffList, "", "  ")
	os.WriteFile("seed_data/staff.json", data, 0644)

	fmt.Printf("   ✅ 员工数据已生成到 seed_data/staff.json\n")
}

func generateCareRecords(db *gorm.DB, count int) {
	fmt.Printf("📋 正在生成 %d 条护理记录数据...\n", count)

	var recordList []CareRecord
	now := time.Now()

	elderlyCount := 150
	staffCount := 60

	for i := 1; i <= count; i++ {
		careType := careTypes[rand.Intn(len(careTypes))]
		elderlyID := uint(rand.Intn(elderlyCount) + 1)
		staffID := uint(rand.Intn(staffCount) + 1)

		record := CareRecord{
			ID:         uint(i),
			ElderlyID:  elderlyID,
			CareType:   careType,
			Content:    randomCareContent(careType),
			CareTime:   randomTime(now.AddDate(-1, 0, 0), now),
			NurseID:    staffID,
			NurseName:  randomName(rand.Intn(2) == 1),
			Images:     randomImages(),
			Evaluation: rand.Intn(3) + 3, // 3-5分
			Tags:       randomTags(),
			Remarks:    randomRemarks(),
			CreatedAt:  randomTime(now.AddDate(-1, 0, 0), now),
		}

		recordList = append(recordList, record)
	}

	data, _ := json.MarshalIndent(recordList, "", "  ")
	os.WriteFile("seed_data/care_records.json", data, 0644)

	fmt.Printf("   ✅ 护理记录数据已生成到 seed_data/care_records.json\n")
}

func generateHealthData(db *gorm.DB, count int) {
	fmt.Printf("❤️  正在生成 %d 条健康数据...\n", count)

	var healthList []HealthData
	now := time.Now()

	elderlyCount := 150

	for i := 1; i <= count; i++ {
		elderlyID := uint(rand.Intn(elderlyCount) + 1)
		recordTime := randomTime(now.AddDate(-1, 0, 0), now)

		// 生成合理的健康数据
		systolic := 110 + rand.Intn(40)   // 收缩压 110-150
		diastolic := 70 + rand.Intn(20)    // 舒张压 70-90

		health := HealthData{
			ID:            uint(i),
			ElderlyID:     elderlyID,
			RecordDate:    recordTime.Format("2006-01-02"),
			RecordTime:    recordTime.Format("15:04"),
			BloodPressure: fmt.Sprintf("%d/%d", systolic, diastolic),
			HeartRate:     60 + rand.Intn(40),  // 60-100
			Temperature:   36.0 + rand.Float64()*1.5, // 36.0-37.5
			BloodSugar:    4.0 + rand.Float64()*6.0,   // 4.0-10.0
			Weight:        40.0 + rand.Float64()*40.0, // 40.0-80.0
			SPO2:          92 + rand.Intn(8),         // 92-100
			Recorder:      randomName(rand.Intn(2) == 1),
			Remarks:       randomHealthRemarks(),
			CreatedAt:     recordTime,
		}

		healthList = append(healthList, health)
	}

	data, _ := json.MarshalIndent(healthList, "", "  ")
	os.WriteFile("seed_data/health_data.json", data, 0644)

	fmt.Printf("   ✅ 健康数据已生成到 seed_data/health_data.json\n")
}

func generateBills(db *gorm.DB, count int) {
	fmt.Printf("💰 正在生成 %d 条账单数据...\n", count)

	var billList []Bill
	now := time.Now()

	elderlyCount := 150

	for i := 1; i <= count; i++ {
		billType := billTypes[rand.Intn(len(billTypes))]
		elderlyID := uint(rand.Intn(elderlyCount) + 1)

		// 根据账单类型生成金额
		baseAmount := map[string]float64{
			"床位费": 2800,
			"护理费": 500 + rand.Float64()*1000,
			"伙食费": 600,
			"医疗费": 100 + rand.Float64()*5000,
			"康复费": 200 + rand.Float64()*1000,
			"用品费": 50 + rand.Float64()*500,
			"其他费用": 100 + rand.Float64()*1000,
		}

		amount := baseAmount[billType]
		if amount == 0 {
			amount = 100 + rand.Float64()*1000
		}

		billDate := randomDate(now.AddDate(-6, 0, 0), now)
		dueDate := billDate.AddDate(0, 0, 10)

		statuses := []string{"paid", "unpaid", "overdue"}
		status := statuses[rand.Intn(len(statuses))]

		bill := Bill{
			ID:       uint(i),
			BillNo:   fmt.Sprintf("B%s%03d", now.Format("200601"), i),
			ElderlyID: elderlyID,
			BillType: billType,
			Amount:   float64(int(amount*100) / 100),
			Period:   fmt.Sprintf("%s年%d月", billDate.Format("2006"), int(billDate.Month())),
			BillDate: billDate.Format("2006-01-02"),
			DueDate:  dueDate.Format("2006-01-02"),
			Status:   status,
		}

		if status == "paid" {
			bill.PaymentDate = randomDate(billDate, now.AddDate(0, 0, 5)).Format("2006-01-02")
			bill.PaymentMethod = []string{"微信支付", "支付宝", "现金", "银行转账"}[rand.Intn(4)]
		}

		billList = append(billList, bill)
	}

	data, _ := json.MarshalIndent(billList, "", "  ")
	os.WriteFile("seed_data/bills.json", data, 0644)

	fmt.Printf("   ✅ 账单数据已生成到 seed_data/bills.json\n")
}

func generateCareTasks(db *gorm.DB, count int) {
	fmt.Printf("✅ 正在生成 %d 条护理任务数据...\n", count)

	var taskList []CareTask
	now := time.Now()

	elderlyCount := 150
	staffCount := 60

	for i := 1; i <= count; i++ {
		taskType := taskTypes[rand.Intn(len(taskTypes))]
		elderlyID := uint(rand.Intn(elderlyCount) + 1)
		staffID := uint(rand.Intn(staffCount) + 1)

		scheduledDate := randomDate(now, now.AddDate(0, 0, 7))
		hour := 8 + rand.Intn(10) // 8-18点

		statuses := []string{"pending", "in_progress", "completed"}
		status := statuses[rand.Intn(len(statuses))]

		priorities := []string{"normal", "important", "urgent"}
		priority := priorities[rand.Intn(len(priorities))]

		task := CareTask{
			ID:            uint(i),
			ElderlyID:     elderlyID,
			TaskType:      taskType,
			Description:   randomCareContent(taskType),
			ScheduledDate: scheduledDate.Format("2006-01-02"),
			ScheduledTime: fmt.Sprintf("%02d:00", hour),
			NurseID:       staffID,
			Status:        status,
			Priority:      priority,
			CreatedAt:     randomTime(now.AddDate(-1, 0, 0), now),
		}

		taskList = append(taskList, task)
	}

	data, _ := json.MarshalIndent(taskList, "", "  ")
	os.WriteFile("seed_data/care_tasks.json", data, 0644)

	fmt.Printf("   ✅ 护理任务数据已生成到 seed_data/care_tasks.json\n")
}

func generateVisitAppointments(db *gorm.DB, count int) {
	fmt.Printf("📅 正在生成 %d 条探视预约数据...\n", count)

	var visitList []VisitAppointment
	now := time.Now()

	elderlyCount := 150

	for i := 1; i <= count; i++ {
		elderlyID := uint(rand.Intn(elderlyCount) + 1)

		visitDate := randomDate(now, now.AddDate(0, 1, 0))
		hour := 9 + rand.Intn(8)*2 // 9, 11, 13, 15, 17点

		statuses := []string{"pending", "approved", "completed", "cancelled"}
		status := statuses[rand.Intn(len(statuses))]

		visit := VisitAppointment{
			ID:            uint(i),
			AppointmentNo: fmt.Sprintf("VA%s%03d", now.Format("200601"), i),
			ElderlyID:     elderlyID,
			VisitorName:   randomName(rand.Intn(2) == 1),
			VisitorPhone:  generatePhone(),
			Relationship: relationships[rand.Intn(len(relationships))],
			VisitType:     []string{"现场探访", "视频探访"}[rand.Intn(2)],
			VisitDate:     visitDate.Format("2006-01-02"),
			VisitTime:     fmt.Sprintf("%02d:00", hour),
			VisitorCount:  1 + rand.Intn(3),
			Status:        status,
			Remarks:       randomRemarks(),
			CreatedAt:     randomTime(now.AddDate(-1, 0, 0), now),
		}

		visitList = append(visitList, visit)
	}

	data, _ := json.MarshalIndent(visitList, "", "  ")
	os.WriteFile("seed_data/visits.json", data, 0644)

	fmt.Printf("   ✅ 探视预约数据已生成到 seed_data/visits.json\n")
}

// 辅助函数
func randomName(isMale bool) string {
	surname := surnames[rand.Intn(len(surnames))]
	var name string
	if isMale {
		name = maleNames[rand.Intn(len(maleNames))]
	} else {
		name = femaleNames[rand.Intn(len(femaleNames))]
	}
	return surname + name
}

func generateIDCard(year int, isMale bool) string {
	areaCode := []string{"110101", "310101", "440101", "500101", "610101"}
	area := areaCode[rand.Intn(len(areaCode))]
	birthday := fmt.Sprintf("%d%02d%02d", year, 1+rand.Intn(12), 1+rand.Intn(28))
	sequence := rand.Intn(999)
	genderCode := sequence%10
	if !isMale {
		genderCode = (genderCode + 1) % 10
	}
	return area + birthday + fmt.Sprintf("%03d", genderCode) + "X"
}

func generatePhone() string {
	prefixes := []string{"130", "131", "132", "133", "135", "136", "137", "138", "139", "150", "151", "152", "153", "155", "156", "157", "158", "159", "180", "181", "182", "183", "185", "186", "187", "188", "189"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	return prefix + fmt.Sprintf("%08d", rand.Intn(100000000))
}

func randomAddress() string {
	cities := []string{"北京市", "上海市", "广州市", "深圳市", "南京市", "杭州市", "成都市", "武汉市"}
	districts := []string{"朝阳区", "海淀区", "浦东新区", "福田区", "玄武区", "西湖区", "武侯区", "武昌区"}
	city := cities[rand.Intn(len(cities))]
	district := districts[rand.Intn(len(districts)]
	return fmt.Sprintf("%s%s%s%d号", city, district, "幸福路", 1+rand.Intn(999))
}

func randomBedNumber() string {
	building := buildings[rand.Intn(len(buildings))]
	floor := floors[rand.Intn(len(floors))]
	room := roomNumbers[rand.Intn(len(roomNumbers))]
	return fmt.Sprintf("%s%s%s", building, floor, room)
}

func randomDate(start, end time.Time) time.Time {
	delta := end.Sub(start)
	randomDelta := time.Duration(rand.Int63n(int64(delta)))
	return start.Add(randomDelta)
}

func randomTime(start, end time.Time) time.Time {
	delta := end.Sub(start)
	randomDelta := time.Duration(rand.Int63n(int64(delta)))
	return start.Add(randomDelta)
}

func randomCareContent(careType string) string {
	contents := careContents[careType]
	if len(contents) == 0 {
		return careType + "服务"
	}
	return contents[rand.Intn(len(contents))]
}

func randomImages() string {
	count := rand.Intn(4)
	if count == 0 {
		return ""
	}
	images := make([]string, count)
	for i := 0; i < count; i++ {
		images[i] = fmt.Sprintf("/uploads/care_%d.jpg", rand.Intn(1000))
	}
	return fmt.Sprintf("[%s]", images[0])
}

func randomTags() string {
	allTags := []string{"服务热情", "专业细致", "耐心体贴", "技术精湛", "态度友好", "响应及时", "沟通顺畅"}
	count := rand.Intn(4)
	if count == 0 {
		return ""
	}
	tags := make([]string, count)
	for i := 0; i < count; i++ {
		tags[i] = allTags[rand.Intn(len(allTags))]
	}
	return fmt.Sprintf("[%s]", tags[0])
}

func randomRemarks() string {
	remarks := []string{
		"老人状态良好，配合度佳",
		"今日情绪稳定，无异常",
		"继续观察，加强护理",
		"老人食欲正常，睡眠良好",
		"按时完成护理任务",
		"家属表示满意",
		"",
		"",
	}
	return remarks[rand.Intn(len(remarks))]
}

func randomHealthRemarks() string {
	remarks := []string{
		"各项指标正常",
		"血压稍高，建议继续监测",
		"血糖正常，注意饮食",
		"心率正常，继续观察",
		"体温正常，无发热症状",
		"",
		"",
	}
	return remarks[rand.Intn(len(remarks))]
}
