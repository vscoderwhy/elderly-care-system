package service

import (
	"elderly-care-system/internal/repository"
	"time"
)

type StatsService struct {
	elderlyRepo *repository.ElderlyRepository
	careRepo    *repository.CareRepository
	billRepo    *repository.BillRepository
	userRepo    *repository.UserRepository
	roomRepo    *repository.RoomRepository
}

func NewStatsService(
	elderlyRepo *repository.ElderlyRepository,
	careRepo *repository.CareRepository,
	billRepo *repository.BillRepository,
	userRepo *repository.UserRepository,
	roomRepo *repository.RoomRepository,
) *StatsService {
	return &StatsService{
		elderlyRepo: elderlyRepo,
		careRepo:    careRepo,
		billRepo:    billRepo,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
	}
}

type DashboardStats struct {
	ElderlyCount      int            `json:"elderly_count"`       // 在院老人数
	TodayCareCount    int64          `json:"today_care_count"`    // 今日护理次数
	PendingBills      int64          `json:"pending_bills"`       // 待缴费账单
	StaffCount        int64          `json:"staff_count"`         // 在职员工
	PendingServices   int64          `json:"pending_services"`    // 待处理服务请求
	BedOccupancyRate  float64        `json:"bed_occupancy_rate"`  // 入住率
	RecentCares       []CareRecord   `json:"recent_cares"`        // 最近护理记录
	PendingTasks      []ServiceTask  `json:"pending_tasks"`       // 待办任务
}

type CareRecord struct {
	ID         uint   `json:"id"`
	ElderlyName string `json:"elderly_name"`
	CareItem   string `json:"care_item"`
	StaffName  string `json:"staff_name"`
	RecordedAt string `json:"recorded_at"`
}

type ServiceTask struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	ElderlyName string `json:"elderly_name"`
	Status      string `json:"status"`
	RequestedAt string `json:"requested_at"`
}

// GetDashboardStats 获取Dashboard统计数据
func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 在院老人数
	elderlyCount, err := s.elderlyRepo.Count()
	if err == nil {
		stats.ElderlyCount = elderlyCount
	}

	// 今日护理次数
	todayCare, err := s.careRepo.CountToday()
	if err == nil {
		stats.TodayCareCount = todayCare
	}

	// 待缴费账单
	pendingBills, err := s.billRepo.CountPending()
	if err == nil {
		stats.PendingBills = pendingBills
	}

	// 在职员工数
	staffCount, err := s.userRepo.CountStaff()
	if err == nil {
		stats.StaffCount = staffCount
	}

	// 待处理服务请求
	pendingServices, err := s.careRepo.CountPendingServices()
	if err == nil {
		stats.PendingServices = pendingServices
	}

	// 入住率
	occupancy, err := s.roomRepo.GetOccupancyRate()
	if err == nil {
		stats.BedOccupancyRate = occupancy
	}

	// 最近护理记录
	recentCares, err := s.careRepo.GetRecentRecords(5)
	if err == nil {
		for _, r := range recentCares {
			stats.RecentCares = append(stats.RecentCares, CareRecord{
				ID:          r.ID,
				ElderlyName: r.Elderly.Name,
				CareItem:    r.CareItem.Name,
				StaffName:   r.Staff.Nickname,
				RecordedAt:  r.RecordedAt.Format("2006-01-02 15:04"),
			})
		}
	}

	// 待办任务（待处理的服务请求）
	pendingTasks, err := s.careRepo.GetPendingServiceRequests(5)
	if err == nil {
		for _, t := range pendingTasks {
			elderlyName := ""
			if t.Elderly != nil {
				elderlyName = t.Elderly.Name
			}
			stats.PendingTasks = append(stats.PendingTasks, ServiceTask{
				ID:          t.ID,
				Type:        t.Type,
				ElderlyName: elderlyName,
				Status:      t.Status,
				RequestedAt: t.RequestedAt.Format("2006-01-02 15:04"),
			})
		}
	}

	return stats, nil
}

// BedOccupancyStat 床位入住统计
type BedOccupancyStat struct {
	BuildingName string  `json:"building_name"`
	TotalBeds    int     `json:"total_beds"`
	OccupiedBeds int     `json:"occupied_beds"`
	OccupancyRate float64 `json:"occupancy_rate"`
}

// GetBedOccupancy 获取各楼栋入住率
func (s *StatsService) GetBedOccupancy() ([]BedOccupancyStat, error) {
	data, err := s.roomRepo.GetBuildingOccupancy()
	if err != nil {
		return nil, err
	}
	var stats []BedOccupancyStat
	for _, d := range data {
		stats = append(stats, BedOccupancyStat{
			BuildingName:  d.BuildingName,
			TotalBeds:     d.TotalBeds,
			OccupiedBeds:  d.OccupiedBeds,
			OccupancyRate: d.OccupancyRate,
		})
	}
	return stats, nil
}

// CareStat 护理统计
type CareStat struct {
	Date       string `json:"date"`
	CareCount  int64  `json:"care_count"`
}

// GetCareStats 获取近7天护理统计
func (s *StatsService) GetCareStats() ([]CareStat, error) {
	stats := []CareStat{}
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		count, err := s.careRepo.CountByDate(date)
		if err != nil {
			count = 0
		}
		stats = append(stats, CareStat{
			Date:      dateStr,
			CareCount: count,
		})
	}

	return stats, nil
}

// FinanceStat 财务统计
type FinanceStat struct {
	TotalIncome    float64 `json:"total_income"`     // 总收入
	MonthlyIncome  float64 `json:"monthly_income"`   // 本月收入
	PendingAmount  float64 `json:"pending_amount"`   // 待收金额
}

// GetFinanceStats 获取财务统计
func (s *StatsService) GetFinanceStats() (*FinanceStat, error) {
	data, err := s.billRepo.GetFinanceStats()
	if err != nil {
		return nil, err
	}
	return &FinanceStat{
		TotalIncome:   data.TotalIncome,
		MonthlyIncome: data.MonthlyIncome,
		PendingAmount: data.PendingAmount,
	}, nil
}
