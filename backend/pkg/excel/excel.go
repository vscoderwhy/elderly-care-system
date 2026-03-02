package excel

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
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

// ExportElderlyList 导出老人列表
func ExportElderlyList(data []ElderlyReport) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet, err := f.NewSheet("老人列表")
	if err != nil {
		return nil, err
	}

	// 设置表头
	headers := []string{"姓名", "性别", "年龄", "房间", "床位", "护理等级", "健康状态", "入住日期", "紧急联系人", "联系电话"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, &excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: "#4F81BD", Pattern: 1},
	})
	}

	// 填充数据
	for row, record := range data {
		rowNum := row + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), record.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), record.Gender)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("%d", record.Age))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), record.Room)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), record.Bed)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), record.CareLevel)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), record.HealthStatus)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", rowNum), record.AdmissionDate)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", rowNum), record.EmergencyContact)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", rowNum), record.EmergencyPhone)
	}

	// 自动调整列宽
	for i := range headers {
		col, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetColWidth(sheet, col, 15)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// ExportCareRecords 导出护理记录
func ExportCareRecords(data []CareRecordReport) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet, err := f.NewSheet("护理记录")
	if err != nil {
		return nil, err
	}

	headers := []string{"老人姓名", "护理项目", "护理员", "记录时间", "备注"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(sheet, cell, header)
	}

	for row, record := range data {
		rowNum := row + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), record.ElderlyName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), record.CareItem)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), record.StaffName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), record.RecordedAt)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), record.Notes)
	}

	for i := range headers {
		col, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetColWidth(sheet, col, 20)
	}

	return f.WriteToBuffer()
}

// ExportHealthData 导出健康数据
func ExportHealthData(data []HealthDataReport) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet, err := f.NewSheet("健康数据")
	if err != nil {
		return nil, err
	}

	headers := []string{"老人姓名", "记录类型", "数值", "单位", "记录时间", "记录人"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(sheet, cell, header)
	}

	for row, record := range data {
		rowNum := row + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), record.ElderlyName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), record.RecordType)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), record.Value)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), record.Unit)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), record.RecordedAt)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), record.RecorderName)
	}

	return f.WriteToBuffer()
}

// ExportFinance 导出财务报表
func ExportFinance(data []FinanceReport) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet, err := f.NewSheet("财务报表")
	if err != nil {
		return nil, err
	}

	headers := []string{"账单号", "老人姓名", "金额", "状态", "创建时间", "支付时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(sheet, cell, header)
	}

	for row, record := range data {
		rowNum := row + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), record.BillNo)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), record.ElderlyName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), record.TotalAmount)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), record.Status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), record.CreatedAt)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), record.PaidAt)
	}

	return f.WriteToBuffer()
}

// GenerateMonthlyReport 生成月度综合报表
func GenerateMonthlyReport(year int, month int, data map[string]interface{}) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet, err := f.NewSheet("月度报表")
	if err != nil {
		return nil, err
	}

	// 标题
	title := fmt.Sprintf("%d年%d月 养老院运营月报", year, month)
	f.SetCellValue(sheet, "A1", title)
	f.MergeCell(sheet, "A1", "F1", title)

	// 统计数据
	row := 3
	for key, value := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), key)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%v", value))
		row++
	}

	return f.WriteToBuffer()
}
