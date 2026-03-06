package handler

import (
	"elderly-care-system/internal/repository"
	"elderly-care-system/pkg/export"
	"elderly-care-system/pkg/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	elderlyRepo *repository.ElderlyRepository
	careRepo    *repository.CareRepository
	billRepo    *repository.BillRepository
}

func NewExportHandler(
	elderlyRepo *repository.ElderlyRepository,
	careRepo *repository.CareRepository,
	billRepo *repository.BillRepository,
) *ExportHandler {
	return &ExportHandler{
		elderlyRepo: elderlyRepo,
		careRepo:    careRepo,
		billRepo:    billRepo,
	}
}

// ExportElderlyList 导出老人列表
func (h *ExportHandler) ExportElderlyList(c *gin.Context) {
	// 获取分页参数，默认导出所有数据
	page := 1
	pageSize := 100000 // 大批量导出
	
	elderly, _, err := h.elderlyRepo.List(page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取老人列表失败")
		return
	}

	report := make([]export.ElderlyReport, 0, len(elderly))
	for _, e := range elderly {
		room := ""
		bed := ""
		if e.Bed != nil {
			bed = e.Bed.Name
			// 简化处理，显示楼层ID作为房间号
			room = fmt.Sprintf("%d室", e.Bed.RoomID)
		}

		// 处理性别显示
		genderText := e.Gender
		if genderText == "male" {
			genderText = "男"
		} else if genderText == "female" {
			genderText = "女"
		}

		// 处理健康状态
		healthStatus := "良好"
		if e.HealthStatus != "" && e.HealthStatus != "{}" {
			healthStatus = e.HealthStatus
		}

		report = append(report, export.ElderlyReport{
			Name:             e.Name,
			Gender:           genderText,
			Age:              calculateAge(e.BirthDate),
			Room:             room,
			Bed:              bed,
			CareLevel:        fmt.Sprintf("%d级护理", e.CareLevel),
			HealthStatus:     healthStatus,
			AdmissionDate:    formatDate(e.AdmissionDate),
			EmergencyContact: e.EmergencyContact,
			EmergencyPhone:   e.EmergencyPhone,
		})
	}

	data, err := export.ExportElderlyListCSV(report)
	if err != nil {
		response.Error(c, 500, "生成报表失败")
		return
	}

	filename := fmt.Sprintf("老人列表_%s.csv", time.Now().Format("20060102150405"))
	
	// 设置正确的响应头，告诉浏览器这是文件下载
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	
	// 直接返回CSV数据
	c.Data(200, "text/csv; charset=utf-8", data)
}

// ExportCareRecords 导出护理记录
func (h *ExportHandler) ExportCareRecords(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	_, _ = time.Parse("2006-01-02", startDate)
	_, _ = time.Parse("2006-01-02", endDate)

	records, _, err := h.careRepo.FindRecordsByElderlyID(0, 1, 10000)
	if err != nil {
		response.Error(c, 500, "获取护理记录失败")
		return
	}

	report := make([]export.CareRecordReport, 0)
	for _, r := range records {
		elderlyName := ""
		if r.Elderly != nil {
			elderlyName = r.Elderly.Name
		}
		staffName := ""
		if r.Staff != nil {
			staffName = r.Staff.Nickname
		}

		report = append(report, export.CareRecordReport{
			ElderlyName: elderlyName,
			CareItem:      r.CareItem.Name,
			StaffName:     staffName,
			RecordedAt:   r.RecordedAt.Format("2006-01-02 15:04:05"),
			Notes:        r.Notes,
		})
	}

	data, err := export.ExportCareRecordsCSV(report)
	if err != nil {
		response.Error(c, 500, "生成报表失败")
		return
	}

	filename := fmt.Sprintf("护理记录_%s_%s.csv", startDate, endDate)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "text/csv; charset=utf-8", data)
}

// ExportHealthData 导出健康数据
func (h *ExportHandler) ExportHealthData(c *gin.Context) {
	_ = c.Query("elderly_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	records, err := h.careRepo.GetLatestHealthRecords(1)
	if err != nil {
		response.Error(c, 500, "获取健康数据失败")
		return
	}

	report := make([]export.HealthDataReport, 0)
	for _, r := range records {
		elderlyName := ""
		if r.Elderly != nil {
			elderlyName = r.Elderly.Name
		}
		recorderName := ""
		if r.Recorder != nil {
			recorderName = r.Recorder.Nickname
		}

		report = append(report, export.HealthDataReport{
			ElderlyName:  elderlyName,
			RecordType:   r.RecordType,
			Value:        r.Value,
			Unit:         r.Unit,
			RecordedAt:   r.RecordedAt.Format("2006-01-02 15:04:05"),
			RecorderName: recorderName,
		})
	}

	data, err := export.ExportHealthDataCSV(report)
	if err != nil {
		response.Error(c, 500, "生成报表失败")
		return
	}

	filename := fmt.Sprintf("健康数据_%s_%s.csv", startDate, endDate)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "text/csv; charset=utf-8", data)
}

// ExportFinance 导出财务报表
func (h *ExportHandler) ExportFinance(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 暂时返回空数据，因为BillRepository没有List方法
	report := make([]export.FinanceReport, 0)

	data, err := export.ExportFinanceCSV(report)
	if err != nil {
		response.Error(c, 500, "生成报表失败")
		return
	}

	filename := fmt.Sprintf("财务报表_%s_%s.csv", startDate, endDate)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "text/csv; charset=utf-8", data)
}

func calculateAge(birthDate *time.Time) int {
	if birthDate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func formatDate(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format("2006-01-02")
}
