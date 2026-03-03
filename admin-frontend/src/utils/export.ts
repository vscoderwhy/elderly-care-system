import * as XLSX from 'xlsx'

// 导出接口
interface ExportOptions {
  filename?: string
  sheetName?: string
  header?: string[]
  data: any[]
  columns?: Column[]
}

interface Column {
  key: string
  title: string
  width?: number
  formatter?: (value: any, row: any) => any
}

/**
 * 导出数据为 Excel
 */
export const exportExcel = (options: ExportOptions) => {
  const {
    filename = '导出数据',
    sheetName = 'Sheet1',
    header = [],
    data,
    columns = []
  } = options

  // 创建工作簿
  const workbook = XLSX.utils.book_new()

  // 准备导出数据
  let exportData: any[][] = []

  if (columns.length > 0) {
    // 使用列配置
    // 表头
    exportData.push(columns.map(col => col.title))

    // 数据行
    data.forEach(row => {
      const dataRow = columns.map(col => {
        const value = row[col.key]
        return col.formatter ? col.formatter(value, row) : value
      })
      exportData.push(dataRow)
    })
  } else if (header.length > 0) {
    // 使用表头配置
    exportData.push(header)

    // 数据行（假设数据是对象数组）
    if (data.length > 0 && typeof data[0] === 'object') {
      data.forEach(row => {
        const dataRow = header.map(() => '')
        Object.keys(row).forEach((key, index) => {
          if (index < header.length) {
            dataRow[index] = row[key]
          }
        })
        exportData.push(dataRow)
      })
    }
  } else {
    // 直接使用数据
    exportData = data
  }

  // 创建工作表
  const worksheet = XLSX.utils.aoa_to_sheet(exportData)

  // 设置列宽
  if (columns.length > 0) {
    const colWidths = columns.map(col => ({
      wch: col.width || 15
    }))
    worksheet['!cols'] = colWidths
  }

  // 添加工作表到工作簿
  XLSX.utils.book_append_sheet(workbook, worksheet, sheetName)

  // 生成文件名（带时间戳）
  const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  const fullFilename = `${filename}_${timestamp}.xlsx`

  // 导出文件
  XLSX.writeFile(workbook, fullFilename)
}

/**
 * 导出老人名单
 */
export const exportElderlyList = (data: any[]) => {
  const columns: Column[] = [
    { key: 'name', title: '姓名', width: 12 },
    { key: 'gender', title: '性别', width: 8 },
    { key: 'age', title: '年龄', width: 8 },
    { key: 'idCard', title: '身份证号', width: 20 },
    { key: 'phone', title: '联系电话', width: 15 },
    { key: 'careLevel', title: '护理等级', width: 12 },
    { key: 'bedNumber', title: '床位号', width: 12 },
    { key: 'checkInDate', title: '入住日期', width: 15 },
    { key: 'familyName', title: '家属联系人', width: 12 },
    { key: 'familyPhone', title: '家属电话', width: 15 },
    { key: 'healthStatus', title: '健康状况', width: 15 },
    { key: 'remark', title: '备注', width: 20 }
  ]

  exportExcel({
    filename: '老人名单',
    sheetName: '老人信息',
    columns,
    data
  })
}

/**
 * 导出护理记录
 */
export const exportCareRecords = (data: any[]) => {
  const columns: Column[] = [
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'careType', title: '护理类型', width: 15 },
    { key: 'careItem', title: '护理项目', width: 20 },
    { key: 'nurseName', title: '护理员', width: 12 },
    { key: 'careTime', title: '护理时间', width: 18 },
    { key: 'duration', title: '时长(分钟)', width: 12 },
    { key: 'description', title: '护理内容', width: 30 },
    { key: 'result', title: '护理结果', width: 20 },
    { key: 'evaluation', title: '评价', width: 15 },
    { key: 'status', title: '状态', width: 10 }
  ]

  exportExcel({
    filename: '护理记录',
    sheetName: '护理记录',
    columns,
    data
  })
}

/**
 * 导出财务报表
 */
export const exportFinancialReport = (data: any[]) => {
  const columns: Column[] = [
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'billNo', title: '账单号', width: 18 },
    { key: 'billType', title: '费用类型', width: 12 },
    { key: 'amount', title: '金额(元)', width: 12 },
    { key: 'billDate', title: '账单日期', width: 15 },
    { key: 'dueDate', title: '应付日期', width: 15 },
    { key: 'paidDate', title: '支付日期', width: 15 },
    { key: 'paymentMethod', title: '支付方式', width: 12 },
    { key: 'status', title: '状态', width: 10 },
    { key: 'remark', title: '备注', width: 20 }
  ]

  exportExcel({
    filename: '财务报表',
    sheetName: '账单明细',
    columns,
    data
  })
}

/**
 * 导出健康数据
 */
export const exportHealthData = (data: any[]) => {
  const columns: Column[] = [
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'recordDate', title: '记录日期', width: 15 },
    { key: 'bloodPressure', title: '血压(mmHg)', width: 15 },
    { key: 'heartRate', title: '心率(次/分)', width: 15 },
    { key: 'bloodSugar', title: '血糖(mmol/L)', width: 15 },
    { key: 'temperature', title: '体温(℃)', width: 12 },
    { key: 'weight', title: '体重(kg)', width: 12 },
    { key: 'oxygenSaturation', title: '血氧(%)', width: 12 },
    { key: 'recordBy', title: '记录人', width: 12 },
    { key: 'remark', title: '备注', width: 20 }
  ]

  exportExcel({
    filename: '健康数据',
    sheetName: '健康记录',
    columns,
    data
  })
}

