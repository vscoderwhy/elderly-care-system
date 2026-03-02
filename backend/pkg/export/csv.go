package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// ElderlyReport 老人信息报表
type ElderlyReport struct {
	Name             string
	Gender           string
	Age              int
	Room            string
	Bed             string
	CareLevel       string
	HealthStatus    string
	AdmissionDate   string
	EmergencyContact string
	EmergencyPhone  string
}

// CareRecordReport 护理记录报表
type CareRecordReport struct {
	ElderlyName   string
	CareItem       string
	StaffName      string
	RecordedAt     string
	Notes          string
}

// HealthDataReport 健康数据报表
type HealthDataReport struct {
	ElderlyName  string
	RecordType   string
	Value        string
	Unit         string
	RecordedAt   string
	RecorderName string
}

// FinanceReport 财务报表
type FinanceReport struct {
	BillNo       string
	ElderlyName  string
	TotalAmount  string
	Status       string
	CreatedAt    string
	PaidAt       string
}

// ExportElderlyListCSV 导出老人列表为CSV
func ExportElderlyListCSV(data []ElderlyReport) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{"姓名", "性别", "年龄", "房间", "床位", "护理等级", "健康状态", "入住日期", "紧急联系人", "联系电话"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, record := range data {
		row := []string{
			record.Name,
			record.Gender,
			fmt.Sprintf("%d", record.Age),
			record.Room,
			record.Bed,
			record.CareLevel,
			record.HealthStatus,
			record.AdmissionDate,
			record.EmergencyContact,
			record.EmergencyPhone,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// ExportCareRecordsCSV 导出护理记录CSV
func ExportCareRecordsCSV(data []CareRecordReport) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"老人姓名", "护理项目", "护理员", "记录时间", "备注"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, record := range data {
		row := []string{
			record.ElderlyName,
			record.CareItem,
			record.StaffName,
			record.RecordedAt,
			record.Notes,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// ExportHealthDataCSV 导出健康数据CSV
func ExportHealthDataCSV(data []HealthDataReport) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"老人姓名", "记录类型", "数值", "单位", "记录时间", "记录人"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, record := range data {
		row := []string{
			record.ElderlyName,
			record.RecordType,
			record.Value,
			record.Unit,
			record.RecordedAt,
			record.RecorderName,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// ExportFinanceCSV 导出财务报表CSV
func ExportFinanceCSV(data []FinanceReport) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"账单号", "老人姓名", "金额", "状态", "创建时间", "支付时间"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, record := range data {
		row := []string{
			record.BillNo,
			record.ElderlyName,
			record.TotalAmount,
			record.Status,
			record.CreatedAt,
			record.PaidAt,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// GenerateMonthlyReportCSV 生成月度综合报表CSV
func GenerateMonthlyReportCSV(year int, month int, data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 标题行
	_ = writer.Write([]string{"统计项目", "数值"})

	// 统计数据
	for key, value := range data {
		_ = writer.Write([]string{key, fmt.Sprintf("%v", value)})
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// MedicationReport 用药记录报表
type MedicationReport struct {
	ElderlyName    string
	MedicationName string
	Dosage         string
	Frequency      string
	StartDate      string
	EndDate        string
	Status         string
}

// ExportMedicationRecordsCSV 导出用药记录CSV
func ExportMedicationRecordsCSV(data []MedicationReport) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"老人姓名", "药品名称", "剂量", "频率", "开始日期", "结束日期", "状态"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, record := range data {
		row := []string{
			record.ElderlyName,
			record.MedicationName,
			record.Dosage,
			record.Frequency,
			record.StartDate,
			record.EndDate,
			record.Status,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}
