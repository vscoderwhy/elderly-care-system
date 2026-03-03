package handler

import (
	"elderly-care-system/internal/repository"
	"elderly-care-system/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	elderlyRepo *repository.ElderlyRepository
	careRepo    *repository.CareRepository
	billRepo    *repository.BillRepository
	roomRepo    *repository.RoomRepository
}

func NewStatisticsHandler(
	elderlyRepo *repository.ElderlyRepository,
	careRepo *repository.CareRepository,
	billRepo *repository.BillRepository,
	roomRepo *repository.RoomRepository,
) *StatisticsHandler {
	return &StatisticsHandler{
		elderlyRepo: elderlyRepo,
		careRepo:    careRepo,
		billRepo:    billRepo,
		roomRepo:    roomRepo,
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *StatisticsHandler) GetDashboardStats(c *gin.Context) {
	// 获取老人总数
	elderly, total, _ := h.elderlyRepo.List(1, 10000)

	// 获取床位统计
	bedStats, _ := h.roomRepo.GetBedStats()

	// 获取本月护理记录
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	careRecords, _, _ := h.careRepo.FindRecordsByDateRange(startOfMonth, time.Now())

	// 计算护理等级分布
	careLevelDist := make(map[int]int)
	for _, e := range elderly {
		careLevelDist[e.CareLevel]++
	}

	// 计算性别分布
	genderDist := make(map[string]int)
	for _, e := range elderly {
		genderDist[e.Gender]++
	}

	// 获取待处理预警
	alertStats := make(map[string]int)
	// TODO: 从alert repository获取

	// 计算入住率，避免除零
	var occupancyRate float64
	if bedStats["total"] > 0 {
		occupancyRate = float64(bedStats["occupied"]) / float64(bedStats["total"]) * 100
	}

	response.Success(c, gin.H{
		"elderly_total":      total,
		"bed_total":          bedStats["total"],
		"bed_occupied":       bedStats["occupied"],
		"bed_available":      bedStats["total"] - bedStats["occupied"],
		"occupancy_rate":     occupancyRate,
		"care_records_today": len(careRecords),
		"care_level_dist":    careLevelDist,
		"gender_dist":        genderDist,
		"alerts":             alertStats,
		"updated_at":         time.Now(),
	})
}

// GetOccupancyTrend 获取入住率趋势（最近30天）
func (h *StatisticsHandler) GetOccupancyTrend(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if num, err := strconv.Atoi(d); err == nil && num > 0 && num <= 365 {
			days = num
		}
	}

	trend := make([]gin.H, 0, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		// TODO: 从数据库获取当天的入住数据
		trend = append(trend, gin.H{
			"date":         date.Format("2006-01-02"),
			"occupied":     0,
			"total":        0,
			"occupancy_rate": 0,
		})
	}

	response.Success(c, trend)
}

// GetHealthTrend 获取健康数据趋势
func (h *StatisticsHandler) GetHealthTrend(c *gin.Context) {
	elderlyID := c.Query("elderly_id")
	days := 30

	records, err := h.careRepo.GetHealthTrend(elderlyID, days)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, records)
}

// GetFinanceStats 获取财务统计数据
func (h *StatisticsHandler) GetFinanceStats(c *gin.Context) {
	// 获取本月数据
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// TODO: 从bill repository获取财务数据
	stats := gin.H{
		"month_income":     0,
		"month_expense":    0,
		"pending_bills":    0,
		"overdue_bills":    0,
		"collection_rate":  0,
		"period": gin.H{
			"start": startOfMonth.Format("2006-01-02"),
			"end":   endOfMonth.Format("2006-01-02"),
		},
	}

	response.Success(c, stats)
}

// GetCareStats 获取护理统计
func (h *StatisticsHandler) GetCareStats(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	endParsed, _ := time.Parse("2006-01-02", endDate)
	end := endParsed.AddDate(0, 0, 1)

	records, _, err := h.careRepo.FindRecordsByDateRange(start, end)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	// 按护理项目统计
	itemStats := make(map[string]int)
	staffStats := make(map[string]int)

	for _, r := range records {
		if r.CareItem.Name != "" {
			itemStats[r.CareItem.Name]++
		}
		if r.Staff != nil {
			staffStats[r.Staff.Nickname]++
		}
	}

	response.Success(c, gin.H{
		"total_records": len(records),
		"item_stats":    itemStats,
		"staff_stats":   staffStats,
		"period": gin.H{
			"start": startDate,
			"end":   endDate,
		},
	})
}

// GetElderlyAgeDistribution 获取老人年龄分布
func (h *StatisticsHandler) GetElderlyAgeDistribution(c *gin.Context) {
	elderly, _, _ := h.elderlyRepo.List(1, 10000)

	// 年龄分组
	ageGroups := map[string]int{
		"60-69": 0,
		"70-79": 0,
		"80-89": 0,
		"90+":   0,
	}

	for _, e := range elderly {
		age := calculateAge(e.BirthDate)
		if age >= 90 {
			ageGroups["90+"]++
		} else if age >= 80 {
			ageGroups["80-89"]++
		} else if age >= 70 {
			ageGroups["70-79"]++
		} else if age >= 60 {
			ageGroups["60-69"]++
		}
	}

	response.Success(c, ageGroups)
}

// GetMonthlyReport 获取月度综合报表
func (h *StatisticsHandler) GetMonthlyReport(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	report := gin.H{
		"period": gin.H{
			"year":  year,
			"month": month,
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
		// TODO: 添加更多统计数据
		"elderly": gin.H{
			"total":        0,
			"new_admitted": 0,
			"discharged":   0,
			"deceased":     0,
		},
		"care": gin.H{
			"total_records": 0,
			"completion_rate": 0,
		},
		"finance": gin.H{
			"income":       0,
			"expense":       0,
			"pending_amount": 0,
		},
		"health": gin.H{
			"abnormal_count": 0,
			"emergency_count": 0,
		},
	}

	response.Success(c, report)
}