/**
 * 导出护理任务
 */
export const exportNursingTasks = (data: any[]) => {
  const columns: Column[] = [
    { key: 'taskNo', title: '任务编号', width: 15 },
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'taskType', title: '任务类型', width: 15 },
    { key: 'taskContent', title: '任务内容', width: 30 },
    { key: 'nurseName', title: '执行人', width: 12 },
    { key: 'planTime', title: '计划时间', width: 18 },
    { key: 'executeTime', title: '执行时间', width: 18 },
    { key: 'duration', title: '时长(分钟)', width: 12 },
    { key: 'status', title: '状态', width: 10 },
    { key: 'remark', title: '备注', width: 20 }
  ]

  exportExcel({
    filename: '护理任务',
    sheetName: '任务列表',
    columns,
    data
  })
}

/**
 * 导出考勤记录
 */
export const exportAttendanceRecords = (data: any[]) => {
  const columns: Column[] = [
    { key: 'employeeName', title: '姓名', width: 12 },
    { key: 'department', title: '部门', width: 15 },
    { key: 'position', title: '职位', width: 12 },
    { key: 'date', title: '日期', width: 15 },
    { key: 'checkInTime', title: '上班时间', width: 15 },
    { key: 'checkOutTime', title: '下班时间', width: 15 },
    { key: 'workHours', title: '工作时长', width: 12 },
    { key: 'status', title: '状态', width: 10 },
    { key: 'lateMinutes', title: '迟到(分钟)', width: 12 },
    { key: 'earlyMinutes', title: '早退(分钟)', width: 12 }
  ]

  exportExcel({
    filename: '考勤记录',
    sheetName: '考勤明细',
    columns,
    data
  })
}

/**
 * 导出库存清单
 */
export const exportInventoryList = (data: any[]) => {
  const columns: Column[] = [
    { key: 'itemNo', title: '物品编号', width: 15 },
    { key: 'itemName', title: '物品名称', width: 20 },
    { key: 'category', title: '分类', width: 12 },
    { key: 'specification', title: '规格', width: 15 },
    { key: 'unit', title: '单位', width: 8 },
    { key: 'quantity', title: '库存数量', width: 12 },
    { key: 'safetyStock', title: '安全库存', width: 12 },
    { key: 'unitPrice', title: '单价(元)', width: 12 },
    { key: 'totalValue', title: '总金额(元)', width: 12 },
    { key: 'supplier', title: '供应商', width: 15 },
    { key: 'lastPurchaseDate', title: '最后采购日期', width: 15 }
  ]

  exportExcel({
    filename: '库存清单',
    sheetName: '库存明细',
    columns,
    data
  })
}

/**
 * 导出访视记录
 */
export const exportVisitRecords = (data: any[]) => {
  const columns: Column[] = [
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'visitorName', title: '访客姓名', width: 12 },
    { key: 'relation', title: '关系', width: 10 },
    { key: 'visitDate', title: '访视日期', width: 15 },
    { key: 'visitTime', title: '访视时间', width: 15 },
    { key: 'duration', title: '时长(分钟)', width: 12 },
    { key: 'purpose', title: '来访目的', width: 20 },
    { key: 'gift', title: '携带物品', width: 15 },
    { key: 'remark', title: '备注', width: 20 },
    { key: 'recordBy', title: '记录人', width: 12 }
  ]

  exportExcel({
    filename: '访视记录',
    sheetName: '访视明细',
    columns,
    data
  })
}

/**
 * 导出多Sheet数据
 */
export const exportMultiSheet = (sheets: {
  name: string
  data: any[][]
  columns?: Column[]
}[], filename = '导出数据') => {
  const workbook = XLSX.utils.book_new()

  sheets.forEach(sheet => {
    let exportData: any[][] = []

    if (sheet.columns && sheet.columns.length > 0) {
      // 表头
      exportData.push(sheet.columns.map(col => col.title))

      // 数据行
      sheet.data.forEach((row: any) => {
        const dataRow = sheet.columns!.map(col => {
          const value = row[col.key]
          return col.formatter ? col.formatter(value, row) : value
        })
        exportData.push(dataRow)
      })
    } else {
      exportData = sheet.data
    }

    const worksheet = XLSX.utils.aoa_to_sheet(exportData)

    if (sheet.columns) {
      worksheet['!cols'] = sheet.columns.map(col => ({
        wch: col.width || 15
      }))
    }

    XLSX.utils.book_append_sheet(workbook, worksheet, sheet.name)
  })

  const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  XLSX.writeFile(workbook, `${filename}_${timestamp}.xlsx`)
}

export default {
  exportExcel,
  exportElderlyList,
  exportCareRecords,
  exportFinancialReport,
  exportHealthData,
  exportNursingTasks,
  exportAttendanceRecords,
  exportInventoryList,
  exportVisitRecords,
  exportMultiSheet
}
